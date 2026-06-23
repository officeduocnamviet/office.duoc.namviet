// Package http là ADAPTER vào của identity: DTO + handler Huma + đăng ký route.
// Nó dịch HTTP <-> use-case (app) và map lỗi nghiệp vụ (apperr) sang envelope
// qua humax.FromAppErr. Nằm dưới internal/ nên module khác không import được.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// AuthService là cổng use-case mà handler cần (interface để test handler bằng
// fake, không cần service thật/DB).
type AuthService interface {
	Login(ctx context.Context, email, password string) (app.Tokens, error)
	Refresh(ctx context.Context, rawRefresh string) (app.Tokens, error)
	Logout(ctx context.Context, rawRefresh string) error
	Me(ctx context.Context, userID string) (app.MeResult, error)
}

// Register đăng ký toàn bộ operation /v1/auth/* lên huma.API. verifier dùng cho
// middleware xác thực route /me (qua platform/authn — 1 enforcement point).
func Register(api huma.API, svc AuthService, verifier *authn.Verifier) {
	registerLogin(api, svc)
	registerRefresh(api, svc)
	registerLogout(api, svc)
	registerMe(api, svc, verifier)
}

// ---- DTO ----

type loginInput struct {
	Body struct {
		Email    string `json:"email" format:"email" doc:"Email đăng nhập"`
		Password string `json:"password" minLength:"1" doc:"Mật khẩu"`
	}
}

type tokensBody struct {
	AccessToken  string `json:"access_token" doc:"JWT truy cập (ES256, TTL ngắn)"`
	RefreshToken string `json:"refresh_token" doc:"Refresh token opaque (xoay vòng)"`
	TokenType    string `json:"token_type" doc:"Loại token" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" doc:"TTL access token (giây)"`
}

type tokensOutput struct {
	Body tokensBody
}

type refreshInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" minLength:"1" doc:"Refresh token đang giữ"`
	}
}

type logoutInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" minLength:"1" doc:"Refresh token cần thu hồi (cả family)"`
	}
}

type logoutOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

type meOutput struct {
	Body struct {
		UserID      string   `json:"user_id"`
		Email       string   `json:"email"`
		UserType    string   `json:"user_type"`
		Permissions []string `json:"permissions"`
	}
}

func toTokensOutput(t app.Tokens) *tokensOutput {
	out := &tokensOutput{}
	out.Body = tokensBody{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    t.ExpiresIn,
	}
	return out
}

// ---- Handlers ----

func registerLogin(api huma.API, svc AuthService) {
	huma.Register(api, huma.Operation{
		OperationID:   "auth-login",
		Method:        http.MethodPost,
		Path:          "/v1/auth/login",
		Summary:       "Đăng nhập bằng email + mật khẩu",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *loginInput) (*tokensOutput, error) {
		tok, err := svc.Login(ctx, in.Body.Email, in.Body.Password)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		return toTokensOutput(tok), nil
	})
}

func registerRefresh(api huma.API, svc AuthService) {
	huma.Register(api, huma.Operation{
		OperationID:   "auth-refresh",
		Method:        http.MethodPost,
		Path:          "/v1/auth/refresh",
		Summary:       "Xoay vòng refresh token (phát access + refresh mới)",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *refreshInput) (*tokensOutput, error) {
		tok, err := svc.Refresh(ctx, in.Body.RefreshToken)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		return toTokensOutput(tok), nil
	})
}

func registerLogout(api huma.API, svc AuthService) {
	huma.Register(api, huma.Operation{
		OperationID:   "auth-logout",
		Method:        http.MethodPost,
		Path:          "/v1/auth/logout",
		Summary:       "Đăng xuất (thu hồi cả family refresh token)",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *logoutInput) (*logoutOutput, error) {
		if err := svc.Logout(ctx, in.Body.RefreshToken); err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &logoutOutput{}
		out.Body.Status = "ok"
		return out, nil
	})
}

func registerMe(api huma.API, svc AuthService, verifier *authn.Verifier) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-me",
		Method:      http.MethodGet,
		Path:        "/v1/auth/me",
		Summary:     "Thông tin + quyền của user hiện tại",
		Tags:        []string{"auth"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{authn.HumaMiddleware(api, verifier)},
	}, func(ctx context.Context, _ *struct{}) (*meOutput, error) {
		claims, ok := authn.ClaimsFromContext(ctx)
		if !ok {
			// Middleware đã chặn; phòng vệ.
			return nil, humax.FromAppErr(errUnauthenticated())
		}
		me, err := svc.Me(ctx, claims.UserID())
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &meOutput{}
		out.Body.UserID = me.UserID
		out.Body.Email = me.Email
		out.Body.UserType = me.UserType
		out.Body.Permissions = me.Permissions
		return out, nil
	})
}
