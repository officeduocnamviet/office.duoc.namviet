// Package http là ADAPTER vào của customers: DTO + handler Huma + đăng ký route.
// Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, và bảo vệ mọi route bằng authz.RequirePermissionHuma
// ("customers.read"). Tiền (công nợ, hạn mức) ra DTO dạng CHUỖI thập phân
// (money.String) — KHÔNG float — để FE không mất chính xác. Nằm dưới internal/
// nên module khác không import.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/customers/app"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Service là cổng use-case mà handler cần (interface để test handler bằng fake).
type Service interface {
	ListCustomers(ctx context.Context, q app.ListCustomersQuery) (app.ListCustomersResult, error)
	GetCustomer(ctx context.Context, id int64) (domain.Customer, error)
}

const permRead = "customers.read"

// Register đăng ký toàn bộ operation /v1/customers lên huma.API. Mọi route yêu
// cầu quyền customers.read (verify token + check perm qua
// authz.RequirePermissionHuma — 1 enforcement point).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListCustomers(api, svc, guard)
	registerGetCustomer(api, svc, guard)
}

// ---- DTO ----

// debtDTO trình bày công nợ: amount = con số đã chọn nguồn (ưu tiên live), kèm
// cả live/static + source để FE hiểu độ tin cậy. Tiền là chuỗi thập phân.
type debtDTO struct {
	Amount string `json:"amount" doc:"Công nợ dùng để hiển thị (chuỗi thập phân, không float)"`
	Live   string `json:"live" doc:"Công nợ tính LIVE từ đơn chưa tất toán"`
	Static string `json:"static" doc:"Công nợ cột tĩnh customers.current_debt (stale)"`
	Source string `json:"source" doc:"Nguồn của amount: live | static"`
}

type b2bDTO struct {
	TaxCode      string `json:"tax_code" doc:"Mã số thuế (MST)"`
	DebtLimit    string `json:"debt_limit" doc:"Hạn mức công nợ (chỉ hiển thị, KHÔNG enforce — credit-limit đang OFF)"`
	PaymentTerm  int    `json:"payment_term" doc:"Kỳ hạn thanh toán (số ngày); 0 = không khai"`
	SalesStaffID string `json:"sales_staff_id" doc:"Nhân viên phụ trách (uuid); rỗng = chưa gán"`
}

type customerDTO struct {
	ID      int64   `json:"id"`
	Code    string  `json:"customer_code"`
	Name    string  `json:"name"`
	Type    string  `json:"customer_type" doc:"B2C | B2B"`
	Phone   string  `json:"phone"`
	Email   string  `json:"email"`
	Address string  `json:"address"`
	Status  string  `json:"status"`
	Debt    debtDTO `json:"debt"`
	B2B     *b2bDTO `json:"b2b" doc:"Hồ sơ doanh nghiệp; null nếu khách B2C"`
}

func toDebtDTO(d domain.DebtSnapshot) debtDTO {
	return debtDTO{
		Amount: d.Amount.String(),
		Live:   d.Live.String(),
		Static: d.Static.String(),
		Source: string(d.Source),
	}
}

func toCustomerDTO(c domain.Customer) customerDTO {
	dto := customerDTO{
		ID:      c.ID,
		Code:    c.Code,
		Name:    c.Name,
		Type:    string(c.Type),
		Phone:   c.Phone,
		Email:   c.Email,
		Address: c.Address,
		Status:  c.Status,
		Debt:    toDebtDTO(c.Debt),
	}
	if c.B2B != nil {
		dto.B2B = &b2bDTO{
			TaxCode:      c.B2B.TaxCode,
			DebtLimit:    c.B2B.DebtLimit.String(),
			PaymentTerm:  c.B2B.PaymentTerm,
			SalesStaffID: c.B2B.SalesStaffID,
		}
	}
	return dto
}

// ---- Inputs/Outputs ----

type listCustomersInput struct {
	Cursor string `query:"cursor" doc:"Con trỏ trang (opaque); rỗng = trang đầu"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"100" doc:"Số khách mỗi trang (mặc định 20)"`
	Type   string `query:"customer_type" enum:"B2B,B2C" doc:"Lọc loại khách; rỗng = tất cả"`
	Status string `query:"status" doc:"Lọc trạng thái (vd active); rỗng = tất cả"`
	Q      string `query:"q" doc:"Tìm theo tên/SĐT/mã KH/MST"`
}

type listCustomersOutput struct {
	Body struct {
		Items      []customerDTO `json:"items"`
		NextCursor string        `json:"next_cursor" doc:"Con trỏ trang kế; rỗng = hết"`
	}
}

type getCustomerInput struct {
	ID int64 `path:"id" doc:"ID khách hàng (bigint)"`
}

type getCustomerOutput struct {
	Body struct {
		Customer customerDTO `json:"customer"`
	}
}

// ---- Handlers ----

func registerListCustomers(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "customers-list",
		Method:      http.MethodGet,
		Path:        "/v1/customers",
		Summary:     "Danh sách khách hàng (keyset pagination + filter)",
		Tags:        []string{"customers"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listCustomersInput) (*listCustomersOutput, error) {
		res, err := svc.ListCustomers(ctx, app.ListCustomersQuery{
			Cursor: in.Cursor,
			Limit:  in.Limit,
			Type:   in.Type,
			Status: in.Status,
			Query:  in.Q,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listCustomersOutput{}
		out.Body.Items = make([]customerDTO, 0, len(res.Items))
		for _, c := range res.Items {
			out.Body.Items = append(out.Body.Items, toCustomerDTO(c))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

func registerGetCustomer(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "customers-get",
		Method:      http.MethodGet,
		Path:        "/v1/customers/{id}",
		Summary:     "Chi tiết khách hàng + hồ sơ B2B + công nợ (nguồn live)",
		Tags:        []string{"customers"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *getCustomerInput) (*getCustomerOutput, error) {
		c, err := svc.GetCustomer(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &getCustomerOutput{}
		out.Body.Customer = toCustomerDTO(c)
		return out, nil
	})
}
