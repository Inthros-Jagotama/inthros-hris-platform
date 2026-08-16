# Development Plan — Competency 360 Module

## 1. Tujuan

Mengembangkan modul Competency menjadi **Competency 360 Assessment** yang menilai kompetensi seorang employee berdasarkan multi-rater feedback, meliputi Self, Superior/Manager, Peer, dan Subordinate.

Modul harus tetap kompatibel dengan competency management yang sudah ada dan menggunakan **Approval Engine** yang sudah tersedia di HRIS sebagai mekanisme approval, tanpa membuat approval workflow baru di dalam modul Competency.

---

## 2. Prinsip Arsitektur

### 2.1 Employee adalah Subject Assessment

Competency 360 dilakukan terhadap **employee**, bukan terhadap organization.

Dalam HRIS saat ini:

- `organization` merepresentasikan position.
- Organization/Position digunakan untuk menentukan competency requirement.
- Employee merupakan subject yang dinilai.
- Rater merupakan employee lain yang memberikan penilaian.

Relasi utama:

```text
Organization / Position
        |
        | competency requirement
        v
Employee
        |
        | assessment subject
        v
Competency 360
        |
        +-- Self
        +-- Superior
        +-- Peer
        +-- Subordinate
```

### 2.2 Existing Competency Management Dipertahankan

Tabel existing tetap menjadi fondasi:

- `competencies`
- `competence_values`
- `competency_values`
- `competency_events`
- `competency_event_targets`
- `competency_scores`
- `competency_score_details`
- `job_family_competencies`

Struktur tersebut tidak dihapus secara langsung. Enhancement dilakukan dengan menyesuaikan bagian yang masih organization-centric dan menambahkan layer assessment/rater/response.

### 2.3 Approval Engine

Approval Engine yang sudah tersedia di HRIS digunakan untuk:

- approval assessment setup bila diperlukan oleh business process;
- approval/finalization hasil assessment bila diperlukan;
- approval development plan jika nanti diintegrasikan dengan development process.

Modul Competency 360 **tidak membuat tabel atau engine approval sendiri**.

---

# 3. Gap Struktur Existing

Berdasarkan `008_competency.sql`, fondasi competency sudah tersedia.

Namun terdapat gap untuk 360 assessment:

| Kebutuhan | Kondisi |
|---|---|
| Competency master | Existing |
| Competency level | Existing, perlu review |
| Competency requirement | Existing melalui Job Family |
| Assessment event/period | Existing |
| Employee sebagai subject | Existing tetapi masih bercampur dengan organization |
| Rater | Belum tersedia |
| Rater type | Belum tersedia |
| Assessment template | Belum tersedia |
| Assessment question/indicator | Belum tersedia |
| Rating scale | Belum eksplisit |
| Individual response | Belum tersedia |
| Rater weighting | Belum tersedia |
| Aggregation 360 | Belum tersedia |
| Self vs others comparison | Belum tersedia |
| Competency gap | Existing |
| Weighted gap | Existing |
| Development linkage | Belum tersedia |

---

# 4. Perubahan Existing Table

## 4.1 `competencies`

Tetap menjadi master competency.

Review/extend agar mendukung:

- active/inactive;
- competency category;
- competency type;
- behavioral/functional distinction;
- competency definition.

Pertimbangkan penambahan:

```text
status
category/type
```

Jika field/cluster yang ada masih digunakan untuk kebutuhan lama, jangan langsung dihapus.

---

## 4.2 `competency_values`

Review fungsi tabel ini.

Saat ini terdapat:

```text
type
level
name
slug
code
description
```

Pastikan tabel dapat menjadi master level/rating yang reusable.

Jika `type + level` terlalu membatasi kebutuhan 360, evaluasi constraint:

```text
uk_compval_type_level
```

agar beberapa competency/framework dapat menggunakan level yang sama tanpa konflik.

---

## 4.3 `competence_values`

Tabel legacy perlu dipertahankan selama masih digunakan oleh existing competency process.

Lakukan analisis dependency sebelum migration.

Target jangka panjang:

```text
competence_values
        |
        | legacy compatibility
        v
existing process

competency_values
        |
        v
new structured competency engine
```

Jangan menghapus tabel legacy pada fase awal enhancement.

---

# 5. Assessment Template

Tambahkan konsep template agar HR dapat membuat konfigurasi assessment.

## 5.1 `competency_assessment_templates`

Contoh:

```text
Standard 360
Leadership 360
Employee 360
Managerial 360
```

Field utama:

```text
id
name
code
description
status
scale_id
created_by
updated_by
created_at
updated_at
```

## 5.2 Template Competencies

Relasikan template dengan competency yang dinilai.

```text
competency_assessment_template_competencies

template_id
competency_id
required_level
weight
sort_order
```

Template tidak menggantikan competency requirement dari position/job family.

