package authn

import (
	"net/http"
	"strings"

	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx"
)

// Middleware trả một chi/net-http middleware: đọc header Authorization, verify
// Bearer JWT (ES256), nạp *Claims vào context. Thiếu/sai token → 401 envelope
// {data:null, error:{code:"unauthorized", ...}} và DỪNG chuỗi handler.
//
// Đây là điểm vào danh tính dùng chung cho mọi route cần đăng nhập.
func Middleware(v *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := bearerToken(r)
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "thiếu hoặc sai Bearer token")
				return
			}
			claims, err := v.Verify(raw)
			if err != nil {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "token không hợp lệ")
				return
			}
			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// bearerToken trích token từ "Authorization: Bearer <token>" (không phân biệt
// hoa thường ở từ khoá Bearer).
func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	const prefix = "bearer "
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return "", false
	}
	tok := strings.TrimSpace(h[len(prefix):])
	if tok == "" {
		return "", false
	}
	return tok, true
}
