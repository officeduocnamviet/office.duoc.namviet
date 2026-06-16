package medical_visits

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// MedicalVisit represents the medical_visits table
type MedicalVisit struct {
	ID                 string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AppointmentID      *string     `gorm:"type:uuid" json:"appointment_id,omitempty"`
	CustomerID         int64       `gorm:"type:bigint;not null" json:"customer_id"`
	DoctorID           *string     `gorm:"type:uuid" json:"doctor_id,omitempty"`
	Temperature        *float64    `gorm:"type:numeric" json:"temperature,omitempty"`
	Pulse              *int        `gorm:"type:integer" json:"pulse,omitempty"`
	SpO2               *int        `gorm:"type:integer" json:"sp02,omitempty"`
	BPSystolic         *int        `gorm:"type:integer" json:"bp_systolic,omitempty"`
	BPDiastolic        *int        `gorm:"type:integer" json:"bp_diastolic,omitempty"`
	Weight             *float64    `gorm:"type:numeric" json:"weight,omitempty"`
	Height             *float64    `gorm:"type:numeric" json:"height,omitempty"`
	Symptoms           *string     `gorm:"type:text" json:"symptoms,omitempty"`
	ExaminationSummary *string     `gorm:"type:text" json:"examination_summary,omitempty"`
	Diagnosis          *string     `gorm:"type:text" json:"diagnosis,omitempty"`
	ICDCode            *string     `gorm:"type:text" json:"icd_code,omitempty"`
	DoctorNotes        *string     `gorm:"type:text" json:"doctor_notes,omitempty"`
	RedFlags           roles.JSONB `gorm:"type:jsonb;default:'[]'::jsonb" json:"red_flags,omitempty"`
	Status             string      `gorm:"type:text;default:'in_progress'" json:"status"`
	CreatedBy          *string     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy          *string     `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt          *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt          *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt          *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
