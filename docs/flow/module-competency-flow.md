# Alur Competency Management & 360 Assessment (Runbook)

Dokumen ini menjelaskan **cara pakai** modul **Competency** — dictionary kompetensi,
values, events, scoring, assessment templates, rater assignment, employee self-assessment,
manager assessment, approval, results, gap analysis, dan reports.

- Lokasi kode: `backend/internal/modules/competency/`
- Halaman UI: `frontend/tenant/src/views/modules/competency/`
- Module slug: `competency`

---

## 1. Ringkasan Alur End-to-End

```
SETUP MASTER                     360 ASSESSMENT                        REPORTING
┌──────────────────┐   ┌────────────────────────────────────┐   ┌──────────────────────┐
│ Competencies      │   │ Create Event → Add Targets          │   │ Employee Result      │
│ Values (legacy +  │──▶│ Assign Raters (suggested/manual)    │──▶│ Employee Gap         │
│   structured)     │   │ Employee Self-Assessment            │   │ Employee Report      │
│ Rating Scales     │   │ Manager Assessment                  │   │ Manager Report       │
│ Indicators        │   │ Submit → Approval                   │   │ HR Report            │
│ Assessment        │   │ Result Calculation                  │   └──────────────────────┘
│   Templates       │   └────────────────────────────────────┘
└──────────────────┘
```

---

## 2. Entitas Utama

| Entitas | Deskripsi |
|---|---|
| Competency | Dictionary kompetensi (nama, deskripsi, kategori) |
| CompetenceValue | Nilai kompetensi (legacy) |
| CompetencyValue | Nilai kompetensi (structured — nama, range, deskripsi per level) |
| CompetencyEvent | Event assessment (nama, periode, template, status) |
| CompetencyEventTarget | Karyawan yang dinilai dalam event |
| CompetencyScore | Skor kompetensi per karyawan |
| CompetencyScoreDetail | Detail skor per kompetensi |
| RatingScale | Rating scale untuk 360 assessment (nama, items dengan range + label) |
| AssessmentTemplate | Template 360 assessment (nama, competencies, rater types) |
| Indicator | Indikator penilaian dalam template |
| Rater | Rater assignment per target (type: SELF/PEER/SUPERIOR/SUBORDINATE/EXTERNAL) |

---

## 3. SETUP — Master Data

| Menu | Endpoint | Deskripsi |
|---|---|---|
| Competencies | `GET/POST /competency/competencies` | CRUD dictionary kompetensi |
| Values (legacy) | `GET/POST /competency/competence-values` | Legacy competency values |
| Values (structured) | `GET/POST /competency/values` | Structured competency values |
| Rating Scales | `GET/POST /competency/rating-scales` | Rating scale untuk 360 assessment |
| Indicators | `GET/POST /competency/indicators` | Indikator penilaian |
| Templates | `GET/POST /competency/templates` | Template 360 assessment |

---

## 4. EVENTS — Assessment Event

1. **Buat event** — `POST /competency/events`: name, period, assessment_template_id, description.
2. **Tambah target** — `POST /competency/event-targets`: event_id, employee_id (karyawan yang akan dinilai).
3. **Assign rater** — `POST /competency/event-targets/:id/raters`: target_id, rater_id, rater_type (SELF/PEER/SUPERIOR/SUBORDINATE/EXTERNAL).
4. **Suggested rater** — `GET /competency/event-targets/:id/suggested-raters`: rekomendasi rater berdasarkan org tree.
5. **Daftar events** — `GET /competency/events`.
6. **Daftar targets** — `GET /competency/event-targets`.

---

## 5. SCORING — Traditional

1. **Buat skor** — `POST /competency/scores`: event_target_id, overall_score.
2. **Detail skor** — `POST /competency/score-details`: score_id, competency_id, score, notes.
3. **Daftar skor** — `GET /competency/scores`.
4. **Detail per score** — `GET /competency/scores/:id/details`.

---

## 6. 360 ASSESSMENT — Employee Self-Assessment

1. **Lihat assessment saya** — `GET /competency/my-assessments`: daftar assessment yang perlu saya isi.
2. **Detail assessment** — `GET /competency/my-assessments/:id`: daftar indikator + pertanyaan.
3. **Simpan jawaban** — `POST /competency/my-assessments/:id/responses`: simpan jawaban per indikator.
4. **Submit** — `POST /competency/my-assessments/:id/submit`: finalisasi jawaban.

---

## 7. 360 ASSESSMENT — Manager Assessment

