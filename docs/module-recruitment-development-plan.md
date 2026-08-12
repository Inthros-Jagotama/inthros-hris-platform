# Recruitment & Onboarding — Development Plan

> **Baseline:** `015_recruitment.sql`
> **Goal:** mengembangkan Recruitment dari simple ATS menjadi **Integrated Recruitment Module** yang terhubung dengan Workforce Intelligence, Organization/Position, Module Approval, Employee, Employee Movement, Competency, Training, Performance, Career Intelligence, Career Path, dan Succession Planning.

---

# 1. Executive Summary

Struktur existing sudah mempunyai fondasi ATS:

```text
job_requisitions
candidates
job_applications
interviews
onboarding_task_templates
employee_onboardings
onboarding_task_items
```

`job_requisitions` saat ini sudah menyimpan `organization_id`, title, department, employment type, salary range, slots, status, requester, approval, dan target start date. `candidates` menyimpan profil dasar kandidat dan sumber kandidat. `job_applications` menghubungkan kandidat dengan requisition dan menyimpan status pipeline dasar. `interviews` sudah mendukung interviewer, jadwal, lokasi/link, score, dan feedback.

Namun target sistem adalah:

```text
Workforce Intelligence
        │
        ▼
Workforce Gap / Hiring Need
        │
        ▼
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
                │
                ▼
             Training
                │
                ▼
           Performance
                │
                ▼
       Career Intelligence
```

---

# 2. Business Principles

## 2.1 Recruitment bukan master Position

Karena konsep HRIS:

```text
Organization = Position
```

Recruitment tidak membuat master position/organization baru.

Recruitment menggunakan:

```text
organization_id
```

atau `position_id` jika pada implementasi Organization/Position sudah dipisahkan.

Master tetap berada pada module Organization/Job Management.

---

## 2.2 Recruitment bukan Approval Engine

Approval seluruh proses mengikuti **Module Approval**.

Minimal:

```text
Job Requisition Approval
Offer Approval
```

Jangan membuat approval engine baru di Recruitment.

Field legacy seperti:

```text
approved_by
```

boleh dipertahankan sementara untuk backward compatibility, tetapi bukan lagi source of truth approval.

---

## 2.3 Recruitment bukan Employee Movement

Untuk kandidat internal:

```text
Recruitment
    ↓
Selection
    ↓
Accepted
    ↓
Employee Movement
```

Recruitment tidak membuat employee baru untuk internal candidate.

---

## 2.4 Recruitment bukan Training

Recruitment hanya menghasilkan kebutuhan onboarding/development.

```text
Recruitment
    ↓
Employee
    ↓
Onboarding
    ↓
Training
```

Training tetap menjadi source of truth untuk training.

---

## 2.5 Career Intelligence bersifat strategis

Career Intelligence sudah dirancang sebagai strategic layer untuk career path, gap analysis, talent map, dan succession; Employee Movement tetap transactional.

Recruitment menggunakan output Career Intelligence untuk internal candidate recommendation, bukan mengeksekusi career movement.

---

# 3. Target Module Boundary

```text
RECRUITMENT
├── Workforce Hiring Need
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

---

# 4. Existing Database Baseline

## 4.1 `job_requisitions`

Existing:

```text
id
organization_id
title
department
employment_type
location
min_salary
max_salary
description
requirements
responsibilities
slots_available
slots_filled
status
requested_by
approved_by
target_start_date
closed_at
created_at
updated_at
```

Struktur tersebut sudah cukup sebagai dasar, tetapi perlu diubah agar terintegrasi dengan master Organization/Position dan Workforce Intelligence.

---

## 4.2 `candidates`

Existing:

```text
id
first_name
last_name
email
phone
address
current_company
current_title
resume_url
portfolio_url
linkedin_url
source
notes
created_at
updated_at
```

Saat ini kandidat masih berupa profil dasar.

---

## 4.3 `job_applications`

Existing:

```text
id
requisition_id
candidate_id
status
applied_at
screened_at
shortlisted_at
offered_at
accepted_at
rejected_at
withdrawn_at
rejection_reason
notes
created_at
updated_at
```

Perlu ditingkatkan menjadi pipeline yang memiliki stage history dan assessment/interview result.

---

## 4.4 `interviews`

Existing:

```text
id
application_id
interviewer_id
stage
scheduled_at
duration_minutes
location
meeting_link
status
score
feedback
completed_at
created_at
updated_at
```

Struktur sudah mendukung interview dasar, tetapi belum mendukung multi-interviewer scorecard secara terstruktur.

---

# 5. Target Database Architecture

```text
job_requisitions
      │
      ├── job_requisition_requirements
      ├── job_requisition_competencies
      └── recruitment approval
      │
      ▼
