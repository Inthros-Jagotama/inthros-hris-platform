# Analisa Perhitungan Job Management Score

> 🔗 **Index dokumentasi:** [`docs/README.md`](README.md)  
> **Terkait:** [`analisis-blueprint-vs-existing.md`](analisis-blueprint-vs-existing.md) · [`platform-architecture-design.md`](platform-architecture-design.md)

> **Dokumen hasil analisa** — mempelajari cara perhitungan *Job Management Score* dari
> `docs/backlog/JobValueCalculator.php` (implementasi legacy PHP/Laravel) dan membandingkannya
> dengan implementasi backend Go saat ini (`backend/internal/modules/jobmanagement/`).
>
> Tanggal analisa: 04 Agustus 2026
> Status: ✅ **Implementasi selesai** (05 Agu 2026) — `calculator.go` sudah dibuat dan
> terintegrasi (recalc otomatis tiap section disimpan, `is_complete`/`completed_at` dipersist,
> breakdown skor ditampilkan dari field `components` di halaman score).
> Keputusan user: **Communication & Influencing Skills level 1→1, 2→3, 3→6** (sudah sesuai
> `MAP_COMMUNICATION`); **bobot (weight) TIDAK mengubah skor** — rumus legacy tetap dipakai.

---

## 1. Ringkasan Eksekutif

`JobValueCalculator.php` adalah mesin kalkulasi **Job Value / Skor Jabatan** dari aplikasi legacy
(Laravel). Kalkulator ini menghitung skor dari **10 komponen jabatan** berdasarkan *level* yang
dipilih pada tiap section job management, lalu:

1. **Memetakan level → poin** melalui 5 tabel mapping (`MAP_DEFAULT`, `MAP_EXTENDED`,
   `MAP_LINEAR_5`, `MAP_LINEAR_8`, `MAP_COMMUNICATION`).
2. **Menggabungkan komponen** (sebagian dengan perkalian, sebagian penjumlahan).
3. **Menghitung 2 total skor**: `job_value_with_financial` dan `job_value_without_financial`
   (tergantung apakah jabatan memiliki wewenang keuangan / `is_authorized`).
4. **Menyimpan hasil** ke tabel `job_management_scores` via `updateOrCreate(organization_id)`,
   lengkap dengan rincian komponen (`components`), poin sub-komponen (`sub_component_points`),
   waktu kalkulasi (`calculated_at`) dan status kelengkapan (`is_complete`).

**Kondisi backend Go saat ini**: endpoint `PUT /job-management/scores/org/:orgId`
(`Service.UpsertJobScore`) **belum menjalankan mesin kalkulasi** — ia hanya menyimpan nilai yang
dikirimkan client (`job_value_with_financial`, `job_value_without_financial`, dll.) apa adanya.
Artinya mesin kalkulasi dari PHP **belum dipindahkan** ke Go.

---

## 2. Sumber yang Dianalisa

| File | Peran |
|---|---|
| `docs/backlog/JobValueCalculator.php` | Mesin kalkulasi legacy (referensi utama) |
| `backend/internal/modules/jobmanagement/model.go` | Model GORM Go: `JobScore` (9.17), `JobValue`, dll. |
| `backend/internal/modules/jobmanagement/dto.go` | DTO `UpdateJobScoreRequest` / `JobScoreResponse` |
| `backend/internal/modules/jobmanagement/service.go` | `UpsertJobScore`, `GetJobScoreByOrganization`, `ListJobScores` |
| `backend/internal/modules/jobmanagement/repository.go` | `UpsertJobScore`, `FindJobScoreByOrganizationID` |
| `backend/internal/pkg/migrator/migrations/tenant/mysql/009_job_management.sql` | Skema tabel 9.1–9.18 |
| Migration 033–042 | Seed `job_management_values` (tipe & level) |
| `frontend/tenant/src/views/modules/job/sections/JobScoreSection.vue` | Halaman skor (konsumsi API) |
| `docs/backlog/JobManagementValuesTableSeeder.php` | Seeder nilai `job_management_values` legacy — struktur `type` vs `code`, level & deskripsi tiap tipe |

---

## 3. Tabel Pemetaan Level → Poin (Mapping Tables)

Kalkulator memakai 5 tabel mapping. Jika `level` tidak ditemukan di tabel, nilai dipatok ke
level tertinggi yang tersedia (`mapPoints` fallback ke `max(array_keys($map))`).

### 3.1 `MAP_DEFAULT` — dipakai untuk mayoritas komponen (5 level)

| Level | 1 | 2 | 3 | 4 | 5 |
|---|---|---|---|---|---|
| Poin | 1 | 3 | 6 | 10 | 15 |

### 3.2 `MAP_EXTENDED` — komponen 8 level

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
|---|---|---|---|---|---|---|---|---|
| Poin | 1 | 3 | 6 | 10 | 15 | 21 | 28 | 36 |

### 3.3 `MAP_LINEAR_5` — 5 level linier

| Level | 1 | 2 | 3 | 4 | 5 |
|---|---|---|---|---|---|
| Poin | 1 | 2 | 3 | 4 | 5 |

### 3.4 `MAP_LINEAR_8` — 8 level linier

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
|---|---|---|---|---|---|---|---|---|
| Poin | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |

### 3.5 `MAP_COMMUNICATION` — komunikasi & influencing (3 level)

| Level | 1 | 2 | 3 |
|---|---|---|---|
| Poin | 1 | 3 | 6 |

---

## 4. Komponen Perhitungan

Setiap organisasi (`organization_id`) dihitung dari data section job management terkait.
Berikut rincian tiap komponen.

### 4.1 Pendidikan & Pengalaman (`calculateEducationExperience`)

**Sumber data**: `job_management_education_experiences` (1 record per org)
→ relasi `educationGrade` dan `experienceGrade` (nilai level dari `job_management_values`).

