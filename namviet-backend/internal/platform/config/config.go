package config

import (
	"fmt"
	"strings"
)

// Config gom toàn bộ cấu hình runtime đọc từ môi trường.
type Config struct {
	AppEnv       string
	HTTPAddr     string
	DatabaseURL  string
	OTLPEndpoint string

	// Dành cho Plan 01 (auth/RBAC) — chưa dùng ở Phase 0.
	JWTPublicKeyPEM  string
	JWTPrivateKeyPEM string
}

// Load đọc cấu hình từ hàm getenv (thường là os.Getenv) để dễ test.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		AppEnv:           def(getenv("APP_ENV"), "development"),
		HTTPAddr:         def(getenv("HTTP_ADDR"), ":8080"),
		DatabaseURL:      getenv("DATABASE_URL"),
		OTLPEndpoint:     getenv("OTLP_ENDPOINT"),
		JWTPublicKeyPEM:  getenv("JWT_PUBLIC_KEY_PEM"),
		JWTPrivateKeyPEM: getenv("JWT_PRIVATE_KEY_PEM"),
	}
	if strings.EqualFold(cfg.AppEnv, "production") && cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required in production")
	}
	return cfg, nil
}

func def(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