job_applications
      │
      ├── application_stage_histories
      ├── screening
      ├── assessments
      ├── interviews
      │      ├── interviewers
      │      └── scorecards
      │
      └── job_offers
```

Candidate:

```text
candidates
   ├── candidate_educations
   ├── candidate_work_experiences
   ├── candidate_skills
   ├── candidate_certifications
   ├── candidate_documents
   └── candidate_consents
```

---

# 6. Phase 0 — Preparation

## Task

- Review migration existing.
- Identifikasi model/entity existing.
- Identifikasi repository/service/controller existing.
- Identifikasi route dan permission existing.
- Identifikasi pola UUID/CHAR(36).
- Identifikasi tenant DB pattern.
- Identifikasi module approval API.
- Identifikasi Organization/Position API.
- Identifikasi Employee API.
- Identifikasi Competency API.
- Identifikasi Workforce Intelligence API.
- Identifikasi Career Intelligence API.
- Identifikasi Employee Movement API.
- Identifikasi Training API.
- Identifikasi Performance API.

## Rule

Semua PK baru menggunakan:

```text
UUID / CHAR(36)
```

mengikuti pola database existing.

---

# 7. Phase 1 — Refactor Job Requisition

## 7.1 Objective

Mengubah requisition dari lowongan sederhana menjadi hiring request yang terintegrasi.

## Enhancement fields

Tambahkan:

```text
requisition_number
workforce_gap_id       nullable
workforce_plan_id      nullable
reason_type
priority
position_id            nullable jika position terpisah
approval_status
opened_at
closed_at
```

### `reason_type`

```text
NEW_POSITION
REPLACEMENT
BACKFILL
WORKFORCE_GAP
EXPANSION
INTERNAL_MOVEMENT
```

### `priority`

```text
LOW
MEDIUM
HIGH
URGENT
```

---

# 8. Organization / Position Integration

## Rule

Jangan menyimpan:

```text
department = "IT"
title = "Software Engineer"
```

sebagai master bebas jika data tersebut sudah tersedia pada Organization/Position.

Gunakan:

```text
organization_id
position_id
```

sesuai struktur master existing.

Recruitment membaca:

```text
Position Title
Organization
Employment Type
Hierarchy
Competency Requirement
Career Path
```

dari master masing-masing.

---

# 9. Phase 2 — Workforce Intelligence Integration

## Flow

```text
Workforce Intelligence
        ↓
Required Workforce
        ↓
Current Workforce
        ↓
Workforce Gap
        ↓
Hiring Recommendation
        ↓
Job Requisition
```

Contoh:

```text
Required: 6
Current: 4
Gap: 2
```

Recruitment:

```text
slots_available = 2
reason_type = WORKFORCE_GAP
```

## Data reference

```text
workforce_gap_id
workforce_plan_id
```

nullable agar requisition manual tetap dapat dibuat.

---

# 10. Phase 3 — Module Approval Integration

## Job Requisition

```text
DRAFT
  ↓
SUBMITTED
  ↓
APPROVAL
  ↓
APPROVED / REJECTED
  ↓
OPEN
```

Approval diproses oleh Module Approval.

## Offer

```text
OFFER_DRAFT
  ↓
SUBMITTED
  ↓
APPROVAL
  ↓
APPROVED
  ↓
