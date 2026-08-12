# Recruitment & Onboarding — Development Plan

> 📅 Revisi struktur: 2026-08-12 (sinkron dengan template plan modul lain) · Status: **PROPOSAL — Integrated Recruitment Module (operasional)** (backend ATS dasar ✅ selesai Juli 2026 — README: 31 Jul, dashboard: 26 Jul; FE ❌ placeholder, integrasi operasional ⏳ belum dieksekusi)
> ✅ **Fakta aktual (audit 2026-08-12):** modul ini **bukan greenfield** — backend ATS dasar sudah diimplementasikan penuh (7 entity, 33 endpoint, 75 test) dan FE masih placeholder "Coming soon". Bagian "target" di dokumen ini (offer, stage history, screening, assessment, scorecard, approval, candidate enhancement, dst.) adalah **rencana enhancement**, bukan status.
> 🔎 **Sumber:** struktur tabel `015_recruitment.sql` (mysql + postgres) + audit `backend/internal/modules/recruitment/` (model.go, service.go, handler.go, routes.go, module.go) + `frontend/tenant/src/views/modules/Recruitment.vue` + `frontend/tenant/src/router/index.js` + cross-reference `docs/module-notification-plan.md` (§5/§9: "Recruitment belum tersentuh" untuk integrasi approval/notifier) + `docs/module-recruitment-strategic-layer-plan.md` (rumah item strategic layer yang dipisah) + `docs/go-module-architecture-report.md` + `docs/project-completion-dashboard.md`.
> 📊 **Progres implementasi (per 2026-08-12):** ✅ 1) Backend ATS lengkap — 7 GORM entity (`JobRequisition`, `Candidate`, `JobApplication`, `Interview`, `OnboardingTaskTemplate`, `EmployeeOnboarding`, `OnboardingTaskItem`) + enum status · ✅ 2) 33 endpoint CRUD/pipeline di 7 resource group · ✅ 3) Seeder 10 onboarding task template default · ✅ 4) 75 test (handler 28 + repository 27 + service 20) · ✅ 5) pipeline aplikasi (status + timestamp otomatis + auto `slots_filled` saat ACCEPTED) · ❌ 6) Frontend masih placeholder ("Coming soon") — hanya route/menu/locale/dashboard card · ⏳ 7) Integrasi operasional dua arah dengan modul lain (Module Approval, Notifier, Employee, Employee Movement) — **belum ada**; Employee 🔶 sebagian (onboarding menunjuk `employee_id` tanpa FK) · 🚫 8) **Scoping 2026-08-12:** Recruitment = **module operasional** — strategic layer (Workforce Intelligence, Career Intelligence, Succession, Performance, Training, Quality of Hire) **dipisah dari plan ini** — out of scope, dikelola modul masing-masing (§5.2).
> ⏳ **Sisa TODO (per review 2026-08-12):** seluruh Gap §7 (G-1 s.d. G-12) — prioritas P0: integrasi approval (G-1), enhancement requisition (G-2), offer management (G-3), Recruitment → Employee/Movement (G-4), pipeline stage history (G-5), halaman FE penuh (G-12).
> ✅ **G-1 selesai (2026-08-12):** requisition → Central Approval (migration 093, interface `ApprovalEngine`, `SubmitRequisition` DRAFT→SUBMITTED, push-callback APPROVED→OPEN / REJECTED / CANCELLED, endpoint `POST /recruitment/requisitions/:id/submit`, wiring main.go) — lihat §G-1. Bagian offer workflow menunggu G-3 (entity `job_offers` belum ada).
> 🔧 **Catatan konsistensi docs:** `project-completion-dashboard.md` masih mencatat plan ini sebagai "📋 Proposal — belum dieksekusi; backend ATS dasar sudah ada, FE masih Coming soon" — setelah revisi ini, baris tersebut sebaiknya di-update mengikuti ringkasan status di header di atas.

---

# 1. Objective

Mengembangkan Recruitment dari **simple ATS** (sudah ada) menjadi **Integrated Recruitment Module (operasional)** yang terhubung dengan Organization/Position, Module Approval, Employee, Employee Movement, Onboarding, dan Competency (untuk candidate matching). **Layer strategis** (Workforce Intelligence, Career Intelligence, Succession, Performance, Training, Quality of Hire) **tidak termasuk scope** modul ini — lihat §5.2.

```text
Job Requisition
        │
        ▼
Module Approval
        │
        ▼
Recruitment Pipeline
        │
        ├── Candidate
        ├── Screening
        ├── Assessment
        ├── Interview
        └── Selection
                │
                ▼
              Offer
                │
          ┌─────┴─────┐
          ▼           ▼
      External      Internal
       Candidate     Employee
          │           │
          ▼           ▼
      Employee    Employee Movement
          │           │
          └─────┬─────┘
                ▼
            Onboarding
```

Status per bagian:

- **ATS dasar (CRUD requisition/candidate/application/interview/onboarding)** — ✅ sudah diimplementasikan (lihat §3.1).
- **Integrated Recruitment (approval, offer, stage history, screening, assessment, scorecard, candidate enhancement, integrasi operasional)** — 🔶 sebagian: **G-1 approval requisition ✅** (2026-08-12); sisanya rencana (lihat Gap Analysis §7).

---

# 2. Existing Database Structure

Sumber: `backend/internal/pkg/migrator/migrations/tenant/mysql/015_recruitment.sql` (+ postgres). Terdapat **7 tabel**, semua PK `CHAR(36)`, timestamp `TIMESTAMP(6)`, tanpa soft delete (`deleted_at`).

## 2.1 Tabel `job_requisitions` (lowongan)

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `organization_id` | CHAR(36) NN | org pemilik lowongan (index `idx_req_org`) |
| `title` | VARCHAR(255) NN | |
| `department` | VARCHAR(150) NULL | teks bebas (bukan master) |
| `employment_type` | VARCHAR(50) NULL | |
| `location` | VARCHAR(255) NULL | |
| `min_salary` / `max_salary` | DECIMAL(15,2) | default 0 |
| `description`, `requirements`, `responsibilities` | TEXT NULL | |
| `slots_available` | INT NN default 1 | |
| `slots_filled` | INT NN default 0 | |
| `status` | VARCHAR(20) NN default `DRAFT` | enum model: `DRAFT, OPEN, IN_PROGRESS, FILLED, CANCELLED` |
| `requested_by` / `approved_by` | CHAR(36) NULL | legacy (bukan source of truth approval) |
| `target_start_date` | DATE NULL | |
| `closed_at` | BIGINT NULL default 0 | |

Index: `idx_req_org (organization_id)`.

