// Package app là tầng use-case của vat: điều phối PHÁT HÀNH hoá đơn VAT (cấp số
// gapless theo serial + dựng/cân HĐ + persist) và đọc HĐ. Mở/commit transaction
// ở đây (hoặc nhận tx từ caller orders để gộp atomic với giao hàng). Domain
// không thấy tx.
package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

// InvoiceStore là PORT GHI HĐ bound tới MỘT transaction (do caller orders hoặc
// TxManager truyền vào). Adapter postgres implement bằng appdb.Queries.WithTx(tx).
// Đây là port ở TẦNG APP (không domain) vì gắn với điều phối transaction +
// cấp số gapless (FOR UPDATE) — domain THUẦN không được biết pgx/tx (arch_test
// chặn).
type InvoiceStore interface {
	// FindIssuedByOrder tìm HĐ đã 'issued' cho order_code (idempotency 1 đơn 1 HĐ).
	// Có → (*IssuedInvoice, nil); không có → (nil, nil). Lỗi hạ tầng → (_, err).
	// KHÔNG nạp lines (service nạp riêng khi cần trả về đầy đủ).
	FindIssuedByOrder(ctx context.Context, orderCode string) (*domain.IssuedInvoice, error)
	// NextInvoiceNo cấp số HĐ GAPLESS cho serial: bảo đảm dòng serial tồn tại
	// (EnsureSerial idempotent với mauSo), khoá dòng serial (FOR UPDATE), trả số
	// kế tiếp rồi tăng next_no — TẤT CẢ trong tx hiện hành (atomic). Tuần tự hoá
	// các tx cùng serial nên số LIÊN TỤC, không trùng/không nhảy.
	NextInvoiceNo(ctx context.Context, serial, mauSo string) (int64, error)
	// InsertInvoice ghi header HĐ (status 'issued'), trả id (uuid sinh app-side).
	// invoiceNo đã cấp gapless. Trùng (race trên UNIQUE) → duplicate=true để service
	// đọc lại HĐ cũ. Tiền lấy từ inv (đã cân ở domain).
	InsertInvoice(ctx context.Context, inv domain.Invoice, invoiceNo int64) (id string, duplicate bool, err error)
	// InsertLine ghi một dòng HĐ (line_no theo thứ tự, đã tính ở domain).
	InsertLine(ctx context.Context, invoiceID string, l domain.InvoiceLine) error
	// GetInvoiceWithLines nạp HĐ theo id KÈM lines (để trả về sau khi phát hành).
	GetInvoiceWithLines(ctx context.Context, id string) (domain.IssuedInvoice, error)
}

// InvoiceStoreFromTx dựng một InvoiceStore bound tới tx. TxManager / IssuePort
// dùng để lấy store cho transaction hiện hành.
type InvoiceStoreFromTx func(tx pgx.Tx) InvoiceStore

// TxManager mở/commit một transaction cho trường hợp phát hành ĐỘC LẬP
// (IssueInvoiceInOwnTx). Khi orders phát hành trong tx giao hàng của HỌ, chúng
// gọi thẳng IssuePort.IssueInvoice(ctx, tx, ...) (KHÔNG qua TxManager) để gộp
// atomic.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

// InvoiceReader là PORT ĐỌC HĐ (bind pool) cho HTTP. Tách khỏi InvoiceStore
// (ghi/tx) để đường đọc không cần transaction.
type InvoiceReader interface {
	// ListInvoices trả một trang HĐ theo keyset (created_at DESC, id DESC), lọc
	// optional order_code/status. Hết → slice rỗng. KHÔNG nạp lines (chỉ header).
	ListInvoices(ctx context.Context, f domain.InvoiceFilter) ([]domain.IssuedInvoice, error)
	// GetInvoice trả một HĐ theo id KÈM lines. Không thấy → apperr.NotFound.
	GetInvoice(ctx context.Context, id string) (domain.IssuedInvoice, error)
}
