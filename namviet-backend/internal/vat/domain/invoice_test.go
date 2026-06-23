package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

func mny(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return m
}

func rate(s string) decimal.Decimal { return decimal.RequireFromString(s) }

const issueDay = "2026-06-20"

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("date %q: %v", s, err)
	}
	return d
}

func oneLine(t *testing.T) []domain.LineInput {
	t.Helper()
	return []domain.LineInput{{
		ProductID: 1, Description: "Paracetamol",
		Quantity: mny(t, "10"), UnitPrice: mny(t, "10000"), VATRate: rate("0.08"),
	}}
}

func TestStatus_Valid(t *testing.T) {
	for _, s := range []domain.Status{domain.StatusDraft, domain.StatusIssued, domain.StatusCancelled} {
		if !s.Valid() {
			t.Errorf("%q phải hợp lệ", s)
		}
	}
	for _, s := range []string{"", "ISSUED", "void", "x"} {
		if domain.Status(s).Valid() {
			t.Errorf("%q phải KHÔNG hợp lệ", s)
		}
	}
}

func TestBuildInvoice_Happy_ComputeAndBalance(t *testing.T) {
	// 2 dòng, thuế suất khác nhau (8% và 10%) — vat_rate là INPUT từng dòng.
	lines := []domain.LineInput{
		{ProductID: 1, Description: "A", Quantity: mny(t, "10"), UnitPrice: mny(t, "10000"), VATRate: rate("0.08")},
		{ProductID: 2, Description: "B", Quantity: mny(t, "3"), UnitPrice: mny(t, "50000"), VATRate: rate("0.10")},
	}
	inv, err := domain.BuildInvoice("HD123", "0312345678", "C26TYY", mustDate(t, issueDay), lines)
	if err != nil {
		t.Fatalf("hợp lệ phải PASS: %v", err)
	}
	// Dòng 1: 10*10000=100000; vat=8000. Dòng 2: 3*50000=150000; vat=15000.
	if inv.Subtotal.String() != "250000" {
		t.Fatalf("subtotal = %s, want 250000", inv.Subtotal)
	}
	if inv.VATAmount.String() != "23000" {
		t.Fatalf("vat_amount = %s, want 23000", inv.VATAmount)
	}
	if inv.Total.String() != "273000" {
		t.Fatalf("total = %s, want 273000", inv.Total)
	}
	// total = subtotal + vat (ép cân).
	if !inv.Total.Equal(inv.Subtotal.Add(inv.VATAmount)) {
		t.Fatal("total phải = subtotal + vat")
	}
	if len(inv.Lines) != 2 || inv.Lines[0].LineNo != 1 || inv.Lines[1].LineNo != 2 {
		t.Fatalf("line_no phải theo thứ tự 1,2: %+v", inv.Lines)
	}
	if inv.Lines[0].LineAmount.String() != "100000" || inv.Lines[0].LineVAT.String() != "8000" {
		t.Fatalf("dòng 1 sai: amount=%s vat=%s", inv.Lines[0].LineAmount, inv.Lines[0].LineVAT)
	}
}

