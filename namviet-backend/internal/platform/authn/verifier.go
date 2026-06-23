package authn

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Verifier verify access JWT bằng public key EC P-256, PIN cứng thuật toán
// ES256. Bất kỳ token nào ký bằng alg khác (HS256, none, RS256...) đều bị từ
// chối — chống tấn công alg-confusion.
type Verifier struct {
	pub *ecdsa.PublicKey
}

// NewVerifier tạo Verifier từ public key EC P-256.
func NewVerifier(pub *ecdsa.PublicKey) *Verifier {
	return &Verifier{pub: pub}
}

// ErrInvalidToken là lỗi chung khi token không hợp lệ (sai chữ ký, hết hạn,
// sai alg...). Cố ý KHÔNG rò chi tiết để tránh lộ thông tin; middleware map
// sang 401 envelope.
var ErrInvalidToken = errors.New("token không hợp lệ")

// Verify parse + verify chuỗi token, trả Claims nếu hợp lệ. Pin ES256 qua
// WithValidMethods + kiểm tra kiểu key trong keyfunc.
func (v *Verifier) Verify(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256"}), // CHỈ ES256
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(nowFunc),
	)
	tok, err := parser.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// Phòng vệ kép: dù WithValidMethods đã chặn, vẫn ép kiểu method ECDSA.
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("%w: alg không phải ECDSA", ErrInvalidToken)
		}
		return v.pub, nil
	})
	if err != nil || !tok.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
