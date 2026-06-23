package orders_test

// Integration test cho ORCHESTRATION (P4b): ShipOrder / RecordPayment / CreatePosSale.
// Dựng Service với DEPS THẬT (inventory.Deductor + accounting.Poster + vat.IssuePort +
// finance.RecordPort) trên testcontainers Postgres → kiểm gộp atomic CROSS-MODULE:
// trừ kho FEFO + phát HĐ VAT + post sổ kép + ghi phiếu thu trong 1 transaction.
// Trọng tâm: (1) happy path đủ hiệu ứng; (2) ROLLBACK nguyên tử khi thiếu tồn.
// Tái dùng qty()/mustMoney() (write_integration_test.go) cùng package.

import (
	"context"
	"fmt"
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
	"github.com/Maneva-AI/namviet-backend/internal/orders"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/vat"
)

const orchWarehouse = 9

// setupOrch dựng pool + seed (kho/lô/tài khoản/kỳ/quỹ) + Service với DEPS THẬT.
func setupOrch(t *testing.T) (*orders.Service, *pgxpool.Pool) {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	seedOrch(t, pool)
	svc := orders.New(pool, orders.Deps{
		Deductor: inventory.NewDeductor(pool),
		Poster:   accounting.NewModule(pool).Poster(),
		Issuer:   vat.NewModule(pool).IssuePort(),
		Recorder: finance.NewModule(pool).RecordPort(),
	})
	return svc, pool
}

// seedOrch nạp dữ liệu nền cho orchestration:
//   - kho 9 (active); sản phẩm 701 (giá vốn 7000, tồn 100), 702 (14000, tồn 100),
//     703 (5000, tồn 1 — dùng test thiếu tồn). Mỗi sp 1 lô (FEFO đơn giản, COGS
//     tất định). KHÔNG seed products (test ref schema không ép FK order_items→products).
//   - chart_of_accounts: đúng mã DefaultRules dùng (131/111/112/511/3331/632/1561).
//   - accounting_periods: kỳ THÁNG HIỆN TẠI 'open' (khớp entry_date = time.Now()).
//   - fund_accounts(1): quỹ tiền mặt.
func seedOrch(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	now := time.Now()
	stmts := []string{
		`INSERT INTO public.warehouses (id, key, name, unit, type, status) VALUES
			(9, 'kho-orch', 'Kho Orchestration', 'Hộp', 'central', 'active')`,
		`INSERT INTO public.product_inventory (id, product_id, warehouse_id, stock_quantity) VALUES
			(7701, 701, 9, 100), (7702, 702, 9, 100), (7703, 703, 9, 1)`,
		`INSERT INTO public.batches (id, product_id, batch_code, expiry_date, inbound_price) VALUES
			(7010, 701, 'LOT-701', '2027-12-31', 7000),
			(7020, 702, 'LOT-702', '2027-12-31', 14000),
			(7030, 703, 'LOT-703', '2027-12-31', 5000)`,
		`INSERT INTO public.inventory_batches (id, warehouse_id, product_id, batch_id, quantity) VALUES
			(8010, 9, 701, 7010, 100), (8020, 9, 702, 7020, 100), (8030, 9, 703, 7030, 1)`,
		`INSERT INTO public.chart_of_accounts (account_code, name, type, balance_type, allow_posting) VALUES
			('131','Phải thu khách hàng','Tài sản','DEBIT', true),
			('111','Tiền mặt','Tài sản','DEBIT', true),
			('112','Tiền gửi ngân hàng','Tài sản','DEBIT', true),
			('511','Doanh thu bán hàng','Doanh thu','CREDIT', true),
			('3331','Thuế GTGT đầu ra','Nợ phải trả','CREDIT', true),
			('632','Giá vốn hàng bán','Chi phí','DEBIT', true),
			('1561','Giá mua hàng hoá','Tài sản','DEBIT', true)`,
		fmt.Sprintf(`INSERT INTO app.accounting_periods (id, year, month, status) VALUES
			(gen_random_uuid(), %d, %d, 'open')`, now.Year(), int(now.Month())),
		`INSERT INTO public.fund_accounts (id, name, type, balance) VALUES (1, 'Quỹ tiền mặt', 'cash', 0)`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("seedOrch: %v\nSQL: %s", err, s)
		}
	}
}

