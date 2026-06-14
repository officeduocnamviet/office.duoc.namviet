package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/namviet/backend-core/internal/features/users"
	"github.com/namviet/backend-core/internal/platform/firebase"
	"github.com/namviet/backend-core/internal/platform/supabase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/namviet/backend-core/docs" // Uncommented after swag init
)

// @title Nam Viet ERP API
// @version 1.0
// @description Backend API for Nam Viet ERP System
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Connect DB & Init Firebase
	supabase.InitDB()
	firebase.InitFirebase()

	// 2. Setup Gin Router
	r := gin.Default()

	// Enable CORS (Basic config)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 3. Register Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 4. Setup Routes
	api := r.Group("/api")
	users.RegisterRoutes(api)

	// 5. Start Server
	log.Println("Server running on port 8080...")
	r.Run(":8080")
}