// orchestration_ports.go — P4b: PORT cho luồng GHI gộp cross-module (ShipOrder /
// RecordPayment / CreatePosSale). orders/app ĐƯỢC PHÉP import port + domain type
// của inventory/accounting/finance/vat (monolith, tầng app — CLAUDE.md §8). Ta
// định nghĩa interface CONSUMER hẹp ở đây (khớp đúng chữ ký port thật của các
// module) để: (1) test bằng fake, (2) không ràng buộc cứng struct cụ thể. Domain
// orders GIỮ THUẦN — các port này gắn pgx.Tx nên ở TẦNG APP.
package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	accountingdomain "github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	financeapp "github.com/Maneva-AI/namviet-backend/internal/finance/app"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	inventorydomain "github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	vatapp "github.com/Maneva-AI/namviet-backend/internal/vat/app"
	vatdomain "github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

// Deductor là port nội bộ trừ kho FEFO (khớp inventory.DeductPort). orders gọi
// trong tx của họ để gộp atomic; trả các lô tiêu thụ kèm giá vốn (per-unit) cho COGS.
type Deductor interface {
	DeductFEFO(ctx context.Context, tx pgx.Tx, warehouseID, productID int64, qty inventorydomain.Quantity) ([]inventorydomain.ConsumedBatch, error)
}

// Poster là port nội bộ post bút toán kép (khớp accounting.Poster). orders post
// từng JournalEntry (đã cân) trong tx của họ — sổ luôn khớp sự kiện, rollback cùng.
type Poster interface {
	Post(ctx context.Context, tx pgx.Tx, e accountingdomain.JournalEntry) (string, error)
}

// InvoiceIssuer là port nội bộ phát hành HĐ VAT (khớp vat.IssuePort). orders gọi
// trong tx giao hàng để gộp atomic với post sổ TAX.
type InvoiceIssuer interface {
	IssueInvoice(ctx context.Context, tx pgx.Tx, p vatapp.IssueParams) (vatdomain.IssuedInvoice, error)
}

// PaymentRecorder là port nội bộ ghi phiếu THU idempotent (khớp finance.RecordPort).
// orders/POS gọi trong tx của họ để gộp atomic với trừ kho + post sổ.
type PaymentRecorder interface {
	RecordPaymentIn(ctx context.Context, tx pgx.Tx, p financeapp.RecordPaymentInParams) (financedomain.Payment, bool, error)
}

// OrchHeader là header đơn đã KHOÁ (FOR UPDATE) cần cho điều phối: mã đơn (ref/HĐ),
// loại đơn (B2B/B2C), khách (cho MST/công nợ), trạng thái xử lý, tổng phải trả.
type OrchHeader struct {
	ID            string
	Code          string
	CustomerID    *int64
	OrderType     string // B2B | B2C
	Status        string // trạng thái xử lý hiện tại (PENDING/CONFIRMED/...)
	FinalAmount   money.Money
	TotalAmount   money.Money
	PaymentStatus string
}

// OrchestrationStore là PORT GHI orders BOUND-TX riêng cho P4b: khoá+đọc header
// đầy đủ, đọc dòng hàng (cho trừ kho + HĐ), tính lại đã-thu trong tx, cập nhật
// payment_status + status. Tách khỏi OrderStore (P4a tạo đơn/state machine đơn
// giản) để giữ port mỗi cái một việc rõ ràng. Mọi thao tác chạy trong CÙNG tx.
type OrchestrationStore interface {
	// GetHeaderForUpdate khoá dòng đơn (FOR UPDATE) + trả header đầy đủ. Không thấy
	// → (_, false, nil). Lỗi hạ tầng → (_, _, err).
	GetHeaderForUpdate(ctx context.Context, orderID string) (OrchHeader, bool, error)
	// ListLinesForOrder trả các dòng hàng (chưa soft-delete) của đơn — product/qty/
	// uom/đơn giá để trừ kho FEFO + dựng dòng HĐ. Hết → slice rỗng.
	ListLinesForOrder(ctx context.Context, orderID string) ([]OrchLine, error)
	// SumPaidInTx trả tổng ĐÃ THU (sổ thực tế) của đơn theo code TRONG tx hiện hành
	// (phiếu vừa ghi cùng tx ĐƯỢC tính). Dùng tính lại payment_status.
	SumPaidInTx(ctx context.Context, orderCode string) (money.Money, error)
	// UpdatePaymentStatus đặt payment_status (unpaid/partial/paid). Trả số dòng đổi.
	UpdatePaymentStatus(ctx context.Context, orderID, paymentStatus string) (rows int64, err error)
	// UpdateStatus đổi trạng thái xử lý (guard status cũ). Trả số dòng đổi — 0 nghĩa
	// là đơn đã đổi bởi luồng khác (service map Conflict).
	UpdateStatus(ctx context.Context, orderID, expected, next string) (rows int64, err error)
	// ListUnpaidOrdersByCustomer trả đơn CHƯA tất toán của khách, CŨ NHẤT trước,
	// KHOÁ FOR UPDATE — để phân bổ phiếu lump-sum tuần tự (mục 55). Kèm Final + đã-thu
	// hiện tại (tính phần còn thiếu). Hết → slice rỗng.
	ListUnpaidOrdersByCustomer(ctx context.Context, customerID int64) ([]UnpaidOrder, error)
	// InsertAllocation ghi/cộng dồn MỘT dòng phân bổ phiếu (paymentID) → đơn
	// (orderCode) số tiền amount. Trong CÙNG tx với ghi phiếu + cập nhật payment_status.
	InsertAllocation(ctx context.Context, paymentID int64, orderCode string, amount money.Money) error
}

// UnpaidOrder là một đơn chưa tất toán (cho phân bổ lump-sum): mã + tổng phải trả +
// đã-thu hiện tại → phần còn thiếu = Final - Paid.
type UnpaidOrder struct {
	ID    string
	Code  string
	Final money.Money
	Paid  money.Money
}

// OrchLine là một dòng hàng đã đọc để điều phối (trừ kho + HĐ). Tiền/lượng decimal.
type OrchLine struct {
	ProductID int64
	Quantity  inventorydomain.Quantity
	UOM       string
	UnitPrice money.Money
	Discount  money.Money
	LineTotal money.Money
}

// OrchestrationStoreFromTx dựng OrchestrationStore bound tới tx (adapter cấp).
type OrchestrationStoreFromTx func(tx pgx.Tx) OrchestrationStore
