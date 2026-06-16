package audit_logs

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// SystemAuditLog represents the system_audit_logs table
type SystemAuditLog struct {
	ID        string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    *string      `gorm:"type:uuid" json:"user_id,omitempty"`
	Action    string       `gorm:"type:text;not null" json:"action"`
	TableName string       `gorm:"type:text;not null" json:"table_name"`
	RecordID  *string      `gorm:"type:text" json:"record_id,omitempty"`
	OldData   *roles.JSONB `gorm:"type:jsonb" json:"old_data,omitempty"`
	NewData   *roles.JSONB `gorm:"type:jsonb" json:"new_data,omitempty"`
	IPAddress *string      `gorm:"type:text" json:"ip_address,omitempty"`
	UserAgent *string      `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
