// Package server sở hữu vòng đời HTTP và wiring route của edge. Tầng HTTP chính
// thức là Huma v2 (code-first OpenAPI 3.1) gắn lên chi router; mọi response giữ
// envelope {data,error} qua package httpx/humax.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Pinger là port kiểm tra sức khoẻ phụ thuộc (thường là *pgxpool.Pool). Tách
// interface để test health không cần DB thật.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Deps gom mọi phụ thuộc runtime — wiring tường minh, dễ test.
type Deps struct {
	// Pool là pool Postgres runtime (có thể nil khi chưa cấu hình DB).
	Pool *pgxpool.Pool
	// Pinger override nguồn health-check (mặc định lấy từ Pool nếu nil). Hữu ích
	// cho test inject pinger lỗi mà không cần container.
	Pinger Pinger
	// Modules là các hàm đăng ký route của bounded context (vd identity). Edge
	// (cmd/api) wiring closure cho mỗi module rồi truyền vào đây — server KHÔNG
	// phụ thuộc ngược lên module nào, giữ đúng chiều phụ thuộc.
	Modules []func(huma.API)
}

// pinger chọn nguồn ping: ưu tiên Deps.Pinger, sau đó Pool, cuối cùng nil.
func (d Deps) pinger() Pinger {
	if d.Pinger != nil {
		return d.Pinger
	}
	if d.Pool != nil {
		return d.Pool
	}
	return nil
}

// NewRouter dựng chi router + middleware nền, gắn huma.API và đăng ký operations.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	api := humax.New(r, "Nam Viet Backend API", "1.0.0")
	registerRoutes(api, d)
	return r
}

// API dựng huma.API trên một chi router cho trước rồi đăng ký toàn bộ operations.
// Tách khỏi NewRouter để cmd/api có thể introspect OpenAPI (dump spec) mà không
// cần khởi động server.
func API(r chi.Router, d Deps) huma.API {
	api := humax.New(r, "Nam Viet Backend API", "1.0.0")
	registerRoutes(api, d)
	return api
}

// NewHTTPServer bọc router với timeout an toàn.
func NewHTTPServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// Run khởi động server và shutdown graceful khi ctx hủy.
func Run(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
