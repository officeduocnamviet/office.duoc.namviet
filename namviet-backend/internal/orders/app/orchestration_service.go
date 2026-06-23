// orchestration_service.go — P4b: 3 use-case GHI GỘP CROSS-MODULE, mỗi cái = 1
// transaction ATOMIC (platform/db.WithinTx qua TxManager). Mọi primitive
// (DeductFEFO / IssueInvoice / Poster.Post / RecordPaymentIn) nhận CÙNG pgx.Tx →
// lỗi bất kỳ bước nào ROLLBACK CẢ CỤM (không trừ kho dở, không HĐ mồ côi, không
// bút toán lệch sự kiện). orders/app điều phối; KHÔNG sửa logic module khác (chỉ
// GỌI port). Tiền decimal toàn tuyến (CẤM float).
//
// GIẢ ĐỊNH ĐÃ GẮN CỜ (kế toán/BA xác nhận trước cutover — xem
// docs/superpowers/specs/2026-06-22-tt133-posting-rules.md §5/§6):
//   - VAT sổ INTERNAL tính từ GIÁ THỰC của đơn (realInternalVAT: Σ line_total ×
//     thuế suất, RoundVND mỗi dòng) — ĐỘC LẬP giá HĐ. Sổ TAX dùng VAT/doanh thu của
//     HĐ (invoice.Subtotal/VATAmount). Khi FE override giá HĐ ≠ giá thực, hai sổ ghi
//     KHÁC nhau ĐÚNG dual-ledger. Khi không override + cùng thuế suất → hai sổ bằng nhau.
//     ⚠️ HẠN CHẾ: HĐ (P5) chưa mô hình hoá chiết khấu dòng → khi đơn có discount,
//     subtotal HĐ (Σ qty×giá) khác final đơn (sau CK); cần kế toán xác nhận HĐ có
//     thể hiện chiết khấu không (mở rộng vat.LineInput sau).
//   - COGS sổ INTERNAL = Σ (qty lô × inbound_price lô) — inbound_price là giá vốn
//     MỖI ĐƠN VỊ (xác nhận từ inventory write_repo + seed test: ConsumedBatch.
//     InboundPrice giữ nguyên giá per-unit của lô, KHÔNG nhân số lượng). Sổ TAX
//     KHÔNG ghi COGS (posting.DefaultRules.TaxRecordsCOGS=false).
//   - MST khách HĐ: input CustomerTaxCode (FE/orders lấy từ customers.b2b_metadata
//     qua port sau — cross-context, KHÔNG đụng bảng customers ở đây). B2B bắt buộc
//     có MST; B2C POS chỉ xuất HĐ khi khách yêu cầu (IssueInvoice=true + có MST).
//   - Serial HĐ: input từ FE (ký hiệu hoá đơn theo dải đăng ký với CQT).
package app

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	accountingdomain "github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/posting"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	financeapp "github.com/Maneva-AI/namviet-backend/internal/finance/app"
	financedomain "github.com/Maneva-AI/namviet-backend/internal/finance/domain"
	inventorydomain "github.com/Maneva-AI/namviet-backend/internal/inventory/domain"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
	vatapp "github.com/Maneva-AI/namviet-backend/internal/vat/app"
	vatdomain "github.com/Maneva-AI/namviet-backend/internal/vat/domain"
)

// Orchestrator là use-case GHI GỘP cross-module của orders (P4b). Giữ các port
// nội bộ của inventory/accounting/vat/finance + store orders bound-tx + Rules
// định khoản (posting.DefaultRules mặc định). nil port → use-case tương ứng trả
// Internal (an toàn cho test từng phần).
type Orchestrator struct {
	orchFromTx OrchestrationStoreFromTx
	txm        TxManager
	deductor   Deductor
	poster     Poster
	issuer     InvoiceIssuer
	recorder   PaymentRecorder
	rules      posting.Rules
	now        func() time.Time
}

// NewOrchestrator dựng Orchestrator. rules thường là posting.DefaultRules (mã TK +
// cờ per-book gom 1 chỗ). now cho phép test bơm thời gian; nil → time.Now.
func NewOrchestrator(
	orchFromTx OrchestrationStoreFromTx,
	txm TxManager,
	deductor Deductor,
	poster Poster,
	issuer InvoiceIssuer,
	recorder PaymentRecorder,
	rules posting.Rules,
) *Orchestrator {
	return &Orchestrator{
		orchFromTx: orchFromTx,
		txm:        txm,
		deductor:   deductor,
		poster:     poster,
		issuer:     issuer,
		recorder:   recorder,
		rules:      rules,
		now:        time.Now,
	}
}

// VATLine là thuế suất + (tuỳ chọn) đơn giá HĐ cho MỘT dòng khi xuất HĐ. FE cung
// cấp khi giao hàng. UnitPrice rỗng (zero + !HasPrice) → dùng đơn giá đơn (sổ
// TAX = sổ INTERNAL về giá). VATRate là decimal (vd 0.08), KHÔNG float.
type VATLine struct {
	VATRate     decimal.Decimal
	UnitPrice   money.Money
	HasPrice    bool
	Description string
}

// ShipInput là tham số giao hàng (CONFIRMED→SHIPPING). WarehouseID là kho xuất.
// CustomerTaxCode là MST khách cho HĐ VAT (B2B bắt buộc). Serial là ký hiệu HĐ.
// VATLines song song với dòng đơn (theo THỨ TỰ dòng đọc ra); nếu rỗng → mỗi dòng
// VAT 0% (không khuyến nghị B2B — nên truyền đủ). FundIsBank/OnCredit định vế tiền
// post sổ (B2B công nợ → OnCredit=true mặc định).
type ShipInput struct {
	OrderID         string
	WarehouseID     int64
	CustomerTaxCode string
	Serial          string
	MauSo           string
	VATLines        []VATLine
}

