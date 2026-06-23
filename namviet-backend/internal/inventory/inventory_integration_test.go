package inventory_test

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

	"github.com/Maneva-AI/namviet-backend/internal/inventory"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

type fixture struct {
	svc  *inventory.Service
	pool *pgxpool.Pool
	key  *ecdsa.PrivateKey
}

// seed nạp dữ liệu test vào các bảng public.* kho (id tường minh vì schema tham
// chiếu không có serial). Tạo:
//   - 2 kho active (1,2) + 1 kho đã đóng (3, deleted_at).
//   - product_inventory: tồn product 100 ở kho 1 và 2; tồn product 100 ở kho ĐÃ
//     ĐÓNG 3 (phải bị loại); 1 dòng warehouse_id NULL (phải bị loại); product 200.
//   - batches: 3 lô của product 100 với hạn dùng khác nhau (FEFO) + 1 lô đã xóa +
//     1 lô hết tồn (quantity 0, phải bị loại).
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.warehouses (id, key, name, unit, type, status) VALUES
			(1, 'kho-tong', 'Kho Tổng', 'Hộp', 'central', 'active'),
			(2, 'cn-q1', 'Chi nhánh Q1', 'Hộp', 'retail', 'active'),
			(3, 'kho-dong', 'Kho Đã Đóng', 'Hộp', 'retail', 'active')`,
		`UPDATE public.warehouses SET deleted_at = now() WHERE id = 3`,
		// Tồn tổng: product 100 ở kho 1 (numeric có thập phân) + kho 2; product 100
		// ở kho ĐÓNG 3 (loại qua join); 1 dòng warehouse NULL (loại); product 200.
		`INSERT INTO public.product_inventory (id, product_id, warehouse_id, stock_quantity, min_stock, max_stock, shelf_location) VALUES
			(1001, 100, 1, 120.5, 10, 500, 'A1'),
			(1002, 100, 2, 30, 5, 200, 'B2'),
			(1003, 100, 3, 999, 0, 0, 'X'),
			(1004, 100, NULL, 7, 0, 0, 'no-wh'),
			(1005, 200, 1, 45, 0, 0, 'C3')`,
		// Lô của product 100: 3 lô hạn dùng khác nhau (chèn KHÔNG theo thứ tự hạn).
		`INSERT INTO public.batches (id, product_id, batch_code, expiry_date, manufacturing_date, inbound_price) VALUES
			(501, 100, 'LOT-FAR', '2027-12-31', '2025-01-01', 9000.50),
			(502, 100, 'LOT-NEAR', '2026-02-15', NULL, 9100),
			(503, 100, 'LOT-MID', '2026-09-30', '2025-06-01', 9050),
			(504, 100, 'LOT-DELETED', '2026-01-01', NULL, 1),
			(505, 100, 'LOT-EMPTY', '2026-03-01', NULL, 1)`,
		`UPDATE public.batches SET deleted_at = now() WHERE id = 504`,
		// Tồn theo lô tại kho 1. Lô 504 (đã xóa) và 505 (tồn 0) phải bị loại khỏi FEFO.
		`INSERT INTO public.inventory_batches (id, warehouse_id, product_id, batch_id, quantity) VALUES
			(9001, 1, 100, 501, 40),
			(9002, 1, 100, 502, 25.5),
			(9003, 1, 100, 503, 10),
			(9004, 1, 100, 504, 99),
			(9005, 1, 100, 505, 0),
			(9006, 2, 100, 502, 5)`,
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
	return fixture{svc: inventory.New(pool), pool: pool, key: key}
}

func TestIntegration_ListWarehouses_ExcludeClosed(t *testing.T) {
	fx := setup(t)
	ws, err := fx.svc.ListWarehouses(context.Background(), "", 100)
	if err != nil {
		t.Fatalf("ListWarehouses: %v", err)
	}
	if len(ws) != 2 { // loại kho 3 đã đóng (deleted_at)
		t.Fatalf("warehouses = %d, want 2 (loại kho đã đóng)", len(ws))
	}
	for _, w := range ws {
		if w.ID == 3 {
			t.Fatal("kho đã đóng (3) không được trả")
		}
	}
}

func TestIntegration_ListStock_FilterAndCloseAndNull(t *testing.T) {
	fx := setup(t)
	// Theo product 100: dòng 1001 (kho1), 1002 (kho2). Loại 1003 (kho đóng) và
	// 1004 (warehouse NULL).
	pid := int64(100)
	res, err := fx.svc.ListStock(context.Background(), app.ListStockQuery{Limit: 50, ProductID: &pid})
	if err != nil {
		t.Fatalf("ListStock: %v", err)
	}
	if len(res.Items) != 2 {
		t.Fatalf("stock product 100 = %d, want 2 (loại kho đóng + warehouse NULL); ids=%v",
			len(res.Items), stockIDs(res.Items))
	}
	// quantity numeric thập phân giữ chính xác (không float).
	var found bool
	for _, s := range res.Items {
		if s.ID == 1001 {
			found = true
			if s.Quantity.String() != "120.5" {
				t.Fatalf("quantity 1001 = %q, want 120.5 (không mất chính xác)", s.Quantity.String())
			}
		}
	}
	if !found {
		t.Fatal("thiếu dòng tồn 1001")
	}

	// Lọc theo kho 1 + product 100 → chỉ 1001.
	wid := int64(1)
	r2, _ := fx.svc.ListStock(context.Background(), app.ListStockQuery{Limit: 50, ProductID: &pid, WarehouseID: &wid})
	if len(r2.Items) != 1 || r2.Items[0].ID != 1001 {
		t.Fatalf("lọc kho 1 sai: %v", stockIDs(r2.Items))
	}
}

func TestIntegration_ListStock_Keyset(t *testing.T) {
	fx := setup(t)
	// Tất cả product, limit 2 → trang đầy + cursor; trang sau tiếp tục theo id.
	p1, _ := fx.svc.ListStock(context.Background(), app.ListStockQuery{Limit: 2})
	if len(p1.Items) != 2 {
		t.Fatalf("trang 1 = %d, want 2; ids=%v", len(p1.Items), stockIDs(p1.Items))
	}
	if p1.NextCursor == "" {
		t.Fatal("trang 1 đầy phải có NextCursor")
	}
	p2, _ := fx.svc.ListStock(context.Background(), app.ListStockQuery{Limit: 2, Cursor: p1.NextCursor})
	// id tăng dần, không lặp lại id trang trước.
	for _, s := range p2.Items {
		if s.ID <= p1.Items[len(p1.Items)-1].ID {
			t.Fatalf("keyset lặp id %d <= cursor", s.ID)
		}
	}
}

func TestIntegration_ListBatchesFEFO_Order(t *testing.T) {
	fx := setup(t)
	// Product 100, mọi kho: lô còn tồn = 501(2027-12),502(2026-02 kho1 + kho2),
	// 503(2026-09). Loại 504 (đã xóa) + 505 (tồn 0). FEFO theo hạn tăng dần:
	// 502(2026-02) trước, rồi 503(2026-09), rồi 501(2027-12).
	batches, err := fx.svc.ListBatchesFEFO(context.Background(), 100, nil)
	if err != nil {
		t.Fatalf("ListBatchesFEFO: %v", err)
	}
	// 502 xuất hiện 2 lần (kho1 + kho2) → tổng 4 dòng còn tồn.
	if len(batches) != 4 {
		t.Fatalf("batches = %d, want 4 (loại đã xóa + tồn 0); codes=%v", len(batches), batchCodes(batches))
	}
	// Dòng đầu phải là lô hết hạn sớm nhất (502, 2026-02-15).
	if batches[0].BatchID != 502 {
		t.Fatalf("FEFO dòng đầu BatchID = %d, want 502 (hết hạn sớm nhất)", batches[0].BatchID)
	}
	// Hạn dùng không giảm dần (đảm bảo ASC).
	for i := 1; i < len(batches); i++ {
		if batches[i].ExpiryDate.Before(batches[i-1].ExpiryDate) {
			t.Fatalf("FEFO không tăng dần ở vị trí %d: %v", i, expiries(batches))
		}
	}
	// Không có lô đã xóa (504) hoặc tồn 0 (505).
	for _, b := range batches {
		if b.BatchID == 504 || b.BatchID == 505 {
			t.Fatalf("lô %d (đã xóa/tồn 0) không được trả", b.BatchID)
		}
		if !b.Quantity.IsPositive() {
			t.Fatalf("lô %d tồn không dương", b.BatchID)
		}
	}

	// Lọc theo kho 1 → loại dòng 502 ở kho 2 → còn 3 dòng.
	wid := int64(1)
	r1, _ := fx.svc.ListBatchesFEFO(context.Background(), 100, &wid)
	if len(r1) != 3 {
		t.Fatalf("lọc kho 1 = %d, want 3", len(r1))
	}
}

// ---- HTTP end-to-end (envelope + authz) ----

func TestIntegration_HTTP_Inventory_EnvelopeAndAuthz(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "inventory.read")

	// Thiếu token → 401.
	noTok := doReq(t, handler, "/v1/warehouses", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Sai quyền → 403.
	wrong := doReq(t, handler, "/v1/warehouses", signToken(t, fx.key, "catalog.read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Có quyền → 200 + envelope {data:{items}, error:null}.
	ok := doReq(t, handler, "/v1/warehouses", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("warehouses = %d, want 200; body=%s", ok.code, ok.body)
	}
	var wenv struct {
		Data struct {
			Items []struct {
				ID int64 `json:"id"`
			} `json:"items"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(ok.body), &wenv); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, ok.body)
	}
	if wenv.Error != nil || len(wenv.Data.Items) != 2 {
		t.Fatalf("warehouses envelope sai: %s", ok.body)
	}

	// Stock: quantity ra chuỗi thập phân (không float).
	sres := doReq(t, handler, "/v1/inventory/stock?product_id=100&warehouse_id=1", tok)
	if sres.code != http.StatusOK {
		t.Fatalf("stock = %d; body=%s", sres.code, sres.body)
	}
	var senv struct {
		Data struct {
			Items []struct {
				ID       int64  `json:"id"`
				Quantity string `json:"quantity"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(sres.body), &senv); err != nil {
		t.Fatalf("decode stock: %v (body=%s)", err, sres.body)
	}
	if len(senv.Data.Items) != 1 || senv.Data.Items[0].Quantity != "120.5" {
		t.Fatalf("stock quantity chuỗi decimal sai: %s", sres.body)
	}

	// Batches FEFO: dòng đầu là lô hết hạn sớm nhất (LOT-NEAR).
	bres := doReq(t, handler, "/v1/inventory/batches?product_id=100", tok)
	if bres.code != http.StatusOK {
		t.Fatalf("batches = %d; body=%s", bres.code, bres.body)
	}
	var benv struct {
		Data struct {
			Items []struct {
				BatchCode    string `json:"batch_code"`
				ExpiryDate   string `json:"expiry_date"`
				InboundPrice string `json:"inbound_price"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(bres.body), &benv); err != nil {
		t.Fatalf("decode batches: %v (body=%s)", err, bres.body)
	}
	if len(benv.Data.Items) != 4 || benv.Data.Items[0].BatchCode != "LOT-NEAR" {
		t.Fatalf("FEFO HTTP sai (dòng đầu phải LOT-NEAR): %s", bres.body)
	}
	if benv.Data.Items[0].ExpiryDate != "2026-02-15" {
		t.Fatalf("expiry_date format sai: %q", benv.Data.Items[0].ExpiryDate)
	}

	// Batches thiếu product_id (required) → 422 validation.
	missing := doReq(t, handler, "/v1/inventory/batches", tok)
	if missing.code != http.StatusUnprocessableEntity {
		t.Fatalf("thiếu product_id = %d, want 422; body=%s", missing.code, missing.body)
	}
}

// ---- helpers ----

func stockIDs(items []domain.StockItem) []int64 {
	out := make([]int64, len(items))
	for i, s := range items {
		out[i] = s.ID
	}
	return out
}

func batchCodes(items []domain.Batch) []string {
	out := make([]string, len(items))
	for i, b := range items {
		out[i] = b.BatchCode
	}
	return out
}

func expiries(items []domain.Batch) []string {
	out := make([]string, len(items))
	for i, b := range items {
		out[i] = b.ExpiryDate.Format("2006-01-02")
	}
	return out
}

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	inventory.RegisterRoutes(api, fx.svc, verifier)
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
