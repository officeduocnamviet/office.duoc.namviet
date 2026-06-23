// Package inventory là COMPOSITION ROOT của bounded context inventory: wiring
// adapter postgres + app use-case + adapter http, rồi export "mặt tiền" cho edge
// (Service + RegisterRoutes — ĐỌC) và cho module khác (Deductor — port nội bộ trừ
// kho FEFO trong tx nghiệp vụ của họ). Module khác chỉ chạm package này hoặc port
// mà app/domain export — KHÔNG chạm repo/internal.
//
// ĐƯỜNG ĐỌC: repo bind thẳng pool (không tx). ĐƯỜNG GHI (P2 — trừ kho FEFO): repo
// ghi bind tới tx do caller/TxManager truyền (writerFromTx) + advisory lock chống
// bán âm. KHÔNG có REST POST công khai cho trừ kho — orders/POS gọi Deductor.
package inventory

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	inventoryhttp "github.com/Maneva-AI/namviet-backend/internal/inventory/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Service là use-case đọc inventory mà edge (HTTP) dùng.
type Service = app.Service

// DeductPort là port nội bộ trừ kho FEFO mà orders/POS gọi trong tx của họ.
type DeductPort = app.DeductPort

// StockInPort là port nội bộ NHẬP KHO (tạo lô + tăng tồn) mà purchasing gọi trong
// tx của họ (gộp atomic với post sổ + chi NCC).
type StockInPort = app.StockInPort

// New dựng Service đọc đầy đủ từ pool Postgres. Repo đọc bind thẳng pool (đường
// đọc không cần transaction).
func New(pool *pgxpool.Pool) *Service {
	repo := postgres.NewRepo(appdb.New(pool))
	return app.New(repo)
}

// NewDeductor dựng use-case GHI (trừ kho FEFO) từ pool Postgres.
//   - Repo ghi (StockWriter) bind tới tx do caller/TxManager truyền (writerFromTx).
//   - TxManager cho DeductFEFOInOwnTx (trừ kho độc lập).
//
// Trả *app.Deductor (thoả DeductPort). orders/POS lấy port này để trừ kho trong
// tx nghiệp vụ của họ (gộp atomic với post sổ + ghi tiền).
func NewDeductor(pool *pgxpool.Pool) *app.Deductor {
	writerFromTx := func(tx pgx.Tx) app.StockWriter {
		return postgres.NewWriteRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	return app.NewDeductor(writerFromTx, txm)
}

// NewStockInner dựng use-case GHI (NHẬP KHO — tạo lô + tăng tồn) từ pool Postgres.
//   - Repo ghi (StockWriter) bind tới tx do caller/TxManager truyền (writerFromTx).
//   - TxManager cho StockInOwnTx (nhập kho độc lập).
//
// Trả *app.StockInner (thoả StockInPort). purchasing lấy port này để nhập kho trong
// tx nghiệp vụ của họ (gộp atomic với post sổ + chi NCC).
func NewStockInner(pool *pgxpool.Pool) *app.StockInner {
	writerFromTx := func(tx pgx.Tx) app.StockWriter {
		return postgres.NewWriteRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	return app.NewStockInner(writerFromTx, txm)
}

// RegisterRoutes mount các operation ĐỌC /v1/warehouses|inventory/* lên huma.API.
// verifier dùng để verify token + ép quyền inventory.read. Edge (platform/server
// qua cmd/api) gọi hàm này.
func RegisterRoutes(api huma.API, svc *Service, verifier *authn.Verifier) {
	inventoryhttp.Register(api, svc, verifier)
}
