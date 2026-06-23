package finance_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/finance/app"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
)

// recordOutInTx chạy RecordPaymentOut trong MỘT tx (mô phỏng purchasing gộp atomic).
func recordOutInTx(ctx context.Context, pool *pgxpool.Pool, rec *app.Recorder, p app.RecordPaymentOutParams) (domain.Payment, bool, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return domain.Payment{}, false, err
	}
	got, created, rerr := rec.RecordPaymentOut(ctx, tx, p)
	if rerr != nil {
		_ = tx.Rollback(ctx)
		return domain.Payment{}, false, rerr
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Payment{}, false, err
	}
	return got, created, nil
}

func countOut(t *testing.T, pool *pgxpool.Pool, poCode string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM public.finance_transactions WHERE ref_type='purchase_order' AND ref_id=$1 AND flow='out' AND deleted_at IS NULL`,
		poCode).Scan(&n); err != nil {
		t.Fatalf("countOut: %v", err)
	}
	return n
}

// TestIntegration_RecordPaymentOut_Happy: phiếu CHI thủ công → 1 row (flow='out',
// ref_type='purchase_order'), trigger TRỪ số dư quỹ MỘT lần. Quỹ 2 khởi 1000 → 1000-300.
func TestIntegration_RecordPaymentOut_Happy(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	p := app.RecordPaymentOutParams{
		RecordPaymentOut: domain.RecordPaymentOut{
			POCode: "PO00000001", Amount: mustMoney(t, "300"), FundAccountID: 2, BookType: domain.BookBoth,
		},
		IdemKey: "pay-out-1",
	}
	got, created, err := recordOutInTx(ctx, pool, mod.Recorder(), p)
	if err != nil || !created {
		t.Fatalf("ghi phiếu chi: err=%v created=%v", err, created)
	}
	if got.Flow != "out" || got.RefType != "purchase_order" || got.RefID != "PO00000001" {
		t.Fatalf("phiếu chi sai: flow=%q ref_type=%q ref_id=%q", got.Flow, got.RefType, got.RefID)
	}
	if got := fundBalance(t, pool, 2); !decEqual(t, got, "700") {
		t.Fatalf("số dư quỹ 2 = %s, want 700 (1000-300)", got)
	}
}

// TestIntegration_RecordPaymentOut_Idempotent: replay cùng idem key → 1 row,
// created=false, KHÔNG trừ tiền 2 lần.
func TestIntegration_RecordPaymentOut_Idempotent(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	p := app.RecordPaymentOutParams{
		RecordPaymentOut: domain.RecordPaymentOut{
			POCode: "PO00000002", Amount: mustMoney(t, "200"), FundAccountID: 2, BookType: domain.BookInternal,
		},
		IdemKey: "pay-out-dup",
	}
	first, c1, err := recordOutInTx(ctx, pool, mod.Recorder(), p)
	if err != nil || !c1 {
		t.Fatalf("lần 1: err=%v created=%v", err, c1)
	}
	second, c2, err := recordOutInTx(ctx, pool, mod.Recorder(), p)
	if err != nil {
		t.Fatalf("lần 2: %v", err)
	}
	if c2 {
		t.Fatal("lần 2 phải created=false (idempotent hit)")
	}
	if first.ID != second.ID {
		t.Fatalf("idempotent: id khác %d != %d", first.ID, second.ID)
	}
	if n := countOut(t, pool, "PO00000002"); n != 1 {
		t.Fatalf("phiếu chi PO00000002 = %d, want 1", n)
	}
	// 1000 - 200 = 800 (KHÔNG 600).
	if got := fundBalance(t, pool, 2); !decEqual(t, got, "800") {
		t.Fatalf("số dư = %s, want 800 (KHÔNG trừ đôi)", got)
	}
}
