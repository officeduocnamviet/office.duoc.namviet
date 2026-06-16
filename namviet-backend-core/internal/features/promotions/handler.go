package promotions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllPromotionsHandler
// @Summary Get All Promotions
// @Description Retrieve a list of all promotions
// @Tags Promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Promotion
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /promotions [get]
func GetAllPromotionsHandler(c *gin.Context) {
	promotions, err := GetAllPromotionsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, promotions)
}

// GetPromotionHandler
// @Summary Get Promotion by ID
// @Description Retrieve a specific promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Promotion ID"
// @Success 200 {object} Promotion
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /promotions/{id} [get]
func GetPromotionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	promotion, err := GetPromotionByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, promotion)
}

// CreatePromotionHandler
// @Summary Create a new Promotion
// @Description Create a new promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePromotionRequest true "Promotion Details"
// @Success 201 {object} Promotion
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /promotions [post]
func CreatePromotionHandler(c *gin.Context) {
	var req CreatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promotion, err := CreatePromotionService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, promotion)
}

// UpdatePromotionHandler
// @Summary Update Promotion
// @Description Update promotion details
// @Tags Promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Promotion ID"
// @Param request body UpdatePromotionRequest true "Update Details"
// @Success 200 {object} Promotion
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /promotions/{id} [put]
func UpdatePromotionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promotion, err := UpdatePromotionService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, promotion)
}

// DeletePromotionHandler
// @Summary Delete Promotion
// @Description Delete a promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Promotion ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /promotions/{id} [delete]
func DeletePromotionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeletePromotionService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
