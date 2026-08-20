# Alur Pengisian Job Management (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Job Management** — setup master
data (job titles, values, clusters), pengisian analisis jabatan (objectives, identifications,
responsibilities, authorities, activities, risks, relationships, financials), kompetensi
(potency, scoring, competency groups), dan dashboard — pola runbook seperti
[`module-attendance-flow.md`](module-attendance-flow.md) & [`module-training-flow.md`](module-training-flow.md).

- Modul core: **Job Management** (`jobmanagement`) — depends on `organization`
- Lokasi kode: `backend/internal/modules/jobmanagement/` · `frontend/tenant/src/views/modules/job/` & `jobvalues/`
- Permission: `jobmanagement.view`, `jobmanagement.create`, `jobmanagement.update`, `jobmanagement.delete`
- Sidebar: **Core HR → Job Management** (dropdown: Job Values Mapping, Job Management)

---

## 1. Ringkasan Alur End-to-End

```
SETUP (sekali)                  DATA JABATAN (per org)              KOMPETENSI & SKOR
┌──────────────────┐   ┌──────────────────────────────────┐   ┌──────────────────────────┐
│ Job Titles (9.1) │   │ Objectives (9.4)                 │   │ Potency Competencies     │
│  └─ Title Sub    │──▶│ Identifications (9.5)            │   │   (9.16)                 │
│ Job Values (9.3) │   │ Responsibilities (9.6)           │──▶│ Competency Groups (9.18) │
│  └─ Clusters     │   │ Education Experiences (9.7)      │   │ Job Scores (9.17)        │
└──────────────────┘   │ HR Authorities (9.8)             │   │   (scoring per org)      │
                       │ Operational Authorities (9.9)    │   └──────────────────────────┘
                       │ Working Activities (9.10)         │
                       │ Working Risks (9.11)              │
                       │ Relationships (9.12 + Details)    │
                       │ Subordinate Controls (9.13)       │
                       │ Assets (9.14)                     │
                       │ Financials (9.15)                  │
                       └──────────────────────────────────┘
```

- **Tidak ada status state-machine** — semua data jabatan bersifat CRUD biasa (create/update/delete).
- **Scope per organisasi** — mayoritas entitas di-scoping ke `organization_id` (posisi jabatan tertentu).
- **Job Values** adalah master dropdown values yang dipakai sebagai referensi di banyak entitas (tipe: environment, hazard, asset, cash, authority, impact, frequency, relationship, education, experience, dll.).
- **Job Titles** adalah master nama jabatan dengan sub-level (title → sub-title → values).

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Job Title | `job_management_titles` | Master nama jabatan (nama, deskripsi, status aktif) |
| Job Title Sub | `job_management_title_subs` | Sub-jabatan di bawah title (hierarchical) |
| Job Value | `job_management_values` | Master values untuk dropdown (tipe: environment, hazard, education, experience, dll.) |
| Value Cluster | `job_management_value_clusters` | Mapping type ↔ cluster kompetensi |
| Job Objective | `job_management_objectives` | Tujuan/jabatan per posisi organisasi |
| Job Identification | `job_management_identifications` | Identitas jabatan (grading, nomenclature) |
| Job Responsibility | `job_management_responsibilities` | Tanggung jawab (main task, activities, outputs, success indicators) |
| Job Education Experience | `job_management_education_experiences` | Pendidikan & pengalaman yang dibutuhkan (relasi ke value education + experience) |
| Education Major (pivot) | `job_management_majors` | Pivot many-to-many ke master jurusan |
| Job Family (pivot) | `job_management_job_family` | Pivot many-to-many ke master bidang pekerjaan |
| Job HR Authority | `job_management_hr_authorities` | Kewenangan SDM per posisi |
| Job Operational Authority | `job_management_operational_authorities` | Kewenangan operasional per posisi |
| Job Working Activity | `job_management_working_activities` | Aktivitas kerja (relasi ke job value) |
| Job Working Risk | `job_management_working_risks` | Risiko kerja (environment + hazard dari job values) |
| Job Relationship | `job_management_relationships` | Hubungan kerja (type, frequency dari job values) |
| Relationship Detail | `job_management_relationship_details` | Rincian aktivitas per hubungan (org ref + activity) |
| Job Subordinate Control | `job_management_subordinate_controls` | Bawahan yang dikendalikan |
| Job Asset | `job_management_assets` | Aset jabatan (type + authority dari job values) |
| Job Financial | `job_management_financials` | Keuangan jabatan (cash, authority, impact dari job values + authorized flag) |
| Job Potency Competency | `job_management_potency_competencies` | Kompetensi potensi (value level + competency + weight) |
| Job Score | `job_management_scores` | Skor jabatan (agregat per organisasi: with/without financial, components JSON) |
| Job Competency Group | `job_management_competency_groups` | Bobot kompetensi per organisasi per kategori |