```
education_points  = MAP_DEFAULT(level_education)
experience_points = MAP_DEFAULT(level_experience)
score             = education_points > 0 DAN experience_points > 0
                    ? education_points * experience_points
                    : 0
```

> ⚠️ **Catatan mapping ke Go**: di legacy, `educationGrade`/`experienceGrade` diambil dari
> level tabel `job_management_values` dengan tipe `education` & `experience`.
> Di Go, kolomnya bernama `EducationID` / `ExperienceID` (lihat migration 045 & 046) dan
> **tidak ada seed tipe `education`/`experience` di migration 033–042** — namun seeder legacy
> mendefinisikan keduanya (5 level: SMA/D3/S1/S2/S3 dan 0–2 th/3–5 th/6–8 th/9–11 th/>12 th,
> lihat 8.6). Perlu di-seed ulang atau verifikasi data aktual.

### 4.2 Potensi Psikologi (`calculatePotentials`)

**Sumber data**: `job_management_potency_competencies` (banyak record per org)
→ `jobManagementValue.type == 'Psychological'`.

```
levels        = [level dari semua record bertipe 'Psychological']
average_level = ceil(rata-rata levels)   // dibulatkan ke atas
score         = average_level ada ? MAP_DEFAULT(average_level) : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go tidak ada tipe `Psychological` di `job_management_values`
> — di legacy `Psychological` adalah **nama grup** (`type`), sedangkan di Go grup tsb dipecah
> menjadi **slug per kompetensi** (lihat 8.6). Tipe psikologi di Go (migration 041):
> `kecerdasan` (5 level), `innovation_creativity`, `self_confidence`, `flexibility`, `tenacity`,
> `continuous_learning` (masing-masing 8 level). Kalkulator Go harus mengumpulkan level dari
> **keenam slug tersebut** (gabungan = setara `Psychological` legacy), lalu rata-rata (ceil).

### 4.3 Kompetensi (`calculateCompetencies`) — Technical × Managerial × Communication

**Sumber data**: `job_management_potency_competencies`
- Technical: record dengan `jobManagementValue.type == 'Technical'` **atau**
  `competency.field == 'Technical Competency'` (via relasi `competency` + `competencyValue`).
- Managerial: record dengan `jobManagementValue.type == 'Managerial'` **atau**
  `competency.field == 'Manajerial'`.
- Communication: record dengan `jobManagementValue.type == 'Communicating & Influencing Skill'`
  (atau descriptions mengandung teks tsb) — **jika tidak ditemukan, level default = 1**.

```
technical_average    = ceil(rata-rata level Technical)     → technical_points    = MAP_EXTENDED(average)
managerial_average   = ceil(rata-rata level Managerial)    → managerial_points   = MAP_DEFAULT(average)
communication_points = MAP_COMMUNICATION(level_comm ?? 1)   // default level 1 → 1 poin

score = technical_points > 0 DAN managerial_points > 0 DAN communication_points > 0
        ? technical_points * managerial_points * communication_points
        : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go tidak ada tipe `Technical`/`Managerial`/`Communicating &
> Influencing Skill`. Migration 042 memakai tipe per kompetensi:
> - **16 tipe technical** (masing-masing 8 level): `competency_based_human_resources_management`,
>   `competency_development`, `people_development`, `career_management`, `hr_assessment`,
>   `recruitement_selection`, `job_analysis_evaluation`, `organizational_development`,
>   `human_resources_information_system`, `workload_analysis`, `performance_apraisal`,
>   `remuneration_manajemen`, `reward_punisment_management`, `health_safety_environment`,
>   `hubungan_industrial`, `budgeting`.
> - **6 tipe managerial** (masing-masing 5 level): `integrity`, `achievement_orientation`,
>   `building_partnership`, `planning_organizing`, `leadership`, `developing_others`.
> - **Communication**: tipe `communicating_influencing_skill` (3 level: Berkomunikasi, Alasan,
>   Perubahan Perilaku) — migration 039.
>
> Backend Go perlu mengelompokkan tipe mana yang "Technical" vs "Managerial" (daftar tetap
> di atas — di legacy keduanya adalah nama grup `type`, di Go dipecah jadi slug per kompetensi,
> lihat 8.6), atau via kolom `job_management_competency_groups` (9.18) yang punya
> `category ENUM('technical','managerial')`.

### 4.4 Pemecahan Masalah (`calculateProblemSolving`) — Environment × Challenge

**Sumber data**: `job_management_potency_competencies`
- Filter record dengan `jobManagementValue.type == 'Problem Solving & Decision Making'`.
- Environment: record yang `descriptions` mengandung teks `'Environment'`.
- Challenge: record yang `descriptions` mengandung `'Chalenge'` atau `'Challenge'`.

```
environment_points = MAP_EXTENDED(level_environment)
challenge_points   = MAP_DEFAULT(level_challenge)
score              = environment_points > 0 DAN challenge_points > 0
                     ? environment_points * challenge_points
                     : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go, tipe `thinking_environment` (8 level) dan
> `thinking_chalenge` (5 level) — migration 040. Pemetaan langsung: Environment →
> `thinking_environment`, Challenge → `thinking_chalenge`.

### 4.5 Kewenangan Keuangan (`calculateFinancialAuthority`)

**Sumber data**: `job_management_financials` (1 record per org)
→ relasi `cashValue`, `authorityValue`, `impactValue`.

```
has_authority = (bool) is_authorized

money_points      = has_authority ? MAP_EXTENDED(level_cash)      : 0
authority_points  = MAP_EXTENDED(level_authority)                  // selalu dihitung
impact_points     = MAP_EXTENDED(level_impact)                     // selalu dihitung

if has_authority:
    score = money_points > 0 DAN authority_points > 0 DAN impact_points > 0
            ? money_points * authority_points * impact_points
            : 0