SENT
```

Recruitment hanya membaca status approval dan meneruskan workflow.

---

# 11. Phase 4 — Job Requisition Requirements

Tambahkan:

```text
job_requisition_requirements
```

Contoh:

```text
education
experience
age
language
availability
other
```

Fields:

```text
id
requisition_id
requirement_type
name
description
minimum_value
maximum_value
is_required
sort_order
created_at
updated_at
```

Semua ID UUID.

---

# 12. Phase 5 — Job Requisition Competency

Tambahkan:

```text
job_requisition_competencies
```

Fields:

```text
id
requisition_id
competency_id
required_level
is_required
weight
created_at
updated_at
```

Contoh:

```text
PHP           Level 4   Required
Laravel       Level 4   Required
PostgreSQL    Level 3   Required
Leadership    Level 2   Optional
```

Ini menjadi basis Candidate Matching.

---

# 13. Phase 6 — Candidate Enhancement

Existing candidate dipertahankan.

Tambahkan:

```text
candidate_number
status
candidate_type
source_id
consent_status
```

### Candidate Type

```text
EXTERNAL
INTERNAL
```

### Source Type

```text
DIRECT
CAREER_SITE
REFERRAL
LINKEDIN
AGENCY
INTERNAL_MOBILITY
CAREER_POOL
OTHER
```

---

# 14. Candidate Education

Buat:

```text
candidate_educations
```

Fields:

```text
id
candidate_id
institution
education_level
field_of_study
start_date
end_date
grade
description
created_at
updated_at
```

---

# 15. Candidate Work Experience

Buat:

```text
candidate_work_experiences
```

Fields:

```text
id
candidate_id
company
position
start_date
end_date
is_current
description
created_at
updated_at
```

---

# 16. Candidate Skills

Buat:

```text
candidate_skills
```

Fields:

```text
id
candidate_id
skill_name
skill_level
years_experience
created_at
updated_at
```

Jika Skill Master sudah tersedia, gunakan:

```text
skill_id
```

jangan membuat master duplicate.

---

# 17. Candidate Certification

Buat:

```text
candidate_certifications
```

Fields:

```text
id
candidate_id
certification_id nullable
name
issuer
certificate_number
issued_date
expired_date
document_id nullable
created_at
updated_at
```

---

# 18. Candidate Documents

Buat:

```text
candidate_documents
```

Jenis:

```text
RESUME
COVER_LETTER
CERTIFICATE
PORTFOLIO
IDENTITY
OTHER
```

Jangan menyimpan file binary langsung pada database.

Gunakan file storage/reference mengikuti pola existing application.

---

# 19. Candidate Internal

Internal candidate harus mendukung:

```text
candidate_type = INTERNAL
employee_id = existing employee
```

Flow:

```text
Employee
   ↓
Career Intelligence
   ↓
Eligible Candidate
   ↓
Recruitment Application
```

Tidak membuat employee baru.

---

# 20. Phase 7 — Recruitment Pipeline

Existing `job_applications.status` diubah menjadi current stage.

Buat:

```text
recruitment_stages
```

Default:

```text
APPLIED
SCREENING
SHORTLISTED
ASSESSMENT
INTERVIEW
FINAL_REVIEW
OFFER
HIRED
REJECTED
WITHDRAWN
```

Buat:

```text
job_application_stage_histories
```

Fields:

```text
id
application_id
from_stage_id nullable
to_stage_id
changed_by
changed_at
notes
```

Tujuan:

```text
Audit Trail
Time to Stage
Time to Hire
Pipeline Analytics
```

---

# 21. Phase 8 — Screening

Buat:

```text
application_screenings
```

Fields:

```text
id
application_id
screened_by
screened_at
score
result
notes
```

Result:

```text
PASS
FAIL
HOLD
```

---

# 22. Phase 9 — Assessment

Buat:

```text
recruitment_assessments
assessment_participants
assessment_results
```

Jenis:

```text
TECHNICAL
PSYCHOLOGICAL
COGNITIVE
PERSONALITY
CASE_STUDY
CODING
LANGUAGE
OTHER
```

Assessment dapat menghasilkan:

```text
score
result
recommendation
```

---

# 23. Phase 10 — Interview Enhancement

Existing interview dipertahankan.

Tambahkan konsep:

```text
interviewers
interview_scorecards
interview_scorecard_items
```

Satu interview dapat memiliki banyak interviewer.

```text
Interview
 ├── HR
 ├── User
 └── Manager
