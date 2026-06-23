package domain

// Status là trạng thái xử lý đơn — kiểu domain khớp CHECK của public.orders
// (status ∈ PENDING/CONFIRMED/SHIPPING/COMPLETED/CANCELLED/REFUNDED, xác nhận từ
// migration init 20260613080000_init_erp_schema.sql). Dùng cho state machine
// THUẦN (CanTransition) + draft. Đường ĐỌC cũ (Order.Status) giữ string để khỏi
// đụng mapping read — adapter set/đọc cột text trực tiếp.
type Status string

const (
	// StatusPending — đơn vừa tạo, chưa duyệt (default DB).
	StatusPending Status = "PENDING"
	// StatusConfirmed — đã duyệt; chưa xuất kho.
	StatusConfirmed Status = "CONFIRMED"
	// StatusShipping — đang giao (đã trừ kho FEFO — chuyển vào trạng thái này thuộc P4b).
	StatusShipping Status = "SHIPPING"
	// StatusCompleted — giao xong, đơn hoàn tất.
	StatusCompleted Status = "COMPLETED"
	// StatusCancelled — huỷ đơn (terminal).
	StatusCancelled Status = "CANCELLED"
	// StatusRefunded — hoàn trả (terminal; nghiệp vụ đảo sổ/hoàn kho thuộc P4b/sau).
	StatusRefunded Status = "REFUNDED"
)

// String trả giá trị chuỗi (để adapter ghi cột status text).
func (s Status) String() string { return string(s) }

// Valid trả true nếu s là một trạng thái hợp lệ (khớp CHECK ở DB). Bảo vệ trước
// khi ghi: không bao giờ để giá trị ngoài tập lọt vào cột.
func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusShipping,
		StatusCompleted, StatusCancelled, StatusRefunded:
		return true
	default:
		return false
	}
}

// transitions là đồ thị chuyển trạng thái HỢP LỆ (đơn-hướng, tiến theo luồng bán
// hàng). Quy ước (spec §P4):
//
//	PENDING   → CONFIRMED | CANCELLED
//	CONFIRMED → SHIPPING  | CANCELLED
//	SHIPPING  → COMPLETED
//	COMPLETED → REFUNDED
//	CANCELLED, REFUNDED   → (terminal, không chuyển tiếp)
//
// KHÔNG cho nhảy bước (PENDING→COMPLETED) hay lùi (CONFIRMED→PENDING). Tự-chuyển
// (from==to) KHÔNG hợp lệ (không phải một bước chuyển).
//
// Lưu ý phạm vi: chuyển SHIPPING (cần trừ kho FEFO + post sổ) và REFUNDED (cần
// đảo bút toán + hoàn kho) là HỢP LỆ về STATE nhưng use-case thực thi cần
// primitive cross-module — HOÃN sang P4b. P4a chỉ dùng các chuyển không đụng
// kho/tiền/sổ: PENDING→CONFIRMED, SHIPPING→COMPLETED, →CANCELLED.
var transitions = map[Status][]Status{
	StatusPending:   {StatusConfirmed, StatusCancelled},
	StatusConfirmed: {StatusShipping, StatusCancelled},
	StatusShipping:  {StatusCompleted},
	StatusCompleted: {StatusRefunded},
	StatusCancelled: {},
	StatusRefunded:  {},
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
