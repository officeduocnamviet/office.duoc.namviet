// Package postgres — phần GHI (P4a): implement app.OrderStore trên appdb.Queries
// (bind tx do caller/TxManager truyền). GHI bảng public.orders + order_items kế
// thừa (strangler, ADR 0001). Tiền NUMERIC <-> common/money (KHÔNG float); quantity
// INTEGER <- domain.Quantity (lấy IntPart, đã validate > 0 ở domain). Mã đơn app
// sinh từ app.order_code_seq + tiền tố; id uuid v7 app-side. Đổi trạng thái dùng
// FOR UPDATE. Idempotency tạo đơn qua app.order_idempotency. KHÔNG đụng kho/tiền/sổ.
package postgres

import (
	"context"
	"errors"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// pgUniqueViolation là SQLSTATE 23505 (unique_violation) — trùng code đơn
// (ux_orders_code_alive) khi đua sinh mã (rất hiếm với sequence riêng).
const pgUniqueViolation = "23505"

// WriteRepo implement app.OrderStore trên appdb.Queries (bind tx). Mọi thao tác
// GHI + sinh mã + FOR UPDATE chạy trong CÙNG tx (atomic theo use-case).
type WriteRepo struct{ q *appdb.Queries }

// NewWriteRepo tạo repo ghi từ một *appdb.Queries (đã bind tx).
func NewWriteRepo(q *appdb.Queries) *WriteRepo { return &WriteRepo{q: q} }

func (r *WriteRepo) FindByIdemKey(ctx context.Context, idemKey string) (string, bool, error) {
	row, err := r.q.FindOrderByIdemKey(ctx, idemKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return row.OrderID, true, nil
}

func (r *WriteRepo) NextCodeSeq(ctx context.Context) (int64, error) {
	return r.q.NextOrderCodeSeq(ctx)
}

func (r *WriteRepo) InsertOrder(ctx context.Context, o app.NewOrderRow) (domain.Order, bool, error) {
	d := o.Draft
	var creatorID *string
	if d.CreatorID != "" {
		c := d.CreatorID
		creatorID = &c
	}
	row, err := r.q.InsertOrder(ctx, appdb.InsertOrderParams{
		ID:          o.ID,
		Code:        o.Code,
		CustomerID:  d.CustomerID,
		CreatorID:   creatorID,
		Status:      d.Status.String(),
		OrderType:   d.OrderType.String(),
		TotalAmount: moneyToNumeric(d.TotalAmount),
		FinalAmount: moneyToNumeric(d.FinalAmount),
		Note:        strPtr(d.Note),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.Order{}, true, nil // trùng code → service cấp số mới
		}
		return domain.Order{}, false, err
	}
	return insertOrderRowToDomain(row), false, nil
}

func (r *WriteRepo) InsertItem(ctx context.Context, orderID string, l domain.ComputedLine) error {
	return r.q.InsertOrderItem(ctx, appdb.InsertOrderItemParams{
		ID:        id.NewString(),
		OrderID:   orderID,
		ProductID: l.ProductID,
		Quantity:  int32(l.Quantity.Decimal().IntPart()), // INTEGER ở DB; domain đã validate > 0
		Uom:       l.UOM,
		UnitPrice: moneyToNumeric(l.UnitPrice),
		Discount:  moneyToNumeric(l.Discount),
		TotalLine: moneyToNumeric(l.LineTotal),
	})
}

func (r *WriteRepo) BindIdemKey(ctx context.Context, idemKey, orderID, orderCode string) (bool, error) {
	// ON CONFLICT DO NOTHING ở SQL. Để biết có chèn được không, đọc lại: nếu bản
	// ghi trỏ về orderID của ta thì inserted=true; nếu trỏ đơn khác thì luồng khác
	// thắng đua → inserted=false.
	if err := r.q.InsertOrderIdemKey(ctx, appdb.InsertOrderIdemKeyParams{
		IdemKey:   idemKey,
		OrderID:   orderID,
		OrderCode: orderCode,
	}); err != nil {
		return false, err
	}
	row, err := r.q.FindOrderByIdemKey(ctx, idemKey)
	if err != nil {
		return false, err
	}
	return row.OrderID == orderID, nil
}

func (r *WriteRepo) GetForUpdate(ctx context.Context, orderID string) (domain.Status, bool, error) {
	row, err := r.q.GetOrderStatusForUpdate(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return domain.Status(row.Status), true, nil
}

func (r *WriteRepo) UpdateStatus(ctx context.Context, orderID string, expected, next domain.Status) (int64, error) {
	return r.q.UpdateOrderStatus(ctx, appdb.UpdateOrderStatusParams{
		OrderID:        orderID,
		NewStatus:      next.String(),
		ExpectedStatus: expected.String(),
	})
}

// GetCreated nạp lại đơn (master + lines) theo id — tái dùng query ĐỌC (GetOrderByID
// + ListOrderLines) trên CÙNG tx để thấy bản vừa ghi. Map dùng helper chung với
// đường đọc (repo.go). Không thấy → apperr.NotFound.
func (r *WriteRepo) GetCreated(ctx context.Context, orderID string) (app.CreatedOrder, error) {
	row, err := r.q.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.CreatedOrder{}, apperr.NotFound("đơn hàng không tồn tại")
		}
		return app.CreatedOrder{}, err
	}
	order := orderRowToDomain(orderRow{
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
	})
	lineRows, err := r.q.ListOrderLines(ctx, orderID)
	if err != nil {
		return app.CreatedOrder{}, err
	}
	lines := make([]domain.OrderLine, 0, len(lineRows))
	for _, lr := range lineRows {
		lines = append(lines, domain.OrderLine{
			ID:         lr.ID,
			ProductID:  lr.ProductID,
			Quantity:   domain.QuantityFromInt(int64(lr.Quantity)),
			UOM:        lr.Uom,
			UnitPrice:  numericToMoney(lr.UnitPrice),
			Discount:   numericToMoney(lr.Discount),
			LineTotal:  numericToMoney(lr.TotalLine),
			IsGift:     derefBool(lr.IsGift),
			BatchNo:    derefStr(lr.BatchNo),
			ExpiryDate: lr.ExpiryDate.Time,
			HasExpiry:  lr.ExpiryDate.Valid,
			Note:       derefStr(lr.Note),
		})
	}
	return app.CreatedOrder{Order: order, Lines: lines}, nil
}

// insertOrderRowToDomain map InsertOrderRow (RETURNING) sang domain.Order. Payment
// để zero (đơn vừa tạo chưa có phiếu thu — suy diễn ở đường đọc khi cần).
func insertOrderRowToDomain(row appdb.InsertOrderRow) domain.Order {
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
	})
}

// moneyToNumeric chuyển money.Money sang pgtype.Numeric KHÔNG qua float: lấy
// coefficient (big.Int) + exponent từ decimal nền. Dùng khi GHI tiền (đối xứng
// với numericToMoney ở repo.go).
func moneyToNumeric(m money.Money) pgtype.Numeric {
	d := m.Decimal()
	return pgtype.Numeric{
		Int:   new(big.Int).Set(d.Coefficient()),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

// Đảm bảo WriteRepo thoả port app ở compile-time.
var _ app.OrderStore = (*WriteRepo)(nil)
