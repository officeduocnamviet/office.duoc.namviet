// Package postgres là ADAPTER ra DB của purchasing: implement app.Store +
// app.ListReader trên appdb.Queries (bind tx do caller/TxManager truyền). GHI/ĐỌC
// bảng MỚI app.purchase_orders + purchase_order_items (goose 00007). Nhập kho/post
// sổ/chi NCC KHÔNG ở đây — purchasing/app gọi PORT module khác với CÙNG tx (gộp
// atomic). Tiền NUMERIC <-> common/money; quantity/vat_rate NUMERIC <-> decimal
// (KHÔNG float). Nằm dưới internal/ nên module khác KHÔNG import được.
package postgres

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	inventorydomain "github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/app"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

// pgUniqueViolation là SQLSTATE 23505 — trùng code PO (ux_purchase_orders_code_alive).
const pgUniqueViolation = "23505"

// Repo implement app.Store + app.ListReader trên appdb.Queries (bind tx). Mọi thao
// tác GHI + sinh mã + FOR UPDATE chạy trong CÙNG tx (atomic theo use-case).
type Repo struct{ q *appdb.Queries }

// NewRepo tạo repo từ một *appdb.Queries (đã bind tx).
func NewRepo(q *appdb.Queries) *Repo { return &Repo{q: q} }

func (r *Repo) NextCodeSeq(ctx context.Context) (int64, error) {
	return r.q.NextPurchaseOrderCodeSeq(ctx)
}

func (r *Repo) InsertPO(ctx context.Context, row app.NewPORow) (domain.PurchaseOrder, bool, error) {
	d := row.Draft
	res, err := r.q.InsertPurchaseOrder(ctx, appdb.InsertPurchaseOrderParams{
		ID:           row.ID,
		Code:         row.Code,
		SupplierID:   d.SupplierID,
		SupplierName: strPtr(d.SupplierName),
		Status:       d.Status.String(),
		TotalAmount:  moneyToNumeric(d.TotalAmount),
		VatAmount:    moneyToNumeric(d.VATAmount),
		Note:         strPtr(d.Note),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.PurchaseOrder{}, true, nil // trùng code → service cấp số mới
		}
		return domain.PurchaseOrder{}, false, err
	}
	return domain.PurchaseOrder{
		ID:           res.ID,
		Code:         res.Code,
		SupplierID:   res.SupplierID,
		SupplierName: derefStr(res.SupplierName),
		Status:       res.Status,
		TotalAmount:  numericToMoney(res.TotalAmount),
		VATAmount:    numericToMoney(res.VatAmount),
		Note:         derefStr(res.Note),
		LockVersion:  res.LockVersion,
	}, false, nil
}

func (r *Repo) InsertPOItem(ctx context.Context, poID string, l domain.ComputedLine) error {
	return r.q.InsertPurchaseOrderItem(ctx, appdb.InsertPurchaseOrderItemParams{
		ID:                id.NewString(),
		PoID:              poID,
		LineNo:            int32(l.LineNo),
		ProductID:         l.ProductID,
		Quantity:          decimalToNumeric(l.Quantity),
		UnitCost:          moneyToNumeric(l.UnitCost),
		VatRate:           decimalToNumeric(l.VATRate),
		BatchCode:         optStr(l.BatchCode),
		ExpiryDate:        datePtr(l.ExpiryDate),
		ManufacturingDate: datePtr(l.ManufacturingDate),
		LineTotal:         moneyToNumeric(l.LineTotal),
	})
}

func (r *Repo) GetCreated(ctx context.Context, poID string) (app.CreatedPO, error) {
	row, err := r.q.GetPurchaseOrderByID(ctx, poID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.CreatedPO{}, apperr.NotFound("đơn mua không tồn tại")
		}
		return app.CreatedPO{}, err
	}
	po := domain.PurchaseOrder{
		ID:           row.ID,
		Code:         row.Code,
		SupplierID:   row.SupplierID,
		SupplierName: derefStr(row.SupplierName),
		Status:       row.Status,
		TotalAmount:  numericToMoney(row.TotalAmount),
		VATAmount:    numericToMoney(row.VatAmount),
		Note:         derefStr(row.Note),
		LockVersion:  row.LockVersion,
	}
	lineRows, err := r.q.ListPurchaseOrderLines(ctx, poID)
	if err != nil {
		return app.CreatedPO{}, err
	}
	lines := make([]domain.PurchaseLine, 0, len(lineRows))
	for _, lr := range lineRows {
		lines = append(lines, domain.PurchaseLine{
			ID:                lr.ID,
			LineNo:            int(lr.LineNo),
			ProductID:         lr.ProductID,
			Quantity:          numericToDecimal(lr.Quantity),
			UnitCost:          numericToMoney(lr.UnitCost),
			VATRate:           numericToDecimal(lr.VatRate),
			BatchCode:         derefStr(lr.BatchCode),
			ExpiryDate:        datePtrOut(lr.ExpiryDate),
			ManufacturingDate: datePtrOut(lr.ManufacturingDate),
			LineTotal:         numericToMoney(lr.LineTotal),
		})
	}
	return app.CreatedPO{PO: po, Lines: lines}, nil
}

