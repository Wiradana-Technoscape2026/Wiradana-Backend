CREATE TABLE IF NOT EXISTS user_cooperative_membership (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID        NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
    cooperative_id UUID        NOT NULL REFERENCES cooperative(id) ON DELETE CASCADE,
    member_id      UUID        REFERENCES member(id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, cooperative_id)
);

-- Populate from existing app_user rows
INSERT INTO user_cooperative_membership (user_id, cooperative_id, member_id)
SELECT id, cooperative_id, member_id
FROM app_user
ON CONFLICT (user_id, cooperative_id) DO NOTHING;
