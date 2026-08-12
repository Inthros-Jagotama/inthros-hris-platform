# Recruitment — Strategic Layer Integration Plan

> 📅 Versi plan: 2026-08-12 · Status: **PROPOSAL — Strategic Layer (Workforce Intelligence / Career Intelligence) ↔ Recruitment** — bagian ini **dipisah dari `module-recruitment-development-plan.md`** (keputusan scoping 2026-08-12: Recruitment = module **operasional**; layer strategis dikelola di plan ini).
> ✅ **Fakta aktual (audit 2026-08-12):** kedua module strategis sudah terimplementasi penuh di sisi backend — **Workforce Intelligence** (gap analysis supply vs demand, projections/hiring needs, `analytics/recruitment`, `candidate-search`) dan **Career Intelligence** (talent maps, career interests, career paths, gap analysis, succession plans). Yang **belum ada** di sisi strategis: integrasi dua arah terstruktur dengan Recruitment (workforce gap → requisition, expected hires → remaining gap, internal candidate via career path, succession fallback, Quality of Hire).
> 🔎 **Sumber:** `docs/module-recruitment-development-plan.md` §5.2 (Out of Scope) + `docs/workforce-intelligence-training-enhancement-plan.md` + `docs/module-career-intelligence-plan.md` + audit `backend/internal/modules/workforceintelligence/` (routes.go: `analytics/recruitment` L42, `candidate-search` L10, people-analytics L95-101) + `backend/internal/modules/careerintelligence/` + `docs/project-completion-dashboard.md`.
> 📊 **Progres per 2026-08-12:** ✅ Workforce Intelligence backend (7 entity, 68 endpoint, 108 test — gap analysis, projections, candidate-search, people-analytics) · ✅ Career Intelligence backend (5 sub-module, 21 endpoint: talent maps, interests, career paths, gap analysis, succession) · 🔶 Workforce Intelligence sebagian mengonsumsi data recruitment via `analytics/recruitment` & `candidate-search` (konsumsi sepihak, tanpa flow balik) · ⏳ Integrasi dua arah dengan Recruitment belum ada.
> ⏳ **Sisa TODO:** Gap §6 — tersisa S-5, S-6, S-7 (S-1 workforce gap → requisition ✅, S-2 expected hires → remaining gap ✅, S-3 candidate search & recruitment analytics ✅, S-4 internal candidate via career path ✅ — 12 Agu 2026); prioritas P0 berikutnya: S-5 succession fallback.
> 🔧 **Catatan konsistensi docs:** plan ini adalah **"rumah" bagi seluruh item strategis** yang dipisah dari `module-recruitment-development-plan.md` §5.2. Jangan kembalikan item di bawah ini ke plan Recruitment — Recruitment hanya menyediakan data operasional.

---

# 1. Objective

Menjembatani **layer strategis** (Workforce Intelligence & Career Intelligence) dengan **Recruitment (operasional)**:

```text
STRATEGIC LAYER (plan ini)
├── Workforce Intelligence
│   ├── Workforce Gap (supply vs demand)
│   ├── Hiring Need / Expected Hires
│   ├── Remaining Workforce Gap
│   └── Quality of Hire (probation + performance + retention)
│
├── Career Intelligence
│   ├── Career Path → Internal Candidate Eligibility
│   ├── Talent Map / Interest
│   └── Succession Plan → fallback external recruitment
│
└── Training (handoff pasca-onboarding)
        │
        ▼
OPERATIONAL LAYER (module-recruitment-development-plan.md)
├── Job Requisition
├── Recruitment Pipeline (candidate → application → interview → offer)
├── Employee / Employee Movement
└── Onboarding
```

Prinsip pembagian responsibility (keputusan scoping 2026-08-12):

```text
Workforce Intelligence → menentukan kebutuhan workforce & hiring need (strategis)
Career Intelligence    → menganalisis career, talent, gap, succession (strategis)
Recruitment            → mencari dan memilih kandidat (operasional — plan terpisah)
```

---

# 2. Boundary — Strategic vs Operational

| Domain | Pemilik | Module/Plan |
|---|---|---|
| Workforce forecasting, workforce gap, hiring need | Workforce Intelligence | plan ini §4.1 (S-1) + `workforce-intelligence-training-enhancement-plan.md` |
| Expected hires → remaining workforce gap | Workforce Intelligence | plan ini §4.2 (S-2) |
| Quality of Hire (probation + performance + retention) | Workforce / Career Intelligence | plan ini §4.6 (S-6) |
| Career path, talent map, interest, succession, eligibility internal candidate | Career Intelligence | plan ini §4.4-§4.5 (S-4/S-5) + `module-career-intelligence-plan.md` |
| Pelaksanaan training & development | Module Training | plan ini §4.7 (S-7) |
| **Requisition, candidate, application, interview, offer, onboarding** | **Recruitment (operasional)** | **`module-recruitment-development-plan.md`** |

