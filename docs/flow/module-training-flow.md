# Alur Pengisian Training & Development (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Training & Development** — setup
master data (categories, courses, providers, trainers), perencanaan (plans, needs, requests),
penyelenggaraan (sessions, attendance, assessments), sertifikasi, evaluasi & effectiveness,
serta reports & history — pola runbook seperti
[`module-leave-flow.md`](module-leave-flow.md) & [`module-reimbursement-flow.md`](module-reimbursement-flow.md).

- Plan pengembangan: `module-training-development-plan.md` — ✅ P0–P2 selesai (2026-08-11)
- Lokasi kode: `backend/internal/modules/training/` · `frontend/tenant/src/views/modules/training/`
- Daftar endpoint + contoh curl: [`../api/api-usage-guide.md`](../api/api-usage-guide.md) → §8.2 (tabel Training)

---

## 1. Ringkasan Alur End-to-End

```
SETUP (sekali)                  PERENCANAAN                       PENYELENGGARAAN                 PASCA-TRAINING
┌───────────────────┐   ┌──────────────────────────┐   ┌──────────────────────────────┐   ┌──────────────────────────┐
│ Categories        │   │ Plan: DRAFT → ACTIVE      │   │ Session: DRAFT → SCHEDULED   │   │ Certificate (generate    │
│ Courses           │   │         → ARCHIVED        │   │   → REGISTRATION_OPEN → FULL │   │   dari participant       │
│ Providers         │   │ Need: OPEN → PLANNED      │   │   → IN_PROGRESS → COMPLETED  │   │   COMPLETED + certif.)   │
│ Trainers          │──▶│         → FULFILLED       │──▶│ Attendance + Assessment       │──▶│ Evaluation Answers       │
│                   │   │ Request: DRAFT → SUBMITTED │   │ Participant: completion       │   │ Effectiveness Assessment │
│                   │   │   → PENDING_APPROVAL      │   │   NOT_STARTED → IN_PROGRESS  │   │ Reports + History        │
│                   │   │   → APPROVED / REJECTED   │   │   → COMPLETED / FAILED       │   │                          │
└───────────────────┘   └──────────────────────────┘   └──────────────────────────────┘   └──────────────────────────┘
```

- **Session status:** `DRAFT → SCHEDULED → REGISTRATION_OPEN → FULL → IN_PROGRESS → COMPLETED` · terminal: `CANCELLED`
- **Participant registration:** `NOMINATED → REQUESTED → APPROVED → REGISTERED` · `WAITLISTED`, `CANCELLED`
- **Participant completion:** `NOT_STARTED → IN_PROGRESS → COMPLETED` · `FAILED`
- **Training request:** `DRAFT → SUBMITTED → PENDING_APPROVAL → APPROVED` · terminal: `REJECTED`, `CANCELLED`
- **Training need:** `OPEN → PLANNED → FULFILLED` · terminal: `CANCELLED`
- **Training plan:** `DRAFT → ACTIVE → ARCHIVED`

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Training Category | `training_categories` | Kategori course (kode, nama, deskripsi, aktif) |
| Training Course | `training_courses` | Course/kelas training (kode, nama, durasi, biaya, bersertifikat, vendor eksternal) |
| Course Objective | `training_course_objectives` | Tujuan pembelajaran per course |
| Course Competency | `training_course_competencies` | Kompetensi yang dibutuhkan per course |
| Course Prerequisite | `training_course_prerequisites` | Prasyarat course lain sebelum mengikuti course ini |
| Training Provider | `training_providers` | Penyedia training (internal/eksternal, vendor, kontak) |
| Training Trainer | `training_trainers` | Pelatih/instruktur (nama, spesialisasi, tarif) |
| Training Session | `training_sessions` | Jadwal pelaksanaan course (tanggal, lokasi, kuota, status) |
| Session Trainer | `training_session_trainers` | Penugasan trainer ke session |
| Training Participant | `training_participants` | Pendaftaran karyawan ke session (registrasi + completion) |
| Training Material | `training_materials` | Materi/makalah per session (judul, URL file) |
| Training Attendance | `training_attendances` | Kehadiran per participant per session (check-in, status) |
| Training Assessment | `training_assessments` | Penilaian/post-test per session |
| Assessment Result | `training_assessment_results` | Hasil penilaian per participant |
| Training Evaluation | `training_evaluations` | Evaluasi training (legacy, sebelum P2) |
| Training Evaluation Form | `training_evaluation_forms` | Form evaluasi (nama, deskripsi, aktif) |
| Evaluation Question | `training_evaluation_questions` | Pertanyaan per form (type: RATING/TEXT/SINGLE_CHOICE/MULTIPLE_CHOICE) |
| Evaluation Answer | `training_evaluation_answers` | Jawaban per participant per pertanyaan |
| Effectiveness Assessment | `training_effectiveness_assessments` | Pengukuran efektivitas (before/after/effectiveness score) |
| Training Certification | `training_certifications` | Master sertifikasi (kode, nama, badan penerbit, masa berlaku) |
| Training Certificate | `training_certificates` | Sertifikat terbit dari participant COMPLETED + certification |
| Training Plan | `training_plans` | Rencana training tahunan (kode, nama, tahun, status) |
| Plan Item | `training_plan_items` | Item perencanaan (course, target tanggal, peserta, biaya, prioritas) |
| Training Need | `training_needs` | Kebutuhan training (sumber, deskripsi, status) |
| Training Request | `training_requests` | Pengajuan training (via Central Approval) |
| Training Mandatory | `training_mandatories` | Training wajib per employee/position |
| Session Cost | `training_session_costs` | Biaya per session (kategori, deskripsi, jumlah) |
| Training Document | `training_documents` | Dokumen pendukung per session (judul, URL file) |