func (r *Repo) GetHeaderForUpdate(ctx context.Context, poID string) (app.POHeader, bool, error) {
	row, err := r.q.GetPurchaseOrderHeaderForUpdate(ctx, poID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.POHeader{}, false, nil
		}
		return app.POHeader{}, false, err
	}
	return app.POHeader{
		ID:           row.ID,
		Code:         row.Code,
		SupplierID:   row.SupplierID,
		SupplierName: derefStr(row.SupplierName),
		Status:       row.Status,
		TotalAmount:  numericToMoney(row.TotalAmount),
		VATAmount:    numericToMoney(row.VatAmount),
		Note:         derefStr(row.Note),
		LockVersion:  row.LockVersion,
	}, true, nil
}

func (r *Repo) ListLines(ctx context.Context, poID string) ([]app.POLine, error) {
	rows, err := r.q.ListPurchaseOrderLines(ctx, poID)
	if err != nil {
		return nil, err
	}
	out := make([]app.POLine, 0, len(rows))
	for _, lr := range rows {
		out = append(out, app.POLine{
			ID:                lr.ID,
			LineNo:            int(lr.LineNo),
			ProductID:         lr.ProductID,
			Quantity:          inventorydomain.QuantityFromDecimal(numericToDecimal(lr.Quantity)),
			UnitCost:          numericToMoney(lr.UnitCost),
			VATRate:           money.FromDecimal(numericToDecimal(lr.VatRate)),
			BatchCode:         derefStr(lr.BatchCode),
			ExpiryDate:        datePtrOut(lr.ExpiryDate),
			ManufacturingDate: datePtrOut(lr.ManufacturingDate),
			LineTotal:         numericToMoney(lr.LineTotal),
		})
	}
	return out, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, poID, expected, next string) (int64, error) {
	return r.q.UpdatePurchaseOrderStatus(ctx, appdb.UpdatePurchaseOrderStatusParams{
		ID:             poID,
		NewStatus:      next,
		ExpectedStatus: expected,
	})
}

func (r *Repo) FindByIdemKey(ctx context.Context, idemKey string) (string, bool, error) {
	row, err := r.q.FindPurchaseOrderByIdemKey(ctx, idemKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return row.PoID, true, nil
}

func (r *Repo) BindIdemKey(ctx context.Context, idemKey, poID, poCode string) (bool, error) {
	if err := r.q.InsertPurchaseOrderIdemKey(ctx, appdb.InsertPurchaseOrderIdemKeyParams{
		IdemKey: idemKey,
		PoID:    poID,
		PoCode:  poCode,
	}); err != nil {
		return false, err
	}
	row, err := r.q.FindPurchaseOrderByIdemKey(ctx, idemKey)
	if err != nil {
		return false, err
	}
	return row.PoID == poID, nil
}

// ListPOs implement app.ListReader: trang PO (keyset created_at DESC, id DESC). Trả
// thêm created_at của dòng cuối để app sinh cursor (không lộ pgtype ra app).
func (r *Repo) ListPOs(ctx context.Context, f app.POFilter) ([]domain.PurchaseOrder, time.Time, error) {
	params := appdb.ListPurchaseOrdersParams{
		SupplierID: f.SupplierID,
		HasCursor:  f.HasCursor,
		AfterID:    f.AfterID,
		RowLimit:   f.Limit,
	}
	if f.Status != "" {
		s := f.Status
		params.Status = &s
	}
	if f.HasCursor {
		params.AfterCreatedAt = pgtype.Timestamptz{Time: f.AfterCreatedAt, Valid: true}
	}
	rows, err := r.q.ListPurchaseOrders(ctx, params)
	if err != nil {
		return nil, time.Time{}, err
	}
	out := make([]domain.PurchaseOrder, 0, len(rows))
	var last time.Time
	for _, row := range rows {
		out = append(out, domain.PurchaseOrder{
			ID:           row.ID,
			Code:         row.Code,
			SupplierID:   row.SupplierID,
			SupplierName: derefStr(row.SupplierName),
			Status:       row.Status,
			TotalAmount:  numericToMoney(row.TotalAmount),
			VATAmount:    numericToMoney(row.VatAmount),
			Note:         derefStr(row.Note),
			LockVersion:  row.LockVersion,
		})
		last = row.CreatedAt.Time
	}
	return out, last, nil
}

// Đảm bảo Repo thoả port app ở compile-time.
var (
	_ app.Store      = (*Repo)(nil)
	_ app.ListReader = (*Repo)(nil)
)

// ---- mapping helpers (tiền/lượng NUMERIC <-> decimal — KHÔNG float) ----

func numericToMoney(n pgtype.Numeric) money.Money {
	return money.FromDecimal(numericToDecimal(n))
}

func numericToDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid || n.NaN || n.Int == nil {
		return decimal.Zero
	}
	return decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
}

func moneyToNumeric(m money.Money) pgtype.Numeric {
	d := m.Decimal()
	return pgtype.Numeric{Int: new(big.Int).Set(d.Coefficient()), Exp: d.Exponent(), Valid: true}
}

func decimalToNumeric(d decimal.Decimal) pgtype.Numeric {
	return pgtype.Numeric{Int: new(big.Int).Set(d.Coefficient()), Exp: d.Exponent(), Valid: true}
}

func datePtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func datePtrOut(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optStr(s string) *string { return strPtr(s) }

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
