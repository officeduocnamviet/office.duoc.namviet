package app

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
)

// accessTTL là thời gian sống access token (ngắn — ARCHITECTURE.md §9).
const accessTTL = 15 * time.Minute

// refreshTokenBytes là độ dài entropy của refresh token opaque (32 byte).
const refreshTokenBytes = 32

// TokenIssuer phát access JWT (ES256) và refresh token opaque. Dùng CHUNG
// authn.Claims để cấu trúc token phát ra khớp tuyệt đối với bộ verify ở
// platform/authn — một nguồn sự thật cho schema claims.
type TokenIssuer struct {
	priv *ecdsa.PrivateKey
	now  func() time.Time
}

// NewTokenIssuer tạo issuer từ private key EC P-256.
func NewTokenIssuer(priv *ecdsa.PrivateKey) *TokenIssuer {
	return &TokenIssuer{priv: priv, now: time.Now}
}

// IssueAccess phát access JWT ES256 với claims sub/user_type/perms/exp/iat/jti.
func (ti *TokenIssuer) IssueAccess(userID, userType string, perms []string) (string, error) {
	now := ti.now()
	claims := authn.Claims{
		UserType:    userType,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
			ID:        id.NewString(), // jti
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return tok.SignedString(ti.priv)
}

// NewRefreshToken sinh một refresh token opaque (random 32 byte, base64url) +
// hash sha256 của nó để lưu DB. Token thô chỉ trả về client một lần; DB chỉ giữ
// hash (không hồi phục được token từ hash).
func NewRefreshToken() (raw, hash string, err error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, HashRefreshToken(raw), nil
}

// HashRefreshToken băm token thô bằng sha256 → hex-less base64 cố định. Dùng cả
// khi phát (lưu hash) lẫn khi tra cứu (hash token client gửi rồi GetByHash).
func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
