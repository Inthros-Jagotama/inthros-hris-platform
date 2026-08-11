# Training & Development Management — Operational Development Plan

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

# 2. Kondisi Database Existing

File `016_training.sql` saat ini memiliki 7 tabel:

1. `training_categories`
2. `training_courses`
3. `training_sessions`
4. `training_participants`
5. `training_materials`
6. `training_evaluations`
7. `training_certificates`

Struktur existing sudah mencakup catalog, session, participant, attendance sederhana, score, material, evaluation, dan certificate.

Namun masih terdapat beberapa keterbatasan untuk operational Training & Development yang lengkap:

- `training_courses.external_vendor` masih berupa text;
- trainer masih disimpan sebagai `trainer_name`;
- session hanya menggunakan `start_date` dan `end_date`;
- attendance hanya berada pada `training_participants.attendance_status`;
- score hanya satu field;
- belum ada enrollment/request workflow;
- belum ada assessment detail;
- belum ada learning objective;
- belum ada relasi course dengan competency;
- belum ada prerequisite;
- belum ada training provider master;
- belum ada training plan;
- belum ada training request;
- belum ada mandatory training;
- belum ada training cost detail;
- evaluation masih satu rating + feedback;
- belum ada training effectiveness;
- certificate belum memiliki master certification;
- belum ada struktur training history yang lebih kaya.

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

---

# 4. In-House dan External Training

Training harus mendukung dua bentuk penyelenggaraan:

```text
IN_HOUSE
EXTERNAL
```

## 4.1 In-House

Contoh:

```text
Course:
Leadership Basic

Delivery:
IN_HOUSE

Trainer:
Employee internal

Location:
Training Room A
```

## 4.2 External

Contoh:

```text
Course:
Leadership Advanced

Delivery:
EXTERNAL

Provider:
ABC Training Institute

Trainer:
External Trainer

Location:
ABC Training Center
```

Flow database tetap sama:

```text
Course
   ↓
Session
   ↓
Participant
   ↓
Attendance
   ↓
Assessment
   ↓
Evaluation
   ↓
Certificate
```

Yang berubah adalah provider, trainer, lokasi, biaya, dan informasi penyelenggara.

---

# 5. Target Menu

```text
Training & Development
│
├── Dashboard
│
├── Training Catalog
│   ├── Categories
│   ├── Courses
│   ├── Objectives
│   ├── Competencies
│   └── Prerequisites
│
├── Training Planning
│   ├── Training Plans
│   ├── Training Needs
│   └── Training Requests
│
├── Training Sessions
│   ├── Scheduled
│   ├── Ongoing
│   └── Completed
│
├── Participants
├── Attendance
├── Assessments
├── Evaluations
├── Certifications
├── Training History
└── Reports
```

---

# 6. Database Enhancement

Semua ID mengikuti pola existing menggunakan:

```sql
CHAR(36)
```

dan tetap menggunakan:

```sql
deleted_at
created_at
updated_at
```

sesuai pola file existing.

---

## 6.1 `training_categories`

**Existing sudah cukup.**

Tetap:

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

`external_vendor` sebaiknya **dihapus dari konsep master course** atau tidak lagi digunakan sebagai sumber utama.

Alasannya satu course dapat diselenggarakan:

```text
Session A → In-house
Session B → External Provider A
Session C → External Provider B
```

Vendor adalah karakteristik session/provider, bukan course.

### Tambahkan

```text
course_type
delivery_type
is_mandatory
```

Contoh:

```text
course_type:
TECHNICAL
SOFT_SKILL
COMPLIANCE
MANAGEMENT
CERTIFICATION
OTHER
```

`delivery_type` dapat berupa default/preferred:

```text
IN_HOUSE
EXTERNAL
BOTH
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

Contoh:

```text
Leadership Basic

1. Memahami leadership fundamentals
2. Mampu melakukan delegation
3. Mampu melakukan coaching
```

---

# 9. `training_course_competencies`

Menghubungkan course dengan competency.

```text
id
course_id
competency_id
target_level
deleted_at
created_at
updated_at
```

Contoh:

```text
Leadership Training
    ↓
Leadership Competency
    ↓
Target Level 4
```

Relasi ini menjadi salah satu sumber data untuk Career Intelligence.

---

# 10. `training_course_prerequisites`

Menyimpan syarat sebelum mengikuti course.

Minimal:

```text
id
course_id
prerequisite_type
prerequisite_id
is_required
deleted_at
created_at
updated_at
```

Contoh prerequisite:

```text
Course
Competency
Certification
Experience
```

---

# 11. `training_providers`

Untuk external training provider.

```text
id
code
name
type
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

