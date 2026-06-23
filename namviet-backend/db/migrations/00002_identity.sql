-- +goose Up
-- Module identity: users + RBAC (roles/permissions) + refresh tokens xoay vòng.
-- Schema app đã tạo ở migration 00001. citext cho email không phân biệt hoa
-- thường; pgcrypto cho gen_random_uuid().
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE app.users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         citext UNIQUE NOT NULL,
    password_hash text NOT NULL,
    hash_algo     text NOT NULL CHECK (hash_algo IN ('argon2id', 'bcrypt')),
    user_type     text NOT NULL DEFAULT 'staff',
    is_active     boolean NOT NULL DEFAULT true,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE app.roles (
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text UNIQUE NOT NULL,
    name text NOT NULL
);

CREATE TABLE app.permissions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code        text UNIQUE NOT NULL,
    description text
);

CREATE TABLE app.role_permissions (
    role_id       uuid NOT NULL REFERENCES app.roles (id) ON DELETE CASCADE,
    permission_id uuid NOT NULL REFERENCES app.permissions (id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE app.user_roles (
    user_id uuid NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    role_id uuid NOT NULL REFERENCES app.roles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE app.refresh_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    token_hash text UNIQUE NOT NULL,
    family_id  uuid NOT NULL,
    used       boolean NOT NULL DEFAULT false,
    revoked    boolean NOT NULL DEFAULT false,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- Tra cứu quyền theo user (join user_roles → role_permissions → permissions).
CREATE INDEX idx_user_roles_user ON app.user_roles (user_id);
-- Revoke nguyên family khi phát hiện reuse cần quét theo family_id.
CREATE INDEX idx_refresh_tokens_family ON app.refresh_tokens (family_id);

-- +goose Down
DROP TABLE IF EXISTS app.refresh_tokens;
DROP TABLE IF EXISTS app.user_roles;
DROP TABLE IF EXISTS app.role_permissions;
DROP TABLE IF EXISTS app.permissions;
DROP TABLE IF EXISTS app.roles;
DROP TABLE IF EXISTS app.users;
-- KHÔNG drop schema app (migration 00001 sở hữu) và KHÔNG drop extension
-- citext/pgcrypto (có thể dùng bởi module khác).
