// Package postgres là ADAPTER ra phía cơ sở dữ liệu của customers: implement port
// domain.Repository bằng query sinh từ sqlc (appdb) và map row <-> entity domain.
// ĐỌC bảng public.* kế thừa (strangler-fig, ADR 0001). Nằm dưới internal/ nên
// module khác KHÔNG import được. Customers read-mostly → repo chỉ có thao tác
// đọc, bind thẳng pool (không tx).
//
// Hai việc "dễ sai" nằm ở adapter này (không ở domain/SQL):
//  1. parse b2b_metadata jsonb → domain.B2BProfile (DebtLimit qua money, KHÔNG
//     float; chấp nhận debt_limit là number hoặc string trong JSON).
//  2. chọn nguồn công nợ: query đã COALESCE live về numeric hợp lệ (0 khi không
//     có đơn) nên LiveDebt.Valid luôn true ở DB này → luôn ưu tiên LIVE
//     (domain.SelectDebt). Cột tĩnh current_debt chỉ là dự phòng phòng thủ.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Repo implement domain.Repository trên appdb.Queries (bind pool).
type Repo struct{ q *appdb.Queries }

// NewRepo tạo repo từ một *appdb.Queries (đã bind pool).
func NewRepo(q *appdb.Queries) *Repo { return &Repo{q: q} }

func (r *Repo) ListCustomers(ctx context.Context, f domain.CustomerFilter) ([]domain.Customer, error) {
	rows, err := r.q.ListCustomers(ctx, appdb.ListCustomersParams{
		AfterID:      f.AfterID,
		RowLimit:     f.Limit,
		CustomerType: strPtr(string(f.Type)),
		Status:       strPtr(f.Status),
		Q:            strPtr(f.Query),
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Customer, 0, len(rows))
	for _, row := range rows {
		out = append(out, customerRowToDomain(customerRow(row)))
	}
	return out, nil
}

func (r *Repo) GetCustomerByID(ctx context.Context, id int64) (domain.Customer, error) {
	row, err := r.q.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Customer{}, apperr.NotFound("khách hàng không tồn tại")
		}
		return domain.Customer{}, err
	}
	return customerRowToDomain(customerRow(row)), nil
}

// ---- mapping row <-> domain ----

// customerRow là tập cột chung của ListCustomersRow và GetCustomerByIDRow (cùng
// SELECT). Dùng một struct trung gian để map một lần, tránh lặp hai hàm giống
// nhau (cùng pattern với catalog.productRow).
type customerRow struct {
	ID           int64
	CustomerCode *string
	Name         string
	CustomerType string
	Phone        *string
	Email        *string
	Address      *string
	Status       string
	B2bMetadata  []byte
	CurrentDebt  pgtype.Numeric
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
	LiveDebt     pgtype.Numeric
}

func customerRowToDomain(c customerRow) domain.Customer {
	typ := domain.CustomerType(c.CustomerType)

	static := numericToMoney(c.CurrentDebt)
	live := numericToMoney(c.LiveDebt)
	// Query COALESCE live về numeric hợp lệ (0 khi không có đơn). LiveDebt.Valid
	// = true nghĩa là có con số live đáng tin → ưu tiên LIVE (domain quyết).
	debt := domain.SelectDebt(live, static, c.LiveDebt.Valid)

	cust := domain.Customer{
		ID:        c.ID,
		Code:      derefStr(c.CustomerCode),
		Name:      c.Name,
		Type:      typ,
		Phone:     derefStr(c.Phone),
		Email:     derefStr(c.Email),
		Address:   derefStr(c.Address),
		Status:    c.Status,
		Debt:      debt,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
	// Chỉ gắn B2BProfile cho khách B2B (B2C không có hồ sơ doanh nghiệp).
	if typ == domain.TypeB2B {
		if p := parseB2BMetadata(c.B2bMetadata); p != nil {
			cust.B2B = p
		} else {
			// B2B nhưng metadata rỗng → vẫn trả profile rỗng để FE biết là B2B.
			cust.B2B = &domain.B2BProfile{}
		}
	}
	return cust
}

// parseB2BMetadata parse jsonb b2b_metadata → domain.B2BProfile. nil/rỗng/"{}"
// → nil (caller quyết gắn profile rỗng hay không). DebtLimit qua money.FromString
// (decimal, KHÔNG float).
//
// Dữ liệu cũ BẨN: debt_limit có thể là number, string số ("50000000"), hoặc rác
// ("abc"); cùng trường nhiều kiểu giữa các bản ghi. Vì vậy đọc từng trường qua
// json.RawMessage + decode phòng thủ — một trường lỗi KHÔNG làm hỏng cả profile
// (chỉ bỏ qua trường đó), và một giá trị không phải-object cho cả metadata thì
// trả nil thay vì panic.
func parseB2BMetadata(raw []byte) *domain.B2BProfile {
	if len(raw) == 0 {
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || len(fields) == 0 {
		return nil
	}

	taxCode := jsonString(fields["tax_code"])
	salesStaff := jsonString(fields["sales_staff_id"])
	debtLimit, hasDebt := jsonMoney(fields["debt_limit"])
	paymentTerm, hasTerm := jsonInt(fields["payment_term"])

	if taxCode == "" && salesStaff == "" && !hasDebt && !hasTerm {
		return nil
	}
	p := &domain.B2BProfile{
		TaxCode:      taxCode,
		SalesStaffID: salesStaff,
		PaymentTerm:  paymentTerm,
		DebtLimit:    money.Zero(),
	}
	if hasDebt {
		p.DebtLimit = debtLimit
	}
	return p
}

// jsonString đọc một RawMessage thành string. Chấp nhận cả JSON string ("x") lẫn
// số (123 → "123"). Rỗng/null/không decode được → "".
func jsonString(r json.RawMessage) string {
	if len(r) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(r, &s); err == nil {
		return s
	}
	var n json.Number
	if err := json.Unmarshal(r, &n); err == nil {
		return n.String()
	}
	return ""
}

// jsonMoney đọc debt_limit (number HOẶC string số) thành money.Money. Rác/null →
// (Zero, false) để caller biết "không có giá trị hợp lệ".
func jsonMoney(r json.RawMessage) (money.Money, bool) {
	s := jsonString(r)
	if s == "" {
		return money.Zero(), false
	}
	m, err := money.FromString(s)
	if err != nil {
		return money.Zero(), false
	}
	return m, true
}

// jsonInt đọc payment_term (number hoặc string số) thành int. Rác/null →
// (0, false).
func jsonInt(r json.RawMessage) (int, bool) {
	s := jsonString(r)
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

// numericToMoney chuyển pgtype.Numeric (do sqlc sinh) sang money.Money KHÔNG đi
// qua float: dựng decimal trực tiếp từ mantissa (big.Int) * 10^Exp. NULL/NaN →
// Zero (cột tiền NULL coi như 0 tiền).
func numericToMoney(n pgtype.Numeric) money.Money {
	if !n.Valid || n.NaN || n.Int == nil {
		return money.Zero()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return money.FromDecimal(d)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Đảm bảo Repo thoả port domain ở compile-time.
var _ domain.Repository = (*Repo)(nil)
