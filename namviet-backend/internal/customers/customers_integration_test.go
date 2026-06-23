package customers_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/customers"
	"github.com/Maneva-AI/namviet-backend/internal/customers/app"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

type fixture struct {
	svc  *customers.Service
	pool *pgxpool.Pool
	key  *ecdsa.PrivateKey
}

// seed nạp dữ liệu test vào public.customers + public.orders (id tường minh).
// Khách:
//
//	201 B2B "Cty Alpha" — MST + debt_limit + payment_term trong b2b_metadata;
//	    đơn: 1tr (unpaid) + 2tr (partial) + 5tr (paid, KHÔNG tính nợ) → live = 3tr.
//	    current_debt tĩnh = 999tr (stale, KHÔNG được dùng).
//	202 B2B "Cty Beta" — metadata rỗng; KHÔNG có đơn → live = 0.
//	203 B2C "Chị Lan"  — không B2B; 1 đơn unpaid 500k → live = 500k.
//	204 B2C đã soft-delete (không được trả).
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.customers
			(id, customer_code, name, customer_type, phone, email, status, b2b_metadata, current_debt)
		 VALUES
			(201, 'KH-ALPHA', 'Cty Alpha', 'B2B', '0901111111', 'alpha@x.vn', 'active',
			 '{"tax_code":"0312345678","debt_limit":50000000,"payment_term":30,"sales_staff_id":"staff-uuid"}'::jsonb,
			 999000000),
			(202, 'KH-BETA', 'Cty Beta', 'B2B', '0902222222', 'beta@x.vn', 'active',
			 '{}'::jsonb, 0),
			(203, 'KH-LAN', 'Chị Lan', 'B2C', '0903333333', 'lan@x.vn', 'active',
			 '{}'::jsonb, 0),
			(204, 'KH-DEL', 'Khách đã xóa', 'B2C', '0904444444', 'del@x.vn', 'active',
			 '{}'::jsonb, 0)`,
		`UPDATE public.customers SET deleted_at = now() WHERE id = 204`,
		`INSERT INTO public.orders (code, customer_id, status, order_type, final_amount, payment_status) VALUES
			('HD-A1', 201, 'COMPLETED', 'B2B', 1000000, 'unpaid'),
			('HD-A2', 201, 'COMPLETED', 'B2B', 2000000, 'partial'),
			('HD-A3', 201, 'COMPLETED', 'B2B', 5000000, 'paid'),
			('HD-A4', 201, 'COMPLETED', 'B2B', 8000000, 'unpaid'),
			('HD-L1', 203, 'COMPLETED', 'B2C', 500000, 'unpaid')`,
		// Đơn HD-A4 bị soft-delete → KHÔNG tính nợ. Live Alpha = 1tr + 2tr = 3tr.
		`UPDATE public.orders SET deleted_at = now() WHERE code = 'HD-A4'`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("seed: %v\nSQL: %s", err, s)
		}
	}
}

func setup(t *testing.T) fixture {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	seed(t, pool)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return fixture{svc: customers.New(pool), pool: pool, key: key}
}

func TestIntegration_ListCustomers_SoftDeleteAndOrder(t *testing.T) {
	fx := setup(t)
	res, err := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 50})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if len(res.Items) != 3 { // loại 204 soft-deleted
		t.Fatalf("đếm = %d, want 3 (loại soft-delete)", len(res.Items))
	}
	for _, c := range res.Items {
		if c.ID == 204 {
			t.Fatal("khách soft-deleted 204 không được trả")
		}
	}
	// keyset id ASC: 201, 202, 203.
	if res.Items[0].ID != 201 || res.Items[2].ID != 203 {
		t.Fatalf("thứ tự keyset sai: %+v", idsOf(res.Items))
	}
}

func TestIntegration_ListCustomers_FilterTypeAndQuery(t *testing.T) {
	fx := setup(t)
	// Filter B2B → 201, 202.
	b2b, _ := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 50, Type: "B2B"})
	if len(b2b.Items) != 2 {
		t.Fatalf("B2B count = %d, want 2", len(b2b.Items))
	}
	// Filter B2C → 203.
	b2c, _ := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 50, Type: "B2C"})
	if len(b2c.Items) != 1 || b2c.Items[0].ID != 203 {
		t.Fatalf("B2C sai: %+v", idsOf(b2c.Items))
	}
	// Tìm theo MST (trong b2b_metadata) → chỉ Alpha (201).
	q, _ := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 50, Query: "0312345678"})
	if len(q.Items) != 1 || q.Items[0].ID != 201 {
		t.Fatalf("search MST sai: %+v", idsOf(q.Items))
	}
	// Tìm theo SĐT → Beta (202).
	qp, _ := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 50, Query: "0902222222"})
	if len(qp.Items) != 1 || qp.Items[0].ID != 202 {
		t.Fatalf("search phone sai: %+v", idsOf(qp.Items))
	}
}

func TestIntegration_ListCustomers_Keyset(t *testing.T) {
	fx := setup(t)
	p1, _ := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 2})
	if len(p1.Items) != 2 || p1.Items[0].ID != 201 || p1.Items[1].ID != 202 {
		t.Fatalf("trang 1 sai: %+v", idsOf(p1.Items))
	}
	if p1.NextCursor == "" {
		t.Fatal("trang 1 đầy phải có NextCursor")
	}
	p2, _ := fx.svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 2, Cursor: p1.NextCursor})
	if len(p2.Items) != 1 || p2.Items[0].ID != 203 {
		t.Fatalf("trang 2 sai: %+v", idsOf(p2.Items))
	}
	if p2.NextCursor != "" {
		t.Fatal("trang 2 hết phải rỗng NextCursor")
	}
}

// Nghiệp vụ LÕI: công nợ live tính từ orders (đơn chưa tất toán, loại paid +
// soft-deleted), ưu tiên hơn cột tĩnh stale.
func TestIntegration_GetCustomer_LiveDebtFromOrders(t *testing.T) {
	fx := setup(t)
	c, err := fx.svc.GetCustomer(context.Background(), 201)
	if err != nil {
		t.Fatalf("GetCustomer: %v", err)
	}
	if c.Debt.Source != domain.DebtSourceLive {
		t.Fatalf("source = %q, want live", c.Debt.Source)
	}
	// 1tr (unpaid) + 2tr (partial) = 3tr; KHÔNG gồm 5tr paid và 8tr soft-deleted.
	if c.Debt.Amount.String() != "3000000" {
		t.Fatalf("live debt = %q, want 3000000", c.Debt.Amount.String())
	}
	// Cột tĩnh stale vẫn được giữ để minh bạch, nhưng KHÔNG phải con số chính.
	if c.Debt.Static.String() != "999000000" {
		t.Fatalf("static phải = 999000000 (giữ nguyên), got %q", c.Debt.Static.String())
	}
	if c.Debt.Amount.Equal(c.Debt.Static) {
		t.Fatal("amount KHÔNG được dùng cột tĩnh stale")
	}
}

func TestIntegration_GetCustomer_B2BProfileParsedFromMetadata(t *testing.T) {
	fx := setup(t)
	c, _ := fx.svc.GetCustomer(context.Background(), 201)
	if !c.IsB2B() || c.B2B == nil {
		t.Fatalf("201 phải B2B có profile: %+v", c)
	}
	if c.B2B.TaxCode != "0312345678" {
		t.Fatalf("MST = %q", c.B2B.TaxCode)
	}
	if c.B2B.DebtLimit.String() != "50000000" {
		t.Fatalf("debt_limit = %q, want 50000000", c.B2B.DebtLimit.String())
	}
	if c.B2B.PaymentTerm != 30 {
		t.Fatalf("payment_term = %d, want 30", c.B2B.PaymentTerm)
	}
}

func TestIntegration_GetCustomer_B2CNoProfileAndOwnDebt(t *testing.T) {
	fx := setup(t)
	c, _ := fx.svc.GetCustomer(context.Background(), 203)
	if c.IsB2B() || c.B2B != nil {
		t.Fatalf("203 là B2C, không có B2B profile: %+v", c)
	}
	if c.Debt.Amount.String() != "500000" {
		t.Fatalf("B2C live debt = %q, want 500000", c.Debt.Amount.String())
	}
}

func TestIntegration_GetCustomer_NoOrders_ZeroLiveDebt(t *testing.T) {
	fx := setup(t)
	c, _ := fx.svc.GetCustomer(context.Background(), 202)
	// Không có đơn → live = 0 (COALESCE), vẫn dùng nguồn live (0 là sự thật).
	if c.Debt.Source != domain.DebtSourceLive {
		t.Fatalf("source = %q, want live", c.Debt.Source)
	}
	if !c.Debt.Amount.IsZero() {
		t.Fatalf("không đơn phải nợ 0, got %q", c.Debt.Amount.String())
	}
}

func TestIntegration_GetCustomer_NotFound(t *testing.T) {
	fx := setup(t)
	_, err := fx.svc.GetCustomer(context.Background(), 9999)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("kind = %v, want NotFound", apperr.KindOf(err))
	}
	// Soft-deleted cũng là NotFound.
	_, err = fx.svc.GetCustomer(context.Background(), 204)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("soft-deleted kind = %v, want NotFound", apperr.KindOf(err))
	}
}

// ---- HTTP end-to-end (envelope + authz) ----

func TestIntegration_HTTP_ListCustomers_Envelope(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "customers.read")

	// Thiếu token → 401 envelope.
	noTok := doReq(t, handler, "/v1/customers", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Sai quyền → 403.
	wrong := doReq(t, handler, "/v1/customers", signToken(t, fx.key, "catalog.read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Có quyền → 200 + envelope {data:{items,next_cursor}, error:null}.
	ok := doReq(t, handler, "/v1/customers?limit=2", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("ok = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Items []struct {
				ID   int64 `json:"id"`
				Debt struct {
					Amount string `json:"amount"`
					Source string `json:"source"`
				} `json:"debt"`
			} `json:"items"`
			NextCursor string `json:"next_cursor"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil {
		t.Fatalf("decode envelope: %v (body=%s)", err, ok.body)
	}
	if env.Error != nil {
		t.Fatalf("error phải null: %s", ok.body)
	}
	if len(env.Data.Items) != 2 {
		t.Fatalf("items = %d, want 2; body=%s", len(env.Data.Items), ok.body)
	}
	// Item đầu (201) debt live = 3tr dạng chuỗi decimal.
	if env.Data.Items[0].Debt.Amount != "3000000" || env.Data.Items[0].Debt.Source != "live" {
		t.Fatalf("debt envelope sai: %+v", env.Data.Items[0].Debt)
	}
	if env.Data.NextCursor == "" {
		t.Fatal("next_cursor phải có (trang đầy)")
	}
}

