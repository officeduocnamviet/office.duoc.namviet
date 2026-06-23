package app

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

// ListReader là PORT ĐỌC danh sách PO (keyset created_at DESC, id DESC). Bound-tx
// như Store (đọc trong tx ngắn) để dùng chung adapter.
type ListReader interface {
	ListPOs(ctx context.Context, f POFilter) ([]domain.PurchaseOrder, time.Time, error)
}

// POFilter là điều kiện lọc + keyset cho ListPOs.
type POFilter struct {
	Status         string
	SupplierID     *int64
	Limit          int32
	HasCursor      bool
	AfterCreatedAt time.Time
	AfterID        string
}

// ListPOsQuery là input đọc danh sách PO đã giải mã ở edge.
type ListPOsQuery struct {
	Cursor     string
	Limit      int32
	Status     string
	SupplierID *int64
}

// ListPOsResult là một trang PO + cursor trang kế (rỗng nếu hết).
type ListPOsResult struct {
	Items      []domain.PurchaseOrder
	NextCursor string
}

// ListReaderFromTx dựng ListReader bound tới tx (adapter cấp). Tách khỏi Store để
// đường đọc danh sách độc lập, nhưng cùng adapter.
type ListReaderFromTx func(tx pgx.Tx) ListReader

// ListPOs trả một trang PO (keyset created_at DESC, id DESC). Tự decode cursor,
// chuẩn hoá limit, sinh NextCursor nếu trang đầy.
func (s *Service) ListPOs(ctx context.Context, q ListPOsQuery) (ListPOsResult, error) {
	if s.txm == nil || s.listFromTx == nil {
		return ListPOsResult{}, apperr.Internal("Service chưa cấu hình đủ cho ListPOs")
	}
	nano, afterID, derr := decodeCursor(q.Cursor)
	if derr != nil {
		return ListPOsResult{}, apperr.Validation("cursor không hợp lệ")
	}
	limit := normalizeLimit(q.Limit)
	f := POFilter{Status: q.Status, SupplierID: q.SupplierID, Limit: limit}
	if q.Cursor != "" {
		f.AfterCreatedAt = time.Unix(0, nano).UTC()
		f.AfterID = afterID
		f.HasCursor = true
	}

	var items []domain.PurchaseOrder
	var lastCreatedAt time.Time
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		got, last, gerr := s.listFromTx(tx).ListPOs(ctx, f)
		if gerr != nil {
			return apperr.Internal("đọc danh sách PO lỗi").WithCause(gerr)
		}
		items = got
		lastCreatedAt = last
		return nil
	})
	if err != nil {
		return ListPOsResult{}, err
	}
	res := ListPOsResult{Items: items}
	if int32(len(items)) == limit && limit > 0 {
		last := items[len(items)-1]
		res.NextCursor = encodeCursor(lastCreatedAt.UnixNano(), last.ID)
	}
	return res, nil
}

func normalizeLimit(l int32) int32 {
	switch {
	case l <= 0:
		return defaultLimit
	case l > maxLimit:
		return maxLimit
	default:
		return l
	}
}

// encodeCursor mã hoá (createdAtNano, id) thành cursor opaque (base64) — keyset
// (created_at DESC, id DESC) ổn định kể cả khi nhiều PO cùng mốc thời gian.
func encodeCursor(createdAtNano int64, id string) string {
	raw := strconv.FormatInt(createdAtNano, 10) + "|" + id
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// decodeCursor giải mã cursor về (createdAtNano, id). Rỗng = trang đầu → (0, "").
// Sai định dạng → lỗi (handler trả 422, KHÔNG im lặng coi trang đầu).
func decodeCursor(cursor string) (int64, string, error) {
	if cursor == "" {
		return 0, "", nil
	}
	rawBytes, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, "", err
	}
	raw := string(rawBytes)
	sep := strings.IndexByte(raw, '|')
	if sep < 0 {
		return 0, "", errors.New("cursor thiếu dấu phân tách")
	}
	nano, err := strconv.ParseInt(raw[:sep], 10, 64)
	if err != nil {
		return 0, "", err
	}
	idPart := raw[sep+1:]
	if idPart == "" {
		return 0, "", errors.New("cursor thiếu id")
	}
	return nano, idPart, nil
}
