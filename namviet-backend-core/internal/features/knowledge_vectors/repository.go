package knowledge_vectors

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Medical Knowledge Vectors
func GetAllMedicalKnowledgeVectors() ([]MedicalKnowledgeVector, error) {
	var results []MedicalKnowledgeVector
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetMedicalKnowledgeVectorByID(id string) (*MedicalKnowledgeVector, error) {
	var result MedicalKnowledgeVector
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medical knowledge vector not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateMedicalKnowledgeVector(data *MedicalKnowledgeVector) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateMedicalKnowledgeVector(data *MedicalKnowledgeVector) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteMedicalKnowledgeVector(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&MedicalKnowledgeVector{}).Error
}

// Product Vectors
func GetAllProductVectors() ([]ProductVector, error) {
	var results []ProductVector
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetProductVectorByID(id string) (*ProductVector, error) {
	var result ProductVector
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product vector not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateProductVector(data *ProductVector) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateProductVector(data *ProductVector) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteProductVector(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&ProductVector{}).Error
}
