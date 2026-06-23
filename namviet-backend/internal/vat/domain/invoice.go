// Package domain là LÕI THUẦN của bounded context vat: hoá đơn GTGT (VAT) của
// hệ Nam Việt. Entity hoá đơn (Invoice) + dòng hoá đơn (InvoiceLine) cùng bất
// biến tính/cân tiền THUẦN (Validate/Compute). KHÔNG import pgx/http/huma/
// framework (ARCHITECTURE.md §3) — chỉ stdlib + shared kernel trung lập
// (common/money + shopspring/decimal). Phụ thuộc một chiều: adapters → app →
// domain.
//
// BỐI CẢNH NGHIỆP VỤ (đúng đắn tiền > mọi thứ — ARCHITECTURE.md §7):
//   - HĐ VAT 100% đơn B2B; MST khách (CustomerTaxCode) BẮT BUỘC — thiếu = 422.
//   - HĐ thuộc SỔ TAX, theo GIÁ HOÁ ĐƠN (unit_price trên HĐ, có thể KHÁC giá
//     bán thật — dual-ledger). Domain KHÔNG biết sổ INTERNAL.
//   - Tiền = common/money (decimal) ↔ NUMERIC. CẤM float.
//   - vat_rate là INPUT từng dòng (vd 0.05/0.08/0.10) — KHÔNG hardcode.
//
// QUY TẮC LÀM TRÒN VAT (để Σ luôn cân):
//   - VAT từng dòng = ROUND_VND(line_amount × vat_rate) — làm tròn về ĐỒNG
//     (scale-0, HALF-UP) NGAY Ở TỪNG DÒNG (per-line rounding).
//   - subtotal = Σ line_amount; vat_amount = Σ line_vat (đã tròn); total =
//     subtotal + vat_amount. Vì tổng các số đã-tròn nên Σ luôn cân khít với
//     header — KHÔNG có sai số làm tròn lệch giữa dòng và tổng.
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// Status là trạng thái phát hành của một hoá đơn (khớp CHECK status ở DB).
type Status string

const (
	// StatusDraft — hoá đơn nháp (chưa cấp số/chưa phát hành). P5 hiện không dùng
	// (mọi HĐ phát hành ngay), chừa cho luồng nháp sau.
	StatusDraft Status = "draft"
	// StatusIssued — hoá đơn ĐÃ phát hành (đã cấp số gapless). Trạng thái mặc định
	// khi IssueInvoice.
	StatusIssued Status = "issued"
	// StatusCancelled — hoá đơn ĐÃ huỷ (defer: luồng huỷ/thay thế làm sau).
	StatusCancelled Status = "cancelled"
)

// Valid trả true nếu s là trạng thái hợp lệ.
func (s Status) Valid() bool {
	return s == StatusDraft || s == StatusIssued || s == StatusCancelled
}

// String trả giá trị chuỗi (adapter ghi cột status text).
func (s Status) String() string { return string(s) }

// Lỗi domain THUẦN (app map sang apperr.Validation → 422).
var (
	// ErrEmptyTaxCode — MST khách rỗng (HĐ VAT B2B bắt buộc có MST).
	ErrEmptyTaxCode = errors.New("MST khách hàng (customer_tax_code) bắt buộc cho hoá đơn VAT")
	// ErrNoLines — hoá đơn không có dòng nào.
	ErrNoLines = errors.New("hoá đơn phải có ít nhất một dòng")
	// ErrEmptyOrderCode — order_code rỗng (HĐ luôn gắn một đơn).
	ErrEmptyOrderCode = errors.New("mã đơn (order_code) không được rỗng")
	// ErrEmptySerial — ký hiệu HĐ (serial) rỗng (cần để cấp số gapless theo serial).
	ErrEmptySerial = errors.New("ký hiệu hoá đơn (serial) không được rỗng")
)

// LineInput là tham số THUẦN dựng một dòng hoá đơn (trước khi tính VAT/cân).
// Tiền = money decimal; vat_rate là decimal (KHÔNG float). app dựng từ DTO rồi
// gọi vào BuildInvoice.
type LineInput struct {
	// ProductID là id sản phẩm (public.products.id, bigint). 0 nếu dòng tự do.
	ProductID int64
	// Description diễn giải dòng (tên hàng trên HĐ).
	Description string
	// Quantity số lượng (money decimal — đủ chính xác, KHÔNG float). > 0.
	Quantity money.Money
	// UnitPrice đơn giá HOÁ ĐƠN (theo sổ TAX — có thể khác giá bán thật). >= 0.
	UnitPrice money.Money
	// VATRate thuế suất dòng (vd 0.08 = 8%). >= 0, là INPUT — KHÔNG hardcode.
	VATRate decimal.Decimal
}

// validate kiểm một dòng input THUẦN: số lượng/đơn giá/thuế suất không âm.
func (l LineInput) validate(idx int) error {
	if l.Quantity.IsNegative() {
		return fmt.Errorf("dòng %d: số lượng (quantity) không được âm", idx+1)
	}
	if l.UnitPrice.IsNegative() {
		return fmt.Errorf("dòng %d: đơn giá (unit_price) không được âm", idx+1)
	}
	if l.VATRate.IsNegative() {
		return fmt.Errorf("dòng %d: thuế suất (vat_rate) không được âm", idx+1)
	}
	return nil
}

