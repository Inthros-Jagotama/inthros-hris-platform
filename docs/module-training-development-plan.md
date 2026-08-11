# Training & Development Management — Operational Development Plan

> **Status dokumen**: diselaraskan dengan arsitektur project yang sudah ada
> (module SDK, SQL migration tenant, Central Approval Engine, RBAC tenant, FE tenant).
> Phase pengembangan dipisah menjadi track **BE** dan **FE** per fase (lihat §42).
> **Status implementasi**: `P0-BE` ✅ SELESAI (2026-08-11), `P0-FE` ✅ SELESAI (2026-08-11),
> `P1-BE` ✅ SELESAI (2026-08-11), `P1-FE` ✅ SELESAI (2026-08-11) — sisanya belum dikerjakan.
> Lihat **§47 Progress Tracker** untuk status tiap phase.
>
> **Konvensi**: setiap kali satu phase selesai diimplementasikan, status phase tsb
> **wajib di-update** di §47 + dicatat di §46 — dilakukan segera setelah selesai.

## 1. Tujuan

Module **Training & Development** digunakan untuk mengelola proses operasional pelatihan karyawan secara end-to-end:

- master katalog pelatihan;
- perencanaan dan pengajuan training;
- penyelenggaraan training;
- training **in-house maupun external**;
- peserta dan enrollment;
- attendance;
- materi;
- assessment;
- evaluasi;
- completion;
- sertifikat;
- training history;
- reporting.

Module ini **tidak menjadi module Intelligence**. Data training nantinya menjadi sumber data untuk:
- **Workforce Intelligence**: analisis kebutuhan dan kondisi workforce;
- **Career Intelligence**: analisis pengembangan individu, competency gap, career path, talent, dan succession.

---

## 2. Kondisi Existing (sesuai arsitektur saat ini)

### 2.1 Backend module sudah terdaftar

Module `backend/internal/modules/training/` **sudah terdaftar** di `backend/cmd/server/main.go`:

```text
ModuleName     = "Training & Development Management"
ModuleSlug     = "training"
ModuleVersion  = "1.0.0"
Priority       = 13 (tenant modules, setelah reimbursement, sebelum workforceintelligence)
IsCore         = false
DependsOn      = [employee, organization, setting]
```

`module.go` saat ini sudah meng-export `Menus` dengan route `/admin/training/*` — **perlu disinkronkan** ke route FE tenant `/training/*` (lihat §5 dan §37).

### 2.2 Database (migration `016_training.sql`)

7 tabel sudah ada (MySQL & PostgreSQL):

1. `training_categories`
2. `training_courses`
3. `training_sessions`
4. `training_participants`
5. `training_materials`
6. `training_evaluations`
7. `training_certificates`

Model GORM di `model.go` sinkron dengan migration (uuid.UUID `CHAR(36)`, `deleted_at`/`created_at`/`updated_at`, `BeforeCreate` hook).

### 2.3 API yang sudah ada

Route group existing: **`/api/v1/tenant/trainings`** (bukan `/training`) dengan CRUD lengkap:

```text
POST/GET/PUT/DELETE /trainings/categories[/:id]
POST/GET/PUT/DELETE /trainings/courses[/:id]
POST/GET/PUT/DELETE /trainings/sessions[/:id]  + PUT /trainings/sessions/:id/status
POST/GET/PUT/DELETE /trainings/participants[/:id]
POST/GET/PUT/DELETE /trainings/materials[/:id]
POST/GET/PUT/DELETE /trainings/evaluations[/:id]
POST/GET/PUT/DELETE /trainings/certificates[/:id]
```

### 2.4 Frontend tenant

- Route FE: `/training` → `views/modules/Training.vue` (masih placeholder "Coming soon").
- Sidebar: menu **Training** masuk section **Talent** (`Sidebar.vue` → `talentItems`, permission `training.view`).
- Locale: key `training` sudah ada di `en.json`/`id.json` (title, description, categories, courses, sessions, participants, evaluations, certificates).
- Gate module: `activeModules.hasModule('training')`.

### 2.5 Keterbatasan existing (tetap relevan)

- `training_courses.external_vendor` masih text;
- trainer masih `trainer_name`;
- session hanya `start_date`/`end_date`;
- attendance hanya `training_participants.attendance_status`;
- score satu field;
- belum ada enrollment/request workflow, assessment detail, learning objective, relasi course–competency, prerequisite, provider master, training plan/need/request, mandatory training, cost detail, evaluation form, effectiveness, certification master, training history kaya.

---

# 3. Prinsip Arsitektur

## 3.1 Operational Module

Training bertanggung jawab terhadap transaksi:

```text
Training Catalog
      ↓
Training Planning
      ↓
Training Request
      ↓
Training Session
      ↓
Enrollment
      ↓
Attendance
      ↓
Assessment
      ↓
Completion
      ↓
Evaluation
      ↓
Certificate
```

## 3.2 Intelligence

Tidak membuat module Intelligence baru.

```text
Training
   │
   ├── Workforce Intelligence
   │      └── workforce/training analysis
   │
   └── Career Intelligence
          └── individual career/development analysis
```

Training hanya menyediakan data operasional.

## 3.3 Prinsip Adaptasi ke Arsitektur Existing

1. **Migration SQL baru, bukan edit file lama.** Enhancement database dibuat sebagai file migration baru
   `backend/internal/pkg/migrator/migrations/tenant/postgres/088_training_core.sql` dst.
   (daftar lengkap: §6 ↔ §42).
   File `016_training.sql` dan model GORM lama tidak diubah strukturnya secara retrofit; perubahan
   kolom (mis. `training_courses`) dilakukan lewat `ALTER TABLE` di migration baru. Versi selanjutnya:
   `088_*`, `089_*`, dst (latest saat dokumen ini ditulis = `087_employeemovement_cancellation`).
2. **Model GORM tetap satu-satunya sumber truth kode** — setiap tabel baru wajib punya struct model
   + `BeforeCreate` hook + `AutoMigrate` di `module.go` `Migrate()`, mengikuti pola entitas existing
   di `model.go` (uuid.UUID `CHAR(36)`, index `idx_trn_*`).
3. **Approval memakai Central Approval Engine**, bukan engine baru. Pola narrow-interface + adapter
   yang sama dengan leave/reimbursement/attendance/employeemovement (§15 dan §45).
4. **RBAC lewat `seed_rbac.go`** — permission baru ditambahkan ke `tenantRBACResources()`
   (format `resource.action`, idempotent) + `Permissions()` di `module.go` (§36).
5. **Route FE memakai pola hub + sub-halaman** seperti Workforce Intelligence / Career Intelligence
   (`/training` hub, `/training/courses`, dst) — bukan menu tree server (§37).
6. **Seeder master data** mengikuti pola existing: SQL seeder migration (pola `033_*`–`053_*`)
   atau via `Seed()` module; permission di-seed lewat `seed_rbac.go` (§40).

---

# 4. In-House dan External Training

Training harus mendukung dua bentuk penyelenggaraan:

```text
IN_HOUSE
EXTERNAL
```

## 4.1 In-House

