# Development Plan — Workforce Intelligence Enhancement

> ⚠️ **Status aktual (2026-08-21): dokumen ini usang.** Diverifikasi langsung ke kode — backend keempat fitur (Headcount Planning, Risk Dashboard, Executive Dashboard, Scenario Planning) route-nya sudah ada, bukan "belum tersedia" seperti klaim §2.2 di bawah. Yang benar-benar belum ada adalah **frontend**-nya (4 halaman Vue) — progres: **Headcount Planning ✅ FE selesai, Risk Dashboard ✅ FE selesai (+ 2 gap backend diperbaiki), Executive Dashboard ⏳ backend ternyata stub penuh (lihat §2.4, FE ditunda), Scenario Planning belum dicek**. Detail lengkap di §2.3/§2.4 (baru).

## 1. Overview

Dokumen ini merupakan development plan untuk menyelesaikan fitur Workforce Intelligence yang belum tersedia:

1. Headcount Planning
2. Risk Dashboard
3. Executive Dashboard
4. Scenario Planning

Fitur yang sudah tersedia dan menjadi data source:

- Candidate Search
- Recruitment Analytics
- Quality of Hire
- Training Analysis

Development harus memanfaatkan module HRIS yang sudah tersedia dan menghindari duplikasi data maupun business logic.

---

# 2. Existing Workforce Intelligence

## 2.1 Completed

| Feature | Status |
|---|---|
| Candidate Search | DONE |
| Recruitment Analytics | DONE |
| Quality of Hire | DONE |
| Training Analysis | DONE |

## 2.2 Planned (per dokumen asli — sudah tidak akurat, lihat §2.3)

| Feature | Priority |
|---|---|
| Headcount Planning | P0 |
| Risk Dashboard | P1 |
| Executive Dashboard | P1 |
| Scenario Planning | P2 |

## 2.3 Status Aktual per 2026-08-21 (koreksi)

Audit langsung ke `backend/internal/modules/workforceintelligence/` (`model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go` — plus test lengkap: `handler_test.go`, `repository_test.go`, `service_test.go`, `helpers_test.go`) dan `frontend/tenant/src/views/modules/workforce-intelligence/`.

| Feature | Backend | Frontend |
|---|---|---|
| Candidate Search | ✅ | ✅ `CandidateSearch.vue` |
| Recruitment Analytics | ✅ | ✅ `RecruitmentAnalytics.vue` |
| Quality of Hire | ✅ | ✅ `QualityOfHire.vue` |
| Training Analysis | ✅ (`GET /analytics/learning`) | ✅ — tapi bukan halaman workforce-intelligence sendiri, card-nya di hub `WorkforceIntelligence.vue` mengarah ke `/training/reports` (halaman modul Training) |
| **Headcount Planning** | ✅ **lengkap** (model `WorkforcePlanningHeadcount`, migration `017_workforce_intelligence`, repo/service/handler penuh) | ✅ **SELESAI (2026-08-21)** — `HeadcountPlanning.vue` |
| **Risk Dashboard** | ✅ **lengkap** — diperbaiki 2026-08-21: endpoint create ditambahkan (tidak ada sama sekali sebelumnya) + 4 widget detail diganti dari hardcoded demo ke data real (lihat §2.4) | ✅ **SELESAI (2026-08-21)** — `RiskDashboard.vue` |
| **Executive Dashboard** | ⚠️ **route ada, tapi SEMUA 8 endpoint stub hardcoded** — lihat §2.4, belum layak dipakai | ⏳ **ditunda** — FE tidak dibangun sampai backend diperbaiki (supaya tidak menampilkan angka palsu) |
| **Scenario Planning** | route ada, belum diaudit apakah real/stub | ❌ belum ada — card "Coming soon" di hub |

**Endpoint lengkap yang sudah aktif** (`backend/internal/modules/workforceintelligence/routes.go`, prefix `/api/v1/tenant/workforce-intelligence`):

