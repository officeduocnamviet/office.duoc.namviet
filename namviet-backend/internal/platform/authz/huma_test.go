package authz_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// signToken ký một access JWT ES256 với perms cho trước (mô phỏng identity issuer).
func signToken(t *testing.T, key *ecdsa.PrivateKey, perms ...string) string {
	t.Helper()
	claims := authn.Claims{
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "u1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

// buildProtected dựng một huma.API với 1 route /v1/cat được bảo vệ bởi
// RequirePermissionHuma(perm).
func buildProtected(t *testing.T, perm string, key *ecdsa.PrivateKey) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&key.PublicKey)
	huma.Register(api, huma.Operation{
		OperationID: "cat-protected",
		Method:      http.MethodGet,
		Path:        "/v1/cat",
		Middlewares: huma.Middlewares{authz.RequirePermissionHuma(api, verifier, perm)},
	}, func(_ context.Context, _ *struct{}) (*struct{ Body struct{ OK bool } }, error) {
		out := &struct{ Body struct{ OK bool } }{}
		out.Body.OK = true
		return out, nil
	})
	return r
}

func do(t *testing.T, h http.Handler, bearer string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/cat", nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRequirePermissionHuma_Allow(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	h := buildProtected(t, "catalog.read", key)
	rec := do(t, h, signToken(t, key, "catalog.read"))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestRequirePermissionHuma_NoToken_401(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	h := buildProtected(t, "catalog.read", key)
	rec := do(t, h, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body=%s", rec.Code, rec.Body.String())
	}
	assertHumaCode(t, rec.Body.String(), "unauthorized")
}

func TestRequirePermissionHuma_MissingPerm_403(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	h := buildProtected(t, "catalog.read", key)
	rec := do(t, h, signToken(t, key, "orders:read"))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body=%s", rec.Code, rec.Body.String())
	}
	assertHumaCode(t, rec.Body.String(), "forbidden")
}

func assertHumaCode(t *testing.T, body, want string) {
	t.Helper()
	var env struct {
		Error *struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("unmarshal: %v (body=%s)", err, body)
	}
	if env.Error == nil || !strings.EqualFold(env.Error.Code, want) {
		t.Fatalf("code = %v, want %q (body=%s)", env.Error, want, body)
	}
}
