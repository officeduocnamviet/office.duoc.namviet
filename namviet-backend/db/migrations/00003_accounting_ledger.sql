-- Sổ kế toán kép NỀN (schema app, ADR 0002). Đây là object MỚI do backend sở
-- hữu (KHÔNG phải bảng public.* kế thừa). Double-entry: mỗi dòng đúng MỘT vế
-- (debit XOR credit) > 0; Σdebit = Σcredit mỗi bút toán (ép ở SERVICE + constraint
-- trigger DEFERRABLE phòng thủ ở DB). Append-only: KHÔNG UPDATE/DELETE dòng sổ ở
-- tầng app. 2 sổ INTERNAL/TAX tách biệt (mỗi entry mang đúng 1 book). Gating kỳ:
-- chỉ post vào kỳ đang 'open'. Tiền = NUMERIC(20,0) (VND, scale-0).
--
-- PK uuid: DEFAULT gen_random_uuid() chỉ là LƯỚI AN TOÀN; app sinh uuid v7
-- (common/id) và truyền vào INSERT (time-ordered, tốt cho index theo thời gian,
-- không phụ thuộc pgcrypto trên prod).

-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE app.accounting_periods (
    id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    year     int  NOT NULL,
    month    int  NOT NULL CHECK (month BETWEEN 1 AND 12),
    status   text NOT NULL DEFAULT 'open' CHECK (status IN ('open','closed')),
    UNIQUE (year, month)
);

CREATE TABLE app.journal_entries (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    book        text NOT NULL CHECK (book IN ('INTERNAL','TAX')),
    entry_date  date NOT NULL,
    period_id   uuid NOT NULL REFERENCES app.accounting_periods(id),
    source_type text NOT NULL DEFAULT '',
    source_id   text NOT NULL DEFAULT '',
    memo        text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_journal_entries_book_period ON app.journal_entries (book, period_id);
CREATE INDEX idx_journal_entries_source ON app.journal_entries (source_type, source_id);
-- Keyset đọc /v1/accounting/entries theo (created_at DESC, id DESC).
CREATE INDEX idx_journal_entries_created ON app.journal_entries (created_at DESC, id DESC);

CREATE TABLE app.journal_entry_lines (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id     uuid NOT NULL REFERENCES app.journal_entries(id) ON DELETE CASCADE,
    line_no      int  NOT NULL,
    account_code text NOT NULL,
    debit        numeric(20,0) NOT NULL DEFAULT 0 CHECK (debit  >= 0),
    credit       numeric(20,0) NOT NULL DEFAULT 0 CHECK (credit >= 0),
    -- Mỗi dòng đúng MỘT vế > 0: (debit=0) XOR (credit=0).
    CHECK ((debit = 0) <> (credit = 0))
);
CREATE INDEX idx_journal_entry_lines_entry ON app.journal_entry_lines (entry_id);
CREATE INDEX idx_journal_entry_lines_account ON app.journal_entry_lines (account_code);

-- assert_entry_balanced: trigger phòng thủ kiểm Σdebit = Σcredit của một entry
-- (và entry không rỗng). Là CONSTRAINT TRIGGER DEFERRABLE INITIALLY DEFERRED nên
-- chạy lúc COMMIT — cho phép insert nhiều dòng trong cùng tx rồi mới cân tổng.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.assert_entry_balanced() RETURNS trigger AS $$
DECLARE
    d   numeric(20,0);
    c   numeric(20,0);
    eid uuid := COALESCE(NEW.entry_id, OLD.entry_id);
BEGIN
    -- Entry có thể đã bị xoá (cascade) trước khi trigger deferred chạy → bỏ qua.
    IF NOT EXISTS (SELECT 1 FROM app.journal_entries WHERE id = eid) THEN
        RETURN NULL;
    END IF;
    SELECT COALESCE(SUM(debit),0), COALESCE(SUM(credit),0)
      INTO d, c FROM app.journal_entry_lines WHERE entry_id = eid;
    IF d <> c THEN
        RAISE EXCEPTION 'but toan % lech: Sigma debit=% <> Sigma credit=%', eid, d, c
            USING ERRCODE = 'check_violation';
    END IF;
    IF d = 0 THEN
        RAISE EXCEPTION 'but toan % rong (khong co dong)', eid
            USING ERRCODE = 'check_violation';
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE CONSTRAINT TRIGGER trg_journal_entry_balanced
    AFTER INSERT OR UPDATE OR DELETE ON app.journal_entry_lines
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE FUNCTION app.assert_entry_balanced();

-- +goose Down
DROP TRIGGER IF EXISTS trg_journal_entry_balanced ON app.journal_entry_lines;
DROP FUNCTION IF EXISTS app.assert_entry_balanced();
DROP TABLE IF EXISTS app.journal_entry_lines;
DROP TABLE IF EXISTS app.journal_entries;
DROP TABLE IF EXISTS app.accounting_periods;
