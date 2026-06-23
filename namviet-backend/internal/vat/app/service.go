package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

const (
	defaultLimit = 50
	maxLimit     = 200
	dateLayout   = "2006-01-02"
)

// Service là use-case của vat. Nó implement port IssuePort (phát hành HĐ trong
// tx của caller orders) + cung cấp use-case đọc HĐ. Đường PHÁT HÀNH KHÔNG có
// REST POST công khai ở P5 — chỉ orders gọi IssueInvoice(ctx, tx, params) trong
// tx giao hàng của họ để gộp atomic (HĐ + giao hàng + post sổ TAX cùng tx).
type Service struct {
	storeFromTx InvoiceStoreFromTx
	txm         TxManager
	read        InvoiceReader
}

// New dựng Service. storeFromTx bind InvoiceStore tới một tx; txm để
// IssueInvoiceInOwnTx mở tx riêng; read là đường đọc (bind pool). Có thể truyền
// nil cho thành phần không dùng (vd test).
func New(storeFromTx InvoiceStoreFromTx, txm TxManager, read InvoiceReader) *Service {
	return &Service{storeFromTx: storeFromTx, txm: txm, read: read}
}

// IssueParams là input phát hành một HĐ VAT. vat_rate nằm trong từng LineInput
// (input — KHÔNG hardcode). IssueDate rỗng (zero) → service dùng ngày hiện tại.
type IssueParams struct {
	OrderCode       string
	CustomerTaxCode string
	Serial          string
	// MauSo (mẫu số HĐ, optional) — gắn khi tạo dòng serial lần đầu (EnsureSerial).
	MauSo     string
	IssueDate time.Time
	Lines     []domain.LineInput
}

// IssueInvoice PHÁT HÀNH một HĐ VAT cho đơn TRONG transaction tx do CALLER
// (orders) truyền (gộp atomic với giao hàng + post sổ TAX). Quy trình:
//  1. IDEMPOTENCY (1 đơn 1 HĐ): đã có HĐ 'issued' cho order_code → TRẢ HĐ cũ
//     (no-op, không phát hành trùng nếu orders retry).
//  2. Dựng + CÂN HĐ thuần ở domain (BuildInvoice): MST bắt buộc/≥1 dòng/không âm,
//     tính line_amount + line_vat (làm tròn VND) + Σ subtotal/vat/total cân khít.
//     Sai → Validation (422).
//  3. Cấp số GAPLESS: EnsureSerial → NextInvoiceNo(serial) khoá dòng serial
//     (FOR UPDATE) → số liên tục không trùng/không nhảy (tuần tự hoá theo serial).
//  4. INSERT header (status='issued') + lines trong tx truyền vào. Trùng do race
//     (UNIQUE order_code WHERE issued) → đọc lại HĐ cũ (idempotent).
//  5. (DEFER) Phát hành điện tử qua provider — chỉ chừa port, KHÔNG gọi ở P5.
//
// Trả IssuedInvoice (kèm lines). Lỗi đã là apperr (map envelope ở http).
func (s *Service) IssueInvoice(ctx context.Context, tx pgx.Tx, p IssueParams) (domain.IssuedInvoice, error) {
	store := s.storeFromTx(tx)

	// (1) Idempotency: 1 đơn 1 HĐ — đã phát hành thì trả HĐ cũ (kèm lines).
	if existing, err := store.FindIssuedByOrder(ctx, p.OrderCode); err != nil {
		return domain.IssuedInvoice{}, apperr.Internal("tra cứu HĐ theo đơn lỗi").WithCause(err)
	} else if existing != nil {
		return store.GetInvoiceWithLines(ctx, existing.ID)
	}

	// (2) Dựng + cân HĐ thuần ở domain. IssueDate rỗng → hôm nay.
	issueDate := p.IssueDate
	if issueDate.IsZero() {
		issueDate = time.Now()
	}
	inv, err := domain.BuildInvoice(p.OrderCode, p.CustomerTaxCode, p.Serial, issueDate, p.Lines)
	if err != nil {
		return domain.IssuedInvoice{}, apperr.Validation(err.Error())
	}

	// (3) Cấp số GAPLESS theo serial (EnsureSerial idempotent → khoá + bump).
	invoiceNo, err := store.NextInvoiceNo(ctx, inv.Serial, p.MauSo)
	if err != nil {
		return domain.IssuedInvoice{}, apperr.Internal("cấp số hoá đơn lỗi").WithCause(err)
	}

	// (4) INSERT header + lines. Trùng (race trên UNIQUE order_code WHERE issued)
	// → đọc lại HĐ cũ trả về (idempotent no-op).
	invoiceID, duplicate, err := store.InsertInvoice(ctx, inv, invoiceNo)
	if err != nil {
		return domain.IssuedInvoice{}, apperr.Internal("ghi hoá đơn lỗi").WithCause(err)
	}
	if duplicate {
		existing, err := store.FindIssuedByOrder(ctx, p.OrderCode)
		if err != nil {
			return domain.IssuedInvoice{}, apperr.Internal("đọc lại HĐ trùng lỗi").WithCause(err)
		}
		if existing == nil {
			return domain.IssuedInvoice{}, apperr.Internal("HĐ trùng nhưng không tìm thấy bản ghi cũ")
		}
		return store.GetInvoiceWithLines(ctx, existing.ID)
	}
	for _, l := range inv.Lines {
		if err := store.InsertLine(ctx, invoiceID, l); err != nil {
			return domain.IssuedInvoice{}, apperr.Internal("ghi dòng hoá đơn lỗi").WithCause(err)
		}
	}
	return store.GetInvoiceWithLines(ctx, invoiceID)
}

