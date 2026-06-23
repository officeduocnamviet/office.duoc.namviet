-- SCHEMA THAM CHIẾU (strangler-fig, ADR 0001) — KHÔNG phải goose migration.
--
-- File này CHỈ mô tả các bảng public.* ĐANG TỒN TẠI trong Postgres của Supabase
-- để (1) sqlc type-check query catalog ở compile-time và (2) dbtest materialize
-- bảng + seed dữ liệu trong testcontainers (DB test trống). TUYỆT ĐỐI KHÔNG chạy
-- qua goose lên prod — backend chỉ ĐỌC các bảng kế thừa này, không sở hữu chúng.
--
-- Cột lấy nguyên văn từ office.duoc.namviet/database_schema.md (doc). Đây là
-- nguồn doc, PHẢI verify lại với prod thật (pg_dump/REST) khi có creds — prod
-- hay lệch migration. id = bigint (KHÔNG uuid). Tiền = numeric (→ common/money).

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.categories (
    id         bigint PRIMARY KEY,
    name       text NOT NULL,
    slug       text NOT NULL,
    parent_id  bigint,
    status     text NOT NULL DEFAULT 'active',
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS public.manufacturers (
    id         bigint PRIMARY KEY,
    name       text NOT NULL,
    slug       text NOT NULL,
    country    text,
    logo_url   text,
    status     text DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS public.products (
    id                     bigint PRIMARY KEY,
    name                   text NOT NULL,
    sku                    text,
    barcode                text,
    description            text,
    active_ingredient      text,
    image_url              text,
    status                 text NOT NULL DEFAULT 'active',
    fts                    tsvector,
    category_id            bigint,
    manufacturer_id        bigint,
    category_name          text,
    manufacturer_name      text,
    distributor_id         bigint,
    invoice_price          numeric DEFAULT 0,
    actual_cost            numeric NOT NULL DEFAULT 0,
    wholesale_unit         text DEFAULT 'Hộp',
    retail_unit            text DEFAULT 'Vỉ',
    conversion_factor      integer DEFAULT 1,
    wholesale_margin_value numeric DEFAULT 0,
    wholesale_margin_type  text DEFAULT '%',
    retail_margin_value    numeric DEFAULT 0,
    retail_margin_type     text DEFAULT '%',
    items_per_carton       integer DEFAULT 1,
    carton_weight          numeric DEFAULT 0,
    carton_dimensions      text,
    purchasing_policy      text DEFAULT 'ALLOW_LOOSE',
    registration_number    text,
    packing_spec           text,
    stock_management_type  text DEFAULT 'lot_date',
    wholesale_margin_rate  numeric DEFAULT 0,
    retail_margin_rate     numeric DEFAULT 0,
    usage_instructions     jsonb DEFAULT '{}'::jsonb,
    stock_status           text DEFAULT 'in_stock',
    product_images         text[] DEFAULT '{}'::text[],
    updated_by             uuid,
    created_at             timestamptz DEFAULT now(),
    updated_at             timestamptz DEFAULT now(),
    deleted_at             timestamptz
);

CREATE TABLE IF NOT EXISTS public.product_units (
    id              bigint PRIMARY KEY,
    product_id      bigint,
    unit_name       text NOT NULL,
    conversion_rate integer DEFAULT 1,
    barcode         text,
    is_base         boolean DEFAULT false,
    is_direct_sale  boolean DEFAULT true,
    price_cost      numeric DEFAULT 0,
    price_sell      numeric DEFAULT 0,
    unit_type       text DEFAULT 'retail',
    price           numeric DEFAULT 0,
    created_at      timestamptz DEFAULT now(),
    updated_at      timestamptz DEFAULT now(),
    deleted_at      timestamptz
);
