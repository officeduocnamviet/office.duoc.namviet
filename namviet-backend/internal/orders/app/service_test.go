package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/orders/app"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

// fakeRepo là fake (không phải mock) của domain.Repository — ưu tiên fake cho
// port (ARCHITECTURE.md §11).
type fakeRepo struct {
	orders     []domain.Order
	lines      []domain.OrderLine
	getErr     error
	linesErr   error
	lastFilter domain.OrderFilter
	getSeen    string
	linesSeen  string
}

func (f *fakeRepo) ListOrders(_ context.Context, flt domain.OrderFilter) ([]domain.Order, error) {
	f.lastFilter = flt
	return f.orders, nil
}

func (f *fakeRepo) GetOrderByID(_ context.Context, id string) (domain.Order, error) {
	f.getSeen = id
	if f.getErr != nil {
		return domain.Order{}, f.getErr
	}
	return domain.Order{ID: id, Code: "HD1"}, nil
}

func (f *fakeRepo) ListLines(_ context.Context, orderID string) ([]domain.OrderLine, error) {
	f.linesSeen = orderID
	if f.linesErr != nil {
		return nil, f.linesErr
	}
	return f.lines, nil
}

func TestListOrders_DefaultLimit(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	if _, err := svc.ListOrders(context.Background(), app.ListOrdersQuery{}); err != nil {
		t.Fatalf("ListOrders: %v", err)
	}
	if f.lastFilter.Limit != 20 {
		t.Fatalf("limit mặc định = %d, want 20", f.lastFilter.Limit)
	}
	if f.lastFilter.AfterCreatedAt != 0 || f.lastFilter.AfterID != "" {
		t.Fatalf("trang đầu phải không có after: %+v", f.lastFilter)
	}
}

func TestListOrders_LimitClamp(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	_, _ = svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 9999})
	if f.lastFilter.Limit != 100 {
		t.Fatalf("limit clamp = %d, want 100", f.lastFilter.Limit)
	}
}

func TestListOrders_FiltersPassthrough(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	cid := int64(201)
	_, err := svc.ListOrders(context.Background(), app.ListOrdersQuery{
		CustomerID:    &cid,
		Status:        "COMPLETED",
		PaymentStatus: "partial",
		FromDate:      111,
		ToDate:        222,
	})
	if err != nil {
		t.Fatalf("ListOrders: %v", err)
	}
	if f.lastFilter.CustomerID == nil || *f.lastFilter.CustomerID != 201 {
		t.Fatalf("customer filter sai: %+v", f.lastFilter.CustomerID)
	}
	if f.lastFilter.Status != "COMPLETED" || f.lastFilter.PaymentStatus != "partial" {
		t.Fatalf("status/payment filter sai: %+v", f.lastFilter)
	}
	if f.lastFilter.FromDate != 111 || f.lastFilter.ToDate != 222 {
		t.Fatalf("date filter sai: %+v", f.lastFilter)
	}
}

func TestListOrders_BadCursor_Validation(t *testing.T) {
	svc := app.New(&fakeRepo{})
	_, err := svc.ListOrders(context.Background(), app.ListOrdersQuery{Cursor: "###"})
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("kind = %v, want Validation", apperr.KindOf(err))
	}
}

func TestListOrders_NextCursor_WhenFull(t *testing.T) {
	ts := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	orders := make([]domain.Order, 20)
	for i := range orders {
		orders[i] = domain.Order{ID: "uuid-" + string(rune('a'+i)), CreatedAt: ts}
	}
	f := &fakeRepo{orders: orders}
	svc := app.New(f)
	res, _ := svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 20})
	if res.NextCursor == "" {
		t.Fatal("trang đầy phải có NextCursor")
	}
	// Cursor phải decode được và trỏ về phần tử cuối — feed lại vào filter.
	_, err := svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 20, Cursor: res.NextCursor})
	if err != nil {
		t.Fatalf("dùng lại cursor lỗi: %v", err)
	}
	if f.lastFilter.AfterID != "uuid-"+string(rune('a'+19)) {
		t.Fatalf("AfterID từ cursor sai: %q", f.lastFilter.AfterID)
	}
	if f.lastFilter.AfterCreatedAt != ts.UnixNano() {
		t.Fatalf("AfterCreatedAt từ cursor sai: %d", f.lastFilter.AfterCreatedAt)
	}
}

func TestListOrders_NoNextCursor_WhenPartial(t *testing.T) {
	f := &fakeRepo{orders: []domain.Order{{ID: "a"}, {ID: "b"}}}
	svc := app.New(f)
	res, _ := svc.ListOrders(context.Background(), app.ListOrdersQuery{Limit: 20})
	if res.NextCursor != "" {
		t.Fatalf("trang chưa đầy không được có NextCursor, got %q", res.NextCursor)
	}
}

func TestGetOrder_LoadsOrderAndLines(t *testing.T) {
	f := &fakeRepo{lines: []domain.OrderLine{{ID: "l1"}, {ID: "l2"}}}
	svc := app.New(f)
	d, err := svc.GetOrder(context.Background(), "ord-7")
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	if d.Order.ID != "ord-7" || f.getSeen != "ord-7" {
		t.Fatalf("phải gọi GetOrderByID với ord-7: got=%q seen=%q", d.Order.ID, f.getSeen)
	}
	if f.linesSeen != "ord-7" || len(d.Lines) != 2 {
		t.Fatalf("phải nạp lines của ord-7: seen=%q n=%d", f.linesSeen, len(d.Lines))
	}
}

func TestGetOrder_NotFoundPassthrough(t *testing.T) {
	f := &fakeRepo{getErr: apperr.NotFound("không có đơn")}
	svc := app.New(f)
	_, err := svc.GetOrder(context.Background(), "x")
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("kind = %v, want NotFound", apperr.KindOf(err))
	}
}
