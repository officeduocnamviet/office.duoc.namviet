package manufacturers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllManufacturersHandler
// @Summary Get All Manufacturers
// @Description Retrieve a list of all manufacturers
// @Tags Manufacturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Manufacturer
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /manufacturers [get]
func GetAllManufacturersHandler(c *gin.Context) {
	manufacturers, err := GetAllManufacturersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, manufacturers)
}

// GetManufacturerHandler
// @Summary Get Manufacturer by ID
// @Description Retrieve a specific manufacturer
// @Tags Manufacturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Manufacturer ID"
// @Success 200 {object} Manufacturer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /manufacturers/{id} [get]
func GetManufacturerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	manufacturer, err := GetManufacturerByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, manufacturer)
}

// CreateManufacturerHandler
// @Summary Create a new Manufacturer
// @Description Create a new manufacturer
// @Tags Manufacturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateManufacturerRequest true "Manufacturer Details"
// @Success 201 {object} Manufacturer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /manufacturers [post]
func CreateManufacturerHandler(c *gin.Context) {
	var req CreateManufacturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manufacturer, err := CreateManufacturerService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, manufacturer)
}

// UpdateManufacturerHandler
// @Summary Update Manufacturer
// @Description Update manufacturer details
// @Tags Manufacturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Manufacturer ID"
// @Param request body UpdateManufacturerRequest true "Update Details"
// @Success 200 {object} Manufacturer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /manufacturers/{id} [put]
func UpdateManufacturerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateManufacturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manufacturer, err := UpdateManufacturerService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, manufacturer)
}
