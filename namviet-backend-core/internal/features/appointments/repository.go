package appointments

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllAppointments() ([]Appointment, error) {
	var results []Appointment
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetAppointmentByID(id string) (*Appointment, error) {
	var result Appointment
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateAppointment(data *Appointment) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateAppointment(data *Appointment) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteAppointment(id string) error {
	db := supabase.DB
	return db.Model(&Appointment{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
