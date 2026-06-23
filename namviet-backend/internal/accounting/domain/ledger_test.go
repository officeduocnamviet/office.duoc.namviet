package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// Bút toán cân điển hình TT133: bán hàng có VAT (Dr 131 / Cr 511 + Cr 3331).
func TestJournalEntry_Validate_Balanced(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookInternal, Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(110000)},
		{AccountCode: "511", Credit: money.FromInt(100000)},
		{AccountCode: "3331", Credit: money.FromInt(10000)},
	}}
	if err := e.Validate(); err != nil {
		t.Fatalf("entry cân phải hợp lệ, lỗi: %v", err)
	}
	if !e.TotalDebit().Equal(money.FromInt(110000)) {
		t.Fatalf("TotalDebit = %s, want 110000", e.TotalDebit())
	}
}

func TestJournalEntry_Validate_Unbalanced(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookInternal, Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(110000)},
		{AccountCode: "511", Credit: money.FromInt(100000)},
	}}
	if err := e.Validate(); err == nil {
		t.Fatal("entry lệch phải trả lỗi")
	}
}

func TestEntryLine_Validate_BothSides(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookTax, Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(100), Credit: money.FromInt(100)},
		{AccountCode: "511", Credit: money.FromInt(100)},
	}}
	if err := e.Validate(); err == nil {
		t.Fatal("line có cả 2 vế phải trả lỗi")
	}
}

func TestEntryLine_Validate_NoSide(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookInternal, Lines: []domain.EntryLine{
		{AccountCode: "131"}, // debit=0 và credit=0 → không vế nào
		{AccountCode: "511", Credit: money.FromInt(0)},
	}}
	if err := e.Validate(); err == nil {
		t.Fatal("line không có vế nào > 0 phải trả lỗi")
	}
}

func TestEntryLine_Validate_NegativeSide(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookInternal, Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(-100)},
		{AccountCode: "511", Credit: money.FromInt(-100)},
	}}
	if err := e.Validate(); err == nil {
		t.Fatal("vế âm phải trả lỗi (CHECK >= 0 ở DB; domain chặn trước)")
	}
}

func TestEntryLine_Validate_EmptyAccount(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookInternal, Lines: []domain.EntryLine{
		{AccountCode: "", Debit: money.FromInt(100)},
		{AccountCode: "511", Credit: money.FromInt(100)},
	}}
	if err := e.Validate(); err == nil {
		t.Fatal("account_code rỗng phải trả lỗi")
	}
}

func TestJournalEntry_Validate_BadBook(t *testing.T) {
	e := domain.JournalEntry{Book: "FOO", Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(100)},
		{AccountCode: "511", Credit: money.FromInt(100)},
	}}
	if err := e.Validate(); err == nil {
		t.Fatal("book sai phải trả lỗi")
	}
}

func TestJournalEntry_Validate_NoLines(t *testing.T) {
	e := domain.JournalEntry{Book: domain.BookInternal}
	if err := e.Validate(); err == nil {
		t.Fatal("entry không có dòng phải trả lỗi")
	}
}

// 2 sổ độc lập: cùng nghiệp vụ, INTERNAL theo giá thực, TAX theo giá HĐ — số tiền
// KHÁC nhau, mỗi entry tự cân, KHÔNG sync.
func TestJournalEntry_TwoBooks_Independent(t *testing.T) {
	internal := domain.JournalEntry{Book: domain.BookInternal, Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(110000)},
		{AccountCode: "511", Credit: money.FromInt(100000)},
		{AccountCode: "3331", Credit: money.FromInt(10000)},
	}}
	tax := domain.JournalEntry{Book: domain.BookTax, Lines: []domain.EntryLine{
		{AccountCode: "131", Debit: money.FromInt(99000)}, // giá HĐ khác giá thực
		{AccountCode: "511", Credit: money.FromInt(90000)},
		{AccountCode: "3331", Credit: money.FromInt(9000)},
	}}
	if err := internal.Validate(); err != nil {
		t.Fatalf("INTERNAL phải hợp lệ: %v", err)
	}
	if err := tax.Validate(); err != nil {
		t.Fatalf("TAX phải hợp lệ: %v", err)
	}
	if internal.TotalDebit().Equal(tax.TotalDebit()) {
		t.Fatal("2 sổ phải độc lập số tiền (không sync)")
	}
}