```http
GET    /candidate-search

GET    /planning/headcounts
POST   /planning/headcounts
GET    /planning/headcounts/:id
PUT    /planning/headcounts/:id
DELETE /planning/headcounts/:id

GET    /planning/forecasts
POST   /planning/forecasts
GET    /planning/forecasts/:id
PUT    /planning/forecasts/:id
DELETE /planning/forecasts/:id

GET    /planning/gap-analysis
GET    /planning/projections

GET    /kpi
GET    /kpi/summary
GET    /kpi/:code

GET    /analytics/headcount
GET    /analytics/attendance
GET    /analytics/leave
GET    /analytics/overtime
GET    /analytics/payroll
GET    /analytics/performance
GET    /analytics/learning
GET    /analytics/recruitment
GET    /analytics/quality-of-hire
GET    /analytics/movement

GET    /capacity/dashboard
GET    /capacity/utilization
GET    /capacity/forecast
GET    /capacity/bottlenecks

GET    /cost/summary
GET    /cost/payroll
GET    /cost/per-employee
GET    /cost/per-department
GET    /cost/budget-vs-actual

GET    /risk/dashboard
GET    /risk/indicators
GET    /risk/indicators/:id
PUT    /risk/indicators/:id
GET    /risk/high-turnover
GET    /risk/retirement
GET    /risk/contract-expiry
GET    /risk/high-absenteeism

GET    /executive/summary
GET    /executive/growth
GET    /executive/cost-trend
GET    /executive/attrition-trend
GET    /executive/capacity
GET    /executive/hiring-progress
GET    /executive/risk-overview
GET    /executive/health-score

GET    /scenarios
POST   /scenarios
GET    /scenarios/:id
PUT    /scenarios/:id
DELETE /scenarios/:id
POST   /scenarios/:id/run
POST   /scenarios/:id/clone

GET    /health/dashboard
GET    /health/scores
GET    /health/scores/:id
GET    /health/span-of-control
GET    /health/succession

GET    /people-analytics/training-vs-performance
GET    /people-analytics/overtime-vs-productivity
GET    /people-analytics/attendance-vs-performance
GET    /people-analytics/compensation-vs-turnover
GET    /people-analytics/source-vs-retention
GET    /people-analytics/career-progression
GET    /people-analytics/learning-effectiveness
```

**Catatan skema**: nama tabel aktual (migration `017_workforce_intelligence`, mysql & postgres) sedikit berbeda dari yang direncanakan di §11/§35/§36 dokumen ini (`workforce_planning_headcounts` bukan `headcount_plans`, `workforce_risk_indicators`, `workforce_scenarios`, dst) — bukan masalah, cuma penamaan berbeda dari rencana awal; skema fungsional sudah setara.

**Kesimpulan**: §3–§41 dokumen ini (arsitektur, formula, API plan, dst.) sudah **terwujud di backend** — bacaan bagus untuk memahami intent desainnya, tapi bukan lagi "to-do", melainkan dokumentasi retroaktif dari apa yang sudah ada. Satu-satunya pekerjaan nyata yang tersisa dari dokumen ini adalah **§42 Frontend Structure** — 4 halaman Vue (Headcount Planning, Risk Dashboard, Executive Dashboard, Scenario Planning) yang mengonsumsi endpoint di atas.

---

# 3. Target Architecture

```text
                         WORKFORCE INTELLIGENCE
                                  │
       ┌──────────────────────────┼──────────────────────────┐
       │                          │                          │
       ▼                          ▼                          ▼
Candidate Search          Recruitment Analytics       Training Analysis
       │                          │                          │
       └──────────────────────────┼──────────────────────────┘
                                  │
                                  ▼
                         Workforce Data Layer
                                  │
             ┌────────────────────┼────────────────────┐
             │                    │                    │
             ▼                    ▼                    ▼
     Headcount Planning      Risk Dashboard      Executive Dashboard
             │                    │                    │
             └────────────────────┼────────────────────┘
                                  │
                                  ▼
                         Scenario Planning
```

Scenario Planning berada di layer paling atas karena menggunakan data dari Headcount, Risk, Recruitment, Training, Talent, dan cost.

---

## 2.4 TODO — Executive Dashboard backend (catatan untuk pekerjaan selanjutnya)

**Temuan (2026-08-21):** sebelum membangun FE Executive Dashboard, audit ke `service.go`/`handler.go` menunjukkan **kedelapan endpoint di bawah `/executive/*` adalah stub hardcoded** — pola identik dengan `GetRiskDetail` sebelum diperbaiki (lihat log perbaikan Risk Dashboard di atas), tapi cakupannya lebih luas (8 endpoint vs 4 widget).

