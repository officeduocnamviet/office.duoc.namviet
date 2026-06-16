package inventory

import (
	"errors"
	"github.com/namviet/backend-core/internal/platform/supabase"
)

func GetInventoryService(warehouseID int64, productID int64) ([]InventoryBatch, error) {
	return GetInventory(warehouseID, productID)
}

func CreateTransactionService(req CreateTransactionRequest) (*InventoryTransaction, error) {
	if req.Type != "IN" && req.Type != "OUT" {
		return nil, errors.New("invalid transaction type, must be IN or OUT")
	}

	db := supabase.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	delta := req.Quantity
	if req.Type == "OUT" {
		delta = -req.Quantity
	}

	transaction := &InventoryTransaction{
		WarehouseID:   req.WarehouseID,
		ProductID:     req.ProductID,
		BatchID:       req.BatchID,
		Type:          req.Type,
		Quantity:      req.Quantity,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
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
