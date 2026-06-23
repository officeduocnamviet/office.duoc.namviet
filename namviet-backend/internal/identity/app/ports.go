package app

import (
	"context"

	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
)

// Repos gom các repository của identity bound tới CÙNG một transaction. Tầng
// app dùng nó bên trong TxManager.WithinTx để các thao tác (đọc user, xoay
// refresh token, cập nhật hash) chạy nguyên tử trong một tx — domain không
// thấy tx, app điều phối transaction (ARCHITECTURE.md §6).
type Repos struct {
	Users  domain.UserRepository
	Tokens domain.RefreshTokenRepository
	Roles  domain.RoleRepository
}

// TxManager mở/commit transaction và cấp bộ Repos bound tới tx cho closure.
// Adapter postgres implement bằng platform/db.WithinTx + Queries.WithTx(tx).
// Đây là port ở TẦNG APP (không phải domain) vì nó thuộc về điều phối use-case.
type TxManager interface {
	// WithinTx chạy fn trong một transaction; commit nếu fn trả nil, rollback
	// nếu lỗi/panic.
	WithinTx(ctx context.Context, fn func(r Repos) error) error
}
