package identity_test

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

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/identity"
	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// fixture gom mọi thứ một test cần.
type fixture struct {
	svc    *identity.Service
	pool   *pgxpool.Pool
	key    *ecdsa.PrivateKey
	userID string
	pw     string
}

// setup spin DB, seed user (argon2id) + role "admin" + permission "users:read".
func setup(t *testing.T) fixture {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	ctx := context.Background()

	q := appdb.New(pool)
	hash, algo, err := app.NewPasswordHasher().Hash("correct-horse")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	userID, err := q.InsertUser(ctx, appdb.InsertUserParams{
		Email: "staff@nv.vn", PasswordHash: hash, HashAlgo: string(algo),
		UserType: "staff", IsActive: true,
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	roleID, err := q.InsertRole(ctx, appdb.InsertRoleParams{Code: "admin", Name: "Quản trị"})
	if err != nil {
		t.Fatalf("seed role: %v", err)
	}
	permID, err := q.InsertPermission(ctx, appdb.InsertPermissionParams{Code: "users:read"})
	if err != nil {
		t.Fatalf("seed perm: %v", err)
	}
	if err := q.AssignPermissionToRole(ctx, appdb.AssignPermissionToRoleParams{RoleID: roleID, PermissionID: permID}); err != nil {
		t.Fatalf("assign perm: %v", err)
	}
	if err := q.AssignRoleToUser(ctx, appdb.AssignRoleToUserParams{UserID: userID, RoleID: roleID}); err != nil {
		t.Fatalf("assign role: %v", err)
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return fixture{
		svc:    identity.New(pool, identity.NewTokenIssuer(key)),
		pool:   pool,
		key:    key,
		userID: userID,
		pw:     "correct-horse",
	}
}

func TestIntegration_Login_OK(t *testing.T) {
	fx := setup(t)
	tok, err := fx.svc.Login(context.Background(), "staff@nv.vn", fx.pw)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if tok.AccessToken == "" || tok.RefreshToken == "" {
		t.Fatal("thiếu token")
	}
}

func TestIntegration_Login_WrongPassword(t *testing.T) {
	fx := setup(t)
	_, err := fx.svc.Login(context.Background(), "staff@nv.vn", "sai")
	if apperr.KindOf(err) != apperr.KindUnauthorized {
		t.Fatalf("kind = %v, want Unauthorized", apperr.KindOf(err))
	}
}

func TestIntegration_RefreshRotation(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()

	login, err := fx.svc.Login(ctx, "staff@nv.vn", fx.pw)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	refreshed, err := fx.svc.Refresh(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if refreshed.RefreshToken == login.RefreshToken {
		t.Fatal("refresh token phải xoay")
	}
	// Token cũ đã used; tổng 2 token cùng family.
	if used := countTokens(t, fx.pool, "used = true"); used != 1 {
		t.Fatalf("used = %d, want 1", used)
	}
	if total := countTokens(t, fx.pool, "true"); total != 2 {
		t.Fatalf("tổng = %d, want 2", total)
	}
	// Token mới refresh được tiếp (chuỗi hợp lệ).
	if _, err := fx.svc.Refresh(ctx, refreshed.RefreshToken); err != nil {
		t.Fatalf("refresh token mới phải hợp lệ: %v", err)
	}
}

func TestIntegration_ReuseDetected_RevokesFamily(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()

	login, _ := fx.svc.Login(ctx, "staff@nv.vn", fx.pw)
	if _, err := fx.svc.Refresh(ctx, login.RefreshToken); err != nil {
		t.Fatalf("refresh 1: %v", err)
	}
	// Dùng LẠI token cũ (đã used) → reuse detected.
	_, err := fx.svc.Refresh(ctx, login.RefreshToken)
	ae, ok := apperr.AsError(err)
	if !ok || ae.Code != "refresh_reuse_detected" {
		t.Fatalf("err = %v, want refresh_reuse_detected", err)
	}
	// Cả family bị revoke.
	if notRevoked := countTokens(t, fx.pool, "revoked = false"); notRevoked != 0 {
		t.Fatalf("còn %d token chưa revoke", notRevoked)
	}
}

func TestIntegration_Logout_RevokesFamily(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()

	login, _ := fx.svc.Login(ctx, "staff@nv.vn", fx.pw)
	if err := fx.svc.Logout(ctx, login.RefreshToken); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if notRevoked := countTokens(t, fx.pool, "revoked = false"); notRevoked != 0 {
		t.Fatalf("logout chưa revoke hết, còn %d", notRevoked)
	}
	if _, err := fx.svc.Refresh(ctx, login.RefreshToken); err == nil {
		t.Fatal("refresh sau logout phải lỗi")
	}
}

func TestIntegration_Me(t *testing.T) {
	fx := setup(t)
	me, err := fx.svc.Me(context.Background(), fx.userID)
	if err != nil {
		t.Fatalf("Me: %v", err)
	}
	if me.Email != "staff@nv.vn" || me.UserType != "staff" {
		t.Fatalf("me = %+v", me)
	}
	if len(me.Permissions) != 1 || me.Permissions[0] != "users:read" {
		t.Fatalf("perms = %v", me.Permissions)
	}
}

// TestIntegration_HTTP_LoginThenMe kiểm end-to-end qua HTTP: login lấy token,
// gọi /v1/auth/me với Bearer → 200 + đúng quyền; thiếu token → 401 envelope.
func TestIntegration_HTTP_LoginThenMe(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)

	// Login qua HTTP.
	loginResp := doJSON(t, handler, http.MethodPost, "/v1/auth/login", "",
		`{"email":"staff@nv.vn","password":"correct-horse"}`)
	if loginResp.code != http.StatusOK {
		t.Fatalf("login status = %d, body=%s", loginResp.code, loginResp.body)
	}
	var loginData struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(loginResp.body), &loginData); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if loginData.Data.AccessToken == "" {
		t.Fatalf("login không trả access token: %s", loginResp.body)
	}

	// /me KHÔNG token → 401.
	noTok := doJSON(t, handler, http.MethodGet, "/v1/auth/me", "", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("me không token status = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// /me CÓ token → 200 + perms.
	withTok := doJSON(t, handler, http.MethodGet, "/v1/auth/me", loginData.Data.AccessToken, "")
	if withTok.code != http.StatusOK {
		t.Fatalf("me có token status = %d, want 200; body=%s", withTok.code, withTok.body)
	}
	if !strings.Contains(withTok.body, "users:read") || !strings.Contains(withTok.body, "staff@nv.vn") {
		t.Fatalf("me body thiếu dữ liệu mong đợi: %s", withTok.body)
	}
}

// ---- HTTP helpers ----

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	identity.RegisterRoutes(api, fx.svc, verifier)
	return r
}

type resp struct {
	code int
	body string
}

func doJSON(t *testing.T, h http.Handler, method, path, bearer, body string) resp {
	t.Helper()
	var rdr *strings.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	} else {
		rdr = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return resp{code: rec.Code, body: rec.Body.String()}
}

func assertErrCode(t *testing.T, body, want string) {
	t.Helper()
	var env struct {
		Error *struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("decode envelope: %v (body=%s)", err, body)
	}
	if env.Error == nil || env.Error.Code != want {
		t.Fatalf("error.code = %v, want %q (body=%s)", env.Error, want, body)
	}
}

func countTokens(t *testing.T, pool *pgxpool.Pool, where string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		"SELECT count(*) FROM app.refresh_tokens WHERE "+where).Scan(&n); err != nil {
		t.Fatalf("count tokens: %v", err)
	}
	return n
}
