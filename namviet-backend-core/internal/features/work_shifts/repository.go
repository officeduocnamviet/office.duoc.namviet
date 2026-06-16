package work_shifts

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Work Shifts
func GetAllWorkShifts() ([]WorkShift, error) {
	var results []WorkShift
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetWorkShiftByID(id int64) (*WorkShift, error) {
	var result WorkShift
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("work shift not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateWorkShift(data *WorkShift) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateWorkShift(data *WorkShift) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteWorkShift(id int64) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&WorkShift{}).Error
}

// Shift Assignments
func GetAllShiftAssignments() ([]ShiftAssignment, error) {
	var results []ShiftAssignment
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetShiftAssignmentByID(id int64) (*ShiftAssignment, error) {
	var result ShiftAssignment
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shift assignment not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateShiftAssignment(data *ShiftAssignment) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateShiftAssignment(data *ShiftAssignment) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteShiftAssignment(id int64) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&ShiftAssignment{}).Error
}

// Shift Handovers
func GetAllShiftHandovers() ([]ShiftHandover, error) {
	var results []ShiftHandover
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetShiftHandoverByID(id string) (*ShiftHandover, error) {
	var result ShiftHandover
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shift handover not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateShiftHandover(data *ShiftHandover) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateShiftHandover(data *ShiftHandover) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteShiftHandover(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&ShiftHandover{}).Error
}
