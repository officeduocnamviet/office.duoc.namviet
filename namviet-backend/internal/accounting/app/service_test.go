package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/app"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// fakeStore là EntryStore giả lập (in-memory) để test logic gate ở service KHÔNG
// cần DB. periodOpen điều khiển có kỳ mở không; postable là tập account_code được
// phép hạch toán; insertedLines/insertedEntries đếm để assert.
type fakeStore struct {
	periodOpen     bool
	postable       map[string]bool
	insertedEntry  int
	insertedLines  int
	failInsertLine bool // mô phỏng trigger cân Σ FAIL lúc commit
}

func (f *fakeStore) OpenPeriodOf(_ context.Context, _ string) (string, bool, error) {
	if !f.periodOpen {
		return "", false, nil
	}
	return "period-2026-06", true, nil
}

func (f *fakeStore) AccountPostable(_ context.Context, code string) (bool, error) {
	return f.postable[code], nil
}

func (f *fakeStore) InsertEntry(_ context.Context, _ domain.JournalEntry, _ string) (string, error) {
	f.insertedEntry++
	return "entry-1", nil
}

func (f *fakeStore) InsertLine(_ context.Context, _ string, _ int32, _ domain.EntryLine) error {
	if f.failInsertLine {
		return errors.New("trigger check_violation: but toan lech")
	}
	f.insertedLines++
	return nil
}

// balanced trả một bút toán cân TT133 (Dr 131 / Cr 511 + 3331).
func balanced(book domain.Book) domain.JournalEntry {
	return domain.JournalEntry{
		Book:      book,
		EntryDate: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(110000)},
			{AccountCode: "511", Credit: money.FromInt(100000)},
			{AccountCode: "3331", Credit: money.FromInt(10000)},
		},
	}
}

func allPostable() map[string]bool {
	return map[string]bool{"131": true, "511": true, "3331": true, "632": true, "156": true}
}

// newSvc dựng service với một fakeStore cho mọi tx (bỏ qua tx trong test).
func newSvc(store *fakeStore) *app.Service {
	return app.New(
		func(_ pgx.Tx) app.EntryStore { return store },
		nil, // TxManager không dùng trong các test Post(ctx, tx, ...)
		nil, // ReadStore không dùng ở đây
	)
}

func TestPost_Success(t *testing.T) {
	store := &fakeStore{periodOpen: true, postable: allPostable()}
	svc := newSvc(store)
	id, err := svc.Post(context.Background(), nil, balanced(domain.BookInternal))
	if err != nil {
		t.Fatalf("post cân + kỳ mở + TK hợp lệ phải thành công: %v", err)
	}
	if id != "entry-1" {
		t.Fatalf("entryID = %q, want entry-1", id)
	}
	if store.insertedEntry != 1 || store.insertedLines != 3 {
		t.Fatalf("inserted entry=%d lines=%d, want 1/3", store.insertedEntry, store.insertedLines)
	}
}

func TestPost_PeriodClosed_Conflict(t *testing.T) {
	store := &fakeStore{periodOpen: false, postable: allPostable()}
	svc := newSvc(store)
	_, err := svc.Post(context.Background(), nil, balanced(domain.BookInternal))
	if err == nil {
		t.Fatal("post vào kỳ khoá phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("kỳ khoá phải là Conflict, got %v", apperr.KindOf(err))
	}
	if store.insertedEntry != 0 {
		t.Fatal("kỳ khoá KHÔNG được insert entry")
	}
}

func TestPost_AccountNotPostable_Unprocessable(t *testing.T) {
	store := &fakeStore{periodOpen: true, postable: map[string]bool{"131": true, "511": true}} // thiếu 3331
	svc := newSvc(store)
	_, err := svc.Post(context.Background(), nil, balanced(domain.BookInternal))
	if err == nil {
		t.Fatal("post vào TK không cho hạch toán phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("TK không postable phải Validation/Unprocessable, got %v", apperr.KindOf(err))
	}
}

func TestPost_Unbalanced_Unprocessable(t *testing.T) {
	store := &fakeStore{periodOpen: true, postable: allPostable()}
	svc := newSvc(store)
	bad := balanced(domain.BookInternal)
	bad.Lines = bad.Lines[:2] // bỏ dòng 3331 → lệch
	_, err := svc.Post(context.Background(), nil, bad)
	if err == nil {
		t.Fatal("post bút toán lệch phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("lệch Σ phải Validation/Unprocessable, got %v", apperr.KindOf(err))
	}
	if store.insertedEntry != 0 {
		t.Fatal("lệch Σ KHÔNG được insert (chặn ở Validate trước DB)")
	}
}

func TestPost_BadBook_Unprocessable(t *testing.T) {
	store := &fakeStore{periodOpen: true, postable: allPostable()}
	svc := newSvc(store)
	bad := balanced("FOO")
	_, err := svc.Post(context.Background(), nil, bad)
	if err == nil {
		t.Fatal("book sai phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("book sai phải Validation, got %v", apperr.KindOf(err))
	}
}

// 2 sổ độc lập: post INTERNAL rồi TAX cùng nghiệp vụ → cả 2 thành công, đếm 2
// entry (không sync, mỗi entry riêng).
func TestPost_TwoBooks_Independent(t *testing.T) {
	store := &fakeStore{periodOpen: true, postable: allPostable()}
	svc := newSvc(store)
	if _, err := svc.Post(context.Background(), nil, balanced(domain.BookInternal)); err != nil {
		t.Fatalf("INTERNAL: %v", err)
	}
	if _, err := svc.Post(context.Background(), nil, balanced(domain.BookTax)); err != nil {
		t.Fatalf("TAX: %v", err)
	}
	if store.insertedEntry != 2 {
		t.Fatalf("2 sổ phải sinh 2 entry độc lập, got %d", store.insertedEntry)
	}
}
