package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

// TestStatus_Valid kiểm tập trạng thái hợp lệ khớp CHECK của public.orders.
func TestStatus_Valid(t *testing.T) {
	valid := []domain.Status{
		domain.StatusPending, domain.StatusConfirmed, domain.StatusShipping,
		domain.StatusCompleted, domain.StatusCancelled, domain.StatusRefunded,
	}
	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("status %q phải hợp lệ", s)
		}
	}
	for _, s := range []domain.Status{"", "pending", "DONE", "draft"} {
		if domain.Status(s).Valid() {
			t.Errorf("status %q KHÔNG được hợp lệ", s)
		}
	}
}

// TestCanTransition_Matrix kiểm TOÀN BỘ ma trận chuyển trạng thái (P4a — chỉ các
// chuyển KHÔNG cần cross-module: PENDING→CONFIRMED→SHIPPING→COMPLETED tiến theo
// luồng; PENDING/CONFIRMED→CANCELLED; chặn nhảy bước/lùi). Chuyển SHIPPING (cần
// trừ kho) và REFUNDED (cần đảo sổ/hoàn kho) HOÃN sang P4b — CanTransition vẫn
// CHO PHÉP về mặt state (nghiệp vụ ràng buộc ở use-case khi đủ primitive), nhưng
// P4a chỉ expose use-case Confirm/Complete/Cancel.
func TestCanTransition_Matrix(t *testing.T) {
	// allowed[from] = tập to hợp lệ.
	allowed := map[domain.Status][]domain.Status{
		domain.StatusPending:   {domain.StatusConfirmed, domain.StatusCancelled},
		domain.StatusConfirmed: {domain.StatusShipping, domain.StatusCancelled},
		domain.StatusShipping:  {domain.StatusCompleted},
		domain.StatusCompleted: {domain.StatusRefunded},
		domain.StatusCancelled: {},
		domain.StatusRefunded:  {},
	}
	all := []domain.Status{
		domain.StatusPending, domain.StatusConfirmed, domain.StatusShipping,
		domain.StatusCompleted, domain.StatusCancelled, domain.StatusRefunded,
	}
	for _, from := range all {
		okSet := map[domain.Status]bool{}
		for _, to := range allowed[from] {
			okSet[to] = true
		}
		for _, to := range all {
			got := domain.CanTransition(from, to)
			want := okSet[to]
			if got != want {
				t.Errorf("CanTransition(%s, %s) = %v, want %v", from, to, got, want)
			}
		}
		// Tự chuyển về chính nó luôn bị chặn (no-op không phải transition hợp lệ).
		if domain.CanTransition(from, from) {
			t.Errorf("CanTransition(%s, %s) phải false (tự chuyển)", from, from)
		}
	}
}

// TestCanTransition_RejectsSkipAndBackward khẳng định một số case quan trọng:
// nhảy bước (PENDING→COMPLETED) và lùi (CONFIRMED→PENDING) đều bị chặn.
func TestCanTransition_RejectsSkipAndBackward(t *testing.T) {
	bad := []struct{ from, to domain.Status }{
		{domain.StatusPending, domain.StatusShipping},    // nhảy bước
		{domain.StatusPending, domain.StatusCompleted},   // nhảy bước
		{domain.StatusConfirmed, domain.StatusPending},   // lùi
		{domain.StatusConfirmed, domain.StatusCompleted}, // nhảy bước
		{domain.StatusShipping, domain.StatusConfirmed},  // lùi
		{domain.StatusCompleted, domain.StatusShipping},  // lùi
		{domain.StatusCancelled, domain.StatusConfirmed}, // từ terminal
	}
	for _, c := range bad {
		if domain.CanTransition(c.from, c.to) {
			t.Errorf("CanTransition(%s, %s) phải bị chặn", c.from, c.to)
		}
	}
}
