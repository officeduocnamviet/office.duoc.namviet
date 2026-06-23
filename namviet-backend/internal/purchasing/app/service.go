// service.go — use-case GHI của purchasing (mục 54). 4 use-case, mỗi cái = 1
// transaction ATOMIC (platform/db.WithinTx qua TxManager). Cross-module
// (StockIn / Poster.Post / RecordPaymentOut) nhận CÙNG pgx.Tx → lỗi bất kỳ bước nào
// ROLLBACK CẢ CỤM (không nhập kho dở, không bút toán mồ côi). purchasing/app điều
// phối; KHÔNG sửa logic module khác (chỉ GỌI port). Tiền decimal toàn tuyến (CẤM float).
//
// GUARD REPLAY (bài học review chiều bán — BẮT BUỘC tránh): replay cùng
// Idempotency-Key sau commit KHÔNG được nhập kho lại / post bút toán lại / chi tiền
// lại. Dùng (1) state machine guard (UpdateStatus theo status cũ — đã đổi → 0 dòng →
// Conflict/no-op) và (2) cờ `created` từ finance (PaySupplier CHỈ post khi created).
package app

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	accountingdomain "github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/posting"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	financeapp "github.com/Maneva-AI/namviet-backend/internal/finance/app"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing/domain"
)

// Service là use-case GHI GỘP cross-module của purchasing. Giữ store PO bound-tx +
// port inventory/accounting/finance + Rules định khoản (posting.DefaultRules). nil
// port → use-case tương ứng trả Internal (an toàn test từng phần).
type Service struct {
	storeFromTx StoreFromTx
	listFromTx  ListReaderFromTx
	txm         TxManager
	stockIn     StockInner
	poster      Poster
	payer       SupplierPayer
	rules       posting.Rules
	now         func() time.Time
}

// New dựng Service. rules thường là posting.DefaultRules. now cho test bơm thời
// gian; nil → time.Now.
func New(storeFromTx StoreFromTx, listFromTx ListReaderFromTx, txm TxManager, stockIn StockInner, poster Poster, payer SupplierPayer, rules posting.Rules) *Service {
	return &Service{
		storeFromTx: storeFromTx,
		listFromTx:  listFromTx,
		txm:         txm,
		stockIn:     stockIn,
		poster:      poster,
		payer:       payer,
		rules:       rules,
		now:         time.Now,
	}
}

// ---- CreatePO (tạo PO draft) ----

// CreatePOInput là input tạo PO đã giải mã ở edge (tiền/lượng decimal — KHÔNG float).
type CreatePOInput struct {
	SupplierID   *int64
	SupplierName string
	Note         string
	Lines        []domain.DraftLine
	IdemKey      string
}

// CreatePO tạo MỘT PO (status draft) + các dòng hàng trong 1 transaction. Idempotent
// theo Idempotency-Key (1 key → 1 PO; trả PO cũ nếu replay). Trả CreatedPO + created.
func (s *Service) CreatePO(ctx context.Context, in CreatePOInput) (CreatedPO, bool, error) {
	if s.txm == nil || s.storeFromTx == nil {
		return CreatedPO{}, false, apperr.Internal("Service chưa cấu hình đủ cho CreatePO")
	}
	// Dựng draft THUẦN (validate + tính tiền). Lỗi domain → Validation/422.
	draft, derr := domain.NewDraft(domain.DraftInput{
		SupplierID:   in.SupplierID,
		SupplierName: strings.TrimSpace(in.SupplierName),
		Note:         in.Note,
		Lines:        in.Lines,
	})
	if derr != nil {
		return CreatedPO{}, false, apperr.Validation(derr.Error())
	}

	idemKey := strings.TrimSpace(in.IdemKey)
	var out CreatedPO
	var created bool
	txErr := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)
		c, wasCreated, cerr := s.createPO(ctx, store, draft, idemKey)
		if cerr != nil {
			return cerr
		}
		out = c
		created = wasCreated
		return nil
	})
	if txErr != nil {
		if txErr == errIdemRace {
			// Đua key → đọc lại PO của luồng thắng (tx mới).
			c, rerr := s.createIdempotentReadback(ctx, idemKey)
			return c, false, rerr
		}
		return CreatedPO{}, false, txErr
	}
	return out, created, nil
}

