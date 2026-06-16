package shipping_partners

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// ShippingPartner represents the shipping_partners table
type ShippingPartner struct {
	ID                  string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Code                string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name                string      `gorm:"type:varchar(255);not null" json:"name"`
	PartnerType         string      `gorm:"type:varchar(50);not null" json:"partner_type"`
	APIConfig           roles.JSONB `gorm:"type:jsonb" json:"api_config,omitempty"`
	TrackingURLTemplate *string     `gorm:"type:text" json:"tracking_url_template,omitempty"`
	Status              string      `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt           *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt           *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt           *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
