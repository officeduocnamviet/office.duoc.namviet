// Package domain là LÕI THUẦN của bounded context catalog: entity sản phẩm và
// danh mục liên quan + PORT interface. KHÔNG import pgx/http/huma/framework
// (ARCHITECTURE.md §3). Chỉ stdlib + shared kernel trung lập (common/money).
// Phụ thuộc đi một chiều: adapters → app → domain.
//
// Catalog là context "nhẹ" (đa số read) nên domain mỏng: chủ yếu là dữ liệu
// product master + giá CƠ SỞ. Giá theo từng khách hàng KHÔNG thuộc catalog (DEFER
// sang module customers/pricing).
package domain

import (
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// Product là aggregate gốc của catalog. ID là int64 vì bảng public.products dùng
// khoá bigint (KHÔNG uuid). category_name/manufacturer_name là cột CACHE trên
// products (database_schema.md) — dùng trực tiếp, tránh join nóng khi list.
type Product struct {
	ID               int64
	Name             string
	SKU              string // có thể rỗng (cột nullable)
	Barcode          string
	Status           string
	CategoryID       *int64 // nil nếu chưa phân loại
	ManufacturerID   *int64
	CategoryName     string // cache trên products
	ManufacturerName string // cache trên products
	InvoicePrice     money.Money
	ActualCost       money.Money
	WholesaleUnit    string // UOM legacy (database_schema.md)
	RetailUnit       string
	ConversionFactor int32 // quy đổi sỉ/lẻ; mặc định 1
	Images           []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ProductUnit là một đơn vị tính đa cấp (Hộp/Vỉ/Viên) của một product, đọc từ
// public.product_units. PriceSell/PriceCost là giá định mức của ĐVT (tiền →
// money.Money, không float).
type ProductUnit struct {
	ID             int64
	ProductID      int64
	UnitName       string
	ConversionRate int32 // tỷ lệ quy đổi so với đơn vị nhỏ nhất
	Barcode        string
	IsBase         bool // đơn vị cơ sở (nhỏ nhất)
	IsDirectSale   bool // cho phép bán lẻ ĐVT này
	UnitType       string
	PriceCost      money.Money
	PriceSell      money.Money
	Price          money.Money
}

// Category là một danh mục phân loại sản phẩm (public.categories).
type Category struct {
	ID       int64
	Name     string
	Slug     string
	ParentID *int64
	Status   string
}

// Manufacturer là một hãng / nhà sản xuất (public.manufacturers).
type Manufacturer struct {
	ID      int64
	Name    string
	Slug    string
	Country string
	LogoURL string
	Status  string
}