## 2.2 Tabel `candidates`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `first_name` / `last_name` | VARCHAR(100) NN | |
| `email` | VARCHAR(255) NN | UNIQUE `idx_cand_email` — validasi duplicate email di service |
| `phone` | VARCHAR(50) NULL | |
| `address` | TEXT NULL | |
| `current_company` / `current_title` | VARCHAR(255) NULL | |
| `resume_url` / `portfolio_url` / `linkedin_url` | TEXT NULL | |
| `source` | VARCHAR(50) NN default `direct` | teks bebas (bukan master) |
| `notes` | TEXT NULL | |

Profil dasar (tanpa sub-tabel edukasi/pengalaman/skill/sertifikat/dokumen).

## 2.3 Tabel `job_applications` (pipeline)

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `requisition_id` | CHAR(36) NN | index `idx_app_req` |
| `candidate_id` | CHAR(36) NN | index `idx_app_cand` |
| `status` | VARCHAR(50) NN default `NEW` | enum model `CandidateStatus`: `NEW, SCREENED, SHORTLISTED, INTERVIEWED, OFFERED, ACCEPTED, REJECTED, WITHDRAWN` (index `idx_app_status`) |
| `applied_at` | BIGINT NN default 0 | epoch nanos |
| `screened_at`, `shortlisted_at`, `offered_at`, `accepted_at`, `rejected_at`, `withdrawn_at` | BIGINT NULL default 0 | timestamp per stage, diisi otomatis oleh service |
| `rejection_reason`, `notes` | TEXT NULL | |

## 2.4 Tabel `interviews`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `application_id` | CHAR(36) NN | index `idx_int_app` |
| `interviewer_id` | CHAR(36) NN | **single interviewer** (multi-interviewer belum didukung) |
| `stage` | VARCHAR(50) NN | default `FIRST_INTERVIEW` |
| `scheduled_at` | BIGINT NN default 0 | |
| `duration_minutes` | INT NN default 60 | |
| `location` | VARCHAR(255) NULL, `meeting_link` TEXT NULL | |
| `status` | VARCHAR(20) NN default `SCHEDULED` | `SCHEDULED, COMPLETED, CANCELLED, RESCHEDULED` |
| `score` | DECIMAL(5,2) NULL | single score |
| `feedback` | TEXT NULL | |
| `completed_at` | BIGINT NULL default 0 | |

## 2.5 Tabel `onboarding_task_templates`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `name` | VARCHAR(255) NN | |
| `description` | TEXT NULL | |
| `category` | VARCHAR(50) NULL | |
| `day_offset` | INT NN default 0 | |
| `assigned_role` | VARCHAR(50) NULL | |
| `is_mandatory` | TINYINT(1) NN default 1 | |

## 2.6 Tabel `employee_onboardings`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `employee_id` | CHAR(36) NN | index `idx_onb_emp` (tanpa FK constraint) |
| `application_id` | CHAR(36) NN | index `idx_onb_app` |
| `start_date` | DATE NN | |
| `status` | VARCHAR(20) NN default `PENDING` | `PENDING → ... → COMPLETED` (service set `completed_at`) |
| `buddy_id` | CHAR(36) NULL | |
| `completed_at` | BIGINT NULL default 0 | |
| `notes` | TEXT NULL | |

## 2.7 Tabel `onboarding_task_items`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `employee_onboarding_id` | CHAR(36) NN | index `idx_onb_task_item` |
| `template_id` | CHAR(36) NULL | referensi template (nullable untuk task custom) |
| `name` | VARCHAR(255) NN | |
| `description` | TEXT NULL | |
| `assigned_to` | CHAR(36) NULL | |
| `due_date` | BIGINT NULL default 0 | |
| `is_completed` | TINYINT(1) NN default 0 | |
| `completed_at` | BIGINT NULL default 0 | |
| `notes` | TEXT NULL | |

Relasi inti:

```text
job_requisitions ◄── job_applications ──► candidates
                              │
                              ▼
                         interviews
                              │
                              ▼
              ┌───────────────┴──────────────┐
              ▼                              ▼
   employee_onboardings ──► onboarding_task_items
              │                    │
              ▼                    ▼
        (employee_id)       (template_id → onboarding_task_templates)
```

---

# 3. Status Aktual

## 3.1 Backend — ✅ SUDAH IMPLEMENTASI (per 31 Jul 2026; audit 2026-08-12)

Modul `backend/internal/modules/recruitment/` (±2.800 baris kode non-test; total ±4.700 dengan test) sudah lengkap sebagai **ATS dasar**:

- **Model** (`model.go`): 7 GORM entity + enum `RequisitionStatus` (`DRAFT/OPEN/IN_PROGRESS/FILLED/CANCELLED`), `CandidateStatus` (`NEW/SCREENED/SHORTLISTED/INTERVIEWED/OFFERED/ACCEPTED/REJECTED/WITHDRAWN`), `InterviewStatus` (`SCHEDULED/COMPLETED/CANCELLED/RESCHEDULED`); semua PK `CHAR(36)` + `BeforeCreate` UUID.
- **Repository / Service / Handler**: CRUD lengkap per entity. Fitur service yang perlu dicatat:
  - `CreateApplication` → status `NEW` + `applied_at`, validasi requisition & candidate ada.
  - `UpdateApplicationStatus` → set timestamp stage otomatis (`ScreenedAt/ShortlistedAt/OfferedAt/AcceptedAt/RejectedAt/WithdrawnAt`); saat **ACCEPTED** otomatis `slots_filled++` dan requisition jadi `FILLED` bila slot terpenuhi.
  - `CreateEmployeeOnboarding` → **auto-generate task items** dari seluruh template aktif.
  - `UpdateInterview` dengan status `COMPLETED` → set `completed_at`; `UpdateEmployeeOnboarding` status `COMPLETED` → set `completed_at`; `UpdateOnboardingTaskItem` `is_completed=true` → set `completed_at`.
- **Routes** (`routes.go`): **33 endpoint** di 7 resource group (prefix `/api/v1/tenant/recruitment`):

| Resource | Endpoint |
|---|---|
| Requisitions | `POST/GET /requisitions`, `GET/PUT/DELETE /requisitions/:id` (5) |
| Candidates | `POST/GET /candidates`, `GET/PUT/DELETE /candidates/:id` (5) |
| Applications | `POST/GET /applications`, `GET/DELETE /applications/:id`, `PUT /applications/:id/status` (5) |
| Interviews | `POST/GET /interviews`, `GET/PUT/DELETE /interviews/:id` (5) |
| Onboarding Task Templates | `POST/GET /onboarding-task-templates`, `PUT/DELETE /onboarding-task-templates/:id` (4) |
| Employee Onboardings | `POST/GET /employee-onboardings`, `GET/PUT/DELETE /employee-onboardings/:id` (5) |
| Onboarding Task Items | `POST /onboarding-task-items`, `GET /employee-onboardings/:id/task-items`, `PUT/DELETE /onboarding-task-items/:id` (4) |

