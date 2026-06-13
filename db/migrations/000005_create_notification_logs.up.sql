CREATE TABLE IF NOT EXISTS notification_log (
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

-- Untuk query riwayat notifikasi per koperasi
CREATE INDEX IF NOT EXISTS idx_notif_coop_created
ON notification_log (cooperative_id, created_at DESC);

-- Untuk query notifikasi per member
CREATE INDEX IF NOT EXISTS idx_notif_member_created
ON notification_log (member_id, created_at DESC);

-- Mencegah notifikasi duplikat untuk event yang sama
CREATE UNIQUE INDEX IF NOT EXISTS idx_notif_dedup
ON notification_log (ref_id, event_type)
WHERE ref_id IS NOT NULL;