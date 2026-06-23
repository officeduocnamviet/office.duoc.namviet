-- SCHEMA THAM CHIẾU (strangler-fig, ADR 0001) — KHÔNG phải goose migration.
--
-- Mô tả bảng public.chart_of_accounts ĐANG TỒN TẠI (Postgres Supabase, monorepo
-- office.duoc.namviet) — CÂY TÀI KHOẢN theo Thông tư 133. Bounded context
-- ACCOUNTING chỉ ĐỌC bảng này (account_code/type/balance_type/allow_posting) để
-- gate hạch toán: chỉ post vào tài khoản có allow_posting = true. Mục đích file:
-- (1) sqlc type-check query accounting ở compile-time; (2) dbtest materialize
-- bảng + seed trong testcontainers (DB test trống). TUYỆT ĐỐI KHÔNG chạy lên prod
-- — backend chỉ ĐỌC bảng kế thừa này, KHÔNG sở hữu.
--
-- Cột lấy nguyên văn từ office.duoc.namviet/database_schema.md (mục
-- chart_of_accounts). LƯU Ý tên cột THẬT: là `name` (KHÔNG `account_name`) và
-- `type` (KHÔNG `account_type`). PHẢI verify lại với prod thật (pg_dump/REST
-- Object.keys) khi có creds.
--
-- Object MỚI ở schema app (journal_entries/_lines/accounting_periods) do goose
-- migration 00003 tạo — sqlc đã thấy qua db/migrations; KHÔNG khai lại ở đây.

CREATE SCHEMA IF NOT EXISTS public;

-- chart_of_accounts: hệ thống tài khoản kế toán (cây TK TT133). PK uuid. Chỉ
-- tài khoản có allow_posting = true mới được hạch toán trực tiếp (TK tổng hợp
-- cấp cha allow_posting = false). balance_type ∈ {DEBIT, CREDIT, BOTH} (số dư
-- thường ở bên nào). status text tự do ('active'...) — soft-delete qua deleted_at.
CREATE TABLE IF NOT EXISTS public.chart_of_accounts (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code  text NOT NULL,                    -- số hiệu TK (vd 111, 131, 511, 632, 3331)
    name          text NOT NULL,                    -- tên tài khoản
    parent_id     uuid,                             -- TK cấp cha (cây)
    type          text NOT NULL,                    -- loại TK (Tài sản, Nợ, Vốn...)
    balance_type  text NOT NULL,                    -- DEBIT | CREDIT | BOTH
    status        text NOT NULL DEFAULT 'active',
    allow_posting boolean NOT NULL DEFAULT true,    -- cho phép hạch toán trực tiếp
    created_at    timestamptz DEFAULT now(),
    updated_at    timestamptz DEFAULT now(),
    deleted_at    timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_chart_of_accounts_code
    ON public.chart_of_accounts (account_code) WHERE deleted_at IS NULL;
