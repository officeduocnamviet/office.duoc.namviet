package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// envelope dùng để assert shape contract chung.
type envelope struct {
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func doGet(t *testing.T, h http.Handler, path string) (*httptest.ResponseRecorder, envelope) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var e envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
		t.Fatalf("body không phải JSON %q: %v", rec.Body.String(), err)
	}
	return rec, e
}

// okPinger / failPinger giả lập DB health.
type okPinger struct{}

func (okPinger) Ping(context.Context) error { return nil }

type failPinger struct{}

func (failPinger) Ping(context.Context) error { return errors.New("down") }

func TestHealth_OK_Envelope(t *testing.T) {
	// Pinger nil → bỏ qua ping, vẫn ok (giống Phase 0 khi chưa cấu hình DB).
	h := NewRouter(Deps{})
	rec, e := doGet(t, h, "/v1/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if e.Error != nil {
		t.Errorf("error phải null khi ok, got %+v", e.Error)
	}
	var data struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(e.Data, &data); err != nil {
		t.Fatalf("data decode lỗi: %v (raw=%s)", err, e.Data)
	}
	if data.Status != "ok" {
		t.Errorf("data.status = %q, want ok", data.Status)
	}
}

func TestHealth_OK_WithPinger(t *testing.T) {
	h := NewRouter(Deps{Pinger: okPinger{}})
	rec, e := doGet(t, h, "/v1/health")
	if rec.Code != http.StatusOK || e.Error != nil {
		t.Fatalf("want 200 ok envelope, got %d err=%+v body=%s", rec.Code, e.Error, "")
	}
}

func TestHealth_DBDown_503Envelope(t *testing.T) {
	h := NewRouter(Deps{Pinger: failPinger{}})
	rec, e := doGet(t, h, "/v1/health")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503; body=%s", rec.Code, rec.Body.String())
	}
	if string(e.Data) != "null" {
		t.Errorf("data phải null khi 503, got %s", e.Data)
	}
	if e.Error == nil || e.Error.Code != "internal" {
		t.Fatalf("error = %+v, want code=internal (503→internal)", e.Error)
	}
}

// TestDemoError_NotFoundEnvelope kiểm operation lỗi đi qua router thật.
func TestDemoError_NotFoundEnvelope(t *testing.T) {
	h := NewRouter(Deps{})
	rec, e := doGet(t, h, "/v1/_demo/not-found")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	if string(e.Data) != "null" || e.Error == nil || e.Error.Code != "not_found" {
		t.Fatalf("want not_found envelope, got data=%s err=%+v", e.Data, e.Error)
	}
}

func TestDemoValidation_422Envelope(t *testing.T) {
	h := NewRouter(Deps{})
	req := httptest.NewRequest(http.MethodPost, "/v1/_demo/echo",
		strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
	var e envelope
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if e.Error == nil || e.Error.Code != "validation_error" {
		t.Fatalf("want validation_error, got %+v; body=%s", e.Error, rec.Body.String())
	}
}
