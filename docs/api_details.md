# api_details.md — Backend Endpoint Reference

Dokumentasi semua endpoint yang terdaftar di `internal/delivery/http/route/route.go`. Update file ini setiap ada endpoint baru.

> **Envelope response** semua endpoint:
> - Sukses: `{ "success": true, "data": { ... } }`
> - Sukses list: `{ "success": true, "data": [...], "meta": { "total": N } }`
> - Error: `{ "success": false, "error": { "code": "...", "message": "..." } }`

**Status legend:** ✅ Live · 🔲 Belum diimplementasi

---

## Ringkasan Semua Endpoint

| Method | Path | Auth | Status |
|---|---|---|---|
| GET | `/health` | — | ✅ |
| POST | `/api/v1/auth/login` | — | ✅ |
| POST | `/api/v1/auth/register/pengurus` | — | ✅ |
| POST | `/api/v1/members` | pengurus | ✅ |
| GET | `/api/v1/members` | pengurus | ✅ |
| GET | `/api/v1/members/:id` | pengurus | ✅ |
| PUT | `/api/v1/members/:id` | pengurus | ✅ |
| GET | `/api/v1/members/:id/savings` | pengurus | ✅ |
| POST | `/api/v1/members/:id/savings` | pengurus | ✅ |
| GET | `/api/v1/savings/summary` | pengurus | ✅ |
| GET | `/api/v1/savings/transactions` | pengurus | ✅ |
| GET | `/api/v1/loan-config` | pengurus | ✅ |
| PUT | `/api/v1/loan-config` | pengurus | ✅ |
| GET | `/api/v1/loan-applications` | pengurus | ✅ |
| POST | `/api/v1/loan-applications` | pengurus | ✅ |
| POST | `/api/v1/loan-applications/:id/approve` | pengurus | ✅ |
| POST | `/api/v1/loan-applications/:id/reject` | pengurus | ✅ |
| GET | `/api/v1/loans` | pengurus | ✅ |
| GET | `/api/v1/loans/:id` | pengurus | ✅ |
| POST | `/api/v1/installments/:id/pay` | pengurus | ✅ |
| POST | `/api/v1/integrations/adins/ocr/ktp` | pengurus | ✅ |
| POST | `/api/v1/integrations/adins/credit-scoring` | pengurus | ✅ |
| GET | `/api/v1/dashboard` | pengurus | ✅ |
| GET | `/api/v1/portal/loan-applications` | anggota | ✅ |
| POST | `/api/v1/portal/loan-applications` | anggota | ✅ |
| GET | `/api/v1/portal/loans` | anggota | ✅ |
| GET | `/api/v1/portal/loans/:id` | anggota | ✅ |
| GET/POST | `/api/v1/inventory/*` | pengurus + modul aktif | 🔲 Tier 3, belum diimplementasi |

---

## Health Check

### GET /health
Server health check. Tidak butuh auth.

**Response 200:**
```json
{ "success": true, "data": "ok" }
```

---

## Auth

### POST /api/v1/auth/register/pengurus
Daftarkan akun pengurus baru untuk sebuah koperasi. Endpoint ini public dan mengembalikan JWT agar akun bisa langsung dipakai.

**Auth:** tidak diperlukan (public endpoint)

