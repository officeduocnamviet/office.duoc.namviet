package domain_test

import (
	"testing"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
)

func date(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

// TestSortFEFO_ExpiryAscending kiểm tra quy tắc FEFO: lô hết hạn TRƯỚC đứng TRƯỚC.
func TestSortFEFO_ExpiryAscending(t *testing.T) {
	batches := []domain.Batch{
		{InventoryBatchID: 1, BatchID: 1, ExpiryDate: date("2027-12-31")},
		{InventoryBatchID: 2, BatchID: 2, ExpiryDate: date("2026-01-15")},
		{InventoryBatchID: 3, BatchID: 3, ExpiryDate: date("2026-08-20")},
	}
	domain.SortFEFO(batches)
	want := []int64{2, 3, 1} // theo expiry tăng dần
	for i, w := range want {
		if batches[i].InventoryBatchID != w {
			t.Fatalf("FEFO sai ở vị trí %d: got InventoryBatchID=%d, want %d (order=%v)",
				i, batches[i].InventoryBatchID, w, order(batches))
		}
	}
}

// TestSortFEFO_TieBreakDeterministic kiểm tra tie-break khi CÙNG hạn dùng: theo
// BatchID rồi InventoryBatchID (ổn định, khớp ORDER BY của query).
func TestSortFEFO_TieBreakDeterministic(t *testing.T) {
	exp := date("2026-06-30")
	batches := []domain.Batch{
		{InventoryBatchID: 50, BatchID: 9, ExpiryDate: exp},
		{InventoryBatchID: 10, BatchID: 5, ExpiryDate: exp},
		{InventoryBatchID: 40, BatchID: 5, ExpiryDate: exp}, // cùng batch 5, id lớn hơn id=10
	}
	domain.SortFEFO(batches)
	// Cùng expiry → BatchID ASC (5 trước 9); trong batch 5 → InventoryBatchID ASC (10 trước 40).
	want := []int64{10, 40, 50}
	for i, w := range want {
		if batches[i].InventoryBatchID != w {
			t.Fatalf("tie-break sai ở vị trí %d: got %d, want %d (order=%v)",
				i, batches[i].InventoryBatchID, w, order(batches))
		}
	}
}

// TestSortFEFO_Empty không panic trên slice rỗng/nil.
func TestSortFEFO_Empty(t *testing.T) {
	domain.SortFEFO(nil)
	domain.SortFEFO([]domain.Batch{})
}

func order(bs []domain.Batch) []int64 {
	out := make([]int64, len(bs))
	for i, b := range bs {
		out[i] = b.InventoryBatchID
	}
	return out
}