```

Setiap interviewer memiliki scorecard sendiri.

---

# 24. Interview Scorecard

Contoh:

```text
Technical Skill       30%
Problem Solving       20%
Communication         20%
Leadership            15%
Culture Fit           15%
```

Score:

```text
1 - 5
```

atau:

```text
0 - 100
```

Normalisasi dilakukan di service layer.

---

# 25. Phase 11 — Candidate Matching

Candidate Matching menggunakan:

```text
Job Requirement
+
Competency
+
Education
+
Experience
+
Skill
+
Certification
+
Assessment
+
Interview
```

Output:

```text
candidate_match_score
```

Contoh:

```text
Budi       92%
Andi       87%
Dedi       76%
```

Jangan menggunakan match score sebagai keputusan otomatis.

Recruiter tetap dapat override dengan alasan yang tercatat.

---

# 26. Phase 12 — Offer Management

Buat:

```text
job_offers
```

Fields minimal:

```text
id
application_id
offer_number
employment_type
salary
allowances
benefits
start_date
expiry_date
status
sent_at
accepted_at
rejected_at
created_at
updated_at
```

Status:

```text
DRAFT
PENDING_APPROVAL
APPROVED
SENT
ACCEPTED
REJECTED
EXPIRED
WITHDRAWN
```

---

# 27. Offer Approval

Flow:

```text
Recruiter
   ↓
Offer Draft
   ↓
Module Approval
   ↓
Approved
   ↓
Offer Sent
   ↓
Candidate
   ↓
Accepted
```

---

# 28. Phase 13 — Recruitment → Employee

External candidate:

```text
Candidate
   ↓
Application
   ↓
Offer Accepted
   ↓
Create Employee
```

Simpan reference:

```text
employee.recruited_from_application_id
```

atau equivalent reference sesuai Employee module.

Tujuan:

```text
Employee
   ↓
Application
   ↓
Requisition
   ↓
Position
```

dapat ditelusuri kembali.

---

# 29. Phase 14 — Internal Recruitment → Employee Movement

Internal:

```text
Existing Employee
       ↓
Internal Application
       ↓
Selection
       ↓
Offer / Decision
       ↓
Employee Movement
       ↓
New Organization / Position
```

Recruitment tidak membuat employee baru.

Employee Movement tetap menjadi transactional execution layer.

---

# 30. Phase 15 — Recruitment → Onboarding

Saat offer accepted:

```text
job_application
        ↓
employee
        ↓
employee_onboarding
```

Existing `employee_onboardings` sudah memiliki:

```text
employee_id
application_id
start_date
status
buddy_id
```

sehingga fondasinya sudah mendukung hubungan Recruitment → Onboarding.

---

# 31. Onboarding Template Enhancement

Existing template memiliki:

```text
name
description
category
day_offset
assigned_role
is_mandatory
```



Enhance agar template dapat dibedakan berdasarkan:

```text
organization_id nullable
position_id nullable
employment_type nullable
```

Contoh:

```text
Software Engineer
    ↓
Laptop
Repository Access
Security Training
Technical Orientation
Team Introduction
```

---

# 32. Phase 16 — Recruitment → Training

Setelah employee onboarding:

```text
Employee
   ↓
Position
   ↓
Required Training
   ↓
Training Module
```

Training tetap operational source of truth.

Recruitment/Onboarding hanya membuat reference atau requirement.

---

# 33. Phase 17 — Recruitment → Competency

Position:

```text
Software Engineer
```

Requirement:

```text
PHP Level 4
Laravel Level 4
PostgreSQL Level 3
```

Recruitment menggunakan competency requirement dari Position/Job Management atau requisition-specific override.

Candidate assessment dapat dibandingkan:

```text
Candidate Competency
vs
Position Competency
```

---

# 34. Phase 18 — Recruitment → Career Intelligence

Career Intelligence dapat memberikan internal candidate:

```text
Position Vacancy
      ↓
Career Path
      ↓
Gap Analysis
      ↓
Eligible Employees
```

Contoh:

```text
IT Supervisor

Budi    94%
Andi    87%
Dedi    76%
```

Recruiter dapat:

```text
Invite to Apply
```

atau:

```text
Create Internal Application
```

Recruitment tidak menentukan career eligibility sendiri.

---

# 35. Phase 19 — Recruitment → Succession Planning

Jika position adalah critical/key position:

```text
Succession Plan
      ↓
Successor Candidates
      ↓
Recruitment
```

Jika successor internal tidak tersedia:

```text
Succession Gap
      ↓
External Recruitment
```

Ini memungkinkan recruitment menjadi fallback untuk succession planning.

---

# 36. Phase 20 — Recruitment → Performance

Setelah employee hired:

```text
Recruitment
   ↓
Employee
   ↓
