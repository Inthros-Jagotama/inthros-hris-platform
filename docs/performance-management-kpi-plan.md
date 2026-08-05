Berikut versi yang sudah disesuaikan dengan standar **UUID sebagai Primary Key** di seluruh tabel. Saya juga menyesuaikan foreign key agar semuanya menggunakan `uuid` sehingga konsisten dengan arsitektur modern (Laravel + PostgreSQL).

````markdown
# KPI Module Development Plan

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

# Phase 5 - Future Enhancement

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

---

# Final Database Structure

```text
performance_periods

performance_perspectives

performance_templates

performance_indicators

performance_evaluations

performance_evaluation_details

performance_progress

performance_comments

performance_attachments

performance_ratings

performance_indicator_formulas

performance_logs
```

---

# Design Principles

- Semua Primary Key menggunakan UUID
- Semua Foreign Key menggunakan UUID
- Organization = Position = 1 Employee
- Tidak menggunakan Assignment KPI
- KPI mengikuti Organization
- Snapshot KPI saat evaluasi
- Histori KPI tidak berubah walaupun Template berubah
- Soft Delete menggunakan deleted_at
- Mendukung Audit Trail
- Siap diintegrasikan dengan Competency Management, Career Path, Succession Planning, Talent Management, dan Workforce Management
````
