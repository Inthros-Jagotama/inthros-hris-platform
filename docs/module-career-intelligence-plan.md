# Career Intelligence & Talent Management — Development Plan

> 📅 Versi plan: 2026-08-21 · Status: **BACKEND ✅ TERIMPLEMENTASI PENUH · FRONTEND ✅ SELESAI (5/5 sub-module)**
> 🔎 Berdasarkan audit modul `backend/internal/modules/careerintelligence`, migration `018_career_intelligence` + `086_career_paths` + `157_career_talentmap_settings` (mysql & postgres), dan `frontend/tenant/src/views/modules/career-intelligence/*.vue`.
> 🔀 **Pemisahan transactional vs strategical (2026-08-10):** Career Paths **pindah penuh** dari modul Employee Movement ke module ini — endpoint `/career-intelligence/paths` (ladder-style `name` + `steps[]`), Employee Movement hanya **membaca** `career_paths`/`career_path_steps` untuk promotion eligibility. Lihat log §7.3.
> ✅ **Update 2026-08-21:** seluruh 5 halaman FE selesai (Career Paths, Succession Plans, Career Interests, Gap Analysis, Talent Maps). Talent Map juga tidak lagi input manual — performance & potential dihitung otomatis dari skor final Performance Management + skor Competency terakhir milik employee (item roadmap §9 no. 8, dulu opsional, sekarang selesai). Lihat log §7.5.
> ✅ **Update 2026-08-21 (lanjutan):** Gap Analysis backend tidak lagi stub — `GetGapAnalysis` sekarang query data competency real (syarat dari `competency_score_details` milik assessment terakhir target organization, level aktual dari assessment terakhir employee). Keterbatasan: hanya akurat untuk target organization yang sudah pernah di-assess (tidak ada tabel "syarat kompetensi per posisi" independen di skema ini) — target yang belum pernah di-assess mengembalikan error jelas, bukan hasil 0%-gap yang menyesatkan. Lihat log §7.6.
> ✅ **Update 2026-08-21 (lanjutan lagi):** `career_path_requirements` selesai — ambang batas eligibility promosi (performance/competency/OKR, sebelumnya hardcode global `>=80`) sekarang bisa dikonfigurasi **per langkah career path** lewat 3 kolom nullable baru di `career_path_steps` (migration 158). Lihat log §7.7.
> ⏳ **Sisa TODO:** (1) integrasi Notification untuk event talent mapping/succession (opsional) · (2) tabel "syarat kompetensi per posisi" independen (opsional — supaya Gap Analysis bisa dipakai untuk posisi yang belum pernah di-assess) — keduanya enhancement opsional, bukan blocker.

---

# 1. Objective

Membangun modul **Career Intelligence & Talent Management** yang menangani aspek **strategical** manajemen talenta — konfigurasi, perencanaan, penilaian, dan rekomendasi karier. Modul ini terdiri dari 5 sub-module:

```text
Career Intelligence & Talent Management
    │
    ├── Talent Maps (9-box grid)      → penilaian performa × potensi
    ├── Career Interests              → aspirasi karier karyawan
    ├── Career Paths                  → konfigurasi jenjang karier (planning)
    ├── Gap Analysis                  → kesenjangan kompetensi menuju target posisi
    └── Succession Plans              → rencana suksesi posisi kunci
```

Prinsip utama:

* **Career Intelligence = STRATEGICAL** — ia menyediakan **konfigurasi, perencanaan, dan rekomendasi** (career path, eligibility, talent grid, succession). Ia **tidak pernah mengeksekusi transaksi employment**.
* **Employee Movement = TRANSACTIONAL** — eksekusi promosi/mutasi/demosi/offboarding tetap milik modul Employee Movement.
* Seluruh Primary Key menggunakan UUID/CHAR(36) sesuai pola database existing.
* Modul bersifat **tenant-scoped**, mengikuti pola `NewTenantDBResolver` seperti seluruh modul lain.
* `DependsOn`: `organization`, `employee`, `jobmanagement`, `competency`.

---

# 2. Status Aktual (per 2026-08-21)

## 2.1 Backend — ✅ TERIMPLEMENTASI PENUH

| Sub-module | Status | Endpoint | Catatan |
|---|---|---|---|
| Talent Maps (9-box) | ✅ | 10 | CRUD + grid + profil employee + **generate otomatis** (`/talent-maps/generate`) + settings ambang batas (`/talent-maps/settings`) |
| Career Interests | ✅ | 3 | Create + List + by employee (tidak ada Update/Delete — lihat §8.4) |
| Career Paths (ladder) | ✅ | 5 | **pindah dari EM** — `/career-intelligence/paths` |
| Gap Analysis | ✅ | 1 | `/paths/gap-analysis` (stub kompetensi, lihat catatan di atas) |
| Succession Plans | ✅ | 6 | CRUD + `/successions/gaps` |

Module terdaftar di `main.go`; AutoMigrate menyertakan seluruh 6 model (`CareerTalentMap`, `CareerInterest`, `CareerPath`, `CareerPathStep`, `CareerSuccessionPlan`, `TalentMapSettings`).

## 2.2 Frontend — ✅ SELESAI (5/5)

| Halaman | Route | Status |
|---|---|---|
| Talent Maps | `/career-intelligence/talent-maps` | ✅ **SELESAI (2026-08-21)** (`TalentMaps.vue` — grid 9-box interaktif + generate otomatis + koreksi manual + settings ambang batas) |
| Career Paths | `/career-intelligence/paths` | ✅ **SELESAI** (`CareerPaths.vue` — CRUD ladder-style lengkap) |
| Succession Plans | `/career-intelligence/succession-plans` | ✅ **SELESAI (2026-08-21)** (`SuccessionPlans.vue` — CRUD penuh) |
| Career Interests | `/career-intelligence/interests` | ✅ **SELESAI (2026-08-21)** (`CareerInterests.vue` — create + list, sesuai kapasitas backend) |
| Gap Analysis | `/career-intelligence/gap-analysis` | ✅ **SELESAI (2026-08-21)** (`GapAnalysis.vue`) |
| Succession Gaps (S-5) | `/career-intelligence/successions` | ✅ (`SuccessionGaps.vue`, dibangun sebelumnya) |
| Career Intelligence (shell) | `/career-intelligence` | ✅ hub dengan card navigasi ke seluruh sub-module (`CareerIntelligence.vue`) |

