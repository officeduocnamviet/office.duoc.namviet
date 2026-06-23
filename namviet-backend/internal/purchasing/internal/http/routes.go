// Package http là ADAPTER vào của purchasing (mục 54): DTO + handler Huma + đăng ký
// route. Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr. Guard quyền: ĐỌC (purchasing.read), GHI (purchasing.write) — 1
// enforcement point. Tiền/lượng vào/ra DTO dạng CHUỖI thập phân — KHÔNG float. Mỗi
// use-case GHI = 1 transaction atomic ở tầng app. Tên DTO prefix "po"/"purchase" để
// UNIQUE toàn cục (Huma panic nếu trùng schema giữa các package http).
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/app"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

const (
	permRead  = "purchasing.read"
	permWrite = "purchasing.write"
)

// Service là cổng use-case mà handler cần (interface để test bằng fake + tránh
// import cycle). Facade *purchasing.Service thoả structurally.
type Service interface {
	CreatePO(ctx context.Context, in app.CreatePOInput) (app.CreatedPO, bool, error)
	ConfirmPO(ctx context.Context, poID string) (domain.PurchaseOrder, error)
	CancelPO(ctx context.Context, poID string) (domain.PurchaseOrder, error)
	ReceivePO(ctx context.Context, in app.ReceiveInput) (app.ReceivedPO, error)
	PaySupplier(ctx context.Context, in app.PaySupplierInput) (app.PaySupplierResult, error)
	GetPO(ctx context.Context, poID string) (app.CreatedPO, error)
	ListPOs(ctx context.Context, q app.ListPOsQuery) (app.ListPOsResult, error)
}

// Register mount toàn bộ operation /v1/purchase-orders* lên huma.API: GET (đọc,
// guard purchasing.read) + POST (ghi, guard purchasing.write).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	read := authz.RequirePermissionHuma(api, verifier, permRead)
	write := authz.RequirePermissionHuma(api, verifier, permWrite)
	registerListPOs(api, svc, read)
	registerGetPO(api, svc, read)
	registerCreatePO(api, svc, write)
	registerConfirmPO(api, svc, write)
	registerCancelPO(api, svc, write)
	registerReceivePO(api, svc, write)
	registerPaySupplier(api, svc, write)
}

// ---- DTO ----

type poDTO struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	SupplierID   *int64 `json:"supplier_id,omitempty"`
	SupplierName string `json:"supplier_name,omitempty"`
	Status       string `json:"status" doc:"draft|ordered|received|paid|cancelled"`
	TotalAmount  string `json:"total_amount" doc:"Tổng tiền hàng ex-VAT (chuỗi thập phân)"`
	VATAmount    string `json:"vat_amount" doc:"Tổng VAT đầu vào (chuỗi thập phân)"`
	Note         string `json:"note,omitempty"`
	LockVersion  int32  `json:"lock_version"`
}

type poLineDTO struct {
	ID                string `json:"id"`
	LineNo            int    `json:"line_no"`
	ProductID         int64  `json:"product_id"`
	Quantity          string `json:"quantity" doc:"Số lượng nhập (chuỗi thập phân)"`
	UnitCost          string `json:"unit_cost" doc:"Đơn giá nhập per-unit (chuỗi thập phân)"`
	VATRate           string `json:"vat_rate" doc:"Thuế suất VAT dòng (vd \"0.08\")"`
	BatchCode         string `json:"batch_code,omitempty"`
	ExpiryDate        string `json:"expiry_date,omitempty" doc:"YYYY-MM-DD"`
	ManufacturingDate string `json:"manufacturing_date,omitempty" doc:"YYYY-MM-DD"`
	LineTotal         string `json:"line_total" doc:"Thành tiền dòng ex-VAT (chuỗi thập phân)"`
}

type poDetailOutput struct {
	Body struct {
		PurchaseOrder poDTO       `json:"purchase_order"`
		Lines         []poLineDTO `json:"lines"`
	}
}

type poActionInput struct {
	ID string `path:"id" doc:"ID đơn mua (uuid)"`
}

type poActionOutput struct {
	Body struct {
		PurchaseOrder poDTO `json:"purchase_order"`
	}
}

func toPODTO(p domain.PurchaseOrder) poDTO {
	return poDTO{
		ID:           p.ID,
		Code:         p.Code,
		SupplierID:   p.SupplierID,
		SupplierName: p.SupplierName,
		Status:       p.Status,
		TotalAmount:  p.TotalAmount.String(),
		VATAmount:    p.VATAmount.String(),
		Note:         p.Note,
		LockVersion:  p.LockVersion,
	}
}

