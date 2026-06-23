package app

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/orders/domain"
)

// WriteService là use-case GHI của orders (P4a): tạo đơn + đổi trạng thái ĐƠN
// GIẢN. Mỗi use-case = 1 transaction (platform/db.WithinTx qua TxManager). KHÔNG
// đụng kho/tiền/sổ (ShipOrder/RecordPayment/POS = P4b). storeFromTx bind OrderStore
// tới tx; txm mở/commit tx. nil component → use-case tương ứng trả Internal (test).
type WriteService struct {
	storeFromTx OrderStoreFromTx
	txm         TxManager
}

// NewWrite dựng WriteService.
func NewWrite(storeFromTx OrderStoreFromTx, txm TxManager) *WriteService {
	return &WriteService{storeFromTx: storeFromTx, txm: txm}
}

// CreateOrderInput là input tạo đơn đã giải mã ở edge (tiền decimal — KHÔNG float).
// IdemKey là Idempotency-Key (header) — rỗng nghĩa là không idempotent (mỗi gọi
// tạo đơn mới). Lines map sang domain.DraftLine.
type CreateOrderInput struct {
	CustomerID *int64
	OrderType  string
	CreatorID  string
	Note       string
	Lines      []domain.DraftLine
	IdemKey    string
}

// CreateOrder tạo MỘT đơn (status PENDING) + các dòng hàng trong 1 transaction.
// Quy trình:
//  1. Dựng + validate + tính tiền THUẦN ở domain (NewDraft) → lỗi = Validation/422.
//  2. (Nếu có IdemKey) tra app.order_idempotency → đã tạo thì TRẢ đơn cũ (no-op).
//  3. Sinh code (sequence + tiền tố) + id (uuid v7); INSERT orders + order_items.
//  4. (Nếu có IdemKey) BindIdemKey; nếu key đã bị luồng khác chiếm (đua) → đọc lại
//     đơn của luồng thắng (idempotent). INSERT trùng code (đua sinh mã) → thử lại.
//
// Tiền decimal toàn tuyến. Trả CreatedOrder (đơn + lines).
func (s *WriteService) CreateOrder(ctx context.Context, in CreateOrderInput) (CreatedOrder, error) {
	if s.txm == nil || s.storeFromTx == nil {
		return CreatedOrder{}, apperr.Internal("WriteService chưa cấu hình đủ cho CreateOrder")
	}

	// (1) Dựng draft THUẦN (validate + tính tiền). Lỗi domain → Validation.
	draft, err := domain.NewDraft(domain.DraftInput{
		CustomerID: in.CustomerID,
		OrderType:  domain.OrderType(strings.TrimSpace(in.OrderType)),
		CreatorID:  strings.TrimSpace(in.CreatorID),
		Note:       in.Note,
		Lines:      in.Lines,
	})
	if err != nil {
		return CreatedOrder{}, apperr.Validation(err.Error())
	}

	idemKey := strings.TrimSpace(in.IdemKey)
	var out CreatedOrder
	txErr := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)

		// (2) Idempotency: đã tạo theo key → trả đơn cũ (no-op).
		if idemKey != "" {
			if existingID, found, ferr := store.FindByIdemKey(ctx, idemKey); ferr != nil {
				return apperr.Internal("tra cứu idempotency tạo đơn lỗi").WithCause(ferr)
			} else if found {
				got, gerr := store.GetCreated(ctx, existingID)
				if gerr != nil {
					return apperr.Internal("đọc lại đơn idempotent lỗi").WithCause(gerr)
				}
				out = got
				return nil
			}
		}

		// (3) Sinh mã + id, INSERT đơn + dòng.
		created, cerr := s.insertNewOrder(ctx, store, draft)
		if cerr != nil {
			return cerr
		}

		// (4) Bind idempotency key TRONG cùng tx. Nếu key đã bị luồng khác commit
		// (đua cùng key đồng thời, rất hiếm) → inserted=false: trả errIdemRace để
		// WithinTx ROLLBACK (huỷ đơn thừa này) → caller readback đơn của luồng thắng
		// trong tx mới. Bảo đảm 1 key ⇒ đúng 1 đơn.
		if idemKey != "" {
			inserted, berr := store.BindIdemKey(ctx, idemKey, created.Order.ID, created.Order.Code)
			if berr != nil {
				return apperr.Internal("ghi idempotency tạo đơn lỗi").WithCause(berr)
			}
			if !inserted {
				return errIdemRace
			}
		}
		out = created
		return nil
	})
	if txErr != nil {
		if txErr == errIdemRace {
			// Tx vừa rollback (đơn thừa bị huỷ). Đọc lại đơn của luồng thắng.
			return s.createIdempotentReadback(ctx, idemKey)
		}
		return CreatedOrder{}, txErr
	}
	return out, nil
}

// errIdemRace báo hiệu đua idempotency-key (luồng khác thắng) để rollback + readback.
var errIdemRace = apperr.Conflict("đua idempotency key — đọc lại đơn của luồng thắng")

// createIdempotentReadback đọc lại đơn của luồng thắng (sau khi rollback đơn thừa).
func (s *WriteService) createIdempotentReadback(ctx context.Context, idemKey string) (CreatedOrder, error) {
	var out CreatedOrder
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
			return apperr.Internal("đọc lại đơn idempotent (readback) lỗi").WithCause(gerr)
		}
		out = got
		return nil
	})
	if err != nil {
		return CreatedOrder{}, err
	}
	return out, nil
}