- **Module** (`module.go`): slug `recruitment`, version `1.0.0`, `DependsOn: organization, employee, setting`; permissions `recruitment.view/create/update/delete/interview/onboard`; menu `/admin/recruitment` (Requisitions, Candidates, Applications, Interviews, Onboarding); `AutoMigrate` 7 model; **seed 10 onboarding template default** (Personal Data Completion, Employment Contract Signing, IT Account Provisioning, BPJS Registration, Bank Account Setup, Orientation & Department Introduction, Policy Review & Acknowledgment, First Week Check-in, Training Plan Setup, Probation Review Preparation).
- **Wiring** (`cmd/server/main.go:854`): `recruitment.NewModule(dbManager, l)` — **plain mount, tanpa ApprovalEngine dan tanpa Notifier**.
- **Test**: 75 = `handler_test.go` 28 + `repository_test.go` 27 + `service_test.go` 20. (`README` menulis "66 unit test (27/23/16)" — angka akurat per audit file adalah **75**; selisih karena penghitungan berubah seiring penambahan test.)
- **RBAC**: `tenantseed/seed_rbac.go` mendaftarkan resource `recruitment` dengan aksi `view, create, update, delete` — **`interview` dan `onboard` ada di module.go tapi tidak di-seed** (perlu sinkronisasi saat enhancement).

## 3.2 Frontend — ❌ BELUM (placeholder)

- `frontend/tenant/src/views/modules/Recruitment.vue` — hanya placeholder **"Coming soon"** (spinner + judul statis).
- Yang **sudah ada**: route `/recruitment` (router/index.js:349-354, meta title/desc bilingual), item sidebar grup Talent (`Sidebar.vue:365`, module-gated `recruitment` + permission `recruitment.view`), locale keys `recruitment.*` EN/ID, kartu quick-access di Dashboard (`Dashboard.vue:200`).
- Tidak ada halaman: requisition list/form, candidate profile, aplikasi pipeline, interview, onboarding, analytics.

## 3.3 Integrasi lintas modul — ⏳ BELUM

| Integrasi | Status | Catatan |
|---|---|---|
| Module Approval (requisition/offer) | ⏳ Belum | notification plan §5/§9: "Recruitment belum tersentuh — belum ada integrasi approval sama sekali"; tidak ada `approval_instance_id` |
| Notification (`Notifier`) | ⏳ Belum | tidak ada interface `Notifier`/`SetNotifier` di service |
| Workforce Intelligence | 🚫 Out of scope | WI mengonsumsi candidates via `GET /workforce-intelligence/candidate-search` & `analytics/recruitment` (konsumsi sepihak di sisi WI); recruitment tidak mengirim data kembali (§5.2) |
| Employee | 🔶 Sebagian | `employee_onboardings.employee_id` (tanpa FK); tidak ada pembuatan employee dari offer/accepted application |
| Employee Movement | ⏳ Belum | tidak ada alur internal candidate → movement |
| Competency | ⏳ Belum | referensi untuk candidate matching (G-9) |
| Career Intelligence / Succession / Performance / Training | 🚫 Out of scope | strategic layer — dikelola plan modul masing-masing (§5.2) |

## 3.4 Selisih plan lama vs aktual (penting)

| Bagian plan lama | Kondisi aktual |
|---|---|
| §4.3 "job_applications perlu ditingkatkan menjadi pipeline" | **Sudah pipeline** — `status` + timestamp per stage + auto slot. Gap yang tersisa: **stage history** (audit trail), bukan status itu sendiri |
| §4.4 "interviews belum mendukung multi-interviewer scorecard" | Benar — single interviewer + single `score`/`feedback`; gap: interviewers + scorecards |
| §23-§24 interview enhancement | Belum ada `interviewers`/`interview_scorecards`/`interview_scorecard_items` |
| §26-§27 offer management | **Belum ada entity `job_offers` sama sekali** |
| §20 recruitment_stages + stage history | **Belum ada** kedua tabel |
| §13-§19 candidate enhancement | Belum ada sub-tabel education/work experience/skills/certification/document; `source` masih teks bebas |
| §7 requisition enhancement | Belum ada `requisition_number`, `reason_type`, `priority`, `position_id`, `approval_status` (plan lama juga menyebut `workforce_gap_id`/`workforce_plan_id` — kini **out of scope** WI, §5.2) |
| §31 onboarding template enhancement | Belum ada scope `organization_id/position_id/employment_type` pada template |
| §46 seeder | Seed onboarding template ✅; seeder stage/source/assessment-type belum ada |

---

# 4. Business Principles (target — dipertahankan dari plan asli)

## 4.1 Recruitment bukan master Position

Karena konsep HRIS `Organization = Position`, Recruitment tidak membuat master position/organization baru — menggunakan `organization_id` (sudah ada) atau `position_id` (bila Organization/Position dipisahkan). Master tetap di module Organization/Job Management.

## 4.2 Recruitment bukan Approval Engine

Approval seluruh proses mengikuti **Module Approval** — minimal `Job Requisition Approval` dan `Offer Approval`. Field legacy `approved_by` boleh dipertahankan untuk backward compatibility, tetapi bukan source of truth approval.

## 4.3 Recruitment bukan Employee Movement

Untuk kandidat internal (`Recruitment → Selection → Accepted → Employee Movement`), Recruitment tidak membuat employee baru.

## 4.4 Recruitment bukan Training

Recruitment hanya menghasilkan kebutuhan onboarding/development; Training tetap source of truth training.

## 4.5 Recruitment = module operasional (strategic layer di luar scope)

Recruitment adalah **module operasional**: mencari & memilih kandidat, memproses aplikasi → interview → offer → onboarding. Layer strategis — **Workforce Intelligence** (kebutuhan workforce / hiring need, forecasting) dan **Career Intelligence** (career path, gap analysis, talent map, succession, eligibility internal) — **di luar scope** modul ini dan dikelola modul masing-masing (§5.2). Recruitment hanya **mengonsumsi output**-nya secara referensial, tidak mengeksekusi.

---

# 5. Target Module Boundary

## 5.1 Cakupan module (operasional)

