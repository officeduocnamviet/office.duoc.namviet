-- name: GetIdempotencyKey :one
SELECT key, request_hash, response_status, response_body, state
FROM app.idempotency_keys
WHERE key = $1;

-- name: InsertIdempotencyKey :exec
-- Begin: chèn key trạng thái in_progress; nếu key đã tồn tại thì bỏ qua
-- (ON CONFLICT DO NOTHING) — giữ đúng behavior chống double-POST.
INSERT INTO app.idempotency_keys (key, request_hash, state)
VALUES ($1, $2, 'in_progress')
ON CONFLICT (key) DO NOTHING;

-- name: CompleteIdempotencyKey :exec
UPDATE app.idempotency_keys
SET response_status = $2, response_body = $3, state = 'done', completed_at = now()
WHERE key = $1;
