package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
	"github.com/Maneva-AI/namviet-backend/internal/identity/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
)

// TestMarkUsed_AtomicClaim chứng minh CLAIM NGUYÊN TỬ của bước xoay vòng refresh
// token: MarkUsed lần đầu thành công (token chuyển used=true), lần hai trên CÙNG
// token phải trả apperr Kind=Conflict. Đây là nền tảng chống refresh-token
// rotation race: hai request đồng thời cùng một token chỉ một được "claim".
func TestMarkUsed_AtomicClaim(t *testing.T) {
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	ctx := context.Background()

	q := appdb.New(pool)

	// Seed một user (refresh_tokens.user_id là FK NOT NULL).
	userID, err := q.InsertUser(ctx, appdb.InsertUserParams{
		Email: "race@nv.vn", PasswordHash: "x", HashAlgo: "argon2id",
		UserType: "staff", IsActive: true,
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	repo := postgres.NewRefreshTokenRepo(q)
	tokenID := id.NewString()
	if err := repo.Insert(ctx, domain.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		TokenHash: "hash-atomic-claim",
		FamilyID:  id.NewString(),
		ExpiresAt: time.Now().Add(time.Hour),
	}); err != nil {
		t.Fatalf("insert token: %v", err)
	}

	// Lần 1: claim thành công.
	if err := repo.MarkUsed(ctx, tokenID); err != nil {
		t.Fatalf("MarkUsed lần 1 = %v, want nil", err)
	}

	// Lần 2: token đã used → mất race → Conflict.
	err = repo.MarkUsed(ctx, tokenID)
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("MarkUsed lần 2 kind = %v, want Conflict (err=%v)", apperr.KindOf(err), err)
	}
}