---

## 3. TAHAP 1 — SETUP Master Data (dikerjakan sekali)

### A. Job Titles (9.1 + 9.2)

Menu **Core HR → Job Management** → card "Job Management" → halaman utama.

- CRUD nama jabatan: nama, deskripsi, status aktif.
- **Expand** baris title untuk melihat & mengelola **Title Subs** (sub-jabatan):
  - Tambah sub: nama, deskripsi.
  - Setiap sub-title nantinya bisa dihubungkan ke Job Values.
- Endpoint: `GET/POST /titles`, `GET/PUT/DELETE /titles/:id`
- Sub: `GET/POST /titles/:id/subs`, `GET/PUT/DELETE /titles/:id/subs/:subId`

### B. Job Values (9.3)

Menu **Core HR → Job Management** → card "Job Values Mapping" → halaman `JobValuesIndex.vue`.

- **Master dropdown values** — berbagai tipe (type): `environment`, `hazard`, `asset`, `cash`,
  `authority`, `impact`, `frequency`, `relationship`, `education`, `experience`, dll.
- Setiap value punya: nama, tipe, type group, urutan (sort_order).
- **Tree view** tersedia: `GET /values/tree` — hierarchical view semua values per tipe.
- **Value Cluster mapping** — map type ↔ cluster kompetensi:
  - `GET /values/clusters/:type` — lihat cluster untuk type tertentu.
  - `PUT /values/clusters/:type` — update mapping cluster.
- Endpoint: `GET/POST /values`, `GET /values/tree`, `GET/PUT/DELETE /values/:id`
- Cluster: `GET /values/clusters/:type`, `PUT /values/clusters/:type`

---

## 4. TAHAP 2 — DATA JABATAN (per posisi organisasi)

Semua entitas di bawah ini di-scoping ke **organization_id** — artinya setiap entitas terkait
ke posisi jabatan tertentu dalam struktur organisasi. Data diisi berdasarkan **nomenclature**
(nama jabatan) dan **full_code** (kode jabatan) dari posisi yang dipilih.

### A. Job Objectives (9.4)

Tujuan/jabatan per posisi.

- Isi: organization_id, nomenclature, full_code, objective (text bebas).
- Endpoint: `GET/POST /objectives`, `GET/PUT/DELETE /objectives/:id`
- Query param: `?organization_id=...` untuk filter per posisi.

### B. Job Identifications (9.5)

Identitas jabatan — menghubungkan posisi dengan grading.

- Isi: organization_id, nomenclature, full_code, **grading_id** (dari master Grading di Settings).
- Endpoint: `GET/POST /identifications`, `GET/PUT/DELETE /identifications/:id`

### C. Job Responsibilities (9.6)

Tanggung jawab jabatan — deskripsi detail pekerjaan.

- Isi: organization_id, nomenclature, full_code:
  - **Main Task** — tugas utama.
  - **Activities** — aktivitas pendukung.
  - **Outputs** — keluaran yang diharapkan.
  - **Success Indicators** — indikator keberhasilan.
- Endpoint: `GET/POST /responsibilities`, `GET/PUT/DELETE /responsibilities/:id`

### D. Job Education Experiences (9.7)

Pendidikan & pengalaman yang dibutuhkan untuk posisi.

- Isi: organization_id, nomenclature, full_code:
  - **education_id** → relasi ke Job Value (type=education).
  - **experience_id** → relasi ke Job Value (type=experience).
  - **Majors** (pivot) → many-to-many ke master `education_majors` (di Settings).
  - **Job Families** (pivot) → many-to-many ke master `job_families` (di Settings).
- Endpoint: `GET/POST /education-experiences`, `GET/PUT/DELETE /education-experiences/:id`

### E. Job HR Authorities (9.8)

Kewenangan SDM per posisi.

- Isi: organization_id, nomenclature, full_code, description.
- Endpoint: `GET/POST /hr-authorities`, `GET/PUT/DELETE /hr-authorities/:id`

### F. Job Operational Authorities (9.9)

Kewenangan operasional per posisi.