---

## 3. TAHAP 1 — SETUP Master Data (dikerjakan sekali)

### A. Training Categories

Menu **Training → Categories** (`/training/categories`, `TrainingCategories.vue`).

- CRUD kategori course: kode, nama, deskripsi, status aktif.
- Endpoint: `GET/POST /categories`, `GET/PUT/DELETE /categories/:id`.
- Permission: `training.settings.create/update/delete` (GET bebas untuk semua dengan `training.view`).

### B. Training Courses

Menu **Training → Courses** (`/training/courses`, `TrainingCourses.vue`).

- CRUD course: kode, nama, deskripsi, durasi (jam), skor minimum, biaya, bersertifikat, vendor eksternal, kategori.
- **Sub-resource per course** (di tab detail):
  - **Objectives** — tujuan pembelajaran: `GET/POST /courses/:id/objectives`, `PUT/DELETE /course-objectives/:id`
  - **Competencies** — kompetensi yang dibutuhkan: `GET/POST /courses/:id/competencies`, `DELETE /course-competencies/:id`
  - **Prerequisites** — prasyarat course lain: `GET/POST /courses/:id/prerequisites`, `DELETE /course-prerequisites/:id`
- Endpoint: `GET/POST /courses`, `GET/PUT/DELETE /courses/:id`.
- Permission: `training.settings.create/update/delete`.

### C. Training Providers

Menu **Training → Providers** (`/training/providers`, `TrainingProviders.vue`).

- CRUD penyedia training: nama, tipe (INTERNAL/EKSTERNAL), vendor, kontak, alamat, catatan.
- Endpoint: `GET/POST /providers`, `GET/PUT/DELETE /providers/:id`.
- Permission: `training.settings.create/update/delete`.

### D. Training Trainers

Menu **Training → Trainers** (`/training/trainers`, `TrainingTrainers.vue`).

- CRUD pelatih: nama, email, telepon, spesialisasi, tarif per jam, bio, status aktif.
- Endpoint: `GET/POST /trainers`, `GET/PUT/DELETE /trainers/:id`.
- Permission: `training.settings.create/update/delete`.

### E. Certification Master (P2)

Menu **Training → Certificates** → tab **Certifications** (`/training/certificates`, `TrainingCertificates.vue`).

