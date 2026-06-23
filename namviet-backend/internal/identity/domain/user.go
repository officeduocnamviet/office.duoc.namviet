// Package domain là LÕI THUẦN của bounded context identity: entity, value
// object, business rule và PORT interface. KHÔNG import pgx/http/huma/framework
// (ARCHITECTURE.md §3). Phụ thuộc đi một chiều: adapters → app → domain.
package domain

import "time"

// HashAlgo là thuật toán băm mật khẩu được hỗ trợ. Khớp CHECK constraint của
// app.users (chỉ 'argon2id' | 'bcrypt').
type HashAlgo string

const (
	// HashArgon2id là thuật toán băm mặc định cho mật khẩu mới.
	HashArgon2id HashAlgo = "argon2id"
	// HashBcrypt là thuật toán legacy nhập từ GoTrue (lazy rehash sang argon2id).
	HashBcrypt HashAlgo = "bcrypt"
)

// User là aggregate gốc của identity. Email là định danh đăng nhập (citext,
// duy nhất). PasswordHash + HashAlgo lưu băm; domain không bao giờ giữ plaintext.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	HashAlgo     HashAlgo
	UserType     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CanLogin trả true nếu user được phép đăng nhập (đang active). Business rule
// thuần, không phụ thuộc hạ tầng.
func (u User) CanLogin() bool {
	return u.IsActive
}
