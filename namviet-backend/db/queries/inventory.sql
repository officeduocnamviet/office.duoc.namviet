-- Inventory (LÁT ĐỌC, ADR 0001): ĐỌC bảng public.* kế thừa. CHƯA có đường GHI
-- (trừ kho FEFO/nhập/chuyển/kiểm kê HOÃN — cần khoá tranh chấp + tx tiền, làm
-- sau module orders). quantity/stock_quantity là NUMERIC → map sang kiểu decimal
-- (domain.Quantity) ở repo, KHÔNG dùng float. inbound_price (giá vốn lô) là
-- numeric → common/money ở repo, KHÔNG float.

-- name: ListWarehouses :many
-- Danh sách kho/chi nhánh còn hoạt động (lọc deleted_at IS NULL). Lọc optional
-- status. Sắp theo id ASC cho ổn định. Không phân trang (số kho nhỏ, có chặn).
SELECT
    id, key, name, unit, address, type, code, manager, phone, status,
    company_id, outlet_type, created_at
FROM public.warehouses
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY id ASC
LIMIT @row_limit::int;

-- name: ListStock :many
-- Tồn TỔNG theo kho (public.product_inventory.stock_quantity). Keyset theo id ASC
-- (id > after_id). Lọc optional product_id và/hoặc warehouse_id. Chỉ lấy dòng có
-- product_id và warehouse_id NOT NULL (dòng tồn hợp lệ) — cột này nullable ở DB.
-- KHÔNG có deleted_at ở product_inventory; nhưng loại tồn của kho đã đóng
-- (warehouses.deleted_at) qua JOIN để không lộ tồn kho không còn dùng.
SELECT
    pi.id,
    pi.product_id,
    pi.warehouse_id,
    pi.stock_quantity,
    pi.min_stock,
    pi.max_stock,
    pi.shelf_location,
    pi.updated_at
FROM public.product_inventory pi
JOIN public.warehouses w ON w.id = pi.warehouse_id AND w.deleted_at IS NULL
WHERE pi.id > @after_id::bigint
  AND pi.product_id IS NOT NULL
  AND (sqlc.narg('product_id')::bigint IS NULL OR pi.product_id = sqlc.narg('product_id')::bigint)
  AND (sqlc.narg('warehouse_id')::bigint IS NULL OR pi.warehouse_id = sqlc.narg('warehouse_id')::bigint)
ORDER BY pi.id ASC
LIMIT @row_limit::int;

-- name: ListBatchesFEFO :many
-- Danh sách LÔ còn tồn của một product theo FEFO (First-Expired-First-Out): sắp
-- xếp HẠN DÙNG (batches.expiry_date) TĂNG DẦN — lô hết hạn trước đứng trước. Dữ
-- liệu nền cho xuất kho FEFO sau này. Join inventory_batches (tồn theo lô tại kho)
-- với batches (hạn dùng + giá nhập). Lọc lô đã xóa (batches.deleted_at IS NULL),
-- chỉ lấy tồn > 0. Optional lọc warehouse_id. Tie-break theo batch_id rồi
-- inventory_batches.id để thứ tự ổn định (deterministic).
SELECT
    ib.id            AS inventory_batch_id,
    ib.warehouse_id,
    ib.product_id,
    ib.batch_id,
    ib.quantity,
    b.batch_code,
    b.expiry_date,
    b.manufacturing_date,
    b.inbound_price
FROM public.inventory_batches ib
JOIN public.batches b ON b.id = ib.batch_id AND b.deleted_at IS NULL
WHERE ib.product_id = @product_id::bigint
  AND ib.quantity > 0
  AND (sqlc.narg('warehouse_id')::bigint IS NULL OR ib.warehouse_id = sqlc.narg('warehouse_id')::bigint)
ORDER BY b.expiry_date ASC, ib.batch_id ASC, ib.id ASC;

-- ============================================================================
-- ĐƯỜNG GHI (P2 — trừ kho FEFO, có khoá tranh chấp). Chạy TRONG tx do caller
-- (orders/POS) truyền — gộp atomic với post sổ + ghi tiền. GHI bảng public.*
-- (strangler, ADR 0001): inventory_batches.quantity + product_inventory.stock_quantity.
-- lock_version các bảng này là migrate-later → dùng SELECT ... FOR UPDATE +
-- pg_advisory_xact_lock để tuần tự hoá, chống bán âm khi đồng thời. KHÔNG float
-- (quantity NUMERIC → decimal ở repo).
-- ============================================================================

-- name: LockWarehouseProduct :exec
-- Khoá tranh chấp ĐẦU TX cho một (warehouse, product): pg_advisory_xact_lock giữ
-- tới khi tx commit/rollback. Tuần tự hoá MỌI giao dịch trừ kho cùng (kho,sp) →
-- chỉ một tx trừ tại một thời điểm (chống bán âm). hashtext gộp 2 khoá thành 1
-- chuỗi 'warehouse:product' rồi băm sang bigint cho advisory lock (1 tham số).
SELECT pg_advisory_xact_lock(
    hashtext(@warehouse_id::bigint || ':' || @product_id::bigint)
);

-- name: ListBatchesForDeductFEFO :many
-- Các lô CÒN TỒN của (warehouse, product) theo FEFO (expiry ASC) ĐỂ TRỪ, giữ dòng
-- bằng FOR UPDATE OF ib (khoá dòng tồn-theo-lô trong tx — không khoá batches để
-- khỏi chặn đọc lô). Lọc lô đã xóa (batches.deleted_at IS NULL) + tồn > 0. Khác
-- ListBatchesFEFO (đọc): BẮT BUỘC warehouse_id (trừ kho theo đúng 1 kho) + FOR
-- UPDATE. Tie-break batch_id rồi ib.id cho thứ tự ổn định.
SELECT
    ib.id            AS inventory_batch_id,
    ib.warehouse_id,
    ib.product_id,
    ib.batch_id,
    ib.quantity,
    b.batch_code,
    b.expiry_date,
    b.manufacturing_date,
    b.inbound_price
