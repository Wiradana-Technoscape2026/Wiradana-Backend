# backend/CLAUDE.md — Wiradana Backend

> Panduan kerja untuk AI agent & developer di direktori `backend/`. Baca ini **sebelum** menyentuh kode apa pun.
> Jika ada konflik antar dokumen: `docs/DECISIONS.md` > `docs/api_planning.md` > file ini.

---

## 0. Wajib Baca Tiap Sesi

Sebelum mulai, buka dan baca:
1. `docs/PROGRESS.md` — apa yang sudah/belum jadi. Jangan rebuild yang sudah ada.
2. `docs/DECISIONS.md` — keputusan yang menyimpang dari spec awal. Ikuti keputusan terbaru.

Dokumen referensi yang sering dibutuhkan:
- **`backend/docs/api_details.md`** — **endpoint yang sudah live**: path, method, params, request body, response shape, error codes. Update setiap endpoint baru selesai.
- `docs/api_planning.md` — **sumber kebenaran field, endpoint, enum**. Shape request/response, aturan bisnis, urutan build.
- `docs/ad_ins_mockup.md` — algoritma scoring deterministik, strategi OCR mock, interface gateway.
- `docs/member-portal.md` — dua role, matriks RBAC, enforcement JWT.
- `backend/docs/ERD.md` — skema DB, relasi antar tabel.
- `backend/docs/be_implementation.md` — implementasi teknis detail (clean arch, rumus, contoh angka).

---

## 1. Orientasi Cepat

**Stack:** Go (Fiber) · GORM · PostgreSQL · golang-migrate · Viper · go-playground/validator · golang-jwt/v5

**Pola arsitektur:** Clean Architecture ala [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture).

```
backend/
  cmd/
    main.go               # entry point: init config, DB, wire dependency, start Fiber
  internal/
    config/               # Viper — baca config.json / env
    domain/               # struct entitas domain (pure Go, no ORM tag)
    repository/           # implementasi DB (GORM)
    usecase/              # logika bisnis (tidak tahu HTTP/DB tech)
    delivery/
      http/               # Fiber handler + routing
        middleware/       # JWT auth, RequireRole, RequireModule
    gateway/
      adins/              # interface OCRGateway + ScoringGateway + implementasi mock
  db/
    migrations/           # file .sql golang-migrate (up & down)
    seeds/                # seed data realistis (≥20 anggota)
```

**Aturan lapisan:**
- `usecase` hanya bergantung pada `domain` + interface (`repository`, `gateway`). Tidak boleh import `fiber`, `gorm`, atau package `delivery`/`gateway` secara langsung.
- `repository` tidak tahu logika bisnis. `delivery/http` tidak tahu SQL.
- Swap implementasi (misalnya ganti mock Ad Ins → real) = **ganti binding di `main.go`, tidak sentuh usecase**.

---

## 2. Aturan Tidak Boleh Dilanggar

### 2.1 Uang — integer Rupiah
- Semua nominal uang = **`int64` (bigint PostgreSQL)**. Tidak ada `float64`, tidak ada `decimal`.
- Di dalam Go: `int64`. Di DB: `BIGINT`. Di JSON: integer (mis. `5000000`).
- Pembulatan pakai `math.Round` saat menghitung angsuran/SHU.

```go
// BENAR
type LoanApplication struct {
    Amount int64 `json:"amount"`
}
// SALAH — jangan
type LoanApplication struct {
    Amount float64 `json:"amount"`
}
```

### 2.2 Tenant — cooperative_id dari JWT
- Semua tabel bisnis punya kolom `cooperative_id UUID NOT NULL`.
- `cooperative_id` **selalu** diambil dari klaim JWT di middleware — **tidak pernah dari body request**.
- Setiap query DB wajib include `WHERE cooperative_id = ?`. Tidak boleh ada query tanpa filter tenant.

```go
// Ambil cooperative_id dari context (di-set middleware)
coopID := c.Locals("cooperative_id").(string)
members, err := r.memberRepo.FindAll(ctx, coopID, filter)
```

### 2.3 Insert-Only — transaksi keuangan
Tabel berikut tidak boleh punya endpoint UPDATE/DELETE:
- `savings_transaction`
- `payment`
- `inventory_movement` (jika Tier 3)

Tidak ada method `Update`/`Delete` di repository-nya. Koreksi = insert transaksi pembalik.