Contoh:

```text
ABC Training Institute
Provider Type: EXTERNAL
```

Untuk in-house:

```text
provider_id = NULL
```

---

# 12. `training_trainers`

Trainer tidak lagi hanya berupa `trainer_name`.

```text
id
type
employee_id
provider_id
name
email
phone
bio
is_active
deleted_at
created_at
updated_at
```

Type:

```text
INTERNAL
EXTERNAL
```

Internal:

```text
employee_id = employee UUID
provider_id = NULL
```

External:

```text
employee_id = NULL
provider_id = provider UUID
```

---

# 13. `training_session_trainers`

Satu session dapat memiliki lebih dari satu trainer.

```text
id
session_id
trainer_id
role
deleted_at
created_at
updated_at
```

Contoh:

```text
Trainer 1 → Main Trainer
Trainer 2 → Assistant Trainer
```

---

# 14. `training_sessions`

Existing:

```text
id
course_id
session_code
trainer_name
location
start_date
end_date
max_quota
status
```

### Enhancement

Tambahkan:

```text
delivery_type
provider_id
start_datetime
end_datetime
meeting_url
registration_deadline
```

Hindari penggunaan `trainer_name` sebagai sumber utama karena trainer sekarang menggunakan `training_trainers`.

Status minimal:

```text
DRAFT
SCHEDULED
REGISTRATION_OPEN
FULL
ONGOING
COMPLETED
CANCELLED
```

Delivery:

```text
IN_HOUSE
EXTERNAL
ONLINE
HYBRID
```

Catatan:
`IN_HOUSE` dan `EXTERNAL` adalah tipe penyelenggara. `ONLINE`, `HYBRID`, dan `ONSITE` lebih tepat sebagai mode delivery jika diperlukan. Jika ingin fleksibel, pisahkan:

```text
provider_type:
IN_HOUSE
EXTERNAL

delivery_mode:
ONSITE
ONLINE
HYBRID
SELF_PACED
```

---

# 15. Training Request

Tambahkan:

```text
training_requests
```

Field:

```text
id
employee_id
course_id
requested_date
reason
priority
status
approval_instance_id
deleted_at
created_at
updated_at
```

Status:

```text
DRAFT
SUBMITTED
APPROVED
REJECTED
CANCELLED
```

Approval menggunakan **Module Approval** yang sudah menjadi module terpusat.

Flow:

```text
Employee / Manager
      ↓
Training Request
      ↓
Approval
      ↓
Approved
      ↓
Enrollment
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
status
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

Contoh:

```text
Training Plan 2026

Q1:
Leadership
Technical

Q2:
Management
Compliance
```

---

# 17. Training Needs

Tambahkan:

```text
training_needs
```

Tujuannya mencatat kebutuhan training secara operasional.

Field:

```text
id
employee_id
organization_id
position_id
course_id
reason
priority
source_type
source_id
status
deleted_at
created_at
updated_at
```

`source_type` dapat menunjuk sumber kebutuhan seperti:

```text
MANUAL
PERFORMANCE
COMPETENCY
CAREER
SUCCESSION
COMPLIANCE
WORKFORCE
```

Training Need bukan Intelligence module. Ini adalah kebutuhan training yang masuk ke operational Training.

---

# 18. Enrollment / Participant

Existing `training_participants` perlu diperkuat.

Existing:

```text
id
session_id
employee_id
attendance_status
score
completed_at
```

### Ubah konsep

`training_participants` menjadi enrollment/participant record.

Tambahkan:

```text
registration_status
registered_at
approved_at
completion_status
completion_date
final_score
passed
remarks
```

Status registration:

```text
NOMINATED
REQUESTED
APPROVED
REGISTERED
WAITLISTED
CANCELLED
```

Completion:

```text
NOT_STARTED
IN_PROGRESS
COMPLETED
FAILED
```

Tambahkan unique constraint untuk mencegah employee terdaftar dua kali pada session yang sama.

---

# 19. Training Attendance

Jangan menggunakan:

```text
training_participants.attendance_status
```

sebagai satu-satunya attendance.

Untuk training multi-day/session, tambahkan:

```text
training_attendances
```

Field:

```text
id
participant_id
attendance_date
check_in
check_out
status
remarks
deleted_at
created_at
updated_at
```

Status:

```text
PRESENT
ABSENT
LATE
EXCUSED
```

Contoh:

```text
Day 1 → PRESENT
Day 2 → PRESENT
Day 3 → ABSENT
Day 4 → PRESENT
```

`attendance_status` pada participant dapat dipertahankan sementara untuk compatibility, tetapi sumber detail attendance harus `training_attendances`.

---

# 20. Training Materials

Existing:

```text
id
session_id
title
file_url
file_type
sort_order
```

Sudah cukup sebagai basic implementation.

Enhancement opsional:

```text
description
is_required
available_from
```

Material bisa berupa:

```text
PDF
PPT
DOC
VIDEO
LINK
OTHER
```

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
type
max_score
passing_score
attempt_limit
is_required
deleted_at
created_at
updated_at
```

