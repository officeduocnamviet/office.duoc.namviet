package clinical_queues

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllClinicalQueuesHandler
// @Summary Get All Clinical Queues
// @Description Retrieve a list of all clinical queues
// @Tags Clinical Queues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ClinicalQueue
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clinical-queues [get]
func GetAllClinicalQueuesHandler(c *gin.Context) {
	queues, err := GetAllClinicalQueuesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, queues)
}

// GetClinicalQueueHandler
// @Summary Get Clinical Queue by ID
// @Description Retrieve a specific clinical queue
// @Tags Clinical Queues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Clinical Queue ID"
// @Success 200 {object} ClinicalQueue
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clinical-queues/{id} [get]
func GetClinicalQueueHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	queue, err := GetClinicalQueueByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, queue)
}

// CreateClinicalQueueHandler
// @Summary Create a new Clinical Queue
// @Description Create a new clinical queue
// @Tags Clinical Queues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateClinicalQueueRequest true "Clinical Queue Details"
// @Success 201 {object} ClinicalQueue
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clinical-queues [post]
func CreateClinicalQueueHandler(c *gin.Context) {
	var req CreateClinicalQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queue, err := CreateClinicalQueueService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, queue)
}

// UpdateClinicalQueueHandler
// @Summary Update Clinical Queue
// @Description Update clinical queue details
// @Tags Clinical Queues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Clinical Queue ID"
// @Param request body UpdateClinicalQueueRequest true "Update Details"
// @Success 200 {object} ClinicalQueue
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clinical-queues/{id} [put]
func UpdateClinicalQueueHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateClinicalQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queue, err := UpdateClinicalQueueService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, queue)
}

// DeleteClinicalQueueHandler
// @Summary Delete Clinical Queue
// @Description Soft delete a clinical queue
// @Tags Clinical Queues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Clinical Queue ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clinical-queues/{id} [delete]
func DeleteClinicalQueueHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteClinicalQueueService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
