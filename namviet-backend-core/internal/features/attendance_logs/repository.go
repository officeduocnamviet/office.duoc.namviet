package attendance_logs

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllAttendanceLogs() ([]AttendanceLog, error) {
	var results []AttendanceLog
	db := supabase.DB
	err := db.Order("check_in_time DESC").Find(&results).Error
	return results, err
}

func GetAttendanceLogByID(id string) (*AttendanceLog, error) {
	var result AttendanceLog
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attendance log not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateAttendanceLog(data *AttendanceLog) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateAttendanceLog(data *AttendanceLog) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteAttendanceLog(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&AttendanceLog{}).Error
}
