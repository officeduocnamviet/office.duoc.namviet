package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

// TestDerivePaymentStatus kiểm quy tắc suy diễn payment_status THUẦN: unpaid khi
// chưa thu, partial khi thu một phần, paid khi đủ/thừa; biên đơn 0 đồng = paid.
func TestDerivePaymentStatus(t *testing.T) {
	cases := []struct {
		name        string
		final, paid string
		want        domain.PaymentStatus
	}{
		{"chưa thu", "1000000", "0", domain.PaymentUnpaid},
		{"thu một phần", "2000000", "500000", domain.PaymentPartial},
		{"thu đủ", "5000000", "5000000", domain.PaymentPaid},
		{"thu thừa → paid", "1000000", "1200000", domain.PaymentPaid},
		{"thu sát đủ thiếu 1đ → partial", "1000000", "999999", domain.PaymentPartial},
		{"đơn 0 đồng, thu 0 → paid", "0", "0", domain.PaymentPaid},
		{"paid âm (phòng thủ) → unpaid", "1000000", "-5", domain.PaymentUnpaid},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := domain.DerivePaymentStatus(mustMoney(t, c.final), mustMoney(t, c.paid))
			if got != c.want {
				t.Fatalf("DerivePaymentStatus(%s, %s) = %q, want %q", c.final, c.paid, got, c.want)
			}
		})
	}
}
