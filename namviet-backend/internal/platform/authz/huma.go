package authz

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
)

// RequirePermissionHuma là middleware cho operation Huma (đăng ký qua
// Operation.Middlewares) gộp XÁC THỰC + ÉP QUYỀN trong một bước, giữ "1
// enforcement point" (ARCHITECTURE.md §9):
//   - thiếu/sai Bearer token → 401 unauthorized;
//   - token hợp lệ nhưng thiếu quyền perm → 403 forbidden;
//   - hợp lệ → nạp *Claims vào context (ClaimsFromContext) rồi gọi next.
//
// Dùng cho route đọc của các module (vd catalog: RequirePermissionHuma(api,
// verifier, "catalog.read")). Cùng nguồn verify như authn.HumaMiddleware.
func RequirePermissionHuma(api huma.API, v *authn.Verifier, perm string) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		raw, ok := bearerFromHuma(ctx)
		if !ok {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "thiếu hoặc sai Bearer token")
			return
		}
		claims, err := v.Verify(raw)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "token không hợp lệ")
			return
		}
		if !claims.HasPermission(perm) {
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "thiếu quyền: "+perm)
			return
		}
		next(huma.WithValue(ctx, authn.CtxKey(), claims))
	}
}

// bearerFromHuma trích token sau "Bearer " (không phân biệt hoa thường) từ header
// Authorization của Huma context.
func bearerFromHuma(ctx huma.Context) (string, bool) {
	h := ctx.Header("Authorization")
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
