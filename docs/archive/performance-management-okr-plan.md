# OKR Module Development Plan

> 🔗 **Arsip dokumentasi:** [`docs/README.md`](../README.md) — dokumen ini adalah referensi historis dari modul yang sudah selesai diimplementasikan.


## Objective

Membangun modul Objective & Key Results (OKR) yang terintegrasi dengan HRIS Platform.

Karakteristik modul:

- Organization = Position = 1 Employee
- Objective mengikuti Organization
- Tidak menggunakan Assignment
- Mendukung Multiple Objective
- Mendukung Multiple Key Result
- Mendukung Progress Check-in
- Mendukung Evidence Attachment
- Menggunakan Snapshot saat Evaluation
- Seluruh Primary Key menggunakan UUID

---

# Phase 1 - Shared Master

Menggunakan master yang sudah ada pada modul Performance (KPI).

| Table | Status | Description |
|-------|--------|-------------|
| performance_periods | Reuse | Periode evaluasi (Annual, Semester, Quarterly) |
| performance_ratings | Reuse | Rating hasil evaluasi (Outstanding, Excellent, etc) |
| performance_logs | Reuse | Audit trail aktivitas |

Tidak perlu membuat ulang tabel master di atas.

---

# Phase 2 - New Tables

## 1. okr_templates

Template OKR untuk setiap Organization.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| organization_id | uuid FK | NO | - | Organization pemilik template |
| period_id | uuid FK | YES | NULL | Performance Period |
| name | varchar(255) | NO | - | Nama Template |
| description | text | YES | NULL | Deskripsi |
| status | smallint | NO | 0 | 0=Draft, 1=Active, 2=Inactive |
| effective_date | date | YES | NULL | Berlaku mulai |
| expired_date | date | YES | NULL | Berlaku sampai |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |
| deleted_at | timestamp | YES | NULL | Soft Delete |

**Indexes:**
- `idx_okr_tpl_org` on (organization_id)
- `idx_okr_tpl_period` on (period_id)
- `idx_okr_tpl_status` on (status)

---

## 2. okr_objectives

Objective yang ingin dicapai.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| template_id | uuid FK | NO | - | OKR Template |
| code | varchar(50) | YES | NULL | Kode Objective |
| title | varchar(255) | NO | - | Judul Objective |
| description | text | YES | NULL | Deskripsi |
| weight | decimal(5,2) | NO | 0 | Bobot Objective (total semua objective = 100%) |
| sort_order | integer | NO | 0 | Urutan tampilan |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |
| deleted_at | timestamp | YES | NULL | Soft Delete |

**Indexes:**
- `idx_okr_obj_template` on (template_id)
- `idx_okr_obj_sort` on (sort_order)

---

## 3. okr_key_results

Target terukur dari setiap Objective.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| objective_id | uuid FK | NO | - | Objective |
| code | varchar(50) | YES | NULL | Kode Key Result |
| title | varchar(255) | NO | - | Judul Key Result |
| description | text | YES | NULL | Deskripsi |
| target_type | varchar(30) | NO | NUMBER | Number, Currency, Percentage, Duration, Boolean |
| target_value | decimal(18,2) | NO | 0 | Nilai Target |
| unit | varchar(50) | YES | NULL | Satuan (%, Hari, Jam, Rupiah, Unit) |
| formula_type | varchar(30) | NO | HIGHER_BETTER | Formula kalkulasi |
| weight | decimal(5,2) | NO | 0 | Bobot dalam Objective |
| minimum_score | decimal(5,2) | NO | 0 | Skor minimum |
| maximum_score | decimal(5,2) | NO | 100 | Skor maksimum |
| sort_order | integer | NO | 0 | Urutan tampilan |
| is_required | boolean | NO | true | Key Result wajib |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |
| deleted_at | timestamp | YES | NULL | Soft Delete |

**Indexes:**
- `idx_okr_kr_objective` on (objective_id)
- `idx_okr_kr_sort` on (sort_order)

---

## 4. okr_evaluations

