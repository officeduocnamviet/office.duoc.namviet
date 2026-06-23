// Package app — tầng use-case của inventory. LÁT ĐỌC giữ nguyên (Service mỏng,
// bind pool). LÁT GHI (P2 — trừ kho FEFO) thêm ở đây: port GHI bound-tx
// (StockWriter) + TxManager + use-case Deductor. Mở/điều phối transaction ở app;
// domain THUẦN không thấy tx (arch_test chặn). Mẫu theo accounting (Poster).
package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

// StockWriter là PORT GHI tồn kho bound tới MỘT transaction (do caller hoặc
// TxManager truyền vào). Adapter postgres implement bằng appdb.Queries.WithTx(tx).
// Là port ở TẦNG APP (không domain) vì gắn pgx.Tx + điều phối khoá/transaction —
// domain THUẦN không được biết pgx (arch_test chặn). Mọi thao tác chạy trong CÙNG
// tx để gộp atomic với nghiệp vụ orders/POS (trừ kho + post sổ + ghi tiền cùng
// rollback nếu lỗi).
type StockWriter interface {
	// LockWarehouseProduct lấy pg_advisory_xact_lock theo (warehouse, product) —
	// giữ tới hết tx, tuần tự hoá MỌI giao dịch trừ kho cùng (kho,sp). PHẢI gọi ĐẦU
	// use-case trước khi đọc tồn, để chống bán âm khi đồng thời.
	LockWarehouseProduct(ctx context.Context, warehouseID, productID int64) error
	// ListBatchesForDeductFEFO trả các lô còn tồn của (warehouse, product) theo FEFO
	// (expiry ASC), đã FOR UPDATE giữ dòng trong tx. Đã lọc lô deleted + tồn 0. Hết
	// → slice rỗng.
	ListBatchesForDeductFEFO(ctx context.Context, warehouseID, productID int64) ([]domain.Batch, error)
	// DeductBatch trừ qty khỏi MỘT dòng tồn-theo-lô (inventory_batches.id). Guard DB
	// quantity >= qty đảm bảo không ghi âm dù race (phòng thủ).
	DeductBatch(ctx context.Context, inventoryBatchID int64, qty domain.Quantity) error
	// DeductTotal trừ qty khỏi tồn TỔNG (product_inventory.stock_quantity) của
	// (warehouse, product) — đúng tổng đã trừ qua các lô.
	DeductTotal(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) error

	// ---- NHẬP KHO (StockIn — đối xứng trừ kho, dùng cho purchasing) ----

	// InsertBatch tạo MỘT lô mới (public.batches) với inbound_price = giá nhập
	// per-unit. Trả batchID để dòng tồn-theo-lô tham chiếu. expiry/mfg nullable.
	InsertBatch(ctx context.Context, productID int64, batchCode string, expiry, mfg *time.Time, inboundPrice money.Money) (batchID int64, err error)
	// InsertInventoryBatch tạo dòng tồn-theo-lô (public.inventory_batches) cho
	// (warehouse, product, batch) vừa tạo với quantity = lượng nhập.
	InsertInventoryBatch(ctx context.Context, warehouseID, productID, batchID int64, qty domain.Quantity) error
	// AddStockTotal cộng dồn tồn TỔNG (product_inventory.stock_quantity) cho
	// (warehouse, product) ĐÃ CÓ dòng. Trả số dòng đổi: 0 = chưa có dòng → caller
	// InsertStockTotal tạo mới (UPSERT thủ công, bảng không có unique (kho,sp)).
	AddStockTotal(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) (rows int64, err error)
	// InsertStockTotal tạo dòng tồn tổng mới cho (warehouse, product) CHƯA CÓ.
	InsertStockTotal(ctx context.Context, warehouseID, productID int64, qty domain.Quantity) error
}

// StockWriterFromTx dựng một StockWriter bound tới tx. TxManager / Deductor dùng để
// lấy writer cho transaction hiện hành. Tách thành func để app không phụ thuộc
// cứng vào cách khởi tạo repo (adapter cung cấp).
type StockWriterFromTx func(tx pgx.Tx) StockWriter

// TxManager mở/commit một transaction cho trường hợp trừ kho ĐỘC LẬP
// (DeductFEFOInOwnTx). Adapter implement bằng platform/db.WithinTx. Khi orders/POS
// trừ kho trong tx nghiệp vụ của HỌ, chúng gọi thẳng Deductor.DeductFEFO(ctx, tx,
// ...) (KHÔNG qua TxManager) để gộp atomic.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
