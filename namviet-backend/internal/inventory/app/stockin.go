package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

// StockInner là use-case GHI NHẬP KHO của inventory: tạo lô + tăng tồn (đối xứng
// Deductor/DeductFEFO). Nó implement port nội bộ StockInPort (purchasing gọi StockIn
// trong tx của họ để gộp atomic với post sổ + chi NCC). KHÔNG có REST POST công khai
// — giống Deductor / accounting.Poster.
type StockInner struct {
	writerFromTx StockWriterFromTx
	txm          TxManager
}

// NewStockInner dựng StockInner. writerFromTx bind StockWriter tới một tx; txm để
// StockInOwnTx mở tx riêng. Có thể truyền nil cho thành phần không dùng.
func NewStockInner(writerFromTx StockWriterFromTx, txm TxManager) *StockInner {
	return &StockInner{writerFromTx: writerFromTx, txm: txm}
}

// StockIn NHẬP qty đơn vị product vào warehouse như MỘT LÔ mới, TRONG transaction tx
// do CALLER truyền (gộp atomic với nghiệp vụ purchasing). Quy trình:
//  1. pg_advisory_xact_lock theo (warehouse, product) — tuần tự hoá thay đổi tồn
//     cùng (kho,sp) (đồng nhất DeductFEFO; chống đua id COALESCE(max+1) khi nhập song
//     song cùng product).
//  2. INSERT public.batches (inbound_price = giá nhập per-unit) → batchID.
//  3. INSERT public.inventory_batches (warehouse,product,batch,quantity=qty).
//  4. UPSERT public.product_inventory (cộng dồn nếu có dòng (kho,sp); tạo mới nếu chưa).
//
// inboundPrice là GIÁ NHẬP per-unit (= unit_cost dòng PO) — KHÔNG nhân số lượng (đối
// xứng ConsumedBatch.InboundPrice của DeductFEFO; COGS sau này = qty × inbound_price).
// Trả batchID lô vừa tạo. qty <= 0 → Validation (không nhập 0/âm). Lỗi đã là apperr.
func (s *StockInner) StockIn(
	ctx context.Context,
	tx pgx.Tx,
	warehouseID, productID int64,
	batchCode string,
	expiry, mfg *time.Time,
	qty domain.Quantity,
	inboundPrice money.Money,
) (int64, error) {
	if !qty.IsPositive() {
		return 0, apperr.Validation("số lượng nhập phải lớn hơn 0")
	}
	if warehouseID <= 0 || productID <= 0 {
		return 0, apperr.Validation("warehouse_id/product_id không hợp lệ (phải > 0)")
	}

	w := s.writerFromTx(tx)

	// 1. Khoá tranh chấp ĐẦU TX — tuần tự hoá thay đổi tồn cùng (kho,sp).
	if err := w.LockWarehouseProduct(ctx, warehouseID, productID); err != nil {
		return 0, apperr.Internal("khoá tồn kho lỗi").WithCause(err)
	}

	// 2. Tạo lô mới (inbound_price = giá nhập per-unit).
	batchID, err := w.InsertBatch(ctx, productID, batchCode, expiry, mfg, inboundPrice)
	if err != nil {
		return 0, apperr.Internal("tạo lô nhập kho lỗi").WithCause(err)
	}

	// 3. Tồn theo lô tại kho.
	if err := w.InsertInventoryBatch(ctx, warehouseID, productID, batchID, qty); err != nil {
		return 0, apperr.Internal("ghi tồn lô lỗi").WithCause(err)
	}

	// 4. Tồn tổng — UPSERT thủ công: cộng dồn nếu đã có dòng, ngược lại tạo mới.
	rows, err := w.AddStockTotal(ctx, warehouseID, productID, qty)
	if err != nil {
		return 0, apperr.Internal("cộng tồn tổng lỗi").WithCause(err)
	}
	if rows == 0 {
		if err := w.InsertStockTotal(ctx, warehouseID, productID, qty); err != nil {
			return 0, apperr.Internal("tạo tồn tổng lỗi").WithCause(err)
		}
	}
	return batchID, nil
}

// StockInOwnTx nhập kho trong một transaction RIÊNG (mở/commit ở đây). Dùng khi nhập
// kho độc lập, không gộp với nghiệp vụ khác (vd công cụ điều chỉnh, test).
func (s *StockInner) StockInOwnTx(
	ctx context.Context,
	warehouseID, productID int64,
	batchCode string,
	expiry, mfg *time.Time,
	qty domain.Quantity,
	inboundPrice money.Money,
) (int64, error) {
	if s.txm == nil {
		return 0, apperr.Internal("TxManager chưa cấu hình cho StockInOwnTx")
	}
	var batchID int64
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		id, err := s.StockIn(ctx, tx, warehouseID, productID, batchCode, expiry, mfg, qty, inboundPrice)
		if err != nil {
			return err
		}
		batchID = id
		return nil
	})
	if err != nil {
		return 0, err
	}
	return batchID, nil
}

// Đảm bảo StockInner thoả StockInPort ở compile-time (port nội bộ cho purchasing).
var _ StockInPort = (*StockInner)(nil)

// StockInPort là PORT NỘI BỘ: module purchasing nhập kho trong tx nghiệp vụ của họ
// qua đây (gộp atomic với post sổ + chi NCC). KHÔNG có REST POST công khai. Định
// nghĩa ở app (không domain) vì nhận pgx.Tx — domain THUẦN không biết tx.
type StockInPort interface {
	StockIn(ctx context.Context, tx pgx.Tx, warehouseID, productID int64, batchCode string, expiry, mfg *time.Time, qty domain.Quantity, inboundPrice money.Money) (batchID int64, err error)
}
