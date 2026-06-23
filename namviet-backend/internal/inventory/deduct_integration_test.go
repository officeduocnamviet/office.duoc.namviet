package inventory_test

import (
	"context"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

// qtyFrom dựng domain.Quantity từ chuỗi (decimal — KHÔNG float).
func qtyFrom(t *testing.T, s string) domain.Quantity {
	t.Helper()
	q, err := domain.QuantityFromString(s)
	if err != nil {
		t.Fatalf("qty %q: %v", s, err)
	}
	return q
}

// deductInTx chạy DeductFEFO trong MỘT tx (mô phỏng caller orders/POS gộp atomic)
// và commit. Trả lô tiêu thụ + lỗi (rollback nếu lỗi → không ghi gì).
func deductInTx(ctx context.Context, pool *pgxpool.Pool, d *app.Deductor, whID, pID int64, qty domain.Quantity) ([]domain.ConsumedBatch, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	consumed, derr := d.DeductFEFO(ctx, tx, whID, pID, qty)
	if derr != nil {
		_ = tx.Rollback(ctx)
		return nil, derr
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return consumed, nil
}

// totalStock đọc tồn TỔNG (product_inventory.stock_quantity) của (warehouse,
// product) — để khẳng định KHÔNG bán âm.
func totalStock(t *testing.T, pool *pgxpool.Pool, whID, pID int64) string {
	t.Helper()
	var s string
	err := pool.QueryRow(context.Background(),
		`SELECT stock_quantity::text FROM public.product_inventory WHERE warehouse_id=$1 AND product_id=$2`,
		whID, pID).Scan(&s)
	if err != nil {
		t.Fatalf("totalStock: %v", err)
	}
	return s
}

// batchQty đọc tồn của một dòng inventory_batches.
func batchQty(t *testing.T, pool *pgxpool.Pool, invBatchID int64) string {
	t.Helper()
	var s string
	err := pool.QueryRow(context.Background(),
		`SELECT quantity::text FROM public.inventory_batches WHERE id=$1`, invBatchID).Scan(&s)
	if err != nil {
		t.Fatalf("batchQty: %v", err)
	}
	return s
}

// TestIntegration_DeductFEFO_HappyMultiBatch: trừ kho 1, product 100 lượng 50 →
// trải qua nhiều lô theo FEFO (502 hết hạn sớm nhất trước, rồi 503, rồi 501).
// Seed (xem inventory_integration_test.go): kho1 product100 có lô 502(25.5),
// 503(10), 501(40) còn tồn; tổng 75.5. Trừ 50 → 502 trọn(25.5) + 503 trọn(10) +
// 501 một phần(14.5).
func TestIntegration_DeductFEFO_HappyMultiBatch(t *testing.T) {
	fx := setup(t)
	d := inventory.NewDeductor(fx.pool)
	ctx := context.Background()

	consumed, err := deductInTx(ctx, fx.pool, d, 1, 100, qtyFrom(t, "50"))
	if err != nil {
		t.Fatalf("DeductFEFO happy phải thành công: %v", err)
	}
	// FEFO: 502(25.5) → 503(10) → 501(14.5).
	if len(consumed) != 3 {
		t.Fatalf("consumed = %d lô, want 3; %+v", len(consumed), consumed)
	}
	if consumed[0].BatchID != 502 || !consumed[0].Quantity.Equal(qtyFrom(t, "25.5")) {
		t.Fatalf("lô 1 (FEFO sớm nhất) sai: batch=%d qty=%s", consumed[0].BatchID, consumed[0].Quantity)
	}
	if consumed[1].BatchID != 503 || !consumed[1].Quantity.Equal(qtyFrom(t, "10")) {
		t.Fatalf("lô 2 sai: batch=%d qty=%s", consumed[1].BatchID, consumed[1].Quantity)
	}
	if consumed[2].BatchID != 501 || !consumed[2].Quantity.Equal(qtyFrom(t, "14.5")) {
		t.Fatalf("lô 3 (một phần) sai: batch=%d qty=%s", consumed[2].BatchID, consumed[2].Quantity)
	}
	// inbound_price mang theo cho COGS (lô 501 = 9000.50).
	if !consumed[2].InboundPrice.Equal(mustMoney(t, "9000.50")) {
		t.Fatalf("inbound_price lô 501 = %s, want 9000.50", consumed[2].InboundPrice)
	}

	// Tồn tổng kho1 product100 ban đầu 120.5 → còn 70.5. Tồn lô: 502→0, 503→0,
	// 501→25.5.
	if got := totalStock(t, fx.pool, 1, 100); !decEqual(got, "70.5") {
		t.Fatalf("tồn tổng sau trừ = %s, want 70.5", got)
	}
	if got := batchQty(t, fx.pool, 9002); !decEqual(got, "0") { // lô 502 kho1
		t.Fatalf("lô 502 sau trừ = %s, want 0", got)
	}
	if got := batchQty(t, fx.pool, 9001); !decEqual(got, "25.5") { // lô 501 kho1
		t.Fatalf("lô 501 sau trừ = %s, want 25.5", got)
	}
}

// TestIntegration_DeductFEFO_Insufficient_NoWrite: cần nhiều hơn tổng tồn →
// Conflict + KHÔNG ghi gì (rollback sạch). Tồn tổng + tồn lô giữ nguyên.
func TestIntegration_DeductFEFO_Insufficient_NoWrite(t *testing.T) {
	fx := setup(t)
	d := inventory.NewDeductor(fx.pool)
	ctx := context.Background()

	before := totalStock(t, fx.pool, 1, 100)
	beforeBatch := batchQty(t, fx.pool, 9002)

	// kho1 product100 tổng các lô = 75.5; cần 1000 → thiếu.
	_, err := deductInTx(ctx, fx.pool, d, 1, 100, qtyFrom(t, "1000"))
	if err == nil {
		t.Fatal("trừ quá tồn phải lỗi (Conflict)")
	}
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("thiếu tồn phải Conflict, got %v (%v)", apperr.KindOf(err), err)
	}
	// Rollback sạch: không dòng nào thay đổi.
	if got := totalStock(t, fx.pool, 1, 100); !decEqual(got, before) {
		t.Fatalf("tồn tổng đổi sau lỗi (phải rollback sạch): %s != %s", got, before)
	}
	if got := batchQty(t, fx.pool, 9002); !decEqual(got, beforeBatch) {
		t.Fatalf("tồn lô đổi sau lỗi (phải rollback sạch): %s != %s", got, beforeBatch)
	}
}

// TestIntegration_DeductFEFO_SkipDeletedAndEmpty: lô đã xóa (504, đã set
// deleted_at) và lô tồn 0 (9005) KHÔNG được tính vào FEFO. kho1 product100 chỉ trừ
// từ 502/503/501.
func TestIntegration_DeductFEFO_SkipDeletedAndEmpty(t *testing.T) {
	fx := setup(t)
	d := inventory.NewDeductor(fx.pool)
	ctx := context.Background()

	// Lô 504 (đã xóa) có tồn 99 ở dòng 9004; lô 505 (tồn 0) ở 9005. Nếu chúng lọt
	// vào FEFO, kết quả/thứ tự sẽ sai. Trừ 5 → chỉ chạm 502 (hết hạn sớm nhất).
	consumed, err := deductInTx(ctx, fx.pool, d, 1, 100, qtyFrom(t, "5"))
	if err != nil {
		t.Fatalf("DeductFEFO phải thành công: %v", err)
	}
	if len(consumed) != 1 || consumed[0].BatchID != 502 {
		t.Fatalf("phải chỉ trừ lô 502 (FEFO, bỏ deleted/empty); got %+v", consumed)
	}
	// Lô đã xóa (9004) giữ nguyên 99.
	if got := batchQty(t, fx.pool, 9004); !decEqual(got, "99") {
		t.Fatalf("lô đã xóa (9004) bị trừ nhầm: %s", got)
	}
}

// TestIntegration_DeductFEFO_Race: ca ĐUA — 2 tx ĐỒNG THỜI cùng trừ (kho9,sp900)
// có tồn VỪA ĐỦ cho 1 → đúng MỘT thành công, MỘT nhận Conflict; tổng trừ KHÔNG
// vượt tồn (KHÔNG âm). advisory lock tuần tự hoá 2 tx; tx thua thấy tồn đã hết →
// PlanFEFO thiếu → Conflict.
func TestIntegration_DeductFEFO_Race(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()

	// Scenario đua riêng: kho 9, product 900, 1 lô tồn ĐÚNG 1 đơn vị.
	seedRace(t, fx.pool)

	d := inventory.NewDeductor(fx.pool)

	const n = 2
	var wg sync.WaitGroup
	errs := make([]error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			_, errs[idx] = deductInTx(ctx, fx.pool, d, 9, 900, qtyFrom(t, "1"))
		}(i)
	}
	wg.Wait()

	// Đúng 1 thành công, 1 Conflict.
	var ok, conflict int
	for _, e := range errs {
		switch {
		case e == nil:
			ok++
		case apperr.KindOf(e) == apperr.KindConflict:
			conflict++
		default:
			t.Fatalf("lỗi lạ ngoài Conflict: %v", e)
		}
	}
	if ok != 1 || conflict != 1 {
		t.Fatalf("race: ok=%d conflict=%d, want ok=1 conflict=1 (errs=%v)", ok, conflict, errs)
	}

	// KHÔNG bán âm: tồn tổng = 0 (không âm), tồn lô = 0.
	if got := totalStock(t, fx.pool, 9, 900); !decEqual(got, "0") {
		t.Fatalf("race: tồn tổng = %s, want 0 (không âm)", got)
	}
	if got := batchQty(t, fx.pool, 9900); !decEqual(got, "0") {
		t.Fatalf("race: tồn lô = %s, want 0 (không âm)", got)
	}
}

