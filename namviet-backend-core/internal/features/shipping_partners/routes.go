package shipping_partners

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the shipping partners routes
func RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/shipping-partners")
	{
		group.GET("", GetAllShippingPartnersHandler)
		group.GET("/:id", GetShippingPartnerHandler)
		group.POST("", CreateShippingPartnerHandler)
		group.PUT("/:id", UpdateShippingPartnerHandler)
		group.DELETE("/:id", DeleteShippingPartnerHandler)
	}
}
