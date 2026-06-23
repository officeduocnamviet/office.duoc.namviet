package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
)

func TestConstructors_Kind(t *testing.T) {
	cases := []struct {
		name string
		err  *apperr.Error
		want apperr.Kind
	}{
		{"not_found", apperr.NotFound("x"), apperr.KindNotFound},
		{"conflict", apperr.Conflict("x"), apperr.KindConflict},
		{"validation", apperr.Validation("x"), apperr.KindValidation},
		{"unauthorized", apperr.Unauthorized("x"), apperr.KindUnauthorized},
		{"forbidden", apperr.Forbidden("x"), apperr.KindForbidden},
		{"internal", apperr.Internal("x"), apperr.KindInternal},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.err.Kind != c.want {
				t.Fatalf("Kind = %v, want %v", c.err.Kind, c.want)
			}
		})
	}
}

func TestUnauthorizedCode_StableCode(t *testing.T) {
	e := apperr.UnauthorizedCode("refresh_reuse_detected", "reuse")
	if e.Kind != apperr.KindUnauthorized {
		t.Fatalf("Kind = %v, want Unauthorized", e.Kind)
	}
	if e.Code != "refresh_reuse_detected" {
		t.Fatalf("Code = %q", e.Code)
	}
}

func TestWithCause_UnwrapAndMessage(t *testing.T) {
	cause := errors.New("boom")
	e := apperr.Internal("op thất bại").WithCause(cause)

	if !errors.Is(e, cause) {
		t.Fatal("errors.Is không thấy cause")
	}
	if got := e.Error(); got != "op thất bại: boom" {
		t.Fatalf("Error() = %q", got)
	}
}

func TestAsError_AndKindOf(t *testing.T) {
	wrapped := fmt.Errorf("ngữ cảnh: %w", apperr.NotFound("user"))

	got, ok := apperr.AsError(wrapped)
	if !ok {
		t.Fatal("AsError không trích được apperr qua %w")
	}
	if got.Kind != apperr.KindNotFound {
		t.Fatalf("Kind = %v", got.Kind)
	}
	if apperr.KindOf(wrapped) != apperr.KindNotFound {
		t.Fatalf("KindOf = %v", apperr.KindOf(wrapped))
	}

	// Lỗi thường → KindInternal, AsError = false.
	plain := errors.New("plain")
	if _, ok := apperr.AsError(plain); ok {
		t.Fatal("AsError nên false với lỗi thường")
	}
	if apperr.KindOf(plain) != apperr.KindInternal {
		t.Fatal("KindOf lỗi thường nên là Internal")
	}
}
