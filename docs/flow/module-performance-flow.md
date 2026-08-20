# Alur Performance Management (Runbook)

Dokumen ini menjelaskan **cara pakai** modul **Performance Management** — setup master
(periods, ratings, perspectives, formulas, components), KPI (template, indicators, evaluations
dengan workflow BSC 2-phase), OKR (templates, objectives, key results, evaluations),
scoring engine, dan dashboards.

- Lokasi kode: `backend/internal/modules/performance/`
- Halaman UI: `frontend/tenant/src/views/modules/performance/`
- Module slug: `performance`

---

## 1. Ringkasan Alur End-to-End

```
SETUP MASTER                     KPI                                 OKR
┌──────────────────┐   ┌──────────────────────────────┐   ┌──────────────────────────────┐
│ Periods           │   │ Templates (BSC)              │   │ Templates                     │
│ Ratings           │──▶│ Indicators (target/actual)    │──▶│ Objectives + Key Results      │
│ Perspectives      │   │ Evaluations:                  │   │ Evaluations:                  │
│ Formulas          │   │   DRAFT → TARGET_SUBMITTED    │   │   DRAFT → SUBMITTED           │
│ Components        │   │   → TARGET_APPROVED           │   │   → APPROVED → COMPLETED      │
│ Org Components    │   │   → SUBMITTED → APPROVED      │   │                               │
│                   │   │   → COMPLETED                  │   │                               │
└──────────────────┘   └──────────────────────────────┘   └──────────────────────────────┘
                              │                                    │
                              ▼                                    ▼
                    Scoring Engine                         Scoring Engine
                    (component weights)                    (objective weights)
```

---

## 2. Entitas Utama

| Entitas | Deskripsi |
|---|---|
| PerformancePeriod | Periode evaluasi (nama, tanggal mulai-akhir, status) |
| PerformanceRating | Master rating (nama, range, bobot) |
| PerformancePerspective | BSC perspective (Financial, Customer, Internal Process, Learning & Growth) |
| PerformanceIndicatorFormula | Rumus kalkulasi indikator |
| PerformanceComponent | Komponen scoring (KPI, Program, Subordinate — 3 seeded, fixed) |
| PerformanceOrganizationComponent | Bobot komponen per organisasi |
| PerformanceTemplate | Template KPI (BSC, dengan perspectives & indicators) |
| PerformanceIndicator | Indikator KPI dalam template |
| PerformanceEvaluation | Evaluasi KPI per employee per period |
| PerformanceEvaluationDetail | Detail item evaluasi (indicator + target + actual + score) |
| PerformanceEvaluationProgramItem | Program item yang dibuat employee sendiri |
| PerformanceTarget | Target per evaluation detail |
| PerformanceProgress | Progress update per detail |
| PerformanceComment | Komentar pada evaluasi |
| PerformanceAttachment | Lampiran pada evaluasi |
| PerformanceLog | Audit trail perubahan evaluasi |
| OKRTemplate | Template OKR |
| OKRObjective | Objective dalam template OKR |
| OKRKeyResult | Key result per objective |
| OKREvaluation | Evaluasi OKR per employee |
| OKREvaluationDetail | Detail evaluasi OKR (key result + actual + score) |
| OKRProgress | Progress update OKR |
| OKRComment | Komentar OKR |
| OKRAttachment | Lampiran OKR |

---

## 3. SETUP — Master Data

| Menu | Endpoint | Deskripsi |
|---|---|---|
| Periods | `GET/POST /performance/periods` | Periode evaluasi (nama, start/end date) |
| Ratings | `GET/POST /performance/ratings` | Master rating (nama, range, bobot) |
| Perspectives | `GET/POST /performance/kpi/perspectives` | BSC perspectives |
| Formulas | `GET/POST /performance/indicator-formulas` | Rumus kalkulasi indikator |
| Components | `GET /performance/kpi/components` | Komponen scoring (read-only, 3 seeded) |
| Org Components | `POST /performance/kpi/organization-components` | Bobot komponen per organisasi |

---

## 4. KPI — Template & Indicators

1. **Buat template** — `POST /performance/kpi/templates`: nama, period_id, description, org scope.
2. **Duplikat template** — `POST /performance/kpi/templates/:id/duplicate`.
3. **Tambah indikator** — `POST /performance/kpi/indicators`: template_id, perspective_id, indicator name, formula, target.
4. **Scope organisasi** — `GET /performance/kpi/templates/organization-scope`: lihat template yang berlaku untuk organisasi tertentu.

---

## 5. KPI — Evaluation Workflow (2-Phase)

### Phase 1: Target Setting