else:
    score = authority_points > 0 DAN impact_points > 0
            ? authority_points * impact_points
            : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go, `JobFinancial` memakai kolom `JobManagementValueCashID`,
> `JobManagementValueAuthorityID`, `JobManagementValueImpactID`, dan `IsAuthorized`.
> Tipe seed yang relevan: `authority` (8 level, 035), `authority_unauthorized` (8 level, 037),
> `impact_unauthorized` (6 level, 036). **Tipe `cash` dan `impact` tidak ter-seed di
> migration 033–042** — namun terdefinisi di seeder legacy: `cash` = Jumlah Uang (5 level),
> `impact` = Dampak dengan wewenang keuangan (6 level) — lihat 8.6. Perlu di-seed ulang.

### 4.6 Kewenangan Aset (`calculateAssetAuthority`)

**Sumber data**: `job_management_assets` (1 record per org)
→ relasi `assetValue` dan `assetAuthority`.

```
asset_value_points      = MAP_LINEAR_8(level_asset)      // 1..8 → 1..8 poin
asset_authority_points  = MAP_DEFAULT(level_asset_authority)

score = value_points > 0 DAN authority_points > 0 ? value_points * authority_points : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go, `JobAsset` memakai `JobManagementValueAssetID` &
> `JobManagementValueAuthorityID`. Tipe seed yang relevan: `asset_authority` (6 level, 038).
> **Tipe `asset` tidak ter-seed** di migration — namun terdefinisi di seeder legacy
> (`Nilai Asset`, 8 level: 0–1 Jt → >1 M) — lihat 8.6. Perlu di-seed ulang.

### 4.7 Kendali Bawahan (`calculateSubordinateControl`)

**Sumber data**: `job_management_subordinate_controls` (1 record per org)
→ relasi `value`.

```
points = MAP_DEFAULT(level_value)
score  = points
```

> Mapping Go: `JobSubordinateControl.JobManagementValueID` → tipe `subordinate` (5 level, 034).

### 4.8 Ruang Lingkup Kerja (`calculateWorkScope`) — Scope × Frequency

**Sumber data**: `job_management_relationships` (1 record per org)
→ relasi `scopeValue` (JobManagementValueRelationship) dan `frequencyValue`.

```
scope_points     = MAP_DEFAULT(level_relationship)
frequency_points = MAP_LINEAR_5(level_frequency)

score = scope_points > 0 DAN frequency_points > 0 ? scope_points * frequency_points : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go, `JobRelationship` memakai
> `JobManagementValueRelationshipID` & `JobManagementValueFrequencyID`.
> **Tipe `relationship` dan `frequency` tidak ter-seed** di migration 033–042 — namun terdefinisi
> di seeder legacy: `relationship` = Lingkup Hubungan Kerja (5 level), `frequency` = Frekuensi
> Hubungan Kerja (5 level) — lihat 8.6. Perlu di-seed ulang.

### 4.9 Aktivitas Kerja (`calculateWorkActivity`)

**Sumber data**: `job_management_working_activities` (1 record per org) → relasi `value`.

```
points = MAP_DEFAULT(level_value)
score  = points
```

> Mapping Go: `JobWorkingActivity.JobManagementValueID` → tipe `activity` (5 level, 033).

### 4.10 Risiko Kerja (`calculateWorkRisk`) — Environment × Hazard

**Sumber data**: `job_management_working_risks` (1 record per org)
→ relasi `environmentValue` dan `hazardValue`.

```
environment_points = MAP_LINEAR_5(level_environment)
hazard_points      = MAP_LINEAR_5(level_hazard)

score = environment_points > 0 DAN hazard_points > 0 ? environment_points * hazard_points : 0
```

> ⚠️ **Catatan mapping ke Go**: di Go, `JobWorkingRisk` memakai
> `JobManagementValueEnvironmentID` & `JobManagementValueHazardID`. Migration 032 mengganti
> tipe `hazard` → `risk`. Frontend meminta `type=environment` dan `type=risk` — **tipe
> `environment` dan `risk` tidak ter-seed** di migration 033–042, namun terdefinisi di seeder
> legacy: `environment` = Lingkungan Kerja (5 level), `risk` = Resiko/Bahaya (5 level) — lihat 8.6.

---

## 5. Agregasi Total Skor

### 5.1 Base Score (tanpa komponen keuangan)

```
base_score =
    education_score
  + competency_aggregate   // technical_points * managerial_points * communication_points
  + asset_score
  + subordinate_score
  + relationship_score
  + activity_score
  + risk_score
```

Dimana:

```
competency_aggregate =
    competency_base_score        // = score komponen 4.3 (perkalian)
  + potential_score              // = score komponen 4.2
  + problem_solving_score        // = score komponen 4.4
```

> Catatan: meskipun komponen 4.2, 4.3, 4.4 dihitung terpisah (masing-masing dari record
> potency competencies dengan tipe berbeda), **ketiganya dijumlahkan** menjadi
> `competency_aggregate` sebelum ditambahkan ke `base_score`.

### 5.2 Total Skor (with / without financial)

```
if has_financial_authority:
    with_financial    = base_score + financial_score      // 4.5 (perkalian 3 faktor)
    without_financial = 0
else:
    with_financial    = 0
    without_financial = base_score + financial_score      // 4.5 (perkalian 2 faktor)
```

### 5.3 Sub-Components (`sub_component_points`) — dipakai untuk `is_complete`

```jsonc
{
  "education":                0,   // education_points
  "experience":               0,   // experience_points
  "potential":                0,   // potential_score
  "competency_technical":     0,   // technical_points
  "competency_managerial":    0,   // managerial_points
  "competency_communication": 0,   // communication_points
  "competency_total":         0,   // competency_aggregate
  "problem_solving":          0,   // problem_solving_score
  "financial_with_authority": 0,   // financial_score jika has_authority, else 0
  "financial_without_authority": 0, // financial_score jika !has_authority, else 0
  "asset_management":         0,   // asset_score
  "subordinate_control":      0,   // subordinate_score
  "work_scope":               0,   // relationship_score
  "work_activity":            0,   // activity_score
  "work_risk":                0    // risk_score
}
```

