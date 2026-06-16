package roles

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB slice for GORM
type JSONB []string

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		// Sometimes postgres driver returns string
		str, ok := value.(string)
		if !ok {
			return errors.New("type assertion to []byte or string failed")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, &j)
}

// Role represents the roles table
type Role struct {
	ID          string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string     `gorm:"type:text;not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Permissions JSONB      `gorm:"type:jsonb;default:'[]'::jsonb" json:"permissions"`
	CreatedAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
