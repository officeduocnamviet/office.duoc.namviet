// Package money cung cấp kiểu Money bọc shopspring/decimal cho MỌI giá trị tiền
// trong hệ thống (ARCHITECTURE.md §3, §4: "Tiền = common/money (decimal) ↔
// NUMERIC. CẤM float"). Đây là shared kernel trung lập domain — chỉ phụ thuộc
// shopspring/decimal (thư viện số thập phân chính xác, KHÔNG phải hạ tầng
// pgx/http), nên domain thuần được phép dùng. Việc chuyển pgtype.Numeric (pgx)
// <-> Money nằm ở tầng adapter (repo), không ở domain.
package money

import "github.com/shopspring/decimal"

// Money là một lượng tiền chính xác thập phân. Bọc decimal.Decimal để không lộ
// kiểu thư viện ra API domain và để chặn lỡ tay dùng float ở money path.
type Money struct {
	d decimal.Decimal
}

// Zero trả lượng tiền 0.
func Zero() Money {
	return Money{d: decimal.Zero}
}

// FromDecimal bọc một decimal.Decimal có sẵn (vd repo đã chuyển từ NUMERIC).
func FromDecimal(d decimal.Decimal) Money {
	return Money{d: d}
}

// FromInt tạo Money từ một số nguyên (vd 110000 = 110.000 VND, scale-0). Tiện cho
// dựng bút toán / test với số tiền tròn. KHÔNG đi qua float.
func FromInt(v int64) Money {
	return Money{d: decimal.NewFromInt(v)}
}

// FromString parse một chuỗi số thập phân (vd "12345.67") thành Money. Chuỗi
// rỗng coi là 0 (cột NUMERIC NULL/để trống map về 0 tiền). Lỗi parse trả error
// để caller xử lý — KHÔNG nuốt lỗi thành 0 ngầm.
func FromString(s string) (Money, error) {
	if s == "" {
		return Zero(), nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, err
	}
	return Money{d: d}, nil
}

// Decimal trả giá trị decimal nền (để adapter chuyển ngược về NUMERIC khi ghi).
func (m Money) Decimal() decimal.Decimal { return m.d }

// String trả biểu diễn thập phân chuẩn (không số 0 thừa), vd "0", "12345.67".
func (m Money) String() string { return m.d.String() }

// Sub trả m - other (số thập phân chính xác, KHÔNG float). Dùng cho "còn nợ" =
// final_amount - đã_thu. Kết quả có thể ÂM (thu thừa) — caller quyết cách trình
// bày; money KHÔNG tự kẹp về 0 (để khỏi che lỗi/nhập trùng phiếu thu).
func (m Money) Sub(other Money) Money { return Money{d: m.d.Sub(other.d)} }

// Add trả m + other (số thập phân chính xác, KHÔNG float). Dùng cộng dồn Σdebit /
// Σcredit khi cân bút toán kép.
func (m Money) Add(other Money) Money { return Money{d: m.d.Add(other.d)} }

// Mul trả m * factor (số thập phân chính xác, KHÔNG float). factor là decimal
// (vd thuế suất 0.08). KHÔNG tự làm tròn — kết quả giữ nguyên scale để caller
// chủ động làm tròn (RoundVND) đúng quy tắc nghiệp vụ. Dùng tính VAT từng dòng
// (line_amount * vat_rate) trước khi làm tròn về đồng.
func (m Money) Mul(factor decimal.Decimal) Money { return Money{d: m.d.Mul(factor)} }

// RoundVND làm tròn về số nguyên đồng (scale 0) theo quy tắc HALF-UP (làm tròn
// nửa lên — 0.5 → 1, áp cho cả số dương). VND không có đơn vị nhỏ hơn đồng nên
// mọi tiền VND lưu/ghi sổ phải scale-0. Dùng sau Mul khi tính VAT từng dòng để
// Σ line_vat cân với vat_amount header. decimal.Round dùng half-away-from-zero
// (HALF_UP), phù hợp thông lệ làm tròn hoá đơn VN.
func (m Money) RoundVND() Money { return Money{d: m.d.Round(0)} }

// IsPositive trả true nếu lượng tiền > 0. Dùng kiểm mỗi dòng bút toán có đúng
// MỘT vế (debit hoặc credit) > 0.
func (m Money) IsPositive() bool { return m.d.IsPositive() }

// IsZero trả true nếu lượng tiền bằng 0.
func (m Money) IsZero() bool { return m.d.IsZero() }

// IsNegative trả true nếu lượng tiền < 0 (vd còn-nợ âm khi thu thừa).
func (m Money) IsNegative() bool { return m.d.IsNegative() }

// Equal so sánh BẰNG giá trị (12.50 == 12.5).
func (m Money) Equal(other Money) bool { return m.d.Equal(other.d) }
