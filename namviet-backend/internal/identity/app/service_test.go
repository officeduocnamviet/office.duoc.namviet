package app_test

import (
	"context"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
)

// ---- Fakes (in-memory) cho port domain + TxManager. Fakes > mocks (§11). ----

type fakeUsers struct {
	byEmail map[string]domain.User
	byID    map[string]domain.User
}

func newFakeUsers(users ...domain.User) *fakeUsers {
	f := &fakeUsers{byEmail: map[string]domain.User{}, byID: map[string]domain.User{}}
	for _, u := range users {
		f.byEmail[u.Email] = u
		f.byID[u.ID] = u
	}
	return f
}

func (f *fakeUsers) GetByEmail(_ context.Context, email string) (domain.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return domain.User{}, apperr.NotFound("user")
	}
	return u, nil
}

func (f *fakeUsers) GetByID(_ context.Context, idv string) (domain.User, error) {
	u, ok := f.byID[idv]
	if !ok {
		return domain.User{}, apperr.NotFound("user")
	}
	return u, nil
}

func (f *fakeUsers) UpdatePasswordHash(_ context.Context, userID, hash string, algo domain.HashAlgo) error {
	u := f.byID[userID]
	u.PasswordHash = hash
	u.HashAlgo = algo
	f.byID[userID] = u
	f.byEmail[u.Email] = u
	return nil
}

type fakeRoles struct{ perms map[string][]string }

func (f *fakeRoles) PermissionCodesForUser(_ context.Context, userID string) ([]string, error) {
	return f.perms[userID], nil
}

type fakeTokens struct {
	byHash map[string]domain.RefreshToken
	byID   map[string]domain.RefreshToken
}

func newFakeTokens() *fakeTokens {
	return &fakeTokens{byHash: map[string]domain.RefreshToken{}, byID: map[string]domain.RefreshToken{}}
}

func (f *fakeTokens) Insert(_ context.Context, t domain.RefreshToken) error {
	f.byHash[t.TokenHash] = t
	f.byID[t.ID] = t
	return nil
}

func (f *fakeTokens) GetByHash(_ context.Context, h string) (domain.RefreshToken, error) {
	t, ok := f.byHash[h]
	if !ok {
		return domain.RefreshToken{}, apperr.NotFound("token")
	}
	return t, nil
}

func (f *fakeTokens) MarkUsed(_ context.Context, idv string) error {
	t := f.byID[idv]
	t.Used = true
	f.byID[idv] = t
	f.byHash[t.TokenHash] = t
	return nil
}

func (f *fakeTokens) RevokeFamily(_ context.Context, familyID string) error {
	for k, t := range f.byID {
		if t.FamilyID == familyID {
			t.Revoked = true
			f.byID[k] = t
			f.byHash[t.TokenHash] = t
		}
	}
	return nil
}

// fakeTxM chạy fn ngay với cùng bộ fake repos (không tx thật) — đủ cho unit test
// logic điều phối; tính nguyên tử kiểm ở integration test.
type fakeTxM struct {
	users  domain.UserRepository
	tokens domain.RefreshTokenRepository
	roles  domain.RoleRepository
}

func (m fakeTxM) WithinTx(_ context.Context, fn func(r app.Repos) error) error {
	return fn(app.Repos{Users: m.users, Tokens: m.tokens, Roles: m.roles})
}

// ---- Helpers dựng service ----

func newService(t *testing.T) (*app.AuthService, *fakeUsers, *fakeTokens, domain.User, string) {
	t.Helper()
	hasher := app.NewPasswordHasher()
	hash, algo, err := hasher.Hash("correct-horse")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	user := domain.User{
		ID: "u1", Email: "a@nv.vn", PasswordHash: hash, HashAlgo: algo,
		UserType: "staff", IsActive: true,
	}
	users := newFakeUsers(user)
	tokens := newFakeTokens()
	roles := &fakeRoles{perms: map[string][]string{"u1": {"users:read"}}}
	txm := fakeTxM{users: users, tokens: tokens, roles: roles}
	svc := app.NewAuthService(users, roles, txm, hasher, app.NewTokenIssuer(newKey(t)))
	return svc, users, tokens, user, "correct-horse"
}

// ---- Tests ----

