package domain_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

func dec(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("decimal %q: %v", s, err)
	}
	return d
}

func mny(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return m
}

func TestNewDraft_Happy_TinhTienVaVAT(t *testing.T) {
	d, err := domain.NewDraft(domain.DraftInput{
		SupplierName: "NCC A",
		Lines: []domain.DraftLine{
			{ProductID: 100, Quantity: dec(t, "10"), UnitCost: mny(t, "1000"), VATRate: dec(t, "0.08")},
			{ProductID: 200, Quantity: dec(t, "2.5"), UnitCost: mny(t, "2000"), VATRate: dec(t, "0.1")},
		},
	})
	if err != nil {
		t.Fatalf("NewDraft happy phải thành công: %v", err)
	}
	if d.Status != domain.StatusDraft {
		t.Fatalf("status draft mới = %q, want draft", d.Status)
	}
	// line1: 10*1000 = 10000 ; vat = 10000*0.08 = 800
	// line2: 2.5*2000 = 5000  ; vat = 5000*0.1 = 500
	if !d.Lines[0].LineTotal.Equal(mny(t, "10000")) {
		t.Errorf("line1 total = %s, want 10000", d.Lines[0].LineTotal)
	}
	if !d.Lines[0].VATAmount.Equal(mny(t, "800")) {
		t.Errorf("line1 vat = %s, want 800", d.Lines[0].VATAmount)
	}
	if !d.Lines[1].LineTotal.Equal(mny(t, "5000")) {
		t.Errorf("line2 total = %s, want 5000", d.Lines[1].LineTotal)
	}
	if d.Lines[0].LineNo != 1 || d.Lines[1].LineNo != 2 {
		t.Errorf("line_no phải 1,2; got %d,%d", d.Lines[0].LineNo, d.Lines[1].LineNo)
	}
	// total = 10000 + 5000 = 15000 ; vat = 800 + 500 = 1300
	if !d.TotalAmount.Equal(mny(t, "15000")) {
		t.Errorf("total = %s, want 15000", d.TotalAmount)
	}
	if !d.VATAmount.Equal(mny(t, "1300")) {
		t.Errorf("vat total = %s, want 1300", d.VATAmount)
	}
}

func TestNewDraft_Validations(t *testing.T) {
	cases := []struct {
		name string
		in   domain.DraftInput
		want error
	}{
		{"no lines", domain.DraftInput{}, domain.ErrNoLines},
		{"bad product", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 0, Quantity: dec(t, "1"), UnitCost: mny(t, "1")},
		}}, domain.ErrInvalidProductID},
		{"zero qty", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: dec(t, "0"), UnitCost: mny(t, "1")},
		}}, domain.ErrQuantityNotPos},
		{"neg qty", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: dec(t, "-1"), UnitCost: mny(t, "1")},
		}}, domain.ErrQuantityNotPos},
		{"neg cost", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: dec(t, "1"), UnitCost: mny(t, "-1")},
		}}, domain.ErrUnitCostNegative},
		{"neg vat", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: dec(t, "1"), UnitCost: mny(t, "1"), VATRate: dec(t, "-0.1")},
		}}, domain.ErrVATRateNegative},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := domain.NewDraft(tc.in)
			if err != tc.want {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestNewDraft_ZeroVATRate_NoVAT(t *testing.T) {
	d, err := domain.NewDraft(domain.DraftInput{
		Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: dec(t, "3"), UnitCost: mny(t, "5000")}, // vat_rate = 0
		},
	})
	if err != nil {
		t.Fatalf("NewDraft: %v", err)
	}
	if !d.VATAmount.IsZero() {
		t.Errorf("vat phải 0 khi vat_rate=0; got %s", d.VATAmount)
	}
	if !d.TotalAmount.Equal(mny(t, "15000")) {
		t.Errorf("total = %s, want 15000", d.TotalAmount)
	}
}
