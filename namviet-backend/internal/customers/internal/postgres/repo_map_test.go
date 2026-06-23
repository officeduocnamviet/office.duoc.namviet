package postgres

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
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
		{"integer", num(1500000, 0), "1500000"},
		{"two decimals", num(123456, -2), "1234.56"},
		{"negative", num(-505, -1), "-50.5"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := numericToMoney(c.in)
			if got.String() != c.want {
				t.Fatalf("numericToMoney = %q, want %q", got.String(), c.want)
			}
		})
	}
}

func TestNumericToMoney_NoFloatPrecisionLoss(t *testing.T) {
	got := numericToMoney(num(1000000000000001, -2)) // 10000000000000.01
	if got.String() != "10000000000000.01" {
		t.Fatalf("mất chính xác: %q", got.String())
	}
}

// Quy tắc chọn nguồn công nợ ở mức row: live hợp lệ (kể cả 0) → ưu tiên live; live
// invalid (phòng thủ) → rơi về cột tĩnh.
func TestCustomerRowToDomain_DebtPrefersLive(t *testing.T) {
	c := customerRowToDomain(customerRow{
		ID:           5,
		Name:         "Cty A",
		CustomerType: "B2B",
		CurrentDebt:  num(999999, 0),  // cột tĩnh stale
		LiveDebt:     num(1500000, 0), // live (ưu tiên)
		B2bMetadata:  []byte(`{}`),
	})
	if c.Debt.Source != domain.DebtSourceLive {
		t.Fatalf("source = %q, want live", c.Debt.Source)
	}
	if c.Debt.Amount.String() != "1500000" {
		t.Fatalf("amount = %q, want 1500000 (live)", c.Debt.Amount.String())
	}
	if c.Debt.Static.String() != "999999" {
		t.Fatalf("phải giữ static = 999999, got %q", c.Debt.Static.String())
	}
}

func TestCustomerRowToDomain_DebtFallbackStaticWhenLiveInvalid(t *testing.T) {
	c := customerRowToDomain(customerRow{
		ID:           6,
		Name:         "Cty B",
		CustomerType: "B2B",
		CurrentDebt:  num(750000, 0),
		LiveDebt:     pgtype.Numeric{Valid: false}, // không có dữ liệu live
		B2bMetadata:  []byte(`{}`),
	})
	if c.Debt.Source != domain.DebtSourceStatic {
		t.Fatalf("source = %q, want static", c.Debt.Source)
	}
	if c.Debt.Amount.String() != "750000" {
		t.Fatalf("amount = %q, want 750000 (static)", c.Debt.Amount.String())
	}
}

// B2C không gắn B2BProfile; B2B có metadata → parse đúng MST/hạn mức (tiền decimal).
func TestCustomerRowToDomain_B2CNoProfile(t *testing.T) {
	c := customerRowToDomain(customerRow{
		ID: 1, Name: "Chị Lan", CustomerType: "B2C",
		LiveDebt: num(0, 0), B2bMetadata: []byte(`{"tax_code":"x"}`),
	})
	if c.IsB2B() {
		t.Fatal("B2C không được IsB2B")
	}
	if c.B2B != nil {
		t.Fatalf("B2C không được có B2BProfile, got %+v", c.B2B)
	}
}

func TestCustomerRowToDomain_B2BProfileParsed(t *testing.T) {
	raw := []byte(`{"tax_code":"0312345678","debt_limit":50000000,"payment_term":30,"sales_staff_id":"abc-uuid"}`)
	c := customerRowToDomain(customerRow{
		ID: 2, Name: "Cty C", CustomerType: "B2B",
		LiveDebt: num(0, 0), B2bMetadata: raw,
	})
	if c.B2B == nil {
		t.Fatal("B2B phải có profile")
	}
	if c.B2B.TaxCode != "0312345678" {
		t.Fatalf("tax_code = %q", c.B2B.TaxCode)
	}
	if c.B2B.DebtLimit.String() != "50000000" {
		t.Fatalf("debt_limit = %q, want 50000000", c.B2B.DebtLimit.String())
	}
	if c.B2B.PaymentTerm != 30 {
		t.Fatalf("payment_term = %d, want 30", c.B2B.PaymentTerm)
	}
	if c.B2B.SalesStaffID != "abc-uuid" {
		t.Fatalf("sales_staff_id = %q", c.B2B.SalesStaffID)
	}
}

// debt_limit lưu dạng CHUỖI trong JSON cũ vẫn phải parse chính xác (không float).
func TestParseB2BMetadata_DebtLimitAsString(t *testing.T) {
	p := parseB2BMetadata([]byte(`{"tax_code":"123","debt_limit":"12345678.90"}`))
	if p == nil {
		t.Fatal("phải parse được")
	}
	if p.DebtLimit.String() != "12345678.9" {
		t.Fatalf("debt_limit chuỗi = %q, want 12345678.9", p.DebtLimit.String())
	}
}

func TestParseB2BMetadata_EmptyVariants(t *testing.T) {
	for _, raw := range [][]byte{nil, {}, []byte(`{}`), []byte(`{"tax_code":""}`)} {
		if p := parseB2BMetadata(raw); p != nil {
			t.Fatalf("metadata rỗng %q phải trả nil, got %+v", raw, p)
		}
	}
}

func TestParseB2BMetadata_DirtyDebtLimitIgnored(t *testing.T) {
	// debt_limit bẩn (không parse được) → bỏ qua trường đó, profile vẫn còn MST.
	p := parseB2BMetadata([]byte(`{"tax_code":"123","debt_limit":"abc"}`))
	if p == nil || p.TaxCode != "123" {
		t.Fatalf("profile phải giữ tax_code: %+v", p)
	}
	if !p.DebtLimit.IsZero() {
		t.Fatalf("debt_limit bẩn phải về 0, got %q", p.DebtLimit.String())
	}
}

// B2B nhưng metadata rỗng → vẫn gắn profile rỗng (FE biết là B2B).
func TestCustomerRowToDomain_B2BEmptyMetadataStillProfile(t *testing.T) {
	c := customerRowToDomain(customerRow{
		ID: 3, Name: "Cty D", CustomerType: "B2B",
		LiveDebt: num(0, 0), B2bMetadata: []byte(`{}`),
	})
	if c.B2B == nil {
		t.Fatal("B2B metadata rỗng vẫn phải có profile (rỗng)")
	}
	if c.B2B.TaxCode != "" || !c.B2B.DebtLimit.Equal(money.Zero()) {
		t.Fatalf("profile rỗng sai: %+v", c.B2B)
	}
}

// Đảm bảo Repo vẫn thoả port domain (compile-time guard cũng có trong repo.go).
var _ domain.Repository = (*Repo)(nil)
