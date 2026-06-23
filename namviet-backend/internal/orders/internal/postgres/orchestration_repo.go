// orchestration_repo.go — phần GHI ORCHESTRATION (P4b): implement
// app.OrchestrationStore trên appdb.Queries (bind tx do caller/TxManager truyền).
// CHỈ chạm bảng public.orders + order_items + đọc finance_transactions (tính lại
// đã-thu). Trừ kho/HĐ/sổ/phiếu thu KHÔNG ở đây — orders/app gọi PORT module khác
// với CÙNG tx (gộp atomic). Tiền NUMERIC <-> common/money (KHÔNG float); quantity
// INTEGER -> inventory/domain.Quantity (decimal).
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	inventorydomain "github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// OrchestrationRepo implement app.OrchestrationStore trên appdb.Queries (bind tx).
// Mọi thao tác chạy trong CÙNG tx — gộp atomic với trừ kho + HĐ + post sổ + phiếu thu.
type OrchestrationRepo struct{ q *appdb.Queries }

// NewOrchestrationRepo tạo repo orchestration từ một *appdb.Queries (đã bind tx).
func NewOrchestrationRepo(q *appdb.Queries) *OrchestrationRepo { return &OrchestrationRepo{q: q} }

// GetHeaderForUpdate khoá dòng đơn (FOR UPDATE) + trả header đầy đủ. Không thấy →
// (_, false, nil).
func (r *OrchestrationRepo) GetHeaderForUpdate(ctx context.Context, orderID string) (app.OrchHeader, bool, error) {
	row, err := r.q.GetOrderHeaderForUpdate(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.OrchHeader{}, false, nil
		}
		return app.OrchHeader{}, false, err
	}
	return app.OrchHeader{
		ID:            row.ID,
		Code:          row.Code,
		CustomerID:    row.CustomerID,
		OrderType:     row.OrderType,
		Status:        row.Status,
		FinalAmount:   numericToMoney(row.FinalAmount),
		TotalAmount:   numericToMoney(row.TotalAmount),
		PaymentStatus: derefStr(row.PaymentStatus),
	}, true, nil
}

// ListLinesForOrder trả dòng hàng (chưa soft-delete) — product/qty/uom/đơn giá để
// trừ kho FEFO + dựng dòng HĐ. quantity INTEGER -> inventory Quantity (decimal).
func (r *OrchestrationRepo) ListLinesForOrder(ctx context.Context, orderID string) ([]app.OrchLine, error) {
	rows, err := r.q.ListOrderLines(ctx, orderID)
	if err != nil {
		return nil, err
	}
	out := make([]app.OrchLine, 0, len(rows))
	for _, lr := range rows {
		out = append(out, app.OrchLine{
			ProductID: lr.ProductID,
			Quantity:  inventorydomain.QuantityFromInt(int64(lr.Quantity)),
			UOM:       lr.Uom,
			UnitPrice: numericToMoney(lr.UnitPrice),
			Discount:  numericToMoney(lr.Discount),
			LineTotal: numericToMoney(lr.TotalLine),
		})
	}
	return out, nil
}

// SumPaidInTx trả tổng ĐÃ THU (sổ thực tế) của đơn theo code TRONG tx hiện hành.
func (r *OrchestrationRepo) SumPaidInTx(ctx context.Context, orderCode string) (money.Money, error) {
	num, err := r.q.SumOrderPaidInTx(ctx, orderCode)
	if err != nil {
		return money.Zero(), err
	}
	return numericToMoney(num), nil
}

// UpdatePaymentStatus đặt payment_status (unpaid/partial/paid). Trả số dòng đổi.
func (r *OrchestrationRepo) UpdatePaymentStatus(ctx context.Context, orderID, paymentStatus string) (int64, error) {
	return r.q.UpdateOrderPaymentStatus(ctx, appdb.UpdateOrderPaymentStatusParams{
		OrderID:       orderID,
		PaymentStatus: paymentStatus,
	})
}

// UpdateStatus đổi trạng thái xử lý (guard status cũ). Trả số dòng đổi.
func (r *OrchestrationRepo) UpdateStatus(ctx context.Context, orderID, expected, next string) (int64, error) {
	return r.q.UpdateOrderStatus(ctx, appdb.UpdateOrderStatusParams{
		OrderID:        orderID,
		NewStatus:      next,
		ExpectedStatus: expected,
	})
}

// ListUnpaidOrdersByCustomer trả đơn CHƯA tất toán của khách, CŨ NHẤT trước, KHOÁ
// FOR UPDATE (phân bổ lump-sum tuần tự). Kèm Final + đã-thu hiện tại.
func (r *OrchestrationRepo) ListUnpaidOrdersByCustomer(ctx context.Context, customerID int64) ([]app.UnpaidOrder, error) {
	rows, err := r.q.ListUnpaidOrdersByCustomerForUpdate(ctx, customerID)
	if err != nil {
		return nil, err
	}
	out := make([]app.UnpaidOrder, 0, len(rows))
	for _, row := range rows {
		out = append(out, app.UnpaidOrder{
			ID:    row.ID,
			Code:  row.Code,
			Final: numericToMoney(row.FinalAmount),
			Paid:  numericToMoney(row.PaidAmount),
		})
	}
	return out, nil
}

// InsertAllocation ghi/cộng dồn dòng phân bổ phiếu (paymentID) → đơn (orderCode).
func (r *OrchestrationRepo) InsertAllocation(ctx context.Context, paymentID int64, orderCode string, amount money.Money) error {
	return r.q.InsertOrderAllocation(ctx, appdb.InsertOrderAllocationParams{
		PaymentID: paymentID,
		OrderCode: orderCode,
		Amount:    moneyToNumeric(amount),
	})
}

// Đảm bảo OrchestrationRepo thoả port app ở compile-time.
var _ app.OrchestrationStore = (*OrchestrationRepo)(nil)
