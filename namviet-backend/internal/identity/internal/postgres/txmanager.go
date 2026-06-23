package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db"
)

// TxManager implement app.TxManager: mở một transaction qua platform/db.WithinTx
// rồi cấp một bộ Repos bind tới CÙNG tx (qua appdb.Queries.WithTx). Đây là chỗ
// DUY NHẤT của identity mở/commit transaction; app chỉ điều phối closure.
type TxManager struct {
	pool *appdbPool
}

// appdbPool gom *pgxpool.Pool (đã có db.Beginner) + factory tạo Queries từ tx.
// Tách nhỏ để dễ test, nhưng thực tế chỉ cần pool.
type appdbPool struct {
	beginner db.Beginner
	newQ     func(pgx.Tx) reposFactory
}

// reposFactory dựng 3 repo từ một *appdb.Queries (bind tx). Định nghĩa như một
// func để txmanager không phụ thuộc cứng vào việc khởi tạo repo.
type reposFactory func() app.Repos

// NewTxManager tạo TxManager từ pool runtime. Mỗi WithinTx sẽ bind Queries vào
// tx hiện hành.
func NewTxManager(pool db.Beginner, queriesFromTx func(pgx.Tx) app.Repos) *TxManager {
	return &TxManager{pool: &appdbPool{
		beginner: pool,
		newQ: func(tx pgx.Tx) reposFactory {
			return func() app.Repos { return queriesFromTx(tx) }
		},
	}}
}

// WithinTx chạy fn trong một transaction, cấp Repos bound tới tx.
func (m *TxManager) WithinTx(ctx context.Context, fn func(r app.Repos) error) error {
	return db.WithinTx(ctx, m.pool.beginner, func(tx pgx.Tx) error {
		repos := m.pool.newQ(tx)()
		return fn(repos)
	})
}

// Đảm bảo thoả port app ở compile-time.
var _ app.TxManager = (*TxManager)(nil)