### 5.4 Status Kelengkapan (`is_complete`)

`isResultComplete` memeriksa **semua sub-komponen > 0** kecuali pasangan finansial
(hanya butuh salah satu `financial_with_authority` **atau** `financial_without_authority` > 0):

```
required = [education, experience, potential, competency_technical,
            competency_managerial, problem_solving, asset_management,
            subordinate_control, work_scope, work_activity, work_risk]
            // financial_with_authority / without → OR (salah satu > 0)

is_complete = semua required > 0 DAN (financial_with > 0 ATAU financial_without > 0)
```

---

## 6. Struktur Penyimpanan

### 6.1 Tabel `job_management_scores` (Go, migration 009)

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | UUID |
| `organization_id` | CHAR(36) UNIQUE NOT NULL | 1 skor per organisasi |
| `job_value_with_financial` | BIGINT UNSIGNED | total dengan wewenang keuangan |
| `job_value_without_financial` | BIGINT UNSIGNED | total tanpa wewenang keuangan |
| `has_financial_authority` | TINYINT(1) | flag wewenang keuangan |
| `components` | JSON NULL | rincian komponen (hasil `components`) |
| `sub_component_points` | JSON NULL | poin sub-komponen |
| `calculated_at` | TIMESTAMP NULL | waktu kalkulasi |
| `created_at` / `updated_at` | TIMESTAMP | timestamp |

> Catatan: tabel Go **tidak punya kolom `is_complete`/`completed_at`** yang ada di PHP legacy
> (`persistResults` mengisi `is_complete` dan `completed_at`). Perlu dipertimbangkan saat
> memindahkan mesin kalkulasi.

### 6.2 Model GORM Go (`JobScore`)

```go
type JobScore struct {
    ID                      uuid.UUID
    OrganizationID          *uuid.UUID  // uniqueIndex
    JobValueWithFinancial   uint64
    JobValueWithoutFinancial uint64
    HasFinancialAuthority   bool
    Components              *string  // JSON
    SubComponentPoints      *string  // JSON
    CalculatedAt            *time.Time
    CreatedAt, UpdatedAt    time.Time
}
```

### 6.3 DTO Go (`UpdateJobScoreRequest`)

```go
type UpdateJobScoreRequest struct {
    JobValueWithFinancial    *uint64  `json:"job_value_with_financial"`
    JobValueWithoutFinancial *uint64  `json:"job_value_without_financial"`
    HasFinancialAuthority    *bool    `json:"has_financial_authority"`
    Components               *string  `json:"components"`
    SubComponentPoints       *string  `json:"sub_component_points"`
}
```

---

## 7. Kondisi Implementasi Go Saat Ini (Update 05 Agu 2026)

### 7.1 Mesin kalkulasi ✅ SUDAH ADA — `calculator.go`

`internal/modules/jobmanagement/calculator.go` — port penuh dari `JobValueCalculator.php`:

- `NewCalculator(repo)` + `CalculateForOrganization(ctx, orgID)` — hitung 1 organisasi.
- `CalculateForOrganizationIDs(ctx, ids)` — batch (mirip `calculateForOrganizationIds` legacy),
  dipakai untuk re-kalkulasi massal.
- `loadCalcData` — load batch semua data section per org (edu-exp, potency, financials, assets,
  subordinates, relations, activities, risks) + semua `job_management_values` dalam 1 query.
- `calculateForSingleOrganization` — agregasi → `JobScoreResult` (semua komponen + total skor +
  `subComponentPoints` + `IsComplete`).
- Mapping level→poin sebagai konstanta Go (MAP_DEFAULT, MAP_EXTENDED, MAP_LINEAR_5,
  MAP_LINEAR_8, MAP_COMMUNICATION) + `mapPoints` fallback ke level tertinggi.
- **`mapCommunication = {1:1, 2:3, 3:6}`** — sesuai keputusan user.
- Perhitungan komponen: `calcEducationExperience`, `calcPotentials`, `calcCompetencies`,
  `calcProblemSolving`, `calcFinancialAuthority`, `calcAssetAuthority`,
  `calcSubordinateControl`, `calcWorkScope`, `calcWorkActivity`, `calcWorkRisk`.
- **Bobot (`weight`) tidak dibaca** dalam rumus — sesuai keputusan user (rumus legacy tetap).

### 7.2 Recalculate otomatis ✅ HOOK di tiap section

- `Service.RecalculateJobScore(ctx, orgID)` memanggil kalkulator → `scoreFromResult` →
  `repo.UpsertJobScore`. (service.go ~2220)
- `RecalculateJobScores(ctx, orgIDs)` — batch massal.
- **Hook otomatis**: semua `CreateJobXxx` / `UpdateJobXxx` / `DeleteJobXxx` untuk section yang
  memengaruhi skor (education-experiences, potency-competencies, financials, assets,
  subordinate-controls, relationships, working-activities, working-risks) memanggil
  `recalculateScore` setelah sukses — bukan hanya tombol recalculate.
- Endpoint: `PUT /api/v1/tenant/job-management/scores/org/:orgId` tetap ada untuk recalculate
  manual (body kosong = hitung ulang).

### 7.3 `is_complete` / `completed_at` ✅ DIPERSIST

`scoreFromResult` mengisi `IsComplete` + `CompletedAt` (hanya saat `is_complete = true`)
ke kolom migration 051 (`is_complete` TINYINT/SMALLINT + `completed_at` TIMESTAMP NULL).
`toJobScoreResponse` mengekspos keduanya ke API → dipakai badge di frontend.

### 7.4 Frontend `JobScoreSection.vue` (Update 05 Agu 2026)