> 🚫 Recruitment **tidak memiliki logika strategis** — ia hanya menyediakan **data operasional** (requisition, pipeline, offer, hire, onboarding) yang dikonsumsi layer strategis.

---

# 3. Status Aktual — Module Strategis

## 3.1 Workforce Intelligence — ✅ BACKEND LENGKAP (25 Jul 2026)

Modul `backend/internal/modules/workforceintelligence/` (7 entity, **68 endpoint, 108 test**):

- **Gap analysis (supply vs demand)** — `analytics/workforce-gap` → `SURPLUS / SHORTAGE / OPTIMAL` per department & overall.
- **Projections** — headcount proyeksi, `hiring_needed`, retirement counts, growth rate.
- **Recruitment analytics** — `GET /workforce-intelligence/analytics/recruitment` (konsumsi sepihak data recruitment).
- **Candidate search** — `GET /workforce-intelligence/candidate-search` (posisi kosong + kandidat recruitment).
- **People analytics** — `source-vs-retention`, `career-progression`, `training-vs-performance`, dsb. (dasar Quality of Hire).
- Permissions: `workforceintelligence.view/create/update/delete` (seeded `tenantseed/seed_rbac.go`).

## 3.2 Career Intelligence — ✅ BACKEND LENGKAP (26 Jul 2026)

Modul `backend/internal/modules/careerintelligence/` (5 sub-module, 21 endpoint):

- Talent Maps (9-box grid) · Career Interests · Career Paths (ladder) · Gap Analysis (`/paths/gap-analysis`) · Succession Plans.
- Permissions: 10 permission granular (`career-intelligence.*`, lihat `module-career-intelligence-plan.md` §2.3).

## 3.3 Integrasi strategis ↔ Recruitment — ⏳ BELUM

| Integrasi | Status | Catatan |
|---|---|---|
| Workforce gap / hiring need → requisition | ⏳ Belum | `job_requisitions` belum punya `workforce_gap_id` / `workforce_plan_id` (dipisah dari plan Recruitment, G-2) |
| Expected hires → remaining workforce gap | ⏳ Belum | tidak ada alur balik dari recruitment (offers accepted) ke WI |
| Candidate search (WI membaca candidates) | 🔶 Sebagian | sudah ada endpoint, konsumsi sepihak |
| Internal candidate via career path | ⏳ Belum | tidak ada `Invite to Apply` / internal application dari CI |
| Succession fallback → external recruitment | ⏳ Belum | succession plan tidak menautkan ke requisition |
| Quality of Hire | ⏳ Belum | data tersedia (source analytics), metrik agregat belum dibangun |
| Onboarding → training handoff | ⏳ Belum | kebutuhan training tidak diteruskan ke Training module |

---

# 4. Gap Analysis & Enhancement Plan

> ⚠️ Seluruh gap di bawah adalah **rencana** (belum dieksekusi per 2026-08-12) dan **dimiliki layer strategis** — implementasi di sisi Recruitment mengikuti plan Recruitment bila endpoint/kolom baru dibutuhkan.

## S-1 🔴 WORKFORCE GAP → REQUISITION (WI → Recruitment)

**Status: ✅ Selesai (12 Agu 2026).** WI sudah menghitung gap (supply vs demand) & `hiring_needed`; requisition kini dapat menautkan ke gap tersebut.

**Implementasi (eksekusi kolaboratif):**
- Migration `091_recruitment_workforce_gap` (mysql + postgres + down): `reason_type` VARCHAR(30) + `workforce_gap_id` CHAR(36) nullable + `workforce_plan_id` CHAR(36) nullable pada `job_requisitions`.
- Flow: `Workforce Gap (Required 6 − Current 4 = Gap 2)` → **WI menghasilkan hiring need** → **Recruitment membuat requisition** dengan `slots_available = hiring need`, `reason_type = WORKFORCE_GAP`.
- Recruitment **membaca** hiring need dari WI via interface narrow `recruitment.WorkforceGapProvider` (`HiringGapForOrganization`), di-wire di `cmd/server/main.go` melalui `workforceGapAdapter{wiSvc}` — Recruitment tidak menghitung gap sendiri.
- Auto-resolve: saat `reason_type = WORKFORCE_GAP` dan `slots_available` tidak ditentukan eksplisit, slots diisi dari hiring need WI (fail-safe: provider nil/error → slots default).
- `reason_type` pada requisition mencakup `WORKFORCE_GAP` (+ `NEW_POSITION`, `REPLACEMENT`, `EXPANSION`).
- OpenAPI: `CreateRequisitionRequest`, `UpdateRequisitionRequest`, `RequisitionResponse` + `reason_type`/`workforce_gap_id`/`workforce_plan_id`.
- Test: +7 unit test service (auto-resolve, provider nil fallback, provider error, explicit slots menang, update reason, unrelated-update tidak menimpa slots, clear gap id).