// vatLine8 dựng VATLine 8% (không override đơn giá → dùng đơn giá đơn).
func vatLine8() app.VATLine {
	return app.VATLine{VATRate: decimal.RequireFromString("0.08")}
}

// stockOf đọc tồn tổng product_inventory của (kho 9, product).
func stockOf(t *testing.T, pool *pgxpool.Pool, productID int64) string {
	t.Helper()
	var q string
	if err := pool.QueryRow(context.Background(),
		`SELECT stock_quantity::text FROM public.product_inventory WHERE warehouse_id=9 AND product_id=$1`, productID).Scan(&q); err != nil {
		t.Fatalf("stockOf %d: %v", productID, err)
	}
	return q
}

// countJournalEntries đếm bút toán đã post cho một đơn (source_id = code).
func countJournalEntries(t *testing.T, pool *pgxpool.Pool, code string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM app.journal_entries WHERE source_type='order' AND source_id=$1`, code).Scan(&n); err != nil {
		t.Fatalf("countJournalEntries: %v", err)
	}
	return n
}

// countUnbalancedEntries đếm bút toán của đơn mà Σdebit ≠ Σcredit (phải = 0).
func countUnbalancedEntries(t *testing.T, pool *pgxpool.Pool, code string) int {
	t.Helper()
	var n int
	err := pool.QueryRow(context.Background(), `
		SELECT count(*) FROM (
			SELECT l.entry_id FROM app.journal_entry_lines l
			JOIN app.journal_entries e ON e.id = l.entry_id
			WHERE e.source_type='order' AND e.source_id=$1
			GROUP BY l.entry_id HAVING SUM(l.debit) <> SUM(l.credit)
		) x`, code).Scan(&n)
	if err != nil {
		t.Fatalf("countUnbalancedEntries: %v", err)
	}
	return n
}

// countInvoices đếm HĐ VAT đã phát hành cho một đơn (order_code).
func countInvoices(t *testing.T, pool *pgxpool.Pool, code string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM app.sales_invoices WHERE order_code=$1`, code).Scan(&n); err != nil {
		t.Fatalf("countInvoices: %v", err)
	}
	return n
}

// paymentStatusOf đọc payment_status của đơn theo id.
func paymentStatusOf(t *testing.T, pool *pgxpool.Pool, id string) string {
	t.Helper()
	var ps string
	if err := pool.QueryRow(context.Background(),
		`SELECT payment_status FROM public.orders WHERE id=$1::uuid`, id).Scan(&ps); err != nil {
		t.Fatalf("paymentStatusOf: %v", err)
	}
	return ps
}

