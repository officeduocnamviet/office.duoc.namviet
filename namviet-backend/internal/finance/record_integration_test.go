package finance_test

import (
	"context"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/finance"
	"github.com/Maneva-AI/namviet-backend/internal/finance/app"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
)

// setup spin Postgres test (đã apply schema tham chiếu public.* gồm
// fund_accounts/finance_transactions + trigger số dư), seed 1 quỹ + 1 đơn.
func setup(t *testing.T) (*finance.Module, *pgxpool.Pool) {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	seed(t, pool)
	return finance.NewModule(pool), pool
}

// seed: quỹ tiền mặt id=1 balance khởi 0; 1 đơn HD001 (chỉ để có code đối chiếu —
// finance ghi ref_id=code, không FK cứng).
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.fund_accounts (id, name, type, balance) VALUES (1, 'Quỹ tiền mặt', 'cash', 0)`,
		`INSERT INTO public.fund_accounts (id, name, type, balance) VALUES (2, 'Ngân hàng Timo', 'bank', 1000)`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("seed: %v\nSQL: %s", err, s)
		}
	}
}

func mustMoney(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return m
}

// recordInTx chạy RecordPaymentIn trong MỘT tx (mô phỏng orders/POS gộp atomic) +
// commit. Rollback nếu lỗi.
func recordInTx(ctx context.Context, pool *pgxpool.Pool, rec *app.Recorder, p app.RecordPaymentInParams) (domain.Payment, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return domain.Payment{}, err
	}
	got, _, rerr := rec.RecordPaymentIn(ctx, tx, p)
	if rerr != nil {
		_ = tx.Rollback(ctx)
		return domain.Payment{}, rerr
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Payment{}, err
	}
	return got, nil
}

// fundBalance đọc số dư quỹ (do trigger duy trì).
func fundBalance(t *testing.T, pool *pgxpool.Pool, id int64) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(),
		`SELECT balance::text FROM public.fund_accounts WHERE id=$1`, id).Scan(&s); err != nil {
		t.Fatalf("fundBalance: %v", err)
	}
	return s
}

// countAlive đếm phiếu còn sống của một đơn.
func countAlive(t *testing.T, pool *pgxpool.Pool, orderCode string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM public.finance_transactions WHERE ref_type='order' AND ref_id=$1 AND deleted_at IS NULL`,
		orderCode).Scan(&n); err != nil {
		t.Fatalf("countAlive: %v", err)
	}
	return n
}

func decEqual(t *testing.T, a, b string) bool {
	t.Helper()
	ma, ea := money.FromString(a)
	mb, eb := money.FromString(b)
	if ea != nil || eb != nil {
		return a == b
	}
	return ma.Equal(mb)
}

// TestIntegration_RecordPaymentIn_Happy: ghi phiếu thu thủ công → 1 row đúng
// (flow='in', book_type, ref_id=code, amount, status='completed'); trigger cộng số
// dư quỹ ĐÚNG MỘT LẦN (Go KHÔNG tự update balance).
func TestIntegration_RecordPaymentIn_Happy(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()

	p := app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode:     "HD001",
			Amount:        mustMoney(t, "150000"),
			FundAccountID: 1,
			BookType:      domain.BookBoth,
		},
		IdemKey: "manual-key-1",
	}
	got, err := recordInTx(ctx, pool, mod.Recorder(), p)
	if err != nil {
		t.Fatalf("ghi phiếu happy phải thành công: %v", err)
	}
	if got.ID == 0 {
		t.Fatal("id phiếu phải do DB sinh (>0)")
	}
	if got.Flow != "in" || got.Status != "completed" || got.RefType != "order" || got.RefID != "HD001" {
		t.Fatalf("phiếu sai: flow=%q status=%q ref_type=%q ref_id=%q", got.Flow, got.Status, got.RefType, got.RefID)
	}
	if got.BookType != domain.BookBoth {
		t.Fatalf("book_type = %q, want BOTH", got.BookType)
	}
	if !got.Amount.Equal(mustMoney(t, "150000")) {
		t.Fatalf("amount = %s, want 150000", got.Amount)
	}
	// Trigger cộng số dư đúng 1 lần (KHÔNG double): 0 + 150000 = 150000.
	if got := fundBalance(t, pool, 1); !decEqual(t, got, "150000") {
		t.Fatalf("balance quỹ = %s, want 150000 (trigger cộng 1 lần)", got)
	}
	if n := countAlive(t, pool, "HD001"); n != 1 {
		t.Fatalf("số phiếu HD001 = %d, want 1", n)
	}
}