```text
RECRUITMENT
├── Job Requisition
├── Recruitment Approval
├── Candidate
├── Application
├── Screening
├── Assessment
├── Interview
├── Selection
├── Offer
└── Recruitment Analytics

ONBOARDING
├── Onboarding Template
├── Employee Onboarding
├── Onboarding Tasks
└── New Employee Preparation
```

Jika Onboarding nantinya menjadi module terpisah, tabel onboarding dapat dipindahkan secara bertahap tanpa mengubah flow Recruitment.

## 5.2 Out of Scope — Strategic Layer (dikelola plan `module-recruitment-strategic-layer-plan.md`)

> Keputusan scoping (2026-08-12): Recruitment = **module operasional**. Item berikut **dipisah dari plan ini** karena berada di layer strategis (Workforce Intelligence / Career Intelligence) atau di luar tanggung jawab Recruitment — seluruh detail enhancement, gap, dan API-nya dikelola di **[`docs/module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md)** (S-1 s.d. S-7):

| Item strategis | Pemilik | Plan pengelola |
|---|---|---|
| Workforce forecasting, workforce gap / hiring need, "Required − Current − Expected Hires = Remaining Gap" | Workforce Intelligence | [`module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md) S-1/S-2 + `docs/workforce-intelligence-training-enhancement-plan.md` |
| Career path, career interest, talent map (9-box), gap analysis karier, succession planning, eligibility internal candidate | Career Intelligence | [`module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md) S-4/S-5 + `docs/module-career-intelligence-plan.md` |
| Quality of Hire (probation + performance + retention) | Workforce / Career Intelligence | [`module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md) S-6 |
| Pelaksanaan training & development | Module Training | [`module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md) S-7 |

Recruitment tetap menyediakan **data operasional** yang dibutuhkan layer strategis (requisition, pipeline, offer, hire, onboarding), tetapi **tidak** memiliki logika strategis apa pun. Seluruh item di tabel atas didokumentasikan secara lengkap di **[`docs/module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md)**.

---

# 6. Target Database Architecture

```text
job_requisitions
      │
      ├── job_requisition_requirements
      ├── job_requisition_competencies
      └── recruitment approval (approval_instance_id)
      │
      ▼
job_applications
      │
      ├── application_stage_histories
      ├── application_screenings
      ├── recruitment_assessments (+ participants/results)
      ├── interviews
      │      ├── interviewers
      │      └── interview_scorecards (+ items)
      │
      └── job_offers

candidates
   ├── candidate_educations
   ├── candidate_work_experiences
   ├── candidate_skills
   ├── candidate_certifications
   ├── candidate_documents
   └── candidate_consents
```

---

# 7. Gap Analysis & Enhancement Plan

> ⚠️ Prioritas diurutkan berdasarkan dampak bisnis. Seluruh gap di bawah adalah **rencana** (belum dieksekusi per 2026-08-12).

## G-1 ✅ MODULE APPROVAL INTEGRATION (requisition ✅ · offer → G-3)

**Status: ✅ Selesai (2026-08-12) untuk requisition.** Bagian offer menunggu G-3 (entity `job_offers` belum ada).

**Yang diimplementasikan (requisition):**
- **Migration `093_recruitment_approval`** (pg + mysql, up/down idempotent): `approval_instance_id` CHAR(36) NULL di `job_requisitions`.
- **Model:** field `ApprovalInstanceID *uuid.UUID` + status constants baru `SUBMITTED`, `REJECTED` (transisi hanya via alur approval — tidak bisa di-set manual lewat `UpdateRequisitionRequest.status`).
- **Service:** interface `ApprovalEngine` narrow (`CreateApprovalInstance`, `GetApprovalInstanceStatus`, `GetActiveFlowIDForModule`) + `SetApprovalEngine`; `SubmitRequisition` (DRAFT → SUBMITTED + instance dibuat + `approval_instance_id` tersimpan); `HandleApprovalStatusChange` (push-callback: hanya SUBMITTED diproses — idempotent; APPROVED → OPEN, REJECTED → REJECTED, CANCELLED → CANCELLED, unknown no-op).
- **Handler/Routes:** `POST /recruitment/requisitions/:id/submit` dengan `approval.EmitRoutingError` (RoutingError → 400 bilingual, pola employeemovement).
- **Wiring `cmd/server/main.go`:** `recruitmentSvc.SetApprovalEngine(sharedApprovalEngine)` + `approvalSvc.RegisterStatusHandler("recruitment", ...)` (sebelum module mount Priority 11).
- **Test:** +10 test service (submit: creates instance / auto-resolve flow / not-draft / no-engine / no-flow; handler: approved→open / rejected / cancelled / not-submitted noop / unknown noop).

**Workflow requisition (aktual):**

```text
DRAFT → SUBMITTED → (Module Approval) → APPROVED/REJECTED → OPEN
```

Catatan: hasil approval tersimpan di instance Approval module (source of truth); `approval_instance_id` hanya referensi. `approved_by` dipertahankan sebagai legacy. Note reject tidak dipersist ke kolom `requirements` (teks persyaratan job asli) — note tersedia di Approval instance.

**Sisa G-1 (offer):** workflow offer `OFFER_DRAFT → SUBMITTED → APPROVED → SENT` dieksekusi bersama G-3 (entity `job_offers`).

**Ref:** plan asli §2.2, §10, §27.

## G-2 🔴 JOB REQUISITION ENHANCEMENT (master position + operational fields)

**Status: ⏳ Belum.** `organization_id` sudah ada; sisanya belum.

**Rencana (enhancement fields):**

```text
requisition_number
reason_type            NEW_POSITION | REPLACEMENT | BACKFILL | EXPANSION | INTERNAL_MOVEMENT
priority               LOW | MEDIUM | HIGH | URGENT
position_id            nullable (bila position terpisah)
approval_status
opened_at
```

- Jangan menyimpan `department`/`title` sebagai master bebas bila sudah tersedia di Organization/Position; Recruitment membaca `Position Title / Organization / Employment Type / Hierarchy / Competency Requirement` dari master.
- Business rule: requisition tidak dapat dibuka sebelum approval selesai; `slots_filled <= slots_available`.
- `workforce_gap_id` / `workforce_plan_id` **tidak termasuk** — workforce gap / hiring need adalah output Workforce Intelligence (out of scope, §5.2).

**Ref:** plan asli §7, §8, §50.

## G-3 🔴 OFFER MANAGEMENT (entity baru)

**Status: ⏳ Belum — entity `job_offers` tidak ada sama sekali.**

**Rencana (fields minimal):**

```text
id, application_id, offer_number, employment_type, salary, allowances, benefits,
start_date, expiry_date, status, sent_at, accepted_at, rejected_at, approval_instance_id, timestamps
```