Performance
```

Career/Workforce Intelligence dapat mengukur:

```text
Quality of Hire
```

berdasarkan:

```text
Hiring Source
Recruitment Score
Probation Result
Performance Result
Retention
```

Recruitment tidak mengubah Performance score.

---

# 37. Recruitment Analytics

Minimal metrics:

```text
Open Requisitions
Applications
Candidates
Shortlisted
Interviews
Offers
Hires
Rejected
Withdrawn
```

Advanced:

```text
Time to Hire
Time to Fill
Time to Stage
Offer Acceptance Rate
Application Conversion Rate
Source Conversion
Candidate Match Score
Quality of Hire
```

---

# 38. Workforce Intelligence Output

Recruitment menyediakan data:

```text
Open Positions
Recruitment Pipeline
Expected Hires
Accepted Offers
Filled Positions
```

Workforce Intelligence dapat menghitung:

```text
Required Workforce
-
Current Workforce
-
Expected Hires
=
Remaining Workforce Gap
```

Contoh:

```text
Required       100
Current         90
Accepted Offer   2

Remaining Gap = 8
```

---

# 39. Career Intelligence Output

Recruitment menerima:

```text
Career Eligibility
Career Path
Competency Gap
Succession Candidates
Talent Pool
```

Recruitment dapat menampilkan:

```text
Internal Candidates
External Candidates
```

dengan source yang jelas.

---

# 40. Permissions

Rekomendasi:

```text
recruitment.view
recruitment.requisition.view
recruitment.requisition.manage
recruitment.requisition.submit
recruitment.candidate.view
recruitment.candidate.manage
recruitment.application.view
recruitment.application.manage
recruitment.screening.manage
recruitment.assessment.manage
recruitment.interview.view
recruitment.interview.manage
recruitment.offer.view
recruitment.offer.manage
recruitment.offer.submit
recruitment.analytics.view
recruitment.onboarding.view
recruitment.onboarding.manage
```

Permission harus mengikuti pola permission module existing.

---

# 41. Authorization

## Requisition

User dapat:

```text
View
Create
Edit
Submit
```

sesuai organization scope.

## Recruiter

```text
Candidate
Application
Screening
Interview
Offer
```

## Hiring Manager

```text
View assigned requisition
Review candidate
Interview
Give recommendation
```

## HR

```text
Full Recruitment
Onboarding
Offer
```

## Employee

Untuk internal recruitment:

```text
View own application
Withdraw application
```

sesuai business rule.

---

# 42. API Plan

## Requisition

```http
GET    /recruitment/requisitions
POST   /recruitment/requisitions
GET    /recruitment/requisitions/{id}
PUT    /recruitment/requisitions/{id}
POST   /recruitment/requisitions/{id}/submit
POST   /recruitment/requisitions/{id}/close
```

## Candidate

```http
GET    /recruitment/candidates
POST   /recruitment/candidates
GET    /recruitment/candidates/{id}
PUT    /recruitment/candidates/{id}
```

## Application

```http
GET    /recruitment/applications
POST   /recruitment/applications
GET    /recruitment/applications/{id}
POST   /recruitment/applications/{id}/stage
POST   /recruitment/applications/{id}/screen
```

## Interview

```http
GET    /recruitment/interviews
POST   /recruitment/interviews
PUT    /recruitment/interviews/{id}
POST   /recruitment/interviews/{id}/complete
```

## Assessment

```http
GET    /recruitment/assessments
POST   /recruitment/assessments
POST   /recruitment/applications/{id}/assessments
```

## Offer

```http
GET    /recruitment/offers
POST   /recruitment/offers
POST   /recruitment/offers/{id}/submit
POST   /recruitment/offers/{id}/send
POST   /recruitment/offers/{id}/accept
POST   /recruitment/offers/{id}/reject
```

---

# 43. Frontend Plan

## Recruitment Dashboard

Widgets:

```text
Open Requisitions
Candidates
Applications
Interviews
Offers
Hires
Time to Hire
```

## Requisition

```text
List
Create
Detail
Edit
Approval Status
Pipeline
```

## Candidate

```text
Candidate List
Candidate Profile
Resume
Experience
Education
Skills
Certification
Applications
```

## Application

Pipeline:

```text
Applied
   ↓
Screening
   ↓
Assessment
   ↓
Interview
   ↓
Final Review
   ↓
Offer
   ↓
