package inventory

import (
	"errors"
	"github.com/namviet/backend-core/internal/platform/supabase"
)

func GetInventoryService(warehouseID int64, productID int64) ([]InventoryBatch, error) {
	return GetInventory(warehouseID, productID)
}

func CreateTransactionService(req CreateTransactionRequest) (*InventoryTransaction, error) {
	validTypes := map[string]bool{"inbound": true, "outbound": true, "transfer": true, "adjustment": true, "return": true}
	if !validTypes[req.Type] {
		return nil, errors.New("invalid transaction type, must be inbound, outbound, transfer, adjustment, or return")
	}

	db := supabase.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	delta := req.Quantity
	if req.Type == "outbound" {
		delta = -req.Quantity
	}

	transaction := &InventoryTransaction{
		WarehouseID:   req.WarehouseID,
		ProductID:     req.ProductID,
		BatchID:       req.BatchID,
		Type:          req.Type,
		ActionGroup:   req.ActionGroup,
		Quantity:      req.Quantity,
		UnitPrice:     req.UnitPrice,
		RefID:         req.RefID,
		Description:   req.Description,
		PartnerID:     req.PartnerID,
	}

	if err := CreateTransactionRecord(tx, transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Only update batch quantity if batch_id is provided
	if req.BatchID != nil {
		if err := UpdateInventoryBatchQuantity(tx, req.WarehouseID, req.ProductID, *req.BatchID, delta); err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// For NamViet system, batches might be mandatory for inventory, but we leave it optional here
		// depending on business rules.
		// Usually if OUT, we must specify batch, if IN, we might create a batch first.
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return transaction, nil
}
