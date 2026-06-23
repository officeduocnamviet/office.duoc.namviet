-- Sinh mã đơn (orders.code) + idempotency tạo đơn — object MỚI do backend sở hữu
-- (schema app, ADR 0002). public.orders.code là TEXT NOT NULL UNIQUE NHƯNG KHÔNG
-- có default DB → APP TỰ SINH. P4a dựng sequence riêng + tiền tố riêng để mã backend
-- sinh KHÔNG đụng UNIQUE với mã lịch sử ERP (vd 'HD...').
--
-- ⚠️ CỜ CẢNH BÁO PROD (xác nhận trước cutover):
--   1) Quy ước mã đơn THẬT của Nam Việt (tiền tố 'DH' vs 'HD', độ rộng zero-pad,
--      có gắn ngày không) — đang đề xuất 'DH' + zero-pad(8). Cấu hình 1 chỗ ở
--      internal/orders (app code generator). KẾ TOÁN/BA xác nhận.
--   2) Khi go-live phải SET sequence start > max số trong các mã hiện có cùng tiền
--      tố để KHÔNG đụng UNIQUE(code) với dữ liệu cũ:
--          SELECT setval('app.order_code_seq', <max_hiện_có> + 1, false);
--      (sequence app riêng + tiền tố riêng đã tránh đụng mã 'HD' lịch sử; bước này
--       chỉ cần khi tiền tố trùng tập mã cũ.)

-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

-- order_code_seq: cấp số tăng dần cho mã đơn (app đọc nextval rồi ghép tiền tố +
-- zero-pad). Sequence bảo đảm KHÔNG trùng, an toàn đua (mỗi nextval một giá trị).
CREATE SEQUENCE IF NOT EXISTS app.order_code_seq AS bigint START WITH 1 INCREMENT BY 1;

-- order_idempotency: chống TẠO ĐƠN trùng theo Idempotency-Key. Một key → một đơn.
-- Tạo đơn lần 2 cùng key → trả đơn cũ (đọc order_id/order_code). idem_key là khoá
-- idempotency tầng app do client gửi (header Idempotency-Key). order_id/order_code
-- trỏ về public.orders đã tạo (uuid/text — KHÔNG FK chéo schema, chỉ giữ ID).
CREATE TABLE app.order_idempotency (
    idem_key   text PRIMARY KEY CHECK (length(trim(idem_key)) > 0),
    order_id   uuid NOT NULL,
    order_code text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS app.order_idempotency;
DROP SEQUENCE IF EXISTS app.order_code_seq;
