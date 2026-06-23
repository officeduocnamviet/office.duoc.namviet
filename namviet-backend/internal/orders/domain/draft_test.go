package domain_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

func qty(n int64) domain.Quantity { return domain.QuantityFromInt(n) }

// line tiện dựng DraftLine cho test (tiền từ chuỗi — KHÔNG float).
func line(t *testing.T, productID int64, q int64, uom, unitPrice, discount string) domain.DraftLine {
	t.Helper()
	return domain.DraftLine{
		ProductID: productID,
		Quantity:  qty(q),
		UOM:       uom,
		UnitPrice: mustMoney(t, unitPrice),
		Discount:  mustMoney(t, discount),
	}
}

// TestNewDraft_ComputesTotals: total_line mỗi dòng = quantity*unit_price - discount;
// total_amount (trước chiết khấu) = Σ quantity*unit_price; final_amount (sau chiết
// khấu) = Σ total_line. Tiền decimal, không float.
func TestNewDraft_ComputesTotals(t *testing.T) {
	d, err := domain.NewDraft(domain.DraftInput{
		OrderType: domain.OrderTypeB2B,
		Lines: []domain.DraftLine{
			line(t, 100, 3, "Hộp", "10000", "5000"), // 3*10000-5000 = 25000
			line(t, 200, 2, "Vỉ", "20000", "0"),     // 2*20000-0    = 40000
		},
	})
	if err != nil {
		t.Fatalf("NewDraft hợp lệ phải thành công: %v", err)
	}
	if len(d.Lines) != 2 {
		t.Fatalf("lines = %d, want 2", len(d.Lines))
	}
	if got := d.Lines[0].LineTotal.String(); got != "25000" {
		t.Errorf("line0 total = %q, want 25000", got)
	}
	if got := d.Lines[1].LineTotal.String(); got != "40000" {
		t.Errorf("line1 total = %q, want 40000", got)
	}
	// total_amount = Σ quantity*unit_price = 30000 + 40000 = 70000 (TRƯỚC chiết khấu).
	if got := d.TotalAmount.String(); got != "70000" {
		t.Errorf("total_amount = %q, want 70000", got)
	}
	// final_amount = Σ total_line = 25000 + 40000 = 65000 (SAU chiết khấu).
	if got := d.FinalAmount.String(); got != "65000" {
		t.Errorf("final_amount = %q, want 65000", got)
	}
}

func TestNewDraft_DefaultsOrderTypeB2C(t *testing.T) {
	d, err := domain.NewDraft(domain.DraftInput{
		Lines: []domain.DraftLine{line(t, 1, 1, "Cái", "1000", "0")},
	})
	if err != nil {
		t.Fatalf("NewDraft: %v", err)
	}
	if d.OrderType != domain.OrderTypeB2C {
		t.Errorf("order_type mặc định = %q, want B2C", d.OrderType)
	}
}

func TestNewDraft_Errors(t *testing.T) {
	cases := []struct {
		name string
		in   domain.DraftInput
	}{
		{"không có dòng", domain.DraftInput{Lines: nil}},
		{"quantity 0", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: domain.ZeroQty(), UOM: "Cái", UnitPrice: mustMoney(t, "1000"), Discount: money.Zero()},
		}}},
		{"unit_price âm", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: qty(1), UOM: "Cái", UnitPrice: mustMoney(t, "-1"), Discount: money.Zero()},
		}}},
		{"discount âm", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: qty(1), UOM: "Cái", UnitPrice: mustMoney(t, "1000"), Discount: mustMoney(t, "-1")},
		}}},
		{"product_id <= 0", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 0, Quantity: qty(1), UOM: "Cái", UnitPrice: mustMoney(t, "1000"), Discount: money.Zero()},
		}}},
		{"order_type sai", domain.DraftInput{OrderType: "X", Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: qty(1), UOM: "Cái", UnitPrice: mustMoney(t, "1000"), Discount: money.Zero()},
		}}},
		{"discount > thành tiền dòng", domain.DraftInput{Lines: []domain.DraftLine{
			{ProductID: 1, Quantity: qty(1), UOM: "Cái", UnitPrice: mustMoney(t, "1000"), Discount: mustMoney(t, "2000")},
		}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := domain.NewDraft(c.in); err == nil {
				t.Fatalf("NewDraft(%s) phải lỗi", c.name)
			}
		})
	}
}

// TestNewDraft_StatusIsPending: đơn mới luôn ở PENDING (chưa duyệt).
func TestNewDraft_StatusIsPending(t *testing.T) {
	d, err := domain.NewDraft(domain.DraftInput{
		Lines: []domain.DraftLine{line(t, 1, 1, "Cái", "1000", "0")},
	})
	if err != nil {
		t.Fatalf("NewDraft: %v", err)
	}
	if d.Status != domain.StatusPending {
		t.Errorf("status đơn mới = %q, want PENDING", d.Status)
	}
}
