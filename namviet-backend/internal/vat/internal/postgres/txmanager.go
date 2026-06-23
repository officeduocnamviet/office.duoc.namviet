package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/platform/db"
	"github.com/Maneva-AI/namviet-backend/internal/vat/app"
)

// TxManager implement app.TxManager: mở một transaction qua platform/db.WithinTx
// rồi cấp pgx.Tx cho closure. Dùng cho Service.IssueInvoiceInOwnTx (phát hành HĐ
// ĐỘC LẬP). Khi orders phát hành trong tx giao hàng của họ, chúng gọi thẳng
// IssuePort.IssueInvoice(ctx, tx, ...) (KHÔNG qua đây) để gộp atomic.
type TxManager struct {
	beginner db.Beginner
}

// NewTxManager tạo TxManager từ pool runtime (db.Beginner).
func NewTxManager(beginner db.Beginner) *TxManager {
	return &TxManager{beginner: beginner}
}

// WithinTx chạy fn trong một transaction; commit nếu fn trả nil, rollback nếu lỗi.
func (m *TxManager) WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return db.WithinTx(ctx, m.beginner, fn)
}

var _ app.TxManager = (*TxManager)(nil)
