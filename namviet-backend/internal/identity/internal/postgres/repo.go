// Package postgres là ADAPTER ra phía cơ sở dữ liệu của identity: implement các
// port domain bằng query sinh từ sqlc (appdb) và map giữa row <-> entity domain.
// Nằm dưới internal/ nên module khác KHÔNG import được (compiler chặn ranh giới).
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// UserRepo implement domain.UserRepository trên appdb.Queries (pool hoặc tx).
type UserRepo struct{ q *appdb.Queries }

// NewUserRepo tạo repo từ một *appdb.Queries (đã bind pool hoặc tx).
func NewUserRepo(q *appdb.Queries) *UserRepo { return &UserRepo{q: q} }

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, apperr.NotFound("user không tồn tại")
		}
		return domain.User{}, err
	}
	return userToDomain(row), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, apperr.NotFound("user không tồn tại")
		}
		return domain.User{}, err
	}
	return userToDomain(row), nil
}

func (r *UserRepo) UpdatePasswordHash(ctx context.Context, userID, hash string, algo domain.HashAlgo) error {
	return r.q.UpdateUserPasswordHash(ctx, appdb.UpdateUserPasswordHashParams{
		ID:           userID,
		PasswordHash: hash,
		HashAlgo:     string(algo),
	})
}

// RefreshTokenRepo implement domain.RefreshTokenRepository.
type RefreshTokenRepo struct{ q *appdb.Queries }

// NewRefreshTokenRepo tạo repo từ một *appdb.Queries.
func NewRefreshTokenRepo(q *appdb.Queries) *RefreshTokenRepo { return &RefreshTokenRepo{q: q} }

func (r *RefreshTokenRepo) Insert(ctx context.Context, t domain.RefreshToken) error {
	return r.q.InsertRefreshToken(ctx, appdb.InsertRefreshTokenParams{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenHash: t.TokenHash,
		FamilyID:  t.FamilyID,
		ExpiresAt: toTimestamptz(t.ExpiresAt),
	})
}

func (r *RefreshTokenRepo) GetByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	row, err := r.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, apperr.NotFound("refresh token không tồn tại")
		}
		return domain.RefreshToken{}, err
	}
	return tokenToDomain(row), nil
}

func (r *RefreshTokenRepo) MarkUsed(ctx context.Context, id string) error {
	// Claim nguyên tử: UPDATE ... WHERE id=$1 AND used=false. 0 dòng ảnh hưởng
	// nghĩa là token đã được dùng (mất race với một request refresh đồng thời,
	// hoặc reuse) → trả Conflict để caller xử lý như reuse-detection.
	rows, err := r.q.MarkRefreshTokenUsed(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperr.Conflict("refresh token đã được dùng")
	}
	return nil
}

func (r *RefreshTokenRepo) RevokeFamily(ctx context.Context, familyID string) error {
	return r.q.RevokeRefreshTokenFamily(ctx, familyID)
}

// RoleRepo implement domain.RoleRepository.
type RoleRepo struct{ q *appdb.Queries }

// NewRoleRepo tạo repo từ một *appdb.Queries.
func NewRoleRepo(q *appdb.Queries) *RoleRepo { return &RoleRepo{q: q} }

func (r *RoleRepo) PermissionCodesForUser(ctx context.Context, userID string) ([]string, error) {
	codes, err := r.q.PermissionCodesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if codes == nil {
		return []string{}, nil
	}
	return codes, nil
}

// ---- mapping row <-> domain ----

func userToDomain(r appdb.AppUser) domain.User {
	return domain.User{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		HashAlgo:     domain.HashAlgo(r.HashAlgo),
		UserType:     r.UserType,
		IsActive:     r.IsActive,
		CreatedAt:    fromTimestamptz(r.CreatedAt),
		UpdatedAt:    fromTimestamptz(r.UpdatedAt),
	}
}

func tokenToDomain(r appdb.AppRefreshToken) domain.RefreshToken {
	return domain.RefreshToken{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		FamilyID:  r.FamilyID,
		Used:      r.Used,
		Revoked:   r.Revoked,
		ExpiresAt: fromTimestamptz(r.ExpiresAt),
		CreatedAt: fromTimestamptz(r.CreatedAt),
	}
}

func toTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func fromTimestamptz(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Đảm bảo các repo thoả port domain ở compile-time.
var (
	_ domain.UserRepository         = (*UserRepo)(nil)
	_ domain.RefreshTokenRepository = (*RefreshTokenRepo)(nil)
	_ domain.RoleRepository         = (*RoleRepo)(nil)
)