```text
Course:        Leadership Basic
Delivery:      IN_HOUSE
Trainer:       Employee internal
Location:      Training Room A
```

## 4.2 External

```text
Course:        Leadership Advanced
Delivery:      EXTERNAL
Provider:      ABC Training Institute
Trainer:       External Trainer
Location:      ABC Training Center
```

Flow database tetap sama (Course → Session → Participant → Attendance → Assessment → Evaluation → Certificate).
Yang berubah adalah provider, trainer, lokasi, biaya, dan informasi penyelenggara.

**Adaptasi arsitektur**: pisahkan dua konsep di model `training_sessions`:

```text
provider_type   IN_HOUSE | EXTERNAL        → penyelenggara
delivery_mode   ONSITE | ONLINE | HYBRID | SELF_PACED   → mode penyampaian
```

`IN_HOUSE`/`EXTERNAL` bukan mode delivery — mereka tipe penyelenggara. `ONLINE`/`HYBRID`/`ONSITE`
adalah mode delivery. `training_courses.delivery_type` menjadi nilai *default/preferred* course
(`IN_HOUSE`/`EXTERNAL`/`BOTH`), sedangkan nilai aktual dipecah per session.

---

# 5. Target Menu (diselaraskan dengan FE tenant)

```text
Training & Development  (/training  → hub ber-card, pola Workforce/Career Intelligence)
│
├── Dashboard           (di dalam halaman hub)
├── Catalog             (/training/catalog)
│   ├── Categories
│   ├── Courses
│   ├── Objectives      (sub-resource course)
│   ├── Competencies    (sub-resource course)
│   ├── Prerequisites   (sub-resource course)
│   ├── Providers       (master penyelenggara external)
│   └── Trainers        (master trainer internal/external)
├── Planning            (/training/planning)
│   ├── Training Plans
│   ├── Training Needs
│   └── Training Requests
├── Sessions            (/training/sessions)
│   ├── Scheduled
│   ├── Ongoing
│   └── Completed
├── Participants        (/training/sessions/:id/participants)
├── Attendance          (tab pada detail session)
├── Assessments         (tab pada detail session)
├── Evaluations         (tab pada detail session)
├── Certifications      (/training/certificates)
├── Training History    (/training/history)
└── Reports             (/training/reports)
```

Pola route FE: sub-halaman didaftarkan di `frontend/tenant/src/router/index.js` dengan
`meta: { module: 'training', backRoute: '/training', ... }`, dan menu sidebar **tetap satu item**
"Training" di section Talent (pola hub Movements & Contracts).

---

# 6. Database Enhancement — Strategi Migration

Semua ID mengikuti pola existing:

```sql
CHAR(36)
```

dan tetap menggunakan:

```sql
deleted_at
created_at
updated_at
```

sesuai pola file existing. Nama index mengikuti pola `idx_trn_*`.

**Strategi**: setiap fase = satu file migration baru, selaras dengan track BE di §42:

```text
088_training_core.sql       → P0-BE: providers, trainers, session_trainers, attendances,
                              assessments, assessment_results
                              + ALTER courses, sessions, participants, materials
089_training_planning.sql   → P1-BE: plans, plan_items, needs, requests, course objectives,
                              course competencies, course prerequisites, mandatories,
                              session_costs, documents
090_training_advanced.sql   → P2-BE: evaluation_forms/questions/answers, effectiveness,
                              certifications + ALTER certificates
```

Catatan: `ALTER TABLE` yang merubah kolom lama (mis. menambah `course_type`,
`delivery_type`, `is_mandatory` di `training_courses`) memakai `ADD COLUMN IF NOT EXISTS`
dengan default, agar data existing tidak rusak. Kolom `external_vendor` dan `trainer_name`
**di-deprecate** (tidak di-drop di fase awal) untuk kompatibilitas.

---

## 6.1 `training_categories`

**Existing sudah cukup.**

```text
id
code
name
description
is_active
deleted_at
created_at
updated_at
```

Enhancement opsional:

```text
sort_order
```

---

# 7. `training_courses`

Existing:

```text
id
category_id
code
name
description
duration_hour
min_score
cost
is_certified
external_vendor
is_active
deleted_at
created_at
updated_at
```

### Perubahan

`external_vendor` di-deprecate dari master course — vendor adalah karakteristik session/provider,
bukan course (satu course bisa diselenggarakan in-house di session A dan external provider di session B).

### Tambahkan (via `ALTER TABLE` di migration baru)

```text
course_type      TECHNICAL | SOFT_SKILL | COMPLIANCE | MANAGEMENT | CERTIFICATION | OTHER
delivery_type    IN_HOUSE | EXTERNAL | BOTH          (default/preferred)
is_mandatory     boolean
```

---

# 8. `training_course_objectives`

Menyimpan learning objectives.

```text
id
course_id
objective
sort_order
deleted_at
created_at
updated_at
```

---

# 9. `training_course_competencies`

Menghubungkan course dengan competency. `competency_id` mengacu ke tabel kompetensi dari
**module competency** (`competencies`), bukan tabel baru.

```text
id
course_id
competency_id     → competencies.id
target_level
deleted_at
created_at
updated_at
```

Relasi ini menjadi salah satu sumber data untuk Career Intelligence.

---

# 10. `training_course_prerequisites`

```text
id
course_id
prerequisite_type    COURSE | COMPETENCY | CERTIFICATION | EXPERIENCE
prerequisite_id
is_required
deleted_at
created_at
updated_at
```

---

# 11. `training_providers`

```text
id
code
name
type               INTERNAL | EXTERNAL
contact_name
email
phone
address
website
is_active
deleted_at
created_at
updated_at
```

In-house: tidak perlu record provider (`provider_id = NULL` di session).

---

# 12. `training_trainers`

Trainer tidak lagi hanya `trainer_name`.

```text
id
type               INTERNAL | EXTERNAL
employee_id        → employees.id (INTERNAL, nullable)
provider_id        → training_providers.id (EXTERNAL, nullable)
name
email
phone
bio
is_active
deleted_at
created_at
updated_at
```

---

# 13. `training_session_trainers`

Satu session dapat memiliki lebih dari satu trainer.

```text
id
session_id
trainer_id
role               MAIN | ASSISTANT
deleted_at
created_at
updated_at
```

---

# 14. `training_sessions`

Existing:

```text
id
course_id
session_code
trainer_name       → deprecated, ganti ke training_trainers + training_session_trainers
location
start_date
end_date
max_quota
status
```

### Enhancement (via `ALTER TABLE`)

```text
provider_type        IN_HOUSE | EXTERNAL
delivery_mode        ONSITE | ONLINE | HYBRID | SELF_PACED
provider_id          → training_providers.id (wajib jika EXTERNAL)
start_datetime       (pertahankan start_date/end_date sebagai proyeksi, datetime untuk presisi)
end_datetime
meeting_url
registration_deadline
```

### Status session — selaraskan dengan enum Go existing

Existing di `model.go`: `SCHEDULED | IN_PROGRESS | COMPLETED | CANCELLED`.

Tambah (via migration + const baru) menjadi:

