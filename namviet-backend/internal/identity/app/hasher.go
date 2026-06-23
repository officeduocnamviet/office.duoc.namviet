package app

import (
	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"

	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
)

// PasswordHasher băm + verify mật khẩu. Hash MỚI luôn dùng argon2id với tham số
// chốt ở ARCHITECTURE.md §9 (m=19456, t=2, p=1). Verify hỗ trợ cả argon2id lẫn
// bcrypt (legacy nhập từ GoTrue) — token bcrypt verify đúng sẽ báo needsRehash
// để use-case nâng cấp dần (lazy rehash), KHÔNG UPDATE auth.users trực tiếp.
type PasswordHasher struct {
	params *argon2id.Params
}

// NewPasswordHasher tạo hasher với tham số argon2id chuẩn.
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		params: &argon2id.Params{
			Memory:      19456, // 19 MiB
			Iterations:  2,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

// Hash băm plaintext bằng argon2id. Trả chuỗi PHC encode + thuật toán argon2id.
func (h *PasswordHasher) Hash(plain string) (hash string, algo domain.HashAlgo, err error) {
	encoded, err := argon2id.CreateHash(plain, h.params)
	if err != nil {
		return "", "", err
	}
	return encoded, domain.HashArgon2id, nil
}

// Verify so khớp plaintext với hash đã lưu theo algo.
//   - argon2id: so khớp PHC.
//   - bcrypt: so khớp bcrypt; nếu đúng → needsRehash=true (cần nâng cấp argon2id).
//
// Trả ok=false (không lỗi) khi mật khẩu sai. err chỉ khác nil khi hash hỏng
// định dạng / thuật toán không hỗ trợ.
func (h *PasswordHasher) Verify(plain, hash string, algo domain.HashAlgo) (ok, needsRehash bool, err error) {
	switch algo {
	case domain.HashArgon2id:
		match, err := argon2id.ComparePasswordAndHash(plain, hash)
		if err != nil {
			return false, false, err
		}
		return match, false, nil
	case domain.HashBcrypt:
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
		if err != nil {
			// Sai mật khẩu là trường hợp thường, không phải lỗi hệ thống.
			return false, false, nil
		}
		// Đúng mật khẩu nhưng hash là bcrypt legacy → cần rehash sang argon2id.
		return true, true, nil
	default:
		return false, false, errUnsupportedAlgo(algo)
	}
}
