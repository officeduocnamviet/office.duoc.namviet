package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

// Service là use-case của accounting. Nó implement port Poster (Post bút toán
// trong tx của caller) + cung cấp use-case đọc sổ. Đường POST KHÔNG có REST công
// khai ở P0 — chỉ module orders/finance/vat gọi Post(ctx, tx, entry) trong tx
// nghiệp vụ của họ để gộp atomic (sổ luôn khớp sự kiện).
type Service struct {
	storeFromTx EntryStoreFromTx
	txm         TxManager
	read        ReadStore
}

// New dựng Service. storeFromTx bind EntryStore tới một tx; txm để PostInOwnTx mở
// tx riêng; read là đường đọc (bind pool). Có thể truyền nil cho thành phần không
// dùng trong ngữ cảnh cụ thể (vd test).
func New(storeFromTx EntryStoreFromTx, txm TxManager, read ReadStore) *Service {
	return &Service{storeFromTx: storeFromTx, txm: txm, read: read}
}

// Post ghi một bút toán kép vào sổ TRONG transaction tx do CALLER truyền (để gộp
// atomic với nghiệp vụ orders/finance). Quy trình gate:
//  1. Validate() thuần (book hợp lệ, mỗi dòng đúng 1 vế > 0, Σdebit=Σcredit) →
//     lệch/sai = Unprocessable (Validation).
//  2. Kỳ kế toán của entry_date phải đang 'open' → khoá = Conflict.
//  3. Mọi account_code phải allow_posting = true → không = Unprocessable.
//  4. InsertEntry + InsertLine theo line_no (cùng tx). Trigger DEFERRABLE cân Σ
//     phòng thủ lúc commit (lệch lọt qua => commit FAIL).
//
// Trả entryID (uuid). Lỗi đã là apperr (map envelope ở tầng http).
func (s *Service) Post(ctx context.Context, tx pgx.Tx, e domain.JournalEntry) (string, error) {
	if err := e.Validate(); err != nil {
		// Bất biến cấu trúc/cân Σ sai về mặt nghiệp vụ → 422.
		return "", apperr.Validation(err.Error())
	}

	store := s.storeFromTx(tx)

	periodID, open, err := store.OpenPeriodOf(ctx, e.EntryDate.Format(dateLayout))
	if err != nil {
		return "", apperr.Internal("đọc kỳ kế toán lỗi").WithCause(err)
	}
	if !open {
		return "", apperr.Conflict("kỳ kế toán của ngày " + e.EntryDate.Format(dateLayout) + " chưa mở hoặc đã khoá")
	}

	for _, l := range e.Lines {
		ok, err := store.AccountPostable(ctx, l.AccountCode)
		if err != nil {
			return "", apperr.Internal("đọc cây tài khoản lỗi").WithCause(err)
		}
		if !ok {
			return "", apperr.Validation("tài khoản " + l.AccountCode + " không cho hạch toán trực tiếp (allow_posting=false hoặc không tồn tại)")
		}
	}

	entryID, err := store.InsertEntry(ctx, e, periodID)
	if err != nil {
		return "", apperr.Internal("ghi bút toán lỗi").WithCause(err)
	}
	for i, l := range e.Lines {
		if err := store.InsertLine(ctx, entryID, int32(i+1), l); err != nil {
			return "", apperr.Internal("ghi dòng bút toán lỗi").WithCause(err)
		}
	}
	return entryID, nil
}

// PostInOwnTx ghi một bút toán trong một transaction RIÊNG (mở/commit ở đây). Dùng
// khi post độc lập, không gộp với nghiệp vụ khác. Lỗi commit (vd trigger cân Σ
// FAIL) trả về sau khi rollback.
func (s *Service) PostInOwnTx(ctx context.Context, e domain.JournalEntry) (string, error) {
	if s.txm == nil {
		return "", apperr.Internal("TxManager chưa cấu hình cho PostInOwnTx")
	}
	var entryID string
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		id, err := s.Post(ctx, tx, e)
		if err != nil {
			return err
		}
		entryID = id
		return nil
	})
	if err != nil {
		return "", err
	}
	return entryID, nil
}

// ListEntriesQuery là input đọc danh sách bút toán đã giải mã ở edge.
type ListEntriesQuery struct {
	Cursor   string
	Limit    int32
	Book     string
	PeriodID string
}

// ListEntriesResult là một trang bút toán + cursor trang kế (rỗng nếu hết).
type ListEntriesResult struct {
	Items      []domain.JournalEntryRecord
	NextCursor string
}

// ListEntries trả một trang bút toán (keyset created_at DESC, id DESC). Tự decode
// cursor, chuẩn hoá limit, sinh NextCursor nếu trang đầy.
func (s *Service) ListEntries(ctx context.Context, q ListEntriesQuery) (ListEntriesResult, error) {
	afterNano, afterID, err := decodeCursor(q.Cursor)
	if err != nil {
		return ListEntriesResult{}, apperr.Validation("cursor không hợp lệ")
	}
	limit := normalizeLimit(q.Limit)
	f := domain.EntryFilter{Limit: limit, Book: q.Book, PeriodID: q.PeriodID}
	if q.Cursor != "" {
		f.AfterCreatedAt = time.Unix(0, afterNano).UTC()
		f.AfterID = afterID
		f.HasCursor = true
	}

	items, err := s.read.ListEntries(ctx, f)
	if err != nil {
		return ListEntriesResult{}, err
	}
	res := ListEntriesResult{Items: items}
	if int32(len(items)) == limit && limit > 0 {
		last := items[len(items)-1]
		res.NextCursor = encodeCursor(last.CreatedAt.UnixNano(), last.ID)
	}
	return res, nil
}

// GetEntry trả một bút toán theo id kèm các dòng.
func (s *Service) GetEntry(ctx context.Context, id string) (domain.JournalEntryRecord, error) {
	return s.read.GetEntry(ctx, id)
}

const dateLayout = "2006-01-02"

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

// Đảm bảo Service thoả Poster ở compile-time (port nội bộ cho module khác).
var _ Poster = (*Service)(nil)

// Poster là PORT NỘI BỘ: module orders/finance/vat post bút toán trong tx nghiệp
// vụ của họ qua đây (gộp atomic). KHÔNG có REST POST công khai ở P0. Định nghĩa ở
// app (không domain) vì nhận pgx.Tx — domain THUẦN không biết tx.
type Poster interface {
	Post(ctx context.Context, tx pgx.Tx, e domain.JournalEntry) (entryID string, err error)
}
