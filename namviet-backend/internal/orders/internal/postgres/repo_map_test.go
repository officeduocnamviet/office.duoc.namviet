package postgres

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

func num(mantissa int64, exp int32) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(mantissa), Exp: exp, Valid: true}
}

func TestNumericToMoney_NoFloat(t *testing.T) {
	cases := []struct {
		name string
		in   pgtype.Numeric
		want string
	}{
		{"invalid->zero", pgtype.Numeric{Valid: false}, "0"},
		{"nan->zero", pgtype.Numeric{NaN: true, Valid: true}, "0"},
		{"integer VND", num(1000000, 0), "1000000"},
		{"2 decimal", num(1234567, -2), "12345.67"},
		{"big no precision loss", num(10000000000001, -2), "100000000000.01"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := numericToMoney(c.in).String(); got != c.want {
				t.Fatalf("numericToMoney = %q, want %q", got, c.want)
			}
		})
	}
}

func TestNanoToTimestamptz(t *testing.T) {
	if nanoToTimestamptz(0).Valid {
		t.Fatal("0 nano phải = invalid (NULL → trang đầu)")
	}
	ts := time.Date(2026, 6, 1, 10, 30, 0, 0, time.UTC)
	got := nanoToTimestamptz(ts.UnixNano())
	if !got.Valid || !got.Time.Equal(ts) {
		t.Fatalf("nanoToTimestamptz round-trip sai: %+v", got)
	}
}

// orderRowToDomain phải đặt Remaining = Final - Paid (suy diễn domain), KHÔNG ép
// float, và giữ NULL tiền → 0.
func TestOrderRowToDomain_PaymentDerived(t *testing.T) {
	row := orderRow{
		ID:          "ord-1",
		Code:        "HD1",
		Status:      "COMPLETED",
		OrderType:   "B2B",
		FinalAmount: num(2000000, 0),
		PaidAmount:  num(500000, 0),
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	o := orderRowToDomain(row)
	if o.Final.String() != "2000000" {
		t.Fatalf("Final = %q", o.Final.String())
	}
	if o.Payment.Paid.String() != "500000" {
		t.Fatalf("Paid = %q", o.Payment.Paid.String())
	}
	if o.Payment.Remaining.String() != "1500000" {
		t.Fatalf("Remaining = %q, want 1500000", o.Payment.Remaining.String())
	}
}

func TestOrderRowToDomain_NullPaidIsZero(t *testing.T) {
	row := orderRow{ID: "x", Code: "HD2", FinalAmount: num(1000000, 0), PaidAmount: pgtype.Numeric{Valid: false}}
	o := orderRowToDomain(row)
	if !o.Payment.Paid.IsZero() {
		t.Fatalf("paid NULL phải = 0, got %q", o.Payment.Paid.String())
	}
	if o.Payment.Remaining.String() != "1000000" {
		t.Fatalf("Remaining = %q, want 1000000 (chưa thu)", o.Payment.Remaining.String())
	}
}

// Đảm bảo Repo vẫn thoả port domain (compile-time guard cũng có trong repo.go).
var _ domain.Repository = (*Repo)(nil)