// InvoiceLine là một dòng hoá đơn ĐÃ TÍNH: line_amount = quantity × unit_price;
// line_vat = ROUND_VND(line_amount × vat_rate). Tiền = money decimal, scale-0
// (VND) sau làm tròn.
type InvoiceLine struct {
	LineNo      int32
	ProductID   int64
	Description string
	Quantity    money.Money
	UnitPrice   money.Money
	VATRate     decimal.Decimal
	// LineAmount = quantity × unit_price (giá trị dòng trước thuế).
	LineAmount money.Money
	// LineVAT = ROUND_VND(line_amount × vat_rate) (thuế dòng, đã làm tròn về đồng).
	LineVAT money.Money
}

// Invoice là một HOÁ ĐƠN GTGT ĐÃ TÍNH (chưa cấp số/chưa persist). Subtotal/
// VATAmount/Total đã ép cân từ các dòng. Cấp số (Serial/InvoiceNo) + id + status
// + ngày là việc của tầng app/adapter — domain chỉ lo cấu trúc + cân tiền.
type Invoice struct {
	OrderCode       string
	CustomerTaxCode string
	Serial          string
	IssueDate       time.Time
	Subtotal        money.Money
	VATAmount       money.Money
	Total           money.Money
	Lines           []InvoiceLine
}

// BuildInvoice là HÀM THUẦN dựng + cân một hoá đơn từ input: validate cấu trúc
// (MST/serial/order_code/≥1 dòng/không âm), tính từng dòng (line_amount,
// line_vat làm tròn về đồng), rồi ép tổng (subtotal/vat_amount/total) bằng cộng
// money decimal — KHÔNG float. Trả Invoice cân khít hoặc lỗi domain (app map
// 422). KHÔNG cấp số ở đây.
func BuildInvoice(orderCode, customerTaxCode, serial string, issueDate time.Time, lines []LineInput) (Invoice, error) {
	if strings.TrimSpace(orderCode) == "" {
		return Invoice{}, ErrEmptyOrderCode
	}
	if strings.TrimSpace(customerTaxCode) == "" {
		return Invoice{}, ErrEmptyTaxCode
	}
	if strings.TrimSpace(serial) == "" {
		return Invoice{}, ErrEmptySerial
	}
	if len(lines) == 0 {
		return Invoice{}, ErrNoLines
	}

	subtotal := money.Zero()
	vatTotal := money.Zero()
	out := make([]InvoiceLine, 0, len(lines))
	for i, in := range lines {
		if err := in.validate(i); err != nil {
			return Invoice{}, err
		}
		lineAmount := in.Quantity.Mul(in.UnitPrice.Decimal()) // quantity × unit_price
		lineVAT := lineAmount.Mul(in.VATRate).RoundVND()      // ROUND_VND per-line
		subtotal = subtotal.Add(lineAmount)
		vatTotal = vatTotal.Add(lineVAT)
		out = append(out, InvoiceLine{
			LineNo:      int32(i + 1),
			ProductID:   in.ProductID,
			Description: in.Description,
			Quantity:    in.Quantity,
			UnitPrice:   in.UnitPrice,
			VATRate:     in.VATRate,
			LineAmount:  lineAmount,
			LineVAT:     lineVAT,
		})
	}

	inv := Invoice{
		OrderCode:       strings.TrimSpace(orderCode),
		CustomerTaxCode: strings.TrimSpace(customerTaxCode),
		Serial:          strings.TrimSpace(serial),
		IssueDate:       issueDate,
		Subtotal:        subtotal,
		VATAmount:       vatTotal,
		Total:           subtotal.Add(vatTotal),
		Lines:           out,
	}
	// Phòng thủ: tổng phải khít (total = subtotal + vat_amount). Vì ta tự dựng từ
	// Σ nên luôn đúng — assert để bắt mọi sai sót logic về sau (khớp CHECK DB).
	if !inv.Total.Equal(inv.Subtotal.Add(inv.VATAmount)) {
		return Invoice{}, fmt.Errorf("hoá đơn lệch tổng: total=%s <> subtotal=%s + vat=%s",
			inv.Total, inv.Subtotal, inv.VATAmount)
	}
	return inv, nil
}

// IssuedInvoice là một hoá đơn ĐÃ PHÁT HÀNH (đọc lại từ DB): cấu trúc Invoice +
// định danh/persisted (ID, InvoiceNo, Status, CreatedAt). Tách khỏi Invoice
// (thứ caller dựng để phát hành) để đường đọc và đường ghi không lẫn lộn.
type IssuedInvoice struct {
	ID        string
	Invoice   // nhúng: OrderCode/CustomerTaxCode/Serial/IssueDate/Subtotal/VATAmount/Total/Lines
	InvoiceNo int64
	Status    Status
	CreatedAt time.Time
}

// InvoiceFilter gom điều kiện lọc + keyset pagination cho danh sách HĐ đọc.
// Keyset theo (created_at DESC, id DESC). Tiêu chí optional.
type InvoiceFilter struct {
	AfterCreatedAt time.Time
	AfterID        string
	HasCursor      bool
	Limit          int32
	// OrderCode lọc theo mã đơn; rỗng = mọi đơn.
	OrderCode string
	// Status lọc theo trạng thái; rỗng = mọi trạng thái.
	Status string
}
