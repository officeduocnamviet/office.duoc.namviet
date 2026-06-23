package orders_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

// Dùng chung fixture/setup/buildHTTP/signToken/httpResp/assertErrCode định nghĩa ở
// orders_integration_test.go (cùng package orders_test). File này CHỈ thêm helper
// GHI riêng (doPost, extractCode, qty, mustMoney, twoLineInput).

func qty(n int64) domain.Quantity { return domain.QuantityFromInt(n) }

func mustMoney(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return m
}

func twoLineInput(t *testing.T, idemKey string) app.CreateOrderInput {
	t.Helper()
	return app.CreateOrderInput{
		OrderType: "B2B",
		Note:      "đơn test",
		IdemKey:   idemKey,
		Lines: []domain.DraftLine{
			{ProductID: 100, Quantity: qty(3), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: mustMoney(t, "5000")},
			{ProductID: 200, Quantity: qty(2), UOM: "Vỉ", UnitPrice: mustMoney(t, "20000"), Discount: money.Zero()},
		},
	}
}

// ---- use-case (qua Service.Write/Read) ----

// TestIntegration_CreateOrder_PersistsHeaderAndItems: tạo đơn → ghi orders +
// order_items đúng; sinh code (tiền tố DH + zero-pad); total/final đúng quy ước
// (total trước chiết khấu = 70000; final sau = 65000); status PENDING.
func TestIntegration_CreateOrder_PersistsHeaderAndItems(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()

	res, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, ""))
	if err != nil {
		t.Fatalf("CreateOrder phải thành công: %v", err)
	}
	o := res.Order
	if o.Status != domain.StatusPending.String() {
		t.Errorf("status = %q, want PENDING", o.Status)
	}
	if !strings.HasPrefix(o.Code, "DH") || len(o.Code) != 10 {
		t.Errorf("code = %q, want tiền tố DH + 8 chữ số", o.Code)
	}
	if o.Total.String() != "70000" {
		t.Errorf("total_amount = %q, want 70000 (trước chiết khấu)", o.Total.String())
	}
	if o.Final.String() != "65000" {
		t.Errorf("final_amount = %q, want 65000 (sau chiết khấu)", o.Final.String())
	}
	if len(res.Lines) != 2 {
		t.Fatalf("lines = %d, want 2", len(res.Lines))
	}

	// Đọc lại qua đường ĐỌC: phải khớp + có 2 dòng trong DB.
	detail, err := fx.svc.Read().GetOrder(ctx, o.ID)
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	if detail.Order.Code != o.Code || len(detail.Lines) != 2 {
		t.Fatalf("đọc lại lệch: code=%s lines=%d", detail.Order.Code, len(detail.Lines))
	}
	// Kiểm trực tiếp DB: orders + order_items đúng số dòng, total_line dòng 0 = 25000.
	var cnt int
	if err := fx.pool.QueryRow(ctx,
		`SELECT count(*) FROM public.order_items WHERE order_id=$1::uuid AND deleted_at IS NULL`, o.ID).Scan(&cnt); err != nil {
		t.Fatalf("count items: %v", err)
	}
	if cnt != 2 {
		t.Fatalf("order_items trong DB = %d, want 2", cnt)
	}
	var lineTotal string
	if err := fx.pool.QueryRow(ctx,
		`SELECT total_line::text FROM public.order_items WHERE order_id=$1::uuid AND product_id=100`, o.ID).Scan(&lineTotal); err != nil {
		t.Fatalf("line total: %v", err)
	}
	if lineTotal != "25000" {
		t.Errorf("total_line product 100 = %q, want 25000", lineTotal)
	}
}

// TestIntegration_CreateOrder_GeneratesDistinctCodes: 2 đơn liên tiếp → mã KHÁC
// nhau (sequence tăng).
func TestIntegration_CreateOrder_GeneratesDistinctCodes(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	a, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, ""))
	if err != nil {
		t.Fatalf("create a: %v", err)
	}
	b, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, ""))
	if err != nil {
		t.Fatalf("create b: %v", err)
	}
	if a.Order.Code == b.Order.Code {
		t.Fatalf("2 đơn phải khác mã, đều = %s", a.Order.Code)
	}
}

