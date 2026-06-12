# be_implementation.md — Backend Wiradana (Go + PostgreSQL, Clean Architecture)

> Plan implementasi backend, dibagi **BE-1** & **BE-2**. Patuhi `api_planning.md` (kontrak field), `ERD.md` (skema), `member-portal.md` (role), `ad_ins_mockup.md` (mock), `inventory-module.md` (modul opsional).
> Prinsip: bunga **flat**, saldo/stok = **SUM** (bukan kolom), transaksi keuangan **INSERT-ONLY**, tenant via **JWT** (`cooperative_id` tidak dari body), dua role **pengurus & anggota**.

---

## 0. Arsitektur (mengikuti github.com/khannedy/golang-clean-architecture)

**Alur:** Delivery (HTTP) → Use Case → Repository → DB; Use Case → Gateway → sistem eksternal (Ad Ins mock). Layer **Entity** (model bisnis) & **Model** (request/response).

**Tech stack (samakan dengan referензi):**
- **GoFiber** (HTTP), **GORM** (ORM), **Viper** (config `config.json`), **golang-migrate** (migrasi), **go-playground/validator** (validasi), **logrus** (log), **JWT** (`golang-jwt/jwt/v5`), **bcrypt**.
- DB: **PostgreSQL** (referensi pakai MySQL — kita pakai Postgres karena butuh `JSONB` & cocok dengan ERD). Kafka/worker **tidak dipakai** (skip; tidak relevan untuk 30 jam).

### Struktur folder
```
cmd/
  web/main.go                 # entrypoint HTTP server
internal/
  config/                     # viper, fiber, gorm, logrus, validator bootstrap
  entity/                     # struct tabel (GORM models) — sesuai ERD
  model/                      # request & response DTO — sesuai api_planning §2
    converter/                # entity <-> model
  repository/                 # akses DB per entity (CRUD + query SUM)
  usecase/                    # business logic (scoring, schedule, shu, dll)
  gateway/
    adins/                    # mock OCR & Credit Scoring (lihat ad_ins_mockup.md)
  delivery/http/
    controller/               # handler per modul
    middleware/               # auth (JWT->ctx), RequireRole, RequireModule
    route/                    # daftar route + guard
db/migrations/                # 001_init.sql, 005_inventory.sql, ...
config.json
```

**Aturan untuk AI agent:**
1. Setiap controller ambil `cooperative_id` & `role` (& `member_id`) dari `c.Locals(...)` yang di-set middleware dari JWT.
2. SEMUA query repository difilter `cooperative_id`.
3. Response pakai envelope `{success, data}` / `{success:false, error}`.
4. Field request/response disalin PERSIS dari `api_planning.md` — jangan ubah nama.

---

## 1. Pembagian Kerja

| | BE-1 | BE-2 |
|---|---|---|
| Owner | config bootstrap, migrasi/entity, Auth+RBAC, Members, OCR gateway, Savings, Dashboard, Portal anggota | Loan (config, application, scoring gateway, approve, schedule), Payments, SHU, Modules, Inventory (Tier 3) |
| Catatan | jam 0–2 buat entity & migrasi bareng | tunggu skema dari BE-1, lalu paralel |

---

## 2. Skema DB (db/migrations/001_init.sql) — OWNER: BE-1

Sesuai ERD.md. Tipe patokan (uang = `bigint` integer Rupiah):

```sql
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
  member_id uuid REFERENCES member(id),          -- NULL untuk pengurus
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
  flat_rate_monthly numeric(5,2) NOT NULL DEFAULT 1.5,   -- persen/bulan
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
```
> Migrasi inventory (`005_inventory.sql`) ada di `inventory-module.md §1` — jalankan hanya jika mengerjakan modul Tier 3.

---

## 3. Auth & RBAC (OWNER: BE-1, PRIORITAS — semua bergantung)

### 3.1 Middleware auth
- Parse `Authorization: Bearer`, verifikasi JWT. Klaim wajib: `user_id`, `cooperative_id`, `role`, `member_id` (boleh kosong untuk pengurus).
- Set ke `c.Locals`. Reject 401 jika invalid.

### 3.2 RequireRole & RequireModule
```go
func RequireRole(roles ...string) fiber.Handler { /* cek c.Locals("role") ∈ roles, else 403 */ }
func RequireModule(key string) fiber.Handler   { /* cek coop_module(key).enabled, else 403 */ }
```
- Endpoint pengurus → `RequireRole("pengurus")`.
- Endpoint `/portal/*` → `RequireRole("anggota")`.
- Endpoint `/inventory/*` → `RequireRole("pengurus")` + `RequireModule("inventory")`.