- CRUD master sertifikasi: kode, nama, badan penerbit, masa berlaku (bulan), perpanjangan wajib, status aktif.
- Endpoint: `GET/POST /certifications`, `GET/PUT/DELETE /certifications/:id`.

---

## 4. TAHAP 2 — PERENCANAAN (Planning)

### A. Training Plans

Menu **Training → Planning** (`/training/plans`, `TrainingPlans.vue`).

- Buat rencana training tahunan: kode, nama, tahun, deskripsi, status.
- **Status plan:** `DRAFT → ACTIVE → ARCHIVED`
- Expand baris plan untuk melihat & mengelola **Plan Items**:
  - Tambah item: pilih course, target tanggal, target peserta, estimasi biaya, prioritas (LOW/MEDIUM/HIGH/URGENT).
  - Endpoint: `GET/POST /plans/:id/items`, `PUT/DELETE /plan-items/:id`.
- Filter: tahun, status.
- Endpoint: `GET/POST /plans`, `GET/PUT/DELETE /plans/:id`.

### B. Training Needs

Menu **Training → Needs** (`/training/needs`, `TrainingNeeds.vue`).

- Catat kebutuhan training dari berbagai sumber: PERFORMANCE_GAP, ONBOARDING, COMPLIANCE, MANUAL.
- **Status need:** `OPEN → PLANNED → FULFILLED` · terminal: `CANCELLED`
- Isi: sumber, employee, course, deskripsi, alasan, status.
- Endpoint: `GET/POST /needs`, `GET/PUT/DELETE /needs/:id`.
- Integrasi: **Recruitment → Onboarding** menghasilkan Training Need sumber `ONBOARDING` saat onboarding selesai.

### C. Training Requests (via Central Approval)

Menu **Training → Requests** (`/training/requests`, `TrainingRequests.vue`).

- Employee mengajukan training → pilih course + session (atau buat baru) + alasan.
- **Status request:** `DRAFT → SUBMITTED → PENDING_APPROVAL → APPROVED` · terminal: `REJECTED`, `CANCELLED`
- **Submit** → instance approval dibuat di Central Approval (modul `training_request`).
- **Approval** → keputusan dipropagasi via callback:
  - `APPROVED` → participant otomatis dibuat (auto-enroll) bila session ada.
  - `REJECTED` → request ditolak.
- **Cancel** → `CANCELLED` (hanya dari `DRAFT` atau `REJECTED`).
- Endpoint: `GET/POST /requests`, `GET /requests/:id`, `POST /requests/:id/submit`, `POST /requests/:id/cancel`.
- Integrasi: **Central Approval** — modul `training_request`.

---

## 5. TAHAP 3 — PENYELENGGARAAN (Sessions)

### A. Buat & Kelola Session

Menu **Training → Sessions** (`/training/sessions`, `TrainingSessions.vue`).

- Buat session dari course: pilih course, judul, tanggal mulai/akhir, lokasi, kuota, status.
- **Status session:** `DRAFT → SCHEDULED → REGISTRATION_OPEN → FULL → IN_PROGRESS → COMPLETED` · terminal: `CANCELLED`
- **Ubah status** via tombol/tombol aksi — setiap perubahan tercatat.
- **Detail session** (`/training/sessions/:id`, `TrainingSessionDetail.vue`):
  - Tab **Trainers** — tambah/hapus trainer dari session.
  - Tab **Participants** — daftar peserta + update registrasi/completion.
  - Tab **Attendance** — tandai kehadiran per peserta.
  - Tab **Assessments** — buat assessment (nama, tipe) + submit hasil per peserta (skor).
  - Tab **Materials** — tambah materi/makalah (judul, URL file).
  - Tab **Costs** — biaya per session (kategori, deskripsi, jumlah).
  - Tab **Documents** — dokumen pendukung (judul, URL file).
  - Tab **Evaluation** — form evaluasi + pertanyaan + submit jawaban per peserta.
  - Tab **Effectiveness** —ukur efektivitas (before/after/effectiveness score) per peserta.
