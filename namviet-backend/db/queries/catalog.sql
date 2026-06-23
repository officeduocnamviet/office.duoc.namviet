-- Catalog (read-mostly, ADR 0001): ĐỌC bảng public.* kế thừa. Mọi query lọc
-- deleted_at IS NULL (soft-delete). Tiền (invoice_price/actual_cost/price*) là
-- numeric → map sang common/money ở repo, KHÔNG dùng float. KHÔNG select cột
-- fts (tsvector) — chỉ dùng nó để filter, không scan ra Go.

-- name: ListProducts :many
-- Keyset pagination theo id ASC: chỉ lấy id > after_id. after_id = 0 cho trang
-- đầu (id bigint luôn > 0). Lấy LIMIT phần tử; app quyết next-cursor. Filter
-- optional category_id + q (ILIKE name/sku). status mặc định lọc ở app qua @status.
SELECT
    id, name, sku, barcode, status,
    category_id, manufacturer_id, category_name, manufacturer_name,
    invoice_price, actual_cost,
    wholesale_unit, retail_unit, conversion_factor,
    product_images, created_at, updated_at
FROM public.products
WHERE deleted_at IS NULL
  AND id > @after_id::bigint
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('category_id')::bigint IS NULL OR category_id = sqlc.narg('category_id')::bigint)
  AND (
        sqlc.narg('q')::text IS NULL
        OR name ILIKE '%' || sqlc.narg('q')::text || '%'
        OR sku ILIKE '%' || sqlc.narg('q')::text || '%'
      )
ORDER BY id ASC
LIMIT @row_limit::int;

-- name: GetProductByID :one
SELECT
    id, name, sku, barcode, status,
    category_id, manufacturer_id, category_name, manufacturer_name,
    invoice_price, actual_cost,
    wholesale_unit, retail_unit, conversion_factor,
    product_images, created_at, updated_at
FROM public.products
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListProductUnits :many
-- Đơn vị tính đa cấp của 1 product (Hộp/Vỉ/Viên). Sắp base trước rồi theo rate.
SELECT
    id, product_id, unit_name, conversion_rate, barcode,
    is_base, is_direct_sale, price_cost, price_sell, unit_type, price
FROM public.product_units
WHERE product_id = $1 AND deleted_at IS NULL
ORDER BY is_base DESC, conversion_rate ASC, id ASC;

-- name: ListCategories :many
SELECT id, name, slug, parent_id, status
FROM public.categories
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY name ASC, id ASC;

-- name: ListManufacturers :many
SELECT id, name, slug, country, logo_url, status
FROM public.manufacturers
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY name ASC, id ASC;
