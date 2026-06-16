package integrations

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// ThirdPartyConnection represents the third_party_connections table
type ThirdPartyConnection struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PartnerName string     `gorm:"type:text;not null" json:"partner_name"`
	APIKey      *string    `gorm:"type:text" json:"api_key,omitempty"`
	SecretKey   *string    `gorm:"type:text" json:"secret_key,omitempty"`
	WebhookURL  *string    `gorm:"type:text" json:"webhook_url,omitempty"`
	Status      string     `gorm:"type:text;default:'active'" json:"status"`
	CreatedAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}

// WebhookLog represents the webhook_logs table
type WebhookLog struct {
	ID             string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PartnerID      *string      `gorm:"type:uuid" json:"partner_id,omitempty"`
	EventType      string       `gorm:"type:text;not null" json:"event_type"`
	Payload        *roles.JSONB `gorm:"type:jsonb" json:"payload,omitempty"`
	ResponseStatus *int         `gorm:"type:integer" json:"response_status,omitempty"`
	ResponseBody   *string      `gorm:"type:text" json:"response_body,omitempty"`
	CreatedAt      *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
