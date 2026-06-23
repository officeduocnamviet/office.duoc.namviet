// Package migrations nhúng các file goose .sql và cung cấp helper apply/rollback
// chạy bằng thư viện pressly/goose (không cần CLI). Dùng cho cả runtime (nếu cần
// auto-migrate ở môi trường dev) lẫn integration test (testcontainers).
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var fsys embed.FS

// schemaFS nhúng các SCHEMA THAM CHIẾU public.* (ADR 0001 strangler-fig). Đây
// KHÔNG phải goose migration — chỉ DDL mô tả bảng kế thừa Supabase để (1) sqlc
// type-check và (2) integration test materialize bảng trong testcontainers (DB
// test trống). TUYỆT ĐỐI KHÔNG apply lên prod.
//
//go:embed schema/*.sql
var schemaFS embed.FS

// FS trả filesystem nhúng chứa thư mục migrations (dùng cho kiểm tra trong test).
func FS() embed.FS { return fsys }

// ApplyPublicSchema chạy toàn bộ DDL trong db/schema/*.sql (các bảng public.*
// tham chiếu) lên db. Dùng CHO INTEGRATION TEST sau khi đã goose Up: harness
// dbtest gọi hàm này để vật chất hoá bảng public.* trong container test trống.
// Các file dùng IF NOT EXISTS nên idempotent. KHÔNG gọi ở runtime prod.
func ApplyPublicSchema(ctx context.Context, db *sql.DB) error {
	entries, err := fs.ReadDir(schemaFS, "schema")
	if err != nil {
		return fmt.Errorf("read dir schema: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := schemaFS.ReadFile("schema/" + e.Name())
		if err != nil {
			return fmt.Errorf("read schema %s: %w", e.Name(), err)
		}
		if _, err := db.ExecContext(ctx, string(b)); err != nil {
			return fmt.Errorf("apply schema %s: %w", e.Name(), err)
		}
	}
	return nil
}

// fsSub trả filesystem đã loại prefix "migrations/" để goose thấy các .sql ở root.
func fsSub() (fs.FS, error) {
	sub, err := fs.Sub(fsys, "migrations")
	if err != nil {
		return nil, fmt.Errorf("fs.Sub migrations: %w", err)
	}
	return sub, nil
}

// Up apply toàn bộ migration còn thiếu lên db (database/sql, driver pgx stdlib).
func Up(ctx context.Context, db *sql.DB) error {
	p, err := provider(db)
	if err != nil {
		return err
	}
	if _, err := p.Up(ctx); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

// Down rollback migration mới nhất (một bước).
func Down(ctx context.Context, db *sql.DB) error {
	p, err := provider(db)
	if err != nil {
		return err
	}
	if _, err := p.Down(ctx); err != nil {
		return fmt.Errorf("goose down: %w", err)
	}
	return nil
}

// DownTo rollback về version cho trước (0 = gỡ sạch).
func DownTo(ctx context.Context, db *sql.DB, version int64) error {
	p, err := provider(db)
	if err != nil {
		return err
	}
	if _, err := p.DownTo(ctx, version); err != nil {
		return fmt.Errorf("goose down-to %d: %w", version, err)
	}
	return nil
}

func provider(db *sql.DB) (*goose.Provider, error) {
	sub, err := fsSub()
	if err != nil {
		return nil, err
	}
	p, err := goose.NewProvider(goose.DialectPostgres, db, sub)
	if err != nil {
		return nil, fmt.Errorf("new goose provider: %w", err)
	}
	return p, nil
}
