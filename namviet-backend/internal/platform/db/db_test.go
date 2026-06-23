package db

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnectAndPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration in -short")
	}
	ctx := context.Background()
	// Wait CHUẨN cho Postgres: ForListeningPort fire SỚM (port mở trước khi DB sẵn
	// sàng nhận query → Ping flaky). Đợi log "ready to accept connections" LẦN 2
	// (postgres log lần 1 lúc init bootstrap, lần 2 khi thật sự sẵn sàng) + port.
	ctr, err := tcpg.Run(ctx, "postgres:18",
		tcpg.WithDatabase("test"), tcpg.WithUsername("test"), tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithStartupTimeoutDefault(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	defer func() { _ = ctr.Terminate(ctx) }()

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	pool, err := Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
}
