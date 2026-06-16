package batches

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBatchesByProductIDHandler
// @Summary Get Batches by Product ID
// @Description Retrieve a list of batches for a specific product
// @Tags Batches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {array} Batch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id}/batches [get]
func GetBatchesByProductIDHandler(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID format"})
		return
	}
	batches, err := GetBatchesByProductIDService(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, batches)
}

// GetBatchHandler
// @Summary Get Batch by ID
// @Description Retrieve a specific batch
// @Tags Batches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Batch ID"
// @Success 200 {object} Batch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /batches/{id} [get]
func GetBatchHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	batch, err := GetBatchByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, batch)
}

// CreateBatchHandler
// @Summary Create a new Batch
// @Description Create a new batch
// @Tags Batches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBatchRequest true "Batch Details"
// @Success 201 {object} Batch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batches [post]
func CreateBatchHandler(c *gin.Context) {
	var req CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch, err := CreateBatchService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, batch)
}

// UpdateBatchHandler
// @Summary Update Batch
// @Description Update batch details
// @Tags Batches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Batch ID"
// @Param request body UpdateBatchRequest true "Update Details"
// @Success 200 {object} Batch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batches/{id} [put]
func UpdateBatchHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch, err := UpdateBatchService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, batch)
}
