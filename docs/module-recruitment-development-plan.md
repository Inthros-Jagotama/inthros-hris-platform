# Recruitment & Onboarding — Development Plan

> 📅 Revisi struktur: 2026-08-12 (sinkron dengan template plan modul lain) · Status: **PROPOSAL — Integrated Recruitment Module (operasional)** (backend ATS dasar ✅ selesai Juli 2026 — README: 31 Jul, dashboard: 26 Jul; FE ❌ placeholder, integrasi operasional ⏳ belum dieksekusi)
> ✅ **Fakta aktual (audit 2026-08-12):** modul ini **bukan greenfield** — backend ATS dasar sudah diimplementasikan penuh (7 entity, 33 endpoint, 75 test) dan FE masih placeholder "Coming soon". Bagian "target" di dokumen ini (offer, stage history, screening, assessment, scorecard, approval, candidate enhancement, dst.) adalah **rencana enhancement**, bukan status.
> 🔎 **Sumber:** struktur tabel `015_recruitment.sql` (mysql + postgres) + audit `backend/internal/modules/recruitment/` (model.go, service.go, handler.go, routes.go, module.go) + `frontend/tenant/src/views/modules/Recruitment.vue` + `frontend/tenant/src/router/index.js` + cross-reference `docs/module-notification-plan.md` (§5/§9: "Recruitment belum tersentuh" untuk integrasi approval/notifier) + `docs/module-recruitment-strategic-layer-plan.md` (rumah item strategic layer yang dipisah) + `docs/go-module-architecture-report.md` + `docs/project-completion-dashboard.md`.
> 📊 **Progres implementasi (per 2026-08-12):** ✅ 1) Backend ATS lengkap — 7 GORM entity (`JobRequisition`, `Candidate`, `JobApplication`, `Interview`, `OnboardingTaskTemplate`, `EmployeeOnboarding`, `OnboardingTaskItem`) + enum status · ✅ 2) 33 endpoint CRUD/pipeline di 7 resource group · ✅ 3) Seeder 10 onboarding task template default · ✅ 4) 75 test (handler 28 + repository 27 + service 20) · ✅ 5) pipeline aplikasi (status + timestamp otomatis + auto `slots_filled` saat ACCEPTED) · ❌ 6) Frontend masih placeholder ("Coming soon") — hanya route/menu/locale/dashboard card · ⏳ 7) Integrasi operasional dua arah dengan modul lain (Module Approval, Notifier, Employee, Employee Movement) — **belum ada**; Employee 🔶 sebagian (onboarding menunjuk `employee_id` tanpa FK) · 🚫 8) **Scoping 2026-08-12:** Recruitment = **module operasional** — strategic layer (Workforce Intelligence, Career Intelligence, Succession, Performance, Training, Quality of Hire) **dipisah dari plan ini** — out of scope, dikelola modul masing-masing (§5.2).
> ⏳ **Sisa TODO (per review 2026-08-13):** G-1 s.d. G-8 selesai; G-9 sebagian (requirement/competency tables ✅, match score ⏳) — sisa gap: G-9 match score, G-10 onboarding scoping, G-11 analytics, G-12 FE penuh.
> ✅ **G-9 sub-project 1 selesai (2026-08-13):** requisition requirements + competencies (migration 105 tabel `job_requisition_requirements` + `job_requisition_competencies`; `requirement_type` free-text, `required_level`/`weight` tidak di-enforce skala DB) — fondasi data untuk candidate matching, lihat §G-9. Algoritma `candidate_match_score` (sub-project 2, lintas Requirement+Competency+Education+Experience+Skill+Certification+Assessment+Interview) belum dieksekusi.
> ✅ **G-7 selesai (2026-08-13):** sub-project 1 application screening (migration 102 tabel `application_screenings`; many-per-application seperti Interview) + sub-project 2 assessment (migration 103 tabel `recruitment_assessments` + `assessment_participants` — 2 tabel bukan 3, hasil digabung ke participant karena 1:1); keduanya CRUD murni tanpa auto-transition status/stage-history — recruiter tetap update status manual — lihat §G-7.
> ✅ **G-8 selesai (2026-08-13):** interview multi-interviewer & scorecard (migration 104 tabel `interviewers` — melengkapi `interviews.interviewer_id` existing, backward compat — + `interview_scorecard_items` — kriteria bebas per interview tanpa master; endpoint `POST /interviews/:id/complete` menghitung weighted average scorecard items ke `Interview.Score` existing) — lihat §G-8.
> ✅ **G-1 selesai (2026-08-12):** requisition → Central Approval (migration 093, interface `ApprovalEngine`, `SubmitRequisition` DRAFT→SUBMITTED, push-callback APPROVED→OPEN / REJECTED / CANCELLED, endpoint `POST /recruitment/requisitions/:id/submit`, wiring main.go) — lihat §G-1. Bagian offer workflow menunggu G-3 (entity `job_offers` belum ada).
> ✅ **G-2 selesai (2026-08-12):** requisition enhancement (migration 094: `requisition_number` auto REQ-YYYYMM-XXXXXXXX, `priority` LOW/MEDIUM/HIGH/URGENT default MEDIUM, `position_id` referensi master position, `opened_at` diset otomatis saat OPEN) — lihat §G-2. `approval_status` TIDAK ditambahkan (G-1 sudah meng-cover via status requisition + approval_instance_id).
> ✅ **G-3 selesai (2026-08-12):** offer management (migration 095 tabel `job_offers`; workflow DRAFT → PENDING_APPROVAL → APPROVED → SENT → ACCEPTED/REJECTED/EXPIRED/WITHDRAWN via Central Approval modul `recruitment_offer`; accept menautkan application ACCEPTED + `slots_filled`++ dengan guard idempotensi; expired guard) — lihat §G-3. BE + FE lengkap (128 test).
> ✅ **G-4 selesai (2026-08-12):** Recruitment → Employee / Employee Movement (migration 096: `employee.recruited_from_application_id` + `candidates.candidate_type`/`employee_id`; offer eksternal diterima → Employee module membuat employee baru dengan referensi; kandidat INTERNAL → diteruskan ke Employee Movement promotion/mutation; guard idempotensi no-duplicate-handoff) — lihat §G-4. BE lengkap (128 test), FE onboarding tetap manual.
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
- **Integrated Recruitment (approval, offer, stage history, screening, assessment, scorecard, candidate enhancement, integrasi operasional)** — 🔶 sebagian: **G-1 approval requisition ✅** + **G-2 requisition enhancement ✅** + **G-3 offer management ✅** + **G-4 Recruitment → Employee/Movement ✅** + **G-5 stage history ✅** + **G-6 candidate enhancement ✅** + **G-7 screening & assessment ✅** + **G-8 interview scorecard ✅** + **G-9 requirement/competency (sub-1) ✅** (2026-08-12/13); sisanya (G-9 match score, G-10 s.d. G-12) rencana (lihat Gap Analysis §7).

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

