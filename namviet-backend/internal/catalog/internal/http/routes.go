// Package http là ADAPTER vào của catalog: DTO + handler Huma + đăng ký route.
// Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, và bảo vệ mọi route bằng authz.RequirePermissionHuma
// ("catalog.read"). Tiền ra DTO dạng CHUỖI thập phân (money.String) — KHÔNG float
// — để FE không mất chính xác. Nằm dưới internal/ nên module khác không import.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/catalog/app"
	"github.com/Maneva-AI/namviet-backend/internal/catalog/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Service là cổng use-case mà handler cần (interface để test handler bằng fake).
type Service interface {
	ListProducts(ctx context.Context, q app.ListProductsQuery) (app.ListProductsResult, error)
	GetProduct(ctx context.Context, id int64) (app.ProductDetail, error)
	ListCategories(ctx context.Context, status string) ([]domain.Category, error)
	ListManufacturers(ctx context.Context, status string) ([]domain.Manufacturer, error)
}

const permRead = "catalog.read"

// Register đăng ký toàn bộ operation /v1/products|categories|manufacturers lên
// huma.API. Mọi route yêu cầu quyền catalog.read (verify token + check perm qua
// authz.RequirePermissionHuma — 1 enforcement point).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListProducts(api, svc, guard)
	registerGetProduct(api, svc, guard)
	registerListCategories(api, svc, guard)
	registerListManufacturers(api, svc, guard)
}

// ---- DTO ----

type productDTO struct {
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	SKU              string   `json:"sku"`
	Barcode          string   `json:"barcode"`
	Status           string   `json:"status"`
	CategoryID       *int64   `json:"category_id"`
	ManufacturerID   *int64   `json:"manufacturer_id"`
	CategoryName     string   `json:"category_name"`
	ManufacturerName string   `json:"manufacturer_name"`
	InvoicePrice     string   `json:"invoice_price" doc:"Giá hóa đơn VAT (chuỗi thập phân, không float)"`
	ActualCost       string   `json:"actual_cost" doc:"Giá vốn thực tế (chuỗi thập phân)"`
	WholesaleUnit    string   `json:"wholesale_unit"`
	RetailUnit       string   `json:"retail_unit"`
	ConversionFactor int32    `json:"conversion_factor"`
	Images           []string `json:"images"`
}

type unitDTO struct {
	ID             int64  `json:"id"`
	UnitName       string `json:"unit_name"`
	ConversionRate int32  `json:"conversion_rate"`
	Barcode        string `json:"barcode"`
	IsBase         bool   `json:"is_base"`
	IsDirectSale   bool   `json:"is_direct_sale"`
	UnitType       string `json:"unit_type"`
	PriceCost      string `json:"price_cost"`
	PriceSell      string `json:"price_sell"`
	Price          string `json:"price"`
}

type categoryDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID *int64 `json:"parent_id"`
	Status   string `json:"status"`
}

type manufacturerDTO struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Country string `json:"country"`
	LogoURL string `json:"logo_url"`
	Status  string `json:"status"`
}

func toProductDTO(p domain.Product) productDTO {
	return productDTO{
		ID:               p.ID,
		Name:             p.Name,
		SKU:              p.SKU,
		Barcode:          p.Barcode,
		Status:           p.Status,
		CategoryID:       p.CategoryID,
		ManufacturerID:   p.ManufacturerID,
		CategoryName:     p.CategoryName,
		ManufacturerName: p.ManufacturerName,
		InvoicePrice:     p.InvoicePrice.String(),
		ActualCost:       p.ActualCost.String(),
		WholesaleUnit:    p.WholesaleUnit,
		RetailUnit:       p.RetailUnit,
		ConversionFactor: p.ConversionFactor,
		Images:           p.Images,
	}
}

func toUnitDTO(u domain.ProductUnit) unitDTO {
	return unitDTO{
		ID:             u.ID,
		UnitName:       u.UnitName,
		ConversionRate: u.ConversionRate,
		Barcode:        u.Barcode,
		IsBase:         u.IsBase,
		IsDirectSale:   u.IsDirectSale,
		UnitType:       u.UnitType,
		PriceCost:      u.PriceCost.String(),
		PriceSell:      u.PriceSell.String(),
		Price:          u.Price.String(),
	}
}

// ---- Inputs/Outputs ----

