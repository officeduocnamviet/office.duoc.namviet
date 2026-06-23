// Package app là tầng use-case của inventory. LÁT NÀY CHỈ ĐỌC nên app MỎNG: nó
// chuẩn hoá input (limit/cursor), gọi port domain.Repository, áp quy tắc FEFO
// thuần (domain.SortFEFO) cho phòng hờ, và trả kết quả. Không mở transaction (đọc
// thuần). Theo ARCHITECTURE.md §3 module nhẹ được phép gộp logic mỏng — nhưng vẫn
// giữ domain thuần + port.
//
// HOÃN (sẽ thêm sau module orders, cần khoá tranh chấp + tx tiền): trừ kho FEFO,
// nhập kho, chuyển kho, kiểm kê, phân bổ giá vốn. Khi đó app này lên "full" 3 lớp
// với TxManager — hiện chưa cần (chống over-engineer).
package app

import (
	"context"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

const (
	// defaultLimit là số bản ghi mỗi trang khi client không truyền.
	defaultLimit = 50
	// maxLimit chặn trên để tránh quét toàn bảng trong một request.
	maxLimit = 200
)

// Service là use-case đọc inventory mà edge (HTTP) dùng.
type Service struct {
	repo domain.Repository
}

// New dựng Service từ một repository (port domain).
func New(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// ListWarehouses trả các kho còn hoạt động (lọc status nếu truyền). limit chuẩn
// hoá ở đây.
func (s *Service) ListWarehouses(ctx context.Context, status string, limit int32) ([]domain.Warehouse, error) {
	return s.repo.ListWarehouses(ctx, domain.WarehouseFilter{
		Limit:  normalizeLimit(limit),
		Status: status,
	})
}

// ListStockQuery là input đã được edge giải mã (cursor thô, app tự decode để giữ
// quy ước pagination ở một nơi). ProductID/WarehouseID nil = không lọc tiêu chí đó.
type ListStockQuery struct {
	Cursor      string
	Limit       int32
	ProductID   *int64
	WarehouseID *int64
}

// ListStockResult là một trang tồn theo kho + cursor trang kế (rỗng nếu hết).
type ListStockResult struct {
	Items      []domain.StockItem
	NextCursor string
}

// ListStock trả một trang tồn TỔNG theo kho (keyset). Tự decode cursor, chuẩn hoá
// limit, và sinh NextCursor = id của phần tử cuối nếu trang đầy (có thể còn nữa).
func (s *Service) ListStock(ctx context.Context, q ListStockQuery) (ListStockResult, error) {
	afterID, err := pagination.DecodeID(q.Cursor)
	if err != nil {
		return ListStockResult{}, apperr.Validation("cursor không hợp lệ")
	}
	limit := normalizeLimit(q.Limit)

	items, err := s.repo.ListStock(ctx, domain.StockFilter{
		AfterID:     afterID,
		Limit:       limit,
		ProductID:   q.ProductID,
		WarehouseID: q.WarehouseID,
	})
	if err != nil {
		return ListStockResult{}, err
	}

	res := ListStockResult{Items: items}
	if int32(len(items)) == limit && limit > 0 {
		res.NextCursor = pagination.EncodeID(items[len(items)-1].ID)
	}
	return res, nil
}

// ListBatchesFEFO trả các lô CÒN TỒN của một product, sắp xếp FEFO (hạn dùng tăng
// dần). Repo đã ORDER BY expiry ASC; app gọi domain.SortFEFO để áp quy tắc FEFO
// như nguồn chân lý ở domain (phòng khi nguồn dữ liệu đổi thứ tự) và để chuẩn bị
// tái dùng cho xuất kho FEFO sau này. Optional lọc theo kho.
func (s *Service) ListBatchesFEFO(ctx context.Context, productID int64, warehouseID *int64) ([]domain.Batch, error) {
	batches, err := s.repo.ListBatchesFEFO(ctx, productID, warehouseID)
	if err != nil {
		return nil, err
	}
	domain.SortFEFO(batches)
	return batches, nil
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