Fungsinya menentukan **apa yang dinilai dalam assessment tertentu**.

---

# 6. Assessment Indicator / Question

Tambahkan tabel untuk behavioral indicator atau pertanyaan.

## 6.1 `competency_indicators`

Contoh:

```text
Competency: Leadership

- Memberikan arahan yang jelas
- Mampu mengambil keputusan
- Memberikan feedback
- Mengembangkan anggota tim
```

Field:

```text
id
competency_id
code
statement
description
status
sort_order
```

## 6.2 Template Indicator

Jika satu competency dapat mempunyai indikator berbeda untuk template berbeda:

```text
competency_assessment_template_indicators

template_id
indicator_id
weight
sort_order
```

---

# 7. Rating Scale

Buat rating scale yang reusable.

## `competency_rating_scales`

Contoh:

```text
1 = Sangat Tidak Memenuhi
2 = Tidak Memenuhi
3 = Memenuhi
4 = Melebihi Ekspektasi
5 = Sangat Baik
```

Detail:

```text
competency_rating_scale_items

scale_id
value
label
description
weight
```

Pastikan rating dapat digunakan untuk:

- competency assessment;
- behavioral indicator;
- aggregation;
- report.

---

# 8. Competency Assessment Subject

`competency_event_targets` perlu diarahkan menjadi subject assessment berbasis employee.

Existing:

```text
competency_event_id
organization_id
employee_id
```

Masalah utama:

```text
UNIQUE (competency_event_id, organization_id)
```

tidak cocok karena satu position/organization dapat memiliki banyak employee.

Ubah konsep unique menjadi minimal:

```text
UNIQUE (competency_event_id, employee_id)
```

`organization_id` tetap dapat dipertahankan sebagai snapshot/reference posisi saat assessment.

Rekomendasi konsep:

```text
competency_event
       |
       v
assessment_subject
       |
       +-- employee_id
       +-- organization_id (snapshot/reference)
```

---

# 9. Rater Assignment

Tambahkan:

## `competency_assessment_raters`

Fungsi:

> menentukan siapa yang menilai employee.

Field:

```text
id
competency_event_target_id
rater_employee_id
rater_type
weight
status
assigned_at
submitted_at
created_at
updated_at
```

Rater type:

```text
self
superior
peer
subordinate
```

Pertimbangkan tambahan:

```text
other
```

untuk stakeholder tertentu.

Constraint penting:

- employee tidak boleh menjadi rater untuk dirinya sendiri kecuali `self`;
- rater harus aktif pada saat assignment;
- duplicate rater tidak diperbolehkan pada assessment yang sama;
- rater type harus konsisten dengan relationship/assignment.

---

# 10. Rater Weight

Tambahkan konfigurasi weight per rater type.

Contoh:

```text
Self        10%
Superior    40%
Peer        30%
Subordinate 20%
```

Sebaiknya konfigurasi dapat dibuat di level template:

```text
competency_assessment_template_rater_types

template_id
rater_type
weight
min_rater
max_rater
required
anonymous
```

Validasi:

```text
total weight = 100%
```

---

# 11. Assessment Response

Tambahkan penyimpanan jawaban mentah.

## `competency_assessment_responses`

```text
id
rater_id
indicator_id
rating_value
comment
submitted_at
created_at
updated_at
```

Jika assessment mendukung komentar per competency, tambahkan layer komentar terpisah atau scope yang jelas.

Data response harus immutable setelah submit/finalization, kecuali dilakukan melalui mekanisme correction/reopen yang terkontrol.

---

# 12. Assessment Lifecycle

Workflow utama:

```text
DRAFT
  |
  v
OPEN
  |
  v
RATER ASSIGNMENT
  |
  v
ASSESSMENT IN PROGRESS
  |
  v
SUBMISSION
  |
  v
VALIDATION
  |
  v
APPROVAL ENGINE
  |
  v
CALCULATION
  |
  v
FINALIZED
  |
  v
REPORT
```

Status final perlu disesuaikan dengan convention status pada HRIS.

---

# 13. Approval Engine Integration

Gunakan Approval Engine existing.

Assessment module mengirimkan approval request dengan context:

```text
module = competency
entity = competency_event / assessment result
entity_id = ...
```

Approval engine menangani:

- approval level;
- approver;
- status;
- approval history;
- rejection;
- resubmission.

Competency module hanya bertanggung jawab terhadap business state:

```text
assessment submitted
assessment awaiting approval
assessment approved
assessment rejected
assessment finalized
```

Jangan membuat:

```text
competency_approvals
competency_approval_steps
competency_approval_histories
```

jika fungsi tersebut sudah disediakan Approval Engine.

---

# 14. Calculation Engine

Implementasikan calculation service terpisah.

Input:

```text
Rater responses
+
Rater type
+
Rater weight
+
Competency/indicator weight
```

