package integrations

import "github.com/namviet/backend-core/internal/features/roles"

// Connection DTOs
type CreateConnectionRequest struct {
	PartnerName string  `json:"partner_name" binding:"required"`
	APIKey      *string `json:"api_key"`
	SecretKey   *string `json:"secret_key"`
	WebhookURL  *string `json:"webhook_url"`
}

type UpdateConnectionRequest struct {
	PartnerName *string `json:"partner_name"`
	APIKey      *string `json:"api_key"`
	SecretKey   *string `json:"secret_key"`
	WebhookURL  *string `json:"webhook_url"`
	Status      *string `json:"status"`
}

// WebhookLog DTOs
type CreateWebhookLogRequest struct {
	PartnerID      *string      `json:"partner_id"`
	EventType      string       `json:"event_type" binding:"required"`
	Payload        *roles.JSONB `json:"payload"`
	ResponseStatus *int         `json:"response_status"`
	ResponseBody   *string      `json:"response_body"`
}
