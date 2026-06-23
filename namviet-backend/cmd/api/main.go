// Command api là entrypoint mỏng theo idiom Mat Ryer: main() chỉ gọi run().
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Maneva-AI/namviet-backend/internal/accounting"
	"github.com/Maneva-AI/namviet-backend/internal/catalog"
	"github.com/Maneva-AI/namviet-backend/internal/customers"
	"github.com/Maneva-AI/namviet-backend/internal/finance"
	"github.com/Maneva-AI/namviet-backend/internal/identity"
	"github.com/Maneva-AI/namviet-backend/internal/inventory"
	"github.com/Maneva-AI/namviet-backend/internal/orders"
	"github.com/Maneva-AI/namviet-backend/internal/purchasing"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/config"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db"
	"github.com/Maneva-AI/namviet-backend/internal/platform/logging"
	"github.com/Maneva-AI/namviet-backend/internal/platform/server"
	"github.com/Maneva-AI/namviet-backend/internal/platform/telemetry"
	"github.com/Maneva-AI/namviet-backend/internal/vat"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "namviet-backend"

func main() {
	if err := run(context.Background(), os.Args[1:], os.Getenv, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, getenv func(string) string, stdout io.Writer) error {
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	dumpOpenAPI := fs.Bool("dump-openapi", false, "in OpenAPI 3.1 (YAML) ra stdout rồi thoát")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// -dump-openapi: dựng API không cần DB, in spec YAML rồi thoát. Dùng để sinh
	// artifact api/openapi.yaml và TS client.
	if *dumpOpenAPI {
		return dumpSpec(stdout)
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(getenv)
	if err != nil {
		return err
	}

	// Telemetry: no-op khi OTLP_ENDPOINT rỗng.
	shutdown, err := telemetry.Init(ctx, serviceName, cfg.OTLPEndpoint)
	if err != nil {
		return err
	}
	defer func() { _ = shutdown(context.Background()) }()

	logger := logging.New(serviceName, cfg.OTLPEndpoint != "")
	slog.SetDefault(logger)

	// Nạp cặp khoá JWT (ES256). Rỗng → ephemeral (cảnh báo: dev-only).
	keys, ephemeral, err := authn.LoadKeyPair(cfg.JWTPrivateKeyPEM, cfg.JWTPublicKeyPEM)
	if err != nil {
		return err
	}
	if ephemeral {
		logger.WarnContext(ctx, "JWT key chưa cấu hình — dùng khoá EPHEMERAL (chỉ dev; token mất hiệu lực khi restart)")
	}

	// Nối DB chỉ khi có DATABASE_URL.
	var pool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err = db.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
	}

	router := server.NewRouter(server.Deps{
		Pool:    pool,
		Modules: buildModules(pool, keys),
	})
	handler := otelhttp.NewHandler(router, "http")
	srv := server.NewHTTPServer(cfg.HTTPAddr, handler)

	logger.InfoContext(ctx, "khởi động", "addr", cfg.HTTPAddr, "env", cfg.AppEnv, "db", pool != nil)
	return server.Run(ctx, srv)
}

// buildModules wiring route của từng bounded context thành closure cho server.
// Đây là composition root: cmd/api biết mọi module; server không phụ thuộc
// ngược. pool có thể nil (vd khi dump-openapi) — register chỉ dựng schema, không
// chạy query.
func buildModules(pool *pgxpool.Pool, keys authn.KeyPair) []func(huma.API) {
	issuer := identity.NewTokenIssuer(keys.Private)
	verifier := authn.NewVerifier(keys.Public)
	authSvc := identity.New(pool, issuer)
	catalogSvc := catalog.New(pool)
	customersSvc := customers.New(pool)
	inventorySvc := inventory.New(pool)
	accountingMod := accounting.NewModule(pool)
	financeMod := finance.NewModule(pool)
	vatMod := vat.NewModule(pool)
	// orders ORCHESTRATION (P4b) cần port nội bộ của inventory/accounting/vat/finance
	// → dựng chúng TRƯỚC rồi inject. orders gọi các port này trong CÙNG tx giao dịch
	// để gộp atomic (trừ kho + post sổ + HĐ + phiếu thu).
	ordersSvc := orders.New(pool, orders.Deps{
		Deductor: inventory.NewDeductor(pool),
		Poster:   accountingMod.Poster(),
		Issuer:   vatMod.IssuePort(),
		Recorder: financeMod.RecordPort(),
	})
	// purchasing (mục 54 — chiều MUA): nhập kho (inventory.StockIn) + post sổ
	// (accounting.Poster) + chi NCC (finance.RecordOutPort) trong CÙNG tx (gộp atomic).
	purchasingSvc := purchasing.New(pool, purchasing.Deps{
		StockIn: inventory.NewStockInner(pool),
		Poster:  accountingMod.Poster(),
		Payer:   financeMod.RecordOutPort(),
	})

	return []func(huma.API){
		func(api huma.API) { identity.RegisterRoutes(api, authSvc, verifier) },
		func(api huma.API) { catalog.RegisterRoutes(api, catalogSvc, verifier) },
		func(api huma.API) { customers.RegisterRoutes(api, customersSvc, verifier) },
		func(api huma.API) { inventory.RegisterRoutes(api, inventorySvc, verifier) },
		func(api huma.API) { orders.RegisterRoutes(api, ordersSvc, verifier) },
		func(api huma.API) { purchasing.RegisterRoutes(api, purchasingSvc, verifier) },
		func(api huma.API) { accountingMod.RegisterRoutes(api, verifier) },
		func(api huma.API) { financeMod.RegisterRoutes(api, verifier) },
		func(api huma.API) { vatMod.RegisterRoutes(api, verifier) },
	}
}

// dumpSpec dựng huma.API (không cần DB) và in OpenAPI 3.1 dạng YAML. Dùng khoá
// ephemeral + pool nil chỉ để sinh schema /v1/auth/*.
func dumpSpec(w io.Writer) error {
	r := chi.NewRouter()
	keys, _, err := authn.LoadKeyPair("", "")
	if err != nil {
		return err
	}
	api := server.API(r, server.Deps{Modules: buildModules(nil, keys)})
	spec, err := api.OpenAPI().YAML()
	if err != nil {
		return fmt.Errorf("openapi yaml: %w", err)
	}
	// Header comment YAML hợp lệ: nêu rõ đây là spec OpenAPI 3.1 (Huma sort key
	// alphabet nên trường `openapi:` nằm sâu trong body).
	if _, err = io.WriteString(w, "# openapi: 3.1 — Nam Viet Backend API (sinh tự động, đừng sửa tay)\n"); err != nil {
		return err
	}
	_, err = w.Write(spec)
	return err
}
