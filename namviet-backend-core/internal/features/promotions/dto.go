package promotions

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreatePromotionRequest struct {
	Code      string      `json:"code" binding:"required"`
	Name      string      `json:"name" binding:"required"`
	Rules     roles.JSONB `json:"rules" binding:"required"`
	StartDate time.Time   `json:"start_date" binding:"required"`
	EndDate   time.Time   `json:"end_date" binding:"required"`
}

type UpdatePromotionRequest struct {
	Name      *string      `json:"name"`
	Rules     *roles.JSONB `json:"rules"`
	StartDate *time.Time   `json:"start_date"`
	EndDate   *time.Time   `json:"end_date"`
	Status    *string      `json:"status"`
}
