package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/db"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
)

func TestWithinTx_CommitsOnSuccess(t *testing.T) {
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	ctx := context.Background()

	err := db.WithinTx(ctx, pool, func(tx pgx.Tx) error {
		_, e := tx.Exec(ctx,
			`INSERT INTO app.idempotency_keys (key, request_hash, state) VALUES ($1,$2,'done')`,
			"k-commit", "h")
		return e
	})
	if err != nil {
		t.Fatalf("WithinTx: %v", err)
	}

	if !rowExists(t, pool, "k-commit") {
		t.Fatal("commit không lưu hàng")
	}
}

func TestWithinTx_RollsBackOnError(t *testing.T) {
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	ctx := context.Background()

	sentinel := errors.New("nghiệp vụ thất bại")
	err := db.WithinTx(ctx, pool, func(tx pgx.Tx) error {
		if _, e := tx.Exec(ctx,
			`INSERT INTO app.idempotency_keys (key, request_hash, state) VALUES ($1,$2,'done')`,
			"k-rollback", "h"); e != nil {
			return e
		}
		return sentinel // ép rollback
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, muốn sentinel", err)
	}

	if rowExists(t, pool, "k-rollback") {
		t.Fatal("rollback nhưng hàng vẫn còn")
	}
}

func rowExists(t *testing.T, pool interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, key string) bool {
	t.Helper()
	var exists bool
	err := pool.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM app.idempotency_keys WHERE key=$1)`, key).Scan(&exists)
	if err != nil {
		t.Fatalf("query exists: %v", err)
	}
	return exists
}
