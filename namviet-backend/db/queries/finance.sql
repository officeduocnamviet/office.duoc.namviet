-- Finance (P3 — ĐƯỜNG GHI phiếu thu + ĐỌC phiếu của đơn, ADR 0001 strangler):
-- GHI/ĐỌC bảng public.finance_transactions kế thừa Supabase. amount = NUMERIC →
-- common/money ở repo, KHÔNG float. Liên kết đơn↔phiếu thu qua
-- (ref_type='order' AND ref_id = orders.code) — ref_id là TEXT (mã đơn), KHÔNG uuid.
--
-- ⚠️ SỐ DƯ QUỸ: PROD có trigger on_finance_transaction_change tự cộng
-- fund_accounts.balance khi phiếu sang status='completed'. Vì vậy KHÔNG có query
-- nào UPDATE fund_accounts.balance ở đây — tránh double-count. Go chỉ INSERT phiếu
-- 'completed', trigger lo số dư.
--
-- IDEMPOTENCY: chống ghi phiếu thu TRÙNG (cộng tiền 2 lần):
--   * Phiếu webhook: dedup theo bank_reference_id (FindAlivePaymentByBankRef).
--   * Phiếu thủ công: dedup theo code idempotent sinh từ idem key
--     (FindAlivePaymentByCode). Unique partial index (deleted_at IS NULL) ở DB
--     chặn race (insert lần 2 → unique violation → app re-SELECT trả phiếu cũ).

-- name: FindAlivePaymentByCode :one
-- Tìm phiếu còn sống (chưa xoá mềm) theo code — để dedup phiếu thủ công idempotent.
-- Không thấy → ErrNoRows (chưa có, được phép insert).
SELECT
    id, code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type,
    created_by, created_at
FROM public.finance_transactions
WHERE code = @code::text
  AND deleted_at IS NULL;

-- name: FindAlivePaymentByBankRef :one
-- Tìm phiếu còn sống theo bank_reference_id — để dedup phiếu thu tự động (webhook).
-- Không thấy → ErrNoRows (chưa có, được phép insert).
SELECT
    id, code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type,
    created_by, created_at
FROM public.finance_transactions
WHERE bank_reference_id = @bank_reference_id::text
  AND deleted_at IS NULL;

-- name: InsertPaymentIn :one
-- Ghi MỘT phiếu THU (flow='in') vào public.finance_transactions. id do DB sinh
-- (IDENTITY) → KHÔNG truyền id. status truyền 'completed' (tiền đã thực vào quỹ) →
-- trigger prod tự cộng fund_accounts.balance MỘT LẦN (Go KHÔNG tự cộng). bank_ref
-- NULL cho phiếu thủ công. RETURNING toàn bộ cột để map về domain.Payment.
INSERT INTO public.finance_transactions (
    code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type, created_by
) VALUES (
    @code::text, now(), 'in', @business_type::text, @amount::numeric,
    @fund_account_id::bigint, @ref_type::text, @ref_id::text,
    sqlc.narg('description')::text, @status::text,
    sqlc.narg('bank_reference_id')::text, @book_type::text,
    sqlc.narg('created_by')::uuid
)
RETURNING
    id, code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type,
    created_by, created_at;

-- name: InsertPaymentOut :one
-- Ghi MỘT phiếu CHI (flow='out') vào public.finance_transactions — trả NCC mua hàng
-- (mục 54), đối xứng InsertPaymentIn. id do DB sinh (IDENTITY). status 'completed' =
-- đã chi khỏi quỹ → trigger prod TRỪ fund_accounts.balance MỘT LẦN (Go KHÔNG tự trừ).
-- bank_ref NULL cho phiếu thủ công. RETURNING toàn bộ cột để map về domain.Payment.
INSERT INTO public.finance_transactions (
    code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type, created_by
) VALUES (
    @code::text, now(), 'out', @business_type::text, @amount::numeric,
    @fund_account_id::bigint, @ref_type::text, @ref_id::text,
    sqlc.narg('description')::text, @status::text,
    sqlc.narg('bank_reference_id')::text, @book_type::text,
    sqlc.narg('created_by')::uuid
)
RETURNING
    id, code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type,
    created_by, created_at;

-- name: ListPaymentsByRef :many
-- Liệt kê phiếu thu/chi còn sống trỏ về MỘT chứng từ gốc (ref_type, ref_id) — cho
-- HTTP đọc /v1/finance/transactions của 1 đơn. Sắp transaction_date DESC, id DESC
-- (mới nhất trước, ổn định). Trả mọi flow/status để FE thấy đủ (UI tự lọc nếu cần).
SELECT
    id, code, transaction_date, flow, business_type, amount, fund_account_id,
    ref_type, ref_id, description, status, bank_reference_id, book_type,
    created_by, created_at
FROM public.finance_transactions
WHERE ref_type = @ref_type::text
  AND ref_id = @ref_id::text
  AND deleted_at IS NULL
ORDER BY transaction_date DESC, id DESC
LIMIT @row_limit::int;

-- name: ConfirmPaymentReceipt :execrows
-- Thủ quỹ "Xác nhận đã thu" (thanh toán 2 bước, spec mục 55): chuyển phiếu THU từ
-- 'pending' (đã thu từ khách, chưa vào quỹ) → 'completed' (vào quỹ). Trigger prod
-- cộng fund_accounts.balance MỘT LẦN tại bước này. CHỈ tác động phiếu đang 'pending'
-- (idempotent: gọi lại trên phiếu đã 'completed' → 0 dòng, KHÔNG cộng đôi). flow='in'.
UPDATE public.finance_transactions
SET status = 'completed', updated_at = now()
WHERE id = @payment_id::bigint
  AND status = 'pending'
  AND lower(flow) = 'in'
  AND deleted_at IS NULL;
