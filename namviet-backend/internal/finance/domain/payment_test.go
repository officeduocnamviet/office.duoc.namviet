package domain_test

import (
	"errors"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
)

func m(t *testing.T, s string) money.Money {
	t.Helper()
	v, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return v
}

func TestBookType_Valid(t *testing.T) {
	for _, bt := range []domain.BookType{domain.BookInternal, domain.BookTax, domain.BookBoth} {
		if !bt.Valid() {
			t.Errorf("%q phải hợp lệ", bt)
		}
	}
	for _, bt := range []domain.BookType{"", "internal", "X", "both"} {
		if domain.BookType(bt).Valid() {
			t.Errorf("%q phải KHÔNG hợp lệ (case-sensitive, chỉ INTERNAL/TAX/BOTH)", bt)
		}
	}
}

func TestRecordPaymentIn_Validate_Happy(t *testing.T) {
	p := domain.RecordPaymentIn{
		OrderCode:     "HD123",
		Amount:        m(t, "150000"),
		FundAccountID: 1,
		BookType:      domain.BookBoth,
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("phiếu hợp lệ phải PASS: %v", err)
	}
}

func TestRecordPaymentIn_Validate_AmountNotPositive(t *testing.T) {
	for _, amt := range []string{"0", "-1", "-150000.50"} {
		p := domain.RecordPaymentIn{OrderCode: "HD1", Amount: m(t, amt), FundAccountID: 1, BookType: domain.BookBoth}
		err := p.Validate()
		if !errors.Is(err, domain.ErrAmountNotPositive) {
			t.Errorf("amount %q phải lỗi ErrAmountNotPositive, got %v", amt, err)
		}
	}
}

func TestRecordPaymentIn_Validate_EmptyOrderCode(t *testing.T) {
	for _, code := range []string{"", "   ", "\t"} {
		p := domain.RecordPaymentIn{OrderCode: code, Amount: m(t, "1"), FundAccountID: 1, BookType: domain.BookBoth}
		if err := p.Validate(); !errors.Is(err, domain.ErrEmptyOrderCode) {
			t.Errorf("order_code %q phải lỗi ErrEmptyOrderCode, got %v", code, err)
		}
	}
}

func TestRecordPaymentIn_Validate_InvalidBookType(t *testing.T) {
	p := domain.RecordPaymentIn{OrderCode: "HD1", Amount: m(t, "1"), FundAccountID: 1, BookType: "BANK"}
	if err := p.Validate(); !errors.Is(err, domain.ErrInvalidBookType) {
		t.Errorf("book_type sai phải lỗi ErrInvalidBookType, got %v", err)
	}
}

// Thứ tự gate: order_code rỗng kiểm TRƯỚC amount (báo lỗi cụ thể đầu tiên gặp).
func TestRecordPaymentIn_Validate_OrderCodeBeforeAmount(t *testing.T) {
	p := domain.RecordPaymentIn{OrderCode: "", Amount: m(t, "0"), FundAccountID: 1, BookType: domain.BookBoth}
	if err := p.Validate(); !errors.Is(err, domain.ErrEmptyOrderCode) {
		t.Errorf("order_code rỗng phải báo trước, got %v", err)
	}
}
