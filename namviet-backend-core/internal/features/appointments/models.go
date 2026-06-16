package appointments

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// Appointment represents the appointments table
type Appointment struct {
	ID              string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CustomerID      int64       `gorm:"type:bigint;not null" json:"customer_id"`
	DoctorID        *string     `gorm:"type:uuid" json:"doctor_id,omitempty"`
	RoomID          *int64      `gorm:"type:bigint" json:"room_id,omitempty"`
	ServiceType     string      `gorm:"type:text;not null" json:"service_type"`
	AppointmentTime time.Time   `gorm:"type:timestamp with time zone;not null" json:"appointment_time"`
	CheckInTime     *time.Time  `gorm:"type:timestamp with time zone" json:"check_in_time,omitempty"`
	Status          string      `gorm:"type:text;default:'pending'" json:"status"`
	Symptoms        roles.JSONB `gorm:"type:jsonb;default:'[]'::jsonb" json:"symptoms,omitempty"`
	Note            *string     `gorm:"type:text" json:"note,omitempty"`
	CreatedBy       *string     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt       *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt       *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt       *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
