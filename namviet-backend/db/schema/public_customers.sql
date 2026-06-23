-- SCHEMA THAM CHIẾU (strangler-fig, ADR 0001) — KHÔNG phải goose migration.
--
-- Mô tả các bảng public.* ĐANG TỒN TẠI (Postgres Supabase, nhánh monorepo
-- office.duoc.namviet) liên quan bounded context customers, để (1) sqlc
-- type-check query customers ở compile-time và (2) dbtest materialize bảng + seed
-- trong testcontainers (DB test trống). TUYỆT ĐỐI KHÔNG chạy qua goose lên prod —
-- backend chỉ ĐỌC các bảng kế thừa này, không sở hữu chúng.
--
-- Cột lấy nguyên văn từ:
--   office.duoc.namviet/database_schema.md (doc) và
--   office.duoc.namviet/namviet-infrastructure/supabase/migrations/20260613080000_init_erp_schema.sql.
-- Đây là nguồn DOC/MIGRATION-CỦA-NHÁNH-NÀY, PHẢI verify lại với prod thật
-- (pg_dump/REST Object.keys) khi có creds — prod hay lệch migration.
--
-- GHI CHÚ NGHIỆP VỤ QUAN TRỌNG (đọc kỹ trước khi đụng debt):
--   * Thông tin B2B (tax_code/MST, debt_limit, payment_term, sales_staff) nằm
--     TRONG customers.b2b_metadata jsonb — KHÔNG có bảng customers_b2b ở DB này
--     (bảng customers_b2b + view actual_current_debt thuộc DB ERP nhánh KHÁC, xem
--     ADR 0002: hai DB là hai nhánh). DÙNG ĐÚNG cái DB này thể hiện: jsonb.
--   * Công nợ: customers.current_debt là cột TĨNH (stale). DB này KHÔNG có view
--     live actual_current_debt. Nguồn LIVE ở đây = tổng final_amount của các đơn
--     chưa tất toán (payment_status <> 'paid') tính từ public.orders.
--   * CAVEAT GHOST-DEBT: public.orders KHÔNG có cột paid_amount (xác nhận từ
--     init migration + doc; ADR 0002 §B coi đây là khoảng trống). Vì vậy debt
--     live KHÔNG trừ được phần đã trả của đơn 'partial' → có thể PHÌNH nợ. KHÔNG
--     tự "sửa" số; chỉ trình bày và ghi chú rõ.
--   * THAM CHIẾU (KHÔNG copy được): DB ERP nhánh khác định nghĩa công nợ live
--     bằng view b2b_customer_debt_view = SUM(orders.final_amount) cho đơn
--     status IN (PACKED,SHIPPING,DELIVERED,COMPLETED) TRỪ tổng phiếu thu
--     finance_transactions(flow=in,completed). DB monorepo NÀY chưa có view đó,
--     chưa có finance_transactions linkage chuẩn, chưa có paid_amount → ta xấp xỉ
--     bằng payment_status<>'paid'. Khi cutover/migrate cần reconcile theo định
--     nghĩa ERP (bổ sung paid_amount + view) — đây là việc của ADR 0002 §B.
--   * id customers = bigint (KHÔNG uuid). Tiền = numeric (→ common/money).

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.customers (
    id             bigint PRIMARY KEY,
    customer_code  text,
    name           text NOT NULL,
    customer_type  text NOT NULL DEFAULT 'B2C', -- 'B2C' | 'B2B'
    phone          text,
    email          text,
    address        text,
    status         text NOT NULL DEFAULT 'active', -- 'active' | 'inactive' | 'banned'
    dob            date,
    gender         text,
    cccd           text,
    loyalty_points integer DEFAULT 0,
    -- b2b_metadata chứa: tax_code, debt_limit, payment_term, sales_staff_id... (jsonb).
    b2b_metadata   jsonb DEFAULT '{}'::jsonb,
    current_debt   numeric DEFAULT 0, -- TĨNH (stale) — không dùng làm nguồn live.
    updated_by     uuid,
    created_at     timestamptz DEFAULT now(),
    updated_at     timestamptz DEFAULT now(),
    deleted_at     timestamptz
);

-- orders: chỉ các cột cần để TÍNH công nợ LIVE (tổng final_amount đơn chưa tất
-- toán). KHÔNG có paid_amount ở DB này (ghost-debt caveat). id = uuid.
CREATE TABLE IF NOT EXISTS public.orders (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code           text NOT NULL,
    customer_id    bigint,
    creator_id     uuid,
    status         text NOT NULL DEFAULT 'PENDING',
    order_type     text NOT NULL DEFAULT 'B2C',
    total_amount   numeric DEFAULT 0,
    final_amount   numeric DEFAULT 0, -- tổng tiền khách phải trả
    payment_status text DEFAULT 'unpaid', -- 'unpaid' | 'partial' | 'paid'
    note           text,
    created_at     timestamptz DEFAULT now(),
    updated_at     timestamptz DEFAULT now(),
    deleted_at     timestamptz
);