Status: `DRAFT, PENDING_APPROVAL, APPROVED, SENT, ACCEPTED, REJECTED, EXPIRED, WITHDRAWN`.

Flow: `Recruiter → Offer Draft → Module Approval → Approved → Offer Sent → Candidate → Accepted`.

Business rules: hanya candidate eligible yang menerima offer; offer expired tidak dapat diterima; offer accepted menghasilkan Employee (external) atau Employee Movement (internal).

**Ref:** plan asli §26-§27.

## G-4 🔴 RECRUITMENT → EMPLOYEE / EMPLOYEE MOVEMENT

**Status: ⏳ Belum.**

**Rencana:**
- **External:** offer accepted → create employee; simpan reference `employee.recruited_from_application_id` (atau equivalent sesuai Employee module) agar `Employee → Application → Requisition → Position` dapat ditelusuri.
- **Internal (`candidate_type = INTERNAL`, `employee_id`):** Recruitment tidak membuat employee baru — hasil seleksi diteruskan ke **Employee Movement** (`Recruitment → Selection → Accepted → Employee Movement → New Organization/Position`).
- Onboarding existing sudah mendukung `employee_onboardings.application_id` + `employee_id` (fondasi Recruitment → Onboarding ada).

**Ref:** plan asli §19, §28-§30.

## G-5 🔴 PIPELINE STAGE HISTORY (audit trail & time metrics)

**Status: ⏳ Sebagian.** `job_applications.status` + timestamp per stage sudah ada; **history/audit trail belum ada**.

**Rencana:**
- `recruitment_stages` (master stage, seeded): `APPLIED, SCREENING, SHORTLISTED, ASSESSMENT, INTERVIEW, FINAL_REVIEW, OFFER, HIRED, REJECTED, WITHDRAWN`.
- `job_application_stage_histories`: `id, application_id, from_stage_id nullable, to_stage_id, changed_by, changed_at, notes`.
- Tujuan: audit trail, `Time to Stage`, `Time to Hire`, pipeline analytics.
- Setiap perubahan stage wajib menulis history (bukan sekadar update `status`).
- Validasi transisi stage (jangan asal set status; gunakan transition service + history).

**Ref:** plan asli §20, §50, §57.

## G-6 🟡 CANDIDATE ENHANCEMENT (profil terstruktur + internal candidate)

**Status: ⏳ Belum.** Profil dasar (email/phone/URL/source teks) sudah ada; sub-tabel & tipe kandidat belum.

**Rencana:**
- Kolom tambahan `candidates`: `candidate_number`, `status`, `candidate_type` (`EXTERNAL/INTERNAL`), `source_id` (master source), `consent_status`.
- Tabel baru: `candidate_educations`, `candidate_work_experiences`, `candidate_skills` (pakai `skill_id` bila Skill Master ada, jangan duplicate master), `candidate_certifications` (`certification_id` nullable), `candidate_documents` (jenis `RESUME/COVER_LETTER/CERTIFICATE/PORTFOLIO/IDENTITY/OTHER` — simpan referensi file, bukan binary), `candidate_consents`.
- Internal candidate: `candidate_type = INTERNAL` + `employee_id`; tidak membuat employee baru.

**Ref:** plan asli §13-§19.

## G-7 🟡 SCREENING & ASSESSMENT

**Status: ⏳ Belum — kedua tabel belum ada.**

**Rencana:**
- `application_screenings`: `id, application_id, screened_by, screened_at, score, result (PASS/FAIL/HOLD), notes`.
- Assessment: `recruitment_assessments`, `assessment_participants`, `assessment_results` — jenis `TECHNICAL, PSYCHOLOGICAL, COGNITIVE, PERSONALITY, CASE_STUDY, CODING, LANGUAGE, OTHER`; hasil `score, result, recommendation`.

**Ref:** plan asli §21-§22.

## G-8 🟡 INTERVIEW MULTI-INTERVIEWER & SCORECARD

**Status: ⏳ Belum.** Interview tunggal (satu `interviewer_id`, satu `score`, `feedback`) sudah ada.

**Rencana:**
- `interviewers` (satu interview → banyak interviewer: HR + User + Manager).
- `interview_scorecards` + `interview_scorecard_items` — contoh bobot: `Technical Skill 30%, Problem Solving 20%, Communication 20%, Leadership 15%, Culture Fit 15%`; skala `1-5` atau `0-100`, normalisasi di service layer.
- Endpoint `POST /interviews/:id/complete` untuk menutup interview + aggregate score.

**Ref:** plan asli §23-§24.

## G-9 🟡 CANDIDATE MATCHING & REQUISITION COMPETENCY

**Status: ⏳ Belum.**

**Rencana:**
- `job_requisition_requirements`: `id, requisition_id, requirement_type, name, description, minimum_value, maximum_value, is_required, sort_order`.
- `job_requisition_competencies`: `id, requisition_id, competency_id, required_level, is_required, weight` — basis candidate matching (contoh: `PHP L4 Required`, `Laravel L4 Required`, `PostgreSQL L3 Required`, `Leadership L2 Optional`).
- Candidate Matching memakai `Job Requirement + Competency + Education + Experience + Skill + Certification + Assessment + Interview` → `candidate_match_score` (contoh Budi 92% / Andi 87% / Dedi 76%). Match score **bukan keputusan otomatis** — recruiter dapat override dengan alasan tercatat.
- Bandingkan `Candidate Competency vs Position Competency` dari Job Management/Position.

**Ref:** plan asli §11-§12, §25, §33.

## G-10 🟢 ONBOARDING ENHANCEMENT (template scoped)

**Status: ⏳ Sebagian.** Template + auto task items sudah ada; scope per org/position/employment_type belum.

**Rencana:**
- Tambah `organization_id nullable`, `position_id nullable`, `employment_type nullable` pada `onboarding_task_templates` agar template dapat dibedakan (mis. Software Engineer → Laptop, Repository Access, Security Training, Technical Orientation, Team Introduction).
- Alur tetap: `job_application → employee → employee_onboarding`; `CreateEmployeeOnboarding` memilih template yang cocok (fallback template global).

**Ref:** plan asli §30-§31.

## G-11 🟡 RECRUITMENT ANALYTICS

**Status: ⏳ Belum.** Tidak ada endpoint analytics recruitment sendiri.

**Rencana:**
- Minimal: `Open Requisitions, Applications, Candidates, Shortlisted, Interviews, Offers, Hires, Rejected, Withdrawn`.
- Advanced: `Time to Hire, Time to Fill, Time to Stage, Offer Acceptance Rate, Application Conversion Rate, Source Conversion, Candidate Match Score`.
- Data operasional recruitment tetap tersedia untuk dibaca Workforce Intelligence (konsumsi sepihak di sisi WI — out of scope plan ini, §5.2).

