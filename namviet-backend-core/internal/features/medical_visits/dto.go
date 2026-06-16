package medical_visits

import "github.com/namviet/backend-core/internal/features/roles"

type CreateMedicalVisitRequest struct {
	AppointmentID *string `json:"appointment_id"`
	CustomerID    int64   `json:"customer_id" binding:"required"`
	DoctorID      *string `json:"doctor_id"`
	Symptoms      *string `json:"symptoms"`
}

type UpdateMedicalVisitRequest struct {
	Temperature        *float64     `json:"temperature"`
	Pulse              *int         `json:"pulse"`
	SpO2               *int         `json:"sp02"`
	BPSystolic         *int         `json:"bp_systolic"`
	BPDiastolic        *int         `json:"bp_diastolic"`
	Weight             *float64     `json:"weight"`
	Height             *float64     `json:"height"`
	Symptoms           *string      `json:"symptoms"`
	ExaminationSummary *string      `json:"examination_summary"`
	Diagnosis          *string      `json:"diagnosis"`
	ICDCode            *string      `json:"icd_code"`
	DoctorNotes        *string      `json:"doctor_notes"`
	RedFlags           *roles.JSONB `json:"red_flags"`
	Status             *string      `json:"status"`
}
