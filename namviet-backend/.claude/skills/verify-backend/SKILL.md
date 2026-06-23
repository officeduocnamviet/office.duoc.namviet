---
name: verify-backend
description: Use when about to claim the Go backend work is done, before committing, or when verifying a module/PR — runs the full quality gate (sqlc diff, build, vet, full test including integration) and confirms green before any success claim
---

# Verify Backend — Cổng chất lượng

Chạy ĐẦY ĐỦ cổng verify TRƯỚC khi nói "xong" / commit. Không claim done nếu chưa chạy + xanh thật (dán output).

## Cổng (chạy hết, theo thứ tự)
```bash
cd <root namviet-backend>
sqlc generate && sqlc diff   # generated khớp schema/queries (exit 0)
go build ./...               # compile sạch
go vet ./...                 # vet sạch
go test ./... -count=1 -p 1  # FULL test gồm integration (testcontainers, cần Docker)
```
- `-p 1`: serialize package để testcontainers ổn định trên Docker Desktop/Windows.
- Nếu đổi route: `make openapi` để đồng bộ `api/openapi.yaml`.

## Điều kiện PASS
- 4 lệnh đều exit 0; test integration **chạy thật** (không `-short`, không skip), mọi package `ok`.
- Bằng chứng: dán dòng `ok ...` của các package integration (identity, db, idempotency...).

## CẤM
- Claim "xong/chạy được" khi mới chạy `-short` (integration bị skip) hoặc chưa chạy.
- Bỏ qua `sqlc diff` (generated lệch → CI/đối tác fail).
- "Tin báo cáo của subagent" mà không tự verify với code đụng tiền/bảo mật.
