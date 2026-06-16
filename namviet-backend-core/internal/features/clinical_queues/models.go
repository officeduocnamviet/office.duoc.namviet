package clinical_queues

import (
	"time"
)

// ClinicalQueue represents the clinical_queues table
type ClinicalQueue struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AppointmentID *string    `gorm:"type:uuid" json:"appointment_id,omitempty"`
	CustomerID    int64      `gorm:"type:bigint;not null" json:"customer_id"`
	DoctorID      *string    `gorm:"type:uuid" json:"doctor_id,omitempty"`
	QueueNumber   int        `gorm:"type:integer;not null" json:"queue_number"`
	Status        string     `gorm:"type:text;default:'waiting'" json:"status"`
	PriorityLevel string     `gorm:"type:text;default:'normal'" json:"priority_level"`
	CheckedInAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"checked_in_at,omitempty"`
	UpdatedAt     *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt     *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