Hired
```

## Interview

```text
Calendar
Schedule
Interviewer
Scorecard
Feedback
Result
```

## Offer

```text
Offer Draft
Approval
Offer Preview
Send
Acceptance
```

---

# 44. Candidate Profile

Candidate detail ideal:

```text
Profile
├── Personal
├── Contact
├── Resume
├── Education
├── Experience
├── Skills
├── Certifications
├── Applications
├── Assessments
├── Interviews
└── Offers
```

Untuk internal:

```text
Current Employee
Current Organization
Current Position
Performance
Competency
Career Eligibility
```

akses data mengikuti permission.

---

# 45. Recruitment Pipeline UI

Contoh:

```text
┌──────────┬───────────┬────────────┬────────────┬─────────┐
│ Applied  │ Screening │ Assessment │ Interview  │ Offer   │
├──────────┼───────────┼────────────┼────────────┼─────────┤
│ Budi     │ Andi      │ Dedi       │ Sari       │ Joko    │
│ Asep     │ Rina      │            │ Toni       │         │
└──────────┴───────────┴────────────┴────────────┴─────────┘
```

Drag/drop hanya mengubah stage melalui backend transition service.

---

# 46. Seeder Plan

Seeder dapat dibuat untuk:

```text
Recruitment Stage
Recruitment Source
Requirement Type
Assessment Type
Interview Stage
Offer Status
Requisition Reason
Priority
```

Contoh:

```text
RecruitmentStageSeeder
RecruitmentSourceSeeder
RecruitmentAssessmentTypeSeeder
RecruitmentRequirementTypeSeeder
```

Seeder mengikuti pola existing project.

Jangan membuat seeder untuk transactional data:

```text
candidate
application
interview
offer
```

kecuali untuk development/demo seeder.

---

# 47. Migration Plan

## Migration 1

Enhance:

```text
job_requisitions
```

## Migration 2

Enhance:

```text
candidates
```

## Migration 3

Enhance:

```text
job_applications
```

## Migration 4

Create:

```text
job_requisition_requirements
job_requisition_competencies
```

## Migration 5

Create:

```text
recruitment_stages
job_application_stage_histories
```

## Migration 6

Create:

```text
candidate_educations
candidate_work_experiences
candidate_skills
candidate_certifications
candidate_documents
```

## Migration 7

Create:

```text
application_screenings
```

## Migration 8

Create:

```text
recruitment_assessments
assessment_participants
assessment_results
```

## Migration 9

Enhance interview:

```text
interviewers
interview_scorecards
interview_scorecard_items
```

## Migration 10

Create:

```text
job_offers
```

## Migration 11

Onboarding enhancement:

```text
organization/position-specific template
```

---

# 48. Data Migration

Existing:

```text
job_requisitions.title
job_requisitions.department
job_requisitions.employment_type
```

harus dipetakan ke master yang tersedia.

Candidate existing:

```text
source
```

tetap dipertahankan dan dimigrasikan ke source master jika source master dibuat.

Application existing:

```text
status
```

dipetakan ke:

```text
recruitment_stages
```

History lama tidak boleh dibuat secara fiktif jika tidak tersedia.

---

# 49. Backend Architecture

Rekomendasi:

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
│   ├── SubmitRequisition
│   ├── MoveApplicationStage
│   ├── ScheduleInterview
│   ├── CompleteInterview
│   ├── CreateOffer
│   └── AcceptOffer
│
└── Integration
    ├── Organization
    ├── Workforce Intelligence
    ├── Approval
    ├── Employee
    ├── Movement
    ├── Competency
    ├── Training
    └── Career Intelligence
```

---

# 50. Business Rules

## Requisition

1. Requisition harus memiliki organization/position.
2. Requisition tidak dapat dibuka sebelum approval selesai.
3. Slots tidak boleh negatif.
4. `slots_filled <= slots_available`.
5. Requisition dapat ditutup jika seluruh slot terpenuhi atau dibatalkan.

## Application

1. Candidate dapat memiliki banyak application.
2. Candidate tidak boleh duplicate application aktif untuk requisition yang sama.
3. Stage transition harus tervalidasi.
4. Setiap perubahan stage harus memiliki history.

## Interview

1. Interview harus memiliki application.
2. Interview harus memiliki interviewer.
3. Interview yang completed dapat memiliki score.
4. Multi-interviewer harus didukung.

## Offer

1. Hanya candidate yang eligible yang dapat menerima offer.
2. Offer harus melalui approval.
3. Offer accepted menghasilkan employee atau Employee Movement sesuai candidate type.
4. Offer expired tidak dapat diterima.