// TestIntegration_RecordPaymentIn_IdempotentManual: gọi 2 lần CÙNG idem key →
// CHỈ 1 row, balance KHÔNG cộng đôi, lần 2 trả phiếu cũ (cùng id).
func TestIntegration_RecordPaymentIn_IdempotentManual(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	p := app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode:     "HD002",
			Amount:        mustMoney(t, "200000"),
			FundAccountID: 1,
			BookType:      domain.BookInternal,
		},
		IdemKey: "manual-dup",
	}
	first, err := recordInTx(ctx, pool, mod.Recorder(), p)
	if err != nil {
		t.Fatalf("lần 1: %v", err)
	}
	second, err := recordInTx(ctx, pool, mod.Recorder(), p)
	if err != nil {
		t.Fatalf("lần 2 (trùng key) phải no-op thành công: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("idempotent: id khác nhau %d != %d (đã tạo phiếu mới!)", first.ID, second.ID)
	}
	if n := countAlive(t, pool, "HD002"); n != 1 {
		t.Fatalf("số phiếu HD002 = %d, want 1 (không tạo trùng)", n)
	}
	// Balance KHÔNG cộng đôi: 0 + 200000 = 200000 (không 400000).
	if got := fundBalance(t, pool, 1); !decEqual(t, got, "200000") {
		t.Fatalf("balance = %s, want 200000 (KHÔNG cộng đôi)", got)
	}
}

// TestIntegration_RecordPaymentIn_IdempotentWebhook: 2 phiếu CÙNG bank_reference_id
// (webhook replay) → CHỈ 1 row, balance không cộng đôi. Quỹ 2 khởi 1000 → 1000+500.
func TestIntegration_RecordPaymentIn_IdempotentWebhook(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	bankRef := "FT2026XYZ123"
	p := app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode:     "HD003",
			Amount:        mustMoney(t, "500"),
			FundAccountID: 2,
			BookType:      domain.BookBoth,
			BankRef:       &bankRef,
		},
	}
	first, err := recordInTx(ctx, pool, mod.Recorder(), p)
	if err != nil {
		t.Fatalf("webhook lần 1: %v", err)
	}
	second, err := recordInTx(ctx, pool, mod.Recorder(), p)
	if err != nil {
		t.Fatalf("webhook replay phải no-op: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("webhook idempotent: id khác %d != %d", first.ID, second.ID)
	}
	if n := countAlive(t, pool, "HD003"); n != 1 {
		t.Fatalf("số phiếu HD003 = %d, want 1", n)
	}
	if got := fundBalance(t, pool, 2); !decEqual(t, got, "1500") {
		t.Fatalf("balance quỹ 2 = %s, want 1500 (1000 + 500 một lần)", got)
	}
}

// TestIntegration_RecordPaymentIn_AmountNotPositive: amount<=0 → 422 Validation,
// KHÔNG ghi gì, balance giữ nguyên.
func TestIntegration_RecordPaymentIn_AmountNotPositive(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	before := fundBalance(t, pool, 1)
	for _, amt := range []string{"0", "-100"} {
		p := app.RecordPaymentInParams{
			RecordPaymentIn: domain.RecordPaymentIn{
				OrderCode: "HD004", Amount: mustMoney(t, amt), FundAccountID: 1, BookType: domain.BookBoth,
			},
			IdemKey: "neg-" + amt,
		}
		_, err := recordInTx(ctx, pool, mod.Recorder(), p)
		if apperr.KindOf(err) != apperr.KindValidation {
			t.Fatalf("amount %q phải Validation(422), got %v", amt, err)
		}
	}
	if n := countAlive(t, pool, "HD004"); n != 0 {
		t.Fatalf("không được ghi phiếu khi amount<=0; có %d", n)
	}
	if got := fundBalance(t, pool, 1); !decEqual(t, got, before) {
		t.Fatalf("balance đổi sau lỗi: %s != %s", got, before)
	}
}

