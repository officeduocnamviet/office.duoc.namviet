// Package http là ADAPTER vào của finance: DTO + handler Huma + đăng ký route
// ĐỌC. Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, bảo vệ route bằng authz.RequirePermissionHuma ("finance.read").
// Tiền amount ra DTO dạng CHUỖI thập phân — KHÔNG float. CHỈ có route ĐỌC; ghi
// phiếu thu là port nội bộ (app.RecordPort) cho orders/POS — KHÔNG expose REST
// POST. Nằm dưới internal/.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Service là cổng use-case ĐỌC mà handler cần (interface để test bằng fake).
type Service interface {
	ListOrderPayments(ctx context.Context, orderCode string, limit int32) ([]domain.Payment, error)
}

// WriteService là cổng use-case GHI cho HTTP: thủ quỹ "Xác nhận đã thu" (thanh toán
// 2 bước). Tách interface để test bằng fake.
type WriteService interface {
	ConfirmReceiptInOwnTx(ctx context.Context, paymentID int64) (bool, error)
}

const (
	permRead  = "finance.read"
	permWrite = "finance.write"
)

// Register đăng ký các operation ĐỌC /v1/finance/* lên huma.API. Mọi route yêu cầu
// quyền finance.read (verify token + check perm qua RequirePermissionHuma).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListTransactions(api, svc, guard)
}

// RegisterWrite đăng ký operation GHI /v1/finance/* (thủ quỹ xác nhận thu) — guard
// quyền finance.write (khác đường đọc finance.read).
func RegisterWrite(api huma.API, svc WriteService, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permWrite)
	registerConfirmReceipt(api, svc, guard)
}

// ---- DTO ----

type financeTransactionDTO struct {
	ID            int64  `json:"id"`
	Code          string `json:"code"`
	Flow          string `json:"flow" doc:"Dòng tiền: 'in' (thu) | 'out' (chi)"`
	BusinessType  string `json:"business_type"`
	Amount        string `json:"amount" doc:"Số tiền (chuỗi thập phân, không float)"`
	FundAccountID int64  `json:"fund_account_id"`
	RefType       string `json:"ref_type"`
	RefID         string `json:"ref_id" doc:"= orders.code khi ref_type='order'"`
	Status        string `json:"status" doc:"'completed' = đã thực vào quỹ"`
	BookType      string `json:"book_type" doc:"INTERNAL | TAX | BOTH"`
	BankRef       string `json:"bank_reference_id,omitempty"`
	Description   string `json:"description,omitempty"`
}

func toPaymentDTO(p domain.Payment) financeTransactionDTO {
	d := financeTransactionDTO{
		ID:            p.ID,
		Code:          p.Code,
		Flow:          p.Flow,
		BusinessType:  p.BusinessType,
		Amount:        p.Amount.String(),
		FundAccountID: p.FundAccountID,
		RefType:       p.RefType,
		RefID:         p.RefID,
		Status:        p.Status,
		BookType:      string(p.BookType),
	}
	if p.BankRef != nil {
		d.BankRef = *p.BankRef
	}
	if p.Description != nil {
		d.Description = *p.Description
	}
	return d
}

// ---- Inputs/Outputs ----

type listTransactionsInput struct {
	RefType string `query:"ref_type" enum:"order" doc:"Loại chứng từ gốc (hiện chỉ 'order')"`
	RefID   string `query:"ref_id" required:"true" doc:"Mã chứng từ gốc — = orders.code khi ref_type='order'"`
	Limit   int32  `query:"limit" minimum:"1" maximum:"200" doc:"Số phiếu mỗi trang (mặc định 50)"`
}

type listTransactionsOutput struct {
	Body struct {
		Items []financeTransactionDTO `json:"items"`
	}
}

// ---- Handlers ----

func registerListTransactions(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "finance-list-transactions",
		Method:      http.MethodGet,
		Path:        "/v1/finance/transactions",
		Summary:     "Danh sách phiếu thu/chi của một chứng từ gốc (đơn hàng)",
		Tags:        []string{"finance"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listTransactionsInput) (*listTransactionsOutput, error) {
		// Hiện chỉ hỗ trợ tra theo đơn (ref_type='order'); ref_type rỗng coi như 'order'.
		items, err := svc.ListOrderPayments(ctx, in.RefID, in.Limit)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listTransactionsOutput{}
		out.Body.Items = make([]financeTransactionDTO, 0, len(items))
		for _, p := range items {
			out.Body.Items = append(out.Body.Items, toPaymentDTO(p))
		}
		return out, nil
	})
}

type confirmReceiptInput struct {
	ID int64 `path:"id" doc:"ID phiếu thu (finance_transactions.id, bigint)"`
}

type confirmReceiptOutput struct {
	Body struct {
		Confirmed bool `json:"confirmed" doc:"true = vừa xác nhận vào quỹ (pending→completed); false = đã completed trước đó (idempotent)"`
	}
}

func registerConfirmReceipt(api huma.API, svc WriteService, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "finance-confirm-receipt",
		Method:        http.MethodPost,
		Path:          "/v1/finance/receipts/{id}/confirm",
		Summary:       "Thủ quỹ xác nhận đã thu (pending→completed, tiền vào quỹ) — thanh toán 2 bước",
		Tags:          []string{"finance"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *confirmReceiptInput) (*confirmReceiptOutput, error) {
		confirmed, err := svc.ConfirmReceiptInOwnTx(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &confirmReceiptOutput{}
		out.Body.Confirmed = confirmed
		return out, nil
	})
}
