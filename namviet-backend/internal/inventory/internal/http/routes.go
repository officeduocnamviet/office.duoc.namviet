// Package http là ADAPTER vào của inventory: DTO + handler Huma + đăng ký route.
// Dịch HTTP <-> use-case (app), map lỗi nghiệp vụ (apperr) sang envelope qua
// humax.FromAppErr, và bảo vệ mọi route bằng authz.RequirePermissionHuma
// ("inventory.read"). Số lượng tồn (quantity) và giá vốn (money) ra DTO dạng
// CHUỖI thập phân — KHÔNG float — để FE không mất chính xác. Nằm dưới internal/
// nên module khác không import.
package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authz"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

// Service là cổng use-case mà handler cần (interface để test handler bằng fake).
type Service interface {
	ListWarehouses(ctx context.Context, status string, limit int32) ([]domain.Warehouse, error)
	ListStock(ctx context.Context, q app.ListStockQuery) (app.ListStockResult, error)
	ListBatchesFEFO(ctx context.Context, productID int64, warehouseID *int64) ([]domain.Batch, error)
}

const permRead = "inventory.read"

// Register đăng ký toàn bộ operation /v1/warehouses|inventory/* lên huma.API. Mọi
// route yêu cầu quyền inventory.read (verify token + check perm qua
// authz.RequirePermissionHuma — 1 enforcement point).
func Register(api huma.API, svc Service, verifier *authn.Verifier) {
	guard := authz.RequirePermissionHuma(api, verifier, permRead)
	registerListWarehouses(api, svc, guard)
	registerListStock(api, svc, guard)
	registerListBatches(api, svc, guard)
}

// ---- DTO ----