// TestIntegration_Orch_B2BLifecycle: create → confirm → ship → record payment.
// Kiểm đủ hiệu ứng ship (trừ kho, HĐ, sổ kép cân) + payment_status partial→paid.
func TestIntegration_Orch_B2BLifecycle(t *testing.T) {
	svc, pool := setupOrch(t)
	ctx := context.Background()

	// Tạo đơn B2B: 701 x2 (=20000) + 702 x3 (=60000) → final 80000 (không CK), CHƯA VAT.
	created, err := svc.CreateOrder(ctx, app.CreateOrderInput{
		OrderType: "B2B",
		Lines: []domain.DraftLine{
			{ProductID: 701, Quantity: qty(2), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()},
			{ProductID: 702, Quantity: qty(3), UOM: "Hộp", UnitPrice: mustMoney(t, "20000"), Discount: money.Zero()},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	code := created.Order.Code
	if created.Order.Final.String() != "80000" {
		t.Fatalf("final = %s, want 80000", created.Order.Final.String())
	}
	if _, err := svc.ConfirmOrder(ctx, created.Order.ID); err != nil {
		t.Fatalf("ConfirmOrder: %v", err)
	}

	// Ship: trừ kho FEFO + HĐ VAT (8%) + post sổ kép.
	shipped, err := svc.ShipOrder(ctx, app.ShipInput{
		OrderID:         created.Order.ID,
		WarehouseID:     orchWarehouse,
		CustomerTaxCode: "0312345678",
		Serial:          "C26TAA",
		MauSo:           "1",
		VATLines:        []app.VATLine{vatLine8(), vatLine8()},
	})
	if err != nil {
		t.Fatalf("ShipOrder phải thành công: %v", err)
	}
	if shipped.Order.Status != domain.StatusShipping.String() {
		t.Errorf("status sau ship = %q, want SHIPPING", shipped.Order.Status)
	}
	// COGS = 2×7000 + 3×14000 = 56000 (giá vốn per-unit của lô).
	if shipped.COGS.String() != "56000" {
		t.Errorf("COGS = %s, want 56000", shipped.COGS.String())
	}
	// HĐ: subtotal 80000, VAT 8% = 6400, total 86400.
	if shipped.Invoice.Subtotal.String() != "80000" || shipped.Invoice.VATAmount.String() != "6400" {
		t.Errorf("HĐ subtotal/vat = %s/%s, want 80000/6400", shipped.Invoice.Subtotal.String(), shipped.Invoice.VATAmount.String())
	}

	// Tồn bị trừ: 701 còn 98, 702 còn 97.
	if got := stockOf(t, pool, 701); got != "98" {
		t.Errorf("tồn 701 = %s, want 98", got)
	}
	if got := stockOf(t, pool, 702); got != "97" {
		t.Errorf("tồn 702 = %s, want 97", got)
	}
	// Sổ kép: 3 bút toán (INTERNAL doanh thu + INTERNAL giá vốn + TAX doanh thu;
	// TAX giá vốn bị tắt theo cờ). Tất cả CÂN (trigger DB đã ép — kiểm lại cho chắc).
	if n := countJournalEntries(t, pool, code); n != 3 {
		t.Errorf("số bút toán = %d, want 3", n)
	}
	if n := countUnbalancedEntries(t, pool, code); n != 0 {
		t.Errorf("có %d bút toán LỆCH Σ (phải 0)", n)
	}
	if n := countInvoices(t, pool, code); n != 1 {
		t.Errorf("số HĐ = %d, want 1", n)
	}

	// Thu một phần 30000 → partial.
	if _, err := svc.RecordPayment(ctx, app.RecordPaymentInput{
		OrderID: created.Order.ID, Amount: mustMoney(t, "30000"), FundAccountID: 1,
		BookType: financedomain.BookBoth, IdemKey: "pay-1",
	}); err != nil {
		t.Fatalf("RecordPayment 1: %v", err)
	}
	if ps := paymentStatusOf(t, pool, created.Order.ID); ps != "partial" {
		t.Errorf("payment_status sau thu 30000 = %q, want partial", ps)
	}
	// Thu nốt 50000 → đủ 80000 → paid.
	res2, err := svc.RecordPayment(ctx, app.RecordPaymentInput{
		OrderID: created.Order.ID, Amount: mustMoney(t, "50000"), FundAccountID: 1,
		BookType: financedomain.BookBoth, IdemKey: "pay-2",
	})
	if err != nil {
		t.Fatalf("RecordPayment 2: %v", err)
	}
	if res2.PaymentStatus != "paid" {
		t.Errorf("payment_status sau thu đủ = %q, want paid", res2.PaymentStatus)
	}
}

// TestIntegration_Orch_TwoStagePayment: thanh toán 2 bước (spec mục 55). NV giao
// hàng thu tiền mặt (Collected=true → phiếu 'pending') → nợ KH giảm NGAY
// (payment_status=paid vì đã-thu đếm cả pending), nhưng phiếu CHƯA 'completed'. Thủ
// quỹ "Xác nhận đã thu" (ConfirmReceipt) → phiếu 'completed'; gọi lại = idempotent.
// (DB test KHÔNG có trigger số dư prod → chỉ kiểm vòng đời status + payment_status.)
func TestIntegration_Orch_TwoStagePayment(t *testing.T) {
	svc, pool := setupOrch(t)
	ctx := context.Background()
	finRec := finance.NewModule(pool).Recorder()

	created, err := svc.CreateOrder(ctx, app.CreateOrderInput{
		OrderType: "B2B",
		Lines: []domain.DraftLine{
			{ProductID: 701, Quantity: qty(2), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()},
			{ProductID: 702, Quantity: qty(3), UOM: "Hộp", UnitPrice: mustMoney(t, "20000"), Discount: money.Zero()},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if _, err := svc.ConfirmOrder(ctx, created.Order.ID); err != nil {
		t.Fatalf("ConfirmOrder: %v", err)
	}
	if _, err := svc.ShipOrder(ctx, app.ShipInput{
		OrderID: created.Order.ID, WarehouseID: orchWarehouse, CustomerTaxCode: "0312345678",
		Serial: "C26TAA", MauSo: "1", VATLines: []app.VATLine{vatLine8(), vatLine8()},
	}); err != nil {
		t.Fatalf("ShipOrder: %v", err)
	}

	// NV thu tiền mặt đủ 80000, Collected=true → phiếu pending; nợ giảm ngay → paid.
	res, err := svc.RecordPayment(ctx, app.RecordPaymentInput{
		OrderID: created.Order.ID, Amount: mustMoney(t, "80000"), FundAccountID: 1,
		BookType: financedomain.BookBoth, Collected: true, IdemKey: "collect-1",
	})
	if err != nil {
		t.Fatalf("RecordPayment (collected): %v", err)
	}
	if res.PaymentStatus != "paid" {
		t.Errorf("payment_status sau NV thu (pending) = %q, want paid (pending tính cho công nợ)", res.PaymentStatus)
	}
	// Phiếu phải là 'pending' (chưa vào quỹ).
	var st string
	if err := pool.QueryRow(ctx, `SELECT status FROM public.finance_transactions WHERE id=$1`, res.Payment.ID).Scan(&st); err != nil {
		t.Fatalf("đọc status phiếu: %v", err)
	}
	if st != "pending" {
		t.Errorf("phiếu Collected phải 'pending', got %q", st)
	}

	// Thủ quỹ xác nhận → 'completed'.
	confirmed, err := finRec.ConfirmReceiptInOwnTx(ctx, res.Payment.ID)
	if err != nil {
		t.Fatalf("ConfirmReceipt: %v", err)
	}
	if !confirmed {
		t.Fatal("ConfirmReceipt lần đầu phải confirmed=true")
	}
	if err := pool.QueryRow(ctx, `SELECT status FROM public.finance_transactions WHERE id=$1`, res.Payment.ID).Scan(&st); err != nil {
		t.Fatalf("đọc status sau confirm: %v", err)
	}
	if st != "completed" {
		t.Errorf("sau ConfirmReceipt phải 'completed', got %q", st)
	}
	// Idempotent: xác nhận lại → false (không cộng đôi).
	again, err := finRec.ConfirmReceiptInOwnTx(ctx, res.Payment.ID)
	if err != nil {
		t.Fatalf("ConfirmReceipt lần 2: %v", err)
	}
	if again {
		t.Error("ConfirmReceipt lần 2 phải false (idempotent)")
	}
}

// TestIntegration_Orch_LumpSumAllocation: thu 1 CỤC từ khách (spec mục 55) → phân
// bổ cho các đơn chưa tất toán CŨ NHẤT trước. 3 đơn cùng khách (200k/300k/500k =
// 1tr); thu 600k → đơn1 200k (paid), đơn2 300k (paid), đơn3 100k (partial), leftover 0.
func TestIntegration_Orch_LumpSumAllocation(t *testing.T) {
	svc, pool := setupOrch(t)
	ctx := context.Background()
	cust := int64(555)

	// Tạo 3 đơn cùng khách, theo thứ tự (đơn1 cũ nhất). final = qty×10000.
	mk := func(qtyN int64) string {
		c, err := svc.CreateOrder(ctx, app.CreateOrderInput{
			CustomerID: &cust, OrderType: "B2B",
			Lines: []domain.DraftLine{{ProductID: 701, Quantity: qty(qtyN), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()}},
		})
		if err != nil {
			t.Fatalf("CreateOrder qty=%d: %v", qtyN, err)
		}
		return c.Order.ID
	}
	o1 := mk(20) // final 200k
	o2 := mk(30) // final 300k
	o3 := mk(50) // final 500k

	// Thu 1 cục 600k → phân bổ oldest-first.
	res, err := svc.RecordLumpSumPayment(ctx, app.LumpSumInput{
		CustomerID: cust, Amount: mustMoney(t, "600000"), FundAccountID: 1,
		BookType: financedomain.BookBoth, IdemKey: "lump-1",
	})
	if err != nil {
		t.Fatalf("RecordLumpSumPayment: %v", err)
	}
	if len(res.Allocations) != 3 {
		t.Fatalf("số dòng phân bổ = %d, want 3 (%+v)", len(res.Allocations), res.Allocations)
	}
	// Oldest-first: đơn1 200k, đơn2 300k, đơn3 100k.
	wantAmt := []string{"200000", "300000", "100000"}
	for i, a := range res.Allocations {
		if a.Amount.String() != wantAmt[i] {
			t.Errorf("phân bổ[%d] = %s, want %s", i, a.Amount.String(), wantAmt[i])
		}
	}
	if res.Leftover.String() != "0" {
		t.Errorf("leftover = %s, want 0", res.Leftover.String())
	}
	// payment_status: đơn1/đơn2 paid, đơn3 partial.
	if ps := paymentStatusOf(t, pool, o1); ps != "paid" {
		t.Errorf("đơn1 payment_status = %q, want paid", ps)
	}
	if ps := paymentStatusOf(t, pool, o2); ps != "paid" {
		t.Errorf("đơn2 payment_status = %q, want paid", ps)
	}
	if ps := paymentStatusOf(t, pool, o3); ps != "partial" {
		t.Errorf("đơn3 payment_status = %q, want partial", ps)
	}
	// Đã-thu đơn3 = 100k (qua allocation, app.order_paid_amount).
	d3, err := svc.GetOrder(ctx, o3)
	if err != nil {
		t.Fatalf("GetOrder o3: %v", err)
	}
	if d3.Order.Payment.Paid.String() != "100000" {
		t.Errorf("đã-thu đơn3 = %s, want 100000 (qua phân bổ)", d3.Order.Payment.Paid.String())
	}
}

// TestIntegration_Orch_RecordPayment_ReplayNoDouble: gọi RecordPayment 2 LẦN cùng
// Idempotency-Key (replay) → KHÔNG nhân đôi bút toán PAYMENT_IN, KHÔNG nhân đôi
// "đã thu". (Bug review #5: trước đây replay post lại bút toán.)
func TestIntegration_Orch_RecordPayment_ReplayNoDouble(t *testing.T) {
	svc, pool := setupOrch(t)
	ctx := context.Background()
	created, err := svc.CreateOrder(ctx, app.CreateOrderInput{
		OrderType: "B2B",
		Lines:     []domain.DraftLine{{ProductID: 701, Quantity: qty(2), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()}},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if _, err := svc.ConfirmOrder(ctx, created.Order.ID); err != nil {
		t.Fatalf("ConfirmOrder: %v", err)
	}
	if _, err := svc.ShipOrder(ctx, app.ShipInput{OrderID: created.Order.ID, WarehouseID: orchWarehouse, CustomerTaxCode: "0312345678", Serial: "C26TAA", VATLines: []app.VATLine{vatLine8()}}); err != nil {
		t.Fatalf("ShipOrder: %v", err)
	}
	code := created.Order.Code
	afterShip := countJournalEntries(t, pool, code) // 3: INTERNAL rev+COGS, TAX rev

	pay := app.RecordPaymentInput{OrderID: created.Order.ID, Amount: mustMoney(t, "20000"), FundAccountID: 1, BookType: financedomain.BookBoth, IdemKey: "rp-replay"}
	if _, err := svc.RecordPayment(ctx, pay); err != nil {
		t.Fatalf("RecordPayment lần 1: %v", err)
	}
	afterPay1 := countJournalEntries(t, pool, code) // +2 PAYMENT_IN (INTERNAL+TAX)
	if afterPay1 != afterShip+2 {
		t.Fatalf("sau thu lần 1 = %d bút toán, want %d (+2 PAYMENT_IN)", afterPay1, afterShip+2)
	}
	// Replay cùng key → KHÔNG post thêm bút toán, KHÔNG cộng đôi đã-thu.
	if _, err := svc.RecordPayment(ctx, pay); err != nil {
		t.Fatalf("RecordPayment lần 2 (replay): %v", err)
	}
	afterPay2 := countJournalEntries(t, pool, code)
	if afterPay2 != afterPay1 {
		t.Errorf("REPLAY nhân đôi bút toán: %d → %d (phải giữ %d)", afterPay1, afterPay2, afterPay1)
	}
	d, err := svc.GetOrder(ctx, created.Order.ID)
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	if d.Order.Payment.Paid.String() != "20000" {
		t.Errorf("đã-thu sau replay = %s, want 20000 (KHÔNG cộng đôi)", d.Order.Payment.Paid.String())
	}
}

// TestIntegration_Orch_LumpSum_ReplayNoDouble: RecordLumpSumPayment 2 LẦN cùng
// Idempotency-Key → KHÔNG nhân đôi allocation/đã-thu. (Bug review #4/#6/#7.)
func TestIntegration_Orch_LumpSum_ReplayNoDouble(t *testing.T) {
	svc, _ := setupOrch(t)
	ctx := context.Background()
	cust := int64(777)
	o, err := svc.CreateOrder(ctx, app.CreateOrderInput{
		CustomerID: &cust, OrderType: "B2B",
		Lines: []domain.DraftLine{{ProductID: 701, Quantity: qty(20), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()}},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	lump := app.LumpSumInput{CustomerID: cust, Amount: mustMoney(t, "150000"), FundAccountID: 1, BookType: financedomain.BookBoth, IdemKey: "ls-replay"}
	if _, err := svc.RecordLumpSumPayment(ctx, lump); err != nil {
		t.Fatalf("lump lần 1: %v", err)
	}
	if _, err := svc.RecordLumpSumPayment(ctx, lump); err != nil {
		t.Fatalf("lump lần 2 (replay): %v", err)
	}
	d, err := svc.GetOrder(ctx, o.Order.ID)
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	// final 200k, thu 150k 1 lần → đã-thu 150k (KHÔNG 300k dù gọi 2 lần).
	if d.Order.Payment.Paid.String() != "150000" {
		t.Errorf("đã-thu sau replay lump = %s, want 150000 (KHÔNG cộng đôi allocation)", d.Order.Payment.Paid.String())
	}
}

// TestIntegration_Orch_ShipRollback_Atomic: ship đơn có dòng THIẾU tồn → Conflict +
// TOÀN BỘ rollback: không trừ kho dòng nào, không HĐ, không bút toán, status giữ CONFIRMED.
func TestIntegration_Orch_ShipRollback_Atomic(t *testing.T) {
	svc, pool := setupOrch(t)
	ctx := context.Background()

	// Đơn: 701 x1 (đủ tồn) + 703 x5 (chỉ có 1 → THIẾU). 701 đứng trước để chứng minh
	// nó KHÔNG bị trừ khi 703 fail (rollback cả cụm).
	created, err := svc.CreateOrder(ctx, app.CreateOrderInput{
		OrderType: "B2B",
		Lines: []domain.DraftLine{
			{ProductID: 701, Quantity: qty(1), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()},
			{ProductID: 703, Quantity: qty(5), UOM: "Hộp", UnitPrice: mustMoney(t, "5000"), Discount: money.Zero()},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	code := created.Order.Code
	if _, err := svc.ConfirmOrder(ctx, created.Order.ID); err != nil {
		t.Fatalf("ConfirmOrder: %v", err)
	}

	_, err = svc.ShipOrder(ctx, app.ShipInput{
		OrderID:         created.Order.ID,
		WarehouseID:     orchWarehouse,
		CustomerTaxCode: "0312345678",
		Serial:          "C26TAA",
		VATLines:        []app.VATLine{vatLine8(), vatLine8()},
	})
	if err == nil {
		t.Fatal("ShipOrder thiếu tồn phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("thiếu tồn phải Conflict, got %v (%v)", apperr.KindOf(err), err)
	}

	// ROLLBACK nguyên tử: 701 vẫn 100 (KHÔNG bị trừ), 703 vẫn 1.
	if got := stockOf(t, pool, 701); got != "100" {
		t.Errorf("tồn 701 sau rollback = %s, want 100 (không bị trừ)", got)
	}
	if got := stockOf(t, pool, 703); got != "1" {
		t.Errorf("tồn 703 sau rollback = %s, want 1", got)
	}
	// Không HĐ mồ côi, không bút toán, status vẫn CONFIRMED.
	if n := countInvoices(t, pool, code); n != 0 {
		t.Errorf("có %d HĐ sau rollback (phải 0)", n)
	}
	if n := countJournalEntries(t, pool, code); n != 0 {
		t.Errorf("có %d bút toán sau rollback (phải 0)", n)
	}
	var status string
	if err := pool.QueryRow(ctx, `SELECT status FROM public.orders WHERE id=$1::uuid`, created.Order.ID).Scan(&status); err != nil {
		t.Fatalf("đọc status: %v", err)
	}
	if status != domain.StatusConfirmed.String() {
		t.Errorf("status sau rollback = %q, want CONFIRMED", status)
	}
}

// TestIntegration_Orch_PosSale_Atomic: bán lẻ tại quầy 1 nhịp → COMPLETED + paid +
// trừ kho + HĐ + sổ. Toàn bộ trong 1 tx.
func TestIntegration_Orch_PosSale_Atomic(t *testing.T) {
	svc, pool := setupOrch(t)
	ctx := context.Background()

	res, err := svc.CreatePosSale(ctx, app.PosSaleInput{
		Lines: []domain.DraftLine{
			{ProductID: 701, Quantity: qty(4), UOM: "Hộp", UnitPrice: mustMoney(t, "10000"), Discount: money.Zero()},
		},
		WarehouseID:     orchWarehouse,
		FundAccountID:   1,
		IdemKey:         "pos-1",
		IssueInvoice:    true,
		CustomerTaxCode: "0312345678",
		Serial:          "C26TAA",
		MauSo:           "1",
		VATLines:        []app.VATLine{vatLine8()},
	})
	if err != nil {
		t.Fatalf("CreatePosSale phải thành công: %v", err)
	}
	if res.Order.Status != domain.StatusCompleted.String() {
		t.Errorf("status POS = %q, want COMPLETED", res.Order.Status)
	}
	if res.Order.PaymentStatus != "paid" {
		t.Errorf("payment_status POS = %q, want paid", res.Order.PaymentStatus)
	}
	if res.Invoice == nil {
		t.Fatal("POS issue_invoice=true phải có HĐ")
	}
	// 701: 100 - 4 = 96.
	if got := stockOf(t, pool, 701); got != "96" {
		t.Errorf("tồn 701 sau POS = %s, want 96", got)
	}
	code := res.Order.Code
	// Sổ POS = 3 bút toán: INTERNAL doanh thu (Dr 111 tiền/Cr 511/Cr 3331) + INTERNAL
	// giá vốn + TAX doanh thu. POS thu NGAY → tiền nằm TRONG bút toán SALE (Dr 111),
	// KHÔNG post PAYMENT_IN riêng (đúng: B2C không có AR 131). Tất cả cân Σ.
	if n := countUnbalancedEntries(t, pool, code); n != 0 {
		t.Errorf("POS có %d bút toán LỆCH Σ (phải 0)", n)
	}
	if n := countJournalEntries(t, pool, code); n < 3 {
		t.Errorf("POS số bút toán = %d, want >= 3", n)
	}
}
