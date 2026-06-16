package roles

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// GetAllRoles fetches all roles from the database
func GetAllRoles() ([]Role, error) {
	var roles []Role
	db := supabase.DB
	err := db.Find(&roles).Error
	return roles, err
}

// GetRoleByID fetches a single role by ID
func GetRoleByID(id string) (*Role, error) {
	var role Role
	db := supabase.DB
	err := db.First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// CreateRole inserts a new role into the database
func CreateRole(role *Role) error {
	db := supabase.DB
	return db.Create(role).Error
}

// UpdateRole updates an existing role in the database
func UpdateRole(role *Role) error {
	db := supabase.DB
	return db.Save(role).Error
}

// DeleteRole removes a role from the database
func DeleteRole(id string) error {
	db := supabase.DB
	return db.Delete(&Role{}, "id = ?", id).Error
}
