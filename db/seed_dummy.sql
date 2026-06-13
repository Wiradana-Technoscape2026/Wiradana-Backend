BEGIN;

INSERT INTO cooperative (id, name, type, created_at)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'Koperasi Padiwangi', 'simpan_pinjam', NOW()),
  ('22222222-2222-2222-2222-222222222222', 'Koperasi Sawargi', 'simpan_pinjam', NOW());

INSERT INTO member (
  id, cooperative_id, nik, full_name, address, phone_number, birth_date, status, custom_attributes, joined_at
)
VALUES
  -- Padiwangi (18 anggota)
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   '3201010101010001', 'Budi Santosa', 'Jl. Melati No. 10, Bandung', '081234567890', '1990-01-15', 'aktif', '{}'::jsonb, NOW() - INTERVAL '18 months'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111',
   '3201010101010002', 'Siti Aminah', 'Jl. Kenanga No. 21, Bandung', '081234567891', '1992-06-20', 'aktif', '{}'::jsonb, NOW() - INTERVAL '12 months'),
  ('cccccccc-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111',
   '3201010101010003', 'Dedi Rahman', 'Jl. Mawar No. 5, Bandung', '081234567892', '1988-11-02', 'aktif', '{}'::jsonb, NOW() - INTERVAL '8 months'),
  ('a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', '11111111-1111-1111-1111-111111111111',
   '3201010101010004', 'Asep Kurniawan', 'Jl. Cempaka No. 3, Bandung', '081234567893', '1987-03-22', 'aktif', '{}'::jsonb, NOW() - INTERVAL '24 months'),
  ('a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', '11111111-1111-1111-1111-111111111111',
   '3201010101010005', 'Sri Wahyuni', 'Jl. Dahlia No. 8, Bandung', '081234567894', '1993-08-10', 'aktif', '{}'::jsonb, NOW() - INTERVAL '20 months'),
  ('a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', '11111111-1111-1111-1111-111111111111',
   '3201010101010006', 'Hendra Gunawan', 'Jl. Flamboyan No. 12, Bandung', '081234567895', '1985-12-05', 'aktif', '{}'::jsonb, NOW() - INTERVAL '16 months'),
  ('a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', '11111111-1111-1111-1111-111111111111',
   '3201010101010007', 'Eni Rahayu', 'Jl. Bougenville No. 4, Bandung', '081234567896', '1995-05-17', 'aktif', '{}'::jsonb, NOW() - INTERVAL '14 months'),
  ('a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', '11111111-1111-1111-1111-111111111111',
   '3201010101010008', 'Yusuf Hidayat', 'Jl. Kamboja No. 9, Bandung', '081234567897', '1991-02-28', 'aktif', '{}'::jsonb, NOW() - INTERVAL '12 months'),
  ('a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', '11111111-1111-1111-1111-111111111111',
   '3201010101010009', 'Dewi Kartika', 'Jl. Anyelir No. 6, Bandung', '081234567898', '1996-09-14', 'aktif', '{}'::jsonb, NOW() - INTERVAL '10 months'),
  ('a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', '11111111-1111-1111-1111-111111111111',
   '3201010101010010', 'Ahmad Fauzi', 'Jl. Tulip No. 2, Bandung', '081234567899', '1989-07-30', 'aktif', '{}'::jsonb, NOW() - INTERVAL '9 months'),
  ('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', '11111111-1111-1111-1111-111111111111',
   '3201010101010011', 'Nining Suryani', 'Jl. Seroja No. 15, Bandung', '081234567800', '1994-04-06', 'aktif', '{}'::jsonb, NOW() - INTERVAL '8 months'),
  ('b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', '11111111-1111-1111-1111-111111111111',
   '3201010101010012', 'Bambang Wahyudi', 'Jl. Teratai No. 11, Bandung', '081234567801', '1986-10-19', 'aktif', '{}'::jsonb, NOW() - INTERVAL '7 months'),
  ('b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', '11111111-1111-1111-1111-111111111111',
   '3201010101010013', 'Tuti Indrawati', 'Jl. Lavender No. 7, Bandung', '081234567802', '1997-01-23', 'aktif', '{}'::jsonb, NOW() - INTERVAL '6 months'),
  ('b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', '11111111-1111-1111-1111-111111111111',
   '3201010101010014', 'Rudi Hermawan', 'Jl. Aster No. 18, Bandung', '081234567803', '1990-06-11', 'aktif', '{}'::jsonb, NOW() - INTERVAL '6 months'),
  ('b5b5b5b5-b5b5-b5b5-b5b5-b5b5b5b5b5b5', '11111111-1111-1111-1111-111111111111',
   '3201010101010015', 'Wati Setiawan', 'Jl. Lily No. 20, Bandung', '081234567804', '1992-11-08', 'aktif', '{}'::jsonb, NOW() - INTERVAL '5 months'),
  ('b6b6b6b6-b6b6-b6b6-b6b6-b6b6b6b6b6b6', '11111111-1111-1111-1111-111111111111',
   '3201010101010016', 'Joko Susilo', 'Jl. Sakura No. 1, Bandung', '081234567805', '1988-03-15', 'aktif', '{}'::jsonb, NOW() - INTERVAL '4 months'),
  ('b7b7b7b7-b7b7-b7b7-b7b7-b7b7b7b7b7b7', '11111111-1111-1111-1111-111111111111',
   '3201010101010017', 'Endang Purwati', 'Jl. Melur No. 13, Bandung', '081234567806', '1995-08-27', 'aktif', '{}'::jsonb, NOW() - INTERVAL '3 months'),
  ('b8b8b8b8-b8b8-b8b8-b8b8-b8b8b8b8b8b8', '11111111-1111-1111-1111-111111111111',
   '3201010101010018', 'Heri Pranoto', 'Jl. Padi No. 5, Bandung', '081234567807', '1993-12-02', 'aktif', '{}'::jsonb, NOW() - INTERVAL '2 months'),
  -- Sawargi (2 anggota)
  ('dddddddd-dddd-dddd-dddd-dddddddddddd', '22222222-2222-2222-2222-222222222222',
   '3302020202020001', 'Rina Lestari', 'Jl. Anggrek No. 7, Bogor', '081345678901', '1994-03-11', 'aktif', '{}'::jsonb, NOW() - INTERVAL '10 months'),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '22222222-2222-2222-2222-222222222222',
   '3302020202020002', 'Agus Pratama', 'Jl. Teratai No. 14, Bogor', '081345678902', '1991-09-28', 'aktif', '{}'::jsonb, NOW() - INTERVAL '6 months');