// ShippedOrder là kết quả ShipOrder: đơn (header sau cập nhật SHIPPING) + HĐ đã
// phát hành + các lô đã trừ (cho FE hiển thị/đối soát) + giá vốn tổng.
type ShippedOrder struct {
	Order    domain.Order
	Invoice  vatdomain.IssuedInvoice
	Consumed []inventorydomain.ConsumedBatch
	COGS     money.Money
}

// ShipOrder giao hàng cho một đơn CONFIRMED trong MỘT transaction atomic:
//  1. Khoá đơn (FOR UPDATE) + đọc header; phải đang CONFIRMED (sai → Conflict).
//  2. Đọc dòng hàng; mỗi dòng DeductFEFO(tx, kho, sp, qty) → gom ConsumedBatch.
//     COGS += Σ (qty × inbound_price). Thiếu tồn → DeductFEFO trả Conflict → TOÀN
//     TX ROLLBACK (không trừ dòng nào, không HĐ, không bút toán, status giữ nguyên).
//  3. Phát hành HĐ VAT (B2B 100%): IssueInvoice(tx, ...) — subtotal/vat/total sạch.
//  4. Dựng SaleInput → rules.SaleEntries → post từng entry (INTERNAL revenue+COGS,
//     TAX revenue) trong tx. Mỗi entry tự cân Σ (domain + trigger DB).
//  5. UpdateStatus CONFIRMED→SHIPPING (guard). Trả đơn + HĐ + lô tiêu thụ.
func (o *Orchestrator) ShipOrder(ctx context.Context, in ShipInput) (ShippedOrder, error) {
	if o.txm == nil || o.orchFromTx == nil || o.deductor == nil || o.poster == nil || o.issuer == nil {
		return ShippedOrder{}, apperr.Internal("Orchestrator chưa cấu hình đủ cho ShipOrder")
	}
	if strings.TrimSpace(in.OrderID) == "" {
		return ShippedOrder{}, apperr.Validation("id đơn không được rỗng")
	}
	if in.WarehouseID <= 0 {
		return ShippedOrder{}, apperr.Validation("warehouse_id không hợp lệ (phải > 0)")
	}

	var out ShippedOrder
	err := o.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := o.orchFromTx(tx)

		// (1) Khoá + đọc header. Phải đang CONFIRMED.
		hdr, found, gerr := store.GetHeaderForUpdate(ctx, in.OrderID)
		if gerr != nil {
			return apperr.Internal("đọc header đơn lỗi").WithCause(gerr)
		}
		if !found {
			return apperr.NotFound("đơn hàng không tồn tại")
		}
		if hdr.Status != domain.StatusConfirmed.String() {
			return apperr.Conflict("chỉ giao được đơn đang CONFIRMED (hiện " + hdr.Status + ")")
		}
		// State machine THUẦN: CONFIRMED→SHIPPING phải hợp lệ (phòng thủ).
		if !domain.CanTransition(domain.Status(hdr.Status), domain.StatusShipping) {
			return apperr.Conflict("không thể chuyển đơn từ " + hdr.Status + " sang SHIPPING")
		}

		// (2) Đọc dòng + trừ kho FEFO từng dòng (gộp atomic). COGS = Σ qty×inbound_price.
		lines, lerr := store.ListLinesForOrder(ctx, in.OrderID)
		if lerr != nil {
			return apperr.Internal("đọc dòng hàng lỗi").WithCause(lerr)
		}
		if len(lines) == 0 {
			return apperr.Conflict("đơn không có dòng hàng để giao")
		}

		cogs := money.Zero()
		consumed := make([]inventorydomain.ConsumedBatch, 0, len(lines))
		for _, l := range lines {
			batches, derr := o.deductor.DeductFEFO(ctx, tx, in.WarehouseID, l.ProductID, l.Quantity)
			if derr != nil {
				// Thiếu tồn (Conflict) / lỗi khác → trả nguyên (đã là apperr) → tx ROLLBACK.
				return derr
			}
			for _, b := range batches {
				// COGS per-unit: qty lô × giá vốn nhập per-unit của lô.
				cogs = cogs.Add(b.InboundPrice.Mul(b.Quantity.Decimal()))
			}
			consumed = append(consumed, batches...)
		}

		// (3) Phát hành HĐ VAT (B2B 100% đơn; MST bắt buộc — domain vat ép).
		invLines := buildInvoiceLines(lines, in.VATLines)
		invoice, ierr := o.issuer.IssueInvoice(ctx, tx, vatapp.IssueParams{
			OrderCode:       hdr.Code,
			CustomerTaxCode: in.CustomerTaxCode,
			Serial:          in.Serial,
			MauSo:           in.MauSo,
			IssueDate:       o.now(),
			Lines:           invLines,
		})
		if ierr != nil {
			return ierr // đã là apperr (vd MST rỗng → Validation) → tx ROLLBACK.
		}

		// (4) Dựng SaleInput + post bút toán 2 sổ.
		//  - Sổ INTERNAL (giá THỰC): doanh thu = final_amount (Σ line_total ex-VAT);
		//    VAT = realInternalVAT (tính từ giá thực × thuế suất từng dòng) → nhất
		//    quán dù FE override giá HĐ (sửa lệch #2/#3 review).
		//  - Sổ TAX (giá HĐ): doanh thu/VAT lấy từ invoice (Subtotal/VATAmount).
		internalLineAmts := make([]money.Money, len(lines))
		for i, l := range lines {
			internalLineAmts[i] = l.LineTotal
		}
		internalVAT := realInternalVAT(internalLineAmts, in.VATLines)
		entries := o.rules.SaleEntries(posting.SaleInput{
			SourceID:        hdr.Code,
			Date:            o.now(),
			OnCredit:        true,  // B2B bán chịu (công nợ) — RecordPayment ghi thu sau.
			FundIsBank:      false, // không dùng khi OnCredit=true (vế nợ là 131).
			RevenueExVAT:    hdr.FinalAmount,
			VATAmount:       internalVAT,
			COGS:            cogs,
			HasInvoice:      true,
			TaxRevenueExVAT: invoice.Subtotal,
			TaxVATAmount:    invoice.VATAmount,
			TaxCOGS:         cogs,
		})
		if perr := o.postEntries(ctx, tx, store, entries); perr != nil {
			return perr
		}

		// (5) Đổi trạng thái CONFIRMED→SHIPPING (guard status cũ).
		rows, uerr := store.UpdateStatus(ctx, in.OrderID, hdr.Status, domain.StatusShipping.String())
		if uerr != nil {
			return apperr.Internal("cập nhật trạng thái giao hàng lỗi").WithCause(uerr)
		}
		if rows == 0 {
			return apperr.Conflict("đơn đã đổi trạng thái bởi thao tác khác")
		}

		out = ShippedOrder{
			Order:    orchHeaderToOrder(hdr, domain.StatusShipping.String(), hdr.PaymentStatus),
			Invoice:  invoice,
			Consumed: consumed,
			COGS:     cogs,
		}
		return nil
	})
	if err != nil {
		return ShippedOrder{}, err
	}
	return out, nil
}