// TestIntegration_RecordPaymentIn_ManualNeedsIdemKey: phiếu thủ công (không
// bank_ref) thiếu IdemKey → 422 (tránh cộng tiền lặp vì không có khoá chống trùng).
func TestIntegration_RecordPaymentIn_ManualNeedsIdemKey(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	p := app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode: "HD005", Amount: mustMoney(t, "100"), FundAccountID: 1, BookType: domain.BookBoth,
		},
	}
	_, err := recordInTx(ctx, pool, mod.Recorder(), p)
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("thiếu idem key phải Validation(422), got %v", err)
	}
}

// TestIntegration_RecordPaymentIn_RaceSameKey: 2 tx ĐỒNG THỜI cùng idem key →
// CHỈ 1 row, KHÔNG cộng đôi, cả hai đều thành công (1 insert + 1 đọc-lại no-op).
func TestIntegration_RecordPaymentIn_RaceSameKey(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	p := app.RecordPaymentInParams{
		RecordPaymentIn: domain.RecordPaymentIn{
			OrderCode:     "HD006",
			Amount:        mustMoney(t, "300000"),
			FundAccountID: 1,
			BookType:      domain.BookBoth,
		},
		IdemKey: "race-key",
	}

	const n = 2
	var wg sync.WaitGroup
	ids := make([]int64, n)
	errs := make([]error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			got, err := recordInTx(ctx, pool, mod.Recorder(), p)
			ids[idx] = got.ID
			errs[idx] = err
		}(i)
	}
	wg.Wait()

	for i, e := range errs {
		if e != nil {
			t.Fatalf("race tx %d phải thành công (insert hoặc no-op), got %v", i, e)
		}
	}
	// Cùng phiếu: 2 tx trả cùng id.
	if ids[0] != ids[1] {
		t.Fatalf("race: 2 id khác nhau %d != %d (tạo trùng phiếu!)", ids[0], ids[1])
	}
	if c := countAlive(t, pool, "HD006"); c != 1 {
		t.Fatalf("race: số phiếu HD006 = %d, want 1 (unique index chặn)", c)
	}
	// Balance cộng đúng MỘT lần: 300000.
	if got := fundBalance(t, pool, 1); !decEqual(t, got, "300000") {
		t.Fatalf("race: balance = %s, want 300000 (KHÔNG cộng đôi)", got)
	}
}

// TestIntegration_ListOrderPayments: đọc phiếu của một đơn (route GET dùng).
func TestIntegration_ListOrderPayments(t *testing.T) {
	mod, pool := setup(t)
	ctx := context.Background()
	for i, amt := range []string{"100", "200"} {
		p := app.RecordPaymentInParams{
			RecordPaymentIn: domain.RecordPaymentIn{
				OrderCode: "HD007", Amount: mustMoney(t, amt), FundAccountID: 1, BookType: domain.BookBoth,
			},
			IdemKey: "list-" + amt,
		}
		if _, err := recordInTx(ctx, pool, mod.Recorder(), p); err != nil {
			t.Fatalf("ghi phiếu %d: %v", i, err)
		}
	}
	items, err := mod.Recorder().ListOrderPayments(ctx, "HD007", 50)
	if err != nil {
		t.Fatalf("ListOrderPayments: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("phiếu HD007 = %d, want 2", len(items))
	}
	for _, it := range items {
		if it.RefID != "HD007" || it.Flow != "in" {
			t.Fatalf("phiếu sai: %+v", it)
		}
	}
}