func TestLogin_OK(t *testing.T) {
	svc, _, _, _, pw := newService(t)
	tok, err := svc.Login(context.Background(), "a@nv.vn", pw)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if tok.AccessToken == "" || tok.RefreshToken == "" {
		t.Fatal("thiếu access/refresh token")
	}
	if tok.ExpiresIn != 900 {
		t.Fatalf("ExpiresIn = %d, want 900", tok.ExpiresIn)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	_, err := svc.Login(context.Background(), "a@nv.vn", "sai")
	if apperr.KindOf(err) != apperr.KindUnauthorized {
		t.Fatalf("err kind = %v, want Unauthorized", apperr.KindOf(err))
	}
}

func TestLogin_UnknownEmail_SameError(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	_, err := svc.Login(context.Background(), "khong-ton-tai@nv.vn", "x")
	if apperr.KindOf(err) != apperr.KindUnauthorized {
		t.Fatalf("err kind = %v, want Unauthorized (chống dò tài khoản)", apperr.KindOf(err))
	}
}

func TestLogin_InactiveUserRejected(t *testing.T) {
	svc, users, _, _, pw := newService(t)
	u := users.byID["u1"]
	u.IsActive = false
	users.byID["u1"] = u
	users.byEmail[u.Email] = u

	_, err := svc.Login(context.Background(), "a@nv.vn", pw)
	if apperr.KindOf(err) != apperr.KindUnauthorized {
		t.Fatalf("user inactive phải bị từ chối, kind=%v", apperr.KindOf(err))
	}
}

func TestRefresh_Rotation(t *testing.T) {
	svc, _, _, _, pw := newService(t)
	login, err := svc.Login(context.Background(), "a@nv.vn", pw)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	refreshed, err := svc.Refresh(context.Background(), login.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if refreshed.RefreshToken == login.RefreshToken {
		t.Fatal("refresh token phải xoay (khác token cũ)")
	}
	if refreshed.AccessToken == "" {
		t.Fatal("thiếu access token mới")
	}
}

func TestRefresh_ReuseDetected_RevokesFamily(t *testing.T) {
	svc, _, tokens, _, pw := newService(t)
	login, _ := svc.Login(context.Background(), "a@nv.vn", pw)

	// Lần 1: hợp lệ.
	if _, err := svc.Refresh(context.Background(), login.RefreshToken); err != nil {
		t.Fatalf("refresh lần 1: %v", err)
	}
	// Lần 2: DÙNG LẠI token cũ (đã used) → reuse detected.
	_, err := svc.Refresh(context.Background(), login.RefreshToken)
	if err == nil {
		t.Fatal("reuse phải lỗi")
	}
	ae, ok := apperr.AsError(err)
	if !ok || ae.Code != "refresh_reuse_detected" {
		t.Fatalf("code = %v, want refresh_reuse_detected", err)
	}

	// Toàn bộ family bị thu hồi → mọi token cùng family đều Revoked.
	for _, tk := range tokens.byID {
		if !tk.Revoked {
			t.Fatalf("token %s chưa bị revoke sau reuse", tk.ID)
		}
	}
}

func TestRefresh_InvalidToken(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	_, err := svc.Refresh(context.Background(), "khong-ton-tai")
	ae, ok := apperr.AsError(err)
	if !ok || ae.Kind != apperr.KindUnauthorized {
		t.Fatalf("err = %v, want unauthorized", err)
	}
}

func TestLogout_RevokesFamily_Idempotent(t *testing.T) {
	svc, _, tokens, _, pw := newService(t)
	login, _ := svc.Login(context.Background(), "a@nv.vn", pw)

	if err := svc.Logout(context.Background(), login.RefreshToken); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	for _, tk := range tokens.byID {
		if !tk.Revoked {
			t.Fatal("logout chưa revoke token")
		}
	}
	// Logout lại / token lạ → idempotent, không lỗi.
	if err := svc.Logout(context.Background(), "khong-ton-tai"); err != nil {
		t.Fatalf("logout idempotent phải không lỗi: %v", err)
	}
}

func TestMe(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	me, err := svc.Me(context.Background(), "u1")
	if err != nil {
		t.Fatalf("Me: %v", err)
	}
	if me.Email != "a@nv.vn" || me.UserType != "staff" {
		t.Fatalf("me = %+v", me)
	}
	if len(me.Permissions) != 1 || me.Permissions[0] != "users:read" {
		t.Fatalf("perms = %v", me.Permissions)
	}
}

func TestLogin_BcryptLazyRehash(t *testing.T) {
	// User có hash bcrypt legacy; login đúng → service rehash sang argon2id.
	hasher := app.NewPasswordHasher()
	bcryptHash := mustBcrypt(t, "legacy-pw")
	user := domain.User{
		ID: "u2", Email: "legacy@nv.vn", PasswordHash: bcryptHash, HashAlgo: domain.HashBcrypt,
		UserType: "staff", IsActive: true,
	}
	users := newFakeUsers(user)
	tokens := newFakeTokens()
	roles := &fakeRoles{perms: map[string][]string{}}
	txm := fakeTxM{users: users, tokens: tokens, roles: roles}
	svc := app.NewAuthService(users, roles, txm, hasher, app.NewTokenIssuer(newKey(t)))

	if _, err := svc.Login(context.Background(), "legacy@nv.vn", "legacy-pw"); err != nil {
		t.Fatalf("Login legacy: %v", err)
	}
	got := users.byID["u2"]
	if got.HashAlgo != domain.HashArgon2id {
		t.Fatalf("sau login bcrypt, algo = %q, want argon2id (lazy rehash)", got.HashAlgo)
	}
	if got.PasswordHash == bcryptHash {
		t.Fatal("hash chưa được nâng cấp")
	}
}
