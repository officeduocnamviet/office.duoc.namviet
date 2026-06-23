package posting_test

import (
	"testing"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/posting"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

var testDate = time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)

// lineByAccount tìm dòng theo account_code (giúp assert đúng vế/đúng số tiền).
func lineByAccount(e domain.JournalEntry, code string) (domain.EntryLine, bool) {
	for _, l := range e.Lines {
		if l.AccountCode == code {
			return l, true
		}
	}
	return domain.EntryLine{}, false
}

// assertBalanced ép bút toán hợp lệ + cân Σ qua chính domain.Validate (3 lớp
// phòng thủ — đây là lớp 1).
func assertBalanced(t *testing.T, e domain.JournalEntry, wantBook domain.Book) {
	t.Helper()
	if e.Book != wantBook {
		t.Fatalf("book = %q, muốn %q", e.Book, wantBook)
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("bút toán phải cân/hợp lệ, lỗi: %v", err)
	}
	if e.SourceType != "order" {
		t.Errorf("SourceType = %q, muốn \"order\"", e.SourceType)
	}
}

// assertDebit / assertCredit kiểm một dòng đúng vế + đúng số tiền, vế kia = 0.
func assertDebit(t *testing.T, e domain.JournalEntry, code string, want money.Money) {
	t.Helper()
	l, ok := lineByAccount(e, code)
	if !ok {
		t.Fatalf("thiếu dòng tài khoản %s", code)
	}
	if !l.Debit.Equal(want) {
		t.Errorf("Dr %s = %s, muốn %s", code, l.Debit, want)
	}
	if !l.Credit.IsZero() {
		t.Errorf("Dr %s không được có vế Có (=%s)", code, l.Credit)
	}
}

func assertCredit(t *testing.T, e domain.JournalEntry, code string, want money.Money) {
	t.Helper()
	l, ok := lineByAccount(e, code)
	if !ok {
		t.Fatalf("thiếu dòng tài khoản %s", code)
	}
	if !l.Credit.Equal(want) {
		t.Errorf("Cr %s = %s, muốn %s", code, l.Credit, want)
	}
	if !l.Debit.IsZero() {
		t.Errorf("Cr %s không được có vế Nợ (=%s)", code, l.Debit)
	}
}

// --- BuildSaleRevenue ---

// INTERNAL, bán chịu, có VAT (cờ InternalRecordsVAT=true): Dr 131=110000 /
// Cr 511=100000 / Cr 3331=10000. Cân.
func TestBuildSaleRevenue_Internal_OnCredit_WithVAT(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildSaleRevenue(domain.BookInternal, true, false,
		money.FromInt(100000), money.FromInt(10000), "HD001", testDate)

	assertBalanced(t, e, domain.BookInternal)
	if len(e.Lines) != 3 {
		t.Fatalf("muốn 3 dòng (131,511,3331), có %d", len(e.Lines))
	}
	assertDebit(t, e, "131", money.FromInt(110000))
	assertCredit(t, e, "511", money.FromInt(100000))
	assertCredit(t, e, "3331", money.FromInt(10000))
	if e.SourceID != "HD001" {
		t.Errorf("SourceID = %q, muốn HD001", e.SourceID)
	}
}

// INTERNAL, thu ngay tiền mặt: vế nợ là 111 (không 131).
func TestBuildSaleRevenue_Internal_CashNow(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildSaleRevenue(domain.BookInternal, false, false,
		money.FromInt(100000), money.FromInt(10000), "HD002", testDate)

	assertBalanced(t, e, domain.BookInternal)
	assertDebit(t, e, "111", money.FromInt(110000))
	if _, ok := lineByAccount(e, "131"); ok {
		t.Error("thu ngay không được dùng 131 (phải thu)")
	}
}

// INTERNAL, thu ngay chuyển khoản: vế nợ là 112 (ngân hàng).
func TestBuildSaleRevenue_Internal_BankNow(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildSaleRevenue(domain.BookInternal, false, true,
		money.FromInt(100000), money.FromInt(10000), "HD003", testDate)

	assertBalanced(t, e, domain.BookInternal)
	assertDebit(t, e, "112", money.FromInt(110000))
	if _, ok := lineByAccount(e, "111"); ok {
		t.Error("thu ngân hàng không được dùng 111 (tiền mặt)")
	}
}

