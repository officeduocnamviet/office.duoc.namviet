package domain

import "context"

// Các PORT do domain ĐỊNH NGHĨA, adapter (postgres) implement
// ("accept interfaces, return structs"). Tất cả method nhận context.Context;
// việc gom nhiều thao tác vào MỘT transaction là trách nhiệm của tầng app
// (qua platform/db.WithinTx + TxManager ở app) — domain không biết tx/SQL.

// UserRepository truy xuất/ghi user.
type UserRepository interface {
	// GetByEmail trả user theo email (citext, không phân biệt hoa thường).
	// Không thấy → apperr.NotFound.
	GetByEmail(ctx context.Context, email string) (User, error)
	// GetByID trả user theo id. Không thấy → apperr.NotFound.
	GetByID(ctx context.Context, id string) (User, error)
	// UpdatePasswordHash cập nhật băm + thuật toán (dùng cho lazy rehash).
	UpdatePasswordHash(ctx context.Context, userID, hash string, algo HashAlgo) error
}

// RefreshTokenRepository quản lý vòng đời refresh token (opaque, xoay vòng).
// Các thao tác xoay vòng (GetByHash → MarkUsed → Insert) phải chạy trong CÙNG
// một transaction; tầng app bind repo này tới tx trước khi gọi.
type RefreshTokenRepository interface {
	// Insert lưu một refresh token mới (đã hash).
	Insert(ctx context.Context, t RefreshToken) error
	// GetByHash trả token theo hash. Không thấy → apperr.NotFound.
	GetByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	// MarkUsed CLAIM NGUYÊN TỬ token (theo id) là đã dùng — bước xoay vòng hợp
	// lệ. Chỉ thành công khi token còn used=false; nếu token đã được dùng (mất
	// race với một request refresh đồng thời, hoặc reuse) → trả apperr
	// Kind=Conflict. Caller (app.Refresh) xử lý Conflict GIỐNG reuse-detection:
	// thu hồi cả family. Chữ ký giữ nguyên — chỉ bổ sung ngữ nghĩa Conflict.
	MarkUsed(ctx context.Context, id string) error
	// RevokeFamily thu hồi toàn bộ token cùng family (reuse-detection / logout).
	RevokeFamily(ctx context.Context, familyID string) error
}

// RoleRepository tra cứu RBAC.
type RoleRepository interface {
	// PermissionCodesForUser trả danh sách mã quyền (permission.code) hiệu lực
	// của user qua user_roles → role_permissions → permissions. Không có quyền
	// → slice rỗng (không lỗi).
	PermissionCodesForUser(ctx context.Context, userID string) ([]string, error)
}
