// Package app — phần GHI: điều phối use-case tạo đơn + đổi trạng thái ĐƠN GIẢN
// (P4a). Mở/commit transaction ở đây (platform/db.WithinTx) và gọi port OrderStore
// bound-tx. Domain không thấy tx (arch_test chặn). Mẫu theo accounting/finance/vat
// (storeFromTx + TxManager). KHÔNG đụng kho/tiền/sổ — ShipOrder/RecordPayment/POS
// = P4b.
package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

// CreatedOrder là kết quả tạo đơn: đơn (master) + các dòng đã ghi. Trả cho HTTP
// để dựng response.
type CreatedOrder struct {
	Order domain.Order
	Lines []domain.OrderLine
}

// OrderStore là PORT GHI bound tới MỘT transaction (do caller/TxManager truyền).
// Adapter postgres implement bằng appdb.Queries.WithTx(tx). Là port TẦNG APP
// (không domain) vì gắn pgx.Tx + sinh mã (sequence) + FOR UPDATE — domain THUẦN
// không biết tx (arch_test chặn). Mọi thao tác chạy trong CÙNG tx (atomic).
type OrderStore interface {
	// FindByIdemKey tra cứu đơn đã tạo theo Idempotency-Key (chống tạo trùng). Có
	// → (orderID, true, nil); không có → ("", false, nil). Lỗi hạ tầng → (_,_,err).
	FindByIdemKey(ctx context.Context, idemKey string) (orderID string, found bool, err error)
	// NextCodeSeq cấp số tăng dần kế tiếp cho mã đơn (app ghép tiền tố + zero-pad).
	NextCodeSeq(ctx context.Context) (int64, error)
	// InsertOrder ghi header đơn (status PENDING). id/code do app sinh. Trùng code
	// (UNIQUE) → duplicate=true để service xử lý (đề phòng đua sinh mã). Trả đơn đã
	// ghi (đường đọc dựng từ đây).
	InsertOrder(ctx context.Context, o NewOrderRow) (saved domain.Order, duplicate bool, err error)
	// InsertItem ghi một dòng hàng (id app sinh).
	InsertItem(ctx context.Context, orderID string, l domain.ComputedLine) error
	// BindIdemKey ghi ánh xạ Idempotency-Key → đơn. ON CONFLICT DO NOTHING → trả
	// inserted=false nếu key đã tồn tại (luồng khác thắng đua).
	BindIdemKey(ctx context.Context, idemKey, orderID, orderCode string) (inserted bool, err error)
	// GetForUpdate khoá dòng đơn (FOR UPDATE) + trả status hiện tại. Không thấy →
	// (_, false, nil). Lỗi hạ tầng → (_,_,err).
	GetForUpdate(ctx context.Context, orderID string) (current domain.Status, found bool, err error)
	// UpdateStatus đổi trạng thái (guard status cũ trong WHERE). Trả số dòng đổi —
	// 0 nghĩa là đơn đã đổi trạng thái bởi luồng khác (service map Conflict).
	UpdateStatus(ctx context.Context, orderID string, expected, next domain.Status) (rows int64, err error)
	// GetCreated nạp lại đơn (master + lines) theo id để trả về sau khi tạo (kèm
	// created_at/updated_at do DB sinh). Dùng cả cho idempotent hit (đọc đơn cũ).
	GetCreated(ctx context.Context, orderID string) (CreatedOrder, error)
}

// NewOrderRow là dữ liệu header để INSERT một đơn (id/code app đã sinh + tổng tiền
// đã tính ở domain). Tách struct để port không lộ chi tiết appdb.
type NewOrderRow struct {
	ID    string
	Code  string
	Draft domain.Draft
}

// OrderStoreFromTx dựng OrderStore bound tới tx. TxManager/WriteService dùng để
// lấy store cho transaction hiện hành.
type OrderStoreFromTx func(tx pgx.Tx) OrderStore

// TxManager mở/commit một transaction (platform/db.WithinTx). Một use-case GHI =
// một transaction.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
