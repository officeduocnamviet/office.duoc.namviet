# DB Design Review — Nam Việt (DB cũ ERP vs DB mới monorepo) → Target Schema Enterprise

> Mục tiêu: so sánh 2 nguồn schema, chỉ ra khoảng trống thiết kế cấp enterprise, và đề xuất **target schema** để hệ ÍT phải sửa/nâng cấp DB về sau. **Read-only review** — không sửa DB.
> Người soạn: vai trò `data-architect`. Ngày: 2026-06-22.
> Chuẩn đích tham chiếu: `namviet-backend/ARCHITECTURE.md` §6/§7, spec `docs/superpowers/specs/2026-06-15-go-backend-migration-design.md` §7.

---

## 1. Tổng quan, phương pháp & giới hạn nguồn

**Hai nguồn được so sánh:**

| Nhãn | Đường dẫn | Bản chất | Có FK/Index/Constraint? |
|---|---|---|---|
| **DB cũ (ERP)** | `nam_viet_erp/supabase/schema.sql` | DDL hợp nhất pg_dump — tables, cột, **enum types**, FK, index (trgm/gin), CHECK, generated columns, RLS, triggers, functions | **CÓ đầy đủ** |
| **DB mới (doc)** | `office.duoc.namviet/database_schema.md` | Export thô dạng doc: chỉ `Table / col: type (Nullable) DEFAULT` | **KHÔNG có** FK/index/constraint/enum → "không xác định được từ doc" |

**Phương pháp:** đối chiếu từng bảng core (products, orders, order_items, batches/inventory, finance_transactions, fund_accounts, chart_of_accounts/accounting_journals, customers/customers_b2b, sales_invoices/finance_invoices), soi tiền/kế toán/toàn vẹn/concurrency, rồi map vào chuẩn đích Go backend.

**Giới hạn quan trọng (đọc trước khi tin báo cáo):**
1. DB mới là **doc-only**: thiếu FK/index/constraint/enum → nhiều nhận định "DB mới thiếu X" thực ra là **"doc không thể hiện X"**, KHÔNG khẳng định prod thiếu. Phải verify với prod (`pg_dump`/`pg_proc`/REST `Object.keys`) trước khi hành động — đúng tinh thần memory *"Verify schema/route prod trước khi nói OK"*.
2. Hai nguồn **không phải hai phiên bản của cùng một DB**: DB cũ là **ERP nội bộ** (HR/clinical/asset/PO/supplier/finance_invoices/customers_b2b…), DB mới (doc) thiên về **vận hành + AI + clinical**, dùng `accounting_journals` thay vì `chart_of_accounts`-only. Khác biệt phản ánh **hai nhánh sản phẩm/thời điểm**, không phải migration tuyến tính. → "diff" dưới đây mang tính **so sánh hai mô hình**, không phải changelog.
3. DB cũ `schema.sql` **chất lượng kỹ thuật cao hơn** doc DB mới ở mặt thể hiện (có enum, generated col, CHECK, partial unique index). Đây là baseline tốt để giữ khi strangler.

---

## 2. Diff cấu trúc (so sánh hai mô hình)

### 2.1 Bảng CHỈ có ở DB mới (doc) — không thấy ở ERP schema.sql
| Bảng | Nhóm | Ảnh hưởng / Ghi chú |
|---|---|---|
| `accounting_journals` | Kế toán | **Khác mô hình** — DB cũ chỉ có `chart_of_accounts`, không có bảng journal. Đây là 1-dòng-debit/credit (xem §3 — KHÔNG phải double-entry chuẩn) |
| `agent_workflows`, `ai_agent_memories`, `medical_knowledge_vectors`, `product_vectors`, `chat_sessions`, `chat_messages` | AI/Chatbot | Defer Phase 2. `embedding USER-DEFINED` = pgvector → cần extension + index ivfflat/hnsw (doc không thể hiện) |
| `appointments`, `clinical_queues`, `medical_visits`, `customer_vaccination_records` | Clinical | Defer Phase 2 |
| `attendance_logs`, `employment_contracts`, `payrolls`, `payroll_items(_v2)`, `shift_assignments`, `shift_handovers`, `work_shifts` | HR | Defer Phase 2 |
| `companies`, `marketing_campaigns`, `promotions`, `customer_vouchers`, `approval_requests/steps`, `assets` | Hỗ trợ | Defer / hoặc map vào catalog/customers |
| `system_audit_logs` + `_2026_06`, `_2026_07` | Audit | **DB mới ĐÃ partition theo tháng** (range partition) — điểm cộng future-proof (xem §3 Partition) |
| `book_type`, `is_posted` trên `finance_transactions` | Kế toán/Finance | DB mới có 2-sổ flag ở finance; DB cũ **chưa có** `book_type`/`is_posted` (theo schema.sql) |

