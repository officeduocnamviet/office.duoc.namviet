package idempotency

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// pgxStore là một Store dùng *pgxpool.Pool với các truy vấn sinh bởi sqlc
// (package appdb) trên bảng app.idempotency_keys. SQL-first, kiểm compile-time,
// không raw string trong code Go.
type pgxStore struct {
	q   *appdb.Queries
	ctx context.Context
}

// NewPgxStore tạo Store dựa trên pool. ctx dùng cho mọi query (gắn deadline ở
// tầng gọi nếu cần).
func NewPgxStore(ctx context.Context, pool *pgxpool.Pool) Store {
	return &pgxStore{q: appdb.New(pool), ctx: ctx}
}

func (s *pgxStore) Get(key string) (Record, bool, error) {
	row, err := s.q.GetIdempotencyKey(s.ctx, key)
	if errors.Is(err, pgx.ErrNoRows) {
		return Record{}, false, nil
	}
	if err != nil {
		return Record{}, false, err
	}
	rec := Record{State: row.State, Body: row.ResponseBody}
	if row.ResponseStatus != nil {
		rec.Status = int(*row.ResponseStatus)
	}
	return rec, true, nil
}

func (s *pgxStore) Begin(key, hash string) error {
	// ON CONFLICT (key) DO NOTHING nằm trong query sqlc — chống double-POST.
	return s.q.InsertIdempotencyKey(s.ctx, appdb.InsertIdempotencyKeyParams{
		Key:         key,
		RequestHash: hash,
	})
}

func (s *pgxStore) Complete(key string, status int, body []byte) error {
	st := int32(status)
	return s.q.CompleteIdempotencyKey(s.ctx, appdb.CompleteIdempotencyKeyParams{
		Key:            key,
		ResponseStatus: &st,
		ResponseBody:   body,
	})
}
