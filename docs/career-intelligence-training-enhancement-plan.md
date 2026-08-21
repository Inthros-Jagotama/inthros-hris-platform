# Career Intelligence — Training Enhancement Plan

> ✅ **Status implementasi (2026-08-21): P0 SELESAI, dianggap cukup untuk saat ini.**
> Dikerjakan via `docs/superpowers/plans/2026-08-21-career-intelligence-training-integration.md` (6 task, semua ter-commit & lulus test):
> - Employee Training Profile ✅ · Training History for Career ✅ · Training & Competency Gap ✅ · Career Training Recommendation ✅
> - Endpoint baru: `GET /training/employees/:employeeId/summary`, `GET /training/courses/by-competency`, `GET /career-intelligence/employees/:employeeId/training-profile`, `GET /career-intelligence/employees/:employeeId/training-recommendations` — semua read-only, tanpa migrasi baru.
> - FE: `GapAnalysis.vue` menampilkan training profile + rekomendasi training per competency gap.
> - **Belum dikerjakan (keputusan: cukup untuk saat ini, bukan blocker):** Feature 5 (Career Development Plan integration) dan Feature 6/§11/§25 (multi-factor Career Readiness score) — keduanya butuh tabel/entity baru + keputusan produk (bentuk development plan, bobot readiness) yang di luar scope P0. Feature 9/10/12 (Talent Mapping evidence, Succession readiness, Certification eligibility) juga belum disentuh (P1/P2 di §27).
> - Verifikasi end-to-end di browser belum bisa dilakukan — tenant dev belum punya data competency assessment real.

## 1. Tujuan

Enhancement ini menambahkan kemampuan **Training & Development** ke dalam **Career Intelligence** tanpa membuat module Intelligence baru.

Fokus Career Intelligence adalah menjawab:

> **"Apakah training yang dimiliki dan dijalani karyawan mendukung perkembangan karier, competency gap, career readiness, talent, dan succession?"**

Training tetap menjadi **operational source of truth**.

Career Intelligence hanya menggunakan data Training sebagai **career development evidence**.

---

# 2. Scope

Enhancement mencakup:

1. Employee Training Profile
2. Training History for Career
3. Training & Competency Gap
4. Career Training Recommendation
5. Career Development Plan Integration
6. Career Readiness
7. Certification & Career Eligibility
8. Training Effectiveness as Career Evidence
9. Talent Mapping Evidence
10. Succession Candidate Readiness
11. Development Progress
12. Career Training Analytics

Tidak termasuk:

- course CRUD;
- session management;
- participant registration;
- attendance transaction;
- assessment transaction;
- certificate issuance.

Semua proses tersebut tetap berada di **Training & Development**.

---

# 3. Prinsip Arsitektur

```text
Training & Development
        │
        │ operational data
        ▼
Career Intelligence
        │
        ├── Career Path
        ├── Competency Gap
        ├── Career Recommendation
        ├── Career Readiness
        ├── Talent Mapping
        └── Succession Planning
```

Training completion **tidak otomatis berarti employee competent**.

Training adalah evidence.

---

# 4. Data Source

Career Intelligence menggunakan:

```text
Employee
Organization
Position
Performance
Competency
Career Path
Career Assessment
Potential
Training Course
Training Course Competency
Training Participant
Training Assessment Result
Training Effectiveness
Training Certificate
Employee Certification
Training Mandatory
```

---

# 5. Feature 1 — Employee Training Profile

Pada profile Career Intelligence tampilkan:

```text
Employee: Budi

Training Summary
----------------
Total Training       12
Completed            11
Failed                1
Training Hours       86
Average Score        87
Certification         3
Competency Training   8
```

Tujuannya untuk melihat development history.

---

# 6. Feature 2 — Career Training History

Tampilkan training yang relevan dengan career path.

Contoh:

```text
Career Target:
IT Supervisor
```

Training:

```text
Leadership Development
Completed
Score: 92
Relevant: YES

Project Management
Completed
Score: 88
Relevant: YES

Office Safety
Completed
Score: 90
Relevant: NO
```

Tidak semua training harus menjadi career evidence.

