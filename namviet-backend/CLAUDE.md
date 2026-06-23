# Nam Việt Backend (Go) — Hướng dẫn cho Claude

> File này là **GUARDRAIL**: mọi task phải bám đúng design ban đầu. ĐỌC TRƯỚC khi làm bất cứ việc gì.
> Thiết kế chi tiết: [ARCHITECTURE.md](ARCHITECTURE.md) (conventions code) + [docs/superpowers/specs/2026-06-15-go-backend-migration-design.md](../docs/superpowers/specs/2026-06-15-go-backend-migration-design.md) (spec gốc).

## Dự án là gì
Backend Go thay thế toàn bộ Supabase functions cho hệ **chuỗi nhà thuốc Nam Việt**: orders, inventory (FEFO), finance/thanh toán, **kế toán sổ kép TT133** (2 sổ INTERNAL/TAX), VAT e-invoice. Phục vụ 3 FE (ERP + 2 portal B2B) qua REST.

## Quyết định kiến trúc đã CHỐT (không tự ý đổi — nếu cần đổi phải hỏi)
- **Modular monolith** Go (1 module, `cmd/api` + `internal/<context>` hexagon). Tách service sau = đổi adapter.
- **REST + Huma v2** (OpenAPI 3.1 code-first). Response **envelope `{data,error}`** cho cả success/error (tương thích `safeRpc` FE) qua `internal/platform/httpx/humax`.
- **pgx v5 + sqlc** (SQL check compile-time). **KHÔNG ORM** (cấm GORM).
- **goose** migration, schema `app`.
- **Auth trong Go**: JWT ES256 (pin alg) + refresh opaque xoay vòng + reuse-detection; argon2id + import bcrypt; **RBAC table-driven, 1 enforcement point** (`platform/authz`). Bỏ Supabase Auth/RLS.
- **Giữ Postgres**; data tạm ở Supabase, self-host là dự án Phase 2 (cutover logical replication).
- **Strangler-fig**: Go API mimic `{data,error}` + map tên-hàm→route → FE gần như không sửa.
- Repo đích: monorepo `officeduocnamviet/office.duoc.namviet` (thay `namviet-backend-core` cũ). Core cũ chỉ là **tham khảo** (bản đồ domain + shape bảng), KHÔNG phải nền code.

## RULE BẮT BUỘC (vi phạm = sai design)
1. **Pragmatic DDD, KHÔNG over-engineer.** Mặc định: gọi service qua port + 1 transaction. KHÔNG event-sourcing/CQRS/saga/outbox trừ khi có lý do rõ. Module nhẹ (catalog read) gộp domain+app; chỉ module tiền đủ 3 lớp.
2. **Domain THUẦN.** `internal/<ctx>/domain/` chỉ import `context`/stdlib, KHÔNG pgx/huma/http/platform. Mỗi context có `domain/arch_test.go` chặn import hạ tầng.
3. **Tiền = `common/money` (decimal) ↔ NUMERIC. CẤM `float`** ở money path (lint `forbidigo` chặn).
4. **Double-entry**: journal_entries + lines, `Σdebit=Σcredit` ép trong DB tx; append-only, sửa = bút toán đảo. **2 sổ INTERNAL/TAX KHÔNG sync.**
5. **TDD**: red→green→refactor. **Unit + integration (testcontainers) PASS cùng commit.** Không write side-effect lên dữ liệu thật khi test.
6. **Lỗi qua `common/apperr`** (Code ổn định) → map envelope ở `httpx/humax`. Không trả lỗi tuỳ tiện.
7. **Repo thuần Go** — KHÔNG Node/`node_modules`. Chỉ publish `api/openapi.yaml`; FE tự sinh TS client.
8. **Commit ngắn gọn tiếng Việt, KHÔNG `Co-Authored-By`.** App DB role không DROP/TRUNCATE. Không commit secret.
9. **KHÔNG push** lên remote cho tới khi có git access org `officeduocnamviet` (hiện gh = `Andrew-Tran-Maneva`, chưa có quyền). Build + commit local thôi.
10. **Verify trước khi nói xong**: `sqlc diff` + `go build` + `go vet` + `go test ./... -count=1 -p 1` (full, gồm integration) đều xanh. Dùng skill `verify-backend`.
11. **Giữ kiến trúc luôn sống.** `ARCHITECTURE.md` là source of truth — mọi thay đổi kiến trúc (ranh giới context, pattern, quyết định) PHẢI cập nhật `ARCHITECTURE.md` + ghi **ADR** ở `docs/adr/` (vai trò `architect`). Auto-memory chỉ giữ **con trỏ + quyết định chốt**, KHÔNG chép toàn văn (tránh stale). Khi recall memory phải verify lại với `ARCHITECTURE.md`.
12. **Thiết kế DB chuẩn enterprise (ADR 0002 + `docs/db-review-target-schema.md`).** Bảng MỚI ở `app`: PK uuid v7; tiền `NUMERIC` scale-0/`(20,4)`; `status` = enum/CHECK (không text tự do); `lock_version`+`updated_at` ở aggregate tiền/kho; CHECK amount/qty≥0; partition-ready cho ledger; kế toán = journal+lines cân Σ. **Migrate bảng `public.*` cũ** (thêm CHECK/enum, normalize b2b_metadata, lock_version, thống nhất quantity) là **breaking → phải verify schema PROD thật trước + có ADR/kế hoạch**; **KHÔNG đổi PK bigint→uuid** bảng cũ.

