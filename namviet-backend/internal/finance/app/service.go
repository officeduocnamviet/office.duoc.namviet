package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
)

const (
	defaultLimit = 50
	maxLimit     = 200
	// codePrefix gắn đầu mã phiếu thu sinh tự động (phân biệt với mã ERP nhập tay).
	codePrefix = "PT"
	// codePrefixOut gắn đầu mã phiếu CHI (trả NCC mua hàng — mục 54).
	codePrefixOut = "PC"
)

// Recorder là use-case GHI của finance: ghi phiếu THU idempotent vào
// finance_transactions TRONG tx của caller. Implement port nội bộ RecordPort
// (orders/POS gọi RecordPaymentIn trong tx của họ để gộp atomic với trừ kho + post
// sổ). KHÔNG có REST POST công khai — giống accounting.Poster / inventory.DeductPort.
type Recorder struct {
	writerFromTx PaymentWriterFromTx
	txm          TxManager
	read         PaymentReader
}

// New dựng Recorder. writerFromTx bind PaymentWriter tới một tx; txm để
// RecordPaymentInOwnTx mở tx riêng; read là đường đọc (bind pool). Có thể truyền
// nil cho thành phần không dùng (vd test).
func New(writerFromTx PaymentWriterFromTx, txm TxManager, read PaymentReader) *Recorder {
	return &Recorder{writerFromTx: writerFromTx, txm: txm, read: read}
}

// RecordPaymentInParams là input ghi phiếu THU. IdemKey BẮT BUỘC cho phiếu THỦ
// CÔNG (BankRef == nil) — dùng sinh code idempotent chống cộng tiền 2 lần. Phiếu
// TỰ ĐỘNG (BankRef != nil) dedup theo bank_reference_id (IdemKey có thể rỗng).
type RecordPaymentInParams struct {
	domain.RecordPaymentIn
	// IdemKey khoá idempotency tầng app (phiếu thủ công). Trùng key → trả phiếu cũ.
	IdemKey string
}

// RecordPaymentIn ghi MỘT phiếu THU cho đơn TRONG transaction tx do CALLER truyền
// (gộp atomic với nghiệp vụ orders/POS). Quy trình:
//  1. Validate THUẦN (amount>0, book_type hợp lệ, order_code không rỗng) → 422.
//  2. IDEMPOTENCY (chống cộng tiền 2 lần):
//     - BankRef != nil (webhook): tìm phiếu theo bank_reference_id → có thì TRẢ phiếu
//     cũ (no-op). code sinh từ bank_ref.
//     - BankRef == nil (thủ công): IdemKey bắt buộc → code sinh từ IdemKey; tìm theo
//     code → có thì TRẢ phiếu cũ.
//  3. INSERT phiếu (flow='in', status='completed'). Trùng (race, unique index) →
//     re-SELECT phiếu cũ trả về. Trigger prod tự cộng fund_accounts.balance MỘT
//     LẦN — Go KHÔNG tự cộng (tránh double-count).
//
// Trả (Payment, created, error). created=true CHỈ khi VỪA INSERT phiếu MỚI;
// false khi idempotent hit (dedup tìm thấy phiếu cũ HOẶC race-duplicate). Caller
// orchestration PHẢI dùng created để CHỈ post bút toán / phân bổ / trừ kho MỘT LẦN
// — replay cùng Idempotency-Key (created=false) KHÔNG được lặp các side-effect đó
// (nếu không sẽ nhân đôi sổ kế toán + "đã thu"). Lỗi đã là apperr.
func (r *Recorder) RecordPaymentIn(ctx context.Context, tx pgx.Tx, p RecordPaymentInParams) (domain.Payment, bool, error) {
	if err := p.Validate(); err != nil {
		return domain.Payment{}, false, apperr.Validation(err.Error())
	}

	hasBankRef := p.BankRef != nil && strings.TrimSpace(*p.BankRef) != ""
	if !hasBankRef && strings.TrimSpace(p.IdemKey) == "" {
		// Phiếu thủ công không có khoá chống trùng nào → từ chối (tránh cộng tiền lặp).
		return domain.Payment{}, false, apperr.Validation("phiếu thu thủ công cần Idempotency-Key để chống trùng")
	}

	w := r.writerFromTx(tx)

	// Chuẩn hoá: code idempotent (từ bank_ref nếu có, ngược lại từ idem key).
	var dedupKey string
	if hasBankRef {
		dedupKey = "bank:" + strings.TrimSpace(*p.BankRef)
	} else {
		dedupKey = "idem:" + strings.TrimSpace(p.IdemKey)
	}
	code := makeCode(dedupKey)

	// (2) Dedup TRƯỚC khi insert: phiếu webhook theo bank_ref; thủ công theo code.
	if hasBankRef {
		existing, err := w.FindAliveByBankRef(ctx, strings.TrimSpace(*p.BankRef))
		if err != nil {
			return domain.Payment{}, false, apperr.Internal("tra cứu phiếu theo bank_ref lỗi").WithCause(err)
		}
		if existing != nil {
			return *existing, false, nil // idempotent hit — KHÔNG tạo mới
		}
	} else {
		existing, err := w.FindAliveByCode(ctx, code)
		if err != nil {
			return domain.Payment{}, false, apperr.Internal("tra cứu phiếu theo code lỗi").WithCause(err)
		}
		if existing != nil {
			return *existing, false, nil // idempotent hit — KHÔNG tạo mới
		}
	}

	// (3) INSERT (flow='in', status='completed'). Trùng do race → re-SELECT phiếu cũ.
	// Chuẩn hoá OrderCode (trim) khi ghi để ref_id KHỚP đúng lookup đọc lại sau này.
	rec := p.RecordPaymentIn
	rec.OrderCode = strings.TrimSpace(rec.OrderCode)
	saved, duplicate, err := w.InsertPaymentIn(ctx, code, rec)
	if err != nil {
		return domain.Payment{}, false, apperr.Internal("ghi phiếu thu lỗi").WithCause(err)
	}
	if duplicate {
		// Hai luồng cùng key chạy đua: phiếu kia thắng. Đọc lại trả phiếu cũ (no-op).
		var existing *domain.Payment
		if hasBankRef {
			existing, err = w.FindAliveByBankRef(ctx, strings.TrimSpace(*p.BankRef))
		} else {
			existing, err = w.FindAliveByCode(ctx, code)
		}
		if err != nil {
			return domain.Payment{}, false, apperr.Internal("đọc lại phiếu trùng lỗi").WithCause(err)
		}
		if existing == nil {
			return domain.Payment{}, false, apperr.Internal("phiếu trùng nhưng không tìm thấy bản ghi cũ")
		}
		return *existing, false, nil // race-duplicate cũng là idempotent hit
	}
	return *saved, true, nil // VỪA tạo mới
}

