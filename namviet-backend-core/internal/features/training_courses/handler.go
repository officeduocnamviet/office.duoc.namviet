package training_courses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllTrainingCoursesHandler
// @Summary Get All Training Courses
// @Description Retrieve a list of all training courses
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} TrainingCourse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /training-courses [get]
func GetAllTrainingCoursesHandler(c *gin.Context) {
	items, err := GetAllTrainingCoursesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetTrainingCourseHandler
// @Summary Get Training Course by ID
// @Description Retrieve a specific training course
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Course ID"
// @Success 200 {object} TrainingCourse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /training-courses/{id} [get]
func GetTrainingCourseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetTrainingCourseByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateTrainingCourseHandler
// @Summary Create a new Training Course
// @Description Create a new training course
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTrainingCourseRequest true "Course Details"
// @Success 201 {object} TrainingCourse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /training-courses [post]
func CreateTrainingCourseHandler(c *gin.Context) {
	var req CreateTrainingCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateTrainingCourseService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateTrainingCourseHandler
// @Summary Update Training Course
// @Description Update training course details
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Course ID"
// @Param request body UpdateTrainingCourseRequest true "Update Details"
// @Success 200 {object} TrainingCourse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /training-courses/{id} [put]
func UpdateTrainingCourseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateTrainingCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateTrainingCourseService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteTrainingCourseHandler
// @Summary Delete Training Course
// @Description Delete a training course
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Course ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /training-courses/{id} [delete]
func DeleteTrainingCourseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteTrainingCourseService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
