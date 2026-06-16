package appointments

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreateAppointmentRequest struct {
	CustomerID      int64       `json:"customer_id" binding:"required"`
	DoctorID        *string     `json:"doctor_id"`
	RoomID          *int64      `json:"room_id"`
	ServiceType     string      `json:"service_type" binding:"required"`
	AppointmentTime time.Time   `json:"appointment_time" binding:"required"`
	Symptoms        roles.JSONB `json:"symptoms"`
	Note            *string     `json:"note"`
}

type UpdateAppointmentRequest struct {
	DoctorID        *string      `json:"doctor_id"`
	RoomID          *int64       `json:"room_id"`
	ServiceType     *string      `json:"service_type"`
	AppointmentTime *time.Time   `json:"appointment_time"`
	CheckInTime     *time.Time   `json:"check_in_time"`
	Status          *string      `json:"status"`
	Symptoms        *roles.JSONB `json:"symptoms"`
	Note            *string      `json:"note"`
}
