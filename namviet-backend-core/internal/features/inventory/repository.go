package inventory

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetInventory(warehouseID int64, productID int64) ([]InventoryBatch, error) {
	var results []InventoryBatch
	db := supabase.DB
	query := db.Model(&InventoryBatch{})
	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}
	err := query.Find(&results).Error
	return results, err
}

func CreateTransactionRecord(tx *gorm.DB, data *InventoryTransaction) error {
	return tx.Create(data).Error
}

func UpdateInventoryBatchQuantity(tx *gorm.DB, warehouseID, productID, batchID int64, delta int) error {
	var invBatch InventoryBatch
	err := tx.First(&invBatch, "warehouse_id = ? AND product_id = ? AND batch_id = ?", warehouseID, productID, batchID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new if IN
		if delta < 0 {
			return errors.New("insufficient stock")
		}
		invBatch = InventoryBatch{
			WarehouseID: warehouseID,
			ProductID:   productID,
			BatchID:     batchID,
			Quantity:    delta,
		}
		return tx.Create(&invBatch).Error
	} else if err != nil {
		return err
	}

	newQuantity := invBatch.Quantity + delta
	if newQuantity < 0 {
		return errors.New("insufficient stock")
	}

	invBatch.Quantity = newQuantity
	return tx.Save(&invBatch).Error
}
