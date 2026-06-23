-- Purchasing (mục 54 — chiều MUA, đối xứng orders chiều BÁN). GHI/ĐỌC bảng MỚI ở
-- schema app (purchase_orders + purchase_order_items, goose 00007). Tiền NUMERIC →
-- common/money ở repo (KHÔNG float); quantity NUMERIC(20,4) → inventory/domain.Quantity
-- (decimal); vat_rate NUMERIC(6,4) → shopspring/decimal. Mã PO app sinh từ
-- app.purchase_order_code_seq + tiền tố. Đổi trạng thái dùng FOR UPDATE + guard
-- status cũ (state machine THUẦN ở domain). Nhập kho/post sổ/chi NCC KHÔNG ở đây —
-- purchasing/app gọi PORT inventory/accounting/finance với CÙNG tx (gộp atomic).

-- name: NextPurchaseOrderCodeSeq :one
-- Cấp số tăng dần cho mã PO (app ghép tiền tố + zero-pad). An toàn đua.
SELECT nextval('app.purchase_order_code_seq')::bigint AS seq;

-- name: InsertPurchaseOrder :one
-- Ghi header PO (status draft). id uuid v7 + code do app truyền. Trùng code (đua sinh
-- mã, rất hiếm) → unique violation (ux_purchase_orders_code_alive) → app cấp số mới.
INSERT INTO app.purchase_orders (
    id, code, supplier_id, supplier_name, status, total_amount, vat_amount, note
) VALUES (
    @id::uuid, @code::text, sqlc.narg('supplier_id')::bigint,
    sqlc.narg('supplier_name')::text, @status::text,
    @total_amount::numeric, @vat_amount::numeric, sqlc.narg('note')::text
)
RETURNING id, code, supplier_id, supplier_name, status, total_amount, vat_amount,
          note, lock_version, created_at, updated_at;

-- name: InsertPurchaseOrderItem :exec
-- Ghi một dòng hàng PO. id uuid v7 do app truyền. line_no thứ tự dòng (UNIQUE po_id,line_no).
INSERT INTO app.purchase_order_items (
    id, po_id, line_no, product_id, quantity, unit_cost, vat_rate,
    batch_code, expiry_date, manufacturing_date, line_total
) VALUES (
    @id::uuid, @po_id::uuid, @line_no::int, @product_id::bigint,
    @quantity::numeric, @unit_cost::numeric, @vat_rate::numeric,
    sqlc.narg('batch_code')::text, sqlc.narg('expiry_date')::date,
    sqlc.narg('manufacturing_date')::date, @line_total::numeric
);

-- name: GetPurchaseOrderHeaderForUpdate :one
-- Khoá dòng PO (FOR UPDATE) + trả header cho điều phối (confirm/receive/pay). Không
-- thấy / đã xoá mềm → ErrNoRows.
SELECT id, code, supplier_id, supplier_name, status, total_amount, vat_amount,
       note, lock_version, created_at, updated_at
FROM app.purchase_orders
WHERE id = @id::uuid AND deleted_at IS NULL
FOR UPDATE;

-- name: ListPurchaseOrderLines :many
-- Các dòng hàng của một PO (theo line_no) — product/qty/unit_cost/vat/lô để nhập kho
-- + post sổ. Hết → slice rỗng.
SELECT id, po_id, line_no, product_id, quantity, unit_cost, vat_rate,
       batch_code, expiry_date, manufacturing_date, line_total
FROM app.purchase_order_items
WHERE po_id = @po_id::uuid
ORDER BY line_no ASC;

-- name: UpdatePurchaseOrderStatus :execrows
-- Đổi trạng thái PO (guard status cũ + bump lock_version). Trả số dòng đổi — 0 nghĩa
-- là PO đã đổi trạng thái bởi luồng khác (service map Conflict).
UPDATE app.purchase_orders
SET status = @new_status::text,
    lock_version = lock_version + 1,
    updated_at = now()
WHERE id = @id::uuid
  AND status = @expected_status::text
  AND deleted_at IS NULL;

-- name: GetPurchaseOrderByID :one
-- Đọc lại PO (header) theo id — đường ĐỌC (GET detail) + readback sau ghi. Không
-- thấy → ErrNoRows.
SELECT id, code, supplier_id, supplier_name, status, total_amount, vat_amount,
       note, lock_version, created_at, updated_at
FROM app.purchase_orders
WHERE id = @id::uuid AND deleted_at IS NULL;

-- name: ListPurchaseOrders :many
-- Danh sách PO còn sống (keyset created_at DESC, id DESC). Lọc optional status +
-- supplier_id. LIMIT phân trang.
SELECT id, code, supplier_id, supplier_name, status, total_amount, vat_amount,
       note, lock_version, created_at, updated_at
FROM app.purchase_orders
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('supplier_id')::bigint IS NULL OR supplier_id = sqlc.narg('supplier_id')::bigint)
  AND (
        NOT @has_cursor::boolean
        OR (created_at, id) < (@after_created_at::timestamptz, @after_id::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT @row_limit::int;

-- name: InsertPurchaseOrderIdemKey :exec
-- Bind Idempotency-Key tạo PO → 1 PO. ON CONFLICT DO NOTHING; app đọc lại để biết
-- thắng đua hay không (tái dùng pattern app.order_idempotency của orders).
INSERT INTO app.purchase_order_idempotency (idem_key, po_id, po_code)
VALUES (@idem_key::text, @po_id::uuid, @po_code::text)
ON CONFLICT (idem_key) DO NOTHING;

-- name: FindPurchaseOrderByIdemKey :one
-- Tra PO đã tạo theo Idempotency-Key (idempotent tạo PO). Không thấy → ErrNoRows.
SELECT po_id, po_code FROM app.purchase_order_idempotency
WHERE idem_key = @idem_key::text;
