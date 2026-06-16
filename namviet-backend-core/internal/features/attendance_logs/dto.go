package attendance_logs

import "time"

type CreateAttendanceLogRequest struct {
	UserID      string    `json:"user_id" binding:"required"`
	BranchID    *int64    `json:"branch_id"`
	CheckInTime time.Time `json:"check_in_time" binding:"required"`
	CheckInIP   *string   `json:"check_in_ip"`
	CheckInLat  *float64  `json:"check_in_lat"`
	CheckInLng  *float64  `json:"check_in_lng"`
}

type UpdateAttendanceLogRequest struct {
	CheckOutTime *time.Time `json:"check_out_time"`
	CheckOutIP   *string    `json:"check_out_ip"`
	CheckOutLat  *float64   `json:"check_out_lat"`
	CheckOutLng  *float64   `json:"check_out_lng"`
	Status       *string    `json:"status"`
}