// insertNewOrder sinh code + id, INSERT header + lines. Trùng code (đua sinh mã,
// rất hiếm) → thử cấp số mới (tối đa vài lần) để không fail oan. Delegate sang
// helper package-level insertOrderViaStore (tái dùng cho POS — P4b).
func (s *WriteService) insertNewOrder(ctx context.Context, store OrderStore, draft domain.Draft) (CreatedOrder, error) {
	return insertOrderViaStore(ctx, store, draft)
}

// insertOrderViaStore sinh code + id, INSERT header + lines qua OrderStore (bound
// tx). Dùng chung cho CreateOrder (P4a) và CreatePosSale (P4b) để KHÔNG lặp logic
// sinh mã/insert. Trùng code (đua sinh mã, rất hiếm) → cấp số mới thử lại tối đa
// vài lần. Trả CreatedOrder (đơn + lines đọc lại).
func insertOrderViaStore(ctx context.Context, store OrderStore, draft domain.Draft) (CreatedOrder, error) {
	const maxCodeRetry = 5
	for attempt := 0; attempt < maxCodeRetry; attempt++ {
		seq, serr := store.NextCodeSeq(ctx)
		if serr != nil {
			return CreatedOrder{}, apperr.Internal("cấp số mã đơn lỗi").WithCause(serr)
		}
		code := formatOrderCode(seq)
		orderID := id.NewString()
		saved, duplicate, ierr := store.InsertOrder(ctx, NewOrderRow{ID: orderID, Code: code, Draft: draft})
		if ierr != nil {
			return CreatedOrder{}, apperr.Internal("ghi đơn lỗi").WithCause(ierr)
		}
		if duplicate {
			continue // mã trùng (đua) → cấp số mới, thử lại
		}
		for _, l := range draft.Lines {
			if lerr := store.InsertItem(ctx, saved.ID, l); lerr != nil {
				return CreatedOrder{}, apperr.Internal("ghi dòng hàng lỗi").WithCause(lerr)
			}
		}
		got, gerr := store.GetCreated(ctx, saved.ID)
		if gerr != nil {
			return CreatedOrder{}, apperr.Internal("đọc lại đơn vừa tạo lỗi").WithCause(gerr)
		}
		return got, nil
	}
	return CreatedOrder{}, apperr.Internal("không cấp được mã đơn duy nhất sau nhiều lần thử")
}

// ConfirmOrder duyệt đơn: PENDING→CONFIRMED. ConfirmOrder/CompleteOrder/CancelOrder
// dùng chung transition (load FOR UPDATE → CanTransition → UPDATE guard).
func (s *WriteService) ConfirmOrder(ctx context.Context, orderID string) (domain.Order, error) {
	return s.transition(ctx, orderID, domain.StatusConfirmed)
}

// CompleteOrder hoàn tất đơn: SHIPPING→COMPLETED. (Chuyển vào SHIPPING cần trừ kho
// → P4b; ở P4a chỉ chuyển COMPLETED khi đơn đã ở SHIPPING.)
func (s *WriteService) CompleteOrder(ctx context.Context, orderID string) (domain.Order, error) {
	return s.transition(ctx, orderID, domain.StatusCompleted)
}

// CancelOrder huỷ đơn: PENDING/CONFIRMED→CANCELLED. (Huỷ đơn đã trừ kho/post sổ
// cần đảo bút toán + hoàn kho → P4b.)
func (s *WriteService) CancelOrder(ctx context.Context, orderID string) (domain.Order, error) {
	return s.transition(ctx, orderID, domain.StatusCancelled)
}

// transition là khung chung đổi trạng thái trong 1 tx: load + FOR UPDATE → kiểm
// CanTransition (state machine THUẦN) → UPDATE guard status cũ. Sai trạng thái /
// không cho chuyển → Conflict. Không thấy đơn → NotFound. Trả đơn đã cập nhật.
func (s *WriteService) transition(ctx context.Context, orderID string, to domain.Status) (domain.Order, error) {
	if s.txm == nil || s.storeFromTx == nil {
		return domain.Order{}, apperr.Internal("WriteService chưa cấu hình đủ cho đổi trạng thái")
	}
	if strings.TrimSpace(orderID) == "" {
		return domain.Order{}, apperr.Validation("id đơn không được rỗng")
	}
	var out domain.Order
	err := s.txm.WithinTx(ctx, func(tx pgx.Tx) error {
		store := s.storeFromTx(tx)
		current, found, gerr := store.GetForUpdate(ctx, orderID)
		if gerr != nil {
			return apperr.Internal("đọc trạng thái đơn lỗi").WithCause(gerr)
		}
		if !found {
			return apperr.NotFound("đơn hàng không tồn tại")
		}
		if !domain.CanTransition(current, to) {
			return apperr.Conflict("không thể chuyển đơn từ " + current.String() + " sang " + to.String())
		}
		rows, uerr := store.UpdateStatus(ctx, orderID, current, to)
		if uerr != nil {
			return apperr.Internal("cập nhật trạng thái đơn lỗi").WithCause(uerr)
		}
		if rows == 0 {
			// Đơn đã đổi trạng thái giữa lúc đọc và ghi (dù đã FOR UPDATE — phòng thủ).
			return apperr.Conflict("đơn đã đổi trạng thái bởi thao tác khác")
		}
		got, derr := store.GetCreated(ctx, orderID)
		if derr != nil {
			return apperr.Internal("đọc lại đơn sau cập nhật lỗi").WithCause(derr)
		}
		out = got.Order
		return nil
	})
	if err != nil {
		return domain.Order{}, err
	}
	return out, nil
}
