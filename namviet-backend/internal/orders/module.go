// Package orders là COMPOSITION ROOT của bounded context orders: wiring adapter
// postgres (repo ĐỌC bind pool + repo GHI bound-tx + TxManager) + app use-case
// (đọc + ghi) + adapter http (GET đọc + POST ghi), rồi export "mặt tiền" cho edge.
// Module khác chỉ chạm package này hoặc port mà domain/app export — KHÔNG chạm
// repo/internal.
//
// P4a — NỀN ĐƯỜNG GHI: tạo đơn (PENDING) + state machine ĐƠN GIẢN (Confirm/
// Complete/Cancel) trong 1 transaction; mã đơn app tự sinh (app.order_code_seq +
// tiền tố), idempotent theo Idempotency-Key. KHÔNG đụng kho/tiền/sổ.
//
// HOÃN sang P4b (cần primitive cross-module): ShipOrder (CONFIRMED→SHIPPING, trừ
// kho FEFO + post sổ INTERNAL + phát HĐ VAT + post sổ TAX), RecordPayment (ghi
// finance_transactions + post PAYMENT_IN + suy payment_status), POS atomic
// (CreatePosSale), Refund (đảo bút toán + hoàn kho).
//
// GIẢ ĐỊNH TÍNH TIỀN (P4a, xác nhận khi tới P4b):
//   - line_total   = quantity * unit_price - discount       (mỗi dòng)
//   - total_amount = Σ (quantity * unit_price)              (TRƯỚC chiết khấu)
//   - final_amount = Σ line_total                            (SAU chiết khấu)
//   - CHƯA gồm VAT: total/final là tiền hàng ex-VAT. P4b khi post sổ TAX + phát HĐ
//     phải tính VAT riêng (vat.IssueInvoice nhận vat_rate từng dòng) — KHÔNG suy
//     ngược từ final_amount. CẦN kế toán/BA xác nhận: final_amount có gồm VAT
//     không (ảnh hưởng cách post sổ TAX ở P4b).
package orders

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/posting"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	ordershttp "github.com/Maneva-AI/namviet-backend/internal/orders/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/orders/internal/postgres"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Deps gom PORT NỘI BỘ của các bounded context khác mà orchestration (P4b) cần để
// gộp atomic trong tx giao dịch của orders: trừ kho FEFO (inventory.Deductor), post
// bút toán (accounting.Poster), phát HĐ VAT (vat.IssuePort), ghi phiếu THU
// (finance.RecordPort). cmd/api dựng các port này từ module tương ứng rồi truyền
// vào. Có thể để nil khi dump-openapi (pool nil) — orchestration route vẫn đăng ký
// schema nhưng không chạy.
type Deps struct {
	Deductor app.Deductor
	Poster   app.Poster
	Issuer   app.InvoiceIssuer
	Recorder app.PaymentRecorder
}

// Service là MẶT TIỀN của module orders (facade): gói use-case ĐỌC (List/Get) +
// GHI nền (Create + đổi trạng thái — P4a) + ORCHESTRATION cross-module (ShipOrder/
// RecordPayment/CreatePosSale — P4b). Giữ tên Service để wiring ổn định.
type Service struct {
	read        *app.Service
	write       *app.WriteService
	orch        *app.Orchestrator
	storeFromTx app.OrderStoreFromTx // dùng lại cho CreatePosSale (tạo đơn trong tx POS)
}

