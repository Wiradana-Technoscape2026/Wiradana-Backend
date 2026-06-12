BEGIN;

INSERT INTO cooperative (id, name, type, created_at)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'Koperasi Padiwangi', 'simpan_pinjam', NOW()),
  ('22222222-2222-2222-2222-222222222222', 'Koperasi Sawargi', 'simpan_pinjam', NOW());

INSERT INTO member (
  id, cooperative_id, nik, full_name, address, phone_number, birth_date, status, custom_attributes, joined_at
)
VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   '3201010101010001', 'Budi Santosa', 'Jl. Melati No. 10, Bandung', '081234567890', '1990-01-15', 'aktif', '{}'::jsonb, NOW() - INTERVAL '18 months'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111',
   '3201010101010002', 'Siti Aminah', 'Jl. Kenanga No. 21, Bandung', '081234567891', '1992-06-20', 'aktif', '{}'::jsonb, NOW() - INTERVAL '12 months'),
  ('cccccccc-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111',
   '3201010101010003', 'Dedi Rahman', 'Jl. Mawar No. 5, Bandung', '081234567892', '1988-11-02', 'aktif', '{}'::jsonb, NOW() - INTERVAL '8 months'),
  ('dddddddd-dddd-dddd-dddd-dddddddddddd', '22222222-2222-2222-2222-222222222222',
   '3302020202020001', 'Rina Lestari', 'Jl. Anggrek No. 7, Bogor', '081345678901', '1994-03-11', 'aktif', '{}'::jsonb, NOW() - INTERVAL '10 months'),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '22222222-2222-2222-2222-222222222222',
   '3302020202020002', 'Agus Pratama', 'Jl. Teratai No. 14, Bogor', '081345678902', '1991-09-28', 'aktif', '{}'::jsonb, NOW() - INTERVAL '6 months');

