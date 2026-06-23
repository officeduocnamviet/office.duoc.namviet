package posting

import (
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// sourceTypeOrder là nhãn nguồn cho mọi bút toán sinh từ một đơn hàng (SALE /
// COGS / PAYMENT_IN). SourceID = mã đơn (orders.code, vd "HD123").
const sourceTypeOrder = "order"

// BuildSaleRevenue dựng bút toán GHI NHẬN DOANH THU cho MỘT sổ (book) từ một sự
// kiện bán hàng. Hàm THUẦN, trả một JournalEntry ĐÃ CÂN (Σnợ = Σcó).
//
// Vế NỢ (1 dòng):
//   - onCredit=true  → AR (131, bán chịu B2B / công nợ).
//   - onCredit=false → Cash (111) hoặc Bank (112) theo fundIsBank (thu ngay POS).
//   - Số nợ = revenueExVAT + (vatAmount NẾU sổ này ghi VAT).
//
// Vế CÓ:
//   - Revenue (511) = revenueExVAT (luôn có).
//   - VATOutput (3331) = vatAmount — CHỈ khi sổ ghi VAT (INTERNAL theo cờ
//     InternalRecordsVAT; TAX luôn ghi) VÀ vatAmount > 0 (domain cấm dòng 0 đồng).
//
// Ghi chú số tiền: với sổ INTERNAL truyền giá bán thật + VAT thực; với sổ TAX
// truyền số theo HÓA ĐƠN (invoice_price + VAT HĐ). Hàm không quan tâm nguồn số —
// chỉ dựng bút toán cân theo book được chỉ định.
func (r Rules) BuildSaleRevenue(
	book domain.Book,
	onCredit bool,
	fundIsBank bool,
	revenueExVAT money.Money,
	vatAmount money.Money,
	sourceID string,
	date time.Time,
) domain.JournalEntry {
	withVAT := r.recordsVAT(book) && vatAmount.IsPositive()

	// Vế nợ (tiền/phải thu) = doanh thu + VAT (nếu sổ ghi VAT).
	debitAmount := revenueExVAT
	if withVAT {
		debitAmount = debitAmount.Add(vatAmount)
	}
	debitAccount := r.fundAccount(fundIsBank)
	if onCredit {
		debitAccount = r.Accounts.AR
	}

	lines := []domain.EntryLine{
		{AccountCode: debitAccount, Debit: debitAmount},
		{AccountCode: r.Accounts.Revenue, Credit: revenueExVAT},
	}
	if withVAT {
		lines = append(lines, domain.EntryLine{AccountCode: r.Accounts.VATOutput, Credit: vatAmount})
	}

	return domain.JournalEntry{
		Book:       book,
		EntryDate:  date,
		SourceType: sourceTypeOrder,
		SourceID:   sourceID,
		Memo:       "Ghi nhận doanh thu bán hàng",
		Lines:      lines,
	}
}

// BuildCOGS dựng bút toán GIÁ VỐN HÀNG BÁN cho MỘT sổ: Dr 632 / Cr 1561 = cogs.
// Hàm THUẦN. Trả (entry, true) nếu sổ này ghi giá vốn; (zero, false) nếu KHÔNG
// (vd sổ TAX khi TaxRecordsCOGS=false) → caller bỏ qua, không post.
//
// cogs là tổng giá vốn (Σ inbound_price các lô FEFO đã xuất, sổ INTERNAL; cơ sở
// hóa đơn nếu sổ TAX có ghi). cogs phải > 0 (không có giá vốn → không sự kiện).
func (r Rules) BuildCOGS(
	book domain.Book,
	cogs money.Money,
	sourceID string,
	date time.Time,
) (domain.JournalEntry, bool) {
	if !r.recordsCOGS(book) {
		return domain.JournalEntry{}, false
	}
	return domain.JournalEntry{
		Book:       book,
		EntryDate:  date,
		SourceType: sourceTypeOrder,
		SourceID:   sourceID,
		Memo:       "Ghi nhận giá vốn hàng bán",
		Lines: []domain.EntryLine{
			{AccountCode: r.Accounts.COGS, Debit: cogs},
			{AccountCode: r.Accounts.Inventory, Credit: cogs},
		},
	}, true
}

// BuildPaymentIn dựng bút toán THU TIỀN của khách (giảm phải thu): Dr 111/112
// (theo fundIsBank) / Cr 131 = amount. Hàm THUẦN, trả entry ĐÃ CÂN.
//
// Dùng cho sự kiện PAYMENT_IN (thu công nợ B2B). amount > 0. Sổ do caller quyết
// (phiếu thu INTERNAL/BOTH → sổ INTERNAL; TAX/BOTH → sổ TAX) — P4 gọi 1 lần/sổ.
func (r Rules) BuildPaymentIn(
	book domain.Book,
	fundIsBank bool,
	amount money.Money,
	sourceID string,
	date time.Time,
) domain.JournalEntry {
	return domain.JournalEntry{
		Book:       book,
		EntryDate:  date,
		SourceType: sourceTypeOrder,
		SourceID:   sourceID,
		Memo:       "Thu tiền khách hàng",
		Lines: []domain.EntryLine{
			{AccountCode: r.fundAccount(fundIsBank), Debit: amount},
			{AccountCode: r.Accounts.AR, Credit: amount},
		},
	}
}

// SaleInput gom tham số một sự kiện BÁN HÀNG để dựng CẢ BỘ bút toán hai sổ. Tiện
// cho P4: gọi 1 lần, nhận tất cả entry đã cân rồi Post lần lượt trong tx.
//
// INTERNAL theo giá/VAT/giá-vốn THỰC; TAX theo số HÓA ĐƠN. HasInvoice=false (POS
// B2C không xuất HĐ) → KHÔNG sinh entry sổ TAX.
type SaleInput struct {
	SourceID   string    // mã đơn (orders.code)
	Date       time.Time // ngày ghi sổ (entry_date)
	OnCredit   bool      // true = bán chịu (Dr 131); false = thu ngay (Dr 111/112)
	FundIsBank bool      // khi thu ngay: true = ngân hàng (112), false = tiền mặt (111)

	// Sổ INTERNAL (giá thực).
	RevenueExVAT money.Money // doanh thu chưa VAT theo giá bán thật
	VATAmount    money.Money // VAT đầu ra thực
	COGS         money.Money // giá vốn thật (Σ lô FEFO)

	// Sổ TAX (chỉ khi HasInvoice). Bỏ trống nếu không xuất HĐ.
	HasInvoice      bool        // true = có HĐ VAT → sinh entry sổ TAX
	TaxRevenueExVAT money.Money // doanh thu chưa VAT theo HĐ
	TaxVATAmount    money.Money // VAT đầu ra theo HĐ
	TaxCOGS         money.Money // giá vốn cơ sở HĐ (chỉ dùng nếu TaxRecordsCOGS=true)
}

// SaleEntries dựng TẤT CẢ bút toán của một sự kiện bán hàng theo Rules, trả các
// JournalEntry ĐÃ CÂN sẵn sàng cho P4 Post lần lượt trong tx. Thứ tự ổn định:
// [INTERNAL doanh thu, INTERNAL giá vốn, (nếu HĐ) TAX doanh thu, (nếu cờ) TAX giá vốn].
//
// Bút toán nào sổ không ghi (vd COGS sổ TAX khi TaxRecordsCOGS=false) tự động bị
// bỏ qua — không xuất hiện trong kết quả.
func (r Rules) SaleEntries(in SaleInput) []domain.JournalEntry {
	entries := make([]domain.JournalEntry, 0, 4)

	// Sổ INTERNAL (luôn có): doanh thu + giá vốn theo giá thực.
	entries = append(entries, r.BuildSaleRevenue(
		domain.BookInternal, in.OnCredit, in.FundIsBank, in.RevenueExVAT, in.VATAmount, in.SourceID, in.Date,
	))
	if cogs, ok := r.BuildCOGS(domain.BookInternal, in.COGS, in.SourceID, in.Date); ok {
		entries = append(entries, cogs)
	}

	// Sổ TAX (chỉ khi có HĐ VAT): doanh thu theo HĐ + giá vốn nếu cờ bật.
	if in.HasInvoice {
		entries = append(entries, r.BuildSaleRevenue(
			domain.BookTax, in.OnCredit, in.FundIsBank, in.TaxRevenueExVAT, in.TaxVATAmount, in.SourceID, in.Date,
		))
		if cogs, ok := r.BuildCOGS(domain.BookTax, in.TaxCOGS, in.SourceID, in.Date); ok {
			entries = append(entries, cogs)
		}
	}

	return entries
}