Header penilaian OKR.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| employee_id | uuid FK | NO | - | Employee yang dievaluasi |
| organization_id | uuid FK | NO | - | Organization saat evaluasi |
| period_id | uuid FK | NO | - | Performance Period |
| template_id | uuid FK | YES | NULL | OKR Template yang digunakan |
| status | varchar(20) | NO | DRAFT | DRAFT, KR_SUBMITTED, KR_APPROVED, SUBMITTED, APPROVED (legacy), COMPLETED, REJECTED |
| submitted_at | timestamp | YES | NULL | Waktu submit assessment |
| submitted_by | uuid FK | YES | NULL | User yang submit |
| approved_at | timestamp | YES | NULL | Waktu approval akhir |
| approved_by | uuid FK | YES | NULL | User yang approve |
| kr_approval_instance_id | uuid FK | YES | NULL | Approval instance checkpoint Key Result (Phase 7) |
| assessment_approval_instance_id | uuid FK | YES | NULL | Approval instance checkpoint assessment (Phase 7) |
| kr_submitted_at | timestamp | YES | NULL | Waktu submit proposal Key Result (Phase 7) |
| final_score | decimal(5,2) | NO | 0 | Nilai akhir |
| rating_id | uuid FK | YES | NULL | Performance Rating |
| reviewer_notes | text | YES | NULL | Catatan reviewer |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |
| deleted_at | timestamp | YES | NULL | Soft Delete |

**Indexes:**
- `idx_okr_eval_employee` on (employee_id)
- `idx_okr_eval_org` on (organization_id)
- `idx_okr_eval_period` on (period_id)
- `idx_okr_eval_status` on (status)

---

## 5. okr_evaluation_details

Snapshot Objective saat evaluasi + Key Result yang **diusulkan karyawan** saat fase DRAFT (Phase 7). Evaluasi dibuat hanya dengan referensi `template_id` — Objective TIDAK disalin sebagai detail. Detail dibuat per proposal Key Result lewat `POST .../evaluations/:id/key-results` (judul/bobot Objective & Key Result di-snapshot ke detail).

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| evaluation_id | uuid FK | NO | - | Evaluation header |
| objective_id | uuid FK | YES | NULL | Reference Objective (nullable jika deleted) |
| key_result_id | uuid FK | YES | NULL | Reference Key Result (nullable jika deleted) |
| objective_title | varchar(255) | NO | - | Snapshot judul Objective |
| key_result_title | varchar(255) | NO | - | Snapshot judul Key Result |
| objective_weight | decimal(5,2) | NO | 0 | Snapshot bobot Objective |
| key_result_weight | decimal(5,2) | NO | 0 | Snapshot bobot Key Result |
| target_value | decimal(18,2) | NO | 0 | Snapshot target |
| target_type | varchar(30) | NO | NUMBER | Snapshot tipe target |
| unit | varchar(50) | YES | NULL | Snapshot satuan |
| formula_type | varchar(30) | NO | HIGHER_BETTER | Snapshot formula |
| actual_value | decimal(18,2) | NO | 0 | Nilai aktual |
| achievement | decimal(5,2) | NO | 0 | Persentase pencapaian |
| score | decimal(5,2) | NO | 0 | Nilai/skor |
| remarks | text | YES | NULL | Catatan reviewer |
| sort_order | integer | NO | 0 | Urutan tampilan |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |

**Indexes:**
- `idx_okr_detail_eval` on (evaluation_id)
- `idx_okr_detail_obj` on (objective_id)
- `idx_okr_detail_kr` on (key_result_id)

---

## 6. okr_progress

Progress Check-in mingguan/bulanan.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| evaluation_detail_id | uuid FK | NO | - | Evaluation Detail |
| progress_date | date | NO | - | Tanggal check-in |
| actual_value | decimal(18,2) | NO | 0 | Nilai aktual saat check-in |
| achievement | decimal(5,2) | NO | 0 | Persentase pencapaian |
| notes | text | YES | NULL | Catatan progress |
| created_by | uuid FK | NO | - | User yang input |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |

**Indexes:**
- `idx_okr_prog_detail` on (evaluation_detail_id)
- `idx_okr_prog_date` on (progress_date)

---

## 7. okr_comments

Komentar antara Employee dan Reviewer.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| evaluation_id | uuid FK | NO | - | Evaluation header |
| parent_id | uuid FK | YES | NULL | Parent comment (untuk reply) |
| comment | text | NO | - | Isi komentar |
| created_by | uuid FK | NO | - | User yang komentar |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |

**Indexes:**
- `idx_okr_comment_eval` on (evaluation_id)
- `idx_okr_comment_parent` on (parent_id)
- `idx_okr_comment_user` on (created_by)

---

## 8. okr_attachments

Evidence pencapaian Key Result.

| Field | Type | Nullable | Default | Description |
|-------|------|----------|---------|-------------|
| id | uuid PK | NO | - | Primary Key |
| evaluation_detail_id | uuid FK | NO | - | Evaluation Detail |
| file_path | varchar(500) | NO | - | Path file di storage |
| file_name | varchar(255) | NO | - | Nama file asli |
| file_type | varchar(100) | YES | NULL | MIME type |
| file_size | bigint | YES | NULL | Ukuran file (bytes) |
| description | text | YES | NULL | Deskripsi attachment |
| uploaded_by | uuid FK | NO | - | User yang upload |
| created_at | timestamp | NO | - | |
| updated_at | timestamp | NO | - | |

