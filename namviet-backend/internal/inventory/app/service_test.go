package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/app"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

func date(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

// fakeRepo là port giả (ARCHITECTURE.md §11: fakes > mocks cho port) để unit test
// use-case không cần DB. Ghi lại tham số nhận được để assert chuẩn hoá ở app.
type fakeRepo struct {
	stock        []domain.StockItem
	batches      []domain.Batch
	warehouses   []domain.Warehouse
	gotStock     domain.StockFilter
	gotWarehouse domain.WarehouseFilter
}

func (f *fakeRepo) ListWarehouses(_ context.Context, fl domain.WarehouseFilter) ([]domain.Warehouse, error) {
	f.gotWarehouse = fl
	return f.warehouses, nil
}

func (f *fakeRepo) ListStock(_ context.Context, fl domain.StockFilter) ([]domain.StockItem, error) {
	f.gotStock = fl
	return f.stock, nil
}

func (f *fakeRepo) ListBatchesFEFO(_ context.Context, _ int64, _ *int64) ([]domain.Batch, error) {
	return f.batches, nil
}

func TestListStock_NormalizesLimitAndDecodesCursor(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)

	// limit 0 → default 50; cursor mã hoá id=42 → AfterID 42.
	cur := pagination.EncodeID(42)
	if _, err := svc.ListStock(context.Background(), app.ListStockQuery{Cursor: cur}); err != nil {
		t.Fatalf("ListStock: %v", err)
	}
	if f.gotStock.Limit != 50 {
		t.Fatalf("limit mặc định = %d, want 50", f.gotStock.Limit)
	}
	if f.gotStock.AfterID != 42 {
		t.Fatalf("AfterID = %d, want 42 (decode cursor)", f.gotStock.AfterID)
	}

	// limit quá lớn → clamp 200.
	_, _ = svc.ListStock(context.Background(), app.ListStockQuery{Limit: 9999})
	if f.gotStock.Limit != 200 {
		t.Fatalf("limit clamp = %d, want 200", f.gotStock.Limit)
	}
}

func TestListStock_BadCursor_Validation(t *testing.T) {
	svc := app.New(&fakeRepo{})
	if _, err := svc.ListStock(context.Background(), app.ListStockQuery{Cursor: "!!!not-base64"}); err == nil {
		t.Fatal("cursor sai phải trả lỗi validation, không nuốt")
	}
}

func TestListStock_NextCursorWhenFull(t *testing.T) {
	f := &fakeRepo{stock: []domain.StockItem{{ID: 1}, {ID: 2}}}
	svc := app.New(f)
	// limit 2, trả đúng 2 → trang đầy → có NextCursor = id phần tử cuối (2).
	res, _ := svc.ListStock(context.Background(), app.ListStockQuery{Limit: 2})
	if res.NextCursor == "" {
		t.Fatal("trang đầy phải có NextCursor")
	}
	gotID, _ := pagination.DecodeID(res.NextCursor)
	if gotID != 2 {
		t.Fatalf("NextCursor decode = %d, want 2", gotID)
	}

	// trang chưa đầy → hết, NextCursor rỗng.
	f2 := &fakeRepo{stock: []domain.StockItem{{ID: 1}}}
	res2, _ := app.New(f2).ListStock(context.Background(), app.ListStockQuery{Limit: 2})
	if res2.NextCursor != "" {
		t.Fatal("trang chưa đầy phải rỗng NextCursor")
	}
}

func TestListBatchesFEFO_ServiceSortsByExpiry(t *testing.T) {
	// Repo trả thứ tự lộn xộn → service phải sắp FEFO (expiry tăng dần).
	f := &fakeRepo{batches: []domain.Batch{
		{InventoryBatchID: 1, BatchID: 1, ExpiryDate: date("2027-01-01")},
		{InventoryBatchID: 2, BatchID: 2, ExpiryDate: date("2026-01-01")},
	}}
	out, err := app.New(f).ListBatchesFEFO(context.Background(), 100, nil)
	if err != nil {
		t.Fatalf("ListBatchesFEFO: %v", err)
	}
	if len(out) != 2 || out[0].InventoryBatchID != 2 {
		t.Fatalf("service phải sắp FEFO: got %v", []int64{out[0].InventoryBatchID, out[1].InventoryBatchID})
	}
}

func TestListWarehouses_NormalizesLimit(t *testing.T) {
	f := &fakeRepo{}
	_, _ = app.New(f).ListWarehouses(context.Background(), "active", 0)
	if f.gotWarehouse.Limit != 50 {
		t.Fatalf("warehouse limit mặc định = %d, want 50", f.gotWarehouse.Limit)
	}
	if f.gotWarehouse.Status != "active" {
		t.Fatalf("status truyền sai: %q", f.gotWarehouse.Status)
	}
}
