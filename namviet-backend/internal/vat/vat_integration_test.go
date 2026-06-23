package vat_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
	"github.com/Maneva-AI/namviet-backend/internal/vat"
	"github.com/Maneva-AI/namviet-backend/internal/vat/app"
	"github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

const (
	taxCode = "0312345678"
	serial  = "C26TYY"
)

type fixture struct {
	mod  *vat.Module
	pool *pgxpool.Pool
	key  *ecdsa.PrivateKey
}

func setup(t *testing.T) fixture {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return fixture{mod: vat.NewModule(pool), pool: pool, key: key}
}

func mny(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return m
}

func line(t *testing.T, qty, price, rate string) domain.LineInput {
	t.Helper()
	return domain.LineInput{
		ProductID: 1, Description: "Thuốc",
		Quantity: mny(t, qty), UnitPrice: mny(t, price), VATRate: decimal.RequireFromString(rate),
	}
}

// issueInTx chạy IssuePort.IssueInvoice trong MỘT tx (mô phỏng orders gộp atomic
// với giao hàng) và commit. Trả HĐ + lỗi (gồm lỗi commit).
func issueInTx(t *testing.T, fx fixture, p app.IssueParams) (domain.IssuedInvoice, error) {
	t.Helper()
	ctx := context.Background()
	tx, err := fx.pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	inv, ierr := fx.mod.IssuePort().IssueInvoice(ctx, tx, p)
	if ierr != nil {
		_ = tx.Rollback(ctx)
		return domain.IssuedInvoice{}, ierr
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.IssuedInvoice{}, err
	}
	return inv, nil
}

func params(order string, lines ...domain.LineInput) app.IssueParams {
	return app.IssueParams{
		OrderCode: order, CustomerTaxCode: taxCode, Serial: serial, MauSo: "1",
		IssueDate: time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC), Lines: lines,
	}
}

// ---- Phát hành (qua IssuePort, tx caller) ----

func TestIntegration_Issue_Happy(t *testing.T) {
	fx := setup(t)
	p := params("HD123", line(t, "10", "10000", "0.08"), line(t, "3", "50000", "0.10"))
	inv, err := issueInTx(t, fx, p)
	if err != nil {
		t.Fatalf("phát hành happy phải thành công: %v", err)
	}
	if inv.ID == "" || inv.Status != domain.StatusIssued {
		t.Fatalf("HĐ sai: id=%q status=%q", inv.ID, inv.Status)
	}
	if inv.InvoiceNo != 1 {
		t.Fatalf("số HĐ đầu tiên phải = 1, got %d", inv.InvoiceNo)
	}
	// subtotal=100000+150000=250000; vat=8000+15000=23000; total=273000.
	if inv.Subtotal.String() != "250000" || inv.VATAmount.String() != "23000" || inv.Total.String() != "273000" {
		t.Fatalf("tổng sai: sub=%s vat=%s total=%s", inv.Subtotal, inv.VATAmount, inv.Total)
	}
	if len(inv.Lines) != 2 || inv.Lines[0].LineNo != 1 || inv.Lines[1].LineNo != 2 {
		t.Fatalf("lines sai: %+v", inv.Lines)
	}
	// total = subtotal + vat (ép cân).
	if !inv.Total.Equal(inv.Subtotal.Add(inv.VATAmount)) {
		t.Fatal("total phải = subtotal + vat")
	}
}

// Cấp số GAPLESS: phát hành nhiều HĐ tuần tự cùng serial → số 1,2,3 liên tục.
func TestIntegration_Issue_GaplessSequential(t *testing.T) {
	fx := setup(t)
	for i := 1; i <= 5; i++ {
		inv, err := issueInTx(t, fx, params("HD"+itoa(i), line(t, "1", "100000", "0.08")))
		if err != nil {
			t.Fatalf("phát hành i=%d: %v", i, err)
		}
		if inv.InvoiceNo != int64(i) {
			t.Fatalf("số HĐ i=%d phải = %d (gapless), got %d", i, i, inv.InvoiceNo)
		}
	}
}