Training menjadi career-relevant terutama jika:

```text
Course
   ↓
Competency
```

atau course ditetapkan sebagai development activity.

---

# 7. Feature 3 — Training & Competency Gap

Contoh:

```text
Career Target:
IT Supervisor
```

Required competency:

```text
Leadership             Level 4
Project Management     Level 3
Communication          Level 4
```

Current:

```text
Leadership             Level 2
Project Management     Level 3
Communication          Level 3
```

Gap:

```text
Leadership             Gap 2
Communication          Gap 1
```

Career Intelligence mencari course yang terkait:

```text
training_course_competencies
```

Hasil:

```text
Recommended Development

Leadership Development
Communication Skills
```

---

# 8. Feature 4 — Career Training Recommendation

Recommendation berdasarkan:

```text
Career Target
Required Competency
Current Competency
Competency Gap
Course Competency Mapping
Training History
Certification Requirement
```

Contoh:

```text
Target Position:
IT Supervisor
```

Recommendation:

| Training | Competency | Gap | Priority |
|---|---|---:|---|
| Leadership Development | Leadership | 2 | High |
| Communication | Communication | 1 | Medium |

---

# 9. Training Recommendation Score

Recommendation dapat menggunakan:

```text
Competency Gap
Career Requirement
Course Relevance
Training History
Assessment Result
Certification Requirement
```

Contoh:

```text
Leadership Development
Career Relevance     100
Competency Gap        95
Course Relevance      95

Recommendation Score 96
```

Formula harus configurable pada Career Intelligence dan tidak hard-coded pada Training.

---

# 10. Feature 5 — Career Development Plan

Training menjadi salah satu development action.

Contoh:

```text
Career Goal
IT Supervisor
       │
       ▼
Competency Gap
Leadership 2 → 4
       │
       ▼
Development Action
Leadership Development
       │
       ▼
Target
December 2026
```

Jika Career Intelligence sudah memiliki Development Plan:

**Jangan membuat development plan kedua di Training.**

Training hanya menjadi referensi:

```text
development_plan_item
        ↓
training_course_id
```

---

# 11. Feature 6 — Career Readiness

Training merupakan salah satu faktor career readiness.

```text
Career Readiness
├── Performance
├── Competency
├── Potential
├── Experience
├── Training
├── Certification
└── Career Assessment
```

Training evidence:

```text
Required Training
Completed Training
Training Score
Assessment Result
Training Effectiveness
Certification
```

Contoh:

```text
Target:
IT Supervisor

Required Training:
Leadership Development      ✓
Project Management          ✓
Supervisor Certification    ✗

Training Readiness:
80%
```

Nilai final Career Readiness tetap dihitung menggunakan faktor lain.

---

# 12. Feature 7 — Certification & Career Eligibility

Certification dapat menjadi requirement posisi.

Contoh:

```text
Target:
IT Security Manager

Required:
CISSP
```

Employee:

```text
Budi
CISSP
Status: VALID
```

Result:

```text
Certification Requirement = PASSED
```

Jika expired:

```text
CISSP
Status: EXPIRED
```

Result:

```text
Certification Requirement = GAP
```

Career Intelligence kemudian dapat merekomendasikan:

```text
Certification Renewal
```

atau training terkait.

---

# 13. Feature 8 — Training Effectiveness as Career Evidence

Completion saja tidak cukup.

Contoh:

```text
Training:
Leadership Development

Completed:
YES

Score:
92

Effectiveness:
LOW
```

Career Intelligence tidak boleh menganggap:

```text
Completed = Competent
```

Evidence harus dipisahkan:

```text
Training Completion
       +
Assessment
       +
Competency Assessment
       +
Effectiveness
```

Competency level tetap berasal dari **Competency Management / assessment process**.

---

# 14. Feature 9 — Competency Development Evidence

Course:

```text
Leadership Development
```

Mapping:

```text
Leadership
Target Level 4
```

Employee:

```text
Before:
Leadership Level 2

Training:
Completed

Post Assessment:
90

After Competency Assessment:
Level 4
```

Career Intelligence dapat menampilkan:

