package shipping_partners

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllShippingPartnersHandler
// @Summary Get All Shipping Partners
// @Description Retrieve a list of all shipping partners
// @Tags Shipping Partners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ShippingPartner
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shipping-partners [get]
func GetAllShippingPartnersHandler(c *gin.Context) {
	items, err := GetAllShippingPartnersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetShippingPartnerHandler
// @Summary Get Shipping Partner by ID
// @Description Retrieve a specific shipping partner
// @Tags Shipping Partners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Shipping Partner ID"
// @Success 200 {object} ShippingPartner
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /shipping-partners/{id} [get]
func GetShippingPartnerHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetShippingPartnerByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateShippingPartnerHandler
// @Summary Create a new Shipping Partner
// @Description Create a new shipping partner
// @Tags Shipping Partners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateShippingPartnerRequest true "Shipping Partner Details"
// @Success 201 {object} ShippingPartner
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shipping-partners [post]
func CreateShippingPartnerHandler(c *gin.Context) {
	var req CreateShippingPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateShippingPartnerService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateShippingPartnerHandler
// @Summary Update Shipping Partner
// @Description Update shipping partner details
// @Tags Shipping Partners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Shipping Partner ID"
// @Param request body UpdateShippingPartnerRequest true "Update Details"
// @Success 200 {object} ShippingPartner
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shipping-partners/{id} [put]
func UpdateShippingPartnerHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateShippingPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateShippingPartnerService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteShippingPartnerHandler
// @Summary Delete Shipping Partner
// @Description Soft delete a shipping partner
// @Tags Shipping Partners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Shipping Partner ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shipping-partners/{id} [delete]
func DeleteShippingPartnerHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteShippingPartnerService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
