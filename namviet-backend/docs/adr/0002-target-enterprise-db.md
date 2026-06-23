# ADR 0002 — Thiết kế DB target enterprise (giảm sửa DB về sau)

**Ngày:** 2026-06-17 · **Trạng thái:** Accepted · **Nguồn:** `docs/db-review-target-schema.md` (so sánh DB cũ ERP vs DB mới monorepo)

## Context
Review DB cũ (`nam_viet_erp/supabase/schema.sql`) vs DB mới (`office.duoc.namviet/database_schema.md`) phát hiện các khoảng trống enterprise (xem doc đầy đủ). Mục tiêu: thiết kế chuẩn để **ít phải nâng cấp/sửa DB về sau**. Lưu ý: `database_schema.md` chỉ liệt kê cột/kiểu/nullable/default — KHÔNG có FK/index/enum/constraint; 2 DB là 2 nhánh khác nhau. **Mọi thao tác đụng bảng cũ phải verify schema PROD thật trước.**

## Decision

### A. Object MỚI do backend sở hữu (schema `app`) — LÀM NGAY, additive
- **Kế toán double-entry chuẩn**: `journal_entries` (header) + `journal_entry_lines` (N dòng, **mỗi dòng debit XOR credit** bằng CHECK) + `account_balances` (projection) + `accounting_periods`. Cột `book` ∈ {INTERNAL, TAX}. Ép **Σdebit=Σcredit** trong transaction (constraint/trigger). **Append-only**, sửa = bút toán đảo. (Thay mô hình 1-dòng `accounting_journals` của DB hiện hữu.)
- **Tiền**: `NUMERIC` — scale-0 cho VND đã chốt, `NUMERIC(20,4)` cho đơn giá/VAT/proration. CẤM float (đã có lint). Round-trip test.
- **PK**: **uuid v7** cho mọi bảng mới (monotonic, hợp ledger/audit).
- **Concurrency**: `lock_version` (optimistic) + `updated_at` cho aggregate tiền/kho; SERIALIZABLE + retry-on-40001; claim nguyên tử (UPDATE có điều kiện).
- **Toàn vẹn miền**: `status`/loại dùng **enum PG hoặc CHECK** (không `text` tự do); CHECK `amount >= 0`/`qty >= 0` nơi hợp lý; FK đầy đủ.
- **VAT**: `sales_invoices` header + lines (module vat).
- **Quy mô**: keyset/cursor pagination + index phù hợp; **thiết kế partition-ready** (theo thời gian) cho ledger/transaction lớn ngay từ đầu.

### B. Bảng `public.*` hiện hữu — MIGRATE SAU (breaking, cần kế hoạch + verify prod)
- Thêm CHECK/enum-CHECK; thêm `lock_version`/`updated_at` cho bảng tiền/kho; thống nhất kiểu `quantity` (int vs numeric); tách `customers.b2b_metadata jsonb` → quan hệ chuẩn hoá (debt_limit/payment_term/sales_staff) như `customers_b2b` của ERP; bổ sung `orders.paid_amount` nếu DB mới thiếu (chống ghost debt).
- **KHÔNG đổi PK `bigint`→`uuid`** ở bảng cũ (phá FK toàn hệ). Giữ nguyên, chỉ chuẩn hoá quanh nó.

### C. Strangler — ĐỌC AS-IS
products/orders/inventory_batches/finance_transactions/customers... giữ nguyên; Go đọc/ghi qua **sqlc reference schema** `db/schema/public_*.sql` (ADR 0001).

## Consequences
- (+) Lõi tiền/kế toán đúng chuẩn ngay từ phần mới → ít phải sửa DB lớn về sau; ranh giới rõ `app.*` (mới, chuẩn) vs `public.*` (kế thừa, migrate dần).
- (−) Tồn tại song song 2 mô hình kế toán (1-dòng cũ ở public vs journal+lines mới ở app) trong giai đoạn chuyển — cần kế hoạch migrate + reconcile khi cắt sang.
- (−) Các quyết định "migrate sau" phụ thuộc **verify schema prod thật** (chưa có psql/MCP) — coi là tiền đề bắt buộc.
