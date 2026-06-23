// Package app là tầng use-case của catalog. Catalog read-only nên app MỎNG: nó
// chuẩn hoá input (limit/cursor/status), gọi port domain.Repository, và ghép kết
// quả (vd product + units). Không mở transaction (đọc thuần). Theo ARCHITECTURE.md
// §3 module nhẹ được phép gộp logic mỏng — nhưng vẫn giữ domain thuần + port.
package app

import (
	"context"

	"github.com/Maneva-AI/namviet-backend/internal/catalog/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
)

const (
	// defaultLimit là số sản phẩm mỗi trang khi client không truyền.
	defaultLimit = 20
	// maxLimit chặn trên để tránh quét toàn bảng trong một request.
	maxLimit = 100
)

// Service là use-case đọc catalog mà edge (HTTP) dùng.
type Service struct {
	repo domain.Repository
}

// New dựng Service từ một repository (port domain).
func New(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// ListProductsQuery là input đã được edge giải mã (cursor → AfterID ở handler,
// nhưng ta nhận cursor thô để app tự decode, giữ quy ước pagination ở một nơi).
type ListProductsQuery struct {
	Cursor     string
	Limit      int32
	Status     string
	CategoryID *int64
	Query      string
}

// ListProductsResult là một trang sản phẩm + cursor trang kế (rỗng nếu hết).
type ListProductsResult struct {
	Items      []domain.Product
	NextCursor string
}

// ListProducts trả một trang sản phẩm theo keyset. Tự decode cursor, chuẩn hoá
// limit, và sinh NextCursor = id của phần tử cuối nếu trang đầy (có thể còn nữa).
func (s *Service) ListProducts(ctx context.Context, q ListProductsQuery) (ListProductsResult, error) {
	afterID, err := pagination.DecodeID(q.Cursor)
	if err != nil {
		return ListProductsResult{}, apperr.Validation("cursor không hợp lệ")
	}
	limit := normalizeLimit(q.Limit)

	items, err := s.repo.ListProducts(ctx, domain.ProductFilter{
		AfterID:    afterID,
		Limit:      limit,
		Status:     q.Status,
		CategoryID: q.CategoryID,
		Query:      q.Query,
	})
	if err != nil {
		return ListProductsResult{}, err
	}

	res := ListProductsResult{Items: items}
	// Trang đầy → có thể còn trang sau; phát cursor từ id phần tử cuối. Trang
	// chưa đầy → hết dữ liệu, NextCursor rỗng.
	if int32(len(items)) == limit && limit > 0 {
		res.NextCursor = pagination.EncodeID(items[len(items)-1].ID)
	}
	return res, nil
}

// ProductDetail gom product + các đơn vị tính của nó.
type ProductDetail struct {
	Product domain.Product
	Units   []domain.ProductUnit
}

// GetProduct trả chi tiết một sản phẩm kèm đơn vị tính. Không thấy product →
// apperr.NotFound (từ repo).
func (s *Service) GetProduct(ctx context.Context, id int64) (ProductDetail, error) {
	p, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}
	units, err := s.repo.ListUnits(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}
	return ProductDetail{Product: p, Units: units}, nil
}

// ListCategories trả danh mục (mặc định chỉ status 'active' nếu onlyActive).
func (s *Service) ListCategories(ctx context.Context, status string) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx, status)
}

// ListManufacturers trả hãng / nhà sản xuất.
func (s *Service) ListManufacturers(ctx context.Context, status string) ([]domain.Manufacturer, error) {
	return s.repo.ListManufacturers(ctx, status)
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
