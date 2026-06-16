package categories

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// GetAllCategories fetches all categories
func GetAllCategories() ([]Category, error) {
	var categories []Category
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&categories).Error
	return categories, err
}

// GetCategoryByID fetches a single category
func GetCategoryByID(id int64) (*Category, error) {
	var category Category
	db := supabase.DB
	err := db.First(&category, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// CreateCategory inserts a new category
func CreateCategory(category *Category) error {
	db := supabase.DB
	return db.Create(category).Error
}

// UpdateCategory updates an existing category
func UpdateCategory(category *Category) error {
	db := supabase.DB
	return db.Save(category).Error
}
