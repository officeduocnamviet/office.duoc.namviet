package users

import (
	"time"
)

// User represents the users table
type User struct {
	ID        string     `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string     `gorm:"type:text" json:"email"`
	FullName  string     `gorm:"type:text;column:full_name" json:"full_name"`
	Phone     string     `gorm:"type:text" json:"phone"`
	Status    string     `gorm:"type:text" json:"status"`
	RoleID    string     `gorm:"type:uuid;column:role_id" json:"role_id"`
	CompanyID string     `gorm:"type:uuid;column:company_id" json:"company_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

// UserFCMToken represents the user_fcm_tokens table
type UserFCMToken struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"type:uuid;column:user_id" json:"user_id"`
	Token      string    `gorm:"type:text;column:token;unique" json:"token"`
	DeviceInfo string    `gorm:"type:text;column:device_info" json:"device_info"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}
