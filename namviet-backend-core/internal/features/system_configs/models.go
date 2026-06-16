package system_configs

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// SystemConfig represents the system_configs table
type SystemConfig struct {
	ID          string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ConfigKey   string      `gorm:"type:text;uniqueIndex;not null" json:"config_key"`
	ConfigValue roles.JSONB `gorm:"type:jsonb;not null" json:"config_value"`
	Description *string     `gorm:"type:text" json:"description,omitempty"`
	UpdatedAt   *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	UpdatedBy   *string     `gorm:"type:uuid" json:"updated_by,omitempty"`
}
