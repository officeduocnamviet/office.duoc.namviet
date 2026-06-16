package training_courses

import (
	"time"

	"gorm.io/gorm"
)

// TrainingCourse represents the training_courses table
type TrainingCourse struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string         `gorm:"type:text;not null" json:"title"`
	ContentType  string         `gorm:"type:text;not null" json:"content_type"`
	ContentURL   *string        `gorm:"type:text" json:"content_url,omitempty"`
	PassingScore *int           `gorm:"type:integer" json:"passing_score,omitempty"`
	Status       string         `gorm:"type:text;default:'active';not null" json:"status"`
	CreatedAt    *time.Time     `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