```text
DRAFT
SCHEDULED
REGISTRATION_OPEN
FULL
IN_PROGRESS          (existing ONGOING → pakai IN_PROGRESS sesuai konstanta Go yang ada)
COMPLETED
CANCELLED
```

`PUT /trainings/sessions/:id/status` yang sudah ada dipakai sebagai transisi status.

---

# 15. Training Request — pakai Central Approval Engine

Tambahkan `training_requests`:

```text
id
employee_id
course_id
requested_date
reason
priority            LOW | MEDIUM | HIGH | URGENT
status              DRAFT | SUBMITTED | PENDING_APPROVAL | APPROVED | REJECTED | CANCELLED
approval_instance_id → approval_instances.id (nullable, diisi saat submit)
deleted_at
created_at
updated_at
```

Status `PENDING_APPROVAL` ditambahkan (bukan hanya `SUBMITTED`) mengikuti pola
leave/reimbursement — setelah `CreateApprovalInstance` berhasil, status request = `PENDING_APPROVAL`.

**Approval menggunakan Central Approval Engine** (module `approval`), bukan approval engine baru:

- Modul flow approval: `training_request` (pola penamaan module approval: `leave`, `payroll`,
  `reimbursement`, `attendance`, `employeemovement`, `performance_kpi_target`, ...).
- Service training mengekspos narrow interface `ApprovalEngine` yang sama dengan modul lain
  (`CreateApprovalInstance`, `GetApprovalInstanceStatus`, `GetActiveFlowIDForModule`,
  `CancelApprovalInstance`) — lihat §45 wiring di `main.go`.
- Status callback didorong (push) oleh approval module: `RegisterStatusHandler("training_request", ...)`
  → handler mengubah status request & auto-enroll peserta saat `APPROVED`.

Flow:

```text
Employee / Manager
      ↓
Training Request
      ↓
Central Approval (flow module = training_request)
      ↓
APPROVED → register ke session (training_participants)
REJECTED → tidak menjadi participant aktif
```

---

# 16. Training Plan

Tambahkan:

```text
training_plans
training_plan_items
```

## `training_plans`

```text
id
code
name
year
description
status               DRAFT | ACTIVE | ARCHIVED
deleted_at
created_at
updated_at
```

## `training_plan_items`

```text
id
training_plan_id
course_id
target_date
target_participants
estimated_cost
priority
deleted_at
created_at
updated_at
```

---

# 17. Training Needs

`training_needs` — kebutuhan training secara operasional (bukan Intelligence).

```text
id
employee_id          (nullable untuk kebutuhan org/posisi)
organization_id      → organizations.id
position_id          → positions.id
course_id
reason
priority
source_type          MANUAL | PERFORMANCE | COMPETENCY | CAREER | SUCCESSION | COMPLIANCE | WORKFORCE
source_id
status               OPEN | PLANNED | FULFILLED | CANCELLED
deleted_at
created_at
updated_at
```

---

# 18. Enrollment / Participant

Existing `training_participants`:

```text
id
session_id
employee_id
attendance_status
score
completed_at
```

### Ubah konsep menjadi enrollment/participant record (via `ALTER TABLE`)

```text
registration_status    NOMINATED | REQUESTED | APPROVED | REGISTERED | WAITLISTED | CANCELLED
registered_at
approved_at
completion_status      NOT_STARTED | IN_PROGRESS | COMPLETED | FAILED
completion_date
final_score
passed
remarks
```

Tambahkan unique constraint `(session_id, employee_id)` untuk mencegah duplikasi peserta.

Catatan: kolom `attendance_status` dipertahankan sementara (compatibility), tetapi sumber detail
attendance ke depan adalah `training_attendances` (§19).

---

# 19. Training Attendance

Untuk training multi-day/session, tambahkan `training_attendances`:

```text
id
participant_id
attendance_date
check_in
check_out
status               PRESENT | ABSENT | LATE | EXCUSED
remarks
deleted_at
created_at
updated_at
```

`LATE` ditambahkan ke enum `AttendanceStatus` Go yang existing
(`PRESENT | ABSENT | EXCUSED`).

---

# 20. Training Materials

Existing sudah cukup sebagai basic implementation. Enhancement opsional:

```text
description
is_required
available_from
```

`file_url` diisi hasil upload via endpoint tenant `/api/v1/tenant/uploads`
(pola attachment yang sudah ada, file diserve publik dari `/uploads`).

---

# 21. Training Assessments

Tambahkan:

```text
training_assessments
training_assessment_results
```

## `training_assessments`

```text
id
session_id
name
type               PRE_TEST | POST_TEST | FINAL | PRACTICAL | OTHER
max_score
passing_score
attempt_limit
is_required
deleted_at
created_at
updated_at
```

## `training_assessment_results`

```text
id
assessment_id
participant_id
score
passed
attempt
completed_at
deleted_at
created_at
updated_at
```

---

# 22. Training Evaluation

Existing `training_evaluations` (rating + feedback) dipertahankan sementara. Untuk versi lengkap:

```text
training_evaluation_forms
training_evaluation_questions
training_evaluation_answers
```

## Evaluation Form

```text
id
session_id
name
is_active
```

## Questions

```text
id
form_id
question
question_type      RATING | TEXT | SINGLE_CHOICE | MULTIPLE_CHOICE
sort_order
is_required
```

## Answers

```text
id
question_id
participant_id
answer
```

---

# 23. Training Effectiveness

`training_effectiveness_assessments` — mengukur dampak setelah training selesai.

```text
id
participant_id
assessment_date
assessor_employee_id
before_score
after_score
effectiveness_score
remarks
deleted_at
created_at
updated_at
```

Dapat dilakukan setelah periode tertentu (mis. 30/60/90 hari).

---

# 24. Certification

Existing `training_certificates` menyimpan `participant_id`, `certificate_no`,
`issued_date`, `expiry_date`.

Tambahkan master `training_certifications`:

```text
id
code
name
issuing_body
validity_period_month
renewal_required
is_active
deleted_at
created_at
updated_at
```

Tambahkan pada `training_certificates`:

```text
certification_id        → training_certifications.id (nullable)
certificate_file_url
```

---

# 25. Mandatory Training

`training_mandatories` — training wajib berdasarkan target.

```text
id
course_id
organization_id        (nullable)
position_id            (nullable)
employment_status_id   → employment_statuses.id (nullable)
due_days
validity_period_month
is_active
deleted_at
created_at
updated_at
```

Target boleh kosong sebagian — record berlaku untuk kombinasi yang terisi.
Dashboard: Required / Completed / Pending / Overdue.

---

# 26. Training Cost

`training_session_costs` — biaya aktual per session.

```text
id
session_id
cost_type            TRAINER | PROVIDER | VENUE | MATERIAL | CERTIFICATION | TRAVEL | ACCOMMODATION | OTHER
description
amount
deleted_at
created_at
updated_at
```

Dihitung: Total Cost, Cost per Participant, Cost per Training Hour.

---

# 27. Training Documents

```text
training_documents
```

```text
id
session_id
document_type        PROPOSAL | QUOTATION | ATTENDANCE_SHEET | INVOICE | CONTRACT | TRAINING_REPORT | OTHER
file_name
file_url
uploaded_by
deleted_at
created_at
updated_at
```

