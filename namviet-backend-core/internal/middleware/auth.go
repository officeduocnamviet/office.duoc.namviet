package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth checks JWT Token or Master Token for authentication
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// MASTER TOKEN cho phép test API mà chưa cần form Login
		if tokenString == "namviet-admin-super-key" {
			c.Set("userID", "00000000-0000-0000-0000-000000000001") // Mock Admin ID
			c.Set("role", "admin")
			c.Next()
			return
		}

		// Xử lý JWT thật (Mock Secret for now)
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("namviet-secret-key-1234"), nil
		})

		if token != nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("userID", claims["sub"])
				c.Set("role", claims["role"])
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
	}
}