### 2.2 Bảng CHỈ có ở DB cũ (ERP) — không thấy ở doc DB mới
| Bảng | Ảnh hưởng |
|---|---|
| `customers_b2b`, `customer_b2b_contacts` | **CAO** — Portal B2B phụ thuộc. Doc DB mới gộp B2B vào `customers.b2b_metadata jsonb` + `customer_type` → **mất quan hệ chuẩn hoá** (debt_limit, payment_term, sales_staff_id, bank…). Xem §3 JSONB |
| `finance_invoices`, `finance_invoice_allocations` | **CAO** — hoá đơn NCC + phân bổ vào PO. Có CHECK `allocated_amount >= 0` |
| `purchase_orders`, `purchase_order_items`, `suppliers`, `inventory_receipts`, `inventory_receipt_items` | **CAO** — toàn bộ purchasing/nhập kho (Phase 1.5 trong spec) |
| `transaction_categories`, `banks`, `shipping_partners/rules`, `service_packages/items`, `prescription_templates` | TB | Một số đã có ở DB mới dưới tên khác |
| Enum types: `account_type`, `account_balance_type`, `transaction_flow`, `business_type`, `order_status`, `stock_management_type`, `customer_b2c_type`, `fund_account_type`… | **CAO** | DB mới (doc) khai báo các cột này là **`text` tự do** → mất ràng buộc miền giá trị. Xem §3 Toàn vẹn |

### 2.3 Cột/kiểu/nullable đổi ở bảng trùng tên (core)

| Bảng | DB cũ (schema.sql) | DB mới (doc) | Mức |
|---|---|---|---|
| **products** | PK `bigint`; `status text`; **không** category_name/manufacturer_name cache (chỉ có manufacturer_name, category_name ở cuối); có generated `fts tsvector`; CHECK purchasing_policy/items_per_carton | PK `bigint`; thêm `category_id`/`manufacturer_id` **lẫn** cache `category_name`/`manufacturer_name`; thêm `usage_instructions jsonb`, `product_images ARRAY`, `stock_status`, margin_rate | TB — cache cột (xem §3 Denormalize) |
| **orders** | PK `uuid`; `status` → **enum `order_status`** (DRAFT/QUOTE/CONFIRMED…); có `paid_amount`, `shipping_fee`, `discount_amount`, `quote_expires_at`, `delivery_*` | PK `uuid`; `status text` tự do + `payment_status text`; **thiếu `paid_amount`** trong doc | **CAO** — mất enum + (doc) thiếu paid_amount ⇒ "ghost debt" như memory `project_b2b_debt_display_source` cảnh báo |
| **order_items** | PK `uuid`; `base_quantity` & `total_line` = **GENERATED ALWAYS STORED**; CHECK qty>0, unit_price>=0 | PK `uuid`; `base_quantity`/`total_line` là cột **thường** (doc không nói generated); thêm `quantity_picked`, `quantity_returned`, `conversion_factor` | TB — DB cũ ép tính toán ở DB (tốt hơn) |
| **batches** | PK `bigint`; cột tối thiểu (product_id, batch_code, expiry, mfg, inbound_price) | thêm `updated_at`, `deleted_at` (soft delete) | Thấp |
| **inventory_batches** | `quantity integer` | `quantity numeric` | TB — kiểu lệch (xem §3) |
| **product_inventory** | `stock_quantity integer`; ít cột vị trí | `stock_quantity numeric`; thêm location_cabinet/row/slot, updated_by | TB |
| **finance_transactions** | PK `bigint`; enum `flow`/`business_type`/`status`; CHECK amount>0; `created_by DEFAULT auth.uid()` | PK `bigint`; `flow/business_type/status text`; thêm `book_type`, `is_posted`, `bank_reference_id`, `target_bank_info jsonb`, `ref_advance_id` | TB — DB mới giàu hơn về finance flags; nhưng mất enum + (doc) **không thấy CHECK amount>0** |
| **fund_accounts** | PK `bigint`; enum `type`/`status`; `balance numeric` | thêm `account_id text` (map tài khoản kế toán) | Thấp |
| **customers** | PK `bigint`; `type` enum `customer_b2c_type`; partial **unique index** (phone CaNhan / name ToChuc); tax_code, allergies, medical_history | PK `bigint`; `customer_type text DEFAULT 'B2C'`; `b2b_metadata jsonb`, `current_debt numeric`; **không** thấy unique index (doc) | **CAO** — gộp B2B vào jsonb + mất unique |
| **chart_of_accounts** | `type`/`balance_type`/`status` = **enum**; `parent_id` self-FK; updated_at trigger | `type`/`balance_type text`; thêm `deleted_at` | TB |
| **sales_invoices** | **Không có trong schema.sql** (ERP có `finance_invoices` cho NCC) | **Không có trong doc DB mới** | — sales_invoices là **bảng MỚI** mà spec yêu cầu cho module `vat` (chưa tồn tại ở cả hai) |