## 2.3 Permissions

Module mengekspos 10 permission (`module.go`). **Catatan implementasi:** route FE untuk sub-module ini sebagian besar digate memakai permission level-module `careerintelligence.view` (string tanpa tanda hubung, hasil generate `tenantseed/seed_rbac.go` dari key `careerintelligence`) — **bukan** string di bawah ini yang berasal dari `module.go`'s `Permissions()`, yang ternyata tidak dipakai untuk seeding RBAC aktual. Lihat catatan di §8.5.

```text
career-intelligence.view
career-intelligence.talent-map.manage   career-intelligence.talent-map.view
career-intelligence.interest.manage     career-intelligence.interest.view
career-intelligence.path.manage         career-intelligence.path.view
career-intelligence.succession.manage   career-intelligence.succession.view
career-intelligence.gap-analysis.view
```

---

# 3. Existing Database Structure

## 3.1 `career_talent_maps` (migration 018)

9-box grid penilaian performa × potensi.

```text
id            CHAR(36) PK
employee_id   CHAR(36)      — karyawan yang dinilai (idx_ctm_employee)
period        CHAR(7)       — format '2026-Q1' (idx_ctm_period)
performance   VARCHAR(20)   — LOW / MEDIUM / HIGH
potential     VARCHAR(20)   — LOW / MEDIUM / HIGH
grid_position VARCHAR(30)   — 9-BOX-1 s.d. 9-BOX-9
notes         TEXT
assessor_id   CHAR(36)      — penilai (manager/HR)
assessed_at   DATE
created_at / updated_at / deleted_at (idx_ctm_deleted_at)
```

## 3.2 `career_interests` (migration 018)

Aspirasi karier karyawan.

```text
id               CHAR(36) PK
employee_id      CHAR(36)      — karyawan (idx_ci_employee)
interest_type    VARCHAR(50)   — LEADERSHIP / SPECIALIST / INTERNATIONAL / ENTREPRENEUR
target_position  VARCHAR(100)
target_department VARCHAR(100)
motivation       TEXT
readiness_level  VARCHAR(20)   — NOW / 1_YEAR / 2_3_YEARS / 3_PLUS
is_active        BOOLEAN
recorded_at      DATE
created_at / updated_at / deleted_at (idx_ci_deleted_at)
```

## 3.3 `career_paths` + `career_path_steps` — ✅ SKEMA TERPADU (018 + 086)

> **Keputusan unifikasi (2026-08-10, log §7.2):** `career_paths` adalah **SATU sumber kebenaran** untuk Career Intelligence (018) **dan** Employee Movement (P1-10). Kolom edge CI lama (`source_title_id, target_title_id, path_type, typical_tenure, requirements, competencies, certifications` + `idx_cp_source/target`) **dihapus dari header** — atribut dipindah ke `career_path_steps`. Edge CI direpresentasikan sebagai **path 2-langkah** (step 1 = source, step 2 = target + atribut CI).

**`career_paths`** — header jenjang (9 kolom terpadu):

```text
id           CHAR(36) PK
name         VARCHAR(100) UNIQUE (uk_career_paths_name)
description  TEXT
is_active    BOOLEAN
created_by   CHAR(36)
updated_by   CHAR(36)
created_at / updated_at / deleted_at (idx_cp_deleted_at)
```

**`career_path_steps`** — langkah jenjang (12 kolom terpadu):

```text
id                    CHAR(36) PK
career_path_id        CHAR(36) FK → career_paths ON DELETE CASCADE
position_id           CHAR(36)   — tanpa FK (pola employee_movements.from_/to_position_id); eksistensi divalidasi service (idx_career_path_steps_position)
sequence              INT        — urutan langkah
minimum_service_months INT       — EM: syarat masa kerja
requirements          TEXT       — EM: persyaratan
path_type             VARCHAR(30)  — CI: PROMOTION / LATERAL / DEMOTION / CROSSFUNCTIONAL (pada step target)
typical_tenure        INT          — CI: months
competencies          TEXT         — CI: JSON list competency IDs
certifications        TEXT         — CI
created_at / updated_at
UNIQUE (career_path_id, sequence)  — uk_career_path_steps_sequence
UNIQUE (career_path_id, position_id) — uk_career_path_steps_position
```

## 3.4 `career_succession_plans` (migration 018)

Rencana suksesi posisi kunci.

```text
id               CHAR(36) PK
position_id      CHAR(36)      — posisi kunci (idx_csp_position)
successor_id     CHAR(36)      — calon suksesor (idx_csp_successor)
readiness_level  VARCHAR(20)   — READY_NOW / READY_1YR / READY_2YR / POTENTIAL
priority_order   INT           — prioritas antar suksesor
target_date      DATE
development_plan TEXT
notes            TEXT
status           VARCHAR(20)   — ACTIVE / COMPLETED / REMOVED
created_at / updated_at / deleted_at (idx_csp_deleted_at)
```

---

# 4. Career Path (strategical planning) — ✅ SELESAI

> Bagian ini dipindahkan dari `module-movement-plan.md` §12.9 (di-archive: `docs/archive/`) (enhancement plan Employee Movement) karena kepemilikan career path kini berada di modul ini (keputusan pemisahan transactional vs strategical, log §7.3).

