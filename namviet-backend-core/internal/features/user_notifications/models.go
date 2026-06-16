package user_notifications

import (
	"time"
)

// UserFCMToken represents the user_fcm_tokens table
type UserFCMToken struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     string     `gorm:"type:uuid;not null" json:"user_id"`
	FCMToken   string     `gorm:"type:text;not null" json:"fcm_token"`
	DeviceID   *string    `gorm:"type:text" json:"device_id,omitempty"`
	DeviceType *string    `gorm:"type:text" json:"device_type,omitempty"`
	CreatedAt  *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt  *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}

// UserSocialMapping represents the user_social_mappings table
type UserSocialMapping struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         string     `gorm:"type:uuid;not null" json:"user_id"`
	SocialProvider string     `gorm:"type:text;not null" json:"social_provider"`
	SocialID       string     `gorm:"type:text;not null" json:"social_id"`
	SocialName     *string    `gorm:"type:text" json:"social_name,omitempty"`
	SocialAvatar   *string    `gorm:"type:text" json:"social_avatar,omitempty"`
	CreatedAt      *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
