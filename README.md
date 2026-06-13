# Wiradana — Backend

> **TechnoScape 2026** · BNCC × Ad Ins · Tim S2U
>
> Digitalisasi penuh operasional Koperasi Simpan Pinjam (KSP): dari pencatatan anggota, simpanan, pinjaman, hingga SHU — menggantikan Excel dengan satu sistem terpusat yang dapat diaudit.

**Live Demo:** [https://wiradana.vercel.app/](https://wiradana.vercel.app/)

---

## Daftar Isi

1. [Gambaran Proyek](#1-gambaran-proyek)
2. [Akun Demo](#2-akun-demo)
3. [Arsitektur — Clean Architecture](#3-arsitektur--clean-architecture)
4. [Struktur Direktori](#4-struktur-direktori)
5. [Tech Stack](#5-tech-stack)
6. [Design System Backend](#6-design-system-backend)
7. [Fitur yang Dibangun](#7-fitur-yang-dibangun)
8. [Integrasi Ad Ins — Framing Jujur](#8-integrasi-ad-ins--framing-jujur)
9. [Database Design](#9-database-design)
10. [API Endpoints](#10-api-endpoints)
11. [Cara Menjalankan](#11-cara-menjalankan)
12. [Unit Test](#12-unit-test)

---

## 1. Gambaran Proyek

**Wiradana** adalah platform SaaS multi-tenant untuk digitalisasi Koperasi Simpan Pinjam (KSP), dibangun dalam sprint 30 jam untuk hackathon TechnoScape 2026. Sistem menggantikan pencatatan manual berbasis Excel yang menjadi *pain point* literal koperasi Padiwangi:

> *"Data-data pengurus masih disimpan ke dalam Excel dan tercatat secara manual sehingga tidak lengkap."*

**Full flow KSP berjalan end-to-end:**

```
Registrasi anggota (+ KTP OCR)
  → Simpanan (pokok / wajib / sukarela, ledger insert-only)
    → Pengajuan pinjaman (+ credit scoring otomatis)
      → Approval pengurus → Jadwal angsuran (flat rate, auto-generated)
        → Pembayaran angsuran → Deteksi tunggakan
          → SHU otomatis per anggota
            → Portal mandiri anggota
```

Setiap langkah di atas **dibangun nyata** (bukan mock/prototipe UI) dan berjalan di database PostgreSQL yang sesungguhnya.

---

## 2. Akun Demo

> Registrasi pengurus dan anggota baru **tidak tersedia di UI** untuk menjaga data demo tetap konsisten. Gunakan akun berikut untuk mencoba seluruh fitur.

**URL:** [https://wiradana.vercel.app/](https://wiradana.vercel.app/)

| Role | Identifier | Password |
|---|---|---|
| Pengurus | `pengurus@padiwangi.id` | `password123` |
| Anggota | `3174096112900001` | `password123` |

- **Pengurus** — akses penuh: kelola anggota, simpanan, pinjaman, approval, SHU, laporan, dan modul inventory.
- **Anggota** (login pakai NIK `3174096112900001`) — portal mandiri: lihat saldo simpanan, jadwal angsuran, estimasi SHU, dan ajukan pinjaman baru.

---

## 3. Arsitektur — Clean Architecture

Backend mengikuti pola [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture), yang mengadaptasi **Clean Architecture** Robert C. Martin untuk Go.

### Prinsip Inti

```
┌─────────────────────────────────────────────┐
│              Delivery Layer                  │
│   (HTTP Handler, Route, Middleware)          │
│   Tahu: Fiber, JWT, HTTP status code        │
│   Tidak tahu: SQL, logika bisnis            │
├─────────────────────────────────────────────┤
│              Usecase Layer                   │
│   (Logika bisnis murni)                     │
│   Tahu: domain rules, interface             │
│   Tidak tahu: HTTP, SQL, ORM                │
├─────────────────────────────────────────────┤
│           Repository Layer                   │
│   (Akses data via GORM)                     │
│   Tahu: SQL, GORM                           │
│   Tidak tahu: HTTP, logika bisnis           │
├─────────────────────────────────────────────┤
│           Gateway Layer                      │
│   (Integrasi eksternal: OCR, Scoring, WA)   │
│   Interface di sini → implementasi bisa swap│
└─────────────────────────────────────────────┘
```

### Aturan Ketergantungan (Dependency Rule)

Setiap lapisan **hanya** bergantung ke arah dalam:

- `delivery/http` → memanggil `usecase` via interface
- `usecase` → memanggil `repository` dan `gateway` via interface
- `repository` → mengakses DB via GORM
- `gateway` → mengakses layanan eksternal (OCR API, WA API)

**Konsekuensi nyata:** mengganti implementasi (contoh: mock Ad Ins scoring → real Ad Ins API) = **hanya ubah binding di `main.go`**, tanpa sentuh usecase sama sekali.

```go
// cmd/web/main.go — Dependency injection manual, tanpa framework DI
scoringGateway := adins.NewMockScoringGateway()
// Di produksi, baris di atas diganti:
// scoringGateway := adins.NewRealAdInsScoringGateway(apiKey, baseURL)
// Usecase tidak berubah sama sekali.

loanUsecase := usecase.NewLoanUsecase(loanRepo, memberRepo, scoringGateway)
```

### Aliran Data per Request

```
HTTP Request
  → Fiber Router
    → Middleware (AuthMiddleware → RequireRole → RequireModule)
      → Controller (parse request, panggil usecase)
        → Usecase (jalankan aturan bisnis, panggil repo/gateway)
          → Repository (query GORM → PostgreSQL)
          → Gateway (panggil API eksternal)
        ← Usecase (kembalikan domain struct / error)
      ← Controller (wrap ke envelope JSON)
  ← HTTP Response { "success": true, "data": {...} }
```

---

## 4. Struktur Direktori

```
backend/
├── cmd/
│   ├── web/
│   │   └── main.go               # Entry point: init config, DB, DI, start Fiber
│   └── migrate/
│       └── main.go               # Migration runner (golang-migrate)
├── internal/
│   ├── config/                   # Viper config, Fiber setup, GORM init, logger, validator
│   ├── entity/                   # GORM model (struct dengan tag ORM & JSON)
│   │   ├── cooperative.go
│   │   ├── app_user.go
│   │   ├── member.go
│   │   ├── savings_transaction.go
│   │   ├── loan_config.go
│   │   ├── loan_application.go
│   │   ├── credit_assessment.go
│   │   ├── loan.go
│   │   ├── installment_schedule.go
│   │   ├── payment.go
│   │   ├── shu_period.go
│   │   ├── shu_distribution.go
│   │   ├── coop_module.go
│   │   ├── inventory.go          # InventoryFieldDef, Product, Movement
│   │   ├── notification_log.go
│   │   ├── idempotency_key.go
│   │   └── user_cooperative_membership.go
│   ├── model/                    # Request/Response DTO + Converter
│   │   ├── auth_model.go
│   │   ├── member_model.go
│   │   ├── savings_model.go
│   │   ├── loan_model.go
│   │   ├── shu_model.go
│   │   ├── sync_model.go
│   │   ├── dashboard_model.go
│   │   └── converter/            # entity → response model
│   ├── repository/               # GORM implementasi (satu file per domain)
│   ├── usecase/                  # Logika bisnis (satu file per domain)
│   │   ├── schedule.go           # GenerateSchedule (flat rate)
│   │   ├── schedule_test.go      # Unit test angsuran flat
│   │   ├── shu_test.go           # Unit test rumus SHU
│   │   ├── savings_test.go       # Unit test aturan simpanan
│   │   └── sync_test.go          # Unit test offline sync
│   ├── delivery/
│   │   └── http/
│   │       ├── controller/       # Fiber handler (satu file per domain)
│   │       │   └── response.go   # Helper OK() / Fail() — envelope terpusat
│   │       ├── middleware/
│   │       │   ├── auth.go       # JWT parse → set Locals
│   │       │   └── rbac_test.go  # Unit test middleware RBAC
│   │       └── route/
│   │           └── route.go      # Semua routing + DI wiring
│   ├── gateway/
│   │   ├── adins/
│   │   │   ├── adins.go          # Interface OCRGateway + ScoringGateway
│   │   │   ├── ktp_ocr_gateway.go # Real: api.co.id OCR integration
│   │   │   ├── ocr_mock.go       # Fallback mock OCR
│   │   │   ├── scoring_mock.go   # Mock credit scoring (deterministik)
│   │   │   └── scoring_mock_test.go
│   │   └── notification/
│   │       ├── notification.go   # Interface NotificationGateway
│   │       ├── mock.go           # Mock: log ke DB
│   │       └── whatsapp.go       # WhatsApp Cloud API sandbox
│   └── scheduler/
│       └── notification_scheduler.go  # Goroutine harian jam 08:00 WIB
├── db/
│   ├── migrations/               # SQL files (golang-migrate up/down)
│   └── seed_dummy.sql            # Seed: 1 koperasi, pengurus, anggota, data demo
├── go.mod
├── go.sum
└── config.json                   # Konfigurasi (jangan di-commit dengan secret)
```

---

## 5. Tech Stack

| Komponen | Teknologi | Versi |
|---|---|---|
| Language | Go | 1.25.6 |
| Web Framework | [GoFiber v2](https://github.com/gofiber/fiber) | v2.52.13 |
| ORM | [GORM](https://gorm.io) + `gorm.io/driver/postgres` | v1.31.1 |
| Database | PostgreSQL (Supabase) | — |
| Migrasi | [golang-migrate](https://github.com/golang-migrate/migrate) | v4.19.1 |
| Konfigurasi | [Viper](https://github.com/spf13/viper) | v1.21.0 |
| Validasi | [go-playground/validator](https://github.com/go-playground/validator) | v10.30.3 |
| JWT | [golang-jwt/v5](https://github.com/golang-jwt/jwt) | v5.3.1 |
| Password | `golang.org/x/crypto/bcrypt` | — |
| UUID | `github.com/google/uuid` | v1.6.0 |
| Logger | [Logrus](https://github.com/sirupsen/logrus) | v1.9.4 |
| Excel Export | [excelize/v2](https://github.com/xuri/excelize) | v2.10.1 |
| OCR Eksternal | [api.co.id](https://api.co.id) KTP OCR | — |
| Deploy | Railway (BE + DB) | — |

---

## 6. Design System Backend

### 5.1 Envelope Response

Semua endpoint menggunakan format envelope terpusat:

```json
// Sukses data tunggal
{ "success": true, "data": { ... } }

// Sukses list
{ "success": true, "data": [ ... ], "meta": { "total": 42 } }

// Error
{ "success": false, "error": { "code": "VALIDATION_ERROR", "message": "NIK sudah terdaftar" } }
```

Helper terpusat di `internal/delivery/http/controller/response.go` — tidak ada envelope inline di handler.

### 5.2 Kode Error Standar

| Code | HTTP | Kapan |
|---|---|---|
| `VALIDATION_ERROR` | 400 | Input tidak valid, format salah |
| `UNAUTHORIZED` | 401 | Token tidak ada / expired / invalid |
| `FORBIDDEN` | 403 | Role salah, atau akses data milik orang lain |
| `NOT_FOUND` | 404 | Resource tidak ada di tenant ini |
| `CONFLICT` | 409 | NIK duplikat, pokok sudah disetor, status salah |
| `INTERNAL` | 500 | Error tidak terduga |

### 5.3 Aturan Uang — Integer Rupiah

Semua nominal uang adalah **`int64` (Go) / `BIGINT` (PostgreSQL)**. Tidak ada `float64` atau `decimal` untuk uang. Pembulatan angsuran menggunakan `math.Round`.

```go
// Benar
type LoanApplication struct {
    Amount int64 `json:"amount"` // mis. 5000000 (Rp 5.000.000)
}
```

### 5.4 Multi-Tenancy — Tenant-Aware

Setiap tabel bisnis memiliki kolom `cooperative_id UUID NOT NULL`. `cooperative_id` **selalu** diambil dari klaim JWT — tidak pernah dari body request.

```go
// Di setiap handler:
coopID := c.Locals("cooperative_id").(string)
// Setiap query DB wajib filter WHERE cooperative_id = ?
```

### 5.5 JWT Claims

```go
type JWTClaims struct {
    UserID        string `json:"user_id"`
    CooperativeID string `json:"cooperative_id"`
    Role          string `json:"role"`      // "pengurus" | "anggota"
    MemberID      string `json:"member_id"` // "" jika pengurus
    jwt.RegisteredClaims
}
```

### 5.6 Middleware Stack

```
Fiber Router
  └── AuthMiddleware         parse token → set Locals (user_id, cooperative_id, role, member_id)
      └── RequireRole("pengurus")   guard endpoint pengurus
      └── RequireRole("anggota")    guard endpoint portal anggota
          └── RequireModule("inventory")  guard modul opsional
```

Implementasi: `internal/delivery/http/middleware/auth.go`

### 5.7 Insert-Only Ledger

Tabel transaksi keuangan (`savings_transaction`, `payment`, `inventory_movement`) adalah **insert-only**:
- Tidak ada endpoint UPDATE/DELETE
- Koreksi = insert transaksi pembalik (reversing entry)
- **Saldo & stok = `SUM(amount)`**, bukan kolom yang di-update

Menjawab pain "saldo tidak sinkron" sekaligus memberi audit trail gratis.

### 5.8 Enum — Lowercase String

Semua enum dikirim dan diterima sebagai string lowercase:

```
role:            pengurus | anggota
savings_type:    pokok | wajib | sukarela
direction:       setor | tarik
loan_status:     aktif | lunas | menunggak
app_status:      pending | approved | rejected
grade:           A | B | C | D
recommendation:  approve | review | reject
installment:     belum_bayar | lunas | terlambat
```

---

## 7. Fitur yang Dibangun

Semua fitur di bawah ini **dibangun nyata dan berjalan** — bukan prototipe UI atau mock data.

### 6.1 Auth & Multi-Cooperative Membership

- Login dengan JWT yang membawa `cooperative_id`, `role`, dan `member_id` dalam klaim
- Satu user bisa anggota di **lebih dari satu koperasi** (junction table `user_cooperative_membership`)
- Jika anggota punya >1 koperasi → response `requires_cooperative_selection: true` + list koperasi → client pilih → re-issue scoped JWT via `POST /auth/select-cooperative`
- bcrypt untuk password hashing; `password_hash` tidak pernah dikembalikan di response

### 6.2 Manajemen Anggota

- CRUD anggota dengan validasi NIK 16 digit, unik per koperasi
- Status keanggotaan: `aktif` / `nonaktif` / `keluar`
- `custom_attributes` JSONB — data tambahan per koperasi tanpa migrasi (Dynamic Entity Builder)
- `phone_number` opsional untuk keperluan notifikasi
- Ringkasan simpanan (`savings_summary`) otomatis terhitung per anggota

### 6.3 KTP OCR — Integrasi Real

- `POST /api/v1/integrations/adins/ocr/ktp` (multipart form)
- Menggunakan **real OCR API** dari [api.co.id](https://api.co.id) — bukan mock deterministik
- Response: NIK, nama lengkap, alamat, tanggal lahir, confidence, `source: "API_CO_ID"`
- Auto-fill form registrasi anggota di frontend, mengeliminasi input manual yang jadi akar "data tidak lengkap"

### 6.4 Manajemen Simpanan (Ledger)

- Catat setoran & penarikan untuk 3 jenis simpanan
- Aturan bisnis yang di-enforce di usecase:
  - Simpanan pokok hanya boleh disetor **satu kali** per anggota
  - Simpanan pokok & wajib **tidak bisa ditarik**
  - Penarikan sukarela tidak boleh melebihi saldo
- Saldo dihitung dengan `SUM` — tidak ada kolom saldo yang di-update
- `GET /savings/summary` — aggregate koperasi (total pokok/wajib/sukarela)
- `GET /savings/transactions` — semua transaksi dengan nama anggota, paginated + filterable

### 6.5 Konfigurasi Pinjaman

- Parameter pinjaman dikonfigurasi **per koperasi** (sesuai rapat anggota/AD·ART)
- Field: `flat_rate_monthly` (%), `max_plafond` (Rp), `penalty_daily` (Rp/hari)
- Perubahan config hanya berlaku untuk pinjaman **baru** — pinjaman lama tidak terpengaruh (rate di-snapshot ke tabel `loan` saat approve)

### 6.6 Pinjaman — Origination → Management → Collection

**Origination (pengajuan):**
- Validasi: anggota harus `aktif`, jumlah ≤ max_plafond, tenor valid
- Sistem otomatis menghitung 5 fitur scoring dari data internal
- Panggil mock Credit Scoring → simpan `CreditAssessment` (snapshot keputusan)
- Status pengajuan: `pending` → `approved` / `rejected`

**Management (setelah approve):**
- Snapshot `flat_rate_monthly` dari `LoanConfig` ke `Loan`
- Auto-generate jadwal angsuran (bulk insert semua periode sekaligus)
- **Rumus angsuran flat:**
  ```
  Bunga per bulan      = round(Pokok × flat_rate_monthly%)
  Angsuran Pokok       = Pokok / tenor_bulan
  Total per bulan      = Angsuran Pokok + Bunga (tetap setiap bulan)
  Periode terakhir     = sisa pokok agar total pas (bukan round)
  ```

**Collection (pembayaran):**
- `POST /installments/:id/pay` → INSERT Payment → recompute status angsuran & status pinjaman
- Mendukung pembayaran **parsial** (partial payment)
- Status pinjaman otomatis berubah: `aktif` → `menunggak` (ada keterlambatan) → `lunas` (semua selesai)

**Export:**
- `GET /loans/:id/export` → file Excel jadwal angsuran dengan styling (judul, summary card, tabel berwarna)

### 6.7 Credit Scoring — Rule-Based (Integration-Ready)

- `POST /api/v1/integrations/adins/credit-scoring` (juga dipanggil internal saat origination)
- Algoritma **deterministik** — input sama menghasilkan output sama
- 5 fitur dihitung dari data internal koperasi:

| Fitur | Bobot | Definisi |
|---|---|---|
| Riwayat ketepatan bayar | 35% | % angsuran tepat waktu dari pinjaman sebelumnya |
| Rasio simpanan / pinjaman | 25% | Total simpanan ÷ jumlah diajukan |
| Lama keanggotaan | 15% | Bulan sejak `joined_at`, dinormalisasi ke 36 bulan |
| Konsistensi simpanan wajib | 15% | % bulan tersetor dalam 12 bulan terakhir |
| Rasio beban angsuran | 10% | Angsuran baru ÷ kapasitas (proxy dari simpanan) |

```
Skor = Σ (fitur_ternormalisasi × bobot)    → 0–100
Grade A (≥80): Approve, limit penuh
Grade B (70–79): Approve, limit standar
Grade C (60–69): Review manual
Grade D (<60): Reject
```

Response menyertakan `source: "MOCK_ADINS_SCORING"` — jujur bahwa ini mock, *integration-ready* untuk swap ke Credit Scoring Ad Ins resmi di produksi.

### 6.8 SHU Otomatis (Sisa Hasil Usaha)

- Buat periode SHU dengan total SHU dan alokasi % (konfigurable per koperasi sesuai AD/ART)
- `POST /shu-periods/:id/calculate`:
  - Hitung total simpanan seluruh anggota aktif
  - Hitung total jasa pinjaman dalam periode
  - Distribusi per anggota: `jasa_modal` (proporsional simpanan) + `jasa_usaha` (proporsional bunga bayar)
  - Bulk insert `ShuDistribution` per anggota
- Menggantikan perhitungan SHU manual yang biasanya memakan waktu berminggu-minggu

### 6.9 Dashboard Pengurus

- Anggota aktif vs total
- Outstanding pinjaman aktif (jumlah + nilai)
- Angsuran jatuh tempo dalam 7 hari ke depan (dengan detail per angsuran)
- Pengajuan pinjaman pending (dengan grade scoring)
- Ringkasan simpanan koperasi (total 3 jenis)

### 6.10 Portal Anggota (Dua Role)

Endpoint `/portal/*` khusus role `anggota`:
- `GET /portal/me` — profil diri + ringkasan simpanan
- `GET /portal/loans` & `/portal/loans/:id` — daftar & detail pinjaman sendiri + jadwal angsuran
- `GET /portal/shu` — estimasi SHU sendiri
- `POST /portal/loan-applications` — ajukan pinjaman (masuk antrian approval pengurus)
- `GET /portal/loan-applications` — status pengajuan sendiri
- `GET /portal/cooperatives` — daftar koperasi yang diikuti (untuk fitur multi-koperasi)
- `GET /portal/dashboard?cooperative_id=...` — dashboard anggota per-koperasi

Backend **memaksa** `member_id` dari token — anggota tidak bisa mengakses data orang lain:
```go
memberID := c.Locals("member_id").(string) // dari JWT, bukan dari body
if loan.MemberID != memberID {
    return Fail(c, 403, "FORBIDDEN", "Bukan pinjaman Anda")
}
```

### 6.11 Modul Inventory — Plug-and-Play (Tier 3)

Modul tambahan yang bisa diaktifkan/dinonaktifkan pengurus (`PUT /modules/inventory`):

- **Dynamic Entity Builder via JSONB**: setiap koperasi mendefinisikan field kustom sendiri (`GET/POST/DELETE /inventory/field-defs`) tanpa migrasi — contoh: "Kadar Air" untuk koperasi beras
- `GET/POST /inventory/products` — master produk/komoditas dengan atribut kustom
- `POST /inventory/products/:id/movements` — catat masuk/keluar stok (insert-only)
- Stok = `SUM(masuk) − SUM(keluar)` — tidak ada kolom stok yang di-update
- Semua endpoint di-guard `RequireModule("inventory")`

### 6.12 Notifikasi

- Scheduler goroutine berjalan harian jam 08:00 WIB
- Trigger otomatis: H-3 jatuh tempo, H-1, hari-H, dan angsuran menunggak
- Event hook: konfirmasi setoran simpanan + konfirmasi pencairan pinjaman
- Dua mode (konfigurable via `WA_MODE`):
  - `mock`: log ke tabel `notification_log` (aman untuk demo offline)
  - `sandbox`: kirim ke WhatsApp Cloud API Meta ke nomor terdaftar
- `GET /notifications/logs` — riwayat notifikasi
- `POST /notifications/trigger` — paksa kirim untuk kebutuhan demo

### 6.13 Offline-First Sync

- `GET /api/v1/sync/pull?since=<RFC3339>` — full & delta pull, di-filter per role
- `POST /api/v1/sync/push` — batch mutations dengan idempotency key (deduplication aman)
- Mutation yang didukung: anggota, simpanan, pengajuan pinjaman, pembayaran, konfigurasi, dan inventory
- Unit test: duplicate, unknown type, missing field, invalid ID, cursor pull — semua pass

### 6.14 Laporan Excel

- `GET /loans/:id/export` — jadwal angsuran dalam format Excel dengan styling (header berwarna, summary card, tabel terformat)

---

## 8. Integrasi Ad Ins — Framing Jujur

Wiradana membangun **integration-ready architecture**: titik integrasi dirancang agar mudah disambungkan ke produk Ad Ins di produksi, dengan mock yang meniru kontrak input/output secara jujur.

| Komponen | Status | Source Field |
|---|---|---|
| KTP OCR / e-KYC | **Real** — api.co.id | `"API_CO_ID"` |
| Credit Scoring | **Mock deterministik** — rule-based | `"MOCK_ADINS_SCORING"` |
| WhatsApp Notifikasi | Mock log + sandbox opsional | `"MOCK_WA"` |

**Untuk men-swap mock → real (di produksi):**
```go
// cmd/web/main.go — hanya ubah baris ini:
// Dari:
scoringGateway := adins.NewMockScoringGateway()
// Ke:
scoringGateway := adins.NewRealAdInsScoringGateway(cfg.AdIns.APIKey, cfg.AdIns.BaseURL)
// Usecase tidak berubah sama sekali.
```

Interface gateway di `internal/gateway/adins/adins.go` adalah kontrak yang tidak boleh diubah — implementasi bisa diganti bebas.

---

## 9. Database Design

PostgreSQL dengan 15 tabel. Semua tabel bisnis memiliki `cooperative_id` (multi-tenancy). Transaksi keuangan bersifat insert-only.

### Entity Relationship (Ringkasan)

```
COOPERATIVE ─┬─── APP_USER ──────────── USER_COOPERATIVE_MEMBERSHIP
             ├─── MEMBER ──┬─────────── SAVINGS_TRANSACTION (insert-only)
             │             ├─────────── LOAN_APPLICATION ── CREDIT_ASSESSMENT
             │             │                    └── LOAN ── INSTALLMENT_SCHEDULE
             │             │                                        └── PAYMENT (insert-only)
             │             └─────────── SHU_DISTRIBUTION
             ├─── LOAN_CONFIG
             ├─── SHU_PERIOD ────────── SHU_DISTRIBUTION
             ├─── COOP_MODULE
             ├─── INVENTORY_FIELD_DEF
             ├─── INVENTORY_PRODUCT ─── INVENTORY_MOVEMENT (insert-only)
             └─── NOTIFICATION_LOG
```

### Keputusan Desain Utama

| Keputusan | Alasan |
|---|---|
| Saldo simpanan = SUM, bukan kolom | Insert-only → saldo selalu dapat ditelusuri, tidak pernah out-of-sync |
| Rate di-snapshot ke tabel `loan` saat approve | Perubahan config tidak mempengaruhi pinjaman lama (integritas historis) |
| `CREDIT_ASSESSMENT` terpisah dari `LOAN_APPLICATION` | Setiap keputusan approve/reject tersimpan dengan fitur + source — dapat dipertanggungjawabkan |
| UUID di-generate di Go (BeforeCreate), bukan dari DB | GORM tahu ID langsung setelah insert tanpa round-trip |
| Dua DSN: `dsn` (pooler) + `migration_dsn` (direct) | golang-migrate membutuhkan session-level advisory lock yang tidak kompatibel dengan pgbouncer transaction mode |

---

## 10. API Endpoints

### Public
```
POST /api/v1/auth/login
POST /api/v1/auth/register/pengurus
POST /api/v1/auth/select-cooperative
```

### Pengurus (Bearer token, role: pengurus)
```
GET    /api/v1/members
POST   /api/v1/members
GET    /api/v1/members/:id
PUT    /api/v1/members/:id
GET    /api/v1/members/:id/savings
POST   /api/v1/members/:id/savings
GET    /api/v1/savings/summary
GET    /api/v1/savings/transactions

GET    /api/v1/loan-config
PUT    /api/v1/loan-config

GET    /api/v1/loan-applications
POST   /api/v1/loan-applications
POST   /api/v1/loan-applications/:id/approve
POST   /api/v1/loan-applications/:id/reject

GET    /api/v1/loans
GET    /api/v1/loans/:id
GET    /api/v1/loans/:id/export    ← file Excel
POST   /api/v1/installments/:id/pay

GET    /api/v1/shu-periods
POST   /api/v1/shu-periods
POST   /api/v1/shu-periods/:id/calculate

GET    /api/v1/dashboard
GET    /api/v1/modules
PUT    /api/v1/modules/:key

GET    /api/v1/notifications/logs
POST   /api/v1/notifications/trigger

GET    /api/v1/sync/pull
POST   /api/v1/sync/push

GET    /api/v1/inventory/field-defs          (RequireModule: inventory)
POST   /api/v1/inventory/field-defs
DELETE /api/v1/inventory/field-defs/:id
GET    /api/v1/inventory/products
POST   /api/v1/inventory/products
PUT    /api/v1/inventory/products/:id
POST   /api/v1/inventory/products/:id/movements
```

### Portal Anggota (Bearer token, role: anggota)
```
GET    /api/v1/portal/me
GET    /api/v1/portal/loans
GET    /api/v1/portal/loans/:id
GET    /api/v1/portal/shu
POST   /api/v1/portal/loan-applications
GET    /api/v1/portal/loan-applications
GET    /api/v1/portal/cooperatives
GET    /api/v1/portal/dashboard
```

### Integration Points (dapat diakses untuk demo)
```
POST   /api/v1/integrations/adins/ocr/ktp          (multipart: image)
POST   /api/v1/integrations/adins/credit-scoring   (json)
```

---

## 11. Cara Menjalankan

### Prasyarat

- Go 1.25+
- PostgreSQL (atau Supabase)
- `golang-migrate` CLI (opsional — bisa pakai `cmd/migrate`)

### Konfigurasi

Salin `.env.example` ke `.env` dan isi:

```env
DATABASE_DSN=postgres://user:pass@host:6543/db?sslmode=require
DATABASE_MIGRATION_DSN=postgres://user:pass@host:5432/db?sslmode=require
JWT_SECRET=your-secret-key
OCR_API_KEY=your-api-co-id-key
WA_MODE=mock                    # mock | sandbox
WA_TOKEN=                       # opsional, untuk mode sandbox
WA_PHONE_ID=                    # opsional, untuk mode sandbox
```

`config.json` digunakan Viper — lihat `.env.example` untuk semua field yang diperlukan.

### Migrasi Database

```bash
go run ./cmd/migrate/
```

### Seed Data Demo

```bash
psql $DATABASE_MIGRATION_DSN -f db/seed_dummy.sql
```

Seed mencakup: 1 koperasi Padiwangi, 1 akun pengurus, 2 akun anggota (budi & siti), 20+ member, riwayat simpanan 12 bulan, pinjaman mix status (aktif/lunas/menunggak), dan 1 SHU period final.

**Kredensial demo** (lihat juga [§2 Akun Demo](#2-akun-demo)):
- Pengurus: `pengurus@padiwangi.id` / `password123`
- Anggota (login pakai NIK): `3174096112900001` / `password123`

### Jalankan Server

```bash
go run ./cmd/web/
# Server berjalan di http://localhost:8080
```

---

## 12. Unit Test

Unit test tersedia untuk komponen dengan rumus/logika bisnis kritis:

```bash
go test ./internal/...
```

| File Test | Kasus Diuji |
|---|---|
| `usecase/schedule_test.go` | Angsuran flat: principal=5jt, rate=1.5%, tenor=12 → interest=75.000, total_due periode 1=491.666 |
| `usecase/shu_test.go` | Rumus jasa_modal (proporsional simpanan), pool 30%, totalSimpanan=0 |
| `usecase/savings_test.go` | Pokok duplikat→error; tarik wajib/pokok→error; sukarela melebihi saldo→error |
| `usecase/sync_test.go` | InvalidID rejected, UnknownType rejected, Duplicate, MissingField, Pull cursor |
| `gateway/adins/scoring_mock_test.go` | Score≈80→Grade A/approve; Score<60→Grade D/reject |
| `delivery/http/middleware/rbac_test.go` | anggota→pengurus endpoint→403; pengurus→portal→403; token invalid→401 |

Semua unit test yang menyangkut rumus keuangan **dicocokkan dengan contoh angka manual** sebelum dipakai di UI.

---

## Referensi Arsitektur

- [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture) — referensi struktur dan pola layering
- [GoFiber Documentation](https://docs.gofiber.io)
- [GORM Documentation](https://gorm.io/docs)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- UU No. 25/1992 tentang Perkoperasian — dasar rumus SHU
