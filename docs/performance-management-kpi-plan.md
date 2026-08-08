Berikut versi yang sudah disesuaikan dengan standar **UUID sebagai Primary Key** di seluruh tabel. Saya juga menyesuaikan foreign key agar semuanya menggunakan `uuid` sehingga konsisten dengan arsitektur modern (Laravel + PostgreSQL).

````markdown
# KPI Module Development Plan

> 🔗 **Arsip dokumentasi:** [`docs/README.md`](../README.md) — dokumen ini adalah referensi historis dari modul yang sudah selesai diimplementasikan.


## Objective

Membangun modul Performance Management (KPI) yang:

- Organization = Position = 1 Employee
- KPI mengikuti Organization
- Tidak memerlukan Assignment KPI
- Mendukung Balanced Scorecard
- Mendukung Review Tahunan maupun Semester
- Mendukung Progress Monitoring
- Mendukung Evidence Attachment
- Menggunakan Snapshot KPI saat evaluasi
- Seluruh Primary Key menggunakan UUID

---

# Phase 1 - Database Enhancement

## 1. Update performance_templates

### Add

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| organization_id | uuid FK | Organization pemilik template |
| period_id | uuid FK | Performance Period |
| name | varchar(255) | Nama Template |
| description | text | Deskripsi |
| status | smallint | Draft / Active / Inactive |
| effective_date | date | Berlaku mulai |
| expired_date | date nullable | Berlaku sampai |
| created_at | timestamp | |
| updated_at | timestamp | |
| deleted_at | timestamp nullable | Soft Delete |

---

## 2. Update performance_indicators

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| template_id | uuid FK | Template KPI |
| perspective_id | uuid FK | Perspective |
| code | varchar(50) | Kode KPI |
| name | varchar(255) | Nama KPI |
| description | text | Deskripsi KPI |
| weight | decimal(5,2) | Bobot |
| target_type | varchar(30) | Number, Currency, Percentage, Duration, Boolean |
| target_value | decimal(18,2) | Target |
| unit | varchar(50) | %, Hari, Jam, Rupiah, Unit |
| formula_type | varchar(30) | Manual, Higher Better, Lower Better, Range |
| minimum_score | decimal(5,2) | Minimum Score |
| maximum_score | decimal(5,2) | Maximum Score |
| sort_order | integer | Urutan |
| is_required | boolean | KPI wajib |
| created_at | timestamp | |
| updated_at | timestamp | |
| deleted_at | timestamp nullable | Soft Delete |

---

## 3. Update performance_evaluations

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| employee_id | uuid FK | Employee |
| organization_id | uuid FK | Organization saat evaluasi |
| period_id | uuid FK | Performance Period |
| status | varchar(20) | Draft, Submitted, Approved, Completed |
| submitted_at | timestamp nullable | Submit |
| approved_at | timestamp nullable | Approval |
| final_score | decimal(5,2) | Nilai Akhir |
| rating_id | uuid FK | Rating |
| created_at | timestamp | |
| updated_at | timestamp | |
| deleted_at | timestamp nullable | Soft Delete |

---

## 4. Update performance_evaluation_details

Snapshot KPI saat evaluasi.

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| evaluation_id | uuid FK | Evaluation |
| indicator_id | uuid FK | KPI Indicator |
| indicator_name | varchar(255) | Snapshot Nama KPI |
| weight | decimal(5,2) | Snapshot Bobot |
| target | decimal(18,2) | Snapshot Target |
| actual | decimal(18,2) | Nilai Aktual |
| achievement | decimal(5,2) | Persentase Pencapaian |
| score | decimal(5,2) | Nilai KPI |
| remarks | text | Catatan Reviewer |
| created_at | timestamp | |
| updated_at | timestamp | |

---

# Phase 2 - New Tables

## performance_progress

Monitoring pencapaian KPI.

| Field | Type |
|--------|------|
| id | uuid PK |
| evaluation_detail_id | uuid FK |
| progress_date | date |
| actual_value | decimal(18,2) |
| achievement | decimal(5,2) |
| notes | text |
| created_by | uuid FK |
| created_at | timestamp |
| updated_at | timestamp |

---

## performance_comments

Komentar antara Employee dan Reviewer.

