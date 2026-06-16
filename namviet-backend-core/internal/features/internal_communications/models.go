package internal_communications

import (
	"time"
)

// InternalChannel represents the internal_channels table
type InternalChannel struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:text;not null" json:"name"`
	Type      string     `gorm:"type:text;default:'group';not null" json:"type"`
	CreatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}

// InternalMessage represents the internal_messages table
type InternalMessage struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChannelID int64      `gorm:"type:bigint;not null" json:"channel_id"`
	SenderID  string     `gorm:"type:uuid;not null" json:"sender_id"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	CreatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
