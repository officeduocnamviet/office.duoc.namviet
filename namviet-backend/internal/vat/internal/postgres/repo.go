// Package postgres là ADAPTER ra phía cơ sở dữ liệu của vat: implement các port
// app (InvoiceStore ghi/tx + cấp số gapless, InvoiceReader đọc) bằng query sinh
// từ sqlc (appdb) và map row <-> entity domain. GHI object MỚI ở schema app
// (sales_invoices/_lines/invoice_serials). Nằm dưới internal/ nên module khác
// KHÔNG import được. Tiền = NUMERIC <-> common/money decimal, KHÔNG float
// (numericToMoney/moneyToNumeric dựng từ big.Int). vat_rate = NUMERIC(6,4) <->
// decimal.Decimal. PK uuid sinh app-side (common/id v7).
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
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
	"github.com/Maneva-AI/namviet-backend/internal/vat/app"
	"github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

// pgUniqueViolation là SQLSTATE 23505 (unique_violation) — HĐ trùng order_code
// (unique index một phần WHERE status='issued') khi 2 tx đua phát hành cùng đơn.
const pgUniqueViolation = "23505"

// InvoiceRepo implement app.InvoiceStore trên appdb.Queries (bind tx do caller/
// TxManager truyền). Mọi thao tác GHI + cấp số gapless chạy trong tx này (gộp
// atomic với nghiệp vụ giao hàng của orders).
type InvoiceRepo struct{ q *appdb.Queries }

// NewInvoiceRepo tạo repo ghi từ một *appdb.Queries (đã bind tx).
func NewInvoiceRepo(q *appdb.Queries) *InvoiceRepo { return &InvoiceRepo{q: q} }

// FindIssuedByOrder tìm HĐ đã 'issued' cho order_code. Không thấy → (nil, nil).
func (r *InvoiceRepo) FindIssuedByOrder(ctx context.Context, orderCode string) (*domain.IssuedInvoice, error) {
	row, err := r.q.FindIssuedInvoiceByOrder(ctx, orderCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	inv := headerToIssued(row)
	return &inv, nil
}

// NextInvoiceNo cấp số HĐ GAPLESS: EnsureSerial (idempotent) bảo đảm dòng serial
// tồn tại, rồi NextInvoiceNo khoá dòng (FOR UPDATE qua UPDATE ... RETURNING) +
// tăng next_no — atomic trong tx hiện hành. UPDATE giữ row-lock tới hết tx nên
// các tx cùng serial tuần tự hoá → số liên tục, không trùng/không nhảy.
func (r *InvoiceRepo) NextInvoiceNo(ctx context.Context, serial, mauSo string) (int64, error) {
	if err := r.q.EnsureSerial(ctx, appdb.EnsureSerialParams{Serial: serial, MauSo: strNarg(mauSo)}); err != nil {
		return 0, err
	}
	return r.q.NextInvoiceNo(ctx, serial)
}

// InsertInvoice ghi header HĐ (id sinh app-side uuid v7, status 'issued'). Trùng
// do race trên UNIQUE (order_code WHERE issued) hoặc (serial,invoice_no) → 23505
// → (_, true, nil) để service đọc lại HĐ cũ (idempotent).
func (r *InvoiceRepo) InsertInvoice(ctx context.Context, inv domain.Invoice, invoiceNo int64) (string, bool, error) {
	invoiceID := id.NewString()
	got, err := r.q.InsertSalesInvoice(ctx, appdb.InsertSalesInvoiceParams{
		ID:              invoiceID,
		OrderCode:       inv.OrderCode,
		CustomerTaxCode: inv.CustomerTaxCode,
		Serial:          inv.Serial,
		InvoiceNo:       invoiceNo,
		IssueDate:       dateToPg(inv.IssueDate),
		Subtotal:        moneyToNumeric(inv.Subtotal),
		VatAmount:       moneyToNumeric(inv.VATAmount),
		Total:           moneyToNumeric(inv.Total),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return "", true, nil // race: HĐ kia thắng → service đọc lại
		}
		return "", false, err
	}
	return got, false, nil
}

// InsertLine ghi một dòng HĐ (id sinh app-side). line_no/line_amount/line_vat đã
// tính ở domain.
func (r *InvoiceRepo) InsertLine(ctx context.Context, invoiceID string, l domain.InvoiceLine) error {
	return r.q.InsertSalesInvoiceLine(ctx, appdb.InsertSalesInvoiceLineParams{
		ID:          id.NewString(),
		InvoiceID:   invoiceID,
		LineNo:      l.LineNo,
		ProductID:   l.ProductID,
		Description: l.Description,
		Quantity:    moneyToNumeric(l.Quantity),
		UnitPrice:   moneyToNumeric(l.UnitPrice),
		VatRate:     decimalToNumeric(l.VATRate),
		LineAmount:  moneyToNumeric(l.LineAmount),
		LineVat:     moneyToNumeric(l.LineVAT),
	})
}

// GetInvoiceWithLines nạp HĐ theo id kèm lines (dùng trong tx ghi để trả về).
func (r *InvoiceRepo) GetInvoiceWithLines(ctx context.Context, invoiceID string) (domain.IssuedInvoice, error) {
	return getWithLines(ctx, r.q, invoiceID)
}

// ReadRepo implement app.InvoiceReader (bind pool). Đường đọc thuần, không tx.
type ReadRepo struct{ q *appdb.Queries }

// NewReadRepo tạo repo đọc từ một *appdb.Queries (đã bind pool).
func NewReadRepo(q *appdb.Queries) *ReadRepo { return &ReadRepo{q: q} }

func (r *ReadRepo) ListInvoices(ctx context.Context, f domain.InvoiceFilter) ([]domain.IssuedInvoice, error) {
	params := appdb.ListSalesInvoicesParams{
		RowLimit:  f.Limit,
		OrderCode: strNarg(f.OrderCode),
		Status:    strNarg(f.Status),
	}
	if f.HasCursor {
		params.AfterCreatedAt = tsToPg(f.AfterCreatedAt)
		aid := f.AfterID
		params.AfterID = &aid
	}
	rows, err := r.q.ListSalesInvoices(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]domain.IssuedInvoice, 0, len(rows))
	for _, row := range rows {
		out = append(out, headerToIssued(row))
	}
	return out, nil
}