Upload via `/api/v1/tenant/uploads` (pola attachment existing).

---

# 28. Training History

Tidak perlu tabel `employee_training_histories` — histori dibangun dari:

```text
training_participants
+ training_sessions
+ training_courses
+ training_assessments
+ training_certificates
```

Endpoint read-only `GET /trainings/history?employee_id=...` cukup men-join tabel di atas.

---

# 29. Business Flow

## 29.1 In-House

```text
Training Need → Training Plan → Course → Create Session → provider_type=IN_HOUSE
→ Internal Trainer → Open Registration → Participant → Attendance
→ Assessment → Completion → Evaluation → Certificate
```

## 29.2 External

```text
Training Need → Training Plan → Course → Create Session → provider_type=EXTERNAL
→ Provider → External Trainer → Open Registration → Participant → Attendance
→ Assessment → Completion → Evaluation → Certificate
```

---

# 30. Training Request Flow

```text
Employee / Manager → Training Request → Central Approval (training_request)
→ APPROVED → register ke session (training_participants)
REJECTED → tidak menjadi participant aktif
```

---

# 31. Training Planning Flow

```text
Training Need → Training Plan → Course → Session → Participant
```

Planning tidak wajib untuk session ad-hoc.

---

# 32. Validation

## Course

- category harus aktif;
- code unik;
- course tidak boleh dipakai jika inactive.

## Session

- course harus aktif;
- start datetime <= end datetime;
- registration deadline < start datetime;
- quota > 0;
- provider wajib untuk EXTERNAL, tidak wajib untuk IN_HOUSE;
- trainer harus sesuai type (INTERNAL → employee, EXTERNAL → provider);
- session code unik.

## Participant

- employee harus aktif;
- tidak boleh duplicate pada session (unique constraint);
- quota dicek (`CountParticipantsBySession` sudah ada di repository);
- prerequisite terpenuhi;
- mandatory training sesuai target.

## Assessment

- score <= max score;
- passing score <= max score;
- assessment required diselesaikan sebelum completion.

## Certificate

- certificate number unik;
- issued date valid;
- expiry date >= issued date;
- hanya untuk participant yang memenuhi completion requirement.

---

# 33. Status Lifecycle (selaras enum Go)

## Course

```text
ACTIVE | INACTIVE
```

## Session

```text
DRAFT | SCHEDULED | REGISTRATION_OPEN | FULL | IN_PROGRESS | COMPLETED | CANCELLED
```

## Participant

```text
registration: NOMINATED | REQUESTED | APPROVED | REGISTERED | WAITLISTED | CANCELLED
completion:   NOT_STARTED | IN_PROGRESS | COMPLETED | FAILED
attendance:   PRESENT | ABSENT | LATE | EXCUSED
```

## Training Request

```text
DRAFT | SUBMITTED | PENDING_APPROVAL | APPROVED | REJECTED | CANCELLED
```

---

# 34. API Plan (diselaraskan dengan route group existing `/trainings`)

Route group tetap **`/api/v1/tenant/trainings`** (bukan `/training`) agar tidak memutus endpoint
yang sudah ada. Gaya route mengikuti existing: **flat + query param filter** untuk resource anak
(pola `ListParticipants?session_id=`, `ListMaterials?session_id=`), dengan beberapa sub-resource
nested hanya jika alami (trainer session, attendance, assessment session).

> Konstrain Gin: route statis dideklarasikan SEBELUM route dinamis `/:id`
> (contoh: `GET /trainings/history` dan `GET /trainings/sessions/summary` sebelum
> `GET /trainings/sessions/:id`).

## Catalog (existing + enhancement)

```http
GET    /trainings/categories
POST   /trainings/categories
GET    /trainings/categories/{id}
PUT    /trainings/categories/{id}
DELETE /trainings/categories/{id}

GET    /trainings/courses
POST   /trainings/courses
GET    /trainings/courses/{id}
PUT    /trainings/courses/{id}
DELETE /trainings/courses/{id}

GET    /trainings/courses/{id}/objectives        (baru)
POST   /trainings/courses/{id}/objectives        (baru)
PUT    /trainings/course-objectives/{id}         (baru)
DELETE /trainings/course-objectives/{id}         (baru)

GET    /trainings/courses/{id}/competencies      (baru)
POST   /trainings/courses/{id}/competencies      (baru)
DELETE /trainings/course-competencies/{id}       (baru)

GET    /trainings/courses/{id}/prerequisites     (baru)
POST   /trainings/courses/{id}/prerequisites     (baru)
DELETE /trainings/course-prerequisites/{id}      (baru)

GET    /trainings/providers                      (baru)
POST   /trainings/providers                      (baru)
GET    /trainings/providers/{id}                 (baru)
PUT    /trainings/providers/{id}                 (baru)
DELETE /trainings/providers/{id}                 (baru)

GET    /trainings/trainers                       (baru)
POST   /trainings/trainers                       (baru)
GET    /trainings/trainers/{id}                  (baru)
PUT    /trainings/trainers/{id}                  (baru)
DELETE /trainings/trainers/{id}                  (baru)
```

## Planning (baru)

```http
GET    /trainings/plans
POST   /trainings/plans
GET    /trainings/plans/{id}
PUT    /trainings/plans/{id}
DELETE /trainings/plans/{id}

GET    /trainings/plans/{id}/items
POST   /trainings/plans/{id}/items
PUT    /trainings/plan-items/{id}
DELETE /trainings/plan-items/{id}

GET    /trainings/needs
POST   /trainings/needs
PUT    /trainings/needs/{id}

GET    /trainings/requests
POST   /trainings/requests
GET    /trainings/requests/{id}
POST   /trainings/requests/{id}/submit           (buat approval instance)
POST   /trainings/requests/{id}/cancel
```

## Sessions (existing + enhancement)

```http
GET    /trainings/sessions
POST   /trainings/sessions
GET    /trainings/sessions/{id}
PUT    /trainings/sessions/{id}
DELETE /trainings/sessions/{id}
PUT    /trainings/sessions/{id}/status           (existing — transisi status)

POST   /trainings/sessions/{id}/publish          (baru — DRAFT → REGISTRATION_OPEN)
POST   /trainings/sessions/{id}/cancel           (baru)

GET    /trainings/sessions/{id}/trainers         (baru)
POST   /trainings/sessions/{id}/trainers         (baru)
DELETE /trainings/session-trainers/{id}          (baru)
```

## Participants (existing + enrollment)

```http
GET    /trainings/participants?session_id=&employee_id=   (existing)
POST   /trainings/participants                            (existing — register)
GET    /trainings/participants/{id}                       (existing)
PUT    /trainings/participants/{id}                       (existing)
DELETE /trainings/participants/{id}                       (existing)
```

## Attendance (baru)

```http
GET    /trainings/sessions/{id}/attendance
POST   /trainings/sessions/{id}/attendance
PUT    /trainings/attendances/{id}
```

## Assessment (baru)

```http
GET    /trainings/sessions/{id}/assessments
POST   /trainings/sessions/{id}/assessments
POST   /trainings/assessments/{id}/results
```

