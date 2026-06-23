package idempotency

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type memStore struct{ m map[string]Record }

func (s *memStore) Get(key string) (Record, bool, error) {
	r, ok := s.m[key]
	return r, ok, nil
}

func (s *memStore) Begin(key, hash string) error {
	s.m[key] = Record{State: "in_progress"}
	return nil
}

func (s *memStore) Complete(key string, status int, body []byte) error {
	s.m[key] = Record{State: "done", Status: status, Body: body}
	return nil
}

func TestMiddleware_ReplaysCompleted(t *testing.T) {
	store := &memStore{m: map[string]Record{
		"abc": {State: "done", Status: 201, Body: []byte(`{"data":{"id":1},"error":null}`)},
	}}
	calls := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(200)
	})
	h := Middleware(store)(next)

	req := httptest.NewRequest(http.MethodPost, "/v1/orders", nil)
	req.Header.Set("Idempotency-Key", "abc")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if calls != 0 {
		t.Errorf("next should not run for completed key; calls=%d", calls)
	}
	if rec.Code != 201 {
		t.Errorf("status = %d, want replayed 201", rec.Code)
	}
	if rec.Body.String() != `{"data":{"id":1},"error":null}` {
		t.Errorf("body = %q, want replayed body", rec.Body.String())
	}
}

func TestMiddleware_PassesThroughWithoutKey(t *testing.T) {
	store := &memStore{m: map[string]Record{}}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(200)
	})
	h := Middleware(store)(next)
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if !called {
		t.Error("next should run when no Idempotency-Key")
	}
}

func TestMiddleware_FirstCallRunsAndCompletes(t *testing.T) {
	store := &memStore{m: map[string]Record{}}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"data":{"id":9},"error":null}`))
	})
	h := Middleware(store)(next)

	req := httptest.NewRequest(http.MethodPost, "/v1/orders", nil)
	req.Header.Set("Idempotency-Key", "key-1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	got, ok, _ := store.Get("key-1")
	if !ok || got.State != "done" {
		t.Fatalf("store should hold done record, got %+v ok=%v", got, ok)
	}
	if got.Status != 201 || string(got.Body) != `{"data":{"id":9},"error":null}` {
		t.Errorf("completed record = %+v", got)
	}
}
