package time_attendance

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreateTimeAttendanceRequest struct {
	EmployeeID    string       `json:"employee_id" binding:"required"`
	Date          time.Time    `json:"date" binding:"required"`
	CheckIn       *time.Time   `json:"check_in"`
	CheckOut      *time.Time   `json:"check_out"`
	Status        *string      `json:"status"`
	ShiftType     *string      `json:"shift_type"`
	OvertimeHours *float64     `json:"overtime_hours"`
	Location      *roles.JSONB `json:"location"`
	DeviceInfo    *roles.JSONB `json:"device_info"`
	Note          *string      `json:"note"`
}

type UpdateTimeAttendanceRequest struct {
	CheckIn       *time.Time   `json:"check_in"`
	CheckOut      *time.Time   `json:"check_out"`
	Status        *string      `json:"status"`
	ShiftType     *string      `json:"shift_type"`
	OvertimeHours *float64     `json:"overtime_hours"`
	Location      *roles.JSONB `json:"location"`
	DeviceInfo    *roles.JSONB `json:"device_info"`
	Note          *string      `json:"note"`
}