- Menampilkan **breakdown komponen dari field `components`** (JSON nested dari DB), bukan lagi
  flat `sub_component_points` — per komponen menampilkan poin (level → poin) + skor, urut
  sesuai navigasi sub-form (Pendidikan & Pengalaman → Potensi → Kompetensi → Problem Solving →
  Keuangan → Aset → Bawahan → Hubungan Kerja → Aktivitas → Risiko).
- `competencies` memakai `base_score` agar tidak dobel hitung dengan potentials & problem_solving.
- Footer total with/without financial. Badge `is_complete` (Lengkap/Belum Lengkap).
- **Ringkasan skor sticky** di form utama (`JobScoreSummary.vue` di `JobManagementForm.vue`)
  agar selalu terlihat saat berpindah navigasi.

---

## 8. Gap & Temuan Penting (Untuk Implementasi Go)

### 8.1 Mesin kalkulasi ✅ SUDAH ADA di Go
`calculator.go` (05 Agu 2026) — port penuh dari legacy, lihat 7.1. Dipanggil otomatis saat
section disimpan (hook di Create/Update/Delete tiap section, lihat 7.2) **dan** di endpoint
recalculate.

### 8.2 Perbedaan penamaan tipe `job_management_values` (legacy vs Go)

| Komponen | Tipe legacy (PHP) | Tipe Go (migration) | Status seed Go |
|---|---|---|---|
| Pendidikan | level dari `educationGrade` (tipe `education`) | `EducationID` → tipe `education` | ✅ seed 049 (5 level: SMP, SMA, D3, S1, S2/S3) |
| Pengalaman | level dari `experienceGrade` (tipe `experience`) | `ExperienceID` → tipe `experience` | ✅ seed 049 (5 level: 0–2, 3–5, 6–8, 9–11, >12 th) |
| Potensi psikologi | `Psychological` | `kecerdasan`, `innovation_creativity`, `self_confidence`, `flexibility`, `tenacity`, `continuous_learning` | ✅ seed 041 |
| Kompetensi technical | `Technical` | 16 tipe (042) | ✅ seed 042 |
| Kompetensi managerial | `Managerial` | 6 tipe (042) | ✅ seed 042 |
| Communication | `Communicating & Influencing Skill` | `communicating_influencing_skill` | ✅ seed 039 |
| Problem solving | `Problem Solving & Decision Making` (descriptions `Environment`/`Chalenge`) | `thinking_environment`, `thinking_chalenge` | ✅ seed 040 |
| Kewenangan keuangan | `cash`, `authority`, `impact` | `authority`, `authority_unauthorized`, `impact_unauthorized` | ✅ `authority` (035), `authority_unauthorized` (037), `impact_unauthorized` (036), `cash` & `impact` (050) |
| Aset | `asset`, `asset_authority` | `asset_authority` | ✅ `asset_authority` (038), `asset` (050) |
| Bawahan | `subordinate` | `subordinate` | ✅ seed 034 |
| Ruang lingkup | `relationship`, `frequency` | `JobManagementValueRelationshipID` / `JobManagementValueFrequencyID` | ✅ seed 050 |
| Aktivitas | `activity` | `activity` | ✅ seed 033 |
| Risiko | `environment`, `hazard` | `environment`, `risk` | ✅ seed 050 |

> **Semua tipe sudah ter-seed** di migration Go: `education` & `experience` (049), dan
> `asset`, `cash`, `impact`, `environment`, `risk`, `relationship`, `frequency` (050) —
> level & deskripsi mengikuti seeder legacy (`JobManagementValuesTableSeeder.php`, lihat 8.6).
> Tidak ada lagi tipe referensi yang belum ter-seed.

### 8.3 Kecerdasan (Intelligence) sebagai potensi psikologi
Berdasarkan keputusan sebelumnya di halaman Kompetensi Potensi, **Kecerdasan termasuk potensi
psikologi** (tipe `kecerdasan`, 5 level) dan disimpan di `job_management_potency_competencies`
**tanpa `competency_id`** (tidak ada kompetensi master untuk Kecerdasan). Kalkulator Go harus
ikut menghitung level `kecerdasan` dalam rata-rata potensi psikologi.

### 8.4 Kolom `is_complete` / `completed_at` ✅ Sudah ada (migration 051)
Migration **051_job_management_scores_is_complete** (MySQL + Postgres, up/down) menambah
kolom `is_complete` (TINYINT(1)/SMALLINT, NOT NULL DEFAULT 0) dan `completed_at`
(TIMESTAMP NULL) pada `job_management_scores`. Kalkulator mengisi keduanya: `completed_at`
hanya diisi saat `is_complete = true` (sama seperti legacy `persistResults`).

### 8.5 Detail hubungan kerja (relationship details)
Section hubungan kerja kini punya detail banyak-per-relationship (`job_management_relationship_details`,
migration 048). Komponen 4.8 hanya memakai `scope × frequency` dari record relationship utama —
**detail tidak memengaruhi skor** di legacy. Perlu konfirmasi apakah detail tetap tidak dihitung.

### 8.6 Struktur data legacy: `type` vs `code` (dari `JobManagementValuesTableSeeder.php`)
Seeder legacy **tidak** menyimpan slug tipe secara langsung di kolom `type` — ia memakai
**dua kolom berbeda**:

- **`type`** = **nama grup** besar (label section), contoh: `Psychological`, `Technical`,
  `Managerial`, `Problem Solving & Decision Making`, `Communicating & Influencing Skill`,
  `Pendidikan`, `Pengalaman Kerja`, `Nilai Asset`, `Jumlah Uang`, `Lingkungan Kerja`, dst.
- **`code`** = **slug** per item (dari `$codeMap`, atau di-generate dari deskripsi bila tidak ada
  di `$codeMap`), contoh: `kecerdasan`, `innovation_creativity`, `thinking_environment`, dst.