**Ref:** plan asli §37.

## G-12 🔴 FRONTEND RECRUITMENT (halaman penuh)

**Status: ❌ Placeholder ("Coming soon").**

**Rencana (mengikuti pola FE modul lain — bilingual + dark mode + `ConfirmDeleteDialog` + skeleton):**
- **Hub Recruitment** (`Recruitment.vue` ditulis ulang — pola `AttendanceAdmin.vue`/`Training.vue`): kartu sub-menu Requisitions / Candidates / Applications / Interviews / Offers / Onboarding + summary cards (Open Requisitions, Candidates, Applications, Interviews, Offers, Hires, Time to Hire).
- **Requisitions**: list + create/edit dialog + approval status badge + close action.
- **Candidates**: list + profile (personal, contact, resume, education, experience, skills, certifications, applications); internal candidate → current employee/org/position (akses mengikuti permission).
- **Applications**: pipeline board (Applied → Screening → Assessment → Interview → Final Review → Offer → Hired) — drag/drop hanya memanggil transition service backend; detail + stage history + screening/assessment/interview/offer list.
- **Interviews**: kalender + schedule + interviewer + scorecard + feedback + result.
- **Offers**: draft → approval → preview → send → acceptance.
- **Onboarding**: templates + employee onboardings + task items checklist.
- **Notifications**: deep-link tipe notifikasi recruitment (setelah G-1/G-3 ada).

**Ref:** plan asli §43-§45.

---

# 8. API Plan

## 8.1 Existing (33 endpoint — sudah ada)

```http
## Requisitions
GET    /api/v1/tenant/recruitment/requisitions
POST   /api/v1/tenant/recruitment/requisitions
GET    /api/v1/tenant/recruitment/requisitions/{id}
PUT    /api/v1/tenant/recruitment/requisitions/{id}
DELETE /api/v1/tenant/recruitment/requisitions/{id}

## Candidates
GET    /api/v1/tenant/recruitment/candidates
POST   /api/v1/tenant/recruitment/candidates
GET    /api/v1/tenant/recruitment/candidates/{id}
PUT    /api/v1/tenant/recruitment/candidates/{id}
DELETE /api/v1/tenant/recruitment/candidates/{id}

## Applications
GET    /api/v1/tenant/recruitment/applications
POST   /api/v1/tenant/recruitment/applications
GET    /api/v1/tenant/recruitment/applications/{id}
PUT    /api/v1/tenant/recruitment/applications/{id}/status
DELETE /api/v1/tenant/recruitment/applications/{id}

## Interviews
GET    /api/v1/tenant/recruitment/interviews
POST   /api/v1/tenant/recruitment/interviews
GET    /api/v1/tenant/recruitment/interviews/{id}
PUT    /api/v1/tenant/recruitment/interviews/{id}
DELETE /api/v1/tenant/recruitment/interviews/{id}

## Onboarding Task Templates
GET    /api/v1/tenant/recruitment/onboarding-task-templates
POST   /api/v1/tenant/recruitment/onboarding-task-templates
PUT    /api/v1/tenant/recruitment/onboarding-task-templates/{id}
DELETE /api/v1/tenant/recruitment/onboarding-task-templates/{id}

## Employee Onboardings
GET    /api/v1/tenant/recruitment/employee-onboardings
POST   /api/v1/tenant/recruitment/employee-onboardings
GET    /api/v1/tenant/recruitment/employee-onboardings/{id}
PUT    /api/v1/tenant/recruitment/employee-onboardings/{id}
DELETE /api/v1/tenant/recruitment/employee-onboardings/{id}

## Onboarding Task Items
GET    /api/v1/tenant/recruitment/employee-onboardings/{id}/task-items
POST   /api/v1/tenant/recruitment/onboarding-task-items
PUT    /api/v1/tenant/recruitment/onboarding-task-items/{id}
DELETE /api/v1/tenant/recruitment/onboarding-task-items/{id}
```

## 8.2 Target tambahan (rencana — per Gap §7)

```http
POST   /recruitment/requisitions/{id}/submit        ← G-1 (kirim ke Module Approval)
POST   /recruitment/requisitions/{id}/close         ← G-2
GET    /recruitment/requisitions/{id}/requirements  ← G-9
GET    /recruitment/requisitions/{id}/competencies  ← G-9
POST   /recruitment/applications/{id}/stage         ← G-5 (transition + history)
POST   /recruitment/applications/{id}/screen        ← G-7
GET    /recruitment/applications/{id}/history       ← G-5
GET    /recruitment/candidates/{id}/profile         ← G-6 (educations/experiences/skills/certs/documents)
POST   /recruitment/applications/{id}/assessments   ← G-7
GET    /recruitment/assessments                     ← G-7
POST   /recruitment/interviews/{id}/complete        ← G-8
GET    /recruitment/interviews/{id}/scorecards      ← G-8
GET    /recruitment/offers                          ← G-3
POST   /recruitment/offers                          ← G-3
POST   /recruitment/offers/{id}/submit|send|accept|reject  ← G-1/G-3
GET    /recruitment/analytics/*                     ← G-11
```

---

# 9. Permissions & Authorization (target)

## 9.1 Permissions aktual

Module `module.go` mendaftarkan 6 permission: `recruitment.view/create/update/delete/interview/onboard`. Tenant RBAC seed (`seed_rbac.go`) hanya menyediakan 4: `recruitment.view/create/update/delete` — **perlu sinkronisasi** (`interview`/`onboard` ditambahkan saat enhancement dieksekusi).

## 9.2 Target permissions granular (rekomendasi — plan asli §40)

```text
recruitment.view
recruitment.requisition.view / manage / submit
recruitment.candidate.view / manage
recruitment.application.view / manage
recruitment.screening.manage
recruitment.assessment.manage
recruitment.interview.view / manage
recruitment.offer.view / manage / submit
recruitment.analytics.view
recruitment.onboarding.view / manage
```

Permission harus mengikuti pola permission module existing (resource + action, seeded di `tenantseed/`).

## 9.3 Authorization roles (rekomendasi — plan asli §41)

| Role | Cakupan |
|---|---|
| **Requester (User)** | Create/edit/submit requisition sesuai organization scope |
| **Recruiter** | Candidate, Application, Screening, Interview, Offer |
| **Hiring Manager** | View assigned requisition, review candidate, interview, beri rekomendasi |
| **HR** | Full Recruitment + Onboarding + Offer |
| **Employee (internal recruitment)** | View own application, withdraw application |