---

## 3. Vấn đề thiết kế & khoảng trống enterprise (kèm severity)

### 3.1 Tiền (Money) — **HIGH**
- DB cũ: **nhất quán `numeric`** cho mọi cột tiền (amount, balance, unit_price, actual_cost, allocated_amount…). **Không thấy float/double/real** ở money path. ✅ baseline tốt.
- **Vấn đề: `numeric` KHÔNG khai scale** (raw `numeric`, không `numeric(p,s)`). Chuẩn đích spec §7 yêu cầu **`NUMERIC` scale-0 cho VND chốt** + **`NUMERIC(20,4)` cho đơn giá/VAT/proration**. → Hiện tại scale không bị ép ⇒ rủi ro lệch làm tròn khi nhiều tầng tính.
- **`inventory_batches.quantity`/`product_inventory.stock_quantity` lệch kiểu** giữa `integer` (cũ) và `numeric` (doc mới) → quyết định 1 kiểu (thuốc lẻ có thể phân số → `numeric`).
- **Severity HIGH** vì là money path; nhưng **không phát hiện float** nên không Critical.

### 3.2 Kế toán / Double-entry — **CRITICAL**
- DB mới `accounting_journals` = **1 dòng** `(account_debit, account_credit, amount)` → đây là mô hình **"single-row T-account"**, **KHÔNG phải double-entry chuẩn**. Hệ quả:
  - Không biểu diễn được bút toán **nhiều nợ / nhiều có** (1 phiếu chi phân bổ nhiều TK).
  - **Không có ràng buộc Σdebit = Σcredit** ở cấp chứng từ (mỗi dòng tự cân vì chỉ 2 vế) — nhưng mất khả năng kiểm Σ khi mở rộng.
  - **Thiếu hẳn chiều sổ INTERNAL/TAX** ở `accounting_journals` (cột `book_type` chỉ nằm ở `finance_transactions`). Spec §7.2 bắt buộc **2 sổ tách biệt KHÔNG sync**.
  - Không có `account_balances`, `accounting_periods` (period close, carry-forward 5xx/6xx/7xx/8xx→911→4212) — spec yêu cầu.
