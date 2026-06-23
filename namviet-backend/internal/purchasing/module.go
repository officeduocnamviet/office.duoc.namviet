// Package purchasing là COMPOSITION ROOT của bounded context purchasing (mua hàng &
// nhập kho, spec mục 54 — chiều MUA, đối xứng orders chiều BÁN): wiring adapter
// postgres (repo bound-tx + TxManager) + app use-case (orchestration cross-module) +
// adapter http, rồi export "mặt tiền" cho edge (Service + RegisterRoutes).
//
// CORE (đủ vòng đời tiền mua hàng): CreatePO (draft) → ConfirmPO (ordered) →
// ReceivePO (received: nhập kho inventory.StockIn + post sổ Dr 1561+133/Cr 331) →
// PaySupplier (paid: chi NCC finance.RecordPaymentOut + post Dr 331/Cr 111/112). Mỗi
// use-case GHI = 1 transaction atomic; cross-module nhận CÙNG tx (gộp atomic, rollback
// cả cụm). GUARD REPLAY: state machine guard + cờ created từ finance.
//
// HOÃN (mục 54 nâng cao): auto-tạo PO khi tồn<min, chương trình NCC, hợp đồng, upload
// HĐ tự điền. Supplier dùng supplier_id bigint + supplier_name text (KHÔNG FK — verify
// prod khi có supplier entity thật).
package purchasing

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/posting"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/app"
	purchasinghttp "github.com/Maneva-AI/namviet-backend/internal/purchasing/internal/http"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/internal/postgres"
)

// Deps gom PORT NỘI BỘ của các bounded context khác mà orchestration cần để gộp
// atomic trong tx giao dịch của purchasing: nhập kho (inventory.StockInPort), post
// bút toán (accounting.Poster), ghi phiếu CHI (finance.RecordOutPort). cmd/api dựng
// các port này từ module tương ứng rồi truyền vào. Có thể để nil khi dump-openapi.
type Deps struct {
	StockIn app.StockInner
	Poster  app.Poster
	Payer   app.SupplierPayer
}

// Service là MẶT TIỀN của module purchasing (facade): gói use-case ĐỌC (List/Get) +
// GHI cross-module (CreatePO/ConfirmPO/CancelPO/ReceivePO/PaySupplier).
type Service struct {
	svc *app.Service
}

// New dựng Service đầy đủ từ pool Postgres + Deps (port cross-module).
//   - storeFromTx/listFromTx bind repo tới tx do TxManager truyền.
//   - posting.DefaultRules: mã TK TT133 (1561/133/331/111/112) + cờ per-book.
func New(pool *pgxpool.Pool, deps Deps) *Service {
	storeFromTx := func(tx pgx.Tx) app.Store {
		return postgres.NewRepo(appdb.New(pool).WithTx(tx))
	}
	listFromTx := func(tx pgx.Tx) app.ListReader {
		return postgres.NewRepo(appdb.New(pool).WithTx(tx))
	}
	txm := postgres.NewTxManager(pool)
	svc := app.New(storeFromTx, listFromTx, txm, deps.StockIn, deps.Poster, deps.Payer, posting.DefaultRules)
	return &Service{svc: svc}
}

// App trả use-case (dùng cho test/đọc nội bộ).
func (s *Service) App() *app.Service { return s.svc }

// RegisterRoutes mount toàn bộ operation /v1/purchase-orders* lên huma.API.
func RegisterRoutes(api huma.API, svc *Service, verifier *authn.Verifier) {
	purchasinghttp.Register(api, svc.svc, verifier)
}