// errIdemRace báo hiệu đua idempotency-key (luồng khác thắng) để rollback + readback.
var errIdemRace = apperr.Conflict("đua idempotency key — đọc lại PO của luồng thắng")

// createPO tạo PO trong tx hiện hành, idempotent theo idemKey. Trả (CreatedPO,
// created, error). created=false khi idem hit (PO đã tồn tại).
func (s *Service) createPO(ctx context.Context, store Store, draft domain.Draft, idemKey string) (CreatedPO, bool, error) {
	if idemKey != "" {
		if existingID, found, ferr := store.FindByIdemKey(ctx, idemKey); ferr != nil {
			return CreatedPO{}, false, apperr.Internal("tra cứu idempotency tạo PO lỗi").WithCause(ferr)
		} else if found {
			got, gerr := store.GetCreated(ctx, existingID)
			return got, false, gerr // idem hit → KHÔNG tạo mới
		}
	}
	created, cerr := s.insertNewPO(ctx, store, draft)
	if cerr != nil {
		return CreatedPO{}, false, cerr
	}
	if idemKey != "" {
		inserted, berr := store.BindIdemKey(ctx, idemKey, created.PO.ID, created.PO.Code)
		if berr != nil {
			return CreatedPO{}, false, apperr.Internal("ghi idempotency tạo PO lỗi").WithCause(berr)
		}
		if !inserted {
			return CreatedPO{}, false, errIdemRace
		}
	}
	return created, true, nil
}

// insertNewPO sinh code + id, INSERT header + lines. Trùng code (đua sinh mã) → cấp
// số mới thử lại tối đa vài lần.
func (s *Service) insertNewPO(ctx context.Context, store Store, draft domain.Draft) (CreatedPO, error) {
	const maxCodeRetry = 5
	for attempt := 0; attempt < maxCodeRetry; attempt++ {
		seq, serr := store.NextCodeSeq(ctx)
		if serr != nil {
			return CreatedPO{}, apperr.Internal("cấp số mã PO lỗi").WithCause(serr)
		}
		code := formatPOCode(seq)
		poID := id.NewString()
		saved, duplicate, ierr := store.InsertPO(ctx, NewPORow{ID: poID, Code: code, Draft: draft})
		if ierr != nil {
			return CreatedPO{}, apperr.Internal("ghi PO lỗi").WithCause(ierr)
		}
		if duplicate {
			continue // mã trùng (đua) → cấp số mới, thử lại
		}
		for _, l := range draft.Lines {
			if lerr := store.InsertPOItem(ctx, saved.ID, l); lerr != nil {
				return CreatedPO{}, apperr.Internal("ghi dòng PO lỗi").WithCause(lerr)
			}
		}
		got, gerr := store.GetCreated(ctx, saved.ID)
		if gerr != nil {
			return CreatedPO{}, apperr.Internal("đọc lại PO vừa tạo lỗi").WithCause(gerr)
		}
		return got, nil
	}
	return CreatedPO{}, apperr.Internal("không cấp được mã PO duy nhất sau nhiều lần thử")
}

// createIdempotentReadback đọc lại PO của luồng thắng (sau khi rollback PO thừa).
func (s *Service) createIdempotentReadback(ctx context.Context, idemKey string) (CreatedPO, error) {
	var out CreatedPO
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)
		existingID, found, ferr := store.FindByIdemKey(ctx, idemKey)
		if ferr != nil {
			return apperr.Internal("tra cứu idempotency (readback) lỗi").WithCause(ferr)
		}
		if !found {
			return apperr.Internal("idempotency key biến mất sau đua")
		}
		got, gerr := store.GetCreated(ctx, existingID)
		if gerr != nil {
			return apperr.Internal("đọc lại PO idempotent (readback) lỗi").WithCause(gerr)
		}
		out = got
		return nil
	})
	if err != nil {
		return CreatedPO{}, err
	}
	return out, nil
}

