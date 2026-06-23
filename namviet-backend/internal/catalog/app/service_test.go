package app_test

import (
	"context"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/catalog/app"
	"github.com/Maneva-AI/namviet-backend/internal/catalog/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
)

// fakeRepo là fake (không phải mock) của domain.Repository — ưu tiên fake cho
// port (ARCHITECTURE.md §11).
type fakeRepo struct {
	products    []domain.Product
	units       []domain.ProductUnit
	getErr      error
	lastFilter  domain.ProductFilter
	getByIDSeen int64
}

func (f *fakeRepo) ListProducts(_ context.Context, flt domain.ProductFilter) ([]domain.Product, error) {
	f.lastFilter = flt
	return f.products, nil
}
func (f *fakeRepo) GetProductByID(_ context.Context, id int64) (domain.Product, error) {
	f.getByIDSeen = id
	if f.getErr != nil {
		return domain.Product{}, f.getErr
	}
	return domain.Product{ID: id, Name: "X"}, nil
}
func (f *fakeRepo) ListUnits(_ context.Context, _ int64) ([]domain.ProductUnit, error) {
	return f.units, nil
}
func (f *fakeRepo) ListCategories(_ context.Context, _ string) ([]domain.Category, error) {
	return nil, nil
}
func (f *fakeRepo) ListManufacturers(_ context.Context, _ string) ([]domain.Manufacturer, error) {
	return nil, nil
}

func TestListProducts_DefaultLimit(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	if _, err := svc.ListProducts(context.Background(), app.ListProductsQuery{}); err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	if f.lastFilter.Limit != 20 {
		t.Fatalf("limit mặc định = %d, want 20", f.lastFilter.Limit)
	}
	if f.lastFilter.AfterID != 0 {
		t.Fatalf("AfterID trang đầu = %d, want 0", f.lastFilter.AfterID)
	}
}

func TestListProducts_LimitClamp(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	_, _ = svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 9999})
	if f.lastFilter.Limit != 100 {
		t.Fatalf("limit clamp = %d, want 100", f.lastFilter.Limit)
	}
}

func TestListProducts_CursorDecodedToAfterID(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	cur := pagination.EncodeID(42)
	_, _ = svc.ListProducts(context.Background(), app.ListProductsQuery{Cursor: cur})
	if f.lastFilter.AfterID != 42 {
		t.Fatalf("AfterID = %d, want 42", f.lastFilter.AfterID)
	}
}

func TestListProducts_BadCursor_Validation(t *testing.T) {
	svc := app.New(&fakeRepo{})
	_, err := svc.ListProducts(context.Background(), app.ListProductsQuery{Cursor: "###"})
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("kind = %v, want Validation", apperr.KindOf(err))
	}
}

func TestListProducts_NextCursor_WhenFull(t *testing.T) {
	// Trả đúng `limit` phần tử → trang đầy → có NextCursor = id phần tử cuối.
	prods := make([]domain.Product, 20)
	for i := range prods {
		prods[i] = domain.Product{ID: int64(i + 1)}
	}
	f := &fakeRepo{products: prods}
	svc := app.New(f)
	res, _ := svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 20})
	if res.NextCursor == "" {
		t.Fatal("trang đầy phải có NextCursor")
	}
	gotID, _ := pagination.DecodeID(res.NextCursor)
	if gotID != 20 {
		t.Fatalf("NextCursor id = %d, want 20", gotID)
	}
}

func TestListProducts_NoNextCursor_WhenPartial(t *testing.T) {
	f := &fakeRepo{products: []domain.Product{{ID: 1}, {ID: 2}}} // < limit
	svc := app.New(f)
	res, _ := svc.ListProducts(context.Background(), app.ListProductsQuery{Limit: 20})
	if res.NextCursor != "" {
		t.Fatalf("trang chưa đầy không được có NextCursor, got %q", res.NextCursor)
	}
}

func TestGetProduct_AttachesUnits(t *testing.T) {
	f := &fakeRepo{units: []domain.ProductUnit{{ID: 1, UnitName: "Viên"}}}
	svc := app.New(f)
	d, err := svc.GetProduct(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetProduct: %v", err)
	}
	if d.Product.ID != 7 || len(d.Units) != 1 || d.Units[0].UnitName != "Viên" {
		t.Fatalf("detail sai: %+v", d)
	}
}

func TestGetProduct_NotFoundPassthrough(t *testing.T) {
	f := &fakeRepo{getErr: apperr.NotFound("không có")}
	svc := app.New(f)
	_, err := svc.GetProduct(context.Background(), 99)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("kind = %v, want NotFound", apperr.KindOf(err))
	}
}
