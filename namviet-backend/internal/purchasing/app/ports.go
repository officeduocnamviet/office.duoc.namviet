// Package app là tầng use-case của purchasing (mua hàng & nhập kho, mục 54 — chiều
// MUA, đối xứng orders chiều BÁN). Điều phối GHI cross-module trong 1 transaction
// ATOMIC: ghi PO (bảng app) + nhập kho (inventory.StockIn) + post sổ
// (accounting.Poster) + chi NCC (finance.RecordOutPort). Mở/commit tx ở đây qua
// platform/db.WithinTx (TxManager); cross-module nhận CÙNG tx → lỗi bất kỳ bước nào
// ROLLBACK CẢ CỤM. Domain THUẦN không thấy tx (arch_test chặn).
package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	accountingdomain "github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	financeapp "github.com/Maneva-AI/namviet-backend/internal/finance/app"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	inventorydomain "github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

// TxManager mở/commit một transaction cho mỗi use-case GHI. Adapter implement bằng
// platform/db.WithinTx.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

// ---- Cross-module ports (purchasing/app ĐƯỢC import port + domain type của
// inventory/accounting/finance — monolith tầng app, CLAUDE.md §8). Định nghĩa
// interface CONSUMER hẹp ở đây (khớp đúng chữ ký port thật) để test bằng fake. ----

// StockInner là port nội bộ NHẬP KHO (khớp inventory.StockInPort). purchasing gọi
// trong tx của họ để gộp atomic; tạo lô (inbound_price = unit_cost) + tăng tồn.
type StockInner interface {
	StockIn(ctx context.Context, tx pgx.Tx, warehouseID, productID int64, batchCode string, expiry, mfg *time.Time, qty inventorydomain.Quantity, inboundPrice money.Money) (batchID int64, err error)
}

// Poster là port nội bộ post bút toán kép (khớp accounting.Poster). purchasing post
// từng JournalEntry (đã cân) trong tx của họ.
type Poster interface {
	Post(ctx context.Context, tx pgx.Tx, e accountingdomain.JournalEntry) (string, error)
}

// SupplierPayer là port nội bộ ghi phiếu CHI NCC idempotent (khớp finance.RecordOutPort).
// purchasing gọi trong tx của họ để gộp atomic với post sổ chi.
type SupplierPayer interface {
	RecordPaymentOut(ctx context.Context, tx pgx.Tx, p financeapp.RecordPaymentOutParams) (financedomain.Payment, bool, error)
}

// ---- Store ports (GHI/ĐỌC bảng app.purchase_orders, bound-tx) ----

// POHeader là header PO đã KHOÁ (FOR UPDATE) cần cho điều phối confirm/receive/pay.
type POHeader struct {
	ID           string
	Code         string
	SupplierID   *int64
	SupplierName string
	Status       string
	TotalAmount  money.Money
	VATAmount    money.Money
	Note         string
	LockVersion  int32
}

// POLine là một dòng PO đã đọc (cho nhập kho + post sổ). Quantity decimal qua
// inventory/domain.Quantity (cùng nền decimal — StockIn nhận kiểu này).
type POLine struct {
	ID                string
	LineNo            int
	ProductID         int64
	Quantity          inventorydomain.Quantity
	UnitCost          money.Money
	VATRate           money.Money // dùng money để giữ decimal (vat_rate scale-4); chỉ nhân, không cộng tiền
	BatchCode         string
	ExpiryDate        *time.Time
	ManufacturingDate *time.Time
	LineTotal         money.Money
}

// Store là PORT GHI/ĐỌC purchasing bound-tx. Adapter postgres implement bằng
// appdb.Queries.WithTx(tx). Mọi thao tác chạy trong CÙNG tx (gộp atomic với nhập kho
// + post sổ + chi NCC).
type Store interface {
	// NextCodeSeq cấp số mã PO (sequence). InsertPO ghi header (status draft) — trùng
	// code (đua sinh mã) → duplicate=true để service cấp số mới.
	NextCodeSeq(ctx context.Context) (int64, error)
	InsertPO(ctx context.Context, row NewPORow) (domain.PurchaseOrder, bool, error)
	InsertPOItem(ctx context.Context, poID string, l domain.ComputedLine) error
	// GetCreated nạp lại PO (header + lines) theo id — readback sau ghi / GET detail.
	GetCreated(ctx context.Context, poID string) (CreatedPO, error)
	// GetHeaderForUpdate khoá dòng PO (FOR UPDATE) + trả header. Không thấy → (_, false, nil).
	GetHeaderForUpdate(ctx context.Context, poID string) (POHeader, bool, error)
	// ListLines trả dòng hàng của PO (theo line_no) — nhập kho + post sổ. Hết → rỗng.
	ListLines(ctx context.Context, poID string) ([]POLine, error)
	// UpdateStatus đổi trạng thái (guard status cũ + bump lock_version). Trả số dòng đổi.
	UpdateStatus(ctx context.Context, poID, expected, next string) (rows int64, err error)
	// FindByIdemKey / BindIdemKey: idempotency tạo PO (1 key → 1 PO).
	FindByIdemKey(ctx context.Context, idemKey string) (poID string, found bool, err error)
	BindIdemKey(ctx context.Context, idemKey, poID, poCode string) (inserted bool, err error)
}

// NewPORow là tham số INSERT một PO (header) đã có id + code.
type NewPORow struct {
	ID    string
	Code  string
	Draft domain.Draft
}

// CreatedPO là PO vừa tạo/đọc lại (header + lines) cho đường trả về.
type CreatedPO struct {
	PO    domain.PurchaseOrder
	Lines []domain.PurchaseLine
}

// StoreFromTx dựng Store bound tới tx (adapter cấp).
type StoreFromTx func(tx pgx.Tx) Store