| Endpoint | Fungsi | Status |
|---|---|---|
| `/executive/summary` | `GetExecutiveSummary` | Sebagian real — `TotalHC` dari `GetActiveEmployeeCount`; `HCGrowth`, `AttritionRate`, `AvgCost`, `UtilizationRate`, `HealthScore` semua angka literal (komentar kode sendiri: `// Would be computed from historical data`) |
| `/executive/growth` | `GetExecutiveGrowth` | Sebagian real — `Current` dari live HC count; `Trend: []` kosong, `Change: 5.2` literal |
| `/executive/cost-trend` | `GetExecutiveCostTrend` | **Stub penuh** — `Current: 8500000000`, `Change: 8.3`, tidak ada query ke data payroll sama sekali |
| `/executive/attrition-trend` | `GetExecutiveAttritionTrend` | **Stub penuh** — `Current: 12.3`, `Change: -1.1`, tidak ada query |
| `/executive/capacity` | `GetExecutiveCapacity` | **Stub penuh** — `UtilizationRate: 78.5`, `AvailableHC: 1245`, dll, semua literal |
| `/executive/hiring-progress` | inline di handler (bukan lewat service) | **Stub penuh** — `Planned: 50, InProgress: 25, Completed: 15, Total: 50` hardcoded langsung di `handler.go` |
| `/executive/risk-overview` | `GetExecutiveRiskOverview` | **Stub penuh** — `TotalRisks: 24, HighRiskCount: 7, CriticalCount: 2` literal, padahal `s.repo.ListRiskIndicators` (data real, dipakai Risk Dashboard) sudah tersedia di file yang sama dan tidak dipanggil di sini |
| `/executive/health-score` | `GetExecutiveHealthScore` | **Stub penuh** — `Score: 72.5` sama persis dengan angka literal `HealthScore` di `GetExecutiveSummary` (indikasi copy-paste, bukan hasil hitung); tidak delegasi ke `/health/dashboard` yang datanya sudah real (`WorkforceHealthScore`, tersimpan di DB) |

Tambahan: tidak satu pun dari 8 handler menerima query param (`period`, `organization_id`) meski beberapa response field (`Period`) menyiratkan seharusnya bisa difilter.

