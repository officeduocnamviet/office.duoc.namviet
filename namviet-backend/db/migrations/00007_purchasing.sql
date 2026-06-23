-- Mua hàng & nhập kho (spec system_features.md mục 54) — chiều MUA, đối xứng chiều
-- BÁN (orders). Object MỚI do backend sở hữu (schema app, ADR 0002): PK uuid v7,
-- tiền NUMERIC scale-0 (VND) / (20,4) cho qty/vat_rate, status = CHECK (không text
-- tự do), lock_version + updated_at cho aggregate tiền, CHECK qty/cost >= 0.
--
-- ⚠️ CỜ CẢNH BÁO PROD (verify trước cutover):
--   1) KHÔNG có bảng public.suppliers / public.purchase_orders trên Supabase hiện
--      tại (xác nhận từ office.duoc.namviet/database_schema.md). supplier_id (bigint)
--      + supplier_name (text) lưu THẲNG, KHÔNG FK — khi có supplier entity thật phải
--      verify + thêm FK/nguồn dữ liệu. manufacturers / vendor_product_mappings có
--      tồn tại nhưng KHÔNG phải "nhà cung cấp mua hàng" → không gắn vội.
--   2) Nhập kho ghi public.batches (inbound_price = unit_cost per-unit) +
--      public.inventory_batches + public.product_inventory (cột thật phải verify).
--   3) Quy ước mã PO (tiền tố "PO" + zero-pad) cấu hình ở internal/purchasing (app
--      code generator) — kế toán/BA xác nhận.
--
-- HOÃN (mục 54 nâng cao, ghi rõ trong docs): auto-tạo PO khi tồn < min, chương
-- trình NCC, hợp đồng, upload HĐ tự điền. CORE = đủ vòng đời PO: draft → ordered →
-- received (nhập kho + post sổ Dr 1561+133/Cr 331) → paid (chi NCC, Dr 331/Cr 111/112).

-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

-- order_code_seq cho mã đơn BÁN đã có (00005). PO dùng sequence RIÊNG để mã PO
-- (tiền tố "PO") độc lập mã đơn bán. nextval an toàn đua (mỗi gọi một giá trị).
CREATE SEQUENCE IF NOT EXISTS app.purchase_order_code_seq AS bigint START WITH 1 INCREMENT BY 1;

-- purchase_orders: đơn mua hàng (header). uuid v7 (app sinh). code app sinh từ
-- purchase_order_code_seq + tiền tố. supplier_id/supplier_name KHÔNG FK (cờ #1).
-- total_amount/vat_amount NUMERIC(20,0) (VND scale-0). status CHECK (state machine
-- THUẦN ép ở domain; CHECK là phòng thủ DB). lock_version + updated_at cho aggregate
-- tiền (ADR 0002). Soft-delete deleted_at (đồng nhất quy ước alive).
CREATE TABLE app.purchase_orders (
    id            uuid PRIMARY KEY,
    code          text NOT NULL,
    supplier_id   bigint,
    supplier_name text,
    status        text NOT NULL DEFAULT 'draft'
                  CHECK (status IN ('draft', 'ordered', 'received', 'paid', 'cancelled')),
    total_amount  numeric(20,0) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    vat_amount    numeric(20,0) NOT NULL DEFAULT 0 CHECK (vat_amount >= 0),
    note          text,
    lock_version  integer NOT NULL DEFAULT 0,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz
);
-- code DUY NHẤT cho PO còn sống (chống tạo trùng mã; app sinh từ sequence + tiền tố,
-- race/lỗi sinh → unique violation thay vì 2 PO cùng mã).
CREATE UNIQUE INDEX ux_purchase_orders_code_alive
    ON app.purchase_orders (code) WHERE deleted_at IS NULL;
-- Tra PO theo nhà cung cấp (danh sách/đối soát).
CREATE INDEX idx_purchase_orders_supplier ON app.purchase_orders (supplier_id);

-- purchase_order_items: dòng hàng của một PO. uuid v7 (app sinh). po_id FK CASCADE
-- (dòng thuộc header, cùng vòng đời). quantity NUMERIC(20,4) > 0; unit_cost giá nhập
-- per-unit NUMERIC(20,0) >= 0 (→ inbound_price lô khi nhập kho). vat_rate NUMERIC(6,4)
-- (vd 0.08). batch_code/expiry/mfg cho lô nhập. line_total = qty*unit_cost (VND).
-- UNIQUE(po_id, line_no) — thứ tự dòng ổn định.
CREATE TABLE app.purchase_order_items (
    id                 uuid PRIMARY KEY,
    po_id              uuid NOT NULL REFERENCES app.purchase_orders (id) ON DELETE CASCADE,
    line_no            integer NOT NULL,
    product_id         bigint NOT NULL,
    quantity           numeric(20,4) NOT NULL CHECK (quantity > 0),
    unit_cost          numeric(20,0) NOT NULL CHECK (unit_cost >= 0),
    vat_rate           numeric(6,4) NOT NULL DEFAULT 0 CHECK (vat_rate >= 0),
    batch_code         text,
    expiry_date        date,
    manufacturing_date date,
    line_total         numeric(20,0) NOT NULL CHECK (line_total >= 0),
    CONSTRAINT uq_poi_po_line UNIQUE (po_id, line_no)
);
CREATE INDEX idx_poi_po ON app.purchase_order_items (po_id);

-- purchase_order_idempotency: chống TẠO PO trùng theo Idempotency-Key (1 key → 1
-- PO). Tái dùng pattern app.order_idempotency (chiều bán). po_id/po_code trỏ
-- purchase_orders (KHÔNG FK cứng — giữ ID, đồng nhất quy ước order_idempotency).
CREATE TABLE app.purchase_order_idempotency (
    idem_key   text PRIMARY KEY CHECK (length(trim(idem_key)) > 0),
    po_id      uuid NOT NULL,
    po_code    text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS app.purchase_order_idempotency;
DROP TABLE IF EXISTS app.purchase_order_items;
DROP TABLE IF EXISTS app.purchase_orders;
DROP SEQUENCE IF EXISTS app.purchase_order_code_seq;
