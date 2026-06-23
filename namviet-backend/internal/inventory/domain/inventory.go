// Package domain là LÕI THUẦN của bounded context inventory: entity kho hàng,
// tồn theo kho, lô hàng (FEFO) + value object Quantity/Money + PORT interface.
// KHÔNG import pgx/http/huma/framework (ARCHITECTURE.md §3). Chỉ stdlib + shared
// kernel trung lập (common/money) + shopspring/decimal (qua Quantity). Phụ thuộc
// đi một chiều: adapters → app → domain.
//
// LÁT NÀY CHỈ ĐỌC (ADR 0001, strangler-fig): liệt kê kho, tồn theo product/kho,
// lô theo FEFO. Mọi đường GHI (trừ kho FEFO, nhập, chuyển kho, kiểm kê, phân bổ
// giá vốn) ĐƯỢC HOÃN có chủ đích — làm SAU module orders vì cần khoá tranh chấp
// (concurrency) + transaction tiền. Vì vậy domain inventory hiện chỉ có entity
// đọc + rule sắp xếp FEFO thuần, CHƯA có invariant ghi/trừ tồn.
package domain

import (
	"errors"
	"sort"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// ErrInsufficientStock là lỗi DOMAIN khi tổng tồn KHẢ DỤNG của các lô không đủ
// cho nhu cầu trừ kho. Service map sang apperr.Conflict (KHÔNG cho tồn âm). Là
// sentinel để caller errors.Is — domain thuần không biết apperr.
var ErrInsufficientStock = errors.New("inventory: tồn không đủ để trừ")

// ConsumedBatch là MỘT lô bị tiêu thụ trong kế hoạch trừ kho FEFO: định danh lô
// (BatchID + dòng tồn-theo-lô InventoryBatchID để adapter UPDATE đúng dòng),
// lượng trừ ở lô đó (Quantity, decimal — có thể là MỘT PHẦN tồn lô ở lô cuối), và
// giá vốn nhập của lô (InboundPrice) — mang theo để orders post bút toán COGS sau
// (Σ inbound_price × qty). KHÔNG float.
type ConsumedBatch struct {
	BatchID          int64
	InventoryBatchID int64
	Quantity         Quantity
	InboundPrice     money.Money
}

// Warehouse là một kho / chi nhánh / cửa hàng (public.warehouses). ID là int64 vì
// bảng dùng khoá bigint (KHÔNG uuid). CompanyID là uuid dạng chuỗi (nullable →
// rỗng nếu chưa gán công ty con). Bảng KHÔNG có updated_at.
type Warehouse struct {
	ID         int64
	Key        string
	Name       string
	Unit       string // đơn vị quản lý tồn gốc (vd "Hộp")
	Address    string
	Type       string // 'retail' | 'wholesale' | 'central'...
	Code       string
	Manager    string
	Phone      string
	Status     string // 'active' | 'inactive'...
	CompanyID  string // uuid dạng chuỗi; rỗng = chưa gán
	OutletType string
	CreatedAt  time.Time
}

// StockItem là tồn TỔNG của một product tại một kho (public.product_inventory).
// Quantity là lượng tồn (NUMERIC → Quantity, KHÔNG float). MinStock/MaxStock là
// ngưỡng cảnh báo (int, nullable → 0 khi không khai).
type StockItem struct {
	ID            int64
	ProductID     int64
	WarehouseID   int64
	Quantity      Quantity
	MinStock      int32
	MaxStock      int32
	ShelfLocation string
	UpdatedAt     time.Time
}

// Batch là một LÔ còn tồn của một product tại một kho, đã ghép hạn dùng + giá
// nhập từ public.batches. Đây là dữ liệu nền cho xuất kho FEFO sau này.
//   - InventoryBatchID: khoá dòng tồn-theo-lô (public.inventory_batches.id).
//   - BatchID: khoá lô (public.batches.id).
//   - Quantity: tồn của lô tại kho (NUMERIC → Quantity, KHÔNG float).
//   - ExpiryDate: hạn dùng (DATE) — khoá sắp xếp FEFO.
//   - InboundPrice: giá vốn nhập của lô (NUMERIC → money.Money, KHÔNG float).
type Batch struct {
	InventoryBatchID  int64
	WarehouseID       int64
	ProductID         int64
	BatchID           int64
	Quantity          Quantity
	BatchCode         string
	ExpiryDate        time.Time
	ManufacturingDate time.Time // zero time nếu NULL
	HasManufacturing  bool      // false nếu manufacturing_date NULL
	InboundPrice      money.Money
}

// SortFEFO sắp xếp danh sách lô theo FEFO (First-Expired-First-Out): hạn dùng
// (ExpiryDate) TĂNG DẦN — lô hết hạn TRƯỚC đứng TRƯỚC (ưu tiên xuất). Tie-break
// theo BatchID rồi InventoryBatchID để thứ tự ổn định (deterministic), khớp ORDER
// BY của query. Hàm THUẦN, ổn định (sort.SliceStable), sửa slice tại chỗ — để
// unit test rule sắp xếp KHÔNG cần DB. Repo đã ORDER BY sẵn; hàm này là nguồn
// chân lý của quy tắc FEFO ở domain (tái dùng khi xuất kho sau).
func SortFEFO(batches []Batch) {
	sort.SliceStable(batches, func(i, j int) bool {
		a, b := batches[i], batches[j]
		if !a.ExpiryDate.Equal(b.ExpiryDate) {
			return a.ExpiryDate.Before(b.ExpiryDate)
		}
		if a.BatchID != b.BatchID {
			return a.BatchID < b.BatchID
		}
		return a.InventoryBatchID < b.InventoryBatchID
	})
}

// PlanFEFO lập KẾ HOẠCH tiêu thụ THUẦN cho việc trừ kho theo FEFO: cộng dồn tồn
// các lô trong available THEO THỨ TỰ slice (caller PHẢI đã sắp FEFO — repo ORDER BY
// expiry ASC, hoặc gọi SortFEFO trước) tới khi đủ need. Lô đầy đủ tiêu thụ TRỌN
// VẸN; lô CUỐI có thể tiêu thụ MỘT PHẦN (tồn lô > nhu cầu còn lại → chỉ lấy phần
// còn thiếu). Trả danh sách ConsumedBatch (kèm InboundPrice mỗi lô cho COGS).
//
// Hàm THUẦN (chỉ stdlib + Quantity/money decimal), KHÔNG chạm DB — để unit-test
// quy tắc tiêu thụ không cần Postgres (giống SortFEFO). available coi như đã lọc
// (lô deleted/tồn 0 đã bị query loại); PlanFEFO chỉ làm số học tồn.
//
// need <= 0 → kế hoạch RỖNG, không lỗi (no-op hợp lệ). Tổng tồn < need →
// ErrInsufficientStock + plan nil (KHÔNG cho tồn âm; KHÔNG tiêu thụ một phần rồi
// báo thiếu — all-or-nothing, để caller rollback sạch).
func PlanFEFO(available []Batch, need Quantity) ([]ConsumedBatch, error) {
	if !need.IsPositive() {
		return nil, nil
	}

	remaining := need
	plan := make([]ConsumedBatch, 0, len(available))
	for _, b := range available {
		if !remaining.IsPositive() {
			break
		}
		if !b.Quantity.IsPositive() {
			continue // phòng thủ: lô không còn tồn (query đã loại, nhưng giữ an toàn)
		}
		// Lấy min(tồn lô, còn cần): lô cuối tiêu thụ một phần khi tồn lô > remaining.
		take := b.Quantity
		if b.Quantity.Cmp(remaining) > 0 {
			take = remaining
		}
		plan = append(plan, ConsumedBatch{
			BatchID:          b.BatchID,
			InventoryBatchID: b.InventoryBatchID,
			Quantity:         take,
			InboundPrice:     b.InboundPrice,
		})
		remaining = remaining.Sub(take)
	}

	if remaining.IsPositive() {
		// Quét hết lô mà vẫn còn cần → tổng tồn không đủ. KHÔNG trả plan dở dang.
		return nil, ErrInsufficientStock
	}
	return plan, nil
}
