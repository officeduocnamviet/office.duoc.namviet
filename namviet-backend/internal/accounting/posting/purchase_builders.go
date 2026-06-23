package posting

import (
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// sourceTypePurchase là nhãn nguồn cho bút toán sinh từ một đơn MUA (PURCHASE /
// SUPPLIER_PAYMENT). SourceID = mã PO (purchase_orders.code).
const sourceTypePurchase = "purchase"

// BuildPurchaseReceipt dựng bút toán NHẬP KHO (mua hàng) cho MỘT sổ (book) — đối
// xứng BuildSaleRevenue. Hàm THUẦN, trả một JournalEntry ĐÃ CÂN.
//
// Vế NỢ:
//   - Inventory (1561) = inventoryCost (giá vốn hàng nhập, ex-VAT) — luôn có.
//   - InputVAT (133) = vatAmount — CHỈ khi sổ ghi VAT (INTERNAL theo cờ
//     InternalRecordsVAT; TAX luôn ghi) VÀ vatAmount > 0 (domain cấm dòng 0 đồng).
//
// Vế CÓ (1 dòng):
//   - Payable (331) = inventoryCost + (vatAmount NẾU sổ này ghi VAT) — công nợ NCC.
//
// Sổ INTERNAL truyền giá nhập thật; sổ TAX truyền số theo HĐ MUA. Hàm không quan tâm
// nguồn số — chỉ dựng bút toán cân theo book chỉ định (đối xứng dual-ledger chiều bán).
func (r Rules) BuildPurchaseReceipt(
	book domain.Book,
	inventoryCost money.Money,
	vatAmount money.Money,
	sourceID string,
	date time.Time,
) domain.JournalEntry {
	withVAT := r.recordsVAT(book) && vatAmount.IsPositive()

	// Vế có (phải trả NCC) = giá vốn + VAT (nếu sổ ghi VAT).
	payable := inventoryCost
	if withVAT {
		payable = payable.Add(vatAmount)
	}

	lines := []domain.EntryLine{
		{AccountCode: r.Accounts.Inventory, Debit: inventoryCost},
	}
	if withVAT {
		lines = append(lines, domain.EntryLine{AccountCode: r.Accounts.InputVAT, Debit: vatAmount})
	}
	lines = append(lines, domain.EntryLine{AccountCode: r.Accounts.Payable, Credit: payable})

	return domain.JournalEntry{
		Book:       book,
		EntryDate:  date,
		SourceType: sourceTypePurchase,
		SourceID:   sourceID,
		Memo:       "Nhập kho mua hàng",
		Lines:      lines,
	}
}

// BuildSupplierPayment dựng bút toán CHI TRẢ người bán (giảm công nợ): Dr 331 /
// Cr 111/112 (theo fundIsBank) = amount. Hàm THUẦN, trả entry ĐÃ CÂN. Đối xứng
// BuildPaymentIn (chiều bán). amount > 0.
func (r Rules) BuildSupplierPayment(
	book domain.Book,
	fundIsBank bool,
	amount money.Money,
	sourceID string,
	date time.Time,
) domain.JournalEntry {
	return domain.JournalEntry{
		Book:       book,
		EntryDate:  date,
		SourceType: sourceTypePurchase,
		SourceID:   sourceID,
		Memo:       "Chi trả người bán",
		Lines: []domain.EntryLine{
			{AccountCode: r.Accounts.Payable, Debit: amount},
			{AccountCode: r.fundAccount(fundIsBank), Credit: amount},
		},
	}
}

// PurchaseInput gom tham số một sự kiện MUA HÀNG (nhập kho) để dựng CẢ BỘ bút toán
// hai sổ — đối xứng SaleInput. INTERNAL theo giá nhập THỰC; TAX theo số HĐ MUA.
// HasInvoice=false (nhập không có HĐ mua) → KHÔNG sinh entry sổ TAX.
type PurchaseInput struct {
	SourceID string    // mã PO (purchase_orders.code)
	Date     time.Time // ngày ghi sổ (entry_date)

	// Sổ INTERNAL (giá nhập thực).
	InventoryCost money.Money // giá vốn hàng nhập ex-VAT (Σ qty × unit_cost)
	VATAmount     money.Money // VAT đầu vào thực

	// Sổ TAX (chỉ khi HasInvoice). Bỏ trống nếu không có HĐ mua.
	HasInvoice       bool        // true = có HĐ VAT mua → sinh entry sổ TAX
	TaxInventoryCost money.Money // giá vốn ex-VAT theo HĐ
	TaxVATAmount     money.Money // VAT đầu vào theo HĐ
}

// PurchaseEntries dựng TẤT CẢ bút toán của một sự kiện nhập kho mua hàng theo Rules,
// trả các JournalEntry ĐÃ CÂN sẵn sàng cho purchasing Post lần lượt trong tx. Thứ tự
// ổn định: [INTERNAL nhập kho, (nếu HĐ) TAX nhập kho].
//
// KHÁC chiều bán: chiều MUA KHÔNG có "giá vốn" riêng (nhập kho CHÍNH LÀ ghi tăng
// 1561) → mỗi sổ chỉ MỘT entry (Dr 1561 [+133]/Cr 331). Sổ TAX chỉ sinh khi có HĐ mua.
func (r Rules) PurchaseEntries(in PurchaseInput) []domain.JournalEntry {
	entries := make([]domain.JournalEntry, 0, 2)

	// Sổ INTERNAL (luôn có): nhập kho theo giá thực.
	entries = append(entries, r.BuildPurchaseReceipt(
		domain.BookInternal, in.InventoryCost, in.VATAmount, in.SourceID, in.Date,
	))

	// Sổ TAX (chỉ khi có HĐ VAT mua): nhập kho theo số HĐ.
	if in.HasInvoice {
		entries = append(entries, r.BuildPurchaseReceipt(
			domain.BookTax, in.TaxInventoryCost, in.TaxVATAmount, in.SourceID, in.Date,
		))
	}

	return entries
}
