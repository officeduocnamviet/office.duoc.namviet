package agent_workflows

import "github.com/namviet/backend-core/internal/features/roles"

type CreateAgentWorkflowRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description *string     `json:"description"`
	TriggerType string      `json:"trigger_type" binding:"required"`
	Steps       roles.JSONB `json:"steps" binding:"required"`
	IsActive    *bool       `json:"is_active"`
}

type UpdateAgentWorkflowRequest struct {
	Name        *string      `json:"name"`
	Description *string      `json:"description"`
	TriggerType *string      `json:"trigger_type"`
	Steps       *roles.JSONB `json:"steps"`
	IsActive    *bool        `json:"is_active"`
}