**Request Body** (`application/json`):
```json
{
  "cooperative_id": "11111111-1111-1111-1111-111111111111",
  "email": "pengurus.baru@wiradana.local",
  "password": "password123"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `cooperative_id` | string | ya | UUID koperasi tujuan |
| `email` | string | ya | Email akun, harus unik |
| `password` | string | ya | Password plaintext, minimal 8 karakter |

**Response 200:**
```json
{
  "success": true,
  "data": {
    "token": "<JWT>",
    "user": {
      "id": "uuid",
      "email": "pengurus.baru@wiradana.local",
      "role": "pengurus",
      "member_id": null
    }
  }
}
```

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Body tidak valid, UUID salah, atau password terlalu pendek |
| 409 | `CONFLICT` | Email sudah terdaftar |

### POST /api/v1/auth/login
Login user (pengurus atau anggota). Mengembalikan JWT token yang harus disertakan di semua request selanjutnya sebagai `Authorization: Bearer <token>`.

**Auth:** tidak diperlukan (public endpoint)

**Request Body** (`application/json`):
```json
{
  "email": "pengurus@padiwangi.id",
  "password": "password123"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `email` | string | ya | Email akun |
| `password` | string | ya | Password plaintext |

**Response 200:**
```json
{
  "success": true,
  "data": {
    "token": "<JWT>",
    "user": {
      "id": "uuid",
      "email": "pengurus@padiwangi.id",
      "role": "pengurus",
      "member_id": null
    }
  }
}
```

> `member_id` berisi UUID jika role `anggota`, bernilai `null` jika role `pengurus`.

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Field kosong atau format salah |
| 401 | `UNAUTHORIZED` | Email/password salah |

---

## Members

> Semua endpoint member memerlukan role **pengurus**.
> `cooperative_id` diambil otomatis dari JWT — **tidak perlu dan tidak boleh dikirim di body**.

### POST /api/v1/members
Daftarkan anggota baru ke koperasi.

**Auth:** Bearer token (role: `pengurus`)

**Request Body** (`application/json`):
```json
{
  "nik": "3175010101900001",
  "full_name": "Budi Santosa",
  "address": "Jl. Melati No. 10, Bandung",
  "birth_date": "1990-01-01",
  "custom_attributes": {
    "pekerjaan": "Petani",
    "agama": "Islam"
  }
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `nik` | string | **ya** | Tepat 16 digit angka, unik per koperasi |
| `full_name` | string | **ya** | Nama lengkap |
| `address` | string | tidak | Alamat tinggal |
| `birth_date` | string | tidak | Format `"YYYY-MM-DD"` |
| `custom_attributes` | object | tidak | Data tambahan bebas (JSONB), default `{}` |

**Response 200:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "nik": "3175010101900001",
    "full_name": "Budi Santosa",
    "address": "Jl. Melati No. 10, Bandung",
    "birth_date": "1990-01-01",
    "status": "aktif",
    "custom_attributes": { "pekerjaan": "Petani" },
    "joined_at": "2026-06-12T10:00:00Z",
    "savings_summary": {
      "pokok": 0,
      "wajib": 0,
      "sukarela": 0,
      "total": 0
    }
  }
}
```

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | NIK bukan 16 digit / `full_name` kosong |
| 409 | `CONFLICT` | NIK sudah terdaftar di koperasi ini |

---

### GET /api/v1/members
Daftar semua anggota koperasi. Mendukung pencarian dan filter status. Hasil diurutkan `joined_at DESC`.

**Auth:** Bearer token (role: `pengurus`)

**Query Params:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `search` | string | tidak | Cari di `full_name` atau `nik` (case-insensitive, partial match) |
| `status` | string | tidak | Filter: `aktif` \| `nonaktif` \| `keluar` |

Contoh: `GET /api/v1/members?search=budi&status=aktif`

**Response 200:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "nik": "3175010101900001",
      "full_name": "Budi Santosa",
      "address": "Jl. Melati No. 10, Bandung",
      "birth_date": "1990-01-01",
      "status": "aktif",
      "custom_attributes": {},
      "joined_at": "2026-06-12T10:00:00Z",
      "savings_summary": {
        "pokok": 500000,
        "wajib": 300000,
        "sukarela": 100000,
        "total": 900000
      }
    }
  ],
  "meta": { "total": 1 }
}
```

---

### GET /api/v1/members/:id
Detail satu anggota beserta ringkasan simpanan real-time.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID anggota |

**Response 200:** sama dengan satu objek di `GET /api/v1/members`

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 404 | `NOT_FOUND` | Anggota tidak ditemukan di koperasi ini |

---

### PUT /api/v1/members/:id
Update data anggota. **NIK tidak bisa diubah.** Semua field bersifat opsional — hanya field yang dikirim yang akan diupdate.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID anggota |

**Request Body** (`application/json`):
```json
{
  "full_name": "Budi Santosa Updated",
  "address": "Jl. Baru No. 5, Bandung",
  "status": "nonaktif",
  "custom_attributes": { "pekerjaan": "Wiraswasta" }
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `full_name` | string | tidak | Nama baru |
| `address` | string | tidak | Alamat baru |
| `status` | string | tidak | `aktif` \| `nonaktif` \| `keluar` |
| `custom_attributes` | object | tidak | **Menimpa seluruh** custom attributes yang ada |

**Response 200:** MemberResponse dengan data terbaru

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | `status` bukan nilai enum yang valid |
| 404 | `NOT_FOUND` | Anggota tidak ditemukan |

---

## Savings

> Semua endpoint simpanan memerlukan role **pengurus**.
> `cooperative_id` diambil otomatis dari JWT — tidak perlu dikirim di body.
> Semua transaksi **INSERT-ONLY** — tidak ada PUT/DELETE.

### GET /api/v1/members/:id/savings
Daftar seluruh transaksi simpanan seorang anggota. Diurutkan `created_at DESC`.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID anggota |

**Response 200:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "member_id": "uuid",
      "savings_type": "wajib",
      "direction": "setor",
      "amount": 150000,
      "created_at": "2026-06-09T10:00:00Z"
    }
  ],
  "meta": { "total": 12 }
}
```

| Field | Tipe | Nilai |
|---|---|---|
| `savings_type` | string | `pokok` \| `wajib` \| `sukarela` |
| `direction` | string | `setor` \| `tarik` |
| `amount` | int64 | Nominal dalam Rupiah |

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 404 | `NOT_FOUND` | Anggota tidak ditemukan di koperasi ini |

---

### POST /api/v1/members/:id/savings
Catat transaksi simpanan (setor atau tarik) untuk seorang anggota.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID anggota |

**Request Body** (`application/json`):
```json
{
  "savings_type": "wajib",
  "direction": "setor",
  "amount": 150000
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `savings_type` | string | **ya** | `pokok` \| `wajib` \| `sukarela` |
| `direction` | string | **ya** | `setor` \| `tarik` |
| `amount` | int64 | **ya** | Nominal Rupiah, harus > 0 |

**Response 200:** SavingsTransactionResponse (satu objek seperti di list)

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Field tidak valid atau `amount` ≤ 0 |
| 404 | `NOT_FOUND` | Anggota tidak ditemukan |
| 409 | `CONFLICT` | Simpanan `pokok` sudah pernah disetor |
| 409 | `CONFLICT` | `pokok` atau `wajib` tidak bisa ditarik |
| 409 | `CONFLICT` | Saldo `sukarela` tidak mencukupi untuk penarikan |

---

### GET /api/v1/savings/summary
Total simpanan seluruh koperasi, dikelompokkan per jenis. Dipakai untuk kartu ringkasan di halaman Simpanan.

**Auth:** Bearer token (role: `pengurus`)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "pokok": 15000000,
    "wajib": 8400000,
    "sukarela": 3200000,
    "total": 26600000
  }
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `pokok` | int64 | Total setor pokok seluruh anggota |
| `wajib` | int64 | Total setor wajib − tarik wajib seluruh anggota |
| `sukarela` | int64 | Total setor sukarela − tarik sukarela seluruh anggota |
| `total` | int64 | `pokok + wajib + sukarela` |