## 4.1 Konsep

Career Path adalah **konfigurasi/perencanaan jenjang karier**, bukan movement transaction. Ia menjawab pertanyaan: *"dari posisi X, apa langkah promosi berikutnya, berapa lama masa kerja minimal, dan persyaratannya?"*

```text
career_paths            ← header jenjang (nama, deskripsi, aktif)
    │
    └── career_path_steps  ← langkah berurutan
             ├── position_id              (posisi pada langkah ini)
             ├── sequence                 (urutan jenjang)
             ├── minimum_service_months   (syarat masa kerja)
             └── requirements             (persyaratan lain)
```

### Contoh jenjang

```text
Staff
  ↓
Senior Staff
  ↓
Supervisor
  ↓
Manager
  ↓
Senior Manager
```

## 4.2 Semantik create — Ladder-style penuh (endpoint utama FE)

Endpoint `/career-intelligence/paths` menerima bentuk **ladder**:

```http
POST /api/v1/tenant/career-intelligence/paths
Content-Type: application/json

{
  "name": "Staff to Supervisor",
  "description": "Jenjang career IT",
  "is_active": true,
  "steps": [
    { "position_id": "<uuid>", "sequence": 1, "minimum_service_months": 12 },
    { "position_id": "<uuid>", "sequence": 2 }
  ]
}
```

Update bersemantik **full-replace**: klien mengirim daftar steps lengkap; server menghapus seluruh steps lama lalu insert ulang dalam satu transaksi.

## 4.3 Validasi (service layer)

```text
Create / Update Career Path
       ↓
minimal 1 step
       ↓
sequence unik per path
       ↓
posisi unik per path
       ↓
semua position_id merujuk posisi yang BENAR-BENAR ADA
       ↓
(enrich position_name batch)
```

- Error validasi di-return sebagai **400 (validation)**, bukan 500 (constraint cryptic) — unique constraint tetap ada di DB sebagai jaring pengaman.
- `position_id` sengaja **tanpa FK** — konsisten dengan `employee_movements.from_/to_position_id`; eksistensi posisi divalidasi di service via `GetPositionNamesByIDs`.

## 4.4 Konsumen (read-only)

Employee Movement **membaca** career path untuk promotion eligibility:

```text
employeemovement.FindCareerPathStepsByPositionID  → path aktif yang memuat posisi employee
employeemovement.findPromotionNextStep            → step berikutnya (sequence + minimum_service_months)
```

EM **tidak** menulis ke `career_paths`/`career_path_steps` (CRUD dihapus dari EM, log §7.3).

---

# 5. Sub-Module Detail

## 5.1 Talent Maps (9-box grid)

Penilaian dua dimensi: **performance** (LOW/MEDIUM/HIGH) × **potential** (LOW/MEDIUM/HIGH) → `grid_position` 9-BOX-1 s.d. 9-BOX-9.

- `GetTalentGrid` — agregasi per period → 9 kuadran (label + daftar employee).
- `GetEmployeeTalentProfile` — histori penilaian + posisi kini seorang employee.
- Service: `CreateTalentMap` · `GetTalentMapByID` · `ListTalentMaps` (filter period/employee + pagination) · `UpdateTalentMap` · `DeleteTalentMap` · `GetTalentGrid` · `GetEmployeeTalentProfile`.
- Repository: `CreateTalentMap` · `FindTalentMapByID` · `ListTalentMaps` · `UpdateTalentMap` · `DeleteTalentMap` · `GetTalentGrid` · `GetTalentHistoryByEmployee`.

## 5.2 Career Interests

Aspirasi karier karyawan (jenis minat, posisi/departemen target, readiness, motivasi).

- Service: `CreateCareerInterest` · `ListCareerInterests` · `GetEmployeeCareerInterests`.
- Repository: `CreateCareerInterest` · `FindCareerInterestByID` · `ListCareerInterests` · `GetInterestsByEmployee`.

## 5.3 Career Paths — ✅ (detail lengkap di §4)

## 5.4 Gap Analysis

Kesenjangan kompetensi karyawan menuju target posisi.

```http
GET /api/v1/tenant/career-intelligence/paths/gap-analysis?employee_id=<uuid>&target_title_id=<uuid>
```

- Service: `GetGapAnalysis` — matched skills vs total required → `gap_percentage` + rekomendasi. **Real data sejak 2026-08-21** (lihat log §7.6): syarat dari `competency_score_details` milik assessment terakhir target organization, level aktual employee dari assessment terakhir mereka sendiri. Error jelas jika target organization belum pernah di-assess.
- Repository: `GetPositionTitle` · `GetOrgCompetencyRequirements` · `GetEmployeeCompetencyLevels`.

## 5.5 Succession Plans

Rencana suksesi posisi kunci: posisi → suksesor, readiness, prioritas, target date, development plan, status.

- Service: `CreateSuccessionPlan` · `ListSuccessionPlans` · `GetSuccessionPlanByID` · `UpdateSuccessionPlan` · `DeleteSuccessionPlan`.
- Repository: `CreateSuccessionPlan` · `FindSuccessionPlanByID` · `ListSuccessionPlans` · `UpdateSuccessionPlan` · `DeleteSuccessionPlan`.

---

# 6. API Plan — ✅ SEMUA ENDPOINT TERIMPLEMENTASI

> Semua path di bawah group `/api/v1/tenant/career-intelligence`.

## Talent Maps

```http
GET    /talent-maps?period&employee_id&page&per_page
POST   /talent-maps
GET    /talent-maps/grid?period
GET    /talent-maps/employee/:employeeId
GET    /talent-maps/:id
PUT    /talent-maps/:id
DELETE /talent-maps/:id
```

## Career Interests

```http
GET    /interests?employee_id&page&per_page
POST   /interests
GET    /interests/employee/:employeeId
```