## 3.2 Frontend — 🔶 SEBAGIAN (hub + beberapa halaman sudah, sisanya placeholder)

> ✅ **Update 2026-08-12 (sore):** menyusul commit `20a9b74` (G-1 submit), `1f1f1e7`/`b2a80d9` (G-2 priority/requisition_number/opened_at), `91ca3f9` (G-3 `Offers.vue`), `85ed0a6` (G-4 `Onboarding.vue`), `9d2b3c0` (badge From Offer di `Employees.vue`) — frontend Recruitment **bukan lagi placeholder tunggal**.

- Halaman **sudah ada**: `Requisitions.vue` (list/create/edit + submit + priority/requisition_number/opened_at), `Offers.vue` (draft → submit → approval → send → accept/reject/withdraw), `Onboarding.vue` (list employee onboarding + badge "From Offer" + status PENDING→IN_PROGRESS→COMPLETED + create dialog dengan auto-suggest employee dari offer diterima), `CandidateSearch.vue`, `InternalCandidates.vue`, `RecruitmentAnalytics.vue` (lihat §5.2/strategic layer plan untuk dua yang terakhir).
- `Recruitment.vue` (hub): sudah menampilkan kartu menu ke halaman-halaman di atas; sisanya (Interview, Screening, Assessment, Candidate profile) masih kartu **"Coming soon"**.
- Yang sudah ada dari awal: route `/recruitment` (router/index.js), item sidebar grup Talent (`Sidebar.vue`, module-gated `recruitment` + permission `recruitment.view`), locale keys `recruitment.*` EN/ID, kartu quick-access di Dashboard (`Dashboard.vue`).
- Tidak ada halaman: candidate profile terstruktur (education/experience/skills/cert/document), aplikasi pipeline (kanban board), interview (kalender/scorecard), screening/assessment.

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
| §20 recruitment_stages + stage history | **Sudah ada** (G-5, 2026-08-12) — 2 tabel + state machine + endpoint history |
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

> ⚠️ Prioritas diurutkan berdasarkan dampak bisnis. G-1 s.d. G-4 ✅ (2026-08-12); gap di bawah sisanya rencana (belum dieksekusi).

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

## G-2 ✅ JOB REQUISITION ENHANCEMENT (master position + operational fields)

**Status: ✅ Selesai (2026-08-12).**

**Yang diimplementasikan:**
- **Migration `094_recruitment_requisition_enhancement`** (pg + mysql, up/down idempotent): 4 kolom di `job_requisitions` — `requisition_number` VARCHAR(50) NULL, `priority` VARCHAR(10) NOT NULL DEFAULT 'MEDIUM', `position_id` CHAR(36) NULL, `opened_at` BIGINT NULL.
- **`requisition_number`** — auto-generated `REQ-YYYYMM-XXXXXXXX` saat create (8 hex char UUID, pola nomor dokumen `TRN-*` training); bisa override via create/update.
- **`priority`** — LOW | MEDIUM | HIGH | URGENT (binding `oneof`), default MEDIUM bila client kosong.
- **`position_id`** — referensi master position (tabel `positions`); tanpa FK/validasi karena modul Organization tidak mengekspos CRUD position (referensi saja).
- **`opened_at`** — unix nano, diset otomatis saat requisition **bertransisi menjadi OPEN** (via approval APPROVED callback G-1 maupun update status manual); nilai 0 (GORM read-back default) diperlakukan sebagai belum dibuka di response.
- **`approval_status` TIDAK ditambahkan** — G-1 sudah meng-cover via status requisition (SUBMITTED/OPEN/REJECTED) + `approval_instance_id`; menambah kolom terpisah = duplikasi.
- `reason_type` enum existing dipertahankan (sudah mencakup NEW_POSITION/REPLACEMENT/EXPANSION/WORKFORCE_GAP/SUCCESSION_GAP) — BACKFILL/INTERNAL_MOVEMENT dari rencana awal tidak ditambahkan agar tidak merusak data & FE existing.
- **Test:** +7 service test (auto-gen number + format, explicit number/priority/position, default priority, opened_at on open, approval approved → opened_at, clear position_id, set priority+number).

**Business rules:** requisition dibuka melalui approval (G-1 single approval path); `slots_filled <= slots_available` sudah di-enforce (auto-increment saat ACCEPTED + auto-FILLED, §G-1 existing). `workforce_gap_id`/`workforce_plan_id` tidak termasuk scope G-2 (output Workforce Intelligence).

**Ref:** plan asli §7, §8, §50.

## G-3 ✅ OFFER MANAGEMENT (entity baru)

**Status: ✅ Selesai (2026-08-12) — BE lengkap; FE menunggu G-12.**

**Yang diimplementasikan:**
- **Migration `095_recruitment_offer`** (pg + mysql, up/down idempotent): tabel baru `job_offers` — `id`, `application_id` (FK → `job_applications`), `offer_number`, `employment_type`, `salary`, `allowances`, `benefits`, `start_date`, `expiry_date`, `status` (default DRAFT), `sent_at`/`accepted_at`/`rejected_at` (bigint unix nano), `approval_instance_id` (CHAR(36) NULL), timestamps; index `idx_offer_app` + `idx_offer_status`.
- **Model:** `JobOffer` + enum `OfferStatus` (`DRAFT, PENDING_APPROVAL, APPROVED, SENT, ACCEPTED, REJECTED, EXPIRED, WITHDRAWN`).
- **Service:**
  - `CreateOffer` — auto `offer_number` `OFF-YYYYMM-XXXXXXXX` (8 hex char, pola G-2) + validasi application exists.
  - `UpdateOffer` / `DeleteOffer` — hanya offer **DRAFT** yang bisa diedit/dihapus.
  - `SubmitOffer` — DRAFT → PENDING_APPROVAL via Central Approval modul **`recruitment_offer`** + auto-resolve flow aktif (pola G-1 requisition).
  - `HandleOfferApprovalStatusChange` — push-callback (hanya PENDING_APPROVAL diproses, idempotent): APPROVED → APPROVED, REJECTED → REJECTED, CANCELLED → WITHDRAWN.
  - `SendOffer` — APPROVED → SENT (+`sent_at`).
  - `AcceptOffer` — SENT → ACCEPTED (+`accepted_at`) dengan **guard expired** (offer melewati `expiry_date` tidak dapat diterima → otomatis EXPIRED); penerimaan menautkan balik ke application (→ ACCEPTED) dan **`slots_filled`++** requisition (→ FILLED bila penuh) dengan **guard idempotensi**: tidak double-count bila application sudah ACCEPTED (jalur manual `UpdateApplicationStatus` / offer lain).
  - `RejectOffer` — SENT → REJECTED (+`rejected_at`, catatan kandidat menolak).
  - `WithdrawOffer` — DRAFT/APPROVED → WITHDRAWN (recruiter menarik).
