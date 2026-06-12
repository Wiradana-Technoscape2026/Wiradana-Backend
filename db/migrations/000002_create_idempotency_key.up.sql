CREATE TABLE idempotency_key (
    key            UUID PRIMARY KEY,
    cooperative_id UUID NOT NULL REFERENCES cooperative(id),
    user_id        UUID NOT NULL REFERENCES app_user(id),
    mutation_type  VARCHAR(100) NOT NULL,
    result_id      UUID,
    status         VARCHAR(20) NOT NULL,
    error_message  TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at   TIMESTAMPTZ
);
CREATE INDEX idx_idempotency_key_coop ON idempotency_key(cooperative_id, created_at);
