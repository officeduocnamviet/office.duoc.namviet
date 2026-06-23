# Kiến trúc & Code Conventions — Nam Việt Backend (Go)

> Tài liệu chốt cách viết code. Mọi module/PR tuân theo. Nguồn: spec `docs/superpowers/specs/2026-06-15-go-backend-migration-design.md` + đã được CP1 (tầng API) chứng minh chạy thật. Cập nhật 2026-06-17.

## 0. Nguyên tắc (theo thứ tự ưu tiên)
1. **Chia bounded context (DDD) cho CHUẨN — ưu tiên số 1.** Ranh giới đúng quan trọng hơn mọi pattern.
2. **PRAGMATIC, KHÔNG over-engineer.** Không event-sourcing / CQRS / saga / outbox theo mặc định. Mặc định = gọi service qua port + **1 transaction**. Pattern nâng cao chỉ thêm KHI THẬT CẦN, có lý do rõ.
3. **Đúng đắn tiền/kế toán > mọi thứ.** Không `float` cho tiền; double-entry phải cân; bất biến ép ở DB. (Đây là chỗ DUY NHẤT được phép "kỹ" hơn mức tối thiểu.)
4. **Có `base`/`common` + pattern chuẩn.** Mọi module viết GIỐNG NHAU (cùng template) → dễ kiểm soát, dễ check lỗi, dễ mở rộng.
5. **TDD**: red → green → refactor. Unit + integration PASS cùng commit.
6. **Repo thuần Go** (không Node) · **không nợ kỹ thuật** (làm tới đâu xong tới đó).

---

## 1. Kiểu kiến trúc
**Modular monolith**: 1 Go module, 1 `cmd/api`, mỗi **bounded context** = 1 thư mục dưới `internal/`, hình hexagon **mỏng vừa đủ**. Tách service sau = đổi adapter.

## 2. Bounded contexts & Context Map (PHẦN QUAN TRỌNG NHẤT)

**Phase 1 — lõi tiền/vận hành (8 context):**

| Context (`internal/`) | Aggregate gốc | Sở hữu (bảng/khái niệm) | Mức độ "kỹ" |
|---|---|---|---|
| **identity** | User, Role | users, roles, permissions, refresh_tokens | vừa |
| **catalog** | Product | products, product_units (UOM), categories, manufacturers, prices, promotions | nhẹ (đa số read) |
| **customers** | Customer | customers, companies, customer_records, hồ sơ công nợ/hạn mức | vừa |
| **inventory** | Batch / StockItem | warehouses, batches (FEFO), product_inventory, movements, landed-cost | **kỹ** (concurrency) |
| **orders** | Order | orders, order_items; áp giá + voucher; 3 state machine | **kỹ** |
| **finance** | Payment / FundAccount | finance_transactions, fund_accounts, phân bổ thanh toán, công nợ, đối soát bank | **kỹ** |
| **accounting** | JournalEntry | chart_of_accounts, journal_entries + _lines, account_balances, periods; 2 sổ INTERNAL/TAX | **kỹ nhất** |
| **vat** | Invoice | sales_invoices, phát hành e-invoice, tính VAT | **kỹ** |

> `purchasing` (PO/supplier/costing) = Phase 1.5, nuôi `inventory.landed-cost`.
> **Defer (Phase 2):** HR (employees/payrolls/attendance/shifts) · clinical (medical_visits/queues/appointments) · AI (agent_workflows/ai_memories/vectors/chats) · marketing/notifications/training/uploads.

**Context map (quan hệ — gọi QUA PORT, không đụng bảng nhau):**
```
orders ──reads pricing──▶ catalog
orders ──checks debt───▶ customers
orders ──reserve/deduct▶ inventory
finance ─allocates────▶ orders (cập nhật paid)   finance ─updates debt─▶ customers
finance, orders, inventory ──post──▶ accounting (journal)   orders ─issue──▶ vat
mọi context ──authz──▶ identity
```
Quy tắc: **không FK chéo schema, không import repo của context khác.** Cần dữ liệu context khác → gọi **port interface** mà context đó export (vd `catalog.PriceQuery`, `inventory.StockReserver`). Quan hệ chéo dùng **ID**, không nhúng entity của nhau.

**Chia context "cho chuẩn":** mỗi context = một năng lực nghiệp vụ tự trị, có ngôn ngữ riêng (vd "Order" ở orders ≠ "Customer record" ở customers). Nếu hai thứ luôn thay đổi cùng nhau và cùng transaction → cùng context; nếu chỉ tham chiếu → tách context + port.