- `level`, `descriptions`, `note`, `sort` — `sort` diisi 1..N per grup untuk urutan tampil
  (hanya untuk grup kompetensi).

> **Catatan asal slug**: untuk grup di luar `$codeMap` (Pendidikan, Jumlah Uang, Nilai Asset,
> dst.), `code` di-generate dari teks deskripsi (mis. `'0 - 500 Jt'` → `'0_500_jt'`) — seeder
> **tidak** menghasilkan slug `education`/`cash`/`asset`. Pemetaan ke nama tipe Go
> (`education`, `cash`, `asset`, `impact`, `environment`, `risk`, `relationship`, `frequency`)
> adalah **keputusan penamaan sisi Go**, bukan turunan langsung dari seeder.
>
> **Dikecualikan dari skor**: grup `Jurusan` & `Bidang Pekerjaan` di seeder hanya berisi level 1
> (opsi pilihan, bukan level skor) — di Go dipindah ke tabel pivot
> `job_management_majors` / `job_management_job_family` dan **tidak ikut dihitung** dalam
> perhitungan skor.

**Konsekuensi untuk Go**: di migration Go 033–042, slug `code` legacy **diangkat menjadi kolom
`type`** (contoh: `type='kecerdasan'`, `type='thinking_environment'`). Sebaliknya, nama grup
legacy (`Psychological`, `Technical`, dst.) **tidak ada** di data Go. Artinya kalkulator Go
**tidak bisa** memfilter `type == 'Psychological'` seperti PHP — harus memakai **daftar tetap
slug tipe** per kelompok (lihat tabel di bawah).

#### Pemetaan `$codeMap` legacy (deskripsi kompetensi → slug)

| Grup legacy (`type`) | Deskripsi → code (slug) |
|---|---|
| `Psychological` | Kecerdasan→`kecerdasan`, Innovation & Creativity→`innovation_creativity`, Self Confidence→`self_confidence`, Flexibility→`flexibility`, Tenacity→`tenacity`, Continuous Learning→`continuous_learning` |
| `Technical` | 16 kompetensi → 16 slug (lengkap di migration 042) |
| `Managerial` | Integrity, Achievement Orientation, Building Partnership, Planning & Organizing, Leadership, Developing Others → 6 slug |
| `Problem Solving & Decision Making` | Thinking Environment→`thinking_environment`, Thinking Chalenge→`thinking_chalenge` |
| `Communicating & Influencing Skill` | → `communicating_influencing_skill` |

#### Level per tipe (konfirmasi dari seeder)

| Grup legacy | Level | Deskripsi level (ringkas) | Tipe Go terkait |
|---|---|---|---|
| `Pendidikan` | 1–5 | SMA, D3, S1, S2, S3 (keputusan user 04 Agu: SMP, SMA, D3, S1, S2/S3) | `education` ✅ seed 049 |
| `Pengalaman Kerja` | 1–5 | 0–2 th, 3–5 th, 6–8 th, 9–11 th, >12 th | `experience` ✅ seed 049 |
| `Psychological` — Kecerdasan | 1–5 | Kurang, Cukup, Rata-rata, Diatas rata-rata, Istimewa | `kecerdasan` (seed 041) |
| `Psychological` — lainnya (5) | 1–8 | Primary → Unique Authority | `innovation_creativity`, `self_confidence`, `flexibility`, `tenacity`, `continuous_learning` (seed 041) |
| `Technical` (16) | 1–8 | Primary → Unique Authority | 16 slug (seed 042) |
| `Managerial` (6) | 1–5 | Task, Supervisory, Managerial, Diverse Managerial, Total Managerial | 6 slug (seed 042) |
| `Problem Solving` — Env | 1–8 | Berulang-ulang → Didefinisikan Secara Abstrak | `thinking_environment` (seed 040) |
| `Problem Solving` — Challenge | 1–5 | Berulang-ulang → Belum dipetakan | `thinking_chalenge` (seed 040) |
| `Communicating & Influencing Skill` | 1–3 | Berkomunikasi, Alasan, Perubahan Perilaku | `communicating_influencing_skill` (seed 039) |
| `Jumlah Uang` | 1–5 | 0–500 Jt → >10 M | `cash` ✅ seed 050 |
| `Wewenang` (Memiliki) | 1–8 | instruksi langsung → panduan sangat luas | `authority` (seed 035) |
| `Memiliki Wewenang` (Tidak) | 1–8 | Dikendalikan ketat → Dipandu strategis | `authority_unauthorized` (seed 037) |
| `Dampak … (Memiliki Wewenang Keuangan)` | 1–6 | insidentil → multi tim/strategis | `impact` ✅ seed 050 |
| `Dampak … (Tidak Memiliki Wewenang Keuangan)` | 1–6 | PAnciliary → Berpengaruh | `impact_unauthorized` (seed 036) |
| `Nilai Asset` | 1–8 | 0–1 Jt → >1 M | `asset` ✅ seed 050 |
| `Wewenang Asset` | 1–6 | Tidak ada Aset → Menyetujui Aset | `asset_authority` (seed 038) |
| `Total Bawahan` | 1–5 | Very Small → Very Large | `subordinate` (seed 034) |
| `Aktifitas Fisik` | 1–5 | Banyak duduk → pengawasan | `activity` (seed 033) |
| `Lingkungan Kerja` | 1–5 | Tenang → sangat menekan | `environment` ✅ seed 050 |
| `Resiko/Bahaya` | 1–5 | Risiko minimum → bahaya mematikan | `risk` ✅ seed 050 |
| `Lingkup Hubungan Kerja` | 1–5 | Unit Kerja → Internasional | `relationship` ✅ seed 050 |
| `Frekuensi Hubungan Kerja` | 1–5 | Sesekali → Sangat Sering | `frequency` ✅ seed 050 |