// TestIntegration_CreateOrder_Idempotent: cùng Idempotency-Key → 1 đơn (gọi 2 lần
// trả cùng id/code, DB chỉ có 1 đơn cho key đó).
func TestIntegration_CreateOrder_Idempotent(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	const key = "idem-key-abc-123"

	first, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, key))
	if err != nil {
		t.Fatalf("create lần 1: %v", err)
	}
	second, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, key))
	if err != nil {
		t.Fatalf("create lần 2 (idempotent) phải thành công: %v", err)
	}
	if first.Order.ID != second.Order.ID || first.Order.Code != second.Order.Code {
		t.Fatalf("idempotent phải trả CÙNG đơn: %s/%s vs %s/%s",
			first.Order.ID, first.Order.Code, second.Order.ID, second.Order.Code)
	}
	// Chỉ 1 đơn trong DB cho key này (qua ánh xạ + count orders cùng code).
	var orders int
	if err := fx.pool.QueryRow(ctx,
		`SELECT count(*) FROM public.orders WHERE code=$1 AND deleted_at IS NULL`, first.Order.Code).Scan(&orders); err != nil {
		t.Fatalf("count orders: %v", err)
	}
	if orders != 1 {
		t.Fatalf("idempotent key phải cho đúng 1 đơn, có %d", orders)
	}
}

// TestIntegration_CreateOrder_ValidationError: input sai (không có dòng) → Validation.
func TestIntegration_CreateOrder_ValidationError(t *testing.T) {
	fx := setup(t)
	_, err := fx.svc.Write().CreateOrder(context.Background(), app.CreateOrderInput{OrderType: "B2C"})
	if err == nil {
		t.Fatal("tạo đơn không dòng phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("phải Validation, got %v", apperr.KindOf(err))
	}
}

// TestIntegration_StatusMachine_HappyPath: PENDING→CONFIRMED, rồi (giả lập đã
// SHIPPING bằng update trực tiếp vì ShipOrder = P4b) →COMPLETED.
func TestIntegration_StatusMachine_HappyPath(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	res, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, ""))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	id := res.Order.ID

	confirmed, err := fx.svc.Write().ConfirmOrder(ctx, id)
	if err != nil {
		t.Fatalf("confirm phải thành công: %v", err)
	}
	if confirmed.Status != domain.StatusConfirmed.String() {
		t.Fatalf("sau confirm = %q, want CONFIRMED", confirmed.Status)
	}

	// Đưa sang SHIPPING (ShipOrder thật là P4b — ở đây set trực tiếp để test Complete).
	if _, err := fx.pool.Exec(ctx,
		`UPDATE public.orders SET status='SHIPPING' WHERE id=$1::uuid`, id); err != nil {
		t.Fatalf("set SHIPPING: %v", err)
	}
	completed, err := fx.svc.Write().CompleteOrder(ctx, id)
	if err != nil {
		t.Fatalf("complete phải thành công: %v", err)
	}
	if completed.Status != domain.StatusCompleted.String() {
		t.Fatalf("sau complete = %q, want COMPLETED", completed.Status)
	}
}

// TestIntegration_CancelOrder: PENDING→CANCELLED.
func TestIntegration_CancelOrder(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	res, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, ""))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	cancelled, err := fx.svc.Write().CancelOrder(ctx, res.Order.ID)
	if err != nil {
		t.Fatalf("cancel phải thành công: %v", err)
	}
	if cancelled.Status != domain.StatusCancelled.String() {
		t.Fatalf("sau cancel = %q, want CANCELLED", cancelled.Status)
	}
}

