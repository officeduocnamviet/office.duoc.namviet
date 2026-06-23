package idempotency_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	migrations "github.com/Maneva-AI/namviet-backend/db"
	"github.com/Maneva-AI/namviet-backend/internal/platform/idempotency"
)

// setupDB khởi container postgres:18, apply migrations bằng goose, trả pool +
// cleanup. Skip khi -short.
func setupDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("skip integration in -short")
	}
	ctx := context.Background()
	ctr, err := tcpg.Run(ctx, "postgres:18",
		tcpg.WithDatabase("test"), tcpg.WithUsername("test"), tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Apply migrations qua goose dùng database/sql + pgx stdlib.
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open sql: %v", err)
	}
	if err := migrations.Up(ctx, sqlDB); err != nil {
		t.Fatalf("migrate up: %v", err)
	}
	_ = sqlDB.Close()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestPgxStore_FullLifecycle(t *testing.T) {
	ctx := context.Background()
	pool := setupDB(t)
	store := idempotency.NewPgxStore(ctx, pool)

	const key = "order-create-123"

	// Chưa có key → Get trả không tồn tại.
	if _, ok, err := store.Get(key); err != nil || ok {
		t.Fatalf("Get trước Begin: ok=%v err=%v, want (false,nil)", ok, err)
	}

	// Begin → tạo bản ghi in_progress.
	if err := store.Begin(key, "hash-1"); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	rec, ok, err := store.Get(key)
	if err != nil || !ok {
		t.Fatalf("Get sau Begin: ok=%v err=%v", ok, err)
	}
	if rec.State != "in_progress" {
		t.Errorf("state = %q, want in_progress", rec.State)
	}

	// Complete → set done + lưu status/body.
	body := []byte(`{"data":{"id":1},"error":null}`)
	if err := store.Complete(key, 201, body); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	rec, ok, err = store.Get(key)
	if err != nil || !ok {
		t.Fatalf("Get sau Complete: ok=%v err=%v", ok, err)
	}
	if rec.State != "done" {
		t.Errorf("state = %q, want done", rec.State)
	}
	if rec.Status != 201 {
		t.Errorf("status = %d, want 201", rec.Status)
	}
	// response_body là jsonb → Postgres chuẩn hoá whitespace; so sánh ngữ nghĩa
	// JSON thay vì byte-exact.
	if !jsonEqual(t, rec.Body, body) {
		t.Errorf("body = %s, want tương đương %s", rec.Body, body)
	}
}

// jsonEqual so sánh hai JSON theo ngữ nghĩa (bỏ qua khác biệt whitespace/thứ tự).
func jsonEqual(t *testing.T, a, b []byte) bool {
	t.Helper()
	var av, bv any
	if err := json.Unmarshal(a, &av); err != nil {
		t.Fatalf("unmarshal a: %v", err)
	}
	if err := json.Unmarshal(b, &bv); err != nil {
		t.Fatalf("unmarshal b: %v", err)
	}
	return reflect.DeepEqual(av, bv)
}

func TestPgxStore_BeginIsIdempotent(t *testing.T) {
	ctx := context.Background()
	pool := setupDB(t)
	store := idempotency.NewPgxStore(ctx, pool)

	const key = "replay-key"
	if err := store.Begin(key, "h1"); err != nil {
		t.Fatalf("Begin#1: %v", err)
	}
	if err := store.Complete(key, 200, []byte(`{"data":{"ok":true},"error":null}`)); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	// Begin lần 2 trên key đã done KHÔNG được ghi đè (ON CONFLICT DO NOTHING).
	if err := store.Begin(key, "h2"); err != nil {
		t.Fatalf("Begin#2: %v", err)
	}
	rec, ok, err := store.Get(key)
	if err != nil || !ok {
		t.Fatalf("Get: ok=%v err=%v", ok, err)
	}
	if rec.State != "done" || rec.Status != 200 {
		t.Errorf("Begin lần 2 đã ghi đè record done: %+v", rec)
	}
}