Type:

```text
PRE_TEST
POST_TEST
FINAL
PRACTICAL
OTHER
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

Contoh:

```text
Pre Test  = 60
Post Test = 90
Improvement = +30
```

---

# 22. Training Evaluation

Existing:

```text
training_evaluations
```

hanya:

```text
rating
feedback
```

Untuk versi lengkap, buat:

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
question_type
sort_order
is_required
```

Question type:

```text
RATING
TEXT
SINGLE_CHOICE
MULTIPLE_CHOICE
```

## Answers

```text
id
question_id
participant_id
answer
```

Contoh:

```text
Trainer knowledge       5
Material quality        4
Relevance               5
Facility                4
Overall satisfaction    5
```

Existing `training_evaluations` dapat dipertahankan sementara atau dimigrasikan ke struktur form baru.

---

# 23. Training Effectiveness

Tambahkan:

```text
training_effectiveness_assessments
```

Tujuan:

Mengukur apakah training memberikan dampak setelah training selesai.

Field:

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

Contoh:

```text
Competency Before: Level 2
Competency After : Level 4
```

Assessment effectiveness dapat dilakukan setelah periode tertentu, misalnya 30/60/90 hari.

---

# 24. Certification

Existing:

```text
training_certificates
```

sudah menyimpan:

```text
participant_id
certificate_no
issued_date
expiry_date
```

Untuk versi lengkap, tambahkan master certification:

```text
training_certifications
```

Field:

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

Kemudian certificate dapat mengacu ke certification.

Tambahkan pada certificate:

```text
certification_id
certificate_file_url
```

---

# 25. Mandatory Training

Tambahkan:

```text
training_mandatories
```

Tujuan:

Menentukan training wajib berdasarkan:

```text
Organization
Position
Employment Status
```

Field:

```text
id
course_id
organization_id
position_id
employment_status_id
due_days
validity_period_month
is_active
deleted_at
created_at
updated_at
```

Contoh:

```text
Position:
Safety Officer

Mandatory:
K3 Training
```

Dashboard:

```text
Required
Completed
Pending
Overdue
```

---

# 26. Training Cost

Jangan hanya menggunakan `training_courses.cost` sebagai biaya final.

Biaya aktual dapat berasal dari:

```text
Trainer
Provider
Venue
Material
Certification
Travel
Accommodation
Other
```

Tambahkan:

```text
training_session_costs
```

Field:

```text
id
session_id
cost_type
description
amount
deleted_at
created_at
updated_at
```

Contoh:

```text
Trainer        10,000,000
Venue           5,000,000
Material        1,000,000
Certification   1,000,000
Travel          3,000,000
```

Kemudian dapat dihitung:

```text
Total Cost
Cost per Participant
Cost per Training Hour
```

---

# 27. Training Documents

Untuk dokumen administratif:

```text
training_documents
```

Field:

```text
id
session_id
document_type
file_name
file_url
uploaded_by
deleted_at
created_at
updated_at
```

Contoh:

```text
Proposal
Quotation
Attendance Sheet
Invoice
Contract
Training Report
```

---

# 28. Training History

Tidak perlu membuat tabel `employee_training_histories` jika seluruh histori dapat dibangun dari:

```text
training_participants
+
training_sessions
+
training_courses
+
training_assessments
+
training_certificates
```

Employee dapat melihat:

```text
Training History

Leadership Basic
10 Jan 2026
Completed
Score 90
Certificate

K3 Training
15 Mar 2026
Completed
Score 88
Certificate
```

Jika ada kebutuhan snapshot historis, baru tambahkan field snapshot yang relevan.

---

# 29. Business Flow

## 29.1 In-House

```text
Training Need
     ↓
Training Plan
     ↓
Course
     ↓
Create Session
     ↓
Delivery = IN_HOUSE
     ↓
Internal Trainer
     ↓
Open Registration
     ↓
Participant
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

## 29.2 External

```text
Training Need
     ↓
Training Plan
     ↓
Course
     ↓
Create Session
     ↓
Delivery = EXTERNAL
     ↓
Provider
     ↓
External Trainer
     ↓
Open Registration
     ↓
