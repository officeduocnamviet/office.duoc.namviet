package httpx

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestWriteData(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteData(rec, 200, map[string]any{"id": 7})

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var got struct {
		Data  map[string]any `json:"data"`
		Error *APIError      `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.Error != nil {
		t.Errorf("error should be null, got %+v", got.Error)
	}
	if got.Data["id"].(float64) != 7 {
		t.Errorf("data.id = %v, want 7", got.Data["id"])
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, 404, "not_found", "đơn không tồn tại")

	var got struct {
		Data  any       `json:"data"`
		Error *APIError `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.Data != nil {
		t.Errorf("data should be null on error, got %v", got.Data)
	}
	if got.Error == nil || got.Error.Code != "not_found" {
		t.Fatalf("error.code = %+v, want not_found", got.Error)
	}
	if got.Error.Message != "đơn không tồn tại" {
		t.Errorf("error.message = %q", got.Error.Message)
	}
}
