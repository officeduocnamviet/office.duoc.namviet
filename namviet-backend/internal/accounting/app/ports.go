// Package app là tầng use-case của accounting: điều phối POST bút toán (gate kỳ
// mở + allow_posting + cân Σ) và đọc sổ. Mở/commit transaction ở đây (hoặc nhận
// tx từ caller để gộp atomic với nghiệp vụ orders/finance). Domain không thấy tx.
package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
)

// EntryStore là PORT GHI sổ bound tới MỘT transaction (do caller hoặc TxManager
// truyền vào). Adapter postgres implement bằng appdb.Queries.WithTx(tx). Đây là
// port ở TẦNG APP (không domain) vì gắn với điều phối transaction — domain THUẦN
// không được biết pgx/tx (arch_test chặn).
type EntryStore interface {
	// OpenPeriodOf trả id kỳ kế toán đang 'open' chứa ngày date (theo year/month).
	// Trả (id, true, nil) nếu có kỳ mở; ("", false, nil) nếu KHÔNG có kỳ mở (chưa
	// mở hoặc đã khoá) — service map thành Conflict. Lỗi hạ tầng → (_, _, err).
	OpenPeriodOf(ctx context.Context, date string) (periodID string, open bool, err error)
	// AccountPostable trả true nếu tài khoản account_code tồn tại và allow_posting.
	// account_code không tồn tại → (false, nil) (service map Unprocessable).
	AccountPostable(ctx context.Context, accountCode string) (bool, error)
	// InsertEntry ghi header bút toán, trả id (uuid sinh app-side). periodID là kỳ
	// đã gate. entryID rỗng để adapter tự sinh.
	InsertEntry(ctx context.Context, e domain.JournalEntry, periodID string) (entryID string, err error)
	// InsertLine ghi một dòng bút toán (line_no theo thứ tự). Constraint trigger
	// DEFERRABLE cân Σ lúc commit — nếu lệch, commit của tx sẽ FAIL.
	InsertLine(ctx context.Context, entryID string, lineNo int32, l domain.EntryLine) error
}

// EntryStoreFromTx dựng một EntryStore bound tới tx. TxManager / Poster dùng để
// lấy store cho transaction hiện hành. Tách thành func để app không phụ thuộc
// cứng vào cách khởi tạo repo (adapter cung cấp).
type EntryStoreFromTx func(tx pgx.Tx) EntryStore

// TxManager mở/commit một transaction cho trường hợp post ĐỘC LẬP (PostInOwnTx).
// Adapter implement bằng platform/db.WithinTx. Khi orders/finance post bút toán
// trong tx nghiệp vụ của HỌ, chúng gọi thẳng Poster.Post(ctx, tx, entry) (KHÔNG
// qua TxManager) để gộp atomic.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

// ReadStore là PORT ĐỌC sổ (bind pool) cho HTTP. Tách khỏi EntryStore (ghi/tx)
// để đường đọc không cần transaction.
type ReadStore interface {
	// ListEntries trả một trang bút toán theo keyset (created_at DESC, id DESC),
	// lọc optional book/periodID. Hết → slice rỗng.
	ListEntries(ctx context.Context, f domain.EntryFilter) ([]domain.JournalEntryRecord, error)
	// GetEntry trả một bút toán theo id KÈM các dòng. Không thấy → apperr.NotFound.
	GetEntry(ctx context.Context, id string) (domain.JournalEntryRecord, error)
}