**Indexes:**
- `idx_okr_attach_detail` on (evaluation_detail_id)
- `idx_okr_attach_user` on (uploaded_by)

---

# Phase 3 - Business Process

## OKR Setup Flow

```text
Performance Period
        │
        ▼
Organization
        │
        ▼
OKR Template  ── hanya Objective (title/weight), Key Results TIDAK lagi di template (Phase 7)
        │
        ├──────────────────┐
        ▼                  ▼
Objective 1            Objective 2
```

---

## OKR Evaluation Flow (dua fase, Phase 7)

```text
Employee
      │
      ▼
Create Evaluation ──────────► hanya mereferensikan TemplateID (Objective diterima via template_id)
      │
      ▼
Propose Key Results ────────► POST .../evaluations/:id/key-results (hanya saat DRAFT)
      │
      ▼
Submit Key Results ─────────► DRAFT → KR_SUBMITTED (routed approval modul okr_key_result)
      │
      ▼
Approve Key Results ────────► KR_SUBMITTED → KR_APPROVED ("OKR Active")
      │
      ▼
Input Actual Value           (gated ke status KR_APPROVED)
      │
      ▼
Calculate Achievement ──────► Achievement = (Actual / Target) × 100
      │
      ▼
Calculate Score ────────────► Score = (Weight × Achievement) / 100
      │
      ▼
Submit Assessment ──────────► KR_APPROVED → SUBMITTED (routed approval modul okr_assessment)
      │
      ▼
Approve Assessment ─────────► SUBMITTED → COMPLETED langsung (auto-complete)
```

---

## Weekly Check-in Flow

```text
Evaluation Detail
        │
        ▼
Add Progress Entry
        │
        ▼
Input Actual Value
        │
        ▼
Auto-Calculate Achievement
        │
        ▼
Update Dashboard
```

---

## Review Workflow (superseded — lihat Phase 7)

> ⚠️ Diagram di bawah adalah rencana awal (single-checkpoint). Implementasi final memakai status machine dua-fase — lihat "Phase 7 - Two-Phase Flow" untuk versi yang benar-benar berjalan.

```text
┌─────────┐
│  DRAFT  │
└────┬────┘
     │ Employee Submit
     ▼
┌──────────┐
│SUBMITTED │
└────┬─────┘
     │ Manager Review
     ▼
┌──────────┐     ┌──────────┐
│ APPROVED │ ←── │ REVISION │ (Optional)
└────┬─────┘     └──────────┘
     │
     ▼
┌───────────┐
│ COMPLETED │
└───────────┘
```

---

# Phase 4 - Dashboard

## Employee Dashboard

| Metric | Description |
|--------|-------------|
| Objective Progress | Progress per Objective |
| Key Result Progress | Progress per Key Result |
| Overall Achievement | Total pencapaian |
| Final Score | Nilai akhir |
| Rating | Rating berdasarkan skor |
| Evidence Count | Jumlah evidence yang diupload |
| Check-in Status | Status check-in mingguan |

---

## Manager Dashboard

| Metric | Description |
|--------|-------------|
| Team OKR Overview | Ringkasan OKR tim |
| Objective Progress | Progress per objective tim |
| Check-in Compliance | Kepatuhan check-in tim |
| Need Review | Evaluasi yang perlu direview |
| Overdue Review | Review yang terlambat |
| Team Achievement | Rata-rata pencapaian tim |

---

## HR Dashboard

| Metric | Description |
|--------|-------------|
| OKR Completion Rate | Persentase OKR selesai |
| Rating Distribution | Distribusi rating |
| Objective Achievement | Rata-rata pencapaian objective |
| Key Result Achievement | Rata-rata pencapaian key result |
| Achievement per Organization | Pencapaian per organisasi |
| Trend per Period | Tren pencapaian per periode |
| Top Performers | Karyawan dengan skor tertinggi |
| Bottom Performers | Karyawan yang perlu perhatian |

---

# Phase 5 - Future Enhancement

- Mid Year Review
- OKR Alignment (Cascading dari Company → Department → Individual)
- OKR Versioning (Track perubahan objective/key result)
- Reminder Notification (Check-in, Deadline)
- Multi Level Approval
- Bonus Integration
- Performance Improvement Plan (PIP) Integration
- Succession Planning Integration
- Career Path Integration
- Talent Management Integration
- 360 Feedback Integration
- OKR Scoring Calibration

---

# Phase 6 - Master Data & Seeder

## Seeder Overview