### 2.4 Password
- `password_hash` (bcrypt) **tidak pernah dikembalikan** di response JSON apa pun.
- Field `AppUser` yang dikembalikan: `id`, `email`, `role`, `member_id` — titik.

### 2.5 Envelope Response
Semua endpoint menggunakan format:
```json
// Sukses
{ "success": true, "data": { ... } }
// Sukses list
{ "success": true, "data": [ ... ], "meta": { "total": 42 } }
// Error
{ "success": false, "error": { "code": "VALIDATION_ERROR", "message": "NIK sudah terdaftar" } }
```

Buat helper terpusat — jangan tulis envelope inline di setiap handler:
```go
// delivery/http/response.go
func OK(c *fiber.Ctx, data any) error {
    return c.JSON(fiber.Map{"success": true, "data": data})
}
func Fail(c *fiber.Ctx, status int, code, message string) error {
    return c.Status(status).JSON(fiber.Map{
        "success": false,
        "error": fiber.Map{"code": code, "message": message},
    })
}
```

### 2.6 Enum — string lowercase
Lihat `api_planning.md §1` untuk daftar lengkap. Gunakan konstanta Go:
```go
const (
    RolePengurus = "pengurus"
    RoleAnggota  = "anggota"

    MemberStatusAktif    = "aktif"
    MemberStatusNonAktif = "nonaktif"
    MemberStatusKeluar   = "keluar"

    SavingsTypePokok    = "pokok"
    SavingsTypeWajib    = "wajib"
    SavingsSukarela     = "sukarela"
    DirectionSetor      = "setor"
    DirectionTarik      = "tarik"

    LoanStatusAktif      = "aktif"
    LoanStatusLunas      = "lunas"
    LoanStatusMenunggak  = "menunggak"

    AppStatusPending   = "pending"
    AppStatusApproved  = "approved"
    AppStatusRejected  = "rejected"

    GradeA = "A"; GradeB = "B"; GradeC = "C"; GradeD = "D"
    RecommApprove = "approve"; RecommReview = "review"; RecommReject = "reject"

    InstallmentBelumBayar = "belum_bayar"
    InstallmentLunas      = "lunas"
    InstallmentTerlambat  = "terlambat"

    ShuDraft = "draft"; ShuFinal = "final"
)
```

---

## 3. Auth & JWT

**Klaim JWT yang wajib ada:**
```go
type JWTClaims struct {
    UserID        string `json:"user_id"`
    CooperativeID string `json:"cooperative_id"`
    Role          string `json:"role"`        // "pengurus" | "anggota"
    MemberID      string `json:"member_id"`   // "" jika pengurus
    jwt.RegisteredClaims
}
```

**Middleware stack (urutan penting):**
```
Fiber router
  └── AuthMiddleware          // validasi token → set Locals (user_id, cooperative_id, role, member_id)
      └── RequireRole("pengurus")    // guard endpoint pengurus
      └── RequireRole("anggota")     // guard endpoint portal
          └── RequireModule("inventory")  // guard modul opsional
```

**`RequireRole` contoh:**
```go
func RequireRole(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Locals("role") != role {
            return Fail(c, 403, "FORBIDDEN", "Akses ditolak")
        }
        return c.Next()
    }
}
```

**Endpoint portal anggota — enforce member_id dari token:**
```go
// Handler portal: JANGAN ambil member_id dari body/param
memberID := c.Locals("member_id").(string)
// Jika anggota coba akses loan milik orang lain:
if loan.MemberID != memberID {
    return Fail(c, 403, "FORBIDDEN", "Bukan pinjaman Anda")
}
```

---

## 4. Database — Konvensi Migration & GORM

### Naming migration
```
{timestamp}_{deskripsi_singkat}.up.sql
{timestamp}_{deskripsi_singkat}.down.sql
// Contoh:
20260601120000_create_app_users.up.sql
20260601120001_create_members.up.sql
```
Jalankan: `migrate -path db/migrations -database $DATABASE_URL up`

### Kolom standar semua tabel bisnis
```sql
id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
cooperative_id  UUID NOT NULL REFERENCES cooperatives(id),
created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
```
Tabel transaksi: tidak ada `updated_at` (insert-only).
Tabel master (member, loan, dll.): tambah `updated_at TIMESTAMPTZ`.

