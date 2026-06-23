// Package idempotency chống double-POST cho các thao tác chuyển tiền: nếu request
// mang Idempotency-Key đã hoàn tất ('done') thì replay nguyên response; nếu chưa
// thì chạy handler và lưu lại kết quả. Store tách qua interface để unit-test không
// cần DB.
package idempotency

import (
	"bytes"
	"net/http"
)

// Record là một bản ghi idempotency đã lưu.
type Record struct {
	State  string // in_progress | done
	Status int
	Body   []byte
}

// Store là kho lưu trạng thái idempotency. Triển khai in-memory cho test, hoặc
// pgxStore cho runtime (xem pgxstore.go).
type Store interface {
	Get(key string) (Record, bool, error)
	Begin(key, hash string) error
	Complete(key string, status int, body []byte) error
}

// capture chặn status + body ghi ra để lưu vào store.
type capture struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (c *capture) WriteHeader(s int) {
	c.status = s
	c.ResponseWriter.WriteHeader(s)
}

func (c *capture) Write(b []byte) (int, error) {
	c.buf.Write(b)
	return c.ResponseWriter.Write(b)
}

// Middleware: key rỗng → đi thẳng; key đã 'done' → replay status+body; ngược lại
// chạy next rồi Complete.
func Middleware(store Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			if rec, ok, _ := store.Get(key); ok && rec.State == "done" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(rec.Status)
				_, _ = w.Write(rec.Body)
				return
			}
			_ = store.Begin(key, "") // hash request có thể bổ sung sau
			c := &capture{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(c, r)
			_ = store.Complete(key, c.status, c.buf.Bytes())
		})
	}
}
