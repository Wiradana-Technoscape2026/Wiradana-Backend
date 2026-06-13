-- Rollback: hapus field transparansi
ALTER TABLE savings_transaction DROP COLUMN IF EXISTS recorded_by;
ALTER TABLE loan_application DROP COLUMN IF EXISTS approved_at;
ALTER TABLE app_user DROP COLUMN IF EXISTS full_name;
