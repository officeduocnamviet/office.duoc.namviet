package domain

import "context"

// OrderFilter gom điều kiện lọc + keyset pagination cho danh sách đơn. Tất cả
// tiêu chí lọc là optional (nil/rỗng = không lọc theo tiêu chí đó). Limit/AfterID
// do tầng app chuẩn hoá. Repo luôn ngầm lọc deleted_at IS NULL (soft-delete).
//
// Keyset theo orders.created_at + id: orders.id là UUID nên KHÔNG keyset bằng id
// tăng dần được như bảng bigint. Ta keyset theo (created_at DESC, id DESC) —
// con trỏ mã hoá thời điểm + id của bản ghi cuối trang trước (xem app).
type OrderFilter struct {
	// AfterCreatedAt + AfterID: lấy đơn "cũ hơn" mốc này (created_at < mốc, hoặc
	// created_at = mốc và id < AfterID). Zero = trang đầu.
	AfterCreatedAt int64 // unix nano của created_at bản ghi cuối trang trước; 0 = trang đầu
	AfterID        string
	// Limit: số bản ghi tối đa mỗi trang (đã chuẩn hoá ở app, > 0).
	Limit int32
	// CustomerID lọc theo khách; nil = mọi khách.
	CustomerID *int64
	// Status lọc theo trạng thái xử lý đơn; rỗng = mọi trạng thái.
	Status string
	// PaymentStatus lọc theo trạng thái thanh toán ('unpaid'|'partial'|'paid');
	// rỗng = tất cả.
	PaymentStatus string
	// FromDate/ToDate lọc theo khoảng created_at (unix nano, đã chuẩn hoá ở app);
	// 0 = không chặn đầu/cuối tương ứng.
	FromDate int64
	ToDate   int64
}

// Repository là PORT đọc dữ liệu orders do domain ĐỊNH NGHĨA; adapter postgres
// implement ("accept interfaces, return structs"). LÁT NÀY CHỈ ĐỌC nên port chỉ
// có thao tác đọc — KHÔNG Insert/Update/Delete (tạo đơn/trừ kho/ghi phiếu thu
// HOÃN sang slice sau). Các thao tác đọc độc lập, không cần transaction nên app
// gọi thẳng repo bind pool.
type Repository interface {
	// ListOrders trả tối đa f.Limit đơn theo keyset (created_at DESC, id DESC),
	// đã lọc deleted_at IS NULL + các tiêu chí trong f. Mỗi đơn KÈM PaymentSummary
	// (đã thu/còn nợ suy diễn từ finance_transactions). Hết → slice rỗng.
	ListOrders(ctx context.Context, f OrderFilter) ([]Order, error)
	// GetOrderByID trả một đơn theo id (uuid, chưa soft-delete) kèm PaymentSummary.
	// Không thấy → apperr.NotFound.
	GetOrderByID(ctx context.Context, id string) (Order, error)
	// ListLines trả các dòng hàng (chưa soft-delete) của một đơn theo thứ tự ổn
	// định (created_at, id). Không có → slice rỗng (không lỗi).
	ListLines(ctx context.Context, orderID string) ([]OrderLine, error)
}