Output:

```text
competency score
overall score
gap
weighted gap
rater distribution
```

Contoh:

```text
Self        4.0 × 10%
Superior    3.0 × 40%
Peer        3.5 × 30%
Subordinate 3.0 × 20%

Final = 3.25
```

---

# 15. Existing `competency_scores`

Existing table:

```text
competency_scores
```

perlu diubah dari konsep:

```text
score per organization
```

menjadi:

```text
score per employee per assessment event
```

Unique key yang direkomendasikan:

```text
UNIQUE (
    competency_event_id,
    employee_id
)
```

Jangan menggunakan:

```text
UNIQUE (organization_id)
```

karena satu position dapat mempunyai banyak employee.

---

# 16. Existing `competency_score_details`

Pertahankan sebagai hasil agregasi competency.

Konsep:

```text
competency_scores
        |
        +-- Communication
        +-- Leadership
        +-- Teamwork
        +-- Problem Solving
```

`employee_level` menjadi hasil calculation engine.

Tambahkan jika diperlukan:

```text
self_score
superior_score
peer_score
subordinate_score
final_score
```

Namun jangan menyimpan terlalu banyak derived data jika dapat dihitung ulang dengan aman.

Gunakan snapshot/finalized result untuk menjaga konsistensi historical assessment.

---

# 17. Gap Analysis

Calculation menghasilkan:

```text
Required Level
Actual Level
Gap
Weighted Gap
```

Contoh:

```text
Leadership
Required = 4.0
360 Score = 3.2
Gap = -0.8
```

Kelompokkan:

### Strength

```text
actual >= required
```

### Development Area

```text
actual < required
```

---

# 18. Self vs Others Analysis

Tambahkan hasil analisis:

```text
Self Score
Others Score
Perception Gap
```

Contoh:

```text
Leadership

Self   = 4.5
Others = 3.1
Gap    = -1.4
```

Gunakan untuk mendeteksi:

- overestimation;
- underestimation;
- alignment;
- development priority.

---

# 19. Anonymity

Untuk Peer dan Subordinate, dukung anonymous feedback.

Contoh report employee:

```text
Peer Average = 3.4
```

bukan:

```text
Peer A = 3.0
Peer B = 3.5
Peer C = 3.7
```

Jika jumlah rater terlalu sedikit, individual identity harus disembunyikan.

Konfigurasi minimum:

```text
anonymous = true
min_rater = 3
```

---

# 20. Reporting

## Employee Report

Menampilkan:

- overall competency;
- competency score;
- required level;
- gap;
- self vs others;
- strength;
- development area;
- comments/feedback.

## Manager Report

Menampilkan:

- employee competency overview;
- competency gap;
- comparison antar employee jika memiliki permission;
- development priority.

## HR Report

Menampilkan:

- competency distribution;
- organization/position competency gap;
- competency heatmap;
- top strengths;
- top development gaps;
- assessment completion;
- rater completion.

---

# 21. Integration dengan Existing HRIS

## Employee

```text
employees
    |
    v
Competency 360 Subject
```

## Organization / Position

```text
organizations
    |
    v
Position Competency Requirement
```

## Job Family

```text
job_family_competencies
    |
    v
Required Competencies
```

## Approval Engine

```text
Competency Assessment
    |
    v
Approval Engine
```

## Training

Competency gap dapat menjadi input training recommendation:

```text
Competency Gap
      |
      v
Development Need
      |
      v
Training Recommendation
```

## Career

Competency result dapat menjadi input career readiness:

```text
Current Competency
        |
        v
Target Position Competency
        |
        v
Career Gap
```

---

# 22. API

Implementasikan API berdasarkan existing API convention HRIS.

Minimal endpoint:

### Master

```text
GET    /competencies
POST   /competencies
GET    /competencies/{id}
PUT    /competencies/{id}
DELETE /competencies/{id}
```

### Template

```text
GET    /competency/templates
POST   /competency/templates
GET    /competency/templates/{id}
PUT    /competency/templates/{id}
```

### Event

```text
GET    /competency/events
POST   /competency/events
GET    /competency/events/{id}
POST   /competency/events/{id}/open
POST   /competency/events/{id}/close
```

### Rater

```text
GET    /competency/assessments/{id}/raters
POST   /competency/assessments/{id}/raters
DELETE /competency/raters/{id}
```

### Assessment

```text
GET    /competency/my-assessments
GET    /competency/my-assessments/{id}
POST   /competency/my-assessments/{id}/responses
POST   /competency/my-assessments/{id}/submit
```

### Result

```text
GET /competency/employees/{employee}/result
GET /competency/employees/{employee}/gap
GET /competency/employees/{employee}/report
```

---

# 23. Authorization

Gunakan authorization/permission system existing.

Minimal permission:

