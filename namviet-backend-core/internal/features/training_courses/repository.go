package training_courses

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllTrainingCourses() ([]TrainingCourse, error) {
	var results []TrainingCourse
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetTrainingCourseByID(id int64) (*TrainingCourse, error) {
	var result TrainingCourse
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("training course not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateTrainingCourse(data *TrainingCourse) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateTrainingCourse(data *TrainingCourse) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteTrainingCourse(id int64) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&TrainingCourse{}).Error
}