---

# 10. Frontend Plan

## 10.1 Status aktual

- Route `/recruitment` + menu sidebar + locale + dashboard card ✅ (lihat §3.2).
- Halaman `Recruitment.vue` ❌ placeholder — belum ada satu pun halaman fungsional.

## 10.2 Target

- **Recruitment Dashboard**: widgets Open Requisitions, Candidates, Applications, Interviews, Offers, Hires, Time to Hire.
- **Requisition**: list, create, detail, edit, approval status, pipeline.
- **Candidate**: list, profile (resume, experience, education, skills, certification, applications); internal → current employee/org/position.
- **Application**: pipeline board (kanban) — drag/drop hanya memanggil backend transition service.
- **Interview**: calendar, schedule, interviewer, scorecard, feedback, result.
- **Offer**: draft, approval, preview, send, acceptance.
- **Onboarding**: templates + per-employee checklist.
- **Analytics**: funnel + time metrics + source analytics.
- Seluruh halaman mengikuti pola FE existing: bilingual `t()` + dark mode + `ConfirmDeleteDialog` + `SkeletonTable`/`useSkeletonPage` + `DateInput`/`TimeInput`/`SelectLabel`/`ToggleSwitch`.

---

# 11. Seeder Plan

- ✅ **Sudah ada**: 10 onboarding task template default (`module.go` `Seed`, idempotent).
- ⏳ **Rencana** (saat G terkait dieksekusi): `RecruitmentStageSeeder`, `RecruitmentSourceSeeder`, `RecruitmentAssessmentTypeSeeder`, `RecruitmentRequirementTypeSeeder` (mengikuti pola existing `tenantseed/`).
- Jangan membuat seeder untuk data transactional (`candidate`, `application`, `interview`, `offer`) kecuali development/demo seeder.

---

# 12. Migration Plan

- ✅ **Sudah ada**: `015_recruitment.sql` (mysql + postgres) — 7 tabel dasar.
- ⏳ **Rencana** (nomor migration menyesuaikan urutan existing — 091+):

```text
M1  Enhance job_requisitions   (requisition_number, reason_type, priority, position_id,
                                 approval_status, opened_at, approval_instance_id)
M2  Enhance candidates         (candidate_number, status, candidate_type, source_id, consent_status)
M3  Create job_offers          (+ approval_instance_id)
M4  Create recruitment_stages + job_application_stage_histories
M5  Create job_requisition_requirements + job_requisition_competencies
M6  Create candidate_educations, candidate_work_experiences, candidate_skills,
     candidate_certifications, candidate_documents, candidate_consents
M7  Create application_screenings
M8  Create recruitment_assessments, assessment_participants, assessment_results
M9  Create interviewers, interview_scorecards, interview_scorecard_items
M10 Enhance onboarding_task_templates (organization_id, position_id, employment_type nullable)
```

Semua migration dibuat berpasangan mysql + postgres (+ `.down.sql` bila pola existing membutuhkan), idempotent.

---

# 13. Data Migration

- `job_requisitions.title/department/employment_type` existing: dipetakan ke master Organization/Position bila tersedia (jangan dihapus — backward compatible).
- `candidates.source` (teks bebas): dipertahankan, dimigrasikan ke source master saat dibuat.
- `job_applications.status` existing: dipetakan ke `recruitment_stages`; **history lama tidak dibuat fiktif** jika tidak tersedia.

---

# 14. Backend Architecture (rekomendasi)

```text
Recruitment
├── Domain
│   ├── Requisition
│   ├── Candidate
│   ├── Application
│   ├── Screening
│   ├── Assessment
│   ├── Interview
│   ├── Offer
│   └── Onboarding
│
├── Application
│   ├── CreateRequisition
│   ├── SubmitRequisition        ← G-1
│   ├── MoveApplicationStage     ← G-5
│   ├── ScheduleInterview
│   ├── CompleteInterview        ← G-8
│   ├── CreateOffer              ← G-3
│   └── AcceptOffer              ← G-3/G-4
│
└── Integration
    ├── Organization
    ├── Approval                 ← G-1
    ├── Employee                 ← G-4
    ├── Movement                 ← G-4
    └── Competency               ← G-9
```

Pola interface narrow (tanpa import modul langsung) + wiring di `cmd/server/main.go` — mengikuti `employeemovement.CareerExecutor`, `leave.AttendanceSessionUpdater`, `payrollApprovalAdapter`, dst.

---

# 15. Business Rules (target)

## Requisition
1. Requisition harus memiliki organization/position.
2. Requisition tidak dapat dibuka sebelum approval selesai.
3. Slots tidak boleh negatif; `slots_filled <= slots_available`.
4. Requisition dapat ditutup jika seluruh slot terpenuhi atau dibatalkan.

## Application
1. Candidate dapat memiliki banyak application.
2. Candidate tidak boleh duplicate application aktif untuk requisition yang sama.
3. Stage transition harus tervalidasi.
4. Setiap perubahan stage harus memiliki history.

## Interview
1. Interview harus memiliki application.
2. Interview harus memiliki interviewer.
3. Interview completed dapat memiliki score.
4. Multi-interviewer harus didukung.

## Offer
1. Hanya candidate yang eligible yang dapat menerima offer.
2. Offer harus melalui approval.
3. Offer accepted menghasilkan employee (external) atau Employee Movement (internal).
4. Offer expired tidak dapat diterima.

## Internal Candidate
- `candidate_type = INTERNAL` harus memiliki `employee_id` dan tidak membuat employee baru.

---

# 16. Testing Plan

## Unit Test (per entity)
- **Requisition**: create, update, submit (G-1), approval status, close, slot validation.
- **Candidate**: create, duplicate email, update, document.
- **Application**: create, duplicate active application, stage transition + history (G-5), rejection, withdrawal.
- **Interview**: schedule, reschedule, complete (G-8), scorecard, multiple interviewers.
- **Offer** (G-3): create, approval, send, accept, reject, expire.
- **Screening/Assessment** (G-7): pass/fail/hold, participant, result.

## Integration Test
```text
Requisition → Module Approval (G-1)       Application → Interview
Application → Offer → Employee (G-3/G-4)  Internal Application → Employee Movement (G-4)
Offer Accepted → Onboarding
```

## E2E Test
- **External Hiring**: `Create Requisition → Submit → Approve → Publish → Candidate Apply → Screening → Assessment → Interview → Final Selection → Offer → Approval → Accepted → Employee Created → Onboarding`.
- **Internal Hiring**: `Internal Application (candidate_type = INTERNAL) → Selection → Offer/Decision → Employee Movement → New Position → Onboarding`.