```text
competency.view
competency.manage

competency_360.view
competency_360.manage
competency_360.create
competency_360.assign_rater
competency_360.assess
competency_360.approve
competency_360.view_result
competency_360.view_report
```

Pastikan employee hanya dapat melihat assessment miliknya sendiri dan hasil yang memang diizinkan.

Rater hanya dapat mengakses assessment yang ditugaskan kepadanya.

---

# 24. Frontend

Menu:

```text
Competency
│
├── Master Competency
├── Competency Framework
├── Assessment Template
├── Assessment Event
├── Rater Management
├── My Assessment
├── Team Assessment
└── Reports
```

### Employee

```text
My Competency 360
    |
    +-- Pending Assessment
    +-- Completed
    +-- My Result
```

### Manager

```text
Team Competency
    |
    +-- Pending
    +-- Employee Result
    +-- Gap Analysis
```

### HR

```text
Competency 360 Management
    |
    +-- Event
    +-- Assignment
    +-- Monitoring
    +-- Approval
    +-- Reports
```

---

# 25. Validation

Backend validation wajib mencakup:

- event harus aktif sebelum assessment dapat diisi;
- employee subject harus valid;
- rater harus valid;
- rater tidak boleh menilai employee yang tidak ditugaskan;
- self hanya boleh menggunakan subject employee sendiri;
- total rater weight harus 100%;
- rating harus berada dalam rating scale;
- response wajib lengkap sebelum submit jika indicator required;
- assessment tidak dapat diubah setelah finalized;
- anonymous result tidak boleh membocorkan identitas rater;
- approval harus menggunakan Approval Engine.

---

# 26. Audit Trail

Catat:

- event dibuat;
- subject ditambahkan;
- rater ditambahkan/dihapus;
- assessment dimulai;
- response disimpan;
- assessment submitted;
- approval requested;
- approval approved/rejected;
- calculation executed;
- result finalized;
- assessment reopened jika fitur tersebut tersedia.

Gunakan audit/logging infrastructure existing jika sudah tersedia.

---

# 27. Backend Implementation

Urutan implementasi:

1. Review existing competency domain/service.
2. Refactor organization-centric score menjadi employee-centric.
3. Buat assessment template.
4. Buat competency indicator.
5. Buat rating scale.
6. Buat rater configuration.
7. Buat rater assignment.
8. Buat assessment response.
9. Implement calculation engine.
10. Implement gap analysis.
11. Integrasikan Approval Engine.
12. Implement finalization.
13. Implement reporting service.
14. Implement authorization.

---

# 28. Database Migration

Migration harus backward-compatible sejauh memungkinkan.

Prioritas:

### Phase 1

Enhance existing:

```text
competencies
competency_values
competency_event_targets
competency_scores
competency_score_details
```

### Phase 2

Create:

```text
competency_assessment_templates
competency_assessment_template_competencies
competency_assessment_template_rater_types
competency_indicators
competency_assessment_template_indicators
competency_rating_scales
competency_rating_scale_items
competency_assessment_raters
competency_assessment_responses
```

### Phase 3

Tambahkan indexes dan constraints.

---

# 29. Testing

## Unit Test

Test:

- rating calculation;
- rater weighting;
- competency weighting;
- gap calculation;
- weighted gap;
- self vs others;
- minimum rater;
- anonymity;
- validation.

## Integration Test

Test:

```text
Event
  -> Subject
  -> Rater
  -> Assessment
  -> Submit
  -> Approval Engine
  -> Calculation
  -> Final Result
```

## Authorization Test

Pastikan:

- employee tidak dapat melihat assessment employee lain;
- rater hanya dapat mengisi assessment yang ditugaskan;
- peer tidak dapat mengubah assessment orang lain;
- manager hanya melihat data sesuai scope permission;
- HR dapat melihat report sesuai permission.

## Frontend Test

Test:

- create template;
- create event;
- assign employee;
- assign rater;
- fill assessment;
- submit;
- approval;
- view result;
- view gap analysis.

---

# 30. Development Phase

## Phase 1 — Analysis & Refactoring

- [x] Review existing competency implementation. (Selesai — §34.1)
- [x] Review relationship Employee → Organization/Position. (Selesai — §34.1: `employee_id` sudah ada di target/scores, tapi unique constraint masih org-centric)
- [x] Review Job Family competency mapping. (Selesai — §34.1: `job_family_competencies` tanpa kolom level/weight & tanpa kode Go)
- [x] Review existing competency score calculation. (Selesai — §34.1: `competency_scores.organization_id` UNIQUE, satu skor per posisi)
- [x] Review Approval Engine API/contract. (Selesai — §34.2: pola performance `ApprovalEngine`, `subscriptionModuleAliases`/`subscriptionModuleSubslots`)
- [x] Tentukan backward compatibility strategy. (Keputusan: §34.7)