type listProductsInput struct {
	Cursor string `query:"cursor" doc:"Con trỏ trang (opaque); rỗng = trang đầu"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"100" doc:"Số sản phẩm mỗi trang (mặc định 20)"`
	// CategoryID 0 (hoặc không truyền) = không lọc; id bigint thật luôn > 0.
	CategoryID int64  `query:"category_id" minimum:"0" doc:"Lọc theo nhóm ngành hàng (0 = tất cả)"`
	Q          string `query:"q" doc:"Tìm theo tên/SKU"`
	Status     string `query:"status" doc:"Lọc trạng thái (vd active); rỗng = tất cả"`
}

type listProductsOutput struct {
	Body struct {
		Items      []productDTO `json:"items"`
		NextCursor string       `json:"next_cursor" doc:"Con trỏ trang kế; rỗng = hết"`
	}
}

type getProductInput struct {
	ID int64 `path:"id" doc:"ID sản phẩm (bigint)"`
}

type getProductOutput struct {
	Body struct {
		Product productDTO `json:"product"`
		Units   []unitDTO  `json:"units"`
	}
}

type listCategoriesInput struct {
	Status string `query:"status" doc:"Lọc trạng thái; rỗng = tất cả"`
}

type listCategoriesOutput struct {
	Body struct {
		Items []categoryDTO `json:"items"`
	}
}

type listManufacturersInput struct {
	Status string `query:"status" doc:"Lọc trạng thái; rỗng = tất cả"`
}

type listManufacturersOutput struct {
	Body struct {
		Items []manufacturerDTO `json:"items"`
	}
}

// ---- Handlers ----

func registerListProducts(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "catalog-list-products",
		Method:      http.MethodGet,
		Path:        "/v1/products",
		Summary:     "Danh sách sản phẩm (keyset pagination + filter)",
		Tags:        []string{"catalog"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listProductsInput) (*listProductsOutput, error) {
		var categoryID *int64
		if in.CategoryID > 0 {
			categoryID = &in.CategoryID
		}
		res, err := svc.ListProducts(ctx, app.ListProductsQuery{
			Cursor:     in.Cursor,
			Limit:      in.Limit,
			Status:     in.Status,
			CategoryID: categoryID,
			Query:      in.Q,
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listProductsOutput{}
		out.Body.Items = make([]productDTO, 0, len(res.Items))
		for _, p := range res.Items {
			out.Body.Items = append(out.Body.Items, toProductDTO(p))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

func registerGetProduct(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "catalog-get-product",
		Method:      http.MethodGet,
		Path:        "/v1/products/{id}",
		Summary:     "Chi tiết sản phẩm + đơn vị tính",
		Tags:        []string{"catalog"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *getProductInput) (*getProductOutput, error) {
		d, err := svc.GetProduct(ctx, in.ID)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &getProductOutput{}
		out.Body.Product = toProductDTO(d.Product)
		out.Body.Units = make([]unitDTO, 0, len(d.Units))
		for _, u := range d.Units {
			out.Body.Units = append(out.Body.Units, toUnitDTO(u))
		}
		return out, nil
	})
}

func registerListCategories(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "catalog-list-categories",
		Method:      http.MethodGet,
		Path:        "/v1/categories",
		Summary:     "Danh sách danh mục",
		Tags:        []string{"catalog"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listCategoriesInput) (*listCategoriesOutput, error) {
		cats, err := svc.ListCategories(ctx, in.Status)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listCategoriesOutput{}
		out.Body.Items = make([]categoryDTO, 0, len(cats))
		for _, c := range cats {
			out.Body.Items = append(out.Body.Items, categoryDTO{
				ID: c.ID, Name: c.Name, Slug: c.Slug, ParentID: c.ParentID, Status: c.Status,
			})
		}
		return out, nil
	})
}

func registerListManufacturers(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "catalog-list-manufacturers",
		Method:      http.MethodGet,
		Path:        "/v1/manufacturers",
		Summary:     "Danh sách hãng / nhà sản xuất",
		Tags:        []string{"catalog"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listManufacturersInput) (*listManufacturersOutput, error) {
		mans, err := svc.ListManufacturers(ctx, in.Status)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listManufacturersOutput{}
		out.Body.Items = make([]manufacturerDTO, 0, len(mans))
		for _, m := range mans {
			out.Body.Items = append(out.Body.Items, manufacturerDTO{
				ID: m.ID, Name: m.Name, Slug: m.Slug, Country: m.Country, LogoURL: m.LogoURL, Status: m.Status,
			})
		}
		return out, nil
	})
}