// RecordPaymentOutParams là input ghi phiếu CHI (trả NCC mua hàng — mục 54). IdemKey
// BẮT BUỘC cho phiếu THỦ CÔNG (BankRef == nil). Đối xứng RecordPaymentInParams.
type RecordPaymentOutParams struct {
	domain.RecordPaymentOut
	IdemKey string
}

// RecordPaymentOut ghi MỘT phiếu CHI (trả NCC) TRONG transaction tx do CALLER truyền
// (gộp atomic với nghiệp vụ purchasing). Đối xứng RecordPaymentIn:
//  1. Validate THUẦN (amount>0, book_type hợp lệ, ref_id không rỗng) → 422.
//  2. IDEMPOTENCY (chống chi tiền 2 lần): BankRef != nil → dedup theo bank_reference_id;
//     BankRef == nil → IdemKey bắt buộc, code sinh từ IdemKey, dedup theo code.
//  3. INSERT phiếu (flow='out'). Trùng (race) → re-SELECT phiếu cũ. Trigger prod tự
//     TRỪ fund_accounts.balance MỘT LẦN — Go KHÔNG tự trừ.
//
// Trả (Payment, created, error). created=true CHỈ khi VỪA INSERT phiếu MỚI; false khi
// idempotent hit (dedup tìm thấy phiếu cũ HOẶC race-duplicate). Caller orchestration
// PHẢI dùng created để CHỈ post bút toán chi MỘT LẦN — replay cùng Idempotency-Key
// (created=false) KHÔNG được post lại (nếu không sẽ nhân đôi sổ + "đã chi"). Mã phiếu
// chi sinh DETERMINISTIC từ idem key với tiền tố RIÊNG (PC) — KHÔNG va chạm phiếu thu.
func (r *Recorder) RecordPaymentOut(ctx context.Context, tx pgx.Tx, p RecordPaymentOutParams) (domain.Payment, bool, error) {
	if err := p.Validate(); err != nil {
		return domain.Payment{}, false, apperr.Validation(err.Error())
	}

	hasBankRef := p.BankRef != nil && strings.TrimSpace(*p.BankRef) != ""
	if !hasBankRef && strings.TrimSpace(p.IdemKey) == "" {
		return domain.Payment{}, false, apperr.Validation("phiếu chi thủ công cần Idempotency-Key để chống trùng")
	}

	w := r.writerFromTx(tx)

	// Code idempotent: bank_ref nếu có, ngược lại idem key. Tiền tố PC (phiếu chi) gắn
	// vào dedupKey để KHÔNG đụng code phiếu thu (PT) cùng idem key.
	var dedupKey string
	if hasBankRef {
		dedupKey = "bank:" + strings.TrimSpace(*p.BankRef)
	} else {
		dedupKey = "idem:" + strings.TrimSpace(p.IdemKey)
	}
	code := makeOutCode(dedupKey)

	// (2) Dedup TRƯỚC khi insert.
	if hasBankRef {
		existing, err := w.FindAliveByBankRef(ctx, strings.TrimSpace(*p.BankRef))
		if err != nil {
			return domain.Payment{}, false, apperr.Internal("tra cứu phiếu chi theo bank_ref lỗi").WithCause(err)
		}
		if existing != nil {
			return *existing, false, nil
		}
	} else {
		existing, err := w.FindAliveByCode(ctx, code)
		if err != nil {
			return domain.Payment{}, false, apperr.Internal("tra cứu phiếu chi theo code lỗi").WithCause(err)
		}
		if existing != nil {
			return *existing, false, nil
		}
	}

	// (3) INSERT (flow='out'). Trùng do race → re-SELECT phiếu cũ.
	rec := p.RecordPaymentOut
	rec.POCode = strings.TrimSpace(rec.POCode)
	saved, duplicate, err := w.InsertPaymentOut(ctx, code, rec)
	if err != nil {
		return domain.Payment{}, false, apperr.Internal("ghi phiếu chi lỗi").WithCause(err)
	}
	if duplicate {
		var existing *domain.Payment
		if hasBankRef {
			existing, err = w.FindAliveByBankRef(ctx, strings.TrimSpace(*p.BankRef))
		} else {
			existing, err = w.FindAliveByCode(ctx, code)
		}
		if err != nil {
			return domain.Payment{}, false, apperr.Internal("đọc lại phiếu chi trùng lỗi").WithCause(err)
		}
		if existing == nil {
			return domain.Payment{}, false, apperr.Internal("phiếu chi trùng nhưng không tìm thấy bản ghi cũ")
		}
		return *existing, false, nil
	}
	return *saved, true, nil
}