## Career Paths (strategical) — ✅ pindah dari Employee Movement (log §7.3)

```http
GET    /paths?page&per_page&keyword
POST   /paths                       (ladder: name + steps[])
GET    /paths/gap-analysis?employee_id&target_title_id     ← sebelum :id (Gin constraint)
GET    /paths/:id
PUT    /paths/:id                   (full-replace steps)
DELETE /paths/:id
```

## Succession Plans

```http
GET    /successions?page&per_page
POST   /successions
GET    /successions/:id
PUT    /successions/:id
DELETE /successions/:id
```

---

# 7. Log Implementasi

> Log berikut dipindahkan dari `module-movement-plan.md` (§3.23, §3.29, §3.30) (di-archive: `docs/archive/`) karena seluruhnya tentang career path yang kini dimiliki modul ini.

## 7.1 Career Path pertama kali dibangun di Employee Movement (migration 086) — 2026-08-10

*dipindah dari module-movement-plan §3.23 (di-archive)*

> ⚠️ **Log historis.** Bagian ini menggambarkan implementasi awal saat career path masih dimiliki Employee
> Movement. Endpoint lama (`/career-paths`) **sudah tidak berlaku** — sejak kepemilikan pindah ke modul ini,
> endpoint aktif adalah `/career-intelligence/paths` (lihat §6). Perombakan ownership dijelaskan di §7.3.

**Career Path** (saat itu plan §12.9) — konfigurasi/planning jenjang karier, bukan movement transaction.

- **Migration 086_career_paths** (mysql + postgres, up + down):
  - `career_paths` — id CHAR(36) PK · `name` UNIQUE · `description` · `is_active` · created_by/updated_by · timestamps.
  - `career_path_steps` — id PK · `career_path_id` FK → `career_paths` **ON DELETE CASCADE** · `position_id` (tanpa FK, pola `employee_movements.from_/to_position_id`) · `sequence` · `minimum_service_months` · `requirements` · **UNIQUE (career_path_id, sequence)** + **UNIQUE (career_path_id, position_id)** + index position_id.
- **Model**: `CareerPath` + `CareerPathStep` (AutoMigrate module ditambah).
- **DTO**: `CreateCareerPathRequest` / `UpdateCareerPathRequest` (full-replace steps) / `CareerPathResponse` (dengan `steps` + `position_name` enrichment) / `PaginatedCareerPathResponse`.
- **Repository**: `CreateCareerPathTx` · `ListCareerPaths` (keyword + pagination) · `FindCareerPathByID` · `ListCareerPathStepsByPathID` (sequence ASC) · `UpdateCareerPathTx` (header + replace seluruh steps, transaksional) · `DeleteCareerPath` (hapus steps eksplisit dulu — konsisten lintas driver, karena SQLite test tidak selalu mengaktifkan FK).
- **Service**: `CreateCareerPath` · `ListCareerPaths` · `GetCareerPathByID` · `UpdateCareerPath` · `DeleteCareerPath` + helper `buildAndValidateCareerPathSteps` (minimal 1 step · sequence unik per path · posisi unik per path · semua `position_id` harus merujuk posisi yang **benar-benar ada**, bukan hanya format UUID) + `enrichCareerPathSteps` (batch `GetPositionNamesByIDs` → `position_name`).
- **Handler + routes**: `GET/POST /career-paths` · `GET/PUT/DELETE /career-paths/:id`.

### Keputusan desain (masih berlaku di CI)

1. **Full-replace steps pada update** — klien mengirim daftar steps lengkap; server menghapus seluruh steps lama lalu insert ulang dalam satu transaksi. Aman karena tidak ada referensi eksternal ke `career_path_steps` (path adalah konfigurasi).
2. **`position_id` tanpa FK** — eksistensi posisi divalidasi di service layer (`GetPositionNamesByIDs`), bukan oleh DB constraint.
3. **UNIQUE (career_path_id, sequence)** mencegah urutan ganda; **UNIQUE (career_path_id, position_id)** mencegah posisi muncul dua kali — keduanya di-validasi pre-insert di service agar error-nya 400 (validation) bukan 500 (constraint cryptic).
4. **`CareerPathID` di-generate di service** (bukan menunggu `BeforeCreate`) agar steps membawa FK yang sama sebelum insert — bug awal yang tertangkap test: steps pernah memakai UUID nol sehingga semua path saling collide.
5. Test list memakai prefix nama unik (`Paginate-*`) karena DB SQLite shared antar-test.

### Validasi

- `go build ./...` ✅ · `go vet` ✅ · `gofmt` bersih ✅
- `go test` module terkait — semua **PASS** ✅
- Test (8): create (steps urut + enrich nama posisi) · duplicate sequence ditolak · posisi tak eksis ditolak · list pagination + keyword · get by id · update full-replace · delete (path + steps bersih) · not-found.

## 7.2 UNIFIKASI schema career_paths (086 + Career Intelligence 018) + apply migration 083–087 — 2026-08-10

*dipindah dari module-movement-plan §3.29 (di-archive)*

> **Latar:** saat verifikasi langsung ke database ditemukan tenant DB hanya 81/86 migration ter-apply (083–087 pending), DAN migration 086 `career_paths` berkonflik dengan tabel `career_paths` dari Career Intelligence (018) yang memakai `CREATE TABLE IF NOT EXISTS` — jika 086 dijalankan apa adanya, kolom `name` tidak akan pernah dibuat. Keputusan user (2026-08-10): **unifikasi penuh** — satu skema `career_paths` + `career_path_steps` untuk kedua modul.

### Yang dikerjakan

