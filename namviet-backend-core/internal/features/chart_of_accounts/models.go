package chart_of_accounts

import (
	"time"
)

// ChartOfAccount represents the chart_of_accounts table
type ChartOfAccount struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AccountCode  string     `gorm:"type:text;not null" json:"account_code"`
	Name         string     `gorm:"type:text;not null" json:"name"`
	ParentID     *string    `gorm:"type:uuid" json:"parent_id,omitempty"`
	Type         string     `gorm:"type:text;not null" json:"type"`
	BalanceType  string     `gorm:"type:text;not null" json:"balance_type"`
	Status       string     `gorm:"type:text;default:'active'" json:"status"`
	AllowPosting bool       `gorm:"type:boolean;default:true" json:"allow_posting"`
	CreatedAt    *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt    *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
