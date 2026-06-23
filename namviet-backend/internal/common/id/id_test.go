package id_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/Maneva-AI/namviet-backend/internal/common/id"
)

func TestNew_IsV7AndUnique(t *testing.T) {
	a := id.New()
	b := id.New()

	if a == uuid.Nil {
		t.Fatal("New() trả uuid.Nil")
	}
	if a == b {
		t.Fatal("hai ID liên tiếp trùng nhau")
	}
	if v := a.Version(); v != 7 {
		t.Fatalf("version = %d, muốn 7 (v7)", v)
	}
}

func TestNewString_Parseable(t *testing.T) {
	s := id.NewString()
	parsed, err := id.Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q): %v", s, err)
	}
	if parsed.String() != s {
		t.Fatalf("round-trip lệch: %q vs %q", parsed.String(), s)
	}
}