- **Handler/Routes:** CRUD `/recruitment/offers` + `POST /offers/:id/{submit,send,accept,reject,withdraw}`.
- **Wiring `cmd/server/main.go`:** `approvalSvc.RegisterStatusHandler("recruitment_offer", ...)`.
- **Test:** +12 service test (create + format OFF- 19 char, invalid application, submit creates instance + auto-resolve flow, approval callback rejected, chain approval→send→accept (application ACCEPTED + slots_filled), expired accept guard, reject, withdraw, update/delete only-draft, **no double-increment slots_filled**). Total recruitment test: **123**.

**Workflow offer (aktual):**

```text
DRAFT → (Submit) → PENDING_APPROVAL → (Module Approval) → APPROVED → (Send) → SENT → (Candidate) → ACCEPTED / REJECTED / EXPIRED
         └────────────── (Withdraw) → WITHDRAWN ──────────────────────────┘
```

Catatan: `offer accepted` menghasilkan Employee (external, G-4) atau Employee Movement (internal, G-4). `accepted` menautkan application ke ACCEPTED + increment `slots_filled` — fondasi G-4.

**Ref:** plan asli §26-§27.

## G-4 ✅ RECRUITMENT → EMPLOYEE / EMPLOYEE MOVEMENT

**Status: ✅ Selesai (2026-08-12) — BE + FE lengkap.** `Onboarding.vue` (commit `85ed0a6`) sudah memakai alur otomatis: dropdown aplikasi ACCEPTED (exclude yang sudah punya onboarding), dropdown employee bertanda ★ From Offer, auto-suggest employee via `recruited_from_application_id` saat aplikasi dipilih. Badge "From Offer" juga tampil di `Employees.vue` (commit `9d2b3c0`).

**Yang diimplementasikan:**
- **Migration `096_recruitment_employee_handoff`** (pg + mysql, up/down idempotent):
  - `employee.recruited_from_application_id` CHAR(36) NULL → referensi balik ke aplikasi recruitment asal (`Employee → Application → Requisition → Position` traceability).
  - `candidates.candidate_type` VARCHAR(20) NOT NULL DEFAULT `'EXTERNAL'` (EXTERNAL | INTERNAL).
  - `candidates.employee_id` CHAR(36) NULL → referensi employee untuk kandidat INTERNAL.
- **Employee module:** model/DTO/response `RecruitedFromApplicationID` — tersimpan saat create + terekspos di response (via `ToResponse`).
- **Recruitment:**
  - `Candidate.CandidateType` (default EXTERNAL) + `Candidate.EmployeeID`; `Create/UpdateCandidateRequest` menerima `candidate_type` + `employee_id` (helper `applyCandidateTypeFields` — referensi employee dibersihkan bila type diubah ke EXTERNAL).
  - Interface narrow `EmployeeProvider.CreateHiredEmployee` + `MovementProvider.CreateHiredMovement` + setters (pola S-1..S-7 — Recruitment TIDAK membuat employee/movement sendiri).
  - **`AcceptOffer` → `handoffHiredEmployee`:**
    - **External** (`candidate_type=EXTERNAL`) → Employee module membuat employee baru (nama/email/phone dari kandidat + `recruited_from_application_id`), `EmployeeID` auto `EMP-XXXXXXXX` di adapter.
    - **Internal** (`candidate_type=INTERNAL` + `employee_id`) → diteruskan ke **Employee Movement** (promotion bila posisi target terisi, mutation bila hanya organisasi; SK `SK-MOV-YYYYMM-XXXXXXXX` auto; effective date dari `offer.start_date`) — bukan employee baru.
    - **Guard transisi status** (`!wasAccepted`): handoff hanya saat aplikasi BARU menjadi ACCEPTED — offer kedua di aplikasi yang sama tidak membuat employee/movement duplikat (idempotensi mirror slots_filled G-3).
    - **Best-effort:** provider nil / error downstream hanya di-log warning — accept offer tidak pernah gagal karenanya.
- **Wiring `cmd/server/main.go`:** `employeeHireAdapter` (employee.Service instance terpisah, pola employeeCareerRepo) + `movementHireAdapter` (employeeMovementSvc) + `SetEmployeeProvider`/`SetMovementProvider` sebelum module mount.
- **Test:** recruitment +5 (external→employee, internal→movement, no-provider fail-safe, candidate internal create/update + employee_id cleared on EXTERNAL, **no duplicate handoff on second offer**); employee +1 (`recruited_from_application_id` persisted). Recruitment total: **128** (handler 28 + repository 27 + service 73).

**Alur aktual:**

```text
External:  Offer Accepted → (Employee module) → Employee baru + recruited_from_application_id
Internal:  Offer Accepted → (Employee Movement) → promotion/mutation DRAFT untuk employee yang sudah ada
```

Catatan: Onboarding existing sudah mendukung `employee_onboardings.application_id` + `employee_id` — employee eksternal hasil G-4 siap dirujuk onboarding. Movement internal dibuat status DRAFT di module movement (HR menyelesaikan alur approval movement).

**Ref:** plan asli §19, §28-§30.

## G-5 ✅ PIPELINE STAGE HISTORY (audit trail & time metrics)

**Status: ✅ Selesai (2026-08-12).**

