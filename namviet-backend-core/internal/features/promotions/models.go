package promotions

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles" // For JSONB
)

// Promotion represents the promotions table
type Promotion struct {
	ID        string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Code      string      `gorm:"type:text;not null" json:"code"`
	Name      string      `gorm:"type:text;not null" json:"name"`
	Rules     roles.JSONB `gorm:"type:jsonb;not null" json:"rules"`
	StartDate time.Time   `gorm:"type:timestamp with time zone;not null" json:"start_date"`
	EndDate   time.Time   `gorm:"type:timestamp with time zone;not null" json:"end_date"`
	Status    *string     `gorm:"type:text;default:'active'" json:"status,omitempty"`
	CreatedAt *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