// ---- ConfirmPO (draft → ordered) ----

// ConfirmPO duyệt/đặt hàng: draft→ordered (FOR UPDATE + CanTransition + guard).
func (s *Service) ConfirmPO(ctx context.Context, poID string) (domain.PurchaseOrder, error) {
	return s.transition(ctx, poID, domain.StatusOrdered)
}

// CancelPO huỷ PO: draft/ordered→cancelled (chưa nhập kho — sau received cần đảo, HOÃN).
func (s *Service) CancelPO(ctx context.Context, poID string) (domain.PurchaseOrder, error) {
	return s.transition(ctx, poID, domain.StatusCancelled)
}

// transition là khung chung đổi trạng thái KHÔNG đụng kho/tiền/sổ (confirm/cancel):
// khoá + đọc header → CanTransition → UpdateStatus guard. Sai trạng thái → Conflict.
func (s *Service) transition(ctx context.Context, poID string, to domain.Status) (domain.PurchaseOrder, error) {
	if s.txm == nil || s.storeFromTx == nil {
		return domain.PurchaseOrder{}, apperr.Internal("Service chưa cấu hình đủ cho đổi trạng thái PO")
	}
	if strings.TrimSpace(poID) == "" {
		return domain.PurchaseOrder{}, apperr.Validation("id PO không được rỗng")
	}
	var out domain.PurchaseOrder
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)
		hdr, found, gerr := store.GetHeaderForUpdate(ctx, poID)
		if gerr != nil {
			return apperr.Internal("đọc header PO lỗi").WithCause(gerr)
		}
		if !found {
			return apperr.NotFound("đơn mua không tồn tại")
		}
		if !domain.CanTransition(domain.Status(hdr.Status), to) {
			return apperr.Conflict("không thể chuyển PO từ " + hdr.Status + " sang " + to.String())
		}
		rows, uerr := store.UpdateStatus(ctx, poID, hdr.Status, to.String())
		if uerr != nil {
			return apperr.Internal("cập nhật trạng thái PO lỗi").WithCause(uerr)
		}
		if rows == 0 {
			return apperr.Conflict("PO đã đổi trạng thái bởi thao tác khác")
		}
		got, derr := store.GetCreated(ctx, poID)
		if derr != nil {
			return apperr.Internal("đọc lại PO sau cập nhật lỗi").WithCause(derr)
		}
		out = got.PO
		return nil
	})
	if err != nil {
		return domain.PurchaseOrder{}, err
	}
	return out, nil
}

// ---- ReceivePO (ordered → received: nhập kho + post sổ) ----

// ReceiveInput là tham số nhận hàng (ordered→received). WarehouseID kho nhập.
// HasInvoice/Tax* cho HĐ mua (sổ TAX) — bỏ trống nếu mua không HĐ.
type ReceiveInput struct {
	POID        string
	WarehouseID int64
	// HasInvoice=true → post thêm bút toán sổ TAX theo số HĐ mua (mặc định = giá thực
	// khi không truyền Tax* riêng). FE/kế toán quyết.
	HasInvoice       bool
	TaxInventoryCost money.Money
	TaxVATAmount     money.Money
}

// ReceivedPO là kết quả ReceivePO: PO (received) + các lô đã tạo + tổng giá vốn nhập.
type ReceivedPO struct {
	PO            domain.PurchaseOrder
	CreatedBatch  []int64 // batchID các lô vừa tạo (đối soát)
	InventoryCost money.Money
	VATAmount     money.Money
}

