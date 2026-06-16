package marketing_campaigns

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreateMarketingCampaignRequest struct {
	Name          string       `json:"name" binding:"required"`
	Description   *string      `json:"description"`
	TargetSegment *roles.JSONB `json:"target_segment"`
	Budget        float64      `json:"budget"`
	StartDate     time.Time    `json:"start_date" binding:"required"`
	EndDate       *time.Time   `json:"end_date"`
}

type UpdateMarketingCampaignRequest struct {
	Name          *string      `json:"name"`
	Description   *string      `json:"description"`
	TargetSegment *roles.JSONB `json:"target_segment"`
	Budget        *float64     `json:"budget"`
	StartDate     *time.Time   `json:"start_date"`
	EndDate       *time.Time   `json:"end_date"`
	Status        *string      `json:"status"`
}