- Endpoint: `GET/POST /sessions`, `GET/PUT/DELETE /sessions/:id`, `PUT /sessions/:id/status`.

### B. Attendance

Di tab **Attendance** pada detail session:

- Tandai kehadiran per participant: status (`PRESENT`/`ABSENT`/`EXCUSED`/`LATE`), catatan, jam check-in/out.
- Endpoint: `GET/POST /sessions/:id/attendance`, `PUT /attendances/:id`.

### C. Assessments

Di tab **Assessments** pada detail session:

- Buat assessment per session (nama, tipe: PRE_TEST/POST_TEST/FINAL/OTHER).
- Submit hasil per participant: skor, pass/fail, catatan.
- Endpoint: `GET/POST /sessions/:id/assessments`, `POST /assessments/:id/results`.

---

## 6. TAHAP 4 — PESERTA (Participants)

### A. Pendaftaran

- **Manual** — `POST /participants`: pilih session + employee + status registrasi.
- **Otomatis (auto-enroll)** — saat training request APPROVED, participant dibuat otomatis.
- **Status registrasi:** `NOMINATED → REQUESTED → APPROVED → REGISTERED` · `WAITLISTED`, `CANCELLED`

### B. Progress Completion

- Update status completion per participant: `NOT_STARTED → IN_PROGRESS → COMPLETED` · `FAILED`.
- Skor akhir, tanggal selesai, catatan.
- **Hanya participant dengan `COMPLETED` yang bisa generate sertifikat.**
- Endpoint: `GET/PUT /participants/:id`.

### C. Training History

Menu **Training → History** (`/training/history`, `TrainingHistory.vue`).

- Pilih employee → lihat riwayat semua training yang pernah diikuti: course, tanggal, skor, status completion, nomor sertifikat.
- Endpoint: `GET /history?employee_id=...`.

---

## 7. TAHAP 5 — PASCA-TRAINING

### A. Evaluation (Form + Answers)

Di tab **Evaluation** pada detail session:

1. **Pilih/buat form evaluasi** — form bisa dipakai ulang lintas session.
2. **Tambah pertanyaan** — type: `RATING` (bintang 1–5), `TEXT` (bebas), `SINGLE_CHOICE` (pilihan tunggal), `MULTIPLE_CHOICE` (pilihan ganda). Pilihan disimpan di JSON `options`.
3. **Submit jawaban per participant** — participant memilih/mengisi pertanyaan.
4. Endpoint: `GET/POST /evaluation-forms`, `GET/PUT/DELETE /evaluation-forms/:form_id`, `GET/POST /evaluation-forms/:form_id/questions`, `PUT/DELETE /evaluation-questions/:id`, `POST /evaluation-forms/:form_id/participants/:participant_id/answers`.

### B. Effectiveness Assessment

Di tab **Effectiveness** pada detail session:

- Ukur efektivitas training per participant: skor sebelum (`before_score`), skor sesudah (`after_score`), skor efektivitas (`effectiveness_score` — dihitung otomatis).
- Endpoint: `GET/POST /effectiveness`, `GET/PUT/DELETE /effectiveness/:id`.

### C. Generate Certificate

Menu **Training → Certificates** → tab **Issued Certificates** (`/training/certificates`, `TrainingCertificates.vue`).

1. **Pilih participant** yang sudah `COMPLETED` + **certification** dari master.
2. **Generate** → sertifikat dibuat otomatis (nomor unik, tanggal terbit, expiry).
3. Update file URL sertifikat bila perlu.
4. Endpoint: `POST /participants/:id/certificate`, `GET/POST /certificates`, `GET/PUT/DELETE /certificates/:id`.

---

## 8. Ringkasan Status & Transisi