---

# 17. Security & Privacy

Candidate data mengandung data pribadi. Wajib:
- role-based access; organization scope; audit log; document access control; consent tracking; secure file storage; masking data pada role tertentu; tidak menampilkan candidate data ke user tanpa permission.
- Data sensitif jangan disimpan jika tidak diperlukan untuk proses recruitment.

---

# 18. Audit Trail

Audit minimal untuk: Requisition, Application, Stage Transition, Screening, Interview, Assessment, Offer, Onboarding.

Contoh:

```text
Application: Budi
SCREENING → INTERVIEW
Changed by: Recruiter
Date: ...
Reason: Passed screening
```

---

# 19. Notification (setelah G-1/G-3)

Integrasi notification (pola `MOVEMENT_*` / `OVERTIME_*`):

```text
REQUISITION_SUBMITTED / APPROVED / REJECTED
INTERVIEW_SCHEDULED / RESCHEDULED
ASSESSMENT_ASSIGNED
OFFER_APPROVAL_REQUIRED / OFFER_APPROVED / OFFER_SENT / OFFER_ACCEPTED
ONBOARDING_STARTED
```

Channel mengikuti notification infrastructure existing (modul `notification`, deep-link di `Notifications.vue`).

---

# 20. Development Priority

## P0 — Core Integrated Recruitment (gap kritis)
1. ✅ ATS dasar (sudah selesai — §3.1)
2. G-1 Module Approval integration
3. G-2 Requisition enhancement (master position)
4. G-3 Offer management
5. G-4 Recruitment → Employee / Employee Movement
6. G-5 Pipeline stage history
7. G-12 Frontend halaman penuh

## P1 — Intelligent Recruitment
8. G-6 Candidate enhancement (profil terstruktur + internal)
9. G-7 Screening & assessment
10. G-8 Interview multi-interviewer & scorecard
11. G-9 Requisition competency + candidate matching
12. G-11 Analytics

## P2 — Advanced
13. Candidate Pool / Tags / Talent Pool
14. Referral management
15. Candidate ranking

> Quality of Hire & Recruitment forecasting dipindah ke layer Workforce / Career Intelligence (out of scope — §5.2).

---

# 21. Recommended Implementation Order

```text
STEP 1  G-1  Module Approval integration (pola employeemovement: interface + wiring + migration)
STEP 2  G-2  Requisition enhancement (master position)
STEP 3  G-5  Pipeline stage history (recruitment_stages + histories)
STEP 4  G-3  Offer management
STEP 5  G-4  Recruitment → Employee / Employee Movement
STEP 6  G-6  Candidate enhancement
STEP 7  G-7  Screening & assessment
STEP 8  G-8  Interview scorecard
STEP 9  G-9  Competency & candidate matching
STEP 10 G-11 Analytics
STEP 11 G-12 Frontend (setelah API lengkap)
STEP 12 Testing & E2E
```

---

# 22. Final Architecture

```text
                 ┌─────────────────┐
                 │ JOB REQUISITION │
                 └────────┬────────┘
                          │
                          ▼
                    MODULE APPROVAL
                          │
                          ▼
                 ┌─────────────────┐
                 │   RECRUITMENT   │
                 ├─────────────────┤
                 │ Candidate       │
                 │ Application     │
                 │ Screening       │
                 │ Assessment      │
                 │ Interview       │
                 │ Selection       │
                 │ Offer           │
                 └────────┬────────┘
                          │
                ┌─────────┴─────────┐
                ▼                   ▼
          External Candidate   Internal Employee
                │                   │
                ▼                   ▼
            EMPLOYEE         EMPLOYEE MOVEMENT
                │                   │
                └─────────┬─────────┘
                          ▼
                     ONBOARDING
```

> 🚫 Layer strategis (Workforce Intelligence, Career Intelligence, Succession, Performance) berada di luar diagram ini — dikelola modul masing-masing (§5.2).

---

# 23. Definition of Done

Status per 2026-08-12 (✅ = sudah, ⬜ = target enhancement). Scope: **operasional** — item strategic layer (WI/CI/Succession/Quality of Hire) tidak tercantum (§5.2):

- [x] Backend ATS dasar (7 entity, 33 endpoint, 75 test).
- [x] Application memiliki pipeline status + timestamp per stage.
- [x] Interview mendukung score/feedback/complete.
- [x] Onboarding template + auto task items.
- [x] Seeder onboarding template.
- [ ] Requisition menggunakan Organization/Position master (`position_id`).
- [ ] Requisition menggunakan Module Approval.
- [ ] Offer menjadi entity sendiri + Module Approval.
- [ ] Stage transition memiliki history.
- [ ] Screening tersedia.
- [ ] Assessment tersedia.
- [ ] Interview mendukung multi-interviewer + scorecard.
- [ ] Candidate memiliki profile terstruktur (education/experience/skills/certification/document).
- [ ] Candidate mendukung internal/external.
- [ ] External candidate dapat menjadi Employee.
- [ ] Internal candidate menggunakan Employee Movement.
- [ ] Offer accepted dapat membuat onboarding.
- [ ] Candidate dapat dinilai terhadap competency requirement.
- [ ] Recruitment analytics tersedia.
- [ ] Permission lengkap + sinkron (module.go vs seed_rbac: `interview`/`onboard`).
- [ ] Audit trail tersedia.
- [ ] Notification tersedia.
- [ ] Frontend Recruitment selesai (semua halaman).
- [ ] Unit/integration/E2E external & internal hiring selesai.
- [ ] Migration & backward compatibility diverifikasi.

---

# 24. Kesimpulan

Target akhir bukan sekadar **ATS** tetapi **Integrated Recruitment (operasional)** yang menjadi penghubung:

```text
Requisition → Approval → Hire → Employee / Movement → Onboarding
```

Prinsip pembagian responsibility (strategic layer di luar scope — §5.2):

```text
Recruitment            → mencari dan memilih kandidat (operasional)
Module Approval        → mengelola approval
Employee               → menjadi master employee
Employee Movement      → mengeksekusi perpindahan internal
Onboarding             → mempersiapkan employee baru
Workforce Intelligence → menentukan kebutuhan workforce   (out of scope — §5.2)
Career Intelligence    → menganalisis career, talent, gap, succession (out of scope — §5.2)
```

Dengan backend ATS yang sudah ada sebagai fondasi, langkah berikutnya difokuskan pada integrasi operasional (Module Approval → Offer → Employee/Movement → pipeline history) dan halaman frontend penuh — lihat Gap Analysis §7 dan urutan implementasi §21.
