package inventory_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

// stockInInTx chạy StockIn trong MỘT tx (mô phỏng purchasing gộp atomic) + commit.
func stockInInTx(ctx context.Context, pool *pgxpool.Pool, s *app.StockInner, whID, pID int64, code string, exp *time.Time, qty domain.Quantity, price money.Money) (int64, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	batchID, serr := s.StockIn(ctx, tx, whID, pID, code, exp, nil, qty, price)
	if serr != nil {
		_ = tx.Rollback(ctx)
		return 0, serr
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return batchID, nil
}

// TestIntegration_StockIn_NewLineAndAccumulate: nhập kho cho (kho,sp) CHƯA có dòng
// tồn tổng → tạo dòng mới; nhập tiếp → cộng dồn. Lô tạo với inbound_price=giá nhập.
func TestIntegration_StockIn_NewLineAndAccumulate(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	s := inventory.NewStockInner(fx.pool)
	exp := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)

	// (kho 1, sp 700) CHƯA có dòng product_inventory → InsertStockTotal.
	b1, err := stockInInTx(ctx, fx.pool, s, 1, 700, "IN-LOT-1", &exp, qtyFrom(t, "30"), mustMoney(t, "1500"))
	if err != nil {
		t.Fatalf("StockIn lần 1: %v", err)
	}
	if got := totalStock(t, fx.pool, 1, 700); !decEqual(got, "30") {
		t.Fatalf("tồn sau nhập 1 = %s, want 30", got)
	}
	if got := batchInbound(t, fx.pool, b1); !decEqual(got, "1500") {
		t.Fatalf("inbound_price lô 1 = %s, want 1500", got)
	}

	// Nhập tiếp lô khác cùng (kho,sp) → AddStockTotal cộng dồn (30 + 20 = 50).
	if _, err := stockInInTx(ctx, fx.pool, s, 1, 700, "IN-LOT-2", &exp, qtyFrom(t, "20"), mustMoney(t, "1600")); err != nil {
		t.Fatalf("StockIn lần 2: %v", err)
	}
	if got := totalStock(t, fx.pool, 1, 700); !decEqual(got, "50") {
		t.Fatalf("tồn sau nhập 2 = %s, want 50 (cộng dồn)", got)
	}
}

// TestIntegration_StockIn_ZeroQty_Validation: qty <= 0 → Validation, KHÔNG ghi gì.
func TestIntegration_StockIn_ZeroQty_Validation(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	s := inventory.NewStockInner(fx.pool)
	exp := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)
	_, err := stockInInTx(ctx, fx.pool, s, 1, 701, "Z", &exp, qtyFrom(t, "0"), mustMoney(t, "100"))
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("qty=0 phải Validation, got %v", err)
	}
}

// TestIntegration_StockIn_ThenDeductFEFO: nhập kho rồi trừ FEFO cùng lô — khẳng định
// StockIn ↔ DeductFEFO đối xứng (lô StockIn tạo được FEFO xuất với inbound_price đúng).
func TestIntegration_StockIn_ThenDeductFEFO(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	si := inventory.NewStockInner(fx.pool)
	d := inventory.NewDeductor(fx.pool)
	exp := time.Date(2027, 3, 31, 0, 0, 0, 0, time.UTC)

	if _, err := stockInInTx(ctx, fx.pool, si, 1, 702, "RT", &exp, qtyFrom(t, "12"), mustMoney(t, "2500")); err != nil {
		t.Fatalf("StockIn: %v", err)
	}
	consumed, err := deductInTx(ctx, fx.pool, d, 1, 702, qtyFrom(t, "5"))
	if err != nil {
		t.Fatalf("DeductFEFO sau StockIn: %v", err)
	}
	if len(consumed) != 1 || !consumed[0].InboundPrice.Equal(mustMoney(t, "2500")) {
		t.Fatalf("FEFO lô StockIn sai: %+v", consumed)
	}
	if got := totalStock(t, fx.pool, 1, 702); !decEqual(got, "7") {
		t.Fatalf("tồn sau nhập 12 trừ 5 = %s, want 7", got)
	}
}

// batchInbound đọc inbound_price của một lô (public.batches).
func batchInbound(t *testing.T, pool *pgxpool.Pool, batchID int64) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(),
		`SELECT inbound_price::text FROM public.batches WHERE id=$1`, batchID).Scan(&s); err != nil {
		t.Fatalf("batchInbound: %v", err)
	}
	return s
}

var _ = pgx.ErrNoRows