// VAT làm tròn từng dòng về ĐỒNG (scale-0). line_amount lẻ → line_vat tròn; Σ
// line_vat = vat_amount header (cân khít, không sai số làm tròn).
func TestBuildInvoice_VATRounding_PerLine(t *testing.T) {
	// 33333 * 0.08 = 2666.64 → làm tròn 2667 (HALF-UP).
	lines := []domain.LineInput{
		{Quantity: mny(t, "1"), UnitPrice: mny(t, "33333"), VATRate: rate("0.08")},
		{Quantity: mny(t, "1"), UnitPrice: mny(t, "33333"), VATRate: rate("0.08")},
	}
	inv, err := domain.BuildInvoice("HD1", "0312345678", "C26TYY", mustDate(t, issueDay), lines)
	if err != nil {
		t.Fatalf("PASS: %v", err)
	}
	// Mỗi dòng vat = 2667; tổng = 5334 (KHÔNG phải round(5333.28)=5333).
	if inv.Lines[0].LineVAT.String() != "2667" {
		t.Fatalf("line vat = %s, want 2667 (per-line round)", inv.Lines[0].LineVAT)
	}
	if inv.VATAmount.String() != "5334" {
		t.Fatalf("vat_amount = %s, want 5334 (Σ dòng đã tròn)", inv.VATAmount)
	}
	if inv.Subtotal.String() != "66666" {
		t.Fatalf("subtotal = %s, want 66666", inv.Subtotal)
	}
	if inv.Total.String() != "72000" {
		t.Fatalf("total = %s, want 72000", inv.Total)
	}
	// vat_amount = Σ line_vat (cân khít với dòng).
	sum := money.Zero()
	for _, l := range inv.Lines {
		sum = sum.Add(l.LineVAT)
	}
	if !sum.Equal(inv.VATAmount) {
		t.Fatalf("Σ line_vat=%s phải = vat_amount=%s", sum, inv.VATAmount)
	}
}

func TestBuildInvoice_ZeroRate_NoVAT(t *testing.T) {
	lines := []domain.LineInput{
		{Quantity: mny(t, "2"), UnitPrice: mny(t, "100000"), VATRate: rate("0")},
	}
	inv, err := domain.BuildInvoice("HD1", "0312345678", "C26TYY", mustDate(t, issueDay), lines)
	if err != nil {
		t.Fatalf("PASS: %v", err)
	}
	if inv.VATAmount.String() != "0" || inv.Total.String() != "200000" {
		t.Fatalf("rate 0 → vat=0 total=200000, got vat=%s total=%s", inv.VATAmount, inv.Total)
	}
}

func TestBuildInvoice_EmptyTaxCode(t *testing.T) {
	for _, mst := range []string{"", "   ", "\t"} {
		_, err := domain.BuildInvoice("HD1", mst, "C26TYY", mustDate(t, issueDay), oneLine(t))
		if !errors.Is(err, domain.ErrEmptyTaxCode) {
			t.Fatalf("MST %q phải lỗi ErrEmptyTaxCode, got %v", mst, err)
		}
	}
}

func TestBuildInvoice_EmptyOrderCode(t *testing.T) {
	_, err := domain.BuildInvoice("  ", "0312345678", "C26TYY", mustDate(t, issueDay), oneLine(t))
	if !errors.Is(err, domain.ErrEmptyOrderCode) {
		t.Fatalf("order_code rỗng phải lỗi, got %v", err)
	}
}

func TestBuildInvoice_EmptySerial(t *testing.T) {
	_, err := domain.BuildInvoice("HD1", "0312345678", "", mustDate(t, issueDay), oneLine(t))
	if !errors.Is(err, domain.ErrEmptySerial) {
		t.Fatalf("serial rỗng phải lỗi, got %v", err)
	}
}

func TestBuildInvoice_NoLines(t *testing.T) {
	_, err := domain.BuildInvoice("HD1", "0312345678", "C26TYY", mustDate(t, issueDay), nil)
	if !errors.Is(err, domain.ErrNoLines) {
		t.Fatalf("không dòng phải lỗi ErrNoLines, got %v", err)
	}
}

func TestBuildInvoice_NegativeLineValues(t *testing.T) {
	cases := []domain.LineInput{
		{Quantity: mny(t, "-1"), UnitPrice: mny(t, "1000"), VATRate: rate("0.08")},
		{Quantity: mny(t, "1"), UnitPrice: mny(t, "-1000"), VATRate: rate("0.08")},
		{Quantity: mny(t, "1"), UnitPrice: mny(t, "1000"), VATRate: rate("-0.08")},
	}
	for i, l := range cases {
		_, err := domain.BuildInvoice("HD1", "0312345678", "C26TYY", mustDate(t, issueDay), []domain.LineInput{l})
		if err == nil {
			t.Fatalf("case %d: giá trị âm phải lỗi", i)
		}
	}
}