1. **Buat evaluasi** — `POST /performance/kpi/evaluations`: employee_id, period_id, template_id → status `DRAFT`.
2. **Tambah detail** — `POST /performance/kpi/evaluation-details`: evaluation_id, indicator_id, target value.
3. **Input target** — `PUT /performance/kpi/evaluation-details/:id/target`.
4. **Submit target** — `POST /performance/kpi/evaluations/:id/submit-target` → `TARGET_SUBMITTED`.
5. **Approve target** — `POST /performance/kpi/evaluations/:id/approve-target` → `TARGET_APPROVED`.
6. **Reject target** — `POST /performance/kpi/evaluations/:id/reject-target` → `DRAFT` (kembali).

### Phase 2: Realization Assessment

1. **Input actual** — `PUT /performance/kpi/evaluation-details/:id/actual` (hanya bisa setelah `TARGET_APPROVED`).
2. **Bulk actuals** — `PUT /performance/kpi/evaluations/:id/actuals`.
3. **Program items** — `POST /performance/kpi/program-items` (employee-authored, bukan dari template).
4. **Submit** — `POST /performance/kpi/evaluations/:id/submit` → `SUBMITTED`.
5. **Approve** — `POST /performance/kpi/evaluations/:id/approve` → `APPROVED`.
6. **Reject** — `POST /performance/kpi/evaluations/:id/reject` → `DRAFT`.
7. **Complete** — `POST /performance/kpi/evaluations/:id/complete` → `COMPLETED`.

---

## 6. KPI — Scoring Engine

1. **Recalculate** — `POST /performance/kpi/evaluations/:id/recalculate`: hitung skor berdasarkan bobot komponen (KPI, Program, Subordinate).
2. **Components** — `GET /performance/kpi/evaluations/:id/components`: lihat skor per komponen.
3. **Update component score** — `PUT /performance/kpi/evaluations/:id/components/:component_id`.
4. **Period recalculate** — `POST /performance/kpi/periods/:period_id/recalculate-scoring`: recalculate semua evaluasi dalam periode.
5. **Progress summary** — `GET /performance/kpi/evaluations/:id/progress-summary`.

---

## 7. OKR — Template & Evaluation

1. **Buat template** — `POST /performance/okr/templates`: nama, deskripsi.
2. **Tambah objective** — `POST /performance/okr/objectives`: template_id, name, description.
3. **Tambah key result** — `POST /performance/okr/key-results`: objective_id, name, target value.
4. **Buat evaluasi** — `POST /performance/okr/evaluations`: employee_id, period_id, template_id.
5. **Submit** — `POST /performance/okr/evaluations/:id/submit` → `SUBMITTED`.
6. **Approve** — `POST /performance/okr/evaluations/:id/approve` → `APPROVED`.
7. **Complete** — `POST /performance/okr/evaluations/:id/complete` → `COMPLETED`.

---

## 8. KPI & OKR — Progress & Comment

- **Progress** — `POST /performance/kpi/progress`: detail_id, progress value, notes.
- **Comments** — `POST /performance/kpi/comments`: evaluation_id, comment text.
- **Attachments** — `POST /performance/kpi/attachments`: evaluation_id, file URL.
- **Logs** — `GET /performance/logs`: audit trail perubahan evaluasi.

---

## 9. KPI & OKR — Dashboard

| Endpoint | Deskripsi |
|---|---|
| `GET /performance/kpi/dashboard/employee/:employee_id` | Dashboard karyawan (evaluasi saya, progress) |
| `GET /performance/kpi/dashboard/manager/:manager_id` | Dashboard manager (bawahan saya, approval pending) |
| `GET /performance/kpi/dashboard/hr` | Dashboard HR (ringkasan semua evaluasi) |
| `GET /performance/okr/dashboard/hr` | Dashboard HR OKR |

---

## 10. Ringkasan Status

| Entitas | Status |
|---|---|
| KPI Evaluation | `DRAFT → TARGET_SUBMITTED → TARGET_APPROVED → SUBMITTED → APPROVED → COMPLETED` |
| OKR Evaluation | `DRAFT → SUBMITTED → APPROVED → COMPLETED` |
| Period | (custom status per kebutuhan) |

---

## 11. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Employee** | Target evaluasi per karyawan |
| **Organization** | Scope template & org component weight |
| **Job Management** | Job title/position reference |
| **Competency** | Competency score sebagai input |
| **Employee Movement** | Performance score dibaca untuk eligibility |
| **Settings** | Grading, employment status reference |

---

## 12. Peta Halaman UI

| Menu | Halaman |
|---|---|
| Performance (hub) | `PerformanceIndex.vue` / `Performance.vue` |
| KPI Hub | `kpi/KPIIndex.vue` |
| KPI Templates | `kpi/KPITemplates.vue` / `kpi/KPITemplateForm.vue` |
| KPI Evaluation Detail | `kpi/KPIEvaluationDetail.vue` |
| KPI Self-Assessment | `kpi/KPISelfAssessment.vue` |
| Periods | `settings/PerformancePeriodsView.vue` |
| OKR Hub | `okr/OKRIndex.vue` |
| OKR Templates | `okr/OKRTemplates.vue` / `okr/OKRTemplateForm.vue` |
| OKR Evaluation Detail | `okr/OKREvaluationDetail.vue` |
| OKR Self-Assessment | `okr/OKRSelfAssessment.vue` |

