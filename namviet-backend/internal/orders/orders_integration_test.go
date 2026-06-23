package orders_test

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
	"github.com/Maneva-AI/namviet-backend/internal/orders"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// id uuid tường minh cho đơn (để GetOrder/keyset xác định được).
const (
	idO1 = "11111111-1111-1111-1111-111111111111" // HD-ORD1, mới nhất
	idO2 = "22222222-2222-2222-2222-222222222222" // HD-ORD2, cũ hơn
	idO3 = "33333333-3333-3333-3333-333333333333" // HD-ORD3, soft-deleted
)

// seed nạp đơn + dòng hàng + finance_transactions vào public.* (DB test trống).
//
//	O1 (HD-ORD1, KH 201, COMPLETED, final 2tr, partial): đã-thu (công nợ) =
//	    500k (completed, INTERNAL) + 300k (completed, BOTH) + 100k (PENDING, INTERNAL)
//	    = ĐÃ THU 900k → còn nợ 1.100.000. (Thanh toán 2 bước, spec mục 55: phiếu
//	    'pending' = NV ĐÃ THU từ khách, nợ giảm NGAY dù chưa nộp quỹ → ĐƯỢC tính.)
//	    Các phiếu KHÔNG được tính:
//	      - 50k book_type='TAX' (không vào sổ thực tế)
//	      - 999k flow='out' (chi, không phải thu)
//	      - 70k flow='in' completed nhưng deleted_at (đã xoá)
//	      - 60k ref_id='HD-OTHER' (trỏ đơn khác)
//	O2 (HD-ORD2, KH 202, PENDING, final 1tr, unpaid): KHÔNG phiếu thu → đã thu
//	    0, còn nợ 1.000.000.
//	O3 (HD-ORD3): soft-deleted → KHÔNG được trả.
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	t1 := "2026-06-01 10:00:00+00" // O2 cũ hơn
	t2 := "2026-06-02 10:00:00+00" // O1 mới hơn
	stmts := []string{
		`INSERT INTO public.orders (id, code, customer_id, creator_id, status, order_type, total_amount, final_amount, payment_status, note, created_at, updated_at) VALUES
			('` + idO1 + `', 'HD-ORD1', 201, NULL, 'COMPLETED', 'B2B', 2100000, 2000000, 'partial', 'don 1', '` + t2 + `', '` + t2 + `'),
			('` + idO2 + `', 'HD-ORD2', 202, NULL, 'PENDING', 'B2B', 1000000, 1000000, 'unpaid', NULL, '` + t1 + `', '` + t1 + `'),
			('` + idO3 + `', 'HD-ORD3', 201, NULL, 'COMPLETED', 'B2B', 500000, 500000, 'unpaid', NULL, '` + t2 + `', '` + t2 + `')`,
		`UPDATE public.orders SET deleted_at = now() WHERE id = '` + idO3 + `'`,
		// Dòng hàng cho O1 (2 dòng); O2 không có dòng.
		`INSERT INTO public.order_items (order_id, product_id, quantity, uom, unit_price, discount, total_line, is_gift, batch_no, expiry_date) VALUES
			('` + idO1 + `', 501, 3, 'Hộp', 500000, 0, 1500000, false, 'LOT-A', '2027-01-31'),
			('` + idO1 + `', 502, 2, 'Vỉ', 250000, 0, 500000, true, NULL, NULL)`,
		// Phiếu thu/chi (finance_transactions). id bigint tường minh; fund_account_id NOT NULL.
		`INSERT INTO public.finance_transactions (id, code, flow, amount, fund_account_id, ref_type, ref_id, status, book_type) VALUES
			(1, 'PT1', 'in',  500000, 1, 'order', 'HD-ORD1', 'completed', 'INTERNAL'),
			(2, 'PT2', 'in',  300000, 1, 'order', 'HD-ORD1', 'completed', 'BOTH'),
			(3, 'PT3', 'in',  100000, 1, 'order', 'HD-ORD1', 'pending',   'INTERNAL'),
			(4, 'PT4', 'in',   50000, 1, 'order', 'HD-ORD1', 'completed', 'TAX'),
			(5, 'PC1', 'out', 999000, 1, 'order', 'HD-ORD1', 'completed', 'INTERNAL'),
			(6, 'PT6', 'in',   60000, 1, 'order', 'HD-OTHER','completed', 'INTERNAL')`,
		// Phiếu IN completed nhưng soft-deleted → không tính.
		`INSERT INTO public.finance_transactions (id, code, flow, amount, fund_account_id, ref_type, ref_id, status, book_type, deleted_at) VALUES
			(7, 'PT7', 'in', 70000, 1, 'order', 'HD-ORD1', 'completed', 'INTERNAL', now())`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("seed: %v\nSQL: %s", err, s)
		}
	}
}

