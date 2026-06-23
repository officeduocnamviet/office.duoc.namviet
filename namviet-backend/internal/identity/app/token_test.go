package app_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
)

func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return k
}

func TestTokenIssuer_IssueThenVerify(t *testing.T) {
	key := newKey(t)
	issuer := app.NewTokenIssuer(key)
	verifier := authn.NewVerifier(&key.PublicKey)

	raw, err := issuer.IssueAccess("user-42", "staff", []string{"users:read", "orders:read"})
	if err != nil {
		t.Fatalf("IssueAccess: %v", err)
	}

	claims, err := verifier.Verify(raw)
	if err != nil {
		t.Fatalf("Verify token vừa phát: %v", err)
	}
	if claims.UserID() != "user-42" {
		t.Fatalf("sub = %q", claims.UserID())
	}
	if claims.UserType != "staff" {
		t.Fatalf("user_type = %q", claims.UserType)
	}
	if !claims.HasPermission("users:read") || !claims.HasPermission("orders:read") {
		t.Fatalf("perms = %v", claims.Permissions)
	}
	if claims.ID == "" {
		t.Fatal("jti rỗng")
	}
}

func TestTokenIssuer_WrongKeyRejected(t *testing.T) {
	issuer := app.NewTokenIssuer(newKey(t))
	otherVerifier := authn.NewVerifier(&newKey(t).PublicKey) // khoá khác

	raw, _ := issuer.IssueAccess("u", "staff", nil)
	if _, err := otherVerifier.Verify(raw); err == nil {
		t.Fatal("token ký bằng khoá khác phải bị từ chối")
	}
}

func TestNewRefreshToken_HashStableAndOpaque(t *testing.T) {
	raw, hash, err := app.NewRefreshToken()
	if err != nil {
		t.Fatalf("NewRefreshToken: %v", err)
	}
	if raw == "" || hash == "" {
		t.Fatal("raw/hash rỗng")
	}
	if raw == hash {
		t.Fatal("hash không được trùng raw")
	}
	// Hash phải xác định (deterministic) để tra cứu GetByHash.
	if app.HashRefreshToken(raw) != hash {
		t.Fatal("HashRefreshToken không ổn định")
	}

	// Hai token khác nhau → hash khác nhau.
	raw2, hash2, _ := app.NewRefreshToken()
	if raw == raw2 || hash == hash2 {
		t.Fatal("hai refresh token không được trùng")
	}
}
