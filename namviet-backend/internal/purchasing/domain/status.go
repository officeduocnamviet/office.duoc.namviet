// Package domain là LÕI THUẦN của bounded context purchasing (mua hàng & nhập kho,
// spec mục 54) — chiều MUA, đối xứng chiều BÁN (orders). Entity PurchaseOrder +
// PurchaseLine + value object Status (state machine THUẦN). KHÔNG import
// pgx/http/huma/framework (ARCHITECTURE.md §3) — chỉ stdlib + shared kernel trung
// lập (common/money). Phụ thuộc một chiều: adapters → app → domain.
package domain

// Status là trạng thái xử lý của một đơn mua (PO) — khớp CHECK của
// app.purchase_orders (status ∈ draft/ordered/received/paid/cancelled). Dùng cho
// state machine THUẦN (CanTransition) + dựng draft. Giá trị LOWERCASE (khác orders
// dùng UPPERCASE) — khớp đúng CHECK migration 00007.
type Status string

const (
	// StatusDraft — PO vừa tạo, chưa đặt hàng (default DB).
	StatusDraft Status = "draft"
	// StatusOrdered — đã gửi đơn đặt cho NCC; chưa nhận hàng.
	StatusOrdered Status = "ordered"
	// StatusReceived — đã nhận & nhập kho (tạo lô + tăng tồn + post sổ Dr 1561+133/Cr 331).
	StatusReceived Status = "received"
	// StatusPaid — đã thanh toán NCC (chi tiền Dr 331/Cr 111/112). Terminal hạnh phúc.
	StatusPaid Status = "paid"
	// StatusCancelled — huỷ PO (terminal). Chỉ huỷ được khi chưa nhập kho (draft/ordered).
	StatusCancelled Status = "cancelled"
)

// String trả giá trị chuỗi (để adapter ghi cột status text).
func (s Status) String() string { return string(s) }

// Valid trả true nếu s là một trạng thái hợp lệ (khớp CHECK ở DB). Bảo vệ trước
// khi ghi: không bao giờ để giá trị ngoài tập lọt vào cột.
func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusOrdered, StatusReceived, StatusPaid, StatusCancelled:
		return true
	default:
		return false
	}
}

// transitions là đồ thị chuyển trạng thái HỢP LỆ (đơn-hướng, tiến theo vòng đời mua
// hàng). Quy ước (spec mục 54):
//
//	draft    → ordered  | cancelled
//	ordered  → received | cancelled
//	received → paid
//	paid, cancelled     → (terminal, không chuyển tiếp)
//
// KHÔNG cho nhảy bước (draft→received) hay lùi. Tự-chuyển (from==to) KHÔNG hợp lệ.
// Sau khi đã received (đã nhập kho thật + post sổ) KHÔNG cho cancelled — huỷ lúc đó
// cần đảo bút toán + xuất kho (HOÃN, ngoài CORE).
var transitions = map[Status][]Status{
	StatusDraft:     {StatusOrdered, StatusCancelled},
	StatusOrdered:   {StatusReceived, StatusCancelled},
	StatusReceived:  {StatusPaid},
	StatusPaid:      {},
	StatusCancelled: {},
}

// CanTransition trả true nếu được phép chuyển from → to theo đồ thị state machine.
// Hàm THUẦN (không DB), là nguồn chân lý cho mọi use-case đổi trạng thái. from/to
// ngoài tập hợp lệ → false.
func CanTransition(from, to Status) bool {
	for _, allowed := range transitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}