- DB cũ thậm chí **không có bảng journal nào** (chỉ `chart_of_accounts`). Memory `project_accounting_module_build` xác nhận hệ hạch toán đang build, **chưa deploy prod**.
- **Đề xuất mô hình target (làm MỚI ở schema `app`, không đụng bảng cũ):**
  ```
  app.journal_entries(id uuid v7 PK, book_type text CHECK in('INTERNAL','TAX'),
      entry_date date, period_id, doc_type, source_ref_type, source_ref_id,
      description, posted_by, posted_at, status CHECK in('draft','posted','reversed'),
      reversed_by_entry_id, lock_version int, created_at)
  app.journal_entry_lines(id uuid v7 PK, entry_id FK, account_code,
      debit NUMERIC(20,4) DEFAULT 0, credit NUMERIC(20,4) DEFAULT 0,
      CHECK ((debit=0) <> (credit=0)),  -- debit XOR credit
      line_no int)
  -- Ràng buộc Σdebit=Σcredit ép TRONG DB tx (deferred constraint trigger hoặc check ở app+trigger)
  -- append-only: cấm UPDATE/DELETE entry đã posted; sửa = bút toán đảo
  app.account_balances(account_code, book_type, period_id, debit_total, credit_total, lock_version)
  app.accounting_periods(id, year, month, status CHECK in('open','closed'), closed_at)
  ```
- **Severity CRITICAL** — đây là rủi ro #1 spec nêu (routing book sai + invariant ledger).

### 3.3 Khóa chính lẫn lộn bigint vs uuid — **MED**
- DB cũ + mới đều lẫn: `products/customers/categories/batches/fund_accounts/finance_transactions = bigint` (identity sequence); `orders/order_items/chart_of_accounts/accounting/appointments = uuid`.
- Hệ quả: FK chéo phải dùng cả `bigint` và `uuid`; bảng mới đẻ ra dễ chọn sai kiểu.
- **Đề xuất:** **bảng MỚI do backend sở hữu (schema `app`) → uuid v7** (helper `common/id`, time-ordered, index-friendly thay vì `gen_random_uuid` v4 ngẫu nhiên gây page-split). **PK cũ giữ nguyên** (bigint sequence) — KHÔNG phá, chỉ tham chiếu bằng ID qua port.
- **Severity MED** (không gây lỗi, chỉ nợ nhất quán).

### 3.4 Soft-delete & audit cols — **MED**
- `deleted_at`: DB mới (doc) khá nhất quán (hầu hết bảng nghiệp vụ có). DB cũ **không đều** (orders/order_items/finance_transactions schema.sql **không** có `deleted_at`; products cũ không có nhưng doc mới có).
- Audit cols: `created_at` phổ biến; `updated_at` thiếu ở vài bảng (batches cũ, order_items cũ); `created_by`/`updated_by` **rời rạc** (products có `updated_by`, orders không có `updated_by`).
- **Đề xuất:** chuẩn hoá **mixin audit** cho bảng MỚI: `created_at/updated_at NOT NULL DEFAULT now()`, `created_by/updated_by uuid`, `deleted_at` (chỉ nơi cần soft-delete). Bảng cũ giữ nguyên (đọc-as-is). Đã có `system_audit_logs` (old/new jsonb) làm audit trail tách biệt telemetry (đúng ARCHITECTURE §10).
- **Severity MED.**

### 3.5 Toàn vẹn: FK / CHECK / UNIQUE / NOT NULL / enum — **HIGH**
- **Enum → text tự do (regression):** DB cũ dùng **enum PG** cho `order_status`, `account_type`, `transaction_flow`, `business_type`, `fund_account_type`, `customer_b2c_type`, `stock_management_type`… DB mới (doc) hạ xuống **`text` + DEFAULT** ⇒ mất ràng buộc miền. Spec không bắt enum PG (Go enforce + envelope), nhưng **tối thiểu cần CHECK constraint** ở các cột status/flow/type để DB tự bảo vệ.
- **CHECK tiền:** DB cũ có `finance_transactions.amount > 0`, `order_items.quantity > 0 / unit_price >= 0`, `allocated_amount >= 0`. Doc DB mới **không thể hiện** các CHECK này (giới hạn nguồn) — phải verify prod; nếu thiếu → bổ sung.
- **UNIQUE:** DB cũ có partial unique tinh tế (`customers` phone-CaNhan / name-ToChuc; `chart_of_accounts.account_code` unique). Cần đảm bảo bảng mới giữ: `account_code` unique, `orders.code` unique, `finance_transactions.code` unique, `products.sku`/`barcode` unique (DB cũ **không** thấy unique trên sku → nên thêm nếu nghiệp vụ cho phép).
- **FK:** doc DB mới không thể hiện FK nào → **không xác định được từ doc**. DB cũ có FK đầy đủ (chart_of_accounts self-FK, finance_invoice_allocations…). Strangler giữ nguyên FK public.
- **Severity HIGH** (toàn vẹn tiền/đơn/kho).