## 3. Layout chuẩn cho MỖI module (template — copy y hệt)
```
internal/<context>/
├── domain/            # THUẦN Go: entity, value object, business rule, PORT interface
│   ├── <aggregate>.go #   (KHÔNG import pgx/http/huma/framework)
│   └── ports.go       #   interface Repository + cổng ra context khác (domain ĐỊNH NGHĨA)
├── app/               # use-case; mở/commit TRANSACTION ở đây; gọi port
│   └── service.go
└── internal/          # adapters (compiler chặn module khác import)
    ├── postgres/       #   repo: sqlc-generated + map domain<->row
    └── http/           #   Huma handler + DTO + routes.go
```
Phụ thuộc: `adapters → app → domain`. Domain không biết SQL/HTTP. Port do **bên tiêu thụ** định nghĩa, adapter implement ("accept interfaces, return structs").
**Module nhẹ** (vd catalog read-only): được phép gộp `domain`+`app` nếu logic mỏng — đừng tạo lớp rỗng cho có (chống over-engineer). **Module "kỹ"** (money) thì đủ 3 lớp.

## 4. `base` / `common` / `platform` (nền dùng chung — để "có base, pattern chuẩn")
- **`internal/platform/`** — hạ tầng kỹ thuật, KHÔNG business: `server`, `config`, `db` (pgxpool+decimal codec + helper `WithinTx`), `httpx`+`httpx/humax` (envelope {data,error}), `idempotency`, `authn`, `authz`, `logging`, `telemetry`.
- **`internal/common/`** — **shared kernel** trung lập + tiện ích để mọi module dùng GIỐNG NHAU (không chứa business của context nào):
  - `money` — kiểu `Money` bọc `shopspring/decimal` (cộng/trừ/nhân, làm tròn VAT có quy tắc). **Bắt buộc dùng cho mọi tiền.**
  - `apperr` — taxonomy lỗi domain (`NotFound`, `Conflict`, `Validation`, `Forbidden`...) → map sang code envelope ở §6.
  - `pagination` (cursor/keyset), `id` (uuid v7 helper), `clock`, `valid` (helper validate dùng chung).
- **Pattern chuẩn lặp lại** (giúp dễ check lỗi/mở rộng): repo nhận `appdb.Queries` + `WithTx`; handler đăng ký qua `huma.Register` trong `routes.go`; lỗi luôn trả qua `apperr`; test có harness testcontainers dùng chung (`internal/platform/db/dbtest`).
> Ranh giới: `common` chỉ chứa thứ **trung lập domain** (Money, lỗi, phân trang). Logic nghiệp vụ KHÔNG bao giờ nằm ở `common` — tránh biến nó thành "god package".

## 5. API (đã chốt ở CP1)
chi v5 + **Huma v2** (code-first OpenAPI 3.1). `Input/Output{Body}` + `huma.Register(...)` trong `routes.go`. Envelope `{data,error}` qua `httpx/humax` (success: transformer bọc; error: `codeError`). Map: 400 `bad_request`, 401 `unauthorized`, 403 `forbidden`, 404 `not_found`, 409 `conflict`, 422 `validation_error`, 5xx `internal`. Validation bằng struct tag. **Idempotency-Key** cho POST/PATCH tiền. Pagination keyset. `/v1`. Action không-CRUD: `POST /v1/orders/{id}:cancel`. Sau khi đổi route: `make openapi`.

## 6. Data access
**pgx v5 + sqlc** (SQL check compile-time, KHÔNG ORM). Query ở `db/queries/<module>.sql`. **Transaction** mở/commit ở `app` qua helper `platform/db.WithinTx`, truyền `Queries.WithTx(tx)` xuống repo; domain không thấy tx. **Nhiều bước tiền = 1 transaction.** Migration goose schema `app`, tên `^\d{14}_[a-z0-9_]+\.sql$`, có test Up/Down.

**Strangler-fig data access (xem ADR 0001):** Backend đọc/ghi **bảng nghiệp vụ ĐANG TỒN TẠI ở schema `public`** của Supabase (products, orders, journal...). sqlc dùng **schema THAM CHIẾU** mô tả `public.*` (`db/schema/public_*.sql`, chỉ để codegen/type-check — KHÔNG phải migration tạo bảng). **CHỈ object MỚI do backend sở hữu** (idempotency_keys, refresh_tokens, app users/roles auth) nằm ở schema `app` qua goose. Integration test materialize schema tham chiếu + seed trong testcontainers. Schema tham chiếu lấy tạm từ `database_schema.md`/models của core cũ; **PHẢI verify lại với PROD thật** (pg_proc/pg_dump REST) trước khi tin — prod hay lệch migration.

