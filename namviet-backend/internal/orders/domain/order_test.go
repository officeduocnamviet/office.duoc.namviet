package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

func mustMoney(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money.FromString(%q): %v", s, err)
	}
	return m
}

func TestComputePayment_Remaining(t *testing.T) {
	cases := []struct {
		name          string
		final, paid   string
		wantRemaining string
		wantPaid      string
	}{
		{"chưa thu", "1000000", "0", "1000000", "0"},
		{"thu một phần", "2000000", "500000", "1500000", "500000"},
		{"tất toán đủ", "5000000", "5000000", "0", "5000000"},
		{"thu thừa → còn nợ âm", "1000000", "1200000", "-200000", "1200000"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ps := domain.ComputePayment(mustMoney(t, c.final), mustMoney(t, c.paid))
			if ps.Remaining.String() != c.wantRemaining {
				t.Errorf("Remaining = %q, want %q", ps.Remaining.String(), c.wantRemaining)
			}
			if ps.Paid.String() != c.wantPaid {
				t.Errorf("Paid = %q, want %q", ps.Paid.String(), c.wantPaid)
			}
			if ps.Final.String() != c.final {
				t.Errorf("Final = %q, want %q", ps.Final.String(), c.final)
			}
		})
	}
}

func TestComputePayment_OverpaidIsNegative(t *testing.T) {
	ps := domain.ComputePayment(mustMoney(t, "100"), mustMoney(t, "150"))
	if !ps.Remaining.IsNegative() {
		t.Fatalf("thu thừa phải cho Remaining âm, got %q", ps.Remaining.String())
	}
}

func TestQuantity_FromInt(t *testing.T) {
	q := domain.QuantityFromInt(3)
	if q.String() != "3" {
		t.Fatalf("QuantityFromInt(3) = %q, want 3", q.String())
	}
	if domain.ZeroQty().String() != "0" || !domain.ZeroQty().IsZero() {
		t.Fatal("ZeroQty phải = 0")
	}
}
