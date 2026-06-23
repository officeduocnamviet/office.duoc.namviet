package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db"
)

// TxManager implement app.TxManager: mở một transaction qua platform/db.WithinTx
// rồi cấp pgx.Tx cho closure. Dùng cho Deductor.DeductFEFOInOwnTx (trừ kho ĐỘC
// LẬP). Khi orders/POS trừ kho trong tx nghiệp vụ của họ, chúng gọi thẳng
// Deductor.DeductFEFO(ctx, tx, ...) (KHÔNG qua đây) để gộp atomic.
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
