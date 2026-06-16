package time_attendance

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllTimeAttendances() ([]TimeAttendance, error) {
	var results []TimeAttendance
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetTimeAttendanceByID(id string) (*TimeAttendance, error) {
	var result TimeAttendance
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("time attendance not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateTimeAttendance(data *TimeAttendance) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateTimeAttendance(data *TimeAttendance) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteTimeAttendance(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&TimeAttendance{}).Error
}