// RecordPaymentInput là tham số ghi thu cho một đơn. Amount > 0. FundAccountID là
// quỹ nhận tiền; FundIsBank chọn 111 (mặt) vs 112 (ngân hàng) khi post. BookType
// chọn sổ phiếu (BOTH cho B2B có HĐ). BankRef cho webhook; IdemKey cho thủ công.
type RecordPaymentInput struct {
	OrderID       string
	Amount        money.Money
	FundAccountID int64
	FundIsBank    bool
	BookType      financedomain.BookType
	BankRef       *string
	IdemKey       string
	CreatedBy     *string
	// Collected (thanh toán 2 bước, spec mục 55): true = NV đã thu từ khách nhưng
	// CHƯA nộp quỹ → phiếu 'pending' (nợ khách giảm NGAY vì đã-thu đếm cả pending;
	// số dư quỹ CHƯA tăng — thủ quỹ ConfirmReceipt sau). false (mặc định) = thu
	// thẳng vào quỹ → 'completed'.
	Collected bool
}

// PaymentResult là kết quả RecordPayment: phiếu đã ghi + payment_status mới của đơn.
type PaymentResult struct {
	Payment       financedomain.Payment
	PaymentStatus string
	Order         domain.Order
}

// RecordPayment ghi MỘT phiếu THU cho đơn + post PAYMENT_IN + tính lại
// payment_status, TRONG 1 transaction atomic:
//  1. Khoá đơn (FOR UPDATE) + đọc header (mã đơn cho ref, tổng phải trả).
//  2. finance.RecordPaymentIn(tx, ...) (idempotent — finance lo dedup, KHÔNG cộng
//     tiền 2 lần). Idempotent hit → vẫn tính lại payment_status (an toàn).
//  3. Post bút toán PAYMENT_IN theo sổ (INTERNAL/BOTH→INTERNAL; TAX/BOTH→TAX).
//  4. Tính lại đã-thu (sổ thực tế) TRONG tx → DerivePaymentStatus → UPDATE đơn.
func (o *Orchestrator) RecordPayment(ctx context.Context, in RecordPaymentInput) (PaymentResult, error) {
	if o.txm == nil || o.orchFromTx == nil || o.recorder == nil || o.poster == nil {
		return PaymentResult{}, apperr.Internal("Orchestrator chưa cấu hình đủ cho RecordPayment")
	}
	if strings.TrimSpace(in.OrderID) == "" {
		return PaymentResult{}, apperr.Validation("id đơn không được rỗng")
	}

	var out PaymentResult
	err := o.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := o.orchFromTx(tx)

		hdr, found, gerr := store.GetHeaderForUpdate(ctx, in.OrderID)
		if gerr != nil {
			return apperr.Internal("đọc header đơn lỗi").WithCause(gerr)
		}
		if !found {
			return apperr.NotFound("đơn hàng không tồn tại")
		}

		bookType := in.BookType
		if !bookType.Valid() {
			bookType = financedomain.BookBoth // mặc định B2B có HĐ VAT.
		}

		// (2) Ghi phiếu thu idempotent (finance lo dedup theo bank_ref/idem key).
		payment, created, rerr := o.recorder.RecordPaymentIn(ctx, tx, financeapp.RecordPaymentInParams{
			RecordPaymentIn: financedomain.RecordPaymentIn{
				OrderCode:     hdr.Code,
				Amount:        in.Amount,
				FundAccountID: in.FundAccountID,
				BookType:      bookType,
				BankRef:       in.BankRef,
				CreatedBy:     in.CreatedBy,
				InitialStatus: paymentInitialStatus(in.Collected),
			},
			IdemKey: in.IdemKey,
		})
		if rerr != nil {
			return rerr // đã là apperr → tx ROLLBACK.
		}

		// (3) Post PAYMENT_IN theo sổ — CHỈ khi phiếu VỪA tạo (created). Replay cùng
		// Idempotency-Key (created=false) KHÔNG post lại → tránh nhân đôi bút toán.
		if created {
			if perr := o.postPaymentEntries(ctx, tx, store, bookType, in.FundIsBank, in.Amount, hdr.Code); perr != nil {
				return perr
			}
		}

		// (4) Tính lại đã-thu TRONG tx (phiếu vừa ghi ĐƯỢC tính) → payment_status.
		paid, serr := store.SumPaidInTx(ctx, hdr.Code)
		if serr != nil {
			return apperr.Internal("tính lại đã thu lỗi").WithCause(serr)
		}
		ps := domain.DerivePaymentStatus(hdr.FinalAmount, paid)
		if _, uerr := store.UpdatePaymentStatus(ctx, in.OrderID, ps.String()); uerr != nil {
			return apperr.Internal("cập nhật payment_status lỗi").WithCause(uerr)
		}

		out = PaymentResult{
			Payment:       payment,
			PaymentStatus: ps.String(),
			Order:         orchHeaderToOrder(hdr, hdr.Status, ps.String()),
		}
		return nil
	})
	if err != nil {
		return PaymentResult{}, err
	}
	return out, nil
}

