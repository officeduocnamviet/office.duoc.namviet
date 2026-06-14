package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

type AuthMiddleware struct {
	db *sql.DB
}

func NewAuthMiddleware(db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

func (m *AuthMiddleware) RequirePermission(requiredPermission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get Token from Header (For real app, verify JWT to get user_id)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Mocking JWT verification: Assume we parsed user_id from token
		// In production, use Supabase JWT library
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// userID := verifyJWT(tokenString)
		
		// For demo/RBAC testing, we can pass user_id in a custom header X-User-Id
		userID := r.Header.Get("X-User-Id")
		if userID == "" {
			http.Error(w, "Unauthorized: Missing X-User-Id for testing", http.StatusUnauthorized)
			return
		}

		// 2. Check Permission in Database
		query := `
			SELECT count(1) 
			FROM users u
			JOIN role_permissions rp ON u.role_id = rp.role_id
			WHERE u.id = $1 AND rp.permission_key = $2
		`
		var count int
		err := m.db.QueryRow(query, userID, requiredPermission).Scan(&count)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count == 0 {
			http.Error(w, "Forbidden: You don't have permission "+requiredPermission, http.StatusForbidden)
			return
		}

		// Pass user_id to request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