// CA ĐUA cấp số cùng serial: N tx ĐỒNG THỜI phát hành (mỗi tx 1 đơn khác nhau) →
// số HĐ LIÊN TỤC 1..N, KHÔNG trùng/KHÔNG nhảy (FOR UPDATE tuần tự hoá theo serial).
func TestIntegration_Issue_RaceGapless(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	const n = 8
	var wg sync.WaitGroup
	nos := make([]int64, n)
	errs := make([]error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			inv, err := issueInTx(t, fx, params("RACE"+itoa(idx), line(t, "1", "100000", "0.08")))
			nos[idx] = inv.InvoiceNo
			errs[idx] = err
		}(i)
	}
	wg.Wait()

	for i, e := range errs {
		if e != nil {
			t.Fatalf("race tx %d phải thành công: %v", i, e)
		}
	}
	// Tập số HĐ phải đúng {1,2,...,n} — không trùng, không nhảy.
	sort.Slice(nos, func(a, b int) bool { return nos[a] < nos[b] })
	for i := 0; i < n; i++ {
		if nos[i] != int64(i+1) {
			t.Fatalf("ca đua cấp số: dãy không liên tục, vị trí %d = %d (want %d); full=%v", i, nos[i], i+1, nos)
		}
	}
	// next_no cuối = n+1 (đúng số đã tiêu thụ, không hụt/không dư).
	var nextNo int64
	if err := fx.pool.QueryRow(ctx, `SELECT next_no FROM app.invoice_serials WHERE serial=$1`, serial).Scan(&nextNo); err != nil {
		t.Fatalf("đọc next_no: %v", err)
	}
	if nextNo != int64(n+1) {
		t.Fatalf("next_no = %d, want %d (gapless)", nextNo, n+1)
	}
}

// Idempotency 1 đơn 1 HĐ: phát hành 2 lần CÙNG order_code → CHỈ 1 HĐ, lần 2 trả
// HĐ cũ (cùng id, cùng số), KHÔNG tiêu thêm số HĐ.
func TestIntegration_Issue_IdempotentSameOrder(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	p := params("HD777", line(t, "2", "100000", "0.08"))
	first, err := issueInTx(t, fx, p)
	if err != nil {
		t.Fatalf("lần 1: %v", err)
	}
	second, err := issueInTx(t, fx, p)
	if err != nil {
		t.Fatalf("lần 2 (trùng đơn) phải no-op thành công: %v", err)
	}
	if first.ID != second.ID || first.InvoiceNo != second.InvoiceNo {
		t.Fatalf("idempotent: HĐ khác nhau %s/#%d != %s/#%d", first.ID, first.InvoiceNo, second.ID, second.InvoiceNo)
	}
	var cnt int
	if err := fx.pool.QueryRow(ctx, `SELECT count(*) FROM app.sales_invoices WHERE order_code=$1`, "HD777").Scan(&cnt); err != nil {
		t.Fatalf("count: %v", err)
	}
	if cnt != 1 {
		t.Fatalf("số HĐ cho HD777 = %d, want 1 (không phát hành trùng)", cnt)
	}
	// next_no chỉ tiêu 1 số (=2): lần 2 dừng ở idempotency trước khi cấp số.
	var nextNo int64
	_ = fx.pool.QueryRow(ctx, `SELECT next_no FROM app.invoice_serials WHERE serial=$1`, serial).Scan(&nextNo)
	if nextNo != 2 {
		t.Fatalf("next_no = %d, want 2 (lần 2 không tiêu số)", nextNo)
	}
}

