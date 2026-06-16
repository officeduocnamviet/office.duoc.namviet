package tests

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/namviet/backend-core/internal/features/companies"
	"github.com/namviet/backend-core/internal/features/roles"
	"github.com/namviet/backend-core/internal/features/users"
	categories "github.com/namviet/backend-core/internal/features/categories"
	customers "github.com/namviet/backend-core/internal/features/customers"
	"github.com/namviet/backend-core/internal/features/inventory"
	orders "github.com/namviet/backend-core/internal/features/orders"
	"github.com/namviet/backend-core/internal/features/products"
	warehouses "github.com/namviet/backend-core/internal/features/warehouses"
	accounting_journals "github.com/namviet/backend-core/internal/features/accounting_journals"
	finance_transactions "github.com/namviet/backend-core/internal/features/finance_transactions"
	fund_accounts "github.com/namviet/backend-core/internal/features/fund_accounts"
	auth_middleware "github.com/namviet/backend-core/internal/middleware"
	"github.com/namviet/backend-core/internal/platform/supabase"
)

// SetupTestRouter initializes the database and returns a Gin engine with routes mounted
func SetupTestRouter() *gin.Engine {
	// Load .env from parent directory
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Note: .env file not found in parent dir, relying on existing env variables")
	}

	// Make sure we have a connection
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL must be set to run tests")
	}

	supabase.InitDB()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Public Routes
	publicAPI := r.Group("/api")
	users.RegisterRoutes(publicAPI)

	// Protected Routes
	protectedAPI := r.Group("/api")
	protectedAPI.Use(auth_middleware.RequireAuth())

	roles.RegisterRoutes(protectedAPI)
	companies.RegisterRoutes(protectedAPI)
	categories.RegisterRoutes(protectedAPI)
	products.RegisterRoutes(protectedAPI)
	warehouses.RegisterRoutes(protectedAPI)
	inventory.RegisterRoutes(protectedAPI)
	customers.RegisterRoutes(protectedAPI)
	orders.RegisterRoutes(protectedAPI)
	accounting_journals.RegisterRoutes(protectedAPI)
	finance_transactions.RegisterRoutes(protectedAPI)
	fund_accounts.RegisterRoutes(protectedAPI)

	return r
}

// PerformRequest is a helper function to perform HTTP requests in tests
func PerformRequest(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	// Use master token to bypass auth middleware
	req.Header.Set("Authorization", "Bearer namviet-admin-super-key")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