- Isi: organization_id, nomenclature, full_code, description.
- Endpoint: `GET/POST /operational-authorities`, `GET/PUT/DELETE /operational-authorities/:id`

### G. Job Working Activities (9.10)

Aktivitas kerja — relasi ke Job Value untuk klasifikasi aktivitas.

- Isi: organization_id, nomenclature, full_code, **job_management_value_id** (relasi ke Job Value).
- Endpoint: `GET/POST /working-activities`, `GET/PUT/DELETE /working-activities/:id`

### H. Job Working Risks (9.11)

Risiko kerja — relasi ke Job Values untuk environment dan hazard.

- Isi: organization_id, nomenclature, full_code:
  - **job_management_value_environment_id** → relasi ke Job Value (type=environment).
  - **job_management_value_hazard_id** → relasi ke Job Value (type=hazard).
- Endpoint: `GET/POST /working-risks`, `GET/PUT/DELETE /working-risks/:id`

### I. Job Relationships (9.12 + 9.12b)

Hubungan kerja — relasi ke Job Values untuk relationship type dan frequency.

- Isi: organization_id, nomenclature, full_code:
  - **job_management_value_relationship_id** → relasi ke Job Value (type=relationship).
  - **job_management_value_frequency_id** → relasi ke Job Value (type=frequency).
- **Relationship Details** (sub-resource) — rincian aktivitas per hubungan:
  - organization_id (ref ke posisi lain yang terlibat), activity (deskripsi ruang lingkup).
  - Endpoint nested: `GET/POST /relationships/:id/details`, `GET/PUT/DELETE /relationships/:id/details/:detailId`
- Endpoint: `GET/POST /relationships`, `GET/PUT/DELETE /relationships/:id`

### J. Job Subordinate Controls (9.13)

Bawahan yang dikendalikan oleh posisi.

- Isi: organization_id, nomenclature, full_code, **job_management_value_id** (relasi ke Job Value).
- Endpoint: `GET/POST /subordinate-controls`, `GET/PUT/DELETE /subordinate-controls/:id`

### K. Job Assets (9.14)

Aset jabatan — relasi ke Job Values untuk asset type dan authority.

- Isi: organization_id, nomenclature, full_code:
  - **job_management_value_asset_id** → relasi ke Job Value (type=asset).
  - **job_management_value_authority_id** → relasi ke Job Value (type=authority).
- Endpoint: `GET/POST /assets`, `GET/PUT/DELETE /assets/:id`

### L. Job Financials (9.15)

Keuangan jabatan — relasi ke Job Values untuk cash, authority, dan impact.

- Isi: organization_id, nomenclature, full_code, **is_authorized** (boolean):
  - **job_management_value_cash_id** → relasi ke Job Value (type=cash).
  - **job_management_value_authority_id** → relasi ke Job Value (type=authority).
  - **job_management_value_impact_id** → relasi ke Job Value (type=impact).
- Endpoint: `GET/POST /financials`, `GET/PUT/DELETE /financials/:id`

---

## 5. TAHAP 3 — KOMPETENSI & SKOR

### A. Job Potency Competencies (9.16)

Kompetensi potensi per posisi — mapping competency ke level value + bobot.

- Isi: organization_id:
  - **job_management_value_id** → level kompetensi (dari Job Value).
  - **competency_id** → relasi ke master Competency (di module Competency).
  - **weight** — bobot kompetensi (decimal).
- Endpoint: `GET/POST /potency-competencies`, `GET/PUT/DELETE /potency-competencies/:id`

### B. Job Scores (9.17)

Skor agregat jabatan — dihitung per posisi organisasi.

- **Updatable** via `PUT /scores/org/:orgId` (upsert):
  - `job_value_with_financial` — total value jabatan dengan otoritas keuangan.
  - `job_value_without_financial` — total value jabatan tanpa otoritas keuangan.
  - `has_financial_authority` — apakah posisi punya otoritas keuangan.
  - `components` (JSON) — komponen skor breakdown.
  - `sub_component_points` (JSON) — poin sub-komponen.
  - `is_complete`, `calculated_at`, `completed_at` — status kelengkapan.
- Endpoint: `GET /scores`, `GET /scores/org/:orgId`, `PUT /scores/org/:orgId`

### C. Job Competency Groups (9.18)

Bobot kompetensi per organisasi per kategori.

- Isi: organization_id, **category** (nama kategori), **weight** (bobot).
- Endpoint: `GET/POST /competency-groups`, `GET/PUT/DELETE /competency-groups/:id`

---

## 6. Dashboard