| Entitas | Status | Transisi |
|---|---|---|
| **Session** | `DRAFT → SCHEDULED → REGISTRATION_OPEN → FULL → IN_PROGRESS → COMPLETED` · `CANCELLED` | manual via `PUT /sessions/:id/status` |
| **Participant (registration)** | `NOMINATED → REQUESTED → APPROVED → REGISTERED` · `WAITLISTED` · `CANCELLED` | manual via `PUT /participants/:id` |
| **Participant (completion)** | `NOT_STARTED → IN_PROGRESS → COMPLETED` · `FAILED` | manual via `PUT /participants/:id` |
| **Plan** | `DRAFT → ACTIVE → ARCHIVED` | manual via `PUT /plans/:id` |
| **Need** | `OPEN → PLANNED → FULFILLED` · `CANCELLED` | manual via `PUT /needs/:id` |
| **Request** | `DRAFT → SUBMITTED → PENDING_APPROVAL → APPROVED` · `REJECTED` · `CANCELLED` | submit → approval engine; callback → `APPROVED`/`REJECTED`; cancel → `CANCELLED` |

---

## 9. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Central Approval** | Instance approval untuk training request (modul `training_request`); callback status → auto-enroll participant bila APPROVED |
| **Notification** | Notifikasi hasil approval (`TRAINING_APPROVED` / `TRAINING_REJECTED`) ke employee |
| **Employee** | Sumber data employee untuk participant & history |
| **Organization** | Sumber data organisasi (posisi, departemen) untuk mandatories & filtering |
| **Competency** | Sumber data kompetensi untuk course competencies |
| **Recruitment** | Handoff onboarding selesai → Training Need sumber `ONBOARDING` |
| **Payroll** | 🚫 Belum terintegrasi (cost report tersedia tapi belum push ke payroll) |

---

## 10. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Training (hub) | `Training.vue` | Kartu menu: Courses / Categories / Sessions / Providers / Trainers / Certificates / History / Reports |
| Training → Courses | `TrainingCourses.vue` | CRUD course + tab detail (objectives, competencies, prerequisites) |
| Training → Categories | `TrainingCategories.vue` | CRUD kategori course |
| Training → Providers | `TrainingProviders.vue` | CRUD penyedia training |
| Training → Trainers | `TrainingTrainers.vue` | CRUD pelatih/instruktur |
| Training → Sessions | `TrainingSessions.vue` | Daftar session + CRUD + ubah status |
| Training → Session Detail | `TrainingSessionDetail.vue` | Detail session + tab: Trainers / Participants / Attendance / Assessments / Materials / Costs / Documents / Evaluation / Effectiveness |
| Training → Participants | `TrainingParticipants.vue` | Daftar semua participant + filter session + ubah registrasi/completion |
| Training → Planning | `TrainingPlans.vue` | Rencana training tahunan + expand plan items |
| Training → Requests | `TrainingRequests.vue` | Pengajuan training (submit/cancel via Central Approval) |
| Training → Needs | `TrainingNeeds.vue` | Kebutuhan training dari berbagai sumber |
| Training → Certificates | `TrainingCertificates.vue` | Tab 1: Master certification CRUD · Tab 2: Issued certificates + generate |
| Training → History | `TrainingHistory.vue` | Riwayat training per employee |
| Training → Reports | `TrainingReports.vue` | Dashboard analytics + Participation / Cost / Compliance reports |
| Approvals (generik) | `Approvals.vue` | Approve/Reject training request |
| Workforce Intel (hub) | `WorkforceIntelligence.vue` | Kartu "Training Analysis" → `/training/reports` |
| Career Intel (hub) | `CareerIntelligence.vue` | Kartu "Development Training" → `/training/history` |

---

## 11. Endpoint API Utama

Semua di bawah `/api/v1/tenant/trainings/`.

