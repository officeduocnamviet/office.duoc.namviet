-- Phân bổ 1 phiếu THU → NHIỀU đơn (spec system_features.md mục 55): khách trả 1
-- cục (vd 10tr) → phân bổ cho nhiều đơn (A1 2tr, A2 3tr, A3 5tr), đơn cũ nhất trả
-- trước. Object MỚI do backend sở hữu (schema app, ADR 0002).
--
-- payment_id trỏ public.finance_transactions.id (bigint) — KHÔNG đặt FK cross-schema
-- (toàn vẹn app-enforce; nhất quán cách app ledger đọc public.chart_of_accounts mà
-- không FK). Tiền NUMERIC(20,0) (VND scale-0). 1 phiếu phân bổ cho 1 đơn TỐI ĐA một
-- dòng (UNIQUE payment_id, order_code).
--
-- "Đã thu mỗi đơn" = phiếu trực tiếp (finance_transactions.ref_id=order.code) +
-- phân bổ (allocations.order_code=order.code). KHÔNG đếm trùng: phiếu lump-sum đặt
-- ref_type='customer' (KHÔNG 'order') nên KHÔNG bị query trực tiếp đếm; chỉ allocation
-- của nó được đếm.

-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE app.finance_transaction_allocations (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    -- payment_id = public.finance_transactions.id (phiếu thu nguồn). KHÔNG FK cross-schema.
    payment_id  bigint NOT NULL,
    -- order_code = public.orders.code (text) — khớp ref_id quy ước đơn↔phiếu.
    order_code  text   NOT NULL,
    amount      numeric(20,0) NOT NULL CHECK (amount > 0),
    created_at  timestamptz NOT NULL DEFAULT now(),
    -- 1 phiếu phân bổ cho 1 đơn tối đa MỘT dòng (cộng dồn nếu cần thì update dòng).
    CONSTRAINT uq_fta_payment_order UNIQUE (payment_id, order_code)
);
-- Tính "đã thu mỗi đơn": SUM theo order_code (index nóng).
CREATE INDEX idx_fta_order ON app.finance_transaction_allocations (order_code);
-- Đối soát 1 phiếu phân bổ ra những đơn nào (icon expand ở UI mục 55).
CREATE INDEX idx_fta_payment ON app.finance_transaction_allocations (payment_id);

-- order_paid_amount: NGUỒN CHÂN LÝ DUY NHẤT tính "đã thu" của một đơn = phiếu THU
-- trực tiếp (finance_transactions.ref_id=order.code) + phân bổ (allocations.order_code).
-- KHÔNG đếm trùng (phiếu lump-sum ref_type='customer' không bị nhánh trực tiếp đếm).
-- Lọc nhất quán: flow='in', status IN (pending,completed) [pending=NV đã thu, mục 55],
-- chưa xoá, book_type sổ thực tế (INTERNAL/BOTH). LANGUAGE plpgsql để defer resolve
-- public.* (test áp public schema SAU migration; plpgsql resolve lúc gọi, không lúc tạo).
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.order_paid_amount(p_order_code text) RETURNS numeric AS $$
DECLARE
    direct numeric;
    alloc  numeric;
BEGIN
    SELECT COALESCE(SUM(ft.amount), 0) INTO direct
    FROM public.finance_transactions ft
    WHERE ft.ref_type = 'order'
      AND ft.ref_id = p_order_code
      AND lower(ft.flow) = 'in'
      AND ft.status IN ('pending', 'completed')
      AND ft.deleted_at IS NULL
      AND ft.book_type IN ('INTERNAL', 'BOTH');

    SELECT COALESCE(SUM(a.amount), 0) INTO alloc
    FROM app.finance_transaction_allocations a
    JOIN public.finance_transactions ft2 ON ft2.id = a.payment_id
    WHERE a.order_code = p_order_code
      AND lower(ft2.flow) = 'in'
      AND ft2.status IN ('pending', 'completed')
      AND ft2.deleted_at IS NULL
      AND ft2.book_type IN ('INTERNAL', 'BOTH');

    RETURN direct + alloc;
END;
$$ LANGUAGE plpgsql STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS app.order_paid_amount(text);
DROP TABLE IF EXISTS app.finance_transaction_allocations;