> **Status terbaru (04 Agu 2026)**: **semua tipe referensi sudah di-seed** di migration Go:
> - **049_seed_job_value_education_experience** — `education` (SMP, SMA, D3, S1, S2/S3) &
>   `experience` (0–2, 3–5, 6–8, 9–11, >12 Tahun), 5 level masing-masing.
> - **050_seed_job_value_cash_impact_environment_risk_relationship_frequency_asset** —
>   `cash` (5), `impact` (6), `environment` (5), `risk` (5), `relationship` (5), `frequency` (5),
>   `asset` (8) — level & deskripsi dari seeder legacy.
> Keduanya tersedia untuk MySQL & Postgres (up/down), idempoten by UUID v4.

### 8.7 Gap struktur tabel & relasi model (hasil verifikasi tambahan)

#### 8.7.1 Relasi GORM untuk kalkulasi **belum tersedia** di sebagian besar model

Kalkulator legacy membaca level via **relasi Eloquent** (`->with(['cashValue', 'authorityValue',
'impactValue'])`, `->with('value')`, dll.). Di Go, verifikasi pada `model.go` & `repository.go`
menunjukkan:

| Model | Relasi GORM ke `job_management_values` | Status |
|---|---|---|
| `JobEducationExperience` | ✅ `Education` (type=education), `Experience` (type=experience), `Majors`, `JobFamilies` | Ada (repository Preload baris 425/446) |
| `JobPotencyCompetency` | ❌ **tidak ada** relasi `JobManagementValue` / `Competency` / `CompetencyValue` | Tidak ada preload di repository |
| `JobFinancial` | ❌ tidak ada relasi `CashValue`/`AuthorityValue`/`ImpactValue` | Tidak ada preload |
| `JobAsset` | ❌ tidak ada relasi `AssetValue`/`AssetAuthority` | Tidak ada preload |
| `JobWorkingRisk` | ❌ tidak ada relasi `EnvironmentValue`/`HazardValue` | Tidak ada preload |
| `JobRelationship` | ❌ tidak ada relasi `ScopeValue`/`FrequencyValue` | Tidak ada preload |
| `JobSubordinateControl` | ❌ tidak ada relasi `Value` | Tidak ada preload |
| `JobWorkingActivity` | ❌ tidak ada relasi `Value` | Tidak ada preload |

> **Implikasi**: kalkulator Go tidak bisa langsung memakai `Preload` seperti PHP. Perlu salah satu:
> (a) menambah relasi GORM + `Preload` di model/repository, atau
> (b) kalkulator melakukan query join manual ke `job_management_values` (via `JobManagementValueID`
> dsb.) untuk mengambil `level` & `type` per record.

#### 8.7.2 Kolom `is_complete` / `completed_at` ✅ Sudah ada (migration 051)

PHP legacy mengisi `is_complete` + `completed_at` di `persistResults`. Di Go, gap ini ditutup
oleh migration **051_job_management_scores_is_complete** (MySQL `TINYINT(1) NOT NULL DEFAULT 0`
+ `TIMESTAMP NULL`; Postgres `SMALLINT NOT NULL DEFAULT 0` + `TIMESTAMP`) dan diisi oleh
`scoreFromResult` (service) dari `JobScoreResult.IsComplete`.

#### 8.7.3 Dua sumber level berbeda: `job_management_values` vs `competency_values`

- `job_management_values` — dipakai semua section job management (tipe slug: `kecerdasan`, dst.).
- `competency_values` — tabel nilai kompetensi module competency (`Type`, `Name`, `Slug`,
  `Level`, `Code`).

Di legacy, `collectLevelsByType` punya **fallback ke `competency.field`** (`Technical Competency` /
`Manajerial`) + `competencyValue.level` saat record tidak bertipe legacy. Di Go,
`JobPotencyCompetency` menyimpan **keduanya**: `job_management_value_id` **dan** `competency_id`.
Kalkulator Go harus mendefinisikan urutan prioritas level. **Usulan prioritas eksplisit**:
1. Pakai `job_management_value_id` → `job_management_values.level` bila ada;
2. Fallback ke `competency_id` → `competencies.field` (Technical/Manajerial) +
   `competency_values.level` bila `job_management_value_id` kosong.
Kasus nyata: record Kecerdasan disimpan **tanpa `competency_id`** (lihat 8.3), jadi
`job_management_value_id` pasti ada — prioritas #1 menanganinya dengan benar.

#### 8.7.4 Tidak ada UNIQUE `(type, level)` di `job_management_values`

Tabel tidak punya constraint unik pada kombinasi `(type, level)` (hanya index `idx_jmv_ref`).
Seeder idempotensi di 033–042 memakai `WHERE NOT EXISTS` / UUID tetap — **bukan** UNIQUE key.
Relevan saat membuat migration seed ulang tipe yang kurang: perlu pola idempoten yang sama.

#### 8.7.5 Penamaan flag: `IsAuthorized` (financial) vs `HasFinancialAuthority` (score)

`job_management_financials.is_authorized` (input) → kalkulator menyalinnya ke
`job_management_scores.has_financial_authority` (hasil). Konsisten, tapi kalkulator wajib
memetakan keduanya.

### 8.8 Ringkasan Gap (checklist untuk implementasi) — UPDATE 05 Agu 2026

