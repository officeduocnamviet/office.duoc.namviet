// Package http là ADAPTER vào của accounting: DTO + handler Huma + đăng ký route
// ĐỌC. Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, bảo vệ mọi route bằng authz.RequirePermissionHuma
// ("accounting.read"). Tiền debit/credit ra DTO dạng CHUỖI thập phân — KHÔNG
// float. APPEND-ONLY: CHỈ có route ĐỌC; post bút toán là port nội bộ (app.Poster)
// cho module orders/finance/vat — KHÔNG expose REST POST ở P0. Nằm dưới internal/.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/app"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Service là cổng use-case ĐỌC mà handler cần (interface để test bằng fake).
type Service interface {
	ListEntries(ctx context.Context, q app.ListEntriesQuery) (app.ListEntriesResult, error)
	GetEntry(ctx context.Context, id string) (domain.JournalEntryRecord, error)
}

const permRead = "accounting.read"

// Register đăng ký các operation ĐỌC /v1/accounting/* lên huma.API. Mọi route yêu
// cầu quyền accounting.read (verify token + check perm qua RequirePermissionHuma).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListEntries(api, svc, guard)
	registerGetEntry(api, svc, guard)
}

// ---- DTO ----

const dateLayout = "2006-01-02"

type entryDTO struct {
	ID         string `json:"id"`
	Book       string `json:"book" doc:"Sổ kế toán: INTERNAL (thực tế) hoặc TAX (thuế)"`
	EntryDate  string `json:"entry_date" doc:"Ngày ghi sổ (YYYY-MM-DD)"`
	PeriodID   string `json:"period_id"`
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	Memo       string `json:"memo"`
	CreatedAt  string `json:"created_at" doc:"Thời điểm tạo (RFC3339)"`
}

type lineDTO struct {
	AccountCode string `json:"account_code"`
	Debit       string `json:"debit" doc:"Số tiền ghi NỢ (chuỗi thập phân, không float)"`
	Credit      string `json:"credit" doc:"Số tiền ghi CÓ (chuỗi thập phân, không float)"`
}

type entryDetailDTO struct {
	entryDTO
	Lines []lineDTO `json:"lines"`
}

func toEntryDTO(r domain.JournalEntryRecord) entryDTO {
	return entryDTO{
		ID:         r.ID,
		Book:       string(r.Book),
		EntryDate:  r.EntryDate.Format(dateLayout),
		PeriodID:   r.PeriodID,
		SourceType: r.SourceType,
		SourceID:   r.SourceID,
		Memo:       r.Memo,
		CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toEntryDetailDTO(r domain.JournalEntryRecord) entryDetailDTO {
	d := entryDetailDTO{entryDTO: toEntryDTO(r)}
	d.Lines = make([]lineDTO, 0, len(r.Lines))
	for _, l := range r.Lines {
		d.Lines = append(d.Lines, lineDTO{
			AccountCode: l.AccountCode,
			Debit:       l.Debit.String(),
			Credit:      l.Credit.String(),
		})
	}
	return d
}

// ---- Inputs/Outputs ----

type listEntriesInput struct {
	Cursor   string `query:"cursor" doc:"Con trỏ trang (opaque); rỗng = trang đầu"`
	Limit    int32  `query:"limit" minimum:"1" maximum:"200" doc:"Số bút toán mỗi trang (mặc định 50)"`
	Book     string `query:"book" enum:"INTERNAL,TAX" doc:"Lọc theo sổ; rỗng = cả hai"`
	PeriodID string `query:"period_id" doc:"Lọc theo kỳ kế toán (uuid); rỗng = mọi kỳ"`
}

type listEntriesOutput struct {
	Body struct {
		Items      []entryDTO `json:"items"`
		NextCursor string     `json:"next_cursor" doc:"Con trỏ trang kế; rỗng = hết"`
	}
}

type getEntryInput struct {
	ID string `path:"id" doc:"ID bút toán (uuid)"`
}

type getEntryOutput struct {
	Body entryDetailDTO
}

// ---- Handlers ----

func registerListEntries(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "accounting-list-entries",
		Method:      http.MethodGet,
		Path:        "/v1/accounting/entries",
		Summary:     "Danh sách bút toán sổ kế toán (keyset pagination)",
		Tags:        []string{"accounting"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listEntriesInput) (*listEntriesOutput, error) {
		res, err := svc.ListEntries(ctx, app.ListEntriesQuery{
			Cursor:   in.Cursor,
			Limit:    in.Limit,
			Book:     in.Book,
			PeriodID: in.PeriodID,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listEntriesOutput{}
		out.Body.Items = make([]entryDTO, 0, len(res.Items))
		for _, r := range res.Items {
			out.Body.Items = append(out.Body.Items, toEntryDTO(r))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

func registerGetEntry(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "accounting-get-entry",
		Method:      http.MethodGet,
		Path:        "/v1/accounting/entries/{id}",
		Summary:     "Chi tiết một bút toán kèm các dòng nợ/có",
		Tags:        []string{"accounting"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *getEntryInput) (*getEntryOutput, error) {
		rec, err := svc.GetEntry(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		return &getEntryOutput{Body: toEntryDetailDTO(rec)}, nil
	})
}
