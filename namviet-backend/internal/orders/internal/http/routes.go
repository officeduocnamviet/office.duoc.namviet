// Package http là ADAPTER vào của orders: DTO + handler Huma + đăng ký route.
// Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, và bảo vệ mọi route bằng authz.RequirePermissionHuma
// ("orders.read"). Tiền (tổng, đơn giá, đã thu, còn nợ) và số lượng ra DTO dạng
// CHUỖI thập phân — KHÔNG float — để FE không mất chính xác. Nằm dưới internal/
// nên module khác không import.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Service là cổng use-case mà handler cần (interface để test handler bằng fake).
type Service interface {
	ListOrders(ctx context.Context, q app.ListOrdersQuery) (app.ListOrdersResult, error)
	GetOrder(ctx context.Context, id string) (app.OrderDetail, error)
}

const permRead = "orders.read"

const dateLayout = "2006-01-02"

// Register đăng ký toàn bộ operation /v1/orders* lên huma.API. Mọi route yêu cầu
// quyền orders.read (verify token + check perm qua authz.RequirePermissionHuma —
// 1 enforcement point).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListOrders(api, svc, guard)
	registerGetOrder(api, svc, guard)
}

// ---- DTO ----

// paymentDTO trình bày suy diễn thanh toán read-only (KHÔNG có paid_amount ở DB).
// Tiền là chuỗi thập phân. remaining có thể ÂM (thu thừa) — FE tự xử lý hiển thị.
type paymentDTO struct {
	Final     string `json:"final" doc:"Tổng tiền khách phải trả (chuỗi thập phân, không float)"`
	Paid      string `json:"paid" doc:"Đã thu — suy diễn từ finance_transactions (phiếu thu hoàn tất, sổ thực tế)"`
	Remaining string `json:"remaining" doc:"Còn nợ = final - paid (có thể âm nếu thu thừa)"`
}

type orderDTO struct {
	ID            string     `json:"id" doc:"ID đơn (uuid)"`
	Code          string     `json:"code" doc:"Mã hóa đơn (VD HD123)"`
	CustomerID    *int64     `json:"customer_id" doc:"ID khách (bigint); null nếu đơn không gắn khách"`
	CreatorID     string     `json:"creator_id" doc:"Người lập (uuid); rỗng nếu không có"`
	Status        string     `json:"status" doc:"Trạng thái xử lý đơn"`
	OrderType     string     `json:"order_type" doc:"B2C | B2B"`
	TotalAmount   string     `json:"total_amount" doc:"Tổng trước chiết khấu (chuỗi thập phân)"`
	FinalAmount   string     `json:"final_amount" doc:"Tổng khách phải trả (chuỗi thập phân)"`
	PaymentStatus string     `json:"payment_status" doc:"unpaid | partial | paid"`
	Note          string     `json:"note"`
	CreatedAt     string     `json:"created_at" doc:"Thời điểm tạo (RFC3339)"`
	Payment       paymentDTO `json:"payment"`
}

type orderLineDTO struct {
	ID         string `json:"id" doc:"ID dòng (uuid)"`
	ProductID  int64  `json:"product_id"`
	Quantity   string `json:"quantity" doc:"Số lượng (chuỗi thập phân, không float)"`
	UOM        string `json:"uom" doc:"Đơn vị tính được chọn"`
	UnitPrice  string `json:"unit_price" doc:"Đơn giá bán (chuỗi thập phân)"`
	Discount   string `json:"discount" doc:"Chiết khấu dòng (chuỗi thập phân)"`
	LineTotal  string `json:"line_total" doc:"Thành tiền dòng (chuỗi thập phân)"`
	IsGift     bool   `json:"is_gift"`
	BatchNo    string `json:"batch_no" doc:"Lô đã chọn; rỗng nếu chưa gán"`
	ExpiryDate string `json:"expiry_date" doc:"Hạn dùng lô (YYYY-MM-DD); rỗng nếu không có"`
	Note       string `json:"note"`
}

func toPaymentDTO(p domain.PaymentSummary) paymentDTO {
	return paymentDTO{
		Final:     p.Final.String(),
		Paid:      p.Paid.String(),
		Remaining: p.Remaining.String(),
	}
}

func toOrderDTO(o domain.Order) orderDTO {
	created := ""
	if !o.CreatedAt.IsZero() {
		created = o.CreatedAt.UTC().Format(time.RFC3339)
	}
	return orderDTO{
		ID:            o.ID,
		Code:          o.Code,
		CustomerID:    o.CustomerID,
		CreatorID:     o.CreatorID,
		Status:        o.Status,
		OrderType:     o.OrderType,
		TotalAmount:   o.Total.String(),
		FinalAmount:   o.Final.String(),
		PaymentStatus: o.PaymentStatus,
		Note:          o.Note,
		CreatedAt:     created,
		Payment:       toPaymentDTO(o.Payment),
	}
}