### 34.7 Keputusan Fase 1 (Final)

- **Slug approval**: `competency_360_assessment` → alias ke subscription `competency` (§34.2).
- **`competence_values` (legacy)**: dipertahankan apa adanya; tidak ada kode baru yang menyentuhnya (§4.3).
- **Backward compatibility tabel existing**: kolom tidak di-drop; hanya unique constraint `competency_event_targets`/`competency_scores` yang diganti ke employee-centric (§8/§15) — aman karena belum ada jalur pengisian aktif (diverifikasi §34.1).
- **Renumbering migrasi**: plan §34.4 menulis mulai `139`, tapi saat implementasi `139_reimbursement_approve_permission` sudah terpakai → migrasi Competency 360 dimulai dari **`140`** (lihat §34.8).

## Phase 2 — Competency Foundation

- [ ] Enhance competency master.
- [ ] Review competency levels.
- [ ] Implement competency indicators.
- [ ] Implement rating scale.
- [ ] Implement competency framework/template.

## Phase 3 — 360 Assessment

- [ ] Implement assessment subject.
- [ ] Implement rater assignment.
- [ ] Implement rater type.
- [ ] Implement rater weighting.
- [ ] Implement assessment response.
- [ ] Implement submission workflow.

## Phase 4 — Approval

- [ ] Integrate Approval Engine.
- [ ] Create approval request.
- [ ] Handle approved state.
- [ ] Handle rejected state.
- [ ] Handle resubmission.
- [ ] Implement finalization.

## Phase 5 — Calculation

- [ ] Implement rater aggregation.
- [ ] Implement competency aggregation.
- [ ] Implement weighted score.
- [ ] Implement required level.
- [ ] Implement gap.
- [ ] Implement weighted gap.
- [ ] Implement self vs others.
- [ ] Implement overall competency score.

## Phase 6 — Reporting

- [ ] Employee report.
- [ ] Manager report.
- [ ] HR report.
- [ ] Competency gap report.
- [ ] Organization/Position competency heatmap.
- [ ] Assessment completion report.

## Phase 7 — Frontend

- [ ] Master competency.
- [ ] Framework/template.
- [ ] Event management.
- [ ] Rater assignment.
- [ ] My assessment.
- [ ] Manager assessment.
- [ ] Result dashboard.
- [ ] HR reporting.

## Phase 8 — Testing & Hardening

- [ ] Backend unit tests.
- [ ] Backend integration tests.
- [ ] Approval integration tests.
- [ ] Authorization tests.
- [ ] Frontend tests.
- [ ] Calculation edge-case tests.
- [ ] Security testing.
- [ ] Performance testing.

---

# 31. Definition of Done

Competency 360 dianggap selesai apabila:

- Employee dapat menjadi subject assessment.
- Position/Organization dapat menentukan competency requirement.
- HR dapat membuat assessment event.
- HR dapat menentukan template.
- Sistem dapat menentukan/menetapkan rater.
- Self, Superior, Peer, dan Subordinate dapat memberikan assessment.
- Rater weighting dapat dikonfigurasi.
- Response tersimpan secara individual.
- Assessment dapat disubmit.
- Approval menggunakan existing Approval Engine.
- Hasil hanya dihitung setelah proses yang ditentukan selesai/approved.
- Sistem menghasilkan competency score.
- Sistem menghasilkan competency gap.
- Sistem menghasilkan self vs others comparison.
- Employee dapat melihat hasil sesuai permission.
- HR dapat melihat aggregated report.
- Anonymous feedback tidak membocorkan identitas rater.
- Assessment yang sudah finalized tidak dapat dimodifikasi tanpa mekanisme resmi.
- Seluruh business-critical flow memiliki automated tests.

---

# 32. Target Architecture

```text
                    COMPETENCY MASTER
                           |
                           v
                 COMPETENCY FRAMEWORK
                           |
                 +---------+---------+
                 |                   |
                 v                   v
          Position Requirement   360 Template
                 |                   |
                 +---------+---------+
                           |
                           v
                   ASSESSMENT EVENT
                           |
                           v
                   EMPLOYEE SUBJECT
                           |
             +-------------+-------------+
             |             |             |
             v             v             v
           Self        Superior        Peer
                                         |
                                   Subordinate
             |             |             |
             +-------------+-------------+
                           |
                           v
                    RATER RESPONSES
                           |
                           v
                   CALCULATION ENGINE
                           |
                           v
                  COMPETENCY RESULT
                           |
             +-------------+-------------+
             |                           |
             v                           v
        GAP ANALYSIS              360 REPORT
             |
       +-----+------+
       |            |
       v            v
    Training      Career
```

---

## 33. Catatan Desain Penting