### GORM model contoh
```go
type Member struct {
    ID               string         `gorm:"primaryKey;type:uuid" json:"id"`
    CooperativeID    string         `gorm:"type:uuid;not null"   json:"-"`
    NIK              string         `gorm:"uniqueIndex:idx_nik_coop;not null" json:"nik"`
    FullName         string         `gorm:"not null"             json:"full_name"`
    Address          string         `json:"address"`
    BirthDate        string         `gorm:"type:date"            json:"birth_date"`
    Status           string         `gorm:"default:aktif"        json:"status"`
    CustomAttributes datatypes.JSON `gorm:"type:jsonb"           json:"custom_attributes"`
    JoinedAt         time.Time      `json:"joined_at"`
    CreatedAt        time.Time      `json:"created_at"`
    UpdatedAt        time.Time      `json:"updated_at"`
}
```
NIK unique **per koperasi** → composite unique index `(nik, cooperative_id)`.

### Saldo = SUM (tidak ada kolom saldo)
```sql
-- Hitung savings_summary anggota
SELECT
    COALESCE(SUM(amount) FILTER (WHERE savings_type='pokok'   AND direction='setor'), 0) AS pokok,
    COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib'   AND direction='setor'), 0)
  - COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib'   AND direction='tarik'), 0) AS wajib,
    COALESCE(SUM(amount) FILTER (WHERE savings_type='sukarela' AND direction='setor'), 0)
  - COALESCE(SUM(amount) FILTER (WHERE savings_type='sukarela' AND direction='tarik'), 0) AS sukarela
FROM savings_transactions
WHERE member_id = $1 AND cooperative_id = $2;
```
**Jangan** simpan saldo di kolom terpisah yang di-UPDATE.

---

## 5. Modul Inti: Aturan Bisnis per Domain

### 5.1 Manajemen Anggota

| Aturan | Implementasi |
|---|---|
| NIK 16 digit, unik per koperasi | Validator tag + DB unique index + tangkap error 23505 → 409 |
| NIK tidak bisa diubah | `PUT /members/:id` tolak field `nik` jika dikirim (atau ignore) |
| Status: `aktif`/`nonaktif`/`keluar` | Enum validation di request |
| `custom_attributes` JSONB | Simpan apa adanya, kembalikan utuh |
| Pengajuan pinjaman hanya boleh anggota `aktif` | Cek di loan usecase sebelum buat aplikasi |

### 5.2 Manajemen Simpanan

| Aturan | Error jika dilanggar |
|---|---|
| Simpanan `pokok` hanya **sekali** setor per member | 409 CONFLICT — cek COUNT sebelum insert |
| `pokok` & `wajib` tidak bisa `tarik` | 409 CONFLICT |
| `sukarela` `tarik` tidak boleh melebihi saldo | 409 CONFLICT — hitung saldo via SUM dulu |
| `amount` harus > 0 | 400 VALIDATION_ERROR |
| INSERT-ONLY — tidak ada update/delete transaksi | Tidak ada endpoint/method update/delete |

```go
// usecase/savings_usecase.go — contoh guard
func (u *savingsUsecase) Record(ctx context.Context, req RecordSavingsReq) (*SavingsTransaction, error) {
    if req.SavingsType == SavingsTypePokok && req.Direction == DirectionSetor {
        count, _ := u.savingsRepo.CountPokok(ctx, req.MemberID, req.CooperativeID)
        if count > 0 {
            return nil, ErrPokokAlreadyRecorded
        }
    }
    if (req.SavingsType == SavingsTypePokok || req.SavingsType == SavingsTypeWajib) &&
        req.Direction == DirectionTarik {
        return nil, ErrCannotWithdrawMandatory
    }
    // ... dst
}
```

### 5.3 Loan Config

- Satu baris per koperasi (`cooperative_id` unique di tabel `loan_config`).
- `flat_rate_monthly` = persen per bulan (mis. `1.5` = 1.5%). Tipe `NUMERIC(5,2)` di DB, tapi kirim/terima sebagai `float64` di JSON — **hanya** untuk nilai rate, bukan untuk uang.
- Perubahan config hanya berlaku pinjaman **baru** — rate di-snapshot ke tabel `loan` saat approve.
- Seed default: `flat_rate_monthly: 1.5`, `max_plafond: 20000000`, `penalty_daily: 5000`.

### 5.4 Loan Origination → Management → Collection

