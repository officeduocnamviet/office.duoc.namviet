package authn_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
)

// newKey sinh cặp khoá EC P-256 cho test.
func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return k
}

// signES256 ký token ES256 với claims cho trước.
func signES256(t *testing.T, key *ecdsa.PrivateKey, c authn.Claims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, c)
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

func validClaims(perms ...string) authn.Claims {
	return authn.Claims{
		UserType:    "staff",
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func TestVerify_OK(t *testing.T) {
	key := newKey(t)
	v := authn.NewVerifier(&key.PublicKey)

	claims, err := v.Verify(signES256(t, key, validClaims("orders:read")))
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.UserID() != "user-123" {
		t.Fatalf("sub = %q", claims.UserID())
	}
	if !claims.HasPermission("orders:read") {
		t.Fatal("mất perm")
	}
}

func TestVerify_Expired(t *testing.T) {
	key := newKey(t)
	v := authn.NewVerifier(&key.PublicKey)

	c := validClaims()
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute)) // đã hết hạn
	if _, err := v.Verify(signES256(t, key, c)); err == nil {
		t.Fatal("token hết hạn phải bị từ chối")
	}
}

// TestVerify_RejectsHS256 chứng minh PIN alg: token HS256 (đối xứng) bị từ chối
// dù attacker thử dùng public key bytes làm secret.
func TestVerify_RejectsHS256(t *testing.T) {
	key := newKey(t)
	v := authn.NewVerifier(&key.PublicKey)

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims())
	hs, err := tok.SignedString([]byte("attacker-secret"))
	if err != nil {
		t.Fatalf("sign hs256: %v", err)
	}
	if _, err := v.Verify(hs); err == nil {
		t.Fatal("token HS256 phải bị từ chối (alg-confusion)")
	}
}

// TestVerify_RejectsNone chứng minh token alg=none bị từ chối.
func TestVerify_RejectsNone(t *testing.T) {
	key := newKey(t)
	v := authn.NewVerifier(&key.PublicKey)

	tok := jwt.NewWithClaims(jwt.SigningMethodNone, validClaims())
	none, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}
	if _, err := v.Verify(none); err == nil {
		t.Fatal("token alg=none phải bị từ chối")
	}
}

func TestMiddleware_MissingToken_401(t *testing.T) {
	key := newKey(t)
	mw := authn.Middleware(authn.NewVerifier(&key.PublicKey))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil) // không header
	mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("không được tới handler khi thiếu token")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	assertEnvelopeErrorCode(t, rec.Body.Bytes(), "unauthorized")
}

func TestMiddleware_ValidToken_PassesClaims(t *testing.T) {
	key := newKey(t)
	mw := authn.Middleware(authn.NewVerifier(&key.PublicKey))

	var gotSub string
	h := mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		c, ok := authn.ClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("middleware không nạp claims")
		}
		gotSub = c.UserID()
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+signES256(t, key, validClaims()))
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if gotSub != "user-123" {
		t.Fatalf("sub = %q", gotSub)
	}
}

func assertEnvelopeErrorCode(t *testing.T, body []byte, wantCode string) {
	t.Helper()
	var env struct {
		Data  any `json:"data"`
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v (body=%s)", err, body)
	}
	if env.Error == nil {
		t.Fatalf("envelope.error nil, body=%s", body)
	}
	if env.Error.Code != wantCode {
		t.Fatalf("error.code = %q, want %q", env.Error.Code, wantCode)
	}
}
