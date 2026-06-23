// Package authn xử lý XÁC THỰC ở edge: verify access JWT (pin thuật toán ES256),
// nạp claims đã verify vào context. Đây là hạ tầng kỹ thuật (platform), không
// chứa business; module identity sở hữu việc PHÁT token, còn package này lo việc
// VERIFY + đưa danh tính vào request context cho mọi context khác dùng chung
// (ARCHITECTURE.md §9: "1 enforcement point", pin alg).
package authn

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims là payload access token đã verify. Khớp với token mà
// identity/app.TokenIssuer phát: sub (user id), user_type, perms (mã quyền),
// jti, exp/iat. Dùng RegisteredClaims của jwt để xử lý exp/iat/sub chuẩn.
type Claims struct {
	UserType    string   `json:"user_type"`
	Permissions []string `json:"perms"`
	jwt.RegisteredClaims
}

// UserID trả subject (id người dùng) dạng chuỗi.
func (c Claims) UserID() string {
	return c.Subject
}

// HasPermission kiểm tra perms có chứa mã quyền perm không (so khớp đầy đủ).
func (c Claims) HasPermission(perm string) bool {
	for _, p := range c.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// ctxKey là khoá riêng để tránh va chạm khi đặt giá trị vào context.
type ctxKey struct{}

// WithClaims gắn claims vào context (dùng nội bộ middleware + test).
func WithClaims(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// ClaimsFromContext lấy claims đã verify khỏi context. ok=false nếu request
// chưa qua middleware xác thực (hoặc token không hợp lệ).
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(*Claims)
	return c, ok
}

// CtxKey trả khoá context (đã bọc) để middleware Huma ở package khác (vd
// platform/authz) gắn claims qua huma.WithValue mà handler vẫn đọc được bằng
// ClaimsFromContext. Giữ ctxKey private trong authn — chỉ lộ giá trị khoá.
func CtxKey() any { return ctxKey{} }

// nowFunc cho phép test ghi đè thời gian. Mặc định time.Now.
var nowFunc = time.Now