| Table | Seeder | Required | Description |
|-------|:------:|:--------:|-------------|
| performance_periods | ❌ | No | Reuse dari KPI, data transaksi |
| performance_ratings | ✅ | Yes | Reuse dari KPI |
| performance_logs | ❌ | No | Reuse dari KPI, audit trail |
| okr_templates | ❌ | No | Berbeda setiap Organization |
| okr_objectives | ❌ | No | Dibuat sesuai kebutuhan |
| okr_key_results | ❌ | No | Dibuat sesuai kebutuhan |
| okr_evaluations | ❌ | No | Data transaksi |
| okr_evaluation_details | ❌ | No | Data transaksi |
| okr_progress | ❌ | No | Data transaksi |
| okr_comments | ❌ | No | Data transaksi |
| okr_attachments | ❌ | No | Data transaksi |

---

## Shared Seeders (dari KPI Module)

### Performance Ratings

| Code | Name | Min Score | Max Score | Color |
|------|------|----------:|----------:|-------|
| OUT | Outstanding | 95.00 | 100.00 | success |
| EXC | Excellent | 85.00 | 94.99 | primary |
| GOO | Good | 75.00 | 84.99 | info |
| FAI | Fair | 60.00 | 74.99 | warning |
| POO | Poor | 0.00 | 59.99 | danger |

### Performance Indicator Formulas

| Code | Name | Description |
|------|------|-------------|
| MANUAL | Manual Score | Nilai diinput manual oleh reviewer |
| HIGHER_BETTER | Higher Better | Semakin tinggi nilai semakin baik |
| LOWER_BETTER | Lower Better | Semakin rendah nilai semakin baik |
| RANGE | Range Score | Nilai berdasarkan rentang tertentu |
| BOOLEAN | Boolean | Ya/Tidak (100% atau 0%) |
| PERCENTAGE | Percentage | Nilai berupa persentase |

---

# Optional Master Data

Untuk fleksibilitas jangka panjang, beberapa field yang saat ini menggunakan Enum dapat dipisahkan menjadi Master Data.

## okr_target_types

| Code | Name |
|------|------|
| NUMBER | Number |
| PERCENTAGE | Percentage |
| CURRENCY | Currency |
| BOOLEAN | Boolean |
| DURATION | Duration |
| SCORE | Score |

---

## okr_units

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
| PROJECT | Project |
| TASK | Task |

---

## okr_statuses

> Status machine implementasi final adalah **dua fase** (Phase 7): `DRAFT → KR_SUBMITTED → KR_APPROVED → SUBMITTED → COMPLETED` (atau `REJECTED` di tiap fase). Tabel di bawah adalah daftar status master data opsional dari rencana awal — hanya sebagian yang dipakai implementasi.

| Code | Name | Description |
|------|------|-------------|
| DRAFT | Draft | Evaluasi baru dibuat — karyawan mengusulkan Key Results |
| KR_SUBMITTED | Key Results Submitted | Proposal Key Results diajukan (Phase 7) |
| KR_APPROVED | Key Results Approved / OKR Active | Proposal Key Results disetujui (Phase 7) |
| SUBMITTED | Submitted | Assessment akhir diajukan |
| APPROVED | Approved | Status legacy — hanya dipakai evaluasi lama |
| COMPLETED | Completed | Selesai (auto-complete saat approve assessment final) |
| REJECTED | Rejected | Ditolak — reject KR → `DRAFT`, reject assessment → `KR_APPROVED` |

---

# Seeder Structure

Struktur dan penamaan mengikuti pola yang sudah ada pada KPI module.

```text
backend/internal/pkg/tenantseed/
└── seed_data.go
    ├── SeedPerformancePerspectives()  // Reuse
    ├── SeedPerformanceRatings()       // Reuse
    └── SeedPerformanceFormulas()      // Reuse
```

OKR tidak memerlukan seeder tambahan karena menggunakan master data dari KPI.

---

# Seeder Priority

## Mandatory (dari KPI)

- PerformancePerspectiveSeeder (jika menggunakan BSC)
- PerformanceRatingSeeder
- PerformanceIndicatorFormulaSeeder

## Optional

- OkrTargetTypeSeeder
- OkrUnitSeeder
- OkrStatusSeeder

---

# API Endpoints

> ℹ️ Endpoint di bawah adalah rencana awal dan sebagian besar path-nya sudah berubah pada implementasi final (prefix `/performance/okr/*`, bukan `/okr/*`, dan beberapa endpoint ditambah/diubah oleh Phase 7-8). Lihat "Phase 7/8" untuk daftar endpoint final.

