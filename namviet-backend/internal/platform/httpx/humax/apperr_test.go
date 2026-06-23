package humax

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
)

func TestFromAppErr_MapsKindToStatusAndCode(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"not_found", apperr.NotFound("x"), http.StatusNotFound, "not_found"},
		{"conflict", apperr.Conflict("x"), http.StatusConflict, "conflict"},
		{"validation", apperr.Validation("x"), http.StatusUnprocessableEntity, "validation_error"},
		{"unauthorized", apperr.Unauthorized("x"), http.StatusUnauthorized, "unauthorized"},
		{"forbidden", apperr.Forbidden("x"), http.StatusForbidden, "forbidden"},
		{"internal", apperr.Internal("x"), http.StatusInternalServerError, "internal"},
		{
			"stable_code_kept",
			apperr.UnauthorizedCode("refresh_reuse_detected", "reuse"),
			http.StatusUnauthorized, "refresh_reuse_detected",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			se := FromAppErr(c.err)
			if se.GetStatus() != c.wantStatus {
				t.Fatalf("status = %d, want %d", se.GetStatus(), c.wantStatus)
			}
			ce, ok := se.(*codeError)
			if !ok {
				t.Fatalf("không phải *codeError: %T", se)
			}
			if ce.ErrBody.Code != c.wantCode {
				t.Fatalf("code = %q, want %q", ce.ErrBody.Code, c.wantCode)
			}
		})
	}
}

func TestFromAppErr_PlainErrorIsInternal(t *testing.T) {
	se := FromAppErr(errors.New("boom"))
	if se.GetStatus() != http.StatusInternalServerError {
		t.Fatalf("status = %d", se.GetStatus())
	}
	ce := se.(*codeError)
	if ce.ErrBody.Code != "internal" {
		t.Fatalf("code = %q", ce.ErrBody.Code)
	}
	// Không rò chi tiết lỗi gốc.
	if ce.ErrBody.Message == "boom" {
		t.Fatal("không được rò message lỗi gốc")
	}
}
