// Package app là tầng use-case của customers. Customers read-mostly nên app
// MỎNG: chuẩn hoá input (limit/cursor/type/status), gọi port domain.Repository,
// và sinh cursor trang kế. Không mở transaction (đọc thuần). Theo ARCHITECTURE.md
// §3 module nhẹ được phép gộp logic mỏng — nhưng vẫn giữ domain thuần + port.
//
// DEFER: tạo/sửa khách, giá-theo-khách. Context này chỉ phục vụ master + đọc
// công nợ (ưu tiên nguồn LIVE).
package app

import (
	"context"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
)

const (
	// defaultLimit là số khách mỗi trang khi client không truyền.
	defaultLimit = 20
	// maxLimit chặn trên để tránh quét toàn bảng trong một request.
	maxLimit = 100
)

// Service là use-case đọc customers mà edge (HTTP) dùng.
type Service struct {
	repo domain.Repository
}

// New dựng Service từ một repository (port domain).
func New(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// ListCustomersQuery là input đã được edge giải mã. Cursor thô để app tự decode,
// giữ quy ước pagination ở một nơi. Type rỗng = mọi loại; chỉ chấp nhận B2B/B2C.
type ListCustomersQuery struct {
	Cursor string
	Limit  int32
	Type   string
	Status string
	Query  string
}

// ListCustomersResult là một trang khách + cursor trang kế (rỗng nếu hết).
type ListCustomersResult struct {
	Items      []domain.Customer
	NextCursor string
}

// ListCustomers trả một trang khách theo keyset. Tự decode cursor, chuẩn hoá
// limit, validate Type, và sinh NextCursor = id phần tử cuối nếu trang đầy.
func (s *Service) ListCustomers(ctx context.Context, q ListCustomersQuery) (ListCustomersResult, error) {
	afterID, err := pagination.DecodeID(q.Cursor)
	if err != nil {
		return ListCustomersResult{}, apperr.Validation("cursor không hợp lệ")
	}
	typ, err := normalizeType(q.Type)
	if err != nil {
		return ListCustomersResult{}, err
	}
	limit := normalizeLimit(q.Limit)

	items, err := s.repo.ListCustomers(ctx, domain.CustomerFilter{
		AfterID: afterID,
		Limit:   limit,
		Type:    typ,
		Status:  q.Status,
		Query:   q.Query,
	})
	if err != nil {
		return ListCustomersResult{}, err
	}

	res := ListCustomersResult{Items: items}
	// Trang đầy → có thể còn trang sau; phát cursor từ id phần tử cuối.
	if int32(len(items)) == limit && limit > 0 {
		res.NextCursor = pagination.EncodeID(items[len(items)-1].ID)
	}
	return res, nil
}

// GetCustomer trả chi tiết một khách (master + B2BProfile nếu có + DebtSnapshot
// nguồn live). Không thấy → apperr.NotFound (từ repo).
func (s *Service) GetCustomer(ctx context.Context, id int64) (domain.Customer, error) {
	return s.repo.GetCustomerByID(ctx, id)
}

// normalizeType validate filter loại khách. Rỗng = không lọc. Chỉ chấp nhận
// B2B/B2C (giá trị khác → Validation, tránh quét nhầm/ô nhiễm filter).
func normalizeType(t string) (domain.CustomerType, error) {
	switch domain.CustomerType(t) {
	case "":
		return "", nil
	case domain.TypeB2B:
		return domain.TypeB2B, nil
	case domain.TypeB2C:
		return domain.TypeB2C, nil
	default:
		return "", apperr.Validation("customer_type chỉ nhận B2B hoặc B2C")
	}
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