type warehouseDTO struct {
	ID         int64  `json:"id"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`
	Address    string `json:"address"`
	Type       string `json:"type"`
	Code       string `json:"code"`
	Manager    string `json:"manager"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
	CompanyID  string `json:"company_id"`
	OutletType string `json:"outlet_type"`
}

type stockDTO struct {
	ID            int64  `json:"id"`
	ProductID     int64  `json:"product_id"`
	WarehouseID   int64  `json:"warehouse_id"`
	Quantity      string `json:"quantity" doc:"Số lượng tồn (chuỗi thập phân, không float)"`
	MinStock      int32  `json:"min_stock"`
	MaxStock      int32  `json:"max_stock"`
	ShelfLocation string `json:"shelf_location"`
}

type batchDTO struct {
	InventoryBatchID  int64  `json:"inventory_batch_id"`
	WarehouseID       int64  `json:"warehouse_id"`
	ProductID         int64  `json:"product_id"`
	BatchID           int64  `json:"batch_id"`
	BatchCode         string `json:"batch_code"`
	Quantity          string `json:"quantity" doc:"Tồn của lô (chuỗi thập phân, không float)"`
	ExpiryDate        string `json:"expiry_date" doc:"Hạn dùng (YYYY-MM-DD)"`
	ManufacturingDate string `json:"manufacturing_date" doc:"Ngày sản xuất (YYYY-MM-DD); rỗng nếu không có"`
	InboundPrice      string `json:"inbound_price" doc:"Giá vốn nhập của lô (chuỗi thập phân)"`
}

const dateLayout = "2006-01-02"

func toWarehouseDTO(w domain.Warehouse) warehouseDTO {
	return warehouseDTO{
		ID: w.ID, Key: w.Key, Name: w.Name, Unit: w.Unit, Address: w.Address,
		Type: w.Type, Code: w.Code, Manager: w.Manager, Phone: w.Phone,
		Status: w.Status, CompanyID: w.CompanyID, OutletType: w.OutletType,
	}
}

func toStockDTO(s domain.StockItem) stockDTO {
	return stockDTO{
		ID:            s.ID,
		ProductID:     s.ProductID,
		WarehouseID:   s.WarehouseID,
		Quantity:      s.Quantity.String(),
		MinStock:      s.MinStock,
		MaxStock:      s.MaxStock,
		ShelfLocation: s.ShelfLocation,
	}
}

func toBatchDTO(b domain.Batch) batchDTO {
	mfg := ""
	if b.HasManufacturing {
		mfg = b.ManufacturingDate.Format(dateLayout)
	}
	return batchDTO{
		InventoryBatchID:  b.InventoryBatchID,
		WarehouseID:       b.WarehouseID,
		ProductID:         b.ProductID,
		BatchID:           b.BatchID,
		BatchCode:         b.BatchCode,
		Quantity:          b.Quantity.String(),
		ExpiryDate:        b.ExpiryDate.Format(dateLayout),
		ManufacturingDate: mfg,
		InboundPrice:      b.InboundPrice.String(),
	}
}

// ---- Inputs/Outputs ----

type listWarehousesInput struct {
	Status string `query:"status" doc:"Lọc trạng thái (vd active); rỗng = tất cả"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"200" doc:"Số kho tối đa (mặc định 50)"`
}

type listWarehousesOutput struct {
	Body struct {
		Items []warehouseDTO `json:"items"`
	}
}

type listStockInput struct {
	Cursor string `query:"cursor" doc:"Con trỏ trang (opaque); rỗng = trang đầu"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"200" doc:"Số bản ghi mỗi trang (mặc định 50)"`
	// 0 (hoặc không truyền) = không lọc; id bigint thật luôn > 0.
	ProductID   int64 `query:"product_id" minimum:"0" doc:"Lọc theo sản phẩm (0 = tất cả)"`
	WarehouseID int64 `query:"warehouse_id" minimum:"0" doc:"Lọc theo kho (0 = tất cả)"`
}

type listStockOutput struct {
	Body struct {
		Items      []stockDTO `json:"items"`
		NextCursor string     `json:"next_cursor" doc:"Con trỏ trang kế; rỗng = hết"`
	}
}

type listBatchesInput struct {
	ProductID   int64 `query:"product_id" required:"true" minimum:"1" doc:"Sản phẩm cần xem lô (bắt buộc)"`
	WarehouseID int64 `query:"warehouse_id" minimum:"0" doc:"Lọc theo kho (0 = tất cả)"`
}

type listBatchesOutput struct {
	Body struct {
		Items []batchDTO `json:"items" doc:"Danh sách lô còn tồn, sắp FEFO (hạn dùng tăng dần)"`
	}
}

// ---- Handlers ----

func registerListWarehouses(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "inventory-list-warehouses",
		Method:      http.MethodGet,
		Path:        "/v1/warehouses",
		Summary:     "Danh sách kho / chi nhánh",
		Tags:        []string{"inventory"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listWarehousesInput) (*listWarehousesOutput, error) {
		ws, err := svc.ListWarehouses(ctx, in.Status, in.Limit)
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listWarehousesOutput{}
		out.Body.Items = make([]warehouseDTO, 0, len(ws))
		for _, w := range ws {
			out.Body.Items = append(out.Body.Items, toWarehouseDTO(w))
		}
		return out, nil
	})
}

func registerListStock(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "inventory-list-stock",
		Method:      http.MethodGet,
		Path:        "/v1/inventory/stock",
		Summary:     "Tồn kho theo sản phẩm và/hoặc kho (keyset pagination)",
		Tags:        []string{"inventory"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listStockInput) (*listStockOutput, error) {
		res, err := svc.ListStock(ctx, app.ListStockQuery{
			Cursor:      in.Cursor,
			Limit:       in.Limit,
			ProductID:   optID(in.ProductID),
			WarehouseID: optID(in.WarehouseID),
		})
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listStockOutput{}
		out.Body.Items = make([]stockDTO, 0, len(res.Items))
		for _, s := range res.Items {
			out.Body.Items = append(out.Body.Items, toStockDTO(s))
		}
		out.Body.NextCursor = res.NextCursor
		return out, nil
	})
}

func registerListBatches(api huma.API, svc Service, guard func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "inventory-list-batches",
		Method:      http.MethodGet,
		Path:        "/v1/inventory/batches",
		Summary:     "Danh sách lô còn tồn của sản phẩm (sắp FEFO theo hạn dùng)",
		Tags:        []string{"inventory"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{guard},
	}, func(ctx context.Context, in *listBatchesInput) (*listBatchesOutput, error) {
		batches, err := svc.ListBatchesFEFO(ctx, in.ProductID, optID(in.WarehouseID))
		if err != nil {
			return nil, humax.FromAppErr(err)
		}
		out := &listBatchesOutput{}
		out.Body.Items = make([]batchDTO, 0, len(batches))
		for _, b := range batches {
			out.Body.Items = append(out.Body.Items, toBatchDTO(b))
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