---

### GET /api/v1/savings/transactions
Daftar seluruh transaksi simpanan koperasi lintas anggota, dengan nama anggota di-join. Diurutkan `created_at DESC`. Mendukung filter dan paginasi.

**Auth:** Bearer token (role: `pengurus`)

**Query Params:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `type` | string | tidak | Filter jenis: `pokok` \| `wajib` \| `sukarela` |
| `direction` | string | tidak | Filter arah: `setor` \| `tarik` |
| `limit` | int | tidak | Jumlah baris per halaman (default: `20`) |
| `page` | int | tidak | Halaman (default: `1`) |

Contoh: `GET /api/v1/savings/transactions?type=wajib&direction=setor&limit=10&page=1`

**Response 200:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "member_id": "uuid",
      "member_name": "Budi Santosa",
      "savings_type": "wajib",
      "direction": "setor",
      "amount": 150000,
      "created_at": "2026-06-09T10:00:00Z"
    }
  ],
  "meta": { "total": 87 }
}
```

---

## Loan Config

> Konfigurasi bunga dan plafond koperasi. Role **pengurus** diperlukan.
> Jika belum ada config, GET akan membuat default otomatis (flat_rate=1.5%, max_plafond=20.000.000, penalty=5.000/hari).

### GET /api/v1/loan-config
Ambil konfigurasi pinjaman koperasi saat ini.

**Auth:** Bearer token (role: `pengurus`)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "flat_rate_monthly": 1.5,
    "max_plafond": 20000000,
    "penalty_daily": 5000
  }
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `id` | string | UUID config |
| `flat_rate_monthly` | float64 | Bunga flat per bulan (%), mis. `1.5` = 1.5% |
| `max_plafond` | int64 | Batas maksimum pinjaman dalam Rupiah |
| `penalty_daily` | int64 | Denda keterlambatan per hari dalam Rupiah |

---

### PUT /api/v1/loan-config
Update konfigurasi pinjaman. Semua field opsional — hanya field yang dikirim yang berubah.

**Auth:** Bearer token (role: `pengurus`)

**Request Body** (`application/json`):
```json
{
  "flat_rate_monthly": 2.0,
  "max_plafond": 25000000,
  "penalty_daily": 10000
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `flat_rate_monthly` | float64 | tidak | Bunga flat baru (harus > 0) |
| `max_plafond` | int64 | tidak | Plafond baru (harus > 0) |
| `penalty_daily` | int64 | tidak | Denda harian baru (boleh 0) |

**Response 200:** LoanConfigResponse (sama dengan GET)

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | `flat_rate_monthly` atau `max_plafond` bernilai ≤ 0 |

---

## Loan Applications

> Pengajuan pinjaman baru. Role **pengurus** diperlukan untuk semua endpoint di section ini.
> Credit scoring dijalankan otomatis saat pengajuan dibuat.

### Tipe Data: LoanApplicationResponse

```json
{
  "id": "uuid",
  "member_id": "uuid",
  "member_name": "Budi Santosa",
  "amount": 5000000,
  "tenor_months": 12,
  "purpose": "Modal tani",
  "status": "pending",
  "approved_by": null,
  "created_at": "2026-06-12T10:00:00Z",
  "assessment": {
    "id": "uuid",
    "application_id": "uuid",
    "score": 80,
    "grade": "A",
    "recommendation": "approve",
    "limit_suggested": 9000000,
    "features": {
      "ketepatan_bayar": 0.95,
      "rasio_simpanan_pinjaman": 0.6,
      "lama_keanggotaan_bulan": 28,
      "konsistensi_simpanan": 0.9,
      "rasio_beban_angsuran": 0.3
    },
    "reasons": ["riwayat pembayaran sangat baik", "simpanan memadai terhadap pinjaman"],
    "source": "MOCK_ADINS_SCORING"
  }
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `status` | string | `pending` \| `approved` \| `rejected` |
| `approved_by` | string\|null | UUID user pengurus yang approve, null jika belum |
| `assessment` | object\|null | Hasil credit scoring, null jika scoring gagal |
| `assessment.grade` | string | `A` (≥80) \| `B` (≥70) \| `C` (≥60) \| `D` (<60) |
| `assessment.recommendation` | string | `approve` \| `review` \| `reject` |
| `assessment.source` | string | Selalu `"MOCK_ADINS_SCORING"` (integration-ready mock) |

---

### GET /api/v1/loan-applications
Daftar semua pengajuan pinjaman di koperasi. Diurutkan `created_at DESC`.

**Auth:** Bearer token (role: `pengurus`)

**Query Params:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `status` | string | tidak | Filter: `pending` \| `approved` \| `rejected` |

Contoh: `GET /api/v1/loan-applications?status=pending`

**Response 200:**
```json
{
  "success": true,
  "data": [ /* array LoanApplicationResponse */ ],
  "meta": { "total": 3 }
}
```

---

### POST /api/v1/loan-applications
Buat pengajuan pinjaman baru untuk anggota. Credit scoring dijalankan otomatis dan hasilnya disimpan di `assessment`.

**Auth:** Bearer token (role: `pengurus`)

**Request Body** (`application/json`):
```json
{
  "member_id": "uuid-anggota",
  "amount": 5000000,
  "tenor_months": 12,
  "purpose": "Modal tani"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `member_id` | string | **ya** | UUID anggota, harus milik koperasi ini |
| `amount` | int64 | **ya** | Jumlah pinjaman dalam Rupiah (> 0) |
| `tenor_months` | int | **ya** | Tenor dalam bulan (> 0) |
| `purpose` | string | tidak | Tujuan penggunaan dana |

**Response 200:** LoanApplicationResponse (dengan `assessment` terisi)

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Field wajib kosong atau format salah |
| 400 | `VALIDATION_ERROR` | `amount` melebihi `max_plafond` di loan-config |
| 404 | `NOT_FOUND` | `member_id` tidak ditemukan di koperasi ini |
| 409 | `CONFLICT` | Status anggota bukan `aktif` |

---

### POST /api/v1/loan-applications/:id/approve
Setujui pengajuan pinjaman. Otomatis membuat `Loan` aktif dan jadwal angsuran flat-rate.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID pengajuan (`loan_application`) |

**Request Body:** tidak diperlukan

**Response 200:**
```json
{
  "success": true,
  "data": {
    "application": { /* LoanApplicationResponse dengan status "approved" */ },
    "loan": { /* LoanResponse dengan schedule penuh */ }
  }
}
```

> Jadwal angsuran (`loan.schedule`) berisi seluruh `tenor_months` periode, dihitung dengan rumus flat-rate dari `loan-config` saat ini.

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 404 | `NOT_FOUND` | Pengajuan tidak ditemukan |
| 409 | `CONFLICT` | Pengajuan bukan dalam status `pending` |

---

### POST /api/v1/loan-applications/:id/reject
Tolak pengajuan pinjaman.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID pengajuan |

**Request Body** (`application/json`, opsional):
```json
{
  "reason": "Simpanan tidak memenuhi syarat"
}
```

**Response 200:** LoanApplicationResponse dengan `status: "rejected"`

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 404 | `NOT_FOUND` | Pengajuan tidak ditemukan |
| 409 | `CONFLICT` | Pengajuan bukan dalam status `pending` |

---

## Loans

> Pinjaman yang sudah disetujui dan aktif. Role **pengurus** diperlukan.

### Tipe Data: LoanResponse

```json
{
  "id": "uuid",
  "application_id": "uuid",
  "member_id": "uuid",
  "member_name": "Budi Santosa",
  "principal": 5000000,
  "flat_rate_monthly": 1.5,
  "tenor_months": 12,
  "status": "aktif",
  "disbursed_at": "2026-06-12",
  "outstanding": 5000000,
  "schedule": [
    {
      "id": "uuid",
      "loan_id": "uuid",
      "period_no": 1,
      "due_date": "2026-07-12",
      "principal_due": 416666,
      "interest_due": 75000,
      "total_due": 491666,
      "paid_amount": 0,
      "status": "belum_bayar"
    }
  ]
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `status` | string | `aktif` \| `lunas` \| `menunggak` |
| `outstanding` | int64 | Sisa pokok yang belum lunas (dihitung real-time dari jadwal) |
| `disbursed_at` | string | Format `"YYYY-MM-DD"` |
| `schedule` | array\|null | Ada di detail (`GET /loans/:id`), tidak ada di list |
| `schedule[].status` | string | `belum_bayar` \| `lunas` \| `terlambat` |
| `schedule[].paid_amount` | int64 | Total pembayaran yang sudah masuk untuk angsuran ini |

---

### GET /api/v1/loans
Daftar semua pinjaman di koperasi. Diurutkan `disbursed_at DESC`. `schedule` tidak disertakan di list.

**Auth:** Bearer token (role: `pengurus`)

**Query Params:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `status` | string | tidak | Filter: `aktif` \| `lunas` \| `menunggak` |

**Response 200:**
```json
{
  "success": true,
  "data": [ /* array LoanResponse tanpa schedule */ ],
  "meta": { "total": 5 }
}
```

---

### GET /api/v1/loans/:id
Detail satu pinjaman beserta jadwal angsuran lengkap.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID pinjaman |

**Response 200:** LoanResponse dengan `schedule` terisi

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 404 | `NOT_FOUND` | Pinjaman tidak ditemukan di koperasi ini |

---

## Installments

### POST /api/v1/installments/:id/pay
Catat pembayaran angsuran. INSERT-ONLY — tidak ada update baris lama. Status angsuran dan pinjaman diperbarui otomatis.

**Auth:** Bearer token (role: `pengurus`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID angsuran (`installment_schedule`) |

**Request Body** (`application/json`):
```json
{
  "amount": 491666,
  "penalty": 0
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `amount` | int64 | **ya** | Jumlah yang dibayar (> 0) |
| `penalty` | int64 | tidak | Denda yang dibayar (≥ 0, default 0) |

**Response 200:**
```json
{
  "success": true,
  "data": {
    "payment": {
      "id": "uuid",
      "schedule_id": "uuid",
      "amount": 491666,
      "penalty": 0,
      "paid_at": "2026-06-12T10:30:00Z"
    },
    "installment": {
      "id": "uuid",
      "loan_id": "uuid",
      "period_no": 1,
      "due_date": "2026-07-12",
      "principal_due": 416666,
      "interest_due": 75000,
      "total_due": 491666,
      "paid_amount": 491666,
      "status": "lunas"
    },
    "loan": { /* LoanResponse terbaru dengan outstanding terupdate */ }
  }
}
```

> **Status logic:**
> - `installment.status → "lunas"` jika `paid_amount >= total_due`
> - `loan.status → "lunas"` jika semua angsuran berstatus `lunas`
> - `loan.status → "menunggak"` jika ada angsuran berstatus `terlambat`

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | `amount` kosong atau ≤ 0 |
| 404 | `NOT_FOUND` | Angsuran tidak ditemukan |

---

## Integrations

> Endpoint integrasi layanan eksternal. Role **pengurus** diperlukan.

### POST /api/v1/integrations/adins/ocr/ktp
Upload foto KTP → ekstrak data via api.co.id → kembalikan field yang sudah di-parse untuk auto-fill form registrasi anggota.

**Auth:** Bearer token (role: `pengurus`)

**Request:** `multipart/form-data`

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `file` | file | **ya** | Foto KTP (JPG/JPEG/PNG, maks. 5 MB) |

Contoh cURL:
```bash
curl -X POST http://localhost:8080/api/v1/integrations/adins/ocr/ktp \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/ktp.jpg"
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "nik": "3175010101010001",
    "full_name": "JOHN DOE",
    "address": "JL. CONTOH NO. 123, RT/RW 001/002, Kel. KELURAHAN, Kec. KECAMATAN, KOTA JAKARTA, DKI JAKARTA",
    "birth_date": "1990-01-01",
    "confidence": 0.97,
    "source": "API_CO_ID"
  }
}
```

| Field Response | Tipe | Keterangan |
|---|---|---|
| `nik` | string | Nomor Induk Kependudukan (16 digit) |
| `full_name` | string | Nama lengkap dari KTP |
| `address` | string | Gabungan: alamat + RT/RW + kelurahan + kecamatan + kota + provinsi |
| `birth_date` | string | Format `"YYYY-MM-DD"` (dikonversi dari `DD-MM-YYYY` di KTP) |
| `confidence` | float | 0.70–0.97, dihitung dari kelengkapan 4 field utama (NIK, nama, alamat, tanggal lahir) |
| `source` | string | Selalu `"API_CO_ID"` — integrasi real dengan api.co.id |

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Field `file` tidak disertakan |
| 400 | `VALIDATION_ERROR` | Tipe file bukan JPG/PNG (cek magic bytes + ekstensi) |
| 400 | `VALIDATION_ERROR` | Ukuran file melebihi 5 MB |
| 422 | `VALIDATION_ERROR` | KTP tidak terbaca (foto buram, gelap, atau miring) |
| 500 | `INTERNAL` | OCR credits habis atau layanan api.co.id down |

---

### POST /api/v1/integrations/adins/credit-scoring
Jalankan credit scoring Ad Ins mock secara langsung dengan input manual. Berguna untuk demo, testing, atau preview skor sebelum membuat pengajuan.

**Auth:** Bearer token (role: `pengurus`)

**Request Body** (`application/json`):
```json
{
  "member_id": "uuid-anggota",
  "features": {
    "ketepatan_bayar": 0.95,
    "rasio_simpanan_pinjaman": 0.6,
    "lama_keanggotaan_bulan": 28,
    "konsistensi_simpanan": 0.9,
    "rasio_beban_angsuran": 0.3
  },
  "jumlah_diajukan": 5000000,
  "tenor_bulan": 12,
  "total_simpanan": 3000000,
  "max_plafond": 20000000
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `member_id` | string | UUID anggota (untuk referensi, tidak divalidasi ke DB) |
| `features` | object | 5 fitur scoring (nilai 0.0–1.0 kecuali `lama_keanggotaan_bulan`) |
| `features.ketepatan_bayar` | float64 | Rasio angsuran yang dibayar tepat waktu (0–1) |
| `features.rasio_simpanan_pinjaman` | float64 | Total simpanan / jumlah pinjaman (0–∞, dikap di 1) |
| `features.lama_keanggotaan_bulan` | float64 | Lama menjadi anggota dalam bulan |
| `features.konsistensi_simpanan` | float64 | Proporsi bulan aktif setor wajib (0–1) |
| `features.rasio_beban_angsuran` | float64 | Angsuran per bulan / pendapatan estimasi (0–1) |
| `jumlah_diajukan` | int64 | Jumlah pinjaman yang diajukan |
| `tenor_bulan` | int | Tenor dalam bulan |
| `total_simpanan` | int64 | Total simpanan anggota untuk menghitung limit |
| `max_plafond` | int64 | Batas plafond koperasi |

**Response 200:**
```json
{
  "success": true,
  "data": {
    "score": 80,
    "grade": "A",
    "recommendation": "approve",
    "limit_rekomendasi": 9000000,
    "reasons": [
      "riwayat pembayaran sangat baik",
      "simpanan memadai terhadap pinjaman",
      "anggota lama dan loyal"
    ],
    "source": "MOCK_ADINS_SCORING"
  }
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `score` | int | Skor 0–100 |
| `grade` | string | `A` (≥80) \| `B` (≥70) \| `C` (≥60) \| `D` (<60) |
| `recommendation` | string | `approve` (A/B) \| `review` (C) \| `reject` (D) |
| `limit_rekomendasi` | int64 | Limit yang direkomendasikan = faktor[grade] × total_simpanan, dikap max_plafond |
| `source` | string | Selalu `"MOCK_ADINS_SCORING"` — mock deterministic, siap swap ke API Ad Ins |

> **Bobot fitur:** ketepatan_bayar (35%) · rasio_simpanan_pinjaman (25%) · konsistensi_simpanan (15%) · lama_keanggotaan (15%) · rasio_beban_angsuran (10%, inverse)

---

## Dashboard

> Ringkasan operasional koperasi. Role **pengurus** diperlukan.

### GET /api/v1/dashboard
Ambil semua data yang dibutuhkan halaman Dashboard: kartu summary, daftar angsuran jatuh tempo/tunggakan, dan pengajuan pinjaman yang menunggu persetujuan.

**Auth:** Bearer token (role: `pengurus`)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "active_members": 9,
    "total_members": 10,
    "total_savings": 11000000,
    "active_loans": 3,
    "active_loans_outstanding": 12900000,
    "overdue_loans": 2,
    "upcoming_installments_count": 30,
    "upcoming_installments": [
      {
        "installment_id": "uuid",
        "loan_id": "uuid",
        "member_name": "Dewi Lestari",
        "period_no": 6,
        "due_date": "2026-06-07",
        "total_due": 566667,
        "status": "terlambat"
      }
    ],
    "pending_applications_count": 3,
    "pending_applications": [
      {
        "id": "uuid",
        "member_name": "Budi Santosa",
        "amount": 5000000,
        "tenor_months": 12,
        "purpose": "Modal tani musim tanam",
        "grade": "A"
      }
    ]
  }
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `active_members` | int64 | Jumlah anggota dengan status `aktif` |
| `total_members` | int64 | Total semua anggota (semua status) |
| `total_savings` | int64 | Total simpanan koperasi (pokok + wajib + sukarela, Rupiah) |
| `active_loans` | int64 | Jumlah pinjaman berstatus `aktif` |
| `active_loans_outstanding` | int64 | Total sisa pokok seluruh pinjaman aktif (Rupiah) |
| `overdue_loans` | int64 | Jumlah pinjaman berstatus `menunggak` |
| `upcoming_installments_count` | int64 | Total angsuran "perlu perhatian" (terlambat + belum_bayar ≤30 hari) |
| `upcoming_installments` | array | Maks 20 baris; terlambat ditampilkan duluan, lalu urut `due_date ASC` |
| `upcoming_installments[].status` | string | `terlambat` \| `belum_bayar` |
| `upcoming_installments[].total_due` | int64 | Nominal angsuran (Rupiah) |
| `pending_applications_count` | int64 | Total pengajuan berstatus `pending` |
| `pending_applications` | array | Maks 10 baris, urut `created_at DESC` |
| `pending_applications[].grade` | string | `A` \| `B` \| `C` \| `D` (dari credit assessment); `""` jika scoring belum ada |

---

## Portal Anggota

> Endpoint untuk anggota koperasi. Role **anggota** diperlukan.
> `member_id` selalu diambil dari klaim JWT — **tidak pernah dari body request**.

### GET /api/v1/portal/loan-applications
Daftar pengajuan pinjaman milik anggota yang sedang login.

**Auth:** Bearer token (role: `anggota`)

**Response 200:**
```json
{
  "success": true,
  "data": [ /* array LoanApplicationResponse milik anggota ini */ ],
  "meta": { "total": 2 }
}
```

---

### POST /api/v1/portal/loan-applications
Ajukan pinjaman baru. `member_id` diambil dari token — tidak perlu dikirim di body.

**Auth:** Bearer token (role: `anggota`)

**Request Body** (`application/json`):
```json
{
  "amount": 3000000,
  "tenor_months": 6,
  "purpose": "Renovasi rumah"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `amount` | int64 | **ya** | Jumlah pinjaman (> 0) |
| `tenor_months` | int | **ya** | Tenor dalam bulan (> 0) |
| `purpose` | string | tidak | Tujuan penggunaan dana |

**Response 200:** LoanApplicationResponse (dengan `assessment` terisi)

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Field wajib kosong atau `amount` melebihi plafond |
| 409 | `CONFLICT` | Status anggota bukan `aktif` |

---

### GET /api/v1/portal/loans
Daftar pinjaman milik anggota yang sedang login. `schedule` tidak disertakan di list.

**Auth:** Bearer token (role: `anggota`)

**Response 200:**
```json
{
  "success": true,
  "data": [ /* array LoanResponse tanpa schedule, hanya milik anggota ini */ ],
  "meta": { "total": 1 }
}
```

---

### GET /api/v1/portal/loans/:id
Detail pinjaman beserta jadwal angsuran lengkap. Hanya bisa mengakses pinjaman milik sendiri.

**Auth:** Bearer token (role: `anggota`)

**Path Param:**
| Param | Keterangan |
|---|---|
| `id` | UUID pinjaman |

**Response 200:** LoanResponse dengan `schedule` terisi

**Error:**
| HTTP | Code | Kapan |
|---|---|---|
| 403 | `FORBIDDEN` | Pinjaman ada tapi bukan milik anggota ini |
| 404 | `NOT_FOUND` | Pinjaman tidak ditemukan |

---

## Inventory — 🔲 Tier 3, Belum Diimplementasi

Group `/api/v1/inventory` terdaftar dengan tiga lapis middleware: `Auth → RequireRole("pengurus") → RequireModule("inventory")`. Endpoint hanya aktif jika modul `inventory` diaktifkan untuk koperasi tersebut di tabel `coop_module`. Dikerjakan hanya jika core (Tier 1+2) solid.

---

## Flow: Registrasi Anggota via Scan KTP

```
1. Pengurus klik "Scan KTP"
2. FE → POST /api/v1/integrations/adins/ocr/ktp  (upload foto KTP)
3. BE → kirim ke api.co.id, parse response
4. FE ← terima { nik, full_name, address, birth_date, confidence }
5. FE auto-fill form registrasi, pengurus koreksi jika perlu
6. FE → POST /api/v1/members  (submit form)
7. BE → insert ke DB, kembalikan MemberResponse
8. FE ← tampilkan detail anggota baru
```

## Flow: Pengajuan & Pencairan Pinjaman

```
1. Pengurus buka form pengajuan, pilih anggota
2. FE → POST /api/v1/loan-applications  { member_id, amount, tenor_months, purpose }
3. BE → validasi anggota aktif + plafond → scoring → INSERT application + credit_assessment
4. FE ← terima LoanApplicationResponse (assessment.recommendation menentukan warna UI)
5. Pengurus review assessment, klik Approve/Reject
6. FE → POST /api/v1/loan-applications/:id/approve
7. BE → UPDATE status → INSERT loan (aktif) → INSERT installment_schedule (tenor_months baris)
8. FE ← terima { application, loan } — loan.schedule berisi jadwal angsuran
```

## Flow: Bayar Angsuran

```
1. Pengurus buka detail pinjaman, pilih angsuran
2. FE → POST /api/v1/installments/:id/pay  { amount, penalty }
3. BE → INSERT payment (INSERT-ONLY) → hitung paid_amount → update status angsuran
4. BE → cek semua angsuran → update loan.status (aktif/lunas/menunggak)
5. FE ← terima { payment, installment, loan } — UI update real-time
```

---

## Catatan Implementasi

- **`cooperative_id`** selalu dari JWT — jangan kirim di body atau query param.
- **`member_id` di portal** selalu dari JWT — jangan kirim di body, abaikan jika dikirim.
- **Uang = `int64` Rupiah** — tidak ada float di DB, API, maupun FE. FE format dengan `formatRp`.
- **Transaksi keuangan INSERT-ONLY** — `payment` tidak pernah di-UPDATE/DELETE. Saldo & status dihitung dari SUM.
- **`outstanding`** dihitung real-time: `principal - SUM(principal_due WHERE status='lunas')`.
- **Credit scoring mock deterministic** — input sama selalu output sama. Field `source: "MOCK_ADINS_SCORING"` menandai ini mock. Swap ke API Ad Ins = ganti implementasi gateway, kontrak tidak berubah.
- **Jadwal angsuran flat:** `interest = round(principal × rate / 100)`, `principal_per_period = principal ÷ tenor` (integer division), period terakhir absorb sisa pembulatan.
- **NIK** harus tepat 16 digit angka. Unik per koperasi (NIK yang sama boleh ada di koperasi berbeda).
- **`savings_summary`** dihitung real-time via `SUM(amount) FILTER (WHERE savings_type=... AND direction=...)` — tidak ada kolom saldo yang di-update.
- **Error codes** yang digunakan: `VALIDATION_ERROR` (400), `UNAUTHORIZED` (401), `FORBIDDEN` (403), `NOT_FOUND` (404), `CONFLICT` (409), `INTERNAL` (500).