## Cost & Documents (baru)

```http
GET    /trainings/sessions/{id}/costs
POST   /trainings/sessions/{id}/costs
PUT    /trainings/session-costs/{id}
DELETE /trainings/session-costs/{id}

GET    /trainings/sessions/{id}/documents
POST   /trainings/sessions/{id}/documents
DELETE /trainings/documents/{id}
```

## Evaluation (baru)

```http
GET    /trainings/sessions/{id}/evaluation-form
POST   /trainings/sessions/{id}/evaluation-form
GET    /trainings/evaluations                       (existing)
POST   /trainings/evaluations                       (existing)
GET    /trainings/evaluations/{id}                  (existing)
PUT    /trainings/evaluations/{id}                  (existing)
DELETE /trainings/evaluations/{id}                  (existing)
```

## Certificate (existing + enhancement)

```http
GET    /trainings/certificates                      (existing)
GET    /trainings/certificates/{id}                 (existing)
POST   /trainings/certificates                      (existing)
POST   /trainings/participants/{id}/certificate     (baru — generate dari completion)
```

## Reports & History (baru, route statis sebelum `/:id`)

```http
GET    /trainings/history?employee_id=...
GET    /trainings/reports/participation
GET    /trainings/reports/cost
GET    /trainings/reports/compliance
GET    /trainings/reports/dashboard
```

---

# 35. Backend Structure

Mengikuti pola repository/service existing project. **Tidak membuat package baru** — semua entitas
tetap di `backend/internal/modules/training/` dengan file yang sama
(`model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`, `module.go`).

```text
training/
├── model.go        → TrainingCategory, Course, CourseObjective, CourseCompetency,
│                     CoursePrerequisite, Provider, Trainer, Session, SessionTrainer,
│                     Plan, PlanItem, Need, Request, Participant, Attendance, Material,
│                     Assessment, AssessmentResult, EvaluationForm, EvaluationQuestion,
│                     EvaluationAnswer, EffectivenessAssessment, Certification, Certificate,
│                     Mandatory, SessionCost, Document
├── dto.go
├── repository.go
├── service.go      → ApprovalEngine narrow interface (submit/cancel request)
├── handler.go
├── routes.go
├── module.go       → AutoMigrate daftar model baru + Permissions() + Menus (route FE)
└── *_test.go       → service/repository/handler/approval_integration_test.go
```

`module.go` `DependsOn` diperluas:

```text
DependsOn = [employee, organization, setting, competency, approval]
```

---

# 36. Authorization

## Permission di `module.go` `Permissions()`

```text
training.view
training.create
training.update
training.delete
training.enroll

training.course.manage
training.session.manage
training.participant.manage
training.attendance.manage
training.assessment.manage
training.evaluation.manage
training.certificate.manage
training.plan.manage
training.request.create
training.request.approve
training.report.view
```

## Seed RBAC (`backend/internal/pkg/tenantseed/seed_rbac.go`)

`tenantRBACResources()` sudah punya resource `training` dengan
`[view, create, update, delete]` (idempotent, deterministik). Perlu ditambahkan action
granular agar permission di atas tersedia di UI RBAC:

```go
{"training", []string{"view", "create", "update", "delete", "enroll",
    "course.manage", "session.manage", "participant.manage", "attendance.manage",
    "assessment.manage", "evaluation.manage", "certificate.manage", "plan.manage",
    "request.create", "request.approve", "report.view"}},
```

> Catatan: `training.request.approve` secara teknis dijalankan oleh Central Approval —
> permission ini mengatur siapa yang boleh *melihat/mentrigger* flow, sedangkan approval
> task approval module yang mengeksekusi.

Approval tetap memakai **Central Approval Engine** — tidak ada approval engine baru di Training.
Flow module baru yang didaftarkan ke approval: `training_request` (termasuk label di
locale `approval.module_names` FE).

---

# 37. Frontend Plan

Pola halaman mengikuti modul operasional yang sudah ada (Leave, Attendance, Movement):

## Hub page `views/modules/Training.vue`

Ganti placeholder "Coming soon" dengan halaman hub ber-card (pola `WorkforceIntelligence.vue` /
`CareerIntelligence.vue` / `EmployeeMovementReports.vue`):

```text
Stat ringkas: Upcoming / Ongoing / Completed training, Total Participants,
Completion Rate, Training Hours, Training Cost, Mandatory Compliance.
Card navigasi: Catalog, Planning, Sessions, Certificates, History, Reports.
```

## Sub-halaman `views/modules/training/`

```text
frontend/tenant/src/views/modules/training/
├── TrainingCourses.vue        (daftar course + detail card: category, duration,
│                               delivery type, provider, certification, prerequisite)
├── TrainingCategories.vue
├── TrainingProviders.vue      (master provider — external)
├── TrainingTrainers.vue       (master trainer — internal/external)
├── TrainingPlanning.vue       (plans + plan items)
├── TrainingNeeds.vue          (training needs + mandatory compliance)
├── TrainingRequests.vue       (request + submit/cancel via approval flow)
├── TrainingSessions.vue       (list: course, delivery, provider, trainer, location,
│                               schedule, quota, participants, status)
├── TrainingSessionDetail.vue  (tabs: Overview, Participants, Attendance, Materials,
│                               Assessment, Evaluation, Certificate, Documents, Cost)
├── TrainingParticipants.vue
├── TrainingCertificates.vue
├── TrainingHistory.vue
└── TrainingReports.vue
```

## Router (`frontend/tenant/src/router/index.js`)

```text
/training                      → Training.vue (hub)
/training/courses              → module: 'training', backRoute: '/training'
/training/categories           → module: 'training'
/training/providers            → module: 'training'
/training/trainers             → module: 'training'
/training/planning             → module: 'training'
/training/requests             → module: 'training'
/training/needs                → module: 'training'
/training/sessions             → module: 'training'
/training/sessions/:id         → detail session
/training/certificates         → module: 'training'
/training/history              → module: 'training'
/training/reports              → module: 'training'
```

## Sidebar (`Sidebar.vue`)

Tetap satu item **Training** di section Talent (`talentItems`, `pi pi-book`,
`/training`, permission `training.view`) — sub-halaman tidak membuat menu tambahan
(pola Movements & Contracts hub).

## Locale (`en.json` / `id.json`)

Perluas key `training`: dashboard, catalog, planning, plans, needs, requests,
providers, trainers, objectives, competencies, prerequisites, attendance,
assessments, evaluations, certifications, history, reports, dan semua status enum.

## Gate module

Router guard `meta.module: 'training'` + `activeModules.hasModule('training')` — sudah mengikuti
pola yang ada.

---

# 38. Reporting

Minimal report (endpoint read-only, `training.report.view`):

### Training Participation

```text
Employee | Organization | Course | Session | Status | Attendance | Score | Completion
```

### Training Cost

```text
Course | Session | Provider | Total Cost | Participant Count | Cost per Participant
```

### Training Compliance

```text
Organization | Position | Employee | Mandatory Course | Due Date | Completion | Status
```

### Training History

```text
Employee | Course | Date | Score | Completion | Certificate
```

