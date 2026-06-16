package medical_visits

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllMedicalVisits() ([]MedicalVisit, error) {
	var results []MedicalVisit
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetMedicalVisitByID(id string) (*MedicalVisit, error) {
	var result MedicalVisit
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medical visit not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateMedicalVisit(data *MedicalVisit) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateMedicalVisit(data *MedicalVisit) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteMedicalVisit(id string) error {
	db := supabase.DB
	return db.Model(&MedicalVisit{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
