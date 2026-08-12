# Workforce Intelligence — Training Enhancement Plan

## 1. Tujuan

Enhancement ini menambahkan kemampuan **Training & Development** ke dalam **Workforce Intelligence** tanpa membuat module Intelligence baru.

Fokus Workforce Intelligence adalah menjawab:

> **Bagaimana kebutuhan, kesiapan, gap, dan investasi training pada workforce organisasi?**

Training tetap menjadi **operational source of truth**. Workforce Intelligence hanya membaca dan menganalisis data training.

---

# 2. Scope

Enhancement mencakup:

1. Training Demand Analysis
2. Workforce Training Gap
3. Organization Training Readiness
4. Mandatory Training Compliance
5. Training Participation Analytics
6. Training Hours Analytics
7. Training Cost Analytics
8. Training Coverage
9. Training Effectiveness pada level organisasi
10. Competency-related Training Analysis
11. Workforce Development Forecast
12. Dashboard dan reporting
13. Data integration dari Training Module

Tidak termasuk:

- CRUD course
- CRUD session
- enrollment operasional
- attendance operasional
- certificate issuance
- training request workflow

Semua transaksi tersebut tetap berada di **Training & Development**.

---

# 3. Prinsip Arsitektur

```text
Training & Development
        │
        │ operational data
        ▼
Workforce Intelligence
        │
        ├── Training Demand
        ├── Training Gap
        ├── Training Coverage
        ├── Training Compliance
        ├── Training Cost
        ├── Training Hours
        └── Workforce Readiness
```

Workforce Intelligence tidak mengubah data transaksi Training.

---

# 4. Data Source

Workforce Intelligence menggunakan:

```text
Employee
Organization
Position
Competency
Performance
Training Course
Training Course Competency
Training Session
Training Participant
Training Attendance
Training Assessment Result
Training Certification
Training Effectiveness
Training Cost
Mandatory Training
```

Data utama dari Training:

```text
training_courses
training_course_competencies
training_sessions
training_participants
training_attendances
training_assessment_results
training_certificates
training_certifications
training_effectiveness_assessments
training_mandatories
training_session_costs
```

---

# 5. Feature 1 — Workforce Training Demand

Tujuan:

Mengetahui kebutuhan training berdasarkan kebutuhan workforce.

Contoh:

```text
Organization: IT

Required Employees       50
Employees Available     45
Workforce Gap             5

Required Competency:
Leadership Level 4

Employees below level:
12

Training Demand:
Leadership Development
Target: 12 employees
```

Dashboard:

```text
Training Demand
├── Total Demand
├── High Priority
├── Medium Priority
├── Low Priority
└── Demand by Organization
```

---

# 6. Feature 2 — Workforce Training Gap

Mengukur gap training pada level organisasi.

Contoh:

```text
Organization: IT

Required Training:
Leadership Development       30
Project Management           25
Security                     40

Completed:
Leadership Development       20
Project Management           18
Security                     35
```

Hasil:

```text
Leadership Gap = 10
Project Management Gap = 7
Security Gap = 5
```

---

# 7. Feature 3 — Training Coverage

Mengukur persentase employee yang sudah mendapatkan training yang dibutuhkan.

Formula:

```text
Training Coverage =
Employees with Required Training
/
Employees Requiring Training
× 100
```

Contoh:

```text
Required Employees = 100
Completed          = 80

Coverage = 80%
```

Breakdown:

```text
Organization
Position
Course
Competency
Period
```

---

# 8. Feature 4 — Mandatory Training Compliance

Menggunakan data `training_mandatories`.

Dashboard:

```text
Mandatory Training Compliance

Required       1,000
Completed        870
Pending           80
Overdue           50

Compliance = 87%
```

Breakdown:

```text
Organization
Position
Employee
Course
Due Date
Status
```

Status:

```text
COMPLETED
PENDING
OVERDUE
NOT_REQUIRED
```

---

# 9. Feature 5 — Training Participation

Analisis jumlah employee yang mengikuti training.

Metrics:

```text
Total Participants
Unique Employees
Completed Training
Failed Training
Cancelled
Waitlisted
```

Breakdown:

