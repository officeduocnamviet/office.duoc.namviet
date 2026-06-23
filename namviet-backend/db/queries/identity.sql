-- name: GetUserByEmail :one
SELECT id, email, password_hash, hash_algo, user_type, is_active, created_at, updated_at
FROM app.users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password_hash, hash_algo, user_type, is_active, created_at, updated_at
FROM app.users
WHERE id = $1;

-- name: UpdateUserPasswordHash :exec
UPDATE app.users
SET password_hash = $2, hash_algo = $3, updated_at = now()
WHERE id = $1;

-- name: InsertRefreshToken :exec
INSERT INTO app.refresh_tokens (id, user_id, token_hash, family_id, expires_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, family_id, used, revoked, expires_at, created_at
FROM app.refresh_tokens
WHERE token_hash = $1;

-- name: MarkRefreshTokenUsed :execrows
-- Claim NGUYÊN TỬ: chỉ đánh dấu used khi token còn used=false. Trả số dòng
-- ảnh hưởng: 1 = claim thành công; 0 = đã used (mất race / reuse) → caller xử lý
-- như reuse-detection. Điều kiện used=false là then chốt chống rotation race.
UPDATE app.refresh_tokens
SET used = true
WHERE id = $1 AND used = false;

-- name: RevokeRefreshTokenFamily :exec
UPDATE app.refresh_tokens
SET revoked = true
WHERE family_id = $1;

-- name: PermissionCodesForUser :many
SELECT DISTINCT p.code
FROM app.user_roles ur
JOIN app.role_permissions rp ON rp.role_id = ur.role_id
JOIN app.permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
ORDER BY p.code;

-- Các query phục vụ SEED dữ liệu test/integration (tạo user + gán quyền).

-- name: InsertUser :one
INSERT INTO app.users (email, password_hash, hash_algo, user_type, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: InsertRole :one
INSERT INTO app.roles (code, name) VALUES ($1, $2) RETURNING id;

-- name: InsertPermission :one
INSERT INTO app.permissions (code, description) VALUES ($1, $2) RETURNING id;

-- name: AssignPermissionToRole :exec
INSERT INTO app.role_permissions (role_id, permission_id) VALUES ($1, $2);

-- name: AssignRoleToUser :exec
INSERT INTO app.user_roles (user_id, role_id) VALUES ($1, $2);