| Field | Type |
|--------|------|
| id | uuid PK |
| evaluation_id | uuid FK |
| employee_id | uuid FK |
| comment | text |
| created_by | uuid FK |
| created_at | timestamp |
| updated_at | timestamp |

---

## performance_attachments

Evidence KPI.

| Field | Type |
|--------|------|
| id | uuid PK |
| evaluation_detail_id | uuid FK |
| file | varchar(255) |
| description | text |
| uploaded_by | uuid FK |
| created_at | timestamp |
| updated_at | timestamp |

---

## performance_ratings

Master Rating.

| Field | Type |
|--------|------|
| id | uuid PK |
| name | varchar(100) |
| min_score | decimal(5,2) |
| max_score | decimal(5,2) |
| color | varchar(20) |
| description | text |
| created_at | timestamp |
| updated_at | timestamp |

Contoh

| Rating | Score |
|---------|------:|
| Outstanding | 95 - 100 |
| Excellent | 85 - 94 |
| Good | 75 - 84 |
| Fair | 60 - 74 |
| Poor | < 60 |

---

## performance_indicator_formulas

Master Formula KPI.

| Field | Type |
|--------|------|
| id | uuid PK |
| name | varchar(100) |
| code | varchar(50) |
| formula_type | varchar(30) |
| expression | text |
| description | text |
| created_at | timestamp |
| updated_at | timestamp |

Contoh Formula

- Higher Better
- Lower Better
- Manual Score
- Percentage
- Boolean
- Range

---

## performance_logs

Audit Trail.

| Field | Type |
|--------|------|
| id | uuid PK |
| evaluation_id | uuid FK |
| action | varchar(100) |
| before | jsonb |
| after | jsonb |
| created_by | uuid FK |
| created_at | timestamp |

---

# Phase 3 - Business Process

## KPI Setup

```text
Performance Period
        │
        ▼
Organization
        │
        ▼
Performance Template
        │
        ▼
Performance Indicator
```

---

## KPI Evaluation

```text
Employee
      │
      ▼
Create Evaluation
      │
      ▼
Snapshot KPI
      │
      ▼
Input Actual
      │
      ▼
Calculate Achievement
      │
      ▼
Calculate Score
      │
      ▼
Performance Rating
      │
      ▼
Completed
```

---

## Progress Monitoring

```text
Evaluation
      │
      ▼
Monthly Progress
      │
      ▼
Dashboard
```

---

## Review Workflow

```text
Draft

↓

Employee Submit

↓

Manager Review

↓

Revision (Optional)

↓

Approved

↓

Completed
```

---

# Phase 4 - Dashboard

## Employee

- KPI Progress
- Achievement
- Final Score
- Rating
- Evidence

---

## Manager

- Team KPI
- KPI Status
- KPI Progress
- Need Review
- Overdue Review

---

## HR

- KPI Completion
- Rating Distribution
- Average Achievement
- Top Performer
- Bottom Performer
- Achievement per Organization
- Trend per Period

---
# Phase 5 - Performance Scoring Configuration

## Objective

Mendukung konfigurasi metode perhitungan nilai Performance pada setiap Organization.

Setiap Organization dapat memiliki:

- Komponen penilaian yang berbeda.
- Bobot penilaian yang berbeda.
- Formula penilaian yang berbeda.

Dengan pendekatan ini setiap jabatan dapat menggunakan metode penilaian yang sesuai dengan tanggung jawabnya.

Contoh:

| Organization | KPI | Work Program | Subordinate KPI |
|--------------|----:|-------------:|----------------:|
| Director | 30% | 30% | 40% |
| Manager | 40% | 40% | 20% |
| Supervisor | 50% | 50% | 0% |
| Staff | 60% | 40% | 0% |

Total bobot setiap Organization harus selalu **100%**.

---

## Database Enhancement

### performance_components

Master komponen penilaian.

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| code | varchar(50) | Unique Code |
| name | varchar(100) | Component Name |
| description | text | Description |
| sort_order | integer | Display Order |
| is_active | boolean | Active Status |
| created_at | timestamp | |
| updated_at | timestamp | |
| deleted_at | timestamp nullable | Soft Delete |

---

### performance_organization_components