- **Migration 086 di-rewrite** (mysql + postgres, up + down) menjadi skema TERPADU (idempotent):
  - `career_paths` = header jenjang: `id, name, description, is_active, created_by, updated_by, created_at, updated_at, deleted_at` + `uk_career_paths_name`.
  - Kolom edge CI lama (`source_title_id, target_title_id, path_type, typical_tenure, requirements, competencies, certifications` + `idx_cp_source/target`) **dihapus** — atribut dipindah ke `career_path_steps`.
  - `career_path_steps` = langkah terpadu: `position_id, sequence, minimum_service_months, requirements` (EM) + `path_type, typical_tenure, competencies, certifications` (CI, pada step target) + `idx_career_path_steps_position` + unique `(career_path_id, sequence)` & `(career_path_id, position_id)` + FK CASCADE.
- **Career Intelligence di-refactor** (model/dto/repo/service):
  - `CareerPath` CI = header (name/description/is_active/soft delete); tambah `CareerPathStep` CI.
  - Edge CI disimpan sebagai **path 2-langkah**: step 1 = source, step 2 = target + atribut edge.
  - `CreateCareerPathTx` (transaksi header+steps), `ListCareerPathStepsByPathID(s)`, `FindCareerPathsBySource` (via step sequence 1), `DeleteCareerPath` (steps hard + header soft).
  - Respons CI backward-compatible: tetap mengekspos `source_title_id/target_title_id/path_type/typical_tenure/requirements/competencies/certifications` (diderivasi dari steps) + tambahan `name` & `steps`.

### Apply migration ke tenant & verifikasi DB

- `installer.exe migrate --company=df687f34-e580-40c5-8935-73180fb5fd3c` → **5 migration applied** (083–087).
- Verifikasi `information_schema`: `schema_migrations` 083–087 hadir · `career_paths` = 9 kolom header terpadu · `career_path_steps` = 12 kolom + index + FK. DB tenant kini **86/86** migration ter-apply.

### Validasi

- `go build ./...` ✅ · `go vet` kedua modul ✅ · `go test ./internal/modules/careerintelligence/ ./internal/modules/employeemovement/` — semua **PASS** ✅.
- DB tenant kini **86/86** migration ter-apply.

### Hasil code review & perbaikan

- **Finding 1 (data loss di migration 086):** tidak ada data `career_paths` pada semua tenant (0 baris) — tidak ada data yang hilang. Risiko teoritis untuk tenant masa depan dengan data CI didokumentasikan.
- **Finding 2 (soft-delete mismatch):** diperbaiki ✅ — `DeletedAt gorm.DeletedAt` + `idx_cp_deleted_at` ditambahkan ke model EM `CareerPath` agar semua query EM otomatis memfilter `deleted_at IS NULL`.
- **Finding 2b (nama unik vs soft delete):** `FindCareerPathByName` (CI) memakai `Unscoped()` agar nama path yang sudah soft-deleted tetap terdeteksi/terpesan.
- **Finding 3 (name collision):** aman — `buildCareerPathName` memakai loop akhiran `-2/-3/...`.
- **Finding 4 (ORDER BY sequence):** aman — semua fetch steps memakai `Order("sequence ASC")`.
- **Finding 5 (validasi asimetris CI vs EM):** diterima sebagai trade-off — CI membuat path 2-step deterministik yang selalu lolos validasi EM; didokumentasikan sebagai kontrak skema terpadu.

## 7.3 PEMISAHAN TRANSACTIONAL vs STRATEGICAL — Career Paths pindah penuh ke Career Intelligence — 2026-08-10

*dipindah dari module-movement-plan §3.30 (di-archive)*

> **Latar:** user meminta career paths berada di modul **career-intelligence** (strategical planning), bukan di **employeemovement** (transactional eksekusi). Setelah unifikasi skema 086, kedua modul berbagi tabel `career_paths`/`career_path_steps` — namun kepemilikan CRUD masih di EM. Keputusan user: **full ownership pindah ke CI** + path `/career-intelligence/paths` + CI mendukung **ladder-style penuh** (`name` + `steps[]`).
>
> ⚠️ Catatan untuk pembaca: bagian ini menyebut edge CI "tetap dipertahankan" — itu kondisi saat log ditulis.
> **Dead code edge CI kemudian dihapus** pada §7.4 (ladder-style menjadi satu-satunya bentuk create).

### Backend Career Intelligence (pemilik baru)

- **dto.go**: tambah `CreateCareerPathStepRequest`, `CreateCareerPathLadderRequest` (name required + description + is_active + steps min 1), `UpdateCareerPathRequest` (full-replace), perluas `CareerPathStepResponse` (`minimum_service_months` + `requirements`) dan `CareerPathResponse` (`description`); alias `PaginatedCareerPathResponse = PaginatedResponse`.
- **repository.go**: tambah `GetPositionNamesByIDs` (batch, tabel `positions`), `UpdateCareerPathTx` (full-replace header + steps, description bisa di-NULL-kan), `ListCareerPaths` + keyword filter.
- **service.go**: tambah `CreateCareerPathLadder` (validasi: min 1 step, sequence unik, posisi unik, posisi benar-benar ada), `GetCareerPathByID`, `UpdateCareerPath` (full-replace); `careerPathToResponse` kini mengisi `description` + `position_name` + min service/requirements (batch resolve).
- **handler.go**: tambah `CreateCareerPathLadder`, `GetCareerPathByID`, `UpdateCareerPath`; `ListCareerPaths` membaca `keyword`.
- **routes.go**: `GET/POST /paths`, `GET /paths/gap-analysis`, `GET/PUT/DELETE /paths/:id` (gap-analysis sebelum `:id` — Gin constraint).
- **module.go**: AutoMigrate kini menyertakan `CareerPathStep`.

### Backend Employee Movement (transactional — read-only)

