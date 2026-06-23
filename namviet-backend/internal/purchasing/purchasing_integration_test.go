package purchasing_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/accounting"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/finance"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/inventory"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/app"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

// fixture wiring purchasing thật (port inventory.StockIn + accounting.Poster +
// finance.RecordOutPort) trên Postgres testcontainers — KHÔNG mock, đúng tinh thần
// "fakes > mocks cho port; real cho SQL".
type fixture struct {
	svc  *purchasing.Service
	pool *pgxpool.Pool
}

func setup(t *testing.T) fixture {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	seed(t, pool)

	accMod := accounting.NewModule(pool)
	finMod := finance.NewModule(pool)
	svc := purchasing.New(pool, purchasing.Deps{
		StockIn: inventory.NewStockInner(pool),
		Poster:  accMod.Poster(),
		Payer:   finMod.RecordOutPort(),
	})
	return fixture{svc: svc, pool: pool}
}

// seed: cây TK TT133 (1561/133/331/111/112), kỳ kế toán mở (entry_date = now() →
// năm/tháng hiện tại), 1 kho, 1 quỹ tiền mặt dư lớn.
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.chart_of_accounts (account_code, name, type, balance_type, allow_posting) VALUES
			('1561','Hàng hoá','Tài sản','DEBIT', true),
			('133','Thuế GTGT được khấu trừ','Tài sản','DEBIT', true),
			('331','Phải trả người bán','Nợ','CREDIT', true),
			('111','Tiền mặt','Tài sản','DEBIT', true),
			('112','Tiền gửi ngân hàng','Tài sản','DEBIT', true)`,
		// Kỳ kế toán mở cho THÁNG HIỆN TẠI (entry_date = now() khi post).
		`INSERT INTO app.accounting_periods (id, year, month, status)
			SELECT gen_random_uuid(), EXTRACT(YEAR FROM now())::int, EXTRACT(MONTH FROM now())::int, 'open'`,
		`INSERT INTO public.warehouses (id, key, name, unit, type, status) VALUES
			(1, 'kho-mua', 'Kho Nhập', 'Hộp', 'central', 'active')`,
		`INSERT INTO public.fund_accounts (id, name, type, balance) VALUES (1, 'Quỹ tiền mặt', 'cash', 100000000)`,
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

func mustDec(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("decimal %q: %v", s, err)
	}
	return d
}

// ---- helpers đọc state ----

func stockTotal(t *testing.T, pool *pgxpool.Pool, whID, pID int64) string {
	t.Helper()
	var s string
	err := pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(stock_quantity),0)::text FROM public.product_inventory WHERE warehouse_id=$1 AND product_id=$2`,
		whID, pID).Scan(&s)
	if err != nil {
		t.Fatalf("stockTotal: %v", err)
	}
	return s
}

func batchCount(t *testing.T, pool *pgxpool.Pool, pID int64) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM public.batches WHERE product_id=$1 AND deleted_at IS NULL`, pID).Scan(&n); err != nil {
		t.Fatalf("batchCount: %v", err)
	}
	return n
}

func inboundPrice(t *testing.T, pool *pgxpool.Pool, batchID int64) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(),
		`SELECT inbound_price::text FROM public.batches WHERE id=$1`, batchID).Scan(&s); err != nil {
		t.Fatalf("inboundPrice: %v", err)
	}
	return s
}

func journalLineCount(t *testing.T, pool *pgxpool.Pool, sourceID, account string) int {
	t.Helper()
	var n int
	err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM app.journal_entry_lines jl
		 JOIN app.journal_entries je ON je.id = jl.entry_id
		 WHERE je.source_type='purchase' AND je.source_id=$1 AND jl.account_code=$2`,
		sourceID, account).Scan(&n)
	if err != nil {
		t.Fatalf("journalLineCount: %v", err)
	}
	return n
}