**Ref:** plan asli recruitment §7-§9 · `workforce-intelligence-training-enhancement-plan.md`.

## S-2 🔴 EXPECTED HIRES → REMAINING WORKFORCE GAP (Recruitment → WI)

**Status: ✅ Selesai (12 Agu 2026).** WI `analytics/recruitment` kini menghitung sisa gap berbasis pipeline.

**Implementasi (kepemilikan WI):**
- WI mengonsumsi `Open Positions, Recruitment Pipeline, Expected Hires (Accepted Offers), Filled Positions` dari `job_requisitions` + `job_applications` (repository read-only, tanpa import modul Recruitment).
- Formula: `Required Workforce − Current Workforce − Expected Hires = Remaining Workforce Gap` (contoh: Required 100, Current 90, Accepted Offer 2 → Remaining Gap 8).
- `GetGapAnalysis` diperkaya: `expected_hires`, `open_positions`, `filled_positions`, `remaining_gap` (shortage − expected hires, min 0).
- `GetRecruitmentAnalytics` diimplementasikan dari repo (sebelumnya placeholder): open positions, expected hires, filled positions, pipeline/funnel per status, remaining gap.
- Handler `GET /workforce-intelligence/analytics/recruitment` kini memanggil service.
- Test: +8 (4 repository + 4 service; total WI 112 → 120). OpenAPI: `GapAnalysisResponse` + `RecruitmentAnalytics` diperkaya.

**Ref:** plan asli recruitment §37-§39.

## S-3 🟢 CANDIDATE SEARCH & RECRUITMENT ANALYTICS (WI konsumsi)

**Status: ✅ SELESAI (12 Agu 2026).**

**Implementasi:**
- `GET /workforce-intelligence/candidate-search` — enhancement: filter **posisi** (`?position` — nama/kode organisasi lowong) + integrasi **internal candidate eligible** (Career Intelligence via interface narrow `InternalEligibilityProvider` + adapter `wiInternalCandidateAdapter` di main.go): per posisi lowong ditampilkan `internal_candidate_count` + `internal_candidates` (employee internal yang eligible via career path, di-dedupe per employee). Provider nil → kosong (fail-safe). **Filter kompetensi ditunda** — kandidat eksternal belum punya data kompetensi (G-9 belum diimplementasikan).
- `GET /workforce-intelligence/analytics/recruitment` — metrik **nyata**: `Time to Hire` (avg accepted_at−applied_at), `Time to Fill` (avg closed_at−created_at requisition FILLED), `Offer Acceptance Rate` (accepted/offered), `Source Conversion` (kandidat→hire per `candidates.source`) + `by_source`; `Candidate Match Score` & `Cost per Hire` tetap placeholder (data kompetensi kandidat G-9 & biaya belum dikumpulkan).
- Tanpa perubahan di sisi Recruitment — konsumsi data tetap sepihak (WI membaca tabel operasional + CI via provider narrow). +12 test (WI 120 → 132); OpenAPI disinkronkan (fix dokumentasi candidate-search yang stale).

**Ref:** plan asli recruitment §37-§39 · `module-recruitment-development-plan.md` G-11.

## S-4 🟡 INTERNAL CANDIDATE VIA CAREER PATH (CI → Recruitment)

**Status: ✅ Selesai (12 Agu 2026).**

**Implementasi (kepemilikan Career Intelligence):**
- CI menyediakan **eligibility internal candidate**: `Position Vacancy → Career Path → Eligible Employees` — employee dengan employment aktif pada source step (semua step sebelum target) dari career path menuju posisi lowongan.
- Endpoint CI: `GET /career-intelligence/paths/{id}/eligible-employees` (service `GetEligibleEmployeesForPath`; target = step terakhir).
- Interface narrow `recruitment.InternalCandidateProvider` (`EligibleEmployeesForPosition`) + adapter `internalCandidateAdapter{ciSvc}` di `cmd/server/main.go` (pola sama S-1).
- Endpoint Recruitment: `GET /recruitment/eligible-internal-candidates?position_id=...` — recruiter membaca eligible employee dari CI; Recruitment **tidak menghitung career eligibility** (fail-safe: provider nil → list kosong).
- Repo CI: `FindCareerPathsByTarget` (path menuju posisi, step terakhir) + `ListEligibleEmployeesByPositionIDs` (employment aktif, non soft-delete).
- Test: +12 (CI 4 repo + 5 service; recruitment 3 service). OpenAPI: 2 endpoint + 2 schema.

