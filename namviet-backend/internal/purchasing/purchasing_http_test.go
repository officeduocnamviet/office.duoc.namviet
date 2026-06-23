package purchasing_test

// HTTP integration test cho module purchasing: full flow create→confirm→receive→pay
// + GET qua httptest + huma (DTO parse + envelope {data,error} + guard authz
// purchasing.write/read) trên Service deps THẬT + testcontainers. Lấp lỗ hổng:
// trước đây purchasing chỉ test ở tầng service, chưa qua HTTP.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing"
)

func phKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return k
}

func phToken(t *testing.T, key *ecdsa.PrivateKey, perms ...string) string {
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

type phResp struct {
	code int
	body string
}

func phDo(t *testing.T, h http.Handler, method, path, body, bearer, idem string) phResp {
	t.Helper()
	var rdr *strings.Reader
	if body == "" {
		rdr = strings.NewReader("")
	} else {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return phResp{code: rec.Code, body: rec.Body.String()}
}

func phBuildHTTP(t *testing.T, fx fixture, key *ecdsa.PrivateKey) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	purchasing.RegisterRoutes(api, fx.svc, authn.NewVerifier(&key.PublicKey))
	return r
}

// TestIntegration_HTTP_Purchasing_FullFlow: create→confirm→receive→pay qua HTTP +
// authz trên route tạo. Khẳng định envelope + chuyển trạng thái tới paid.
func TestIntegration_HTTP_Purchasing_FullFlow(t *testing.T) {
	fx := setup(t)
	key := phKey(t)
	h := phBuildHTTP(t, fx, key)
	write := phToken(t, key, "purchasing.write")
	read := phToken(t, key, "purchasing.read")

	createBody := `{"supplier_name":"NCC HTTP","lines":[
		{"product_id":100,"quantity":"10","unit_cost":"1000","vat_rate":"0.08","batch_code":"LOT-H","expiry_date":"2027-12-31"}]}`

	// Authz route tạo: 401 thiếu token, 403 sai quyền.
	if r := phDo(t, h, http.MethodPost, "/v1/purchase-orders", createBody, "", ""); r.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; %s", r.code, r.body)
	}
	if r := phDo(t, h, http.MethodPost, "/v1/purchase-orders", createBody, read, ""); r.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; %s", r.code, r.body)
	}

	// Tạo PO (201) + envelope.
	rc := phDo(t, h, http.MethodPost, "/v1/purchase-orders", createBody, write, "po-http-1")
	if rc.code != http.StatusCreated {
		t.Fatalf("create = %d, want 201; %s", rc.code, rc.body)
	}
	var created struct {
		Data struct {
			PurchaseOrder struct {
				ID          string `json:"id"`
				Status      string `json:"status"`
				TotalAmount string `json:"total_amount"`
				VATAmount   string `json:"vat_amount"`
			} `json:"purchase_order"`
			Lines []json.RawMessage `json:"lines"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(rc.body), &created); err != nil {
		t.Fatalf("decode create: %v; %s", err, rc.body)
	}
	po := created.Data.PurchaseOrder
	if created.Error != nil || po.Status != "draft" || po.TotalAmount != "10000" || po.VATAmount != "800" || len(created.Data.Lines) != 1 {
		t.Fatalf("envelope tạo PO sai: %s", rc.body)
	}
	id := po.ID

	// Confirm (200).
	if r := phDo(t, h, http.MethodPost, "/v1/purchase-orders/"+id+"/confirm", "", write, ""); r.code != http.StatusOK {
		t.Fatalf("confirm = %d, want 200; %s", r.code, r.body)
	}
	// Receive (200) — nhập kho.
	rr := phDo(t, h, http.MethodPost, "/v1/purchase-orders/"+id+"/receive", `{"warehouse_id":1}`, write, "")
	if rr.code != http.StatusOK {
		t.Fatalf("receive = %d, want 200; %s", rr.code, rr.body)
	}
	// Pay (200) — chi NCC 10800 = total 10000 + vat 800.
	rp := phDo(t, h, http.MethodPost, "/v1/purchase-orders/"+id+"/pay", `{"amount":"10800","fund_account_id":1,"book_type":"INTERNAL"}`, write, "pay-http-1")
	if rp.code != http.StatusOK {
		t.Fatalf("pay = %d, want 200; %s", rp.code, rp.body)
	}
	var paid struct {
		Data struct {
			PurchaseOrder struct {
				Status string `json:"status"`
			} `json:"purchase_order"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(rp.body), &paid); err != nil || paid.Data.PurchaseOrder.Status != "paid" {
		t.Fatalf("pay envelope want status=paid: %s", rp.body)
	}

	// GET detail (200) — đọc qua HTTP, guard purchasing.read.
	rg := phDo(t, h, http.MethodGet, "/v1/purchase-orders/"+id, "", read, "")
	if rg.code != http.StatusOK {
		t.Fatalf("get = %d, want 200; %s", rg.code, rg.body)
	}
	// Tồn đã tăng (nhập kho qua HTTP receive).
	if got := stockTotal(t, fx.pool, 1, 100); !decEqual(got, "10") {
		t.Fatalf("tồn sp100 sau receive HTTP = %s, want 10", got)
	}
}
