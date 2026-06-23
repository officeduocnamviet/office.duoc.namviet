package posting_test

import (
	"go/build"
	"strings"
	"testing"
)

// TestPosting_NoInfraImports là ARCHITECTURE-FITNESS TEST: package posting là lớp
// TEMPLATE THUẦN (chỉ dựng bút toán, không post), nên phải sạch hạ tầng giống
// domain — chỉ stdlib + domain + common/money. KHÔNG pgx/http/huma/app/adapter.
// Nếu lỡ tay import app (đụng tx) hoặc pgx, test này đỏ ngay.
func TestPosting_NoInfraImports(t *testing.T) {
	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		t.Fatalf("import dir posting: %v", err)
	}

	forbidden := []string{
		"github.com/jackc/pgx",          // SQL driver
		"github.com/danielgtaylor/huma", // HTTP framework
		"github.com/go-chi/chi",         // router
		"net/http",                      // HTTP
		"github.com/Maneva-AI/namviet-backend/internal/platform",            // mọi adapter hạ tầng
		"github.com/Maneva-AI/namviet-backend/internal/accounting/app",      // không đụng tx/use-case
		"github.com/Maneva-AI/namviet-backend/internal/accounting/internal", // không đụng adapter
	}

	all := append([]string{}, pkg.Imports...)
	all = append(all, pkg.TestImports...)

	for _, imp := range all {
		for _, bad := range forbidden {
			if imp == bad || strings.HasPrefix(imp, bad) {
				t.Errorf("posting import bị cấm %q (khớp %q) — template phải thuần", imp, bad)
			}
		}
	}
}
