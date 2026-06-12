CREATE TABLE notification_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cooperative_id  UUID NOT NULL REFERENCES cooperative(id),
    member_id       UUID NOT NULL REFERENCES member(id),
    phone_number    TEXT NOT NULL,
    channel         TEXT NOT NULL DEFAULT 'whatsapp',
    event_type      TEXT NOT NULL,
    ref_id          UUID,
    message         TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'sent',
    source          TEXT NOT NULL DEFAULT 'MOCK_WA',
    error_msg       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notif_coop_created ON notification_log (cooperative_id, created_at DESC);

CREATE UNIQUE INDEX idx_notif_dedup ON notification_log (ref_id, event_type, (created_at::date))
    WHERE ref_id IS NOT NULL;
