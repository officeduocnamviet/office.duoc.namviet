package chats

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// ChatSession represents the chat_sessions table
type ChatSession struct {
	ID         string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CustomerID *string    `gorm:"type:uuid" json:"customer_id,omitempty"`
	AgentID    *string    `gorm:"type:uuid" json:"agent_id,omitempty"`
	Status     string     `gorm:"type:text;default:'active'" json:"status"`
	CreatedAt  *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	EndedAt    *time.Time `gorm:"type:timestamp with time zone" json:"ended_at,omitempty"`
}

// ChatMessage represents the chat_messages table
type ChatMessage struct {
	ID          string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SessionID   *string      `gorm:"type:uuid" json:"session_id,omitempty"`
	SenderType  string       `gorm:"type:text;not null" json:"sender_type"`
	SenderID    *string      `gorm:"type:uuid" json:"sender_id,omitempty"`
	Content     string       `gorm:"type:text;not null" json:"content"`
	MessageType string       `gorm:"type:text;default:'text'" json:"message_type"`
	Metadata    *roles.JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt   *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
