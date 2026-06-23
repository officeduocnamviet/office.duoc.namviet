package postgres

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/catalog/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

func num(mantissa int64, exp int32) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(mantissa), Exp: exp, Valid: true}
}

func TestNumericToMoney(t *testing.T) {
	cases := []struct {
		name string
		in   pgtype.Numeric
		want string
	}{
		{"invalid->zero", pgtype.Numeric{Valid: false}, "0"},
		{"nan->zero", pgtype.Numeric{NaN: true, Valid: true}, "0"},
		{"integer", num(123456, 0), "123456"},
		{"two decimals", num(1234567, -2), "12345.67"},
		{"negative", num(-505, -1), "-50.5"},
		{"zero", num(0, 0), "0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := numericToMoney(c.in)
			if got.String() != c.want {
				t.Fatalf("numericToMoney(%v) = %q, want %q", c.in, got.String(), c.want)
			}
		})
	}
}

func TestNumericToMoney_NoFloatPrecisionLoss(t *testing.T) {
	// Giá trị mà float64 không biểu diễn chính xác — phải giữ nguyên qua decimal.
	got := numericToMoney(num(1000000000000001, -2)) // 10000000000000.01
	if got.String() != "10000000000000.01" {
		t.Fatalf("mất chính xác: %q", got.String())
	}
}

func TestProductRowToDomain_Nullables(t *testing.T) {
	p := productRowToDomain(productRow{
		ID:            7,
		Name:          "Paracetamol",
		Status:        "active",
		InvoicePrice:  num(1500000, -2), // 15000.00
		ActualCost:    pgtype.Numeric{Valid: false},
		ProductImages: nil,
		// các con trỏ để nil → phải về giá trị rỗng/mặc định
	})
	if p.SKU != "" || p.Barcode != "" || p.CategoryName != "" {
		t.Fatalf("nullable text phải rỗng: %+v", p)
	}
	if p.CategoryID != nil || p.ManufacturerID != nil {
		t.Fatalf("nullable id phải nil")
	}
	if p.ConversionFactor != 1 {
		t.Fatalf("conversion_factor nil phải default 1, got %d", p.ConversionFactor)
	}
	if p.Images == nil {
		t.Fatalf("Images phải là slice rỗng, không nil (để JSON ra [])")
	}
	if !p.ActualCost.Equal(money.Zero()) {
		t.Fatalf("ActualCost NULL phải = 0")
	}
	if p.InvoicePrice.String() != "15000" {
		t.Fatalf("InvoicePrice = %q, want 15000", p.InvoicePrice.String())
	}
}

func TestUnitRowToDomain_DefaultsAndMoney(t *testing.T) {
	u := unitRowToDomain(appdb.ListProductUnitsRow{
		ID:        3,
		UnitName:  "Viên",
		PriceSell: num(500, 0),
		// IsBase nil → false; IsDirectSale nil → true; ProductID nil → 0
	})
	if u.IsBase {
		t.Fatalf("IsBase nil phải false")
	}
	if !u.IsDirectSale {
		t.Fatalf("IsDirectSale nil phải default true")
	}
	if u.ProductID != 0 {
		t.Fatalf("ProductID nil phải 0")
	}
	if u.PriceSell.String() != "500" {
		t.Fatalf("PriceSell = %q", u.PriceSell.String())
	}
}

// Đảm bảo Repo vẫn thoả port domain (compile-time guard cũng có trong repo.go).
var _ domain.Repository = (*Repo)(nil)