## Bounded contexts (Phase 1) — xem context map ở ARCHITECTURE.md §2
`identity` (mẫu, xong) · `catalog` · `customers` · `inventory` · `orders` · `finance` · `accounting` · `vat`. Defer Phase 2: HR/clinical/AI/marketing.
Mỗi module copy template của `identity`: `internal/<ctx>/{domain, app, internal/{postgres,http}}` + `module.go` (composition root).

## Vai trò (agents — dùng đúng người cho đúng việc)
| Agent | Khi dùng |
|-------|----------|
| `backend-engineer` | Viết/sửa code Go domain (hexagon, sqlc, Huma, TDD) |
| `frontend-engineer` | FE Next.js/React (shadcn-only, dùng OpenAPI client, không mock data) |
| `test-engineer` | Viết test (unit + integration testcontainers, property test ledger) |
| `code-reviewer` | Review trước merge: DDD purity, no-float, security, envelope |
| `architect` | Quyết ranh giới bounded context, context map, viết ADR |
| `pm` | Chia phase/task, acceptance criteria, ưu tiên |
| `ba` | Làm rõ nghiệp vụ pharma/TT133, requirement, edge case (tham khảo core cũ + DB) |
| `data-architect` | Schema/migration (goose), query sqlc, toàn vẹn dữ liệu, dual-ledger |

## Skills
- `new-bounded-context` — scaffold 1 module DDD theo template `identity`.
- `verify-backend` — chạy cổng verify đầy đủ (sqlc/build/vet/full-test).
- `new-migration` — tạo goose migration (schema `app`) + regen sqlc đúng quy tắc.

## Commands
```bash
make test          # go test ./... -race -count=1 -p 1 (full, gồm integration — cần Docker)
make test-short    # unit thôi (skip integration)
make build / vet / lint / vuln
make sqlc          # sqlc generate
make openapi       # sinh api/openapi.yaml (hợp đồng cho FE)
```

## Key files
| File | Vai trò |
|------|---------|
| `ARCHITECTURE.md` | Conventions code (BẮT BUỘC theo) |
| `internal/identity/` | Module MẪU — copy pattern theo đây |
| `internal/platform/` | Nền kỹ thuật (server, db, httpx/humax, authn, authz, idempotency, telemetry) |
| `internal/common/` | Shared kernel trung lập (money, apperr, id) |
| `db/migrations/` `db/queries/` | goose SQL + sqlc queries |
| `api/openapi.yaml` | Hợp đồng API (FE sinh TS client từ đây) |
