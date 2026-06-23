package authz_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
)

// withClaims gắn claims vào request context, mô phỏng kết quả của authn.
func withClaims(req *http.Request, perms ...string) *http.Request {
	c := &authn.Claims{
		Permissions:      perms,
		RegisteredClaims: jwt.RegisteredClaims{Subject: "u1"},
	}
	return req.WithContext(authn.WithClaims(req.Context(), c))
}

func serve(t *testing.T, perm string, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	reached := false
	h := authz.RequirePermission(perm)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(rec, req)
	if rec.Code == http.StatusOK && !reached {
		t.Fatal("status OK nhưng handler không chạy")
	}
	return rec
}

func TestRequirePermission_Allow(t *testing.T) {
	req := withClaims(httptest.NewRequest(http.MethodGet, "/x", nil), "users:read", "orders:read")
	rec := serve(t, "users:read", req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestRequirePermission_Deny_403(t *testing.T) {
	req := withClaims(httptest.NewRequest(http.MethodGet, "/x", nil), "orders:read")
	rec := serve(t, "users:write", req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
	assertCode(t, rec.Body.Bytes(), "forbidden")
}

func TestRequirePermission_NoClaims_401(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil) // không có claims
	rec := serve(t, "users:read", req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	assertCode(t, rec.Body.Bytes(), "unauthorized")
}

func assertCode(t *testing.T, body []byte, want string) {
	t.Helper()
	var env struct {
		Error *struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Error == nil || env.Error.Code != want {
		t.Fatalf("code = %v, want %q (body=%s)", env.Error, want, body)
	}
}
