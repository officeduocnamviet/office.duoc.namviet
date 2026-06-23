package domain

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// Lỗi domain THUẦN khi dựng draft PO (app map sang apperr.Validation → 422).
var (
	ErrNoLines           = errors.New("đơn mua phải có ít nhất một dòng hàng")
	ErrQuantityNotPos    = errors.New("số lượng mỗi dòng phải lớn hơn 0")
	ErrUnitCostNegative  = errors.New("đơn giá nhập không được âm")
	ErrVATRateNegative   = errors.New("thuế suất VAT không được âm")
	ErrInvalidProductID  = errors.New("product_id không hợp lệ (phải > 0)")
)

// DraftLine là input MỘT dòng hàng để dựng PO (tiền + lượng decimal — KHÔNG float).
// UnitCost là giá nhập per-unit (→ inbound_price lô khi nhập kho). VATRate decimal
// (vd 0.08). BatchCode/Expiry/Mfg cho lô nhập (optional ở draft, set khi nhận hàng).
type DraftLine struct {
	ProductID         int64
	Quantity          decimal.Decimal
	UnitCost          money.Money
	VATRate           decimal.Decimal
	BatchCode         string
	ExpiryDate        *time.Time
	ManufacturingDate *time.Time
}

// DraftInput là input dựng một PO mới (chưa có mã/ID — sinh ở app/adapter).
// SupplierID nullable (NCC chưa có entity — supplier_id bigint không FK).
type DraftInput struct {
	SupplierID   *int64
	SupplierName string
	Note         string
	Lines        []DraftLine
}

// ComputedLine là một dòng PO ĐÃ TÍNH thành tiền + VAT dòng. LineTotal = qty*unit_cost
// (ex-VAT, VND scale-0). VATAmount = LineTotal × VATRate (RoundVND mỗi dòng).
type ComputedLine struct {
	LineNo            int
	ProductID         int64
	Quantity          decimal.Decimal
	UnitCost          money.Money
	VATRate           decimal.Decimal
	BatchCode         string
	ExpiryDate        *time.Time
	ManufacturingDate *time.Time
	LineTotal         money.Money // = qty * unit_cost (ex-VAT)
	VATAmount         money.Money // = LineTotal × VATRate (RoundVND)
}

// Draft là một PO HỢP LỆ đã tính tiền, sẵn sàng persist (chưa có code/ID). Status
// luôn draft. TotalAmount = Σ LineTotal (ex-VAT); VATAmount = Σ VATAmount dòng.
type Draft struct {
	SupplierID   *int64
	SupplierName string
	Note         string
	Status       Status
	TotalAmount  money.Money
	VATAmount    money.Money
	Lines        []ComputedLine
}

// NewDraft validate input THUẦN và tính tiền cho một PO mới. Quy ước:
//
//	line_total   = quantity * unit_cost              (mỗi dòng, ex-VAT)
//	vat_dòng     = round(line_total × vat_rate)      (làm tròn VND mỗi dòng)
//	total_amount = Σ line_total                       (ex-VAT)
//	vat_amount   = Σ vat_dòng
//
// Tiền decimal toàn tuyến (CẤM float). Lỗi là lỗi domain thuần (app map Validation/422).
func NewDraft(in DraftInput) (Draft, error) {
	if len(in.Lines) == 0 {
		return Draft{}, ErrNoLines
	}

	total := money.Zero() // Σ line_total (ex-VAT)
	vatTotal := money.Zero()
	computed := make([]ComputedLine, 0, len(in.Lines))
	for i, l := range in.Lines {
		if l.ProductID <= 0 {
			return Draft{}, ErrInvalidProductID
		}
		if !l.Quantity.IsPositive() {
			return Draft{}, ErrQuantityNotPos
		}
		if l.UnitCost.IsNegative() {
			return Draft{}, ErrUnitCostNegative
		}
		if l.VATRate.IsNegative() {
			return Draft{}, ErrVATRateNegative
		}
		// line_total = quantity * unit_cost (decimal nhân decimal — KHÔNG float).
		lineTotal := l.UnitCost.Mul(l.Quantity)
		vatAmount := lineTotal.Mul(l.VATRate).RoundVND()
		computed = append(computed, ComputedLine{
			LineNo:            i + 1,
			ProductID:         l.ProductID,
			Quantity:          l.Quantity,
			UnitCost:          l.UnitCost,
			VATRate:           l.VATRate,
			BatchCode:         l.BatchCode,
			ExpiryDate:        l.ExpiryDate,
			ManufacturingDate: l.ManufacturingDate,
			LineTotal:         lineTotal,
			VATAmount:         vatAmount,
		})
		total = total.Add(lineTotal)
		vatTotal = vatTotal.Add(vatAmount)
	}

	return Draft{
		SupplierID:   in.SupplierID,
		SupplierName: in.SupplierName,
		Note:         in.Note,
		Status:       StatusDraft,
		TotalAmount:  total,
		VATAmount:    vatTotal,
		Lines:        computed,
	}, nil
}

// PurchaseOrder là một PO đã persist (header) ở mức domain (cho đường trả về).
type PurchaseOrder struct {
	ID           string
	Code         string
	SupplierID   *int64
	SupplierName string
	Status       string
	TotalAmount  money.Money
	VATAmount    money.Money
	Note         string
	LockVersion  int32
}

// PurchaseLine là một dòng PO đã persist (cho đường trả về / nhập kho + post sổ).
type PurchaseLine struct {
	ID                string
	LineNo            int
	ProductID         int64
	Quantity          decimal.Decimal
	UnitCost          money.Money
	VATRate           decimal.Decimal
	BatchCode         string
	ExpiryDate        *time.Time
	ManufacturingDate *time.Time
	LineTotal         money.Money
}