```text
Organization
Position
Course
Category
Period
```

---

# 10. Feature 6 — Training Hours

Mengukur learning hours workforce.

Metrics:

```text
Total Training Hours
Average Training Hours / Employee
Training Hours / Organization
Training Hours / Position
Training Hours / Course
```

Contoh:

```text
IT
Employees              100
Training Hours        1,200
Average / Employee       12
```

---

# 11. Feature 7 — Training Cost

Menggunakan `training_session_costs`.

Metrics:

```text
Total Training Cost
Cost per Employee
Cost per Participant
Cost per Course
Cost per Organization
Cost per Training Hour
```

Contoh:

```text
Training Cost          500M
Participants           250
Cost / Participant       2M
```

---

# 12. Feature 8 — In-House vs External Analysis

Karena Training mendukung:

```text
IN_HOUSE
EXTERNAL
```

Workforce Intelligence dapat membandingkan:

```text
In-House
External
```

Metrics:

```text
Participant Count
Training Hours
Cost
Completion Rate
Average Score
Effectiveness
```

Contoh:

| Metric | In-House | External |
|---|---:|---:|
| Participants | 300 | 150 |
| Hours | 2,000 | 1,000 |
| Cost | 200M | 450M |
| Completion | 96% | 91% |

---

# 13. Feature 9 — Competency Training Coverage

Training yang memiliki:

```text
training_course_competencies
```

dapat dianalisis terhadap competency workforce.

Contoh:

```text
Competency:
Leadership

Required Level:
4

Employees:
100

Employees Level < 4:
40

Training completed:
25

Training coverage:
62.5%
```

Dashboard:

```text
Competency
Required
Gap
Training Required
Training Completed
Coverage
```

---

# 14. Feature 10 — Workforce Training Readiness

Mengukur kesiapan workforce berdasarkan training.

Contoh:

```text
Organization: IT

Required Training Completion = 90%
Actual Completion             = 82%

Readiness:
82%
```

Jika digunakan bersama competency:

```text
Training
+
Competency
+
Certification
+
Performance
```

maka Workforce Intelligence dapat menghasilkan:

```text
Workforce Readiness
```

Catatan:

Training bukan satu-satunya faktor readiness.

---

# 15. Feature 11 — Training Effectiveness

Menggunakan data:

```text
training_effectiveness_assessments
```

pada agregasi organisasi.

Metrics:

```text
Average Before Score
Average After Score
Average Improvement
Effectiveness Rate
```

Contoh:

```text
Before Average = 70
After Average  = 84

Improvement = +14
```

Breakdown:

```text
Organization
Position
Course
Competency
Period
```

---

# 16. Feature 12 — Training Forecast

Gunakan historical training demand untuk membantu workforce planning.

Contoh:

```text
2025:
Leadership Training Demand = 80

2026:
Leadership Training Demand = 120

Forecast:
2027 ≈ 150
```

Forecast harus menjadi analytical feature dan tidak mengubah Training Plan secara otomatis.

---

# 17. Feature 13 — Workforce Training Risk

Deteksi risiko:

```text
High Training Gap
Low Training Coverage
Overdue Mandatory Training
Critical Certification Expiry
Low Training Effectiveness
Low Training Participation
```

Contoh:

```text
HIGH RISK

Organization: Production

Mandatory Training Compliance = 72%
Critical Training Gap           = 18%
Overdue Employees               = 25
```

---

# 18. Dashboard

## Workforce Training Overview

```text
Total Training
Participants
Training Hours
Training Cost
Coverage
Compliance
Effectiveness
```

## Organization Comparison

```text
Organization
Employees
Training Hours
Coverage
Compliance
Cost
Effectiveness
```

## Competency Training Gap

```text
Competency
Required Employees
Gap
Training Required
Training Completed
Coverage
```

## Training Investment

```text
In-House Cost
External Cost
Total Cost
Cost / Employee
Cost / Participant
```

---

# 19. Recommended Data Model

Tidak perlu menduplikasi tabel operational Training.

Jika Workforce Intelligence menggunakan analytical/read model, buat tabel projection:

```text
workforce_training_summary
workforce_training_gaps
workforce_training_compliance
workforce_training_costs
```