```text
Development Progress

Leadership
2 → 4
Status: Achieved
```

Training menyediakan evidence, sedangkan competency assessment menetapkan level resmi.

---

# 15. Feature 10 — Talent Mapping

Training menjadi salah satu evidence talent.

Contoh:

```text
Budi

Performance        High
Potential          High
Competency Gap     Low
Training Progress  High
Certification      3
Career Readiness   Ready Soon
```

Talent box tetap dihitung dari model Talent Mapping yang sudah ada.

Training hanya menjadi supporting evidence.

---

# 16. Feature 11 — Succession Planning

Untuk posisi:

```text
IT Manager
```

Succession candidates:

```text
Budi
Andi
Dedi
```

Career Intelligence menampilkan training readiness:

| Candidate | Performance | Competency | Training | Readiness |
|---|---:|---:|---:|---|
| Budi | 91 | 88 | 95 | Ready Now |
| Andi | 87 | 82 | 80 | Ready in 1–2 Years |
| Dedi | 84 | 75 | 70 | Development Required |

Training bukan penentu tunggal succession.

Candidate assessment tetap mempertimbangkan:

```text
Performance
Competency
Potential
Experience
Training
Certification
Career Assessment
```

---

# 17. Feature 12 — Development Progress

Career Intelligence dapat menampilkan:

```text
Development Plan

Leadership Development       COMPLETED
Communication Training       COMPLETED
Project Management           IN PROGRESS
Certification                PENDING
```

Progress:

```text
Completed 2 / 4
Progress = 50%
```

---

# 18. Feature 13 — Career Training Gap

Dashboard employee:

```text
Career Target:
IT Supervisor

Required Training       5
Completed               3
Remaining               2

Required Certification  2
Valid                   1
Gap                     1
```

Gap dapat dikelompokkan:

```text
Training Gap
Certification Gap
Competency Gap
Experience Gap
```

---

# 19. Feature 14 — Career Path Training Requirement

Setiap career transition dapat memiliki training requirement.

Contoh:

```text
Career Path

IT Staff
   ↓
IT Supervisor
```

Requirement:

```text
Competency:
Leadership Level 4
Project Management Level 3

Training:
Leadership Development
Project Management

Certification:
Supervisor Certification
```

Career Intelligence dapat menampilkan:

```text
Career Transition Readiness

Competency       80%
Training         100%
Certification     50%
Experience        90%

Overall:
Ready Soon
```

Bobot harus mengikuti konfigurasi Career Intelligence yang berlaku.

---

# 20. Feature 15 — Training as Development Recommendation

Career Intelligence dapat memberikan:

```text
Recommended Development

1. Leadership Development
   Reason:
   Leadership competency gap = 2

2. Communication Skills
   Reason:
   Target position requires level 4

3. Project Management
   Reason:
   Required certification pending
```

Recommendation tidak otomatis mendaftarkan employee ke training.

User tetap melakukan action melalui Training Module.

---

# 21. Feature 16 — Training vs Career Readiness

Gunakan status:

```text
NOT_REQUIRED
REQUIRED
RECOMMENDED
IN_PROGRESS
COMPLETED
FAILED
EXPIRED
```

Contoh:

```text
Leadership Training
Status: COMPLETED

Project Management
Status: REQUIRED

Certification
Status: EXPIRED
```

---

# 22. Recommended Data Model

Tidak perlu membuat transaction table baru di Career Intelligence untuk menggandakan data Training.

Gunakan:

```text
Training Operational Data
        ↓
Career Intelligence Read Model / Query
```

Jika membutuhkan snapshot atau performance analytics, dapat dibuat projection:

```text
career_training_summary
career_training_gaps
career_training_recommendations
career_training_readiness
```

Namun tabel projection hanya diperlukan jika arsitektur dan volume data membutuhkan analytical/read model.

---

# 23. API Integration

Training menyediakan:

```http
GET /training/employees/{employeeId}/history
GET /training/employees/{employeeId}/summary
GET /training/employees/{employeeId}/certifications
GET /training/employees/{employeeId}/competency-training
GET /training/employees/{employeeId}/mandatory-status
GET /training/courses/{courseId}/competencies
GET /training/employees/{employeeId}/effectiveness
```

