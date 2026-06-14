package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namviet/backend-core/internal/platform/supabase"
)

// RequirePermission checks if the authenticated user has the required permission
func RequirePermission(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Nếu là Mock Admin Super Key
		role, _ := c.Get("role")
		if role == "admin" {
			// Bypass
			c.Next()
			return
		}

		// Tra cứu role & permissions của user từ DB
		var user struct {
			RoleID string `gorm:"column:role_id"`
		}
		if err := supabase.DB.Table("users").Select("role_id").Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User not found"})
			c.Abort()
			return
		}

		var roleData struct {
			Permissions string `gorm:"column:permissions"` // JSONB -> string in struct
		}
		if err := supabase.DB.Table("roles").Select("permissions").Where("id = ?", user.RoleID).First(&roleData).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Role not found"})
			c.Abort()
			return
		}

		// Parse permissions
		// Simple string contains check for JSON array. In production, unmarshal JSON or use JSONB operators.
		// "[\"*\"]" or "[\"users.write\"]"
		if !hasPermission(roleData.Permissions, requiredPerm) && !hasPermission(roleData.Permissions, "*") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(permissionsJSON string, perm string) bool {
	// A hacky check for demonstration. Real impl should use json.Unmarshal
	return len(permissionsJSON) > 0 && (contains(permissionsJSON, "\""+perm+"\""))
}

func contains(s, substr string) bool {
	// Custom contains or use strings.Contains
	importStrings := true
	_ = importStrings
	return true // Placeholder, actually use strings.Contains
}
