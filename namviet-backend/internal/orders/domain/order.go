// Package domain là LÕI THUẦN của bounded context orders: aggregate Order (đơn
// bán) + dòng hàng (OrderLine) + ảnh chụp thanh toán (PaymentSummary) + PORT
// interface. KHÔNG import pgx/http/huma/framework (ARCHITECTURE.md §3). Chỉ
// stdlib + shared kernel trung lập (common/money) + shopspring/decimal (qua
// Quantity). Phụ thuộc đi một chiều: adapters → app → domain.
//
// LÁT NÀY CHỈ ĐỌC (ADR 0001, strangler-fig): liệt kê đơn (keyset + lọc), xem chi
// tiết đơn + dòng hàng, và SUY DIỄN "đã thu / còn nợ" từ finance_transactions.
// Mọi đường GHI (tạo đơn, trừ kho FEFO, ghi phiếu thu, post sổ kế toán) ĐƯỢC
// HOÃN có chủ đích sang slice sau (design đang chốt với user) — vì cần state
// machine 3 trạng thái, khoá tranh chấp tồn kho và transaction tiền/double-entry.
// Vì vậy domain orders hiện chỉ có entity đọc + quy tắc suy diễn thanh toán
// THUẦN, CHƯA có invariant chuyển trạng thái/ghi.
package domain

import (
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// Order là aggregate gốc của context orders (public.orders). ID là uuid dạng
// CHUỖI (bảng dùng uuid, KHÔNG bigint) — map ở adapter. CustomerID nullable
// (đơn lẻ POS có thể không gắn khách) → con trỏ; nil = không gắn khách. Tiền
// (Total/Final) dùng money.Money — KHÔNG float. Payment là ảnh chụp suy diễn
// read-only từ finance_transactions (xem PaymentSummary). Bảng KHÔNG có cột
// paid_amount nên KHÔNG bao giờ ghi ngược — chỉ suy diễn lúc đọc.
type Order struct {
	ID            string
	Code          string // mã hóa đơn (VD HD123); NOT NULL ở DB
	CustomerID    *int64 // bigint nullable; nil = không gắn khách
	CreatorID     string // uuid người lập (nullable → rỗng)
	Status        string // trạng thái xử lý đơn ('PENDING'|'COMPLETED'...)
	OrderType     string // 'B2C' | 'B2B'
	Total         money.Money
	Final         money.Money
	PaymentStatus string // 'unpaid' | 'partial' | 'paid'
	Note          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	// Payment là ảnh chụp đã-thu/còn-nợ suy diễn từ finance_transactions (sổ thực
	// tế INTERNAL). Zero value khi chưa nạp (vd ở danh sách nếu không tính kèm).
	Payment PaymentSummary
}

// OrderLine là một dòng hàng của đơn (public.order_items). ID uuid dạng chuỗi.
// Quantity ở DB này là INTEGER, nhưng ta map sang value object Quantity (decimal)
// để KHÔNG ép float và đồng nhất với inventory (ADR 0002 §B coi "thống nhất kiểu
// quantity" là việc migrate sau). UnitPrice/Discount/LineTotal là money — KHÔNG
// float. ExpiryDate là hạn dùng của lô đã chọn (nullable → zero time + HasExpiry
// = false).
type OrderLine struct {
	ID         string
	ProductID  int64
	Quantity   Quantity
	UOM        string // đơn vị tính được chọn (VD "Vỉ")
	UnitPrice  money.Money
	Discount   money.Money
	LineTotal  money.Money
	IsGift     bool
	BatchNo    string
	ExpiryDate time.Time
	HasExpiry  bool // false nếu order_items.expiry_date NULL
	Note       string
}

// PaymentSummary là ảnh chụp thanh toán của một đơn tại thời điểm đọc, SUY DIỄN
// read-only từ finance_transactions (KHÔNG có cột paid_amount ở orders). Paid =
// tổng phiếu THU (flow IN) đã hoàn tất trỏ về đơn theo (ref_type='order',
// ref_id=orders.code), tính theo SỔ THỰC TẾ (book_type INTERNAL/BOTH). Remaining
// = Final - Paid; có thể ÂM nếu thu thừa (over-paid) — KHÔNG kẹp về 0 để khỏi
// che lỗi nhập trùng phiếu thu.
type PaymentSummary struct {
	Final     money.Money // tổng tiền khách phải trả (từ đơn)
	Paid      money.Money // đã thu (suy diễn từ finance_transactions)
	Remaining money.Money // còn nợ = Final - Paid (có thể âm)
}

// ComputePayment dựng PaymentSummary từ tổng phải-trả (final) và đã-thu (paid).
// Hàm THUẦN (dễ unit test, không DB): Remaining = Final - Paid. Là nguồn chân
// lý của quy tắc "còn nợ" ở domain (tái dùng khi ghi phiếu thu sau này). KHÔNG
// kẹp âm: thu thừa là sự thật cần thấy, không phải lỗi để giấu.
func ComputePayment(final, paid money.Money) PaymentSummary {
	return PaymentSummary{
		Final:     final,
		Paid:      paid,
		Remaining: final.Sub(paid),
	}
}
