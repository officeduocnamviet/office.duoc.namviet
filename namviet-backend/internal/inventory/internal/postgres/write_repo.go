package postgres

import (
	"context"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// WriteRepo implement app.StockWriter trên appdb.Queries (bind tx do caller/
// TxManager truyền). Mọi thao tác GHI tồn (advisory lock, đọc FOR UPDATE, trừ lô,
// trừ tổng) chạy trong CÙNG tx này — gộp atomic với nghiệp vụ orders/POS. GHI bảng
// public.* kế thừa (strangler, ADR 0001). quantity NUMERIC <-> domain.Quantity
// (decimal), KHÔNG float.
type WriteRepo struct{ q *appdb.Queries }

// NewWriteRepo tạo repo ghi từ một *appdb.Queries (đã bind tx).
func NewWriteRepo(q *appdb.Queries) *WriteRepo { return &WriteRepo{q: q} }

// LockWarehouseProduct lấy pg_advisory_xact_lock theo (warehouse, product) — giữ
// tới hết tx, tuần tự hoá trừ kho cùng (kho,sp).
func (r *WriteRepo) LockWarehouseProduct(ctx context.Context, warehouseID, productID int64) error {
	return r.q.LockWarehouseProduct(ctx, appdb.LockWarehouseProductParams{
		WarehouseID: warehouseID,
		ProductID:   productID,
	})
}

// ListBatchesForDeductFEFO trả các lô còn tồn của (warehouse, product) theo FEFO
// (expiry ASC), đã FOR UPDATE giữ dòng. Map row -> domain.Batch (quantity/
// inbound_price qua big.Int, KHÔNG float).
func (r *WriteRepo) ListBatchesForDeductFEFO(ctx context.Context, warehouseID, productID int64) ([]domain.Batch, error) {
	rows, err := r.q.ListBatchesForDeductFEFO(ctx, appdb.ListBatchesForDeductFEFOParams{
		WarehouseID: warehouseID,
		ProductID:   productID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Batch, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Batch{
			InventoryBatchID:  row.InventoryBatchID,
			WarehouseID:       row.WarehouseID,
			ProductID:         row.ProductID,
			BatchID:           row.BatchID,
			Quantity:          numericToQuantity(row.Quantity),
			BatchCode:         row.BatchCode,
			ExpiryDate:        row.ExpiryDate.Time,
			ManufacturingDate: row.ManufacturingDate.Time,
			HasManufacturing:  row.ManufacturingDate.Valid,
			InboundPrice:      numericToMoney(row.InboundPrice),
		})
	}
	return out, nil
}

// DeductBatch trừ qty khỏi một dòng inventory_batches (guard DB quantity >= qty).
func (r *WriteRepo) DeductBatch(ctx context.Context, inventoryBatchID int64, qty domain.Quantity) error {
	return r.q.DeductInventoryBatch(ctx, appdb.DeductInventoryBatchParams{
		InventoryBatchID: inventoryBatchID,
		Qty:              quantityToNumeric(qty),
	})
}

// DeductTotal trừ qty khỏi tồn tổng product_inventory của (warehouse, product).
func (r *WriteRepo) DeductTotal(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) error {
	return r.q.DeductProductInventory(ctx, appdb.DeductProductInventoryParams{
		WarehouseID: warehouseID,
		ProductID:   productID,
		Qty:         quantityToNumeric(qty),
	})
}

// ---- NHẬP KHO (StockIn — đối xứng trừ kho) ----

// InsertBatch tạo MỘT lô mới (public.batches) với inbound_price = giá nhập per-unit.
// Trả batchID. expiry/mfg nullable → pgtype.Date Valid=false.
func (r *WriteRepo) InsertBatch(ctx context.Context, productID int64, batchCode string, expiry, mfg *time.Time, inboundPrice money.Money) (int64, error) {
	return r.q.InsertBatch(ctx, appdb.InsertBatchParams{
		ProductID:         productID,
		BatchCode:         batchCode,
		ExpiryDate:        datePtr(expiry),
		ManufacturingDate: datePtr(mfg),
		InboundPrice:      moneyToNumeric(inboundPrice),
	})
}

// InsertInventoryBatch tạo dòng tồn-theo-lô cho (warehouse, product, batch) vừa tạo.
func (r *WriteRepo) InsertInventoryBatch(ctx context.Context, warehouseID, productID, batchID int64, qty domain.Quantity) error {
	return r.q.InsertInventoryBatch(ctx, appdb.InsertInventoryBatchParams{
		WarehouseID: warehouseID,
		ProductID:   productID,
		BatchID:     batchID,
		Quantity:    quantityToNumeric(qty),
	})
}

// AddStockTotal cộng dồn tồn tổng cho (warehouse, product) đã có dòng. Trả số dòng đổi.
func (r *WriteRepo) AddStockTotal(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) (int64, error) {
	return r.q.AddProductInventoryStock(ctx, appdb.AddProductInventoryStockParams{
		WarehouseID: warehouseID,
		ProductID:   productID,
		Qty:         quantityToNumeric(qty),
	})
}

// InsertStockTotal tạo dòng tồn tổng mới cho (warehouse, product) chưa có.
func (r *WriteRepo) InsertStockTotal(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) error {
	return r.q.InsertProductInventoryStock(ctx, appdb.InsertProductInventoryStockParams{
		WarehouseID: warehouseID,
		ProductID:   productID,
		Qty:         quantityToNumeric(qty),
	})
}

// moneyToNumeric chuyển money.Money sang pgtype.Numeric KHÔNG qua float (đối xứng
// numericToMoney ở repo.go) — dùng khi GHI inbound_price lô.
func moneyToNumeric(m money.Money) pgtype.Numeric {
	d := m.Decimal()
	return pgtype.Numeric{
		Int:   new(big.Int).Set(d.Coefficient()),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

// datePtr chuyển *time.Time → pgtype.Date (nil → Valid=false cho cột date nullable).
func datePtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

// quantityToNumeric chuyển domain.Quantity (decimal) sang pgtype.Numeric KHÔNG qua
// float: lấy coefficient (big.Int) + exponent từ decimal nền. Dùng khi GHI (trừ
// tồn lô/tổng).
func quantityToNumeric(q domain.Quantity) pgtype.Numeric {
	d := q.Decimal()
	return pgtype.Numeric{
		Int:   new(big.Int).Set(d.Coefficient()),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

// Đảm bảo WriteRepo thoả port app ở compile-time.
var _ app.StockWriter = (*WriteRepo)(nil)