## OKR Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/templates` | List templates |
| POST | `/api/v1/tenant/performance/okr/templates` | Create template — hanya mendefinisikan Objective (title/weight); dibatasi ke Organization bawahan efektif (Phase 8) |
| GET | `/api/v1/tenant/performance/okr/templates/:id` | Get template |
| PUT | `/api/v1/tenant/performance/okr/templates/:id` | Update template |
| DELETE | `/api/v1/tenant/performance/okr/templates/:id` | Delete template |
| POST | `/api/v1/tenant/performance/okr/templates/:id/duplicate` | Duplicate template |
| GET | `/api/v1/tenant/performance/okr/templates/objective-scope` | **(Baru, Phase 8)** Resolve eligibility + daftar Organization bawahan efektif pemanggil |

## OKR Objectives

> Baris bertanda *(rencana)* tidak ada di routes.go — hanya rencana awal.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/objectives` | List objectives *(rencana)* |
| POST | `/api/v1/tenant/performance/okr/objectives` | Create objective |
| GET | `/api/v1/tenant/performance/okr/objectives/:id` | Get objective |
| PUT | `/api/v1/tenant/performance/okr/objectives/:id` | Update objective |
| DELETE | `/api/v1/tenant/performance/okr/objectives/:id` | Delete objective |

## OKR Key Results

Master data Key Result di Template (`okr_key_results`) masih ada sebagai endpoint, namun **sudah tidak dipakai** oleh frontend `OKRTemplateForm.vue` sejak Phase 7 — Key Result kini diusulkan karyawan sendiri per Objective saat evaluasi (lihat "OKR Evaluation Key Results" di bawah), bukan didefinisikan HR di Template.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/key-results` | List key results (master data Template) *(rencana)* |
| POST | `/api/v1/tenant/performance/okr/key-results` | Create key result |
| GET | `/api/v1/tenant/performance/okr/key-results/:id` | Get key result |
| PUT | `/api/v1/tenant/performance/okr/key-results/:id` | Update key result |
| DELETE | `/api/v1/tenant/performance/okr/key-results/:id` | Delete key result |

## OKR Evaluation Key Results (Baru, Phase 7)

Key Result yang benar-benar dipakai evaluasi — diusulkan karyawan sendiri per Objective, hanya selama status `DRAFT`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/key-results` | Usulkan Key Result baru di bawah sebuah Objective |
| PUT | `/api/v1/tenant/performance/okr/evaluation-key-results/:id/target` | Edit target Key Result yang diusulkan |
| DELETE | `/api/v1/tenant/performance/okr/evaluation-key-results/:id` | Hapus Key Result yang diusulkan |

## OKR Evaluations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/evaluations` | List evaluations |
| POST | `/api/v1/tenant/performance/okr/evaluations` | Create evaluation with snapshot (hanya snapshot Objective, tanpa Key Result) |
| GET | `/api/v1/tenant/performance/okr/evaluations/:id` | Get evaluation |
| GET | `/api/v1/tenant/performance/okr/evaluations/:id/details` | Get evaluation with details |
| PUT | `/api/v1/tenant/performance/okr/evaluations/:id` | Update evaluation |
| DELETE | `/api/v1/tenant/performance/okr/evaluations/:id` | Delete evaluation |

## OKR Evaluation Details

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/evaluations/:id/details` | List details |
| PUT | `/api/v1/tenant/performance/okr/evaluation-details/:id` | Update actual value — backend-gated ke status `KR_APPROVED` |
| PUT | `/api/v1/tenant/performance/okr/evaluations/:id/actuals` | Bulk update actuals |

## OKR Workflow (dua fase, Phase 7)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/submit-key-results` | **Baru.** `DRAFT` → `KR_SUBMITTED`, routed lewat central approval module `okr_key_result` |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/approve-key-results` | **Baru.** Fallback manual `KR_SUBMITTED` → `KR_APPROVED` ("OKR Active") saat belum ada flow approval dikonfigurasi |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/reject-key-results` | **Baru.** Fallback manual `KR_SUBMITTED` → `DRAFT` |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/recalculate` | Recalculate score |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/submit` | `KR_APPROVED` → `SUBMITTED`, routed lewat central approval module `okr_assessment` |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/approve` | `SUBMITTED` → **`COMPLETED` langsung** (auto-complete, tidak ada langkah manual "Complete" terpisah) |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/reject` | `SUBMITTED` → `KR_APPROVED` (bukan ke `DRAFT` — Key Result yang sudah disetujui tidak ikut ditolak) |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/complete` | Fallback manual lama, hanya reachable untuk evaluasi lama yang masih berstatus `APPROVED` (legacy) |

## OKR Progress

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/performance/okr/progress` | Create progress |
| GET | `/api/v1/tenant/performance/okr/evaluation-details/:id/progress` | List progress |
| GET | `/api/v1/tenant/performance/okr/progress/:id` | Get progress |
| PUT | `/api/v1/tenant/performance/okr/progress/:id` | Update progress |
| DELETE | `/api/v1/tenant/performance/okr/progress/:id` | Delete progress |

