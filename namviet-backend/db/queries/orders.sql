-- Orders (ĐỌC + GHI nền P4a, ADR 0001): ĐỌC bảng public.* kế thừa; GHI tạo đơn +
-- đổi trạng thái ĐƠN GIẢN (xem khối "ĐƯỜNG GHI" cuối file). Trừ kho FEFO, ghi
-- phiếu thu, post sổ, POS atomic (ShipOrder/RecordPayment) HOÃN sang P4b (cần
-- cross-module). Tiền (total_amount/final_amount/unit_price/discount/total_line/amount) là
-- NUMERIC → common/money ở repo, KHÔNG float. order_items.quantity là INTEGER →
-- domain.Quantity (decimal) ở repo.
--
-- "ĐÃ THU" của đơn = tổng phiếu THU đã hoàn tất trỏ về đơn qua
--   (ref_type='order' AND ref_id = orders.code)  [ref_id là TEXT = MÃ đơn]
-- lọc:  lower(flow)='in'  (ERP set 'in' lowercase; so khớp case-insensitive cho
-- an toàn vì doc ghi 'IN'),  status='completed'  (tiền đã thực vào quỹ),
-- deleted_at IS NULL,  book_type IN ('INTERNAL','BOTH')  → SỔ THỰC TẾ (giả định
-- cần xác nhận business). orders KHÔNG có paid_amount → đây là suy diễn read-only.

-- name: ListOrders :many
-- Keyset pagination theo (created_at DESC, id DESC): trang kế lấy đơn "cũ hơn"
-- mốc cursor — (created_at < @after_created_at) OR (created_at = @after_created_at
-- AND id < @after_id). @after_created_at NULL = trang đầu (lấy từ mới nhất). Lọc
-- optional customer_id, status, payment_status, khoảng created_at [from_date,
-- to_date]. Mỗi đơn kèm paid_amount (subquery aggregate finance_transactions).
SELECT
    o.id,
    o.code,
    o.customer_id,
    o.creator_id,
    o.status,
    o.order_type,
    o.total_amount,
    o.final_amount,
    o.payment_status,
    o.note,
    o.created_at,
    o.updated_at,
    app.order_paid_amount(o.code) AS paid_amount
FROM public.orders o
WHERE o.deleted_at IS NULL
  AND (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR o.created_at < sqlc.narg('after_created_at')::timestamptz
        OR (o.created_at = sqlc.narg('after_created_at')::timestamptz
            AND o.id < sqlc.narg('after_id')::uuid)
      )
  AND (sqlc.narg('customer_id')::bigint IS NULL OR o.customer_id = sqlc.narg('customer_id')::bigint)
  AND (sqlc.narg('status')::text IS NULL OR o.status = sqlc.narg('status')::text)
  AND (sqlc.narg('payment_status')::text IS NULL OR o.payment_status = sqlc.narg('payment_status')::text)
  AND (sqlc.narg('from_date')::timestamptz IS NULL OR o.created_at >= sqlc.narg('from_date')::timestamptz)
  AND (sqlc.narg('to_date')::timestamptz IS NULL OR o.created_at <= sqlc.narg('to_date')::timestamptz)
ORDER BY o.created_at DESC, o.id DESC
LIMIT @row_limit::int;

-- name: GetOrderByID :one
-- Một đơn theo id (uuid) chưa soft-delete, kèm paid_amount suy diễn.
SELECT
    o.id,
    o.code,
    o.customer_id,
    o.creator_id,
    o.status,
    o.order_type,
    o.total_amount,
    o.final_amount,
    o.payment_status,
    o.note,
    o.created_at,
    o.updated_at,
    app.order_paid_amount(o.code) AS paid_amount
FROM public.orders o
WHERE o.id = @order_id::uuid AND o.deleted_at IS NULL;

-- name: ListOrderLines :many
-- Dòng hàng của một đơn (chưa soft-delete), thứ tự ổn định (created_at, id).
SELECT
    oi.id,
    oi.product_id,
    oi.quantity,
    oi.uom,
    oi.unit_price,
    oi.discount,
    oi.total_line,
    oi.is_gift,
    oi.batch_no,
    oi.expiry_date,
    oi.note
FROM public.order_items oi
WHERE oi.order_id = @order_id::uuid
  AND oi.deleted_at IS NULL
ORDER BY oi.created_at ASC, oi.id ASC;