**Yang diimplementasikan:**
- **Migration `097_recruitment_stage_history`** (pg + mysql, up/down idempotent): dua tabel baru — `recruitment_stages` (master stage: `id, code, name, sort_order, created_at, updated_at`, `code` UNIQUE; seeded 8 stage per `CandidateStatus`: NEW, SCREENED, SHORTLISTED, INTERVIEWED, OFFERED, ACCEPTED, REJECTED, WITHDRAWN — seed ditulis langsung di migration SQL via `INSERT ... ON CONFLICT (code) DO NOTHING` / `INSERT IGNORE`, lihat catatan seeder di bawah) + `job_application_stage_histories` (audit trail: `id, application_id, from_stage_id nullable, to_stage_id, changed_by, changed_at, notes`; index `idx_ash_app` + `idx_ash_changed_at`).
- **Model:** `RecruitmentStage` + `ApplicationStageHistory`; enum `CandidateStatus` (8 nilai existing) **tidak diubah** — `from_stage_id`/`to_stage_id` pada history adalah kolom UUID terpisah yang mereferensikan `recruitment_stages`, bukan perubahan pada enum itu sendiri.
- **Service:** state machine `transitionApplicationStatus` — validasi transisi (forward jumps allowed, backward/from-terminal rejected, same-status no-op) + mandatory history write; wired ke `UpdateApplicationStatus`, `AcceptOffer`, dan `CreateApplication` (writes initial NEW history row); sentinel error `ErrInvalidStatusTransition`.
- **Handler/Routes:** `GET /recruitment/applications/:id/history` — return `[]ApplicationStageHistory` dengan detail stage + `changed_by` (UUID string aktor, non-nil untuk transisi manual via `UpdateApplicationStatus`; nil untuk transisi sistem seperti `AcceptOffer`).
- **Error mapping:** `ErrInvalidStatusTransition` → HTTP 400 di-handle langsung di `handler.go` (bukan `cmd/server/main.go`, yang tidak disentuh fitur ini).
- **Seeder:** karena `Module.Seed()` (di `module.go`) tidak pernah dipanggil untuk tenant DB (lihat catatan wiring di §3.x/general), 8 stage default di-seed langsung di migration SQL (idempotent per-row via unique `code`) alih-alih lewat kode Go.
- **Test:** +15 test service (transition validation, history write, forward jump, backward reject, same-status noop, initial NEW row, accept offer history, handler route); total recruitment: **143** (handler 28 + repository 27 + service 88).

**Workflow aplikasi (audit trail aktual):**

```text
NEW → SCREENED → SHORTLISTED → INTERVIEWED → OFFERED → ACCEPTED
   ↘ (each stage wrote to job_application_stage_histories)
```

Catatan: history entry nullable `from_stage_id` saat aplikasi baru (initial NEW); setiap transisi tulis row baru dengan user & timestamp otomatis. Validasi mencegah backtrack dan duplikat transisi dalam proses sebelum finalkan.

**Ref:** plan asli §20, §50, §57.

## G-6 🔶 CANDIDATE ENHANCEMENT (profil terstruktur + internal candidate) — sub-project 1/3 ✅ + sub-project 2/3 ✅ + sub-project 3a ✅ + sub-project 3b ✅ (sisa: status/source_id)

**Status: ✅ Selesai (2026-08-13) — semua sub-table G-6 yang direncanakan sudah selesai (sub-project 1+2+3a+3b); sisa scope hanya `candidates.status` + `source_id`, keduanya sudah diputuskan skip/deferred sejak brainstorming sub-project 1 (bukan pekerjaan yang masih dijadwalkan).**

**Sub-project 1 (2026-08-12):** Kolom `candidate_number` + tabel `candidate_educations` + tabel `candidate_work_experiences` ✅ (migration 098).

**Sub-project 2 (2026-08-12):** Tabel `candidate_skills` + tabel `candidate_certifications` ✅ (migration 099); sisa scope (status, source_id, documents, consents) deferred ke sub-project 3.

**Sub-project 3a (2026-08-13):** Tabel `candidate_documents` ✅ (migration 100); sisa scope (status, source_id, consents) masih deferred.

**Sub-project 3b (2026-08-13):** Tabel `candidate_consents` ✅ (migration 101) — LAST dari originally-planned G-6 sub-tables; sisa scope (status, source_id) tetap skipped/deferred seperti keputusan sub-project 1.

**Yang diimplementasikan (sub-project 1):**
- **Migration `098_recruitment_candidate_profile_basics`** (pg + mysql, up/down idempotent): 
  - Kolom `candidate_number` VARCHAR(50) NULL di `candidates` — auto-generated `CAND-YYYYMM-XXXXXXXX` saat create (8 hex char UUID, pola G-2/G-3 requisition_number/offer_number); bisa override via create/update.
  - Tabel baru `candidate_educations` — `id, candidate_id (FK), education_id (FK → setting.educations — level master), institution_name, education_major_id (FK → setting.education_majors), major (free-text fallback), gpa, start_year, end_year, is_highest, notes, created_at, updated_at`; index `idx_ed_cand`.
  - Tabel baru `candidate_work_experiences` — `id, candidate_id (FK), company_name, job_title, employment_type, start_date, end_date, is_current, description, created_at, updated_at`; index `idx_wx_cand`.
- **Model:** `CandidateEducation` + `CandidateWorkExperience` structs + `CandidateNumber` field di `Candidate`; relation fields `Education *setting.Education` dan `EducationMajor *setting.EducationMajor` pada `CandidateEducation` (mirroring `EmployeeEducation` pattern).
- **Service:** CRUD methods untuk kedua entitas + `generateCandidateNumber()` helper; wired ke `CreateCandidate` auto-generate.
- **Handler/Routes:** 8 endpoint CRUD — `POST/GET /recruitment/candidates/:id/educations`, `PUT/DELETE /recruitment/educations/:id`, `POST/GET /recruitment/candidates/:id/work-experiences`, `PUT/DELETE /recruitment/work-experiences/:id` (mirroring `Interview`/`OnboardingTaskTemplate` pattern; path param `:id` bukan `:candidate_id` — Gin panic kalau dua route berbagi posisi segment path dengan nama wildcard berbeda, dan `/candidates/:id` sudah ada duluan).
- **DTO:** `CandidateEducationRequest/Response` + `CandidateWorkExperienceRequest/Response`; add `candidate_number` ke `CandidateResponse`.
- **Test:** +17 test (repository +4: educations + work_experiences CRUD roundtrip; service +9: both entities CRUD + candidate_number auto-gen format + override + explicit update; handler +4: both entities create/list/update/delete). Total recruitment: **161** (handler 36 + repository 34 + service 91).

**Yang diimplementasikan (sub-project 2):**
- **Migration `099_recruitment_candidate_skills_certifications`** (pg + mysql, up/down idempotent):
  - Tabel baru `candidate_skills` — `id, candidate_id (FK), competency_id (FK → competency.Competency — existing master), level SMALLINT NULL (Go: *int), notes TEXT NULL, created_at, updated_at`; index `idx_cand_skill_candidate` + `idx_cand_skill_competency`; NOT NULL constraints + referential integrity.
  - Tabel baru `candidate_certifications` — `id, candidate_id (FK), name VARCHAR(255) NN (free-text, no master), issuing_organization VARCHAR(255) NULL, issue_date DATE NULL, expiry_date DATE NULL, credential_url TEXT NULL, notes TEXT NULL, created_at, updated_at`; index `idx_cand_cert_candidate`; no master reference (Certification Master belum ada di codebase).
