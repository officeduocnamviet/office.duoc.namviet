package system_configs

import "github.com/namviet/backend-core/internal/features/roles"

type CreateSystemConfigRequest struct {
	ConfigKey   string      `json:"config_key" binding:"required"`
	ConfigValue roles.JSONB `json:"config_value" binding:"required"`
	Description *string     `json:"description"`
	UpdatedBy   *string     `json:"updated_by"`
}

type UpdateSystemConfigRequest struct {
	ConfigValue *roles.JSONB `json:"config_value"`
	Description *string      `json:"description"`
	UpdatedBy   *string      `json:"updated_by"`
}
