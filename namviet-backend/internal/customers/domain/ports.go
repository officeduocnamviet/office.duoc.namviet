package domain

import "context"

// CustomerFilter gom điều kiện lọc + keyset pagination cho danh sách khách hàng.
// Tất cả tuỳ chọn (rỗng/nil = không lọc theo tiêu chí đó) trừ Limit/AfterID do
// tầng app chuẩn hoá. Repo luôn ngầm lọc deleted_at IS NULL (soft-delete).
type CustomerFilter struct {
	// AfterID: chỉ lấy customer có id > AfterID (keyset theo id ASC). 0 = trang đầu.
	AfterID int64
	// Limit: số bản ghi tối đa mỗi trang (đã chuẩn hoá ở app, > 0).
	Limit int32
	// Type lọc theo loại khách (TypeB2B/TypeB2C); rỗng = mọi loại.
	Type CustomerType
	// Status lọc theo trạng thái (vd "active"); rỗng = mọi trạng thái.
	Status string
	// Query tìm theo name/phone/customer_code/MST (ILIKE); rỗng = không tìm.
	Query string
}

// Repository là PORT đọc dữ liệu customers do domain ĐỊNH NGHĨA; adapter postgres
// implement ("accept interfaces, return structs"). Customers read-mostly nên port
// chỉ có thao tác đọc — không Insert/Update/Delete (DEFER tạo/sửa sang sau). Các
// thao tác đọc độc lập, không cần transaction nên app gọi thẳng repo bind pool.
type Repository interface {
	// ListCustomers trả tối đa f.Limit khách theo keyset (id > f.AfterID, id ASC),
	// đã lọc deleted_at IS NULL + các tiêu chí trong f, kèm DebtSnapshot (đã chọn
	// nguồn live). Hết dữ liệu → slice rỗng.
	ListCustomers(ctx context.Context, f CustomerFilter) ([]Customer, error)
	// GetCustomerByID trả khách theo id (chưa bị soft-delete) kèm B2BProfile (nếu
	// có) + DebtSnapshot. Không thấy → apperr.NotFound.
	GetCustomerByID(ctx context.Context, id int64) (Customer, error)
}