- **Model:** `CandidateSkill` + `CandidateCertification` structs; relation field `Competency *competency.Competency` pada `CandidateSkill` (mirroring competency referencing pattern; `competencies` table dimiliki modul `competency`, bukan `setting`).
- **Service:** CRUD methods untuk kedua entitas; tidak ada validasi range level (1-5 adalah konvensi terdokumentasi, tidak di-enforce di layer manapun — sesuai keputusan desain).
- **Handler/Routes:** 8 endpoint CRUD — `POST/GET /recruitment/candidates/:id/skills`, `PUT/DELETE /recruitment/skills/:id`, `POST/GET /recruitment/candidates/:id/certifications`, `PUT/DELETE /recruitment/certifications/:id` (identik pola G-6 sub-1 education/work-experience; path param `:id`).
- **DTO:** `CreateCandidateSkillRequest`, `UpdateCandidateSkillRequest`, `CandidateSkillResponse` + `CreateCandidateCertificationRequest`, `UpdateCandidateCertificationRequest`, `CandidateCertificationResponse`.
- **Test:** +10 test (repository +4: skills + certifications CRUD; service +4: create-skill + unknown-candidate guard + unknown-competency guard + create-certification; handler +2: skill + certification create). Total recruitment: **171** (handler 38 + repository 38 + service 95).

**Yang diimplementasikan (sub-project 3a):**
- **Migration `100_candidate_documents`** (pg + mysql, up/down idempotent): tabel baru `candidate_documents` — `id, candidate_id (FK → candidates, ON DELETE CASCADE), document_type VARCHAR(20) NN DEFAULT 'OTHER', name VARCHAR(255) NN, file_url TEXT NN, notes TEXT NULL, created_at, updated_at`; index `idx_cand_doc_candidate`. Referensi saja, bukan binary — file sesungguhnya diupload lewat endpoint generik `POST /api/v1/tenant/uploads` (`backend/internal/pkg/upload`) yang mengembalikan URL; tabel ini hanya menyimpan URL tersebut.
- **Model:** `CandidateDocument` struct (`internal/modules/recruitment/model.go`); `document_type` (`RESUME/COVER_LETTER/CERTIFICATE/PORTFOLIO/IDENTITY/OTHER`) di-enforce via Gin binding `oneof=...` di layer request, bukan constraint DB — pola sama dengan `CandidateType`/`OfferStatus` di modul ini.
- **Service:** CRUD methods `CreateCandidateDocument`, `ListCandidateDocuments`, `UpdateCandidateDocument`, `DeleteCandidateDocument`.
- **Handler/Routes:** 4 endpoint CRUD — `POST/GET /recruitment/candidates/:id/documents`, `PUT/DELETE /recruitment/documents/:id` (identik pola G-6 sub-1/sub-2; path param `:id`).
- **DTO:** `CreateCandidateDocumentRequest`, `UpdateCandidateDocumentRequest`, `CandidateDocumentResponse`.
- **Test:** +9 test. Total recruitment: **180** (handler 41 + repository 41 + service 98).

**Yang diimplementasikan (sub-project 3b):**
- **Migration `101_candidate_consents`** (pg + mysql, up/down idempotent): tabel baru `candidate_consents` — `id, candidate_id (FK → candidates, ON DELETE CASCADE), action VARCHAR(20) NN, notes TEXT NULL, changed_by CHAR(36) NULL, changed_at BIGINT NN, created_at`; index `idx_cand_consent_candidate`. Append-only audit log consent pemrosesan data pribadi (GDPR-style) — **tidak ada `updated_at`**, baris tidak pernah diupdate, hanya di-insert. `changed_by` = staff/recruiter yang mencatat (sistem ini tidak punya portal publik candidate-facing, jadi consent didokumentasikan HR/recruiter, bukan self-service).
- **Model:** `CandidateConsent` struct (`internal/modules/recruitment/model.go`); `action` (`GRANTED/REVOKED`) di-enforce via Gin binding `oneof=...` di layer request, bukan constraint DB — pola sama dengan modul lain.
- **Service/Handler:** **hanya 2 method — `Create` + `List`, BUKAN CRUD 4/5-method seperti sub-project 1/2/3a.** Deviasi desain yang disengaja: sifat append-only audit log berarti tidak ada Update/Delete endpoint sama sekali (mengubah/menghapus riwayat consent akan merusak nilai audit trail-nya).
- **Handler/Routes:** 2 endpoint — `POST /recruitment/candidates/:id/consents`, `GET /recruitment/candidates/:id/consents` (tidak ada `PUT`/`DELETE`).
- **DTO:** `CreateCandidateConsentRequest`, `CandidateConsentResponse` (tidak ada `Update...Request` — konsisten dengan append-only design).
- **Test:** +10 test. Total recruitment: **190** (handler 45 + repository 43 + service 102) — diverifikasi via `go test ./internal/modules/recruitment/... -v -count=1 | grep -c "^--- PASS:"`.

**Rencana (sisa G-6 — skipped/deferred, bukan dijadwalkan ulang):**
- `candidates.status` (availability status ACTIVE/BLACKLISTED/dst.) — **skipped**: tidak jelas kebutuhan bisnis, potensi redundan dengan application-level status.
- `candidates.source_id` + master source — **deferred**: `source` tetap teks bebas; membangun source master adalah effort terpisah dan belum diputuskan prioritasnya.

Semua sub-table G-6 yang direncanakan (educations, work_experiences, skills, certifications, documents, consents) sudah selesai per sub-project 3b.

**Ref:** design spec sub-project 1 `docs/superpowers/specs/2026-08-12-candidate-profile-basics-design.md`; design spec sub-project 2 `docs/superpowers/specs/2026-08-12-candidate-skills-certifications-design.md`; design spec sub-project 3a `docs/superpowers/specs/2026-08-12-candidate-documents-design.md`; sub-project 3b task briefs `.superpowers/sdd/2026-08-12-candidate-consents/`; migration 099 (skills/certifications), migration 100 (documents), migration 101 (consents); plan asli §13-§19.

## G-7 ✅ SCREENING & ASSESSMENT — sub-project 1 (screening) ✅ + sub-project 2 (assessment) ✅

**Status: ✅ Selesai (2026-08-13).** Screening (sub-project 1) dan Assessment (sub-project 2) keduanya selesai.