// ReceivePO nhận hàng cho một PO đang ORDERED trong MỘT transaction atomic:
//  1. Khoá PO (FOR UPDATE) + đọc header; phải đang ordered (sai → Conflict).
//  2. Đọc dòng; mỗi dòng StockIn(tx, kho, sp, batch_code, expiry, mfg, qty,
//     inbound_price=unit_cost) → tạo lô + tăng tồn. InventoryCost += Σ line_total.
//  3. PurchaseEntries → post (INTERNAL Dr 1561+133/Cr 331; nếu HĐ thêm sổ TAX).
//  4. UpdateStatus ordered→received (guard). Replay sau commit: PO đã received →
//     bước (1) chặn (status != ordered → Conflict) → KHÔNG nhập kho/post lại.
func (s *Service) ReceivePO(ctx context.Context, in ReceiveInput) (ReceivedPO, error) {
	if s.txm == nil || s.storeFromTx == nil || s.stockIn == nil || s.poster == nil {
		return ReceivedPO{}, apperr.Internal("Service chưa cấu hình đủ cho ReceivePO")
	}
	if strings.TrimSpace(in.POID) == "" {
		return ReceivedPO{}, apperr.Validation("id PO không được rỗng")
	}
	if in.WarehouseID <= 0 {
		return ReceivedPO{}, apperr.Validation("warehouse_id không hợp lệ (phải > 0)")
	}

	var out ReceivedPO
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)

		// (1) Khoá + đọc header. Phải đang ORDERED.
		hdr, found, gerr := store.GetHeaderForUpdate(ctx, in.POID)
		if gerr != nil {
			return apperr.Internal("đọc header PO lỗi").WithCause(gerr)
		}
		if !found {
			return apperr.NotFound("đơn mua không tồn tại")
		}
		if hdr.Status != domain.StatusOrdered.String() {
			return apperr.Conflict("chỉ nhận hàng được PO đang ordered (hiện " + hdr.Status + ")")
		}
		if !domain.CanTransition(domain.Status(hdr.Status), domain.StatusReceived) {
			return apperr.Conflict("không thể chuyển PO từ " + hdr.Status + " sang received")
		}

		// (2) Đọc dòng + nhập kho từng dòng (gộp atomic).
		lines, lerr := store.ListLines(ctx, in.POID)
		if lerr != nil {
			return apperr.Internal("đọc dòng PO lỗi").WithCause(lerr)
		}
		if len(lines) == 0 {
			return apperr.Conflict("PO không có dòng hàng để nhập")
		}

		inventoryCost := money.Zero()
		vatAmount := money.Zero()
		batchIDs := make([]int64, 0, len(lines))
		for _, l := range lines {
			batchID, serr := s.stockIn.StockIn(ctx, tx, in.WarehouseID, l.ProductID, l.BatchCode,
				l.ExpiryDate, l.ManufacturingDate, l.Quantity, l.UnitCost)
			if serr != nil {
				return serr // đã là apperr → tx ROLLBACK (không nhập dòng nào).
			}
			batchIDs = append(batchIDs, batchID)
			inventoryCost = inventoryCost.Add(l.LineTotal)
			vatAmount = vatAmount.Add(l.LineTotal.Mul(l.VATRate.Decimal()).RoundVND())
		}

		// (3) Post bút toán nhập kho (Dr 1561+133/Cr 331). Sổ TAX nếu có HĐ mua.
		taxInv := in.TaxInventoryCost
		taxVAT := in.TaxVATAmount
		if in.HasInvoice && taxInv.IsZero() && taxVAT.IsZero() {
			// HĐ mua không truyền số riêng → dùng giá thực (sổ TAX = INTERNAL).
			taxInv = inventoryCost
			taxVAT = vatAmount
		}
		entries := s.rules.PurchaseEntries(posting.PurchaseInput{
			SourceID:         hdr.Code,
			Date:             s.now(),
			InventoryCost:    inventoryCost,
			VATAmount:        vatAmount,
			HasInvoice:       in.HasInvoice,
			TaxInventoryCost: taxInv,
			TaxVATAmount:     taxVAT,
		})
		if perr := s.postEntries(ctx, tx, entries); perr != nil {
			return perr
		}

		// (4) ordered→received (guard status cũ).
		rows, uerr := store.UpdateStatus(ctx, in.POID, hdr.Status, domain.StatusReceived.String())
		if uerr != nil {
			return apperr.Internal("cập nhật trạng thái nhận hàng lỗi").WithCause(uerr)
		}
		if rows == 0 {
			return apperr.Conflict("PO đã đổi trạng thái bởi thao tác khác")
		}

		got, derr := store.GetCreated(ctx, in.POID)
		if derr != nil {
			return apperr.Internal("đọc lại PO sau nhận hàng lỗi").WithCause(derr)
		}
		out = ReceivedPO{
			PO:            got.PO,
			CreatedBatch:  batchIDs,
			InventoryCost: inventoryCost,
			VATAmount:     vatAmount,
		}
		return nil
	})
	if err != nil {
		return ReceivedPO{}, err
	}
	return out, nil
}

