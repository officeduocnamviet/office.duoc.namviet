package postgres

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

func num(mantissa int64, exp int32) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(mantissa), Exp: exp, Valid: true}
}

func TestNumericToQuantity(t *testing.T) {
	cases := []struct {
		name string
		in   pgtype.Numeric
		want string
	}{
		{"invalid->zero", pgtype.Numeric{Valid: false}, "0"},
		{"nan->zero", pgtype.Numeric{NaN: true, Valid: true}, "0"},
		{"integer", num(150, 0), "150"},
		{"one decimal", num(125, -1), "12.5"},
		{"zero", num(0, 0), "0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := numericToQuantity(c.in)
			if got.String() != c.want {
				t.Fatalf("numericToQuantity(%v) = %q, want %q", c.in, got.String(), c.want)
			}
		})
	}
}

func TestNumericToQuantity_NoFloatPrecisionLoss(t *testing.T) {
	got := numericToQuantity(num(10000000000001, -2)) // 100000000000.01
	if got.String() != "100000000000.01" {
		t.Fatalf("mất chính xác: %q", got.String())
	}
}

func TestNumericToMoney_CostNoFloat(t *testing.T) {
	got := numericToMoney(num(1234567, -2)) // 12345.67 giá vốn lô
	if got.String() != "12345.67" {
		t.Fatalf("inbound_price = %q, want 12345.67", got.String())
	}
	if !numericToMoney(pgtype.Numeric{Valid: false}).IsZero() {
		t.Fatal("NULL inbound_price phải = 0 tiền")
	}
}

// Đảm bảo Repo vẫn thoả port domain (compile-time guard cũng có trong repo.go).
var _ domain.Repository = (*Repo)(nil)
