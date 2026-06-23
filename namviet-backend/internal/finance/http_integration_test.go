package finance_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Maneva-AI/namviet-backend/internal/finance/app"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// TestIntegration_HTTP_Finance_EnvelopeAndAuthz: route ĐỌC /v1/finance/transactions
// đúng envelope {data,error} + ép quyền finance.read (401 thiếu token, 403 sai
// quyền, 200 + items khi đủ quyền), 422 thiếu ref_id (required).
func TestIntegration_HTTP_Finance_EnvelopeAndAuthz(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()

	// Seed 1 phiếu cho HD100 để route trả 1 item.
	p := app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode: "HD100", Amount: mustMoney(t, "777"), FundAccountID: 1, BookType: domain.BookBoth,
		},
		IdemKey: "http-1",
	}
	if _, err := recordInTx(ctx, pool, mod.Recorder(), p); err != nil {
		t.Fatalf("seed phiếu HTTP: %v", err)
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	mod.RegisterRoutes(api, authn.NewVerifier(&key.PublicKey))

	// Thiếu token → 401.
	noTok := doReq(t, r, "/v1/finance/transactions?ref_id=HD100", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Sai quyền → 403.
	wrong := doReq(t, r, "/v1/finance/transactions?ref_id=HD100", signToken(t, key, "accounting.read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Đủ quyền → 200 + envelope {data:{items:[1]}, error:null}, amount chuỗi decimal.
	ok := doReq(t, r, "/v1/finance/transactions?ref_id=HD100", signToken(t, key, "finance.read"))
	if ok.code != http.StatusOK {
		t.Fatalf("transactions = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Items []struct {
				RefID  string `json:"ref_id"`
				Amount string `json:"amount"`
				Flow   string `json:"flow"`
				Status string `json:"status"`
			} `json:"items"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, ok.body)
	}
	if env.Error != nil || len(env.Data.Items) != 1 {
		t.Fatalf("envelope sai: %s", ok.body)
	}
	it := env.Data.Items[0]
	if it.RefID != "HD100" || it.Amount != "777" || it.Flow != "in" || it.Status != "completed" {
		t.Fatalf("item sai: %+v", it)
	}

	// Thiếu ref_id (required) → 422 validation.
	missing := doReq(t, r, "/v1/finance/transactions", signToken(t, key, "finance.read"))
	if missing.code != http.StatusUnprocessableEntity {
		t.Fatalf("thiếu ref_id = %d, want 422; body=%s", missing.code, missing.body)
	}
}

// ---- helpers HTTP ----

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
	s, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

type httpResp struct {
	code int
	body string
}

func doReq(t *testing.T, h http.Handler, path, bearer string) httpResp {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return httpResp{code: rec.Code, body: rec.Body.String()}
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

// TestIntegration_HTTP_Finance_ConfirmReceipt: route GHI POST /v1/finance/receipts/
// {id}/confirm (thủ quỹ xác nhận thu — 2 bước). Guard finance.write (401/403), 200 +
// confirmed=true cho phiếu pending; gọi lại idempotent → confirmed=false.
func TestIntegration_HTTP_Finance_ConfirmReceipt(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()

	// Seed phiếu PENDING (NV thu, chưa nộp quỹ).
	pay, err := recordInTx(ctx, pool, mod.Recorder(), app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode: "HD200", Amount: mustMoney(t, "5000"), FundAccountID: 1,
			BookType: domain.BookInternal, InitialStatus: domain.StatusPending,
		},
		IdemKey: "confirm-http-1",
	})
	if err != nil {
		t.Fatalf("seed pending: %v", err)
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	mod.RegisterRoutes(api, authn.NewVerifier(&key.PublicKey))
	path := fmt.Sprintf("/v1/finance/receipts/%d/confirm", pay.ID)

	if x := doPostFin(t, r, path, ""); x.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; %s", x.code, x.body)
	}
	if x := doPostFin(t, r, path, signToken(t, key, "finance.read")); x.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; %s", x.code, x.body)
	}
	ok := doPostFin(t, r, path, signToken(t, key, "finance.write"))
	if ok.code != http.StatusOK {
		t.Fatalf("confirm = %d, want 200; %s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Confirmed bool `json:"confirmed"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil || env.Error != nil || !env.Data.Confirmed {
		t.Fatalf("confirm envelope want confirmed=true: %s", ok.body)
	}
	// Idempotent: xác nhận lại → confirmed=false (đã completed, không cộng đôi số dư).
	again := doPostFin(t, r, path, signToken(t, key, "finance.write"))
	if err := json.Unmarshal([]byte(again.body), &env); err != nil || env.Data.Confirmed {
		t.Fatalf("re-confirm phải confirmed=false (idempotent): %s", again.body)
	}
}

// doPostFin gửi POST (không body) — cho route confirm receipt.
func doPostFin(t *testing.T, h http.Handler, path, bearer string) httpResp {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return httpResp{code: rec.Code, body: rec.Body.String()}
}
