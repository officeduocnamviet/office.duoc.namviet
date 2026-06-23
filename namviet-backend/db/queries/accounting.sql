-- Accounting (P0 — sổ kế toán nền). GHI object MỚI ở schema app
-- (journal_entries/_lines/accounting_periods) + ĐỌC public.chart_of_accounts
-- (strangler, ADR 0001). Tiền debit/credit là NUMERIC(20,0) → common/money ở repo
-- (KHÔNG float). PK uuid sinh app-side (common/id v7) và truyền vào INSERT — ổn
-- định, time-ordered, không phụ thuộc pgcrypto trên prod. Post bút toán chạy
-- TRONG tx do caller truyền (gộp atomic với nghiệp vụ orders/finance): service
-- insert entry rồi insert từng line; trigger DEFERRABLE cân Σ lúc commit.

-- name: GetChartAccount :one
-- Đọc một tài khoản trong cây TK theo account_code (chưa soft-delete). Trả
-- allow_posting/balance_type/type để service gate hạch toán (chỉ post vào TK
-- allow_posting = true). Không thấy → sqlc trả pgx.ErrNoRows (repo map NotFound).
SELECT account_code, allow_posting, balance_type, type
FROM public.chart_of_accounts
WHERE account_code = @account_code::text
  AND deleted_at IS NULL;

-- name: GetOpenPeriodByDate :one
-- Tìm kỳ kế toán chứa entry_date (theo year/month của ngày) VÀ đang 'open'. Dùng
-- gate: chỉ post vào kỳ mở. Không có kỳ (chưa mở) hoặc kỳ đã 'closed' → ErrNoRows
-- (service phân biệt "không có kỳ mở" = Conflict kỳ khoá).
SELECT id, year, month, status
FROM app.accounting_periods
WHERE year = @year::int
  AND month = @month::int
  AND status = 'open';

-- name: CreatePeriod :one
-- Mở một kỳ kế toán mới (status 'open'). UNIQUE(year,month) chống trùng. id sinh
-- app-side. Dùng cho seed/khởi tạo kỳ (P0 chưa expose REST quản lý kỳ).
INSERT INTO app.accounting_periods (id, year, month, status)
VALUES (@id::uuid, @year::int, @month::int, COALESCE(sqlc.narg('status')::text, 'open'))
RETURNING id, year, month, status;

-- name: InsertJournalEntry :one
-- Ghi một bút toán (header). book ∈ {INTERNAL,TAX}; period_id là kỳ đang mở (đã
-- gate ở service). id sinh app-side. RETURNING id để service insert lines tiếp.
INSERT INTO app.journal_entries (id, book, entry_date, period_id, source_type, source_id, memo)
VALUES (@id::uuid, @book::text, @entry_date::date, @period_id::uuid,
        @source_type::text, @source_id::text, @memo::text)
RETURNING id;

-- name: InsertJournalEntryLine :exec
-- Ghi một dòng bút toán. Mỗi dòng đúng MỘT vế > 0 (CHECK ở DB); Σdebit=Σcredit
-- kiểm bởi constraint trigger DEFERRABLE lúc commit. id sinh app-side.
INSERT INTO app.journal_entry_lines (id, entry_id, line_no, account_code, debit, credit)
VALUES (@id::uuid, @entry_id::uuid, @line_no::int, @account_code::text,
        @debit::numeric, @credit::numeric);

-- name: GetEntry :one
-- Một bút toán theo id (cho HTTP đọc /v1/accounting/entries/{id}).
SELECT id, book, entry_date, period_id, source_type, source_id, memo, created_at
FROM app.journal_entries
WHERE id = @id::uuid;

-- name: ListEntryLines :many
-- Các dòng của một bút toán, thứ tự theo line_no (ổn định).
SELECT id, entry_id, line_no, account_code, debit, credit
FROM app.journal_entry_lines
WHERE entry_id = @entry_id::uuid
ORDER BY line_no ASC, id ASC;

-- name: ListEntries :many
-- Danh sách bút toán, keyset theo (created_at DESC, id DESC): trang kế lấy entry
-- "cũ hơn" mốc cursor. @after_created_at NULL = trang đầu (từ mới nhất). Lọc
-- optional book và period_id.
SELECT id, book, entry_date, period_id, source_type, source_id, memo, created_at
FROM app.journal_entries
WHERE (
        sqlc.narg('after_created_at')::timestamptz IS NULL
        OR created_at < sqlc.narg('after_created_at')::timestamptz
        OR (created_at = sqlc.narg('after_created_at')::timestamptz
            AND id < sqlc.narg('after_id')::uuid)
      )
  AND (sqlc.narg('book')::text IS NULL OR book = sqlc.narg('book')::text)
  AND (sqlc.narg('period_id')::uuid IS NULL OR period_id = sqlc.narg('period_id')::uuid)
ORDER BY created_at DESC, id DESC
LIMIT @row_limit::int;
