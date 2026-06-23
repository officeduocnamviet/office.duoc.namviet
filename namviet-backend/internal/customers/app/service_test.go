package app_test

import (
	"context"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/pagination"
	"github.com/Maneva-AI/namviet-backend/internal/customers/app"
	"github.com/Maneva-AI/namviet-backend/internal/customers/domain"
)

// fakeRepo là fake (không phải mock) của domain.Repository — ưu tiên fake cho
// port (ARCHITECTURE.md §11).
type fakeRepo struct {
	customers   []domain.Customer
	getErr      error
	lastFilter  domain.CustomerFilter
	getByIDSeen int64
}

func (f *fakeRepo) ListCustomers(_ context.Context, flt domain.CustomerFilter) ([]domain.Customer, error) {
	f.lastFilter = flt
	return f.customers, nil
}

func (f *fakeRepo) GetCustomerByID(_ context.Context, id int64) (domain.Customer, error) {
	f.getByIDSeen = id
	if f.getErr != nil {
		return domain.Customer{}, f.getErr
	}
	return domain.Customer{ID: id, Name: "X", Type: domain.TypeB2B}, nil
}

func TestListCustomers_DefaultLimit(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	if _, err := svc.ListCustomers(context.Background(), app.ListCustomersQuery{}); err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if f.lastFilter.Limit != 20 {
		t.Fatalf("limit mặc định = %d, want 20", f.lastFilter.Limit)
	}
	if f.lastFilter.AfterID != 0 {
		t.Fatalf("AfterID trang đầu = %d, want 0", f.lastFilter.AfterID)
	}
}

func TestListCustomers_LimitClamp(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	_, _ = svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 9999})
	if f.lastFilter.Limit != 100 {
		t.Fatalf("limit clamp = %d, want 100", f.lastFilter.Limit)
	}
}

func TestListCustomers_CursorDecodedToAfterID(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	cur := pagination.EncodeID(42)
	_, _ = svc.ListCustomers(context.Background(), app.ListCustomersQuery{Cursor: cur})
	if f.lastFilter.AfterID != 42 {
		t.Fatalf("AfterID = %d, want 42", f.lastFilter.AfterID)
	}
}

func TestListCustomers_BadCursor_Validation(t *testing.T) {
	svc := app.New(&fakeRepo{})
	_, err := svc.ListCustomers(context.Background(), app.ListCustomersQuery{Cursor: "###"})
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("kind = %v, want Validation", apperr.KindOf(err))
	}
}

func TestListCustomers_TypeFilterPassthrough(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	_, err := svc.ListCustomers(context.Background(), app.ListCustomersQuery{Type: "B2B"})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if f.lastFilter.Type != domain.TypeB2B {
		t.Fatalf("filter type = %q, want B2B", f.lastFilter.Type)
	}
}

func TestListCustomers_BadType_Validation(t *testing.T) {
	svc := app.New(&fakeRepo{})
	_, err := svc.ListCustomers(context.Background(), app.ListCustomersQuery{Type: "WHOLESALE"})
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("kind = %v, want Validation cho type lạ", apperr.KindOf(err))
	}
}

func TestListCustomers_NextCursor_WhenFull(t *testing.T) {
	custs := make([]domain.Customer, 20)
	for i := range custs {
		custs[i] = domain.Customer{ID: int64(i + 1)}
	}
	f := &fakeRepo{customers: custs}
	svc := app.New(f)
	res, _ := svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 20})
	if res.NextCursor == "" {
		t.Fatal("trang đầy phải có NextCursor")
	}
	gotID, _ := pagination.DecodeID(res.NextCursor)
	if gotID != 20 {
		t.Fatalf("NextCursor id = %d, want 20", gotID)
	}
}

func TestListCustomers_NoNextCursor_WhenPartial(t *testing.T) {
	f := &fakeRepo{customers: []domain.Customer{{ID: 1}, {ID: 2}}}
	svc := app.New(f)
	res, _ := svc.ListCustomers(context.Background(), app.ListCustomersQuery{Limit: 20})
	if res.NextCursor != "" {
		t.Fatalf("trang chưa đầy không được có NextCursor, got %q", res.NextCursor)
	}
}

func TestGetCustomer_Passthrough(t *testing.T) {
	f := &fakeRepo{}
	svc := app.New(f)
	c, err := svc.GetCustomer(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetCustomer: %v", err)
	}
	if c.ID != 7 || f.getByIDSeen != 7 {
		t.Fatalf("phải gọi repo với id 7: got id=%d seen=%d", c.ID, f.getByIDSeen)
	}
}

func TestGetCustomer_NotFoundPassthrough(t *testing.T) {
	f := &fakeRepo{getErr: apperr.NotFound("không có")}
	svc := app.New(f)
	_, err := svc.GetCustomer(context.Background(), 99)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("kind = %v, want NotFound", apperr.KindOf(err))
	}
}
