package clinical_queues

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllClinicalQueues() ([]ClinicalQueue, error) {
	var results []ClinicalQueue
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetClinicalQueueByID(id string) (*ClinicalQueue, error) {
	var result ClinicalQueue
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("clinical queue not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateClinicalQueue(data *ClinicalQueue) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateClinicalQueue(data *ClinicalQueue) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteClinicalQueue(id string) error {
	db := supabase.DB
	return db.Model(&ClinicalQueue{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
