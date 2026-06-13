-- Migration: tambah field transparansi pengurus
-- Additive only — tidak menghapus kolom yang ada

-- 1. Tambah full_name ke app_user (nama pengurus/anggota)
ALTER TABLE app_user ADD COLUMN IF NOT EXISTS full_name text;

-- 2. Tambah approved_at ke loan_application (kapan diapprove/direject)
ALTER TABLE loan_application ADD COLUMN IF NOT EXISTS approved_at timestamptz;

-- 3. Tambah recorded_by ke savings_transaction (siapa pengurus yang mencatat)
ALTER TABLE savings_transaction ADD COLUMN IF NOT EXISTS recorded_by uuid REFERENCES app_user(id);

-- Backfill approved_at for existing approved/rejected applications
-- (hanya opsional untuk data lama, set ke created_at sebagai approximation)
UPDATE loan_application
SET approved_at = created_at
WHERE status IN ('approved', 'rejected') AND approved_at IS NULL;