**Rekomendasi perbaikan (belum dikerjakan, scope besar — 4 modul sumber data):**
1. `risk-overview` — paling mudah: tinggal reuse `s.repo.ListRiskIndicators` yang sudah ada (persis pola `GetRiskDashboard`), agregasi by-department/by-category dari `WorkforceRiskIndicator.DepartmentID`.
2. `health-score` — delegasikan ke computation/storage yang sama dengan `GetHealthDashboard`/`ListHealthScores` (§ Organization Health, `WorkforceHealthScore`) alih-alih duplikasi angka hardcoded.
3. `summary`/`growth` — lengkapi bagian yang masih literal (`AttritionRate` dari `employee_movements` movement_type=offboarding, sama seperti perbaikan Risk Dashboard's high-turnover widget; `Trend` histori bulanan perlu query baru, ikuti pola `GetMonthlyHeadcountTrend`).
4. `cost-trend`/`capacity`/`hiring-progress` — butuh riset dulu data apa yang tersedia di modul Payroll (cost) dan Recruitment/Requisition (hiring progress: planned vs filled) sebelum bisa di-wire; belum diriset.

**Keputusan (2026-08-21):** FE Executive Dashboard **ditunda**, tidak dibangun di atas data palsu. Dicatat sebagai TODO terpisah, dikerjakan setelah Scenario Planning diaudit/diselesaikan.

---

# 4. HEADCOUNT PLANNING

## 4.1 Objective

Menyediakan kemampuan untuk:

- Merencanakan kebutuhan employee.
- Membandingkan current HC dengan planned HC.
- Menghitung workforce gap.
- Memprediksi headcount.
- Mengidentifikasi kebutuhan recruitment.
- Mengidentifikasi kebutuhan internal movement.
- Menghubungkan kebutuhan HC dengan talent/succession.
- Mengestimasi manpower cost.

Karena HRIS tidak menggunakan entity `Position`, Headcount Planning menggunakan:

```text
Organization = Position/Jabatan
```

sebagai unit perencanaan.

---

# 5. Headcount Planning Concept

## 5.1 Current Headcount

Current HC dihitung dari employee aktif pada organization.

```text
Current HC =
COUNT(active employees)
```

---

## 5.2 Planned Headcount

Jumlah employee yang direncanakan tersedia.

```text
Organization:
Backend Developer

Current HC = 4
Planned HC = 6
```

---

## 5.3 Forecast Headcount

Formula:

```text
Forecast HC =
Current HC
+ Planned Join
- Planned Exit
+ Internal In
- Internal Out
```

---

## 5.4 Workforce Gap

```text
Gap =
Planned HC - Forecast HC
```

Contoh:

```text
Current HC       100
Planned Join       5
Expected Exit     10
Forecast HC       95
Planned HC       110

Gap               15
```

---

# 6. Headcount Fulfillment

Gap tidak langsung menjadi recruitment.

Sistem harus mengidentifikasi:

```text
Gap
 │
 ├── Internal Movement
 ├── Promotion
 ├── Succession
 └── Recruitment
```

Contoh:

```text
Gap = 15

Internal Movement = 5
Succession        = 3
Recruitment       = 7
```

Maka Recruitment Requirement:

```text
7 employees
```

---

# 7. Headcount Planning Period

Support:

- Annual planning.
- Quarterly planning.
- Monthly forecast.

Minimum:

```text
Planning Year
2027
```

Optional monthly forecast:

```text
Jan
Feb
Mar
...
Dec
```

---

# 8. Attrition Forecast

Headcount Planning harus menggunakan historical employee movement untuk memprediksi attrition.

Data source:

- Employee Movement.
- Resignation.
- Retirement.
- Contract End.
- Termination.

Contoh:

```text
Historical Attrition Rate = 8%

Current HC = 1,000

Expected Attrition =
1,000 × 8%
= 80 employees
```

Forecast:

```text
1,000 - 80 = 920
```

---

# 9. Headcount Cost

Integrasikan dengan Payroll/Salary Structure.

Cost dapat mencakup:

- Basic salary.
- Allowance.
- Employer BPJS.
- Benefit.
- Other employment cost.

Output:

```text
Current Cost
Planned Cost
Incremental Cost
Annual Cost
```

Contoh:

```text
Current HC       100
Planned HC       120

Additional HC     20

Estimated Annual Cost:
Rp X
```

---

# 10. Headcount Dashboard

Dashboard harus menampilkan:

```text
Current HC
Planned HC
Forecast HC
Gap
Expected Join
Expected Exit
Internal Fill
Recruitment Requirement
Estimated Cost
```

Visual:

- HC trend.
- HC by organization.
- Gap by organization.
- Join vs exit.
- Recruitment requirement.
- Cost projection.

---

# 11. Headcount Database

Recommended tables:

```text
headcount_plans
headcount_plan_details
headcount_plan_monthly_details
headcount_plan_fulfillments
```

Use UUID dan mengikuti standard database HRIS existing.

---

# 12. Headcount Integration

Integrasi:

```text
Organization
     ↓
Employee
     ↓
Headcount Planning
     ↓
Employee Movement
     ↓
Career / Talent
     ↓
Recruitment
     ↓
Payroll
```

---

# 13. Headcount Approval

Gunakan existing Approval Engine.

Lifecycle:

```text
DRAFT
 ↓
SUBMITTED
 ↓
IN_REVIEW
 ↓
APPROVED
 ↓
FINALIZED
```

Rejected:

```text
IN_REVIEW
 ↓
REJECTED
 ↓
DRAFT
```

Tidak membuat approval engine baru.

---

# 14. RISK DASHBOARD

## 14.1 Objective

Risk Dashboard digunakan untuk mengidentifikasi risiko workforce yang dapat mengganggu operasional organisasi.

Risk tidak hanya berarti employee akan resign.

Risk dapat mencakup:

- Attrition Risk.
- Critical Position Risk.
- Talent Risk.
- Succession Risk.
- Skill/Competency Risk.
- Recruitment Risk.
- Training Risk.
- Workforce Capacity Risk.
- Retirement Risk.

---

# 15. Workforce Risk Categories

## 15.1 Attrition Risk

Mengidentifikasi employee/organization dengan risiko turnover tinggi.

Input:

- Historical turnover.
- Employee tenure.
- Employee movement.
- Performance.
- Career data jika tersedia.
- Engagement data jika tersedia.
- Compensation data jika tersedia.

Output:

```text
Low
Medium
High
Critical
```

---

# 16. Critical Workforce Risk

Identifikasi organization yang sangat bergantung pada sedikit employee.

Contoh:

```text
Organization:
Senior Database Engineer

Current HC:
2

Critical Employee:
1

Succession:
None
```

Risk:

```text
CRITICAL
```

Karena kehilangan satu employee dapat mengganggu operasional.

---

# 17. Succession Risk

Integrasikan dengan Talent/Succession.

Contoh:

```text
Organization:
IT Manager

Incumbent:
1

Successor:
0

Risk:
HIGH
```

Jika:

```text
Successor Ready Now = 0
Successor Ready Soon = 0
```

maka:

```text
Succession Risk = CRITICAL
```

---

# 18. Competency Risk

Bandingkan:

```text
Required Competency
        vs
Available Competency
```

Contoh:

```text
Required:
Leadership       80
Strategic        80
Technical        85

Actual:
Leadership       65
Strategic        70
Technical        90
```

Gap:

```text
Leadership      -15
Strategic       -10
Technical        +5
```

Risk:

```text
Competency Risk = HIGH
```

---

# 19. Recruitment Risk

Gunakan Recruitment Analytics.

Contoh:

```text
Open Requirement       50
Average Time to Hire   60 days
Available Candidate    12
```

Jika kebutuhan sangat besar tetapi candidate supply rendah:

```text
Recruitment Risk = HIGH
```

---

# 20. Training Risk

Gunakan Training Analysis.

Contoh:

```text
Required Training Completion = 95%

Actual Completion = 70%
```

Risk:

```text
Training Compliance Risk = HIGH
```

---

# 21. Workforce Capacity Risk

Bandingkan:

```text
Required HC
vs
Available HC
```

Contoh:

```text
Required HC = 200
Available HC = 160

Capacity Gap = 40
```

Risk:

```text
Capacity Risk = HIGH
```

---

# 22. Risk Scoring

Gunakan standardized score:

```text
0 - 24    LOW
25 - 49   MEDIUM
50 - 74   HIGH
75 - 100  CRITICAL
```

Setiap risk memiliki:

```text
Risk Score
Risk Level
Risk Factor
Risk Trend
Recommended Action
```

---

# 23. Risk Dashboard

Executive view:

```text
WORKFORCE RISK

Critical       8
High          21
Medium        45
Low          120
```

Risk by category:

```text
Attrition             15
Succession              8
Competency             12
Recruitment             7
Training                9
Capacity               11
```

Visual:

- Risk heatmap.
- Risk trend.
- Risk by organization.
- Risk by job level.
- Risk by risk category.

---

# 24. Risk Heatmap

Contoh:

| Organization | Attrition | Succession | Competency | Recruitment | Overall |
|---|---|---|---|---|---|
| IT | High | Critical | Medium | High | Critical |
| Finance | Low | Medium | Low | Low | Medium |
| Operations | Medium | High | High | Critical | Critical |

Klik organization → drill down employee/risk factor.

---

# 25. Risk Action

Setiap high/critical risk harus dapat memiliki action.

Contoh:

```text
Risk:
Succession Risk - IT Manager

Action:
Create successor development plan

Owner:
HR Manager

Due Date:
2027-03-31

Status:
OPEN
```

Status:

```text
OPEN
IN_PROGRESS
MITIGATED
ACCEPTED
CLOSED
```

---

# 26. EXECUTIVE DASHBOARD

## 26.1 Objective

Executive Dashboard menjadi **single view** kondisi workforce perusahaan.

Dashboard tidak melakukan calculation baru.

Dashboard mengambil hasil dari:

- Headcount Planning.
- Recruitment Analytics.
- Quality of Hire.
- Training Analysis.
- Risk Dashboard.
- Talent/9 Box.
- Employee data.
- Payroll/cost.

---

# 27. Executive KPI

Minimum KPI:

### Workforce

```text
Total Employee
Current Headcount
Planned Headcount
HC Gap
Attrition Rate
```

### Recruitment

```text
Open Requirement
Time to Hire
Cost per Hire
Quality of Hire
```

### Talent

```text
High Potential
Critical Talent
Successor Coverage
Talent Risk
```

### Training

```text
Training Completion
Training Cost
Training Effectiveness
Skill Gap
```

### Risk

```text
Critical Risk
High Risk
Succession Risk
Attrition Risk
```

### Cost

```text
Total Payroll Cost
Training Cost
Recruitment Cost
Projected Workforce Cost
```

---

# 28. Executive Dashboard Layout

```text
WORKFORCE EXECUTIVE DASHBOARD

┌────────────┬────────────┬────────────┬────────────┐
│ Employee   │ HC Gap     │ Attrition  │ Risk       │
│ 1,250       │ +60        │ 7.2%       │ 8 Critical│
└────────────┴────────────┴────────────┴────────────┘

┌───────────────────────┬────────────────────────┐
│ Headcount Trend       │ Workforce Risk         │
│                       │                        │
│ Chart                 │ Heatmap                │
└───────────────────────┴────────────────────────┘

┌───────────────────────┬────────────────────────┐
│ Recruitment            │ Talent                 │
│                       │                        │
│ Requirement / Hiring  │ 9 Box / Successor      │
└───────────────────────┴────────────────────────┘

┌─────────────────────────────────────────────────┐
│ Workforce Cost                                  │
│ Current / Forecast / Budget                     │
└─────────────────────────────────────────────────┘
```

---

# 29. Executive Insights

Dashboard se harus menghasilkan insight, bukan hanya angka.

Contoh:

```text
Workforce Insight

Operations membutuhkan tambahan 60 employee
pada 2027.

30 dapat dipenuhi melalui recruitment,
20 melalui internal movement,
10 melalui succession.

Namun 8 critical workforce risks
teridentifikasi pada Operations.
```

---

# 30. SCENARIO PLANNING

## 30.1 Objective

Scenario Planning digunakan untuk melakukan simulasi:

> "Apa yang terjadi jika kondisi workforce berubah?"

Scenario Planning **tidak mengubah actual HRIS data**.

Semua calculation dilakukan dalam simulation context.

---

# 31. Scenario Types

Minimum scenario:

1. Workforce Growth.
2. Workforce Reduction.
3. Attrition Increase.
4. Recruitment Delay.
5. Budget Reduction.
6. Business Expansion.
7. Business Contraction.
8. Internal Mobility.
9. Automation/Productivity.

---

# 32. Scenario Example — Growth

Current:

```text
HC = 1,000
```

Scenario:

```text
Business Growth = +20%
```

System menghitung:

```text
Required HC = 1,200
Gap = 200
```

Kemudian:

```text
Internal Fill = 50
Succession = 30
Recruitment = 120
```

Estimated cost:

```text
Additional Annual Cost = Rp X
```

---

# 33. Scenario Example — Attrition

Baseline:

```text
Attrition = 7%
```

Scenario:

```text
Attrition = 12%
```

System menghitung impact:

```text
Additional Exit
Additional Recruitment
Additional Cost
Potential Capacity Gap
```

---

# 34. Scenario Example — Budget Reduction

Baseline:

```text
Workforce Budget = Rp 100 M
```

Scenario:

```text
Budget Reduction = 10%
```

System mencari:

```text
Required HC
Possible Hiring Reduction
Critical Organization
Workforce Risk
```

Output:

```text
Budget:
Rp 100 M → Rp 90 M

Potential HC reduction:
45

Critical Risk:
Operations
IT
Customer Service
```

---

# 35. Scenario Structure

Recommended conceptual table:

```text
workforce_scenarios
```

Fields:

```text
id
company_id
name
description
scenario_type
base_period
status
created_by
created_at
updated_at
```

Scenario assumptions:

```text
workforce_scenario_assumptions
```

Contoh:

```text
assumption_type
assumption_value
organization_id
```

Types:

```text
HC_GROWTH
HC_REDUCTION
ATTRITION_RATE
HIRING_DELAY
BUDGET_CHANGE
SALARY_CHANGE
INTERNAL_MOBILITY
```

---

# 36. Scenario Result

```text
workforce_scenario_results
```

Menyimpan:

```text
scenario_id
organization_id
current_hc
projected_hc
hc_gap
recruitment_requirement
internal_fill
estimated_cost
risk_score
```

---

# 37. Scenario Isolation

Scenario tidak boleh mengubah:

- Employee.
- Organization.
- Recruitment.
- Payroll.
- Headcount actual.
- Talent actual.

Scenario hanya melakukan:

```text
Actual Data
    ↓
Copy / Snapshot
    ↓
Scenario Assumption
    ↓
Simulation
    ↓
Scenario Result
```

---

# 38. Scenario Comparison

User dapat membandingkan:

```text
Baseline
vs
Growth 10%
vs
Growth 20%
vs
Budget -10%
```

Contoh:

| Metric | Baseline | Growth 10% | Growth 20% | Budget -10% |
|---|---:|---:|---:|---:|
| HC | 1,000 | 1,100 | 1,200 | 950 |
| HC Gap | 0 | 100 | 200 | -50 |
| Recruitment | 0 | 70 | 120 | 0 |
| Cost | 100M | 112M | 126M | 90M |
| Risk | Low | Medium | High | Critical |

---

# 39. Scenario Dashboard

Tampilkan:

- Scenario name.
- Baseline.
- Assumptions.
- Projected HC.
- HC gap.
- Recruitment requirement.
- Internal fulfillment.
- Cost impact.
- Risk impact.

Charts:

- HC projection.
- Cost projection.
- Recruitment demand.
- Risk score.

---

# 40. Backend Architecture

Recommended service separation:

```text
WorkforceIntelligence
│
├── HeadcountPlanningService
├── WorkforceRiskService
├── ExecutiveDashboardService
└── WorkforceScenarioService
```

Data provider:

```text
WorkforceDataProvider
```

yang mengambil data dari module existing.

Contoh:

```text
WorkforceDataProvider
├── EmployeeProvider
├── OrganizationProvider
├── RecruitmentProvider
├── TrainingProvider
├── TalentProvider
├── PayrollProvider
└── MovementProvider
```

Hindari query langsung ke seluruh module dari controller.

---

# 41. API

## Headcount

```text
GET    /api/workforce/headcount
GET    /api/workforce/headcount/summary
GET    /api/workforce/headcount/by-organization
POST   /api/workforce/headcount/plans
GET    /api/workforce/headcount/plans/{id}
PUT    /api/workforce/headcount/plans/{id}
POST   /api/workforce/headcount/plans/{id}/submit
POST   /api/workforce/headcount/plans/{id}/revision
```

## Risk

```text
GET /api/workforce/risks
GET /api/workforce/risks/summary
GET /api/workforce/risks/by-category
GET /api/workforce/risks/by-organization
GET /api/workforce/risks/{id}
POST /api/workforce/risks/{id}/actions
```

## Executive

```text
GET /api/workforce/executive/summary
GET /api/workforce/executive/workforce
GET /api/workforce/executive/recruitment
GET /api/workforce/executive/talent
GET /api/workforce/executive/training
GET /api/workforce/executive/risk
GET /api/workforce/executive/cost
```

## Scenario

```text
GET    /api/workforce/scenarios
POST   /api/workforce/scenarios
GET    /api/workforce/scenarios/{id}
PUT    /api/workforce/scenarios/{id}
POST   /api/workforce/scenarios/{id}/run
GET    /api/workforce/scenarios/{id}/results
GET    /api/workforce/scenarios/compare
DELETE /api/workforce/scenarios/{id}
```

---

# 42. Frontend Structure

Recommended:

```text
Workforce Intelligence
│
├── Candidate Search
├── Recruitment Analytics
├── Quality of Hire
├── Training Analysis
│
├── Headcount Planning
│   ├── Overview
│   ├── Plans
│   ├── Forecast
│   ├── Gap Analysis
│   └── Cost
│
├── Risk Dashboard
│   ├── Overview
│   ├── Risk Heatmap
│   ├── Risk by Organization
│   ├── Risk Detail
│   └── Mitigation Actions
│
├── Executive Dashboard
│   ├── Workforce
│   ├── Recruitment
│   ├── Talent
│   ├── Training
│   ├── Risk
│   └── Cost
│
└── Scenario Planning
    ├── Scenarios
    ├── Create Scenario
    ├── Simulation
    ├── Results
    └── Comparison
```

---

# 43. Permissions

Recommended:

```text
workforce.headcount.view
workforce.headcount.create
workforce.headcount.update
workforce.headcount.submit
workforce.headcount.approve

workforce.risk.view
workforce.risk.manage
workforce.risk.action

workforce.executive.view

workforce.scenario.view
workforce.scenario.create
workforce.scenario.update
workforce.scenario.run
workforce.scenario.delete
```

---

# 44. Data Security

Workforce Intelligence dapat mengandung sensitive employee information.

Implement:

- Company-level access.
- Organization-level access.
- Role-based access.
- Sensitive field masking.
- Audit log.
- No raw sensitive employee data pada executive dashboard jika tidak diperlukan.
- Aggregation untuk executive view.
- Export permission.

Untuk dashboard executive, gunakan aggregation:

```text
Employee-level data
        ↓
Aggregation
        ↓
Executive Metrics
```

---

# 45. Caching

Dashboard tidak perlu selalu menghitung seluruh data secara realtime.

Gunakan caching untuk:

- Executive KPI.
- Risk summary.
- Headcount summary.
- Recruitment analytics.
- Training metrics.

Contoh:

```text
Workforce Dashboard
       ↓
Cache
       ↓
Data Provider
```

Cache harus di-invalidasi setelah data penting berubah.

---

# 46. Scheduled Jobs

Recommended jobs:

```text
CalculateWorkforceRiskJob
GenerateHeadcountForecastJob
RefreshWorkforceMetricsJob
CalculateAttritionForecastJob
RefreshExecutiveDashboardJob
```

Frequency:

```text
Daily:
Risk / Workforce metrics

Monthly:
Headcount forecast

On demand:
Scenario planning
```

---

# 47. Development Priority

## P0 — Headcount Planning

Implement terlebih dahulu:

```text
Organization
    ↓
Current HC
    ↓
Planned HC
    ↓
Forecast
    ↓
Gap
    ↓
Fulfillment
    ↓
Recruitment Requirement
```

Alasan:

Headcount menjadi foundation untuk Risk, Executive Dashboard, dan Scenario Planning.

---

## P1 — Risk Dashboard

Setelah Headcount tersedia:

```text
Headcount
+
Attrition
+
Talent
+
Competency
+
Recruitment
+
Training
    ↓
Risk Dashboard
```

---

## P1 — Executive Dashboard

Setelah data Headcount dan Risk tersedia:

```text
Headcount
Recruitment
Quality of Hire
Training
Talent
Risk
Cost
    ↓
Executive Dashboard
```

---

## P2 — Scenario Planning

Dibangun terakhir:

```text
Headcount
+
Recruitment
+
Risk
+
Cost
+
Talent
    ↓
Scenario Engine
```

---

# 48. Testing Strategy

## Headcount

Test:

- Current HC.
- Forecast.
- Planned HC.
- Gap.
- Attrition.
- Internal fulfillment.
- Recruitment requirement.
- Cost.
- Approval.
- Revision.

## Risk

Test:

- Risk score.
- Risk classification.
- Attrition risk.
- Succession risk.
- Competency risk.
- Recruitment risk.
- Capacity risk.

## Executive Dashboard

Test:

- KPI aggregation.
- Organization filter.
- Company filter.
- Period filter.
- Cost aggregation.
- Permission.

## Scenario

Test:

- Scenario creation.
- Assumption.
- Simulation.
- Result.
- Comparison.
- Data isolation.

Scenario test harus memastikan:

```text
Running Scenario
≠
Changing Actual HRIS Data
```

---

# 49. Acceptance Criteria

## Headcount Planning

- [ ] Current HC tersedia.
- [ ] Planned HC dapat dibuat.
- [ ] Forecast HC tersedia.
- [ ] Gap otomatis dihitung.
- [ ] Attrition dapat diperhitungkan.
- [ ] Internal fulfillment dapat dicatat.
- [ ] Recruitment requirement otomatis dihitung.
- [ ] Manpower cost tersedia.
- [ ] Approval menggunakan existing Approval Engine.
- [ ] Revision tersedia.
- [ ] Dashboard tersedia.

## Risk Dashboard

- [ ] Risk categories tersedia.
- [ ] Risk score tersedia.
- [ ] Risk level tersedia.
- [ ] Risk heatmap tersedia.
- [ ] Risk by organization tersedia.
- [ ] Critical risk tersedia.
- [ ] Mitigation action tersedia.
- [ ] Risk trend tersedia.

## Executive Dashboard

- [ ] Workforce KPI tersedia.
- [ ] Recruitment KPI tersedia.
- [ ] Talent KPI tersedia.
- [ ] Training KPI tersedia.
- [ ] Risk KPI tersedia.
- [ ] Cost KPI tersedia.
- [ ] Drill-down tersedia.
- [ ] Company/organization filter tersedia.

## Scenario Planning

- [ ] Scenario dapat dibuat.
- [ ] Scenario assumption tersedia.
- [ ] Simulation dapat dijalankan.
- [ ] HC impact tersedia.
- [ ] Recruitment impact tersedia.
- [ ] Cost impact tersedia.
- [ ] Risk impact tersedia.
- [ ] Scenario comparison tersedia.
- [ ] Scenario tidak mengubah actual HRIS data.

---

# 50. Final Integration Model

Target akhir Workforce Intelligence:

```text
                       WORKFORCE INTELLIGENCE
                                │
      ┌───────────────┬─────────┼──────────┬──────────────┐
      │               │         │          │              │
      ▼               ▼         ▼          ▼              ▼
 Candidate        Recruitment  Quality   Training      Headcount
 Search           Analytics    of Hire   Analysis      Planning
      │               │         │          │              │
      └───────────────┴─────────┴──────────┴──────────────┘
                                │
                                ▼
                       WORKFORCE DATA LAYER
                                │
                   ┌────────────┴────────────┐
                   ▼                         ▼
             Risk Dashboard          Executive Dashboard
                   │                         │
                   └────────────┬────────────┘
                                ▼
                       Scenario Planning
                                │
                                ▼
                      Workforce Intelligence
```

---

# 51. Expected Business Outcome

Setelah seluruh fitur selesai, Workforce Intelligence tidak hanya menjadi kumpulan dashboard, tetapi menjadi **decision-support system untuk HR**.

Contoh alur:

```text
Business Plan
     ↓
Headcount Planning
     ↓
HC Gap = 100
     ↓
Talent Internal = 40
     ↓
Recruitment Need = 60
     ↓
Recruitment Analytics
     ↓
Candidate Supply rendah
     ↓
Recruitment Risk HIGH
     ↓
Training / Upskilling
     ↓
Scenario Planning
     ↓
Compare:
- Recruitment 60
- Internal development 40
- Automation
     ↓
Executive Dashboard
     ↓
Management Decision
```

Dengan demikian, empat fitur yang belum selesai memiliki peran yang berbeda:

| Feature | Fungsi Utama |
|---|---|
| **Headcount Planning** | Menentukan kebutuhan workforce |
| **Risk Dashboard** | Menentukan risiko workforce |
| **Executive Dashboard** | Menyajikan kondisi workforce kepada management |
| **Scenario Planning** | Mensimulasikan dampak keputusan workforce |

Urutan development yang direkomendasikan:

**Headcount Planning → Risk Dashboard → Executive Dashboard → Scenario Planning**

karena setiap layer berikutnya dapat menggunakan hasil layer sebelumnya.