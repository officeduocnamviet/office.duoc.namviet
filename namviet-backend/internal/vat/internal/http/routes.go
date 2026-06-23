// Package http là ADAPTER vào của vat: DTO + handler Huma + đăng ký route ĐỌC.
// Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, bảo vệ mọi route bằng authz.RequirePermissionHuma
// ("vat.read"). Tiền (subtotal/vat/total/đơn giá) ra DTO dạng CHUỖI thập phân —
// KHÔNG float. CHỈ có route ĐỌC; phát hành HĐ là port nội bộ (app.IssuePort) cho
// orders — KHÔNG expose REST POST ở P5. Nằm dưới internal/.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
	"github.com/Maneva-AI/namviet-backend/internal/vat/app"
	"github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

// Service là cổng use-case ĐỌC mà handler cần (interface để test bằng fake).
type Service interface {
	ListInvoices(ctx context.Context, q app.ListInvoicesQuery) (app.ListInvoicesResult, error)
	GetInvoice(ctx context.Context, id string) (domain.IssuedInvoice, error)
}

const permRead = "vat.read"

// Register đăng ký các operation ĐỌC /v1/vat/* lên huma.API. Mọi route yêu cầu
// quyền vat.read (verify token + check perm qua RequirePermissionHuma).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListInvoices(api, svc, guard)
	registerGetInvoice(api, svc, guard)
}

// ---- DTO ----

const dateLayout = "2006-01-02"

type vatInvoiceDTO struct {
	ID              string `json:"id"`
	OrderCode       string `json:"order_code" doc:"= orders.code (mã đơn)"`
	CustomerTaxCode string `json:"customer_tax_code" doc:"MST khách hàng (bắt buộc)"`
	Serial          string `json:"serial" doc:"Ký hiệu hoá đơn"`
	InvoiceNo       int64  `json:"invoice_no" doc:"Số hoá đơn (gapless theo serial)"`
	IssueDate       string `json:"issue_date" doc:"Ngày phát hành (YYYY-MM-DD)"`
	Subtotal        string `json:"subtotal" doc:"Tổng tiền hàng trước thuế (chuỗi thập phân, không float)"`
	VATAmount       string `json:"vat_amount" doc:"Tổng thuế GTGT (chuỗi thập phân)"`
	Total           string `json:"total" doc:"Tổng thanh toán = subtotal + vat_amount"`
	Status          string `json:"status" doc:"draft | issued | cancelled"`
	CreatedAt       string `json:"created_at" doc:"Thời điểm tạo (RFC3339)"`
}

type vatLineDTO struct {
	LineNo      int32  `json:"line_no"`
	ProductID   int64  `json:"product_id"`
	Description string `json:"description"`
	Quantity    string `json:"quantity" doc:"Số lượng (chuỗi thập phân)"`
	UnitPrice   string `json:"unit_price" doc:"Đơn giá hoá đơn (chuỗi thập phân)"`
	VATRate     string `json:"vat_rate" doc:"Thuế suất dòng (vd 0.08)"`
	LineAmount  string `json:"line_amount" doc:"Giá trị dòng trước thuế"`
	LineVAT     string `json:"line_vat" doc:"Thuế GTGT dòng (đã làm tròn về đồng)"`
}

type vatInvoiceDetailDTO struct {
	vatInvoiceDTO
	Lines []vatLineDTO `json:"lines"`
}

func toInvoiceDTO(r domain.IssuedInvoice) vatInvoiceDTO {
	return vatInvoiceDTO{
		ID:              r.ID,
		OrderCode:       r.OrderCode,
		CustomerTaxCode: r.CustomerTaxCode,
		Serial:          r.Serial,
		InvoiceNo:       r.InvoiceNo,
		IssueDate:       r.IssueDate.Format(dateLayout),
		Subtotal:        r.Subtotal.String(),
		VATAmount:       r.VATAmount.String(),
		Total:           r.Total.String(),
		Status:          r.Status.String(),
		CreatedAt:       r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toInvoiceDetailDTO(r domain.IssuedInvoice) vatInvoiceDetailDTO {
	d := vatInvoiceDetailDTO{vatInvoiceDTO: toInvoiceDTO(r)}
	d.Lines = make([]vatLineDTO, 0, len(r.Lines))
	for _, l := range r.Lines {
		d.Lines = append(d.Lines, vatLineDTO{
			LineNo:      l.LineNo,
			ProductID:   l.ProductID,
			Description: l.Description,
			Quantity:    l.Quantity.String(),
			UnitPrice:   l.UnitPrice.String(),
			VATRate:     l.VATRate.String(),
			LineAmount:  l.LineAmount.String(),
			LineVAT:     l.LineVAT.String(),
		})
	}
	return d
}

// ---- Inputs/Outputs ----

type listInvoicesInput struct {
	Cursor    string `query:"cursor" doc:"Con trỏ trang (opaque); rỗng = trang đầu"`
	Limit     int32  `query:"limit" minimum:"1" maximum:"200" doc:"Số HĐ mỗi trang (mặc định 50)"`
	OrderCode string `query:"order_code" doc:"Lọc theo mã đơn; rỗng = mọi đơn"`
	Status    string `query:"status" enum:"draft,issued,cancelled" doc:"Lọc theo trạng thái; rỗng = mọi trạng thái"`
}

type listInvoicesOutput struct {
	Body struct {
		Items      []vatInvoiceDTO `json:"items"`
		NextCursor string          `json:"next_cursor" doc:"Con trỏ trang kế; rỗng = hết"`
	}
}

type getInvoiceInput struct {
	ID string `path:"id" doc:"ID hoá đơn (uuid)"`
}

type getInvoiceOutput struct {
	Body vatInvoiceDetailDTO
}

// ---- Handlers ----

func registerListInvoices(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "vat-list-invoices",
		Method:      http.MethodGet,
		Path:        "/v1/vat/invoices",
		Summary:     "Danh sách hoá đơn VAT (keyset pagination)",
		Tags:        []string{"vat"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listInvoicesInput) (*listInvoicesOutput, error) {
		res, err := svc.ListInvoices(ctx, app.ListInvoicesQuery{
			Cursor:    in.Cursor,
			Limit:     in.Limit,
			OrderCode: in.OrderCode,
			Status:    in.Status,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listInvoicesOutput{}
		out.Body.Items = make([]vatInvoiceDTO, 0, len(res.Items))
		for _, r := range res.Items {
			out.Body.Items = append(out.Body.Items, toInvoiceDTO(r))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

func registerGetInvoice(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "vat-get-invoice",
		Method:      http.MethodGet,
		Path:        "/v1/vat/invoices/{id}",
		Summary:     "Chi tiết một hoá đơn VAT kèm các dòng",
		Tags:        []string{"vat"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *getInvoiceInput) (*getInvoiceOutput, error) {
		rec, err := svc.GetInvoice(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		return &getInvoiceOutput{Body: toInvoiceDetailDTO(rec)}, nil
	})
}
