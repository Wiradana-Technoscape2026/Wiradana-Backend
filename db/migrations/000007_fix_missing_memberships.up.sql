-- Repair partial-registration state: member rows exist in multiple cooperatives
-- but user_cooperative_membership was never created for the second+ cooperative
-- because the old code crashed on UNIQUE(email) in app_user before inserting the membership.
--
-- For every (member, app_user) pair where app_user.email = member.nik,
-- ensure a membership row exists. ON CONFLICT DO NOTHING makes this idempotent.
INSERT INTO user_cooperative_membership (id, user_id, cooperative_id, member_id, created_at)
SELECT
    gen_random_uuid(),
    au.id,
    m.cooperative_id,
    m.id,
    NOW()
FROM member m
JOIN app_user au ON au.email = m.nik
ON CONFLICT (user_id, cooperative_id) DO NOTHING;