- **routes.go**: group `/career-paths` DIHAPUS (endpoint pindah ke CI).
- **handler.go**: 5 handler career path dihapus (Create/List/GetByID/Update/Delete).
- **service.go**: 7 method CRUD career path dihapus; **read-only eligibility dipertahankan** (`FindCareerPathStepsByPositionID`, `FindCareerPathByID`, `findPromotionNextStep`).
- **repository.go**: `CreateCareerPathTx`/`ListCareerPaths`/`UpdateCareerPathTx`/`DeleteCareerPath` dihapus; **read methods dipertahankan** (`FindCareerPathByID`, `ListCareerPathStepsByPathID`, `FindCareerPathStepsByPositionID`, `GetPositionNamesByIDs`).
- **dto.go + model ToResponse**: DTO career path + `ToResponse` dihapus.
- **module.go**: AutoMigrate `CareerPath`/`CareerPathStep` dihapus (pemilik skema = CI).

### Test & FE

- CI: handler/service/repo test disesuaikan ke ladder-style (seed tabel `positions`).
- EM: `careerpath_test.go` dihapus; `eligibility_test.go` seeding via gorm langsung.
- FE tenant: `CareerPaths.vue` memakai `/api/v1/tenant/career-intelligence/paths` + namespace i18n **`career_paths.*`** (baru); route `/career-intelligence/paths` module gate **`career-intelligence`**; sidebar item pindah dari Operations (EM) ke **Strategic** (CI, `careerintelligence.view`).

### Validasi

- `go build ./...` ✅ · `go vet` CI + EM ✅ · `go test ./internal/modules/...` — **17/17 modul PASS** ✅ · `npm run build` (tenant FE) ✅.

## 7.4 Dead code cleanup — edge CI dihapus (2026-08-10)

Setelah ladder-style menjadi satu-satunya bentuk create, dead code edge CI dihapus:

- Handler `CreateCareerPath` (edge), service `CreateCareerPath` + `buildCareerPathName` + `mustUUID`, repo `FindCareerPathByName`, DTO `CreateCareerPathRequest` — **dihapus total**.
- `GetPositionTitle` **dipertahankan** (dipakai gap analysis); field `PathType`/`TypicalTenure`/`Competencies`/`Certifications` tetap hidup di model/response (bagian skema terpadu).
- Test edge diganti `TestService_CreateCareerPathLadder_Success`; test handler di-rename konsisten.
- Validasi: `go build ./...` ✅ · `go vet` ✅ · `go test ./internal/modules/...` — **17/17 modul PASS** ✅ · `npm run build` (tenant FE) ✅.

## 7.6 Gap Analysis — real query menggantikan stub hardcoded (2026-08-21)

**Latar:** `GetGapAnalysis` sejak awal hardcode `matchedSkills=3`, `totalRequired=8`, dan 3 rekomendasi tetap — komentar di kode secara eksplisit menyebut "would query competency module in production". User meminta ini diselesaikan.

**Temuan skema:** tidak ada tabel "syarat kompetensi per posisi/organisasi" independen di modul Competency — required level hanya tersimpan per-template assessment (`CompetencyAssessmentTemplateCompetency.RequiredLevel`), yang terhubung ke organisasi hanya secara transien lewat `CompetencyEventTarget.OrganizationID` saat sebuah assessment event dijalankan. Keputusan user: pakai data assessment yang sudah ada (bukan bangun tabel requirement baru) — `competency_scores.organization_id` punya unique index (satu row per organisasi), jadi "syarat kompetensi organisasi X" = `competency_score_details.standard_level` dari assessment terakhir organisasi itu sendiri.

**Implementasi:**
- **Repository baru** (`careerintelligence/repository.go`): `GetOrgCompetencyRequirements(ctx, orgID)` — join `competency_score_details`+`competencies` dari row `competency_scores` terbaru milik org tsb; `GetEmployeeCompetencyLevels(ctx, employeeID)` — map competency_id→employee_level dari assessment terakhir employee (query cross-module langsung ke tabel, pola sama seperti `workforceintelligence.GetQualityOfHire`).
- **Service** (`GetGapAnalysis`): ambil requirements + employee levels → hitung matched/gap per kompetensi → `MatchedSkills`/`TotalRequired`/`GapPercentage` real, `Recommendations` dibangun per kompetensi yang belum terpenuhi (priority HIGH/MEDIUM/LOW berdasar besar gap), `EstimatedTimeline` diturunkan dari gap percentage (3-6/6-12/12-18/18-24 bulan). Target organization yang belum pernah di-assess → error 400 jelas ("belum ada standar kompetensi tercatat"), bukan hasil 0%-gap yang menyesatkan.
- **Handler**: error `GetGapAnalysis` diubah dari 500 ke 400 (kegagalan validasi data, bukan server error).
- **Test**: fixture `seedGapAnalysisCompetencyData` (helpers_test.go) membuat tabel `competencies`/`competency_scores`/`competency_score_details` + data 2 kompetensi (1 terpenuhi employee, 1 tidak) untuk `TestService_GetGapAnalysis_Success`/`TestHandler_GetGapAnalysis_Success`; tambah `TestService_GetGapAnalysis_NoRequirements` untuk kasus target belum pernah di-assess.

### Validasi

- `go build ./...` ✅ · `go vet` ✅ · `go test ./internal/modules/careerintelligence/...` — **PASS** ✅.

## 7.7 `career_path_requirements` — ambang batas eligibility per langkah career path (2026-08-21)

**Latar:** `employeemovement`'s promotion/movement eligibility check membanding skor performance/competency/OKR karyawan terhadap threshold `>=80` yang **hardcode global** (`eligibilityMinPerformanceScore`/`eligibilityMinCompetencyScore`/`eligibilityMinOKRScore` di `service.go`) — sama untuk semua posisi/career path, tidak bisa dikonfigurasi per jenjang. Tenure (`minimum_service_months`) sudah bisa di-override per step sejak awal; skor belum.