1. **Employee adalah subject utama Competency 360.**
2. `organization_id` tetap digunakan untuk reference/snapshot position, bukan sebagai identity utama assessment.
3. `competency_scores` harus dapat menyimpan score banyak employee pada organization/position yang sama.
4. Rater assignment dan response harus menjadi layer terpisah dari final score.
5. Final score merupakan hasil calculation, bukan input langsung.
6. Approval tidak boleh dibuat ulang karena HRIS sudah memiliki Approval Engine.
7. Historical assessment harus tetap konsisten walaupun competency requirement atau position berubah setelah assessment selesai.
8. Perubahan master competency tidak boleh mengubah hasil assessment yang sudah finalized.
9. Gunakan UUID/CHAR(36) mengikuti convention existing HRIS.
10. Integrasi Training dan Career sebaiknya disiapkan melalui `development need`, tetapi implementasinya dapat menjadi fase lanjutan.

---

# 34. Rencana Implementasi BE & FE (Sesuai Struktur Codebase Existing)

> Ditulis setelah eksplorasi kode aktual (`backend/internal/modules/competency/`, migrasi `008_competency.sql`, `backend/internal/modules/performance/`, `backend/internal/modules/approval/service.go`, `frontend/tenant/src/views/modules/competency/`). Bagian §1–§33 di atas adalah desain generik; bagian ini menggroundkan-nya ke kondisi kode nyata per hari ini, mengikuti pola yang sama dipakai `docs/module-attendance-business-travel-development-plan.md` §54. **Belum ada kode yang ditulis** — ini murni rencana.

## 34.1 Kondisi Existing (Faktual, Bukan Asumsi)

Modul `competency` saat ini jauh lebih tipis dari yang tersirat di §3 (tabel "Gap Struktur Existing"):

- **Backend**: `backend/internal/modules/competency/` — 8 file (`model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`, `module.go`, + 3 test file). Semua endpoint di `/api/v1/tenant/competency/*` **hanya CRUD polos** untuk 7 tabel dari migrasi `008_competency.sql`. **Tidak ada** endpoint submit/approve/reject/rater sama sekali — ini genuinely dari nol, bukan enhancement.
- **`competency_event_targets`**: field-nya sudah punya `MissingSelf`/`MissingSuperior`/`MissingPeer`/`MissingSubordinate` (smallint counter) — ini **satu-satunya jejak konsep multi-rater** yang sudah ada, tapi cuma counter, bukan tabel response sungguhan. Unique constraint masih `(competency_event_id, organization_id)` persis seperti yang dikeluhkan §8 — **belum diperbaiki**.
- **`competency_scores`**: `organization_id` masih **UNIQUE** (satu skor per organization/position, bukan per employee) — persis masalah yang dijelaskan §15, **belum diperbaiki**.
- **`job_family_competencies`** (migrasi `002_organization.sql`, bukan `008`): cuma `id, job_family_id, competency_id` — **tidak ada kolom `required_level`/`weight`**, dan **tidak ada satu pun kode Go** (model/repo/service) yang menyentuh tabel ini. Artinya "competency requirement dari Job Family" yang disebut §2.1/§21 **belum benar-benar terhubung ke apapun** — perlu dibangun dari nol, bukan sekadar "direview".
- **Frontend**: `frontend/tenant/src/views/modules/competency/` **cuma ada satu file**, `Competencies.vue` — sidebar (`module.go` baris 66–71) menjanjikan sub-menu Values/Events/Scores yang **tidak ada halamannya sama sekali**. Sama seperti kondisi Reimbursement sebelum dikerjakan (lihat `docs/module-reimbursement-development-plan.md`) — placeholder, bukan modul aktif.
- **Migrasi terbaru di repo saat ini: `138`** (`138_reimbursement_requests_payment_details.sql`). Migrasi Competency 360 baru harus mulai dari **139**.

Implikasi: Phase 1 (§30 "Analysis & Refactoring") di plan generik **bukan sekadar review** — sebagian besar itemnya baru diketahui benar-benar kosong setelah eksplorasi ini, bukan cuma perlu "disesuaikan".

## 34.2 Pola Approval — Reuse Persis dari Performance (KPI/OKR)

Kasus paling mirip dengan kebutuhan Competency 360 (dua checkpoint approval independen pada dokumen yang berkaitan) adalah **Performance module**, bukan modul lain:

```go
// backend/internal/modules/performance/service.go — pola yang harus ditiru
type ApprovalEngine interface {
    CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
    GetActiveFlowIDForModule(ctx context.Context, module string) (string, error)
    // ...
}
const ApprovalModuleKPITarget = "performance_kpi_target"
const ApprovalModuleKPIRealization = "performance_kpi_realization"
```

Untuk Competency 360, definisikan konstanta serupa (nama final terserah tim, contoh):

```text
ApprovalModuleCompetency360Assessment = "competency_360_assessment"   // finalisasi hasil 360 (§13 plan generik)
```

