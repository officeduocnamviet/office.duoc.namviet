package money_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

func TestZero_IsZero(t *testing.T) {
	z := money.Zero()
	if !z.IsZero() {
		t.Fatalf("Zero() phải IsZero")
	}
	if z.String() != "0" {
		t.Fatalf("Zero().String() = %q, want \"0\"", z.String())
	}
}

func TestFromString_RoundTrip(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"0", "0"},
		{"123456", "123456"},
		{"12345.67", "12345.67"},
		{"-50.5", "-50.5"},
		{"0.001", "0.001"},
	}
	for _, c := range cases {
		m, err := money.FromString(c.in)
		if err != nil {
			t.Fatalf("FromString(%q): %v", c.in, err)
		}
		if m.String() != c.want {
			t.Errorf("FromString(%q).String() = %q, want %q", c.in, m.String(), c.want)
		}
	}
}

func TestFromString_Invalid(t *testing.T) {
	if _, err := money.FromString("abc"); err == nil {
		t.Fatal("FromString(\"abc\") phải lỗi")
	}
}

func TestFromString_EmptyIsZero(t *testing.T) {
	m, err := money.FromString("")
	if err != nil {
		t.Fatalf("FromString(\"\"): %v", err)
	}
	if !m.IsZero() {
		t.Fatalf("FromString(\"\") phải = Zero")
	}
}

func TestFromDecimal_And_Decimal(t *testing.T) {
	d := decimal.RequireFromString("99.99")
	m := money.FromDecimal(d)
	if !m.Decimal().Equal(d) {
		t.Fatalf("Decimal() = %s, want %s", m.Decimal(), d)
	}
}

func TestSub(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"1000000", "300000", "700000"},   // còn nợ = final - đã thu
		{"5000000", "5000000", "0"},       // tất toán đủ
		{"1000000", "1500000", "-500000"}, // thu thừa → âm (over-paid)
		{"12345.67", "0.67", "12345"},
	}
	for _, c := range cases {
		a, _ := money.FromString(c.a)
		b, _ := money.FromString(c.b)
		got := a.Sub(b)
		if got.String() != c.want {
			t.Errorf("(%s).Sub(%s) = %q, want %q", c.a, c.b, got.String(), c.want)
		}
	}
}

func TestSub_NoFloatPrecisionLoss(t *testing.T) {
	a, _ := money.FromString("100000000000.01")
	b, _ := money.FromString("0.01")
	if a.Sub(b).String() != "100000000000" {
		t.Fatalf("mất chính xác: %q", a.Sub(b).String())
	}
}

func TestIsNegative(t *testing.T) {
	neg, _ := money.FromString("-1")
	zero := money.Zero()
	pos, _ := money.FromString("1")
	if !neg.IsNegative() {
		t.Fatal("-1 phải IsNegative")
	}
	if zero.IsNegative() || pos.IsNegative() {
		t.Fatal("0 và 1 KHÔNG được IsNegative")
	}
}

func TestFromInt(t *testing.T) {
	if money.FromInt(110000).String() != "110000" {
		t.Fatalf("FromInt(110000).String() = %q", money.FromInt(110000).String())
	}
	if !money.FromInt(0).IsZero() {
		t.Fatal("FromInt(0) phải IsZero")
	}
}

func TestAdd(t *testing.T) {
	a, _ := money.FromString("100000.01")
	b, _ := money.FromString("0.99")
	if a.Add(b).String() != "100001" {
		t.Fatalf("Add mất chính xác: %q", a.Add(b).String())
	}
	// Cộng dồn 3 dòng credit (511 + 3331) cân với 1 dòng debit 131.
	sum := money.FromInt(100000).Add(money.FromInt(10000))
	if !sum.Equal(money.FromInt(110000)) {
		t.Fatalf("Σ = %s, want 110000", sum)
	}
}

func TestIsPositive(t *testing.T) {
	if !money.FromInt(1).IsPositive() {
		t.Fatal("1 phải IsPositive")
	}
	if money.Zero().IsPositive() {
		t.Fatal("0 KHÔNG được IsPositive")
	}
	if money.FromInt(-1).IsPositive() {
		t.Fatal("-1 KHÔNG được IsPositive")
	}
}

func TestMul(t *testing.T) {
	// line_amount * vat_rate (chưa làm tròn — giữ nguyên scale).
	cases := []struct {
		amount string
		rate   string
		want   string
	}{
		{"100000", "0.08", "8000"},   // tròn
		{"100000", "0.1", "10000"},   // tròn
		{"100000", "0.05", "5000"},   // tròn
		{"33333", "0.08", "2666.64"}, // lẻ — chưa làm tròn
		{"100000", "0", "0"},         // rate 0
	}
	for _, c := range cases {
		amt := money.FromDecimal(decimal.RequireFromString(c.amount))
		rate := decimal.RequireFromString(c.rate)
		got := amt.Mul(rate)
		if got.String() != c.want {
			t.Errorf("(%s).Mul(%s) = %q, want %q", c.amount, c.rate, got.String(), c.want)
		}
	}
}

func TestRoundVND(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"8000", "8000"},    // đã tròn
		{"2666.64", "2667"}, // 0.64 → lên
		{"2666.4", "2666"},  // 0.4 → xuống
		{"2666.5", "2667"},  // 0.5 → HALF-UP lên
		{"0.5", "1"},        // 0.5 → 1
		{"0.49", "0"},       // xuống
		{"123456.999", "123457"},
	}
	for _, c := range cases {
		in := money.FromDecimal(decimal.RequireFromString(c.in))
		got := in.RoundVND()
		if got.String() != c.want {
			t.Errorf("RoundVND(%s) = %q, want %q", c.in, got.String(), c.want)
		}
	}
}

func TestEqual(t *testing.T) {
	a, _ := money.FromString("10.50")
	b := money.FromDecimal(decimal.RequireFromString("10.50"))
	c, _ := money.FromString("10.51")
	if !a.Equal(b) {
		t.Fatalf("a phải Equal b")
	}
	if a.Equal(c) {
		t.Fatalf("a KHÔNG được Equal c")
	}
}