Career Intelligence:

```http
GET /career-intelligence/employees/{employeeId}/training-profile
GET /career-intelligence/employees/{employeeId}/training-gaps
GET /career-intelligence/employees/{employeeId}/training-recommendations
GET /career-intelligence/employees/{employeeId}/career-readiness
GET /career-intelligence/employees/{employeeId}/development-progress
```

---

# 24. Dashboard

## Career Development Overview

```text
Training Completed
Training Hours
Average Score
Certification
Competency Training
Development Progress
```

## Career Readiness

```text
Target Position
Required Training
Completed Training
Training Gap
Certification Gap
Competency Gap
Career Readiness
```

## Development Recommendation

```text
Recommended Course
Competency Gap
Career Target
Priority
Status
```

## Succession

```text
Position
Candidate
Training Readiness
Competency Readiness
Certification
Overall Readiness
```

---

# 25. Career Intelligence Training Scoring

Jangan membuat Training Score menjadi Career Readiness secara langsung.

Contoh:

```text
Training Score = 90
```

tidak berarti:

```text
Career Readiness = 90
```

Gunakan Training sebagai salah satu component:

```text
Career Readiness
=
Performance
+
Competency
+
Potential
+
Experience
+
Training
+
Certification
```

Bobot mengikuti konfigurasi Career Intelligence.

---

# 26. Important Business Rules

1. Training completion tidak otomatis menaikkan competency.
2. Training score tidak sama dengan competency score.
3. Certificate valid tidak otomatis berarti employee ready untuk promosi.
4. Career Intelligence hanya menggunakan training yang relevan terhadap career/development.
5. Mandatory compliance tidak selalu merupakan career development.
6. Training recommendation tidak otomatis melakukan enrollment.
7. Succession decision tidak ditentukan hanya berdasarkan training.
8. Competency level resmi berasal dari Competency Assessment.
9. Career Readiness tetap dihitung berdasarkan seluruh faktor yang dikonfigurasi.
10. Training menjadi evidence, bukan decision maker.

---

# 27. Development Priority

## P0

1. Employee Training Profile
2. Training History
3. Training ↔ Competency Analysis
4. Career Training Gap
5. Career Training Recommendation
6. Career Development Plan integration

## P1

7. Career Readiness
8. Certification & Eligibility
9. Development Progress
10. Career Path Training Requirement
11. Training Effectiveness Evidence

## P2

12. Talent Mapping integration
13. Succession Planning integration
14. Advanced recommendation scoring
15. Career training analytics
16. Career readiness projection

---

# 28. Testing

Test:

- training history;
- course competency mapping;
- competency gap;
- career training recommendation;
- development plan integration;
- certification eligibility;
- certificate expiry;
- training effectiveness;
- career readiness;
- succession candidate ranking;
- employee filtering;
- career target filtering.

---

# 29. Integration Flow

```text
Employee
   │
   ├── Performance
   ├── Competency
   ├── Potential
   ├── Experience
   └── Career Path
          │
          ▼
Career Intelligence
          │
          ▼
Competency / Career Gap
          │
          ▼
Training Recommendation
          │
          ▼
Training & Development
          │
          ▼
Training Completion
          │
          ├── Assessment
          ├── Certification
          └── Effectiveness
          │
          ▼
Career Intelligence
          │
          ├── Development Progress
          ├── Career Readiness
          ├── Talent Mapping
          └── Succession Planning
```

---

# 30. Prinsip Final

Career Intelligence menjawab:

> **"Apakah training yang dilakukan membantu karyawan memenuhi kebutuhan kariernya?"**

Training menjawab:

> **"Training apa yang dilakukan, kapan, oleh siapa, dan apa hasilnya?"**

Career Intelligence menggunakan data Training untuk:

```text
Career Path
Competency Gap
Development Recommendation
Career Readiness
Talent Mapping
Succession Planning
```

Tidak dibuat module `Training Intelligence`.

Training tetap menjadi operational source of truth, sedangkan Career Intelligence menjadi layer analisis karier.