1. **Lihat assessment bawahan** — `GET /competency/manager-assessments`: daftar bawahan yang perlu saya nilai.
2. **Isi assessment** — gunakan endpoint yang sama (save responses + submit).

---

## 8. 360 ASSESSMENT — Approval & Results

1. **Submit untuk approval** — `POST /competency/event-targets/:id/submit-approval`: kumpulkan semua jawaban, hitung skor, kirim ke approval.
2. **Employee Result** — `GET /competency/employees/:employee/result`: skor akhir per kompetensi.
3. **Employee Gap** — `GET /competency/employees/:employee/gap`: gap antara skor saat ini vs target.
4. **Employee Report** — `GET /competency/employees/:employee/report`: laporan lengkap per karyawan.
5. **Manager Report** — `GET /competency/reports/manager`: laporan untuk manager.
6. **HR Report** — `GET /competency/reports/hr`: laporan untuk HR (agregasi).

---

## 9. Ringkasan Status

| Entitas | Status |
|---|---|
| Event | `DRAFT → ACTIVE → COMPLETED` |
| Event Target | `PENDING → ASSESSING → SUBMITTED → APPROVED` |
| Assessment (rater) | `DRAFT → SUBMITTED` |

---

## 10. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Employee** | Target assessment, rater assignment |
| **Organization** | Suggested rater dari org tree |
| **Performance** | Competency score sebagai input scoring |
| **Employee Movement** | Competency score dibaca untuk eligibility |
| **Career Intelligence** | Competency requirement pada career path |
| **Training** | Training effectiveness berdasarkan competency gap |

---

## 11. Peta Halaman UI

| Menu | Halaman |
|---|---|
| Competencies (hub) | `Competencies.vue` |
| Rating Scales | `CompetencyValues.vue` |
| Indicators | `CompetencyIndicators.vue` |
| Assessment Templates | `AssessmentTemplates.vue` / `AssessmentTemplateForm.vue` |
| Events | `CompetencyEvents.vue` |
| Rater Assignment | `RaterAssignment.vue` |
| My Assessment | `MyAssessments.vue` |
| Manager Assessment | `ManagerAssessments.vue` |
| Reports | `CompetencyReports.vue` |
| Assessment Result | `AssessmentResult.vue` |

---

## 12. Endpoint API Utama

Semua di bawah `/api/v1/tenant/competency/`.

| Area | Endpoint |
|---|---|
| Competencies | `GET/POST /competencies`, `GET/PUT/DELETE .../:id` |
| Values (legacy) | `GET/POST /competence-values`, `GET/PUT/DELETE .../:id` |
| Values (structured) | `GET/POST /values`, `GET/PUT/DELETE .../:id` |
| Events | `GET/POST /events`, `GET/PUT/DELETE .../:id` |
| Event Targets | `GET/POST /event-targets`, `GET/PUT/DELETE .../:id` |
| Scores | `GET/POST /scores`, `GET/PUT/DELETE .../:id` |
| Score Details | `GET/POST /scores/:id/details`, `GET/PUT/DELETE .../:id` |
| Rating Scales | `GET/POST /rating-scales`, `GET/PUT/DELETE .../:id` |
| Templates | `GET/POST /templates`, `GET/PUT/DELETE .../:id` |
| Template Indicators | `GET/PUT /templates/:id/indicators` |
| Indicators | `GET/POST /indicators`, `GET/PUT/DELETE .../:id` |
| Rater Assignment | `POST /event-targets/:id/raters`, `GET .../:id/raters`, `GET .../:id/suggested-raters`, `DELETE /raters/:id` |
| Manager Assessment | `GET /manager-assessments` |
| My Assessment | `GET /my-assessments`, `GET /my-assessments/:id`, `POST .../:id/responses`, `POST .../:id/submit` |
| Approval | `POST /event-targets/:id/submit-approval` |
| Results & Gap | `GET /employees/:employee/result`, `GET /employees/:employee/gap` |
| Reports | `GET /employees/:employee/report`, `GET /reports/manager`, `GET /reports/hr` |

---

## 13. Catatan Penting

- **Rater types**: SELF (diri sendiri), PEER (rekan kerja), SUPERIOR (atasan), SUBORDINATE (bawahan), EXTERNAL (pihak luar).
- **Suggested raters** diambil dari org tree — memudahkan manager memilih rater yang relevan.
- **Approval** dikumpulkan setelah semua rater selesai mengisi → submit untuk approval → skor final dihitung.
- **Gap analysis** membandingkan skor aktual vs target kompetensi (berdasarkan grading/job family).
- **Competency values** punya 2 versi: legacy (flat) dan structured (per level dengan range).