// ---- PaySupplier (received → paid: chi NCC + post sổ) ----

// PaySupplierInput là tham số chi trả NCC. Amount > 0. FundAccountID quỹ xuất tiền;
// FundIsBank chọn 111 (mặt) vs 112 (ngân hàng). BookType sổ phiếu (BOTH cho HĐ mua).
type PaySupplierInput struct {
	POID          string
	Amount        money.Money
	FundAccountID int64
	FundIsBank    bool
	BookType      financedomain.BookType
	BankRef       *string
	IdemKey       string
	CreatedBy     *string
}

// PaySupplierResult là kết quả PaySupplier: phiếu chi + PO (paid).
type PaySupplierResult struct {
	Payment financedomain.Payment
	PO      domain.PurchaseOrder
}

// PaySupplier chi trả NCC cho một PO đang RECEIVED trong MỘT transaction atomic:
//  1. Khoá PO (FOR UPDATE) + đọc header; phải đang received (sai → Conflict).
//  2. finance.RecordPaymentOut(tx, ...) (idempotent — trả created). Replay cùng
//     Idempotency-Key (created=false) → KHÔNG post sổ lại (bài học review).
//  3. CHỈ post BuildSupplierPayment (Dr 331/Cr 111/112) khi created.
//  4. received→paid (guard). Replay sau commit: PO đã paid → bước (1) chặn (status
//     != received → Conflict) → an toàn kép (cờ created + state machine).
func (s *Service) PaySupplier(ctx context.Context, in PaySupplierInput) (PaySupplierResult, error) {
	if s.txm == nil || s.storeFromTx == nil || s.payer == nil || s.poster == nil {
		return PaySupplierResult{}, apperr.Internal("Service chưa cấu hình đủ cho PaySupplier")
	}
	if strings.TrimSpace(in.POID) == "" {
		return PaySupplierResult{}, apperr.Validation("id PO không được rỗng")
	}
	if !in.Amount.IsPositive() {
		return PaySupplierResult{}, apperr.Validation("số tiền chi phải lớn hơn 0")
	}
	bookType := in.BookType
	if !bookType.Valid() {
		bookType = financedomain.BookBoth // mặc định HĐ mua VAT.
	}

	var out PaySupplierResult
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)

		hdr, found, gerr := store.GetHeaderForUpdate(ctx, in.POID)
		if gerr != nil {
			return apperr.Internal("đọc header PO lỗi").WithCause(gerr)
		}
		if !found {
			return apperr.NotFound("đơn mua không tồn tại")
		}
		if hdr.Status != domain.StatusReceived.String() {
			return apperr.Conflict("chỉ chi trả được PO đang received (hiện " + hdr.Status + ")")
		}

		// (2) Ghi phiếu chi idempotent (finance lo dedup theo bank_ref/idem key).
		payment, created, rerr := s.payer.RecordPaymentOut(ctx, tx, financeapp.RecordPaymentOutParams{
			RecordPaymentOut: financedomain.RecordPaymentOut{
				POCode:        hdr.Code,
				Amount:        in.Amount,
				FundAccountID: in.FundAccountID,
				BookType:      bookType,
				BankRef:       in.BankRef,
				CreatedBy:     in.CreatedBy,
			},
			IdemKey: in.IdemKey,
		})
		if rerr != nil {
			return rerr
		}

		// (3) Post phiếu chi theo sổ — CHỈ khi phiếu VỪA tạo (created). Replay cùng
		// Idempotency-Key (created=false) KHÔNG post lại → tránh nhân đôi bút toán.
		if created {
			if perr := s.postSupplierPaymentEntries(ctx, tx, bookType, in.FundIsBank, in.Amount, hdr.Code); perr != nil {
				return perr
			}
		}

		// (4) received→paid (guard).
		rows, uerr := store.UpdateStatus(ctx, in.POID, hdr.Status, domain.StatusPaid.String())
		if uerr != nil {
			return apperr.Internal("cập nhật trạng thái chi trả lỗi").WithCause(uerr)
		}
		if rows == 0 {
			return apperr.Conflict("PO đã đổi trạng thái bởi thao tác khác")
		}

		got, derr := store.GetCreated(ctx, in.POID)
		if derr != nil {
			return apperr.Internal("đọc lại PO sau chi trả lỗi").WithCause(derr)
		}
		out = PaySupplierResult{Payment: payment, PO: got.PO}
		return nil
	})
	if err != nil {
		return PaySupplierResult{}, err
	}
	return out, nil
}

