package domain

import "time"

// RefreshToken là token làm mới opaque, xoay vòng theo "family". Ta KHÔNG lưu
// token thô — chỉ lưu TokenHash (sha256 của token random 32 byte). Mỗi chuỗi
// refresh thuộc cùng FamilyID; khi phát hiện dùng lại token cũ (Used==true),
// toàn bộ family bị thu hồi (reuse-detection, ARCHITECTURE.md §9).
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	FamilyID  string
	Used      bool
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

// IsActive trả true nếu token còn dùng được tại thời điểm now: chưa revoke,
// chưa hết hạn. (Trạng thái Used được xử lý riêng ở use-case để phân biệt
// reuse với token hợp lệ.)
func (t RefreshToken) IsActive(now time.Time) bool {
	return !t.Revoked && now.Before(t.ExpiresAt)
}