## OKR Comments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/performance/okr/comments` | Create comment |
| GET | `/api/v1/tenant/performance/okr/evaluations/:id/comments` | List comments |
| PUT | `/api/v1/tenant/performance/okr/comments/:id` | Update comment |
| DELETE | `/api/v1/tenant/performance/okr/comments/:id` | Delete comment |

## OKR Attachments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/performance/okr/attachments` | Upload attachment |
| GET | `/api/v1/tenant/performance/okr/evaluation-details/:id/attachments` | List attachments |
| GET | `/api/v1/tenant/performance/okr/attachments/:id` | Get attachment *(rencana)* |
| DELETE | `/api/v1/tenant/performance/okr/attachments/:id` | Delete attachment |

## OKR Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/dashboard/employee/:id` | Employee dashboard *(rencana)* |
| GET | `/api/v1/tenant/performance/okr/dashboard/manager/:id` | Manager dashboard *(rencana)* |
| GET | `/api/v1/tenant/performance/okr/dashboard/hr` | HR dashboard |

---

# Integration Points

Modul OKR terintegrasi dengan:

| Module | Integration |
|--------|-------------|
| Performance (KPI) | Shared master: periods, ratings, formulas, logs |
| Organization | Organization hierarchy untuk cascading |
| Employee | Employee data untuk evaluasi |
| Notification | Reminder check-in dan deadline |
| File Storage | Upload evidence/attachment |
| Audit Trail | Log semua aktivitas |

---

# Design Principles

- Seluruh Primary Key menggunakan UUID
- Seluruh Foreign Key menggunakan UUID
- Organization = Position = 1 Employee
- Objective mengikuti Organization
- Tidak menggunakan Assignment
- Snapshot Objective saat Evaluation; Key Result diusulkan karyawan per Objective (Phase 7)
- Histori tidak berubah walaupun Template diperbarui
- Mendukung Soft Delete
- Mendukung Weekly/Monthly Check-in
- Mendukung Evidence Attachment
- Mendukung Comment Thread
- Mendukung Audit Trail
- Master Data menggunakan Seeder (reuse dari KPI)
- Siap diintegrasikan dengan Performance Management, Competency Management, Career Path, Succession Planning, Talent Management, dan Bonus Management

---

# Implementation Status

| Phase | Status | Completion Date | Notes |
|-------|--------|-----------------|-------|
| Phase 1 - Shared Master | ✅ Completed | - | Reuse dari KPI module |
| Phase 2 - Database | ✅ Completed | 2026-08-05 | Migration 057_okr.sql, Models |
| Phase 3 - Business Logic | ✅ Completed | 2026-08-06 | DTO, Repository, Service, Handler, Routes, Module |
| Phase 4 - Dashboard | ✅ Completed | 2026-08-06 | HR Dashboard endpoint |
| Phase 5 - Future Enhancement | 🔶 Partially completed | 2026-08-08 | Multi Level Approval + OKR Alignment (cascading) selesai — lihat Phase 7/8; sisanya (Mid Year Review, OKR Versioning, Reminder Notification, dst.) masih pending |
| Phase 6 - Seeder | ✅ Completed | - | Reuse dari KPI (tidak perlu seeder baru) |
| Phase 7 - Two-Phase Flow & Central Approval Integration | ✅ Completed | 2026-08-08 | Lihat detail di bawah |
| Phase 8 - Cascading Objective Creation | ✅ Completed | 2026-08-08 | Lihat detail di bawah |

## Module Integration

OKR telah diintegrasikan ke dalam Performance module sebagai sub-module (2026-08-06).

**Struktur:**
```
backend/internal/modules/performance/
├── model.go              # KPI models
├── dto.go                # KPI DTOs  
├── repository.go         # KPI repository
├── service.go            # KPI service
├── handler.go            # KPI handlers
├── okr_model.go          # OKR models (8 entities)
├── okr_dto.go            # OKR DTOs
├── okr_repository.go     # OKR repository (~40 methods)
├── okr_service.go        # OKR service + workflow
├── okr_handler.go        # OKR handlers
├── routes.go             # KPI + OKR routes
└── module.go             # Module init (KPI + OKR)
```

## API Endpoints

**Shared Master:** `/performance/*`
- `/performance/periods/*` - Periode evaluasi
- `/performance/ratings/*` - Rating master
- `/performance/indicator-formulas/*` - Formula master

