package training_courses

type CreateTrainingCourseRequest struct {
	Title        string  `json:"title" binding:"required"`
	ContentType  string  `json:"content_type" binding:"required"`
	ContentURL   *string `json:"content_url"`
	PassingScore *int    `json:"passing_score"`
}

type UpdateTrainingCourseRequest struct {
	Title        *string `json:"title"`
	ContentType  *string `json:"content_type"`
	ContentURL   *string `json:"content_url"`
	PassingScore *int    `json:"passing_score"`
	Status       *string `json:"status"`
}