**Yang diimplementasikan (sub-project 1 — Screening):**
- **Migration `102_recruitment_screening`** (pg + mysql, up/down idempotent): tabel baru `application_screenings` — `id, application_id (FK → job_applications, ON DELETE CASCADE), screened_by CHAR(36) NULL, screened_at BIGINT NN default 0, score DECIMAL(5,2) NULL, result VARCHAR(10) NN default 'HOLD', notes TEXT NULL, created_at, updated_at`; index `idx_screen_app`.
- **Model:** `ApplicationScreening` struct — **many-per-application** (pola sama dengan `Interview`, bukan upsert/unique constraint); `result` (`PASS/FAIL/HOLD`) di-enforce via Gin binding `oneof=...`, bukan constraint DB (pola modul ini).
- **Service:** CRUD murni (`CreateApplicationScreening`, `ListApplicationScreenings`, `UpdateApplicationScreening`, `DeleteApplicationScreening`) — validasi application exists; default `result=HOLD` bila kosong. **Deviasi desain disengaja (dikonfirmasi user saat brainstorming):** create/update **TIDAK** memanggil state machine G-5 (`transitionApplicationStatus`) atau menulis `job_application_stage_histories` — screening murni catatan pendukung, recruiter tetap mengubah status aplikasi manual lewat endpoint status existing.
- **Handler/Routes:** `POST/GET /recruitment/applications/:id/screenings`, `PUT/DELETE /recruitment/screenings/:id` (pola G-6 sub-tables; path param `:id`).
- **DTO:** `CreateApplicationScreeningRequest`, `UpdateApplicationScreeningRequest`, `ApplicationScreeningResponse`.
- **Test:** +11 (repository +3: create/find, list, update+delete; service +6: create, default-result, unknown-application guard, list, update+delete, **no-auto-transition assertion** — status aplikasi tetap NEW setelah screening PASS dibuat; handler +2: create, list). Total recruitment setelah sub-project 1: **201**.

**Yang diimplementasikan (sub-project 2 — Assessment):**
- **Migration `103_recruitment_assessment`** (pg + mysql, up/down idempotent): 2 tabel (bukan 3 seperti rencana awal — `assessment_results` **digabung** ke `assessment_participants` karena hasil 1:1 per participant, tidak ada kebutuhan multi-assessor/multi-attempt; keputusan YAGNI dikonfirmasi saat brainstorming):
  - `recruitment_assessments` — sesi/batch tes: `id, requisition_id (FK nullable → job_requisitions, ON DELETE SET NULL), name, type VARCHAR(20) default OTHER, scheduled_at, location, meeting_link, notes, created_at, updated_at`; index `idx_assess_req`. Tidak terikat 1 application — bisa diikuti banyak kandidat sekaligus (mis. "Technical Test Batch Maret").
  - `assessment_participants` — kandidat (via `application_id`) yang mengikuti sesi: `id, assessment_id (FK → recruitment_assessments, CASCADE), application_id (FK → job_applications, CASCADE), status VARCHAR(20) default INVITED, score, result, recommendation, created_at, updated_at`; unique `(assessment_id, application_id)` — satu kandidat tidak bisa didaftarkan dobel ke sesi yang sama; index `idx_partic_assess` + `idx_partic_app`.
- **Model:** `RecruitmentAssessment` + `AssessmentParticipant`; `type` (`TECHNICAL/PSYCHOLOGICAL/COGNITIVE/PERSONALITY/CASE_STUDY/CODING/LANGUAGE/OTHER`) dan `status` (`INVITED/COMPLETED/NO_SHOW`)/`result` (`PASS/FAIL/HOLD`) di-enforce via Gin binding `oneof=...`, bukan constraint DB (pola modul ini).
- **Service:** `CreateAssessment` (validasi requisition bila diisi, default `type=OTHER`) + `List/Get/Update/DeleteAssessment`; `AddAssessmentParticipant` (validasi assessment + application exists, default `status=INVITED`) + `List/Update/DeleteAssessmentParticipant`. **Sama seperti Screening: TIDAK memicu transisi status `job_applications` atau menulis `job_application_stage_histories`** (G-5) — recruiter tetap mengubah status aplikasi manual.
- **Handler/Routes:** `POST/GET /recruitment/assessments`, `GET/PUT/DELETE /recruitment/assessments/:id`, `POST/GET /recruitment/assessments/:id/participants`, `PUT/DELETE /recruitment/assessment-participants/:id`.
- **Test:** +20 (repository +5: assessment CRUD, participant create+list, participant update+delete; service +12: create, default-type, with-requisition, unknown-requisition guard, list, update+delete, add-participant, unknown-assessment guard, unknown-application guard, list-participants, update+delete-participant, **no-auto-transition assertion**; handler +3: create, invalid-type 400, add+list participants). Total recruitment: **221** (handler 50 + repository 51 + service 120).

**Ref:** plan asli §21-§22.

## G-8 ✅ INTERVIEW MULTI-INTERVIEWER & SCORECARD

**Status: ✅ Selesai (2026-08-13).**

**Yang diimplementasikan:**
- **Migration `104_recruitment_interview_scorecard`** (pg + mysql, up/down idempotent): 2 tabel baru —
  - `interviewers` — `id, interview_id (FK → interviews, CASCADE), employee_id, role VARCHAR(50) NULL, created_at`; unique `(interview_id, employee_id)`; index `idx_interviewer_int`. **Melengkapi** `interviews.interviewer_id` existing (pewawancara utama, **dipertahankan** untuk backward compat — bukan breaking change) dengan pewawancara tambahan (HR + User + Manager).
  - `interview_scorecard_items` — `id, interview_id (FK → interviews, CASCADE), criterion VARCHAR(255), weight DECIMAL(5,2) default 0, score DECIMAL(5,2) NULL, notes, created_at, updated_at`; index `idx_scorecard_int`. **Kriteria bebas per interview, tanpa master** (keputusan brainstorming — YAGNI, konsisten dengan `CandidateSkill.Level` yang juga tidak di-enforce skala-nya via DB).
- **Model:** `Interviewer` + `InterviewScorecardItem`.
- **Service:** `AddInterviewer/ListInterviewers/RemoveInterviewer`; `AddScorecardItem/ListScorecardItems/UpdateScorecardItem/DeleteScorecardItem`; **`CompleteInterview`** (endpoint `POST /interviews/:id/complete`) — set `status=COMPLETED` + `completed_at`, dan menghitung `Interview.Score` **existing** (bukan kolom baru) sebagai weighted average scorecard items: `Σ(score×weight)/Σ(weight)`, item tanpa score dilewati; bila tidak ada item berskor, `Score` tidak diubah (recruiter tetap bisa isi manual seperti alur `UpdateInterview` existing).
- **Handler/Routes:** `POST/GET /interviews/:id/interviewers`, `DELETE /interviewers/:id`, `POST/GET /interviews/:id/scorecard-items`, `PUT/DELETE /scorecard-items/:id`, `POST /interviews/:id/complete`.
- **Test:** +14 (repository +5: interviewer create/list/delete, scorecard item CRUD; service +7: add/list/remove interviewer, unknown-interview guard, scorecard item CRUD, **complete-interview weighted average assertion** (Technical 80×70% + Communication 90×30% = 83), **complete-interview-tanpa-item-keeps-score-unset** assertion; handler +2: add+list interviewers, add scorecard item + complete). Total recruitment: **235** (handler 52 + repository 56 + service 127).