type fixture struct {
	svc  *orders.Service
	pool *pgxpool.Pool
	key  *ecdsa.PrivateKey
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
	return fixture{svc: orders.New(pool, orders.Deps{}), pool: pool, key: key}
}

func TestIntegration_ListOrders_SoftDeleteAndOrder(t *testing.T) {
	fx := setup(t)
	res, err := fx.svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 50})
	if err != nil {
		t.Fatalf("ListOrders: %v", err)
	}
	if len(res.Items) != 2 { // loại O3 soft-deleted
		t.Fatalf("đếm = %d, want 2 (loại soft-delete)", len(res.Items))
	}
	// created_at DESC: O1 (2026-06-02) trước O2 (2026-06-01).
	if res.Items[0].ID != idO1 || res.Items[1].ID != idO2 {
		t.Fatalf("thứ tự created_at DESC sai: %s, %s", res.Items[0].ID, res.Items[1].ID)
	}
}

func TestIntegration_GetOrder_PaidFromFinanceTransactions(t *testing.T) {
	fx := setup(t)
	d, err := fx.svc.GetOrder(context.Background(), idO1)
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	// ĐÃ THU = 500k (INTERNAL) + 300k (BOTH) + 100k (pending, INTERNAL) = 900k; loại
	// TAX/out/deleted/đơn-khác. pending ĐƯỢC tính (NV đã thu, nợ giảm — spec mục 55).
	if d.Order.Payment.Paid.String() != "900000" {
		t.Fatalf("paid = %q, want 900000", d.Order.Payment.Paid.String())
	}
	if d.Order.Payment.Final.String() != "2000000" {
		t.Fatalf("final = %q, want 2000000", d.Order.Payment.Final.String())
	}
	// Còn nợ = 2tr - 900k = 1.100.000.
	if d.Order.Payment.Remaining.String() != "1100000" {
		t.Fatalf("remaining = %q, want 1100000", d.Order.Payment.Remaining.String())
	}
	// 2 dòng hàng.
	if len(d.Lines) != 2 {
		t.Fatalf("lines = %d, want 2", len(d.Lines))
	}
}

func TestIntegration_GetOrder_NoPayments_ZeroPaid(t *testing.T) {
	fx := setup(t)
	d, err := fx.svc.GetOrder(context.Background(), idO2)
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	if !d.Order.Payment.Paid.IsZero() {
		t.Fatalf("không phiếu thu phải paid 0, got %q", d.Order.Payment.Paid.String())
	}
	if d.Order.Payment.Remaining.String() != "1000000" {
		t.Fatalf("remaining = %q, want 1000000", d.Order.Payment.Remaining.String())
	}
	if len(d.Lines) != 0 {
		t.Fatalf("O2 không có dòng hàng, got %d", len(d.Lines))
	}
}

func TestIntegration_GetOrder_NotFound(t *testing.T) {
	fx := setup(t)
	_, err := fx.svc.GetOrder(context.Background(), "99999999-9999-9999-9999-999999999999")
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("kind = %v, want NotFound", apperr.KindOf(err))
	}
	// Soft-deleted cũng là NotFound.
	_, err = fx.svc.GetOrder(context.Background(), idO3)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("soft-deleted kind = %v, want NotFound", apperr.KindOf(err))
	}
}

func TestIntegration_ListOrders_FilterPaymentStatusAndCustomer(t *testing.T) {
	fx := setup(t)
	// payment_status=partial → chỉ O1.
	p, _ := fx.svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 50, PaymentStatus: "partial"})
	if len(p.Items) != 1 || p.Items[0].ID != idO1 {
		t.Fatalf("filter partial sai: %d items", len(p.Items))
	}
	// customer_id=202 → chỉ O2.
	cid := int64(202)
	c, _ := fx.svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 50, CustomerID: &cid})
	if len(c.Items) != 1 || c.Items[0].ID != idO2 {
		t.Fatalf("filter customer sai: %d items", len(c.Items))
	}
}