// PosSaleInput là tham số bán lẻ tại quầy (B2C atomic). Lines là dòng hàng (tái
// dùng domain.DraftLine — tính tiền THUẦN ở NewDraft). WarehouseID kho xuất.
// FundAccountID/FundIsBank quỹ thu. IssueInvoice=true → xuất HĐ (cần MST+serial).
type PosSaleInput struct {
	CustomerID      *int64
	CreatorID       string
	Note            string
	Lines           []domain.DraftLine
	WarehouseID     int64
	FundAccountID   int64
	FundIsBank      bool
	IdemKey         string
	IssueInvoice    bool
	CustomerTaxCode string
	Serial          string
	MauSo           string
	VATLines        []VATLine // song song Lines (theo thứ tự) khi IssueInvoice.
}

// PosSaleResult là kết quả CreatePosSale: đơn (COMPLETED, paid) + phiếu thu + HĐ
// (nếu xuất) + lô tiêu thụ + giá vốn.
type PosSaleResult struct {
	Order    domain.Order
	Lines    []domain.OrderLine
	Payment  financedomain.Payment
	Invoice  *vatdomain.IssuedInvoice
	Consumed []inventorydomain.ConsumedBatch
	COGS     money.Money
}

// CreatePosSale bán lẻ tại quầy trong MỘT transaction atomic: tạo đơn (B2C) →
// DeductFEFO → (HĐ nếu yêu cầu) → RecordPaymentIn (thu đủ) → post toàn bộ bút
// toán (revenue+COGS+payment) → status COMPLETED + payment_status=paid. Cùng
// advisory lock chống bán âm (qua DeductFEFO). Lỗi bước nào → ROLLBACK cả cụm.
//
// Cần OrderStore (P4a) để tạo đơn + OrchestrationStore (P4b) để đổi trạng thái/
// payment_status — cả hai bound CÙNG tx (storeFromTx + orchFromTx trên cùng tx).
func (o *Orchestrator) CreatePosSale(ctx context.Context, storeFromTx OrderStoreFromTx, in PosSaleInput) (PosSaleResult, error) {
	if o.txm == nil || o.orchFromTx == nil || o.deductor == nil || o.poster == nil || o.recorder == nil || storeFromTx == nil {
		return PosSaleResult{}, apperr.Internal("Orchestrator chưa cấu hình đủ cho CreatePosSale")
	}
	if in.WarehouseID <= 0 {
		return PosSaleResult{}, apperr.Validation("warehouse_id không hợp lệ (phải > 0)")
	}
	if in.IssueInvoice && strings.TrimSpace(in.CustomerTaxCode) == "" {
		return PosSaleResult{}, apperr.Validation("xuất HĐ cần MST khách (customer_tax_code)")
	}

	// Dựng draft THUẦN (validate + tính tiền). POS = B2C.
	draft, derr := domain.NewDraft(domain.DraftInput{
		CustomerID: in.CustomerID,
		OrderType:  domain.OrderTypeB2C,
		CreatorID:  strings.TrimSpace(in.CreatorID),
		Note:       in.Note,
		Lines:      in.Lines,
	})
	if derr != nil {
		return PosSaleResult{}, apperr.Validation(derr.Error())
	}

	var out PosSaleResult
	err := o.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		ostore := storeFromTx(tx) // P4a store (tạo đơn)
		store := o.orchFromTx(tx) // P4b store (đổi trạng thái)

		// (1) Tạo đơn (PENDING) + idempotency.
		created, wasCreated, cerr := o.createPosOrder(ctx, ostore, draft, in.IdemKey)
		if cerr != nil {
			return cerr
		}
		// Replay cùng Idempotency-Key: đơn POS đã hoàn tất trước đó → KHÔNG trừ kho /
		// thu tiền / post sổ LẠI (chống nhân đôi). Trả đơn + dòng đã có rồi dừng.
		if !wasCreated {
			out = PosSaleResult{Order: created.Order, Lines: created.Lines}
			return nil
		}
		code := created.Order.Code

		// (2) Trừ kho FEFO từng dòng (advisory lock trong DeductFEFO).
		cogs := money.Zero()
		consumed := make([]inventorydomain.ConsumedBatch, 0, len(draft.Lines))
		for _, l := range draft.Lines {
			batches, dferr := o.deductor.DeductFEFO(ctx, tx, in.WarehouseID, l.ProductID, toInventoryQty(l.Quantity))
			if dferr != nil {
				return dferr // thiếu tồn → ROLLBACK cả cụm (đơn vừa tạo cũng huỷ).
			}
			for _, b := range batches {
				cogs = cogs.Add(b.InboundPrice.Mul(b.Quantity.Decimal()))
			}
			consumed = append(consumed, batches...)
		}

		// (3) HĐ VAT nếu khách yêu cầu (B2C xuất khi cần — Q1).
		var invoicePtr *vatdomain.IssuedInvoice
		vatAmount := money.Zero()
		taxSubtotal := money.Zero()
		if in.IssueInvoice {
			if o.issuer == nil {
				return apperr.Internal("Orchestrator chưa cấu hình issuer cho HĐ POS")
			}
			invLines := buildInvoiceLinesFromDraft(draft.Lines, in.VATLines)
			invoice, ierr := o.issuer.IssueInvoice(ctx, tx, vatapp.IssueParams{
				OrderCode:       code,
				CustomerTaxCode: in.CustomerTaxCode,
				Serial:          in.Serial,
				MauSo:           in.MauSo,
				IssueDate:       o.now(),
				Lines:           invLines,
			})
			if ierr != nil {
				return ierr
			}
			invoicePtr = &invoice
			vatAmount = invoice.VATAmount
			taxSubtotal = invoice.Subtotal
		}

		// (4) Post bút toán bán hàng. POS thu ngay → vế tiền 111/112 (OnCredit=false).
		// Sổ INTERNAL: VAT tính từ giá THỰC (realInternalVAT) — nhất quán dù override
		// giá HĐ. Sổ TAX: VAT/doanh thu từ invoice. Không HĐ → internalVAT=0 (VATLines rỗng).
		posLineAmts := make([]money.Money, len(draft.Lines))
		for i, l := range draft.Lines {
			posLineAmts[i] = l.LineTotal
		}
		internalVAT := realInternalVAT(posLineAmts, in.VATLines)
		entries := o.rules.SaleEntries(posting.SaleInput{
			SourceID:        code,
			Date:            o.now(),
			OnCredit:        false,
			FundIsBank:      in.FundIsBank,
			RevenueExVAT:    draft.FinalAmount,
			VATAmount:       internalVAT,
			COGS:            cogs,
			HasInvoice:      in.IssueInvoice,
			TaxRevenueExVAT: taxSubtotal,
			TaxVATAmount:    vatAmount,
			TaxCOGS:         cogs,
		})
		if perr := o.postEntries(ctx, tx, store, entries); perr != nil {
			return perr
		}

		// (5) Ghi phiếu THU (thu đủ = final_amount + VAT nếu có HĐ). POS thu ngay nên
		// KHÔNG sinh PAYMENT_IN riêng (đã gộp vế tiền vào bút toán SALE ở (4)) — phiếu
		// thu vẫn ghi để có dòng tiền finance_transactions + suy payment_status=paid.
		// Số thu = tổng khách trả: final (ex-VAT) + VAT HĐ (nếu có).
		// Khách trả tiền THỰC = doanh thu thực + VAT thực (internalVAT), KHÔNG phải VAT
		// hóa đơn (khi override giá HĐ thì tiền khách trả theo giá thực). Sửa #3 review.
		amountReceived := draft.FinalAmount.Add(internalVAT)
		bookType := financedomain.BookInternal
		if in.IssueInvoice {
			bookType = financedomain.BookBoth
		}
		// created bỏ qua: CreatePosSale đã chặn replay ở wasCreated (đơn mới mỗi lần);
		// phiếu POS thu ngay, KHÔNG post PAYMENT_IN riêng (vế tiền nằm trong bút toán SALE).
		payment, _, prerr := o.recorder.RecordPaymentIn(ctx, tx, financeapp.RecordPaymentInParams{
			RecordPaymentIn: financedomain.RecordPaymentIn{
				OrderCode:     code,
				Amount:        amountReceived,
				FundAccountID: in.FundAccountID,
				BookType:      bookType,
			},
			IdemKey: posPaymentIdemKey(in.IdemKey, code),
		})
		if prerr != nil {
			return prerr
		}

		// (6) Đổi trạng thái → COMPLETED + payment_status=paid (thu đủ).
		paid, serr := store.SumPaidInTx(ctx, code)
		if serr != nil {
			return apperr.Internal("tính lại đã thu (POS) lỗi").WithCause(serr)
		}
		ps := domain.DerivePaymentStatus(amountReceived, paid)
		if _, uerr := store.UpdatePaymentStatus(ctx, created.Order.ID, ps.String()); uerr != nil {
			return apperr.Internal("cập nhật payment_status (POS) lỗi").WithCause(uerr)
		}
		// PENDING→...→COMPLETED: POS hoàn tất ngay. Đi qua các bước hợp lệ của state
		// machine để không nhảy bước (PENDING→CONFIRMED→SHIPPING→COMPLETED).
		for _, step := range []domain.Status{domain.StatusConfirmed, domain.StatusShipping, domain.StatusCompleted} {
			cur := posPrevStatus(step)
			rows, uerr := store.UpdateStatus(ctx, created.Order.ID, cur.String(), step.String())
			if uerr != nil {
				return apperr.Internal("cập nhật trạng thái POS lỗi").WithCause(uerr)
			}
			if rows == 0 {
				return apperr.Conflict("đơn POS đã đổi trạng thái bởi thao tác khác")
			}
		}

		finalOrder := created.Order
		finalOrder.Status = domain.StatusCompleted.String()
		finalOrder.PaymentStatus = ps.String()
		out = PosSaleResult{
			Order:    finalOrder,
			Lines:    created.Lines,
			Payment:  payment,
			Invoice:  invoicePtr,
			Consumed: consumed,
			COGS:     cogs,
		}
		return nil
	})
	if err != nil {
		return PosSaleResult{}, err
	}
	return out, nil
}