// INTERNAL với cờ InternalRecordsVAT=false: KHÔNG có 3331; vế nợ = doanh thu (chưa
// VAT). Cân: Dr 131=100000 / Cr 511=100000.
func TestBuildSaleRevenue_Internal_NoVATFlag(t *testing.T) {
	r := posting.DefaultRules
	r.InternalRecordsVAT = false
	e := r.BuildSaleRevenue(domain.BookInternal, true, false,
		money.FromInt(100000), money.FromInt(10000), "HD004", testDate)

	assertBalanced(t, e, domain.BookInternal)
	if len(e.Lines) != 2 {
		t.Fatalf("INTERNAL không ghi VAT muốn 2 dòng, có %d", len(e.Lines))
	}
	assertDebit(t, e, "131", money.FromInt(100000)) // KHÔNG cộng VAT
	assertCredit(t, e, "511", money.FromInt(100000))
	if _, ok := lineByAccount(e, "3331"); ok {
		t.Error("cờ InternalRecordsVAT=false thì không được có 3331")
	}
}

// TAX LUÔN ghi VAT bất kể cờ InternalRecordsVAT (cờ chỉ áp INTERNAL).
func TestBuildSaleRevenue_Tax_AlwaysVAT(t *testing.T) {
	r := posting.DefaultRules
	r.InternalRecordsVAT = false // cờ này không ảnh hưởng sổ TAX
	e := r.BuildSaleRevenue(domain.BookTax, true, false,
		money.FromInt(90000), money.FromInt(9000), "HD005", testDate)

	assertBalanced(t, e, domain.BookTax)
	if len(e.Lines) != 3 {
		t.Fatalf("TAX muốn 3 dòng (có VAT), có %d", len(e.Lines))
	}
	assertDebit(t, e, "131", money.FromInt(99000))
	assertCredit(t, e, "511", money.FromInt(90000))
	assertCredit(t, e, "3331", money.FromInt(9000))
}

// VAT = 0 (không thuế): không sinh dòng 3331 dù sổ "ghi VAT" — domain cấm dòng 0
// đồng. Vẫn cân: Dr=Cr=doanh thu.
func TestBuildSaleRevenue_ZeroVAT_NoVATLine(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildSaleRevenue(domain.BookInternal, true, false,
		money.FromInt(100000), money.Zero(), "HD006", testDate)

	assertBalanced(t, e, domain.BookInternal)
	if len(e.Lines) != 2 {
		t.Fatalf("VAT=0 muốn 2 dòng, có %d", len(e.Lines))
	}
	assertDebit(t, e, "131", money.FromInt(100000))
	if _, ok := lineByAccount(e, "3331"); ok {
		t.Error("VAT=0 không được sinh dòng 3331 (dòng 0 đồng vi phạm domain)")
	}
}

// --- BuildCOGS ---

// INTERNAL luôn ghi giá vốn: Dr 632 / Cr 1561 = cogs. Cân.
func TestBuildCOGS_Internal(t *testing.T) {
	r := posting.DefaultRules
	e, ok := r.BuildCOGS(domain.BookInternal, money.FromInt(70000), "HD010", testDate)
	if !ok {
		t.Fatal("INTERNAL phải ghi giá vốn (ok=true)")
	}
	assertBalanced(t, e, domain.BookInternal)
	assertDebit(t, e, "632", money.FromInt(70000))
	assertCredit(t, e, "1561", money.FromInt(70000))
}

// TAX với TaxRecordsCOGS=false (mặc định): KHÔNG sinh bút toán (ok=false).
func TestBuildCOGS_Tax_DefaultSkips(t *testing.T) {
	r := posting.DefaultRules // TaxRecordsCOGS=false
	_, ok := r.BuildCOGS(domain.BookTax, money.FromInt(70000), "HD011", testDate)
	if ok {
		t.Fatal("mặc định sổ TAX không ghi giá vốn (ok phải=false)")
	}
}

// TAX với cờ TaxRecordsCOGS=true: CÓ sinh bút toán giá vốn ở sổ TAX.
func TestBuildCOGS_Tax_FlagOn(t *testing.T) {
	r := posting.DefaultRules
	r.TaxRecordsCOGS = true
	e, ok := r.BuildCOGS(domain.BookTax, money.FromInt(65000), "HD012", testDate)
	if !ok {
		t.Fatal("cờ TaxRecordsCOGS=true thì sổ TAX phải ghi giá vốn")
	}
	assertBalanced(t, e, domain.BookTax)
	assertDebit(t, e, "632", money.FromInt(65000))
	assertCredit(t, e, "1561", money.FromInt(65000))
}

// --- BuildPaymentIn ---

// Thu tiền mặt: Dr 111 / Cr 131 = amount. Cân.
func TestBuildPaymentIn_Cash(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildPaymentIn(domain.BookInternal, false, money.FromInt(50000), "HD020", testDate)
	assertBalanced(t, e, domain.BookInternal)
	assertDebit(t, e, "111", money.FromInt(50000))
	assertCredit(t, e, "131", money.FromInt(50000))
}