// RecordPaymentOutOwnTx ghi phiếu CHI trong một transaction RIÊNG (mở/commit ở đây).
func (r *Recorder) RecordPaymentOutOwnTx(ctx context.Context, p RecordPaymentOutParams) (domain.Payment, error) {
	if r.txm == nil {
		return domain.Payment{}, apperr.Internal("TxManager chưa cấu hình cho RecordPaymentOutOwnTx")
	}
	var out domain.Payment
	err := r.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		got, _, err := r.RecordPaymentOut(ctx, tx, p)
		if err != nil {
			return err
		}
		out = got
		return nil
	})
	if err != nil {
		return domain.Payment{}, err
	}
	return out, nil
}

// RecordPaymentInOwnTx ghi phiếu THU trong một transaction RIÊNG (mở/commit ở
// đây). Dùng khi ghi phiếu ĐỘC LẬP, không gộp với nghiệp vụ khác (vd thu nợ rời,
// đối soát webhook đứng riêng, test).
func (r *Recorder) RecordPaymentInOwnTx(ctx context.Context, p RecordPaymentInParams) (domain.Payment, error) {
	if r.txm == nil {
		return domain.Payment{}, apperr.Internal("TxManager chưa cấu hình cho RecordPaymentInOwnTx")
	}
	var out domain.Payment
	err := r.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		got, _, err := r.RecordPaymentIn(ctx, tx, p)
		if err != nil {
			return err
		}
		out = got
		return nil
	})
	if err != nil {
		return domain.Payment{}, err
	}
	return out, nil
}

