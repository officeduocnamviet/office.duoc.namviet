package warehouses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllWarehousesHandler
// @Summary Get All Warehouses
// @Description Retrieve a list of all warehouses
// @Tags Warehouses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Warehouse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /warehouses [get]
func GetAllWarehousesHandler(c *gin.Context) {
	warehouses, err := GetAllWarehousesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouses)
}

// GetWarehouseHandler
// @Summary Get Warehouse by ID
// @Description Retrieve a specific warehouse
// @Tags Warehouses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Warehouse ID"
// @Success 200 {object} Warehouse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /warehouses/{id} [get]
func GetWarehouseHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	warehouse, err := GetWarehouseByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouse)
}

// CreateWarehouseHandler
// @Summary Create a new Warehouse
// @Description Create a new warehouse
// @Tags Warehouses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWarehouseRequest true "Warehouse Details"
// @Success 201 {object} Warehouse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /warehouses [post]
func CreateWarehouseHandler(c *gin.Context) {
	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouse, err := CreateWarehouseService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, warehouse)
}

// UpdateWarehouseHandler
// @Summary Update Warehouse
// @Description Update warehouse details
// @Tags Warehouses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Warehouse ID"
// @Param request body UpdateWarehouseRequest true "Update Details"
// @Success 200 {object} Warehouse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /warehouses/{id} [put]
func UpdateWarehouseHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouse, err := UpdateWarehouseService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouse)
}

// DeleteWarehouseHandler
// @Summary Delete Warehouse
// @Description Delete a warehouse
// @Tags Warehouses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Warehouse ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /warehouses/{id} [delete]
func DeleteWarehouseHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := DeleteWarehouseService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
