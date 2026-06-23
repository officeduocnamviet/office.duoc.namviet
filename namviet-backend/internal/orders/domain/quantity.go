package domain

import "github.com/shopspring/decimal"

// Quantity là LƯỢNG của một dòng hàng, chính xác thập phân. Cột
// order_items.quantity ở DB này là INTEGER, nhưng ta bọc decimal.Decimal (KHÔNG
// float) để: (1) đồng nhất với inventory.domain.Quantity (cùng quy ước, dễ tái
// dùng khi nối orders↔inventory lúc trừ kho), và (2) an toàn nếu ADR 0002 §B
// thống nhất quantity về NUMERIC sau này. Chỉ phụ thuộc shopspring/decimal (thư
// viện số thập phân, KHÔNG phải hạ tầng pgx/http) nên domain thuần được phép
// dùng — y như common/money. Chuyển int4/NUMERIC (pgx) <-> Quantity nằm ở adapter
// (repo), không ở domain.
type Quantity struct {
	d decimal.Decimal
}

// ZeroQty trả lượng 0.
func ZeroQty() Quantity { return Quantity{d: decimal.Zero} }

// QuantityFromInt bọc một số nguyên (cột order_items.quantity là integer).
func QuantityFromInt(n int64) Quantity { return Quantity{d: decimal.NewFromInt(n)} }

// QuantityFromDecimal bọc một decimal.Decimal có sẵn.
func QuantityFromDecimal(d decimal.Decimal) Quantity { return Quantity{d: d} }

// Decimal trả giá trị decimal nền.
func (q Quantity) Decimal() decimal.Decimal { return q.d }

// String trả biểu diễn thập phân chuẩn (không số 0 thừa), vd "0", "3".
func (q Quantity) String() string { return q.d.String() }

// IsZero trả true nếu lượng bằng 0.
func (q Quantity) IsZero() bool { return q.d.IsZero() }

// Equal so sánh BẰNG giá trị.
func (q Quantity) Equal(other Quantity) bool { return q.d.Equal(other.d) }