-- ============================================================================
-- ĐƯỜNG GHI (P4a — tạo đơn + state machine ĐƠN GIẢN, ADR 0001 strangler). Chạy
-- TRONG tx do app mở (platform/db.WithinTx). GHI bảng public.orders + order_items
-- (mã đơn app tự sinh từ app.order_code_seq + tiền tố; id uuid v7 app-side). Đổi
-- trạng thái dùng SELECT ... FOR UPDATE (lock_version orders là migrate-later).
-- Idempotency tạo đơn qua app.order_idempotency (1 Idempotency-Key → 1 đơn).
-- KHÔNG đụng kho/tiền/sổ (ShipOrder/RecordPayment/POS = P4b). Tiền NUMERIC →
-- common/money ở repo, KHÔNG float. quantity INTEGER (cột thật) → ghi từ int.
-- ============================================================================

-- name: NextOrderCodeSeq :one
-- Cấp số tăng dần kế tiếp cho mã đơn (app ghép tiền tố + zero-pad). Sequence bảo
-- đảm KHÔNG trùng, an toàn đua (mỗi nextval một giá trị duy nhất).
SELECT nextval('app.order_code_seq')::bigint AS seq;

-- name: FindOrderByIdemKey :one
-- Tra cứu đơn đã tạo theo Idempotency-Key (chống tạo trùng). Có → trả order_id +
-- code (app đọc lại đơn cũ). Không có → ErrNoRows (được phép tạo mới).
SELECT order_id, order_code
FROM app.order_idempotency
WHERE idem_key = @idem_key::text;

-- name: InsertOrderIdemKey :exec
-- Ghi ánh xạ Idempotency-Key → đơn vừa tạo. ON CONFLICT DO NOTHING chống race
-- (2 luồng cùng key): luồng thua không ghi đè, app phát hiện qua RowsAffected/
-- re-SELECT để trả đơn của luồng thắng. KHÔNG cộng/ghi gì thêm.
INSERT INTO app.order_idempotency (idem_key, order_id, order_code)
VALUES (@idem_key::text, @order_id::uuid, @order_code::text)
ON CONFLICT (idem_key) DO NOTHING;

-- name: InsertOrder :one
-- Ghi MỘT đơn (status PENDING). id + code do APP sinh (uuid v7 + tiền tố/seq) →
-- truyền vào (KHÔNG dùng default DB). customer_id/creator_id/note nullable. total_
-- amount/final_amount NUMERIC (đã tính ở domain). created_at/updated_at để DB
-- default now(). Trùng code (UNIQUE WHERE deleted_at IS NULL) → unique violation
-- (app xử lý). RETURNING các cột cần để map về domain + đường đọc.
INSERT INTO public.orders (
    id, code, customer_id, creator_id, status, order_type,
    total_amount, final_amount, payment_status, note
) VALUES (
    @id::uuid, @code::text, sqlc.narg('customer_id')::bigint,
    sqlc.narg('creator_id')::uuid, @status::text, @order_type::text,
    @total_amount::numeric, @final_amount::numeric, 'unpaid', sqlc.narg('note')::text
)
RETURNING id, code, customer_id, creator_id, status, order_type,
    total_amount, final_amount, payment_status, note, created_at, updated_at;

-- name: InsertOrderItem :exec
-- Ghi MỘT dòng hàng của đơn. id app sinh (uuid v7). quantity là INTEGER (cột
-- thật). unit_price/discount/total_line NUMERIC (đã tính ở domain). is_gift mặc
-- định false ở P4a (chưa làm hàng tặng). created_at để DB default now().
INSERT INTO public.order_items (
    id, order_id, product_id, quantity, uom, unit_price, discount, total_line, is_gift
) VALUES (
    @id::uuid, @order_id::uuid, @product_id::bigint, @quantity::integer,
    @uom::text, @unit_price::numeric, @discount::numeric, @total_line::numeric, false
);

-- name: GetOrderStatusForUpdate :one
-- Lấy status hiện tại của đơn (chưa xoá mềm) và KHOÁ dòng (FOR UPDATE) trong tx
-- để đổi trạng thái an toàn (tuần tự hoá các thao tác đồng thời trên cùng đơn).
-- Không thấy → ErrNoRows (app map NotFound).
SELECT id, status
FROM public.orders
WHERE id = @order_id::uuid AND deleted_at IS NULL
FOR UPDATE;

-- name: UpdateOrderStatus :execrows
-- Cập nhật trạng thái đơn (+ updated_at). Guard status cũ trong WHERE (optimistic
-- theo giá trị đã đọc dưới FOR UPDATE) — RowsAffected=0 nếu đơn đã đổi trạng thái
-- bởi luồng khác (app map Conflict). KHÔNG ghi đè đơn đã xoá mềm.
UPDATE public.orders
SET status = @new_status::text, updated_at = now()
WHERE id = @order_id::uuid
  AND status = @expected_status::text
  AND deleted_at IS NULL;