Konfigurasi bobot komponen untuk setiap Organization.

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| organization_id | uuid FK | Organization |
| component_id | uuid FK | Performance Component |
| weight | decimal(5,2) | Weight (%) |
| is_enabled | boolean | Active Component |
| sort_order | integer | Display Order |
| created_at | timestamp | |
| updated_at | timestamp | |

---

### performance_evaluation_components

Snapshot hasil perhitungan setiap komponen pada saat Evaluation.

| Field | Type | Description |
|--------|------|-------------|
| id | uuid PK | Primary Key |
| evaluation_id | uuid FK | Performance Evaluation |
| component_id | uuid FK | Performance Component |
| component_name | varchar(100) | Snapshot Component |
| score | decimal(5,2) | Component Score |
| weight | decimal(5,2) | Snapshot Weight |
| final_score | decimal(5,2) | Weighted Score |
| calculated_at | timestamp | Calculation Time |
| created_at | timestamp | |
| updated_at | timestamp | |

---

## Default Components

| Code | Component |
|------|-----------|
| KPI | KPI |
| WORK_PROGRAM | Work Program |
| SUBORDINATE | Subordinate KPI |

---

## Organization Configuration Example

### Director

| Component | Weight |
|-----------|-------:|
| KPI | 30% |
| Work Program | 30% |
| Subordinate KPI | 40% |

### Manager

| Component | Weight |
|-----------|-------:|
| KPI | 40% |
| Work Program | 40% |
| Subordinate KPI | 20% |

### Staff

| Component | Weight |
|-----------|-------:|
| KPI | 60% |
| Work Program | 40% |

---

## Score Calculation

```
Final Score =
Σ(Component Score × Weight)
```

Contoh

| Component | Score | Weight | Result |
|-----------|------:|-------:|-------:|
| KPI | 90 | 40% | 36 |
| Work Program | 85 | 40% | 34 |
| Subordinate KPI | 95 | 20% | 19 |

```
Final Score = 89
```

---

## Subordinate KPI

Subordinate KPI dihitung otomatis dari rata-rata nilai akhir seluruh bawahan langsung (direct subordinate) pada periode yang sama.

Contoh

```text
General Manager
        │
        ├── Manager A (90)
        ├── Manager B (80)
        └── Manager C (85)
```

```
Subordinate KPI

=

(90 + 80 + 85) / 3

=

85
```

Perhitungan dilakukan berdasarkan struktur Organization (`parent_id`) sehingga tidak memerlukan relasi bawahan tambahan.

Organization yang tidak memiliki bawahan dapat menonaktifkan komponen **Subordinate KPI** atau memberikan bobot **0%**.

---

## Validation Rules

- Total bobot setiap Organization harus = 100%.
- Satu komponen hanya boleh muncul satu kali pada setiap Organization.
- Organization tanpa bawahan tidak wajib menggunakan komponen Subordinate KPI.
- Perubahan bobot hanya berlaku untuk Evaluation yang dibuat setelah konfigurasi diubah.
- Bobot disimpan sebagai snapshot pada saat Evaluation.

---

## Business Flow

```text
Organization

↓

Performance Component Configuration

↓

Performance Evaluation

↓

Calculate Component Score

↓

Performance Scoring Engine

↓

Final Score

↓

Performance Rating
```


---

---

# Phase 6 - Master Data & Seeder

Beberapa tabel merupakan **master data** yang direkomendasikan menggunakan Seeder agar implementasi HRIS lebih cepat, konsisten, dan mudah dipelihara.


# Phase 7 - Future Enhancement

- Mid Year Review
- Calibration
- KPI Versioning
- Reminder Notification
- Multi Level Approval
- Bonus Integration
- Performance Improvement Plan (PIP)
- Succession Planning Integration
- Career Path Integration
- Talent Management Integration

## Seeder Overview

