// Package postgres là ADAPTER ra phía cơ sở dữ liệu của catalog: implement port
// domain.Repository bằng query sinh từ sqlc (appdb) và map row <-> entity domain.
// ĐỌC bảng public.* kế thừa (strangler-fig, ADR 0001). Nằm dưới internal/ nên
// module khác KHÔNG import được. Catalog read-only → repo chỉ có thao tác đọc,
// bind thẳng pool (không tx).
package postgres

import (
	"context"
	"errors"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/Maneva-AI/namviet-backend/internal/catalog/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/appdb"
)

// Repo implement domain.Repository trên appdb.Queries (bind pool).
type Repo struct{ q *appdb.Queries }

// NewRepo tạo repo từ một *appdb.Queries (đã bind pool).
func NewRepo(q *appdb.Queries) *Repo { return &Repo{q: q} }

func (r *Repo) ListProducts(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {
	rows, err := r.q.ListProducts(ctx, appdb.ListProductsParams{
		AfterID:    f.AfterID,
		RowLimit:   f.Limit,
		Status:     strPtr(f.Status),
		CategoryID: f.CategoryID,
		Q:          strPtr(f.Query),
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Product, 0, len(rows))
	for _, row := range rows {
		out = append(out, productRowToDomain(productRow(row)))
	}
	return out, nil
}

func (r *Repo) GetProductByID(ctx context.Context, id int64) (domain.Product, error) {
	row, err := r.q.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, apperr.NotFound("sản phẩm không tồn tại")
		}
		return domain.Product{}, err
	}
	return productRowToDomain(productRow(row)), nil
}

func (r *Repo) ListUnits(ctx context.Context, productID int64) ([]domain.ProductUnit, error) {
	rows, err := r.q.ListProductUnits(ctx, &productID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.ProductUnit, 0, len(rows))
	for _, row := range rows {
		out = append(out, unitRowToDomain(row))
	}
	return out, nil
}

func (r *Repo) ListCategories(ctx context.Context, status string) ([]domain.Category, error) {
	rows, err := r.q.ListCategories(ctx, strPtr(status))
	if err != nil {
		return nil, err
	}
	out := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Category{
			ID:       row.ID,
			Name:     row.Name,
			Slug:     row.Slug,
			ParentID: row.ParentID,
			Status:   row.Status,
		})
	}
	return out, nil
}

func (r *Repo) ListManufacturers(ctx context.Context, status string) ([]domain.Manufacturer, error) {
	rows, err := r.q.ListManufacturers(ctx, strPtr(status))
	if err != nil {
		return nil, err
	}
	out := make([]domain.Manufacturer, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Manufacturer{
			ID:      row.ID,
			Name:    row.Name,
			Slug:    row.Slug,
			Country: derefStr(row.Country),
			LogoURL: derefStr(row.LogoUrl),
			Status:  derefStr(row.Status),
		})
	}
	return out, nil
}

// ---- mapping row <-> domain ----

// productRow là tập cột chung của ListProductsRow và GetProductByIDRow (cùng
// SELECT). Dùng một struct trung gian để map một lần, tránh lặp 2 hàm giống nhau.
type productRow struct {
	ID               int64
	Name             string
	Sku              *string
	Barcode          *string
	Status           string
	CategoryID       *int64
	ManufacturerID   *int64
	CategoryName     *string
	ManufacturerName *string
	InvoicePrice     pgtype.Numeric
	ActualCost       pgtype.Numeric
	WholesaleUnit    *string
	RetailUnit       *string
	ConversionFactor *int32
	ProductImages    []string
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
}

func productRowToDomain(p productRow) domain.Product {
	images := p.ProductImages
	if images == nil {
		images = []string{}
	}
	return domain.Product{
		ID:               p.ID,
		Name:             p.Name,
		SKU:              derefStr(p.Sku),
		Barcode:          derefStr(p.Barcode),
		Status:           p.Status,
		CategoryID:       p.CategoryID,
		ManufacturerID:   p.ManufacturerID,
		CategoryName:     derefStr(p.CategoryName),
		ManufacturerName: derefStr(p.ManufacturerName),
		InvoicePrice:     numericToMoney(p.InvoicePrice),
		ActualCost:       numericToMoney(p.ActualCost),
		WholesaleUnit:    derefStr(p.WholesaleUnit),
		RetailUnit:       derefStr(p.RetailUnit),
		ConversionFactor: derefInt32OrOne(p.ConversionFactor),
		Images:           images,
		CreatedAt:        p.CreatedAt.Time,
		UpdatedAt:        p.UpdatedAt.Time,
	}
}

func unitRowToDomain(u appdb.ListProductUnitsRow) domain.ProductUnit {
	var productID int64
	if u.ProductID != nil {
		productID = *u.ProductID
	}
	return domain.ProductUnit{
		ID:             u.ID,
		ProductID:      productID,
		UnitName:       u.UnitName,
		ConversionRate: derefInt32OrOne(u.ConversionRate),
		Barcode:        derefStr(u.Barcode),
		IsBase:         u.IsBase != nil && *u.IsBase,
		IsDirectSale:   u.IsDirectSale == nil || *u.IsDirectSale, // default true
		UnitType:       derefStr(u.UnitType),
		PriceCost:      numericToMoney(u.PriceCost),
		PriceSell:      numericToMoney(u.PriceSell),
		Price:          numericToMoney(u.Price),
	}
}

// numericToMoney chuyển pgtype.Numeric (do sqlc sinh) sang money.Money KHÔNG đi
// qua float: dựng decimal trực tiếp từ mantissa (big.Int) * 10^Exp. NULL/NaN →
// Zero (cột giá NULL coi như 0 tiền).
func numericToMoney(n pgtype.Numeric) money.Money {
	if !n.Valid || n.NaN || n.Int == nil {
		return money.Zero()
	}
	d := decimal.NewFromBigInt(new(big.Int).Set(n.Int), n.Exp)
	return money.FromDecimal(d)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt32OrOne(v *int32) int32 {
	if v == nil {
		return 1
	}
	return *v
}

// Đảm bảo Repo thoả port domain ở compile-time.
var _ domain.Repository = (*Repo)(nil)