Participant
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

---

# 30. Training Request Flow

```text
Employee / Manager
        ↓
Training Request
        ↓
Approval Module
        ↓
Approved
        ↓
Register to Session
        ↓
Training Participant
```

Untuk request yang ditolak:

```text
Rejected
```

Tidak boleh menjadi participant aktif.

---

# 31. Training Planning Flow

```text
Training Need
      ↓
Training Plan
      ↓
Course
      ↓
Session
      ↓
Participant
```

Planning tidak wajib untuk setiap session jika organisasi membutuhkan training ad-hoc.

---

# 32. Validation

## Course

- category harus aktif;
- code unik;
- course tidak boleh digunakan jika inactive.

## Session

- course harus aktif;
- start datetime <= end datetime;
- registration deadline < start datetime;
- quota > 0;
- provider wajib untuk external;
- provider tidak wajib untuk in-house;
- trainer harus sesuai type;
- session code unik.

## Participant

- employee harus aktif;
- employee tidak boleh duplicate pada session;
- quota harus dicek;
- prerequisite harus terpenuhi;
- mandatory training harus sesuai target.

## Assessment

- score tidak boleh melebihi max score;
- passing score <= max score;
- assessment required harus diselesaikan sebelum completion.

## Certificate

- certificate number unik;
- issued date valid;
- expiry date >= issued date;
- certificate hanya dapat dibuat untuk participant yang memenuhi completion requirement.

---

# 33. Status Lifecycle

## Course

```text
ACTIVE
INACTIVE
```

## Session

```text
DRAFT
SCHEDULED
REGISTRATION_OPEN
FULL
ONGOING
COMPLETED
CANCELLED
```

## Participant

```text
NOMINATED
REQUESTED
APPROVED
REGISTERED
WAITLISTED
CANCELLED
COMPLETED
FAILED
```

## Training Request

```text
DRAFT
SUBMITTED
APPROVED
REJECTED
CANCELLED
```

---

# 34. API Plan

## Catalog

```http
GET    /training/categories
POST   /training/categories
GET    /training/categories/{id}
PUT    /training/categories/{id}
DELETE /training/categories/{id}

GET    /training/courses
POST   /training/courses
GET    /training/courses/{id}
PUT    /training/courses/{id}
DELETE /training/courses/{id}
```

## Planning

```http
GET    /training/plans
POST   /training/plans
GET    /training/plans/{id}
PUT    /training/plans/{id}

GET    /training/needs
POST   /training/needs
PUT    /training/needs/{id}

GET    /training/requests
POST   /training/requests
GET    /training/requests/{id}
```

## Sessions

```http
GET    /training/sessions
POST   /training/sessions
GET    /training/sessions/{id}
PUT    /training/sessions/{id}
DELETE /training/sessions/{id}

POST   /training/sessions/{id}/publish
POST   /training/sessions/{id}/cancel
```

## Participants

```http
GET    /training/sessions/{id}/participants
POST   /training/sessions/{id}/participants
DELETE /training/sessions/{id}/participants/{participantId}
```

## Attendance

```http
GET    /training/sessions/{id}/attendance
POST   /training/sessions/{id}/attendance
PUT    /training/attendance/{id}
```

## Assessment

```http
GET    /training/sessions/{id}/assessments
POST   /training/sessions/{id}/assessments
POST   /training/assessments/{id}/results
```

## Evaluation

```http
GET    /training/sessions/{id}/evaluation
POST   /training/sessions/{id}/evaluation
```

## Certificate

```http
GET    /training/certificates
GET    /training/certificates/{id}
POST   /training/participants/{id}/certificate
```

---

# 35. Backend Structure

Ikuti pola repository/service yang digunakan project existing.

```text
Training
├── Category
├── Course
├── CourseObjective
├── CourseCompetency
├── CoursePrerequisite
├── Provider
├── Trainer
├── Session
├── SessionTrainer
├── Plan
├── PlanItem
├── Need
├── Request
├── Participant
├── Attendance
├── Material
├── Assessment
├── AssessmentResult
├── Evaluation
├── EvaluationQuestion
├── EvaluationAnswer
├── EffectivenessAssessment
├── Certification
├── Certificate
├── Mandatory
├── SessionCost
└── Document
```

---

# 36. Authorization

Minimal permission:

