package ai_agent_memories

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreateAIAgentMemoryRequest struct {
	SessionID *string      `json:"session_id"`
	UserID    *string      `json:"user_id"`
	Key       string       `json:"key" binding:"required"`
	Value     roles.JSONB  `json:"value" binding:"required"`
	ExpiresAt *time.Time   `json:"expires_at"`
}

type UpdateAIAgentMemoryRequest struct {
	SessionID *string      `json:"session_id"`
	UserID    *string      `json:"user_id"`
	Key       *string      `json:"key"`
	Value     *roles.JSONB `json:"value"`
	ExpiresAt *time.Time   `json:"expires_at"`
}