// seedRace nạp scenario đua: kho 9 active, product 900 tồn tổng 1, 1 lô tồn 1.
func seedRace(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.warehouses (id, key, name, unit, type, status) VALUES
			(9, 'kho-dua', 'Kho Đua', 'Hộp', 'central', 'active')`,
		`INSERT INTO public.product_inventory (id, product_id, warehouse_id, stock_quantity) VALUES
			(1900, 900, 9, 1)`,
		`INSERT INTO public.batches (id, product_id, batch_code, expiry_date, inbound_price) VALUES
			(590, 900, 'RACE-LOT', '2026-12-31', 5000)`,
		`INSERT INTO public.inventory_batches (id, warehouse_id, product_id, batch_id, quantity) VALUES
			(9900, 9, 900, 590, 1)`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("seedRace: %v\nSQL: %s", err, s)
		}
	}
}

// ---- helpers tiền/decimal ----

func mustMoney(t *testing.T, s string) money.Money {
	t.Helper()
	m, err := money.FromString(s)
	if err != nil {
		t.Fatalf("money %q: %v", s, err)
	}
	return m
}

// decEqual so sánh hai biểu diễn decimal BẰNG GIÁ TRỊ ("70.5" == "70.50"), tránh
// phụ thuộc cách Postgres in scale. Dùng QuantityFromString (decimal, KHÔNG float).
func decEqual(a, b string) bool {
	qa, ea := domain.QuantityFromString(a)
	qb, eb := domain.QuantityFromString(b)
	if ea != nil || eb != nil {
		return a == b
	}
	return qa.Equal(qb)
}

// giữ tham chiếu pgx (helper tx dùng pool.Begin trả pgx.Tx gián tiếp).
var _ = pgx.ErrNoRows
