---
name: new-bounded-context
description: Use when scaffolding a new bounded-context module in the Nam Viet Go backend (catalog, customers, inventory, orders, finance, accounting, vat) — creates the hexagon domain/app/adapters structure following the identity template
---

# New Bounded Context — Nam Việt Go

Scaffold 1 module DDD theo **template `internal/identity/`**. Đọc `ARCHITECTURE.md` §3 trước.

## Arguments
- `context` (bắt buộc): tên context, ví dụ `orders`, `accounting`.
- `complexity` (tuỳ chọn): `light` (read-mostly, gộp domain+app) | `full` (đường tiền, đủ 3 lớp). Mặc định `full` cho money domain.

## CHECKLIST (không bỏ bước)

### 1. Tạo layout (copy từ identity)
```
internal/<context>/
├── module.go                 # composition root: New(...)→Service, RegisterRoutes(api, svc, deps)
├── domain/
│   ├── <aggregate>.go        # entity + value object (THUẦN Go, chỉ stdlib/context)
│   ├── ports.go              # interface Repository + cổng ra context khác (domain định nghĩa)
│   └── arch_test.go          # copy từ identity, đổi đường dẫn forbidden → ép domain thuần
├── app/
│   ├── service.go            # use-case; mở tx qua platform/db.WithinTx + TxManager
│   ├── ports.go              # TxManager + Repos struct
│   └── errors.go             # helper apperr cho context
└── internal/
    ├── postgres/
    │   ├── repo.go           # implement port domain bằng appdb (sqlc); map row<->domain; var _ domain.X = (*Repo)(nil)
    │   └── txmanager.go
    └── http/
        ├── routes.go         # Input/Output DTO + huma.Register
        └── errors.go         # apperr → humax envelope
```
`light`: bỏ thư mục `app/` riêng nếu logic mỏng — nhưng vẫn giữ `domain/` thuần + ports. Đừng tạo lớp rỗng cho có.

### 2. Migration + queries
- `new-migration` (skill) tạo `db/migrations/NNNNN_<context>.sql` (schema `app`).
- `db/queries/<context>.sql` → `make sqlc` → dùng `appdb.Queries`.

### 3. TDD (BẮT BUỘC)
- Domain: unit test rule/invariant (không DB).
- App: unit test use-case với **fake** port.
- Integration: qua harness `internal/platform/db/dbtest` (testcontainers) cho repo/tx/concurrency.
- Money domain: property test bất biến (`Σdebit=Σcredit`, làm tròn VAT) + test concurrency (claim nguyên tử, no double-deduct).
- Viết test FAIL trước → code → PASS.

### 4. Wire + contract
- `module.New(...)` ở `cmd/api/buildModules`; `RegisterRoutes` mount lên huma.API (server KHÔNG phụ thuộc ngược module).
- Route money POST/PATCH gắn Idempotency-Key. Sau đó `make openapi` cập nhật `api/openapi.yaml`.

### 5. Verify
Chạy skill `verify-backend` — phải xanh trước khi commit. Commit ngắn gọn tiếng Việt, không push.

## Anti-patterns (CẤM)
| SAI | ĐÚNG |
|-----|------|
| domain import pgx/huma | domain chỉ stdlib/context + `arch_test.go` chặn |
| float cho tiền | `common/money` (decimal) ↔ NUMERIC |
| raw SQL / GORM | sqlc |
| share bảng / FK chéo context | gọi qua port, truyền ID |
| event-sourcing/CQRS "cho chuẩn" | pragmatic: port + 1 transaction |
| lớp app/domain rỗng cho module nhẹ | gộp lại, đừng over-engineer |