func (r *ReadRepo) GetInvoice(ctx context.Context, invoiceID string) (domain.IssuedInvoice, error) {
	return getWithLines(ctx, r.q, invoiceID)
}

// getWithLines nạp header + lines của một HĐ (dùng chung đường ghi/đọc). Không
// thấy header → apperr.NotFound.
func getWithLines(ctx context.Context, q *appdb.Queries, invoiceID string) (domain.IssuedInvoice, error) {
	row, err := q.GetSalesInvoice(ctx, invoiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.IssuedInvoice{}, apperr.NotFound("hoá đơn không tồn tại")
		}
		return domain.IssuedInvoice{}, err
	}
	inv := headerToIssued(row)
	lines, err := q.ListSalesInvoiceLines(ctx, invoiceID)
	if err != nil {
		return domain.IssuedInvoice{}, err
	}
	inv.Lines = make([]domain.InvoiceLine, 0, len(lines))
	for _, l := range lines {
		inv.Lines = append(inv.Lines, domain.InvoiceLine{
			LineNo:      l.LineNo,
			ProductID:   l.ProductID,
			Description: l.Description,
			Quantity:    numericToMoney(l.Quantity),
			UnitPrice:   numericToMoney(l.UnitPrice),
			VATRate:     numericToDecimal(l.VatRate),
			LineAmount:  numericToMoney(l.LineAmount),
			LineVAT:     numericToMoney(l.LineVat),
		})
	}
	return inv, nil
}

// ---- mapping row <-> domain ----

// headerToIssued map một row header sales_invoices sang IssuedInvoice (không
// lines). Dùng cho cả FindIssuedByOrder/List/Get.
func headerToIssued(row appdb.AppSalesInvoice) domain.IssuedInvoice {
	return domain.IssuedInvoice{
		ID: row.ID,
		Invoice: domain.Invoice{
			OrderCode:       row.OrderCode,
			CustomerTaxCode: row.CustomerTaxCode,
			Serial:          row.Serial,
			IssueDate:       row.IssueDate.Time,
			Subtotal:        numericToMoney(row.Subtotal),
			VATAmount:       numericToMoney(row.VatAmount),
			Total:           numericToMoney(row.Total),
		},
		InvoiceNo: row.InvoiceNo,
		Status:    domain.Status(row.Status),
		CreatedAt: row.CreatedAt.Time,
	}
}

// numericToMoney chuyển pgtype.Numeric sang money.Money KHÔNG qua float: dựng
// decimal từ mantissa (big.Int) * 10^Exp. NULL/NaN → Zero.
func numericToMoney(n pgtype.Numeric) money.Money {
	return money.FromDecimal(numericToDecimal(n))
}

// numericToDecimal chuyển pgtype.Numeric sang decimal.Decimal KHÔNG qua float
// (dùng cho vat_rate — không bọc money vì là tỷ lệ, không phải tiền).
func numericToDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid || n.NaN || n.Int == nil {
		return decimal.Zero
	}
	return decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
}

// moneyToNumeric chuyển money.Money sang pgtype.Numeric KHÔNG qua float.
func moneyToNumeric(m money.Money) pgtype.Numeric {
	return decimalToNumeric(m.Decimal())
}

// decimalToNumeric chuyển decimal.Decimal sang pgtype.Numeric (coefficient +
// exponent, KHÔNG qua float). Dùng cho cả tiền và vat_rate.
func decimalToNumeric(d decimal.Decimal) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   new(big.Int).Set(d.Coefficient()),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

// dateToPg chuyển time.Time → pgtype.Date (chỉ phần ngày).
func dateToPg(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

// tsToPg chuyển time.Time → pgtype.Timestamptz.
func tsToPg(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// strNarg chuyển chuỗi rỗng → nil (không lọc) cho sqlc narg.
func strNarg(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Đảm bảo các repo thoả port app ở compile-time.
var (
	_ app.InvoiceStore  = (*InvoiceRepo)(nil)
	_ app.InvoiceReader = (*ReadRepo)(nil)
)