```text
training.view
training.create
training.update
training.delete

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

Approval tetap menggunakan **Module Approval**, bukan approval engine baru di Training.

---

# 37. Frontend Plan

## Dashboard

Tampilkan:

```text
Upcoming Training
Ongoing Training
Completed Training
Total Participants
Completion Rate
Training Hours
Training Cost
Certificate Expiring
Mandatory Training Compliance
```

## Catalog

Course card:

```text
Course
Category
Duration
Delivery Type
Provider
Certification
Prerequisite
```

## Session

```text
Course
Delivery
Provider
Trainer
Location
Schedule
Quota
Participants
Status
```

## Participant

```text
Employee
Registration Status
Attendance
Pre-Test
Post-Test
Final Score
Completion
Certificate
```

## Training Detail

Tab:

```text
Overview
Participants
Attendance
Materials
Assessment
Evaluation
Certificate
Documents
Cost
```

---

# 38. Reporting

Minimal report:

### Training Participation

```text
Employee
Organization
Course
Session
Status
Attendance
Score
Completion
```

### Training Cost

```text
Course
Session
Provider
Total Cost
Participant Count
Cost per Participant
```

### Training Compliance

```text
Organization
Position
Employee
Mandatory Course
Due Date
Completion
Status
```

### Training History

```text
Employee
Course
Date
Score
Completion
Certificate
```

---

# 39. Integration dengan Modul Lain

## Employee Management

Training mengambil:

```text
employee_id
organization_id
position_id
employment_status
```

## Competency

Training menghasilkan evidence:

```text
Course
   ↓
Competency
   ↓
Training Completion
   ↓
Competency Development
```

## Performance / KPI / OKR

Training dapat menjadi response terhadap performance/competency gap.

```text
Performance
    ↓
Gap
    ↓
Training Need
```

## Career Intelligence

Training digunakan untuk:

```text
Career Path
Career Eligibility
Competency Gap
Development Recommendation
Talent
Succession
```

## Workforce Intelligence

Training digunakan untuk:

```text
Workforce Planning
Workforce Gap
Mandatory Training Compliance
Organization Training Analysis
Training Cost Analysis
```

Tidak ada module Intelligence baru.

---

# 40. Seeders

Master yang cocok dibuat seeder:

```text
Training Categories
Course Types
Delivery Types
Assessment Types
Attendance Status
Participant Status
Training Request Status
Session Status
Training Provider Type
Trainer Type
```

Contoh:

```text
categories:
- Technical
- Soft Skill
- Leadership
- Management
- Compliance
- Safety
- Certification
```

Seeder mengikuti pola seeder existing project.

---

# 41. Testing Plan

## Unit Test

Test:

- course validation;
- session validation;
- quota calculation;
- participant registration;
- duplicate participant;
- prerequisite validation;
- assessment score;
- passing score;
- completion;
- certificate eligibility;
- mandatory training;
- provider validation;
- trainer validation.

## Feature Test

Test flow:

```text
Create Course
→ Create Session
→ Register Employee
→ Attendance
→ Assessment
→ Complete
→ Evaluation
→ Certificate
```

Test kedua:

```text
External Course
→ Provider
→ External Trainer
→ Session
→ Participant
→ Completion
```

Test approval:

```text
Training Request
→ Submit
→ Approval
→ Approved
→ Enrollment
```

---

# 42. Development Priority

## P0 — Core Operational

1. Refactor `training_courses`
2. Provider
3. Trainer
4. Session enhancement
5. Participant/enrollment
6. Attendance detail
7. Material
8. Assessment
9. Completion
10. Certificate

## P1 — Planning & Governance

11. Training Request
12. Training Plan
13. Training Need
14. Mandatory Training
15. Training Cost
16. Training Documents
17. Course Objective
18. Course Competency
19. Prerequisite

## P2 — Advanced Development

20. Evaluation Form
21. Training Effectiveness
22. Certification Master
23. Advanced reporting
24. Workforce Intelligence integration
25. Career Intelligence integration

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
       ├─ Objective          └─ Request            ├─ Attendance
       ├─ Competency                               ├─ Material
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

Struktur existing sudah menjadi fondasi yang baik untuk operational training karena sudah memiliki:

```text
Category
Course
Session
Participant
Material
Evaluation
Certificate
```

Pengembangan utama adalah mengubahnya dari **simple training management** menjadi **end-to-end Training & Development Management** dengan menambahkan:

```text
Training Planning
Training Need
Training Request
Provider
Trainer
Enrollment
Detailed Attendance
Objectives
Competency Mapping
Prerequisite
Assessment
Evaluation
Effectiveness
Mandatory Training
Certification
Cost
Documents
```

Prinsip akhirnya:

```text
Training Module
    = Operational Source of Truth

Workforce Intelligence
    = Workforce-level analysis

Career Intelligence
    = Individual career/development analysis
```

Tidak diperlukan module Intelligence baru.
