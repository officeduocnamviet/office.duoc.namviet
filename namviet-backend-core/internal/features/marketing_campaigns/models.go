package marketing_campaigns

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// MarketingCampaign represents the marketing_campaigns table
type MarketingCampaign struct {
	ID            string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name          string       `gorm:"type:text;not null" json:"name"`
	Description   *string      `gorm:"type:text" json:"description,omitempty"`
	TargetSegment *roles.JSONB `gorm:"type:jsonb" json:"target_segment,omitempty"`
	Budget        float64      `gorm:"type:numeric;default:0;not null" json:"budget"`
	StartDate     time.Time    `gorm:"type:timestamp with time zone;not null" json:"start_date"`
	EndDate       *time.Time   `gorm:"type:timestamp with time zone" json:"end_date,omitempty"`
	Status        string       `gorm:"type:text;default:'draft';not null" json:"status"`
	CreatedAt     *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt     *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}