-- ============================================================================
-- ĐƯỜNG GHI ORCHESTRATION (P4b — gộp cross-module trong 1 tx: ShipOrder /
-- RecordPayment / CreatePosSale). Trừ kho FEFO + phát HĐ VAT + post sổ kép +
-- ghi phiếu thu gọi qua PORT module khác (inventory/vat/accounting/finance) —
-- KHÔNG có SQL ở đây cho các bảng đó. Riêng orders cần: KHOÁ + đọc header đầy đủ
-- (FOR UPDATE) và cập nhật payment_status suy diễn.
-- ============================================================================

-- name: GetOrderHeaderForUpdate :one
-- Lấy header đơn ĐẦY ĐỦ (code/status/order_type/customer_id/final_amount) + KHOÁ
-- dòng (FOR UPDATE) trong tx — dùng cho ShipOrder/RecordPayment để vừa khoá vừa
-- đọc các trường cần điều phối (mã đơn cho ref/HĐ, loại đơn B2B/B2C cho luồng
-- post sổ, tổng phải trả để tính payment_status). Không thấy → ErrNoRows
-- (app map NotFound). KHÔNG khoá theo paid_amount (suy diễn động, đọc riêng).
SELECT id, code, customer_id, creator_id, status, order_type,
    total_amount, final_amount, payment_status
FROM public.orders
WHERE id = @order_id::uuid AND deleted_at IS NULL
FOR UPDATE;

-- name: SumOrderPaidInTx :one
-- Tổng ĐÃ THU của đơn theo mã (ref_id = orders.code) TRONG tx hiện hành — phiếu/
-- allocation vừa ghi cùng tx ĐƯỢC tính (visibility trong tx). Gọi NGUỒN CHÂN LÝ
-- app.order_paid_amount (phiếu trực tiếp + phân bổ; flow='in', status pending+
-- completed, sổ INTERNAL/BOTH). Dùng tính lại payment_status sau RecordPaymentIn.
SELECT app.order_paid_amount(@order_code::text)::numeric AS paid_amount;

-- name: UpdateOrderPaymentStatus :execrows
-- Cập nhật payment_status suy diễn (unpaid/partial/paid) + updated_at. KHÔNG đụng
-- status xử lý đơn. Chỉ ghi đơn còn sống. RowsAffected=0 nếu đơn biến mất (phòng thủ).
UPDATE public.orders
SET payment_status = @payment_status::text, updated_at = now()
WHERE id = @order_id::uuid
  AND deleted_at IS NULL;

-- ============================================================================
-- PHÂN BỔ 1 phiếu THU → NHIỀU đơn (spec mục 55). Phiếu lump-sum (ref_type='customer')
-- phân bổ cho các đơn chưa tất toán, CŨ NHẤT trước. "Đã thu mỗi đơn" = phiếu trực
-- tiếp + allocation (qua app.order_paid_amount).
-- ============================================================================

-- name: InsertOrderAllocation :exec
-- Ghi MỘT dòng phân bổ phiếu (payment_id) cho một đơn (order_code) số tiền amount.
-- 1 phiếu × 1 đơn tối đa một dòng (UNIQUE). ON CONFLICT DO NOTHING → IDEMPOTENT: gọi
-- lại (replay) KHÔNG cộng dồn (tránh nhân đôi "đã thu"). amount > 0 (CHECK ở bảng).
INSERT INTO app.finance_transaction_allocations (payment_id, order_code, amount)
VALUES (@payment_id::bigint, @order_code::text, @amount::numeric)
ON CONFLICT (payment_id, order_code) DO NOTHING;

-- name: ListUnpaidOrdersByCustomerForUpdate :many
-- Đơn CHƯA tất toán (payment_status <> 'paid') của một khách, CŨ NHẤT trước (mục
-- 55) + KHOÁ dòng (FOR UPDATE) để phân bổ tuần tự an toàn. Kèm final_amount + đã-thu
-- hiện tại (app.order_paid_amount) để app tính phần còn thiếu mỗi đơn.
SELECT o.id, o.code, o.final_amount, app.order_paid_amount(o.code) AS paid_amount
FROM public.orders o
WHERE o.customer_id = @customer_id::bigint
  AND o.payment_status <> 'paid'
  AND o.deleted_at IS NULL
ORDER BY o.created_at ASC, o.id ASC
FOR UPDATE;
