// Package app là tầng use-case của finance: điều phối GHI phiếu thu (validate
// thuần + idempotency chống cộng tiền 2 lần) trong transaction của caller. Mở/
// commit transaction ở đây (hoặc nhận tx từ orders/POS để gộp atomic với trừ kho
// + post sổ). Domain không thấy tx (arch_test chặn). Mẫu theo accounting (Poster)
// + inventory (DeductPort).
package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/finance/domain"
)

// PaymentWriter là PORT GHI phiếu thu bound tới MỘT transaction (do caller hoặc
// TxManager truyền). Adapter postgres implement bằng appdb.Queries.WithTx(tx). Là
// port ở TẦNG APP (không domain) vì gắn pgx.Tx + điều phối transaction — domain
// THUẦN không được biết pgx (arch_test chặn). Mọi thao tác chạy trong CÙNG tx để
// gộp atomic với nghiệp vụ orders/POS.
//
// ⚠️ KHÔNG có thao tác UPDATE fund_accounts.balance: PROD có trigger tự cộng số dư
// khi phiếu sang status='completed'. Go chỉ INSERT phiếu completed — trigger lo số
// dư, tránh double-count.
type PaymentWriter interface {
	// FindAliveByBankRef tìm phiếu còn sống theo bank_reference_id (dedup webhook).
	// Không thấy → (nil, nil). Lỗi hạ tầng → (_, err).
	FindAliveByBankRef(ctx context.Context, bankRef string) (*domain.Payment, error)
	// FindAliveByCode tìm phiếu còn sống theo code (dedup phiếu thủ công idempotent).
	// Không thấy → (nil, nil). Lỗi hạ tầng → (_, err).
	FindAliveByCode(ctx context.Context, code string) (*domain.Payment, error)
	// InsertPaymentIn ghi MỘT phiếu THU (flow='in', status='completed') với code đã
	// sinh idempotent. Trả phiếu đã ghi (id do DB sinh). Nếu trùng code/bank_ref
	// (unique partial index) → trả (nil, true, nil) để service re-SELECT phiếu cũ
	// (idempotent no-op, KHÔNG cộng tiền lần 2). Lỗi khác → (_, false, err).
	InsertPaymentIn(ctx context.Context, code string, p domain.RecordPaymentIn) (saved *domain.Payment, duplicate bool, err error)
	// InsertPaymentOut ghi MỘT phiếu CHI (flow='out') với code đã sinh idempotent —
	// trả NCC mua hàng (mục 54), đối xứng InsertPaymentIn. Trùng code/bank_ref → trả
	// (nil, true, nil) để service re-SELECT phiếu cũ (idempotent no-op). Trigger prod
	// TRỪ số dư quỹ MỘT LẦN khi status='completed' — Go KHÔNG tự trừ.
	InsertPaymentOut(ctx context.Context, code string, p domain.RecordPaymentOut) (saved *domain.Payment, duplicate bool, err error)
	// ConfirmReceipt chuyển phiếu THU 'pending' (đã thu từ khách) → 'completed' (vào
	// quỹ) cho thanh toán 2 bước (spec mục 55). Trả số dòng đổi: 1 = vừa xác nhận;
	// 0 = phiếu đã 'completed' / không phải pending (idempotent, KHÔNG cộng đôi số dư).
	ConfirmReceipt(ctx context.Context, paymentID int64) (rows int64, err error)
}

// PaymentWriterFromTx dựng một PaymentWriter bound tới tx. TxManager / Recorder
// dùng để lấy writer cho transaction hiện hành. Tách thành func để app không phụ
// thuộc cứng vào cách khởi tạo repo (adapter cung cấp).
type PaymentWriterFromTx func(tx pgx.Tx) PaymentWriter

// TxManager mở/commit một transaction cho trường hợp ghi phiếu ĐỘC LẬP
// (RecordPaymentInOwnTx). Adapter implement bằng platform/db.WithinTx. Khi
// orders/POS ghi phiếu trong tx nghiệp vụ của HỌ, chúng gọi thẳng
// Recorder.RecordPaymentIn(ctx, tx, ...) (KHÔNG qua TxManager) để gộp atomic.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

// PaymentReader là PORT ĐỌC phiếu (bind pool) cho HTTP. Tách khỏi PaymentWriter
// (ghi/tx) để đường đọc không cần transaction.
type PaymentReader interface {
	// ListByRef trả các phiếu còn sống trỏ về (ref_type, ref_id) — sắp mới nhất
	// trước. Hết → slice rỗng.
	ListByRef(ctx context.Context, refType, refID string, limit int32) ([]domain.Payment, error)
}
