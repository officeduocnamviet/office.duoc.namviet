package uploads

import "github.com/gin-gonic/gin"

// RegisterRoutes registers the upload routes
func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/upload", UploadImage)
}
