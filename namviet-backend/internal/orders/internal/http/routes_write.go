// routes_write.go — ADAPTER GHI của orders (P4a): POST tạo đơn + đổi trạng thái
// ĐƠN GIẢN. Guard quyền "orders.write" (khác đường ĐỌC dùng "orders.read"). Tiền
// vào/ra DTO dạng CHUỖI thập phân — KHÔNG float (FE không mất chính xác). Số lượng
// vào dạng số nguyên (cột quantity INTEGER). Idempotency-Key qua header cho POST
// tạo đơn (1 key → 1 đơn). KHÔNG có ShipOrder/RecordPayment/POS (P4b).
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// WriteService là cổng use-case GHI mà handler cần (interface để test bằng fake).
type WriteService interface {
	CreateOrder(ctx context.Context, in app.CreateOrderInput) (app.CreatedOrder, error)
	ConfirmOrder(ctx context.Context, orderID string) (domain.Order, error)
	CompleteOrder(ctx context.Context, orderID string) (domain.Order, error)
	CancelOrder(ctx context.Context, orderID string) (domain.Order, error)
}

const permWrite = "orders.write"

// RegisterWrite đăng ký các operation GHI /v1/orders* (POST) lên huma.API. Mọi
// route yêu cầu quyền orders.write (verify token + check perm — 1 enforcement point).
func RegisterWrite(api huma.API, svc WriteService, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permWrite)
	registerCreateOrder(api, svc, guard)
	registerConfirmOrder(api, svc, guard)
	registerCompleteOrder(api, svc, guard)
	registerCancelOrder(api, svc, guard)
}

// ---- DTO ----

type createOrderLineDTO struct {
	ProductID int64  `json:"product_id" minimum:"1" doc:"ID sản phẩm (bigint)"`
	Quantity  int64  `json:"quantity" minimum:"1" doc:"Số lượng (số nguyên dương)"`
	UOM       string `json:"uom" minLength:"1" doc:"Đơn vị tính (VD Hộp/Vỉ/Cái)"`
	UnitPrice string `json:"unit_price" doc:"Đơn giá bán (chuỗi thập phân, không float)"`
	// discount optional (omitempty) — rỗng/không gửi = 0.
	Discount string `json:"discount,omitempty" doc:"Chiết khấu dòng (chuỗi thập phân; bỏ trống = 0)"`
}

// createOrderInput — các trường optional dùng ,omitempty để Huma KHÔNG đánh
// required (customer_id null cho POS B2C; creator_id/note có thể trống).
type createOrderInput struct {
	IdemKey string `header:"Idempotency-Key" doc:"Khoá idempotency tạo đơn (1 key → 1 đơn); rỗng = không idempotent"`
	Body    struct {
		CustomerID *int64               `json:"customer_id,omitempty" doc:"ID khách (bigint); bỏ trống nếu đơn không gắn khách (POS B2C)"`
		OrderType  string               `json:"order_type,omitempty" enum:"B2B,B2C" doc:"Loại đơn; bỏ trống = B2C"`
		CreatorID  string               `json:"creator_id,omitempty" doc:"Người lập (uuid); bỏ trống nếu không có"`
		Note       string               `json:"note,omitempty"`
		Lines      []createOrderLineDTO `json:"lines" minItems:"1" doc:"Các dòng hàng (ít nhất 1)"`
	}
}

type orderDetailOutput struct {
	Body struct {
		Order orderDTO       `json:"order"`
		Lines []orderLineDTO `json:"lines"`
	}
}

type orderActionInput struct {
	ID string `path:"id" doc:"ID đơn hàng (uuid)"`
}

type orderActionOutput struct {
	Body struct {
		Order orderDTO `json:"order"`
	}
}

// ---- Handlers ----

func registerCreateOrder(api huma.API, svc WriteService, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "orders-create",
		Method:        http.MethodPost,
		Path:          "/v1/orders",
		Summary:       "Tạo đơn hàng (PENDING) + dòng hàng",
		Tags:          []string{"orders"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createOrderInput) (*orderDetailOutput, error) {
		lines, err := toDraftLines(in.Body.Lines)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		res, err := svc.CreateOrder(ctx, app.CreateOrderInput{
			CustomerID: in.Body.CustomerID,
			OrderType:  in.Body.OrderType,
			CreatorID:  in.Body.CreatorID,
			Note:       in.Body.Note,
			Lines:      lines,
			IdemKey:    in.IdemKey,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &orderDetailOutput{}
		out.Body.Order = toOrderDTO(res.Order)
		out.Body.Lines = make([]orderLineDTO, 0, len(res.Lines))
		for _, l := range res.Lines {
			out.Body.Lines = append(out.Body.Lines, toOrderLineDTO(l))
		}
		return out, nil
	})
}

func registerConfirmOrder(api huma.API, svc WriteService, guard func(huma.Context, func(huma.Context))) {
	registerStatusAction(api, guard, "orders-confirm", "/v1/orders/{id}/confirm",
		"Duyệt đơn (PENDING → CONFIRMED)", svc.ConfirmOrder)
}

func registerCompleteOrder(api huma.API, svc WriteService, guard func(huma.Context, func(huma.Context))) {
	registerStatusAction(api, guard, "orders-complete", "/v1/orders/{id}/complete",
		"Hoàn tất đơn (SHIPPING → COMPLETED)", svc.CompleteOrder)
}

func registerCancelOrder(api huma.API, svc WriteService, guard func(huma.Context, func(huma.Context))) {
	registerStatusAction(api, guard, "orders-cancel", "/v1/orders/{id}/cancel",
		"Huỷ đơn (PENDING/CONFIRMED → CANCELLED)", svc.CancelOrder)
}

// registerStatusAction đăng ký một POST đổi-trạng-thái dùng chung (action không
// CRUD: POST /v1/orders/{id}/<action>). fn là use-case tương ứng (Confirm/Complete/
// Cancel). Trả đơn sau cập nhật.
func registerStatusAction(
	api huma.API,
	guard func(huma.Context, func(huma.Context)),
	opID, path, summary string,
	fn func(ctx context.Context, orderID string) (domain.Order, error),
) {
	huma.Register(api, huma.Operation{
		OperationID:   opID,
		Method:        http.MethodPost,
		Path:          path,
		Summary:       summary,
		Tags:          []string{"orders"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *orderActionInput) (*orderActionOutput, error) {
		o, err := fn(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &orderActionOutput{}
		out.Body.Order = toOrderDTO(o)
		return out, nil
	})
}

// toDraftLines chuyển DTO dòng hàng (tiền chuỗi) sang domain.DraftLine (money
// decimal). Lỗi parse tiền → apperr.Validation (422). discount rỗng = 0.
func toDraftLines(in []createOrderLineDTO) ([]domain.DraftLine, error) {
	out := make([]domain.DraftLine, 0, len(in))
	for _, l := range in {
		unitPrice, err := money.FromString(l.UnitPrice)
		if err != nil {
			return nil, apperr.Validation("unit_price không hợp lệ: " + l.UnitPrice)
		}
		discount, err := money.FromString(l.Discount)
		if err != nil {
			return nil, apperr.Validation("discount không hợp lệ: " + l.Discount)
		}
		out = append(out, domain.DraftLine{
			ProductID: l.ProductID,
			Quantity:  domain.QuantityFromInt(l.Quantity),
			UOM:       l.UOM,
			UnitPrice: unitPrice,
			Discount:  discount,
		})
	}
	return out, nil
}