### 3.6 Denormalize cache — **MED**
- `products.category_name` + `manufacturer_name` **cache song song** với `category_id`/`manufacturer_id` (DB mới có cả hai). DB cũ index theo **tên cache** (`idx_products_category` on `category_name`) → nếu rename category mà không sync ⇒ lệch + index sai.
- **Đề xuất:** giữ cache (đọc nhanh, hợp lý cho catalog read-heavy) **NHƯNG**: (a) index theo `category_id` thay vì name; (b) cơ chế sync rõ ràng (trigger hoặc update qua port `catalog`), không để FE/RPC tự set rời rạc. Memory `project_recent_fixes` cho thấy enrichment đã tập trung — giữ hướng đó.
- **Severity MED.**

### 3.7 Index & pagination cho bảng lớn — **HIGH (verify-gated)**
- DB cũ: `finance_transactions` index tốt (date DESC, fund, partner, ref, status, trgm code/desc). `orders` có index customer/status. `batches` có `(product_id, expiry_date)` cho FEFO. ✅
- **Thiếu/không thấy:** `inventory_transactions` (bảng thẻ kho lớn ở DB mới) — doc không thể hiện index; cần `(warehouse_id, product_id, created_at)` + keyset. `accounting_journals`/journal lines cần `(book_type, entry_date)` + `(account_code, entry_date)`.
- **Pagination:** RPC cũ phần lớn offset-based. Spec §7 yêu cầu **keyset/cursor** cho ledger lớn. → thiết kế cursor `(entry_date, id)` cho journal, `(transaction_date, id)` cho finance.
- **Severity HIGH** cho ledger/inventory_transactions (tăng trưởng vô hạn).

### 3.8 Concurrency / optimistic lock — **HIGH**
- **Cả hai DB đều KHÔNG có `lock_version`/`version`** ở bảng tiền/kho (`fund_accounts.balance`, `product_inventory.stock_quantity`, `inventory_batches.quantity`, `customers.current_debt`). → cập nhật song song dễ lost-update (đã từng có bug lock auth + advisory lock ở `record_manual_payment_received` — memory).
- Spec §7: `SERIALIZABLE` + retry-40001 mặc định, `lock_version` optimistic cho client, `FOR UPDATE` hot-row.
- **Đề xuất:** bảng MỚI tiền/kho thêm `lock_version int NOT NULL DEFAULT 0`; bảng cũ (đọc-as-is) dùng `FOR UPDATE`/advisory lock trong tx Go.
- **Severity HIGH.**

### 3.9 JSONB — chỗ hợp lý vs nên tách bảng — **MED**
- **Hợp lý (giữ):** `usage_instructions` (HDSD bán-cấu trúc), `webhook_logs.payload`, `system_audit_logs.old/new_data`, `chat.llm_meta`, `promotions.rules`, `cash_tally`. Đây là dữ liệu schema-less/append, không query quan hệ.
- **NÊN TÁCH BẢNG:**
  - `customers.b2b_metadata jsonb` (MST, debt_limit, payment_term, công nợ) → **phải là bảng `customers_b2b` chuẩn hoá** (DB cũ đã làm đúng). Query công nợ/hạn mức trên jsonb = chậm + không index + không CHECK. Memory `feedback_portal_rpc_params` + `project_b2b_debt_display_source` đều liên quan B2B → giữ quan hệ chuẩn.
  - `finance_invoices.items_json jsonb` (DB cũ) → cân nhắc tách `finance_invoice_items` nếu cần đối chiếu từng dòng với PO.
