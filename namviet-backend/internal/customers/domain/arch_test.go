package domain_test

import (
	"go/build"
	"strings"
	"testing"
)

// TestDomain_NoInfraImports là ARCHITECTURE-FITNESS TEST (ARCHITECTURE.md §11):
// ép ranh giới hexagon — package domain của customers phải THUẦN, KHÔNG import hạ
// tầng (pgx, huma, net/http, sqlc appdb, các platform adapter, app, adapters).
// Cho phép common/money (shared kernel trung lập domain). Dùng go/build (stdlib).
func TestDomain_NoInfraImports(t *testing.T) {
	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		t.Fatalf("import dir domain: %v", err)
	}

	forbidden := []string{
		"github.com/jackc/pgx",          // SQL driver
		"github.com/danielgtaylor/huma", // HTTP framework
		"github.com/go-chi/chi",         // router
		"net/http",                      // HTTP
		"github.com/Maneva-AI/namviet-backend/internal/platform",           // mọi adapter hạ tầng
		"github.com/Maneva-AI/namviet-backend/internal/customers/app",      // không phụ thuộc ngược lên app
		"github.com/Maneva-AI/namviet-backend/internal/customers/internal", // không phụ thuộc adapter
	}

	all := append([]string{}, pkg.Imports...)
	all = append(all, pkg.TestImports...)

	for _, imp := range all {
		for _, bad := range forbidden {
			if imp == bad || strings.HasPrefix(imp, bad) {
				t.Errorf("domain import bị cấm %q (khớp %q) — domain phải thuần, không phụ thuộc hạ tầng", imp, bad)
			}
		}
	}
}