**Keputusan desain:** tambah 3 kolom nullable ke `career_path_steps` (bukan tabel terpisah — konsisten dengan skema terpadu 086/087 di mana `career_path_steps` sudah jadi tempat gabungan atribut EM+CI per langkah): `min_performance_score`, `min_competency_score`, `min_okr_score` (`DECIMAL(5,2)`, semua opsional). `NULL` = pakai default global; diisi = override untuk promosi ke langkah itu.

**Implementasi:**
- **Migration 158** (mysql + postgres, up + down): `ALTER TABLE career_path_steps ADD COLUMN min_performance_score/min_competency_score/min_okr_score`.
- **Model** (`CareerPathStep`, careerintelligence & employeemovement — dua struct terpisah yang memetakan tabel yang sama): tambah 3 field `*float64`.
- **CI dto/service**: `CreateCareerPathStepRequest`/`CareerPathStepResponse` tambah field (binding `omitempty,gte=0,lte=100`); `buildCareerPathSteps`/`careerPathToResponse` memetakan field baru — sehingga FE `CareerPaths.vue` bisa create/edit/lihat threshold per step lewat endpoint `/career-intelligence/paths` yang sudah ada (tidak ada endpoint baru).
- **EM service (`GetPromotionEligibility`)**: `findPromotionNextStep` sekarang juga membaca 3 kolom baru dari `nextStep` dan meng-override `promo.Min*Score` (default diisi dari konstanta global sebelum kemungkinan di-override). `evaluatePromotionRules` diberi 3 parameter threshold baru (bukan lagi memakai konstanta langsung), dipanggil dengan `promo.Min*Score` yang sudah final. `movement-eligibility` (`evaluateDefaultRules`) **tidak diubah** — endpoint itu general "boleh pindah?" tanpa konteks target step tertentu, jadi tetap threshold global (scope §9 #6 secara semantik adalah untuk promosi ke posisi spesifik, bukan movement umum).
- **FE (`CareerPaths.vue`)**: bagian collapsible "Ambang batas eligibility promosi (opsional)" per step di dialog create/edit, 3 `InputNumber` (placeholder "Default (80)"), otomatis expanded saat edit step yang sudah punya override.

**Verifikasi:**
- `go build ./...` ✅ · `go vet` ✅.
- `go test ./internal/modules/careerintelligence/...` — **PASS** ✅.
- `go test ./internal/modules/employeemovement/...` — 15 test gagal (`organizations.code NOT NULL` dll) **diverifikasi pre-existing** via `git stash` diff (list identik dengan/tanpa perubahan ini, hanya durasi run yang beda). Test integrasi `TestGetPromotionEligibility_WithCareerPath` (yang seharusnya menguji jalur baru ini end-to-end) termasuk yang terblokir bug fixture pre-existing tsb, jadi ditambahkan unit test murni `TestEvaluatePromotionRules_PerStepThresholdOverride` (tanpa DB) untuk memverifikasi logika override benar: skor 70 gagal terhadap default 80, lolos terhadap override 60 — **PASS** ✅.
- `npm run build` (tenant FE) ✅.

---

# 8. Frontend Plan

## 8.1 Career Paths — ✅ SELESAI

`views/modules/CareerPaths.vue` (route `/career-intelligence/paths`, module gate `career-intelligence`):

- **List**: DataTable lazy + pagination, pencarian keyword, SkeletonTable, kolom nama/deskripsi, rantai steps (dengan panah antar posisi), jumlah step, status aktif, updated_at, aksi view/edit/delete.
- **Detail dialog**: visual **career ladder timeline vertikal** — step bernomor dengan connector line, tag masa kerja minimal + requirements, langkah pertama disorot.
- **Dialog create/edit**: form nama + deskripsi + toggle aktif + **editor steps dinamis** (tambah/hapus/reorder naik-turun, posisi via dropdown `organization = position` dari `/organizations`, InputNumber masa kerja minimal, requirements) — sequence dihitung otomatis dari urutan.
- **Delete**: ConfirmDeleteDialog + error handling.
- i18n namespace **`career_paths.*`** (en/id); sidebar di grup **Strategic**.
- Payload create/update cocok dengan DTO backend (sequence = idx+1, opsional pakai `undefined`, full-replace steps; `description: ''` eksplisit untuk clear karena backend pointer-presence).

## 8.2 Career Intelligence shell — ✅ SELESAI

`views/modules/CareerIntelligence.vue` (route `/career-intelligence`) adalah hub dengan card navigasi ke seluruh sub-module (Talent Maps, Career Paths, Succession Plans, Succession Gaps, Gap Analysis, Development Training). Tidak ada lagi card "Coming soon".

## 8.3 Halaman FE — ✅ SEMUA SELESAI (2026-08-21)

| # | Halaman | Route | Endpoint konsumen |
|---|---|---|---|
| 1 | Talent Maps (9-box grid + generate otomatis) | `/career-intelligence/talent-maps` | `/talent-maps`, `/talent-maps/grid`, `/talent-maps/generate`, `/talent-maps/settings`, `/talent-maps/employee/:id` |
| 2 | Career Interests (create + list) | `/career-intelligence/interests` | `/interests`, `/interests/employee/:id` |
| 3 | Gap Analysis | `/career-intelligence/gap-analysis` | `/paths/gap-analysis` |
| 4 | Succession Plans (CRUD) | `/career-intelligence/succession-plans` | `/successions` |

Pola FE mengikuti `CareerPaths.vue` (DataTable lazy + SkeletonTable + ConfirmDeleteDialog + useI18n + module gate `career-intelligence`).

## 8.4 Career Interests — keterbatasan backend (create + list saja)

Backend Career Interests hanya menyediakan `POST /interests`, `GET /interests`, dan `GET /interests/employee/:employeeId` — **tidak ada Update/Delete**. `CareerInterests.vue` dibangun sesuai kapasitas ini: form tambah + tabel daftar (filter per employee) + dialog detail, tanpa tombol edit/hapus. Kalau ke depan dibutuhkan koreksi/penghapusan, perlu ditambah endpoint `PUT`/`DELETE /interests/:id` di backend dulu.

## 8.5 Talent Map — generate otomatis dari Performance Management + Competency (2026-08-21)

Sebelumnya `performance`/`potential` diinput manual (LOW/MEDIUM/HIGH langsung dari user). Sekarang otomatis:

- **`GenerateTalentMap`** (`service.go`): ambil skor final performance evaluation terakhir yang `COMPLETED`/`ACTUAL_APPROVED` via `PerformanceProvider.LatestFinalScore` (reuse `performance.Repository.LatestFinalScoreByEmployee`), dan skor competency terbaru via `CompetencyProvider.LatestScore` (reuse `competency.Repository.LatestScoreByEmployee`) — kedua provider di-wire di `main.go` pakai adapter yang sama persis dengan `employeemovement`'s promotion eligibility (`performanceEligibilityAdapter`/`competencyEligibilityAdapter`, structural typing Go).
- Skor (skala 0-100) dibanding LOW/MEDIUM/HIGH memakai **`TalentMapSettings`** — singleton tenant settings baru (migration `157_career_talentmap_settings`, default low_max=50/high_min=80 untuk performance & potential masing-masing), bisa diubah via `GET/PUT /talent-maps/settings` dan dialog pengaturan (ikon gear) di `TalentMaps.vue`.
- Jika employee belum punya evaluasi performance final atau assessment competency, `GenerateTalentMap` menolak dengan pesan error jelas (bukan silently default).
- CRUD manual (`POST/PUT /talent-maps`) tetap ada di backend untuk fleksibilitas, tapi FE hanya expose **Generate** (create otomatis) + **Koreksi** (edit performance/potential/notes hasil generate, untuk override HR) + **Delete**.
- **Catatan verifikasi (2026-08-21):** dicek langsung ke tenant DB `hris_pt-inthros-jago-utama` — belum ada karyawan yang memenuhi syarat generate: performance evaluation yang ada berstatus `APPROVED` (bukan `COMPLETED`/`ACTUAL_APPROVED`), dan `competency_scores` untuk employee masih kosong. Fitur generate baru bisa dipakai setelah alur approval performance & minimal satu competency assessment selesai untuk karyawan terkait.

---

# 9. Roadmap / TODO

| # | Item | Status | Area |
|---|---|---|---|
| 1 | Career Paths backend (CRUD ladder) | ✅ | BE |
| 2 | Career Paths FE (`CareerPaths.vue`) | ✅ | FE |
| 3 | Talent Maps / Interests / Gap Analysis / Succession Plans backend | ✅ | BE |
| 4 | Halaman FE Talent Maps / Interests / Gap Analysis / Succession Plans | ✅ **SELESAI (2026-08-21)** — semua 4 halaman dibangun | FE |
| 5 | Gap Analysis FE (form employee + target title, hasil + rekomendasi) | ✅ **SELESAI (2026-08-21)** — `GapAnalysis.vue` | FE |
| 6 | `career_path_requirements` (rule eligibility terstruktur) | ✅ **SELESAI (2026-08-21)** — migration 158, lihat §7.7 | DB/BE |
| 7 | Notification untuk event talent mapping / succession / interest | ⏳ opsional | BE/FE |
| 8 | Integrasi hasil performance/competency ke talent grid (input 9-box otomatis) | ✅ **SELESAI (2026-08-21)** — `GenerateTalentMap` + `TalentMapSettings` (migration 157), lihat §8.5 | BE/FE |
| 9 | Gap Analysis backend real (bukan stub hardcoded) | ✅ **SELESAI (2026-08-21)** — query `competency_score_details` real, lihat §7.6 | BE |
| 10 | Tabel syarat kompetensi per posisi independen (opsional, agar Gap Analysis bisa dipakai untuk posisi yang belum pernah di-assess) | ⏳ opsional | DB/BE |

---

# 10. Integrasi dengan Modul Lain

```text
                 Career Intelligence (STRATEGICAL)
                            │
        ┌───────────────────┼───────────────────────┐
        │                   │                       │
        ▼                   ▼                       ▼
   Employee Movement   Performance Mgmt        Organization
   (promotion elig.)   (KPI/OKR/competency)    (positions)
        │                   │                       │
        └───────────────────┼───────────────────────┘
                            ▼
              Konfigurasi & rekomendasi karier
                            │
                            ▼
                 Employee Movement (TRANSACTIONAL)
                 — eksekusi promosi/mutasi/demosi
```

- **Employee Movement** — membaca `career_paths`/`career_path_steps` untuk promotion eligibility (`movement-eligibility`, `promotion-eligibility`). Tidak menulis.
- **Performance Management & Competency** — skor final menjadi **input otomatis** untuk talent mapping (9-box, sejak 2026-08-21, lihat §8.5) dan untuk eligibility promosi (Employee Movement). Career Intelligence tidak menghitung KPI/OKR/competency sendiri, hanya membaca skor final via provider (`PerformanceProvider`/`CompetencyProvider`) lalu membanding LOW/MEDIUM/HIGH.
- **Organization / Job Management** — `position_id` pada steps/talent map/succession merujuk tabel `positions`; eksistensi divalidasi via `GetPositionNamesByIDs`/`GetPositionTitle`.
- **Employee** — `employee_id` pada talent map/interest/succession merujuk karyawan.
- **Recruitment (operasional)** — integrasi strategis (internal candidate via career path, succession fallback → external recruitment) dikelola di **`docs/archive/module-recruitment-strategic-layer-plan.md`** (S-4/S-5); Recruitment hanya mengeksekusi, tidak menentukan eligibility (scoping 2026-08-12).
