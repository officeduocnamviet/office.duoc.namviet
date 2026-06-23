package pagination_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
)

func TestEncodeDecode_RoundTrip(t *testing.T) {
	for _, id := range []int64{1, 42, 999999, 9223372036854775807} {
		c := pagination.EncodeID(id)
		if c == "" {
			t.Fatalf("EncodeID(%d) rỗng", id)
		}
		got, err := pagination.DecodeID(c)
		if err != nil {
			t.Fatalf("DecodeID(%q): %v", c, err)
		}
		if got != id {
			t.Fatalf("round-trip id = %d, want %d", got, id)
		}
	}
}

func TestEncode_Opaque(t *testing.T) {
	// Cursor phải là chuỗi opaque (base64), KHÔNG để lộ id thô dạng "42".
	c := pagination.EncodeID(42)
	if c == "42" {
		t.Fatalf("cursor không được là id thô: %q", c)
	}
}

func TestDecode_Empty(t *testing.T) {
	// Cursor rỗng = trang đầu → id 0, không lỗi.
	got, err := pagination.DecodeID("")
	if err != nil {
		t.Fatalf("DecodeID(\"\"): %v", err)
	}
	if got != 0 {
		t.Fatalf("DecodeID(\"\") = %d, want 0", got)
	}
}

func TestDecode_Invalid(t *testing.T) {
	if _, err := pagination.DecodeID("@@@not-base64@@@"); err == nil {
		t.Fatal("DecodeID chuỗi rác phải lỗi")
	}
	if _, err := pagination.DecodeID("bm90LWFuLWludA=="); err == nil { // base64("not-an-int")
		t.Fatal("DecodeID payload không phải số phải lỗi")
	}
}
