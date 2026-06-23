// routes_orchestration.go — ADAPTER GHI ORCHESTRATION của orders (P4b): POST giao
// hàng (ship), ghi phiếu thu (payments), bán lẻ tại quầy (pos/sales). Guard quyền
// "orders.write". Mỗi use-case = 1 transaction atomic ở tầng app (Orchestrator) —
// adapter này chỉ dịch HTTP <-> input/output. Tiền/lượng ra/vào DTO dạng CHUỖI
// thập phân — KHÔNG float.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	inventorydomain "github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	vatdomain "github.com/Maneva-AI/namviet-backend/internal/vat/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// OrchestrationService là cổng use-case ORCHESTRATION mà handler cần (interface để
// test bằng fake + tránh import cycle: http KHÔNG import package orders). Facade
// *orders.Service thoả structurally.
type OrchestrationService interface {
	ShipOrder(ctx context.Context, in app.ShipInput) (app.ShippedOrder, error)
	RecordPayment(ctx context.Context, in app.RecordPaymentInput) (app.PaymentResult, error)
	CreatePosSale(ctx context.Context, in app.PosSaleInput) (app.PosSaleResult, error)
	RecordLumpSumPayment(ctx context.Context, in app.LumpSumInput) (app.LumpSumResult, error)
}

// RegisterOrchestration đăng ký các operation GHI cross-module (ship/payments/pos/
// lump-sum) lên huma.API. Mọi route yêu cầu quyền orders.write (1 enforcement point).
func RegisterOrchestration(api huma.API, svc OrchestrationService, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permWrite)
	registerShipOrder(api, svc, guard)
	registerRecordPayment(api, svc, guard)
	registerPosSale(api, svc, guard)
	registerLumpSumPayment(api, svc, guard)
}

// ---- DTO chung ----

// orchVATLineDTO là thuế suất + (tuỳ chọn) đơn giá HĐ cho một dòng khi xuất HĐ VAT.
type orchVATLineDTO struct {
	VATRate     string `json:"vat_rate" doc:"Thuế suất dòng (vd \"0.08\" = 8%)"`
	UnitPrice   string `json:"unit_price,omitempty" doc:"Đơn giá HĐ (chuỗi; bỏ trống = dùng đơn giá đơn)"`
	Description string `json:"description,omitempty" doc:"Diễn giải dòng trên HĐ"`
}

// invoiceSummaryDTO là tóm tắt HĐ VAT đã phát hành (cho FE hiển thị/đối soát).
type invoiceSummaryDTO struct {
	ID        string `json:"id"`
	Serial    string `json:"serial" doc:"Ký hiệu HĐ"`
	InvoiceNo int64  `json:"invoice_no" doc:"Số HĐ (gapless theo serial)"`
	Subtotal  string `json:"subtotal" doc:"Tiền hàng chưa VAT (chuỗi thập phân)"`
	VATAmount string `json:"vat_amount" doc:"Tiền VAT (chuỗi thập phân)"`
	Total     string `json:"total" doc:"Tổng tiền HĐ (chuỗi thập phân)"`
	Status    string `json:"status" doc:"draft|issued|cancelled"`
}

// consumedBatchDTO là một lô đã trừ khi giao hàng (cho FE đối soát + tính COGS).
type consumedBatchDTO struct {
	BatchID          int64  `json:"batch_id"`
	InventoryBatchID int64  `json:"inventory_batch_id"`
	Quantity         string `json:"quantity" doc:"Số lượng trừ từ lô (chuỗi thập phân)"`
	InboundPrice     string `json:"inbound_price" doc:"Giá vốn nhập per-unit của lô (chuỗi thập phân)"`
}

// financePaymentDTO là tóm tắt phiếu THU đã ghi (đặt tên khác paymentDTO ở routes.go
// — cái kia là ảnh chụp đã-thu/còn-nợ suy diễn).
type financePaymentDTO struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Amount   string `json:"amount" doc:"Số tiền thu (chuỗi thập phân)"`
	Status   string `json:"status" doc:"pending|completed..."`
	BookType string `json:"book_type" doc:"INTERNAL|TAX|BOTH"`
}

// ---- Ship (giao hàng) ----

type shipOrderInput struct {
	ID   string `path:"id" doc:"ID đơn hàng (uuid)"`
	Body struct {
		WarehouseID     int64        `json:"warehouse_id" minimum:"1" doc:"Kho xuất hàng"`
		CustomerTaxCode string       `json:"customer_tax_code" doc:"MST khách (bắt buộc cho HĐ VAT B2B)"`
		Serial          string       `json:"serial" doc:"Ký hiệu HĐ (dải đăng ký với CQT)"`
		MauSo           string       `json:"mau_so,omitempty" doc:"Mẫu số HĐ (NĐ123)"`
		VATLines        []orchVATLineDTO `json:"vat_lines,omitempty" doc:"Thuế suất/đơn giá HĐ theo THỨ TỰ dòng đơn; nên truyền đủ cho B2B"`
	}
}

