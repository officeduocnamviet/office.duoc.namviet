package catalog_test

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

	"github.com/Maneva-AI/namviet-backend/internal/catalog"
	"github.com/Maneva-AI/namviet-backend/internal/catalog/app"
	"github.com/Maneva-AI/namviet-backend/internal/catalog/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

type fixture struct {
	svc  *catalog.Service
	pool *pgxpool.Pool
	key  *ecdsa.PrivateKey
}

// seed nạp dữ liệu test vào các bảng public.* (id tường minh vì schema tham
// chiếu không có serial). Tạo 3 product active + 1 soft-deleted, 2 category,
// 1 manufacturer, units cho product 1. Giá numeric để test map money.
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.categories (id, name, slug, status) VALUES
			(1, 'Thuốc kê đơn', 'thuoc-ke-don', 'active'),
			(2, 'Thực phẩm chức năng', 'tpcn', 'active'),
			(3, 'Danh mục đã xóa', 'deleted-cat', 'active')`,
		`UPDATE public.categories SET deleted_at = now() WHERE id = 3`,
		`INSERT INTO public.manufacturers (id, name, slug, country, status) VALUES
			(10, 'Traphaco', 'traphaco', 'VN', 'active')`,
		`INSERT INTO public.products
			(id, name, sku, status, category_id, manufacturer_id, category_name, manufacturer_name,
			 invoice_price, actual_cost, wholesale_unit, retail_unit, conversion_factor, product_images)
		 VALUES
			(101, 'Paracetamol 500mg', 'PARA500', 'active', 1, 10, 'Thuốc kê đơn', 'Traphaco',
			 15000.50, 9000.00, 'Hộp', 'Vỉ', 10, '{"a.jpg","b.jpg"}'),
			(102, 'Vitamin C 1000', 'VITC1000', 'active', 2, 10, 'Thực phẩm chức năng', 'Traphaco',
			 50000, 30000, 'Hộp', 'Viên', 100, '{}'),
			(103, 'Amoxicillin', 'AMOX', 'inactive', 1, 10, 'Thuốc kê đơn', 'Traphaco',
			 20000, 12000, 'Hộp', 'Vỉ', 10, '{}'),
			(104, 'Sản phẩm đã xóa', 'DELSKU', 'active', 1, 10, 'Thuốc kê đơn', 'Traphaco',
			 1, 1, 'Hộp', 'Vỉ', 1, '{}')`,
		`UPDATE public.products SET deleted_at = now() WHERE id = 104`,
		`INSERT INTO public.product_units
			(id, product_id, unit_name, conversion_rate, is_base, is_direct_sale, price_sell, unit_type)
		 VALUES
			(1001, 101, 'Vỉ', 1, true, true, 1600.00, 'retail'),
			(1002, 101, 'Hộp', 10, false, true, 15000.00, 'wholesale')`,
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
	return fixture{svc: catalog.New(pool), pool: pool, key: key}
}

func TestIntegration_ListProducts_SoftDeleteAndStatus(t *testing.T) {
	fx := setup(t)
	// Không filter status → vẫn loại soft-deleted (104), còn 101/102/103.
	res, err := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 50})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	if len(res.Items) != 3 {
		t.Fatalf("đếm = %d, want 3 (loại soft-delete)", len(res.Items))
	}
	for _, p := range res.Items {
		if p.ID == 104 {
			t.Fatal("product soft-deleted 104 không được trả")
		}
	}
	// Filter status=active → loại 103 (inactive) → còn 101/102.
	res2, _ := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 50, Status: "active"})
	if len(res2.Items) != 2 {
		t.Fatalf("active count = %d, want 2", len(res2.Items))
	}
}

func TestIntegration_ListProducts_Keyset(t *testing.T) {
	fx := setup(t)
	// Trang 1: limit 2 → 2 phần tử (101,102) + NextCursor.
	p1, _ := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 2})
	if len(p1.Items) != 2 || p1.Items[0].ID != 101 || p1.Items[1].ID != 102 {
		t.Fatalf("trang 1 sai: %+v", ids(p1.Items))
	}
	if p1.NextCursor == "" {
		t.Fatal("trang 1 đầy phải có NextCursor")
	}
	// Trang 2: dùng cursor → còn 103.
	p2, _ := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 2, Cursor: p1.NextCursor})
	if len(p2.Items) != 1 || p2.Items[0].ID != 103 {
		t.Fatalf("trang 2 sai: %+v", ids(p2.Items))
	}
	if p2.NextCursor != "" {
		t.Fatal("trang 2 hết phải rỗng NextCursor")
	}
}

func TestIntegration_ListProducts_FilterCategoryAndQuery(t *testing.T) {
	fx := setup(t)
	cat := int64(2)
	res, _ := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 50, CategoryID: &cat})
	if len(res.Items) != 1 || res.Items[0].ID != 102 {
		t.Fatalf("filter category sai: %+v", ids(res.Items))
	}
	// Tìm theo SKU.
	q, _ := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 50, Query: "PARA"})
	if len(q.Items) != 1 || q.Items[0].ID != 101 {
		t.Fatalf("search sai: %+v", ids(q.Items))
	}
}

func TestIntegration_ListProducts_MoneyMappedNoFloat(t *testing.T) {
	fx := setup(t)
	res, _ := fx.svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 1})
	p := res.Items[0]
	if p.InvoicePrice.String() != "15000.5" {
		t.Fatalf("invoice_price = %q, want 15000.5", p.InvoicePrice.String())
	}
	if p.ActualCost.String() != "9000" {
		t.Fatalf("actual_cost = %q, want 9000", p.ActualCost.String())
	}
	if len(p.Images) != 2 {
		t.Fatalf("images = %v, want 2 phần tử", p.Images)
	}
}

func TestIntegration_GetProduct_WithUnits(t *testing.T) {
	fx := setup(t)
	d, err := fx.svc.GetProduct(context.Background(), 101)
	if err != nil {
		t.Fatalf("GetProduct: %v", err)
	}
	if d.Product.Name != "Paracetamol 500mg" {
		t.Fatalf("product sai: %+v", d.Product)
	}
	if len(d.Units) != 2 || !d.Units[0].IsBase || d.Units[0].UnitName != "Vỉ" {
		t.Fatalf("units sai (base trước): %+v", d.Units)
	}
}

func TestIntegration_GetProduct_NotFound(t *testing.T) {
	fx := setup(t)
	_, err := fx.svc.GetProduct(context.Background(), 999)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("kind = %v, want NotFound", apperr.KindOf(err))
	}
	// Soft-deleted cũng là NotFound.
	_, err = fx.svc.GetProduct(context.Background(), 104)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("soft-deleted kind = %v, want NotFound", apperr.KindOf(err))
	}
}

func TestIntegration_ListCategories_SoftDelete(t *testing.T) {
	fx := setup(t)
	cats, _ := fx.svc.ListCategories(context.Background(), "")
	if len(cats) != 2 { // loại category 3 đã xóa
		t.Fatalf("categories = %d, want 2", len(cats))
	}
}

// ---- HTTP end-to-end (envelope + authz) ----

func TestIntegration_HTTP_ListProducts_Envelope(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "catalog.read")

	// Thiếu token → 401 envelope.
	noTok := doReq(t, handler, "/v1/products", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Thiếu quyền → 403.
	wrong := doReq(t, handler, "/v1/products", signToken(t, fx.key, "orders:read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Có quyền → 200 + envelope {data:{items,next_cursor}, error:null}.
	ok := doReq(t, handler, "/v1/products?limit=2", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("ok = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Items []struct {
				ID           int64  `json:"id"`
				InvoicePrice string `json:"invoice_price"`
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
	if env.Data.Items[0].InvoicePrice != "15000.5" {
		t.Fatalf("price chuỗi decimal sai: %q", env.Data.Items[0].InvoicePrice)
	}
	if env.Data.NextCursor == "" {
		t.Fatal("next_cursor phải có (trang đầy)")
	}
}

func TestIntegration_HTTP_GetProduct_NotFound_Envelope(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	r := doReq(t, handler, "/v1/products/999", signToken(t, fx.key, "catalog.read"))
	if r.code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", r.code, r.body)
	}
	assertErrCode(t, r.body, "not_found")
}

// ---- helpers ----

func ids(ps []domain.Product) []int64 {
	out := make([]int64, len(ps))
	for i, p := range ps {
		out[i] = p.ID
	}
	return out
}

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	catalog.RegisterRoutes(api, fx.svc, verifier)
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
