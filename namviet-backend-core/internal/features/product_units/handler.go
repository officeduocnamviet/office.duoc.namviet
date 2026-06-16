package product_units

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProductUnitHandler
// @Summary Create a new Product Unit
// @Description Create a unit for a specific product
// @Tags Product Units
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Param request body CreateProductUnitRequest true "Product Unit Details"
// @Success 201 {object} ProductUnit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{product_id}/units [post]
func CreateProductUnitHandler(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID format"})
		return
	}

	var req CreateProductUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unit, err := CreateProductUnitService(productID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, unit)
}

// UpdateProductUnitHandler
// @Summary Update Product Unit
// @Description Update product unit details
// @Tags Product Units
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path int true "Unit ID"
// @Param request body UpdateProductUnitRequest true "Update Details"
// @Success 200 {object} ProductUnit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/units/{unit_id} [put]
func UpdateProductUnitHandler(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("unit_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Unit ID format"})
		return
	}
	var req UpdateProductUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unit, err := UpdateProductUnitService(unitID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unit)
}

// DeleteProductUnitHandler
// @Summary Delete Product Unit
// @Description Delete a product unit
// @Tags Product Units
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unit_id path int true "Unit ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/units/{unit_id} [delete]
func DeleteProductUnitHandler(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("unit_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Unit ID format"})
		return
	}
	
	err = DeleteProductUnitService(unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
