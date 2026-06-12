BEGIN;

INSERT INTO cooperative (id, name, type, created_at)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'Koperasi Wiradana', 'simpan_pinjam', NOW()),
  ('22222222-2222-2222-2222-222222222222', 'Koperasi Sawargi', 'simpan_pinjam', NOW());

INSERT INTO member (
  id,
  cooperative_id,
  nik,
  full_name,
  address,
  birth_date,
  status,
  custom_attributes,
  joined_at
)
VALUES
  (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    '11111111-1111-1111-1111-111111111111',
    '3201010101010001',
    'Budi Santoso',
    'Jl. Melati No. 10, Bandung',
    '1990-01-15',
    'aktif',
    '{}'::jsonb,
    NOW() - INTERVAL '18 months'
  ),
  (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '11111111-1111-1111-1111-111111111111',
    '3201010101010002',
    'Siti Aminah',
    'Jl. Kenanga No. 21, Bandung',
    '1992-06-20',
    'aktif',
    '{}'::jsonb,
    NOW() - INTERVAL '12 months'
  ),
  (
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    '11111111-1111-1111-1111-111111111111',
    '3201010101010003',
    'Dedi Rahman',
    'Jl. Mawar No. 5, Bandung',
    '1988-11-02',
    'aktif',
    '{}'::jsonb,
    NOW() - INTERVAL '8 months'
  ),
  (
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    '22222222-2222-2222-2222-222222222222',
    '3302020202020001',
    'Rina Lestari',
    'Jl. Anggrek No. 7, Bogor',
    '1994-03-11',
    'aktif',
    '{}'::jsonb,
    NOW() - INTERVAL '10 months'
  ),
  (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    '22222222-2222-2222-2222-222222222222',
    '3302020202020002',
    'Agus Pratama',
    'Jl. Teratai No. 14, Bogor',
    '1991-09-28',
    'aktif',
    '{}'::jsonb,
    NOW() - INTERVAL '6 months'
  );

INSERT INTO app_user (
  id,
  cooperative_id,
  member_id,
  email,
  password_hash,
  role
)
VALUES
  (
    '99999999-9999-9999-9999-999999999991',
    '11111111-1111-1111-1111-111111111111',
    NULL,
    'pengurus@wiradana.local',
    '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6',
    'pengurus'
  ),
  (
    '99999999-9999-9999-9999-999999999992',
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'budi.santoso@wiradana.local',
    '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6',
    'anggota'
  ),
  (
    '99999999-9999-9999-9999-999999999993',
    '11111111-1111-1111-1111-111111111111',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'siti.aminah@wiradana.local',
    '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6',
    'anggota'
  ),
  (
    '99999999-9999-9999-9999-999999999994',
    '22222222-2222-2222-2222-222222222222',
    NULL,
    'pengurus.sawargi@wiradana.local',
    '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6',
    'pengurus'
  ),
  (
    '99999999-9999-9999-9999-999999999995',
    '22222222-2222-2222-2222-222222222222',
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'rina.lestari@wiradana.local',
    '$2a$10$hlbA.jBgJtM6p0n8rnipKOiYIRrm37By7yPApkMBR/fUzlBJ5PMc6',
    'anggota'
  );

INSERT INTO loan_config (
  id,
  cooperative_id,
  flat_rate_monthly,
  max_plafond,
  penalty_daily
)
VALUES
  (
    '88888888-8888-8888-8888-888888888881',
    '11111111-1111-1111-1111-111111111111',
    1.50,
    20000000,
    0
  ),
  (
    '88888888-8888-8888-8888-888888888882',
    '22222222-2222-2222-2222-222222222222',
    1.75,
    15000000,
    0
  );

INSERT INTO coop_module (
  id,
  cooperative_id,
  module_key,
  enabled
)
VALUES
  (
    '77777777-7777-7777-7777-777777777771',
    '11111111-1111-1111-1111-111111111111',
    'simpan_pinjam',
    TRUE
  ),
  (
    '77777777-7777-7777-7777-777777777772',
    '11111111-1111-1111-1111-111111111111',
    'inventory',
    FALSE
  ),
  (
    '77777777-7777-7777-7777-777777777773',
    '22222222-2222-2222-2222-222222222222',
    'simpan_pinjam',
    TRUE
  ),
  (
    '77777777-7777-7777-7777-777777777774',
    '22222222-2222-2222-2222-222222222222',
    'inventory',
    FALSE
  );

INSERT INTO savings_transaction (
  id,
  cooperative_id,
  member_id,
  savings_type,
  direction,
  amount,
  created_at
)
VALUES
  (
    '66666666-6666-6666-6666-666666666661',
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'pokok',
    'setor',
    50000,
    NOW() - INTERVAL '18 months'
  ),
  (
    '66666666-6666-6666-6666-666666666662',
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'wajib',
    'setor',
    25000,
    NOW() - INTERVAL '2 months'
  ),
  (
    '66666666-6666-6666-6666-666666666663',
    '11111111-1111-1111-1111-111111111111',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'pokok',
    'setor',
    50000,
    NOW() - INTERVAL '12 months'
  ),
  (
    '66666666-6666-6666-6666-666666666664',
    '11111111-1111-1111-1111-111111111111',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'sukarela',
    'setor',
    150000,
    NOW() - INTERVAL '1 months'
  ),
  (
    '66666666-6666-6666-6666-666666666665',
    '22222222-2222-2222-2222-222222222222',
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'pokok',
    'setor',
    50000,
    NOW() - INTERVAL '10 months'
  ),
  (
    '66666666-6666-6666-6666-666666666666',
    '22222222-2222-2222-2222-222222222222',
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    'wajib',
    'setor',
    25000,
    NOW() - INTERVAL '6 months'
  );

COMMIT;
