package ai_agent_memories

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// AIAgentMemory represents the ai_agent_memories table
type AIAgentMemory struct {
	ID        string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SessionID *string     `gorm:"type:uuid" json:"session_id,omitempty"`
	UserID    *string     `gorm:"type:uuid" json:"user_id,omitempty"`
	Key       string      `gorm:"type:text;not null" json:"key"`
	Value     roles.JSONB `gorm:"type:jsonb;not null" json:"value"`
	ExpiresAt *time.Time  `gorm:"type:timestamp with time zone" json:"expires_at,omitempty"`
	CreatedAt *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
