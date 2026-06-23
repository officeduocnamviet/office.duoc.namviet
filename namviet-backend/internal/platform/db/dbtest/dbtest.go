// Package dbtest là harness testcontainers DÙNG CHUNG cho integration test của
// mọi module (ARCHITECTURE.md §11). Nó spin một Postgres 18 thật, chạy toàn bộ
// migration goose (Up), rồi trả về một *pgxpool.Pool đã đăng ký codec decimal
// cùng hàm cleanup. Nhờ tập trung ở một chỗ, mọi test viết GIỐNG NHAU và việc
// đổi version Postgres/chiến lược migrate chỉ sửa một nơi.
//
// Package này chỉ build trong ngữ cảnh test (import testing) nhưng đặt ở thư mục
// thường để các package _test khác import được.
package dbtest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // driver database/sql cho goose
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	migrations "github.com/Maneva-AI/namviet-backend/db"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db"
)

// NewPool spin Postgres 18, apply migration, trả pool + cleanup. Test phải gọi
// t.Cleanup(cleanup) (hoặc defer) để terminate container. Tự skip khi -short.
func NewPool(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	if testing.Short() {
		t.Skip("skip integration trong -short")
	}

	ctx := context.Background()
	ctr, err := tcpg.Run(ctx, "postgres:18",
		tcpg.WithDatabase("test"),
		tcpg.WithUsername("test"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("dbtest: start container: %v", err)
	}

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = ctr.Terminate(ctx)
		t.Fatalf("dbtest: connection string: %v", err)
	}

	// Migrate qua database/sql (driver pgx stdlib) như goose yêu cầu.
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		_ = ctr.Terminate(ctx)
		t.Fatalf("dbtest: sql.Open: %v", err)
	}
	if err := migrations.Up(ctx, sqlDB); err != nil {
		_ = sqlDB.Close()
		_ = ctr.Terminate(ctx)
		t.Fatalf("dbtest: migrate up: %v", err)
	}
	// Vật chất hoá schema THAM CHIẾU public.* (ADR 0001) để integration test của
	// module strangler-fig (catalog...) thấy bảng kế thừa trong container trống.
	// Idempotent (IF NOT EXISTS); module chỉ-app như identity không dùng bảng này
	// nhưng tạo thêm bảng rỗng vô hại.
	if err := migrations.ApplyPublicSchema(ctx, sqlDB); err != nil {
		_ = sqlDB.Close()
		_ = ctr.Terminate(ctx)
		t.Fatalf("dbtest: apply public schema: %v", err)
	}
	_ = sqlDB.Close()

	// Pool runtime (có codec decimal) cho repo.
	pool, err := db.Connect(ctx, dsn)
	if err != nil {
		_ = ctr.Terminate(ctx)
		t.Fatalf("dbtest: connect pool: %v", err)
	}

	cleanup := func() {
		pool.Close()
		_ = ctr.Terminate(context.Background())
	}
	return pool, cleanup
}
