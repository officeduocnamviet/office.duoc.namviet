// Package posting là lớp TEMPLATE BÚT TOÁN (posting templates) của bounded
// context accounting: hàm THUẦN map một SỰ KIỆN nghiệp vụ (bán hàng, giá vốn,
// thu tiền) → []domain.JournalEntry ĐÃ CÂN, mỗi entry đúng MỘT sổ (Book). P4
// (orders) gọi các builder này rồi đưa từng entry qua accounting Poster.Post
// trong tx nghiệp vụ của họ (gộp atomic — sổ luôn khớp sự kiện).
//
// Package này THUẦN như domain: chỉ stdlib (time) + domain + common/money. KHÔNG
// pgx/http/tx — không tự post, chỉ DỰNG bút toán. Mọi mã tài khoản + cờ chính
// sách per-book gom Ở ĐÂY (Rules) — đổi 1 chỗ, P4 KHÔNG phải sửa.
package posting

import "github.com/Maneva-AI/namviet-backend/internal/accounting/domain"

// Accounts là CẤU HÌNH mã tài khoản TT133 dùng để dựng bút toán. Gom một chỗ để
// không hardcode rải rác trong P4 (ARCHITECTURE.md §7 / tt133-posting-rules §4).
//
// MẶC ĐỊNH TT133 — kế toán Nam Việt PHẢI duyệt trước cutover (xem
// docs/superpowers/specs/2026-06-22-tt133-posting-rules.md §5). Khi mã tài khoản
// thật trên public.chart_of_accounts khác mặc định, đổi DefaultRules Ở ĐÂY —
// KHÔNG sửa code orchestration (P4).
type Accounts struct {
	AR        string // Phải thu của khách hàng (bán chịu B2B)
	Cash      string // Tiền mặt (thu ngay tại quầy)
	Bank      string // Tiền gửi ngân hàng (chuyển khoản)
	Revenue   string // Doanh thu bán hàng & cung cấp dịch vụ (chưa VAT)
	VATOutput string // Thuế GTGT đầu ra phải nộp
	COGS      string // Giá vốn hàng bán
	Inventory string // Giá mua hàng hóa (tồn kho — tăng khi nhập, giảm khi xuất bán)
	// ---- Chiều MUA (purchasing, mục 54) — MẶC ĐỊNH TT133, kế toán duyệt ----
	InputVAT string // Thuế GTGT được khấu trừ (133) — VAT đầu vào khi mua có HĐ
	Payable  string // Phải trả người bán (331) — công nợ NCC khi nhập kho
}

// Rules là TOÀN BỘ cấu hình template bút toán: mã tài khoản + cờ chính sách
// per-book. P4 dựng một Rules (mặc định DefaultRules) rồi truyền vào các builder.
//
// Cờ phản ánh quyết định dual-ledger của Nam Việt (tt133-posting-rules §6):
//   - InternalRecordsVAT: sổ INTERNAL (thực tế) CÓ tách 3331 hay không.
//   - TaxRecordsCOGS: sổ TAX (thuế) CÓ ghi giá vốn 632/1561 hay không.
//
// Sổ TAX LUÔN ghi VAT (đó là mục đích sổ thuế) → không cần cờ cho việc đó.
type Rules struct {
	Accounts Accounts
	// InternalRecordsVAT: true = sổ INTERNAL tách dòng 3331 (ghi đủ kinh tế
	// thực, mặc định). false = INTERNAL chỉ ghi tiền↔doanh thu, VAT chỉ ở sổ TAX.
	InternalRecordsVAT bool
	// TaxRecordsCOGS: true = sổ TAX cũng ghi bút toán giá vốn 632/1561. false =
	// sổ TAX CHỈ ghi doanh thu+VAT phục vụ báo cáo GTGT (mặc định).
	TaxRecordsCOGS bool
}

// DefaultRules là MẶC ĐỊNH KỸ THUẬT để build P1+P4 chạy end-to-end
// (tt133-posting-rules §6, 2026-06-22). KHÔNG phải quyết định kế toán cuối cùng —
// kế toán Nam Việt duyệt §5 trước cutover; nếu khác thì sửa Ở ĐÂY, P4 không đổi.
//
//	Mã TK : AR=131, Cash=111, Bank=112, Revenue=511, VATOutput=3331,
//	        COGS=632, Inventory=1561, InputVAT=133, Payable=331
//	        (⚠️ verify với public.chart_of_accounts thật)
//	Cờ    : InternalRecordsVAT=true (INTERNAL có 3331/133),
//	        TaxRecordsCOGS=false (TAX không ghi giá vốn)
var DefaultRules = Rules{
	Accounts: Accounts{
		AR:        "131",
		Cash:      "111",
		Bank:      "112",
		Revenue:   "511",
		VATOutput: "3331",
		COGS:      "632",
		Inventory: "1561",
		InputVAT:  "133", // MẶC ĐỊNH TT133 (Thuế GTGT được khấu trừ), kế toán duyệt
		Payable:   "331", // MẶC ĐỊNH TT133 (Phải trả người bán), kế toán duyệt
	},
	InternalRecordsVAT: true,
	TaxRecordsCOGS:     false,
}

// recordsVAT trả true nếu sổ b ghi VAT theo cấu hình hiện tại: sổ TAX LUÔN ghi
// VAT; sổ INTERNAL theo cờ InternalRecordsVAT.
func (r Rules) recordsVAT(b domain.Book) bool {
	if b == domain.BookTax {
		return true
	}
	return r.InternalRecordsVAT
}

// recordsCOGS trả true nếu sổ b ghi bút toán giá vốn: sổ INTERNAL LUÔN ghi (đó là
// nơi theo dõi P&L thật); sổ TAX theo cờ TaxRecordsCOGS.
func (r Rules) recordsCOGS(b domain.Book) bool {
	if b == domain.BookInternal {
		return true
	}
	return r.TaxRecordsCOGS
}

// fundAccount trả mã tài khoản tiền theo loại quỹ: ngân hàng (112) hay tiền mặt
// (111). Dùng cho vế tiền của bán thu-ngay và của PAYMENT_IN.
func (r Rules) fundAccount(fundIsBank bool) string {
	if fundIsBank {
		return r.Accounts.Bank
	}
	return r.Accounts.Cash
}
