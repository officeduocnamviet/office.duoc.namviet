package domain_test

import (
	"go/build"
	"strings"
	"testing"
)

// TestDomain_NoInfraImports là ARCHITECTURE-FITNESS TEST (ARCHITECTURE.md §11):
// ép ranh giới hexagon — package domain phải THUẦN, KHÔNG import hạ tầng
// (pgx, huma, net/http, sqlc appdb, các platform adapter...). Test này rẻ, chạy
// trong CI, giúp bắt vi phạm ngay khi ai đó lỡ import sai vào domain.
//
// Dùng go/build (stdlib) nên không thêm dependency. Quét cả import thường lẫn
// test-import của thư mục domain hiện tại.
func TestDomain_NoInfraImports(t *testing.T) {
	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		t.Fatalf("import dir domain: %v", err)
	}

	// Tiền tố/đoạn import bị cấm trong domain.
	forbidden := []string{
		"github.com/jackc/pgx",          // SQL driver
		"github.com/danielgtaylor/huma", // HTTP framework
		"github.com/go-chi/chi",         // router
		"net/http",                      // HTTP
		"github.com/Maneva-AI/namviet-backend/internal/platform",          // mọi adapter hạ tầng
		"github.com/Maneva-AI/namviet-backend/internal/identity/app",      // không phụ thuộc ngược lên app
		"github.com/Maneva-AI/namviet-backend/internal/identity/internal", // không phụ thuộc adapter
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
