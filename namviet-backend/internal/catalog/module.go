// Package catalog là COMPOSITION ROOT của bounded context catalog: wiring adapter
// postgres + app use-case + adapter http, rồi export "mặt tiền" (Service +
// RegisterRoutes) cho edge. Module khác chỉ chạm package này hoặc port mà domain
// export — KHÔNG chạm repo/internal. Catalog read-only (ADR 0001): repo bind
// thẳng pool, không TxManager.
package catalog

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/catalog/app"
	cataloghttp "github.com/Maneva-AI/namviet-backend/internal/catalog/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/catalog/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Service là use-case đọc catalog mà edge (HTTP) dùng.
type Service = app.Service

// New dựng Service đầy đủ từ pool Postgres. Repo đọc bind thẳng pool (catalog
// read-only nên không cần transaction).
func New(pool *pgxpool.Pool) *Service {
	repo := postgres.NewRepo(appdb.New(pool))
	return app.New(repo)
}

// RegisterRoutes mount các operation /v1/products|categories|manufacturers lên
// huma.API. verifier dùng để verify token + ép quyền catalog.read. Edge
// (platform/server qua cmd/api) gọi hàm này.
func RegisterRoutes(api huma.API, svc *Service, verifier *authn.Verifier) {
	cataloghttp.Register(api, svc, verifier)
}
