// Package finance là COMPOSITION ROOT của bounded context finance: wiring adapter
// postgres (repo ghi/đọc + TxManager) + app use-case + adapter http, rồi export
// "mặt tiền" cho edge (RegisterRoutes — CHỈ ĐỌC) và cho module khác (RecordPort —
// port nội bộ ghi phiếu THU trong tx nghiệp vụ của họ). Module khác chỉ chạm
// package này hoặc port app export — KHÔNG chạm repo/internal.
//
// P3 chỉ làm đường THU (RecordPaymentIn). KHÔNG có REST POST ghi phiếu — orders/POS
// gọi RecordPort.RecordPaymentIn(ctx, tx, ...) trong tx của họ (gộp atomic với trừ
// kho + post sổ). ⚠️ KHÔNG update fund_accounts.balance ở Go — trigger prod lo số dư.
package finance

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/finance/app"
	financehttp "github.com/Maneva-AI/namviet-backend/internal/finance/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/finance/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// RecordPort là port nội bộ ghi phiếu THU mà orders/POS gọi trong tx của họ.
type RecordPort = app.RecordPort

// RecordOutPort là port nội bộ ghi phiếu CHI (trả NCC mua hàng — mục 54) mà
// purchasing gọi trong tx của họ.
type RecordOutPort = app.RecordOutPort

// Module là mặt tiền của bounded context finance. Giữ Recorder (use-case ghi +
// đọc) đã wiring. cmd/api dựng Module rồi gọi RegisterRoutes; module nghiệp vụ
// khác lấy RecordPort() để ghi phiếu thu.
type Module struct {
	rec *app.Recorder
}

// NewModule dựng Module đầy đủ từ pool Postgres.
//   - Repo ghi (PaymentWriter) bind tới tx do caller/TxManager truyền (writerFromTx).
//   - Repo đọc (PaymentReader) bind thẳng pool.
//   - TxManager cho RecordPaymentInOwnTx (ghi phiếu độc lập).
func NewModule(pool *pgxpool.Pool) *Module {
	read := postgres.NewReadRepo(appdb.New(pool))
	writerFromTx := func(tx pgx.Tx) app.PaymentWriter {
		return postgres.NewWriteRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	return &Module{rec: app.New(writerFromTx, txm, read)}
}

// Recorder trả use-case (ghi phiếu + đọc) — dùng cho test/đọc nội bộ.
func (m *Module) Recorder() *app.Recorder { return m.rec }

// RecordPort trả port nội bộ để module orders/POS ghi phiếu THU trong tx nghiệp
// vụ của họ (gộp atomic). Recorder implement RecordPort.
func (m *Module) RecordPort() app.RecordPort { return m.rec }

// RecordOutPort trả port nội bộ để module purchasing ghi phiếu CHI (trả NCC) trong
// tx nghiệp vụ của họ (gộp atomic). Recorder implement RecordOutPort.
func (m *Module) RecordOutPort() app.RecordOutPort { return m.rec }

// RegisterRoutes mount các operation ĐỌC /v1/finance/* lên huma.API. verifier
// dùng để verify token + ép quyền finance.read.
func (m *Module) RegisterRoutes(api huma.API, verifier *authn.Verifier) {
	financehttp.Register(api, m.rec, verifier)
	financehttp.RegisterWrite(api, m.rec, verifier)
}
