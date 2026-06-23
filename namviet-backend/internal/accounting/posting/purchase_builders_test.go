package posting_test

import (
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/posting"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// TestPurchaseReceipt_Internal_WithVAT: nhập kho có VAT, sổ INTERNAL (InternalRecordsVAT=true)
// → Dr 1561 (giá vốn) + Dr 133 (VAT) / Cr 331 (tổng). Cân Σ.
func TestPurchaseReceipt_Internal_WithVAT(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildPurchaseReceipt(domain.BookInternal, money.FromInt(100000), money.FromInt(8000), "PO001", testDate)
	assertBalancedPurchase(t, e, domain.BookInternal)
	assertDebit(t, e, "1561", money.FromInt(100000))
	assertDebit(t, e, "133", money.FromInt(8000))
	assertCredit(t, e, "331", money.FromInt(108000))
}

// TestPurchaseReceipt_NoVAT: vat=0 → KHÔNG dòng 133; Dr 1561 / Cr 331 = giá vốn.
func TestPurchaseReceipt_NoVAT(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildPurchaseReceipt(domain.BookInternal, money.FromInt(50000), money.Zero(), "PO002", testDate)
	assertBalancedPurchase(t, e, domain.BookInternal)
	if _, ok := lineByAccount(e, "133"); ok {
		t.Error("vat=0 KHÔNG được có dòng 133")
	}
	assertDebit(t, e, "1561", money.FromInt(50000))
	assertCredit(t, e, "331", money.FromInt(50000))
}

// TestPurchaseEntries_NoInvoice_OnlyInternal: không HĐ mua → chỉ 1 entry sổ INTERNAL.
func TestPurchaseEntries_NoInvoice_OnlyInternal(t *testing.T) {
	r := posting.DefaultRules
	entries := r.PurchaseEntries(posting.PurchaseInput{
		SourceID:      "PO003",
		Date:          testDate,
		InventoryCost: money.FromInt(100000),
		VATAmount:     money.FromInt(8000),
		HasInvoice:    false,
	})
	if len(entries) != 1 {
		t.Fatalf("không HĐ mua → 1 entry (INTERNAL); got %d", len(entries))
	}
	assertBalancedPurchase(t, entries[0], domain.BookInternal)
}

// TestPurchaseEntries_WithInvoice_BothBooks: có HĐ mua → 2 entry (INTERNAL giá thực +
// TAX giá HĐ). Mỗi entry cân Σ; số tiền KHÁC nhau đúng dual-ledger.
func TestPurchaseEntries_WithInvoice_BothBooks(t *testing.T) {
	r := posting.DefaultRules
	entries := r.PurchaseEntries(posting.PurchaseInput{
		SourceID:         "PO004",
		Date:             testDate,
		InventoryCost:    money.FromInt(100000),
		VATAmount:        money.FromInt(8000),
		HasInvoice:       true,
		TaxInventoryCost: money.FromInt(90000), // HĐ ghi giá khác giá thực
		TaxVATAmount:     money.FromInt(9000),
	})
	if len(entries) != 2 {
		t.Fatalf("có HĐ mua → 2 entry; got %d", len(entries))
	}
	assertBalancedPurchase(t, entries[0], domain.BookInternal)
	assertBalancedPurchase(t, entries[1], domain.BookTax)
	// Sổ INTERNAL giá thực.
	assertDebit(t, entries[0], "1561", money.FromInt(100000))
	assertCredit(t, entries[0], "331", money.FromInt(108000))
	// Sổ TAX giá HĐ.
	assertDebit(t, entries[1], "1561", money.FromInt(90000))
	assertCredit(t, entries[1], "331", money.FromInt(99000))
}

// TestBuildSupplierPayment: chi trả NCC Dr 331 / Cr 111 (tiền mặt) hoặc 112 (NH). Cân.
func TestBuildSupplierPayment(t *testing.T) {
	r := posting.DefaultRules
	cash := r.BuildSupplierPayment(domain.BookInternal, false, money.FromInt(108000), "PO005", testDate)
	assertBalancedPurchase(t, cash, domain.BookInternal)
	assertDebit(t, cash, "331", money.FromInt(108000))
	assertCredit(t, cash, "111", money.FromInt(108000))

	bank := r.BuildSupplierPayment(domain.BookInternal, true, money.FromInt(108000), "PO005", testDate)
	assertCredit(t, bank, "112", money.FromInt(108000))
}

// assertBalancedPurchase ép bút toán mua cân Σ + SourceType="purchase".
func assertBalancedPurchase(t *testing.T, e domain.JournalEntry, wantBook domain.Book) {
	t.Helper()
	if e.Book != wantBook {
		t.Fatalf("book = %q, muốn %q", e.Book, wantBook)
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("bút toán phải cân/hợp lệ, lỗi: %v", err)
	}
	if e.SourceType != "purchase" {
		t.Errorf("SourceType = %q, muốn \"purchase\"", e.SourceType)
	}
}
