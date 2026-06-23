package config

import "testing"

func TestLoad_Defaults(t *testing.T) {
	get := func(string) string { return "" }
	cfg, err := Load(get)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr default = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("AppEnv default = %q, want development", cfg.AppEnv)
	}
}

func TestLoad_RequiresDatabaseURLInProd(t *testing.T) {
	get := func(k string) string {
		if k == "APP_ENV" {
			return "production"
		}
		return ""
	}
	if _, err := Load(get); err == nil {
		t.Fatal("expected error when DATABASE_URL missing in production")
	}
}

func TestLoad_ProdWithDatabaseURLOK(t *testing.T) {
	env := map[string]string{
		"APP_ENV":      "production",
		"DATABASE_URL": "postgres://user:pass@host:5432/db",
	}
	cfg, err := Load(func(k string) string { return env[k] })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DatabaseURL == "" {
		t.Error("DatabaseURL should be set")
	}
}
