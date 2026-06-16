package shipping_partners

import "github.com/namviet/backend-core/internal/features/roles"

type CreateShippingPartnerRequest struct {
	Code                string       `json:"code" binding:"required"`
	Name                string       `json:"name" binding:"required"`
	PartnerType         string       `json:"partner_type" binding:"required"`
	APIConfig           *roles.JSONB `json:"api_config"`
	TrackingURLTemplate *string      `json:"tracking_url_template"`
}

type UpdateShippingPartnerRequest struct {
	Code                *string      `json:"code"`
	Name                *string      `json:"name"`
	PartnerType         *string      `json:"partner_type"`
	APIConfig           *roles.JSONB `json:"api_config"`
	TrackingURLTemplate *string      `json:"tracking_url_template"`
	Status              *string      `json:"status"`
}