**KPI:** `/performance/kpi/*`
- Templates, Indicators, Perspectives, Evaluations, Progress, Comments, Attachments, Dashboard

**OKR:** `/performance/okr/*`
- Templates, Objectives, Key Results (master data), Evaluations, Evaluation Key Results, Progress, Comments, Attachments, Dashboard

---

# Phase 7 - Two-Phase Flow & Central Approval Integration

## Objective

Alur awal OKR (single-checkpoint `DRAFT → SUBMITTED → APPROVED → COMPLETED`, Key Result didefinisikan HR di Template) tidak sesuai dengan proses yang diinginkan:

```
1. Atasan membuat OKR
2. Bawahan menerima Objective
3. Bawahan menyusun/mengusulkan Key Result
4. Atasan melakukan Review
5. Atasan Approve
6. OKR Active
7. Bawahan melakukan Check-in
8. Bawahan Self Assessment
9. Atasan melakukan Assessment
10. Final Score
```

Phase 7 menerapkan pola yang sama persis dengan enhancement KPI (lihat `docs/performance-management-kpi-plan.md` Phase 8/9): status machine dua-fase, sub-resource Key Result yang diusulkan karyawan sendiri, integrasi modul approval terpusat, dan auto-complete pada approval final.

## Status Machine

```
DRAFT → KR_SUBMITTED → KR_APPROVED → SUBMITTED → COMPLETED (auto saat approve final)
  ↑____________|                          ↑____________|
   (reject Key Result)                  (reject assessment → kembali ke KR_APPROVED)
```

| Langkah pada flow yang diminta | Mekanisme |
|---|---|
| 1. Atasan membuat OKR | Template + Objective (title/weight saja) via `OKRTemplateForm.vue`, dibatasi Organization sesuai Phase 8 |
| 2. Bawahan menerima Objective | Karyawan memilih Template ber-status Active untuk Organization-nya via `OKRSelfAssessment.vue` → snapshot Objective saja (Key Result tidak ikut disalin) |
| 3. Bawahan mengusulkan Key Result | **Baru** — endpoint `POST .../evaluations/:id/key-results`, hanya selama `DRAFT`, bobot Key Result per Objective dibatasi maksimal 100% |
| 4-5. Atasan Review/Approve | Checkpoint baru: `submit-key-results` (`DRAFT`→`KR_SUBMITTED`), routed lewat modul approval `okr_key_result`; approve→`KR_APPROVED`, reject→`DRAFT` |
| 6. OKR Active | Status `KR_APPROVED` |
| 7. Check-in | `okr_progress` (sudah ada sebelumnya), dipakai selama `KR_APPROVED` |
| 8. Self Assessment | Isi `actual_value` per Key Result (backend digerbang ke status `KR_APPROVED`, sebelumnya hanya digerbang di frontend) + tombol batch "Simpan Aktual", lalu `submit` (`KR_APPROVED`→`SUBMITTED`), routed lewat modul approval `okr_assessment` |
| 9. Atasan Assessment | Checkpoint `SUBMITTED` yang sudah ada, kini routed lewat approval terpusat; **approve final langsung `COMPLETED`** (tidak ada tombol "Selesaikan" manual lagi) |
| 10. Final Score | `RecalculateEvaluationScore` (tidak berubah) |

## Integrasi Modul Approval

Dua module slug baru, `okr_key_result` dan `okr_assessment`, di-alias ke subscription module `performance` (sama seperti `performance_kpi_target`/`performance_kpi_realization`) lewat `subscriptionModuleAliases`/`subscriptionModuleSubslots` di `approval/service.go` — tenant cukup subscribe modul `performance`, flow bisa dikonfigurasi generik di slug `performance` dan otomatis dipakai kedua checkpoint OKR maupun KPI.

Perilaku hard-fail-jika-flow-dikonfigurasi-tapi-gagal-resolve (bukan silent fallback) diterapkan identik dengan KPI — lihat `SubmitKeyResults`/`SubmitEvaluation` di `okr_service.go`.

## Deliverables