func TestIntegration_HTTP_GetCustomer_Envelope(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "customers.read")

	ok := doReq(t, handler, "/v1/customers/201", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("ok = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Customer struct {
				ID  int64 `json:"id"`
				B2B *struct {
					TaxCode   string `json:"tax_code"`
					DebtLimit string `json:"debt_limit"`
				} `json:"b2b"`
				Debt struct {
					Amount string `json:"amount"`
				} `json:"debt"`
			} `json:"customer"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, ok.body)
	}
	if env.Data.Customer.B2B == nil || env.Data.Customer.B2B.TaxCode != "0312345678" {
		t.Fatalf("B2B MST sai: %+v", env.Data.Customer.B2B)
	}
	if env.Data.Customer.B2B.DebtLimit != "50000000" {
		t.Fatalf("debt_limit chuỗi sai: %q", env.Data.Customer.B2B.DebtLimit)
	}
	if env.Data.Customer.Debt.Amount != "3000000" {
		t.Fatalf("debt amount sai: %q", env.Data.Customer.Debt.Amount)
	}

	// Không thấy → 404 envelope.
	nf := doReq(t, handler, "/v1/customers/9999", tok)
	if nf.code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", nf.code, nf.body)
	}
	assertErrCode(t, nf.body, "not_found")
}

// ---- helpers ----

func idsOf(cs []domain.Customer) []int64 {
	out := make([]int64, len(cs))
	for i, c := range cs {
		out[i] = c.ID
	}
	return out
}

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	customers.RegisterRoutes(api, fx.svc, verifier)
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