// Thu chuyển khoản: Dr 112 / Cr 131.
func TestBuildPaymentIn_Bank(t *testing.T) {
	r := posting.DefaultRules
	e := r.BuildPaymentIn(domain.BookTax, true, money.FromInt(50000), "HD021", testDate)
	assertBalanced(t, e, domain.BookTax)
	assertDebit(t, e, "112", money.FromInt(50000))
	assertCredit(t, e, "131", money.FromInt(50000))
}

// --- SaleEntries (tổng hợp 2 sổ) ---

// Bán B2B có HĐ, mặc định: INTERNAL{doanh thu+giá vốn}, TAX{doanh thu} = 3 entry.
// Mỗi entry cân; đúng book; TAX không có giá vốn.
func TestSaleEntries_B2B_Default(t *testing.T) {
	r := posting.DefaultRules
	entries := r.SaleEntries(posting.SaleInput{
		SourceID:        "HD100",
		Date:            testDate,
		OnCredit:        true,
		RevenueExVAT:    money.FromInt(100000),
		VATAmount:       money.FromInt(10000),
		COGS:            money.FromInt(70000),
		HasInvoice:      true,
		TaxRevenueExVAT: money.FromInt(90000), // giá HĐ khác giá thực
		TaxVATAmount:    money.FromInt(9000),
		TaxCOGS:         money.FromInt(65000),
	})

	if len(entries) != 3 {
		t.Fatalf("muốn 3 entry (INTERNAL doanh thu+giá vốn, TAX doanh thu), có %d", len(entries))
	}
	// Tất cả phải cân.
	for i, e := range entries {
		if err := e.Validate(); err != nil {
			t.Fatalf("entry[%d] phải cân: %v", i, err)
		}
	}
	// entry[0] INTERNAL doanh thu, entry[1] INTERNAL giá vốn, entry[2] TAX doanh thu.
	assertBalanced(t, entries[0], domain.BookInternal)
	assertDebit(t, entries[0], "131", money.FromInt(110000))
	assertBalanced(t, entries[1], domain.BookInternal)
	assertDebit(t, entries[1], "632", money.FromInt(70000))
	assertBalanced(t, entries[2], domain.BookTax)
	assertDebit(t, entries[2], "131", money.FromInt(99000))
	assertCredit(t, entries[2], "511", money.FromInt(90000))

	// 2 sổ độc lập: doanh thu INTERNAL (110000) ≠ TAX (99000).
	if entries[0].TotalDebit().Equal(entries[2].TotalDebit()) {
		t.Error("2 sổ phải độc lập số tiền (giá thực vs giá HĐ)")
	}
}

// POS B2C thu ngay, không HĐ: chỉ sổ INTERNAL{doanh thu thu-ngay + giá vốn} = 2
// entry, không có entry sổ TAX.
func TestSaleEntries_POS_NoInvoice(t *testing.T) {
	r := posting.DefaultRules
	entries := r.SaleEntries(posting.SaleInput{
		SourceID:     "HD101",
		Date:         testDate,
		OnCredit:     false,
		FundIsBank:   false,
		RevenueExVAT: money.FromInt(50000),
		VATAmount:    money.FromInt(5000),
		COGS:         money.FromInt(30000),
		HasInvoice:   false,
	})

	if len(entries) != 2 {
		t.Fatalf("POS không HĐ muốn 2 entry (chỉ INTERNAL), có %d", len(entries))
	}
	for _, e := range entries {
		if e.Book != domain.BookInternal {
			t.Errorf("POS không HĐ chỉ được sổ INTERNAL, gặp %q", e.Book)
		}
		if err := e.Validate(); err != nil {
			t.Fatalf("entry phải cân: %v", err)
		}
	}
	// Thu ngay tiền mặt → 111, không 131.
	assertDebit(t, entries[0], "111", money.FromInt(55000))
}

// Cờ TaxRecordsCOGS=true: bán có HĐ sinh 4 entry (cả giá vốn sổ TAX).
func TestSaleEntries_TaxCOGSFlagOn(t *testing.T) {
	r := posting.DefaultRules
	r.TaxRecordsCOGS = true
	entries := r.SaleEntries(posting.SaleInput{
		SourceID:        "HD102",
		Date:            testDate,
		OnCredit:        true,
		RevenueExVAT:    money.FromInt(100000),
		VATAmount:       money.FromInt(10000),
		COGS:            money.FromInt(70000),
		HasInvoice:      true,
		TaxRevenueExVAT: money.FromInt(100000),
		TaxVATAmount:    money.FromInt(10000),
		TaxCOGS:         money.FromInt(70000),
	})
	if len(entries) != 4 {
		t.Fatalf("cờ TaxRecordsCOGS=true muốn 4 entry, có %d", len(entries))
	}
	assertBalanced(t, entries[3], domain.BookTax)
	assertDebit(t, entries[3], "632", money.FromInt(70000))
}