// LumpSumInput: thu 1 CỤC từ khách → phân bổ cho các đơn CHƯA tất toán của khách,
// CŨ NHẤT trước (spec mục 55). Phiếu gắn ref_type='customer' (KHÔNG 'order') để
// query đã-thu-trực-tiếp-theo-đơn không đếm trùng; chỉ allocation được đếm.
type LumpSumInput struct {
	CustomerID    int64
	Amount        money.Money
	FundAccountID int64
	FundIsBank    bool
	BookType      financedomain.BookType
	Collected     bool // true = NV thu chưa nộp quỹ (pending); false = vào quỹ ngay (completed)
	BankRef       *string
	IdemKey       string
	CreatedBy     *string
}

// AllocationLine: 1 dòng phân bổ phiếu lump-sum cho 1 đơn.
type AllocationLine struct {
	OrderID   string
	OrderCode string
	Amount    money.Money
}

// LumpSumResult: phiếu đã ghi + các dòng phân bổ + phần dư chưa phân bổ (thu thừa →
// credit khách, KHÔNG gắn đơn nào).
type LumpSumResult struct {
	Payment     financedomain.Payment
	Allocations []AllocationLine
	Leftover    money.Money
}

// RecordLumpSumPayment ghi MỘT phiếu THU lump-sum cho khách + phân bổ cho các đơn
// chưa tất toán CŨ NHẤT trước, TRONG 1 transaction atomic:
//  1. Khoá + đọc đơn chưa tất toán của khách (FOR UPDATE, oldest-first).
//  2. finance.RecordPaymentIn (ref_type='customer', idempotent) → 1 phiếu tổng.
//  3. Phân bổ tuần tự: mỗi đơn nhận min(còn-thiếu-của-đơn, còn-lại-của-phiếu) →
//     InsertAllocation + tính lại payment_status đơn. Dừng khi hết tiền.
//  4. Post PAYMENT_IN (Dr 111/112 / Cr 131) cho TỔNG tiền phiếu.
//  Phần dư (thu thừa) trả ở Leftover — KHÔNG gắn đơn (credit khách).
func (o *Orchestrator) RecordLumpSumPayment(ctx context.Context, in LumpSumInput) (LumpSumResult, error) {
	if o.txm == nil || o.orchFromTx == nil || o.recorder == nil || o.poster == nil {
		return LumpSumResult{}, apperr.Internal("Orchestrator chưa cấu hình đủ cho RecordLumpSumPayment")
	}
	if in.CustomerID <= 0 {
		return LumpSumResult{}, apperr.Validation("customer_id không hợp lệ (phải > 0)")
	}
	if !in.Amount.IsPositive() {
		return LumpSumResult{}, apperr.Validation("số tiền thu phải lớn hơn 0")
	}
	bookType := in.BookType
	if !bookType.Valid() {
		bookType = financedomain.BookBoth
	}

	var out LumpSumResult
	err := o.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := o.orchFromTx(tx)

		// (1) Đơn chưa tất toán của khách, cũ nhất trước (FOR UPDATE).
		unpaid, uerr := store.ListUnpaidOrdersByCustomer(ctx, in.CustomerID)
		if uerr != nil {
			return apperr.Internal("đọc đơn chưa tất toán lỗi").WithCause(uerr)
		}

		// (2) Ghi phiếu lump-sum (ref_type='customer'). Idempotent qua finance.
		custRef := strconv.FormatInt(in.CustomerID, 10)
		payment, created, rerr := o.recorder.RecordPaymentIn(ctx, tx, financeapp.RecordPaymentInParams{
			RecordPaymentIn: financedomain.RecordPaymentIn{
				Amount:        in.Amount,
				FundAccountID: in.FundAccountID,
				BookType:      bookType,
				BankRef:       in.BankRef,
				CreatedBy:     in.CreatedBy,
				InitialStatus: paymentInitialStatus(in.Collected),
				RefType:       financedomain.RefTypeCustomer,
				RefID:         custRef,
			},
			IdemKey: in.IdemKey,
		})
		if rerr != nil {
			return rerr
		}
		// Replay cùng Idempotency-Key: phiếu đã ghi + đã phân bổ + post sổ ở lần trước
		// → KHÔNG phân bổ/post LẠI (chống nhân đôi "đã thu" + bút toán). Trả phiếu cũ.
		if !created {
			out = LumpSumResult{Payment: payment, Allocations: nil, Leftover: money.Zero()}
			return nil
		}

		// (3) Phân bổ tuần tự oldest-first.
		remaining := in.Amount
		allocs := make([]AllocationLine, 0, len(unpaid))
		for _, ord := range unpaid {
			if !remaining.IsPositive() {
				break
			}
			orderRemaining := ord.Final.Sub(ord.Paid)
			if !orderRemaining.IsPositive() {
				continue // đơn đã đủ (phòng thủ) → bỏ qua
			}
			take := orderRemaining
			if orderRemaining.Sub(remaining).IsPositive() { // orderRemaining > remaining → lấy phần còn lại
				take = remaining
			}
			if aerr := store.InsertAllocation(ctx, payment.ID, ord.Code, take); aerr != nil {
				return apperr.Internal("ghi phân bổ phiếu lỗi").WithCause(aerr)
			}
			// Tính lại đã-thu (gồm allocation vừa ghi) → payment_status đơn.
			paid, serr := store.SumPaidInTx(ctx, ord.Code)
			if serr != nil {
				return apperr.Internal("tính lại đã thu (phân bổ) lỗi").WithCause(serr)
			}
			ps := domain.DerivePaymentStatus(ord.Final, paid)
			if _, perr := store.UpdatePaymentStatus(ctx, ord.ID, ps.String()); perr != nil {
				return apperr.Internal("cập nhật payment_status (phân bổ) lỗi").WithCause(perr)
			}
			allocs = append(allocs, AllocationLine{OrderID: ord.ID, OrderCode: ord.Code, Amount: take})
			remaining = remaining.Sub(take)
		}

		// (4) Post PAYMENT_IN cho TỔNG tiền phiếu (Dr 111/112 / Cr 131).
		if perr := o.postPaymentEntries(ctx, tx, store, bookType, in.FundIsBank, in.Amount, payment.Code); perr != nil {
			return perr
		}

		out = LumpSumResult{Payment: payment, Allocations: allocs, Leftover: remaining}
		return nil
	})
	if err != nil {
		return LumpSumResult{}, err
	}
	return out, nil
}