Menu **Core HR → Job Management** → halaman utama `JobManagement.vue`.

- Ringkasan jumlah data per entitas (titles, values, objectives, responsibilities, dll.).
- Endpoint: `GET /dashboard` → response `JobManagementDashboardResponse`.

---

## 7. Ringkasan Status & Transisi

| Entitas | Status | Transisi |
|---|---|---|
| **Semua data jabatan** | — | CRUD biasa, tidak ada state machine |
| **Job Title** | `status` (int8, optional) | Manual via update |
| **Job Score** | `is_complete` (bool) | Upsert via `PUT /scores/org/:orgId` |

> Modul ini tidak memiliki approval flow atau status workflow — semua data diisi langsung
> oleh user dengan permission `jobmanagement.create/update/delete`.

---

## 8. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Organization** | Dependency wajib — semua data jabatan di-scoping ke `organization_id` (posisi dalam struktur organisasi) |
| **Settings (Grading)** | Master Grading → di-link di Job Identifications (`grading_id`) |
| **Settings (Education)** | Master Education → di-link di Job Education Experiences (`education_id`) |
| **Settings (Education Majors)** | Master Jurusan → many-to-many via pivot `job_management_majors` |
| **Settings (Job Families)** | Master Bidang Pekerjaan → many-to-many via pivot `job_management_job_family` |
| **Competency** | Master Competency → di-link di Job Potency Competencies (`competency_id`) |
| **Employee** | 🚫 Belum terintegrasi — job data belum di-push ke employee profil |
| **Recruitment** | 🚫 Belum terintegrasi — job requirements belum dibaca oleh recruitment requisition |
| **Training** | 🚫 Belum terintegrasi — competency groups belum di-push ke training needs |
| **Payroll** | 🚫 Belum terintegrasi — job scores belum menentukan grade/compensation |

---

## 9. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Core HR → Job Management (dropdown) | | Parent menu di sidebar dengan 2 child |
| ├ Job Management (card hub) | `JobManagement.vue` | Dashboard card: ringkasan jumlah data per entitas + card menu ke form |
| ├ Job Management Form | `JobManagementForm.vue` | Form analisis jabatan multi-section: Objectives, Identifications, Responsibilities, Education Experiences, HR Authorities, Operational Authorities, Working Activities, Working Risks, Relationships (with nested Details), Subordinate Controls, Assets, Financials, Potency Competencies, Job Scores, Competency Groups |
| ├ Job Values Mapping | `JobValuesIndex.vue` | List values per tipe + cluster mapping |
| ├ Job Value Form | `JobValuesForm.vue` | Form edit values per type tertentu |
| ├ Job Value Section | `JobValueSection.vue` | Section component untuk value editor |
| ├ Job Value Cluster Card | `JobValueClusterCard.vue` | Card untuk cluster mapping per type |
| ├ Job Objective Section | `JobObjectiveSection.vue` | Section di form: tujuan jabatan |
| ├ Job Identification Section | `JobIdentificationSection.vue` | Section di form: identitas + grading |
| ├ Job Responsibility Section | `JobResponsibilitySection.vue` | Section di form: tanggung jawab |
| ├ Job Edu Exp Section | `JobEduExpSection.vue` | Section di form: pendidikan & pengalaman |
| ├ Job HR Authority Section | `JobHRAuthoritySection.vue` | Section di form: kewenangan SDM |
| ├ Job Op Authority Section | `JobOpAuthoritySection.vue` | Section di form: kewenangan operasional |
| ├ Job Activity Section | `JobActivitySection.vue` | Section di form: aktivitas kerja |
| ├ Job Risk Section | `JobRiskSection.vue` | Section di form: risiko kerja |
| ├ Job Relationship Section | `JobRelationshipSection.vue` | Section di form: hubungan kerja (dengan nested details) |
| ├ Job Subordinate Section | `JobSubordinateSection.vue` | Section di form: bawahan |
| ├ Job Asset Section | `JobAssetSection.vue` | Section di form: aset jabatan |
| ├ Job Financial Section | `JobFinancialSection.vue` | Section di form: keuangan jabatan |
| ├ Job Potency Section | `JobPotencySection.vue` | Section di form: kompetensi potensi |
| ├ Job Score Section | `JobScoreSection.vue` | Section di form: skor jabatan |
| ├ Job Score Summary | `JobScoreSummary.vue` | Ringkasan skor di form |
| └ Job Competency Group Section | `JobCompetencyGroupSection.vue` | Section di form: bobot kompetensi |

---