**Flow `POST /loan-applications`:**
```
1. Validasi: member aktif, amount > 0, amount ≤ max_plafond
2. Hitung features dari data internal (§5.6)
3. Panggil ScoringGateway.Score(...) → simpan CreditAssessment
4. Buat LoanApplication (status=pending) + embed assessment
5. Return
```

**Flow `POST /loan-applications/:id/approve`:**
```
1. Cek status == "pending" → else 409
2. Snapshot flat_rate_monthly dari LoanConfig ke Loan
3. Buat Loan (status=aktif)
4. Generate schedule → INSERT semua Installment sekaligus (bulk insert)
5. approved_by dari JWT user_id
6. Return { application, loan }
```

**Rumus angsuran flat (integer Rupiah):**
```go
func GenerateSchedule(principal int64, rateMonthly float64, tenorMonths int, disbursedAt time.Time) []Installment {
    interestPerPeriod := int64(math.Round(float64(principal) * rateMonthly / 100))
    principalPerPeriod := principal / int64(tenorMonths)
    schedule := make([]Installment, tenorMonths)

    for i := 0; i < tenorMonths; i++ {
        p := principalPerPeriod
        if i == tenorMonths-1 {
            // Periode terakhir: sisa pokok agar total pas
            p = principal - principalPerPeriod*int64(tenorMonths-1)
        }
        schedule[i] = Installment{
            PeriodNo:     i + 1,
            DueDate:      disbursedAt.AddDate(0, i+1, 0).Format("2006-01-02"),
            PrincipalDue: p,
            InterestDue:  interestPerPeriod,
            TotalDue:     p + interestPerPeriod,
            Status:       InstallmentBelumBayar,
        }
    }
    return schedule
}
```

**Contoh wajib cocok (unit test):** principal=5.000.000, rate=1.5%/bln, tenor=12 bulan:
- `interest_per_period` = round(5.000.000 × 0.015) = **75.000**
- `principal_per_period` = 5.000.000 / 12 = **416.666** (periode 1–11), periode 12 = 416.676
- `total_due` periode 1–11 = **491.666**

**Flow `POST /installments/:id/pay`:**
```
1. INSERT Payment
2. Hitung paid_amount = SUM(payments WHERE schedule_id = ...)
3. Update installment.status:
   - paid_amount >= total_due → "lunas"
   - paid_amount > 0 && due_date < today → "terlambat" (partial late)
4. Recompute loan.status:
   - semua installment lunas → "lunas"
   - ada installment terlambat → "menunggak"
   - else → "aktif"
5. Return { payment, installment, loan }
```

### 5.5 SHU

**Flow `POST /shu-periods/:id/calculate`:**
```
1. Cek status == "draft" → else 409
2. Hitung total_simpanan semua anggota aktif (pokok + wajib)
3. Hitung total_jasa_pinjaman = SUM interest dari payments dalam periode
4. Per anggota:
   - jasa_modal = round(simpanan_anggota / total_simpanan × shu_modal_pool)
   - jasa_usaha = round(jasa_dibayar_anggota / total_jasa_koperasi × shu_usaha_pool)
   - total_shu = jasa_modal + jasa_usaha
5. Bulk INSERT ShuDistribution
6. Update ShuPeriod.status = "final"
7. Return { period, distributions[] }
```

**Contoh wajib cocok (unit test):** `be_implementation.md §5.6` — pastikan angka SHU per anggota match sebelum dipakai di FE.

---

## 6. Gateway Ad Ins — Mock Integration

Interface di `internal/gateway/adins/adins.go` — **jangan ubah interface ini**; produksi tinggal buat implementasi baru dan ganti binding di `main.go`.

```go
package adins

import "context"

type OCRGateway interface {
    ExtractKTP(ctx context.Context, image []byte) (KTPResult, error)
}

type ScoringGateway interface {
    Score(ctx context.Context, in ScoringInput) (ScoringResult, error)
}

type KTPResult struct {
    NIK        string  `json:"nik"`
    FullName   string  `json:"full_name"`
    Address    string  `json:"address"`
    BirthDate  string  `json:"birth_date"` // "YYYY-MM-DD"
    Confidence float64 `json:"confidence"`
    Source     string  `json:"source"`     // "MOCK_LITEDMS"
}

type ScoringInput struct {
    MemberID       string             `json:"member_id"`
    Features       map[string]float64 `json:"features"`
    JumlahDiajukan int64              `json:"jumlah_diajukan"`
    TenorBulan     int                `json:"tenor_bulan"`
}

type ScoringResult struct {
    Score            int      `json:"score"`
    Grade            string   `json:"grade"`
    Recommendation   string   `json:"recommendation"` // lowercase: approve|review|reject
    LimitRekomendasi int64    `json:"limit_rekomendasi"`
    Reasons          []string `json:"reasons"`
    Source           string   `json:"source"` // "MOCK_ADINS_SCORING"
}
```