Cukup **satu** checkpoint approval untuk finalisasi (bukan per-rater) — rater submission bukan approval, itu business-state internal modul (§13 plan generik sudah benar soal ini). Kalau nanti dibutuhkan approval terpisah untuk assessment *setup* (event dibuka), tambahkan slug kedua dengan pola yang sama.

Wiring di `main.go` mengikuti persis performance (`SetApprovalEngine` → `RegisterStatusHandler` sebelum `NewModuleWithServices`) — lihat `main.go` sekitar baris 1130–1154 untuk pola yang harus dicontek baris-per-baris.

### ⚠️ Bug Laten yang Sudah Terbukti Terjadi — Jangan Diulang

Saat Business Travel diintegrasikan ke Approval, ditemukan **submit selalu gagal** karena dua map di `backend/internal/modules/approval/service.go` tidak diupdate:

- `subscriptionModuleAliases` (baris ±258–275): slug flow → slug subscription asli. **Wajib** ada entry `"competency_360_assessment": "competency"` — tanpa ini, `ensureModuleSubscribed` akan selalu menolak dengan "module not subscribed" karena `competency_360_assessment` bukan slug subscription asli.
- `subscriptionModuleSubslots` (baris ±190–203): slug subscription → daftar slug flow yang di-unlock di module picker Approval Flow Builder. **Wajib** ada entry `"competency": ["competency_360_assessment"]` — tanpa ini, HR tidak akan pernah bisa membuat flow approval untuk Competency 360 di dropdown, walau module Competency sudah disubscribe.

**Cek saat ini: `competency` belum punya entry di kedua map itu sama sekali** (dikonfirmasi kosong). Ini bukan sesuatu yang perlu "ditemukan lagi" — sudah diketahui sebelum implementasi dimulai, jadi masukkan sebagai satu langkah eksplisit di Phase 4 (§30), bukan ditemukan lewat trial-and-error seperti kasus Business Travel.

## 34.3 Pola Rater Response — Analog Terdekat: Overtime ASSIGNED Flow

Tidak ada modul existing dengan pola "4 rater menilai 1 subject" persis, tapi **Attendance Overtime's dua-alur** (`AttendanceOvertimeRequest.FlowType`: `SELF` vs `ASSIGNED`, migrasi 080, §32b di `docs/module-attendance-plan.md`) adalah analog terdekat untuk konsep "orang lain mengisi form tentang/untuk orang lain, lalu ada approval gate": assignment → pengisian oleh assignee → submit → approval. Pola state machine-nya (`SUBMITTED → PENDING_APPROVAL → APPROVED → WAITING_ACTUAL → ACTUAL_SUBMITTED`) bisa jadi referensi konseptual untuk alur *assignment rater → isi response → submit response*, tapi **tidak bisa dipakai langsung** karena overtime single-assignee, sedangkan 360 butuh multi-rater per subject. Tetap perlu tabel baru (§9/§11 plan generik: `competency_assessment_raters`, `competency_assessment_responses`), bukan reuse tabel overtime.

## 34.4 Urutan Migrasi Konkret (mulai `140`)

Menggabungkan §28 (Database Migration) plan generik dengan penomoran nyata. Catatan: plan awalnya menulis mulai `139`, tapi saat implementasi `139_reimbursement_approve_permission` sudah terpakai → semua nomor digeser +1 (lihat §34.7):

```text
140_competency_job_family_requirement.sql   -- ALTER job_family_competencies: + required_level, + weight
                                              -- (tabel ini SEBELUMNYA tidak dipakai kode apapun — aman diubah)
141_competency_event_target_employee_unique.sql
                                              -- ALTER competency_event_targets: ganti UNIQUE(event,org) → UNIQUE(event,employee)
                                              -- §8 plan generik. BREAKING pada constraint lama — cek dulu data existing
                                              -- (kemungkinan besar masih kosong karena belum ada fitur yang mengisi tabel ini secara aktif)
142_competency_score_employee_unique.sql
                                              -- ALTER competency_scores: drop UNIQUE(organization_id) → UNIQUE(competency_event_id, employee_id)
                                              -- §15 plan generik
143_competency_rating_scales.sql              -- competency_rating_scales, competency_rating_scale_items (§7)
144_competency_assessment_templates.sql        -- competency_assessment_templates, template_competencies,
                                                -- template_rater_types (§5, §10)
145_competency_indicators.sql                   -- competency_indicators, template_indicators (§6)
146_competency_assessment_raters.sql             -- competency_assessment_raters (§9)
147_competency_assessment_responses.sql           -- competency_assessment_responses (§11)
148_competency_assessment_approval_instance.sql    -- + approval_instance_id ke tabel finalisasi assessment
                                                     -- (pola sama migrasi 061/133 — reimbursement/business travel)
```

