package domain

import (
	"errors"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// OrderType là loại đơn — khớp CHECK của public.orders (order_type ∈ B2B/B2C,
// default B2C). B2B = bán buôn/portal (xuất HĐ VAT 100%); B2C = bán lẻ/POS.
type OrderType string

const (
	// OrderTypeB2C — bán lẻ (mặc định).
	OrderTypeB2C OrderType = "B2C"
	// OrderTypeB2B — bán buôn/portal.
	OrderTypeB2B OrderType = "B2B"
)

// String trả giá trị chuỗi (để adapter ghi cột order_type text).
func (t OrderType) String() string { return string(t) }

// Valid trả true nếu t là loại đơn hợp lệ.
func (t OrderType) Valid() bool { return t == OrderTypeB2C || t == OrderTypeB2B }

// Lỗi domain THUẦN khi dựng draft (app map sang apperr.Validation → 422).
var (
	ErrNoLines           = errors.New("đơn phải có ít nhất một dòng hàng")
	ErrInvalidOrderType  = errors.New("order_type không hợp lệ (chỉ B2B/B2C)")
	ErrQuantityNotPos    = errors.New("số lượng mỗi dòng phải lớn hơn 0")
	ErrUnitPriceNegative = errors.New("đơn giá không được âm")
	ErrDiscountNegative  = errors.New("chiết khấu không được âm")
	ErrDiscountTooLarge  = errors.New("chiết khấu dòng không được lớn hơn thành tiền (quantity*unit_price)")
	ErrInvalidProductID  = errors.New("product_id không hợp lệ (phải > 0)")
	ErrEmptyUOM          = errors.New("đơn vị tính (uom) không được rỗng")
)

// DraftLine là input MỘT dòng hàng để dựng đơn (tiền decimal — KHÔNG float).
type DraftLine struct {
	ProductID int64
	Quantity  Quantity
	UOM       string
	UnitPrice money.Money
	Discount  money.Money
}

// DraftInput là input dựng một đơn mới (chưa có mã/ID — sinh ở app/adapter).
// CustomerID nullable (POS B2C vãng lai → nil). OrderType rỗng → mặc định B2C.
type DraftInput struct {
	CustomerID *int64
	OrderType  OrderType
	CreatorID  string
	Note       string
	Lines      []DraftLine
}

// ComputedLine là một dòng hàng ĐÃ TÍNH thành tiền (line_total). Lưu lại Quantity/
// UnitPrice/Discount để adapter ghi nguyên vào order_items.
type ComputedLine struct {
	ProductID int64
	Quantity  Quantity
	UOM       string
	UnitPrice money.Money
	Discount  money.Money
	LineTotal money.Money // = quantity*unit_price - discount
}

// Draft là một đơn HỢP LỆ đã tính tiền, sẵn sàng persist (chưa có code/ID). Status
// luôn PENDING. TotalAmount = Σ quantity*unit_price (TRƯỚC chiết khấu); FinalAmount
// = Σ line_total (SAU chiết khấu).
type Draft struct {
	CustomerID  *int64
	OrderType   OrderType
	CreatorID   string
	Note        string
	Status      Status
	TotalAmount money.Money
	FinalAmount money.Money
	Lines       []ComputedLine
}

// NewDraft validate input THUẦN và tính tiền cho một đơn mới. Quy ước tiền (xem
// spec §8 + giả định ghi rõ ở module.go):
//
//	line_total    = quantity * unit_price - discount   (mỗi dòng)
//	total_amount  = Σ (quantity * unit_price)          (TỔNG TRƯỚC chiết khấu)
//	final_amount  = Σ line_total                        (TỔNG SAU chiết khấu)
//
// CHƯA gồm VAT (xem cờ cảnh báo ở module.go — P4b post sổ TAX dựa giả định này).
// Lỗi là lỗi domain thuần (app map Validation/422).
func NewDraft(in DraftInput) (Draft, error) {
	ot := in.OrderType
	if ot == "" {
		ot = OrderTypeB2C
	}
	if !ot.Valid() {
		return Draft{}, ErrInvalidOrderType
	}
	if len(in.Lines) == 0 {
		return Draft{}, ErrNoLines
	}

	total := money.Zero() // Σ quantity*unit_price (trước chiết khấu)
	final := money.Zero() // Σ line_total (sau chiết khấu)
	computed := make([]ComputedLine, 0, len(in.Lines))
	for _, l := range in.Lines {
		if l.ProductID <= 0 {
			return Draft{}, ErrInvalidProductID
		}
		if l.UOM == "" {
			return Draft{}, ErrEmptyUOM
		}
		if l.Quantity.IsZero() || l.Quantity.Decimal().IsNegative() {
			return Draft{}, ErrQuantityNotPos
		}
		if l.UnitPrice.IsNegative() {
			return Draft{}, ErrUnitPriceNegative
		}
		if l.Discount.IsNegative() {
			return Draft{}, ErrDiscountNegative
		}
		// gross = quantity * unit_price (decimal nhân decimal — KHÔNG float).
		gross := money.FromDecimal(l.Quantity.Decimal().Mul(l.UnitPrice.Decimal()))
		if l.Discount.Sub(gross).IsPositive() {
			// discount > gross → âm thành tiền (vô lý nghiệp vụ).
			return Draft{}, ErrDiscountTooLarge
		}
		lineTotal := gross.Sub(l.Discount)
		computed = append(computed, ComputedLine{
			ProductID: l.ProductID,
			Quantity:  l.Quantity,
			UOM:       l.UOM,
			UnitPrice: l.UnitPrice,
			Discount:  l.Discount,
			LineTotal: lineTotal,
		})
		total = total.Add(gross)
		final = final.Add(lineTotal)
	}

	return Draft{
		CustomerID:  in.CustomerID,
		OrderType:   ot,
		CreatorID:   in.CreatorID,
		Note:        in.Note,
		Status:      StatusPending,
		TotalAmount: total,
		FinalAmount: final,
		Lines:       computed,
	}, nil
}
