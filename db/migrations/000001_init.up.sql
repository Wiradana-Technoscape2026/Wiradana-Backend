CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE cooperative (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  type text NOT NULL DEFAULT 'simpan_pinjam',
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE member (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  nik text NOT NULL,
  full_name text NOT NULL,
  address text,
  birth_date date,
  status text NOT NULL DEFAULT 'aktif' CHECK (status IN ('aktif','nonaktif','keluar')),
  custom_attributes jsonb NOT NULL DEFAULT '{}',
  joined_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (cooperative_id, nik)
);

CREATE TABLE app_user (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  member_id uuid REFERENCES member(id),
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  role text NOT NULL CHECK (role IN ('pengurus','anggota','admin','partner'))
);

CREATE TABLE savings_transaction (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  member_id uuid NOT NULL REFERENCES member(id),
  savings_type text NOT NULL CHECK (savings_type IN ('pokok','wajib','sukarela')),
  direction text NOT NULL CHECK (direction IN ('setor','tarik')),
  amount bigint NOT NULL CHECK (amount > 0),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_savings_member ON savings_transaction(member_id, savings_type);

CREATE TABLE loan_config (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  flat_rate_monthly numeric(5,2) NOT NULL DEFAULT 1.5,
  max_plafond bigint NOT NULL DEFAULT 20000000,
  penalty_daily bigint NOT NULL DEFAULT 0
);

CREATE TABLE loan_application (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  member_id uuid NOT NULL REFERENCES member(id),
  approved_by uuid REFERENCES app_user(id),
  amount bigint NOT NULL CHECK (amount > 0),
  tenor_months int NOT NULL CHECK (tenor_months > 0),
  purpose text,
  status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected')),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_app_status ON loan_application(cooperative_id, status);

CREATE TABLE credit_assessment (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  application_id uuid NOT NULL UNIQUE REFERENCES loan_application(id),
  score int NOT NULL,
  grade text NOT NULL CHECK (grade IN ('A','B','C','D')),
  recommendation text NOT NULL CHECK (recommendation IN ('approve','review','reject')),
  limit_suggested bigint NOT NULL,
  features jsonb NOT NULL DEFAULT '{}',
  reasons jsonb NOT NULL DEFAULT '[]',
  source text NOT NULL DEFAULT 'MOCK_ADINS_SCORING'
);

CREATE TABLE loan (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  application_id uuid NOT NULL UNIQUE REFERENCES loan_application(id),
  member_id uuid NOT NULL REFERENCES member(id),
  principal bigint NOT NULL,
  flat_rate_monthly numeric(5,2) NOT NULL,
  tenor_months int NOT NULL,
  status text NOT NULL DEFAULT 'aktif' CHECK (status IN ('aktif','lunas','menunggak')),
  disbursed_at date NOT NULL DEFAULT current_date
);

CREATE TABLE installment_schedule (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  loan_id uuid NOT NULL REFERENCES loan(id),
  period_no int NOT NULL,
  due_date date NOT NULL,
  principal_due bigint NOT NULL,
  interest_due bigint NOT NULL,
  total_due bigint NOT NULL,
  status text NOT NULL DEFAULT 'belum_bayar' CHECK (status IN ('belum_bayar','lunas','terlambat'))
);
CREATE INDEX idx_inst_loan ON installment_schedule(loan_id, status);

CREATE TABLE payment (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  schedule_id uuid NOT NULL REFERENCES installment_schedule(id),
  amount bigint NOT NULL CHECK (amount > 0),
  penalty bigint NOT NULL DEFAULT 0,
  paid_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE shu_period (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  year int NOT NULL,
  total_shu bigint NOT NULL,
  pct_jasa_modal numeric(5,2) NOT NULL,
  pct_jasa_usaha numeric(5,2) NOT NULL,
  status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','final'))
);

CREATE TABLE shu_distribution (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  shu_period_id uuid NOT NULL REFERENCES shu_period(id),
  member_id uuid NOT NULL REFERENCES member(id),
  jasa_modal bigint NOT NULL,
  jasa_usaha bigint NOT NULL,
  total_shu bigint NOT NULL
);

CREATE TABLE coop_module (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  module_key text NOT NULL,
  enabled boolean NOT NULL DEFAULT false,
  UNIQUE (cooperative_id, module_key)
);
