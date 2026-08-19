# Alur Rekrutmen (Recruitment) — Dokumentasi Lengkap

Dokumen ini menjelaskan alur bisnis modul **Recruitment** dari awal (pengajuan lowongan) sampai
selesai (kandidat diterima & onboarding), termasuk integrasi dengan Central Approval, Notifikasi,
Employee, Employee Movement, dan Training.

- Referensi plan: G-1 … G-12 (offer/application/candidate) — `module-recruitment-development-plan.md` *(di-archive: `docs/archive/`)*; S-1 … S-7 (strategic layer) — [`module-recruitment-strategic-layer-plan.md`](../archive/module-recruitment-strategic-layer-plan.md)
- Lokasi kode: `backend/internal/modules/recruitment/` · `frontend/tenant/src/views/modules/recruitment/`
- Daftar endpoint + contoh curl: [`../api/api-usage-guide.md`](../api/api-usage-guide.md) → §8.2 (tabel Recruitment) & §8.10 (Succession Gaps → requisition)

---

## 1. Ringkasan Alur End-to-End

```
┌────────────────────────────────────────────────────────────────────────────────────┐
│ 1. Pengajuan Lowongan (Requisition)                                                │
│    DRAFT ──Ajukan──▶ SUBMITTED ──Approval──▶ OPEN ──▶ IN_PROGRESS ──▶ FILLED       │
│                                     └──▶ REJECTED / CANCELLED                      │
└────────────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│ 2. Pipeline Kandidat (Application)                                                 │
│    NEW ─▶ SCREENED ─▶ SHORTLISTED ─▶ INTERVIEWED ─▶ OFFERED ─▶ ACCEPTED            │
│      │         │            │             │           └──▶ REJECTED / WITHDRAWN     │
│      │         └── Screening (hasil screener)                                       │
│      └── Assessment (batch, G-7)  ·  Interview (G-8)                                │
│          Penilaian Kandidat (G-12)  ·  Match Score (advisory)                      │
└────────────────────────────────────────────────────────────────────────────────────┘
                                     │ ACCEPTED
                                     ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│ 3. Offer (G-3)                                                                     │
│    DRAFT ──Ajukan──▶ PENDING_APPROVAL ──Approval──▶ APPROVED ──Kirim──▶ SENT       │
│                                                 └──▶ REJECTED                      │
│    SENT ──Kandidat menerima──▶ ACCEPTED  ·  ──menolak──▶ REJECTED                  │
│    SENT ──Kadaluarsa──▶ EXPIRED · ──ditarik──▶ WITHDRAWN                           │
└────────────────────────────────────────────────────────────────────────────────────┘
                                     │ ACCEPTED
                                     ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│ 4. Penciptaan Karyawan / Mutasi (G-4)                                               │
│    Kandidat EXTERNAL → Employee baru (recruited_from_application_id)               │
│    Kandidat INTERNAL → Employee Movement (hasil seleksi)                           │
└────────────────────────────────────────────────────────────────────────────────────┘
                                     ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│ 5. Onboarding (G-4 / S-7)                                                          │
│    Onboarding selesai ──handoff──▶ Training Need (modul Training)                  │
└────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Job Requisition | `job_requisitions` | Pengajuan kebutuhan lowongan (lowongan) |
| Candidate | `candidates` | Profil kandidat (eksternal/internal) |
| Job Application | `job_applications` | Kandidat melamar ke sebuah requisition (pipeline) |
| ApplicationStageHistory | `application_stage_histories` | Riwayat perubahan status aplikasi |
| JobRequisitionRequirement | `job_requisition_requirements` | Requirement lowongan (pendidikan/pengalaman/jurusan/job family) |
| JobRequisitionCompetency | `job_requisition_competencies` | Kompetensi + bobot + required level per requisition (override) |
| ApplicationAssessment | `application_assessments` | Penilaian kandidat (pendidikan/pengalaman/level kompetensi) + skor |
| Application Screening | `application_screenings` | Hasil penyaringan manual per aplikasi |
| Recruitment Assessment | `recruitment_assessments` | Sesi assessment batch (mis. tes psikologi) |
| Assessment Participant | `assessment_participants` | Kandidat peserta sebuah assessment |
| Interview | `interviews` | Jadwal wawancara + interviewer + scorecard |
| Job Offer | `job_offers` | Penawaran kerja ke kandidat |
| Onboarding | `employee_onboardings` | Proses onboarding employee hasil offer |

---

## 3. Alur Detail per Sub-Proses

### 3.1 Pengajuan Lowongan (Requisition) — S-1 / S-5 / G-1

**Status requisition:** `DRAFT → SUBMITTED → OPEN → IN_PROGRESS → FILLED` · terminal: `REJECTED`, `CANCELLED`

1. **Buat (DRAFT)** — HR membuat requisition. `requested_by` diambil dari user login (auth context).
   - **Reason type** menentukan sumber kebutuhan:
     - `WORKFORCE_GAP` (S-1) — diisi dari Workforce Gap Workforce Intelligence
     - `SUCCESSION_GAP` (S-5) — posisi kunci tanpa successor siap (Career Intelligence)
     - `NEW_POSITION`, `REPLACEMENT`, `EXPANSION` — manual
   - Kode otomatis `REQ-{YYYYMM}-{hex}` di-generate backend.
2. **Ajukan (SUBMITTED)** — tombol "Ajukan" → membuat **instance approval** di Central Approval
   (module `recruitment`). Flow aktif modul recruitment di-auto-resolve.
3. **Approval** — keputusan APPROVED/REJECTED/CANCELLED dipropagasi balik via callback:
   - `APPROVED` → status `OPEN` + `opened_at` di-set (timestamp otomatis).
   - `REJECTED` → status `REJECTED`.
4. **Terisi (FILLED)** — otomatis saat `slots_filled >= slots_available` (di-update saat
   kandidat diterima / offer diterima).
5. **Notifikasi** — ke requester:
   - `REQUISITION_SUBMITTED` saat diajukan
   - `REQUISITION_APPROVED` saat disetujui
   - `REQUISITION_REJECTED` saat ditolak
   Approver menerima `APPROVAL_TASK_ASSIGNED` dari approval engine.

> Catatan: requisition berstatus `DRAFT` dapat diubah & dihapus; hanya yang `DRAFT` yang bisa diajukan.

### 3.2 Requirement & Kompetensi Lowongan — G-9

Halaman **Rekrutmen → Lowongan → Requirements & Competencies** (`RequisitionRequirements.vue`)
mendefinisikan standar penilaian lowongan sebelum kandidat dinilai:

- **Card Requirements** — kebutuhan lowongan:
  - Pendidikan minimum · Pengalaman (tahun) · **Education Majors** (jurusan) · **Job Family**
  - Nama jurusan / job family ditampilkan dari master data (bukan id)
- **Card Competencies** — daftar kompetensi yang wajib dimiliki kandidat, dipecah sesuai pola
  Job Management:
  - **Teknis** (Technical Competency) · **Manajerial** · **Lainnya** (kompetensi di luar kedua cluster)
  - Setiap item: kompetensi + **bobot (%)** + **required level** (skala 1–8 Job Management)
  - Item kompetensi bersifat **read-only** dari sisi requisition (tombol tambah/hapus tidak ada)
- **Sinkronisasi dari Job Management** — tombol "Sinkronisasi" menyalin kompetensi potensi
  Job Management organisasi (bobot + level) ke requisition:
  - Kompetensi baru → `POST` · kompetensi yang sudah ada → `PUT` (**backfill** bobot & level)
  - Dipakai juga sebagai fallback: kalau requisition belum punya override sendiri, penilaian &
    match score otomatis memakai kompetensi Job Management organisasi (per `organization_id`).
- Endpoint: `POST/GET /requisitions/:id/requirements`, `PUT/DELETE /requirements/:id`,
  `POST/GET /requisitions/:id/competencies`, `PUT/DELETE /requisition-competencies/:id`.

### 3.3 Profil Kandidat (Candidate) — G-6

- Kandidat **EXTERNAL** (default): data pribadi (nama, email, telepon, alamat, perusahaan saat ini,
  sumber, resume/portfolio/linkedin) + profil terstruktur:
  - Riwayat Pendidikan · Pengalaman Kerja · Skills · Sertifikasi · Dokumen · Persetujuan (consent)
- Kandidat **INTERNAL** (S-4): menunjuk employee yang sudah ada; hasil seleksi diteruskan ke
  **Employee Movement**, bukan membuat employee baru. Sumber kandidat internal yang *eligible*
  diambil dari Career Intelligence (`/career-intelligence/paths/{id}/eligible-employees`) dan
  diekspos lewat `GET /recruitment/eligible-internal-candidates` (S-4).
- Kode kandidat `CAND-{YYYYMM}-{hex}` di-generate backend.

### 3.4 Pipeline Aplikasi (Application) — G-12

**Status aplikasi:** `NEW → SCREENED → SHORTLISTED → INTERVIEWED → OFFERED` · terminal: `ACCEPTED`, `REJECTED`, `WITHDRAWN`

- Satu kandidat bisa punya banyak aplikasi (satu per requisition).
- Perubahan status **hanya maju** (state machine G-5):
  - non-terminal → non-terminal hanya jika urutan lebih tinggi/equal (tidak bisa mundur)
  - status terminal bisa di-set dari status mana pun
  - terminal = final (tidak ada transisi keluar)
- Setiap transisi valid menulis **Stage History** (dari → ke, oleh siapa, kapan, catatan) dan
  timestamp per status (`screened_at`, `shortlisted_at`, dst.).
- Mengubah status lewat UI menampilkan dialog konfirmasi (alasan penolakan wajib saat REJECTED).

### 3.5 Screening

- Screener mencatat hasil penyaringan per aplikasi: `score`, `result` (`PASS` / `FAIL` / `HOLD`), `notes`.
- Hasil screening menjadi pertimbangan untuk majukan status ke `SCREENED` / `SHORTLISTED`.

### 3.6 Assessment — G-7 sub-project 2

- **Assessment = sesi batch** (nama, jenis, jadwal, lokasi, link, catatan) yang diikuti banyak kandidat.
  Jenis: `TECHNICAL`, `PSYCHOLOGICAL`, `COGNITIVE`, `PERSONALITY`, `CASE_STUDY`, `CODING`, `LANGUAGE`, `OTHER`.
- Kandidat ditambahkan sebagai **participant** berstatus `INVITED`.
- Hasil per peserta diisi oleh assessor: `status` (`INVITED` / `COMPLETED` / `NO_SHOW`),
  `score`, `result` (`PASS` / `FAIL` / `HOLD`), `recommendation`.
- UI: menu **Rekrutmen → Assessment** (buat/edit/hapus sesi) + tab Assessment di Detail Aplikasi
  (lihat keikutsertaan kandidat).

### 3.7 Interview — G-8

- Jadwal wawancara per aplikasi: interviewer, stage, jadwal, durasi, lokasi, meeting link.
- Status interview: `SCHEDULED → COMPLETED` · `CANCELLED`, `RESCHEDULED`.
- Mendukung **multi-interviewer** + **scorecard** (kriteria & bobot, nilai per kriteria).
- UI: tab Interview di Detail Aplikasi (buat jadwal + kelola interviewer/scorecard + tandai selesai).

### 3.8 Penilaian Kandidat — G-12

Tab **Penilaian** di Detail Aplikasi (sebelum tab Match Score) adalah penilaian terstruktur
kandidat terhadap kebutuhan requisisi (`RequisitionRequirements.vue` sebagai acuan):

- **Pendidikan** — requirement (mis. Strata 1) + radio *Sesuai / Tidak Sesuai* + catatan
- **Pengalaman** — requirement (mis. 6–8 Tahun) + radio *Sesuai / Tidak Sesuai* + catatan
- **Kompetensi** — daftar kompetensi requirement (nama + required level + bobot) dengan
  dropdown **level kandidat** (Lv.1–8, nama level dari Job Management, bisa dikosongkan)
- **Skor otomatis** (dihitung **server-side**, klien tidak bisa mengatur):
  Pendidikan **20%** + Pengalaman **30%** + Kompetensi **50%** (rasio `level kandidat /
  required level` berbobot — pola Match Score)
- Disimpan 1 baris per aplikasi di `application_assessments` (upsert); level kompetensi
  & breakdown tersimpan JSON, skor selalu dihitung ulang.
- Endpoint: `GET/PUT /applications/:id/assessment`.
- Level yang diisi di sini menjadi sumber **Match Score** bila sudah ada (fallback ke skill
  kandidat) — lihat 3.9.

### 3.9 Match Score (advisory)

- `GET /applications/:id/match-score` menghitung kecocokan kandidat vs posisi
  (requirement/kompetensi) secara on-the-fly — bersifat **advisory**, tidak mengubah status.
- Sumber level kandidat (G-12): level dari **Penilaian Kandidat** (tab Penilaian)
  didahulukan bila sudah diisi; fallback ke **skill kandidat** (tab Skills di Detail Kandidat).
- Kolom **Score** di daftar aplikasi (halaman Applications) menampilkan match score ini.

### 3.10 Offer — G-3

**Status offer:** `DRAFT → PENDING_APPROVAL → APPROVED → SENT → ACCEPTED` · terminal: `REJECTED`, `EXPIRED`, `WITHDRAWN`

1. **Buat (DRAFT)** — offer untuk sebuah aplikasi: tipe pekerjaan, gaji, tunjangan, benefits, tanggal mulai & berlaku.
2. **Ajukan** → instance approval module `recruitment_offer` (flow terpisah dari requisition).
3. **Approved** → offer bisa **dikirim** ke kandidat (`SENT`).
4. **Kandidat merespons:**
   - Terima → `ACCEPTED` (+ `accepted_at`)
   - Tolak → `REJECTED`
   - Lewat masa berlaku → `EXPIRED` · ditarik HR → `WITHDRAWN`

### 3.11 Offer Diterima → Karyawan / Mutasi — G-4

Saat offer **ACCEPTED**, sistem otomatis (dari `AcceptOffer`):

1. Aplikasi terkait maju ke **ACCEPTED** (state machine).
2. Requisition `slots_filled++`; jika penuh → requisition **FILLED** (guard anti double-count).
3. Membuat entitas berdasarkan jenis kandidat:
   - **EXTERNAL** → **Employee baru** dengan referensi `recruited_from_application_id`
     (traceability aplikasi → karyawan).
   - **INTERNAL** → **Employee Movement** hasil seleksi.
4. Notifikasi hasil.

### 3.12 Onboarding — G-4 / S-7

- Employee hasil offer menjalani **onboarding** (task template).
- Saat onboarding **selesai** → handoff ke modul **Training**:
  `CreateOnboardingNeed` membuat Training Need bersumber `ONBOARDING` (S-7).
- UI: menu **Rekrutmen → Onboarding** (daftar + status).

### 3.13 Recruitment Analytics — G-11

- `GET /recruitment/analytics/summary?from=&to=` menghitung metrik pipeline secara on-the-fly
  (tanpa migration baru): total requisition/application/candidate, **time-to-hire**,
  **offer-acceptance-rate**, **application-conversion-rate**, dan **source-conversion**.
- UI: kartu summary di hub **Rekrutmen** (`Recruitment.vue`) — fail-silent (kartu tidak tampil
  bila endpoint error).
- S-3 menambahkan metrik Time to Hire/Fill, Offer Acceptance Rate, dan Source Conversion pada
  analisis Workforce Intelligence yang mengonsumsi data pipeline ini.

---

## 4. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Central Approval** | Instance approval requisition (`recruitment`) & offer (`recruitment_offer`); callback status → requisition/offer; task per approver (USER/ROLE) |
| **Notification** | Notifikasi ke requester: `REQUISITION_SUBMITTED` / `REQUISITION_APPROVED` / `REQUISITION_REJECTED` (best-effort); approver menerima `APPROVAL_TASK_ASSIGNED` dari approval engine. **Catatan:** offer belum punya notifikasi tersendiri (OFFER_* belum diimplementasikan) |
| **Employee** | Membuat employee baru saat offer eksternal diterima |
| **Employee Movement** | Membuat movement saat offer internal diterima |
| **Training** | Handoff onboarding selesai → Training Need (S-7) |
| **Workforce Intelligence** | Sumber `WORKFORCE_GAP` requisition (S-1) |
| **Career Intelligence** | Sumber `SUCCESSION_GAP` (S-5) + eligible internal candidate (S-4) |

---

## 5. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Rekrutmen → Lowongan | `Requisitions.vue` | CRUD + ajukan + status/approval |
| Rekrutmen → Lowongan → Requirements & Competencies | `RequisitionRequirements.vue` | Requirement + kompetensi (teknis/manajerial) + sinkronisasi Job Management (G-9) |
| Rekrutmen → Kandidat | `Candidates.vue` / `CandidateDetail.vue` | Profil kandidat + sub-profil (G-6) |
| Rekrutmen → Aplikasi | `Applications.vue` / `ApplicationDetail.vue` | Pipeline + daftar (kolom score) + detail (history/screening/assessment/interview/penilaian/match score) |
| Rekrutmen → Kandidat Internal | `InternalCandidates.vue` | Eligible via career path (S-4) — `GET /eligible-internal-candidates` |
| Rekrutmen → Assessment | `Assessments.vue` | Sesi assessment batch (G-7) |
| Rekrutmen → Offer | `Offers.vue` | Offer + approval + kirim/terima |
| Rekrutmen → Onboarding | `Onboarding.vue` | Employee hasil offer + handoff training (S-7) |
| Rekrutmen (hub) | `Recruitment.vue` | Kartu menu + **summary cards analytics** (G-11, fail-silent) |

> Cari kandidat eksternal (posisi/filter) ada di modul **Workforce Intelligence**
> (`CandidateSearch.vue`), bukan di halaman recruitment — lihat §4.

---

## 6. Endpoint API Utama

Semua di bawah `/api/v1/tenant/recruitment/`.

| Area | Endpoint |
|---|---|
| Requisition | `POST/GET /requisitions`, `GET/PUT/DELETE /requisitions/:id`, `POST /requisitions/:id/submit` |
| Requirement & Kompetensi | `POST/GET /requisitions/:id/requirements`, `PUT/DELETE /requirements/:id`, `POST/GET /requisitions/:id/competencies`, `PUT/DELETE /requisition-competencies/:id` |
| Candidate | `POST/GET /candidates`, `GET/PUT/DELETE /candidates/:id`, sub-profil `…/educations`, `…/work-experiences`, `…/skills`, `…/certifications`, `…/documents`, `…/consents` |
| Application | `POST/GET /applications`, `PUT /applications/:id/status`, `GET /applications/:id/history`, `GET/PUT /applications/:id/assessment`, `GET /applications/:id/match-score`, `POST/GET /applications/:id/screenings`, `PUT/DELETE /screenings/:id` |
| Internal Candidates | `GET /eligible-internal-candidates` (S-4 — kandidat internal eligible via career path) |
| Assessment | `POST/GET /assessments`, `GET/PUT/DELETE /assessments/:id`, `POST/GET /assessments/:id/participants`, `PUT/DELETE /assessment-participants/:id` |
| Analytics | `GET /analytics/summary` (G-11 — counts + time-to-hire + offer-acceptance-rate + conversion) |
| Interview | `POST/GET /interviews`, `GET/PUT/DELETE /interviews/:id`, `POST/GET /interviews/:id/interviewers`, `POST/GET /interviews/:id/scorecard-items` |
| Offer | `POST/GET /offers`, `GET/PUT/DELETE /offers/:id`, `POST /offers/:id/submit`, `POST /offers/:id/send`, `POST /offers/:id/accept`, `POST /offers/:id/reject`, `POST /offers/:id/withdraw` |
| Onboarding | `POST/GET /employee-onboardings`, `GET/PUT/DELETE /employee-onboardings/:id`, `GET /employee-onboardings/:id/task-items`, `POST /onboarding-task-items` + `PUT/DELETE /onboarding-task-items/:id`, task templates `POST/GET /onboarding-task-templates` + `PUT/DELETE /onboarding-task-templates/:id` |

---

## 7. Catatan Penting

- **Draft vs Ajukan**: requisition & offer yang baru dibuat berstatus DRAFT dan **belum masuk
  approval/notifikasi** sampai tombol "Ajukan" ditekan.
- **Status aplikasi tidak bisa mundur** (state machine G-5) — ditolak oleh backend.
- **Slots filled** dihitung otomatis & anti double-count (satu kandidat = satu slot).
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
