package app

import "fmt"

// Cấu hình sinh mã đơn — MỘT CHỖ DUY NHẤT (đổi quy ước ở đây).
//
// ⚠️ CỜ XÁC NHẬN PROD: quy ước mã đơn THẬT của Nam Việt (tiền tố 'DH' đơn-hàng vs
// 'HD' hoá-đơn lịch sử, độ rộng zero-pad, có gắn ngày không) CẦN kế toán/BA xác
// nhận. Đề xuất hiện tại: tiền tố "DH" + zero-pad 8 chữ số → "DH00000123".
// Sequence app.order_code_seq là riêng (KHÔNG đụng mã 'HD' cũ); khi go-live nếu
// đổi sang tiền tố trùng tập mã cũ thì PHẢI setval start > max hiện có (xem
// migration 00005 + db-review).
const (
	// codePrefix gắn đầu mã đơn (phân biệt mã backend sinh với mã ERP lịch sử).
	codePrefix = "DH"
	// codePadWidth số chữ số zero-pad cho phần số (đủ rộng cho vòng đời dự kiến).
	codePadWidth = 8
)

// formatOrderCode ghép tiền tố + zero-pad(seq) thành mã đơn, vd seq=123 →
// "DH00000123". seq lấy từ app.order_code_seq (duy nhất, an toàn đua).
func formatOrderCode(seq int64) string {
	return fmt.Sprintf("%s%0*d", codePrefix, codePadWidth, seq)
}
