package humax_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// buildAPI dựng một huma.API thật trên chi router với config envelope của humax,
// đăng ký 3 operation đại diện cho 3 nhánh contract cần kiểm.
func buildAPI(t *testing.T) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "Test API", "1.0.0")

	// (1) success: body tuỳ ý → phải bị bọc {data, error:null}.
	type okBody struct {
		Status string `json:"status"`
	}
	type okOut struct {
		Body okBody
	}
	huma.Register(api, huma.Operation{
		OperationID: "get-ok",
		Method:      http.MethodGet,
		Path:        "/ok",
	}, func(_ context.Context, _ *struct{}) (*okOut, error) {
		return &okOut{Body: okBody{Status: "ok"}}, nil
	})

	// (2) error: handler trả 404 → phải thành {data:null, error:{code:not_found}}.
	huma.Register(api, huma.Operation{
		OperationID: "get-missing",
		Method:      http.MethodGet,
		Path:        "/missing",
	}, func(_ context.Context, _ *struct{}) (*okOut, error) {
		return nil, huma.Error404NotFound("đơn không tồn tại")
	})

	// (3) validation: input có field bắt buộc; thiếu → 422 envelope.
	type createIn struct {
		Body struct {
			Name string `json:"name" minLength:"1"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID:   "create-thing",
		Method:        http.MethodPost,
		Path:          "/things",
		DefaultStatus: http.StatusCreated,
	}, func(_ context.Context, in *createIn) (*okOut, error) {
		return &okOut{Body: okBody{Status: in.Body.Name}}, nil
	})

	return r
}

// envelope là shape ta yêu cầu FE thấy.
type envelope struct {
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func decode(t *testing.T, body []byte) envelope {
	t.Helper()
	var e envelope
	if err := json.Unmarshal(body, &e); err != nil {
		t.Fatalf("body không phải JSON hợp lệ %q: %v", string(body), err)
	}
	return e
}

func TestSuccess_WrappedEnvelope(t *testing.T) {
	h := buildAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	e := decode(t, rec.Body.Bytes())
	if e.Error != nil {
		t.Errorf("error phải null khi success, got %+v", e.Error)
	}
	var data struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(e.Data, &data); err != nil {
		t.Fatalf("data không decode được: %v (raw=%s)", err, e.Data)
	}
	if data.Status != "ok" {
		t.Errorf("data.status = %q, want ok", data.Status)
	}
	// Không double-wrap: data KHÔNG được chứa lại key "data".
	var nested map[string]any
	_ = json.Unmarshal(e.Data, &nested)
	if _, dup := nested["data"]; dup {
		t.Errorf("phát hiện double-wrap: data chứa lại 'data': %s", e.Data)
	}
}

func TestError404_Envelope(t *testing.T) {
	h := buildAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	// Contract FE: error trả về application/json (KHÔNG problem+json).
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	e := decode(t, rec.Body.Bytes())
	if string(e.Data) != "null" {
		t.Errorf("data phải null khi lỗi, got %s", e.Data)
	}
	if e.Error == nil {
		t.Fatalf("error phải khác null khi 404")
	}
	if e.Error.Code != "not_found" {
		t.Errorf("error.code = %q, want not_found", e.Error.Code)
	}
	if e.Error.Message != "đơn không tồn tại" {
		t.Errorf("error.message = %q", e.Error.Message)
	}
}

func TestValidation422_Envelope(t *testing.T) {
	h := buildAPI(t)
	// Body JSON hợp lệ nhưng vi phạm schema: name rỗng (< minLength 1) → 422.
	req := httptest.NewRequest(http.MethodPost, "/things",
		strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
	e := decode(t, rec.Body.Bytes())
	if string(e.Data) != "null" {
		t.Errorf("data phải null khi validation lỗi, got %s", e.Data)
	}
	if e.Error == nil || e.Error.Code != "validation_error" {
		t.Fatalf("error.code = %+v, want validation_error; body=%s", e.Error, rec.Body.String())
	}
}

// TestStatusCodeToCode kiểm bảng map status→code đầy đủ theo spec.
func TestStatusCodeToCode(t *testing.T) {
	cases := map[int]string{
		http.StatusBadRequest:          "bad_request",
		http.StatusUnauthorized:        "unauthorized",
		http.StatusForbidden:           "forbidden",
		http.StatusNotFound:            "not_found",
		http.StatusConflict:            "conflict",
		http.StatusUnprocessableEntity: "validation_error",
		http.StatusInternalServerError: "internal",
		http.StatusBadGateway:          "internal", // fallback cho 5xx khác
	}
	for status, want := range cases {
		if got := humax.CodeForStatus(status); got != want {
			t.Errorf("CodeForStatus(%d) = %q, want %q", status, got, want)
		}
	}
}
