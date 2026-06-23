package domain

import "github.com/shopspring/decimal"

// Quantity là một LƯỢNG TỒN KHO chính xác thập phân. Bọc decimal.Decimal vì cột
// quantity/stock_quantity ở DB này là NUMERIC (KHÔNG phải int) — ADR 0002 §B coi
// "thống nhất kiểu quantity (int vs numeric)" là việc migrate sau. Tới lúc đó ta
// PHẢI tôn trọng kiểu thật: map NUMERIC → decimal, TUYỆT ĐỐI KHÔNG ép float (mất
// chính xác). Đây là shared-kernel-cục-bộ của domain inventory; chỉ phụ thuộc
// shopspring/decimal (thư viện số thập phân, KHÔNG phải hạ tầng pgx/http) nên
// domain thuần được phép dùng — y như common/money. Việc chuyển pgtype.Numeric
// (pgx) <-> Quantity nằm ở tầng adapter (repo), không ở domain.
type Quantity struct {
	d decimal.Decimal
}

// ZeroQty trả lượng tồn 0.
func ZeroQty() Quantity {
	return Quantity{d: decimal.Zero}
}

// QuantityFromDecimal bọc một decimal.Decimal có sẵn (vd repo đã chuyển từ NUMERIC).
func QuantityFromDecimal(d decimal.Decimal) Quantity {
	return Quantity{d: d}
}

// QuantityFromInt dựng Quantity từ số nguyên (vd order_items.quantity là INTEGER ở
// DB này). Giữ chính xác qua decimal, KHÔNG float.
func QuantityFromInt(n int64) Quantity {
	return Quantity{d: decimal.NewFromInt(n)}
}

// QuantityFromString parse một chuỗi số thập phân (vd "12.5") thành Quantity.
// Chuỗi rỗng coi là 0 (cột NUMERIC NULL/để trống map về 0). Lỗi parse trả error
// để caller xử lý — KHÔNG nuốt lỗi thành 0 ngầm.
func QuantityFromString(s string) (Quantity, error) {
	if s == "" {
		return ZeroQty(), nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Quantity{}, err
	}
	return Quantity{d: d}, nil
}

// Decimal trả giá trị decimal nền (để adapter chuyển ngược về NUMERIC khi ghi).
func (q Quantity) Decimal() decimal.Decimal { return q.d }

// String trả biểu diễn thập phân chuẩn (không số 0 thừa), vd "0", "12.5".
func (q Quantity) String() string { return q.d.String() }

// IsZero trả true nếu lượng tồn bằng 0.
func (q Quantity) IsZero() bool { return q.d.IsZero() }

// IsPositive trả true nếu lượng tồn > 0 (lô còn hàng để xuất).
func (q Quantity) IsPositive() bool { return q.d.IsPositive() }

// Equal so sánh BẰNG giá trị (12.50 == 12.5).
func (q Quantity) Equal(other Quantity) bool { return q.d.Equal(other.d) }

// Add trả q + other (thập phân chính xác, KHÔNG float). Dùng cộng dồn tồn các lô
// khi lập kế hoạch tiêu thụ FEFO.
func (q Quantity) Add(other Quantity) Quantity { return Quantity{d: q.d.Add(other.d)} }

// Sub trả q - other (thập phân chính xác, KHÔNG float). Dùng tính lượng CÒN THIẾU
// sau khi trừ tồn một lô. Kết quả có thể ÂM (lô thừa so với nhu cầu) — caller (kế
// hoạch FEFO) tự quyết cách dùng; Quantity KHÔNG tự kẹp về 0 (tránh che lỗi).
func (q Quantity) Sub(other Quantity) Quantity { return Quantity{d: q.d.Sub(other.d)} }

// Cmp so sánh: trả -1 nếu q < other, 0 nếu bằng, 1 nếu q > other. Dùng quyết định
// lô cuối tiêu thụ MỘT PHẦN (tồn lô > nhu cầu còn lại) hay TOÀN BỘ.
func (q Quantity) Cmp(other Quantity) int { return q.d.Cmp(other.d) }

// GreaterThanOrEqual trả true nếu q >= other. Dùng kiểm tổng tồn khả dụng có đủ
// nhu cầu trừ kho không (đủ → lập kế hoạch; thiếu → ErrInsufficientStock).
func (q Quantity) GreaterThanOrEqual(other Quantity) bool { return q.d.GreaterThanOrEqual(other.d) }