### 6.1 Mock OCR (`ocr_mock.go`)

Dua strategi (detail: `ad_ins_mockup.md §1`):
- **Strategi A (rekomendasi jika waktu cukup):** Tesseract via `gosseract` → parse regex NIK/Nama/Alamat/Tgl → confidence dari kelengkapan field. Jika confidence < 0.5 → jatuh ke Strategi B.
- **Strategi B (fallback aman demo):** map hash/nama file gambar dummy ke data tetap.

```go
var mockKTPMap = map[string]KTPResult{
    "ktp_budi.jpg": {
        NIK: "3273010101900001", FullName: "Budi Santosa",
        Address: "Jl. Melati No. 10, Bandung", BirthDate: "1990-01-01",
        Confidence: 0.97, Source: "MOCK_LITEDMS",
    },
}
```
**JANGAN pakai foto KTP asli orang lain.** Buat file dummy sendiri.

### 6.2 Mock Credit Scoring (`scoring_mock.go`)

Algoritma deterministik — **input sama → output selalu sama** (tidak ada random).

```go
func (s *MockScoringGateway) Score(ctx context.Context, in ScoringInput) (ScoringResult, error) {
    f := in.Features
    score := math.Round(100 * (
        0.35*f["ketepatan_bayar"] +
        0.25*math.Min(f["rasio_simpanan_pinjaman"], 1) +
        0.15*math.Min(f["lama_keanggotaan_bulan"]/36, 1) +
        0.15*f["konsistensi_simpanan"] +
        0.10*(1-math.Min(f["rasio_beban_angsuran"], 1)),
    ))

    grade, rec := scoreToGrade(int(score))
    limit := limitFromGrade(grade, totalSimpanan, maxPlafond)
    reasons := generateReasons(f)

    return ScoringResult{
        Score: int(score), Grade: grade, Recommendation: rec,
        LimitRekomendasi: limit, Reasons: reasons,
        Source: "MOCK_ADINS_SCORING",
    }, nil
}

func scoreToGrade(score int) (string, string) {
    switch {
    case score >= 80: return "A", "approve"
    case score >= 70: return "B", "approve"
    case score >= 60: return "C", "review"
    default:          return "D", "reject"
    }
}
```

**Contoh unit test wajib cocok** (`ad_ins_mockup.md §2.4`):
- ketepatan=0.95, rasio_simp_pinj=0.6, lama=28, konsistensi=0.9, beban=0.3 → **score≈80, grade=A**

### 6.3 Hitung Features (Loan Usecase → Scoring)

Features dihitung dari data internal sebelum panggil `ScoringGateway`:

```go
func (u *loanUsecase) buildScoringFeatures(ctx context.Context, memberID, coopID string, amount int64, tenor int) map[string]float64 {
    // 1. ketepatan_bayar: proporsi angsuran tepat waktu
    //    = jumlah installment "lunas" tidak terlambat / total installment lunas
    //    Default 1.0 jika belum pernah pinjam

    // 2. rasio_simpanan_pinjaman: total_simpanan / amount (cap 1.0)

    // 3. lama_keanggotaan_bulan: bulan sejak joined_at (cap cap 36 untuk normalisasi)

    // 4. konsistensi_simpanan: proporsi bulan ada setoran wajib dalam 12 bln terakhir

    // 5. rasio_beban_angsuran: angsuran_baru / estimasi_kapasitas
    //    (jika tidak ada data pendapatan, pakai simpanan/12 sebagai proxy)

    return map[string]float64{
        "ketepatan_bayar":          ketepatan,
        "rasio_simpanan_pinjaman":  math.Min(float64(totalSimpanan)/float64(amount), 1.0),
        "lama_keanggotaan_bulan":   float64(lamaBulan),
        "konsistensi_simpanan":     konsistensi,
        "rasio_beban_angsuran":     math.Min(beban, 1.0),
    }
}
```