Namun jika arsitektur saat ini menggunakan query langsung/reporting view, tabel tersebut tidak wajib.

### Recommendation

Gunakan:

```text
Operational DB
       ↓
Query / View / Read Model
       ↓
Workforce Intelligence
```

Hindari duplicate transaction data.

---

# 20. API

```http
GET /workforce-intelligence/training/overview
GET /workforce-intelligence/training/demand
GET /workforce-intelligence/training/gaps
GET /workforce-intelligence/training/coverage
GET /workforce-intelligence/training/compliance
GET /workforce-intelligence/training/participation
GET /workforce-intelligence/training/hours
GET /workforce-intelligence/training/cost
GET /workforce-intelligence/training/effectiveness
GET /workforce-intelligence/training/competencies
GET /workforce-intelligence/training/readiness
GET /workforce-intelligence/training/forecast
GET /workforce-intelligence/training/risks
```

Filter:

```text
organization_id
position_id
course_id
category_id
competency_id
period
delivery_type
```

---

# 21. Integration Flow

```text
Employee
Organization
Position
Competency
Performance
       │
       ▼
Workforce Intelligence
       │
       ▼
Workforce Gap
       │
       ▼
Training Demand
       │
       ▼
Training Module
       │
       ▼
Training Completion
       │
       ▼
Workforce Intelligence
       │
       ├── Coverage
       ├── Compliance
       ├── Cost
       ├── Effectiveness
       └── Readiness
```

---

# 22. Development Priority

## P0

1. Training Overview
2. Training Participation
3. Training Hours
4. Training Coverage
5. Mandatory Compliance
6. Competency Training Gap

## P1

7. Training Cost
8. In-House vs External
9. Training Effectiveness
10. Workforce Readiness
11. Organization Comparison

## P2

12. Training Risk
13. Training Demand Forecast
14. Workforce Training Forecast
15. Advanced analytical scoring

---

# 23. Testing

Test:

- training coverage calculation;
- mandatory compliance;
- overdue calculation;
- training hours;
- cost aggregation;
- in-house vs external;
- competency gap;
- effectiveness;
- organization filtering;
- position filtering;
- period filtering;
- workforce readiness.

---

# 24. Prinsip Final

Workforce Intelligence menjawab:

> **"Apakah workforce kita sudah mendapatkan training yang dibutuhkan untuk memenuhi kebutuhan organisasi?"**

Bukan:

> "Siapa yang harus mengikuti course?"

Keputusan operasional tetap dilakukan di Training & Development.

Output Workforce Intelligence:

```text
Workforce Need
Training Gap
Training Coverage
Training Compliance
Training Investment
Training Effectiveness
Workforce Readiness
Training Risk
Training Forecast
```

---

# 25. Integrasi dengan Recruitment — Strategic Layer

> Referensi silang ke **`docs/module-recruitment-strategic-layer-plan.md`** (scoping 2026-08-12).

Plan ini berfokus pada **training analytics** di Workforce Intelligence. Item strategis WI yang berkaitan dengan **Recruitment (operasional)** — workforce gap → hiring need → requisition, expected hires → remaining workforce gap, candidate search, dan Quality of Hire — **tidak** dikelola di sini, melainkan di **`docs/module-recruitment-strategic-layer-plan.md`**:

```text
Strategic Layer Integration Plan (S-1 s.d. S-7)
├── S-1  Workforce Gap → Requisition          (workforce_gap_id, hiring need)
├── S-2  Expected Hires → Remaining Gap       (WI mengonsumsi accepted offers)
├── S-3  Candidate Search & Recruitment Analytics
└── S-6  Quality of Hire                      (recruitment score + probation + performance + retention)
```

Hubungan kedua dokumen:

```text
Workforce Intelligence
├── Training analytics        → plan ini (workforce-intelligence-training-enhancement-plan.md)
└── Recruitment analytics     → module-recruitment-strategic-layer-plan.md (S-1/S-2/S-3/S-6)
```

Recruitment tetap **module operasional** (plan terpisah: `docs/module-recruitment-development-plan.md`) — ia hanya menyediakan data; WI yang menghitung kebutuhan workforce & hiring need.
