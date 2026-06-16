package time_attendance

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// TimeAttendance represents the time_attendance table
type TimeAttendance struct {
	ID            string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	EmployeeID    string      `gorm:"type:uuid;not null" json:"employee_id"`
	Date          time.Time   `gorm:"type:date;not null" json:"date"`
	CheckIn       *time.Time  `gorm:"type:timestamp with time zone" json:"check_in,omitempty"`
	CheckOut      *time.Time  `gorm:"type:timestamp with time zone" json:"check_out,omitempty"`
	Status        string      `gorm:"type:text;default:'present'" json:"status"`
	ShiftType     string      `gorm:"type:text;default:'morning'" json:"shift_type"`
	OvertimeHours *float64    `gorm:"type:numeric;default:0" json:"overtime_hours,omitempty"`
	Location      roles.JSONB `gorm:"type:jsonb;default:'{}'::jsonb" json:"location,omitempty"`
	DeviceInfo    roles.JSONB `gorm:"type:jsonb;default:'{}'::jsonb" json:"device_info,omitempty"`
	Note          *string     `gorm:"type:text" json:"note,omitempty"`
	CreatedAt     *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt     *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}