INSERT INTO app_user (
  id, cooperative_id, member_id, email, password_hash, role
)
VALUES
  -- Pengurus Padiwangi
  ('99999999-9999-9999-9999-999999999991', '11111111-1111-1111-1111-111111111111',
   NULL, 'pengurus@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'pengurus'),
  -- Anggota Padiwangi (login pakai email)
  ('99999999-9999-9999-9999-999999999992', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'budi@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota'),
  ('99999999-9999-9999-9999-999999999993', '11111111-1111-1111-1111-111111111111',
   'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'siti@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota'),
  -- Pengurus Sawargi
  ('99999999-9999-9999-9999-999999999994', '22222222-2222-2222-2222-222222222222',
   NULL, 'pengurus@sawargi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'pengurus'),
  -- Anggota Sawargi
  ('99999999-9999-9999-9999-999999999995', '22222222-2222-2222-2222-222222222222',
   'dddddddd-dddd-dddd-dddd-dddddddddddd', 'rina@sawargi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota');

INSERT INTO loan_config (
  id, cooperative_id, flat_rate_monthly, max_plafond, penalty_daily
)
VALUES
  ('88888888-8888-8888-8888-888888888881', '11111111-1111-1111-1111-111111111111', 1.50, 20000000, 5000),
  ('88888888-8888-8888-8888-888888888882', '22222222-2222-2222-2222-222222222222', 1.75, 15000000, 5000);

INSERT INTO coop_module (
  id, cooperative_id, module_key, enabled
)
VALUES
  ('77777777-7777-7777-7777-777777777771', '11111111-1111-1111-1111-111111111111', 'simpan_pinjam', TRUE),
  ('77777777-7777-7777-7777-777777777772', '11111111-1111-1111-1111-111111111111', 'inventory', FALSE),
  ('77777777-7777-7777-7777-777777777773', '22222222-2222-2222-2222-222222222222', 'simpan_pinjam', TRUE),
  ('77777777-7777-7777-7777-777777777774', '22222222-2222-2222-2222-222222222222', 'inventory', FALSE);

-- Savings transactions (mix pokok, wajib, sukarela; some anggota with multiple months)
INSERT INTO savings_transaction (
  id, cooperative_id, member_id, savings_type, direction, amount, created_at
)
VALUES
  -- Budi Santoso (pokok + wajib bulanan)
  ('66666666-6666-6666-6666-666666666601', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'pokok', 'setor', 100000, NOW() - INTERVAL '18 months'),
  ('66666666-6666-6666-6666-666666666602', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib', 'setor', 50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666603', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib', 'setor', 50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666604', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib', 'setor', 50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666605', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib', 'setor', 50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666606', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib', 'setor', 50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666607', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib', 'setor', 50000, NOW() - INTERVAL '1 months'),
  -- Siti Aminah (pokok + wajib + sukarela)
  ('66666666-6666-6666-6666-666666666608', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'pokok', 'setor', 100000, NOW() - INTERVAL '12 months'),
  ('66666666-6666-6666-6666-666666666609', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'wajib', 'setor', 50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666610', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'wajib', 'setor', 50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666611', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'wajib', 'setor', 50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666612', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'sukarela', 'setor', 250000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666613', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'sukarela', 'setor', 100000, NOW() - INTERVAL '1 months'),
  -- Dedi Rahman (pokok + wajib)
  ('66666666-6666-6666-6666-666666666614', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'pokok', 'setor', 100000, NOW() - INTERVAL '8 months'),
  ('66666666-6666-6666-6666-666666666615', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'wajib', 'setor', 50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666616', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'sukarela', 'setor', 200000, NOW() - INTERVAL '2 months'),
  -- Rina Lestari (Sawargi)
  ('66666666-6666-6666-6666-666666666617', '22222222-2222-2222-2222-222222222222', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'pokok', 'setor', 150000, NOW() - INTERVAL '10 months'),
  ('66666666-6666-6666-6666-666666666618', '22222222-2222-2222-2222-222222222222', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'wajib', 'setor', 75000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666619', '22222222-2222-2222-2222-222222222222', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'sukarela', 'setor', 500000, NOW() - INTERVAL '2 months'),
  -- Agus Pratama (Sawargi)
  ('66666666-6666-6666-6666-666666666620', '22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'pokok', 'setor', 150000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666621', '22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'wajib', 'setor', 75000, NOW() - INTERVAL '4 months');

-- SHU Period for Padiwangi 2025
INSERT INTO shu_period (
  id, cooperative_id, year, total_shu, pct_jasa_modal, pct_jasa_usaha, status
)
VALUES
  ('shu11111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111',
   2025, 40000000, 20, 25, 'final');

-- SHU Distributions
INSERT INTO shu_distribution (
  id, shu_period_id, member_id, jasa_modal, jasa_usaha, total_shu
)
VALUES
  ('shud1111-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 320000, 166600, 486600),
  ('shud1111-1111-1111-1111-111111111112', 'shu11111-1111-1111-1111-111111111111',
   'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 280000, 133400, 413400),
  ('shud1111-1111-1111-1111-111111111113', 'shu11111-1111-1111-1111-111111111111',
   'cccccccc-cccc-cccc-cccc-cccccccccccc', 200000, 100000, 300000);

-- Loan applications
INSERT INTO loan_application (
  id, cooperative_id, member_id, approved_by, amount, tenor_months, purpose, status, created_at
)
VALUES
  -- Budi - approved -> loan
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '99999999-9999-9999-9999-999999999991',
   5000000, 12, 'Modal tani', 'approved', NOW() - INTERVAL '5 months'),
  -- Siti - rejected
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', '11111111-1111-1111-1111-111111111111',
   'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '99999999-9999-9999-9999-999999999991',
   8000000, 6, 'Usaha warung', 'rejected', NOW() - INTERVAL '3 months'),
  -- Dedi - pending
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac', '11111111-1111-1111-1111-111111111111',
   'cccccccc-cccc-cccc-cccc-cccccccccccc', NULL,
   3000000, 6, 'Renovasi rumah', 'pending', NOW() - INTERVAL '1 months'),
  -- Rina - approved -> loan
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', '22222222-2222-2222-2222-222222222222',
   'dddddddd-dddd-dddd-dddd-dddddddddddd', '99999999-9999-9999-9999-999999999994',
   6000000, 12, 'Pupuk dan bibit', 'approved', NOW() - INTERVAL '4 months');

-- Credit assessments
INSERT INTO credit_assessment (
  id, application_id, score, grade, recommendation, limit_suggested, features, reasons, source
)
VALUES
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 80, 'A', 'approve', 9000000,
   '{"ketepatan_bayar":0.95,"rasio_simpanan_pinjaman":0.6,"lama_keanggotaan_bulan":28,"konsistensi_simpanan":0.9,"rasio_beban_angsuran":0.3}'::jsonb,
   '["riwayat pembayaran sangat baik","simpanan memadai terhadap pinjaman"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', 58, 'D', 'reject', 2000000,
   '{"ketepatan_bayar":0.4,"rasio_simpanan_pinjaman":0.2,"lama_keanggotaan_bulan":12,"konsistensi_simpanan":0.5,"rasio_beban_angsuran":0.9}'::jsonb,
   '["riwayat pembayaran kurang","simpanan tidak memadai"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac', 68, 'C', 'review', 3500000,
   '{"ketepatan_bayar":0.7,"rasio_simpanan_pinjaman":0.5,"lama_keanggotaan_bulan":8,"konsistensi_simpanan":0.7,"rasio_beban_angsuran":0.6}'::jsonb,
   '["simpanan cukup","perlu review manual"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', 75, 'B', 'approve', 7000000,
   '{"ketepatan_bayar":0.85,"rasio_simpanan_pinjaman":0.6,"lama_keanggotaan_bulan":10,"konsistensi_simpanan":0.8,"rasio_beban_angsuran":0.4}'::jsonb,
   '["simpanan cukup","riwayat baik"]'::jsonb, 'MOCK_ADINS_SCORING');

-- Loans (dari approved applications)
INSERT INTO loan (
  id, cooperative_id, application_id, member_id, principal, flat_rate_monthly, tenor_months, status, disbursed_at
)
VALUES
  -- Budi - aktif (5 bulan berjalan, beberapa sudah dibayar)
  ('lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
   5000000, 1.50, 12, 'aktif', (NOW() - INTERVAL '5 months')::date),
  -- Rina - menunggak (ada angsuran terlambat)
  ('lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', '22222222-2222-2222-2222-222222222222',
   'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', 'dddddddd-dddd-dddd-dddd-dddddddddddd',
   6000000, 1.75, 12, 'menunggak', (NOW() - INTERVAL '4 months')::date);

-- Installment schedules (Budi, 12 bulan flat 1.5%)
-- flat_rate=1.5% → interest = round(5jt × 1.5/100) = 75000
-- principal per period = 5jt/12 = 416666, last = 416674 agar total pas
INSERT INTO installment_schedule (
  id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status
)
VALUES
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa01', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, (NOW() - INTERVAL '4 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa02', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2, (NOW() - INTERVAL '3 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa03', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3, (NOW() - INTERVAL '2 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa04', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 4, (NOW() - INTERVAL '1 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa05', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 5, NOW()::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa06', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 6, (NOW() + INTERVAL '1 months')::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa07', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 7, (NOW() + INTERVAL '2 months')::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa08', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 8, (NOW() + INTERVAL '3 months')::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa09', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 9, (NOW() + INTERVAL '4 months')::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa10', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 10, (NOW() + INTERVAL '5 months')::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa11', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 11, (NOW() + INTERVAL '6 months')::date, 416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa12', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 12, (NOW() + INTERVAL '7 months')::date, 416674, 75000, 491674, 'belum_bayar');

-- Installment schedules (Rina, 12 bulan flat 1.75%)
-- interest = round(6jt × 1.75/100) = 105000
-- principal per period = 6jt/12 = 500000
INSERT INTO installment_schedule (
  id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status
)
VALUES
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa13', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 1, (NOW() - INTERVAL '3 months')::date, 500000, 105000, 605000, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa14', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 2, (NOW() - INTERVAL '2 months')::date, 500000, 105000, 605000, 'terlambat'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa15', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 3, (NOW() - INTERVAL '1 months')::date, 500000, 105000, 605000, 'terlambat'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa16', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 4, NOW()::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa17', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 5, (NOW() + INTERVAL '1 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa18', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 6, (NOW() + INTERVAL '2 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa19', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 7, (NOW() + INTERVAL '3 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa20', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 8, (NOW() + INTERVAL '4 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa21', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 9, (NOW() + INTERVAL '5 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa22', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 10, (NOW() + INTERVAL '6 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa23', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 11, (NOW() + INTERVAL '7 months')::date, 500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa24', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 12, (NOW() + INTERVAL '8 months')::date, 500000, 105000, 605000, 'belum_bayar');

-- Payments (Budi paid period 1-4, Rina paid period 1 only)
INSERT INTO payment (
  id, schedule_id, amount, penalty, paid_at
)
VALUES
  -- Budi: 4 payments (period 1-4)
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa01', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa01', 491666, 0, NOW() - INTERVAL '4 months'),
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa02', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa02', 491666, 0, NOW() - INTERVAL '3 months'),
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa03', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa03', 491666, 0, NOW() - INTERVAL '2 months'),
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa04', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa04', 491666, 0, NOW() - INTERVAL '1 months'),
  -- Rina: 1 payment (period 1)
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa05', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa13', 605000, 0, NOW() - INTERVAL '3 months');

COMMIT;
