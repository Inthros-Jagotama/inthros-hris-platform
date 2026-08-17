# Alur Pengisian Competency 360 Assessment (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Competency 360** — penilaian kompetensi
employee berbasis multi-rater (Self, Superior/Manager, Peer, Subordinate) dari setup master sampai
hasil akhir — pola runbook seperti [`module-payroll-user-flow.md`](module-payroll-user-flow.md).

- Plan pengembangan: [`module-competency-360-development-plan.md`](../module-competency-360-development-plan.md) — fase 1–7 selesai, fase 8 (testing/hardening) sebagian
- Lokasi kode: `backend/internal/modules/competency/` · `frontend/tenant/src/views/modules/competency/`
- Daftar endpoint + contoh curl: [`../api/api-usage-guide.md`](../api/api-usage-guide.md) → §8.2 (tabel Competency)

---

## 1. Ringkasan Alur End-to-End

```
SETUP MASTER (sekali)                    EVENT (per periode)                     HASIL
┌──────────────────────┐   ┌─────────────────────────────────────────────┐   ┌──────────────┐
│ Rating Scales        │   │ Event (draft→active)                        │   │ Approved      │
│ Competencies+Indicator│──▶│  └─ Target (employee+template)             │──▶│  → finalisasi │
│ Assessment Templates  │   │      └─ Rater (self/manager/peer/sub)      │   │  → hitung skor│
└──────────────────────┘   │           └─ Isi & submit assessment        │   │  → gap/report │
                           └─────────────────────────────────────────────┘   └──────────────┘
```

- **Status event:** `draft → active → closed`
- **Status target:** `draft → submitted → finalized` (finalized via Central Approval)
- **Status rater:** `assigned → submitted`

---

## 2. Entitas Utama (360)

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Rating Scale | `competency_rating_scales` (+ items) | Skala nilai (mis. 1–5) + label + bobot tiap level |
| Competency | `competencies` | Master kompetensi (cluster teknis/manajerial) |
| Indicator | `competency_indicators` | Pernyataan perilaku per kompetensi |
| Assessment Template | `competency_assessment_templates` (+ indicators) | Paket indikator yang dinilai dalam satu assessment |
| Assessment Event | `competency_events` | Sesi assessment per periode (draft/active/closed) |
| Event Target | `competency_event_targets` | Employee yang dinilai + template + posisi |
| Rater | `competency_event_target_raters` | Penilai per target (self/manager/peer/subordinate), status assigned → submitted |
| Score / Score Detail | `competency_scores` (+ details) | Hasil perhitungan setelah finalisasi |

---

## 3. TAHAP 1 — SETUP MASTER (dikerjakan sekali)

### A. Rating Scales

Menu **Competency → Rating Scales** (`/competencies/values`, `CompetencyValues.vue`).

- Buat skala (kode, nama, status) + item skala: `value` (1–8), `label`, bobot, urutan.
- Endpoint: `GET/POST /rating-scales`, `GET/PUT/DELETE /rating-scales/:id`.

### B. Kompetensi & Indikator

Menu **Competency** (`Competencies.vue` — hub) → **Indicators** (`/competencies/indicators`, `CompetencyIndicators.vue`).

- Master kompetensi: `GET/POST /competencies`, `GET/PUT/DELETE /competencies/:id`.
- Indikator perilaku per kompetensi: `GET/POST /indicators`, `GET/PUT/DELETE /indicators/:id`.

### C. Assessment Templates

Menu **Competency → Templates** (`/competencies/templates`, `AssessmentTemplates.vue` + `AssessmentTemplateForm.vue`).

- Buat template → pilih kompetensi + indikator yang dinilai (`POST /templates`, `POST /templates/:id/indicators`).
- Template menjadi acuan penilaian setiap target.

---

## 4. TAHAP 2 — EVENT & RATER (per periode)

### A. Buat Event

Menu **Competency → Events** (`/competencies/events`, `CompetencyEvents.vue`).

- Buat event assessment: `GET/POST /events`, `GET/PUT/DELETE /events/:id`.
- Status `draft` → `active` saat siap diisi.

### B. Tambah Target

- Target = employee yang dinilai, dengan template + posisi:
  `GET/POST /event-targets`, `GET/PUT/DELETE /event-targets/:id`.
- Target berstatus `draft` sampai rater mengisi & diajukan.

### C. Assign Rater

Menu **Competency → Rater Assignment** (`/competencies/raters`, `RaterAssignment.vue`).

- Tambah rater per target: `POST /event-targets/:id/raters` · lihat `GET /event-targets/:id/raters` · hapus via `DELETE /raters/:id`.
- **Suggested raters** otomatis dari struktur org (`GET /event-targets/:id/suggested-raters`) —
  termasuk **self rater** otomatis (assigned).
- Tipe rater: Self, Superior/Manager, Peer, Subordinate.

---

## 5. TAHAP 3 — PENGISIAN (rater mengisi)

### A. My Assessments (Self & rater lain yang menilai)

Menu **Competency → My Assessment** (`/competencies/my-assessments`, `MyAssessments.vue`).

- Daftar assessment yang ditugaskan ke saya: `GET /my-assessments`.
- Isi per assessment: `GET /my-assessments/:id` + simpan respons `POST /my-assessments/:id/responses`.
- Tombol **Save Draft** (simpan tanpa submit) vs **Submit** (`POST /my-assessments/:id/submit` →
  status rater `submitted`, tidak bisa diubah lagi).