// ---- helpers ----

// realInternalVAT tính VAT sổ INTERNAL từ GIÁ THỰC từng dòng (LineTotal sau chiết
// khấu) × thuế suất từng dòng, làm tròn RoundVND MỖI DÒNG (khớp cách HĐ tính). ĐỘC
// LẬP giá HĐ override (sổ TAX): sổ INTERNAL luôn nhất quán doanh-thu-thực ↔ VAT-thực
// (sửa lệch khi FE override giá HĐ ≠ giá thực — dual-ledger). Khi KHÔNG override +
// cùng thuế suất, kết quả == VAT của HĐ. lineAmounts/rates song song theo thứ tự dòng.
func realInternalVAT(lineAmounts []money.Money, vls []VATLine) money.Money {
	sum := money.Zero()
	for i, amt := range lineAmounts {
		if i >= len(vls) {
			break
		}
		sum = sum.Add(amt.Mul(vls[i].VATRate).RoundVND())
	}
	return sum
}

// createPosOrder tạo đơn POS (PENDING) qua OrderStore (P4a) trong tx hiện hành,
// idempotent theo IdemKey. Tái dùng sinh mã + insert của P4a (KHÔNG lặp logic).
// Trả (CreatedOrder, created, error). created=false khi idem hit (đơn đã tồn tại)
// → CreatePosSale PHẢI dừng (sale đã hoàn tất trước đó), KHÔNG trừ kho/thu/post lại
// (chống nhân đôi khi replay cùng Idempotency-Key).
func (o *Orchestrator) createPosOrder(ctx context.Context, store OrderStore, draft domain.Draft, idemKey string) (CreatedOrder, bool, error) {
	idemKey = strings.TrimSpace(idemKey)
	if idemKey != "" {
		if existingID, found, ferr := store.FindByIdemKey(ctx, idemKey); ferr != nil {
			return CreatedOrder{}, false, apperr.Internal("tra cứu idempotency POS lỗi").WithCause(ferr)
		} else if found {
			got, gerr := store.GetCreated(ctx, existingID)
			return got, false, gerr // idem hit → KHÔNG tạo mới
		}
	}
	created, cerr := insertOrderViaStore(ctx, store, draft)
	if cerr != nil {
		return CreatedOrder{}, false, cerr
	}
	if idemKey != "" {
		inserted, berr := store.BindIdemKey(ctx, idemKey, created.Order.ID, created.Order.Code)
		if berr != nil {
			return CreatedOrder{}, false, apperr.Internal("ghi idempotency POS lỗi").WithCause(berr)
		}
		if !inserted {
			// Đua key (rất hiếm trong cùng tx POS) → trả lỗi để rollback; caller retry.
			return CreatedOrder{}, false, errIdemRace
		}
	}
	return created, true, nil // VỪA tạo mới
}