| Table | Seeder | Required | Description |
|--------|:------:|:--------:|-------------|
| performance_periods | ❌ | No | Data transaksi, dibuat setiap periode |
| performance_perspectives | ✅ | Yes | Master Balanced Scorecard Perspective |
| performance_templates | ❌ | No | Berbeda setiap Organization |
| performance_indicators | ❌ | No | Dibuat sesuai kebutuhan KPI |
| performance_evaluations | ❌ | No | Data transaksi |
| performance_evaluation_details | ❌ | No | Data transaksi |
| performance_progress | ❌ | No | Data transaksi |
| performance_comments | ❌ | No | Data transaksi |
| performance_attachments | ❌ | No | Data transaksi |
| performance_ratings | ✅ | Yes | Master Rating Performance |
| performance_indicator_formulas | ✅ | Yes | Master Formula KPI |
| performance_logs | ❌ | No | Audit Trail |

---

## Performance Perspectives Seeder

Digunakan apabila perusahaan menerapkan Balanced Scorecard.

| Code | Name | Sort Order |
|------|------|-----------:|
| FIN | Financial | 1 |
| CUS | Customer | 2 |
| INT | Internal Process | 3 |
| LRN | Learning & Growth | 4 |

---

## Performance Ratings Seeder

Master rating berdasarkan nilai akhir KPI.

| Code | Name | Min Score | Max Score | Color |
|------|------|----------:|----------:|--------|
| OUT | Outstanding | 95.00 | 100.00 | success |
| EXC | Excellent | 85.00 | 94.99 | primary |
| GOO | Good | 75.00 | 84.99 | info |
| FAI | Fair | 60.00 | 74.99 | warning |
| POO | Poor | 0.00 | 59.99 | danger |

---

## Performance Indicator Formulas Seeder

Formula standar yang dapat digunakan pada setiap KPI.

| Code | Name | Description |
|------|------|-------------|
| MANUAL | Manual Score | Nilai diinput manual oleh reviewer |
| HIGHER | Higher Better | Semakin tinggi nilai semakin baik |
| LOWER | Lower Better | Semakin rendah nilai semakin baik |
| RANGE | Range Score | Nilai berdasarkan rentang tertentu |
| BOOLEAN | Boolean | Ya/Tidak |
| PERCENTAGE | Percentage | Nilai berupa persentase |

---

# Optional Master Data

Untuk fleksibilitas jangka panjang, beberapa field yang saat ini menggunakan Enum dapat dipisahkan menjadi Master Data.

## performance_target_types

| Code | Name |
|------|------|
| NUMBER | Number |
| PERCENTAGE | Percentage |
| CURRENCY | Currency |
| BOOLEAN | Boolean |
| DURATION | Duration |
| SCORE | Score |

---

## performance_units

| Code | Name |
|------|------|
| PERCENT | Percent (%) |
| IDR | Indonesian Rupiah |
| USD | US Dollar |
| DAY | Day |
| HOUR | Hour |
| MINUTE | Minute |
| PERSON | Person |
| UNIT | Unit |
| CASE | Case |
| ITEM | Item |

---

## performance_statuses

Apabila status tidak menggunakan Enum.

| Code | Name |
|------|------|
| DRAFT | Draft |
| SUBMITTED | Submitted |
| REVIEW | Under Review |
| APPROVED | Approved |
| COMPLETED | Completed |
| CANCELLED | Cancelled |

---

## performance_period_types

Digunakan apabila perusahaan mendukung berbagai jenis periode penilaian.

| Code | Name |
|------|------|
| MONTHLY | Monthly |
| QUARTERLY | Quarterly |
| SEMESTER | Semester |
| YEARLY | Yearly |

---

# Seeder Structure

struktur dan penamaan ikuti pola yang sudah ada

---

# Seeder Priority

## Mandatory

- PerformancePerspectiveSeeder
- PerformanceRatingSeeder
- PerformanceIndicatorFormulaSeeder

## Recommended

- PerformanceUnitSeeder

## Optional

- PerformanceTargetTypeSeeder
- PerformanceStatusSeeder
- PerformancePeriodTypeSeeder

---

# Design Principles

- Seluruh Primary Key menggunakan UUID.
- Seluruh Foreign Key menggunakan UUID.
- Organization = Position = 1 Employee.
- KPI mengikuti Organization.
- Tidak menggunakan Assignment KPI.
- Evaluation menggunakan Snapshot KPI.
- Histori KPI tidak berubah walaupun Template berubah.
- Mendukung Soft Delete.
- Mendukung Audit Trail.
- Mendukung Progress Monitoring.
- Master Data menggunakan Seeder agar implementasi lebih cepat dan konsisten.
- Siap diintegrasikan dengan Competency Management, Career Path, Succession Planning, Talent Management, Bonus Management, dan Performance Improvement Plan (PIP).

