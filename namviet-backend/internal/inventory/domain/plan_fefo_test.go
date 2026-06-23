package domain_test

import (
	"errors"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

func qty(t *testing.T, s string) domain.Quantity {
	t.Helper()
	q, err := domain.QuantityFromString(s)
	if err != nil {
		t.Fatalf("qty %q: %v", s, err)
	}
	return q
}

// batch dựng một lô đã ghép tồn + giá nhập (input cho PlanFEFO; thứ tự FEFO do
// caller bảo đảm — PlanFEFO tiêu thụ THEO THỨ TỰ slice).
func batch(invBatchID, batchID int64, quantity domain.Quantity, price int64) domain.Batch {
	return domain.Batch{
		InventoryBatchID: invBatchID,
		BatchID:          batchID,
		Quantity:         quantity,
		InboundPrice:     money.FromInt(price),
	}
}

// TestPlanFEFO_SingleBatchExact: nhu cầu vừa khít một lô → tiêu thụ trọn lô, 1 dòng.
func TestPlanFEFO_SingleBatchExact(t *testing.T) {
	avail := []domain.Batch{batch(9001, 501, qty(t, "10"), 9000)}
	plan, err := domain.PlanFEFO(avail, qty(t, "10"))
	if err != nil {
		t.Fatalf("PlanFEFO đủ phải thành công: %v", err)
	}
	if len(plan) != 1 {
		t.Fatalf("plan = %d dòng, want 1", len(plan))
	}
	if !plan[0].Quantity.Equal(qty(t, "10")) {
		t.Fatalf("tiêu thụ = %s, want 10", plan[0].Quantity)
	}
	if plan[0].InventoryBatchID != 9001 || plan[0].BatchID != 501 {
		t.Fatalf("id lô sai: inv=%d batch=%d", plan[0].InventoryBatchID, plan[0].BatchID)
	}
	if !plan[0].InboundPrice.Equal(money.FromInt(9000)) {
		t.Fatalf("inbound_price phải mang theo để post COGS: %s", plan[0].InboundPrice)
	}
}

// TestPlanFEFO_MultiBatch_OrderPreserved: nhu cầu trải nhiều lô → tiêu thụ ĐÚNG
// THỨ TỰ slice (đã FEFO): lô đầu trọn vẹn, lô kế tiếp tục.
func TestPlanFEFO_MultiBatch_OrderPreserved(t *testing.T) {
	avail := []domain.Batch{
		batch(9001, 501, qty(t, "10"), 9000), // hạn sớm nhất (FEFO trước)
		batch(9002, 502, qty(t, "20"), 9100),
		batch(9003, 503, qty(t, "5"), 9050),
	}
	plan, err := domain.PlanFEFO(avail, qty(t, "25"))
	if err != nil {
		t.Fatalf("PlanFEFO đủ phải thành công: %v", err)
	}
	// 10 (lô1 trọn) + 15 (lô2 một phần) = 25; lô3 không chạm.
	if len(plan) != 2 {
		t.Fatalf("plan = %d dòng, want 2 (lô1 trọn + lô2 một phần)", len(plan))
	}
	if plan[0].InventoryBatchID != 9001 || !plan[0].Quantity.Equal(qty(t, "10")) {
		t.Fatalf("dòng 1 sai: inv=%d qty=%s", plan[0].InventoryBatchID, plan[0].Quantity)
	}
	if plan[1].InventoryBatchID != 9002 || !plan[1].Quantity.Equal(qty(t, "15")) {
		t.Fatalf("dòng 2 (lô cuối một phần) sai: inv=%d qty=%s", plan[1].InventoryBatchID, plan[1].Quantity)
	}
	// Tổng tiêu thụ = nhu cầu.
	total := domain.ZeroQty()
	for _, c := range plan {
		total = total.Add(c.Quantity)
	}
	if !total.Equal(qty(t, "25")) {
		t.Fatalf("Σ tiêu thụ = %s, want 25", total)
	}
}

// TestPlanFEFO_LastBatchPartial: lô cuối tiêu thụ MỘT PHẦN (tồn lô > nhu cầu còn
// lại) — phải trả đúng phần cần, KHÔNG trả trọn lô.
func TestPlanFEFO_LastBatchPartial(t *testing.T) {
	avail := []domain.Batch{
		batch(9001, 501, qty(t, "3"), 9000),
		batch(9002, 502, qty(t, "100"), 9100), // tồn lớn, chỉ lấy một phần
	}
	plan, err := domain.PlanFEFO(avail, qty(t, "7.5"))
	if err != nil {
		t.Fatalf("PlanFEFO đủ phải thành công: %v", err)
	}
	if len(plan) != 2 {
		t.Fatalf("plan = %d dòng, want 2", len(plan))
	}
	if !plan[0].Quantity.Equal(qty(t, "3")) {
		t.Fatalf("lô đầu phải trọn 3: %s", plan[0].Quantity)
	}
	// 7.5 - 3 = 4.5 ở lô cuối (một phần, decimal — không float).
	if !plan[1].Quantity.Equal(qty(t, "4.5")) {
		t.Fatalf("lô cuối tiêu thụ một phần = %s, want 4.5", plan[1].Quantity)
	}
}

// TestPlanFEFO_Insufficient: tổng tồn < nhu cầu → ErrInsufficientStock, KHÔNG plan
// (không cho tồn âm).
func TestPlanFEFO_Insufficient(t *testing.T) {
	avail := []domain.Batch{
		batch(9001, 501, qty(t, "10"), 9000),
		batch(9002, 502, qty(t, "5"), 9100),
	}
	plan, err := domain.PlanFEFO(avail, qty(t, "20")) // cần 20, chỉ có 15
	if !errors.Is(err, domain.ErrInsufficientStock) {
		t.Fatalf("thiếu tồn phải trả ErrInsufficientStock, got %v", err)
	}
	if plan != nil {
		t.Fatalf("thiếu tồn KHÔNG được trả plan, got %d dòng", len(plan))
	}
}

// TestPlanFEFO_EmptyAvailable: không có lô nào còn tồn mà nhu cầu > 0 → thiếu.
func TestPlanFEFO_EmptyAvailable(t *testing.T) {
	_, err := domain.PlanFEFO(nil, qty(t, "1"))
	if !errors.Is(err, domain.ErrInsufficientStock) {
		t.Fatalf("không có lô + cần > 0 phải ErrInsufficientStock, got %v", err)
	}
}

// TestPlanFEFO_ZeroNeed: nhu cầu 0 → plan rỗng, không lỗi (no-op hợp lệ).
func TestPlanFEFO_ZeroNeed(t *testing.T) {
	avail := []domain.Batch{batch(9001, 501, qty(t, "10"), 9000)}
	plan, err := domain.PlanFEFO(avail, domain.ZeroQty())
	if err != nil {
		t.Fatalf("nhu cầu 0 không được lỗi: %v", err)
	}
	if len(plan) != 0 {
		t.Fatalf("nhu cầu 0 → plan rỗng, got %d dòng", len(plan))
	}
}
