# Nam Việt Backend (Go)

Modular monolith (walking skeleton — Phase 0) thay thế dần Supabase functions theo
chiến lược strangler-fig. Xem spec:
`../docs/superpowers/specs/2026-06-15-go-backend-migration-design.md` và plan
`../docs/superpowers/plans/2026-06-15-go-backend-foundation.md`.

Module path: `github.com/Maneva-AI/namviet-backend`. Go 1.26.

## Bố cục

```
cmd/api/                 entrypoint mỏng (Mat Ryer idiom: main -> run)
internal/platform/
  config/                Load(getenv) -> Config
  server/                chi router + /v1/health (Deps{Pool})
  httpx/                 envelope {data,error} tương thích safeRpc FE
  db/                    pgxpool + codec NUMERIC<->shopspring/decimal
  logging/               slog JSON (+ otelslog fanout khi bật OTel)
  telemetry/             OTLP gRPC tracer provider (no-op khi endpoint rỗng)
  idempotency/           Store interface + middleware + pgxstore
db/migrations/           goose *.sql (schema app)
db/queries/ + sqlc.yaml  scaffold sqlc (CHƯA generate)
deploy/                  Dockerfile + docker-compose + Caddy + otel-collector
.github/workflows/ci.yml CI: verify decimal lib, build, test -race, govulncheck, lint
.golangci.yml            forbidigo cấm float ở money path
```

## Dev

- `make run` — chạy server (đọc env, xem `internal/platform/config`). Mặc định `:8080`.
- `make build` — build binary ra `bin/api`.
- `make test` — unit + integration (`-race`). Integration cần Docker (testcontainers).
- `make test-short` — chỉ unit, bỏ qua integration (`-short`).
- `make lint` / `make vuln` — chất lượng + bảo mật supply-chain.

## Trạng thái Phase 0 — phần ĐÃ HOÃN (build-xanh là ưu tiên)

- **Huma v2 / OpenAPI (Task 8 của plan):** HOÃN. `/v1/health` hiện dùng chi thuần
  (vẫn đúng envelope `{data,error}`). Lý do: tránh rủi ro build đỏ / phụ thuộc
  thêm; sẽ tích hợp ở plan sau. `make openapi` vì vậy chưa khả dụng.
- **sqlc generate / code sinh:** HOÃN — `sqlc.yaml` + `db/queries/idempotency.sql`
  chỉ là scaffold, KHÔNG chạy `sqlc generate`, KHÔNG import code sinh vào build.
  Idempotency store dùng Postgres được hiện thực bằng raw SQL trong
  `internal/platform/idempotency/pgxstore.go` (không phụ thuộc code sinh).
- **Task 0 (dump pg_proc PROD):** bỏ qua theo yêu cầu; chỉ giữ `.gitignore`.
