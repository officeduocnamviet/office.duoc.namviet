package domain

import "context"

// ProductFilter gom các điều kiện lọc + keyset pagination cho danh sách sản
// phẩm. Tất cả tùy chọn (nil/rỗng = không lọc theo tiêu chí đó) trừ Limit/AfterID
// do tầng app chuẩn hoá. Repo luôn ngầm lọc deleted_at IS NULL (soft-delete).
type ProductFilter struct {
	// AfterID: chỉ lấy product có id > AfterID (keyset theo id ASC). 0 = trang đầu.
	AfterID int64
	// Limit: số bản ghi tối đa mỗi trang (đã chuẩn hoá ở app, > 0).
	Limit int32
	// Status lọc theo trạng thái kinh doanh (vd "active"); rỗng = mọi trạng thái.
	Status string
	// CategoryID lọc theo nhóm ngành hàng; nil = mọi nhóm.
	CategoryID *int64
	// Query tìm theo name/sku (ILIKE); rỗng = không tìm.
	Query string
}

// Repository là PORT đọc dữ liệu catalog do domain ĐỊNH NGHĨA; adapter postgres
// implement ("accept interfaces, return structs"). Catalog read-only nên port
// chỉ có thao tác đọc — không Insert/Update/Delete. Các thao tác là độc lập,
// không cần transaction (đọc thuần) nên app gọi thẳng repo bind pool.
type Repository interface {
	// ListProducts trả tối đa f.Limit product theo keyset (id > f.AfterID, id ASC),
	// đã lọc deleted_at IS NULL + các tiêu chí trong f. Hết dữ liệu → slice rỗng.
	ListProducts(ctx context.Context, f ProductFilter) ([]Product, error)
	// GetProductByID trả product theo id (chưa bị soft-delete). Không thấy →
	// apperr.NotFound.
	GetProductByID(ctx context.Context, id int64) (Product, error)
	// ListUnits trả các đơn vị tính của một product (base trước), đã lọc
	// deleted_at IS NULL. Không có → slice rỗng (không lỗi).
	ListUnits(ctx context.Context, productID int64) ([]ProductUnit, error)
	// ListCategories trả danh mục (lọc deleted_at IS NULL + status nếu truyền).
	ListCategories(ctx context.Context, status string) ([]Category, error)
	// ListManufacturers trả hãng (lọc deleted_at IS NULL + status nếu truyền).
	ListManufacturers(ctx context.Context, status string) ([]Manufacturer, error)
}