### 3.3 POST /auth/login
- Cari user by email, bcrypt compare. Sukses → JWT (exp 24h) berisi klaim di 3.1. Return `{token, user}` (user termasuk `role` & `member_id`).
- Matriks akses lengkap: `member-portal.md §2`.

---

## 4. BE-1 — Tugas Detail

### 4.1 Members (api_planning §3.2) [pengurus]
- `GET /members?search=&status=`: filter `cooperative_id`; search LIKE `full_name`/`nik`; sertakan `savings_summary`.
- `POST /members`: validasi NIK 16 digit & unik per koperasi (409 jika dobel). `custom_attributes` default `{}`.
  - **Opsional:** body `create_portal_account:true` → buat `app_user` role `anggota`, `member_id` terisi, password awal = NIK (atau default). (lihat member-portal §1.2)
- `PUT /members/:id`: tidak boleh ubah `nik`.

### 4.2 OCR Gateway (mock LiteDMS) (ad_ins_mockup §1)
- Interface `OCRGateway.ExtractKTP(image) KTPResult` di `gateway/adins`.
- Endpoint `POST /integrations/adins/ocr/ktp` (multipart). Hybrid: Tesseract → fallback map. `source:"MOCK_LITEDMS"`.

### 4.3 Savings (INSERT-ONLY) (api_planning §3.3) [pengurus]
- **savings_summary:**
  ```sql
  SELECT savings_type,
         SUM(CASE WHEN direction='setor' THEN amount ELSE -amount END) AS saldo
  FROM savings_transaction WHERE member_id=$1 GROUP BY savings_type;
  ```
- `POST /members/:id/savings` rules:
  - `pokok`+`setor` hanya sekali per member (else 409).
  - `pokok`/`wajib`+`tarik` → 409. `amount > 0`. Tidak ada update/delete.

### 4.4 Dashboard (api_planning §3.8) [pengurus]
- `total_members` (aktif), `total_savings` (SUM setor−tarik), `active_loans`, `overdue_loans`.
- `notifications`: installment `due_date` ≤3 hari & belum lunas.

### 4.5 Portal Anggota (api_planning §3.8 + member-portal §4) [anggota]
- Semua endpoint pakai `member_id` dari JWT — JANGAN dari body/param.
- `GET /portal/me`: member + savings_summary + savings_recent.
- `GET /portal/loans`, `GET /portal/loans/:id` (403 jika bukan miliknya).
- `GET /portal/shu`: estimasi/realisasi SHU dirinya.
- `POST /portal/loan-applications`: jalankan flow origination (4.7) dengan member_id dari token; hasil tetap `pending`.
- `GET /portal/loan-applications`: daftar pengajuan sendiri.

---

## 5. BE-2 — Tugas Detail

### 5.1 Loan Config (api_planning §3.4) [pengurus]
- `GET/PUT /loan-config` untuk koperasi aktif. Seed 1 row.

### 5.2 Scoring Gateway (mock Credit Scoring) (ad_ins_mockup §2)
- Interface `ScoringGateway.Score(in) ScoringResult` di `gateway/adins/scoring_mock.go`.
- **Hitung features dari data internal** (di use case sebelum memanggil gateway):
  ```
  ketepatan_bayar         = installment lunas tepat waktu / total jatuh tempo (default 1.0 jika belum pernah)
  rasio_simpanan_pinjaman = total_simpanan / amount_diajukan        (cap 1.0)
  lama_keanggotaan_bulan  = bulan sejak joined_at                   (normalisasi cap 36)
  konsistensi_simpanan    = bulan ada setoran wajib / bulan keanggotaan
  rasio_beban_angsuran    = angsuran_per_bulan / (total_simpanan/12 + 1)   (cap 1.0)
  ```
- **Skor & grade & limit & reasons:** rumus lengkap + contoh di `ad_ins_mockup.md §2.3–2.4`.
- Endpoint `POST /integrations/adins/credit-scoring` membungkus gateway (tetap diekspos untuk demo "integration point").

### 5.3 Loan Application (Origination) (api_planning §3.5) [pengurus & via portal]
- `POST /loan-applications` (pengurus) / `POST /portal/loan-applications` (anggota, member_id dari JWT):
  1. validasi member `aktif`, `amount <= max_plafond`.
  2. hitung features → `ScoringGateway.Score` → simpan `credit_assessment`.
  3. return application + assessment embed.
- `POST /:id/approve` [pengurus, hanya status `pending`]: set `approved_by`+`approved`; **buat LOAN + generate schedule** (5.4). Return `{application, loan}`.
- `POST /:id/reject` [pengurus].

