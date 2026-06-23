package postgres

import (
	"context"

	"github.com/Maneva-AI/namviet-backend/internal/finance/app"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// ReadRepo implement app.PaymentReader (bind pool). Đường đọc thuần, không tx.
type ReadRepo struct{ q *appdb.Queries }

// NewReadRepo tạo repo đọc từ một *appdb.Queries (đã bind pool).
func NewReadRepo(q *appdb.Queries) *ReadRepo { return &ReadRepo{q: q} }

// ListByRef trả các phiếu còn sống trỏ về (ref_type, ref_id), mới nhất trước.
func (r *ReadRepo) ListByRef(ctx context.Context, refType, refID string, limit int32) ([]domain.Payment, error) {
	rows, err := r.q.ListPaymentsByRef(ctx, appdb.ListPaymentsByRefParams{
		RefType:  refType,
		RefID:    refID,
		RowLimit: limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Payment, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Payment{
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
		})
	}
	return out, nil
}

// Đảm bảo ReadRepo thoả port app ở compile-time.
var _ app.PaymentReader = (*ReadRepo)(nil)
