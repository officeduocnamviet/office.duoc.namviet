package domain

import "time"

// EntryFilter gom điều kiện lọc + keyset pagination cho danh sách bút toán đọc.
// Tất cả tiêu chí là optional. Limit/AfterID do tầng app chuẩn hoá. Keyset theo
// (created_at DESC, id DESC): trang kế lấy bút toán "cũ hơn" mốc cursor.
type EntryFilter struct {
	// AfterCreatedAt: mốc cursor (created_at của bản ghi cuối trang trước). Chỉ áp
	// khi HasCursor = true (zero time hợp lệ ở epoch, nên cần cờ riêng).
	AfterCreatedAt time.Time
	// AfterID: id (uuid) của bản ghi cuối trang trước (tie-break cùng created_at).
	AfterID string
	// HasCursor phân biệt "trang đầu" với mốc cursor hợp lệ.
	HasCursor bool
	// Limit: số bản ghi tối đa mỗi trang (đã chuẩn hoá ở app, > 0).
	Limit int32
	// Book lọc theo sổ ("INTERNAL"/"TAX"); rỗng = cả hai sổ.
	Book string
	// PeriodID lọc theo kỳ kế toán (uuid); rỗng = mọi kỳ.
	PeriodID string
}

// JournalEntryRecord là một bút toán ĐÃ GHI (đọc lại từ sổ): JournalEntry thuần
// + các trường định danh/persisted (ID, PeriodID, CreatedAt). Tách khỏi
// JournalEntry (thứ caller dựng để POST, chưa có id/kỳ) để đường đọc và đường ghi
// không lẫn lộn. Lines rỗng ở danh sách (chỉ nạp khi GetEntry một bút toán).
type JournalEntryRecord struct {
	ID         string
	Book       Book
	EntryDate  time.Time
	PeriodID   string
	SourceType string
	SourceID   string
	Memo       string
	CreatedAt  time.Time
	Lines      []EntryLine
}
