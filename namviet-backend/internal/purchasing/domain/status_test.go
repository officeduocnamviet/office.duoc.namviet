package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

func TestStatus_Valid(t *testing.T) {
	valid := []domain.Status{
		domain.StatusDraft, domain.StatusOrdered, domain.StatusReceived,
		domain.StatusPaid, domain.StatusCancelled,
	}
	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("%q phải hợp lệ", s)
		}
	}
	for _, s := range []domain.Status{"", "DRAFT", "shipping", "foo"} {
		if domain.Status(s).Valid() {
			t.Errorf("%q KHÔNG được hợp lệ", s)
		}
	}
}

func TestCanTransition(t *testing.T) {
	allowed := [][2]domain.Status{
		{domain.StatusDraft, domain.StatusOrdered},
		{domain.StatusDraft, domain.StatusCancelled},
		{domain.StatusOrdered, domain.StatusReceived},
		{domain.StatusOrdered, domain.StatusCancelled},
		{domain.StatusReceived, domain.StatusPaid},
	}
	for _, tc := range allowed {
		if !domain.CanTransition(tc[0], tc[1]) {
			t.Errorf("phải cho chuyển %q → %q", tc[0], tc[1])
		}
	}

	forbidden := [][2]domain.Status{
		{domain.StatusDraft, domain.StatusReceived},  // nhảy bước
		{domain.StatusDraft, domain.StatusPaid},       // nhảy bước
		{domain.StatusOrdered, domain.StatusDraft},    // lùi
		{domain.StatusReceived, domain.StatusCancelled}, // huỷ sau nhập kho (cần đảo — HOÃN)
		{domain.StatusReceived, domain.StatusReceived}, // tự-chuyển
		{domain.StatusPaid, domain.StatusReceived},     // terminal
		{domain.StatusCancelled, domain.StatusOrdered}, // terminal
	}
	for _, tc := range forbidden {
		if domain.CanTransition(tc[0], tc[1]) {
			t.Errorf("KHÔNG được cho chuyển %q → %q", tc[0], tc[1])
		}
	}
}
