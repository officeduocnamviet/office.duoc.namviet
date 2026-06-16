package categories

import "time"

// Category represents the categories table
type Category struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:text;not null" json:"name"`
	Slug      string     `gorm:"type:text;not null" json:"slug"`
	ParentID  *int64     `gorm:"type:bigint" json:"parent_id,omitempty"`
	Status    string     `gorm:"type:text;default:'active'" json:"status"`
	CreatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
