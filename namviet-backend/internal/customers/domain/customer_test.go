package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
)

func mustMoney(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money.FromString(%q): %v", s, err)
	}
	return m
}

func TestIsB2B(t *testing.T) {
	if !(domain.Customer{Type: domain.TypeB2B}).IsB2B() {
		t.Fatal("TypeB2B phải IsB2B() = true")
	}
	if (domain.Customer{Type: domain.TypeB2C}).IsB2B() {
		t.Fatal("TypeB2C phải IsB2B() = false")
	}
}

// SelectDebt là quy tắc nghiệp vụ dễ sai: ưu tiên LIVE (kể cả 0 đồng) khi có dữ
// liệu đơn; chỉ rơi về cột tĩnh stale khi không có dữ liệu live.
func TestSelectDebt_PrefersLiveWhenAvailable(t *testing.T) {
	live := mustMoney(t, "1500000")
	static := mustMoney(t, "999999") // cột tĩnh lệch (stale)
	got := domain.SelectDebt(live, static, true)

	if got.Source != domain.DebtSourceLive {
		t.Fatalf("source = %q, want live", got.Source)
	}
	if !got.Amount.Equal(live) {
		t.Fatalf("Amount = %q, want %q (live)", got.Amount.String(), live.String())
	}
	// Giữ cả hai con số để minh bạch.
	if !got.Live.Equal(live) || !got.Static.Equal(static) {
		t.Fatalf("phải giữ cả Live (%q) và Static (%q)", got.Live.String(), got.Static.String())
	}
}

func TestSelectDebt_LiveZeroVsStaleStatic(t *testing.T) {
	// Khách thật sự KHÔNG nợ (live = 0) nhưng cột tĩnh còn số cũ → vẫn tin live=0.
	live := money.Zero()
	static := mustMoney(t, "500000")
	got := domain.SelectDebt(live, static, true)

	if got.Source != domain.DebtSourceLive {
		t.Fatalf("source = %q, want live", got.Source)
	}
	if !got.Amount.IsZero() {
		t.Fatalf("Amount = %q, want 0 (live thắng cột tĩnh stale)", got.Amount.String())
	}
}

func TestSelectDebt_FallsBackToStaticWhenNoLiveData(t *testing.T) {
	live := money.Zero()
	static := mustMoney(t, "750000")
	got := domain.SelectDebt(live, static, false)

	if got.Source != domain.DebtSourceStatic {
		t.Fatalf("source = %q, want static", got.Source)
	}
	if !got.Amount.Equal(static) {
		t.Fatalf("Amount = %q, want %q (static)", got.Amount.String(), static.String())
	}
}