## 7. Tiền & Double-entry (LÕI — bất khả nhân nhượng, đây là chỗ được "kỹ")
- Tiền = `NUMERIC` ↔ `common/money` (decimal). **CẤM `float`** ở money path (lint chặn). Test round-trip.
- Double-entry: `journal_entries` + `journal_entry_lines` (mỗi dòng debit XOR credit); **`Σdebit=Σcredit` ép TRONG DB tx**; entry **append-only**, sửa = **bút toán đảo**.
- **2 sổ INTERNAL/TAX tách biệt, KHÔNG sync.** Test riêng từng đường posting.
- Concurrency: `SERIALIZABLE` + retry-on-40001; `lock_version` optimistic; `FOR UPDATE` chỉ cho hot-row.
- Idempotency tiền: dedupe theo mã GD bất biến của bank. Pin đúng `shopspring/decimal` + `govulncheck`.
- **Bảng MỚI (schema `app`) chuẩn enterprise (xem `docs/db-review-target-schema.md` + ADR 0002):** PK **uuid v7**; tiền `NUMERIC` scale-0 (VND) / `NUMERIC(20,4)` (đơn giá/VAT/proration); `status`/loại = **enum PG hoặc CHECK** (không `text` tự do); `lock_version` + `updated_at` cho aggregate tiền/kho; CHECK `amount/qty >= 0`; **partition-ready** cho ledger/transaction lớn; kế toán = `journal_entries`+`journal_entry_lines` cân Σ (KHÔNG mô hình 1-dòng). Đụng bảng `public.*` cũ (migrate) phải **verify schema PROD** trước.

## 8. Liên-module, jobs, events (PRAGMATIC — chỉ thêm khi cần)
- **Mặc định**: orchestration gọi service context khác qua **port + 1 transaction**. Đơn giản, dễ debug.
- **River** (queue trên Postgres) cho việc nền (gửi mail/push, đối soát) — enqueue trong cùng tx. Periodic job thay `pg_cron`.
- **Domain event / outbox**: KHÔNG mặc định. Chỉ dùng khi một hành động phải kích nhiều side-effect bất đồng bộ và việc tách giảm rủi ro thật (vd thay cascade trigger phức tạp). Khi dùng: consumer idempotent.
- KHÔNG event-sourcing, KHÔNG CQRS trừ khi có nhu cầu read-model thật (vd báo cáo BCTC nặng) — và phải nêu lý do.

## 9. AuthN/AuthZ
Access JWT **ES256** TTL ngắn (pin alg); refresh opaque xoay vòng + reuse-detection. Hash **argon2id**; import bcrypt GoTrue (lazy rehash), không UPDATE `auth.users`. **RBAC table-driven, 1 enforcement point** (`platform/authz`). Bỏ RLS + per-scope guard.

## 10. Observability
OpenTelemetry + `slog` (otelslog, dùng `*Context`). Redaction PII ở Collector. **Telemetry ≠ audit trail** (audit = sổ cái + bảng audit trong PG).

## 11. Testing (TDD)
- Unit thuần domain (không DB) cho rule/invariant; integration (testcontainers postgres:18) cho repo/tx/concurrency.
- **Fakes > mocks** cho port; real cho SQL. **Architecture-fitness test** (depguard) chặn domain import hạ tầng — chạy CI (rẻ, giúp "dễ check lỗi").
- Property-based test (`rapid`): **chỉ** cho bất biến ledger (Σ cân, làm tròn VAT) — không ép mọi nơi.
- Mọi fix kèm unit+integration PASS cùng commit. Không write lên dữ liệu thật khi test.

## 12. Repo hygiene
Thuần Go (FE sinh TS từ `api/openapi.yaml`). `.gitattributes` EOL=LF. Commit ngắn gọn tiếng Việt, không `Co-Authored-By`. App DB role không DROP/TRUNCATE. Không commit secret.

## 13. Reference
ThreeDots `wild-workouts-go-ddd-example` (DDD/hexagon pragmatic) · Ben Johnson `wtf` (domain-first layout) · Martin Fowler *Accounting Patterns* / Modern Treasury / TigerBeetle (double-entry) · `riverqueue/river`.
> Core cũ `officeduocnamviet/office.duoc.namviet/namviet-backend-core` = **tham khảo breadth (42 feature) + shape bảng thật**, KHÔNG phải nền code.