// ---- helpers ----

// postEntries post lần lượt mỗi JournalEntry qua Poster trong tx (gộp atomic). Lỗi
// post (kỳ khoá / TK không cho hạch toán / lệch Σ) trả nguyên (đã apperr) → ROLLBACK.
func (s *Service) postEntries(ctx context.Context, tx pgx.Tx, entries []accountingdomain.JournalEntry) error {
	for _, e := range entries {
		if _, perr := s.poster.Post(ctx, tx, e); perr != nil {
			return perr
		}
	}
	return nil
}

// postSupplierPaymentEntries post bút toán chi NCC theo sổ: INTERNAL/BOTH → sổ
// INTERNAL; TAX/BOTH → sổ TAX (BOTH post CẢ hai). Đối xứng postPaymentEntries bán hàng.
func (s *Service) postSupplierPaymentEntries(ctx context.Context, tx pgx.Tx, bookType financedomain.BookType, fundIsBank bool, amount money.Money, code string) error {
	date := s.now()
	if bookType == financedomain.BookInternal || bookType == financedomain.BookBoth {
		e := s.rules.BuildSupplierPayment(accountingdomain.BookInternal, fundIsBank, amount, code, date)
		if _, perr := s.poster.Post(ctx, tx, e); perr != nil {
			return perr
		}
	}
	if bookType == financedomain.BookTax || bookType == financedomain.BookBoth {
		e := s.rules.BuildSupplierPayment(accountingdomain.BookTax, fundIsBank, amount, code, date)
		if _, perr := s.poster.Post(ctx, tx, e); perr != nil {
			return perr
		}
	}
	return nil
}

// ---- ĐỌC (cho HTTP GET) ----

// GetPO trả một PO (header + lines) theo id — đọc trong tx ngắn.
func (s *Service) GetPO(ctx context.Context, poID string) (CreatedPO, error) {
	if s.txm == nil || s.storeFromTx == nil {
		return CreatedPO{}, apperr.Internal("Service chưa cấu hình đủ cho GetPO")
	}
	if strings.TrimSpace(poID) == "" {
		return CreatedPO{}, apperr.Validation("id PO không được rỗng")
	}
	var out CreatedPO
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		got, gerr := s.storeFromTx(tx).GetCreated(ctx, poID)
		if gerr != nil {
			return gerr
		}
		out = got
		return nil
	})
	if err != nil {
		return CreatedPO{}, err
	}
	return out, nil
}