---

## 7. Endpoint Integration Points (Tetap Diekspos)

Dua endpoint berikut diekspos sebagai HTTP endpoint meskipun juga dipanggil internal — agar **"integration point" terlihat di demo**:

```
POST /api/v1/integrations/adins/ocr/ktp          (multipart: image)
POST /api/v1/integrations/adins/credit-scoring   (json)
```

Kedua endpoint ini ada di `delivery/http/integration_controller.go`. Mereka langsung memanggil implementasi mock gateway.

---

## 8. Kode Error & HTTP Status

| `code` | HTTP | Kapan |
|---|---|---|
| `VALIDATION_ERROR` | 400 | Input tidak valid (field kosong, format salah, amount ≤ 0) |
| `UNAUTHORIZED` | 401 | Token tidak ada, expired, atau invalid |
| `FORBIDDEN` | 403 | Role tidak sesuai, atau akses data milik orang lain |
| `NOT_FOUND` | 404 | Resource tidak ada di koperasi ini |
| `CONFLICT` | 409 | NIK duplikat, pokok sudah disetor, status salah (mis. approve loan non-pending) |
| `INTERNAL` | 500 | Error tidak terduga |

Tangkap error PostgreSQL `23505` (unique violation) → konversi ke `CONFLICT`:
```go
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) && pgErr.Code == "23505" {
    return Fail(c, 409, "CONFLICT", "NIK sudah terdaftar")
}
```

---

## 9. Routing (Ringkasan dari api_planning.md)

```go
// main routes
api := app.Group("/api/v1")

// Public
api.Post("/auth/login", authHandler.Login)

// Pengurus only
pengurus := api.Use(AuthMiddleware, RequireRole("pengurus"))
pengurus.Get("/members", memberHandler.List)
pengurus.Post("/members", memberHandler.Create)
pengurus.Get("/members/:id", memberHandler.GetByID)
pengurus.Put("/members/:id", memberHandler.Update)
pengurus.Get("/members/:id/savings", savingsHandler.List)
pengurus.Post("/members/:id/savings", savingsHandler.Record)
pengurus.Get("/loan-config", loanConfigHandler.Get)
pengurus.Put("/loan-config", loanConfigHandler.Update)
pengurus.Get("/loan-applications", loanAppHandler.List)
pengurus.Post("/loan-applications", loanAppHandler.Create)
pengurus.Post("/loan-applications/:id/approve", loanAppHandler.Approve)
pengurus.Post("/loan-applications/:id/reject", loanAppHandler.Reject)
pengurus.Get("/loans", loanHandler.List)
pengurus.Get("/loans/:id", loanHandler.GetByID)
pengurus.Post("/installments/:id/pay", installmentHandler.Pay)
pengurus.Get("/shu-periods", shuHandler.ListPeriods)
pengurus.Post("/shu-periods", shuHandler.CreatePeriod)
pengurus.Post("/shu-periods/:id/calculate", shuHandler.Calculate)
pengurus.Get("/dashboard", dashboardHandler.Get)
pengurus.Get("/modules", moduleHandler.List)
pengurus.Put("/modules/:key", moduleHandler.Update)

// Portal anggota
portal := api.Group("/portal").Use(AuthMiddleware, RequireRole("anggota"))
portal.Get("/me", portalHandler.Me)
portal.Get("/loans", portalHandler.Loans)
portal.Get("/loans/:id", portalHandler.LoanDetail)
portal.Get("/shu", portalHandler.SHU)
portal.Post("/loan-applications", portalHandler.ApplyLoan)
portal.Get("/loan-applications", portalHandler.LoanApplications)

// Integration points (diekspos, tidak perlu role guard — atau beri guard sesuai kebutuhan demo)
api.Post("/integrations/adins/ocr/ktp", integrationHandler.OCRKTP)
api.Post("/integrations/adins/credit-scoring", integrationHandler.CreditScoring)

// Tier 3 — Inventory (dengan module guard)
inv := pengurus.Group("/inventory").Use(RequireModule("inventory"))
inv.Get("/field-defs", inventoryHandler.ListFieldDefs)
inv.Post("/field-defs", inventoryHandler.CreateFieldDef)
inv.Delete("/field-defs/:id", inventoryHandler.DeleteFieldDef)
inv.Get("/products", inventoryHandler.ListProducts)
inv.Post("/products", inventoryHandler.CreateProduct)
inv.Put("/products/:id", inventoryHandler.UpdateProduct)
inv.Post("/products/:id/movements", inventoryHandler.RecordMovement)
```