// TestIntegration_InvalidTransition_Conflict: nhảy bước (PENDING→COMPLETED) → 409.
func TestIntegration_InvalidTransition_Conflict(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	res, err := fx.svc.Write().CreateOrder(ctx, twoLineInput(t, ""))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	// PENDING → COMPLETED (nhảy bước) phải Conflict.
	_, err = fx.svc.Write().CompleteOrder(ctx, res.Order.ID)
	if err == nil {
		t.Fatal("PENDING→COMPLETED phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("nhảy bước phải Conflict, got %v (%v)", apperr.KindOf(err), err)
	}
	// Huỷ đơn đã huỷ (CANCELLED→CANCELLED) cũng Conflict (terminal).
	if _, err := fx.svc.Write().CancelOrder(ctx, res.Order.ID); err != nil {
		t.Fatalf("cancel lần 1: %v", err)
	}
	if _, err := fx.svc.Write().CancelOrder(ctx, res.Order.ID); err == nil {
		t.Fatal("cancel đơn đã CANCELLED phải Conflict")
	}
}

// TestIntegration_Confirm_NotFound: id không tồn tại → NotFound.
func TestIntegration_Confirm_NotFound(t *testing.T) {
	fx := setup(t)
	_, err := fx.svc.Write().ConfirmOrder(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("confirm đơn không tồn tại phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("phải NotFound, got %v", apperr.KindOf(err))
	}
}

// ---- HTTP (envelope + authz orders.write) ----

func TestIntegration_HTTP_CreateOrder_AndAuthz(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)

	body := `{"order_type":"B2B","note":"qua http","lines":[
		{"product_id":100,"quantity":3,"uom":"Hộp","unit_price":"10000","discount":"5000"},
		{"product_id":200,"quantity":2,"uom":"Vỉ","unit_price":"20000","discount":"0"}]}`

	// Thiếu token → 401.
	if r := doPost(t, handler, "/v1/orders", body, "", ""); r.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", r.code, r.body)
	}
	// Sai quyền (orders.read) → 403.
	if r := doPost(t, handler, "/v1/orders", body, signToken(t, fx.key, "orders.read"), ""); r.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", r.code, r.body)
	}

	// Có quyền orders.write → 201, envelope {data:{order,lines}}, tiền ra chuỗi.
	tok := signToken(t, fx.key, "orders.write")
	r := doPost(t, handler, "/v1/orders", body, tok, "")
	if r.code != http.StatusCreated {
		t.Fatalf("create = %d, want 201; body=%s", r.code, r.body)
	}
	var env struct {
		Data struct {
			Order struct {
				ID          string `json:"id"`
				Code        string `json:"code"`
				Status      string `json:"status"`
				TotalAmount string `json:"total_amount"`
				FinalAmount string `json:"final_amount"`
			} `json:"order"`
			Lines []struct {
				LineTotal string `json:"line_total"`
			} `json:"lines"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(r.body), &env); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, r.body)
	}
	if env.Error != nil || env.Data.Order.Status != "PENDING" || env.Data.Order.FinalAmount != "65000" || len(env.Data.Lines) != 2 {
		t.Fatalf("envelope tạo đơn sai: %s", r.body)
	}
	orderID := env.Data.Order.ID

	// confirm → 200, status CONFIRMED.
	rc := doPost(t, handler, "/v1/orders/"+orderID+"/confirm", "", tok, "")
	if rc.code != http.StatusOK {
		t.Fatalf("confirm = %d, want 200; body=%s", rc.code, rc.body)
	}
	// cancel sau confirm → 200 (CONFIRMED→CANCELLED hợp lệ).
	rcancel := doPost(t, handler, "/v1/orders/"+orderID+"/cancel", "", tok, "")
	if rcancel.code != http.StatusOK {
		t.Fatalf("cancel = %d, want 200; body=%s", rcancel.code, rcancel.body)
	}
	// complete sau cancel (terminal) → 409.
	rcomplete := doPost(t, handler, "/v1/orders/"+orderID+"/complete", "", tok, "")
	if rcomplete.code != http.StatusConflict {
		t.Fatalf("complete đơn CANCELLED = %d, want 409; body=%s", rcomplete.code, rcomplete.body)
	}
}

// TestIntegration_HTTP_CreateOrder_IdempotentHeader: cùng Idempotency-Key header →
// trả cùng đơn.
func TestIntegration_HTTP_CreateOrder_IdempotentHeader(t *testing.T) {
	fx := setup(t)
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "orders.write")
	body := `{"order_type":"B2C","lines":[{"product_id":1,"quantity":1,"uom":"Cái","unit_price":"1000","discount":"0"}]}`

	r1 := doPost(t, handler, "/v1/orders", body, tok, "key-http-xyz")
	r2 := doPost(t, handler, "/v1/orders", body, tok, "key-http-xyz")
	if r1.code != http.StatusCreated || r2.code != http.StatusCreated {
		t.Fatalf("idempotent http codes = %d/%d, want 201/201", r1.code, r2.code)
	}
	code1 := extractCode(t, r1.body)
	code2 := extractCode(t, r2.body)
	if code1 == "" || code1 != code2 {
		t.Fatalf("idempotent http phải cùng mã đơn: %q vs %q", code1, code2)
	}
}

// ---- helpers GHI (tái dùng buildHTTP/signToken/httpResp từ orders_integration_test.go) ----

func doPost(t *testing.T, h http.Handler, path, body, bearer, idemKey string) httpResp {
	t.Helper()
	var rdr *strings.Reader
	if body == "" {
		rdr = strings.NewReader("")
	} else {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(http.MethodPost, path, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	if idemKey != "" {
		req.Header.Set("Idempotency-Key", idemKey)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return httpResp{code: rec.Code, body: rec.Body.String()}
}

func extractCode(t *testing.T, body string) string {
	t.Helper()
	var env struct {
		Data struct {
			Order struct {
				Code string `json:"code"`
			} `json:"order"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("decode code: %v (body=%s)", err, body)
	}
	return env.Data.Order.Code
}
