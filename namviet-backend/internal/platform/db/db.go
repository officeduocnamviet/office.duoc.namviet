// Package db sở hữu pool kết nối Postgres và codec NUMERIC<->shopspring/decimal.
package db

import (
	"context"
	"fmt"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect tạo pgxpool và đăng ký codec NUMERIC<->shopspring/decimal cho mọi conn
// thông qua AfterConnect. Pgx v5 dùng signature
// `AfterConnect func(context.Context, *pgx.Conn) error`, nên lấy type map trực
// tiếp từ `conn.TypeMap()`.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	cfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	// MaxConns để mặc định pgx (cores*4); tinh chỉnh sau theo tải.
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}
	return pool, nil
}
