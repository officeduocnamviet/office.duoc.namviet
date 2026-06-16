package training_courses

import "time"

func GetAllTrainingCoursesService() ([]TrainingCourse, error) {
	return GetAllTrainingCourses()
}

func GetTrainingCourseByIDService(id int64) (*TrainingCourse, error) {
	return GetTrainingCourseByID(id)
}

func CreateTrainingCourseService(req CreateTrainingCourseRequest) (*TrainingCourse, error) {
	course := &TrainingCourse{
		Title:        req.Title,
		ContentType:  req.ContentType,
		ContentURL:   req.ContentURL,
		PassingScore: req.PassingScore,
		Status:       "active",
	}

	if err := CreateTrainingCourse(course); err != nil {
		return nil, err
	}
	return course, nil
}

func UpdateTrainingCourseService(id int64, req UpdateTrainingCourseRequest) (*TrainingCourse, error) {
	course, err := GetTrainingCourseByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.ContentType != nil {
		course.ContentType = *req.ContentType
	}
	if req.ContentURL != nil {
		course.ContentURL = req.ContentURL
	}
	if req.PassingScore != nil {
		course.PassingScore = req.PassingScore
	}
	if req.Status != nil {
		course.Status = *req.Status
	}

	now := time.Now()
	course.UpdatedAt = &now

	if err := UpdateTrainingCourse(course); err != nil {
		return nil, err
	}
	return course, nil
}

func DeleteTrainingCourseService(id int64) error {
	return DeleteTrainingCourse(id)
}
