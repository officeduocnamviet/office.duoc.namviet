package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Beginner là cổng tối thiểu để mở transaction. *pgxpool.Pool thoả interface
// này; tách interface giúp test/app không phụ thuộc cứng vào pool cụ thể.
type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// WithinTx chạy fn trong MỘT transaction: begin → fn → commit; nếu fn trả lỗi
// hoặc panic thì rollback. Đây là điểm DUY NHẤT mở/commit transaction (theo
// ARCHITECTURE.md §6): tầng app gọi WithinTx và truyền pgx.Tx xuống repo
// (qua Queries.WithTx), domain không bao giờ thấy tx.
//
// Rollback dùng context.Background() có chủ đích: khi ctx gốc đã bị huỷ
// (deadline/cancel) ta vẫn cần rollback để không rò connection.
func WithinTx(ctx context.Context, b Beginner, fn func(tx pgx.Tx) error) (err error) {
	tx, err := b.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(context.Background())
			panic(p) // tái phát panic sau khi đã rollback
		}
		if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// Đảm bảo *pgxpool.Pool thoả Beginner ở compile-time.
var _ Beginner = (*pgxpool.Pool)(nil)