Semua migrasi butuh pasangan `postgres`+`mysql`, plus `.down.sql`, mengikuti konvensi 100% yang sudah dipakai di seluruh repo (lihat `docs/module-attendance-business-travel-development-plan.md` §54.2 untuk detail konvensi penulisan).

**Catatan risiko migrasi 141/142**: ubah UNIQUE constraint pada tabel yang sudah ada di production. Sebelum eksekusi nanti (di luar scope plan ini), **wajib** cek dulu apakah ada data existing yang akan melanggar constraint baru — kemungkinan besar aman karena `competency_event_targets`/`competency_scores` belum punya jalur pengisian aktif di kode manapun saat ini, tapi ini harus diverifikasi faktual saat implementasi, bukan diasumsikan dari sini.

## 34.8 Nomor Migrasi Aktual (hasil implementasi)

Karena `139_reimbursement_approve_permission` sudah terpakai sebelum Competency 360 dikerjakan, seluruh nomor di §34.4 digeser +1 dan diimplementasikan persis sebagai berikut:

```text
140_competency_job_family_requirement
141_competency_event_target_employee_unique
142_competency_score_employee_unique
143_competency_rating_scales
144_competency_assessment_templates
145_competency_indicators
146_competency_assessment_raters
147_competency_assessment_responses
148_competency_assessment_approval_instance
```

Semua punya pasangan `postgres` + `mysql` + `.down.sql`. Migrasi **tidak dijalankan** dalam implementasi ini (di luar scope — hanya file migrasi yang dibuat).

## 34.5 Frontend — Realita vs §24 Plan Generik

§24 menyebut menu lengkap (Master Competency, Framework, Assessment Event, Rater Management, My Assessment, Team Assessment, Reports) seolah-olah tinggal nambah — padahal **hanya `Competencies.vue` yang ada**, dan itu pun perlu diverifikasi isinya (placeholder atau fungsional — belum dicek isinya, baru dikonfirmasi file-nya ada). Rencana halaman baru:

```text
frontend/tenant/src/views/modules/competency/
├── Competencies.vue          # sudah ada — verifikasi isi & lengkapi jika perlu
├── CompetencyValues.vue       # BARU — master rating scale/level (§4.2)
├── CompetencyEvents.vue        # BARU — event/period management
├── AssessmentTemplates.vue      # BARU — §5
├── RaterAssignment.vue           # BARU — §9, per event
├── MyAssessments.vue              # BARU — inbox rater ("My Assessment" §24)
├── AssessmentResult.vue            # BARU — hasil individual + gap (§17-18)
└── CompetencyReports.vue            # BARU — §20 (Employee/Manager/HR report)
```

Pola implementasi ikuti konvensi yang sudah terbukti: `<script setup>`, panggil `@/services/api` langsung tanpa store/API-client layer terpisah, tab manual (bukan PrimeVue TabView) untuk halaman multi-section seperti `BusinessTravelDetail.vue`/`PayrollRunDetail.vue`.

## 34.6 Urutan Kerja yang Disarankan (Revisi dari §30)

Menyelaraskan Phase generik §30 dengan temuan konkret di atas:

1. **Phase 1 (Analysis)** — sudah **sebagian selesai lewat eksplorasi ini**. Sisanya: putuskan nama final slug approval, putuskan apakah `competence_values` (legacy) benar-benar masih dipakai proses lain sebelum disentuh (§4.3 plan generik minta ini dicek, belum dicek di eksplorasi ini).
2. **Phase 2 (Foundation)**: migrasi `142`–`144` (rating scale, template, indicator) + model/repo/service/handler/routes Go standar, tanpa approval dulu (mirip pola "module skeleton dulu, approval belakangan" yang dipakai Business Travel §54.7 urutan kerja).
3. **Phase 3 (360 Assessment core)**: migrasi `139`–`141` (breaking-safe alterations, verifikasi data existing dulu) + `145`–`146` (rater, response).
4. **Phase 4 (Approval)**: migrasi `147`, tambah entry ke `subscriptionModuleAliases`/`subscriptionModuleSubslots` (§34.2 di atas) **sebelum** testing submit — supaya tidak mengulang siklus debug yang sama seperti Business Travel.
5. **Phase 5 (Calculation)**: service kalkulasi terpisah (§14 plan generik) — murni logic, tidak butuh tabel baru di luar yang sudah dibuat.
6. **Phase 6-7 (Reporting, Frontend)**: FE §34.5 di atas, dikerjakan setelah backend Phase 2-5 punya endpoint yang bisa dites.
7. **Phase 8 (Testing)**: ikuti §29 plan generik, tambahkan test khusus untuk kedua map approval (mirip `approval/module_subscription_test.go`'s `TestService_ListAvailableModules_IncludesKPISubModules`, cukup tiru pola itu untuk Competency).