## Internal Candidate

```text
candidate_type = INTERNAL
```

harus memiliki:

```text
employee_id
```

dan tidak membuat employee baru.

---

# 51. Testing Plan

## Unit Test

### Requisition

- create
- update
- submit
- approval status
- close
- slot validation

### Candidate

- create
- duplicate email
- update
- document

### Application

- create
- duplicate active application
- stage transition
- rejection
- withdrawal

### Interview

- schedule
- reschedule
- complete
- scorecard
- multiple interviewers

### Offer

- create
- approval
- send
- accept
- reject
- expire

---

# 52. Integration Test

Test:

```text
Workforce Gap
    ↓
Requisition
```

```text
Requisition
    ↓
Module Approval
```

```text
Application
    ↓
Interview
```

```text
Application
    ↓
Offer
    ↓
Employee
```

```text
Internal Application
    ↓
Employee Movement
```

```text
Offer Accepted
    ↓
Onboarding
```

```text
Employee
    ↓
Training
```

```text
Career Intelligence
    ↓
Internal Candidate
    ↓
Recruitment
```

---

# 53. E2E Test

## External Hiring

```text
Create Workforce Gap
↓
Create Requisition
↓
Submit
↓
Approve
↓
Publish
↓
Candidate Apply
↓
Screening
↓
Assessment
↓
Interview
↓
Final Selection
↓
Offer
↓
Approval
↓
Offer Accepted
↓
Employee Created
↓
Onboarding
↓
Training
```

## Internal Hiring

```text
Position Vacancy
↓
Career Intelligence
↓
Eligible Employee
↓
Internal Application
↓
Selection
↓
Offer / Decision
↓
Employee Movement
↓
New Position
↓
Onboarding / Development
```

---

# 54. Recruitment Analytics

## Funnel

```text
Applications
    ↓
Screened
    ↓
Shortlisted
    ↓
Assessment
    ↓
Interview
    ↓
Offer
    ↓
Accepted
    ↓
Hired
```

Metrics:

```text
Conversion Rate
Rejection Rate
Withdrawal Rate
Offer Acceptance Rate
```

## Time Metrics

```text
Time to Screen
Time to Interview
Time to Offer
Time to Hire
Time to Fill
```

## Source Analytics

```text
Source
Applications
Shortlisted
Interviews
Offers
Hires
Quality of Hire
```

---

# 55. Quality of Hire

Setelah employee masuk:

```text
Recruitment
      ↓
Employee
      ↓
Probation
      ↓
Performance
```

Career/Workforce Intelligence dapat menghitung:

```text
Quality of Hire
```

berdasarkan:

```text
Recruitment Match Score
Interview Score
Assessment Score
Probation Result
Performance
Retention
```

Recruitment menyimpan source data; intelligence menghitung analytical score.

---

# 56. Security & Privacy

Candidate data mengandung data pribadi.

Wajib:

- role-based access;
- organization scope;
- audit log;
- document access control;
- consent tracking;
- secure file storage;
- masking data pada role tertentu;
- tidak menampilkan candidate data ke user tanpa permission.

Data sensitif jangan disimpan jika tidak diperlukan untuk proses recruitment.

---

# 57. Audit Trail

Audit minimal untuk:

```text
Requisition
Application
Stage Transition
Screening
Interview
Assessment
Offer
Onboarding
```

Contoh:

```text
Application
Budi
SCREENING → INTERVIEW
Changed by: Recruiter
Date: ...
Reason: Passed screening
```

---

# 58. Notification

Integrasi notification:

```text
Requisition Submitted
Requisition Approved
Requisition Rejected
Interview Scheduled
Interview Rescheduled
Assessment Assigned
Offer Approval Required
Offer Approved
Offer Sent
Offer Accepted
Onboarding Started
```

Channel mengikuti notification infrastructure existing.

---

# 59. Development Priority

## P0 — Core Integrated Recruitment

1. Requisition refactor
2. Organization/Position integration
3. Module Approval integration
4. Candidate enhancement
5. Application pipeline
6. Stage history
7. Screening
8. Interview enhancement
9. Offer
10. Recruitment → Employee
11. Onboarding integration

## P1 — Intelligent Recruitment

