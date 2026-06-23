// Package postgres là ADAPTER ra phía cơ sở dữ liệu của accounting: implement các
// port app (EntryStore ghi/tx, ReadStore đọc) bằng query sinh từ sqlc (appdb) và
// map row <-> entity domain. GHI object MỚI ở schema app (journal_entries/_lines/
// accounting_periods) + ĐỌC public.chart_of_accounts (strangler, ADR 0001). Nằm
// dưới internal/ nên module khác KHÔNG import được. Tiền debit/credit = NUMERIC
// (20,0) <-> common/money decimal, KHÔNG float (numericToMoney/moneyToNumeric
// dựng từ big.Int, không qua float). PK uuid sinh app-side (common/id v7).
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/app"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// EntryRepo implement app.EntryStore trên appdb.Queries (bind tx do caller/TxManager
// truyền). Mọi thao tác GHI sổ chạy trong tx này (gộp atomic với nghiệp vụ).
type EntryRepo struct{ q *appdb.Queries }

// NewEntryRepo tạo repo ghi từ một *appdb.Queries (đã bind tx).
func NewEntryRepo(q *appdb.Queries) *EntryRepo { return &EntryRepo{q: q} }

// OpenPeriodOf tìm kỳ kế toán đang 'open' chứa ngày date (YYYY-MM-DD). Không có
// kỳ mở (chưa mở/đã khoá) → (_, false, nil). Lỗi parse date → (_, false, err).
func (r *EntryRepo) OpenPeriodOf(ctx context.Context, date string) (string, bool, error) {
	d, err := parseYMD(date)
	if err != nil {
		return "", false, err
	}
	row, err := r.q.GetOpenPeriodByDate(ctx, appdb.GetOpenPeriodByDateParams{
		Year:  int32(d.year),
		Month: int32(d.month),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil // không có kỳ mở
		}
		return "", false, err
	}
	return row.ID, true, nil
}

// AccountPostable trả true nếu account_code tồn tại (chưa soft-delete) và
// allow_posting = true. Không tồn tại → (false, nil).
func (r *EntryRepo) AccountPostable(ctx context.Context, accountCode string) (bool, error) {
	row, err := r.q.GetChartAccount(ctx, accountCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return row.AllowPosting, nil
}

// InsertEntry ghi header bút toán (id sinh app-side uuid v7) trả id.
func (r *EntryRepo) InsertEntry(ctx context.Context, e domain.JournalEntry, periodID string) (string, error) {
	entryID := id.NewString()
	got, err := r.q.InsertJournalEntry(ctx, appdb.InsertJournalEntryParams{
		ID:         entryID,
		Book:       string(e.Book),
		EntryDate:  dateToPg(e.EntryDate),
		PeriodID:   periodID,
		SourceType: e.SourceType,
		SourceID:   e.SourceID,
		Memo:       e.Memo,
	})
	if err != nil {
		return "", err
	}
	return got, nil
}

// InsertLine ghi một dòng bút toán (id sinh app-side). Trigger DEFERRABLE cân Σ
// kiểm lúc commit.
func (r *EntryRepo) InsertLine(ctx context.Context, entryID string, lineNo int32, l domain.EntryLine) error {
	return r.q.InsertJournalEntryLine(ctx, appdb.InsertJournalEntryLineParams{
		ID:          id.NewString(),
		EntryID:     entryID,
		LineNo:      lineNo,
		AccountCode: l.AccountCode,
		Debit:       moneyToNumeric(l.Debit),
		Credit:      moneyToNumeric(l.Credit),
	})
}

// ReadRepo implement app.ReadStore (bind pool). Đường đọc thuần, không tx.
type ReadRepo struct{ q *appdb.Queries }

// NewReadRepo tạo repo đọc từ một *appdb.Queries (đã bind pool).
func NewReadRepo(q *appdb.Queries) *ReadRepo { return &ReadRepo{q: q} }

func (r *ReadRepo) ListEntries(ctx context.Context, f domain.EntryFilter) ([]domain.JournalEntryRecord, error) {
	params := appdb.ListEntriesParams{
		RowLimit: f.Limit,
		Book:     strNarg(f.Book),
		PeriodID: strNarg(f.PeriodID),
	}
	if f.HasCursor {
		params.AfterCreatedAt = tsToPg(f.AfterCreatedAt)
		aid := f.AfterID
		params.AfterID = &aid
	}
	rows, err := r.q.ListEntries(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]domain.JournalEntryRecord, 0, len(rows))
	for _, row := range rows {
		out = append(out, entryRecord(row))
	}
	return out, nil
}

func (r *ReadRepo) GetEntry(ctx context.Context, entryID string) (domain.JournalEntryRecord, error) {
	row, err := r.q.GetEntry(ctx, entryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.JournalEntryRecord{}, apperr.NotFound("bút toán không tồn tại")
		}
		return domain.JournalEntryRecord{}, err
	}
	rec := entryRecord(row)
	lines, err := r.q.ListEntryLines(ctx, entryID)
	if err != nil {
		return domain.JournalEntryRecord{}, err
	}
	rec.Lines = make([]domain.EntryLine, 0, len(lines))
	for _, l := range lines {
		rec.Lines = append(rec.Lines, domain.EntryLine{
			AccountCode: l.AccountCode,
			Debit:       numericToMoney(l.Debit),
			Credit:      numericToMoney(l.Credit),
		})
	}
	return rec, nil
}

// ---- mapping row <-> domain ----

func entryRecord(row appdb.AppJournalEntry) domain.JournalEntryRecord {
	return domain.JournalEntryRecord{
		ID:         row.ID,
		Book:       domain.Book(row.Book),
		EntryDate:  row.EntryDate.Time,
		PeriodID:   row.PeriodID,
		SourceType: row.SourceType,
		SourceID:   row.SourceID,
		Memo:       row.Memo,
		CreatedAt:  row.CreatedAt.Time,
	}
}

// numericToMoney chuyển pgtype.Numeric sang money.Money KHÔNG đi qua float: dựng
// decimal trực tiếp từ mantissa (big.Int) * 10^Exp. NULL/NaN → Zero.
func numericToMoney(n pgtype.Numeric) money.Money {
	if !n.Valid || n.NaN || n.Int == nil {
		return money.Zero()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return money.FromDecimal(d)
}

// moneyToNumeric chuyển money.Money sang pgtype.Numeric KHÔNG qua float: lấy
// coefficient (big.Int) + exponent từ decimal. Dùng khi GHI debit/credit.
func moneyToNumeric(m money.Money) pgtype.Numeric {
	d := m.Decimal()
	return pgtype.Numeric{
		Int:   new(big.Int).Set(d.Coefficient()),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

type ymd struct {
	year  int
	month int
}

// parseYMD parse "YYYY-MM-DD" → year/month (để tra kỳ kế toán). Sai định dạng →
// lỗi (service map Internal — service luôn truyền date từ time.Time đã hợp lệ).
func parseYMD(s string) (ymd, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return ymd{}, fmt.Errorf("parse date %q: %w", s, err)
	}
	return ymd{year: t.Year(), month: int(t.Month())}, nil
}

// dateToPg chuyển time.Time → pgtype.Date (chỉ phần ngày).
func dateToPg(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

// tsToPg chuyển time.Time → pgtype.Timestamptz.
func tsToPg(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// strNarg chuyển chuỗi rỗng → nil (không lọc) cho sqlc narg.
func strNarg(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Đảm bảo các repo thoả port app ở compile-time.
var (
	_ app.EntryStore = (*EntryRepo)(nil)
	_ app.ReadStore  = (*ReadRepo)(nil)
)