// postEntries post lần lượt mỗi JournalEntry qua Poster trong tx (gộp atomic). Lỗi
// post (kỳ khoá / TK không cho hạch toán / lệch Σ) trả nguyên (đã là apperr) → ROLLBACK.
func (o *Orchestrator) postEntries(ctx context.Context, tx pgx.Tx, _ OrchestrationStore, entries []accountingdomain.JournalEntry) error {
	for _, e := range entries {
		if _, perr := o.poster.Post(ctx, tx, e); perr != nil {
			return perr
		}
	}
	return nil
}

// postPaymentEntries post bút toán PAYMENT_IN theo sổ của phiếu: INTERNAL/BOTH →
// post sổ INTERNAL; TAX/BOTH → post sổ TAX (BOTH post CẢ hai).
func (o *Orchestrator) postPaymentEntries(ctx context.Context, tx pgx.Tx, _ OrchestrationStore, bookType financedomain.BookType, fundIsBank bool, amount money.Money, code string) error {
	date := o.now()
	if bookType == financedomain.BookInternal || bookType == financedomain.BookBoth {
		e := o.rules.BuildPaymentIn(accountingdomain.BookInternal, fundIsBank, amount, code, date)
		if _, perr := o.poster.Post(ctx, tx, e); perr != nil {
			return perr
		}
	}
	if bookType == financedomain.BookTax || bookType == financedomain.BookBoth {
		e := o.rules.BuildPaymentIn(accountingdomain.BookTax, fundIsBank, amount, code, date)
		if _, perr := o.poster.Post(ctx, tx, e); perr != nil {
			return perr
		}
	}
	return nil
}