---

# Implementation Status

| Phase | Status | Completion Date | Notes |
|-------|--------|-----------------|-------|
| Phase 1-4 | ✅ Completed | 2026-08-06 | Detail frontend KPI diarsipkan ke `docs/archive/frontend-performance-kpi-plan.md` |
| Phase 5 - Scoring Configuration | ✅ Completed | 2026-08-07 | Backend: models, DTO, repository, service (scoring engine), handler, routes, migration |
| Phase 6 - Seeder | ✅ Completed | 2026-08-07 | Default components (KPI, Work Program, Subordinate KPI) di-seed via `module.go` |
| Phase 7 - Future Enhancement | 🔶 Partially completed | 2026-08-08 | Multi Level Approval selesai (lihat Phase 8); sisanya (Mid Year Review, Calibration, KPI Versioning, dst.) masih pending |
| Phase 8 - Approval Integration & Two-Phase Target/Realization | ✅ Completed | 2026-08-08 | Lihat detail di bawah |
| Phase 9 - Program Component, Bottom-Up Subordinate Scoring & OKR Self-Assessment | ✅ Completed | 2026-08-08 | Lihat detail di bawah |

## Phase 5 Deliverables

| File | Status | Description |
|------|--------|-------------|
| `migrations/tenant/postgres/058_performance_scoring_configuration.sql` (+ down, + mysql) | ✅ | 3 tabel baru: performance_components, performance_organization_components, performance_evaluation_components |
| `backend/internal/modules/performance/model.go` | ✅ | 3 struct baru: PerformanceComponent, PerformanceOrganizationComponent, PerformanceEvaluationComponent |
| `backend/internal/modules/performance/dto.go` | ✅ | DTO CRUD komponen, konfigurasi bobot organisasi, response snapshot komponen evaluasi |
| `backend/internal/modules/performance/repository.go` | ✅ | CRUD komponen, upsert konfigurasi bobot per organisasi, upsert snapshot evaluasi, batch lookup child organization & rata-rata skor |
| `backend/internal/modules/performance/service.go` | ✅ | `CalculateEvaluationComponentScoring` (scoring engine), validasi total bobot = 100%, resolusi skor per kode komponen (KPI/SUBORDINATE otomatis, lainnya manual) |
| `backend/internal/modules/performance/handler.go` | ✅ | HTTP handlers untuk semua endpoint Phase 5 |
| `backend/internal/modules/performance/routes.go` | ✅ | Route registration di bawah `/performance/kpi/*` |
| `backend/internal/modules/performance/module.go` | ✅ | AutoMigrate + seed default 3 komponen |

## Scoring Engine Behavior

- Jika Organization **belum** dikonfigurasi (tidak ada `performance_organization_components` enabled), evaluasi tetap memakai perhitungan KPI murni dari `RecalculateEvaluationScore` (backward compatible, tidak ada breaking change).
- Jika dikonfigurasi, `CalculateEvaluationComponentScoring` mem-validasi total bobot komponen enabled = 100% (menolak jika tidak), lalu menghitung skor per komponen:
  - **KPI**: diambil dari final_score evaluasi itu sendiri (hasil `evaluation_details`).
  - **SUBORDINATE**: rata-rata final_score evaluasi seluruh direct-child Organization (`parent_id`) pada periode yang sama.
  - **Komponen lain (mis. Work Program)**: tidak ada sumber data otomatis — diisi manual oleh reviewer via `PUT /evaluations/:id/components/:component_id`, nilai sebelumnya dipertahankan saat rekalkulasi ulang.
- Final score evaluasi = Σ(component_score × weight/100), disimpan sebagai snapshot per komponen di `performance_evaluation_components`.

## Frontend

✅ Diimplementasikan:
- `PerformanceComponentsView.vue` — master data Performance Components. Dikunci ke 3 komponen tetap (Indikator/KPI, Program, Kinerja Bawahan) — tombol New/Delete dihapus, hanya edit nama/deskripsi/urutan/is_active yang diizinkan.
- `PerformanceScoringConfigView.vue` — konfigurasi bobot komponen per Organization.
- Breakdown skor per komponen ditampilkan di `KPIEvaluationDetail.vue` (tabel "Scoring Components").