---

## 10. Urutan Build (dari api_planning.md §5)

Ikuti urutan ini — jangan skip ke depan jika blok sebelumnya belum selesai:

| Blok (PLANNING.md) | Target Jam | Yang Dibangun |
|---|---|---|
| **A** | J0–4 | Setup repo, Fiber, GORM, PostgreSQL, migrasi tabel inti, JWT login pengurus |
| **B** | J4–10 | Members CRUD + NIK validation, OCR mock endpoint, Savings (setor/tarik + aturan bisnis) |
| **C** | J10–16 | LoanConfig, LoanApplication (+scoring flow), Approve (+generate schedule), Loans, Installments/Pay, Dashboard |
| **D** | J16–24 | SHU periods (+calculate), Portal endpoints (`/portal/*`), notifikasi mock log |
| **E** | J24–28 | Modules toggle, Inventory (Tier 3 jika core solid), polish |

**Checkpoint J16 = titik aman.** Full flow pinjaman harus jalan end-to-end sebelum melanjutkan ke SHU/portal.

---

## 11. Seed Data

Seed minimal sebelum demo (`db/seeds/`):
- 1 koperasi: `Padiwangi`
- 1 akun pengurus: `pengurus@padiwangi.id` / `password123`
- 2 akun anggota: `budi@padiwangi.id`, `siti@padiwangi.id` (untuk demo portal)
- ≥20 member dengan simpanan (mix pokok/wajib/sukarela)
- Pinjaman dengan status mix: beberapa `lunas`, beberapa `aktif`, 1–2 `menunggak`
- Riwayat simpanan 6–12 bulan terakhir agar scoring features punya data
- 1 ShuPeriod `final` dengan distribusi (agar dashboard SHU bisa ditampilkan)

---

## 12. Unit Test Wajib (Sebelum Diintegrasikan ke UI)

Buat unit test di `internal/usecase/` untuk memastikan angka cocok dengan contoh manual:

```go
// TestGenerateSchedule: principal=5.000.000, rate=1.5%, tenor=12
// → interest_per_period=75.000, total_due periode 1=491.666

// TestSHUCalculation: lihat contoh angka di be_implementation.md §5.6
// → jasa_modal + jasa_usaha per anggota harus match

// TestCreditScoring: input dari ad_ins_mockup.md §2.4
// → score≈80, grade="A", recommendation="approve", limit=9.000.000
```

AI agent: **jalankan unit test ini sebelum laporan "feature selesai"** — satu angka salah di rumus bisa merusak kepercayaan demo.

---

## 13. Aturan Kejujuran untuk Naskah Demo

Ini berlaku saat code comment, README, atau dokumentasi yang mungkin dibaca juri:

- Mock Ad Ins = **"integration-ready mock"**, bukan *"integrasi API Ad Ins"*
- Scoring = **"rule-based berbasis data internal"**, bukan *"AI-powered"*
- Multi-tenant = **"tenant-aware via cooperative_id"**, bukan *"RLS aktif"*
- Transaksi = **"insert-only"**, bukan *"audit trail eksplisit"*
- Bunga = **"flat"** saja — jangan sebut "menurun" jika tidak diimplementasikan

---

## 14. Referensi Cepat

| Kebutuhan | Lihat di |
|---|---|
| Shape semua request/response | `docs/api_planning.md §2–3` |
| Rumus angsuran flat + contoh angka | `docs/WIRADANA.md §5.1` + `backend/docs/be_implementation.md §5.4` |
| Rumus SHU + contoh angka | `docs/WIRADANA.md §5.2` + `backend/docs/be_implementation.md §5.6` |
| Algoritma scoring + contoh | `docs/ad_ins_mockup.md §2.3–2.4` |
| Skema DB lengkap | `backend/docs/ERD.md` |
| Matriks RBAC | `docs/member-portal.md §2` |
| Flow portal anggota | `docs/member-portal.md §3–4` |
| Modul Inventory detail | `backend/docs/modules/inventory.md` |
| Status sudah/belum | `docs/PROGRESS.md` |
| Keputusan menyimpang | `docs/DECISIONS.md` |
