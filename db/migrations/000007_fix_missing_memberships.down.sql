-- No destructive rollback: removing memberships that were legitimately missing
-- could lock users out. This migration is safe to leave applied.
SELECT 1;
