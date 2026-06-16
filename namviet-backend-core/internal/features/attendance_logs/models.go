package attendance_logs

import (
	"time"
)

// AttendanceLog represents the attendance_logs table
type AttendanceLog struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID       string     `gorm:"type:uuid;not null" json:"user_id"`
	BranchID     *int64     `gorm:"type:bigint" json:"branch_id,omitempty"`
	CheckInTime  time.Time  `gorm:"type:timestamp with time zone;default:now();not null" json:"check_in_time"`
	CheckInIP    *string    `gorm:"type:text" json:"check_in_ip,omitempty"`
	CheckInLat   *float64   `gorm:"type:numeric" json:"check_in_lat,omitempty"`
	CheckInLng   *float64   `gorm:"type:numeric" json:"check_in_lng,omitempty"`
	CheckOutTime *time.Time `gorm:"type:timestamp with time zone" json:"check_out_time,omitempty"`
	CheckOutIP   *string    `gorm:"type:text" json:"check_out_ip,omitempty"`
	CheckOutLat  *float64   `gorm:"type:numeric" json:"check_out_lat,omitempty"`
	CheckOutLng  *float64   `gorm:"type:numeric" json:"check_out_lng,omitempty"`
	Status       string     `gorm:"type:text;default:'present';not null" json:"status"`
}
