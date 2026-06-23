// Package logging cấu hình slog (JSON ra stdout) và tùy chọn fanout sang OTel.
package logging

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// New trả *slog.Logger ghi JSON ra stdout; nếu otelEnabled, gắn thêm handler OTel
// (log có trace_id khi dùng logger.InfoContext(ctx, ...)).
func New(serviceName string, otelEnabled bool) *slog.Logger {
	stdout := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	if !otelEnabled {
		return slog.New(stdout)
	}
	otelHandler := otelslog.NewHandler(serviceName)
	return slog.New(fanout{stdout, otelHandler})
}

// fanout gửi mỗi record tới nhiều handler.
type fanout []slog.Handler

func (f fanout) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range f {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (f fanout) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range f {
		_ = h.Handle(ctx, r.Clone())
	}
	return nil
}

func (f fanout) WithAttrs(a []slog.Attr) slog.Handler {
	out := make(fanout, len(f))
	for i, h := range f {
		out[i] = h.WithAttrs(a)
	}
	return out
}

func (f fanout) WithGroup(name string) slog.Handler {
	out := make(fanout, len(f))
	for i, h := range f {
		out[i] = h.WithGroup(name)
	}
	return out
}