12. Requisition competency
13. Candidate competency/skills
14. Assessment
15. Candidate matching
16. Workforce Intelligence integration
17. Career Intelligence integration
18. Internal recruitment
19. Employee Movement integration
20. Training integration

## P2 — Advanced Recruitment

21. Candidate Pool
22. Candidate Tags
23. Talent Pool
24. Referral management
25. Candidate ranking
26. Quality of Hire
27. Advanced analytics
28. Recruitment forecasting

---

# 60. Recommended Implementation Order

```text
STEP 1
Database & Migration
        ↓
STEP 2
Models / Repository
        ↓
STEP 3
Requisition
        ↓
STEP 4
Approval
        ↓
STEP 5
Candidate
        ↓
STEP 6
Application Pipeline
        ↓
STEP 7
Screening
        ↓
STEP 8
Interview
        ↓
STEP 9
Assessment
        ↓
STEP 10
Offer
        ↓
STEP 11
Employee Integration
        ↓
STEP 12
Onboarding
        ↓
STEP 13
Competency Integration
        ↓
STEP 14
Workforce Intelligence
        ↓
STEP 15
Career Intelligence
        ↓
STEP 16
Employee Movement
        ↓
STEP 17
Training
        ↓
STEP 18
Analytics
        ↓
STEP 19
Testing
```

---

# 61. Final Architecture

```text
                 WORKFORCE INTELLIGENCE
                          │
                          ▼
                   Workforce Gap
                          │
                          ▼
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
                          │
                          ▼
                       TRAINING
                          │
                          ▼
                     PERFORMANCE
                          │
                          ▼
                 CAREER INTELLIGENCE
                          │
              ┌───────────┼───────────┐
              ▼           ▼           ▼
         Career Path   Talent Map  Succession
```

---

# 62. Definition of Done

Recruitment enhancement dianggap selesai apabila:

- [ ] Requisition menggunakan Organization/Position master.
- [ ] Requisition dapat berasal dari Workforce Gap.
- [ ] Requisition menggunakan Module Approval.
- [ ] Candidate memiliki profile terstruktur.
- [ ] Candidate mendukung internal/external.
- [ ] Application memiliki pipeline.
- [ ] Stage transition memiliki history.
- [ ] Screening tersedia.
- [ ] Assessment tersedia.
- [ ] Interview mendukung multi-interviewer.
- [ ] Interview menggunakan scorecard.
- [ ] Offer menjadi entity sendiri.
- [ ] Offer menggunakan Module Approval.
- [ ] External candidate dapat menjadi Employee.
- [ ] Internal candidate menggunakan Employee Movement.
- [ ] Offer accepted dapat membuat onboarding.
- [ ] Onboarding dapat terhubung dengan Training.
- [ ] Requisition dapat menggunakan competency requirement.
- [ ] Candidate dapat dinilai terhadap competency requirement.
- [ ] Career Intelligence dapat memberikan internal candidate.
- [ ] Workforce Intelligence dapat memberikan hiring need.
- [ ] Recruitment menyediakan data kembali ke Workforce Intelligence.
- [ ] Recruitment menyediakan data kembali ke Career Intelligence.
- [ ] Permission lengkap.
- [ ] Audit trail tersedia.
- [ ] Notification tersedia.
- [ ] Unit test selesai.
- [ ] Integration test selesai.
- [ ] E2E external hiring selesai.
- [ ] E2E internal hiring selesai.
- [ ] Migration dan backward compatibility diverifikasi.

---

# 63. Kesimpulan

Target akhir bukan sekadar:

```text
ATS
```

tetapi:

```text
Integrated Recruitment
```

yang menjadi penghubung antara:

```text
Workforce Planning
        ↓
Hiring
        ↓
Employee
        ↓
Onboarding
        ↓
Training
        ↓
Performance
        ↓
Career
```

Dengan demikian Recruitment menjadi bagian dari **Employee Lifecycle** sekaligus menjadi execution layer untuk memenuhi kebutuhan workforce.

Prinsip pembagian responsibility:

```text
Workforce Intelligence
→ menentukan kebutuhan workforce

Recruitment
→ mencari dan memilih kandidat

Module Approval
→ mengelola approval

Employee
→ menjadi master employee

Employee Movement
→ mengeksekusi perpindahan internal

Onboarding
→ mempersiapkan employee baru

Training
→ mengelola development/training

Performance
→ menilai performance

Career Intelligence
→ menganalisis career, talent, gap, dan succession
```