func TestIntegration_Issue_MissingTaxCode_422(t *testing.T) {
	fx := setup(t)
	p := params("HD1", line(t, "1", "100000", "0.08"))
	p.CustomerTaxCode = "   "
	_, err := issueInTx(t, fx, p)
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("thiếu MST phải Validation(422), got %v", err)
	}
}

func TestIntegration_Issue_NoLines_422(t *testing.T) {
	fx := setup(t)
	_, err := issueInTx(t, fx, params("HD1"))
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("không dòng phải Validation(422), got %v", err)
	}
}

// Ép cân tổng ở DB: chèn thẳng header có total != subtotal+vat → CHECK chặn.
func TestIntegration_DB_TotalMustBalance(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	_, err := fx.pool.Exec(ctx, `
		INSERT INTO app.sales_invoices
		  (order_code, customer_tax_code, serial, invoice_no, issue_date, subtotal, vat_amount, total)
		VALUES ('HDX', '0312345678', 'C26TYY', 1, '2026-06-20', 100000, 8000, 999999)`)
	if err == nil {
		t.Fatal("total lệch phải bị CHECK chặn ở DB")
	}
}

// ---- HTTP đọc (envelope + authz) ----

func TestIntegration_HTTP_VAT_ReadAndAuthz(t *testing.T) {
	fx := setup(t)
	// Phát hành vài HĐ để list.
	var firstID string
	for i := 1; i <= 3; i++ {
		inv, err := issueInTx(t, fx, params("HD"+itoa(i), line(t, "1", "100000", "0.08")))
		if err != nil {
			t.Fatalf("seed HĐ %d: %v", i, err)
		}
		if i == 1 {
			firstID = inv.ID
		}
	}
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "vat.read")

	// Thiếu token → 401.
	noTok := doReq(t, handler, "/v1/vat/invoices", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Sai quyền → 403.
	wrong := doReq(t, handler, "/v1/vat/invoices", signToken(t, fx.key, "orders.read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Có quyền → 200 + envelope, lọc status=issued.
	ok := doReq(t, handler, "/v1/vat/invoices?status=issued", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("invoices = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Items []struct {
				ID        string `json:"id"`
				Total     string `json:"total"`
				InvoiceNo int64  `json:"invoice_no"`
			} `json:"items"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, ok.body)
	}
	if env.Error != nil || len(env.Data.Items) != 3 {
		t.Fatalf("envelope sai (want 3 items): %s", ok.body)
	}
	// total là chuỗi (không float JSON).
	if env.Data.Items[0].Total == "" {
		t.Fatalf("total phải là chuỗi: %s", ok.body)
	}

	// GetInvoice kèm lines.
	g := doReq(t, handler, "/v1/vat/invoices/"+firstID, tok)
	if g.code != http.StatusOK {
		t.Fatalf("get invoice = %d; body=%s", g.code, g.body)
	}
	var genv struct {
		Data struct {
			ID    string `json:"id"`
			Lines []struct {
				LineVAT string `json:"line_vat"`
				VATRate string `json:"vat_rate"`
			} `json:"lines"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(g.body), &genv); err != nil {
		t.Fatalf("decode get: %v (body=%s)", err, g.body)
	}
	if len(genv.Data.Lines) != 1 || genv.Data.Lines[0].LineVAT == "" || genv.Data.Lines[0].VATRate == "" {
		t.Fatalf("lines sai: %s", g.body)
	}

	// id không tồn tại → 404.
	nf := doReq(t, handler, "/v1/vat/invoices/00000000-0000-0000-0000-000000000000", tok)
	if nf.code != http.StatusNotFound {
		t.Fatalf("not found = %d, want 404; body=%s", nf.code, nf.body)
	}
}

// ---- helpers ----

func itoa(i int) string { return string(rune('0'+i/10)) + string(rune('0'+i%10)) }

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	fx.mod.RegisterRoutes(api, verifier)
	return r
}

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

var _ = pgx.ErrNoRows // giữ import pgx nếu helper tx không tham chiếu trực tiếp