- **Severity MED.**

### 3.10 Partition cho ledger/transaction lớn — **MED (future-proof)**
- DB mới **đã partition `system_audit_logs` theo tháng** (`_2026_06`, `_2026_07`) ✅ — chứng tỏ có nhận thức.
- **Đề xuất:** áp range partition theo `entry_date`/`transaction_date` (tháng hoặc năm) cho **journal_entry_lines**, **finance_transactions**, **inventory_transactions** khi volume lớn (nhà thuốc chuỗi → triệu dòng/năm). **Chưa cần ngay** — thiết kế PK/index sao cho **partition được sau mà không breaking** (PK gồm cột phân vùng: `(id, entry_date)`), tránh phải sửa DB lần 2.
- **Severity MED** (future-proof; quyết định "thiết kế-để-partition-được" ngay).

### 3.11 Dual-ledger (invoice_price/actual_cost + sales_invoices) — **đánh giá: SẠCH, giữ**
- `products.invoice_price` (giá HĐ VAT) vs `actual_cost` (giá vốn thực) — tách 2 sổ đúng nguyên tắc memory `project_dual_ledger` (sổ thực tế vs sổ VAT, KHÔNG sync). ✅
- `finance_transactions.book_type ∈ {INTERNAL,TAX,BOTH}` — đúng hướng. Cần **đẩy `book_type` xuống journal** (xem §3.2), không chỉ ở finance.
- `sales_invoices` (cho module vat) **chưa tồn tại** ở cả hai DB → thiết kế mới: header + lines, tách VAT-inclusive→pre-tax theo dòng (làm tròn half-up, line-level), append-only như ledger.
- **Đánh giá:** mô hình dual-ledger **sạch về nguyên tắc**, chỉ thiếu **hiện thực journal 2-sổ** đúng chuẩn.

---

## 4. Khuyến nghị target enterprise (Now / Migrate-later / Read-as-is)

### 4.1 LÀM NGAY — additive, object MỚI ở schema `app` (không breaking)
| # | Hạng mục | Nội dung |
|---|---|---|
| N1 | **Ledger chuẩn** | `app.journal_entries` + `journal_entry_lines` + `account_balances` + `accounting_periods`, có `book_type INTERNAL/TAX`, CHECK debit XOR credit, **Σdebit=Σcredit ép trong DB tx**, append-only (sửa=đảo). Thay `accounting_journals` 1-dòng. |
| N2 | **Money type chốt** | Mọi cột tiền bảng mới: `NUMERIC` scale-0 (VND chốt) / `NUMERIC(20,4)` (đơn giá/VAT/proration). Round-trip test decimal. |
| N3 | **uuid v7** cho mọi PK bảng mới | helper `common/id`; bỏ `gen_random_uuid` v4 cho bảng app mới. |
| N4 | **lock_version** | thêm vào mọi bảng tiền/kho MỚI; SERIALIZABLE+retry ở app. |
| N5 | **idempotency_keys** | đã có (`app.idempotency_keys`) ✅ — dùng cho mọi POST/PATCH tiền + dedupe webhook bank theo mã GD. |
| N6 | **audit mixin** + **CHECK/enum-as-CHECK** cho bảng mới | created/updated/by + status/flow/type CHECK. |
| N7 | **sales_invoices** (module vat) | header+lines, tách VAT theo dòng, append-only. |
| N8 | **Keyset pagination + index** | journal `(book_type,entry_date,id)`, finance `(transaction_date,id)`. |
| N9 | **Thiết kế-để-partition-được** | PK ledger/tx mới gồm cột thời gian để range-partition sau không breaking. |

