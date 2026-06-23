package migrations_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	migrations "github.com/Maneva-AI/namviet-backend/db"
)

// TestMigrations_UpThenDown apply toàn bộ migration rồi rollback sạch trên một
// Postgres thật (testcontainers). Bảo đảm cả nhánh Up lẫn Down của goose không
// lỗi và schema lên/xuống đối xứng.
func TestMigrations_UpThenDown(t *testing.T) {
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
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	// Up: apply tất cả.
	if err := migrations.Up(ctx, db); err != nil {
		t.Fatalf("Up: %v", err)
	}
	// Bảng app.idempotency_keys (00001) + identity (00002) + sổ kế toán (00003)
	// phải tồn tại sau Up.
	for _, tbl := range []string{
		"idempotency_keys", "users", "roles", "permissions", "role_permissions", "user_roles", "refresh_tokens",
		"accounting_periods", "journal_entries", "journal_entry_lines",
		"invoice_serials", "sales_invoices", "sales_invoice_lines",
		"order_idempotency",
		"finance_transaction_allocations",
		"purchase_orders", "purchase_order_items",
	} {
		if !tableExists(t, db, "app", tbl) {
			t.Fatalf("sau Up: bảng app.%s không tồn tại", tbl)
		}
	}
	// Sequence sinh mã đơn (00005) + mã PO (00007) phải tồn tại sau Up.
	if !sequenceExists(t, db, "app", "order_code_seq") {
		t.Fatal("sau Up: sequence app.order_code_seq không tồn tại")
	}
	if !sequenceExists(t, db, "app", "purchase_order_code_seq") {
		t.Fatal("sau Up: sequence app.purchase_order_code_seq không tồn tại")
	}

	// DownTo 0: gỡ sạch mọi migration.
	if err := migrations.DownTo(ctx, db, 0); err != nil {
		t.Fatalf("DownTo 0: %v", err)
	}
	for _, tbl := range []string{"idempotency_keys", "users", "refresh_tokens", "journal_entries", "journal_entry_lines", "accounting_periods", "sales_invoices", "sales_invoice_lines", "invoice_serials", "order_idempotency", "purchase_orders", "purchase_order_items"} {
		if tableExists(t, db, "app", tbl) {
			t.Fatalf("sau Down: bảng app.%s vẫn còn", tbl)
		}
	}
	if sequenceExists(t, db, "app", "order_code_seq") {
		t.Fatal("sau Down: sequence app.order_code_seq vẫn còn")
	}
	if sequenceExists(t, db, "app", "purchase_order_code_seq") {
		t.Fatal("sau Down: sequence app.purchase_order_code_seq vẫn còn")
	}

	// Up lại sau Down → idempotent, không lỗi.
	if err := migrations.Up(ctx, db); err != nil {
		t.Fatalf("Up lại sau Down: %v", err)
	}
	if !tableExists(t, db, "app", "idempotency_keys") {
		t.Fatalf("sau Up lần 2: bảng không tồn tại")
	}
}

func sequenceExists(t *testing.T, db *sql.DB, schema, name string) bool {
	t.Helper()
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.sequences
			WHERE sequence_schema = $1 AND sequence_name = $2
		)`, schema, name).Scan(&exists)
	if err != nil {
		t.Fatalf("query sequence exists: %v", err)
	}
	return exists
}

func tableExists(t *testing.T, db *sql.DB, schema, table string) bool {
	t.Helper()
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = $1 AND table_name = $2
		)`, schema, table).Scan(&exists)
	if err != nil {
		t.Fatalf("query table exists: %v", err)
	}
	return exists
}
