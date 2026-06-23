// Package postgres là ADAPTER ra phía cơ sở dữ liệu của inventory: implement port
// domain.Repository bằng query sinh từ sqlc (appdb) và map row <-> entity domain.
// ĐỌC bảng public.* kế thừa (strangler-fig, ADR 0001). Nằm dưới internal/ nên
// module khác KHÔNG import được. LÁT NÀY CHỈ ĐỌC → repo chỉ có thao tác đọc, bind
// thẳng pool (không tx). Đường ghi/trừ tồn HOÃN sau orders.
package postgres

import (
	"context"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Repo implement domain.Repository trên appdb.Queries (bind pool).
type Repo struct{ q *appdb.Queries }

// NewRepo tạo repo từ một *appdb.Queries (đã bind pool).
func NewRepo(q *appdb.Queries) *Repo { return &Repo{q: q} }

func (r *Repo) ListWarehouses(ctx context.Context, f domain.WarehouseFilter) ([]domain.Warehouse, error) {
	rows, err := r.q.ListWarehouses(ctx, appdb.ListWarehousesParams{
		Status:   strPtr(f.Status),
		RowLimit: f.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Warehouse, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Warehouse{
			ID:         row.ID,
			Key:        row.Key,
			Name:       row.Name,
			Unit:       row.Unit,
			Address:    derefStr(row.Address),
			Type:       row.Type,
			Code:       derefStr(row.Code),
			Manager:    derefStr(row.Manager),
			Phone:      derefStr(row.Phone),
			Status:     row.Status,
			CompanyID:  derefStr(row.CompanyID),
			OutletType: derefStr(row.OutletType),
			CreatedAt:  row.CreatedAt.Time,
		})
	}
	return out, nil
}

func (r *Repo) ListStock(ctx context.Context, f domain.StockFilter) ([]domain.StockItem, error) {
	rows, err := r.q.ListStock(ctx, appdb.ListStockParams{
		AfterID:     f.AfterID,
		RowLimit:    f.Limit,
		ProductID:   f.ProductID,
		WarehouseID: f.WarehouseID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.StockItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.StockItem{
			ID:            row.ID,
			ProductID:     derefInt64(row.ProductID),
			WarehouseID:   derefInt64(row.WarehouseID),
			Quantity:      numericToQuantity(row.StockQuantity),
			MinStock:      derefInt32(row.MinStock),
			MaxStock:      derefInt32(row.MaxStock),
			ShelfLocation: derefStr(row.ShelfLocation),
			UpdatedAt:     row.UpdatedAt.Time,
		})
	}
	return out, nil
}

func (r *Repo) ListBatchesFEFO(ctx context.Context, productID int64, warehouseID *int64) ([]domain.Batch, error) {
	rows, err := r.q.ListBatchesFEFO(ctx, appdb.ListBatchesFEFOParams{
		ProductID:   productID,
		WarehouseID: warehouseID,
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

// ---- mapping row <-> domain ----

// numericToMoney chuyển pgtype.Numeric sang money.Money KHÔNG đi qua float: dựng
// decimal trực tiếp từ mantissa (big.Int) * 10^Exp. NULL/NaN → Zero.
func numericToMoney(n pgtype.Numeric) money.Money {
	if !n.Valid || n.NaN || n.Int == nil {
		return money.Zero()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return money.FromDecimal(d)
}

// numericToQuantity chuyển pgtype.Numeric (cột quantity/stock_quantity NUMERIC)
// sang domain.Quantity KHÔNG đi qua float (giữ chính xác). NULL/NaN → 0.
func numericToQuantity(n pgtype.Numeric) domain.Quantity {
	if !n.Valid || n.NaN || n.Int == nil {
		return domain.ZeroQty()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return domain.QuantityFromDecimal(d)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt32(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func derefInt64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// Đảm bảo Repo thoả port domain ở compile-time.
var _ domain.Repository = (*Repo)(nil)
