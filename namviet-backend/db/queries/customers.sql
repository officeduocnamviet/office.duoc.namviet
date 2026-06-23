-- Customers (read-mostly, ADR 0001): ĐỌC bảng public.* kế thừa. Mọi query lọc
-- deleted_at IS NULL (soft-delete). Tiền (current_debt) là numeric → map sang
-- common/money ở repo, KHÔNG dùng float.
--
-- Công nợ LIVE = tổng final_amount các đơn CHƯA tất toán (payment_status <>
-- 'paid'), tính từ public.orders, loại đơn soft-deleted. Đây là nguồn ưu tiên so
-- với cột tĩnh customers.current_debt (stale). CAVEAT: không trừ được phần đã trả
-- của đơn 'partial' vì orders KHÔNG có paid_amount (ghost-debt) — chỉ trình bày.
-- B2B (tax_code/debt_limit/payment_term) nằm trong b2b_metadata jsonb → trả raw
-- jsonb, repo parse (KHÔNG suy diễn ở SQL).

-- name: ListCustomers :many
-- Keyset pagination theo id ASC: chỉ lấy id > after_id (after_id = 0 cho trang
-- đầu; id bigint luôn > 0). Lọc optional customer_type ('B2B'/'B2C'), status, và
-- q (ILIKE name/phone/customer_code + b2b_metadata->>tax_code).
SELECT
    c.id,
    c.customer_code,
    c.name,
    c.customer_type,
    c.phone,
    c.email,
    c.address,
    c.status,
    c.b2b_metadata,
    c.current_debt,
    c.created_at,
    c.updated_at,
    COALESCE((
        SELECT SUM(o.final_amount)
        FROM public.orders o
        WHERE o.customer_id = c.id
          AND o.deleted_at IS NULL
          AND o.payment_status <> 'paid'
    ), 0)::numeric AS live_debt
FROM public.customers c
WHERE c.deleted_at IS NULL
  AND c.id > @after_id::bigint
  AND (sqlc.narg('customer_type')::text IS NULL OR c.customer_type = sqlc.narg('customer_type')::text)
  AND (sqlc.narg('status')::text IS NULL OR c.status = sqlc.narg('status')::text)
  AND (
        sqlc.narg('q')::text IS NULL
        OR c.name ILIKE '%' || sqlc.narg('q')::text || '%'
        OR c.phone ILIKE '%' || sqlc.narg('q')::text || '%'
        OR c.customer_code ILIKE '%' || sqlc.narg('q')::text || '%'
        OR (c.b2b_metadata ->> 'tax_code') ILIKE '%' || sqlc.narg('q')::text || '%'
      )
ORDER BY c.id ASC
LIMIT @row_limit::int;

-- name: GetCustomerByID :one
SELECT
    c.id,
    c.customer_code,
    c.name,
    c.customer_type,
    c.phone,
    c.email,
    c.address,
    c.status,
    c.b2b_metadata,
    c.current_debt,
    c.created_at,
    c.updated_at,
    COALESCE((
        SELECT SUM(o.final_amount)
        FROM public.orders o
        WHERE o.customer_id = c.id
          AND o.deleted_at IS NULL
          AND o.payment_status <> 'paid'
    ), 0)::numeric AS live_debt
FROM public.customers c
WHERE c.id = $1 AND c.deleted_at IS NULL;
