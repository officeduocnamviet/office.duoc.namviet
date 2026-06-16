package customers

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles" // For JSONB
)

// Customer represents the customers table
type Customer struct {
	ID             int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerCode   *string     `gorm:"type:text" json:"customer_code,omitempty"`
	Name           string      `gorm:"type:text;not null" json:"name"`
	CustomerType   string      `gorm:"type:text;default:'B2C'" json:"customer_type"`
	Phone          *string     `gorm:"type:text" json:"phone,omitempty"`
	Email          *string     `gorm:"type:text" json:"email,omitempty"`
	Address        *string     `gorm:"type:text" json:"address,omitempty"`
	Status         string      `gorm:"type:text;default:'active'" json:"status"`
	DOB            *string     `gorm:"type:date" json:"dob,omitempty"`
	Gender         *string     `gorm:"type:text" json:"gender,omitempty"`
	CCCD           *string     `gorm:"type:text" json:"cccd,omitempty"`
	LoyaltyPoints  *int        `gorm:"type:integer;default:0" json:"loyalty_points,omitempty"`
	B2BMetadata    roles.JSONB `gorm:"type:jsonb;default:'{}'" json:"b2b_metadata,omitempty"`
	CurrentDebt    *float64    `gorm:"type:numeric;default:0" json:"current_debt,omitempty"`
	UpdatedBy      *string     `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt      *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt      *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt      *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