FROM public.inventory_batches ib
JOIN public.batches b ON b.id = ib.batch_id AND b.deleted_at IS NULL
WHERE ib.warehouse_id = @warehouse_id::bigint
  AND ib.product_id = @product_id::bigint
  AND ib.quantity > 0
ORDER BY b.expiry_date ASC, ib.batch_id ASC, ib.id ASC
FOR UPDATE OF ib;

-- name: DeductInventoryBatch :exec
-- Trừ tồn của MỘT dòng tồn-theo-lô (inventory_batches) theo kế hoạch FEFO. Guard
-- quantity >= @qty trong WHERE để KHÔNG bao giờ ghi âm dù có race (advisory lock
-- đã chống, đây là phòng thủ tầng DB). Caller đã FOR UPDATE dòng này trong tx.
UPDATE public.inventory_batches
SET quantity = quantity - @qty::numeric,
    updated_at = now()
WHERE id = @inventory_batch_id::bigint
  AND quantity >= @qty::numeric;

-- name: DeductProductInventory :exec
-- Trừ tồn TỔNG (product_inventory.stock_quantity) của (warehouse, product) đúng
-- tổng đã trừ qua các lô. Cập nhật theo (product_id, warehouse_id) — KHÓA tự nhiên
-- của dòng tồn tổng (FOR UPDATE đã giữ dòng lô; advisory lock tuần tự hoá cả cụm).
UPDATE public.product_inventory
SET stock_quantity = stock_quantity - @qty::numeric,
    updated_at = now()
WHERE product_id = @product_id::bigint
  AND warehouse_id = @warehouse_id::bigint;

-- ============================================================================
-- ĐƯỜNG GHI (NHẬP KHO — StockIn, đối xứng DeductFEFO). purchasing/POS gọi TRONG tx
-- của họ (gộp atomic với post sổ + chi NCC). GHI bảng public.* (strangler, ADR 0001):
-- batches (tạo LÔ mới + inbound_price = giá nhập per-unit), inventory_batches (tồn
-- theo lô tại kho), product_inventory (tồn tổng — UPSERT). Advisory lock theo
-- (warehouse, product) như DeductFEFO. KHÔNG float (numeric → decimal/money ở repo).
-- ============================================================================

-- name: InsertBatch :one
-- Tạo MỘT lô mới (public.batches) khi nhập kho. id do DB sinh? KHÔNG: bảng cũ id =
-- bigint NHƯNG ở schema tham chiếu test không có IDENTITY → dùng nextval của sequence
-- ngầm? Để an toàn cross-env, lấy id = COALESCE(max+1) trong tx (advisory lock đã
-- tuần tự hoá (kho,sp) — chống đua id cho cùng product). inbound_price = giá nhập
-- per-unit của dòng PO. RETURNING id cho inventory_batches tham chiếu.
INSERT INTO public.batches (id, product_id, batch_code, expiry_date, manufacturing_date, inbound_price)
VALUES (
    (SELECT COALESCE(MAX(id), 0) + 1 FROM public.batches),
    @product_id::bigint, @batch_code::text,
    sqlc.narg('expiry_date')::date, sqlc.narg('manufacturing_date')::date,
    @inbound_price::numeric
)
RETURNING id;

-- name: InsertInventoryBatch :exec
-- Tạo dòng tồn-theo-lô (public.inventory_batches) cho (warehouse, product, batch) vừa
-- tạo. id = COALESCE(max+1) (đồng nhất cách InsertBatch — bảng cũ id bigint không
-- IDENTITY ở schema test). quantity = lượng nhập (NUMERIC).
INSERT INTO public.inventory_batches (id, warehouse_id, product_id, batch_id, quantity)
VALUES (
    (SELECT COALESCE(MAX(id), 0) + 1 FROM public.inventory_batches),
    @warehouse_id::bigint, @product_id::bigint, @batch_id::bigint, @quantity::numeric
);

-- Tồn TỔNG nhập kho là UPSERT thủ công (bảng public.product_inventory KHÔNG có unique
-- (warehouse,product) ở schema tham chiếu): repo gọi AddProductInventoryStock trước
-- (cộng dồn nếu đã có dòng), nếu 0 dòng → InsertProductInventoryStock (tạo dòng mới).

-- name: AddProductInventoryStock :execrows
-- Cộng dồn tồn tổng cho (warehouse, product) ĐÃ CÓ dòng. Trả số dòng đổi: 0 → repo
-- INSERT dòng mới (chưa có tồn cho cặp này).
UPDATE public.product_inventory
SET stock_quantity = stock_quantity + @qty::numeric,
    updated_at = now()
WHERE product_id = @product_id::bigint
  AND warehouse_id = @warehouse_id::bigint;

-- name: InsertProductInventoryStock :exec
-- Tạo dòng tồn tổng mới cho (warehouse, product) chưa có. id = COALESCE(max+1)
-- (bảng id bigint không IDENTITY ở schema test). stock_quantity = qty nhập ban đầu.
INSERT INTO public.product_inventory (id, product_id, warehouse_id, stock_quantity)
VALUES (
    (SELECT COALESCE(MAX(id), 0) + 1 FROM public.product_inventory),
    @product_id::bigint, @warehouse_id::bigint, @qty::numeric
);