| # | Gap | Dampak | Status |
|---|---|---|---|
| 1 | ~~Mesin kalkulasi belum ada di Go~~ | Skor tidak pernah terhitung | ✅ **Selesai** — `calculator.go` (7.1) |
| 2 | ~~Relasi GORM ke `job_management_values`~~ | Kalkulator tidak bisa preload level | ✅ **Selesai** — kalkulator memakai join manual di `loadCalcData` (tanpa mengubah model) |
| 3 | ~~Tipe data referensi belum ter-seed~~ | — | ✅ **Sudah di-seed semua** (049, 050, 052/053, 054) |
| 4 | Nama tipe Go = slug, bukan grup legacy | Filter `type == 'Psychological'` tidak jalan | ✅ **Selesai** — daftar tetap slug + kolom `type_group`/`description_group` (052/053) + endpoint tree (lihat 8.9) |
| 5 | ~~Kolom `is_complete` / `completed_at` tidak ada~~ | — | ✅ **Selesai** — migration **051** + persist `scoreFromResult` |
| 6 | Kecerdasan disimpan tanpa `competency_id` | Harus ikut rata-rata potensi | ✅ **Selesai** — `calcPotentials` menghitung slug psikologi termasuk `kecerdasan` |
| 7 | `competency_values` vs `job_management_values` dua sumber level | Prioritas level ambigu | ✅ **Selesai** — prioritas `job_management_value_id`, fallback ke competencies |

### 8.9 Kolom `type_group` & `description_group` (migration 052/053) + cluster mapping (054) + tree endpoint

- **052_job_management_values_type_group** (MySQL + Postgres, up/down): tambah kolom
  `type_group` (VARCHAR) & `description_group` (VARCHAR) pada `job_management_values`.
- **053_seed_job_value_type_group**: isi `type_group` & `description_group` untuk semua tipe
  (contoh: `type_group=psychological`, `type=kecerdasan` → `description_group=Kecerdasan`;
  tipe psikologi lain → label grupnya; technical/managerial → label kompetensi, dst.).
  Idempotent by UUID, tersedia MySQL & Postgres.
- **054_create_job_management_value_clusters**: tabel `job_management_value_clusters` untuk
  memetakan tipe technical/managerial ke cluster kompetensi (setting di halaman mapping Job
  Value — submenu Technical & Managerial memakai endpoint cluster di bawah).
- **Endpoint tree**: `GET /api/v1/tenant/job-management/values/tree` — mengembalikan hierarki
  `type_group → daftar tipe (label = description_group) → options per tipe (level + deskripsi)`
  dengan urutan grup tetap (education, experience, psychological, technical, managerial,
  communication, problem_solving, financial, asset, subordinate, activity, environment, risk,
  relationship, frequency). Dipakai form potensi (filter type_group → multi-select tipe →
  tabel isian per tipe).
- **Endpoint cluster**: `GET/PUT /api/v1/tenant/job-management/values/clusters/:type` —
  membaca/menyimpan mapping cluster (`technical`/`managerial`) → daftar kompetensi, dipakai
  card Kompetensi Teknis/Manajerial (filter dari table `competencies` dengan cluster selain
  Core/Manajerial, option dari `job_management_values` type `technical`/`managerial`).

---

## 9. Rencana Implementasi yang Disarankan — ✅ SEMUA SELESAI (05 Agu 2026)

1. ✅ **`internal/modules/jobmanagement/calculator.go`** — port dari `JobValueCalculator.php`:
   - `CalculateForOrganization(ctx, orgID)` + `CalculateForOrganizationIDs(ctx, ids)` (batch).
   - Mapping level→poin (5 tabel) sebagai konstanta Go.
   - Helper pengelompokan tipe psikologi / technical / managerial (daftar slug tetap).
2. ✅ **Akses data untuk kalkulator** — `loadCalcData` melakukan query join manual ke
   `job_management_values` (tanpa mengubah model GORM).
3. ✅ **Integrasi**:
   - Hook recalculate otomatis di Create/Update/Delete tiap section yang memengaruhi skor
     (`recalculateScore` → `RecalculateJobScore` → `UpsertJobScore`).
   - Endpoint `PUT /scores/org/:orgId` tetap ada untuk recalculate manual (body kosong).
   - Hasil disimpan ke `JobScore` (`components`, `sub_component_points`, `calculated_at`,
     `is_complete`, `completed_at`).
   - Batch `RecalculateJobScores(ctx, orgIDs)` untuk re-kalkulasi massal.
4. ✅ **Seed data referensi** — migration 049, 050, 052/053 (type_group/description_group),
   054 (clusters) untuk MySQL & Postgres.
5. ✅ **Migrasi `is_complete` & `completed_at`** — migration 051 (MySQL + Postgres) + persist
   dari `JobScoreResult.IsComplete` (completed_at hanya saat is_complete).
6. ✅ **Verifikasi data** — tipe referensi tersedia di DB tenant; re-kalkulasi massal dijalankan
   untuk organisasi yang punya data section.

---

## 10. Contoh Perhitungan (Ilustrasi)

Misalkan sebuah jabatan memiliki:

| Komponen | Level | Poin |
|---|---|---|
| Pendidikan | 4 | 10 |
| Pengalaman | 3 | 6 |
| Potensi psikologi (avg) | 4 (ceil) | 10 |
| Technical (avg) | 5 | 15 (MAP_EXTENDED) |
| Managerial (avg) | 3 | 6 |
| Communication | 2 | 3 |
| Problem solving: env 4 / challenge 3 | — | 10 × 6 = 60 |
| Aset (linear8 3 × default 4) | — | 3 × 10 = 30 |
| Bawahan | 4 | 10 |
| Scope × Frequency (4 × 3) | — | 10 × 3 = 30 |
| Aktivitas | 3 | 6 |
| Risiko (env 3 × hazard 3) | — | 3 × 3 = 9 |
| Keuangan (has authority: cash 4 × auth 5 × impact 4) | — | 10 × 15 × 10 = 1500 |

```
competency_aggregate = (15 * 6 * 3) + 10 + 60 = 270 + 70 = 340
base_score = 60(edu) + 340 + 30 + 10 + 30 + 6 + 9 = 485
with_financial    = 485 + 1500 = 1985
without_financial = 0   (karena has_financial_authority = true)
```

---

*Dokumen analisa ini sudah diimplementasikan penuh di Go (05 Agu 2026). Lihat juga
`docs/project-completion-dashboard.md` untuk changelog fitur dan `docs/openapi-report.md`
untuk daftar endpoint lengkap.*