| Area | Endpoint |
|---|---|
| **Categories** | `GET/POST /categories`, `GET/PUT/DELETE /categories/:id` |
| **Courses** | `GET/POST /courses`, `GET/PUT/DELETE /courses/:id` |
| **Course Sub-resources** | `GET/POST /courses/:id/objectives` · `PUT/DELETE /course-objectives/:id` · `GET/POST /courses/:id/competencies` · `DELETE /course-competencies/:id` · `GET/POST /courses/:id/prerequisites` · `DELETE /course-prerequisites/:id` |
| **Providers** | `GET/POST /providers`, `GET/PUT/DELETE /providers/:id` |
| **Trainers** | `GET/POST /trainers`, `GET/PUT/DELETE /trainers/:id` |
| **Sessions** | `GET/POST /sessions`, `GET/PUT/DELETE /sessions/:id`, `PUT /sessions/:id/status` |
| **Session Trainers** | `GET/POST /sessions/:id/trainers`, `DELETE /session-trainers/:id` |
| **Participants** | `GET/POST /participants`, `GET/PUT/DELETE /participants/:id` |
| **Attendance** | `GET/POST /sessions/:id/attendance`, `PUT /attendances/:id` |
| **Assessments** | `GET/POST /sessions/:id/assessments`, `POST /assessments/:id/results` |
| **Materials** | `GET/POST /materials`, `PUT/DELETE /materials/:id` |
| **Session Costs** | `GET/POST /sessions/:id/costs`, `PUT/DELETE /session-costs/:id` |
| **Documents** | `GET/POST /sessions/:id/documents`, `DELETE /documents/:id` |
| **Evaluations (legacy)** | `GET/POST /evaluations`, `GET/PUT/DELETE /evaluations/:id` |
| **Evaluation Forms** | `GET/POST /evaluation-forms`, `GET/PUT/DELETE /evaluation-forms/:form_id` |
| **Evaluation Questions** | `GET/POST /evaluation-forms/:form_id/questions`, `PUT/DELETE /evaluation-questions/:id` |
| **Evaluation Answers** | `GET /evaluation-answers`, `POST /evaluation-forms/:form_id/participants/:participant_id/answers` |
| **Effectiveness** | `GET/POST /effectiveness`, `GET/PUT/DELETE /effectiveness/:id` |
| **Certifications (master)** | `GET/POST /certifications`, `GET/PUT/DELETE /certifications/:id` |
| **Certificates (issued)** | `GET/POST /certificates`, `GET/PUT/DELETE /certificates/:id`, `POST /participants/:id/certificate` |
| **Plans** | `GET/POST /plans`, `GET/PUT/DELETE /plans/:id` |
| **Plan Items** | `GET/POST /plans/:id/items`, `PUT/DELETE /plan-items/:id` |
| **Needs** | `GET/POST /needs`, `GET/PUT/DELETE /needs/:id` |
| **Requests** | `GET/POST /requests`, `GET /requests/:id`, `POST /requests/:id/submit`, `POST /requests/:id/cancel` |
| **Mandatories** | `GET/POST /mandatories`, `GET/PUT/DELETE /mandatories/:id` |
| **History** | `GET /history?employee_id=` |
| **Reports** | `GET /reports/dashboard` · `GET /reports/participation` · `GET /reports/cost` · `GET /reports/compliance` |

---

## 12. Catatan Penting

- **Setup master** (categories, courses, providers, trainers) dilindungi permission `training.settings.*` — hanya admin yang boleh CRUD. GET tetap terbuka untuk semua dengan `training.view`.
- **Training Request → Approval** menggunakan Central Approval Engine — employee mengajukan, approver menindaklanjuti di halaman **Approvals** generik (bukan di halaman training).
- **Auto-enroll** otomatis saat request APPROVED — participant dibuat dengan status `REGISTERED` + completion `NOT_STARTED`.
- **Certificate hanya bisa di-generate** untuk participant yang sudah `COMPLETED` + memiliki certification di master.
- **Effectiveness score** dihitung otomatis sebagai selisih `after_score − before_score`.
- **Compliance report** menghasilkan status `COMPLETED` / `NOT_COMPLETED` berdasarkan apakah employee sudah mengikuti training mandatory.
- **Evaluation Form** bisa dipakai ulang lintas session — form dibuat sekali, di-assign ke session via `/sessions/:id/evaluation-form`.
- **Gin wildcard constraint** — route evaluation-forms menggunakan `:form_id` (bukan `:id`) untuk menghindari konflik wildcard dengan `/evaluation-forms/:form_id/questions`.
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
