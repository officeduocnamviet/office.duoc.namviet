package app

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

// Deductor là use-case GHI của inventory: trừ kho theo FEFO có khoá tranh chấp.
// Nó implement port nội bộ DeductPort (orders/POS gọi DeductFEFO trong tx của họ
// để gộp atomic). KHÔNG có REST POST công khai ở P2 — giống accounting.Poster.
type Deductor struct {
	writerFromTx StockWriterFromTx
	txm          TxManager
}

// NewDeductor dựng Deductor. writerFromTx bind StockWriter tới một tx; txm để
// DeductFEFOInOwnTx mở tx riêng. Có thể truyền nil cho thành phần không dùng.
func NewDeductor(writerFromTx StockWriterFromTx, txm TxManager) *Deductor {
	return &Deductor{writerFromTx: writerFromTx, txm: txm}
}

// DeductFEFO trừ qty đơn vị product khỏi warehouse theo FEFO, TRONG transaction tx
// do CALLER truyền (gộp atomic với nghiệp vụ orders/POS). Quy trình:
//  1. pg_advisory_xact_lock theo (warehouse, product) — tuần tự hoá trừ kho cùng
//     (kho,sp), chống bán âm khi 2 giao dịch đồng thời (khoá giữ tới hết tx).
//  2. Đọc các lô còn tồn theo FEFO (expiry ASC) FOR UPDATE.
//  3. Lập kế hoạch tiêu thụ THUẦN (domain.PlanFEFO) — thiếu → Conflict (KHÔNG âm).
//  4. Áp kế hoạch: trừ từng inventory_batches + trừ product_inventory tổng.
//
// Trả danh sách ConsumedBatch (kèm InboundPrice mỗi lô — orders post COGS sau).
// qty <= 0 → no-op, trả nil (không lỗi). Lỗi đã là apperr (map envelope ở http).
func (d *Deductor) DeductFEFO(ctx context.Context, tx pgx.Tx, warehouseID, productID int64, qty domain.Quantity) ([]domain.ConsumedBatch, error) {
	if !qty.IsPositive() {
		return nil, nil // trừ 0/âm = no-op hợp lệ
	}

	w := d.writerFromTx(tx)

	// 1. Khoá tranh chấp ĐẦU TX — tuần tự hoá trừ kho cùng (kho,sp).
	if err := w.LockWarehouseProduct(ctx, warehouseID, productID); err != nil {
		return nil, apperr.Internal("khoá tồn kho lỗi").WithCause(err)
	}

	// 2. Đọc lô còn tồn theo FEFO, giữ dòng (FOR UPDATE) trong tx.
	available, err := w.ListBatchesForDeductFEFO(ctx, warehouseID, productID)
	if err != nil {
		return nil, apperr.Internal("đọc lô tồn kho lỗi").WithCause(err)
	}

	// 3. Lập kế hoạch tiêu thụ THUẦN (domain) — thiếu tồn → Conflict, KHÔNG âm.
	plan, err := domain.PlanFEFO(available, qty)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientStock) {
			return nil, apperr.Conflict("tồn kho không đủ để trừ (cần " + qty.String() + ")")
		}
		return nil, apperr.Internal("lập kế hoạch trừ kho lỗi").WithCause(err)
	}

	// 4. Áp kế hoạch: trừ từng lô (DB guard quantity >= qty) + trừ tồn tổng.
	for _, c := range plan {
		if err := w.DeductBatch(ctx, c.InventoryBatchID, c.Quantity); err != nil {
			return nil, apperr.Internal("trừ tồn lô lỗi").WithCause(err)
		}
	}
	if err := w.DeductTotal(ctx, warehouseID, productID, qty); err != nil {
		return nil, apperr.Internal("trừ tồn tổng lỗi").WithCause(err)
	}
	return plan, nil
}

// DeductFEFOInOwnTx trừ kho trong một transaction RIÊNG (mở/commit ở đây). Dùng
// khi trừ kho độc lập, không gộp với nghiệp vụ khác (vd test, công cụ điều chỉnh).
func (d *Deductor) DeductFEFOInOwnTx(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) ([]domain.ConsumedBatch, error) {
	if d.txm == nil {
		return nil, apperr.Internal("TxManager chưa cấu hình cho DeductFEFOInOwnTx")
	}
	var consumed []domain.ConsumedBatch
	err := d.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		c, err := d.DeductFEFO(ctx, tx, warehouseID, productID, qty)
		if err != nil {
			return err
		}
		consumed = c
		return nil
	})
	if err != nil {
		return nil, err
	}
	return consumed, nil
}

// Đảm bảo Deductor thoả DeductPort ở compile-time (port nội bộ cho orders/POS).
var _ DeductPort = (*Deductor)(nil)

// DeductPort là PORT NỘI BỘ: module orders/POS trừ kho FEFO trong tx nghiệp vụ
// của họ qua đây (gộp atomic với post sổ + ghi tiền). KHÔNG có REST POST công khai
// ở P2. Định nghĩa ở app (không domain) vì nhận pgx.Tx — domain THUẦN không biết tx.
type DeductPort interface {
	DeductFEFO(ctx context.Context, tx pgx.Tx, warehouseID, productID int64, qty domain.Quantity) ([]domain.ConsumedBatch, error)
}
