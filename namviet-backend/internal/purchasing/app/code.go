package app

import "fmt"

// Cấu hình sinh mã PO — MỘT CHỖ DUY NHẤT (đổi quy ước ở đây).
//
// ⚠️ CỜ XÁC NHẬN PROD: quy ước mã đơn MUA THẬT của Nam Việt (tiền tố, zero-pad, có
// gắn ngày không) CẦN kế toán/BA xác nhận. Đề xuất hiện tại: "PO" + zero-pad 8 chữ
// số → "PO00000123". Sequence app.purchase_order_code_seq riêng (KHÔNG đụng mã đơn
// bán); khi go-live nếu đổi tiền tố trùng tập mã cũ thì PHẢI setval start > max.
const (
	codePrefix   = "PO"
	codePadWidth = 8
)

// formatPOCode ghép tiền tố + zero-pad(seq) thành mã PO, vd seq=123 → "PO00000123".
func formatPOCode(seq int64) string {
	return fmt.Sprintf("%s%0*d", codePrefix, codePadWidth, seq)
}
