package domain

import "context"

// StockFilter gom điều kiện lọc + keyset pagination cho danh sách tồn theo kho.
// Tất cả tiêu chí lọc là optional (nil = không lọc theo tiêu chí đó). Limit/AfterID
// do tầng app chuẩn hoá. Repo ngầm loại tồn của kho đã đóng (warehouses.deleted_at).
type StockFilter struct {
	// AfterID: chỉ lấy dòng tồn có id > AfterID (keyset theo product_inventory.id
	// ASC). 0 = trang đầu.
	AfterID int64
	// Limit: số bản ghi tối đa mỗi trang (đã chuẩn hoá ở app, > 0).
	Limit int32
	// ProductID lọc theo sản phẩm; nil = mọi sản phẩm.
	ProductID *int64
	// WarehouseID lọc theo kho; nil = mọi kho.
	WarehouseID *int64
}

// WarehouseFilter gom điều kiện lọc cho danh sách kho. Số kho nhỏ nên không
// keyset; vẫn có Limit chặn trên ở app để an toàn.
type WarehouseFilter struct {
	// Limit: số kho tối đa trả về (đã chuẩn hoá ở app, > 0).
	Limit int32
	// Status lọc theo trạng thái (vd "active"); rỗng = mọi trạng thái.
	Status string
}

// Repository là PORT đọc dữ liệu inventory do domain ĐỊNH NGHĨA; adapter postgres
// implement ("accept interfaces, return structs"). LÁT NÀY CHỈ ĐỌC nên port chỉ
// có thao tác đọc — KHÔNG Insert/Update/Delete (trừ/nhập/chuyển kho HOÃN sau
// orders). Các thao tác là độc lập, không cần transaction (đọc thuần) nên app gọi
// thẳng repo bind pool.
type Repository interface {
	// ListWarehouses trả các kho còn hoạt động (lọc deleted_at IS NULL + status
	// nếu truyền), tối đa f.Limit, theo id ASC. Hết → slice rỗng (không lỗi).
	ListWarehouses(ctx context.Context, f WarehouseFilter) ([]Warehouse, error)
	// ListStock trả tối đa f.Limit dòng tồn theo keyset (id > f.AfterID, id ASC),
	// đã loại tồn của kho đã đóng + lọc product_id/warehouse_id nếu truyền. Hết →
	// slice rỗng.
	ListStock(ctx context.Context, f StockFilter) ([]StockItem, error)
	// ListBatchesFEFO trả các lô CÒN TỒN (> 0, lô chưa soft-delete) của một
	// product, đã sắp xếp FEFO (expiry_date ASC). Optional lọc theo warehouseID
	// (nil = mọi kho). Không có → slice rỗng (không lỗi).
	ListBatchesFEFO(ctx context.Context, productID int64, warehouseID *int64) ([]Batch, error)
}
