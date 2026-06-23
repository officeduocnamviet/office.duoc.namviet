// Package app là tầng use-case của orders. LÁT NÀY CHỈ ĐỌC nên app MỎNG: nó
// chuẩn hoá input (limit/cursor/filter), gọi port domain.Repository, và sinh
// cursor trang kế. Không mở transaction (đọc thuần). Theo ARCHITECTURE.md §3
// module nhẹ được phép gộp logic mỏng — nhưng vẫn giữ domain thuần + port. Khi
// thêm đường GHI (tạo đơn, state machine 3 trạng thái, trừ kho, ghi phiếu thu,
// post sổ) thì module này LÊN "full" 3 lớp với TxManager + SERIALIZABLE.
//
// HOÃN (slice sau, design đang chốt với user): tạo/sửa/huỷ đơn, áp giá+voucher,
// reserve/trừ kho FEFO, ghi finance_transactions, post journal kế toán.
package app

import (
	"context"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

const (
	// defaultLimit là số đơn mỗi trang khi client không truyền.
	defaultLimit = 20
	// maxLimit chặn trên để tránh quét toàn bảng trong một request.
	maxLimit = 100
)

// Service là use-case đọc orders mà edge (HTTP) dùng.
type Service struct {
	repo domain.Repository
}

// New dựng Service từ một repository (port domain).
func New(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// ListOrdersQuery là input đã được edge giải mã. Cursor thô để app tự decode,
// giữ quy ước pagination ở một nơi. CustomerID nil = không lọc. From/ToDate là
// unix nano (0 = không chặn).
type ListOrdersQuery struct {
	Cursor        string
	Limit         int32
	CustomerID    *int64
	Status        string
	PaymentStatus string
	FromDate      int64
	ToDate        int64
}

// ListOrdersResult là một trang đơn + cursor trang kế (rỗng nếu hết).
type ListOrdersResult struct {
	Items      []domain.Order
	NextCursor string
}

// ListOrders trả một trang đơn theo keyset (created_at DESC, id DESC). Tự decode
// cursor, chuẩn hoá limit, và sinh NextCursor từ (created_at, id) của phần tử
// cuối nếu trang đầy (có thể còn nữa).
func (s *Service) ListOrders(ctx context.Context, q ListOrdersQuery) (ListOrdersResult, error) {
	afterNano, afterID, err := decodeCursor(q.Cursor)
	if err != nil {
		return ListOrdersResult{}, apperr.Validation("cursor không hợp lệ")
	}
	limit := normalizeLimit(q.Limit)

	items, err := s.repo.ListOrders(ctx, domain.OrderFilter{
		AfterCreatedAt: afterNano,
		AfterID:        afterID,
		Limit:          limit,
		CustomerID:     q.CustomerID,
		Status:         q.Status,
		PaymentStatus:  q.PaymentStatus,
		FromDate:       q.FromDate,
		ToDate:         q.ToDate,
	})
	if err != nil {
		return ListOrdersResult{}, err
	}

	res := ListOrdersResult{Items: items}
	if int32(len(items)) == limit && limit > 0 {
		last := items[len(items)-1]
		res.NextCursor = encodeCursor(last.CreatedAt.UnixNano(), last.ID)
	}
	return res, nil
}

// OrderDetail là chi tiết một đơn kèm danh sách dòng hàng.
type OrderDetail struct {
	Order domain.Order
	Lines []domain.OrderLine
}

// GetOrder trả chi tiết một đơn (master + PaymentSummary suy diễn) cùng các dòng
// hàng. Không thấy đơn → apperr.NotFound (từ repo). Đơn không có dòng hàng →
// Lines rỗng (không lỗi).
func (s *Service) GetOrder(ctx context.Context, id string) (OrderDetail, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return OrderDetail{}, err
	}
	lines, err := s.repo.ListLines(ctx, id)
	if err != nil {
		return OrderDetail{}, err
	}
	return OrderDetail{Order: order, Lines: lines}, nil
}

func normalizeLimit(l int32) int32 {
	switch {
	case l <= 0:
		return defaultLimit
	case l > maxLimit:
		return maxLimit
	default:
		return l
	}
}
