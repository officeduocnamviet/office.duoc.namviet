package domain

import "github.com/Maneva-AI/namviet-backend/internal/common/money"

// PaymentStatus là trạng thái thanh toán suy diễn của một đơn — khớp CHECK của
// public.orders (payment_status ∈ unpaid/partial/paid, default unpaid). KHÁC
// Status (trạng thái xử lý đơn): 3 tầng độc lập (finance_transactions ≠
// orders.status ≠ orders.payment_status — [[project_timo_bank_transfer_flow]]).
type PaymentStatus string

const (
	// PaymentUnpaid — chưa thu đồng nào (paid <= 0).
	PaymentUnpaid PaymentStatus = "unpaid"
	// PaymentPartial — đã thu một phần (0 < paid < final).
	PaymentPartial PaymentStatus = "partial"
	// PaymentPaid — đã thu đủ (paid >= final). Thu thừa cũng coi là paid (còn nợ âm).
	PaymentPaid PaymentStatus = "paid"
)

// String trả giá trị chuỗi (để adapter ghi cột payment_status text).
func (p PaymentStatus) String() string { return string(p) }

// DerivePaymentStatus suy trạng thái thanh toán THUẦN từ tổng phải trả (final) và
// đã thu (paid) — nguồn chân lý cho việc cập nhật orders.payment_status sau khi
// ghi phiếu thu. Quy ước (so sánh decimal, KHÔNG float):
//
//	paid <= 0          → unpaid
//	0 < paid < final   → partial
//	paid >= final      → paid (kể cả thu thừa)
//
// Trường hợp biên final <= 0 (đơn 0 đồng — hiếm): mọi paid >= 0 đều coi là paid
// (không còn gì để thu). Hàm THUẦN (dễ unit-test, không DB).
func DerivePaymentStatus(final, paid money.Money) PaymentStatus {
	if !paid.IsPositive() {
		// paid == 0 hoặc âm (không có phiếu thu hợp lệ).
		if !final.IsPositive() {
			return PaymentPaid // đơn 0 đồng: coi như đã tất toán.
		}
		return PaymentUnpaid
	}
	// paid > 0. Còn nợ = final - paid; <= 0 nghĩa là đã đủ (hoặc thừa).
	remaining := final.Sub(paid)
	if remaining.IsPositive() {
		return PaymentPartial
	}
	return PaymentPaid
}