**Ref:** plan asli §23-§24.

## G-9 🔶 CANDIDATE MATCHING & REQUISITION COMPETENCY — sub-project 1 ✅ (requirement/competency tables) — sisa: match score

**Status: 🔶 Sebagian (2026-08-13).** Sub-project 1 (fondasi data: requirement + competency tables) selesai; sub-project 2 (algoritma `candidate_match_score` lintas entity) belum dieksekusi — dipisah karena kompleksitasnya jauh lebih besar (agregasi Requirement+Competency+Education+Experience+Skill+Certification+Assessment+Interview).

**Yang diimplementasikan (sub-project 1):**
- **Migration `105_recruitment_requisition_requirements`** (pg + mysql, up/down idempotent): 2 tabel baru —
  - `job_requisition_requirements` — `id, requisition_id (FK → job_requisitions, CASCADE), requirement_type VARCHAR(50) free-text, name, description, minimum_value/maximum_value DECIMAL(10,2) NULL, is_required BOOLEAN default true, sort_order, created_at, updated_at`; index `idx_reqmt_req`. `requirement_type` **sengaja free-text** (bukan enum) — konsisten dengan `job_requisitions.department`/`candidates.source` yang juga free-text di modul ini.
  - `job_requisition_competencies` — `id, requisition_id (FK CASCADE), competency_id (FK → competencies, modul competency, CASCADE), required_level SMALLINT NULL, is_required BOOLEAN default true, weight DECIMAL(5,2) NULL, created_at, updated_at`; unique `(requisition_id, competency_id)`; index `idx_reqcomp_req` + `idx_reqcomp_comp`. `required_level` tidak di-enforce skala DB — pola sama `CandidateSkill.Level`.
- **Model:** `JobRequisitionRequirement` + `JobRequisitionCompetency`.
- **Service:** CRUD murni untuk keduanya — `Create/List/Update/Delete{RequisitionRequirement,RequisitionCompetency}`; validasi requisition exists (dan competency exists untuk yang kedua, reuse `FindCompetencyByID` existing dari G-6 sub-2).
- **Handler/Routes:** `POST/GET /recruitment/requisitions/:id/requirements`, `PUT/DELETE /recruitment/requirements/:id`, `POST/GET /recruitment/requisitions/:id/competencies`, `PUT/DELETE /recruitment/requisition-competencies/:id`.
- **Test:** +14 (repository +6: requirement CRUD, competency CRUD; service +6: requirement create/unknown-requisition-guard/list-update-delete, competency create/unknown-competency-guard/list-update-delete; handler +2: requirement create+list, competency create+list). Total recruitment: **249** (handler 54 + repository 62 + service 133).

**Rencana (sisa — sub-project 2, match score):**
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

## G-12 🔶 FRONTEND RECRUITMENT (halaman penuh — sebagian sudah)

**Status: 🔶 Sebagian.** Hub + Requisitions + Offers + Onboarding sudah ada (lihat §3.2); Candidates/Applications/Interviews/Screening/Assessment masih "Coming soon".

**Sudah ada:**
- [x] **Hub Recruitment** (`Recruitment.vue`): kartu sub-menu ke halaman yang sudah ada + "Coming soon" untuk sisanya.
- [x] **Requisitions**: list + create/edit dialog + submit + priority/requisition_number/opened_at + status badge.
- [x] **Offers**: draft → submit → approval → send → accept/reject/withdraw.
- [x] **Onboarding**: list employee onboarding + status PENDING→IN_PROGRESS→COMPLETED + badge From Offer + create dialog dengan auto-suggest employee dari offer.

**Rencana (sisa — mengikuti pola FE modul lain: bilingual + dark mode + `ConfirmDeleteDialog` + skeleton):**
- **Candidates**: list + profile (personal, contact, resume, education, experience, skills, certifications, applications); internal candidate → current employee/org/position (akses mengikuti permission). Menunggu G-6.
- **Applications**: pipeline board (Applied → Screening → Assessment → Interview → Final Review → Offer → Hired) — drag/drop hanya memanggil transition service backend; detail + stage history (G-5 ✅) + screening/assessment/interview/offer list. Menunggu G-7/G-8.
- **Interviews**: kalender + schedule + interviewer + scorecard + feedback + result. Menunggu G-8.
- **Screening/Assessment**: menunggu G-7.
- **Summary cards** di hub (Open Requisitions, Candidates, Applications, Interviews, Offers, Hires, Time to Hire): menunggu G-11 (analytics).
- **Notifications**: deep-link tipe notifikasi recruitment (setelah G-19/notification infra terpasang).

**Ref:** plan asli §43-§45.

---

# 8. API Plan

## 8.1 Existing (85 endpoint — sudah ada)

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

## Candidate Educations (G-6 sub-project 1)
POST   /api/v1/tenant/recruitment/candidates/{id}/educations
GET    /api/v1/tenant/recruitment/candidates/{id}/educations
PUT    /api/v1/tenant/recruitment/educations/{id}
DELETE /api/v1/tenant/recruitment/educations/{id}

## Candidate Work Experiences (G-6 sub-project 1)
POST   /api/v1/tenant/recruitment/candidates/{id}/work-experiences
GET    /api/v1/tenant/recruitment/candidates/{id}/work-experiences
PUT    /api/v1/tenant/recruitment/work-experiences/{id}
DELETE /api/v1/tenant/recruitment/work-experiences/{id}

## Candidate Skills (G-6 sub-project 2)
POST   /api/v1/tenant/recruitment/candidates/{id}/skills
GET    /api/v1/tenant/recruitment/candidates/{id}/skills
PUT    /api/v1/tenant/recruitment/skills/{id}
DELETE /api/v1/tenant/recruitment/skills/{id}

## Candidate Certifications (G-6 sub-project 2)
POST   /api/v1/tenant/recruitment/candidates/{id}/certifications
GET    /api/v1/tenant/recruitment/candidates/{id}/certifications
PUT    /api/v1/tenant/recruitment/certifications/{id}
DELETE /api/v1/tenant/recruitment/certifications/{id}

## Candidate Documents (G-6 sub-project 3a)
POST   /api/v1/tenant/recruitment/candidates/{id}/documents
GET    /api/v1/tenant/recruitment/candidates/{id}/documents
PUT    /api/v1/tenant/recruitment/documents/{id}
DELETE /api/v1/tenant/recruitment/documents/{id}

## Candidate Consents (G-6 sub-project 3b — append-only, no PUT/DELETE)
POST   /api/v1/tenant/recruitment/candidates/{id}/consents
GET    /api/v1/tenant/recruitment/candidates/{id}/consents

