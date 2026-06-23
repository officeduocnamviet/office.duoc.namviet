-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE app.idempotency_keys (
    key             text PRIMARY KEY,
    request_hash    text NOT NULL,
    response_status integer,
    response_body   jsonb,
    state           text NOT NULL DEFAULT 'in_progress', -- in_progress | done
    created_at      timestamptz NOT NULL DEFAULT now(),
    completed_at    timestamptz
);

-- +goose Down
DROP TABLE IF EXISTS app.idempotency_keys;
DROP SCHEMA IF EXISTS app;
