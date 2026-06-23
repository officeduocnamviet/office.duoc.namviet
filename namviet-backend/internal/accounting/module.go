// Package accounting là COMPOSITION ROOT của bounded context accounting: wiring
// adapter postgres (repo ghi/đọc + TxManager) + app use-case + adapter http, rồi
// export "mặt tiền" cho edge (RegisterRoutes — CHỈ ĐỌC) và cho module khác
// (Poster — port nội bộ post bút toán trong tx nghiệp vụ của họ). Module khác chỉ
// chạm package này hoặc port app export — KHÔNG chạm repo/internal.
//
// APPEND-ONLY + 2 sổ INTERNAL/TAX không sync: HTTP chỉ đọc; ghi sổ đi qua Poster
// (orders/finance/vat gọi trong tx của họ để gộp atomic — sổ luôn khớp sự kiện).
package accounting

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/app"
	accountinghttp "github.com/Maneva-AI/namviet-backend/internal/accounting/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Module là mặt tiền của bounded context accounting. Giữ Service (use-case đọc +
// post) đã wiring. cmd/api dựng Module rồi gọi RegisterRoutes; module nghiệp vụ
// khác lấy Poster() để post bút toán.
type Module struct {
	svc *app.Service
}

// NewModule dựng Module đầy đủ từ pool Postgres.
//   - Repo ghi (EntryStore) bind tới tx do caller/TxManager truyền (storeFromTx).
//   - Repo đọc (ReadStore) bind thẳng pool.
//   - TxManager cho PostInOwnTx (post độc lập).
func NewModule(pool *pgxpool.Pool) *Module {
	read := postgres.NewReadRepo(appdb.New(pool))
	storeFromTx := func(tx pgx.Tx) app.EntryStore {
		return postgres.NewEntryRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	return &Module{svc: app.New(storeFromTx, txm, read)}
}

// Service trả use-case (đọc sổ) — dùng cho test/đọc nội bộ.
func (m *Module) Service() *app.Service { return m.svc }

// Poster trả port nội bộ để module orders/finance/vat post bút toán trong tx
// nghiệp vụ của họ (gộp atomic). Service implement Poster.
func (m *Module) Poster() app.Poster { return m.svc }

// RegisterRoutes mount các operation ĐỌC /v1/accounting/* lên huma.API. verifier
// dùng để verify token + ép quyền accounting.read.
func (m *Module) RegisterRoutes(api huma.API, verifier *authn.Verifier) {
	accountinghttp.Register(api, m.svc, verifier)
}