// IssueInvoiceInOwnTx phát hành HĐ trong một transaction RIÊNG (mở/commit ở đây).
// Dùng khi phát hành ĐỘC LẬP, không gộp với nghiệp vụ khác (vd test, phát hành
// rời). orders dùng IssueInvoice(ctx, tx, ...) trong tx của họ thay vì hàm này.
func (s *Service) IssueInvoiceInOwnTx(ctx context.Context, p IssueParams) (domain.IssuedInvoice, error) {
	if s.txm == nil {
		return domain.IssuedInvoice{}, apperr.Internal("TxManager chưa cấu hình cho IssueInvoiceInOwnTx")
	}
	var out domain.IssuedInvoice
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		got, err := s.IssueInvoice(ctx, tx, p)
		if err != nil {
			return err
		}
		out = got
		return nil
	})
	if err != nil {
		return domain.IssuedInvoice{}, err
	}
	return out, nil
}

// ListInvoicesQuery là input đọc danh sách HĐ đã giải mã ở edge.
type ListInvoicesQuery struct {
	Cursor    string
	Limit     int32
	OrderCode string
	Status    string
}

// ListInvoicesResult là một trang HĐ + cursor trang kế (rỗng nếu hết).
type ListInvoicesResult struct {
	Items      []domain.IssuedInvoice
	NextCursor string
}

// ListInvoices trả một trang HĐ (keyset created_at DESC, id DESC). Tự decode
// cursor, chuẩn hoá limit, sinh NextCursor nếu trang đầy.
func (s *Service) ListInvoices(ctx context.Context, q ListInvoicesQuery) (ListInvoicesResult, error) {
	afterNano, afterID, err := decodeCursor(q.Cursor)
	if err != nil {
		return ListInvoicesResult{}, apperr.Validation("cursor không hợp lệ")
	}
	limit := normalizeLimit(q.Limit)
	f := domain.InvoiceFilter{Limit: limit, OrderCode: q.OrderCode, Status: q.Status}
	if q.Cursor != "" {
		f.AfterCreatedAt = time.Unix(0, afterNano).UTC()
		f.AfterID = afterID
		f.HasCursor = true
	}
	items, err := s.read.ListInvoices(ctx, f)
	if err != nil {
		return ListInvoicesResult{}, err
	}
	res := ListInvoicesResult{Items: items}
	if int32(len(items)) == limit && limit > 0 {
		last := items[len(items)-1]
		res.NextCursor = encodeCursor(last.CreatedAt.UnixNano(), last.ID)
	}
	return res, nil
}

// GetInvoice trả một HĐ theo id kèm các dòng.
func (s *Service) GetInvoice(ctx context.Context, id string) (domain.IssuedInvoice, error) {
	return s.read.GetInvoice(ctx, id)
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

// Đảm bảo Service thoả IssuePort ở compile-time (port nội bộ cho orders).
var _ IssuePort = (*Service)(nil)

// IssuePort là PORT NỘI BỘ: module orders phát hành HĐ VAT trong tx giao hàng
// của họ qua đây (gộp atomic). KHÔNG có REST POST công khai ở P5. Định nghĩa ở
// app (không domain) vì nhận pgx.Tx — domain THUẦN không biết tx.
type IssuePort interface {
	IssueInvoice(ctx context.Context, tx pgx.Tx, p IssueParams) (domain.IssuedInvoice, error)
}
