-- SCHEMA THAM CHIẾU (strangler-fig, ADR 0001) — KHÔNG phải goose migration.
--
-- Mô tả các bảng public.* ĐANG TỒN TẠI (Postgres Supabase, nhánh monorepo
-- office.duoc.namviet) liên quan bounded context inventory, để (1) sqlc
-- type-check query inventory ở compile-time và (2) dbtest materialize bảng + seed
-- trong testcontainers (DB test trống). TUYỆT ĐỐI KHÔNG chạy qua goose lên prod —
-- backend chỉ ĐỌC các bảng kế thừa này, không sở hữu chúng.
--
-- Cột lấy nguyên văn từ office.duoc.namviet/database_schema.md (doc). Đây là
-- nguồn DOC, PHẢI verify lại với prod thật (pg_dump/REST Object.keys) khi có
-- creds — prod hay lệch migration.
--
-- GHI CHÚ KIỂU/RÀNG BUỘC QUAN TRỌNG (đọc kỹ trước khi đụng kho):
--   * id mọi bảng kho = bigint (KHÔNG uuid), TRỪ inventory_transactions.id = uuid.
--   * quantity / stock_quantity = NUMERIC (KHÔNG phải int) ở DB này. ADR 0002 §B
--     coi "thống nhất kiểu quantity (int vs numeric)" là việc migrate sau →
--     ở đây map NUMERIC → decimal (KHÔNG ép float), không tự đổi kiểu.
--   * Tồn TỔNG theo kho: public.product_inventory.stock_quantity. KHÔNG có
--     deleted_at ở bảng này (không soft-delete dòng tồn) — chỉ lọc warehouse hợp lệ.
--   * Tồn theo LÔ tại kho: public.inventory_batches.quantity (quantity tồn của
--     (warehouse, product, batch)). KHÔNG có deleted_at; hạn dùng/giá nhập nằm ở
--     bảng public.batches (join qua batch_id). FEFO = sort batches.expiry_date ASC.
--   * Lô hàng (hạn dùng + giá nhập): public.batches — CÓ deleted_at (soft-delete
--     lô). FEFO chỉ xét lô deleted_at IS NULL và còn tồn > 0.
--   * Giá vốn lô (inbound_price) = numeric → common/money ở repo (KHÔNG float).

CREATE SCHEMA IF NOT EXISTS public;

-- Kho hàng / chi nhánh / cửa hàng bán lẻ. id = bigint. KHÔNG có updated_at.
-- Có deleted_at (đóng cửa kho) → lọc deleted_at IS NULL.
CREATE TABLE IF NOT EXISTS public.warehouses (
    id          bigint PRIMARY KEY,
    key         text NOT NULL,
    name        text NOT NULL,
    unit        text NOT NULL DEFAULT 'Hộp',
    address     text,
    type        text NOT NULL DEFAULT 'retail',
    latitude    numeric,
    longitude   numeric,
    code        text,
    manager     text,
    phone       text,
    status      text NOT NULL DEFAULT 'active',
    company_id  uuid,
    outlet_type text,
    created_at  timestamptz DEFAULT now(),
    deleted_at  timestamptz
);

-- Tồn kho TỔNG hợp tại từng kho. id = bigint. stock_quantity = NUMERIC.
-- KHÔNG có deleted_at. product_id/warehouse_id nullable ở DB này.
CREATE TABLE IF NOT EXISTS public.product_inventory (
    id              bigint PRIMARY KEY,
    product_id      bigint,
    warehouse_id    bigint,
    stock_quantity  numeric NOT NULL DEFAULT 0,
    min_stock       integer DEFAULT 0,
    max_stock       integer DEFAULT 0,
    shelf_location  text DEFAULT 'Chưa xếp',
    location_cabinet text,
    location_row    text,
    location_slot   text,
    updated_by      uuid,
    updated_at      timestamptz DEFAULT now()
);

-- Lô sản xuất + hạn sử dụng của sản phẩm. id = bigint. CÓ deleted_at.
-- expiry_date = DATE (FEFO sort theo cột này). inbound_price = NUMERIC (giá vốn).
CREATE TABLE IF NOT EXISTS public.batches (
    id                 bigint PRIMARY KEY,
    product_id         bigint NOT NULL,
    batch_code         text NOT NULL,
    expiry_date        date NOT NULL,
    manufacturing_date date,
    inbound_price      numeric DEFAULT 0,
    created_at         timestamptz DEFAULT now(),
    updated_at         timestamptz DEFAULT now(),
    deleted_at         timestamptz
);

-- Tồn theo LÔ tại từng kho (warehouse, product, batch). id = bigint.
-- quantity = NUMERIC. KHÔNG có deleted_at (soft-delete ở bảng batches).
CREATE TABLE IF NOT EXISTS public.inventory_batches (
    id           bigint PRIMARY KEY,
    warehouse_id bigint NOT NULL,
    product_id   bigint NOT NULL,
    batch_id     bigint NOT NULL,
    quantity     numeric NOT NULL DEFAULT 0,
    updated_at   timestamptz DEFAULT now()
);