**Ref:** plan asli recruitment §34-§36 · `module-career-intelligence-plan.md` §4.

## S-5 🟡 SUCCESSION PLANNING → FALLBACK EXTERNAL RECRUITMENT (CI)

**Status: ⏳ Belum.**

**Rencana (kepemilikan Career Intelligence):**
- `Succession Plan → Successor Candidates` (internal); bila successor internal **tidak tersedia** → `Succession Gap → External Recruitment` (recruitment jadi fallback).
- CI menandai posisi kunci tanpa successor siap; WI/CI menautkan ke kebutuhan requisition eksternal (kolaborasi dengan S-1).

**Ref:** plan asli recruitment §34-§36 · `module-career-intelligence-plan.md` §5.5.

## S-6 🟡 QUALITY OF HIRE (WI/CI metrik agregat)

**Status: ⏳ Belum.** Data tersedia (recruitment scores + probation + performance + retention); metrik agregat belum dibangun.

**Rencana (kepemilikan Workforce Intelligence):**
- `Quality of Hire = Recruitment Match Score + Interview Score + Assessment Score + Probation Result + Performance + Retention`.
- Breakdown: `Hiring Source, Requisition, Organization, Position, Period`.
- Output: korelasi source-vs-retention (`people-analytics/source-vs-retention` sudah ada sebagai dasar) diperluas menjadi metrik Quality of Hire.

**Ref:** plan asli recruitment §55 · `workforce-intelligence-training-enhancement-plan.md`.

## S-7 🟢 ONBOARDING → TRAINING HANDOFF

**Status: ⏳ Belum.**

**Rencana (handoff, bukan eksekusi):**
- Setelah onboarding selesai, kebutuhan training/development diteruskan ke **Training module** (`Employee → Position → Required Training`).
- Recruitment/Onboarding **tidak mengeksekusi training** — hanya menghasilkan kebutuhan (prinsip plan Recruitment §4.4).
- Kepemilikan penuh ada di Training module; WI dapat menganalisis hasilnya (training demand/gap — `workforce-intelligence-training-enhancement-plan.md`).

**Ref:** plan asli recruitment §32 · `workforce-intelligence-training-enhancement-plan.md`.

---

# 5. API Plan

## 5.1 Existing (module strategis — sudah ada)

```http
## Workforce Intelligence
GET    /api/v1/tenant/workforce-intelligence/analytics/workforce-gap
GET    /api/v1/tenant/workforce-intelligence/analytics/projections
GET    /api/v1/tenant/workforce-intelligence/analytics/recruitment
GET    /api/v1/tenant/workforce-intelligence/candidate-search
GET    /api/v1/tenant/workforce-intelligence/people-analytics/source-vs-retention

## Career Intelligence
GET    /api/v1/tenant/career-intelligence/paths
GET    /api/v1/tenant/career-intelligence/paths/gap-analysis
GET    /api/v1/tenant/career-intelligence/successions
GET    /api/v1/tenant/career-intelligence/talent-maps
```

## 5.2 Target tambahan (rencana — per Gap §4)

```http
## Workforce Intelligence (kepemilikan WI)
GET    /workforce-intelligence/hiring-plan                       ← S-1/S-2 (hiring need + remaining gap)
GET    /workforce-intelligence/analytics/quality-of-hire         ← S-6
GET    /workforce-intelligence/analytics/recruitment?metric=time-to-hire ← S-3

## Career Intelligence (kepemilikan CI)
GET    /career-intelligence/paths/{id}/eligible-employees        ← S-4 (internal candidate eligibility)
POST   /career-intelligence/successions/{id}/external-recruit    ← S-5 (fallback → kebutuhan requisition)

## Kolaboratif (butuh endpoint kecil di sisi Recruitment)
POST   /recruitment/requisitions/{id}/link-gap                   ← S-1 (tautkan workforce_gap_id)
GET    /recruitment/offers?status=accepted                       ← S-2 (konsumsi accepted offers oleh WI)
```

---

# 6. Permissions & Authorization