---

# 39. Integration dengan Modul Lain

## Employee Management

Training mengambil:

```text
employee_id
organization_id
position_id
employment_status_id
```

## Competency

`training_course_competencies` menghubungkan course ke `competencies.id` (module competency).
Training completion menjadi evidence pengembangan kompetensi.

## Performance / KPI / OKR

Training dapat menjadi response terhadap performance/competency gap
(`training_needs.source_type = PERFORMANCE`).

## Career Intelligence

Training dipakai untuk: Career Path, Career Eligibility, Competency Gap,
Development Recommendation, Talent, Succession — via relasi course–competency.

## Workforce Intelligence

Training dipakai untuk: Workforce Planning, Workforce Gap, Mandatory Training Compliance,
Organization Training Analysis, Training Cost Analysis.

## Approval (Central)

Training Request → flow module `training_request` → callback `RegisterStatusHandler`.

Tidak ada module Intelligence baru.

---

# 40. Seeders

Mengikuti pola existing:

1. **Permission** → `seed_rbac.go` (`tenantRBACResources()`), idempotent.
2. **Master data** → SQL seeder migration (pola `033_*`–`053_*`) atau `Seed()` di `module.go`.
   Kandidat seeder:

```text
training_categories   (Technical, Soft Skill, Leadership, Management, Compliance, Safety, Certification)
course_type / delivery_type / assessment type / attendance status / participant status /
request status / session status / provider type / trainer type  (diseed sebagai data enum master bila ada tabelnya)
```

Seeder category contoh mengikuti pola seeder existing project.

---

# 41. Testing Plan

Test dijalankan dalam fase fiturnya: unit/integration test di track **BE** (P0-BE, P1-BE,
P2-BE) dan validasi UI di track **FE** (pola manual/end-to-end yang sudah ada di modul lain).

## Unit Test (pola existing `service_test.go`, `repository_test.go`, `handler_test.go`)

- course validation; session validation; quota calculation;
- participant registration; duplicate participant (unique constraint);
- prerequisite validation; assessment score; passing score;
- completion; certificate eligibility; mandatory training;
- provider validation; trainer validation.

## Approval Integration Test (pola `approval_integration_test.go` leave/reimbursement/attendance)

- `fakeApprovalEngine` mencatat `CreateApprovalInstance` call;
- submit request → `approval_instance_id` terisi + status `PENDING_APPROVAL`;
- callback approved → status `APPROVED` + auto-enroll participant;
- callback rejected → tidak menjadi participant.

## Feature Test

```text
Create Course → Create Session → Register Employee → Attendance → Assessment
→ Complete → Evaluation → Certificate

External Course → Provider → External Trainer → Session → Participant → Completion

Training Request → Submit → Central Approval → Approved → Enrollment
```

---

# 42. Development Priority — Backend & Frontend (phase dipisah)

Setiap fase memiliki **dua track terpisah**: **BE** (backend) dan **FE** (frontend).
Eksekusi mengikuti urutan per fase — **BE dulu, FE menyusul** (vertical slice per fitur:
FE satu fitur baru dimulai setelah endpoint BE fitur tersebut tersedia):

```text
P0-BE → P0-FE → P1-BE → P1-FE → P2-BE → P2-FE
```

Migration SQL, model GORM, permission seeding, dan wiring `main.go` adalah bagian dari track **BE**.
Router, sidebar, locale, dan halaman Vue adalah bagian dari track **FE**.

## P0 — Core Operational (migration 088)

### P0-BE (Backend) — ✅ SELESAI (2026-08-11)

1. Migration `088_training_core.sql`:
   - tabel baru: `training_providers`, `training_trainers`, `training_session_trainers`;
   - `ALTER training_courses` (course_type, delivery_type, is_mandatory);
   - `ALTER training_sessions` (provider_type, delivery_mode, provider_id,
     start_datetime, end_datetime, meeting_url, registration_deadline);
   - `ALTER training_participants` (registration_status, registered_at, approved_at,
     completion_status, completion_date, final_score, passed, remarks,
     unique(session_id, employee_id));
   - `ALTER training_materials` (description, is_required, available_from);
   - tabel baru: `training_attendances`, `training_assessments`,
     `training_assessment_results`.

   > Catatan: attendance & assessment API selesai di P0-BE, tetapi tab UI-nya baru
   > disajikan di P1-FE (deferral sengaja — P0-FE fokus ke catalog/session/participant).
   > `training_session_costs` & `training_documents` dibuat di migration 089 (P1-BE)
   > bersama API & tab Cost/Documents-nya (vertical slice utuh).
2. Model GORM + `AutoMigrate` di `module.go` untuk semua tabel baru
   (uuid.UUID CHAR(36), index `idx_trn_*`, `BeforeCreate` hook) + deprecate
   `external_vendor`/`trainer_name`.
3. DTO + repository + service + handler + routes:
   - CRUD provider & trainer;
   - CRUD course/session diperluas (field baru + validasi provider_type/delivery_mode,
     session code unik);
   - enrollment participant (registration_status, quota check, duplicate check);
   - attendance detail (`training_attendances`);
   - assessment + assessment results (validasi score <= max_score, passing);
   - materials enhancement (description, is_required, available_from).
4. `Permissions()` di `module.go` diperluas + `seed_rbac.go`
   (`tenantRBACResources()`) menambah action granular training (§36).
5. `Menus` di `module.go` disinkronkan ke route FE `/training/*` (§5).
6. Unit test: course/session/provider/trainer validation, quota, duplicate participant,
   assessment score, completion, certificate eligibility.
7. Feature test BE: Create Course → Session → Register → Attendance → Assessment
   → Complete → Certificate.

### P0-FE (Frontend) — ✅ SELESAI (2026-08-11)

1. Hub page `Training.vue` (ganti placeholder) + stat ringkas
   (Upcoming/Ongoing/Completed, Total Participants, Completion Rate, Training Hours,
   Training Cost). Stat awal memakai endpoint list existing (count dari sessions &
   participants); dashboard analitik penuh baru di P2-FE (`/trainings/reports/dashboard`).
2. `TrainingCourses.vue` + `TrainingCategories.vue` (CRUD + course card: category,
   duration, delivery type, provider, certification, prerequisite) +
   `TrainingProviders.vue` + `TrainingTrainers.vue` (master data).
3. `TrainingSessions.vue` (list + filter status Scheduled/Ongoing/Completed) +
   `TrainingSessionDetail.vue` tabs awal (Overview, Participants, Materials).
4. `TrainingParticipants.vue` (registrasi/enrollment + registration status).
5. Router `/training/*` + `meta.module: 'training'` + locale keys `training.*`
   (en.json/id.json) — core keys (catalog, sessions, participants). Sidebar tetap satu
   item di section Talent.

## P1 — Planning & Governance (migration 089)

### P1-BE (Backend) — ✅ SELESAI (2026-08-11)

1. Migration `089_training_planning.sql`:
   - tabel baru: `training_plans`, `training_plan_items`, `training_needs`,
     `training_requests` (termasuk `approval_instance_id`);
   - tabel baru: `training_course_objectives`, `training_course_competencies`,
     `training_course_prerequisites`, `training_mandatories`;
   - tabel baru: `training_session_costs`, `training_documents`.