// buildInvoiceLines dựng dòng HĐ từ dòng đơn (OrchLine) + VATLines song song. Nếu
// VATLine có override đơn giá (HasPrice) → dùng giá HĐ; nếu không → đơn giá đơn.
// Thiếu VATLine cho dòng → VAT 0% (FE nên truyền đủ cho B2B).
func buildInvoiceLines(lines []OrchLine, vatLines []VATLine) []vatdomain.LineInput {
	out := make([]vatdomain.LineInput, 0, len(lines))
	for i, l := range lines {
		rate := decimal.Zero
		price := l.UnitPrice
		desc := ""
		if i < len(vatLines) {
			rate = vatLines[i].VATRate
			if vatLines[i].HasPrice {
				price = vatLines[i].UnitPrice
			}
			desc = vatLines[i].Description
		}
		out = append(out, vatdomain.LineInput{
			ProductID:   l.ProductID,
			Description: desc,
			Quantity:    money.FromDecimal(l.Quantity.Decimal()),
			UnitPrice:   price,
			VATRate:     rate,
		})
	}
	return out
}

// buildInvoiceLinesFromDraft như buildInvoiceLines nhưng từ ComputedLine (POS tạo
// đơn trong cùng tx — chưa đọc lại OrchLine).
func buildInvoiceLinesFromDraft(lines []domain.ComputedLine, vatLines []VATLine) []vatdomain.LineInput {
	out := make([]vatdomain.LineInput, 0, len(lines))
	for i, l := range lines {
		rate := decimal.Zero
		price := l.UnitPrice
		desc := ""
		if i < len(vatLines) {
			rate = vatLines[i].VATRate
			if vatLines[i].HasPrice {
				price = vatLines[i].UnitPrice
			}
			desc = vatLines[i].Description
		}
		out = append(out, vatdomain.LineInput{
			ProductID:   l.ProductID,
			Description: desc,
			Quantity:    money.FromDecimal(l.Quantity.Decimal()),
			UnitPrice:   price,
			VATRate:     rate,
		})
	}
	return out
}

// orchHeaderToOrder dựng domain.Order (đường trả về) từ header đã khoá + trạng
// thái/payment_status mới. Payment summary để zero (FE đọc lại nếu cần chi tiết).
func orchHeaderToOrder(h OrchHeader, status, paymentStatus string) domain.Order {
	return domain.Order{
		ID:            h.ID,
		Code:          h.Code,
		CustomerID:    h.CustomerID,
		Status:        status,
		OrderType:     h.OrderType,
		Total:         h.TotalAmount,
		Final:         h.FinalAmount,
		PaymentStatus: paymentStatus,
	}
}

// toInventoryQty chuyển orders/domain.Quantity → inventory/domain.Quantity (cùng
// decimal nền, khác kiểu Go vì hai bounded context).
func toInventoryQty(q domain.Quantity) inventorydomain.Quantity {
	return inventorydomain.QuantityFromDecimal(q.Decimal())
}

// paymentInitialStatus map cờ Collected → trạng thái khởi tạo phiếu: Collected=true
// (NV đã thu từ khách, chưa nộp quỹ) → 'pending'; false → 'completed' (thu thẳng vào
// quỹ). Thanh toán 2 bước, spec mục 55.
//
// ⚠️ TINH CHỈNH KẾ TOÁN (chờ kế toán duyệt — tt133-posting-rules): hiện bút toán
// PAYMENT_IN (Dr 111/112 / Cr 131) post NGAY lúc ghi phiếu cho cả pending. Bản
// "chuẩn sách" hơn dùng TK 113 (Tiền đang chuyển): thu (pending) = Dr 113/Cr 131,
// xác nhận (completed) = Dr 111/Cr 113. Số dư QUỸ (fund_accounts.balance) vẫn đúng
// vì trigger prod chỉ bắn ở 'completed'; chỉ thời điểm ghi nhận 111 ở SỔ là sớm
// hơn. Nâng cấp khi kế toán xác nhận có dùng TK 113 hay không.
func paymentInitialStatus(collected bool) string {
	if collected {
		return financedomain.StatusPending
	}
	return financedomain.StatusCompleted
}

// posPaymentIdemKey sinh idem key cho phiếu thu POS: ưu tiên IdemKey của đơn (nếu
// có) ghép code để 1 đơn 1 phiếu; rỗng → dùng code (đơn đã UNIQUE code → an toàn).
func posPaymentIdemKey(orderIdemKey, code string) string {
	k := strings.TrimSpace(orderIdemKey)
	if k != "" {
		return "pos:" + k + ":" + code
	}
	return "pos:" + code
}

// posPrevStatus trả trạng thái NGAY TRƯỚC step trong luồng tiến POS (để guard
// UpdateStatus). PENDING→CONFIRMED→SHIPPING→COMPLETED.
func posPrevStatus(step domain.Status) domain.Status {
	switch step {
	case domain.StatusConfirmed:
		return domain.StatusPending
	case domain.StatusShipping:
		return domain.StatusConfirmed
	case domain.StatusCompleted:
		return domain.StatusShipping
	default:
		return domain.StatusPending
	}
}