---

# Phase 8 - Approval Integration & Two-Phase Target/Realization

## Objective

Selaraskan proses approval KPI dengan modul `approval` terpusat (single source of truth untuk semua workflow approval lintas modul — lihat arsitektur approval module), dan mengubah alur pengajuan evaluasi karyawan menjadi dua tahap agar target disetujui atasan sebelum realisasi bisa diisi.

## Status Machine (baru)

```
DRAFT → TARGET_SUBMITTED → TARGET_APPROVED → SUBMITTED → APPROVED → COMPLETED
  ↑____________|                                  ↑___________|
   (reject target)                                 (reject realization)
```

- **DRAFT**: karyawan mengisi target indikator (+ item Program bila komponen Program aktif untuk Organization-nya).
- **TARGET_SUBMITTED → TARGET_APPROVED**: target diajukan ("Ajukan Target") dan disetujui/ditolak atasan via approval instance terpisah, module slug `performance_kpi_target`.
- **TARGET_APPROVED**: karyawan mengisi aktual indikator + item Program.
- **SUBMITTED → APPROVED → COMPLETED**: realisasi diajukan ("Ajukan Realisasi") dan disetujui/ditolak via approval instance terpisah, module slug `performance_kpi_realization`. Reject pada tahap ini kembali ke `TARGET_APPROVED` (bukan `DRAFT`), karena target sudah final.

Kedua module slug di atas di-alias ke subscription module `performance` (lihat `subscriptionModuleAliases`/`subscriptionModuleSubslots` di `approval/service.go`) sehingga tenant hanya perlu subscribe modul `performance` biasa, dan flow approval bisa dikonfigurasi generik di bawah slug `performance` lalu otomatis dipakai kedua checkpoint (fallback resolusi via `GetActiveFlowByModule`).

## Perilaku Kunci

- **Hard-fail bila flow dikonfigurasi tapi gagal resolve** (mis. seluruh hierarki atasan vakan): submit ditolak, status TIDAK berubah. Fallback manual-approval (tanpa modul approval) HANYA berlaku bila memang tidak ada flow yang dikonfigurasi sama sekali.
- **Walk-up hierarki organisasi**: bila atasan langsung vakan (tidak ada employment aktif), resolusi approver naik terus ke `parent_id` berikutnya sampai ketemu Organization terisi, atau gagal total bila mentok ke akar hierarki tanpa hasil (`resolveSupervisorAssignees` di `approval/service.go`).
- **Zero-assignee guard**: `CreateInstance` di modul approval menolak instance yang landing di step dengan nol approver ter-resolve (dulu diam-diam membuat instance yang tidak bisa disetujui siapapun).
- Field `target_approval_instance_id` / `realization_approval_instance_id` disimpan di `performance_evaluations` untuk melacak instance approval aktif per fase.

## Program Component (baru)

- `PerformanceEvaluationProgramItem` — sub-resource baru per evaluation (`performance_evaluation_program_items`), employee-authored (tidak dari template HR): title, weight, target, actual, formula_type, unit_of_measurement, score.
- Total bobot item Program per evaluation dibatasi maksimal 100%.
- Bisa diubah selama status masih `DRAFT` (fase target); actual hanya bisa diisi saat `TARGET_APPROVED`.
- Skor komponen Program = jumlah (bukan rata-rata) skor seluruh item — simetris dengan cara komponen KPI dihitung dari total skor `evaluation_details`.

## Template Simplification

- `KPITemplateForm.vue` — kolom Target & Unit dihapus dari tabel indikator; HR hanya mendefinisikan perspective, title, weight, formula. Target kini murni input karyawan di fase Ajukan Target.
- Dropdown Organization saat membuat template kini dibatasi hanya Organization di bawah hierarki milik karyawan pembuat (endpoint baru `GET /performance/kpi/templates/organization-scope`, walk-down descendant via `GetDescendantOrganizations`).

## Save UX

