package domain_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

func TestQuantity_FromStringAndString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", "0"},
		{"0", "0"},
		{"12", "12"},
		{"12.5", "12.5"},
		{"12.50", "12.5"}, // chuẩn hoá số 0 thừa
	}
	for _, c := range cases {
		q, err := domain.QuantityFromString(c.in)
		if err != nil {
			t.Fatalf("QuantityFromString(%q): %v", c.in, err)
		}
		if q.String() != c.want {
			t.Fatalf("QuantityFromString(%q).String() = %q, want %q", c.in, q.String(), c.want)
		}
	}
}

func TestQuantity_ParseError(t *testing.T) {
	if _, err := domain.QuantityFromString("abc"); err == nil {
		t.Fatal("parse 'abc' phải lỗi, không nuốt thành 0")
	}
}

func TestQuantity_NoFloatPrecisionLoss(t *testing.T) {
	// Giá trị mà float64 không biểu diễn chính xác — phải giữ nguyên qua decimal.
	q, err := domain.QuantityFromString("100000000000.01")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if q.String() != "100000000000.01" {
		t.Fatalf("mất chính xác: %q", q.String())
	}
}

func TestQuantity_Predicates(t *testing.T) {
	if !domain.ZeroQty().IsZero() {
		t.Fatal("ZeroQty phải IsZero")
	}
	if domain.ZeroQty().IsPositive() {
		t.Fatal("0 không IsPositive")
	}
	pos := domain.QuantityFromDecimal(decimal.RequireFromString("0.5"))
	if !pos.IsPositive() || pos.IsZero() {
		t.Fatal("0.5 phải IsPositive, không IsZero")
	}
	a := domain.QuantityFromDecimal(decimal.RequireFromString("12.50"))
	b := domain.QuantityFromDecimal(decimal.RequireFromString("12.5"))
	if !a.Equal(b) {
		t.Fatal("12.50 phải Equal 12.5")
	}
}
