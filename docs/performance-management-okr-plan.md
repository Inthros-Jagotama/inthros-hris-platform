# OKR Module Development Plan

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
| status | varchar(20) | NO | DRAFT | DRAFT, SUBMITTED, APPROVED, COMPLETED |
| submitted_at | timestamp | YES | NULL | Waktu submit |
| submitted_by | uuid FK | YES | NULL | User yang submit |
| approved_at | timestamp | YES | NULL | Waktu approval |
| approved_by | uuid FK | YES | NULL | User yang approve |
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

Snapshot Objective dan Key Result saat evaluasi.

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
OKR Template
        │
        ├──────────────────┐
        ▼                  ▼
Objective 1            Objective 2
        │                  │
        ▼                  ▼
Key Results          Key Results
```

---

## OKR Evaluation Flow

```text
Employee
      │
      ▼
Create Evaluation ──────────► Snapshot dari Template
      │
      ▼
Input Actual Value
      │
      ▼
Calculate Achievement ──────► Achievement = (Actual / Target) × 100
      │
      ▼
Calculate Score ────────────► Score = (Weight × Achievement) / 100
      │
      ▼
Sum All Scores
      │
      ▼
Assign Rating ──────────────► Berdasarkan Final Score
      │
      ▼
Completed
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

## Review Workflow

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

| Code | Name | Description |
|------|------|-------------|
| DRAFT | Draft | Evaluasi baru dibuat |
| SUBMITTED | Submitted | Sudah disubmit untuk review |
| REVIEW | Under Review | Sedang direview |
| REVISION | Revision | Perlu revisi |
| APPROVED | Approved | Sudah disetujui |
| COMPLETED | Completed | Selesai |
| CANCELLED | Cancelled | Dibatalkan |

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

## OKR Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/okr/templates` | List templates |
| POST | `/api/v1/tenant/okr/templates` | Create template |
| GET | `/api/v1/tenant/okr/templates/:id` | Get template |
| PUT | `/api/v1/tenant/okr/templates/:id` | Update template |
| DELETE | `/api/v1/tenant/okr/templates/:id` | Delete template |
| POST | `/api/v1/tenant/okr/templates/:id/duplicate` | Duplicate template |

## OKR Objectives

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/okr/objectives` | List objectives |
| POST | `/api/v1/tenant/okr/objectives` | Create objective |
| GET | `/api/v1/tenant/okr/objectives/:id` | Get objective |
| PUT | `/api/v1/tenant/okr/objectives/:id` | Update objective |
| DELETE | `/api/v1/tenant/okr/objectives/:id` | Delete objective |

## OKR Key Results

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/okr/key-results` | List key results |
| POST | `/api/v1/tenant/okr/key-results` | Create key result |
| GET | `/api/v1/tenant/okr/key-results/:id` | Get key result |
| PUT | `/api/v1/tenant/okr/key-results/:id` | Update key result |
| DELETE | `/api/v1/tenant/okr/key-results/:id` | Delete key result |

## OKR Evaluations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/okr/evaluations` | List evaluations |
| POST | `/api/v1/tenant/okr/evaluations/snapshot` | Create evaluation with snapshot |
| GET | `/api/v1/tenant/okr/evaluations/:id` | Get evaluation |
| GET | `/api/v1/tenant/okr/evaluations/:id/full` | Get evaluation with details |
| PUT | `/api/v1/tenant/okr/evaluations/:id` | Update evaluation |
| DELETE | `/api/v1/tenant/okr/evaluations/:id` | Delete evaluation |

## OKR Evaluation Details

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/okr/evaluations/:id/details` | List details |
| PUT | `/api/v1/tenant/okr/evaluation-details/:id/actual` | Update actual value |
| PUT | `/api/v1/tenant/okr/evaluations/:id/actuals` | Bulk update actuals |

## OKR Workflow

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/okr/evaluations/:id/recalculate` | Recalculate score |
| POST | `/api/v1/tenant/okr/evaluations/:id/submit` | Submit evaluation |
| POST | `/api/v1/tenant/okr/evaluations/:id/approve` | Approve evaluation |
| POST | `/api/v1/tenant/okr/evaluations/:id/reject` | Reject evaluation |
| POST | `/api/v1/tenant/okr/evaluations/:id/complete` | Complete evaluation |

## OKR Progress

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/okr/progress` | Create progress |
| GET | `/api/v1/tenant/okr/evaluation-details/:id/progress` | List progress |
| GET | `/api/v1/tenant/okr/progress/:id` | Get progress |
| PUT | `/api/v1/tenant/okr/progress/:id` | Update progress |
| DELETE | `/api/v1/tenant/okr/progress/:id` | Delete progress |

## OKR Comments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/okr/comments` | Create comment |
| GET | `/api/v1/tenant/okr/evaluations/:id/comments` | List comments |
| PUT | `/api/v1/tenant/okr/comments/:id` | Update comment |
| DELETE | `/api/v1/tenant/okr/comments/:id` | Delete comment |

## OKR Attachments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/okr/attachments` | Upload attachment |
| GET | `/api/v1/tenant/okr/evaluation-details/:id/attachments` | List attachments |
| GET | `/api/v1/tenant/okr/attachments/:id` | Get attachment |
| DELETE | `/api/v1/tenant/okr/attachments/:id` | Delete attachment |

## OKR Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/okr/dashboard/employee/:id` | Employee dashboard |
| GET | `/api/v1/tenant/okr/dashboard/manager/:id` | Manager dashboard |
| GET | `/api/v1/tenant/okr/dashboard/hr` | HR dashboard |

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
- Snapshot Objective & Key Result saat Evaluation
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
| Phase 5 - Future Enhancement | ⏳ Pending | - | |
| Phase 6 - Seeder | ✅ Completed | - | Reuse dari KPI (tidak perlu seeder baru) |

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
- Templates, Objectives, Key Results, Evaluations, Progress, Comments, Attachments, Dashboard