func toPOLineDTO(l domain.PurchaseLine) poLineDTO {
	d := poLineDTO{
		ID:        l.ID,
		LineNo:    l.LineNo,
		ProductID: l.ProductID,
		Quantity:  l.Quantity.String(),
		UnitCost:  l.UnitCost.String(),
		VATRate:   l.VATRate.String(),
		BatchCode: l.BatchCode,
		LineTotal: l.LineTotal.String(),
	}
	if l.ExpiryDate != nil {
		d.ExpiryDate = l.ExpiryDate.Format("2006-01-02")
	}
	if l.ManufacturingDate != nil {
		d.ManufacturingDate = l.ManufacturingDate.Format("2006-01-02")
	}
	return d
}

// ---- Create PO ----

type poCreateLineDTO struct {
	ProductID         int64  `json:"product_id" minimum:"1" doc:"ID sản phẩm (bigint)"`
	Quantity          string `json:"quantity" doc:"Số lượng nhập (chuỗi thập phân, > 0)"`
	UnitCost          string `json:"unit_cost" doc:"Đơn giá nhập per-unit (chuỗi thập phân, >= 0)"`
	VATRate           string `json:"vat_rate,omitempty" doc:"Thuế suất VAT (vd \"0.08\"; bỏ trống = 0)"`
	BatchCode         string `json:"batch_code,omitempty"`
	ExpiryDate        string `json:"expiry_date,omitempty" doc:"YYYY-MM-DD"`
	ManufacturingDate string `json:"manufacturing_date,omitempty" doc:"YYYY-MM-DD"`
}

type poCreateInput struct {
	IdemKey string `header:"Idempotency-Key" doc:"Khoá idempotency tạo PO (1 key → 1 PO); rỗng = không idempotent"`
	Body    struct {
		SupplierID   *int64            `json:"supplier_id,omitempty" doc:"ID NCC (bigint, KHÔNG FK); bỏ trống nếu chưa có"`
		SupplierName string            `json:"supplier_name,omitempty" doc:"Tên NCC (lưu thẳng)"`
		Note         string            `json:"note,omitempty"`
		Lines        []poCreateLineDTO `json:"lines" minItems:"1" doc:"Các dòng hàng (ít nhất 1)"`
	}
}

