---
name: new-migration
description: Use when adding or changing a Postgres migration in the Go backend — creates a goose migration in schema app with correct naming, regenerates sqlc, and adds an Up/Down test
---

# New Migration — Nam Việt Go (goose + sqlc)

Tạo/sửa migration AN TOÀN cho backend Go. Đọc `ARCHITECTURE.md` §6.

## Arguments
- `description` (bắt buộc): mô tả ngắn snake_case, ví dụ `add_orders` / `journal_entries`.

## CHECKLIST
### 1. Tạo file goose đúng tên
`db/migrations/<NNNNN>_<description>.sql` — số tăng dần (tiếp số lớn nhất hiện có) HOẶC `^\d{14}_[a-z0-9_]+\.sql$`. KHÔNG có ký tự thừa (1 chữ "a" thừa từng làm CLI skip silent).
```sql
-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;
-- ... DDL trong schema app ...

-- +goose Down
-- ... rollback tương ứng ...
```

### 2. Quy tắc nội dung
- Mọi object ở schema **`app`** (KHÔNG đụng `public` của Supabase khi chưa cần).
- **Tiền = `NUMERIC`** (không float/bigint tuỳ tiện). PK ưu tiên `uuid` (gen_random_uuid / uuid v7).
- Bảng tài chính: append-only; ràng buộc bất biến (`Σdebit=Σcredit`) bằng constraint/trigger nếu thuộc ledger. KHÔNG cho phép UPDATE/DELETE dòng tài chính ở tầng app.
- KHÔNG `DROP`/`TRUNCATE` bảng có dữ liệu trong migration app. Bảng UPDATE/DELETE cần PK/REPLICA IDENTITY (chuẩn bị cutover Phase 2).
- **Chuẩn enterprise bảng MỚI (ADR 0002):** PK **uuid v7**; tiền `NUMERIC` scale-0 (VND)/`(20,4)`; `status`/loại = **enum hoặc CHECK** (không text tự do); thêm `lock_version`+`updated_at` cho aggregate tiền/kho; CHECK `amount/qty >= 0`; FK đầy đủ; thiết kế **partition-ready** cho ledger/transaction lớn; kế toán = `journal_entries`+`journal_entry_lines` (mỗi dòng debit XOR credit, ép Σ trong tx) — KHÔNG mô hình 1-dòng.
- **Đụng bảng `public.*` cũ (migrate, breaking):** verify schema PROD thật trước; KHÔNG đổi PK bigint→uuid bảng cũ.

### 3. Sinh lại sqlc
- Thêm query vào `db/queries/<module>.sql` rồi `make sqlc` (`sqlc generate`). `sqlc diff` phải exit 0.

### 4. Test Up/Down (BẮT BUỘC)
Có/regression test migrate Up rồi Down trên testcontainers (mẫu: `db/migrations_test.go`). Chạy `go test ./db/ -run Migrations`.

### 5. Verify
Skill `verify-backend` xanh trước khi commit. Commit ngắn gọn tiếng Việt, không push.

## Anti-patterns (CẤM)
| SAI | ĐÚNG |
|-----|------|
| tên file sai pattern (ký tự thừa) | `NNNNN_snake_case.sql` sạch |
| float cho tiền | NUMERIC |
| DROP/TRUNCATE dữ liệu | thêm cột/migrate dữ liệu an toàn |
| đụng schema `public` Supabase | schema `app` |
| quên Down / quên test | có Up+Down + test |