func toOrderLineDTO(l domain.OrderLine) orderLineDTO {
	expiry := ""
	if l.HasExpiry {
		expiry = l.ExpiryDate.Format(dateLayout)
	}
	return orderLineDTO{
		ID:         l.ID,
		ProductID:  l.ProductID,
		Quantity:   l.Quantity.String(),
		UOM:        l.UOM,
		UnitPrice:  l.UnitPrice.String(),
		Discount:   l.Discount.String(),
		LineTotal:  l.LineTotal.String(),
		IsGift:     l.IsGift,
		BatchNo:    l.BatchNo,
		ExpiryDate: expiry,
		Note:       l.Note,
	}
}

// ---- Inputs/Outputs ----

type listOrdersInput struct {
	Cursor string `query:"cursor" doc:"Con trỏ trang (opaque); rỗng = trang đầu (đơn mới nhất trước)"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"100" doc:"Số đơn mỗi trang (mặc định 20)"`
	// 0 = không lọc; id bigint thật luôn > 0.
	CustomerID    int64  `query:"customer_id" minimum:"0" doc:"Lọc theo khách (0 = tất cả)"`
	Status        string `query:"status" doc:"Lọc trạng thái xử lý đơn; rỗng = tất cả"`
	PaymentStatus string `query:"payment_status" enum:"unpaid,partial,paid" doc:"Lọc trạng thái thanh toán; rỗng = tất cả"`
	FromDate      string `query:"from_date" doc:"Lọc đơn tạo từ ngày (YYYY-MM-DD); rỗng = không chặn"`
	ToDate        string `query:"to_date" doc:"Lọc đơn tạo đến hết ngày (YYYY-MM-DD); rỗng = không chặn"`
}

type listOrdersOutput struct {
	Body struct {
		Items      []orderDTO `json:"items"`
		NextCursor string     `json:"next_cursor" doc:"Con trỏ trang kế; rỗng = hết"`
	}
}

type getOrderInput struct {
	ID string `path:"id" doc:"ID đơn hàng (uuid)"`
}

type getOrderOutput struct {
	Body struct {
		Order orderDTO       `json:"order"`
		Lines []orderLineDTO `json:"lines"`
	}
}

// ---- Handlers ----

func registerListOrders(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "orders-list",
		Method:      http.MethodGet,
		Path:        "/v1/orders",
		Summary:     "Danh sách đơn hàng (keyset pagination + filter)",
		Tags:        []string{"orders"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listOrdersInput) (*listOrdersOutput, error) {
		from, err := parseDayStart(in.FromDate)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		to, err := parseDayEnd(in.ToDate)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		res, err := svc.ListOrders(ctx, app.ListOrdersQuery{
			Cursor:        in.Cursor,
			Limit:         in.Limit,
			CustomerID:    optID(in.CustomerID),
			Status:        in.Status,
			PaymentStatus: in.PaymentStatus,
			FromDate:      from,
			ToDate:        to,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listOrdersOutput{}
		out.Body.Items = make([]orderDTO, 0, len(res.Items))
		for _, o := range res.Items {
			out.Body.Items = append(out.Body.Items, toOrderDTO(o))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

func registerGetOrder(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "orders-get",
		Method:      http.MethodGet,
		Path:        "/v1/orders/{id}",
		Summary:     "Chi tiết đơn hàng + dòng hàng + đã thu/còn nợ (suy diễn)",
		Tags:        []string{"orders"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *getOrderInput) (*getOrderOutput, error) {
		d, err := svc.GetOrder(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &getOrderOutput{}
		out.Body.Order = toOrderDTO(d.Order)
		out.Body.Lines = make([]orderLineDTO, 0, len(d.Lines))
		for _, l := range d.Lines {
			out.Body.Lines = append(out.Body.Lines, toOrderLineDTO(l))
		}
		return out, nil
	})
}

// optID chuyển 0 → nil (không lọc), id > 0 → con trỏ tới id.
func optID(id int64) *int64 {
	if id <= 0 {
		return nil
	}
	return &id
}

// parseDayStart parse YYYY-MM-DD → unix nano tại 00:00:00 UTC (đầu ngày). Rỗng →
// 0 (không chặn). Sai định dạng → apperr.Validation.
func parseDayStart(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return 0, errBadDate("from_date")
	}
	return t.UTC().UnixNano(), nil
}

// parseDayEnd parse YYYY-MM-DD → unix nano tại 23:59:59.999999999 UTC (cuối ngày)
// để bao trùm trọn ngày khi lọc o.created_at <= to_date. Rỗng → 0.
func parseDayEnd(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return 0, errBadDate("to_date")
	}
	end := t.UTC().Add(24*time.Hour - time.Nanosecond)
	return end.UnixNano(), nil
}