func registerCreatePO(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "purchasing-create-po",
		Method:        http.MethodPost,
		Path:          "/v1/purchase-orders",
		Summary:       "Tạo đơn mua hàng (PO, status draft)",
		Tags:          []string{"purchasing"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *poCreateInput) (*poDetailOutput, error) {
		lines, err := toDraftLines(in.Body.Lines)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		res, _, err := svc.CreatePO(ctx, app.CreatePOInput{
			SupplierID:   in.Body.SupplierID,
			SupplierName: in.Body.SupplierName,
			Note:         in.Body.Note,
			Lines:        lines,
			IdemKey:      in.IdemKey,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		return toDetailOutput(res), nil
	})
}

// ---- Confirm / Cancel PO ----

func registerConfirmPO(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "purchasing-confirm-po",
		Method:        http.MethodPost,
		Path:          "/v1/purchase-orders/{id}/confirm",
		Summary:       "Xác nhận/đặt hàng (draft→ordered)",
		Tags:          []string{"purchasing"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *poActionInput) (*poActionOutput, error) {
		po, err := svc.ConfirmPO(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &poActionOutput{}
		out.Body.PurchaseOrder = toPODTO(po)
		return out, nil
	})
}

func registerCancelPO(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "purchasing-cancel-po",
		Method:        http.MethodPost,
		Path:          "/v1/purchase-orders/{id}/cancel",
		Summary:       "Huỷ đơn mua (draft/ordered→cancelled)",
		Tags:          []string{"purchasing"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *poActionInput) (*poActionOutput, error) {
		po, err := svc.CancelPO(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &poActionOutput{}
		out.Body.PurchaseOrder = toPODTO(po)
		return out, nil
	})
}

// ---- Receive PO (nhập kho + post sổ) ----

type poReceiveInput struct {
	ID   string `path:"id" doc:"ID đơn mua (uuid)"`
	Body struct {
		WarehouseID      int64  `json:"warehouse_id" minimum:"1" doc:"Kho nhập hàng"`
		HasInvoice       bool   `json:"has_invoice,omitempty" doc:"true = có HĐ VAT mua → post thêm sổ TAX"`
		TaxInventoryCost string `json:"tax_inventory_cost,omitempty" doc:"Giá vốn ex-VAT theo HĐ mua (chuỗi; bỏ trống = giá thực)"`
		TaxVATAmount     string `json:"tax_vat_amount,omitempty" doc:"VAT đầu vào theo HĐ mua (chuỗi; bỏ trống = VAT thực)"`
	}
}

type poReceiveOutput struct {
	Body struct {
		PurchaseOrder poDTO   `json:"purchase_order"`
		CreatedBatch  []int64 `json:"created_batch" doc:"batchID các lô vừa tạo"`
		InventoryCost string  `json:"inventory_cost" doc:"Tổng giá vốn nhập ex-VAT (chuỗi)"`
		VATAmount     string  `json:"vat_amount" doc:"Tổng VAT đầu vào (chuỗi)"`
	}
}

func registerReceivePO(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "purchasing-receive-po",
		Method:        http.MethodPost,
		Path:          "/v1/purchase-orders/{id}/receive",
		Summary:       "Nhận hàng (ordered→received): nhập kho + post sổ Dr 1561+133/Cr 331 (atomic)",
		Tags:          []string{"purchasing"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *poReceiveInput) (*poReceiveOutput, error) {
		taxInv, err := money.FromString(in.Body.TaxInventoryCost)
		if err != nil {
			return nil, humax.FromAppErr(apperr.Validation("tax_inventory_cost không hợp lệ: " + in.Body.TaxInventoryCost))
		}
		taxVAT, err := money.FromString(in.Body.TaxVATAmount)
		if err != nil {
			return nil, humax.FromAppErr(apperr.Validation("tax_vat_amount không hợp lệ: " + in.Body.TaxVATAmount))
		}
		res, err := svc.ReceivePO(ctx, app.ReceiveInput{
			POID:             in.ID,
			WarehouseID:      in.Body.WarehouseID,
			HasInvoice:       in.Body.HasInvoice,
			TaxInventoryCost: taxInv,
			TaxVATAmount:     taxVAT,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &poReceiveOutput{}
		out.Body.PurchaseOrder = toPODTO(res.PO)
		out.Body.CreatedBatch = res.CreatedBatch
		out.Body.InventoryCost = res.InventoryCost.String()
		out.Body.VATAmount = res.VATAmount.String()
		return out, nil
	})
}

// ---- Pay supplier (chi NCC + post sổ) ----

type poPayInput struct {
	ID      string `path:"id" doc:"ID đơn mua (uuid)"`
	IdemKey string `header:"Idempotency-Key" doc:"Khoá idempotency phiếu chi thủ công (bắt buộc nếu không có bank_ref)"`
	Body    struct {
		Amount        string `json:"amount" doc:"Số tiền chi (chuỗi thập phân, > 0)"`
		FundAccountID int64  `json:"fund_account_id" minimum:"1" doc:"Quỹ/tài khoản XUẤT tiền"`
		FundIsBank    bool   `json:"fund_is_bank,omitempty" doc:"true = ngân hàng (112); false = tiền mặt (111)"`
		BookType      string `json:"book_type,omitempty" enum:"INTERNAL,TAX,BOTH" doc:"Sổ ghi phiếu; bỏ trống = BOTH"`
		BankRef       string `json:"bank_ref,omitempty" doc:"Mã giao dịch ngân hàng (phiếu chi tự động, chống trùng)"`
		CreatedBy     string `json:"created_by,omitempty" doc:"uuid người lập phiếu"`
	}
}

type poPayOutput struct {
	Body struct {
		Payment       poPaymentDTO `json:"payment"`
		PurchaseOrder poDTO        `json:"purchase_order"`
	}
}

type poPaymentDTO struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Amount   string `json:"amount" doc:"Số tiền chi (chuỗi thập phân)"`
	Flow     string `json:"flow" doc:"'out' (phiếu chi)"`
	Status   string `json:"status" doc:"pending|completed"`
	BookType string `json:"book_type" doc:"INTERNAL|TAX|BOTH"`
}

func registerPaySupplier(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "purchasing-pay-supplier",
		Method:        http.MethodPost,
		Path:          "/v1/purchase-orders/{id}/pay",
		Summary:       "Chi trả NCC (received→paid): ghi phiếu chi + post sổ Dr 331/Cr 111/112 (atomic)",
		Tags:          []string{"purchasing"},
		Security:      []map[string][]string{{"bearerAuth": {}}},
		Middlewares:   huma.Middlewares{guard},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, in *poPayInput) (*poPayOutput, error) {
		amount, err := money.FromString(in.Body.Amount)
		if err != nil {
			return nil, humax.FromAppErr(apperr.Validation("amount không hợp lệ: " + in.Body.Amount))
		}
		res, err := svc.PaySupplier(ctx, app.PaySupplierInput{
			POID:          in.ID,
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
		out := &poPayOutput{}
		out.Body.Payment = poPaymentDTO{
			ID:       res.Payment.ID,
			Code:     res.Payment.Code,
			Amount:   res.Payment.Amount.String(),
			Flow:     res.Payment.Flow,
			Status:   res.Payment.Status,
			BookType: res.Payment.BookType.String(),
		}
		out.Body.PurchaseOrder = toPODTO(res.PO)
		return out, nil
	})
}

// ---- GET detail + list ----

type poGetInput struct {
	ID string `path:"id" doc:"ID đơn mua (uuid)"`
}

func registerGetPO(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "purchasing-get-po",
		Method:      http.MethodGet,
		Path:        "/v1/purchase-orders/{id}",
		Summary:     "Chi tiết một đơn mua (header + dòng hàng)",
		Tags:        []string{"purchasing"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *poGetInput) (*poDetailOutput, error) {
		res, err := svc.GetPO(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		return toDetailOutput(res), nil
	})
}

type poListInput struct {
	Cursor     string `query:"cursor" doc:"Cursor trang kế (opaque)"`
	Limit      int32  `query:"limit" minimum:"1" maximum:"200" doc:"Số PO mỗi trang (mặc định 50)"`
	Status     string `query:"status" enum:"draft,ordered,received,paid,cancelled" doc:"Lọc theo trạng thái"`
	SupplierID int64  `query:"supplier_id" doc:"Lọc theo NCC (bigint; 0/bỏ trống = không lọc)"`
}

type poListOutput struct {
	Body struct {
		Items      []poDTO `json:"items"`
		NextCursor string  `json:"next_cursor,omitempty" doc:"Cursor trang kế; rỗng nếu hết"`
	}
}

func registerListPOs(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "purchasing-list-pos",
		Method:      http.MethodGet,
		Path:        "/v1/purchase-orders",
		Summary:     "Danh sách đơn mua (keyset, lọc status/supplier)",
		Tags:        []string{"purchasing"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *poListInput) (*poListOutput, error) {
		var supplierID *int64
		if in.SupplierID > 0 {
			s := in.SupplierID
			supplierID = &s
		}
		res, err := svc.ListPOs(ctx, app.ListPOsQuery{
			Cursor:     in.Cursor,
			Limit:      in.Limit,
			Status:     in.Status,
			SupplierID: supplierID,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &poListOutput{}
		out.Body.Items = make([]poDTO, 0, len(res.Items))
		for _, p := range res.Items {
			out.Body.Items = append(out.Body.Items, toPODTO(p))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

// ---- mapping helpers ----

func toDetailOutput(res app.CreatedPO) *poDetailOutput {
	out := &poDetailOutput{}
	out.Body.PurchaseOrder = toPODTO(res.PO)
	out.Body.Lines = make([]poLineDTO, 0, len(res.Lines))
	for _, l := range res.Lines {
		out.Body.Lines = append(out.Body.Lines, toPOLineDTO(l))
	}
	return out
}

// toDraftLines chuyển DTO (qty/unit_cost/vat_rate/ngày chuỗi) → domain.DraftLine
// (decimal/money/time). Lỗi parse → apperr.Validation (422).
func toDraftLines(in []poCreateLineDTO) ([]domain.DraftLine, error) {
	out := make([]domain.DraftLine, 0, len(in))
	for _, l := range in {
		qty, err := decimal.NewFromString(l.Quantity)
		if err != nil {
			return nil, apperr.Validation("quantity không hợp lệ: " + l.Quantity)
		}
		cost, err := money.FromString(l.UnitCost)
		if err != nil {
			return nil, apperr.Validation("unit_cost không hợp lệ: " + l.UnitCost)
		}
		rate := decimal.Zero
		if l.VATRate != "" {
			rate, err = decimal.NewFromString(l.VATRate)
			if err != nil {
				return nil, apperr.Validation("vat_rate không hợp lệ: " + l.VATRate)
			}
		}
		expiry, err := parseDatePtr(l.ExpiryDate)
		if err != nil {
			return nil, apperr.Validation("expiry_date không hợp lệ (YYYY-MM-DD): " + l.ExpiryDate)
		}
		mfg, err := parseDatePtr(l.ManufacturingDate)
		if err != nil {
			return nil, apperr.Validation("manufacturing_date không hợp lệ (YYYY-MM-DD): " + l.ManufacturingDate)
		}
		out = append(out, domain.DraftLine{
			ProductID:         l.ProductID,
			Quantity:          qty,
			UnitCost:          cost,
			VATRate:           rate,
			BatchCode:         l.BatchCode,
			ExpiryDate:        expiry,
			ManufacturingDate: mfg,
		})
	}
	return out, nil
}

func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