func TestIntegration_ListOrders_Keyset(t *testing.T) {
	fx := setup(t)
	p1, _ := fx.svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 1})
	if len(p1.Items) != 1 || p1.Items[0].ID != idO1 {
		t.Fatalf("trang 1 sai: %+v", p1.Items)
	}
	if p1.NextCursor == "" {
		t.Fatal("trang 1 đầy phải có NextCursor")
	}
	p2, _ := fx.svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 1, Cursor: p1.NextCursor})
	if len(p2.Items) != 1 || p2.Items[0].ID != idO2 {
		t.Fatalf("trang 2 sai: %+v", p2.Items)
	}
	// Trang 2 vừa đầy (len==limit) nên keyset vẫn phát cursor (client phải fetch
	// thêm 1 lần để biết hết — ngữ nghĩa keyset chuẩn, giống customers). Trang 3
	// phải RỖNG items + rỗng cursor.
	if p2.NextCursor == "" {
		t.Fatal("trang 2 đầy vẫn phát NextCursor (ngữ nghĩa keyset)")
	}
	p3, _ := fx.svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 1, Cursor: p2.NextCursor})
	if len(p3.Items) != 0 {
		t.Fatalf("trang 3 phải hết (0 item), got %d", len(p3.Items))
	}
	if p3.NextCursor != "" {
		t.Fatalf("trang 3 hết phải rỗng NextCursor, got %q", p3.NextCursor)
	}
}

// ---- HTTP end-to-end (envelope + authz) ----

func TestIntegration_HTTP_ListOrders_Envelope(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "orders.read")

	// Thiếu token → 401 envelope.
	noTok := doReq(t, handler, "/v1/orders", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Sai quyền → 403.
	wrong := doReq(t, handler, "/v1/orders", signToken(t, fx.key, "catalog.read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Có quyền → 200 + envelope {data:{items,next_cursor}, error:null}.
	ok := doReq(t, handler, "/v1/orders?limit=50", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("ok = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Items []struct {
				ID      string `json:"id"`
				Code    string `json:"code"`
				Payment struct {
					Paid      string `json:"paid"`
					Remaining string `json:"remaining"`
				} `json:"payment"`
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
	// Item đầu (O1) đã thu 900k / còn nợ 1.1tr dạng chuỗi decimal (gồm phiếu pending).
	if env.Data.Items[0].Code != "HD-ORD1" ||
		env.Data.Items[0].Payment.Paid != "900000" ||
		env.Data.Items[0].Payment.Remaining != "1100000" {
		t.Fatalf("payment envelope sai: %+v", env.Data.Items[0])
	}
}

func TestIntegration_HTTP_GetOrder_Envelope(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "orders.read")

	ok := doReq(t, handler, "/v1/orders/"+idO1, tok)
	if ok.code != http.StatusOK {
		t.Fatalf("ok = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Order struct {
				Code    string `json:"code"`
				Payment struct {
					Paid      string `json:"paid"`
					Remaining string `json:"remaining"`
				} `json:"payment"`
			} `json:"order"`
			Lines []struct {
				ProductID int64  `json:"product_id"`
				Quantity  string `json:"quantity"`
				UnitPrice string `json:"unit_price"`
				LineTotal string `json:"line_total"`
			} `json:"lines"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, ok.body)
	}
	if env.Data.Order.Code != "HD-ORD1" || env.Data.Order.Payment.Paid != "900000" {
		t.Fatalf("order envelope sai: %+v", env.Data.Order)
	}
	if len(env.Data.Lines) != 2 {
		t.Fatalf("lines = %d, want 2; body=%s", len(env.Data.Lines), ok.body)
	}
	// Tìm dòng theo product_id (thứ tự dòng không phải hợp đồng cứng — hai dòng
	// seed cùng created_at nên tie-break theo id uuid, không xác định). Dòng 501:
	// 3 Hộp x 500000 = 1.500.000 (chuỗi thập phân, không float).
	found501 := false
	for _, l := range env.Data.Lines {
		if l.ProductID == 501 {
			found501 = true
			if l.Quantity != "3" || l.UnitPrice != "500000" || l.LineTotal != "1500000" {
				t.Fatalf("line 501 envelope sai: %+v", l)
			}
		}
	}
	if !found501 {
		t.Fatalf("thiếu dòng product 501; body=%s", ok.body)
	}

	// Không thấy → 404 envelope.
	nf := doReq(t, handler, "/v1/orders/99999999-9999-9999-9999-999999999999", tok)
	if nf.code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", nf.code, nf.body)
	}
	assertErrCode(t, nf.body, "not_found")
}

func TestIntegration_HTTP_ListOrders_BadDate_Validation(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "orders.read")
	bad := doReq(t, handler, "/v1/orders?from_date=2026-13-99", tok)
	if bad.code != http.StatusUnprocessableEntity {
		t.Fatalf("bad date = %d, want 422; body=%s", bad.code, bad.body)
	}
	assertErrCode(t, bad.body, "validation_error")
}

// ---- helpers ----

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	orders.RegisterRoutes(api, fx.svc, verifier)
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