### 5.4 Schedule Generator (bunga FLAT) — `usecase/schedule.go`
Snapshot `principal`,`flat_rate_monthly`,`tenor` ke `loan`, lalu:
```
bunga_per_bulan = round(principal * flat_rate_monthly / 100)   // tetap
pokok_per_bulan = round(principal / tenor)
for i in 1..tenor:
   principal_due = (i<tenor) ? pokok_per_bulan : principal - pokok_per_bulan*(tenor-1)
   interest_due  = bunga_per_bulan
   total_due     = principal_due + interest_due
   due_date      = disbursed_at + i bulan
```
**Contoh** principal=5.000.000, rate=1.5%, tenor=12 → bunga 75.000/bln tetap, pokok 416.667 (periode 12 = 416.663). Total kembali 5.900.000.

### 5.5 Payments (Collection, INSERT-ONLY) (api_planning §3.6) [pengurus]
- `POST /installments/:id/pay`:
  1. INSERT payment.
  2. `paid_amount = SUM(payment.amount)`; ≥ `total_due` → `lunas`; lewat `due_date` & belum lunas → `terlambat`.
  3. recompute `loan.status`: semua lunas → `lunas`; ada terlambat → `menunggak`; else `aktif`.
- `outstanding` = `principal − SUM(principal_due WHERE status='lunas')`.

### 5.6 SHU (api_planning §3.7) — `usecase/shu.go` [pengurus]
```
shu_modal_pool = total_shu * pct_jasa_modal/100
shu_usaha_pool = total_shu * pct_jasa_usaha/100
per anggota:
  jasa_modal = round(simpanan_anggota / total_simpanan_koperasi * shu_modal_pool)
  jasa_usaha = round(jasa_pinjaman_anggota / total_jasa_pinjaman * shu_usaha_pool)  // 0 jika tak meminjam
  total_shu  = jasa_modal + jasa_usaha
```
**Contoh:** total_shu=40jt, modal 20%; simpanan 8jt dari total 200jt → jasa_modal = 8/200 × 8.000.000 = 320.000.

### 5.7 Modules (api_planning §3.9) [pengurus]
- `GET /modules`, `PUT /modules/:key` (toggle). Seed 2 row: `simpan_pinjam`(true), `inventory`(false).

### 5.8 Inventory (Tier 3, opsional) — lihat `inventory-module.md §4`
- Hanya jika core (Tier 1+2) solid. Field-def CRUD, product CRUD (validasi `custom_attributes` vs field_def), movement insert-only (stok = SUM, keluar>stok → 409). Guard `RequireModule("inventory")`.

---

## 6. Seed Data (OWNER: BE-1 + AI agent)
- 1 cooperative "Padiwangi", 1 user **pengurus**, 1–2 user **anggota** (role anggota, member_id terisi) untuk demo portal.
- (Opsional) cooperative kedua untuk cerita tenant-aware.
- ±20 member, `joined_at` 6–24 bln lalu; simpanan wajib bulanan + pokok sekali; sebagian sukarela.
- 3–5 pinjaman: sebagian lunas, sebagian aktif, 1 menunggak.
- `loan_config` default (1.5%, 20jt). `coop_module` (simpan_pinjam on, inventory off).

---

## 7. Definition of Done BE
- [ ] Login → JWT `cooperative_id`,`role`,`member_id`.
- [ ] RBAC: anggota ditolak (403) di endpoint pengurus; portal pakai member_id dari token.
- [ ] Member CRUD + NIK unik + savings_summary benar.
- [ ] Savings insert-only + aturan pokok/wajib.
- [ ] Application → scoring (gateway) → assessment embed.
- [ ] Approve → loan + schedule flat benar (contoh 5.4).
- [ ] Pay → status installment & loan recompute benar.
- [ ] SHU calculate cocok contoh 5.6.
- [ ] Portal: me/loans/shu/loan-applications jalan & terbatas ke diri sendiri.
- [ ] Dashboard angka konsisten. Semua query terfilter `cooperative_id`.
- [ ] (Tier 3) Inventory modul + guard jalan.

---

## 8. Unit Test Wajib (AI agent bantu — di folder `test/`)
- `schedule_test`: total_due×tenor = principal + bunga×tenor; pembulatan periode terakhir benar.
- `scoring_test`: input ekstrem → grade A & D; contoh ad_ins_mockup §2.4 → skor ~80/A.
- `shu_test`: cocok contoh 5.6 (320.000).
- `savings_test`: pokok kedua → 409; tarik wajib → 409.
- `rbac_test`: anggota akses /members → 403; portal pakai member lain → tetap data sendiri.
