-- VAT (P5 — hoá đơn GTGT). GHI object MỚI ở schema app (sales_invoices/_lines/
-- invoice_serials). Tiền = NUMERIC → common/money ở repo (KHÔNG float). PK uuid
-- sinh app-side (common/id v7) và truyền vào INSERT. Phát hành HĐ chạy TRONG tx
-- do caller (orders) truyền (gộp atomic với giao hàng): service cấp số gapless
-- (FOR UPDATE serial) → insert header → insert từng line.

-- name: FindIssuedInvoiceByOrder :one
-- Idempotency 1 đơn 1 HĐ: tìm HĐ đã 'issued' cho order_code. Có → service trả
-- HĐ cũ (no-op, không phát hành trùng). Không thấy → pgx.ErrNoRows.
SELECT id, order_code, customer_tax_code, serial, invoice_no, issue_date,
       subtotal, vat_amount, total, status, lock_version, created_at
FROM app.sales_invoices
WHERE order_code = @order_code::text
  AND status = 'issued';

-- name: NextInvoiceNo :one
-- Cấp số HĐ GAPLESS theo serial: khoá dòng serial (FOR UPDATE — tuần tự hoá các
-- tx cùng serial), lấy next_no hiện tại làm số HĐ, tăng next_no lên 1 NGAY trong
-- CÙNG tx. Trả số vừa cấp. Phải gọi TRONG tx phát hành (cùng tx insert header)
-- để gapless + atomic: nếu tx rollback thì next_no không bị tiêu hao.
-- Dòng serial phải tồn tại trước (EnsureSerial) — không thì ErrNoRows.
UPDATE app.invoice_serials
SET next_no = next_no + 1
WHERE serial = @serial::text
RETURNING (next_no - 1)::bigint AS issued_no;

-- name: EnsureSerial :exec
-- Tạo dòng serial nếu chưa có (idempotent). next_no khởi 1. mau_so optional.
-- Gọi trước NextInvoiceNo để đảm bảo dòng tồn tại cho FOR UPDATE.
INSERT INTO app.invoice_serials (serial, mau_so, next_no)
VALUES (@serial::text, COALESCE(sqlc.narg('mau_so')::text, ''), 1)
ON CONFLICT (serial) DO NOTHING;

-- name: InsertSalesInvoice :one
-- Ghi header HĐ (status 'issued'). id sinh app-side. invoice_no đã cấp gapless.
-- UNIQUE(serial,invoice_no) + unique index một phần (order_code WHERE issued)
-- chặn trùng (race) → repo bắt 23505. RETURNING id để service insert lines.
INSERT INTO app.sales_invoices (
    id, order_code, customer_tax_code, serial, invoice_no, issue_date,
    subtotal, vat_amount, total, status
) VALUES (
    @id::uuid, @order_code::text, @customer_tax_code::text, @serial::text,
    @invoice_no::bigint, @issue_date::date,
    @subtotal::numeric, @vat_amount::numeric, @total::numeric, 'issued'
)
RETURNING id;

-- name: InsertSalesInvoiceLine :exec
-- Ghi một dòng HĐ (line_no theo thứ tự). Tiền NUMERIC; vat_rate là input. id
-- sinh app-side. UNIQUE(invoice_id,line_no) chống trùng dòng.
INSERT INTO app.sales_invoice_lines (
    id, invoice_id, line_no, product_id, description,
    quantity, unit_price, vat_rate, line_amount, line_vat
) VALUES (
    @id::uuid, @invoice_id::uuid, @line_no::int, @product_id::bigint, @description::text,
    @quantity::numeric, @unit_price::numeric, @vat_rate::numeric,
    @line_amount::numeric, @line_vat::numeric
);

-- name: GetSalesInvoice :one
-- Một HĐ theo id (cho HTTP đọc /v1/vat/invoices/{id}).
SELECT id, order_code, customer_tax_code, serial, invoice_no, issue_date,
       subtotal, vat_amount, total, status, lock_version, created_at
FROM app.sales_invoices
WHERE id = @id::uuid;

-- name: ListSalesInvoiceLines :many
-- Các dòng của một HĐ, thứ tự theo line_no (ổn định).
SELECT id, invoice_id, line_no, product_id, description,
       quantity, unit_price, vat_rate, line_amount, line_vat
FROM app.sales_invoice_lines
WHERE invoice_id = @invoice_id::uuid
ORDER BY line_no ASC, id ASC;

-- name: ListSalesInvoices :many
-- Danh sách HĐ, keyset theo (created_at DESC, id DESC): trang kế lấy HĐ "cũ hơn"
-- mốc cursor. @after_created_at NULL = trang đầu. Lọc optional order_code/status.
SELECT id, order_code, customer_tax_code, serial, invoice_no, issue_date,
       subtotal, vat_amount, total, status, lock_version, created_at
FROM app.sales_invoices
WHERE (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR created_at < sqlc.narg('after_created_at')::timestamptz
        OR (created_at = sqlc.narg('after_created_at')::timestamptz
            AND id < sqlc.narg('after_id')::uuid)
      )
  AND (sqlc.narg('order_code')::text IS NULL OR order_code = sqlc.narg('order_code')::text)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY created_at DESC, id DESC
LIMIT @row_limit::int;