2. Training Request + Central Approval:
   - narrow interface `ApprovalEngine` di service;
   - wiring `main.go`: `NewModuleWithService` + `SetApprovalEngine(sharedApprovalEngine)`
     + `SetNotifier(notificationSvc)` + `RegisterStatusHandler("training_request", ...)` (§45);
   - submit → `CreateApprovalInstance` → status `PENDING_APPROVAL`;
   - callback approved → status `APPROVED` + auto-enroll participant;
     rejected → tidak menjadi participant aktif.
3. CRUD plan/plan-item, need, request (+ submit/cancel), objective/competency/prerequisite
   (sub-resource course), mandatory training + compliance check, `training_session_costs`,
   `training_documents` (upload via `/api/v1/tenant/uploads`).
4. Unit test + approval integration test (`fakeApprovalEngine`, pola
   `approval_integration_test.go` leave/reimbursement/attendance).

### P1-FE (Frontend) — ✅ SELESAI (2026-08-11)

1. `TrainingPlanning.vue` (plans + plan items).
2. Training Request: form request → submit → approval flow `training_request`;
   tampilkan `approval_instance_id` + status; locale `approval.module_names`
   ditambah `"training_request": "Training Request"`.
3. `TrainingSessionDetail.vue` tab tambahan: Attendance, Assessment, Cost, Documents.
4. Training Needs page (source: MANUAL/PERFORMANCE/COMPETENCY/CAREER/SUCCESSION/
   COMPLIANCE/WORKFORCE) + Mandatory compliance view (Required/Completed/Pending/Overdue).
5. Locale keys planning/needs/requests + status enum
   (DRAFT/SUBMITTED/PENDING_APPROVAL/APPROVED/REJECTED/CANCELLED).

## P2 — Advanced Development (migration 090)

### P2-BE (Backend)

1. Migration `090_training_advanced.sql`:
   - tabel baru: `training_evaluation_forms`, `training_evaluation_questions`,
     `training_evaluation_answers`, `training_effectiveness_assessments`;
   - tabel baru: `training_certifications` + `ALTER training_certificates`
     (certification_id, certificate_file_url).
2. Evaluation form + answers; effectiveness assessment (before/after score,
   30/60/90 hari); certification master + generate certificate dari completion.
3. Reporting endpoints: `GET /trainings/reports/participation|cost|compliance`,
   `GET /trainings/reports/dashboard`, `GET /trainings/history?employee_id=`
   (route statis sebelum `/:id` — konstrain Gin).
4. Workforce Intelligence & Career Intelligence integration — kedua module membaca data
   operasional training (course–competency, completion, mandatory compliance) tanpa
   module Intelligence baru.

### P2-FE (Frontend)

1. `TrainingCertificates.vue` (master certification + generate certificate).
2. Evaluation form UI (RATING/TEXT/SINGLE_CHOICE/MULTIPLE_CHOICE) + effectiveness entry.
3. `TrainingHistory.vue` (per employee: course, date, score, completion, certificate).
4. `TrainingReports.vue` (participation, cost, compliance, dashboard analytics).
5. Integrasi analisis: Workforce Intel menampilkan training analysis; Career Intel
   menampilkan development recommendation dari course–competency mapping.
6. Locale keys lanjutan: reports, effectiveness, certification.

---

# 43. Target Operational Architecture

```text
                    TRAINING & DEVELOPMENT
                             │
       ┌─────────────────────┼─────────────────────┐
       │                     │                     │
       ▼                     ▼                     ▼
    CATALOG               PLANNING             EXECUTION
       │                     │                     │
       ├─ Category           ├─ Training Need      ├─ Session
       ├─ Course             ├─ Training Plan      ├─ Participant
       ├─ Objective          └─ Request ──► Central Approval
       ├─ Competency                               ├─ Attendance
       └─ Prerequisite                             └─ Trainer
                             │
                             ▼
                         ASSESSMENT
                             │
                  ┌──────────┴──────────┐
                  ▼                     ▼
             Evaluation            Completion
                                        │
                                        ▼
                                  Certification
                                        │
                         ┌──────────────┴──────────────┐
                         ▼                             ▼
                Workforce Intelligence       Career Intelligence
```

---

# 44. Kesimpulan

Struktur existing sudah menjadi fondasi yang baik untuk operational training
(Category, Course, Session, Participant, Material, Evaluation, Certificate + module terdaftar
di Priority 13 + CRUD API + placeholder FE). Pengembangan utama adalah mengubahnya dari
**simple training management** menjadi **end-to-end Training & Development Management**
dengan menambahkan Planning, Need, Request (approval central), Provider, Trainer, Enrollment,
Detailed Attendance, Objectives, Competency Mapping, Prerequisite, Assessment, Evaluation Form,
Effectiveness, Mandatory Training, Certification, Cost, dan Documents.

Prinsip akhirnya:

```text
Training Module        = Operational Source of Truth
Workforce Intelligence = Workforce-level analysis
Career Intelligence    = Individual career/development analysis
```

Tidak diperlukan module Intelligence baru.

---

# 45. Wiring di `main.go` (contoh pola — saat implementasi)

Mengikuti pola leave/reimbursement/attendance/employeemovement yang sudah ada:

```go
// Construct service up front agar status handler approval bisa didaftarkan.
trainingResolver := training.NewTenantDBResolver(dbManager)
trainingRepo := training.NewRepository(trainingResolver)
trainingSvc := training.NewService(trainingRepo, l.Named("training"))
trainingSvc.SetApprovalEngine(sharedApprovalEngine) // adapter yang sama dgn modul lain
trainingSvc.SetNotifier(notificationSvc)
approvalSvc.RegisterStatusHandler("training_request", func(ctx context.Context,
    documentID uuid.UUID, status approval.InstanceStatus, note string) error {
    return trainingSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
})

// Registrasi berubah dari training.NewModule(dbManager, l) menjadi
// training.NewModuleWithService(l, trainingSvc) — pola sama dgn leave/payroll.
```

Perubahan terkait di `main.go`:
- `training.NewModule(...)` → `training.NewModuleWithService(l, trainingSvc)` (Priority 13 tetap).
- Import module `approval` sudah ada.
- Locale FE `approval.module_names` ditambah `"training_request": "Training Request"`.
- Module approval checker sudah generik (module slug dicek dari company_modules),
  tidak perlu perubahan.

---

# 46. Catatan Revisi

