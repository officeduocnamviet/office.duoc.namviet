// Package authz là ĐIỂM ÉP QUYỀN DUY NHẤT (ARCHITECTURE.md §9: RBAC
// table-driven, 1 enforcement point). Middleware đọc claims đã verify từ
// context (do platform/authn nạp) và kiểm tra quyền theo mã (perm string).
// Không có per-scope guard rải rác, không RLS.
package authz

import (
	"net/http"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx"
)

// RequirePermission trả middleware chặn request nếu claims thiếu quyền perm.
//   - Không có claims trong context (chưa qua authn) → 401 unauthorized.
//   - Có claims nhưng thiếu quyền → 403 forbidden.
//
// Dùng SAU authn.Middleware trong chuỗi middleware của route.
func RequirePermission(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := authn.ClaimsFromContext(r.Context())
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "chưa xác thực")
				return
			}
			if !claims.HasPermission(perm) {
				httpx.WriteError(w, http.StatusForbidden, "forbidden", "thiếu quyền: "+perm)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