type shipOrderOutput struct {
	Body struct {
		Order    orderDTO           `json:"order"`
		Invoice  invoiceSummaryDTO  `json:"invoice"`
		Consumed []consumedBatchDTO `json:"consumed"`
		COGS     string             `json:"cogs" doc:"Tổng giá vốn đã xuất (chuỗi thập phân)"`
	}
}

func registerShipOrder(api huma.API, svc OrchestrationService, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "orders-ship",
		Method:        http.MethodPost,
		Path:          "/v1/orders/{id}/ship",
		Summary:       "Giao hàng (CONFIRMED→SHIPPING): trừ kho FEFO + xuất HĐ VAT + post sổ kép (atomic)",
		Tags:          []string{"orders"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *shipOrderInput) (*shipOrderOutput, error) {
		vatLines, err := toVATLines(in.Body.VATLines)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		res, err := svc.ShipOrder(ctx, app.ShipInput{
			OrderID:         in.ID,
			WarehouseID:     in.Body.WarehouseID,
			CustomerTaxCode: in.Body.CustomerTaxCode,
			Serial:          in.Body.Serial,
			MauSo:           in.Body.MauSo,
			VATLines:        vatLines,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &shipOrderOutput{}
		out.Body.Order = toOrderDTO(res.Order)
		out.Body.Invoice = toInvoiceSummary(res.Invoice)
		out.Body.Consumed = toConsumedDTOs(res.Consumed)
		out.Body.COGS = res.COGS.String()
		return out, nil
	})
}

// ---- RecordPayment (ghi phiếu thu) ----

type recordPaymentInput struct {
	ID      string `path:"id" doc:"ID đơn hàng (uuid)"`
	IdemKey string `header:"Idempotency-Key" doc:"Khoá idempotency phiếu thủ công (bắt buộc nếu không có bank_ref)"`
	Body    struct {
		Amount        string `json:"amount" doc:"Số tiền thu (chuỗi thập phân, > 0)"`
		FundAccountID int64  `json:"fund_account_id" minimum:"1" doc:"Quỹ/tài khoản nhận tiền"`
		FundIsBank    bool   `json:"fund_is_bank,omitempty" doc:"true = ngân hàng (112); false = tiền mặt (111)"`
		BookType      string `json:"book_type,omitempty" enum:"INTERNAL,TAX,BOTH" doc:"Sổ ghi phiếu; bỏ trống = BOTH"`
		BankRef       string `json:"bank_ref,omitempty" doc:"Mã giao dịch ngân hàng (phiếu webhook, chống trùng)"`
		CreatedBy     string `json:"created_by,omitempty" doc:"uuid người lập phiếu"`
	}
}

type recordPaymentOutput struct {
	Body struct {
		Payment       financePaymentDTO `json:"payment"`
		PaymentStatus string            `json:"payment_status" doc:"unpaid|partial|paid (sau khi thu)"`
		Order         orderDTO          `json:"order"`
	}
}

func registerRecordPayment(api huma.API, svc OrchestrationService, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "orders-record-payment",
		Method:        http.MethodPost,
		Path:          "/v1/orders/{id}/payments",
		Summary:       "Ghi phiếu THU cho đơn + post PAYMENT_IN + cập nhật payment_status (atomic)",
		Tags:          []string{"orders"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *recordPaymentInput) (*recordPaymentOutput, error) {
		amount, err := money.FromString(in.Body.Amount)
		if err != nil {
			return nil, humax.FromAppErr(apperr.Validation("amount không hợp lệ: " + in.Body.Amount))
		}
		res, err := svc.RecordPayment(ctx, app.RecordPaymentInput{
			OrderID:       in.ID,
			Amount:        amount,
			FundAccountID: in.Body.FundAccountID,
			FundIsBank:    in.Body.FundIsBank,
			BookType:      financedomain.BookType(in.Body.BookType),
			BankRef:       optStr(in.Body.BankRef),
			IdemKey:       in.IdemKey,
			CreatedBy:     optStr(in.Body.CreatedBy),
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &recordPaymentOutput{}
		out.Body.Payment = toFinancePaymentDTO(res.Payment)
		out.Body.PaymentStatus = res.PaymentStatus
		out.Body.Order = toOrderDTO(res.Order)
		return out, nil
	})
}

// ---- POS sale (bán lẻ tại quầy) ----

type posSaleInput struct {
	IdemKey string `header:"Idempotency-Key" doc:"Khoá idempotency tạo đơn POS (1 key → 1 đơn)"`
	Body    struct {
		CustomerID      *int64               `json:"customer_id,omitempty" doc:"ID khách; bỏ trống nếu khách vãng lai"`
		CreatorID       string               `json:"creator_id,omitempty" doc:"Người bán (uuid)"`
		Note            string               `json:"note,omitempty"`
		Lines           []createOrderLineDTO `json:"lines" minItems:"1" doc:"Các dòng hàng (ít nhất 1)"`
		WarehouseID     int64                `json:"warehouse_id" minimum:"1" doc:"Kho xuất"`
		FundAccountID   int64                `json:"fund_account_id" minimum:"1" doc:"Quỹ thu tiền"`
		FundIsBank      bool                 `json:"fund_is_bank,omitempty"`
		IssueInvoice    bool                 `json:"issue_invoice,omitempty" doc:"true = xuất HĐ VAT (cần MST + serial)"`
		CustomerTaxCode string               `json:"customer_tax_code,omitempty" doc:"MST (bắt buộc nếu issue_invoice)"`
		Serial          string               `json:"serial,omitempty"`
		MauSo           string               `json:"mau_so,omitempty"`
		VATLines        []orchVATLineDTO         `json:"vat_lines,omitempty" doc:"Theo THỨ TỰ dòng khi xuất HĐ"`
	}
}

type posSaleOutput struct {
	Body struct {
		Order    orderDTO           `json:"order"`
		Lines    []orderLineDTO     `json:"lines"`
		Payment  financePaymentDTO  `json:"payment"`
		Invoice  *invoiceSummaryDTO `json:"invoice,omitempty" doc:"HĐ VAT nếu có xuất; null nếu không"`
		Consumed []consumedBatchDTO `json:"consumed"`
		COGS     string             `json:"cogs"`
	}
}

func registerPosSale(api huma.API, svc OrchestrationService, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "pos-create-sale",
		Method:        http.MethodPost,
		Path:          "/v1/pos/sales",
		Summary:       "Bán lẻ tại quầy (atomic): tạo đơn + trừ kho FEFO + thu tiền + (HĐ) + post sổ → COMPLETED",
		Tags:          []string{"orders"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *posSaleInput) (*posSaleOutput, error) {
		lines, err := toDraftLines(in.Body.Lines)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		vatLines, err := toVATLines(in.Body.VATLines)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		res, err := svc.CreatePosSale(ctx, app.PosSaleInput{
			CustomerID:      in.Body.CustomerID,
			CreatorID:       in.Body.CreatorID,
			Note:            in.Body.Note,
			Lines:           lines,
			WarehouseID:     in.Body.WarehouseID,
			FundAccountID:   in.Body.FundAccountID,
			FundIsBank:      in.Body.FundIsBank,
			IdemKey:         in.IdemKey,
			IssueInvoice:    in.Body.IssueInvoice,
			CustomerTaxCode: in.Body.CustomerTaxCode,
			Serial:          in.Body.Serial,
			MauSo:           in.Body.MauSo,
			VATLines:        vatLines,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &posSaleOutput{}
		out.Body.Order = toOrderDTO(res.Order)
		out.Body.Lines = make([]orderLineDTO, 0, len(res.Lines))
		for _, l := range res.Lines {
			out.Body.Lines = append(out.Body.Lines, toOrderLineDTO(l))
		}
		out.Body.Payment = toFinancePaymentDTO(res.Payment)
		if res.Invoice != nil {
			inv := toInvoiceSummary(*res.Invoice)
			out.Body.Invoice = &inv
		}
		out.Body.Consumed = toConsumedDTOs(res.Consumed)
		out.Body.COGS = res.COGS.String()
		return out, nil
	})
}

// ---- Lump-sum payment (thu 1 cục phân bổ nhiều đơn — mục 55) ----

type lumpSumInput struct {
	CustomerID int64  `path:"id" doc:"ID khách hàng (bigint)"`
	IdemKey    string `header:"Idempotency-Key" doc:"Khoá idempotency (phiếu thủ công); bắt buộc nếu không có bank_ref"`
	Body       struct {
		Amount        string `json:"amount" doc:"Số tiền thu 1 cục (chuỗi thập phân, > 0)"`
		FundAccountID int64  `json:"fund_account_id" minimum:"1" doc:"Quỹ/tài khoản nhận tiền"`
		FundIsBank    bool   `json:"fund_is_bank,omitempty"`
		BookType      string `json:"book_type,omitempty" enum:"INTERNAL,TAX,BOTH" doc:"Sổ ghi phiếu; bỏ trống = BOTH"`
		Collected     bool   `json:"collected,omitempty" doc:"true = NV đã thu, chưa nộp quỹ (pending); false = vào quỹ ngay"`
		BankRef       string `json:"bank_ref,omitempty"`
		CreatedBy     string `json:"created_by,omitempty"`
	}
}

type allocationLineDTO struct {
	OrderID   string `json:"order_id"`
	OrderCode string `json:"order_code"`
	Amount    string `json:"amount" doc:"Số tiền phân bổ cho đơn (chuỗi thập phân)"`
}

type lumpSumOutput struct {
	Body struct {
		Payment     financePaymentDTO   `json:"payment"`
		Allocations []allocationLineDTO `json:"allocations" doc:"Phân bổ cho từng đơn (cũ nhất trước)"`
		Leftover    string              `json:"leftover" doc:"Thu thừa chưa phân bổ (credit khách); '0' nếu phân bổ hết"`
	}
}

func registerLumpSumPayment(api huma.API, svc OrchestrationService, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "customers-lump-payment",
		Method:        http.MethodPost,
		Path:          "/v1/customers/{id}/payments",
		Summary:       "Thu 1 cục từ khách → phân bổ cho các đơn chưa tất toán (cũ nhất trước, atomic)",
		Tags:          []string{"orders"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *lumpSumInput) (*lumpSumOutput, error) {
		amount, err := money.FromString(in.Body.Amount)
		if err != nil {
			return nil, humax.FromAppErr(apperr.Validation("amount không hợp lệ: " + in.Body.Amount))
		}
		res, err := svc.RecordLumpSumPayment(ctx, app.LumpSumInput{
			CustomerID:    in.CustomerID,
			Amount:        amount,
			FundAccountID: in.Body.FundAccountID,
			FundIsBank:    in.Body.FundIsBank,
			BookType:      financedomain.BookType(in.Body.BookType),
			Collected:     in.Body.Collected,
			BankRef:       optStr(in.Body.BankRef),
			IdemKey:       in.IdemKey,
			CreatedBy:     optStr(in.Body.CreatedBy),
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &lumpSumOutput{}
		out.Body.Payment = toFinancePaymentDTO(res.Payment)
		out.Body.Allocations = make([]allocationLineDTO, 0, len(res.Allocations))
		for _, a := range res.Allocations {
			out.Body.Allocations = append(out.Body.Allocations, allocationLineDTO{
				OrderID: a.OrderID, OrderCode: a.OrderCode, Amount: a.Amount.String(),
			})
		}
		out.Body.Leftover = res.Leftover.String()
		return out, nil
	})
}

// ---- mapping helpers ----

// toVATLines chuyển DTO (vat_rate/unit_price chuỗi) sang app.VATLine (decimal/money).
// vat_rate rỗng → 0%; unit_price rỗng → HasPrice=false (dùng đơn giá đơn). Lỗi parse
// → apperr.Validation (422).
func toVATLines(in []orchVATLineDTO) ([]app.VATLine, error) {
	out := make([]app.VATLine, 0, len(in))
	for i, l := range in {
		rate := decimal.Zero
		if l.VATRate != "" {
			r, err := decimal.NewFromString(l.VATRate)
			if err != nil {
				return nil, apperr.Validation("vat_rate dòng không hợp lệ: " + l.VATRate)
			}
			rate = r
		}
		vl := app.VATLine{VATRate: rate, Description: l.Description}
		if l.UnitPrice != "" {
			price, err := money.FromString(l.UnitPrice)
			if err != nil {
				return nil, apperr.Validation("unit_price HĐ dòng không hợp lệ: " + l.UnitPrice)
			}
			vl.UnitPrice = price
			vl.HasPrice = true
		}
		_ = i
		out = append(out, vl)
	}
	return out, nil
}

func toInvoiceSummary(inv vatdomain.IssuedInvoice) invoiceSummaryDTO {
	return invoiceSummaryDTO{
		ID:        inv.ID,
		Serial:    inv.Serial,
		InvoiceNo: inv.InvoiceNo,
		Subtotal:  inv.Subtotal.String(),
		VATAmount: inv.VATAmount.String(),
		Total:     inv.Total.String(),
		Status:    inv.Status.String(),
	}
}

func toConsumedDTOs(in []inventorydomain.ConsumedBatch) []consumedBatchDTO {
	out := make([]consumedBatchDTO, 0, len(in))
	for _, b := range in {
		out = append(out, consumedBatchDTO{
			BatchID:          b.BatchID,
			InventoryBatchID: b.InventoryBatchID,
			Quantity:         b.Quantity.String(),
			InboundPrice:     b.InboundPrice.String(),
		})
	}
	return out
}

func toFinancePaymentDTO(p financedomain.Payment) financePaymentDTO {
	return financePaymentDTO{
		ID:       p.ID,
		Code:     p.Code,
		Amount:   p.Amount.String(),
		Status:   p.Status,
		BookType: p.BookType.String(),
	}
}

// optStr trả nil nếu chuỗi rỗng, ngược lại con trỏ tới chuỗi (cho trường optional).
func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
