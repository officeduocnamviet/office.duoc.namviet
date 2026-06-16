package batches

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetBatchesByProductID(productID int64) ([]Batch, error) {
	var results []Batch
	db := supabase.DB
	err := db.Where("product_id = ? AND deleted_at IS NULL", productID).Find(&results).Error
	return results, err
}

func GetBatchByID(id int64) (*Batch, error) {
	var result Batch
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("batch not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateBatch(data *Batch) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateBatch(data *Batch) error {
	db := supabase.DB
	return db.Save(data).Error
}