| File | Perubahan |
|---|---|
| `migrations/tenant/{mysql,postgres}/069_okr_two_phase.sql` (+down) | Tambah `kr_approval_instance_id`, `assessment_approval_instance_id`, `kr_submitted_at` ke `okr_evaluations` |
| `okr_model.go` | Status baru `KR_SUBMITTED`/`KR_APPROVED` + 3 field baru |
| `okr_dto.go` | DTO Key Result yang diusulkan karyawan (`CreateOKREvaluationKeyResultRequest`, dst.) |
| `okr_repository.go` | CRUD Key Result evaluasi, `DeleteOKREvaluationDetail` |
| `okr_service.go` | `CreateEvaluationKeyResult`/`Update.../Delete...`, `SubmitKeyResults`/`ApproveKeyResults`/`RejectKeyResults`, `HandleKeyResultApprovalStatusChange`/`HandleAssessmentApprovalStatusChange`, `SetApprovalEngine`, auto-complete di `ApproveEvaluation` |
| `okr_handler.go`, `routes.go` | Handler + route endpoint baru |
| `module.go`, `cmd/server/main.go` | `NewModuleWithServices` menerima `okrSvc` yang sudah di-wire dengan `sharedApprovalEngine` sebelum modul di-mount |
| `approval/service.go` | Tambah alias `okr_key_result`/`okr_assessment` → `performance` |
| `OKREvaluationDetail.vue` | Rework total: Key Result dikelompokkan per Objective, tombol "Ajukan Key Result"/"Simpan Target"/"Simpan Aktual"/"Ajukan Assessment", hint fase, tidak ada tombol "Selesaikan" manual |
| `OKRIndex.vue` | Filter status tambah `KR_SUBMITTED`/`KR_APPROVED` |
| `Approvals.vue` | Render "submitted data" untuk modul `okr_key_result`/`okr_assessment` di halaman My Tasks (sebelumnya kosong karena tidak ada mapping endpoint) |
| `okr_two_phase_test.go`, `okr_approval_routing_test.go` | Test dua-fase + routing approval (10 test) |

## Bug Fixes Terkait

- `EmployeeName`/`OrganizationName`/`PeriodCode` ada di `OKREvaluationResponse` tapi tidak pernah diisi (pola yang sama berulang seperti KPI) — diperbaiki dengan `enrichEvaluationResponses` batch enrichment.
- `GetOrganizationParentID` (KPI maupun OKR) memakai `Pluck` ke scalar `*string`, ternyata tidak reliable di sqlite test driver ("Scan called without calling Next") — baru ketahuan setelah Phase 8 menulis test yang benar-benar menjalankannya. Diperbaiki ke pola `.Select().Scan()` yang sudah terbukti benar di modul approval.

---

# Phase 8 - Cascading Objective Creation (top-down)

## Objective

Alur pembuatan Objective harus mengikuti hierarki organisasi dari atas ke bawah:

1. Organization paling atas hanya bisa membuat Objective **untuk bawahannya**, bukan untuk dirinya sendiri.
2. Karyawan baru bisa mengusulkan Key Result (self-assessment) setelah Organization-nya menerima Objective — sudah otomatis berlaku lewat pengecekan Template aktif yang ada di `GetMyOKRContext`.
3. Setelah menerima Objective, karyawan tersebut boleh membuat Objective untuk bawahannya sendiri — cascading berlanjut ke bawah.

"Bawahan" = Organization **langsung** (`parent_id`), tapi bila Organization anak langsung vakan (tidak ada employment aktif), resolusi jalan terus turun sampai ketemu Organization yang terisi — pola walk-down yang sama dengan Subordinate KPI scoring. Seseorang wajib sudah punya Objective miliknya sendiri sebelum boleh membuat untuk bawahan — **kecuali** paling atas hierarki (tidak punya parent Organization), yang memulai cascade tanpa syarat ini.

## Implementasi

- `GetEffectiveChildOrganizationIDs` (walk-down, skip Organization vakan) dan `HasActiveTemplateWithObjectives` (cek "sudah menerima Objective") ditambahkan ke `OKRRepository`.
- `GetObjectiveCreationScope` (endpoint baru `GET /performance/okr/templates/objective-scope`) me-resolve eligibility + daftar Organization bawahan efektif pemanggil — dipakai `OKRTemplateForm.vue` untuk mengisi dropdown Organization saat membuat Objective baru (mode edit tetap memakai daftar Organization penuh, gate hanya berlaku saat create).
- `CreateTemplate` backend-enforced: pemanggil harus eligible, dan `organization_id` target harus anggota himpunan bawahan efektif pemanggil (tidak boleh Organization milik sendiri).
- 6 test baru (`okr_objective_scope_test.go`) mencakup: paling atas bisa membuat untuk bawahan langsung, tidak bisa untuk diri sendiri, tidak bisa untuk Organization di luar hierarkinya, karyawan menengah wajib punya Objective sendiri dulu, walk-down melewati Organization vakan, dan user tanpa posisi aktif.

## Template Simplification (ikut Phase 8)

`OKRTemplateForm.vue` disederhanakan: hanya mendefinisikan Objective (title + weight). Tabel Key Result di form Template dihapus total — Key Result kini murni diusulkan karyawan sendiri per Objective saat mengisi evaluasi (Phase 7), bukan didefinisikan HR di Template.