func paymentOutCount(t *testing.T, pool *pgxpool.Pool, poCode string) int {
	t.Helper()
	var n int
	err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM public.finance_transactions WHERE ref_type='purchase_order' AND ref_id=$1 AND flow='out' AND deleted_at IS NULL`,
		poCode).Scan(&n)
	if err != nil {
		t.Fatalf("paymentOutCount: %v", err)
	}
	return n
}

func fundBalance(t *testing.T, pool *pgxpool.Pool, id int64) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(),
		`SELECT balance::text FROM public.fund_accounts WHERE id=$1`, id).Scan(&s); err != nil {
		t.Fatalf("fundBalance: %v", err)
	}
	return s
}

func decEqual(a, b string) bool {
	da, ea := decimal.NewFromString(a)
	db, eb := decimal.NewFromString(b)
	if ea != nil || eb != nil {
		return a == b
	}
	return da.Equal(db)
}

// draftLines: 2 dòng — sp 100 (10 × 1000, vat 8%), sp 200 (5 × 2000, vat 8%).
// public.batches.expiry_date là NOT NULL (lô dược BẮT BUỘC hạn dùng) → set expiry.
func draftLines(t *testing.T) []domain.DraftLine {
	exp := time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)
	mfg := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return []domain.DraftLine{
		{ProductID: 100, Quantity: mustDec(t, "10"), UnitCost: mustMoney(t, "1000"), VATRate: mustDec(t, "0.08"), BatchCode: "LOT-A", ExpiryDate: &exp, ManufacturingDate: &mfg},
		{ProductID: 200, Quantity: mustDec(t, "5"), UnitCost: mustMoney(t, "2000"), VATRate: mustDec(t, "0.08"), BatchCode: "LOT-B", ExpiryDate: &exp},
	}
}

// TestIntegration_FullFlow: CreatePO → Confirm → Receive (nhập kho + post sổ) → Pay
// (chi NCC + post sổ). Khẳng định tồn tăng, lô tạo với inbound_price=unit_cost, sổ
// Dr 1561+133/Cr 331 + Dr 331/Cr 111, số dư quỹ giảm đúng.
func TestIntegration_FullFlow(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	a := fx.svc.App()

	created, wasCreated, err := a.CreatePO(ctx, app.CreatePOInput{
		SupplierName: "NCC A", Lines: draftLines(t), IdemKey: "po-1",
	})
	if err != nil || !wasCreated {
		t.Fatalf("CreatePO: err=%v created=%v", err, wasCreated)
	}
	po := created.PO
	if po.Status != "draft" {
		t.Fatalf("status sau tạo = %q, want draft", po.Status)
	}
	// total = 10*1000 + 5*2000 = 20000 ; vat = 800 + 800 = 1600
	if !po.TotalAmount.Equal(mustMoney(t, "20000")) || !po.VATAmount.Equal(mustMoney(t, "1600")) {
		t.Fatalf("tiền PO sai: total=%s vat=%s", po.TotalAmount, po.VATAmount)
	}

	if _, err := a.ConfirmPO(ctx, po.ID); err != nil {
		t.Fatalf("ConfirmPO: %v", err)
	}

	// Receive: nhập kho 2 lô + post sổ INTERNAL.
	recv, err := a.ReceivePO(ctx, app.ReceiveInput{POID: po.ID, WarehouseID: 1})
	if err != nil {
		t.Fatalf("ReceivePO: %v", err)
	}
	if recv.PO.Status != "received" {
		t.Fatalf("status sau nhận = %q, want received", recv.PO.Status)
	}
	if len(recv.CreatedBatch) != 2 {
		t.Fatalf("tạo %d lô, want 2", len(recv.CreatedBatch))
	}
	// inventoryCost = 20000 ; vat = 1600.
	if !recv.InventoryCost.Equal(mustMoney(t, "20000")) || !recv.VATAmount.Equal(mustMoney(t, "1600")) {
		t.Fatalf("nhập kho sai: cost=%s vat=%s", recv.InventoryCost, recv.VATAmount)
	}
	// Tồn tăng: sp100 = 10, sp200 = 5.
	if got := stockTotal(t, fx.pool, 1, 100); !decEqual(got, "10") {
		t.Fatalf("tồn sp100 = %s, want 10", got)
	}
	if got := stockTotal(t, fx.pool, 1, 200); !decEqual(got, "5") {
		t.Fatalf("tồn sp200 = %s, want 5", got)
	}
	if n := batchCount(t, fx.pool, 100); n != 1 {
		t.Fatalf("sp100 có %d lô, want 1", n)
	}
	// inbound_price = unit_cost per-unit (1000), KHÔNG nhân số lượng.
	if got := inboundPrice(t, fx.pool, recv.CreatedBatch[0]); !decEqual(got, "1000") {
		t.Fatalf("inbound_price lô sp100 = %s, want 1000 (per-unit)", got)
	}
	// Sổ INTERNAL: Dr 1561 + Dr 133 / Cr 331 (1 entry → mỗi TK 1 dòng).
	if journalLineCount(t, fx.pool, po.Code, "1561") != 1 {
		t.Fatalf("thiếu dòng Dr 1561")
	}
	if journalLineCount(t, fx.pool, po.Code, "133") != 1 {
		t.Fatalf("thiếu dòng Dr 133")
	}
	if journalLineCount(t, fx.pool, po.Code, "331") < 1 {
		t.Fatalf("thiếu dòng Cr 331")
	}

	// Pay: chi NCC = total + vat = 21600.
	balBefore := fundBalance(t, fx.pool, 1)
	pay, err := a.PaySupplier(ctx, app.PaySupplierInput{
		POID: po.ID, Amount: mustMoney(t, "21600"), FundAccountID: 1,
		BookType: financedomain.BookInternal, IdemKey: "pay-1",
	})
	if err != nil {
		t.Fatalf("PaySupplier: %v", err)
	}
	if pay.PO.Status != "paid" {
		t.Fatalf("status sau chi = %q, want paid", pay.PO.Status)
	}
	if pay.Payment.Flow != "out" {
		t.Fatalf("phiếu flow = %q, want out", pay.Payment.Flow)
	}
	if n := paymentOutCount(t, fx.pool, po.Code); n != 1 {
		t.Fatalf("phiếu chi PO = %d, want 1", n)
	}
	// Số dư quỹ giảm đúng 21600 (trigger trừ MỘT lần khi flow=out completed).
	wantBal, _ := decimal.NewFromString(balBefore)
	wantBal = wantBal.Sub(decimal.RequireFromString("21600"))
	if got := fundBalance(t, fx.pool, 1); !decEqual(got, wantBal.String()) {
		t.Fatalf("số dư quỹ sau chi = %s, want %s (giảm 21600)", got, wantBal)
	}
	// Sổ chi NCC: Dr 331 / Cr 111.
	if journalLineCount(t, fx.pool, po.Code, "111") != 1 {
		t.Fatalf("thiếu dòng Cr 111 (chi NCC)")
	}
}

// TestIntegration_CreatePO_Idempotent: replay cùng Idempotency-Key → 1 PO, created=false.
func TestIntegration_CreatePO_Idempotent(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	a := fx.svc.App()
	in := app.CreatePOInput{SupplierName: "NCC", Lines: draftLines(t), IdemKey: "dup-key"}
	first, c1, err := a.CreatePO(ctx, in)
	if err != nil || !c1 {
		t.Fatalf("lần 1: err=%v created=%v", err, c1)
	}
	second, c2, err := a.CreatePO(ctx, in)
	if err != nil {
		t.Fatalf("lần 2: %v", err)
	}
	if c2 {
		t.Fatal("lần 2 phải created=false (idempotent hit)")
	}
	if first.PO.ID != second.PO.ID {
		t.Fatalf("idempotent: id khác %s != %s", first.PO.ID, second.PO.ID)
	}
}

// TestIntegration_ReceivePO_ReplayNoDoubleStockIn: gọi ReceivePO lần 2 (PO đã
// received) → Conflict, KHÔNG nhập kho lại / KHÔNG post sổ lại (tồn + dòng sổ giữ nguyên).
func TestIntegration_ReceivePO_ReplayNoDoubleStockIn(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	a := fx.svc.App()

	created, _, err := a.CreatePO(ctx, app.CreatePOInput{SupplierName: "NCC", Lines: draftLines(t), IdemKey: "po-r"})
	if err != nil {
		t.Fatalf("CreatePO: %v", err)
	}
	po := created.PO
	if _, err := a.ConfirmPO(ctx, po.ID); err != nil {
		t.Fatalf("ConfirmPO: %v", err)
	}
	if _, err := a.ReceivePO(ctx, app.ReceiveInput{POID: po.ID, WarehouseID: 1}); err != nil {
		t.Fatalf("ReceivePO lần 1: %v", err)
	}
	stockAfter := stockTotal(t, fx.pool, 1, 100)
	batchesAfter := batchCount(t, fx.pool, 100)

	// Replay: PO đã received → Conflict (state machine guard chặn ở bước 1).
	_, err = a.ReceivePO(ctx, app.ReceiveInput{POID: po.ID, WarehouseID: 1})
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("replay receive phải Conflict, got %v", err)
	}
	// Tồn + lô KHÔNG đổi (không nhập kho lại).
	if got := stockTotal(t, fx.pool, 1, 100); !decEqual(got, stockAfter) {
		t.Fatalf("tồn đổi sau replay (phải giữ): %s != %s", got, stockAfter)
	}
	if got := batchCount(t, fx.pool, 100); got != batchesAfter {
		t.Fatalf("số lô đổi sau replay (phải giữ): %d != %d", got, batchesAfter)
	}
	// Sổ KHÔNG nhân đôi: vẫn đúng 1 dòng Dr 1561.
	if n := journalLineCount(t, fx.pool, po.Code, "1561"); n != 1 {
		t.Fatalf("dòng Dr 1561 = %d sau replay, want 1 (không post lại)", n)
	}
}

// TestIntegration_PaySupplier_ReplayNoDoublePost: gọi PaySupplier lần 2 CÙNG
// Idempotency-Key (sau commit) → KHÔNG post bút toán lại / KHÔNG trừ tiền lại.
// Vì PO đã paid, bước (1) chặn bằng Conflict — khẳng định cờ created + state machine.
func TestIntegration_PaySupplier_ReplayNoDoublePost(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	a := fx.svc.App()

	created, _, err := a.CreatePO(ctx, app.CreatePOInput{SupplierName: "NCC", Lines: draftLines(t), IdemKey: "po-p"})
	if err != nil {
		t.Fatalf("CreatePO: %v", err)
	}
	po := created.PO
	if _, err := a.ConfirmPO(ctx, po.ID); err != nil {
		t.Fatalf("ConfirmPO: %v", err)
	}
	if _, err := a.ReceivePO(ctx, app.ReceiveInput{POID: po.ID, WarehouseID: 1}); err != nil {
		t.Fatalf("ReceivePO: %v", err)
	}
	payIn := app.PaySupplierInput{
		POID: po.ID, Amount: mustMoney(t, "21600"), FundAccountID: 1,
		BookType: financedomain.BookInternal, IdemKey: "pay-dup",
	}
	if _, err := a.PaySupplier(ctx, payIn); err != nil {
		t.Fatalf("PaySupplier lần 1: %v", err)
	}
	balAfter := fundBalance(t, fx.pool, 1)

	// Replay: PO đã paid → Conflict (state machine chặn). KHÔNG trừ tiền/post lại.
	_, err = a.PaySupplier(ctx, payIn)
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("replay pay phải Conflict, got %v", err)
	}
	if n := paymentOutCount(t, fx.pool, po.Code); n != 1 {
		t.Fatalf("phiếu chi = %d sau replay, want 1 (không chi lại)", n)
	}
	if got := fundBalance(t, fx.pool, 1); !decEqual(got, balAfter) {
		t.Fatalf("số dư quỹ đổi sau replay (phải giữ): %s != %s", got, balAfter)
	}
	if n := journalLineCount(t, fx.pool, po.Code, "111"); n != 1 {
		t.Fatalf("dòng Cr 111 = %d sau replay, want 1 (không post lại)", n)
	}
}

// TestIntegration_Receive_WrongState: nhận hàng khi PO đang draft (chưa confirm) → Conflict.
func TestIntegration_Receive_WrongState(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	a := fx.svc.App()
	created, _, err := a.CreatePO(ctx, app.CreatePOInput{SupplierName: "NCC", Lines: draftLines(t), IdemKey: "po-w"})
	if err != nil {
		t.Fatalf("CreatePO: %v", err)
	}
	_, err = a.ReceivePO(ctx, app.ReceiveInput{POID: created.PO.ID, WarehouseID: 1})
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("nhận hàng PO draft phải Conflict, got %v", err)
	}
}
