package app_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
)

// mustBcrypt tạo một bcrypt hash cho test (mô phỏng dữ liệu legacy GoTrue).
func mustBcrypt(t *testing.T, pw string) string {
	t.Helper()
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("gen bcrypt: %v", err)
	}
	return string(b)
}

func TestHasher_Argon2idRoundTrip(t *testing.T) {
	h := app.NewPasswordHasher()
	hash, algo, err := h.Hash("s3cret-pass")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if algo != domain.HashArgon2id {
		t.Fatalf("algo = %q, want argon2id", algo)
	}

	ok, needsRehash, err := h.Verify("s3cret-pass", hash, domain.HashArgon2id)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !ok {
		t.Fatal("đúng mật khẩu nhưng Verify=false")
	}
	if needsRehash {
		t.Fatal("argon2id mới không nên cần rehash")
	}
}

func TestHasher_WrongPassword(t *testing.T) {
	h := app.NewPasswordHasher()
	hash, _, _ := h.Hash("đúng")

	ok, _, err := h.Verify("sai", hash, domain.HashArgon2id)
	if err != nil {
		t.Fatalf("Verify lỗi không mong muốn: %v", err)
	}
	if ok {
		t.Fatal("mật khẩu sai nhưng Verify=true")
	}
}

func TestHasher_BcryptVerifyNeedsRehash(t *testing.T) {
	h := app.NewPasswordHasher()
	// Tạo một bcrypt hash mô phỏng dữ liệu legacy GoTrue.
	bh, err := bcrypt.GenerateFromPassword([]byte("legacy-pw"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("gen bcrypt: %v", err)
	}

	ok, needsRehash, err := h.Verify("legacy-pw", string(bh), domain.HashBcrypt)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !ok {
		t.Fatal("đúng mật khẩu bcrypt nhưng Verify=false")
	}
	if !needsRehash {
		t.Fatal("bcrypt verify đúng phải báo needsRehash=true")
	}

	// Sai mật khẩu với bcrypt → false, không lỗi.
	ok, _, err = h.Verify("sai", string(bh), domain.HashBcrypt)
	if err != nil {
		t.Fatalf("Verify (sai) lỗi: %v", err)
	}
	if ok {
		t.Fatal("bcrypt sai mật khẩu nhưng Verify=true")
	}
}