// New dựng Service đầy đủ từ pool Postgres + Deps (port cross-module cho P4b).
//   - Repo ĐỌC bind thẳng pool (đường đọc không cần transaction).
//   - Repo GHI (OrderStore/OrchestrationStore) bind tới tx do TxManager truyền.
//   - Orchestrator điều phối trừ kho + HĐ + post sổ + phiếu thu trong 1 tx atomic
//     (posting.DefaultRules: mã TK TT133 + cờ per-book mặc định, gắn cờ).
func New(pool *pgxpool.Pool, deps Deps) *Service {
	read := app.New(postgres.NewRepo(appdb.New(pool)))
	storeFromTx := func(tx pgx.Tx) app.OrderStore {
		return postgres.NewWriteRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	write := app.NewWrite(storeFromTx, txm)

	orchFromTx := func(tx pgx.Tx) app.OrchestrationStore {
		return postgres.NewOrchestrationRepo(appdb.New(pool).WithTx(tx))
	}
	orch := app.NewOrchestrator(orchFromTx, txm, deps.Deductor, deps.Poster, deps.Issuer, deps.Recorder, posting.DefaultRules)

	return &Service{read: read, write: write, orch: orch, storeFromTx: storeFromTx}
}

// Read trả use-case ĐỌC (dùng cho test/đọc nội bộ).
func (s *Service) Read() *app.Service { return s.read }

// Write trả use-case GHI (dùng cho test/ghi nội bộ).
func (s *Service) Write() *app.WriteService { return s.write }

// ListOrders/GetOrder delegate sang use-case ĐỌC — giữ API facade ổn định (caller
// cũ gọi thẳng svc.ListOrders/GetOrder không phải đổi).
func (s *Service) ListOrders(ctx context.Context, q app.ListOrdersQuery) (app.ListOrdersResult, error) {
	return s.read.ListOrders(ctx, q)
}

// GetOrder delegate sang use-case ĐỌC.
func (s *Service) GetOrder(ctx context.Context, id string) (app.OrderDetail, error) {
	return s.read.GetOrder(ctx, id)
}

// CreateOrder/ConfirmOrder/CompleteOrder/CancelOrder delegate sang use-case GHI.
func (s *Service) CreateOrder(ctx context.Context, in app.CreateOrderInput) (app.CreatedOrder, error) {
	return s.write.CreateOrder(ctx, in)
}

// ConfirmOrder delegate sang use-case GHI.
func (s *Service) ConfirmOrder(ctx context.Context, id string) (domain.Order, error) {
	return s.write.ConfirmOrder(ctx, id)
}

// CompleteOrder delegate sang use-case GHI.
func (s *Service) CompleteOrder(ctx context.Context, id string) (domain.Order, error) {
	return s.write.CompleteOrder(ctx, id)
}

// CancelOrder delegate sang use-case GHI.
func (s *Service) CancelOrder(ctx context.Context, id string) (domain.Order, error) {
	return s.write.CancelOrder(ctx, id)
}

// ShipOrder (P4b) giao hàng: CONFIRMED→SHIPPING gộp trừ kho FEFO + phát HĐ VAT +
// post sổ kép, atomic 1 tx. Delegate sang Orchestrator.
func (s *Service) ShipOrder(ctx context.Context, in app.ShipInput) (app.ShippedOrder, error) {
	return s.orch.ShipOrder(ctx, in)
}

// RecordPayment (P4b) ghi phiếu THU + post PAYMENT_IN + suy payment_status, atomic.
func (s *Service) RecordPayment(ctx context.Context, in app.RecordPaymentInput) (app.PaymentResult, error) {
	return s.orch.RecordPayment(ctx, in)
}

// CreatePosSale (P4b) bán lẻ tại quầy atomic: tạo đơn + trừ kho + thu tiền + (HĐ) +
// post sổ → COMPLETED/paid. Truyền storeFromTx (P4a) để tạo đơn trong CÙNG tx POS.
func (s *Service) CreatePosSale(ctx context.Context, in app.PosSaleInput) (app.PosSaleResult, error) {
	return s.orch.CreatePosSale(ctx, s.storeFromTx, in)
}

// RecordLumpSumPayment (mục 55) thu 1 cục từ khách → phân bổ cho đơn chưa tất toán
// CŨ NHẤT trước, atomic 1 tx. Delegate sang Orchestrator.
func (s *Service) RecordLumpSumPayment(ctx context.Context, in app.LumpSumInput) (app.LumpSumResult, error) {
	return s.orch.RecordLumpSumPayment(ctx, in)
}

// RegisterRoutes mount toàn bộ operation /v1/orders* + /v1/pos/* lên huma.API: GET
// (đọc, guard orders.read) + POST nền/orchestration (ghi, guard orders.write).
func RegisterRoutes(api huma.API, svc *Service, verifier *authn.Verifier) {
	ordershttp.Register(api, svc.read, verifier)
	ordershttp.RegisterWrite(api, svc.write, verifier)
	ordershttp.RegisterOrchestration(api, svc, verifier)
}
