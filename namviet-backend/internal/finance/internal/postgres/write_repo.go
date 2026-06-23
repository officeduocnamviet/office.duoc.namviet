// Package postgres là ADAPTER ra phía cơ sở dữ liệu của finance: implement port
// app (PaymentWriter ghi/tx, PaymentReader đọc) bằng query sinh từ sqlc (appdb) và
// map row <-> entity domain. GHI/ĐỌC bảng public.finance_transactions kế thừa
// (strangler, ADR 0001). Nằm dưới internal/ nên module khác KHÔNG import được.
// Tiền amount = NUMERIC <-> common/money decimal, KHÔNG float (numericToMoney/
// moneyToNumeric dựng từ big.Int).
//
// ⚠️ KHÔNG có thao tác UPDATE fund_accounts.balance — trigger prod tự cộng số dư
// khi phiếu sang status='completed' (xem db/schema/public_finance.sql). Ghi 2 lần
// = cộng đôi.
package postgres

import (
	"context"
	"errors"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/finance/app"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// pgUniqueViolation là SQLSTATE 23505 (unique_violation) — phiếu trùng code/
// bank_ref khi 2 luồng cùng idem key đua insert.
const pgUniqueViolation = "23505"

// WriteRepo implement app.PaymentWriter trên appdb.Queries (bind tx do caller/
// TxManager truyền). Mọi thao tác chạy trong CÙNG tx — gộp atomic với nghiệp vụ
// orders/POS.
type WriteRepo struct{ q *appdb.Queries }

// NewWriteRepo tạo repo ghi từ một *appdb.Queries (đã bind tx).
func NewWriteRepo(q *appdb.Queries) *WriteRepo { return &WriteRepo{q: q} }

// FindAliveByBankRef tìm phiếu còn sống theo bank_reference_id. Không thấy → (nil, nil).
func (r *WriteRepo) FindAliveByBankRef(ctx context.Context, bankRef string) (*domain.Payment, error) {
	row, err := r.q.FindAlivePaymentByBankRef(ctx, bankRef)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	p := bankRefRowToPayment(row)
	return &p, nil
}

// FindAliveByCode tìm phiếu còn sống theo code. Không thấy → (nil, nil).
func (r *WriteRepo) FindAliveByCode(ctx context.Context, code string) (*domain.Payment, error) {
	row, err := r.q.FindAlivePaymentByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	p := codeRowToPayment(row)
	return &p, nil
}

// InsertPaymentIn ghi phiếu THU (flow='in', status='completed'). Unique violation
// (trùng code/bank_ref do race) → (nil, true, nil) để service re-SELECT phiếu cũ.
func (r *WriteRepo) InsertPaymentIn(ctx context.Context, code string, p domain.RecordPaymentIn) (*domain.Payment, bool, error) {
	row, err := r.q.InsertPaymentIn(ctx, appdb.InsertPaymentInParams{
		Code:            code,
		BusinessType:    domain.BusinessTypeSale,
		Amount:          moneyToNumeric(p.Amount),
		FundAccountID:   p.FundAccountID,
		RefType:         p.EffectiveRefType(),
		RefID:           p.EffectiveRefID(),
		Description:     p.Description,
		Status:          p.EffectiveStatus(),
		BankReferenceID: trimptr(p.BankRef),
		BookType:        p.BookType.String(),
		CreatedBy:       p.CreatedBy,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, true, nil // race: phiếu kia thắng → service đọc lại
		}
		return nil, false, err
	}
	out := insertRowToPayment(row)
	return &out, false, nil
}

// InsertPaymentOut ghi phiếu CHI (flow='out', trả NCC mua hàng — mục 54). Unique
// violation (trùng code/bank_ref do race) → (nil, true, nil) để service re-SELECT.
func (r *WriteRepo) InsertPaymentOut(ctx context.Context, code string, p domain.RecordPaymentOut) (*domain.Payment, bool, error) {
	row, err := r.q.InsertPaymentOut(ctx, appdb.InsertPaymentOutParams{
		Code:            code,
		BusinessType:    domain.BusinessTypePurchase,
		Amount:          moneyToNumeric(p.Amount),
		FundAccountID:   p.FundAccountID,
		RefType:         p.EffectiveRefType(),
		RefID:           p.EffectiveRefID(),
		Description:     p.Description,
		Status:          p.EffectiveStatus(),
		BankReferenceID: trimptr(p.BankRef),
		BookType:        p.BookType.String(),
		CreatedBy:       p.CreatedBy,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, true, nil // race: phiếu kia thắng → service đọc lại
		}
		return nil, false, err
	}
	out := insertOutRowToPayment(row)
	return &out, false, nil
}

// ConfirmReceipt chuyển phiếu THU 'pending' → 'completed' (thủ quỹ xác nhận vào
// quỹ). Trả số dòng đổi (1 = vừa xác nhận; 0 = đã completed/không phải pending →
// idempotent). Trigger prod cộng số dư quỹ tại transition này.
func (r *WriteRepo) ConfirmReceipt(ctx context.Context, paymentID int64) (int64, error) {
	return r.q.ConfirmPaymentReceipt(ctx, paymentID)
}

// ---- mapping row -> domain ----

func insertRowToPayment(row appdb.InsertPaymentInRow) domain.Payment {
	return domain.Payment{
		ID:            row.ID,
		Code:          row.Code,
		Flow:          row.Flow,
		BusinessType:  row.BusinessType,
		Amount:        numericToMoney(row.Amount),
		FundAccountID: row.FundAccountID,
		RefType:       deref(row.RefType),
		RefID:         deref(row.RefID),
		Status:        row.Status,
		BookType:      domain.BookType(row.BookType),
		BankRef:       row.BankReferenceID,
		Description:   row.Description,
		CreatedBy:     row.CreatedBy,
	}
}

func insertOutRowToPayment(row appdb.InsertPaymentOutRow) domain.Payment {
	return domain.Payment{
		ID:            row.ID,
		Code:          row.Code,
		Flow:          row.Flow,
		BusinessType:  row.BusinessType,
		Amount:        numericToMoney(row.Amount),
		FundAccountID: row.FundAccountID,
		RefType:       deref(row.RefType),
		RefID:         deref(row.RefID),
		Status:        row.Status,
		BookType:      domain.BookType(row.BookType),
		BankRef:       row.BankReferenceID,
		Description:   row.Description,
		CreatedBy:     row.CreatedBy,
	}
}

func codeRowToPayment(row appdb.FindAlivePaymentByCodeRow) domain.Payment {
	return domain.Payment{
		ID:            row.ID,
		Code:          row.Code,
		Flow:          row.Flow,
		BusinessType:  row.BusinessType,
		Amount:        numericToMoney(row.Amount),
		FundAccountID: row.FundAccountID,
		RefType:       deref(row.RefType),
		RefID:         deref(row.RefID),
		Status:        row.Status,
		BookType:      domain.BookType(row.BookType),
		BankRef:       row.BankReferenceID,
		Description:   row.Description,
		CreatedBy:     row.CreatedBy,
	}
}

func bankRefRowToPayment(row appdb.FindAlivePaymentByBankRefRow) domain.Payment {
	return domain.Payment{
		ID:            row.ID,
		Code:          row.Code,
		Flow:          row.Flow,
		BusinessType:  row.BusinessType,
		Amount:        numericToMoney(row.Amount),
		FundAccountID: row.FundAccountID,
		RefType:       deref(row.RefType),
		RefID:         deref(row.RefID),
		Status:        row.Status,
		BookType:      domain.BookType(row.BookType),
		BankRef:       row.BankReferenceID,
		Description:   row.Description,
		CreatedBy:     row.CreatedBy,
	}
}

// numericToMoney chuyển pgtype.Numeric sang money.Money KHÔNG qua float: dựng
// decimal từ mantissa (big.Int) * 10^Exp. NULL/NaN → Zero.
func numericToMoney(n pgtype.Numeric) money.Money {
	if !n.Valid || n.NaN || n.Int == nil {
		return money.Zero()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return money.FromDecimal(d)
}

// moneyToNumeric chuyển money.Money sang pgtype.Numeric KHÔNG qua float: lấy
// coefficient (big.Int) + exponent từ decimal. Dùng khi GHI amount.
func moneyToNumeric(m money.Money) pgtype.Numeric {
	d := m.Decimal()
	return pgtype.Numeric{
		Int:   new(big.Int).Set(d.Coefficient()),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

// deref trả giá trị chuỗi từ con trỏ (nil → rỗng) — RefType/RefID nullable ở DB.
func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// trimptr trả con trỏ chuỗi đã trim (nil/rỗng/khoảng-trắng → nil) cho bank_ref.
func trimptr(s *string) *string {
	if s == nil {
		return nil
	}
	t := *s
	if t == "" {
		return nil
	}
	return s
}

// Đảm bảo WriteRepo thoả port app ở compile-time.
var _ app.PaymentWriter = (*WriteRepo)(nil)