### 4.2 MIGRATE SAU — đụng bảng/PK hiện hữu, breaking (cần plan + verify prod + downtime)
| # | Hạng mục | Lý do hoãn |
|---|---|---|
| M1 | Bổ sung **CHECK amount>0 / qty>0 / enum-CHECK status** vào bảng `public` cũ | Cần verify prod (`pg_proc`/dump) xem đã có chưa; thêm CHECK trên bảng có data dirty → validate trước; có thể lock-out như bug RLS 2026-06-14. |
| M2 | Tách `customers.b2b_metadata jsonb` → quan hệ `customers_b2b` | Breaking với code đọc jsonb; nhưng DB cũ ERP **đã có `customers_b2b`** → đối chiếu 2 nhánh trước khi quyết hợp nhất. |
| M3 | Thêm `lock_version`/`deleted_at`/`updated_at` vào bảng cũ tiền/kho | ALTER bảng nóng + sửa mọi writer; làm khi context đó cutover sang Go. |
| M4 | Thống nhất kiểu `quantity` integer↔numeric (inventory) | đổi kiểu cột = rewrite + ảnh hưởng RPC; làm khi inventory cutover. |
| M5 | Đổi PK bigint→uuid (nếu muốn nhất quán tuyệt đối) | **KHÔNG khuyến nghị** — phá FK toàn hệ; giữ bigint cũ, chỉ bảng mới uuid v7. |

### 4.3 ĐỌC-AS-IS khi strangler (giữ nguyên, backend chỉ đọc/ghi qua sqlc schema tham chiếu `public`)
- `products`, `product_units`, `categories`, `manufacturers` (catalog read-heavy) — giữ cache cột, đọc as-is.
- `orders`, `order_items` — giữ generated columns (`base_quantity`, `total_line`) ở DB cũ (tốt); Go ghi qua repo.
- `batches`, `inventory_batches`, `product_inventory` — FEFO đọc `(product_id, expiry_date)` index sẵn.
- `finance_transactions`, `fund_accounts` — ghi qua tx Go + `FOR UPDATE` (chưa có lock_version).
- `customers`, `customers_b2b`, `suppliers`, `purchase_orders` — đọc as-is; purchasing là Phase 1.5.
- `system_audit_logs` (đã partition) — append.

---

## 5. Bảng quyết định

| Hạng mục | Quyết định | Lý do |
|---|---|---|
| Mô hình journal | **Làm mới `journal_entries`+`lines` ở `app`**, bỏ `accounting_journals` 1-dòng | 1-dòng không phải double-entry chuẩn, không cân Σ, thiếu 2 sổ → rủi ro #1 |
| 2 sổ INTERNAL/TAX | **`book_type` trên entry, 2 đường posting tách biệt, KHÔNG sync** | TT133 + memory `project_dual_ledger`; test riêng từng đường |
| Tiền | **NUMERIC scale-0 / (20,4), cấm float** | spec §7; round-trip decimal; không phát hiện float ở DB cũ (giữ tốt) |
| PK bảng mới | **uuid v7** (giữ bigint/uuid cũ) | time-ordered index-friendly; không phá FK cũ |
| Enum cũ → text mới | **Tối thiểu CHECK constraint** (không bắt enum PG) | Go enforce + envelope; nhưng DB phải tự bảo vệ miền giá trị |
| `b2b_metadata jsonb` | **Giữ quan hệ `customers_b2b` chuẩn hoá** (migrate-later nếu hợp nhất) | công nợ/hạn mức cần index+CHECK; jsonb chậm |
| Cache category/manufacturer_name | **Giữ + sync qua port + index theo _id** | read-heavy; tránh lệch + index sai |
| lock_version | **Bảng mới: có; bảng cũ: FOR UPDATE trong tx** | chống lost-update tiền/kho |
| Partition | **Thiết kế-để-partition-được ngay, bật khi volume lớn** | future-proof, tránh sửa DB lần 2 |
| Strangler data | **Object mới ở `app` (goose); bảng nghiệp vụ ở `public` đọc-as-is qua sqlc ref schema** | ADR 0001; verify prod trước khi tin doc |
| Nguồn doc DB mới | **Verify với prod trước mọi hành động** | doc thiếu FK/index/constraint → không kết luận "thiếu" từ doc |

---

*File này là review read-only. KHÔNG commit (repo đang có process khác commit). Mọi nhận định "DB mới thiếu X" cần verify prod schema trước khi migrate.*
