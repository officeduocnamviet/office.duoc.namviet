-- Hoá đơn GTGT (VAT) — object MỚI do backend sở hữu (schema app, ADR 0002).
-- HĐ thuộc SỔ TAX (theo giá hoá đơn — dual-ledger). Cấp số GAPLESS (liên tục,
-- không nhảy số) theo từng ký hiệu (serial) — yêu cầu hoá đơn điện tử VN (NĐ
-- 123/TT78): SELECT ... FOR UPDATE dòng app.invoice_serials → dùng next_no →
-- next_no = next_no + 1 trong CÙNG tx (atomic + tuần tự hoá theo serial).
-- Tiền = NUMERIC(20,0) (VND, scale-0). PK uuid v7 sinh app-side (common/id).
-- 1 đơn 1 HĐ ở giai đoạn này (UNIQUE order_code khi status='issued', xem index
-- một phần bên dưới). Phát hành điện tử qua provider (VNPT/Viettel/MISA) = DEFER.
--
-- PK uuid: DEFAULT gen_random_uuid() chỉ là LƯỚI AN TOÀN; app sinh uuid v7 và
-- truyền vào INSERT (time-ordered, tốt cho index theo thời gian).

-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

-- invoice_serials: cấp số gapless theo ký hiệu HĐ. mau_so = mẫu số HĐ (NĐ123,
-- vd "1"); serial = ký hiệu (vd "C26TYY"). next_no = số HĐ kế tiếp sẽ cấp.
CREATE TABLE app.invoice_serials (
    serial   text   PRIMARY KEY,
    mau_so   text   NOT NULL DEFAULT '',
    next_no  bigint NOT NULL DEFAULT 1 CHECK (next_no >= 1)
);

CREATE TABLE app.sales_invoices (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    -- order_code = orders.code (text, KHÔNG uuid — khớp finance_transactions.ref_id).
    order_code        text   NOT NULL,
    -- MST khách — BẮT BUỘC cho HĐ VAT B2B (NOT NULL + CHECK không rỗng).
    customer_tax_code text   NOT NULL CHECK (length(trim(customer_tax_code)) > 0),
    serial            text   NOT NULL,
    -- invoice_no: số HĐ gapless theo serial (UNIQUE(serial, invoice_no)).
    invoice_no        bigint NOT NULL CHECK (invoice_no >= 1),
    issue_date        date   NOT NULL,
    subtotal          numeric(20,0) NOT NULL CHECK (subtotal   >= 0),
    vat_amount        numeric(20,0) NOT NULL CHECK (vat_amount >= 0),
    total             numeric(20,0) NOT NULL CHECK (total      >= 0),
    status            text   NOT NULL DEFAULT 'issued' CHECK (status IN ('draft','issued','cancelled')),
    lock_version      int    NOT NULL DEFAULT 0,
    created_at        timestamptz NOT NULL DEFAULT now(),
    -- total = subtotal + vat_amount (ép cân ở DB — phòng thủ ngoài service/domain).
    CONSTRAINT chk_invoice_total CHECK (total = subtotal + vat_amount),
    -- Số HĐ duy nhất theo từng ký hiệu (gapless theo serial).
    CONSTRAINT uq_invoice_serial_no UNIQUE (serial, invoice_no)
);
-- Keyset đọc /v1/vat/invoices theo (created_at DESC, id DESC).
CREATE INDEX idx_sales_invoices_created ON app.sales_invoices (created_at DESC, id DESC);
CREATE INDEX idx_sales_invoices_status ON app.sales_invoices (status);
-- 1 đơn 1 HĐ ở giai đoạn này: chặn 2 HĐ 'issued' cùng order_code (idempotency
-- tầng DB; HĐ 'cancelled' không tính để cho phép phát hành lại sau khi huỷ).
CREATE UNIQUE INDEX uq_sales_invoices_order_issued
    ON app.sales_invoices (order_code) WHERE status = 'issued';

CREATE TABLE app.sales_invoice_lines (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id  uuid NOT NULL REFERENCES app.sales_invoices(id) ON DELETE CASCADE,
    line_no     int  NOT NULL,
    product_id  bigint NOT NULL DEFAULT 0,
    description text NOT NULL DEFAULT '',
    quantity    numeric(20,4) NOT NULL CHECK (quantity   >= 0),
    unit_price  numeric(20,0) NOT NULL CHECK (unit_price >= 0),
    -- vat_rate là INPUT từng dòng (vd 0.05/0.08/0.10) — KHÔNG hardcode.
    vat_rate    numeric(6,4)  NOT NULL CHECK (vat_rate   >= 0),
    line_amount numeric(20,0) NOT NULL CHECK (line_amount >= 0),
    line_vat    numeric(20,0) NOT NULL CHECK (line_vat   >= 0),
    UNIQUE (invoice_id, line_no)
);
CREATE INDEX idx_sales_invoice_lines_invoice ON app.sales_invoice_lines (invoice_id);

-- +goose Down
DROP TABLE IF EXISTS app.sales_invoice_lines;
DROP TABLE IF EXISTS app.sales_invoices;
DROP TABLE IF EXISTS app.invoice_serials;
