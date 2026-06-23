// Package vat là COMPOSITION ROOT của bounded context vat: wiring adapter
// postgres (repo ghi/đọc + cấp số gapless + TxManager) + app use-case + adapter
// http, rồi export "mặt tiền" cho edge (RegisterRoutes — CHỈ ĐỌC) và cho module
// orders (IssuePort — port nội bộ phát hành HĐ VAT trong tx giao hàng của họ).
// Module khác chỉ chạm package này hoặc port app export — KHÔNG chạm repo/internal.
//
// HĐ VAT 100% đơn B2B, MST bắt buộc, thuộc SỔ TAX. Cấp số GAPLESS theo serial.
// Phát hành điện tử qua provider (VNPT/Viettel/MISA) = DEFER (chừa port sau).
// HTTP chỉ đọc; phát hành đi qua IssuePort (orders gọi trong tx của họ để gộp
// atomic với giao hàng + post sổ TAX — HĐ luôn khớp sự kiện).
package vat

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
	"github.com/Maneva-AI/namviet-backend/internal/vat/app"
	vathttp "github.com/Maneva-AI/namviet-backend/internal/vat/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/vat/internal/postgres"
)

// IssuePort là port nội bộ phát hành HĐ VAT mà orders gọi trong tx của họ.
type IssuePort = app.IssuePort

// IssueParams là tham số phát hành HĐ (re-export cho orders dựng input).
type IssueParams = app.IssueParams

// Module là mặt tiền của bounded context vat. Giữ Service (use-case phát hành +
// đọc) đã wiring. cmd/api dựng Module rồi gọi RegisterRoutes; module orders lấy
// IssuePort() để phát hành HĐ.
type Module struct {
	svc *app.Service
}

// NewModule dựng Module đầy đủ từ pool Postgres.
//   - Repo ghi (InvoiceStore) bind tới tx do caller/TxManager truyền (storeFromTx).
//   - Repo đọc (InvoiceReader) bind thẳng pool.
//   - TxManager cho IssueInvoiceInOwnTx (phát hành độc lập).
func NewModule(pool *pgxpool.Pool) *Module {
	read := postgres.NewReadRepo(appdb.New(pool))
	storeFromTx := func(tx pgx.Tx) app.InvoiceStore {
		return postgres.NewInvoiceRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	return &Module{svc: app.New(storeFromTx, txm, read)}
}

// Service trả use-case (phát hành + đọc) — dùng cho test/đọc nội bộ.
func (m *Module) Service() *app.Service { return m.svc }

// IssuePort trả port nội bộ để module orders phát hành HĐ VAT trong tx giao hàng
// của họ (gộp atomic). Service implement IssuePort.
func (m *Module) IssuePort() app.IssuePort { return m.svc }

// RegisterRoutes mount các operation ĐỌC /v1/vat/* lên huma.API. verifier dùng
// để verify token + ép quyền vat.read.
func (m *Module) RegisterRoutes(api huma.API, verifier *authn.Verifier) {
	vathttp.Register(api, m.svc, verifier)
}