---

## 13. Endpoint API Utama

### Shared Master Data (`/api/v1/tenant/performance/`)

| Area | Endpoint |
|---|---|
| Periods | `GET/POST /periods`, `GET/PUT/DELETE /periods/:id` |
| Ratings | `GET/POST /ratings`, `GET/PUT/DELETE /ratings/:id` |
| Formulas | `GET/POST /indicator-formulas`, `GET/PUT/DELETE .../:id` |
| Logs | `GET /logs`, `GET /logs/:id` |

### KPI (`/api/v1/tenant/performance/kpi/`)

| Area | Endpoint |
|---|---|
| Perspectives | `GET/POST /perspectives`, `GET/PUT/DELETE .../:id` |
| Templates | `GET/POST /templates`, `GET/PUT/DELETE .../:id`, `POST .../:id/duplicate`, `GET /templates/organization-scope` |
| Indicators | `GET/POST /indicators`, `GET/PUT/DELETE .../:id` |
| Evaluations | `GET/POST /evaluations`, `GET/PUT/DELETE .../:id`, `PUT .../:id/status` |
| Snapshot | `POST /evaluations/snapshot`, `GET /evaluations/:id/full` |
| Details | `POST /evaluation-details`, `GET /evaluations/:id/details`, `PUT/DELETE .../:id` |
| Target/Actual | `PUT /evaluation-details/:id/target`, `PUT .../:id/actual`, `PUT /evaluations/:id/actuals` |
| Program Items | `POST /program-items`, `GET /evaluations/:id/program-items`, `PUT/DELETE .../:id` |
| Progress | `POST /progress`, `GET/PUT/DELETE .../:id` |
| Comments | `POST /comments`, `GET/PUT/DELETE .../:id` |
| Attachments | `POST /attachments`, `GET/PUT/DELETE .../:id` |
| Scoring | `POST /evaluations/:id/recalculate`, `GET .../:id/progress-summary` |
| Workflow | `POST .../submit-target`, `POST .../approve-target`, `POST .../reject-target`, `POST .../submit`, `POST .../approve`, `POST .../reject`, `POST .../complete` |
| Components | `GET /components`, `GET/PUT .../:id`, `POST /organization-components`, `GET/DELETE .../:id` |
| Eval Components | `GET /evaluations/:id/components`, `POST .../calculate-scoring`, `PUT .../:component_id` |
| Period Scoring | `POST /periods/:period_id/recalculate-scoring` |
| Dashboard | `GET /dashboard/employee/:employee_id`, `GET /dashboard/manager/:manager_id`, `GET /dashboard/hr` |
| Context | `GET /my-context` |

### OKR (`/api/v1/tenant/performance/okr/`)

| Area | Endpoint |
|---|---|
| Templates | `GET/POST /templates`, `GET/PUT/DELETE .../:id`, `POST .../:id/duplicate`, `GET .../:id/objectives` |
| Objectives | `GET/POST /objectives`, `GET/PUT/DELETE .../:id` |
| Key Results | `GET/POST /key-results`, `GET/PUT/DELETE .../:id` |
| Evaluations | `GET/POST /evaluations`, `GET/PUT/DELETE .../:id`, `GET .../:id/details` |
| KR Proposal | `POST /evaluations/:id/key-results`, `POST .../submit-key-results`, `POST .../approve-key-results`, `POST .../reject-key-results` |
| Workflow | `POST /evaluations/:id/submit`, `POST .../approve`, `POST .../reject`, `POST .../complete` |
| Details | `PUT /evaluation-details/:id`, `GET .../:id/progress`, `GET .../:id/attachments` |
| Progress | `GET/POST /progress`, `GET/PUT/DELETE .../:id` |
| Comments | `GET /evaluations/:id/comments`, `POST /comments`, `PUT/DELETE .../:id` |
| Attachments | `POST /attachments`, `DELETE .../:id` |
| Dashboard | `GET /dashboard/hr` |
| Context | `GET /my-context` |

---

## 14. Catatan Penting

- **Components** (KPI/Program/Subordinate) di-seed 3 baris, tidak bisa ditambah/dihapus (hanya update nama/bobot/active).
- **KPI Evaluations** punya **2-phase workflow**: target setting dulu, baru realization assessment.
- **OKR Evaluations** workflow lebih sederhana: submit → approve → complete.
- **Program items** dibuat employee sendiri (bukan dari template HR).
- **Scope organisasi** menentukan template mana yang berlaku untuk karyawan tertentu.
- **Scoring engine** menghitung skor akhir berdasarkan bobot komponen per organisasi.
