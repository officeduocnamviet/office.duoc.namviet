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
	roles.RegisterRoutes(api)
	categories.RegisterRoutes(api)
	manufacturers.RegisterRoutes(api)
	products.RegisterRoutes(api)
	product_units.RegisterRoutes(api)
	batches.RegisterRoutes(api)
	warehouses.RegisterRoutes(api)
	inventory.RegisterRoutes(api)
	customers.RegisterRoutes(api)
	orders.RegisterRoutes(api)
	promotions.RegisterRoutes(api)

	// --- MODULE 4: CLINICAL (Y tế & Lâm sàng) ---
	appointments.RegisterRoutes(api)
	clinical_queues.RegisterRoutes(api)
	medical_visits.RegisterRoutes(api)

	// --- MODULE 5: HR & PAYROLL (Nhân sự & Tính lương) ---
	employees.RegisterRoutes(api)
	time_attendance.RegisterRoutes(api)
	payrolls.RegisterRoutes(api)

	// --- MODULE 6: FINANCE & ACCOUNTING (Tài chính & Kế toán) ---
	fund_accounts.RegisterRoutes(api)
	finance_transactions.RegisterRoutes(api)
	chart_of_accounts.RegisterRoutes(api)
	accounting_journals.RegisterRoutes(api)

	// --- MODULE 7: SYSTEM, APPROVALS & INTEGRATIONS (Hệ thống & Tích hợp) ---
	companies.RegisterRoutes(api)
	approvals.RegisterRoutes(api)
	system_configs.RegisterRoutes(api)
	audit_logs.RegisterRoutes(api)
	integrations.RegisterRoutes(api)
	shipping_partners.RegisterRoutes(api)

	// --- MODULE 8: AI & CHATBOT ECOSYSTEM (AI & Chatbot) ---
	agent_workflows.RegisterRoutes(api)
	ai_agent_memories.RegisterRoutes(api)
	chats.RegisterRoutes(api)
	knowledge_vectors.RegisterRoutes(api)

	// --- MODULE 9: ADVANCED HR & OPERATIONS (Nhân sự & Vận hành Nâng cao) ---
	employment_contracts.RegisterRoutes(api)
	attendance_logs.RegisterRoutes(api)
	work_shifts.RegisterRoutes(api)
	training_courses.RegisterRoutes(api)

	// --- MODULE 10: CRM & MARKETING ADVANCED (CRM & Marketing Nâng cao) ---
	marketing_campaigns.RegisterRoutes(api)
	customer_records.RegisterRoutes(api)
	internal_communications.RegisterRoutes(api)
	user_notifications.RegisterRoutes(api)

	// 5. Start Server
	log.Println("Server running on port 8080...")
	r.Run(":8080")
}