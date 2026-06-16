package agent_workflows

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// AgentWorkflow represents the agent_workflows table
type AgentWorkflow struct {
	ID          string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string       `gorm:"type:text;not null" json:"name"`
	Description *string      `gorm:"type:text" json:"description,omitempty"`
	TriggerType string       `gorm:"type:text;not null" json:"trigger_type"`
	Steps       roles.JSONB  `gorm:"type:jsonb;not null" json:"steps"`
	IsActive    bool         `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt   *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}
