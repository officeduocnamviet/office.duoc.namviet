// Package domain là LÕI THUẦN của bounded context accounting: sổ kế toán kép
// (double-entry) theo Thông tư 133. Entity bút toán (JournalEntry) + dòng bút
// toán (EntryLine) + value object Book, cùng bất biến cân sổ THUẦN (Validate).
// KHÔNG import pgx/http/huma/framework (ARCHITECTURE.md §3) — chỉ stdlib + shared
// kernel trung lập (common/money). Phụ thuộc đi một chiều: adapters → app → domain.
//
// BẤT BIẾN KẾ TOÁN (đây là chỗ duy nhất được "kỹ" — ARCHITECTURE.md §7):
//   - Mỗi EntryLine đúng MỘT vế (debit XOR credit) > 0, vế còn lại = 0; không vế
//     nào âm.
//   - Σdebit = Σcredit trong một JournalEntry (cộng bằng money decimal, KHÔNG
//     float). Ép Ở ĐÂY (domain) + service + constraint trigger DB (phòng thủ 3 lớp).
//   - Append-only: bút toán đã ghi KHÔNG sửa/xoá ở tầng app; sửa = bút toán đảo.
//   - 2 sổ INTERNAL/TAX TÁCH BIỆT, KHÔNG sync: mỗi entry mang đúng 1 Book; cùng
//     một nghiệp vụ có thể sinh entry ở CẢ HAI sổ với SỐ TIỀN KHÁC nhau (giá thực
//     vs giá hoá đơn). Validate chỉ cân TỪNG entry — không ràng buộc chéo 2 sổ.
package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// Book là sổ kế toán mà một bút toán thuộc về. 2 sổ song song, KHÔNG sync.
type Book string

const (
	// BookInternal — sổ thực tế: doanh thu theo giá bán thật, giá vốn theo lô
	// FEFO. Dùng cho P&L nội bộ + công nợ thật.
	BookInternal Book = "INTERNAL"
	// BookTax — sổ thuế: chỉ ghi giao dịch CÓ hoá đơn VAT, theo giá trên hoá đơn.
	// Dùng cho báo cáo thuế.
	BookTax Book = "TAX"
)

// Valid trả true nếu b là một sổ hợp lệ (khớp CHECK book ở DB).
func (b Book) Valid() bool { return b == BookInternal || b == BookTax }

// EntryLine là một dòng của bút toán: ghi NỢ (Debit) hoặc CÓ (Credit) một tài
// khoản. Đúng MỘT vế > 0, vế kia = 0 (double-entry). Tiền = money.Money (decimal),
// KHÔNG float.
type EntryLine struct {
	// AccountCode là số hiệu tài khoản tham chiếu public.chart_of_accounts
	// (vd "131","511","632","156","3331"). Service kiểm allow_posting.
	AccountCode string
	// Debit là số tiền ghi NỢ (>= 0). > 0 thì Credit phải = 0.
	Debit money.Money
	// Credit là số tiền ghi CÓ (>= 0). > 0 thì Debit phải = 0.
	Credit money.Money
}

// validate kiểm một dòng đúng MỘT vế > 0, không vế nào âm, có account_code.
func (l EntryLine) validate() error {
	if l.AccountCode == "" {
		return errors.New("dòng bút toán thiếu account_code")
	}
	if l.Debit.IsNegative() || l.Credit.IsNegative() {
		return fmt.Errorf("dòng %s có vế âm (debit/credit phải >= 0)", l.AccountCode)
	}
	debitPos := l.Debit.IsPositive()
	creditPos := l.Credit.IsPositive()
	// XOR: đúng một vế > 0.
	if debitPos == creditPos {
		return fmt.Errorf("dòng %s phải có ĐÚNG MỘT vế > 0 (debit XOR credit)", l.AccountCode)
	}
	return nil
}

// JournalEntry là một BÚT TOÁN kép: tập các EntryLine cân nhau trong một sổ
// (Book) tại một ngày (EntryDate). Append-only. SourceType/SourceID liên kết
// nghiệp vụ nguồn (vd "order"/order code) để truy vết. Memo là diễn giải.
type JournalEntry struct {
	Book       Book
	EntryDate  time.Time
	SourceType string
	SourceID   string
	Memo       string
	Lines      []EntryLine
}

// TotalDebit trả tổng số tiền ghi NỢ của bút toán (cộng decimal, KHÔNG float).
func (e JournalEntry) TotalDebit() money.Money {
	sum := money.Zero()
	for _, l := range e.Lines {
		sum = sum.Add(l.Debit)
	}
	return sum
}

// TotalCredit trả tổng số tiền ghi CÓ của bút toán.
func (e JournalEntry) TotalCredit() money.Money {
	sum := money.Zero()
	for _, l := range e.Lines {
		sum = sum.Add(l.Credit)
	}
	return sum
}

// Validate ép bất biến kế toán THUẦN của bút toán (không cần DB/kỳ — chỉ cấu
// trúc): book hợp lệ; có ít nhất 1 dòng; mỗi dòng đúng MỘT vế > 0; và quan trọng
// nhất Σdebit = Σcredit. Gating kỳ mở + allow_posting nằm ở tầng app (service),
// vì cần dữ liệu ngoài (kỳ kế toán, cây tài khoản) — domain không biết.
func (e JournalEntry) Validate() error {
	if !e.Book.Valid() {
		return fmt.Errorf("book %q không hợp lệ (chỉ INTERNAL hoặc TAX)", e.Book)
	}
	if len(e.Lines) == 0 {
		return errors.New("bút toán phải có ít nhất một dòng")
	}
	for _, l := range e.Lines {
		if err := l.validate(); err != nil {
			return err
		}
	}
	if !e.TotalDebit().Equal(e.TotalCredit()) {
		return fmt.Errorf("bút toán lệch: Σdebit=%s <> Σcredit=%s", e.TotalDebit(), e.TotalCredit())
	}
	return nil
}