// ConfirmReceipt (thủ quỹ "Xác nhận đã thu" — thanh toán 2 bước, spec mục 55)
// chuyển phiếu THU 'pending' → 'completed' TRONG tx của caller (gộp atomic nếu cần).
// Idempotent: phiếu đã 'completed' → confirmed=false, KHÔNG lỗi (không cộng đôi số
// dư). paymentID không tồn tại / không phải pending-in → confirmed=false.
func (r *Recorder) ConfirmReceipt(ctx context.Context, tx pgx.Tx, paymentID int64) (confirmed bool, err error) {
	w := r.writerFromTx(tx)
	rows, cerr := w.ConfirmReceipt(ctx, paymentID)
	if cerr != nil {
		return false, apperr.Internal("xác nhận thu tiền lỗi").WithCause(cerr)
	}
	return rows > 0, nil
}

// ConfirmReceiptInOwnTx xác nhận phiếu THU trong tx RIÊNG (thủ quỹ thao tác độc lập).
func (r *Recorder) ConfirmReceiptInOwnTx(ctx context.Context, paymentID int64) (bool, error) {
	if r.txm == nil {
		return false, apperr.Internal("TxManager chưa cấu hình cho ConfirmReceiptInOwnTx")
	}
	var confirmed bool
	err := r.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		c, e := r.ConfirmReceipt(ctx, tx, paymentID)
		if e != nil {
			return e
		}
		confirmed = c
		return nil
	})
	if err != nil {
		return false, err
	}
	return confirmed, nil
}

// ListOrderPayments trả các phiếu thu/chi của MỘT đơn (ref_type='order',
// ref_id=orderCode) — cho HTTP đọc. Sắp mới nhất trước.
func (r *Recorder) ListOrderPayments(ctx context.Context, orderCode string, limit int32) ([]domain.Payment, error) {
	if strings.TrimSpace(orderCode) == "" {
		return nil, apperr.Validation("ref_id (mã đơn) không được rỗng")
	}
	return r.read.ListByRef(ctx, domain.RefTypeOrder, strings.TrimSpace(orderCode), normalizeLimit(limit))
}

// makeCode sinh mã phiếu thu DETERMINISTIC từ khoá dedup (bank_ref hoặc idem key):
// băm sha256 rồi lấy 24 hex đầu, gắn tiền tố PT. Cùng khoá → cùng code → unique
// index chặn trùng (idempotent). KHÔNG dùng thời gian/ngẫu nhiên (phải lặp lại
// được). Ngắn gọn, đủ phân biệt; KHÔNG va chạm thực tế.
func makeCode(dedupKey string) string {
	sum := sha256.Sum256([]byte(dedupKey))
	return codePrefix + "-" + hex.EncodeToString(sum[:])[:24]
}

// makeOutCode sinh mã phiếu CHI DETERMINISTIC từ khoá dedup, tiền tố PC (phiếu chi)
// — KHÁC PT (phiếu thu) để cùng idem key KHÔNG va chạm code giữa thu và chi.
func makeOutCode(dedupKey string) string {
	sum := sha256.Sum256([]byte("out:" + dedupKey))
	return codePrefixOut + "-" + hex.EncodeToString(sum[:])[:24]
}

func normalizeLimit(l int32) int32 {
	switch {
	case l <= 0:
		return defaultLimit
	case l > maxLimit:
		return maxLimit
	default:
		return l
	}
}

// Đảm bảo Recorder thoả RecordPort ở compile-time (port nội bộ cho orders/POS).
var _ RecordPort = (*Recorder)(nil)

// RecordPort là PORT NỘI BỘ: module orders/POS ghi phiếu THU trong tx nghiệp vụ
// của họ qua đây (gộp atomic với trừ kho + post sổ). KHÔNG có REST POST công khai.
// Định nghĩa ở app (không domain) vì nhận pgx.Tx — domain THUẦN không biết tx.
type RecordPort interface {
	RecordPaymentIn(ctx context.Context, tx pgx.Tx, p RecordPaymentInParams) (domain.Payment, bool, error)
}

// Đảm bảo Recorder thoả RecordOutPort ở compile-time (port nội bộ cho purchasing).
var _ RecordOutPort = (*Recorder)(nil)

// RecordOutPort là PORT NỘI BỘ: module purchasing ghi phiếu CHI (trả NCC) trong tx
// nghiệp vụ của họ qua đây (gộp atomic với post sổ). Đối xứng RecordPort. KHÔNG có
// REST POST công khai.
type RecordOutPort interface {
	RecordPaymentOut(ctx context.Context, tx pgx.Tx, p RecordPaymentOutParams) (domain.Payment, bool, error)
}
