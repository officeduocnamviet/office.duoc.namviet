package authn

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// HumaMiddleware là phiên bản middleware cho operation Huma (dùng cho route đăng
// ký qua huma.Register với Operation.Middlewares). Nó verify Bearer token, nạp
// *Claims vào huma.Context để handler lấy qua ClaimsFromContext(ctx.Context()).
// Thiếu/sai token → ghi 401 envelope qua huma.WriteErr và DỪNG (không gọi next).
//
// Có middleware net/http (Middleware) cho route chi thường và bản Huma này cho
// operation Huma — cùng một nguồn verify, giữ "1 enforcement point".
func HumaMiddleware(api huma.API, v *Verifier) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		raw, ok := bearerFromHuma(ctx)
		if !ok {
			writeUnauthorized(api, ctx, "thiếu hoặc sai Bearer token")
			return
		}
		claims, err := v.Verify(raw)
		if err != nil {
			writeUnauthorized(api, ctx, "token không hợp lệ")
			return
		}
		next(huma.WithValue(ctx, ctxKey{}, claims))
	}
}

func bearerFromHuma(ctx huma.Context) (string, bool) {
	h := ctx.Header("Authorization")
	const prefix = "bearer "
	if len(h) < len(prefix) || !equalFold(h[:len(prefix)], prefix) {
		return "", false
	}
	tok := h[len(prefix):]
	// Cắt khoảng trắng đầu/cuối thủ công (tránh import strings cho việc nhỏ).
	for len(tok) > 0 && (tok[0] == ' ' || tok[0] == '\t') {
		tok = tok[1:]
	}
	for len(tok) > 0 && (tok[len(tok)-1] == ' ' || tok[len(tok)-1] == '\t') {
		tok = tok[:len(tok)-1]
	}
	if tok == "" {
		return "", false
	}
	return tok, true
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

func writeUnauthorized(api huma.API, ctx huma.Context, msg string) {
	// huma.WriteErr dùng huma.NewError (đã được humax gán) → envelope đúng.
	_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, msg)
}