- Kolom status menampilkan `assigned` (info) / `submitted` (success).

### B. Manager Assessments (tinjauan atasan)

Menu **Competency → Manager Assessment** (`/competencies/manager-assessments`, `ManagerAssessments.vue`).

- Penilaian oleh manager untuk bawahan di org tree: `GET /manager-assessments`.

---

## 6. TAHAP 4 — APPROVAL & HASIL

### A. Submit Approval (target → submitted)

- Saat seluruh rater (atau cukup menurut kebijakan) sudah submit, target diajukan:
  `POST /event-targets/:id/submit-approval` → status target `submitted`, instance approval
  dibuat di Central Approval (modul competency — pola backward-compatible performance).
- Approver menindaklanjuti di halaman **Approvals** generik.

### B. Finalisasi & Perhitungan

- Callback approval `APPROVED` → target **finalized** + `finalized_at` (`HandleAssessmentApprovalStatusChange`).
- **Calculation engine** menghitung skor setelah finalisasi (weighted rater + weighting
  konfigurasi): `GET /scores`, `GET /scores/:id/details`, `GET /score-details`.
- Hasil akhir: **competency score**, **gap**, **self vs others comparison**.

### C. Lihat Hasil & Report

| Menu | Halaman | Endpoint |
|---|---|---|
| Competency → Results | `AssessmentResult.vue` | `GET /employees/:employee/result` · `GET /employees/:employee/gap` |
| Competency → Reports | `CompetencyReports.vue` | `GET /reports/hr` (aggregate) · `GET /reports/manager` · `GET /employees/:employee/report` |

- **Anonimitas**: feedback rater tidak membocorkan identitas (kecuali self).
- Assessment yang sudah **finalized** tidak bisa dimodifikasi tanpa mekanisme resmi.

---

## 7. Ringkasan Status

| Entitas | Status | Transisi |
|---|---|---|
| Event | `draft → active → closed` | buat → aktifkan → tutup |
| Target | `draft → submitted → finalized` | tambah rater → submit-approval → approval APPROVED |
| Rater | `assigned → submitted` | assign (auto self + suggested + manual) → submit assessment |

---

## 8. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Organization / Employee** | Struktur org untuk suggested raters & manager (org tree, `supervisor_id`) |
| **Central Approval** | Approval finalisasi target assessment (modul competency) |
| **Job Management** | Referensi kompetensi requirement posisi (kompatibilitas existing) |
| **Notification** | Task approval via approval engine (`APPROVAL_TASK_ASSIGNED`) |

---

## 9. Peta Halaman UI

| Menu | Halaman | Permission (strict) |
|---|---|---|
| Competency (hub) | `Competencies.vue` | Kartu menu per permission |
| → Rating Scales | `CompetencyValues.vue` | `competency.settings.view` |
| → Indicators | `CompetencyIndicators.vue` | `competency.settings.view` |
| → Templates | `AssessmentTemplates.vue` / `AssessmentTemplateForm.vue` | `competency.settings.view` |
| → Events | `CompetencyEvents.vue` | `competency.settings.view` |
| → Rater Assignment | `RaterAssignment.vue` | `competency.settings.view` |
| → My Assessment | `MyAssessments.vue` | `competency.assessment.view` |
| → Manager Assessment | `ManagerAssessments.vue` | `competency.assessment.view` |
| → Results | `AssessmentResult.vue` | `competency.report.view` |
| → Reports | `CompetencyReports.vue` | `competency.report.view` |

---

## 10. Endpoint API Utama

Semua di bawah `/api/v1/tenant/competency/`.

| Area | Endpoint |
|---|---|
| Master | `GET/POST /competencies`, `GET/PUT/DELETE /competencies/:id`, `GET/POST /competence-values`, `GET/PUT/DELETE /indicators/:id` |
| Rating Scale | `GET/POST /rating-scales`, `GET/PUT/DELETE /rating-scales/:id` |
| Template | `GET/POST /templates`, `GET/PUT/DELETE /templates/:id`, `POST/GET /templates/:id/indicators` |
| Event | `GET/POST /events`, `GET/PUT/DELETE /events/:id`, `GET/POST /event-targets`, `GET/PUT/DELETE /event-targets/:id` |
| Rater | `GET /event-targets/:id/suggested-raters`, `POST/GET /event-targets/:id/raters`, `DELETE /raters/:id` |
| Assessment | `GET /my-assessments`, `GET /my-assessments/:id`, `POST /my-assessments/:id/responses`, `POST /my-assessments/:id/submit`, `GET /manager-assessments` |
| Approval | `POST /event-targets/:id/submit-approval` |
| Hasil | `GET /scores`, `GET /scores/:id`, `GET /scores/:id/details`, `GET /score-details`, `GET /employees/:employee/result`, `GET /employees/:employee/gap`, `GET /employees/:employee/report` |
| Report | `GET /reports/hr`, `GET /reports/manager` |

---

## 11. Catatan Penting

- **Rater tidak bisa mengubah** assessment setelah submit (status `submitted`).
- **Hasil dihitung setelah approval** — belum finalized = belum ada skor final.
- **Self rater dibuat otomatis** saat assign (status assigned); tetap wajib diisi & di-submit.
- **Anonimitas rater** dijaga pada aggregated report (kecuali self).
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
