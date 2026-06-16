package warehouses

import "time"

// Warehouse represents the warehouses table
type Warehouse struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:text;not null" json:"name"`
	Type      string     `gorm:"type:text;default:'main'" json:"type"`
	Address   *string    `gorm:"type:text" json:"address,omitempty"`
	Manager   *string    `gorm:"type:text" json:"manager,omitempty"`
	Status    string     `gorm:"type:text;default:'active'" json:"status"`
	CreatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
