package orders_test

// HTTP integration test cho các route GHI ORCHESTRATION (P4b): POST ship / payments
// / pos/sales / customers/{id}/payments. Đi qua httptest + huma (DTO parse + envelope
// {data,error} + guard authz orders.write) trên Service deps THẬT + testcontainers.
// Lấp lỗ hổng: trước đây orchestration chỉ test ở tầng service, chưa qua HTTP.
// Tái dùng signToken/doPost/httpResp (cùng package orders_test).

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Maneva-AI/namviet-backend/internal/orders"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// buildOrchHTTP dựng handler HTTP từ Service deps THẬT + key (verifier từ public key).
func buildOrchHTTP(t *testing.T, svc *orders.Service, key *ecdsa.PrivateKey) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	orders.RegisterRoutes(api, svc, authn.NewVerifier(&key.PublicKey))
	return r
}

func orchKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return k
}

// TestIntegration_HTTP_PosSale: POST /v1/pos/sales — 401 thiếu token, 403 sai quyền,
// 201 + envelope khi orders.write. Tiền ra chuỗi, đơn COMPLETED + paid.
func TestIntegration_HTTP_PosSale(t *testing.T) {
	svc, _ := setupOrch(t)
	key := orchKey(t)
	h := buildOrchHTTP(t, svc, key)
	body := `{"lines":[{"product_id":701,"quantity":4,"uom":"Hộp","unit_price":"10000","discount":"0"}],
		"warehouse_id":9,"fund_account_id":1,"issue_invoice":true,"customer_tax_code":"0312345678",
		"serial":"C26TAA","mau_so":"1","vat_lines":[{"vat_rate":"0.08"}]}`

	if r := doPost(t, h, "/v1/pos/sales", body, "", ""); r.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; %s", r.code, r.body)
	}
	if r := doPost(t, h, "/v1/pos/sales", body, signToken(t, key, "orders.read"), ""); r.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; %s", r.code, r.body)
	}
	r := doPost(t, h, "/v1/pos/sales", body, signToken(t, key, "orders.write"), "")
	if r.code != http.StatusCreated {
		t.Fatalf("pos sale = %d, want 201; %s", r.code, r.body)
	}
	var env struct {
		Data struct {
			Order struct {
				Status        string `json:"status"`
				PaymentStatus string `json:"payment_status"`
			} `json:"order"`
			Invoice *struct {
				InvoiceNo int64 `json:"invoice_no"`
			} `json:"invoice"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(r.body), &env); err != nil {
		t.Fatalf("decode: %v; %s", err, r.body)
	}
	if env.Error != nil || env.Data.Order.Status != "COMPLETED" || env.Data.Order.PaymentStatus != "paid" || env.Data.Invoice == nil {
		t.Fatalf("envelope POS sai: %s", r.body)
	}
}

// TestIntegration_HTTP_ShipThenPayment: tạo+duyệt đơn (svc) → ship qua HTTP → thu qua
// HTTP. Kiểm guard + envelope + chuyển trạng thái.
func TestIntegration_HTTP_ShipThenPayment(t *testing.T) {
	svc, _ := setupOrch(t)
	key := orchKey(t)
	h := buildOrchHTTP(t, svc, key)
	tok := signToken(t, key, "orders.write")

	created, err := svc.CreateOrder(t.Context(), app.CreateOrderInput{
		OrderType: "B2B",
		Lines: []domain.DraftLine{
			{ProductID: 701, Quantity: qty(2), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: mustMoney(t, "0")},
			{ProductID: 702, Quantity: qty(3), UOM: "Hộp", UnitPrice: mustMoney(t, "20000"), Discount: mustMoney(t, "0")},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if _, err := svc.ConfirmOrder(t.Context(), created.Order.ID); err != nil {
		t.Fatalf("ConfirmOrder: %v", err)
	}

	// Ship qua HTTP.
	shipBody := `{"warehouse_id":9,"customer_tax_code":"0312345678","serial":"C26TAA","mau_so":"1",
		"vat_lines":[{"vat_rate":"0.08"},{"vat_rate":"0.08"}]}`
	rs := doPost(t, h, "/v1/orders/"+created.Order.ID+"/ship", shipBody, tok, "")
	if rs.code != http.StatusOK {
		t.Fatalf("ship = %d, want 200; %s", rs.code, rs.body)
	}
	// Thu tiền qua HTTP (80000) → paid.
	rp := doPost(t, h, "/v1/orders/"+created.Order.ID+"/payments", `{"amount":"80000","fund_account_id":1,"book_type":"BOTH"}`, tok, "pay-http-1")
	if rp.code != http.StatusOK {
		t.Fatalf("payment = %d, want 200; %s", rp.code, rp.body)
	}
	var env struct {
		Data struct {
			PaymentStatus string `json:"payment_status"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(rp.body), &env); err != nil || env.Data.PaymentStatus != "paid" {
		t.Fatalf("payment envelope sai (want paid): %s", rp.body)
	}
}

// TestIntegration_HTTP_LumpPayment: tạo 2 đơn 1 khách (svc) → POST /v1/customers/{id}/
// payments qua HTTP → phân bổ, envelope có allocations.
func TestIntegration_HTTP_LumpPayment(t *testing.T) {
	svc, _ := setupOrch(t)
	key := orchKey(t)
	h := buildOrchHTTP(t, svc, key)
	tok := signToken(t, key, "orders.write")
	cust := int64(8888)
	for _, q := range []int64{10, 20} { // final 100k + 200k
		if _, err := svc.CreateOrder(t.Context(), app.CreateOrderInput{
			CustomerID: &cust, OrderType: "B2B",
			Lines: []domain.DraftLine{{ProductID: 701, Quantity: qty(q), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: mustMoney(t, "0")}},
		}); err != nil {
			t.Fatalf("CreateOrder: %v", err)
		}
	}
	r := doPost(t, h, "/v1/customers/8888/payments", `{"amount":"150000","fund_account_id":1,"book_type":"BOTH"}`, tok, "lump-http-1")
	if r.code != http.StatusOK {
		t.Fatalf("lump = %d, want 200; %s", r.code, r.body)
	}
	var env struct {
		Data struct {
			Allocations []struct {
				Amount string `json:"amount"`
			} `json:"allocations"`
			Leftover string `json:"leftover"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(r.body), &env); err != nil {
		t.Fatalf("decode: %v; %s", err, r.body)
	}
	// 150k phân bổ: đơn1 100k (paid) + đơn2 50k (partial), leftover 0.
	if len(env.Data.Allocations) != 2 || env.Data.Allocations[0].Amount != "100000" || env.Data.Allocations[1].Amount != "50000" || env.Data.Leftover != "0" {
		t.Fatalf("lump envelope sai: %s", r.body)
	}
}