## Applications
GET    /api/v1/tenant/recruitment/applications
POST   /api/v1/tenant/recruitment/applications
GET    /api/v1/tenant/recruitment/applications/{id}
PUT    /api/v1/tenant/recruitment/applications/{id}/status
DELETE /api/v1/tenant/recruitment/applications/{id}
GET    /api/v1/tenant/recruitment/applications/{id}/history

## Application Screenings (G-7 sub-project 1 — many-per-application)
POST   /api/v1/tenant/recruitment/applications/{id}/screenings
GET    /api/v1/tenant/recruitment/applications/{id}/screenings
PUT    /api/v1/tenant/recruitment/screenings/{id}
DELETE /api/v1/tenant/recruitment/screenings/{id}

## Recruitment Assessments + Participants (G-7 sub-project 2 — batch session)
POST   /api/v1/tenant/recruitment/assessments
GET    /api/v1/tenant/recruitment/assessments
GET    /api/v1/tenant/recruitment/assessments/{id}
PUT    /api/v1/tenant/recruitment/assessments/{id}
DELETE /api/v1/tenant/recruitment/assessments/{id}
POST   /api/v1/tenant/recruitment/assessments/{id}/participants
GET    /api/v1/tenant/recruitment/assessments/{id}/participants
PUT    /api/v1/tenant/recruitment/assessment-participants/{id}
DELETE /api/v1/tenant/recruitment/assessment-participants/{id}

## Interviews
GET    /api/v1/tenant/recruitment/interviews
POST   /api/v1/tenant/recruitment/interviews
GET    /api/v1/tenant/recruitment/interviews/{id}
PUT    /api/v1/tenant/recruitment/interviews/{id}
DELETE /api/v1/tenant/recruitment/interviews/{id}
POST   /api/v1/tenant/recruitment/interviews/{id}/complete

## Job Requisition Requirements + Competencies (G-9 sub-project 1)
POST   /api/v1/tenant/recruitment/requisitions/{id}/requirements
GET    /api/v1/tenant/recruitment/requisitions/{id}/requirements
PUT    /api/v1/tenant/recruitment/requirements/{id}
DELETE /api/v1/tenant/recruitment/requirements/{id}
POST   /api/v1/tenant/recruitment/requisitions/{id}/competencies
GET    /api/v1/tenant/recruitment/requisitions/{id}/competencies
PUT    /api/v1/tenant/recruitment/requisition-competencies/{id}
DELETE /api/v1/tenant/recruitment/requisition-competencies/{id}

## Interviewers (G-8 — multi-interviewer, melengkapi interviews.interviewer_id)
POST   /api/v1/tenant/recruitment/interviews/{id}/interviewers
GET    /api/v1/tenant/recruitment/interviews/{id}/interviewers
DELETE /api/v1/tenant/recruitment/interviewers/{id}

## Interview Scorecard Items (G-8 — kriteria bebas per interview)
POST   /api/v1/tenant/recruitment/interviews/{id}/scorecard-items
GET    /api/v1/tenant/recruitment/interviews/{id}/scorecard-items
PUT    /api/v1/tenant/recruitment/scorecard-items/{id}
DELETE /api/v1/tenant/recruitment/scorecard-items/{id}

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
GET    /recruitment/requisitions/{id}/requirements  ← G-9 ✅ (sub-project 1, see §8.1)
GET    /recruitment/requisitions/{id}/competencies  ← G-9 ✅ (sub-project 1, see §8.1)
POST   /recruitment/applications/{id}/stage         ← G-5 (transition + history)
POST   /recruitment/applications/{id}/screen        ← G-7 ✅ (implemented as /applications/{id}/screenings, see §8.1)
GET    /recruitment/candidates/{id}/profile         ← G-6 (educations/experiences/skills/certs/documents)
POST   /recruitment/applications/{id}/assessments   ← G-7 ✅ (implemented as /assessments + /assessments/{id}/participants, see §8.1)
GET    /recruitment/assessments                     ← G-7 ✅
POST   /recruitment/interviews/{id}/complete        ← G-8 ✅
GET    /recruitment/interviews/{id}/scorecards      ← G-8 ✅ (implemented as /interviews/{id}/scorecard-items, see §8.1)
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
- Halaman fungsional yang sudah ada: `Requisitions.vue`, `Offers.vue`, `Onboarding.vue`, `CandidateSearch.vue`, `InternalCandidates.vue`, `RecruitmentAnalytics.vue`. `Recruitment.vue` (hub) menampilkan kartu ke halaman-halaman ini + "Coming soon" untuk sisanya (lihat §3.2, G-12).

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

- ✅ **Sudah ada**: `015_recruitment.sql` (mysql + postgres) — 7 tabel dasar; `093`-`105` (G-1 s.d. G-9 sub-1, lihat §7).
- ⏳ **Rencana** (nomor migration menyesuaikan urutan existing — lanjut dari 106):

```text
M10 Enhance onboarding_task_templates (organization_id, position_id, employment_type nullable)  (G-10)
```

Sudah dieksekusi (bukan lagi rencana): M1 (093-094, requisition approval + enhancement), M3 (095, job_offers), M4 (097, recruitment_stages + stage history), M6 (098-101, candidate sub-tables), M7 (102, application_screenings), M8 (103, recruitment_assessments + assessment_participants — 2 tabel, bukan 3, lihat §G-7), M9 (104, interviewers + interview_scorecard_items, lihat §G-8), M5 (105, job_requisition_requirements + job_requisition_competencies, lihat §G-9). M2 (`candidates.status`/`source_id`) diputuskan **skip/deferred** — lihat §G-6.

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
- [x] Stage transition memiliki history.
- [ ] Screening tersedia.
- [ ] Assessment tersedia.
- [ ] Interview mendukung multi-interviewer + scorecard.
- [x] Candidate memiliki profile terstruktur (education/experience/skills/certification/document/consent) — ✅ semua sub-table G-6 selesai: education/experience/skills/certifications/documents/consents (G-6 sub-project 1+2+3a+3b).
- [ ] Candidate mendukung internal/external.
- [ ] External candidate dapat menjadi Employee.
- [ ] Internal candidate menggunakan Employee Movement.
- [x] Offer accepted dapat membuat onboarding (G-4: External → Employee, Internal → Employee Movement; onboarding FE sudah pakai alur auto-suggest).
- [ ] Candidate dapat dinilai terhadap competency requirement.
- [ ] Recruitment analytics tersedia.
- [ ] Permission lengkap + sinkron (module.go vs seed_rbac: `interview`/`onboard`).
- [ ] Audit trail tersedia.
- [ ] Notification tersedia.
- [ ] Frontend Recruitment selesai (semua halaman) — 🔶 sebagian: Requisitions/Offers/Onboarding ✅, Candidates/Applications/Interviews/Screening/Assessment ❌ (lihat G-12).
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
