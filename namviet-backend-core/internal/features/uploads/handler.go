package uploads

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// @Summary Upload an image to Cloudinary
// @Description Uploads a multipart/form-data image file to Cloudinary and returns the optimized public URL.
// @Tags uploads
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to upload"
// @Success 200 {object} map[string]interface{} "{"url": "https://..."}"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /api/upload [post]
func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'image' field in form data"})
		return
	}
	defer file.Close()

	// Check if CLOUDINARY_URL exists
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CLOUDINARY_URL is not configured"})
		return
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary client"})
		return
	}

	// Upload to Cloudinary
	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "namviet/uploads",
		PublicID: header.Filename, // Optional: preserve original name or let it auto-generate
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary", "details": err.Error()})
		return
	}

	// Format URL to include q_auto,f_webp
	// Original URL: https://res.cloudinary.com/cloud_name/image/upload/v12345/folder/file.ext
	// Target URL: https://res.cloudinary.com/cloud_name/image/upload/f_webp,q_auto/v12345/folder/file.ext
	secureURL := resp.SecureURL
	if strings.Contains(secureURL, "/upload/") {
		secureURL = strings.Replace(secureURL, "/upload/", "/upload/f_webp,q_auto/", 1)
	}

	c.JSON(http.StatusOK, gin.H{
		"url": secureURL,
	})
}
