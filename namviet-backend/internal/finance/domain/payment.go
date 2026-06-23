// Package domain là LÕI THUẦN của bounded context finance: phiếu thu/chi
// (Payment) trên sổ quỹ. P3 chỉ làm đường THU (PaymentIn). KHÔNG import
// pgx/http/huma/framework (ARCHITECTURE.md §3) — chỉ stdlib + shared kernel trung
// lập (common/money). Phụ thuộc một chiều: adapters → app → domain.
//
// BẤT BIẾN NGHIỆP VỤ (đúng đắn tiền > mọi thứ — ARCHITECTURE.md §7):
//   - amount phải > 0 (không ghi phiếu 0/âm — soft-delete/đảo phiếu là việc khác).
//   - book_type ∈ {INTERNAL, TAX, BOTH} (khớp giá trị thật finance_transactions).
//   - order_code không rỗng (ref_id = orders.code, liên kết đơn↔phiếu thu qua mã).
//   - Tiền = common/money (decimal) ↔ NUMERIC. CẤM float.
package domain

import (
	"errors"
	"strings"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// BookType là sổ mà một phiếu thu ghi vào. Khác Book của accounting (chỉ
// INTERNAL/TAX): finance_transactions cho phép BOTH = ghi CẢ hai sổ (tiền thật +
// có hoá đơn VAT). Bút toán post sau (P4): INTERNAL/BOTH → sổ INTERNAL;
// TAX/BOTH → sổ TAX.
type BookType string

const (
	// BookInternal — phiếu chỉ vào sổ thực tế (thu thật, không xuất HĐ).
	BookInternal BookType = "INTERNAL"
	// BookTax — phiếu chỉ vào sổ thuế (hiếm; thường đi cùng INTERNAL nên dùng BOTH).
	BookTax BookType = "TAX"
	// BookBoth — phiếu vào CẢ hai sổ: mặc định cho đơn B2B có HĐ VAT (tiền thật + HĐ).
	BookBoth BookType = "BOTH"
)

// Valid trả true nếu bt là một book_type hợp lệ (khớp giá trị thật ở DB cũ).
func (bt BookType) Valid() bool {
	return bt == BookInternal || bt == BookTax || bt == BookBoth
}

// String trả giá trị chuỗi (để adapter ghi cột book_type text).
func (bt BookType) String() string { return string(bt) }

// Các hằng quy ước nghiệp vụ khi ghi phiếu THU cho một đơn (xác nhận từ core cũ +
// ERP migration): flow='in', status='completed' (tiền đã thực vào quỹ — mốc
// trigger prod cộng số dư), ref_type='order' (trỏ về đơn bán qua orders.code).
const (
	FlowIn  = "in"
	FlowOut = "out" // chi tiền (phiếu CHI — trả NCC mua hàng, mục 54)
	// RefTypePurchaseOrder — phiếu CHI gắn đơn MUA (purchase_orders.code). Đối xứng
	// RefTypeOrder (chiều bán); query "đã chi theo PO" lọc ref_type='purchase_order'.
	RefTypePurchaseOrder = "purchase_order"
	// BusinessTypePurchase — nghiệp vụ: chi tiền mua hàng (phiếu CHI cho NCC).
	BusinessTypePurchase = "purchase_payment"
	StatusCompleted      = "completed"
	// RefTypeCustomer — phiếu THU lump-sum gắn KHÁCH (không phải 1 đơn): tiền phân
	// bổ cho nhiều đơn qua app.finance_transaction_allocations. ref_type='customer'
	// để query "đã thu trực tiếp theo đơn" (ref_type='order') KHÔNG đếm trùng.
	RefTypeCustomer = "customer"
	// StatusPending — phiếu thu ĐÃ THU TỪ KHÁCH nhưng CHƯA vào quỹ (vd NV giao hàng
	// thu tiền mặt, chờ thủ quỹ "Xác nhận đã thu"). Nợ khách GIẢM ngay (đã-thu đếm cả
	// pending), nhưng số dư quỹ CHƯA tăng (trigger prod chỉ bắn khi 'completed').
	// Spec system_features.md mục 55.
	StatusPending    = "pending"
	RefTypeOrder     = "order"
	BusinessTypeSale = "sale_receipt" // nghiệp vụ: thu tiền bán hàng
)

// Lỗi domain THUẦN (app map sang apperr.Validation → 422).
var (
	ErrAmountNotPositive    = errors.New("số tiền phiếu thu phải lớn hơn 0")
	ErrInvalidBookType      = errors.New("book_type không hợp lệ (chỉ INTERNAL/TAX/BOTH)")
	ErrEmptyOrderCode       = errors.New("mã đơn (order_code) không được rỗng")
	ErrInvalidInitialStatus = errors.New("trạng thái khởi tạo không hợp lệ (chỉ pending/completed)")
)

// RecordPaymentIn là tham số THUẦN để ghi một phiếu THU cho đơn. KHÔNG chứa kiểu
// hạ tầng (chỉ money + chuỗi/số). app dựng từ DTO/port input rồi gọi Validate.
type RecordPaymentIn struct {
	// OrderCode = orders.code (text) — set vào ref_id. Liên kết đơn↔phiếu thu.
	OrderCode string
	// Amount số tiền thu (> 0). money decimal — KHÔNG float.
	Amount money.Money
	// FundAccountID quỹ/tài khoản nhận tiền (public.fund_accounts.id, bigint).
	FundAccountID int64
	// BookType sổ ghi nhận (INTERNAL/TAX/BOTH). Rỗng → app mặc định BOTH (đơn B2B
	// có HĐ VAT) — nhưng Validate yêu cầu app đã chuẩn hoá trước.
	BookType BookType
	// BankRef (tuỳ chọn) mã giao dịch ngân hàng cho phiếu THU TỰ ĐỘNG (webhook) —
	// dùng chống trùng. nil cho phiếu thủ công.
	BankRef *string
	// Description (tuỳ chọn) diễn giải phiếu.
	Description *string
	// CreatedBy (tuỳ chọn) uuid người lập phiếu (text uuid).
	CreatedBy *string
	// InitialStatus (tuỳ chọn) trạng thái khởi tạo phiếu: rỗng/'completed' = thu
	// thẳng vào quỹ ngay (POS, webhook ngân hàng đã về); 'pending' = đã thu từ khách
	// nhưng chưa vào quỹ (NV giao hàng thu tiền mặt — chờ thủ quỹ xác nhận). Thanh
	// toán 2 bước (spec mục 55): pending → completed bằng ConfirmReceipt.
	InitialStatus string
	// RefType/RefID (tuỳ chọn) ghi đè chứng từ gốc của phiếu. Rỗng → mặc định
	// ('order', OrderCode) cho phiếu thu MỘT đơn (giữ hành vi cũ). Đặt
	// ('customer', <customerID>) cho phiếu LUMP-SUM phân bổ nhiều đơn (mục 55).
	RefType string
	RefID   string
}

// EffectiveRefType trả ref_type thực ghi: RefType nếu khai, ngược lại 'order'.
func (p RecordPaymentIn) EffectiveRefType() string {
	if strings.TrimSpace(p.RefType) != "" {
		return p.RefType
	}
	return RefTypeOrder
}

// EffectiveRefID trả ref_id thực ghi: RefID nếu khai RefType, ngược lại OrderCode.
func (p RecordPaymentIn) EffectiveRefID() string {
	if strings.TrimSpace(p.RefType) != "" {
		return strings.TrimSpace(p.RefID)
	}
	return strings.TrimSpace(p.OrderCode)
}

// EffectiveStatus trả trạng thái khởi tạo thực dùng khi INSERT: 'pending' nếu khai
// rõ, ngược lại mặc định 'completed' (giữ hành vi cũ — thu ngay vào quỹ).
func (p RecordPaymentIn) EffectiveStatus() string {
	if p.InitialStatus == StatusPending {
		return StatusPending
	}
	return StatusCompleted
}

// Validate ép bất biến THUẦN trước khi ghi (app gọi → map Validation/422). KHÔNG
// chạm DB. Chuẩn hoá OrderCode (trim) để so khớp ref_id ổn định.
func (p RecordPaymentIn) Validate() error {
	// Chứng từ gốc bắt buộc: phiếu đơn → OrderCode; phiếu lump-sum → RefID (customer).
	if strings.TrimSpace(p.EffectiveRefID()) == "" {
		return ErrEmptyOrderCode
	}
	if !p.Amount.IsPositive() {
		return ErrAmountNotPositive
	}
	if !p.BookType.Valid() {
		return ErrInvalidBookType
	}
	if p.InitialStatus != "" && p.InitialStatus != StatusPending && p.InitialStatus != StatusCompleted {
		return ErrInvalidInitialStatus
	}
	return nil
}

// RecordPaymentOut là tham số THUẦN để ghi một phiếu CHI (trả NCC mua hàng — mục
// 54). Đối xứng RecordPaymentIn (chiều thu). flow='out'. Chứng từ gốc mặc định
// ('purchase_order', POCode). KHÔNG chứa kiểu hạ tầng (chỉ money + chuỗi/số).
type RecordPaymentOut struct {
	// POCode = purchase_orders.code (text) — set vào ref_id. Liên kết PO↔phiếu chi.
	POCode string
	// Amount số tiền chi (> 0). money decimal — KHÔNG float.
	Amount money.Money
	// FundAccountID quỹ/tài khoản XUẤT tiền (public.fund_accounts.id, bigint).
	FundAccountID int64
	// BookType sổ ghi nhận (INTERNAL/TAX/BOTH). Rỗng → app mặc định.
	BookType BookType
	// BankRef (tuỳ chọn) mã giao dịch ngân hàng (phiếu chi tự động) — dedup chống trùng.
	BankRef *string
	// Description (tuỳ chọn) diễn giải phiếu.
	Description *string
	// CreatedBy (tuỳ chọn) uuid người lập phiếu.
	CreatedBy *string
	// InitialStatus (tuỳ chọn) trạng thái khởi tạo: rỗng/'completed' = đã chi khỏi quỹ
	// (số dư giảm — trigger prod); 'pending' = lệnh chi chờ duyệt (chưa rút quỹ).
	InitialStatus string
	// RefType/RefID (tuỳ chọn) ghi đè chứng từ gốc. Rỗng → ('purchase_order', POCode).
	RefType string
	RefID   string
}

// EffectiveRefType trả ref_type thực ghi: RefType nếu khai, ngược lại 'purchase_order'.
func (p RecordPaymentOut) EffectiveRefType() string {
	if strings.TrimSpace(p.RefType) != "" {
		return p.RefType
	}
	return RefTypePurchaseOrder
}

// EffectiveRefID trả ref_id thực ghi: RefID nếu khai RefType, ngược lại POCode.
func (p RecordPaymentOut) EffectiveRefID() string {
	if strings.TrimSpace(p.RefType) != "" {
		return strings.TrimSpace(p.RefID)
	}
	return strings.TrimSpace(p.POCode)
}

// EffectiveStatus trả trạng thái khởi tạo thực dùng khi INSERT: 'pending' nếu khai
// rõ, ngược lại mặc định 'completed' (đã chi khỏi quỹ).
func (p RecordPaymentOut) EffectiveStatus() string {
	if p.InitialStatus == StatusPending {
		return StatusPending
	}
	return StatusCompleted
}

// Validate ép bất biến THUẦN trước khi ghi phiếu CHI (app gọi → map Validation/422).
func (p RecordPaymentOut) Validate() error {
	if strings.TrimSpace(p.EffectiveRefID()) == "" {
		return ErrEmptyOrderCode
	}
	if !p.Amount.IsPositive() {
		return ErrAmountNotPositive
	}
	if !p.BookType.Valid() {
		return ErrInvalidBookType
	}
	if p.InitialStatus != "" && p.InitialStatus != StatusPending && p.InitialStatus != StatusCompleted {
		return ErrInvalidInitialStatus
	}
	return nil
}

// Payment là một phiếu thu/chi ĐÃ GHI (kết quả RecordPaymentIn). Phản chiếu một
// dòng public.finance_transactions ở mức domain (tiền = money, KHÔNG float).
type Payment struct {
	ID            int64
	Code          string
	Flow          string
	BusinessType  string
	Amount        money.Money
	FundAccountID int64
	RefType       string
	RefID         string
	Status        string
	BookType      BookType
	BankRef       *string
	Description   *string
	CreatedBy     *string
}