- Tidak ada permission baru yang wajib di modul strategis untuk enhancement ini — reuse `workforceintelligence.view` & `career-intelligence.*` yang sudah ada.
- Akses data recruitment oleh WI/CI harus melalui **service/interface narrow** (bukan query langsung lintas modul) — pola `employeemovement.CareerExecutor`.
- Role: `Company Admin` (full), `Manager` (view/create/update), `Employee` (view) — sudah ter-seed di `authz/rbac.go` & `tenantseed/seed_rbac.go`.

---

# 7. Data Model

Tidak membuat tabel strategis baru yang duplikatif — gunakan:

```text
Operational (Recruitment)                Strategic (baca)
├── job_requisitions (+ workforce_gap_id/plan_id setelah S-1)
├── job_applications / candidates        → WI analytics/recruitment & candidate-search
├── interviews / job_offers              → Expected Hires (accepted offers)
└── employee_onboardings                 → Training handoff
```

Projection (bila diperlukan) di sisi WI, mengikuti pola `workforce-intelligence-training-enhancement-plan.md` §19:

```text
workforce_hiring_plan_summary
workforce_hiring_gaps
workforce_quality_of_hire
```

---

# 8. Development Priority

## P0 — Integrasi strategis inti
1. ~~S-1 Workforce gap → requisition (WI → Recruitment)~~ ✅ 12 Agu 2026
2. ~~S-2 Expected hires → remaining workforce gap (Recruitment → WI)~~ ✅ 12 Agu 2026
3. ~~S-4 Internal candidate via career path (CI → Recruitment)~~ ✅ 12 Agu 2026
4. S-5 Succession fallback → external recruitment

## P1 — Metrik & analitik
5. S-6 Quality of Hire
6. ~~S-3 Enhancement candidate-search & recruitment analytics~~ ✅ 12 Agu 2026

## P2 — Handoff & polish
7. S-7 Onboarding → training handoff
8. Dashboard "Hiring Plan vs Actual" (WI)

---

# 9. Recommended Implementation Order

```text
STEP 1  ✅ S-1  Workforce gap → requisition (interface narrow + migration workforce_gap_id) — SELESAI 12 Agu 2026
STEP 2  ✅ S-2  Expected hires → remaining gap (WI konsumsi accepted offers) — SELESAI 12 Agu 2026
STEP 3  ✅ S-4  Internal candidate via career path (eligible-employees + invite to apply) — SELESAI 12 Agu 2026
STEP 4  ✅ S-3  Enhancement candidate-search (filter posisi + internal candidate eligible) & recruitment analytics (Time to Hire/Fill, Offer Acceptance, Source Conversion) — SELESAI 12 Agu 2026
STEP 5  S-5  Succession fallback → external recruitment
STEP 6  S-6  Quality of Hire
STEP 7  S-7  Onboarding → training handoff
STEP 8  Testing & E2E (strategic layer)
```

---

# 10. Definition of Done

- [x] Workforce Intelligence backend (gap analysis, projections, candidate-search, recruitment analytics).
- [x] Career Intelligence backend (talent maps, interests, career paths, gap analysis, succession plans).
- [x] Requisition dapat menautkan `workforce_gap_id` (S-1). ✅ 12 Agu 2026
- [x] WI menghitung `Remaining Workforce Gap` dari expected hires (S-2). ✅ 12 Agu 2026
- [x] CI menyediakan daftar employee eligible untuk internal application (S-4). ✅ 12 Agu 2026
- [x] Candidate search mendukung filter posisi + internal candidate eligible; recruitment analytics menghitung Time to Hire/Fill, Offer Acceptance Rate, Source Conversion (S-3). ✅ 12 Agu 2026
- [ ] Succession gap menghasilkan kebutuhan external recruitment (S-5).
- [ ] Metrik Quality of Hire tersedia (S-6).
- [ ] Kebutuhan training pasca-onboarding diteruskan ke Training module (S-7).
- [ ] Tidak ada logika strategis yang bocor ke Recruitment (boundary terjaga).
- [ ] Integration test: Workforce Gap → Requisition · Accepted Offer → Remaining Gap · Career Path → Internal Application · Succession → External Recruitment.

---

# 11. Kesimpulan

Plan ini adalah **rumah bagi layer strategis** yang dipisah dari `module-recruitment-development-plan.md` (keputusan scoping 2026-08-12):

```text
Workforce Intelligence → menentukan kebutuhan workforce & hiring need (strategis)
Recruitment            → mencari dan memilih kandidat (operasional — plan terpisah)
Career Intelligence    → menganalisis career, talent, gap, succession (strategis)
```

Backend kedua module strategis sudah lengkap; fokus berikutnya adalah **integrasi dua arah** dengan Recruitment (S-1 → S-7) — lihat Gap §4 dan urutan implementasi §9.
