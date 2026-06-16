package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/namviet/backend-core/internal/features/batches"
	"github.com/namviet/backend-core/internal/features/categories"
	"github.com/namviet/backend-core/internal/features/customers"
	"github.com/namviet/backend-core/internal/features/inventory"
	"github.com/namviet/backend-core/internal/features/manufacturers"
	"github.com/namviet/backend-core/internal/features/orders"
	"github.com/namviet/backend-core/internal/features/product_units"
	"github.com/namviet/backend-core/internal/features/products"
	"github.com/namviet/backend-core/internal/features/promotions"
	"github.com/namviet/backend-core/internal/features/roles"
	"github.com/namviet/backend-core/internal/features/users"
	"github.com/namviet/backend-core/internal/features/warehouses"
	"github.com/namviet/backend-core/internal/features/appointments"
	"github.com/namviet/backend-core/internal/features/clinical_queues"
	"github.com/namviet/backend-core/internal/features/medical_visits"
	"github.com/namviet/backend-core/internal/features/employees"
	"github.com/namviet/backend-core/internal/features/time_attendance"
	"github.com/namviet/backend-core/internal/features/payrolls"
	"github.com/namviet/backend-core/internal/features/fund_accounts"
	"github.com/namviet/backend-core/internal/features/finance_transactions"
	"github.com/namviet/backend-core/internal/features/chart_of_accounts"
	"github.com/namviet/backend-core/internal/features/accounting_journals"
	"github.com/namviet/backend-core/internal/features/companies"
	"github.com/namviet/backend-core/internal/features/approvals"
	"github.com/namviet/backend-core/internal/features/system_configs"
	"github.com/namviet/backend-core/internal/features/audit_logs"
	"github.com/namviet/backend-core/internal/features/integrations"
	"github.com/namviet/backend-core/internal/features/shipping_partners"
	"github.com/namviet/backend-core/internal/features/agent_workflows"
	"github.com/namviet/backend-core/internal/features/ai_agent_memories"
	"github.com/namviet/backend-core/internal/features/chats"
	"github.com/namviet/backend-core/internal/features/knowledge_vectors"
	"github.com/namviet/backend-core/internal/features/employment_contracts"
	"github.com/namviet/backend-core/internal/features/attendance_logs"
	"github.com/namviet/backend-core/internal/features/work_shifts"
	"github.com/namviet/backend-core/internal/features/training_courses"
	"github.com/namviet/backend-core/internal/features/marketing_campaigns"
	"github.com/namviet/backend-core/internal/features/customer_records"
	"github.com/namviet/backend-core/internal/features/internal_communications"
	"github.com/namviet/backend-core/internal/features/user_notifications"
	"github.com/namviet/backend-core/internal/platform/firebase"
	"github.com/namviet/backend-core/internal/platform/supabase"
	// Upload routes
	uploads "github.com/namviet/backend-core/internal/features/uploads"

	auth_middleware "github.com/namviet/backend-core/internal/middleware"
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

	// Enable CORS (Secure config)
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:3000":                       true,
			"https://namviet-omnichannel.web.app":         true,
			"https://namviet-omnichannel.firebaseapp.com": true,
		}

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
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
	publicAPI := r.Group("/api")
	users.RegisterRoutes(publicAPI)

	protectedAPI := r.Group("/api")
	protectedAPI.Use(auth_middleware.RequireAuth())

	roles.RegisterRoutes(protectedAPI)
	categories.RegisterRoutes(protectedAPI)
	manufacturers.RegisterRoutes(protectedAPI)
	products.RegisterRoutes(protectedAPI)
	product_units.RegisterRoutes(protectedAPI)
	batches.RegisterRoutes(protectedAPI)
	warehouses.RegisterRoutes(protectedAPI)
	inventory.RegisterRoutes(protectedAPI)
	customers.RegisterRoutes(protectedAPI)
	orders.RegisterRoutes(protectedAPI)
	promotions.RegisterRoutes(protectedAPI)

	// --- MODULE 4: CLINICAL (Y tế & Lâm sàng) ---
	appointments.RegisterRoutes(protectedAPI)
	clinical_queues.RegisterRoutes(protectedAPI)
	medical_visits.RegisterRoutes(protectedAPI)

	// --- MODULE 5: HR & PAYROLL (Nhân sự & Tính lương) ---
	employees.RegisterRoutes(protectedAPI)
	time_attendance.RegisterRoutes(protectedAPI)
	payrolls.RegisterRoutes(protectedAPI)

	// --- MODULE 6: FINANCE & ACCOUNTING (Tài chính & Kế toán) ---
	fund_accounts.RegisterRoutes(protectedAPI)
	finance_transactions.RegisterRoutes(protectedAPI)
	chart_of_accounts.RegisterRoutes(protectedAPI)
	uploads.RegisterRoutes(protectedAPI)
	accounting_journals.RegisterRoutes(protectedAPI)

	// --- MODULE 7: SYSTEM, APPROVALS & INTEGRATIONS (Hệ thống & Tích hợp) ---
	companies.RegisterRoutes(protectedAPI)
	approvals.RegisterRoutes(protectedAPI)
	system_configs.RegisterRoutes(protectedAPI)
	audit_logs.RegisterRoutes(protectedAPI)
	integrations.RegisterRoutes(protectedAPI)
	shipping_partners.RegisterRoutes(protectedAPI)

	// --- MODULE 8: AI & CHATBOT ECOSYSTEM (AI & Chatbot) ---
	agent_workflows.RegisterRoutes(protectedAPI)
	ai_agent_memories.RegisterRoutes(protectedAPI)
	chats.RegisterRoutes(protectedAPI)
	knowledge_vectors.RegisterRoutes(protectedAPI)

	// --- MODULE 9: ADVANCED HR & OPERATIONS (Nhân sự & Vận hành Nâng cao) ---
	employment_contracts.RegisterRoutes(protectedAPI)
	attendance_logs.RegisterRoutes(protectedAPI)
	work_shifts.RegisterRoutes(protectedAPI)
	training_courses.RegisterRoutes(protectedAPI)

	// --- MODULE 10: CRM & MARKETING ADVANCED (CRM & Marketing Nâng cao) ---
	marketing_campaigns.RegisterRoutes(protectedAPI)
	customer_records.RegisterRoutes(protectedAPI)
	internal_communications.RegisterRoutes(protectedAPI)
	user_notifications.RegisterRoutes(protectedAPI)

	// 5. Start Server
	log.Println("Server running on port 8080...")
	r.Run(":8080")
}