## 10. Endpoint API Utama

Semua di bawah `/api/v1/tenant/job-management/`.

| Area | Endpoint |
|---|---|
| **Dashboard** | `GET /dashboard` |
| **Titles (9.1)** | `GET/POST /titles`, `GET/PUT/DELETE /titles/:id` |
| **Title Subs (9.2)** | `GET/POST /titles/:id/subs`, `GET/PUT/DELETE /titles/:id/subs/:subId` |
| **Values (9.3)** | `GET/POST /values`, `GET /values/tree`, `GET/PUT/DELETE /values/:id` |
| **Value Clusters** | `GET /values/clusters/:type`, `PUT /values/clusters/:type` |
| **Objectives (9.4)** | `GET/POST /objectives`, `GET/PUT/DELETE /objectives/:id` |
| **Identifications (9.5)** | `GET/POST /identifications`, `GET/PUT/DELETE /identifications/:id` |
| **Responsibilities (9.6)** | `GET/POST /responsibilities`, `GET/PUT/DELETE /responsibilities/:id` |
| **Education Experiences (9.7)** | `GET/POST /education-experiences`, `GET/PUT/DELETE /education-experiences/:id` |
| **HR Authorities (9.8)** | `GET/POST /hr-authorities`, `GET/PUT/DELETE /hr-authorities/:id` |
| **Operational Authorities (9.9)** | `GET/POST /operational-authorities`, `GET/PUT/DELETE /operational-authorities/:id` |
| **Working Activities (9.10)** | `GET/POST /working-activities`, `GET/PUT/DELETE /working-activities/:id` |
| **Working Risks (9.11)** | `GET/POST /working-risks`, `GET/PUT/DELETE /working-risks/:id` |
| **Relationships (9.12)** | `GET/POST /relationships`, `GET/PUT/DELETE /relationships/:id` |
| **Relationship Details (9.12b)** | `GET/POST /relationships/:id/details`, `GET/PUT/DELETE /relationships/:id/details/:detailId` |
| **Subordinate Controls (9.13)** | `GET/POST /subordinate-controls`, `GET/PUT/DELETE /subordinate-controls/:id` |
| **Assets (9.14)** | `GET/POST /assets`, `GET/PUT/DELETE /assets/:id` |
| **Financials (9.15)** | `GET/POST /financials`, `GET/PUT/DELETE /financials/:id` |
| **Potency Competencies (9.16)** | `GET/POST /potency-competencies`, `GET/PUT/DELETE /potency-competencies/:id` |
| **Scores (9.17)** | `GET /scores`, `GET /scores/org/:orgId`, `PUT /scores/org/:orgId` |
| **Competency Groups (9.18)** | `GET/POST /competency-groups`, `GET/PUT/DELETE /competency-groups/:id` |

---

## 11. Catatan Penting

- **Modul core** (`IsCore: true`) — selalu aktif, dependency ke `organization` saja.
- **Scope organization** — mayoritas entitas di-scoping ke `organization_id` (posisi dalam
  struktur organisasi). User memilih posisi organisasi, lalu mengisi data jabatan untuk
  posisi tersebut di multi-section form (`JobManagementForm.vue`).
- **Job Values sebagai master dropdown** — banyak entitas jabatan meng-relasi ke Job Value
  via `job_management_value_*` FK. Value dibagi per tipe (environment, hazard, asset, cash,
  authority, impact, frequency, relationship, education, experience). Tipe inilah yang
  digunakan untuk klasifikasi.
- **Tree view values** — endpoint `GET /values/tree` mengembalikan hierarki values per tipe
  untuk keperluan UI dropdown hierarkis.
- **Value Cluster** — mapping type ↔ cluster kompetensi diatur di halaman "Mapping Job Value"
  (`PUT /values/clusters/:type`). Cluster ini menentukan kelompok kompetensi untuk scoring.
- **Job Score** — skor dihitung dari komponen JSON (`components`) dan sub-komponen
  (`sub_component_points`). Endpoint `PUT /scores/org/:orgId` melakukan upsert (create or update).
- **No approval flow** — tidak ada integrasi Central Approval di modul ini.
- **Nested routes** — Titles memiliki subs (`/titles/:id/subs`), Relationships memiliki
  details (`/relationships/:id/details`). Keduanya menggunakan wildcard `:id` di parent level.
- **Sidebar** — di bawah "Core HR" section dengan dropdown (2 child: "Job Values Mapping" &
  "Job Management"). Auto-open saat salah satu child aktif.
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