| Tanggal | Perubahan |
| --- | --- |
| 2026-08-11 | Penyesuaian plan dengan arsitektur existing: strategi migration baru (088+), route group `/trainings` + konvensi route existing, integrasi Central Approval Engine (`training_request` + `ApprovalInstanceID`), enum Go (IN_PROGRESS, PENDING_APPROVAL, LATE), permission via `seed_rbac.go`, FE hub + sub-halaman, wiring `main.go`. Belum ada implementasi. |
| 2026-08-11 | §42 dipecah menjadi track BE dan FE per fase (P0-BE → P0-FE → P1-BE → P1-FE → P2-BE → P2-FE) — migration/model/API/test di BE; hub/router/sidebar/locale/halaman Vue di FE. Belum ada implementasi. |
| 2026-08-11 | Review pemisahan BE/FE: sinkronkan §6 ↔ §42 (088_training_core / 089_training_planning / 090_training_advanced), tambah `ALTER training_materials` di 088, pindahkan `training_session_costs` & `training_documents` ke 089 (P1-BE) + endpoint Cost/Documents di §34, tambah halaman FE Providers/Trainers/Needs/Requests di §37, klarifikasi stat hub P0-FE (pakai list existing, dashboard penuh di P2) dan deferral tab attendance/assessment ke P1-FE. Belum ada implementasi. |
| 2026-08-11 | **P0-BE selesai diimplementasikan.** Migration `088_training_core` (MySQL + PostgreSQL, up + down) + model GORM + AutoMigrate + DTO/repo/service/handler/routes (providers, trainers, session_trainers, attendances, assessments, assessment_results + enhancement courses/sessions/participants/materials) + permission granular `seed_rbac.go` + unit test. Build/vet/test ✅; migration tervalidasi di MySQL 8.0.30 & PostgreSQL 18 (idempotent, unique partial, down). Status plan diperbarui (header + §42 + §47 Progress Tracker). |
| 2026-08-11 | **P1-FE selesai diimplementasikan.** Halaman `TrainingPlans.vue` (CRUD plan + plan items expandable per plan), `TrainingRequests.vue` (create DRAFT + submit → PENDING_APPROVAL via Central Approval `training_request` + cancel; status badge + supervisor note + akses action sesuai status), `TrainingNeeds.vue` (CRUD need + filter status/sumber; posisi dari organization-summaries karena `/job-management/positions` tidak ada di project). Tab Costs & Documents di `TrainingSessionDetail.vue` (CRUD + total estimasi biaya). Hub `Training.vue`: menu Planning/Requests/Needs jadi card aktif, coming-soon hanya Certificates/Reports. Router `/training/plans|requests|needs` + locale en/id. `npm run build` ✅ (dist direbuild). Hasil review: bug `plan` undefined di expansion slot diperbaiki (`{data: planRow}`), submit/cancel kirim body `{}` (hindari gagal bind JSON), endpoint positions diganti ke organization-summaries. |
| 2026-08-11 | **P1-BE selesai diimplementasikan.** Migration `089_training_planning` (MySQL + PostgreSQL, up/down — tervalidasi: up + idempotent 2× + down bersih) untuk plans, plan_items, needs, requests (approval_instance_id), course objectives/competencies/prerequisites, mandatories, session_costs, documents. Model GORM + AutoMigrate + DTO binding `oneof` lengkap. Service: CRUD plan/plan-item, need (termasuk DeleteNeed soft-delete), request + submit/cancel via Central Approval (`training_request`) dengan auto-resolve flow + auto-enroll participant saat APPROVED, sub-resource course (objective/competency/prerequisite), mandatory, session costs, documents. Wiring `main.go`: `NewModuleWithService` + `SetApprovalEngine(sharedApprovalEngine)` + `SetNotifier(notificationSvc)` + `RegisterStatusHandler("training_request", ...)`. Permission granular `seed_rbac.go` (`plan.manage`, `request.create`, `request.approve`, `report.view`). Test: `plan_p1_test.go` (fakeApprovalEngine pola leave, 21 test — CRUD + approval integration + auto-enroll + cover competency/cost/document sesuai review). Build/vet/test ✅. |
| 2026-08-11 | **P0-FE selesai diimplementasikan.** Hub `Training.vue` (stat ringkas dari endpoint list + card navigasi + card roadmap P1/P2) + 7 sub-halaman `views/modules/training/` (Categories, Courses, Providers, Trainers, Sessions, SessionDetail tabs Overview/Participants/Materials, Participants) + router `/training/*` dengan `meta.module` + `backRoute` + locale `training.*` lengkap (en/id). `npm run build` ✅ (dist direbuild). Hasil review: filter type provider non-fungsional dihapus (BE `ListProviders` belum dukung `type`), nama course di Participants di-resolve via daftar course, stat hub per_page 500, status CANCELLED tidak tersedia saat create participant. |

---

# 47. Progress Tracker Implementasi

> **Konvensi (wajib)**: setiap kali satu phase selesai diimplementasikan,
> update status phase tsb di tabel ini + catat di §46 — dilakukan segera
> setelah implementasi selesai & tervalidasi, sebelum lanjut ke phase berikutnya.

| Phase | Track | Status | Tanggal | Ringkasan |
| --- | --- | --- | --- | --- |
| P0 | BE | ✅ SELESAI | 2026-08-11 | Migration `088_training_core` (providers, trainers, session_trainers, attendances, assessments, assessment_results + ALTER courses/sessions/participants/materials), model GORM + AutoMigrate, CRUD API provider/trainer + attendance/assessment + enhancement course/session/participant/material, permission granular `seed_rbac.go`, unit test. Build/vet/test ✅, migration tervalidasi MySQL & PostgreSQL. |
| P0 | FE | ✅ SELESAI | 2026-08-11 | Hub `Training.vue` (stat + card nav) + sub-halaman `views/modules/training/` (Courses, Categories, Providers, Trainers, Sessions, SessionDetail tabs Overview/Participants/Materials, Participants), router `/training/*` + `meta.module`, locale en/id lengkap (enum status, type, mode, role). Build ✅ (dist direbuild). |
| P1 | BE | ✅ SELESAI | 2026-08-11 | Migration `089_training_planning` (plans, plan_items, needs, requests, objectives, competencies, prerequisites, mandatories, session_costs, documents) tervalidasi MySQL & PostgreSQL; model GORM + AutoMigrate; CRUD plan/plan-item/need(+delete)/request (+ submit/cancel via Central Approval `training_request` — callback approved → auto-enroll participant); sub-resource course objective/competency/prerequisite + mandatory + session costs + documents; wiring `main.go` (NewModuleWithService + SetApprovalEngine + SetNotifier + RegisterStatusHandler); seed_rbac permission granular; unit test + approval integration test (`fakeApprovalEngine`, 21 test baru). Build/vet/test ✅. |
| P1 | FE | ✅ SELESAI | 2026-08-11 | Halaman `TrainingPlans.vue` (CRUD plan + plan items expandable), `TrainingRequests.vue` (create + submit/cancel via Central Approval + status badge + supervisor note), `TrainingNeeds.vue` (CRUD + filter status/sumber); tab Costs & Documents di `TrainingSessionDetail.vue` (CRUD + total estimasi); hub `Training.vue` menu Planning/Requests/Needs aktif (coming-soon tinggal Certificates/Reports); router `/training/plans|requests|needs` + locale en/id (±120 keys baru). Build ✅ (dist direbuild). |
| P2 | BE | ⬜ BELUM | — | Migration `090_training_advanced` + evaluation form, effectiveness, certification, reporting/history, integrasi Workforce/Career Intelligence. |
| P2 | FE | ⬜ BELUM | — | Certificates, evaluation form UI, History, Reports, analisis di Workforce/Career Intel. |

