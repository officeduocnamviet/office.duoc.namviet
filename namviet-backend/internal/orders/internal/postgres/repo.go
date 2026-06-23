// Package postgres là ADAPTER ra phía cơ sở dữ liệu của orders: implement port
// domain.Repository bằng query sinh từ sqlc (appdb) và map row <-> entity domain.
// ĐỌC bảng public.* kế thừa (strangler-fig, ADR 0001). Nằm dưới internal/ nên
// module khác KHÔNG import được. LÁT NÀY CHỈ ĐỌC → repo chỉ có thao tác đọc, bind
// thẳng pool (không tx). Đường ghi (tạo đơn, ghi phiếu thu, post sổ) HOÃN sang
// slice sau.
//
// Việc "dễ sai" nằm ở adapter này (không ở domain/SQL):
//   - map NUMERIC → money/Quantity KHÔNG đi qua float (dựng decimal từ mantissa).
//   - "đã thu" suy diễn ở SQL (subquery aggregate finance_transactions), repo chỉ
//     dựng PaymentSummary qua domain.ComputePayment (Remaining = Final - Paid).
//   - keyset cursor: app truyền unix nano (int64); repo chuyển sang
//     pgtype.Timestamptz. 0 = trang đầu → NULL (lấy từ mới nhất).
package postgres

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Repo implement domain.Repository trên appdb.Queries (bind pool).
type Repo struct{ q *appdb.Queries }

// NewRepo tạo repo từ một *appdb.Queries (đã bind pool).
func NewRepo(q *appdb.Queries) *Repo { return &Repo{q: q} }

func (r *Repo) ListOrders(ctx context.Context, f domain.OrderFilter) ([]domain.Order, error) {
	rows, err := r.q.ListOrders(ctx, appdb.ListOrdersParams{
		AfterCreatedAt: nanoToTimestamptz(f.AfterCreatedAt),
		AfterID:        strPtr(f.AfterID),
		CustomerID:     f.CustomerID,
		Status:         strPtr(f.Status),
		PaymentStatus:  strPtr(f.PaymentStatus),
		FromDate:       nanoToTimestamptz(f.FromDate),
		ToDate:         nanoToTimestamptz(f.ToDate),
		RowLimit:       f.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Order, 0, len(rows))
	for _, row := range rows {
		out = append(out, orderRowToDomain(orderRow{
			ID:            row.ID,
			Code:          row.Code,
			CustomerID:    row.CustomerID,
			CreatorID:     row.CreatorID,
			Status:        row.Status,
			OrderType:     row.OrderType,
			TotalAmount:   row.TotalAmount,
			FinalAmount:   row.FinalAmount,
			PaymentStatus: row.PaymentStatus,
			Note:          row.Note,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
			PaidAmount:    row.PaidAmount,
		}))
	}
	return out, nil
}

func (r *Repo) GetOrderByID(ctx context.Context, id string) (domain.Order, error) {
	row, err := r.q.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, apperr.NotFound("đơn hàng không tồn tại")
		}
		return domain.Order{}, err
	}
	return orderRowToDomain(orderRow{
		ID:            row.ID,
		Code:          row.Code,
		CustomerID:    row.CustomerID,
		CreatorID:     row.CreatorID,
		Status:        row.Status,
		OrderType:     row.OrderType,
		TotalAmount:   row.TotalAmount,
		FinalAmount:   row.FinalAmount,
		PaymentStatus: row.PaymentStatus,
		Note:          row.Note,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		PaidAmount:    row.PaidAmount,
	}), nil
}

func (r *Repo) ListLines(ctx context.Context, orderID string) ([]domain.OrderLine, error) {
	rows, err := r.q.ListOrderLines(ctx, orderID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.OrderLine, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.OrderLine{
			ID:         row.ID,
			ProductID:  row.ProductID,
			Quantity:   domain.QuantityFromInt(int64(row.Quantity)),
			UOM:        row.Uom,
			UnitPrice:  numericToMoney(row.UnitPrice),
			Discount:   numericToMoney(row.Discount),
			LineTotal:  numericToMoney(row.TotalLine),
			IsGift:     derefBool(row.IsGift),
			BatchNo:    derefStr(row.BatchNo),
			ExpiryDate: row.ExpiryDate.Time,
			HasExpiry:  row.ExpiryDate.Valid,
			Note:       derefStr(row.Note),
		})
	}
	return out, nil
}

// ---- mapping row <-> domain ----

// orderRow là tập cột chung của ListOrdersRow và GetOrderByIDRow (cùng SELECT).
// Dùng một struct trung gian để map một lần, tránh lặp (cùng pattern customers).
type orderRow struct {
	ID            string
	Code          string
	CustomerID    *int64
	CreatorID     *string
	Status        string
	OrderType     string
	TotalAmount   pgtype.Numeric
	FinalAmount   pgtype.Numeric
	PaymentStatus *string
	Note          *string
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
	PaidAmount    pgtype.Numeric
}

func orderRowToDomain(o orderRow) domain.Order {
	final := numericToMoney(o.FinalAmount)
	paid := numericToMoney(o.PaidAmount)
	return domain.Order{
		ID:            o.ID,
		Code:          o.Code,
		CustomerID:    o.CustomerID,
		CreatorID:     derefStr(o.CreatorID),
		Status:        o.Status,
		OrderType:     o.OrderType,
		Total:         numericToMoney(o.TotalAmount),
		Final:         final,
		PaymentStatus: derefStr(o.PaymentStatus),
		Note:          derefStr(o.Note),
		CreatedAt:     o.CreatedAt.Time,
		UpdatedAt:     o.UpdatedAt.Time,
		// "Đã thu" suy diễn ở SQL; domain quyết Remaining = Final - Paid (có thể âm).
		Payment: domain.ComputePayment(final, paid),
	}
}

// numericToMoney chuyển pgtype.Numeric sang money.Money KHÔNG đi qua float: dựng
// decimal trực tiếp từ mantissa (big.Int) * 10^Exp. NULL/NaN → Zero.
func numericToMoney(n pgtype.Numeric) money.Money {
	if !n.Valid || n.NaN || n.Int == nil {
		return money.Zero()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return money.FromDecimal(d)
}

// nanoToTimestamptz chuyển unix nano (app keyset/date filter) sang
// pgtype.Timestamptz. 0 = không có giá trị → invalid (SQL coi như NULL → trang
// đầu / không chặn ngày). KHÔNG đi qua float.
func nanoToTimestamptz(nano int64) pgtype.Timestamptz {
	if nano == 0 {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: time.Unix(0, nano).UTC(), Valid: true}
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

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Đảm bảo Repo thoả port domain ở compile-time.
var _ domain.Repository = (*Repo)(nil)