INSERT INTO app_user (
  id, cooperative_id, member_id, email, password_hash, role, full_name
)
VALUES
  ('99999999-9999-9999-9999-999999999991', '11111111-1111-1111-1111-111111111111',
   NULL, 'pengurus@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'pengurus', 'Ahmad Fauzi'),
  ('99999999-9999-9999-9999-999999999992', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'budi@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota', NULL),
  ('99999999-9999-9999-9999-999999999993', '11111111-1111-1111-1111-111111111111',
   'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'siti@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota', NULL),
  ('99999999-9999-9999-9999-999999999994', '22222222-2222-2222-2222-222222222222',
   NULL, 'pengurus@sawargi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'pengurus', 'Bambang Wijaya'),
  ('99999999-9999-9999-9999-999999999995', '22222222-2222-2222-2222-222222222222',
   'dddddddd-dddd-dddd-dddd-dddddddddddd', 'rina@sawargi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota', NULL),
  ('99999999-9999-9999-9999-999999999996', '11111111-1111-1111-1111-111111111111',
   'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'asep@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota', NULL),
  ('99999999-9999-9999-9999-999999999997', '11111111-1111-1111-1111-111111111111',
   'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'sri@padiwangi.id', '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6', 'anggota', NULL);

INSERT INTO loan_config (id, cooperative_id, flat_rate_monthly, max_plafond, penalty_daily)
VALUES
  ('88888888-8888-8888-8888-888888888881', '11111111-1111-1111-1111-111111111111', 1.50, 20000000, 5000),
  ('88888888-8888-8888-8888-888888888882', '22222222-2222-2222-2222-222222222222', 1.75, 15000000, 5000);

INSERT INTO coop_module (id, cooperative_id, module_key, enabled)
VALUES
  ('77777777-7777-7777-7777-777777777771', '11111111-1111-1111-1111-111111111111', 'simpan_pinjam', TRUE),
  ('77777777-7777-7777-7777-777777777772', '11111111-1111-1111-1111-111111111111', 'inventory', FALSE),
  ('77777777-7777-7777-7777-777777777773', '22222222-2222-2222-2222-222222222222', 'simpan_pinjam', TRUE),
  ('77777777-7777-7777-7777-777777777774', '22222222-2222-2222-2222-222222222222', 'inventory', FALSE);

-- ============================================================
-- SAVINGS TRANSACTIONS (INSERT-ONLY)
-- ============================================================
INSERT INTO savings_transaction (id, cooperative_id, member_id, savings_type, direction, amount, created_at)
VALUES
  -- Budi Santosa: pokok + 6 bln wajib + sukarela
  ('66666666-6666-6666-6666-666666666601', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'pokok',    'setor', 100000, NOW() - INTERVAL '18 months'),
  ('66666666-6666-6666-6666-666666666602', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666603', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666604', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666605', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666606', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666607', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Siti Aminah: pokok + 3 bln wajib + sukarela
  ('66666666-6666-6666-6666-666666666608', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'pokok',    'setor', 100000, NOW() - INTERVAL '12 months'),
  ('66666666-6666-6666-6666-666666666609', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666610', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666611', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666612', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'sukarela', 'setor', 250000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666613', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'sukarela', 'setor', 100000, NOW() - INTERVAL '1 months'),
  -- Dedi Rahman: pokok + 1 bln wajib + sukarela
  ('66666666-6666-6666-6666-666666666614', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'pokok',    'setor', 100000, NOW() - INTERVAL '8 months'),
  ('66666666-6666-6666-6666-666666666615', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666616', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'sukarela', 'setor', 200000, NOW() - INTERVAL '2 months'),
  -- Asep Kurniawan: pokok + 12 bln wajib + sukarela
  ('66666666-6666-6666-6666-666666666622', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'pokok',    'setor', 100000, NOW() - INTERVAL '24 months'),
  ('66666666-6666-6666-6666-666666666623', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '12 months'),
  ('66666666-6666-6666-6666-666666666624', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '11 months'),
  ('66666666-6666-6666-6666-666666666625', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '10 months'),
  ('66666666-6666-6666-6666-666666666626', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '9 months'),
  ('66666666-6666-6666-6666-666666666627', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '8 months'),
  ('66666666-6666-6666-6666-666666666628', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '7 months'),
  ('66666666-6666-6666-6666-666666666629', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666630', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666631', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666632', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666633', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666634', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  ('66666666-6666-6666-6666-666666666635', '11111111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 'sukarela', 'setor', 500000, NOW() - INTERVAL '6 months'),
  -- Sri Wahyuni: pokok + 10 bln wajib
  ('66666666-6666-6666-6666-666666666636', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'pokok',    'setor', 100000, NOW() - INTERVAL '20 months'),
  ('66666666-6666-6666-6666-666666666637', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '10 months'),
  ('66666666-6666-6666-6666-666666666638', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '9 months'),
  ('66666666-6666-6666-6666-666666666639', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '8 months'),
  ('66666666-6666-6666-6666-666666666640', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '7 months'),
  ('66666666-6666-6666-6666-666666666641', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666642', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666643', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666644', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666645', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666646', '11111111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Hendra Gunawan: pokok + 8 bln wajib
  ('66666666-6666-6666-6666-666666666647', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'pokok',    'setor', 100000, NOW() - INTERVAL '16 months'),
  ('66666666-6666-6666-6666-666666666648', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '8 months'),
  ('66666666-6666-6666-6666-666666666649', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '7 months'),
  ('66666666-6666-6666-6666-666666666650', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666651', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666652', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666653', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666654', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666655', '11111111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Eni Rahayu: pokok + 6 bln wajib
  ('66666666-6666-6666-6666-666666666656', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'pokok',    'setor', 100000, NOW() - INTERVAL '14 months'),
  ('66666666-6666-6666-6666-666666666657', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666658', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666659', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666660', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666661', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666662', '11111111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Yusuf Hidayat: pokok + 6 bln wajib
  ('66666666-6666-6666-6666-666666666663', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'pokok',    'setor', 100000, NOW() - INTERVAL '12 months'),
  ('66666666-6666-6666-6666-666666666664', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'wajib',    'setor',  50000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666665', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'wajib',    'setor',  50000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666666', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666667', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666668', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666669', '11111111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Dewi Kartika: pokok + 4 bln wajib
  ('66666666-6666-6666-6666-666666666670', '11111111-1111-1111-1111-111111111111', 'a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', 'pokok',    'setor', 100000, NOW() - INTERVAL '10 months'),
  ('66666666-6666-6666-6666-666666666671', '11111111-1111-1111-1111-111111111111', 'a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666672', '11111111-1111-1111-1111-111111111111', 'a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666673', '11111111-1111-1111-1111-111111111111', 'a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666674', '11111111-1111-1111-1111-111111111111', 'a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Ahmad Fauzi: pokok + 4 bln wajib
  ('66666666-6666-6666-6666-666666666675', '11111111-1111-1111-1111-111111111111', 'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'pokok',    'setor', 100000, NOW() - INTERVAL '9 months'),
  ('66666666-6666-6666-6666-666666666676', '11111111-1111-1111-1111-111111111111', 'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'wajib',    'setor',  50000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666677', '11111111-1111-1111-1111-111111111111', 'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666678', '11111111-1111-1111-1111-111111111111', 'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666679', '11111111-1111-1111-1111-111111111111', 'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Nining Suryani: pokok + 3 bln wajib
  ('66666666-6666-6666-6666-666666666680', '11111111-1111-1111-1111-111111111111', 'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'pokok',    'setor', 100000, NOW() - INTERVAL '8 months'),
  ('66666666-6666-6666-6666-666666666681', '11111111-1111-1111-1111-111111111111', 'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666682', '11111111-1111-1111-1111-111111111111', 'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666683', '11111111-1111-1111-1111-111111111111', 'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Bambang Wahyudi: pokok + 3 bln wajib
  ('66666666-6666-6666-6666-666666666684', '11111111-1111-1111-1111-111111111111', 'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 'pokok',    'setor', 100000, NOW() - INTERVAL '7 months'),
  ('66666666-6666-6666-6666-666666666685', '11111111-1111-1111-1111-111111111111', 'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 'wajib',    'setor',  50000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666686', '11111111-1111-1111-1111-111111111111', 'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666687', '11111111-1111-1111-1111-111111111111', 'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Tuti, Rudi: pokok + 2 bln wajib each
  ('66666666-6666-6666-6666-666666666688', '11111111-1111-1111-1111-111111111111', 'b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 'pokok',    'setor', 100000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666689', '11111111-1111-1111-1111-111111111111', 'b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666690', '11111111-1111-1111-1111-111111111111', 'b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  ('66666666-6666-6666-6666-666666666691', '11111111-1111-1111-1111-111111111111', 'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 'pokok',    'setor', 100000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666692', '11111111-1111-1111-1111-111111111111', 'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 'wajib',    'setor',  50000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666693', '11111111-1111-1111-1111-111111111111', 'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Wati, Joko: pokok + 1 bln wajib each
  ('66666666-6666-6666-6666-666666666694', '11111111-1111-1111-1111-111111111111', 'b5b5b5b5-b5b5-b5b5-b5b5-b5b5b5b5b5b5', 'pokok',    'setor', 100000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666695', '11111111-1111-1111-1111-111111111111', 'b5b5b5b5-b5b5-b5b5-b5b5-b5b5b5b5b5b5', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  ('66666666-6666-6666-6666-666666666696', '11111111-1111-1111-1111-111111111111', 'b6b6b6b6-b6b6-b6b6-b6b6-b6b6b6b6b6b6', 'pokok',    'setor', 100000, NOW() - INTERVAL '4 months'),
  ('66666666-6666-6666-6666-666666666697', '11111111-1111-1111-1111-111111111111', 'b6b6b6b6-b6b6-b6b6-b6b6-b6b6b6b6b6b6', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Endang, Heri: pokok only (baru bergabung)
  ('66666666-6666-6666-6666-666666666698', '11111111-1111-1111-1111-111111111111', 'b7b7b7b7-b7b7-b7b7-b7b7-b7b7b7b7b7b7', 'pokok',    'setor', 100000, NOW() - INTERVAL '3 months'),
  ('66666666-6666-6666-6666-666666666699', '11111111-1111-1111-1111-111111111111', 'b7b7b7b7-b7b7-b7b7-b7b7-b7b7b7b7b7b7', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  ('66666666-6666-6666-6666-6666666666a0', '11111111-1111-1111-1111-111111111111', 'b8b8b8b8-b8b8-b8b8-b8b8-b8b8b8b8b8b8', 'pokok',    'setor', 100000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-6666666666a1', '11111111-1111-1111-1111-111111111111', 'b8b8b8b8-b8b8-b8b8-b8b8-b8b8b8b8b8b8', 'wajib',    'setor',  50000, NOW() - INTERVAL '1 months'),
  -- Rina Lestari & Agus Pratama (Sawargi)
  ('66666666-6666-6666-6666-666666666617', '22222222-2222-2222-2222-222222222222', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'pokok',    'setor', 150000, NOW() - INTERVAL '10 months'),
  ('66666666-6666-6666-6666-666666666618', '22222222-2222-2222-2222-222222222222', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'wajib',    'setor',  75000, NOW() - INTERVAL '5 months'),
  ('66666666-6666-6666-6666-666666666619', '22222222-2222-2222-2222-222222222222', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'sukarela', 'setor', 500000, NOW() - INTERVAL '2 months'),
  ('66666666-6666-6666-6666-666666666620', '22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'pokok',    'setor', 150000, NOW() - INTERVAL '6 months'),
  ('66666666-6666-6666-6666-666666666621', '22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'wajib',    'setor',  75000, NOW() - INTERVAL '4 months');

-- ============================================================
-- SHU PERIOD (Padiwangi 2025, final)
-- ============================================================
INSERT INTO shu_period (id, cooperative_id, year, total_shu, pct_jasa_modal, pct_jasa_usaha, status)
VALUES
  ('shu11111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111',
   2025, 40000000, 20, 25, 'final');

-- SHU Distributions (18 anggota Padiwangi)
-- shu_modal_pool = 8000000, total_simpanan ≈ 5850000
-- shu_usaha_pool = 10000000, paid by members with loans
INSERT INTO shu_distribution (id, shu_period_id, member_id, jasa_modal, jasa_usaha, total_shu)
VALUES
  ('shud0001-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 546000,  1989000, 2535000),
  ('shud0002-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 341000,       0,  341000),
  ('shud0003-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 205000,       0,  205000),
  ('shud0004-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 957000,  1790000, 2747000),
  ('shud0005-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 820000,  3183000, 4003000),
  ('shud0006-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 683000,  2785000, 3468000),
  ('shud0007-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', 615000,   248000,  863000),
  ('shud0008-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a8a8a8a8-a8a8-a8a8-a8a8-a8a8a8a8a8a8', 546000,       0,  546000),
  ('shud0009-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a9a9a9a9-a9a9-a9a9-a9a9-a9a9a9a9a9a9', 410000,       0,  410000),
  ('shud0010-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 342000,       0,  342000),
  ('shud0011-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 274000,       0,  274000),
  ('shud0012-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 274000,       0,  274000),
  ('shud0013-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 274000,       0,  274000),
  ('shud0014-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 274000,       0,  274000),
  ('shud0015-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b5b5b5b5-b5b5-b5b5-b5b5-b5b5b5b5b5b5', 205000,       0,  205000),
  ('shud0016-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b6b6b6b6-b6b6-b6b6-b6b6-b6b6b6b6b6b6', 205000,       0,  205000),
  ('shud0017-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b7b7b7b7-b7b7-b7b7-b7b7-b7b7b7b7b7b7', 137000,       0,  137000),
  ('shud0018-1111-1111-1111-111111111111', 'shu11111-1111-1111-1111-111111111111', 'b8b8b8b8-b8b8-b8b8-b8b8-b8b8b8b8b8b8', 137000,       0,  137000);

-- ============================================================
-- LOAN APPLICATIONS (9 total: 4 existing + 5 new)
-- ============================================================
INSERT INTO loan_application (
  id, cooperative_id, member_id, approved_by, amount, tenor_months, purpose, status, created_at
)
VALUES
  -- Existing 3 Padiwangi
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '99999999-9999-9999-9999-999999999991',
   5000000, 12, 'Modal tani', 'approved', NOW() - INTERVAL '5 months'),
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', '11111111-1111-1111-1111-111111111111',
   'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '99999999-9999-9999-9999-999999999991',
   8000000, 6, 'Usaha warung', 'rejected', NOW() - INTERVAL '3 months'),
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac', '11111111-1111-1111-1111-111111111111',
   'cccccccc-cccc-cccc-cccc-cccccccccccc', NULL,
   3000000, 6, 'Renovasi rumah', 'pending', NOW() - INTERVAL '1 months'),
  -- Asep - approved -> loan lunas
  ('lapp0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', '99999999-9999-9999-9999-999999999991',
   3000000, 6, 'Pengembangan usaha tani', 'approved', NOW() - INTERVAL '8 months'),
  -- Sri - approved -> loan aktif
  ('lapp0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', '99999999-9999-9999-9999-999999999991',
   4000000, 12, 'Modal dagang', 'approved', NOW() - INTERVAL '9 months'),
  -- Hendra - approved -> loan aktif
  ('lapp0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', '99999999-9999-9999-9999-999999999991',
   7000000, 12, 'Pembelian mesin', 'approved', NOW() - INTERVAL '5 months'),
  -- Eni - approved -> loan menunggak
  ('lapp0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7', '99999999-9999-9999-9999-999999999991',
   2500000, 6, 'Keperluan rumah tangga', 'approved', NOW() - INTERVAL '5 months'),
  -- Ahmad - pending
  ('lapp0010-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', NULL,
   4000000, 12, 'Modal usaha beras', 'pending', NOW() - INTERVAL '2 weeks'),
  -- Sawargi: Rina - approved -> loan menunggak
  ('laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', '22222222-2222-2222-2222-222222222222',
   'dddddddd-dddd-dddd-dddd-dddddddddddd', '99999999-9999-9999-9999-999999999994',
   6000000, 12, 'Pupuk dan bibit', 'approved', NOW() - INTERVAL '4 months');

-- ============================================================
-- CREDIT ASSESSMENTS
-- ============================================================
INSERT INTO credit_assessment (
  id, application_id, score, grade, recommendation, limit_suggested, features, reasons, source
)
VALUES
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 80, 'A', 'approve', 9000000,
   '{"ketepatan_bayar":0.95,"rasio_simpanan_pinjaman":0.6,"lama_keanggotaan_hari":840,"konsistensi_simpanan":0.9,"rasio_beban_angsuran":0.3}'::jsonb,
   '["riwayat pembayaran sangat baik","simpanan memadai terhadap pinjaman"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab', 58, 'D', 'reject', 2000000,
   '{"ketepatan_bayar":0.4,"rasio_simpanan_pinjaman":0.2,"lama_keanggotaan_hari":360,"konsistensi_simpanan":0.5,"rasio_beban_angsuran":0.9}'::jsonb,
   '["riwayat pembayaran kurang","simpanan tidak memadai"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac', 68, 'C', 'review', 3500000,
   '{"ketepatan_bayar":0.7,"rasio_simpanan_pinjaman":0.5,"lama_keanggotaan_hari":240,"konsistensi_simpanan":0.7,"rasio_beban_angsuran":0.6}'::jsonb,
   '["simpanan cukup","perlu review manual"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('capp0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'lapp0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 85, 'A', 'approve', 10000000,
   '{"ketepatan_bayar":1.0,"rasio_simpanan_pinjaman":0.7,"lama_keanggotaan_hari":720,"konsistensi_simpanan":1.0,"rasio_beban_angsuran":0.2}'::jsonb,
   '["anggota senior","riwayat simpanan sangat konsisten","simpanan memadai"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('capp0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'lapp0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 78, 'B', 'approve', 8000000,
   '{"ketepatan_bayar":0.9,"rasio_simpanan_pinjaman":0.65,"lama_keanggotaan_hari":600,"konsistensi_simpanan":0.85,"rasio_beban_angsuran":0.35}'::jsonb,
   '["riwayat baik","simpanan memadai"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('capp0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'lapp0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 76, 'B', 'approve', 12000000,
   '{"ketepatan_bayar":0.88,"rasio_simpanan_pinjaman":0.58,"lama_keanggotaan_hari":480,"konsistensi_simpanan":0.88,"rasio_beban_angsuran":0.4}'::jsonb,
   '["riwayat pembayaran baik","konsistensi simpanan baik"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('capp0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'lapp0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 63, 'C', 'review', 3000000,
   '{"ketepatan_bayar":0.6,"rasio_simpanan_pinjaman":0.4,"lama_keanggotaan_hari":420,"konsistensi_simpanan":0.75,"rasio_beban_angsuran":0.65}'::jsonb,
   '["beban angsuran cukup tinggi","perlu review manual"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('capp0010-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'lapp0010-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 65, 'C', 'review', 3500000,
   '{"ketepatan_bayar":1.0,"rasio_simpanan_pinjaman":0.45,"lama_keanggotaan_hari":270,"konsistensi_simpanan":0.78,"rasio_beban_angsuran":0.55}'::jsonb,
   '["anggota baru","konsistensi simpanan baik"]'::jsonb, 'MOCK_ADINS_SCORING'),
  ('caaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', 'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', 75, 'B', 'approve', 7000000,
   '{"ketepatan_bayar":0.85,"rasio_simpanan_pinjaman":0.6,"lama_keanggotaan_hari":300,"konsistensi_simpanan":0.8,"rasio_beban_angsuran":0.4}'::jsonb,
   '["simpanan cukup","riwayat baik"]'::jsonb, 'MOCK_ADINS_SCORING');

-- ============================================================
-- LOANS
-- ============================================================
INSERT INTO loan (
  id, cooperative_id, application_id, member_id, principal, flat_rate_monthly, tenor_months, status, disbursed_at
)
VALUES
  -- Budi: aktif (5 bln berjalan, 4 sudah bayar)
  ('lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
   5000000, 1.50, 12, 'aktif', (NOW() - INTERVAL '5 months')::date),
  -- Asep: lunas (3jt/6mo, semua sudah bayar)
  ('loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'lapp0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4',
   3000000, 1.50, 6, 'lunas', (NOW() - INTERVAL '7 months')::date),
  -- Sri: aktif (4jt/12mo, 8 sudah bayar)
  ('loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'lapp0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5',
   4000000, 1.50, 12, 'aktif', (NOW() - INTERVAL '8 months')::date),
  -- Hendra: aktif (7jt/12mo, 4 sudah bayar)
  ('loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'lapp0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6',
   7000000, 1.50, 12, 'aktif', (NOW() - INTERVAL '4 months')::date),
  -- Eni: menunggak (2.5jt/6mo, 1 bayar, 2 terlambat)
  ('loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
   'lapp0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'a7a7a7a7-a7a7-a7a7-a7a7-a7a7a7a7a7a7',
   2500000, 1.50, 6, 'menunggak', (NOW() - INTERVAL '4 months')::date),
  -- Rina (Sawargi): menunggak
  ('lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', '22222222-2222-2222-2222-222222222222',
   'laaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad', 'dddddddd-dddd-dddd-dddd-dddddddddddd',
   6000000, 1.75, 12, 'menunggak', (NOW() - INTERVAL '4 months')::date);

-- ============================================================
-- INSTALLMENT SCHEDULES
-- ============================================================

-- Budi: 5jt, 1.5%, 12 bulan → interest=75000, principal=416666 (last=416674)
INSERT INTO installment_schedule (id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status)
VALUES
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa01', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1,  (NOW() - INTERVAL '4 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa02', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2,  (NOW() - INTERVAL '3 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa03', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3,  (NOW() - INTERVAL '2 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa04', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 4,  (NOW() - INTERVAL '1 months')::date, 416666, 75000, 491666, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa05', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 5,  NOW()::date,                           416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa06', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 6,  (NOW() + INTERVAL '1 months')::date,  416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa07', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 7,  (NOW() + INTERVAL '2 months')::date,  416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa08', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 8,  (NOW() + INTERVAL '3 months')::date,  416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa09', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 9,  (NOW() + INTERVAL '4 months')::date,  416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa10', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 10, (NOW() + INTERVAL '5 months')::date,  416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa11', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 11, (NOW() + INTERVAL '6 months')::date,  416666, 75000, 491666, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa12', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 12, (NOW() + INTERVAL '7 months')::date,  416674, 75000, 491674, 'belum_bayar');

-- Asep: 3jt, 1.5%, 6 bulan → interest=45000, principal=500000 (semua lunas)
INSERT INTO installment_schedule (id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status)
VALUES
  ('inst0004-aaaa-aaaa-aaaa-000000000001', 'loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, (NOW() - INTERVAL '6 months')::date, 500000, 45000, 545000, 'lunas'),
  ('inst0004-aaaa-aaaa-aaaa-000000000002', 'loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2, (NOW() - INTERVAL '5 months')::date, 500000, 45000, 545000, 'lunas'),
  ('inst0004-aaaa-aaaa-aaaa-000000000003', 'loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3, (NOW() - INTERVAL '4 months')::date, 500000, 45000, 545000, 'lunas'),
  ('inst0004-aaaa-aaaa-aaaa-000000000004', 'loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 4, (NOW() - INTERVAL '3 months')::date, 500000, 45000, 545000, 'lunas'),
  ('inst0004-aaaa-aaaa-aaaa-000000000005', 'loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 5, (NOW() - INTERVAL '2 months')::date, 500000, 45000, 545000, 'lunas'),
  ('inst0004-aaaa-aaaa-aaaa-000000000006', 'loan0004-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 6, (NOW() - INTERVAL '1 months')::date, 500000, 45000, 545000, 'lunas');

-- Sri: 4jt, 1.5%, 12 bulan → interest=60000, principal=333333 (last=333337); 8 lunas, 4 belum
INSERT INTO installment_schedule (id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status)
VALUES
  ('inst0005-aaaa-aaaa-aaaa-000000000001', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1,  (NOW() - INTERVAL '7 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000002', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2,  (NOW() - INTERVAL '6 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000003', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3,  (NOW() - INTERVAL '5 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000004', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 4,  (NOW() - INTERVAL '4 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000005', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 5,  (NOW() - INTERVAL '3 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000006', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 6,  (NOW() - INTERVAL '2 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000007', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 7,  (NOW() - INTERVAL '1 months')::date, 333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000008', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 8,  NOW()::date,                          333333, 60000, 393333, 'lunas'),
  ('inst0005-aaaa-aaaa-aaaa-000000000009', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 9,  (NOW() + INTERVAL '1 months')::date,  333333, 60000, 393333, 'belum_bayar'),
  ('inst0005-aaaa-aaaa-aaaa-000000000010', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 10, (NOW() + INTERVAL '2 months')::date,  333333, 60000, 393333, 'belum_bayar'),
  ('inst0005-aaaa-aaaa-aaaa-000000000011', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 11, (NOW() + INTERVAL '3 months')::date,  333333, 60000, 393333, 'belum_bayar'),
  ('inst0005-aaaa-aaaa-aaaa-000000000012', 'loan0005-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 12, (NOW() + INTERVAL '4 months')::date,  333337, 60000, 393337, 'belum_bayar');

-- Hendra: 7jt, 1.5%, 12 bulan → interest=105000, principal=583333 (last=583337); 4 lunas, 8 belum
INSERT INTO installment_schedule (id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status)
VALUES
  ('inst0006-aaaa-aaaa-aaaa-000000000001', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1,  (NOW() - INTERVAL '3 months')::date, 583333, 105000, 688333, 'lunas'),
  ('inst0006-aaaa-aaaa-aaaa-000000000002', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2,  (NOW() - INTERVAL '2 months')::date, 583333, 105000, 688333, 'lunas'),
  ('inst0006-aaaa-aaaa-aaaa-000000000003', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3,  (NOW() - INTERVAL '1 months')::date, 583333, 105000, 688333, 'lunas'),
  ('inst0006-aaaa-aaaa-aaaa-000000000004', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 4,  NOW()::date,                          583333, 105000, 688333, 'lunas'),
  ('inst0006-aaaa-aaaa-aaaa-000000000005', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 5,  (NOW() + INTERVAL '1 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000006', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 6,  (NOW() + INTERVAL '2 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000007', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 7,  (NOW() + INTERVAL '3 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000008', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 8,  (NOW() + INTERVAL '4 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000009', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 9,  (NOW() + INTERVAL '5 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000010', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 10, (NOW() + INTERVAL '6 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000011', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 11, (NOW() + INTERVAL '7 months')::date,  583333, 105000, 688333, 'belum_bayar'),
  ('inst0006-aaaa-aaaa-aaaa-000000000012', 'loan0006-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 12, (NOW() + INTERVAL '8 months')::date,  583337, 105000, 688337, 'belum_bayar');

-- Eni: 2.5jt, 1.5%, 6 bulan → interest=37500, principal=416666 (last=416670); 1 lunas, 2 terlambat, 3 belum
INSERT INTO installment_schedule (id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status)
VALUES
  ('inst0007-aaaa-aaaa-aaaa-000000000001', 'loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, (NOW() - INTERVAL '3 months')::date, 416666, 37500, 454166, 'lunas'),
  ('inst0007-aaaa-aaaa-aaaa-000000000002', 'loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2, (NOW() - INTERVAL '2 months')::date, 416666, 37500, 454166, 'terlambat'),
  ('inst0007-aaaa-aaaa-aaaa-000000000003', 'loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3, (NOW() - INTERVAL '1 months')::date, 416666, 37500, 454166, 'terlambat'),
  ('inst0007-aaaa-aaaa-aaaa-000000000004', 'loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 4, NOW()::date,                          416666, 37500, 454166, 'belum_bayar'),
  ('inst0007-aaaa-aaaa-aaaa-000000000005', 'loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 5, (NOW() + INTERVAL '1 months')::date,  416666, 37500, 454166, 'belum_bayar'),
  ('inst0007-aaaa-aaaa-aaaa-000000000006', 'loan0007-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 6, (NOW() + INTERVAL '2 months')::date,  416670, 37500, 454170, 'belum_bayar');

-- Rina (Sawargi): 6jt, 1.75%, 12 bulan → interest=105000, principal=500000; 1 lunas, 2 terlambat, 9 belum
INSERT INTO installment_schedule (id, loan_id, period_no, due_date, principal_due, interest_due, total_due, status)
VALUES
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa13', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 1,  (NOW() - INTERVAL '3 months')::date, 500000, 105000, 605000, 'lunas'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa14', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 2,  (NOW() - INTERVAL '2 months')::date, 500000, 105000, 605000, 'terlambat'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa15', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 3,  (NOW() - INTERVAL '1 months')::date, 500000, 105000, 605000, 'terlambat'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa16', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 4,  NOW()::date,                          500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa17', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 5,  (NOW() + INTERVAL '1 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa18', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 6,  (NOW() + INTERVAL '2 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa19', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 7,  (NOW() + INTERVAL '3 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa20', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 8,  (NOW() + INTERVAL '4 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa21', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 9,  (NOW() + INTERVAL '5 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa22', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 10, (NOW() + INTERVAL '6 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa23', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 11, (NOW() + INTERVAL '7 months')::date,  500000, 105000, 605000, 'belum_bayar'),
  ('in11111-aaaa-aaaa-aaaa-aaaaaaaaaa24', 'lnnnnnn-aaaa-aaaa-aaaa-aaaaaaaaaaab', 12, (NOW() + INTERVAL '8 months')::date,  500000, 105000, 605000, 'belum_bayar');

-- ============================================================
-- PAYMENTS (INSERT-ONLY)
-- ============================================================
INSERT INTO payment (id, schedule_id, amount, penalty, paid_at)
VALUES
  -- Budi: 4 payments (period 1-4)
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa01', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa01', 491666, 0, NOW() - INTERVAL '4 months'),
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa02', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa02', 491666, 0, NOW() - INTERVAL '3 months'),
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa03', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa03', 491666, 0, NOW() - INTERVAL '2 months'),
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa04', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa04', 491666, 0, NOW() - INTERVAL '1 months'),
  -- Asep: 6 payments (all periods lunas)
  ('pay0004-aaaa-aaaa-aaaa-000000000001', 'inst0004-aaaa-aaaa-aaaa-000000000001', 545000, 0, NOW() - INTERVAL '6 months'),
  ('pay0004-aaaa-aaaa-aaaa-000000000002', 'inst0004-aaaa-aaaa-aaaa-000000000002', 545000, 0, NOW() - INTERVAL '5 months'),
  ('pay0004-aaaa-aaaa-aaaa-000000000003', 'inst0004-aaaa-aaaa-aaaa-000000000003', 545000, 0, NOW() - INTERVAL '4 months'),
  ('pay0004-aaaa-aaaa-aaaa-000000000004', 'inst0004-aaaa-aaaa-aaaa-000000000004', 545000, 0, NOW() - INTERVAL '3 months'),
  ('pay0004-aaaa-aaaa-aaaa-000000000005', 'inst0004-aaaa-aaaa-aaaa-000000000005', 545000, 0, NOW() - INTERVAL '2 months'),
  ('pay0004-aaaa-aaaa-aaaa-000000000006', 'inst0004-aaaa-aaaa-aaaa-000000000006', 545000, 0, NOW() - INTERVAL '1 months'),
  -- Sri: 8 payments (periods 1-8)
  ('pay0005-aaaa-aaaa-aaaa-000000000001', 'inst0005-aaaa-aaaa-aaaa-000000000001', 393333, 0, NOW() - INTERVAL '7 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000002', 'inst0005-aaaa-aaaa-aaaa-000000000002', 393333, 0, NOW() - INTERVAL '6 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000003', 'inst0005-aaaa-aaaa-aaaa-000000000003', 393333, 0, NOW() - INTERVAL '5 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000004', 'inst0005-aaaa-aaaa-aaaa-000000000004', 393333, 0, NOW() - INTERVAL '4 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000005', 'inst0005-aaaa-aaaa-aaaa-000000000005', 393333, 0, NOW() - INTERVAL '3 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000006', 'inst0005-aaaa-aaaa-aaaa-000000000006', 393333, 0, NOW() - INTERVAL '2 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000007', 'inst0005-aaaa-aaaa-aaaa-000000000007', 393333, 0, NOW() - INTERVAL '1 months'),
  ('pay0005-aaaa-aaaa-aaaa-000000000008', 'inst0005-aaaa-aaaa-aaaa-000000000008', 393333, 0, NOW()),
  -- Hendra: 4 payments (periods 1-4)
  ('pay0006-aaaa-aaaa-aaaa-000000000001', 'inst0006-aaaa-aaaa-aaaa-000000000001', 688333, 0, NOW() - INTERVAL '3 months'),
  ('pay0006-aaaa-aaaa-aaaa-000000000002', 'inst0006-aaaa-aaaa-aaaa-000000000002', 688333, 0, NOW() - INTERVAL '2 months'),
  ('pay0006-aaaa-aaaa-aaaa-000000000003', 'inst0006-aaaa-aaaa-aaaa-000000000003', 688333, 0, NOW() - INTERVAL '1 months'),
  ('pay0006-aaaa-aaaa-aaaa-000000000004', 'inst0006-aaaa-aaaa-aaaa-000000000004', 688333, 0, NOW()),
  -- Eni: 1 payment (period 1), periods 2-3 terlambat tidak dibayar
  ('pay0007-aaaa-aaaa-aaaa-000000000001', 'inst0007-aaaa-aaaa-aaaa-000000000001', 454166, 0, NOW() - INTERVAL '3 months'),
  -- Rina (Sawargi): 1 payment (period 1)
  ('pppppp-aaaa-aaaa-aaaa-aaaaaaaaaa05', 'in11111-aaaa-aaaa-aaaa-aaaaaaaaaa13', 605000, 0, NOW() - INTERVAL '3 months');

COMMIT;