- Input Target dan Aktual di `KPIEvaluationDetail.vue` tidak lagi auto-save per field (on blur). Diganti tombol eksplisit **"Simpan Target"** / **"Simpan Aktual"** yang menyimpan seluruh baris indikator + item Program sekaligus dalam satu aksi.

---

# Phase 9 - Program Component, Bottom-Up Subordinate Scoring & OKR Self-Assessment

## Total Score Formula (dikonfirmasi)

```
Final Score = (Total Skor Indikator × Bobot KPI)
            + (Total Skor Program × Bobot Program)
            + (Rata-rata Skor Bawahan yang COMPLETED × Bobot Kinerja Bawahan)
```

Total skor Indikator dan Program adalah **jumlah**, bukan rata-rata (`CalculateEvaluationComponentScoring` di `service.go`).

## Bottom-Up Subordinate Scoring

- **Walk-up/walk-down melompati Organization vakan**: `ResolveEffectiveSupervisorOrgID` (naik) dan `GetEffectiveChildOrganizationIDs` (turun) — bila Organization langsung di atas/bawah tidak punya employment aktif, resolusi lanjut ke level berikutnya, bukan berhenti/gagal di level pertama.
- **Setiap bawahan efektif selalu ikut jadi pembagi rata-rata** (`GetAverageFinalScore`): yang `COMPLETED` memakai `final_score` aktualnya; yang belum `COMPLETED` — termasuk Organization yang sudah punya karyawan tapi belum pernah membuat target/evaluasi sama sekali — dihitung sebagai skor **0**, bukan dikecualikan dari rata-rata.
- **Propagasi otomatis**: saat sebuah evaluasi `COMPLETED`, skor Subordinate seluruh atasan efektif di sepanjang hierarki (sampai paling atas) dihitung ulang otomatis (`CompleteEvaluation` → `propagateSubordinateScoreUpward`), best-effort (tidak memblokir completion karyawan sendiri bila gagal).
- **Batch recalculation akhir periode**: endpoint `POST /performance/kpi/periods/:period_id/recalculate-scoring` + tombol "Hitung Ulang Skor" di halaman Performance Periods — menghitung ulang seluruh evaluasi dalam satu periode dari bawah ke atas hierarki (`RecalculatePeriodScoring`), untuk dipakai setelah periode penilaian berakhir.

## OKR Self-Assessment

- `GET /performance/okr/my-context` (backend) + `OKRSelfAssessment.vue` (frontend) di `/performance/okr/my-evaluation` — mirror dari alur self-assessment KPI yang sudah ada, memakai template OKR `status=1` (Active) milik Organization karyawan.

## Dashboard & Navigasi

- Sidebar: "Dashboard Kinerja" dipindah ke atas "KPI" pada submenu Performance.
- `PerformanceIndex.vue` — menu card dikelompokkan: **Self-Assessment** (highlight, deskripsi detail), **KPI**, **OKR**, **Shared Settings** (Periods dipindah keluar dari grup KPI karena dipakai bersama KPI & OKR). Card Self-Assessment otomatis muted & non-klik bila karyawan belum punya posisi aktif atau belum ada template aktif untuk Organization-nya.
- **Fix**: ringkasan (Quick Stats) di Dashboard Kinerja sempat selalu kosong karena hanya mencari Performance Period berstatus `active`, padahal Period baru default berstatus `draft`. Ditambahkan fallback `GetLatestPeriod` (Period terbaru apapun statusnya) pada Employee/Manager/HR dashboard.

## Bug Fixes Terkait

- `UpsertEvaluationComponent` menulis `CreatedAt` zero-value saat update (bukan create), menyebabkan error MySQL strict-mode `'0000-00-00'` pada `calculate-scoring` setelah komponen (mis. skor bawahan) pernah tersimpan sebelumnya — diperbaiki dengan mempertahankan `CreatedAt` asli dan `Updates()` kolom spesifik alih-alih `Save()` penuh.
- Approve/Reject pada task approval sempat mengembalikan `403 Forbidden` untuk role Employee (default view-only) karena middleware RBAC generik mewajibkan permission `approval.create` — padahal otorisasi sebenarnya (assignee task pending) sudah divalidasi di `service.SubmitAction`. Endpoint `POST /approval/instances/:id/actions` dikecualikan dari blanket permission check tersebut.