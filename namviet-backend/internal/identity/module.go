// Package identity là COMPOSITION ROOT của bounded context identity: nó wiring
// adapter (postgres) + app use-case + adapter http, rồi export những gì module
// khác / edge cần (AuthService cho HTTP, không export repo). Đây là "mặt tiền"
// duy nhất của context — module khác chỉ chạm package này hoặc port mà nó export.
package identity

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
	identityhttp "github.com/Maneva-AI/namviet-backend/internal/identity/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/identity/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Service là use-case xác thực mà edge (HTTP) dùng. Alias để gói khác không cần
// import package app trực tiếp.
type Service = app.AuthService

// New dựng AuthService đầy đủ từ pool Postgres + private/public key đã nạp.
//   - Repo đọc (Users/Roles) bind thẳng pool.
//   - TxManager bind Queries vào tx cho thao tác refresh nguyên tử.
func New(pool *pgxpool.Pool, issuer *app.TokenIssuer) *Service {
	baseQ := appdb.New(pool)
	users := postgres.NewUserRepo(baseQ)
	roles := postgres.NewRoleRepo(baseQ)

	txm := postgres.NewTxManager(pool, func(tx pgx.Tx) app.Repos {
		q := baseQ.WithTx(tx)
		return app.Repos{
			Users:  postgres.NewUserRepo(q),
			Tokens: postgres.NewRefreshTokenRepo(q),
			Roles:  postgres.NewRoleRepo(q),
		}
	})

	return app.NewAuthService(users, roles, txm, app.NewPasswordHasher(), issuer)
}

// NewTokenIssuer re-export để cmd/server dựng issuer từ key mà không import app.
var NewTokenIssuer = app.NewTokenIssuer

// RegisterRoutes mount các operation /v1/auth/* của module lên huma.API. verifier
// dùng cho middleware xác thực route /me. Edge (platform/server) gọi hàm này.
func RegisterRoutes(api huma.API, svc *Service, verifier *authn.Verifier) {
	identityhttp.Register(api, svc, verifier)
}
