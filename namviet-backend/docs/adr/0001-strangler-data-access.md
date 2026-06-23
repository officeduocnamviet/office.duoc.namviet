# ADR 0001 — Data access trong giai đoạn strangler-fig

**Ngày:** 2026-06-17 · **Trạng thái:** Accepted

## Context
Backend Go thay dần Supabase functions theo strangler-fig, KẾT NỐI cùng Postgres của Supabase. Dữ liệu nghiệp vụ (products, orders, inventory, journal...) đã tồn tại trong schema `public`. Nếu backend tạo bảng `app` rỗng song song, API sẽ trả rỗng — sai (vi phạm "không mock/empty data sản phẩm"). Đồng thời backend có object MỚI của riêng nó (auth users/roles, refresh_tokens, idempotency_keys) không nên nhét vào `public`.

## Decision
- **Đọc/ghi dữ liệu nghiệp vụ hiện hữu ở schema `public`** (không tạo lại). sqlc dùng **schema tham chiếu** `db/schema/public_<area>.sql` mô tả các bảng `public.*` liên quan — CHỈ để type-check/codegen, KHÔNG `CREATE` qua goose lên prod.
- **Object MỚI do backend sở hữu** → schema `app` qua goose migration (như `identity`).
- **Integration test**: harness `dbtest` apply schema tham chiếu (`public.*`) + seed dữ liệu test trong testcontainers (vì DB test trống).
- **Nguồn schema tham chiếu**: tạm lấy từ `database_schema.md`/models của `namviet-backend-core` cũ; **bắt buộc verify lại với PROD** (pg_proc + REST `Object.keys`) khi có creds — prod hay lệch migration/core cũ.

## Consequences
- (+) API trả dữ liệu thật ngay; cutover từng route an toàn; không trùng bảng.
- (+) Ranh giới rõ: `public.*` = dữ liệu kế thừa (đọc/ghi qua sqlc reference); `app.*` = backend sở hữu.
- (−) Phải duy trì schema tham chiếu cho `public.*` và giữ nó khớp prod (rủi ro lệch → verify bắt buộc).
- (−) Ghi vào `public.*` phải tôn trọng trigger/ràng buộc hiện hữu (cascade ngầm) — khi port domain ghi tiền, tái hiện side-effect thành saga tường minh (xem spec §7.4), không dựa trigger.
