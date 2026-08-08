= HRIS Platform — OpenAPI Comprehensive Report (v19) =


**Generated:** 08 August 2026
**Spec Version:** 1.6.3
**Total Paths:** 473
**Total Endpoints (methods):** 821
**Total Schemas:** 513
**Total Tags:** 32

> 🔗 **Index dokumentasi:** [`docs/README.md`](README.md) · **Terkait:** [`api/api-usage-guide.md`](api/api-usage-guide.md) · [`go-module-architecture-report.md`](go-module-architecture-report.md)

## Coverage Summary

| Metric | Coverage | % |
|---|---|---|
| Endpoints with `summary` | 821/821 | 100% |
| Endpoints with `description` | 821/821 | 100% |
| Endpoints with `operationId` | 821/821 | 100% |

## Response Format & Bilingual Support

Semua endpoint mengembalikan response dengan format standar:

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "Created successfully"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource not found"
  }
}
```
> With `Accept-Language: id`: `"message": "Resource tidak ditemukan"`

### Validation Error Response
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": {
      "email": ["Must be a valid email address"],
      "nik": ["Invalid NIK format, must be 16 digits"]
    }
  }
}
```
> With `Accept-Language: id`: `"message": "Validasi gagal"`, `"email": ["Format email tidak valid"]`

### Paginated Response
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Bilingual Support (Bahasa Indonesia & English)

API supports two languages for response messages:

| Header | Language | Description |
|--------|----------|-------------|
| (no `Accept-Language`) | **English** | Default language |
| `Accept-Language: id` | **Bahasa Indonesia** | All messages automatically switch to Indonesian |

This header affects all response messages, including:
- **Success messages** (created, updated, deleted)
- **Error messages** (not found, internal error, validation error)
- **Field-level validation errors**

### Custom Indonesian Validators

Tenant endpoints support validation for Indonesian data formats:

| Tag | Format | Example | Description |
|-----|--------|---------|-------------|
| `nik` | 16 digits | `3273010101900001` | National ID (KTP) |
| `npwp` | 15-16 digits | `0123456789012345` | Tax ID |
| `npwp_format` | XX.XXX.XXX.X-XXX.XXX | `01.234.567.8-901.234` | Tax ID (formatted) |
| `kk` | 16 digits | `1234567890123456` | Family Register |
| `phone_id` | +628/08xx (7-11 digits) | `08123456789` | Phone Number |
| `postal_code` | 5 digits | `12345` | Postal Code |
| `date_id` | YYYY-MM-DD | `2026-12-31` | Date (ISO 8601) |
| `passport` | 1 letter + 8 digits | `A12345678` | Passport |
| `sim` | 12 digits | `123456789012` | Driver License |
| `no_rekening` | 8-20 digits | `1234567890` | Bank Account |

## 1. Endpoints per Module (Tag)

| # | Tag | Endpoints | Paths |
|---|---|---|---|
| 1 | Tenant: Performance Management | 146 | 103 |
| 2 | Tenant: Settings | 107 | 44 |
| 3 | Tenant: Job Management | 96 | 40 |
| 4 | Tenant: Workforce Intelligence & Strategic Pl... | 68 | 58 |
| 5 | Tenant: Payroll & Compensation Engine | 47 | 24 |
| 6 | Tenant: Employees | 36 | 23 |
| 7 | Tenant: Competency Management | 35 | 15 |
| 8 | Tenant: Training & Development Management | 35 | 15 |
| 9 | Tenant: Recruitment & Onboarding (ATS) | 33 | 16 |
| 10 | Tenant: Time & Attendance | 30 | 15 |
| 11 | Tenant: Leave & Time Off | 23 | 12 |
| 12 | Tenant: Career Intelligence | 19 | 11 |
| 13 | Tenant: Organizations | 18 | 11 |
| 14 | Tenant: Approval | 17 | 11 |
| 15 | Tenant: Employee Movement & Career Management | 16 | 10 |
| 16 | Tenant: Reimbursement & Claim | 15 | 7 |
| 17 | Platform: Companies | 11 | 8 |
| 18 | Platform: RBAC Management | 10 | 6 |
| 19 | Platform: Packages | 9 | 6 |
| 20 | Tenant: RBAC Management | 8 | 6 |
| 21 | Platform: Modules | 7 | 5 |
| 22 | Platform: Users | 6 | 3 |
| 23 | Platform: Licenses | 5 | 2 |
| 24 | Platform: Monitoring | 5 | 5 |
| 25 | Health | 4 | 4 |
| 26 | Tenant: Packages | 4 | 4 |
| 27 | Tenant: User Accounts | 3 | 2 |
| 28 | Platform: Auth | 2 | 2 |
| 29 | Public | 2 | 2 |
| 30 | Tenant Auth | 2 | 2 |
| 31 | Tenant: Company | 2 | 1 |
| | **TOTAL** | **821** | **473** |

## 2. Module Detail

### Tenant: Performance Management
**Description:** Performance Management â€” BSC (Balanced Scorecard) based KPI and performance evaluation module. Includes performance periods, BSC perspectives, KPI templates and indicators, employee evaluations, and individual performance targets with full status workflow (DRAFT->PLAN_SUBMITTED->PLAN_APPROVED->ACTUAL_SUBMITTED->ACTUAL_APPROVED->COMPLETED).
**Endpoints:** 146 | **Paths:** 103
**Methods:** DELETE=22 GET=53 POST=42 PUT=29

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/performance/indicator-formulas` | Create indicator formula | Buat formula kalkulasi skor KPI (MANUAL/HIGHER_BETTER/LOWER_BETTER/RANGE). |
| `GET` | `/api/v1/tenant/performance/indicator-formulas` | List indicator formulas | Ambil daftar formula kalkulasi skor KPI dengan pagination. |
| `GET` | `/api/v1/tenant/performance/indicator-formulas/{id}` | Get indicator formula by ID | Ambil satu formula kalkulasi skor KPI. |
| `PUT` | `/api/v1/tenant/performance/indicator-formulas/{id}` | Update indicator formula | Perbarui formula kalkulasi skor KPI. |
| `DELETE` | `/api/v1/tenant/performance/indicator-formulas/{id}` | Delete indicator formula | Hapus satu formula kalkulasi skor KPI. |
| `POST` | `/api/v1/tenant/performance/kpi/attachments` | Create performance attachment | Lampirkan file bukti/dokumen pendukung pada evaluation detail. |
| `GET` | `/api/v1/tenant/performance/kpi/attachments/{id}` | Get performance attachment by ID | Ambil satu lampiran performance berdasarkan ID. |
| `PUT` | `/api/v1/tenant/performance/kpi/attachments/{id}` | Update performance attachment | Perbarui deskripsi lampiran performance. |
| `DELETE` | `/api/v1/tenant/performance/kpi/attachments/{id}` | Delete performance attachment | Hapus satu lampiran performance. |
| `POST` | `/api/v1/tenant/performance/kpi/comments` | Create performance comment | Tambahkan komentar/review pada sebuah performance evaluation. |
| `GET` | `/api/v1/tenant/performance/kpi/comments/{id}` | Get performance comment by ID | Ambil satu komentar performance berdasarkan ID. |
| `PUT` | `/api/v1/tenant/performance/kpi/comments/{id}` | Update performance comment | Perbarui isi komentar performance. |
| `DELETE` | `/api/v1/tenant/performance/kpi/comments/{id}` | Delete performance comment | Hapus satu komentar performance. |
| `GET` | `/api/v1/tenant/performance/kpi/components` | List performance components | Ambil daftar komponen scoring KPI dengan pagination. |
| `GET` | `/api/v1/tenant/performance/kpi/components/{id}` | Get performance component by ID | Ambil detail satu komponen scoring KPI. |
| `PUT` | `/api/v1/tenant/performance/kpi/components/{id}` | Update performance component | Perbarui kode, nama, deskripsi, urutan, atau status aktif komponen scoring. |
| `GET` | `/api/v1/tenant/performance/kpi/dashboard/employee/{employee_id}` | Get KPI employee dashboard | Dashboard KPI untuk employee: daftar KPI, progress, achievement, dan skor evaluasi periode aktif. |
| `GET` | `/api/v1/tenant/performance/kpi/dashboard/hr` | Get KPI HR dashboard | Dashboard KPI untuk HR: ringkasan evaluasi seluruh organisasi (total, status, skor rata-rata, distribusi). |
| `GET` | `/api/v1/tenant/performance/kpi/dashboard/manager/{manager_id}` | Get KPI manager dashboard | Dashboard KPI untuk manager: ringkasan tim, anggota tim beserta skor & status KPI masing-masing. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluation-details` | Create evaluation detail | Add a BSC perspective detail to a performance evaluation, including achievement percentage, weight, and score for that perspective. |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluation-details/{id}` | Update evaluation detail | Update a BSC perspective detail's achievement percentage, weight, score, or description. |
| `DELETE` | `/api/v1/tenant/performance/kpi/evaluation-details/{id}` | Delete evaluation detail | Permanently delete a BSC perspective detail from the evaluation. |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluation-details/{id}/actual` | Update KPI evaluation detail actual | Input nilai aktual (actual) untuk satu detail evaluasi — skor & achievement dihitung ulang. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluation-details/{id}/attachments` | List attachments by evaluation detail | Ambil semua lampiran pada evaluation detail tertentu. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluation-details/{id}/progress` | List progress by evaluation detail | Ambil semua progres KPI yang tercatat untuk evaluation detail tertentu. |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluation-details/{id}/target` | Update evaluation detail target | Input nilai target untuk satu detail evaluasi KPI (fase planning/target). |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations` | Create performance evaluation | Start a new performance evaluation for an employee. Links the employee to a performance period and KPI template for assessment. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations` | List performance evaluations | Retrieve a paginated list of performance evaluations, optionally filtered by employee, organization, period, or status. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/snapshot` | Create KPI evaluation with snapshot | Buat evaluasi KPI baru sekaligus snapshot KPI dari template ke evaluation details (nilai target terkunci). |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}` | Get performance evaluation by ID | Retrieve detailed information about a specific performance evaluation, including its status, scores, and linked targets. |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluations/{id}` | Update performance evaluation | Update evaluation metadata such as supervisor assignment or notes. Only provided fields will be updated. |
| `DELETE` | `/api/v1/tenant/performance/kpi/evaluations/{id}` | Delete performance evaluation | Permanently delete a performance evaluation. Only evaluations in DRAFT status can be deleted. |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluations/{id}/actuals` | Bulk update KPI evaluation actuals | Input nilai aktual sekaligus untuk banyak evaluation detail dalam satu request. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/approve` | Approve KPI evaluation | Setujui evaluasi KPI (status → APPROVED). |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/approve-target` | Approve KPI evaluation target | Setujui fase target evaluasi KPI (status → PLAN_APPROVED). |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/calculate-scoring` | Calculate evaluation component scoring | Jalankan scoring engine: hitung skor tiap komponen dari data terkait (KPI, competency, dll) lalu simpan hasilnya. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/comments` | List comments by evaluation | Ambil semua komentar pada sebuah performance evaluation. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/complete` | Complete KPI evaluation | Selesaikan evaluasi KPI (status → COMPLETED) — hasil akhir terkunci. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/components` | List evaluation component scores | Ambil skor per komponen untuk sebuah evaluasi KPI (hasil scoring engine). |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluations/{id}/components/{component_id}` | Update evaluation component score | Isi skor komponen secara manual (mis. Work Program yang tidak bisa dihitung otomatis — wajib diisi reviewer). |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/details` | List evaluation details by evaluation ID | Retrieve all BSC perspective detail records for a specific performance evaluation, showing achievement per perspective. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/full` | Get KPI evaluation with details | Ambil evaluasi KPI lengkap: detail perspektif, target, progress, komentar, dan lampiran. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/logs` | List audit logs by evaluation | Ambil audit trail perubahan pada sebuah performance evaluation. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/program-items` | List evaluation program items | Ambil daftar program item (program kerja yang diajukan karyawan) untuk sebuah evaluasi KPI. |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/progress-summary` | Get KPI evaluation progress summary | Ringkasan progres realisasi KPI untuk sebuah evaluasi (nilai aktual terbaru & achievement per detail). |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/recalculate` | Recalculate KPI evaluation score | Hitung ulang skor evaluasi KPI dari nilai aktual & formula tiap indikator, lalu simpan hasilnya. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/reject` | Reject KPI evaluation | Tolak evaluasi KPI (status → REJECTED) untuk direvisi. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/reject-target` | Reject KPI evaluation target | Tolak fase target evaluasi KPI untuk direvisi (status → DRAFT). |
| `PUT` | `/api/v1/tenant/performance/kpi/evaluations/{id}/status` | Update evaluation status | Transition a performance evaluation through its workflow: DRAFT -> PLAN_SUBMITTED -> PLAN_APPROVED -> ACTUAL_SUBMITTED -> ACTUAL_APPROVED -> COMPLE... |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/submit` | Submit KPI evaluation | Ajukan evaluasi KPI (status → SUBMITTED) untuk direview atasan. |
| `POST` | `/api/v1/tenant/performance/kpi/evaluations/{id}/submit-target` | Submit KPI evaluation target | Ajukan fase target evaluasi KPI (status → PLAN_SUBMITTED). |
| `GET` | `/api/v1/tenant/performance/kpi/evaluations/{id}/targets` | List performance targets by evaluation ID | Retrieve all KPI targets for a specific performance evaluation, showing planned vs actual achievement for each indicator. |
| `POST` | `/api/v1/tenant/performance/kpi/indicators` | Create KPI indicator | Create a new KPI indicator linked to a template and BSC perspective. Defines target value, weight, and measurement unit. |
| `GET` | `/api/v1/tenant/performance/kpi/indicators` | List KPI indicators | Retrieve a paginated list of KPI indicators, optionally filtered by template or perspective. |
| `GET` | `/api/v1/tenant/performance/kpi/indicators/{id}` | Get KPI indicator by ID | Retrieve a specific KPI indicator by its unique ID, including target value and measurement settings. |
| `PUT` | `/api/v1/tenant/performance/kpi/indicators/{id}` | Update KPI indicator | Update a KPI indicator's title, weight, target value, or measurement unit. Only provided fields will be updated. |
| `DELETE` | `/api/v1/tenant/performance/kpi/indicators/{id}` | Delete KPI indicator | Permanently delete a KPI indicator from its template. |
| `GET` | `/api/v1/tenant/performance/kpi/my-context` | Get my KPI context | Resolusi user login ke employee & Organization saat ini, lalu mengembalikan template PUBLISHED untuk organisasi tersebut (dipakai halaman self-asse... |
| `POST` | `/api/v1/tenant/performance/kpi/organization-components` | Upsert organization component weight | Atur/upsert konfigurasi bobot komponen scoring untuk sebuah organisasi (enable/disable, weight, sort order). |
| `DELETE` | `/api/v1/tenant/performance/kpi/organization-components/{id}` | Delete organization component | Hapus konfigurasi komponen scoring dari organisasi. |
| `GET` | `/api/v1/tenant/performance/kpi/organizations/{organization_id}/components` | List organization components | Ambil daftar komponen scoring yang diaktifkan untuk sebuah organisasi beserta bobotnya. |
| `POST` | `/api/v1/tenant/performance/kpi/periods/{period_id}/recalculate-scoring` | Recalculate period scoring (batch) | Jalankan batch recalculation scoring untuk seluruh evaluasi dalam sebuah periode KPI (bottom-up subordinate scoring / akhir periode) lalu simpan ha... |
| `POST` | `/api/v1/tenant/performance/kpi/perspectives` | Create BSC perspective | Create a new Balanced Scorecard perspective (e.g. Financial, Customer, Internal Process, Learning & Growth). |
| `GET` | `/api/v1/tenant/performance/kpi/perspectives` | List BSC perspectives | Retrieve a paginated list of BSC perspectives used in performance templates. Ordered by sort_order by default. |
| `GET` | `/api/v1/tenant/performance/kpi/perspectives/{id}` | Get BSC perspective by ID | Retrieve a specific BSC perspective by its unique ID. |
| `PUT` | `/api/v1/tenant/performance/kpi/perspectives/{id}` | Update BSC perspective | Update a BSC perspective's name, description, or sort order. Only provided fields will be updated. |
| `DELETE` | `/api/v1/tenant/performance/kpi/perspectives/{id}` | Delete BSC perspective | Permanently delete a BSC perspective from the system. |
| `POST` | `/api/v1/tenant/performance/kpi/program-items` | Create KPI program item | Buat program item (program kerja yang diajukan karyawan sendiri) pada evaluasi KPI — tanpa template HR. |
| `DELETE` | `/api/v1/tenant/performance/kpi/program-items/{id}` | Delete program item | Hapus sebuah program item dari evaluasi KPI. |
| `PUT` | `/api/v1/tenant/performance/kpi/program-items/{id}/actual` | Update program item actual | Input nilai aktual program item — achievement & score dihitung ulang. |
| `PUT` | `/api/v1/tenant/performance/kpi/program-items/{id}/target` | Update program item target | Perbarui judul, formula, atau nilai target sebuah program item. |
| `POST` | `/api/v1/tenant/performance/kpi/progress` | Create performance progress | Catat progres realisasi KPI untuk satu evaluation detail (nilai aktual per tanggal). |
| `GET` | `/api/v1/tenant/performance/kpi/progress/{id}` | Get performance progress by ID | Ambil satu catatan progres KPI berdasarkan ID. |
| `PUT` | `/api/v1/tenant/performance/kpi/progress/{id}` | Update performance progress | Perbarui tanggal, nilai aktual, achievement, atau catatan progres KPI. |
| `DELETE` | `/api/v1/tenant/performance/kpi/progress/{id}` | Delete performance progress | Hapus satu catatan progres KPI. |
| `POST` | `/api/v1/tenant/performance/kpi/targets` | Create performance target | Add an individual KPI target to a performance evaluation, setting the target value and weight for measurement. |
| `PUT` | `/api/v1/tenant/performance/kpi/targets/{id}` | Update performance target | Update a KPI target's planned value, actual achievement, or weight. Setting actual_value triggers automatic achievement percentage calculation. |
| `DELETE` | `/api/v1/tenant/performance/kpi/targets/{id}` | Delete performance target | Permanently delete a KPI target from the evaluation. |
| `POST` | `/api/v1/tenant/performance/kpi/templates` | Create KPI template | Create a new KPI template for an organization. Templates define the structure of performance evaluations including indicators from BSC perspectives. |
| `GET` | `/api/v1/tenant/performance/kpi/templates` | List KPI templates | Retrieve a paginated list of KPI templates, optionally filtered by organization. |
| `GET` | `/api/v1/tenant/performance/kpi/templates/organization-scope` | List KPI template organization scope | Ambil daftar organisasi yang boleh dipilih saat membuat/mengedit KPI template — hanya organisasi turunan dari organisasi user (hierarki org), organ... |
| `GET` | `/api/v1/tenant/performance/kpi/templates/{id}` | Get KPI template by ID | Retrieve a specific KPI template by its unique ID, including associated indicators. |
| `PUT` | `/api/v1/tenant/performance/kpi/templates/{id}` | Update KPI template | Update a KPI template's name, description, or status. Status can be transitioned between DRAFT, PUBLISHED, and ARCHIVED. |
| `DELETE` | `/api/v1/tenant/performance/kpi/templates/{id}` | Delete KPI template | Permanently delete a KPI template. Indicators linked to this template may also be removed. |
| `GET` | `/api/v1/tenant/performance/logs` | List performance audit logs | Ambil daftar audit trail perubahan data performance dengan pagination. |
| `GET` | `/api/v1/tenant/performance/logs/{id}` | Get performance log by ID | Ambil satu audit trail berdasarkan ID. |
| `POST` | `/api/v1/tenant/performance/okr/attachments` | Create OKR attachment | Lampirkan file bukti/dokumen pendukung pada evaluation detail OKR. |
| `DELETE` | `/api/v1/tenant/performance/okr/attachments/{id}` | Delete OKR attachment | Hapus satu lampiran OKR. |
| `POST` | `/api/v1/tenant/performance/okr/comments` | Create OKR comment | Tambahkan komentar/review pada evaluasi OKR (mendukung reply via parent_id). |
| `PUT` | `/api/v1/tenant/performance/okr/comments/{id}` | Update OKR comment | Perbarui isi komentar OKR. |
| `DELETE` | `/api/v1/tenant/performance/okr/comments/{id}` | Delete OKR comment | Hapus komentar OKR. |
| `GET` | `/api/v1/tenant/performance/okr/dashboard/hr` | Get OKR HR dashboard | Dashboard OKR untuk HR: total evaluasi, sebaran status, skor & achievement rata-rata, distribusi rating. |
| `PUT` | `/api/v1/tenant/performance/okr/evaluation-details/{id}` | Update OKR evaluation detail actual | Input nilai aktual satu evaluation detail OKR. |
| `GET` | `/api/v1/tenant/performance/okr/evaluation-details/{id}/attachments` | List OKR attachments by evaluation detail | Ambil daftar lampiran bukti pada satu evaluation detail OKR. |
| `GET` | `/api/v1/tenant/performance/okr/evaluation-details/{id}/progress` | List OKR progress by evaluation detail | Ambil riwayat progres (check-in) untuk satu evaluation detail OKR. |
| `DELETE` | `/api/v1/tenant/performance/okr/evaluation-key-results/{id}` | Delete OKR evaluation key result | Hapus Key Result yang diajukan karyawan, hanya saat evaluasi masih DRAFT. |
| `PUT` | `/api/v1/tenant/performance/okr/evaluation-key-results/{id}/target` | Update OKR evaluation key result target | Perbarui target Key Result yang diajukan karyawan (judul, target, unit, formula, bobot) sebelum disubmit. |
| `POST` | `/api/v1/tenant/performance/okr/evaluations` | Create OKR evaluation with snapshot | Buat evaluasi OKR employee dari template (objective & key results di-snapshot ke evaluation details). |
| `GET` | `/api/v1/tenant/performance/okr/evaluations` | List OKR evaluations | Ambil daftar evaluasi OKR dengan pagination, filter per employee/organisasi/periode/status. |
| `GET` | `/api/v1/tenant/performance/okr/evaluations/{id}` | Get OKR evaluation by ID | Ambil detail evaluasi OKR termasuk details, skor, dan status workflow. |
| `PUT` | `/api/v1/tenant/performance/okr/evaluations/{id}` | Update OKR evaluation | Perbarui status atau catatan reviewer evaluasi OKR. |
| `DELETE` | `/api/v1/tenant/performance/okr/evaluations/{id}` | Delete OKR evaluation | Hapus evaluasi OKR (hanya status DRAFT). |
| `PUT` | `/api/v1/tenant/performance/okr/evaluations/{id}/actuals` | Bulk update OKR evaluation actuals | Input nilai aktual sekaligus untuk banyak evaluation detail OKR. |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/approve` | Approve OKR evaluation | Setujui evaluasi OKR (status → APPROVED). |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/approve-key-results` | Approve OKR key results | Setujui proposal Key Results (status → KR_APPROVED, "OKR Active") — dapat diresolusi otomatis oleh approval flow jika dikonfigurasi. |
| `GET` | `/api/v1/tenant/performance/okr/evaluations/{id}/comments` | List OKR comments by evaluation | Ambil semua komentar (dengan replies) pada sebuah evaluasi OKR. |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/complete` | Complete OKR evaluation | Selesaikan evaluasi OKR (status → COMPLETED). |
| `GET` | `/api/v1/tenant/performance/okr/evaluations/{id}/details` | Get OKR evaluation with details | Ambil evaluasi OKR lengkap dengan seluruh evaluation details. |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/key-results` | Create employee-proposed OKR evaluation key result | Tambahkan Key Result yang diajukan employee di bawah Objective hasil snapshot, hanya saat evaluasi berstatus DRAFT (fase KR proposal). |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/recalculate` | Recalculate OKR evaluation score | Hitung ulang skor & achievement evaluasi OKR dari nilai aktual dan formula tiap key result. |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/reject` | Reject OKR evaluation | Tolak evaluasi OKR (status → REJECTED). |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/reject-key-results` | Reject OKR key results | Tolak proposal Key Results kembali ke DRAFT (status → DRAFT) dengan catatan penolakan, agar karyawan dapat memperbaiki proposal. |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/submit` | Submit OKR evaluation | Ajukan evaluasi OKR (status → SUBMITTED). |
| `POST` | `/api/v1/tenant/performance/okr/evaluations/{id}/submit-key-results` | Submit OKR key results for approval | Ajukan proposal Key Results karyawan untuk persetujuan (status → KR_SUBMITTED) melalui approval flow modul performance. |
| `POST` | `/api/v1/tenant/performance/okr/key-results` | Create OKR key result | Buat key result terukur di dalam sebuah objective (target, unit, formula, bobot). |
| `GET` | `/api/v1/tenant/performance/okr/key-results/{id}` | Get OKR key result by ID | Ambil detail key result OKR. |
| `PUT` | `/api/v1/tenant/performance/okr/key-results/{id}` | Update OKR key result | Perbarui target, unit, formula, bobot, atau status key result. |
| `DELETE` | `/api/v1/tenant/performance/okr/key-results/{id}` | Delete OKR key result | Hapus key result dari objective. |
| `GET` | `/api/v1/tenant/performance/okr/my-context` | Get my OKR context | Konteks OKR user saat ini: apakah punya posisi, employee/organization terkait, dan daftar template OKR yang tersedia (dipakai self-assessment). |
| `POST` | `/api/v1/tenant/performance/okr/objectives` | Create OKR objective | Buat objective baru di dalam sebuah template OKR. |
| `GET` | `/api/v1/tenant/performance/okr/objectives/{id}` | Get OKR objective by ID | Ambil detail objective OKR beserta key results. |
| `PUT` | `/api/v1/tenant/performance/okr/objectives/{id}` | Update OKR objective | Perbarui judul, deskripsi, bobot, atau urutan objective. |
| `DELETE` | `/api/v1/tenant/performance/okr/objectives/{id}` | Delete OKR objective | Hapus objective beserta key results terkait. |
| `GET` | `/api/v1/tenant/performance/okr/objectives/{id}/key-results` | List OKR key results by objective | Ambil daftar key results dalam sebuah objective (non-paginated). |
| `POST` | `/api/v1/tenant/performance/okr/progress` | Create OKR progress | Catat progres (check-in) nilai aktual untuk satu evaluation detail OKR. |
| `GET` | `/api/v1/tenant/performance/okr/progress/{id}` | Get OKR progress by ID | Ambil satu catatan progres OKR. |
| `PUT` | `/api/v1/tenant/performance/okr/progress/{id}` | Update OKR progress | Perbarui tanggal, nilai aktual, atau catatan progres OKR. |
| `DELETE` | `/api/v1/tenant/performance/okr/progress/{id}` | Delete OKR progress | Hapus satu catatan progres OKR. |
| `POST` | `/api/v1/tenant/performance/okr/templates` | Create OKR template | Buat template OKR untuk sebuah organisasi & periode (dengan objective & key results). |
| `GET` | `/api/v1/tenant/performance/okr/templates` | List OKR templates | Ambil daftar template OKR dengan pagination, filter per organisasi/periode. |
| `GET` | `/api/v1/tenant/performance/okr/templates/objective-scope` | Get OKR objective creation scope | Resolusi scope pembuatan objective cascading: apakah employee saat ini boleh membuat objective untuk organisasi bawahan, dan organisasi mana yang m... |
| `GET` | `/api/v1/tenant/performance/okr/templates/{id}` | Get OKR template by ID | Ambil detail template OKR termasuk daftar objective & key results. |
| `PUT` | `/api/v1/tenant/performance/okr/templates/{id}` | Update OKR template | Perbarui nama, periode, status, atau tanggal efektif template OKR. |
| `DELETE` | `/api/v1/tenant/performance/okr/templates/{id}` | Delete OKR template | Hapus template OKR beserta objective & key results terkait. |
| `POST` | `/api/v1/tenant/performance/okr/templates/{id}/duplicate` | Duplicate OKR template | Duplikat template OKR (beserta objective & key results) sebagai template baru. |
| `GET` | `/api/v1/tenant/performance/okr/templates/{id}/objectives` | List OKR objectives by template | Ambil daftar objective dalam sebuah template OKR (non-paginated). |
| `POST` | `/api/v1/tenant/performance/periods` | Create performance period | Create a new performance evaluation period (e.g. Q1 2026). Period type must be one of: MONTHLY, QUARTERLY, SEMESTER, ANNUAL. |
| `GET` | `/api/v1/tenant/performance/periods` | List performance periods | Retrieve a paginated list of performance periods, optionally filtered by year or status. |
| `GET` | `/api/v1/tenant/performance/periods/{id}` | Get performance period by ID | Retrieve detailed information about a specific performance period by its unique UUID. |
| `PUT` | `/api/v1/tenant/performance/periods/{id}` | Update performance period | Update an existing performance period's details. Only provided fields will be updated. |
| `DELETE` | `/api/v1/tenant/performance/periods/{id}` | Delete performance period | Permanently delete a performance period by its unique ID. |
| `POST` | `/api/v1/tenant/performance/ratings` | Create performance rating | Buat skala rating penilaian (mis. A=90-100, B=80-89) untuk konversi skor akhir. |
| `GET` | `/api/v1/tenant/performance/ratings` | List performance ratings | Ambil daftar rating skala penilaian dengan pagination. |
| `GET` | `/api/v1/tenant/performance/ratings/{id}` | Get performance rating by ID | Ambil satu rating skala penilaian. |
| `PUT` | `/api/v1/tenant/performance/ratings/{id}` | Update performance rating | Perbarui rating skala penilaian. |
| `DELETE` | `/api/v1/tenant/performance/ratings/{id}` | Delete performance rating | Hapus satu rating skala penilaian. |

### Tenant: Settings
**Description:** Settings & Master Data Reference -- manage zones, provinces, regencies, districts, villages, educations, religions, marital statuses, relationship types, banks, employment statuses, nationalities, job families, and salary grades. CRUD operations for all tenant reference data, education majors.
**Endpoints:** 107 | **Paths:** 44
**Methods:** DELETE=21 GET=44 POST=21 PUT=21

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/settings/banks` | List all Banks | Retrieve a paginated list of Banks. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/banks` | Create a new bank | Create a new banks record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/banks/{id}` | Get bank by ID | Get detailed information about a specific bank by its ID. |
| `PUT` | `/api/v1/tenant/settings/banks/{id}` | Update bank | Update a bank record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/banks/{id}` | Delete bank | Soft-delete a bank record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/company-holidays` | List all Company Holidays | Retrieve a paginated list of company holidays (reference data for attendance/leave). |
| `POST` | `/api/v1/tenant/settings/company-holidays` | Create a new company holiday | Create a new company holiday entry (e.g. Tahun Baru, Idul Fitri, etc.). |
| `GET` | `/api/v1/tenant/settings/company-holidays/{id}` | Get a company holiday by ID |  |
| `PUT` | `/api/v1/tenant/settings/company-holidays/{id}` | Update a company holiday |  |
| `DELETE` | `/api/v1/tenant/settings/company-holidays/{id}` | Delete a company holiday |  |
| `GET` | `/api/v1/tenant/settings/competencies` | List all Competencies | List competency master data with pagination and optional search filter (matches name, field, cluster, or definition). |
| `POST` | `/api/v1/tenant/settings/competencies` | Create setting competency | Buat data kompetensi di master data settings. |
| `GET` | `/api/v1/tenant/settings/competencies/{id}` | Get setting competency by ID | Ambil satu data kompetensi di settings. |
| `PUT` | `/api/v1/tenant/settings/competencies/{id}` | Update setting competency | Perbarui data kompetensi di settings. |
| `DELETE` | `/api/v1/tenant/settings/competencies/{id}` | Delete setting competency | Hapus data kompetensi di settings. |
| `GET` | `/api/v1/tenant/settings/districts` | List all Districts | Retrieve a paginated list of Districts. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/districts` | Create a new district | Create a new districts record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/districts/{id}` | Get district by ID | Get detailed information about a specific district by its ID. |
| `PUT` | `/api/v1/tenant/settings/districts/{id}` | Update district | Update a district record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/districts/{id}` | Delete district | Soft-delete a district record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/education-majors` | List all Education Majors | Retrieve a paginated list of education majors (fields of study). Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/education-majors` | Create a new education major | Create a new education major record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/education-majors/{id}` | Get education major by ID | Get detailed information about a specific education major by its ID. |
| `PUT` | `/api/v1/tenant/settings/education-majors/{id}` | Update education major | Update an education major record's details including code, name, and sort order. |
| `DELETE` | `/api/v1/tenant/settings/education-majors/{id}` | Delete education major | Soft-delete an education major record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/educations` | List all Educations | Retrieve a paginated list of Educations. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/educations` | Create a new education | Create a new educations record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/educations/{id}` | Get education by ID | Get detailed information about a specific education by its ID. |
| `PUT` | `/api/v1/tenant/settings/educations/{id}` | Update education | Update a education record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/educations/{id}` | Delete education | Soft-delete a education record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/employment-statuses` | List all Employment Statuses | Retrieve a paginated list of Employment Statuses. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/employment-statuses` | Create a new employmentstatus | Create a new employment statuses record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/employment-statuses/{id}` | Get employmentstatus by ID | Get detailed information about a specific employmentstatus by its ID. |
| `PUT` | `/api/v1/tenant/settings/employment-statuses/{id}` | Update employmentstatus | Update a employmentstatus record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/employment-statuses/{id}` | Delete employmentstatus | Soft-delete a employmentstatus record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/gradings` | List all Gradings | Retrieve a paginated list of Gradings (job grading levels). |
| `POST` | `/api/v1/tenant/settings/gradings` | Create a new grading | Create a new grading record (job grading level). |
| `GET` | `/api/v1/tenant/settings/gradings/{id}` | Get grading by ID | Get detailed information about a specific grading by its ID. |
| `PUT` | `/api/v1/tenant/settings/gradings/{id}` | Update grading | Update a grading record's code, name, or sort order. |
| `DELETE` | `/api/v1/tenant/settings/gradings/{id}` | Delete grading | Soft-delete a grading record. |
| `GET` | `/api/v1/tenant/settings/insurances` | List all Insurances | Retrieve a paginated list of Insurance types (reference data for employee insurance). |
| `POST` | `/api/v1/tenant/settings/insurances` | Create a new insurance type | Create a new insurance type entry (e.g. BPJS Kesehatan, BPJS Ketenagakerjaan, etc.). |
| `GET` | `/api/v1/tenant/settings/insurances/{id}` | Get insurance type by ID | Get detailed information about a specific insurance type by its ID. |
| `PUT` | `/api/v1/tenant/settings/insurances/{id}` | Update insurance type | Update an existing insurance type's code, name, or sort order. |
| `DELETE` | `/api/v1/tenant/settings/insurances/{id}` | Delete insurance type | Soft-delete an insurance type record. |
| `GET` | `/api/v1/tenant/settings/job-families` | List all Job Families | Retrieve a paginated list of Job Families. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/job-families` | Create a new jobfamily | Create a new job families record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/job-families/{id}` | Get jobfamily by ID | Get detailed information about a specific jobfamily by its ID. |
| `PUT` | `/api/v1/tenant/settings/job-families/{id}` | Update jobfamily | Update a jobfamily record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/job-families/{id}` | Delete jobfamily | Soft-delete a jobfamily record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/marital-statuses` | List all Marital Statuses | Retrieve a paginated list of Marital Statuses. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/marital-statuses` | Create a new maritalstatus | Create a new marital statuses record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/marital-statuses/{id}` | Get maritalstatus by ID | Get detailed information about a specific maritalstatus by its ID. |
| `PUT` | `/api/v1/tenant/settings/marital-statuses/{id}` | Update maritalstatus | Update a maritalstatus record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/marital-statuses/{id}` | Delete maritalstatus | Soft-delete a maritalstatus record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/nationalities` | List all Nationalities | Retrieve a paginated list of Nationalities. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/nationalities` | Create a new nationality | Create a new nationalities record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/nationalities/{id}` | Get nationality by ID | Get detailed information about a specific nationality by its ID. |
| `PUT` | `/api/v1/tenant/settings/nationalities/{id}` | Update nationality | Update a nationality record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/nationalities/{id}` | Delete nationality | Soft-delete a nationality record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/provinces` | List all Provinces | Retrieve a paginated list of Provinces. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/provinces` | Create a new province | Create a new provinces record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/provinces/{id}` | Get province by ID | Get detailed information about a specific province by its ID. |
| `PUT` | `/api/v1/tenant/settings/provinces/{id}` | Update province | Update a province record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/provinces/{id}` | Delete province | Soft-delete a province record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/ptkps` | List all PTKPs | Retrieve a paginated list of Penghasilan Tidak Kena Pajak (PTKP) records. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/ptkps` | Create a new PTKP | Create a new PTKP record with name, amount, and group. |
| `GET` | `/api/v1/tenant/settings/ptkps/{id}` | Get PTKP by ID | Get detailed information about a specific PTKP record by its ID. |
| `PUT` | `/api/v1/tenant/settings/ptkps/{id}` | Update PTKP | Update a PTKP record's name, amount, or group. |
| `DELETE` | `/api/v1/tenant/settings/ptkps/{id}` | Delete PTKP | Delete a PTKP record. |
| `GET` | `/api/v1/tenant/settings/regencies` | List all Regencies | Retrieve a paginated list of Regencies. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/regencies` | Create a new regency | Create a new regencies record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/regencies/{id}` | Get regency by ID | Get detailed information about a specific regency by its ID. |
| `PUT` | `/api/v1/tenant/settings/regencies/{id}` | Update regency | Update a regency record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/regencies/{id}` | Delete regency | Soft-delete a regency record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/relationship-types` | List all Relationship Types | Retrieve a paginated list of Relationship Types. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/relationship-types` | Create a new relationshiptype | Create a new relationship types record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/relationship-types/{id}` | Get relationshiptype by ID | Get detailed information about a specific relationshiptype by its ID. |
| `PUT` | `/api/v1/tenant/settings/relationship-types/{id}` | Update relationshiptype | Update a relationshiptype record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/relationship-types/{id}` | Delete relationshiptype | Soft-delete a relationshiptype record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/religions` | List all Religions | Retrieve a paginated list of Religions. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/religions` | Create a new religion | Create a new religions record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/religions/{id}` | Get religion by ID | Get detailed information about a specific religion by its ID. |
| `PUT` | `/api/v1/tenant/settings/religions/{id}` | Update religion | Update a religion record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/religions/{id}` | Delete religion | Soft-delete a religion record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/salary-grades` | List all Salary Grades | Retrieve a paginated list of Salary Grades. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/salary-grades` | Create a new salarygrade | Create a new salary grades record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/salary-grades/{id}` | Get salarygrade by ID | Get detailed information about a specific salarygrade by its ID. |
| `PUT` | `/api/v1/tenant/settings/salary-grades/{id}` | Update salarygrade | Update a salarygrade record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/salary-grades/{id}` | Delete salarygrade | Soft-delete a salarygrade record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/ters` | List all TERs | Retrieve a paginated list of Tarif Efektif Rata-rata (TER) records used for PPh21 calculation. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/ters` | Create a new TER | Create a new TER record with group, bruto range, and rate. |
| `GET` | `/api/v1/tenant/settings/ters/{id}` | Get TER by ID | Get detailed information about a specific TER record by its ID. |
| `PUT` | `/api/v1/tenant/settings/ters/{id}` | Update TER | Update a TER record's group, bruto range, or rate. |
| `DELETE` | `/api/v1/tenant/settings/ters/{id}` | Delete TER | Delete a TER record. |
| `GET` | `/api/v1/tenant/settings/villages` | List all Villages | Retrieve a paginated list of Villages. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/villages` | Create a new village | Create a new villages record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/villages/search` | Search villages by name | Search villages by name for autocomplete. Returns matching villages with their province, regency, and district hierarchy names. |
| `GET` | `/api/v1/tenant/settings/villages/{id}` | Get village by ID | Get detailed information about a specific village by its ID. |
| `PUT` | `/api/v1/tenant/settings/villages/{id}` | Update village | Update a village record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/villages/{id}` | Delete village | Soft-delete a village record. Sets the deleted_at timestamp and hides it from standard queries. |
| `GET` | `/api/v1/tenant/settings/villages/{id}/detail` | Get village detail with hierarchy | Get a village by its code including the full hierarchy: village name, district name, regency name, and province name. |
| `GET` | `/api/v1/tenant/settings/zones` | List all Zones | Retrieve a paginated list of Zones. Supports pagination parameters. |
| `POST` | `/api/v1/tenant/settings/zones` | Create a new zone | Create a new zones record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/settings/zones/{id}` | Get zone by ID | Get detailed information about a specific zone by its ID. |
| `PUT` | `/api/v1/tenant/settings/zones/{id}` | Update zone | Update a zone record's details including code, name, and other attributes. |
| `DELETE` | `/api/v1/tenant/settings/zones/{id}` | Delete zone | Soft-delete a zone record. Sets the deleted_at timestamp and hides it from standard queries. |

### Tenant: Job Management
**Description:** Job analysis management including titles, values, objectives, responsibilities, competencies, and scoring
**Endpoints:** 96 | **Paths:** 40
**Methods:** DELETE=18 GET=40 POST=18 PUT=20

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/job-management/assets` | List job assets with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/assets` | Create job asset | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/assets/{id}` | Get job asset by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/assets/{id}` | Update job asset | Update an existing assets record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/assets/{id}` | Delete job asset | Delete a assets record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/competency-groups` | List competency groups by organization | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/competency-groups` | Create competency group | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/competency-groups/{id}` | Get competency group by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/competency-groups/{id}` | Update competency group | Update an existing competency groups record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/competency-groups/{id}` | Delete competency group | Delete a competency groups record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/education-experiences` | List job education experiences with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/education-experiences` | Create job education experience | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/education-experiences/{id}` | Get job education experience by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/education-experiences/{id}` | Update job education experience | Update an existing education experiences record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/education-experiences/{id}` | Delete job education experience | Delete a education experiences record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/financials` | List job financials with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/financials` | Create job financial | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/financials/{id}` | Get job financial by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/financials/{id}` | Update job financial | Update an existing financials record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/financials/{id}` | Delete job financial | Delete a financials record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/hr-authorities` | List HR authorities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/hr-authorities` | Create HR authority | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/hr-authorities/{id}` | Get HR authority by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/hr-authorities/{id}` | Update HR authority | Update an existing hr authorities record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/hr-authorities/{id}` | Delete HR authority | Delete a hr authorities record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/identifications` | List job identifications with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/identifications` | Create a new job identification | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/identifications/{id}` | Get job identification by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/identifications/{id}` | Update job identification | Update an existing identifications record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/identifications/{id}` | Delete job identification | Delete a identifications record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/objectives` | List job objectives with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/objectives` | Create a new job objective | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/objectives/{id}` | Get job objective by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/objectives/{id}` | Update job objective | Update an existing objectives record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/objectives/{id}` | Delete job objective | Delete a objectives record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/operational-authorities` | List operational authorities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/operational-authorities` | Create operational authority | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/operational-authorities/{id}` | Get operational authority by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/operational-authorities/{id}` | Update operational authority | Update an existing operational authorities record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/operational-authorities/{id}` | Delete operational authority | Delete a operational authorities record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/potency-competencies` | List potency competencies with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/potency-competencies` | Create potency competency | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/potency-competencies/{id}` | Get potency competency by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/potency-competencies/{id}` | Update potency competency | Update an existing potency competencies record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/potency-competencies/{id}` | Delete potency competency | Delete a potency competencies record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/relationships` | List relationships with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/relationships` | Create job relationship | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/relationships/{id}` | Get relationship by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/relationships/{id}` | Update relationship | Update an existing relationships record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/relationships/{id}` | Delete relationship | Delete a relationships record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/relationships/{id}/details` | List relationship details | Detail banyak-per-relationship (migration 048). |
| `POST` | `/api/v1/tenant/job-management/relationships/{id}/details` | Create relationship detail | Tambah detail hubungan kerja (work relations + activity in connection). |
| `GET` | `/api/v1/tenant/job-management/relationships/{id}/details/{detailId}` | Get relationship detail by ID | Ambil satu detail hubungan kerja. |
| `PUT` | `/api/v1/tenant/job-management/relationships/{id}/details/{detailId}` | Update relationship detail | Update satu detail hubungan kerja. |
| `DELETE` | `/api/v1/tenant/job-management/relationships/{id}/details/{detailId}` | Delete relationship detail | Hapus satu detail hubungan kerja. |
| `GET` | `/api/v1/tenant/job-management/responsibilities` | List job responsibilities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/responsibilities` | Create a new job responsibility | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/responsibilities/{id}` | Get job responsibility by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/responsibilities/{id}` | Update job responsibility | Update an existing responsibilities record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/responsibilities/{id}` | Delete job responsibility | Delete a responsibilities record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/scores` | List job scores with pagination | Retrieve a paginated list of job management resources. |
| `GET` | `/api/v1/tenant/job-management/scores/org/{orgId}` | Get job score by organization | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/scores/org/{orgId}` | Upsert job score for organization | Hitung ulang skor otomatis (body kosong) lalu simpan. Menghasilkan components, sub_component_points, is_complete, completed_at. |
| `GET` | `/api/v1/tenant/job-management/subordinate-controls` | List subordinate controls with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/subordinate-controls` | Create subordinate control | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/subordinate-controls/{id}` | Get subordinate control by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/subordinate-controls/{id}` | Update subordinate control | Update an existing subordinate controls record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/subordinate-controls/{id}` | Delete subordinate control | Delete a subordinate controls record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/titles` | List job titles with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/titles` | Create a new job title | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/titles/{id}` | Get job title by ID with subs | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/titles/{id}` | Update job title | Update an existing titles record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/titles/{id}` | Delete job title | Delete a titles record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/titles/{titleId}/subs` | List subs under a job title | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/titles/{titleId}/subs` | Create a sub under a job title | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/titles/{titleId}/subs/{subId}` | Get job title sub by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/titles/{titleId}/subs/{subId}` | Update job title sub | Update an existing {subId} record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/titles/{titleId}/subs/{subId}` | Delete job title sub | Delete a {subId} record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/values` | List job values with pagination | Retrieve a paginated list of job management resources, optionally filtered by type. |
| `POST` | `/api/v1/tenant/job-management/values` | Create a new job value | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/values/clusters/{type}` | List cluster mapping for job value type | Mapping tipe technical/managerial → cluster kompetensi dari tabel job_management_value_clusters. |
| `PUT` | `/api/v1/tenant/job-management/values/clusters/{type}` | Update cluster mapping for job value type | Simpan mapping cluster untuk tipe tertentu (technical/managerial). |
| `GET` | `/api/v1/tenant/job-management/values/tree` | Get job values tree | Mengembalikan hierarki type_group → daftar tipe (label = description_group) → options per tipe (level + deskripsi) dengan urutan grup tetap. Dipaka... |
| `GET` | `/api/v1/tenant/job-management/values/{id}` | Get job value by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/values/{id}` | Update job value | Update an existing values record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/values/{id}` | Delete job value | Delete a values record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/working-activities` | List working activities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/working-activities` | Create working activity | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/working-activities/{id}` | Get working activity by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/working-activities/{id}` | Update working activity | Update an existing working activities record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/working-activities/{id}` | Delete working activity | Delete a working activities record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/job-management/working-risks` | List working risks with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/working-risks` | Create working risk | Create a new job management resource. |
| `GET` | `/api/v1/tenant/job-management/working-risks/{id}` | Get working risk by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/working-risks/{id}` | Update working risk | Update an existing working risks record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/job-management/working-risks/{id}` | Delete working risk | Delete a working risks record by its unique ID. This action may be reversible depending on system configuration. |

### Tenant: Workforce Intelligence & Strategic Planning
**Description:** Workforce Intelligence & Strategic Workforce Planning â€” strategic analytics layer for headcount planning, forecasting, gap analysis, KPI monitoring, workforce analytics (headcount, attendance, leave, overtime, payroll, performance, learning, recruitment, movement), capacity planning, cost analytics, risk monitoring, executive dashboards, scenario simulation, organization health scoring, and people analytics (training-vs-performance, overtime-vs-productivity, etc.). Read-only analytics module aggregating data from all operational HR modules.
**Endpoints:** 68 | **Paths:** 58
**Methods:** DELETE=3 GET=56 POST=5 PUT=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/attendance` | Attendance analytics dashboard | Analyze attendance metrics: average attendance rate, late rate, absentee rate, trend, and department breakdown. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/headcount` | Headcount analytics dashboard | Analyze workforce composition: total HC, active HC, distribution by department, employment type, gender, and headcount trend over time. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/learning` | Learning analytics dashboard | Analyze learning and development: completion rate, average score, total training hours, and breakdown by course. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/leave` | Leave analytics dashboard | Analyze leave utilization: average utilization rate, total days taken, breakdown by leave type and department. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/movement` | Movement analytics dashboard | Analyze employee movement: promotion and mutation counts with breakdown by department and movement type. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/overtime` | Overtime analytics dashboard | Analyze overtime patterns: average OT hours, total OT cost, department breakdown, and trend over time. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/payroll` | Payroll analytics dashboard | Analyze payroll metrics: total payroll cost, average salary, breakdown by department and grade, with trend analysis. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/performance` | Performance analytics dashboard | Analyze employee performance: average score, top performer percentage, department breakdown, and score distribution. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/recruitment` | Recruitment analytics dashboard | Analyze recruitment efficiency: time to hire, cost per hire, source breakdown, and recruitment funnel metrics. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/bottlenecks` | Bottleneck analysis | Identify capacity bottlenecks across departments. Flags departments with WARNING or CRITICAL utilization levels. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/dashboard` | Capacity dashboard | Get workforce capacity dashboard: overall utilization rate, available headcount, department breakdown, and bottleneck identification. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/forecast` | Capacity forecast | Get projected capacity forecast: future utilization, current vs projected needed headcount, and capacity gap analysis by department. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/utilization` | Resource utilization rate | Get workforce utilization rate data points. Shows how effectively the workforce is being utilized. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/budget-vs-actual` | Budget vs actual cost comparison | Get budget vs actual workforce cost comparison. Shows budget targets versus actual spending for cost control analysis. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/payroll` | Payroll cost breakdown | Get detailed payroll cost breakdown: total salary, allowances, deductions, BPJS contributions, by grade and component. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/per-department` | Cost by department | Get workforce cost broken down by department for comparison and budget allocation analysis. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/per-employee` | Cost per employee analysis | Get cost per employee metrics: average, median, minimum, and maximum cost per employee with breakdown by department and grade. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/summary` | Cost summary dashboard | Get workforce cost summary: total payroll, total benefit, total labor cost, cost per employee, department breakdown, and budget vs actual comparison. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/attrition-trend` | Executive attrition trend | Executive-level attrition rate trend analysis to monitor talent retention. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/capacity` | Executive capacity overview | Executive-level capacity overview: utilization rate, available HC, active departments, and bottleneck count. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/cost-trend` | Executive cost trend | Executive-level workforce cost trend analysis for budget planning and cost control. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/growth` | Executive workforce growth trend | Executive-level workforce growth trend analysis showing HC change over time. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/health-score` | Executive health score | Executive-level organization health score with span of control, manager ratio, internal hiring rate, succession coverage, and overall status (HEALT... |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/hiring-progress` | Hiring progress tracker | Track hiring progress: planned, in-progress, and completed hires for executive monitoring. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/risk-overview` | Executive risk overview | Executive-level risk overview: total risks, high and critical counts, by department and category. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/summary` | Executive workforce summary | Executive dashboard summary: total HC, HC growth, attrition rate, average cost, utilization rate, and overall health score in a single view. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/dashboard` | Organization health dashboard | Get organization health dashboard: composite health score with span of control, manager ratio, promotion rate, internal hiring rate, succession cov... |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/scores` | List health scores | Retrieve a paginated list of organization health scores. Filter by period and organization unit. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/scores/{id}` | Get health score by ID | Get detailed health score components for a specific score record. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/span-of-control` | Span of control analysis | Analyze span of control: average manager-to-report ratio, healthy range (3:1-7:1), status indicator, and department breakdown. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/succession` | Succession readiness analysis | Analyze succession readiness: roles with identified successors, total key roles, coverage rate, status (HEALTHY/WARNING/CRITICAL), and department b... |
| `GET` | `/api/v1/tenant/workforce-intelligence/kpi` | List KPIs | Retrieve a paginated list of workforce KPIs. Filter by period, dimension, or KPI code. Each KPI includes value, target, unit, and snapshot date. |
| `GET` | `/api/v1/tenant/workforce-intelligence/kpi/summary` | KPI summary (on-target vs below-target) | Get a summary of KPIs showing total count, on-target count, and below-target count for a given period. |
| `GET` | `/api/v1/tenant/workforce-intelligence/kpi/{code}` | Get KPI by code | Get a specific KPI by its unique code (e.g., attrition_rate, turnover_rate, retention_rate, etc.). |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/attendance-vs-performance` | Attendance vs performance correlation | Analyze correlation between attendance rate and employee performance to identify attendance-related performance risks. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/career-progression` | Career progression vs performance correlation | Analyze correlation between career advancement (promotions, movements) and employee performance trends. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/compensation-vs-turnover` | Compensation vs turnover correlation | Analyze correlation between compensation levels and turnover rates to evaluate compensation competitiveness and retention strategies. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/learning-effectiveness` | Learning effectiveness analysis | Analyze the effectiveness of learning programs by correlating training completion, scores, and subsequent performance improvements. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/overtime-vs-productivity` | Overtime vs productivity correlation | Analyze correlation between overtime hours and productivity metrics to identify optimal workload levels. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/source-vs-retention` | Recruitment source vs retention correlation | Analyze correlation between recruitment source (job board, referral, agency, etc.) and employee retention rates to optimize sourcing strategy. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/training-vs-performance` | Training vs performance correlation | Analyze correlation between training participation and employee performance scores. Helps evaluate training ROI and effectiveness. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/forecasts` | List workforce forecasts | Retrieve a paginated list of workforce forecasts. Filter by period, organization, and forecast type (DEMAND, SUPPLY, or HIRING). |
| `POST` | `/api/v1/tenant/workforce-intelligence/planning/forecasts` | Create workforce forecast | Create a new workforce forecast. Supports DEMAND (required headcount), SUPPLY (available headcount), and HIRING (gap to fill) forecast types. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/forecasts/{id}` | Get forecast by ID | Get detailed information about a specific workforce forecast including headcount projection and confidence level. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/planning/forecasts/{id}` | Update forecast | Update an existing workforce forecast's headcount and confidence level. |
| `DELETE` | `/api/v1/tenant/workforce-intelligence/planning/forecasts/{id}` | Delete forecast | Remove a workforce forecast record. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/gap-analysis` | Workforce gap analysis (supply vs demand) | Analyze workforce gaps by comparing supply (available HC) vs demand (required HC). Returns SURPLUS, SHORTAGE, or OPTIMAL status per department and ... |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/headcounts` | List headcount plans (planned vs actual HC per period) | Retrieve a paginated list of headcount plans. Filter by period (YYYY-MM) and/or organization unit. Each record shows planned headcount, actual head... |
| `POST` | `/api/v1/tenant/workforce-intelligence/planning/headcounts` | Create headcount plan | Create a new headcount plan for a specific period and organization. Records the planned headcount target for workforce planning purposes. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/headcounts/{id}` | Get headcount plan by ID | Get detailed information about a specific headcount plan including planned vs actual headcount and variance. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/planning/headcounts/{id}` | Update headcount plan | Update an existing headcount plan's planned headcount or snapshot date. |
| `DELETE` | `/api/v1/tenant/workforce-intelligence/planning/headcounts/{id}` | Delete headcount plan | Remove a headcount plan record from the system. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/projections` | Workforce projections (hiring, retirement, growth) | Get workforce projections including projected headcount, hiring needs, retirement counts, and growth rates for strategic workforce planning. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/contract-expiry` | Contract expiration risk analysis | Analyze contract expiry risk: upcoming contract/PKWT expirations, departments affected, and renewal recommendations. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/dashboard` | Risk dashboard | Get risk dashboard overview: total risks, high/critical risk counts, and list of active risk indicators with their levels and scores. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/high-absenteeism` | High absenteeism risk analysis | Analyze high absenteeism risk: current absenteeism rate vs threshold, affected departments, trend, and intervention recommendations. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/high-turnover` | High turnover risk analysis | Analyze high turnover risk: current turnover rate vs threshold, affected departments, trend, and recommended interventions. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/indicators` | List risk indicators | Retrieve a paginated list of risk indicators. Filter by period and/or risk level (LOW, MEDIUM, HIGH, CRITICAL). |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/indicators/{id}` | Get risk indicator by ID | Get detailed information about a specific risk indicator including current score, threshold, and recommendations. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/risk/indicators/{id}` | Update risk indicator | Update a risk indicator's level and/or recommendation. Used to acknowledge and document mitigation actions. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/retirement` | Retirement risk analysis | Analyze retirement risk: upcoming retirements, impacted roles, succession gaps, and recommendations for knowledge transfer. |
| `GET` | `/api/v1/tenant/workforce-intelligence/scenarios` | List scenarios | Retrieve a paginated list of saved simulation scenarios. Filter by scenario type (NEW_BRANCH, REORG, GROWTH, REDUCTION, RETIREMENT, BUDGET). |
| `POST` | `/api/v1/tenant/workforce-intelligence/scenarios` | Create scenario | Create a new scenario for workforce simulation. Supports NEW_BRANCH, REORG, GROWTH, REDUCTION, RETIREMENT, and BUDGET scenario types. Parameters ar... |
| `GET` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}` | Get scenario by ID | Get detailed information about a specific scenario including its parameters, results (if run), and status (DRAFT, RUNNING, COMPLETED). |
| `PUT` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}` | Update scenario | Update an existing scenario's name, description, and/or parameters. |
| `DELETE` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}` | Delete scenario | Soft-delete a scenario by ID. |
| `POST` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}/clone` | Clone scenario | Clone an existing scenario as a new DRAFT scenario. Useful for creating variations of a simulation without affecting the original. |
| `POST` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}/run` | Run scenario simulation | Execute a scenario simulation. Runs the scenario's parameters through the simulation engine and stores results in the scenario record. |

### Tenant: Payroll & Compensation Engine
**Endpoints:** 47 | **Paths:** 24
**Methods:** DELETE=8 GET=18 POST=12 PUT=9

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/payroll/bpjs-rate-components` | Create BPJS rate component | Create a new bpjs rate components record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/bpjs-rate-components/{id}` | Get BPJS rate component by ID | Retrieve a paginated list of bpjs rate components records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/bpjs-rate-components/{id}` | Update BPJS rate component | Update an existing bpjs rate components record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/bpjs-rate-components/{id}` | Delete BPJS rate component | Delete a bpjs rate components record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/payroll/bpjs-settings` | List BPJS settings | Retrieve a paginated list of bpjs settings records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/bpjs-settings` | Create BPJS setting | Create a new bpjs settings record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/bpjs-settings/{id}` | Get BPJS setting by ID | Retrieve a paginated list of bpjs settings records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/bpjs-settings/{id}` | Update BPJS setting | Update an existing bpjs settings record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/bpjs-settings/{id}` | Delete BPJS setting | Delete a bpjs settings record by its unique ID. This action may be reversible depending on system configuration. |
| `POST` | `/api/v1/tenant/payroll/employee-bank-profiles` | Create employee bank profile | Create a new employee bank profiles record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/employee-bank-profiles/{id}` | Get employee bank profile by ID | Retrieve a paginated list of employee bank profiles records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/employee-bank-profiles/{id}` | Update employee bank profile | Update an existing employee bank profiles record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/employee-bank-profiles/{id}` | Delete employee bank profile | Delete a employee bank profiles record by its unique ID. This action may be reversible depending on system configuration. |
| `POST` | `/api/v1/tenant/payroll/employee-bpjs-profiles` | Create employee BPJS profile | Create a new employee bpjs profiles record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/employee-bpjs-profiles/{id}` | Get employee BPJS profile by ID | Retrieve a paginated list of employee bpjs profiles records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/employee-bpjs-profiles/{id}` | Update employee BPJS profile | Update an existing employee bpjs profiles record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/employee-bpjs-profiles/{id}` | Delete employee BPJS profile | Delete a employee bpjs profiles record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/payroll/employee-payroll-profiles` | List employee payroll profiles | Retrieve a paginated list of employee payroll profiles records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/employee-payroll-profiles` | Create employee payroll profile | Create a new employee payroll profiles record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/employee-payroll-profiles/{id}` | Get employee payroll profile by ID | Retrieve a paginated list of employee payroll profiles records. Supports filtering, sorting, and pagination parameters. |
| `DELETE` | `/api/v1/tenant/payroll/employee-payroll-profiles/{id}` | Delete employee payroll profile | Delete a employee payroll profiles record by its unique ID. This action may be reversible depending on system configuration. |
| `POST` | `/api/v1/tenant/payroll/employee-tax-profiles` | Create employee tax profile | Create a new employee tax profiles record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/employee-tax-profiles/{id}` | Get employee tax profile by ID | Retrieve a paginated list of employee tax profiles records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/employee-tax-profiles/{id}` | Update employee tax profile | Update an existing employee tax profiles record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/employee-tax-profiles/{id}` | Delete employee tax profile | Delete a employee tax profiles record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/payroll/periods` | List payroll periods | Retrieve a paginated list of periods records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/periods` | Create payroll period | Create a new periods record. Validates required fields and returns the created resource with its assigned ID. |
| `PUT` | `/api/v1/tenant/payroll/periods/{id}` | Update payroll period | Update an existing periods record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `GET` | `/api/v1/tenant/payroll/pph21-ptkp-rates` | List PPh21 PTKP rates | Retrieve a paginated list of pph21 ptkp rates records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/pph21-ptkp-rates` | Create PPh21 PTKP rate | Create a new pph21 ptkp rates record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/pph21-settings` | List PPh21 settings | Retrieve a paginated list of pph21 settings records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/pph21-settings` | Create PPh21 setting | Create a new pph21 settings record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/pph21-settings/{id}` | Get PPh21 setting by ID | Retrieve a paginated list of pph21 settings records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/pph21-settings/{id}` | Update PPh21 setting | Update an existing pph21 settings record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/pph21-settings/{id}` | Delete PPh21 setting | Delete a pph21 settings record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/payroll/pph21-tax-brackets` | List PPh21 tax brackets | Retrieve a paginated list of pph21 tax brackets records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/pph21-tax-brackets` | Create PPh21 tax bracket | Create a new pph21 tax brackets record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/runs` | List payroll runs | Retrieve a paginated list of runs records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/runs` | Create payroll run | Create a new runs record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/runs/{id}` | Get payroll run by ID | Retrieve a paginated list of runs records. Supports filtering, sorting, and pagination parameters. |
| `GET` | `/api/v1/tenant/payroll/runs/{id}/approval` | Check payroll run approval status | Retrieve a paginated list of approval records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/runs/{id}/status` | Update payroll run status | Update an existing status record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `GET` | `/api/v1/tenant/payroll/salary-components` | List salary components with pagination | Retrieve a paginated list of salary components records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/salary-components` | Create salary component | Create a new salary components record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/payroll/salary-components/{id}` | Get salary component by ID | Retrieve a paginated list of salary components records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/payroll/salary-components/{id}` | Update salary component | Update an existing salary components record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/payroll/salary-components/{id}` | Delete salary component | Delete a salary components record by its unique ID. This action may be reversible depending on system configuration. |

### Tenant: Employees
**Description:** Employee management with 8 sub-modules (addresses, emergency contacts, families, educations, experiences, documents, insurances, employments)
**Endpoints:** 36 | **Paths:** 23
**Methods:** DELETE=11 GET=2 POST=11 PUT=12

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/employees` | List employees with pagination | Retrieve a paginated list of employees records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/employees` | Create a new employee | Create a new employees record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/employees/{id}` | Get employee by ID with all sub-modules | Retrieve a paginated list of employees records. Supports filtering, sorting, and pagination parameters. |
| `PUT` | `/api/v1/tenant/employees/{id}` | Update employee | Update an existing employees record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/employees/{id}` | Delete employee (hard delete) | Delete a employees record by its unique ID. This action may be reversible depending on system configuration. |
| `POST` | `/api/v1/tenant/employees/{id}/addresses` | Add employee address | Add a new address for an employee. Supports address types: MAIN (primary residence) and DOMICILE (current stay). Includes full address, RT/RW, subd... |
| `PUT` | `/api/v1/tenant/employees/{id}/addresses/{addressId}` | Update employee address | Update an existing employee address record. Can modify address type, full address details, RT/RW, administrative region, postal code, or coordinates. |
| `DELETE` | `/api/v1/tenant/employees/{id}/addresses/{addressId}` | Delete employee address | Remove an employee's address record from the system. |
| `POST` | `/api/v1/tenant/employees/{id}/banks` | Create employee bank account | Create a new bank account record for an employee. References a Bank from settings/banks. |
| `PUT` | `/api/v1/tenant/employees/{id}/banks/{bankId}` | Update employee bank account | Update an existing employee bank account record. |
| `DELETE` | `/api/v1/tenant/employees/{id}/banks/{bankId}` | Delete employee bank account | Delete an employee bank account record. |
| `POST` | `/api/v1/tenant/employees/{id}/documents` | Upload employee document | Upload and attach a document to an employee's profile. Supports document types such as ID card (KTP), NPWP, BPJS, diplomas, certificates, and other... |
| `POST` | `/api/v1/tenant/employees/{id}/documents/upload` | Upload a document file for employee | Upload a document file for an employee. Creates a document record with the uploaded file path. Supports PDF, DOC, DOCX, XLS, XLSX, JPG, PNG, GIF, T... |
| `PUT` | `/api/v1/tenant/employees/{id}/documents/{documentId}` | Update document metadata | Update an employee document's metadata such as document name, type, or description. |
| `DELETE` | `/api/v1/tenant/employees/{id}/documents/{documentId}` | Delete employee document | Remove a document from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/documents/{documentId}/upload` | Replace document file | Replace a document file or update document metadata. If a file is provided, it replaces the existing file on disk. If no file is provided, only nam... |
| `POST` | `/api/v1/tenant/employees/{id}/educations` | Add education record | Add an educational background record for an employee. Includes education level, institution name, major/field of study, graduation year, GPA, and d... |
| `PUT` | `/api/v1/tenant/employees/{id}/educations/{educationId}` | Update education record | Update an employee's educational background record including institution, degree, graduation date, and academic achievements. |
| `DELETE` | `/api/v1/tenant/employees/{id}/educations/{educationId}` | Delete education record | Remove an educational background record from the employee's profile. |
| `POST` | `/api/v1/tenant/employees/{id}/emergency-contacts` | Add emergency contact | Register an emergency contact person for an employee. Includes name, relationship, phone number, and alternative contact information. |
| `PUT` | `/api/v1/tenant/employees/{id}/emergency-contacts/{contactId}` | Update emergency contact | Update an employee's emergency contact details such as name, relationship, or phone number. |
| `DELETE` | `/api/v1/tenant/employees/{id}/emergency-contacts/{contactId}` | Delete emergency contact | Remove an emergency contact record from the employee's profile. |
| `POST` | `/api/v1/tenant/employees/{id}/employments` | Add employment record | Add an employment assignment for an employee. Associates the employee with an organization unit, position, employment status, and includes decision... |
| `PUT` | `/api/v1/tenant/employees/{id}/employments/{employmentId}` | Update employment record | Update an employee's employment assignment including organization, position, status, and decision letter information. |
| `DELETE` | `/api/v1/tenant/employees/{id}/employments/{employmentId}` | Delete employment record | Remove an employment assignment from the employee's profile. |
| `POST` | `/api/v1/tenant/employees/{id}/experiences` | Add work experience | Add a work experience record for an employee. Includes company name, position, start and end dates, job description, and reason for leaving. |
| `PUT` | `/api/v1/tenant/employees/{id}/experiences/{experienceId}` | Update work experience | Update an employee's work experience record including company, position, employment period, and responsibilities. |
| `DELETE` | `/api/v1/tenant/employees/{id}/experiences/{experienceId}` | Delete work experience | Remove a work experience record from the employee's profile. |
| `POST` | `/api/v1/tenant/employees/{id}/families` | Add family member | Add a family member record for an employee. Includes name, relationship (spouse, child, parent, sibling), date of birth, occupation, and dependency... |
| `PUT` | `/api/v1/tenant/employees/{id}/families/{familyId}` | Update family member | Update an employee's family member details including relationship, personal information, and tax dependency status. |
| `DELETE` | `/api/v1/tenant/employees/{id}/families/{familyId}` | Delete family member | Remove a family member record from the employee's profile. |
| `POST` | `/api/v1/tenant/employees/{id}/insurances` | Add insurance (BPJS) | Register an insurance record for an employee. Typically used for BPJS Kesehatan (health) and BPJS Ketenagakerjaan (employment) social security prog... |
| `PUT` | `/api/v1/tenant/employees/{id}/insurances/{insuranceId}` | Update insurance record | Update an employee's insurance record including BPJS participation number, coverage type, and contribution details. |
| `DELETE` | `/api/v1/tenant/employees/{id}/insurances/{insuranceId}` | Delete insurance record | Remove an insurance record from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/photo` | Upload employee profile photo | Upload a profile photo for an employee. Accepts JPG, PNG, GIF, WebP files up to 2MB. The photo is stored on the server and the employee's profile_p... |
| `DELETE` | `/api/v1/tenant/employees/{id}/photo` | Delete employee profile photo | Remove the profile photo from an employee. Deletes the file from the server and clears the profile_picture field. |

### Tenant: Competency Management
**Description:** Competency management including master competencies, values, events, targets, scores, and score details
**Endpoints:** 35 | **Paths:** 15
**Methods:** DELETE=7 GET=14 POST=7 PUT=7

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/competency/competence-values` | List competence values (legacy) | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/competence-values` | Create competence value (legacy) | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/competence-values/{id}` | Get competence value by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/competence-values/{id}` | Update competence value | Update an existing competence values record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/competence-values/{id}` | Delete competence value | Delete a competence values record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/competency/competencies` | List competencies | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/competencies` | Create a new competency | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/competencies/{id}` | Get competency by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/competencies/{id}` | Update competency | Update an existing competencies record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/competencies/{id}` | Delete competency | Delete a competencies record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/competency/event-targets` | List competency event targets | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/event-targets` | Create competency event target | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/event-targets/{id}` | Get event target by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/event-targets/{id}` | Update event target | Update an existing event targets record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/event-targets/{id}` | Delete event target | Delete a event targets record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/competency/events` | List competency events | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/events` | Create competency event | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/events/{id}` | Get competency event by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/events/{id}` | Update competency event | Update an existing events record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/events/{id}` | Delete competency event | Delete a events record by its unique ID. This action may be reversible depending on system configuration. |
| `POST` | `/api/v1/tenant/competency/score-details` | Create score detail | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/score-details/{id}` | Get score detail by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/score-details/{id}` | Update score detail | Update an existing score details record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/score-details/{id}` | Delete score detail | Delete a score details record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/competency/scores` | List competency scores | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/scores` | Create competency score | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/scores/{id}` | Get competency score by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/scores/{id}` | Update competency score | Update an existing scores record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/scores/{id}` | Delete competency score | Delete a scores record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/competency/scores/{scoreId}/details` | List competency score details | Retrieve a paginated list of competency resources. |
| `GET` | `/api/v1/tenant/competency/values` | List competency values (structured) | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/values` | Create competency value | Create a new competency resource. |
| `GET` | `/api/v1/tenant/competency/values/{id}` | Get competency value by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/values/{id}` | Update competency value | Update an existing values record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/competency/values/{id}` | Delete competency value | Delete a values record by its unique ID. This action may be reversible depending on system configuration. |

### Tenant: Training & Development Management
**Description:** End-to-end training and development management including course catalog, session scheduling, participant registration, attendance tracking, materials, evaluations, and certificate issuance
**Endpoints:** 35 | **Paths:** 15
**Methods:** DELETE=7 GET=13 POST=7 PUT=8

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/trainings/categories` | Create training category | Create a new training category (e.g. Technical, Soft Skill, Leadership, Compliance). Categories are used to group related training courses. |
| `GET` | `/api/v1/tenant/trainings/categories` | List training categories | Retrieve a paginated list of training categories, ordered by code. Categories group training courses by subject area (Technical, Soft Skill, Leader... |
| `GET` | `/api/v1/tenant/trainings/categories/{id}` | Get training category by ID | Retrieve detailed information about a specific training category including its code, name, description, and active status. |
| `PUT` | `/api/v1/tenant/trainings/categories/{id}` | Update training category | Update an existing training category. Fields that are not provided will remain unchanged. |
| `DELETE` | `/api/v1/tenant/trainings/categories/{id}` | Delete training category | Soft-delete a training category. The category will be marked as deleted but retained in the database for historical purposes. |
| `POST` | `/api/v1/tenant/trainings/certificates` | Issue training certificate | Issue a certificate to a training participant. Requires a unique certificate number and issued date. Optionally set an expiry date. |
| `GET` | `/api/v1/tenant/trainings/certificates` | List training certificates | Retrieve a paginated list of issued certificates. Filter by participant_id to view all certificates for a specific participant. |
| `GET` | `/api/v1/tenant/trainings/certificates/{id}` | Get training certificate by ID | Retrieve a specific training certificate by its unique ID, including certificate number and validity period. |
| `PUT` | `/api/v1/tenant/trainings/certificates/{id}` | Update training certificate | Update a training certificate's number, issued date, or expiry date. |
| `DELETE` | `/api/v1/tenant/trainings/certificates/{id}` | Delete training certificate | Delete a training certificate record. |
| `POST` | `/api/v1/tenant/trainings/courses` | Create training course | Create a new training course under a specific category. Each course has a unique code and can include duration, cost, and minimum score requirements. |
| `GET` | `/api/v1/tenant/trainings/courses` | List training courses | Retrieve a paginated list of training courses. Optionally filter by category_id to view courses within a specific category. |
| `GET` | `/api/v1/tenant/trainings/courses/{id}` | Get training course by ID | Retrieve detailed information about a specific training course including category, duration, cost, and certification settings. |
| `PUT` | `/api/v1/tenant/trainings/courses/{id}` | Update training course | Update an existing training course. Only provided fields will be updated. |
| `DELETE` | `/api/v1/tenant/trainings/courses/{id}` | Delete training course | Soft-delete a training course. Associated sessions and materials are not deleted. |
| `POST` | `/api/v1/tenant/trainings/evaluations` | Create training evaluation | Submit a training evaluation/feedback for a session. Rating must be between 1 and 5, with optional textual feedback. |
| `GET` | `/api/v1/tenant/trainings/evaluations` | List training evaluations | Retrieve a paginated list of training evaluations. Filter by session_id for session feedback or employee_id for personal evaluation history. |
| `GET` | `/api/v1/tenant/trainings/evaluations/{id}` | Get training evaluation by ID | Retrieve a specific training evaluation including rating and feedback details. |
| `PUT` | `/api/v1/tenant/trainings/evaluations/{id}` | Update training evaluation | Update a training evaluation's rating or feedback text. |
| `DELETE` | `/api/v1/tenant/trainings/evaluations/{id}` | Delete training evaluation | Remove a training evaluation from the system. |
| `POST` | `/api/v1/tenant/trainings/materials` | Create training material | Add a new material (file, document, or resource) to a training session. Supports file URLs and type classification. |
| `GET` | `/api/v1/tenant/trainings/materials` | List training materials by session | List all materials attached to a training session. Requires session_id as a query parameter. Results are ordered by sort_order. |
| `PUT` | `/api/v1/tenant/trainings/materials/{id}` | Update training material | Update a training material's title, file URL, file type, or sort order. |
| `DELETE` | `/api/v1/tenant/trainings/materials/{id}` | Delete training material | Remove a training material from the session. |
| `POST` | `/api/v1/tenant/trainings/participants` | Register training participant | Register an employee as a participant in a training session. Validates that the session is not cancelled and has available quota. |
| `GET` | `/api/v1/tenant/trainings/participants` | List training participants | Retrieve a paginated list of training participants. Filter by session_id to view all participants in a session, or by employee_id to view an employ... |
| `GET` | `/api/v1/tenant/trainings/participants/{id}` | Get training participant by ID | Retrieve participant details including attendance status, score, and completion date. |
| `PUT` | `/api/v1/tenant/trainings/participants/{id}` | Update training participant | Update a participant's attendance status (PRESENT, ABSENT, EXCUSED) or score. Setting a score automatically marks the participant as completed with... |
| `DELETE` | `/api/v1/tenant/trainings/participants/{id}` | Delete training participant | Remove a participant from a training session. |
| `POST` | `/api/v1/tenant/trainings/sessions` | Create training session | Schedule a new training session/class for a course. Defines the trainer, date range, location, and maximum participant quota. |
| `GET` | `/api/v1/tenant/trainings/sessions` | List training sessions | Retrieve a paginated list of training sessions. Supports filtering by course_id and status (SCHEDULED, IN_PROGRESS, COMPLETED, CANCELLED). |
| `GET` | `/api/v1/tenant/trainings/sessions/{id}` | Get training session by ID | Retrieve detailed information about a specific training session including course, trainer, schedule, quota, and current status. |
| `PUT` | `/api/v1/tenant/trainings/sessions/{id}` | Update training session | Update an existing training session's schedule, trainer, location, or quota. |
| `DELETE` | `/api/v1/tenant/trainings/sessions/{id}` | Delete training session | Soft-delete a training session. Participants are not automatically removed. |
| `PUT` | `/api/v1/tenant/trainings/sessions/{id}/status` | Update training session status | Transition a training session through its lifecycle: SCHEDULED -> IN_PROGRESS -> COMPLETED or CANCELLED. Status changes affect participant registra... |

### Tenant: Recruitment & Onboarding (ATS)
**Description:** Recruitment & Onboarding (ATS) â€” job requisitions, candidate management, applications, interviews, and employee onboarding workflows
**Endpoints:** 33 | **Paths:** 16
**Methods:** DELETE=7 GET=12 POST=7 PUT=7

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/recruitment/applications` | Create job application | Submit a candidate's application to a job requisition. Candidate and requisition must exist |
| `GET` | `/api/v1/tenant/recruitment/applications` | List job applications | Retrieve paginated list of applications, optionally filtered by requisition, candidate, or status |
| `GET` | `/api/v1/tenant/recruitment/applications/{id}` | Get job application by ID | Retrieve application details including current status and notes |
| `DELETE` | `/api/v1/tenant/recruitment/applications/{id}` | Delete job application | Permanently delete an application record |
| `PUT` | `/api/v1/tenant/recruitment/applications/{id}/status` | Update application status | Update application status throughout the recruitment pipeline. Automatically updates requisition slots_filled when ACCEPTED |
| `POST` | `/api/v1/tenant/recruitment/candidates` | Create candidate | Register a new candidate. Email must be unique across the system |
| `GET` | `/api/v1/tenant/recruitment/candidates` | List candidates | Retrieve paginated list of candidates with optional search by name or email |
| `GET` | `/api/v1/tenant/recruitment/candidates/{id}` | Get candidate by ID | Retrieve detailed candidate information including contact details and resume links |
| `PUT` | `/api/v1/tenant/recruitment/candidates/{id}` | Update candidate | Update candidate profile fields. Only provided fields will be updated |
| `DELETE` | `/api/v1/tenant/recruitment/candidates/{id}` | Delete candidate | Permanently delete a candidate record |
| `POST` | `/api/v1/tenant/recruitment/employee-onboardings` | Create employee onboarding | Start onboarding for an accepted candidate. Automatically creates task items from active templates |
| `GET` | `/api/v1/tenant/recruitment/employee-onboardings` | List employee onboardings | Retrieve paginated list of employee onboardings, optionally filtered by status |
| `GET` | `/api/v1/tenant/recruitment/employee-onboardings/{id}` | Get employee onboarding by ID | Retrieve onboarding details including start date, buddy, and current status |
| `PUT` | `/api/v1/tenant/recruitment/employee-onboardings/{id}` | Update employee onboarding | Update onboarding details. Setting status to COMPLETED automatically records completion timestamp |
| `DELETE` | `/api/v1/tenant/recruitment/employee-onboardings/{id}` | Delete employee onboarding | Permanently delete an employee onboarding record and its task items |
| `GET` | `/api/v1/tenant/recruitment/employee-onboardings/{id}/task-items` | List onboarding task items | Retrieve all task items for a specific employee onboarding, ordered by due date |
| `POST` | `/api/v1/tenant/recruitment/interviews` | Create interview | Schedule a new interview for a job application with interviewer, stage, and time slot |
| `GET` | `/api/v1/tenant/recruitment/interviews` | List interviews | Retrieve paginated list of interviews, optionally filtered by application or interviewer |
| `GET` | `/api/v1/tenant/recruitment/interviews/{id}` | Get interview by ID | Retrieve interview details including score, feedback, and status |
| `PUT` | `/api/v1/tenant/recruitment/interviews/{id}` | Update interview | Update interview schedule, score, feedback, or status. Setting status to COMPLETED automatically records completion timestamp |
| `DELETE` | `/api/v1/tenant/recruitment/interviews/{id}` | Delete interview | Permanently delete an interview record |
| `POST` | `/api/v1/tenant/recruitment/onboarding-task-items` | Create onboarding task item | Add a custom task item to an employee onboarding. Can optionally link to a template |
| `PUT` | `/api/v1/tenant/recruitment/onboarding-task-items/{id}` | Update onboarding task item | Update task item details. Setting is_completed to true automatically records completion timestamp |
| `DELETE` | `/api/v1/tenant/recruitment/onboarding-task-items/{id}` | Delete onboarding task item | Permanently delete a task item |
| `POST` | `/api/v1/tenant/recruitment/onboarding-task-templates` | Create onboarding task template | Create a reusable task template for employee onboarding (e.g., IT Setup, Contract Signing) |
| `GET` | `/api/v1/tenant/recruitment/onboarding-task-templates` | List onboarding task templates | Retrieve paginated list of task templates, optionally filtered by category |
| `PUT` | `/api/v1/tenant/recruitment/onboarding-task-templates/{id}` | Update onboarding task template | Update a task template properties |
| `DELETE` | `/api/v1/tenant/recruitment/onboarding-task-templates/{id}` | Delete onboarding task template | Permanently delete a task template |
| `POST` | `/api/v1/tenant/recruitment/requisitions` | Create job requisition | Create a new job requisition with position details, salary range, and number of slots available |
| `GET` | `/api/v1/tenant/recruitment/requisitions` | List job requisitions | Retrieve paginated list of job requisitions, optionally filtered by organization and status |
| `GET` | `/api/v1/tenant/recruitment/requisitions/{id}` | Get job requisition by ID | Retrieve detailed job requisition information by UUID |
| `PUT` | `/api/v1/tenant/recruitment/requisitions/{id}` | Update job requisition | Update job requisition fields. Only provided fields will be updated |
| `DELETE` | `/api/v1/tenant/recruitment/requisitions/{id}` | Delete job requisition | Permanently delete a job requisition |

### Tenant: Time & Attendance
**Description:** Time and attendance management including company settings, shifts, employee shift assignments, geofence locations, check-in/check-out events, daily work sessions, overtime requests, and exempt positions
**Endpoints:** 30 | **Paths:** 15
**Methods:** DELETE=4 GET=15 POST=6 PUT=5

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/attendance/employee-shifts` | List employee shift assignments | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/employee-shifts` | Assign a shift to an employee | Create a new employee shifts record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/attendance/employee-shifts/{id}` | Get employee shift assignment by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/employee-shifts/{id}` | Update employee shift assignment | Update an attendance record. |
| `DELETE` | `/api/v1/tenant/attendance/employee-shifts/{id}` | Delete employee shift assignment | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/events` | List attendance events (check-in/out) | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/events` | Create an attendance event (check-in/out) | Create a new events record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/attendance/events/{id}` | Get event by ID | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/exempt-positions` | List exempt positions (positions not requiring attendance) | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/exempt-positions` | Create an exempt position | Create a new exempt positions record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/attendance/exempt-positions/{id}` | Get exempt position by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/exempt-positions/{id}` | Update an exempt position | Update an attendance record. |
| `DELETE` | `/api/v1/tenant/attendance/exempt-positions/{id}` | Delete an exempt position | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/locations` | List attendance locations (geofence) | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/locations` | Create an attendance location (geofence) | Create a new locations record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/attendance/locations/{id}` | Get location by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/locations/{id}` | Update a location | Update an attendance record. |
| `DELETE` | `/api/v1/tenant/attendance/locations/{id}` | Delete a location | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/overtime-requests` | List overtime requests | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/overtime-requests` | Create an overtime request | Create a new overtime requests record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/attendance/overtime-requests/{id}` | Get overtime request by ID | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/sessions` | List daily work sessions | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/sessions/detail` | Get session detail for an employee on a specific date | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/settings` | Get company attendance settings | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/settings` | Upsert company attendance settings | Update an attendance record. |
| `GET` | `/api/v1/tenant/attendance/shifts` | List company shifts | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/shifts` | Create a company shift | Create a new shifts record. Validates required fields and returns the created resource with its assigned ID. |
| `GET` | `/api/v1/tenant/attendance/shifts/{id}` | Get shift by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/shifts/{id}` | Update a shift | Update an attendance record. |
| `DELETE` | `/api/v1/tenant/attendance/shifts/{id}` | Delete a shift | Remove an attendance record. |

### Tenant: Leave & Time Off
**Description:** Leave and time off management including leave types, accrual policies, leave reasons, leave requests, leave request details, and employee leave balances
**Endpoints:** 23 | **Paths:** 12
**Methods:** DELETE=4 GET=11 POST=4 PUT=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/leave/accrual-policies` | List accrual policies with pagination | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/accrual-policies` | Create an accrual policy | Create a new leave resource. |
| `GET` | `/api/v1/tenant/leave/accrual-policies/{id}` | Get accrual policy by ID | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/accrual-policies/{id}` | Update accrual policy | Update an existing accrual policies record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/leave/accrual-policies/{id}` | Delete accrual policy | Delete a accrual policies record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/leave/balances` | List leave balances with pagination | Retrieve a paginated list of leave resources. |
| `GET` | `/api/v1/tenant/leave/balances/employees/{employeeId}/types/{leaveTypeId}` | Get leave balance for specific employee and leave type | Retrieve a paginated list of leave resources. |
| `GET` | `/api/v1/tenant/leave/reasons` | List leave reasons | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/reasons` | Create a leave reason | Create a new leave resource. |
| `GET` | `/api/v1/tenant/leave/reasons/{id}` | Get leave reason by ID | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/reasons/{id}` | Update leave reason | Update an existing reasons record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/leave/reasons/{id}` | Delete leave reason | Delete a reasons record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/leave/requests` | List leave requests with pagination | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/requests` | Create a leave request | Create a new leave resource. |
| `GET` | `/api/v1/tenant/leave/requests/{id}` | Get leave request by ID | Retrieve a paginated list of leave resources. |
| `DELETE` | `/api/v1/tenant/leave/requests/{id}` | Delete leave request | Delete a requests record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/leave/requests/{id}/details` | List leave request details (daily breakdown) | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/requests/{id}/status` | Update leave request status (approve/reject/cancel) | Update an existing status record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `GET` | `/api/v1/tenant/leave/types` | List leave types with pagination | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/types` | Create a new leave type | Create a new leave resource. |
| `GET` | `/api/v1/tenant/leave/types/{id}` | Get leave type by ID | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/types/{id}` | Update leave type | Update an existing types record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/leave/types/{id}` | Delete leave type | Delete a types record by its unique ID. This action may be reversible depending on system configuration. |

### Tenant: Career Intelligence
**Description:** Career Intelligence & Talent Management â€” strategic talent analytics for 9-box talent mapping, career interests tracking, career path gap analysis, and succession planning. Provides talent review data to identify high-potential employees, plan career development, and ensure leadership pipeline readiness.
**Endpoints:** 19 | **Paths:** 11
**Methods:** DELETE=3 GET=10 POST=4 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/career-intelligence/interests` | List career interests | Retrieve a paginated list of career interests. Filter by employee to see an individual's career aspirations. |
| `POST` | `/api/v1/tenant/career-intelligence/interests` | Record career interest | Record a career interest for an employee. Supports interest types: LEADERSHIP, SPECIALIST, INTERNATIONAL, and ENTREPRENEUR with optional target pos... |
| `GET` | `/api/v1/tenant/career-intelligence/interests/employee/{employeeId}` | Get employee career interests | Get all active career interests for a specific employee. Used in talent review and career development planning. |
| `GET` | `/api/v1/tenant/career-intelligence/paths` | List career paths | Retrieve a paginated list of defined career paths between position titles. Includes PROMOTION, LATERAL, DEMOTION, and CROSSFUNCTIONAL path types. |
| `POST` | `/api/v1/tenant/career-intelligence/paths` | Create career path | Define a career path between two position titles. Specifies the path type (PROMOTION/LATERAL/DEMOTION/CROSSFUNCTIONAL), typical tenure, requirement... |
| `GET` | `/api/v1/tenant/career-intelligence/paths/gap-analysis` | Career gap analysis | Analyze the gap between an employee's current qualifications and the requirements of a target position title. Returns matched skills, total require... |
| `DELETE` | `/api/v1/tenant/career-intelligence/paths/{id}` | Delete career path | Hapus satu jalur karier (career path). |
| `GET` | `/api/v1/tenant/career-intelligence/successions` | List succession plans | Retrieve a paginated list of succession plans. Shows which positions have identified successors with readiness levels (READY_NOW, READY_1YR, READY_... |
| `POST` | `/api/v1/tenant/career-intelligence/successions` | Create succession plan | Create a succession plan for a key position. Identifies a potential successor with readiness level, priority order, target date, and development plan. |
| `GET` | `/api/v1/tenant/career-intelligence/successions/{id}` | Get succession plan by ID | Get detailed information about a specific succession plan including successor details, readiness level, and development plan. |
| `PUT` | `/api/v1/tenant/career-intelligence/successions/{id}` | Update succession plan | Update an existing succession plan's readiness level, priority order, target date, development plan, and/or notes. |
| `DELETE` | `/api/v1/tenant/career-intelligence/successions/{id}` | Delete succession plan | Soft-delete a succession plan. |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps` | List talent maps | Retrieve a paginated list of talent mapping entries for 9-box assessment. Filter by assessment period (e.g., 2026-Q1) and/or employee. |
| `POST` | `/api/v1/tenant/career-intelligence/talent-maps` | Create talent map entry | Create a new 9-box talent mapping entry for an employee. Combines performance rating (LOW/MEDIUM/HIGH) and potential rating (LOW/MEDIUM/HIGH) to de... |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps/employee/{employeeId}` | Get employee talent profile | Get a comprehensive talent profile for an employee: current 9-box assessment, historical talent map entries, career interests, and recommended next... |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps/grid` | Get talent map grid (9-box view) | Get the 9-box talent grid overview for a given period. Returns employee counts per quadrant (9-BOX-1 through 9-BOX-9) with labels and descriptions ... |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps/{id}` | Get talent map by ID | Get detailed information about a specific talent mapping entry. |
| `PUT` | `/api/v1/tenant/career-intelligence/talent-maps/{id}` | Update talent map entry | Update an existing talent map entry's performance rating, potential rating, and/or notes. Changes to performance/potential automatically recalculat... |
| `DELETE` | `/api/v1/tenant/career-intelligence/talent-maps/{id}` | Delete talent map entry | Soft-delete a talent map entry. |

### Tenant: Organizations
**Description:** Organization structure (tree hierarchy)
**Endpoints:** 18 | **Paths:** 11
**Methods:** DELETE=2 GET=9 POST=5 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/organization-summaries` | Create organization summary | Create a new organization summary record. Requires code, decree_no, and decree_date. Optional status (default: inactive). Only ONE summary can have... |
| `GET` | `/api/v1/tenant/organization-summaries` | List organization summaries (paginated) | Retrieve a paginated list of organization summaries sorted by creation date (descending). |
| `GET` | `/api/v1/tenant/organization-summaries/stats` | Get organization summary statistics | Get aggregate statistics about organization summaries and organizations. |
| `GET` | `/api/v1/tenant/organization-summaries/{id}` | Get organization summary by ID | Get detailed information about a specific organization summary by its ID. |
| `PUT` | `/api/v1/tenant/organization-summaries/{id}` | Update organization summary | Update an existing organization summary's code, decree_no, decree_date, or status. Only ONE summary can have status=active at a time — setting stat... |
| `DELETE` | `/api/v1/tenant/organization-summaries/{id}` | Delete organization summary | Soft-delete an organization summary. Cannot delete if organizations are still attached to this summary. |
| `GET` | `/api/v1/tenant/organizations` | List organizations or get tree | Get detailed information about an organizational unit and its children. Mendukung ?search= (filter code/full_code/nomenclature), ?summary_id=, ?act... |
| `POST` | `/api/v1/tenant/organizations` | Create organization | Create a new organizations record. Validates required fields and returns the created resource with its assigned ID. |
| `POST` | `/api/v1/tenant/organizations/clone` | Clone current organization tree to a draft version | Create a DRAFT version snapshot of the current organization tree for restructuring simulation. The clone preserves the complete tree structure in a... |
| `GET` | `/api/v1/tenant/organizations/history` | List organization change history | Retrieve a paginated audit trail of all changes made to the organization structure. Supports filtering by specific organization ID. Each entry reco... |
| `POST` | `/api/v1/tenant/organizations/versions` | Create a new version snapshot of the organization tree | Capture a point-in-time snapshot of the entire organization tree structure. The snapshot includes both the hierarchical tree view and a flat list o... |
| `GET` | `/api/v1/tenant/organizations/versions` | List all organization versions | Retrieve a paginated list of all saved organization structure versions. Each version entry includes metadata (name, status, node count) but exclude... |
| `GET` | `/api/v1/tenant/organizations/versions/{id}` | Get organization version detail | Get detailed information about a specific organization version. By default, the snapshot payload is excluded for performance. Append ?snapshot=true... |
| `GET` | `/api/v1/tenant/organizations/versions/{id}/diff/{targetId}` | Compare two organization versions and show differences | Compare two organization versions and produce a detailed diff. The response categorizes changes into ADDED (nodes in target but not source), REMOVE... |
| `POST` | `/api/v1/tenant/organizations/versions/{id}/restore` | Restore organization tree from a version snapshot | Restore the entire organization structure from a version snapshot. This operation atomically: (1) hard-deletes all current organizations, (2) recre... |
| `GET` | `/api/v1/tenant/organizations/{id}` | Get organization by ID | Get detailed information about an organizational unit and its children. |
| `PUT` | `/api/v1/tenant/organizations/{id}` | Update organization | Update organizational unit code, name, or parent assignment. |
| `DELETE` | `/api/v1/tenant/organizations/{id}` | Delete organization | Remove an organizational unit from the hierarchy. |

### Tenant: Approval
**Description:** Approval engine for multi-level workflow
**Endpoints:** 17 | **Paths:** 11
**Methods:** DELETE=2 GET=8 POST=5 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/approval/active-flow` | Get active approval flow by module | Resolusi otomatis alur persetujuan aktif untuk sebuah module (?module=xxx) — dipakai consumer yang ingin auto-resolve flow tanpa memilih flow_id ma... |
| `GET` | `/api/v1/tenant/approval/available-modules` | List available approval modules | Ambil slug module yang aktif/disubscribe tenant — dipakai flow builder agar hanya menampilkan module yang tersedia. |
| `GET` | `/api/v1/tenant/approval/flows` | List approval flows | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/flows` | Create approval flow | Create a new approval resource. |
| `GET` | `/api/v1/tenant/approval/flows/{flowId}` | Get approval flow by ID | Retrieve a paginated list of approval resources. |
| `PUT` | `/api/v1/tenant/approval/flows/{flowId}` | Update approval flow | Update an existing flows record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/approval/flows/{flowId}` | Delete approval flow | Delete a flows record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/approval/flows/{flowId}/steps` | List approval flow steps | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/flows/{flowId}/steps` | Create approval flow step | Create a new approval resource. |
| `PUT` | `/api/v1/tenant/approval/flows/{flowId}/steps/{stepId}` | Update approval flow step | Update an existing steps record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/approval/flows/{flowId}/steps/{stepId}` | Delete approval flow step | Delete a steps record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/approval/instances` | List approval instances | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/instances` | Create approval instance | Create a new approval resource. |
| `GET` | `/api/v1/tenant/approval/instances/{id}` | Get approval instance by ID | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/instances/{id}/actions` | Submit approval action (approve/reject) | Create a new approval resource. |
| `POST` | `/api/v1/tenant/approval/instances/{id}/cancel` | Cancel approval instance | Cancel an active approval instance. This will void all pending tasks and mark the instance as CANCELLED. Only instances in PENDING status can be ca... |
| `GET` | `/api/v1/tenant/approval/tasks/pending` | List my pending approval tasks | Retrieve a paginated list of approval resources. |

### Tenant: Employee Movement & Career Management
**Description:** Employee career movements management including promotions, demotions, mutations, contract extensions (PKWT), retirements, offboarding, and employment contract management
**Endpoints:** 16 | **Paths:** 10
**Methods:** DELETE=2 GET=6 POST=6 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/employee-movements/contracts` | List employee contracts | Retrieve a paginated list of employee movement resources. |
| `POST` | `/api/v1/tenant/employee-movements/contracts` | Create employee contract | Create a new employee movement resource. |
| `GET` | `/api/v1/tenant/employee-movements/contracts/{id}` | Get contract by ID | Retrieve a paginated list of employee movement resources. |
| `PUT` | `/api/v1/tenant/employee-movements/contracts/{id}` | Update contract | Update an existing contracts record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/employee-movements/contracts/{id}` | Delete contract | Delete a contracts record by its unique ID. This action may be reversible depending on system configuration. |
| `GET` | `/api/v1/tenant/employee-movements/employees/{employeeId}/contracts` | List contracts by employee | Retrieve a paginated list of employee movement resources. |
| `GET` | `/api/v1/tenant/employee-movements/employees/{employeeId}/movements` | List movements by employee | Retrieve a paginated list of employee movement resources. |
| `GET` | `/api/v1/tenant/employee-movements/movements` | List employee movements | Retrieve a paginated list of employee movement resources. |
| `POST` | `/api/v1/tenant/employee-movements/movements` | Create employee movement | Create a new employee movement resource. |
| `GET` | `/api/v1/tenant/employee-movements/movements/{id}` | Get movement by ID | Retrieve a paginated list of employee movement resources. |
| `PUT` | `/api/v1/tenant/employee-movements/movements/{id}` | Update movement | Update an existing movements record by its unique ID. Accepts partial updates; only provided fields will be modified. |
| `DELETE` | `/api/v1/tenant/employee-movements/movements/{id}` | Delete movement | Delete a movements record by its unique ID. This action may be reversible depending on system configuration. |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/approve` | Approve movement | Create a new employee movement resource. |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/cancel` | Cancel movement | Create a new employee movement resource. |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/execute` | Execute movement | Create a new employee movement resource. |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/submit` | Submit employee movement for approval | Kirim movement berstatus draft ke alur persetujuan (approval flow) terpusat. Setelah disetujui, approval engine akan mengeksekusi perpindahan. |

### Tenant: Reimbursement & Claim
**Description:** Reimbursement & claim management including reimbursement types, requests, items, and payment processing
**Endpoints:** 15 | **Paths:** 7
**Methods:** DELETE=3 GET=5 POST=3 PUT=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/reimbursements/requests` | Create reimbursement request | Submit a new reimbursement request for approval. Includes the reimbursement type, title, description, and optional notes. The request starts in DRA... |
| `GET` | `/api/v1/tenant/reimbursements/requests` | List reimbursement requests | Retrieve a paginated list of reimbursement requests. Supports filtering by employee, status (DRAFT/SUBMITTED/APPROVED/REJECTED/PAID/CANCELLED), and... |
| `GET` | `/api/v1/tenant/reimbursements/requests/{id}` | Get reimbursement request | Get detailed information about a specific reimbursement request, including its items, status history, and approval timeline. |
| `PUT` | `/api/v1/tenant/reimbursements/requests/{id}` | Update reimbursement request | Update a reimbursement request. Only requests in DRAFT status can be modified. Changes to title, description, or items are allowed before submission. |
| `DELETE` | `/api/v1/tenant/reimbursements/requests/{id}` | Delete reimbursement request | Delete a reimbursement request. Only requests in DRAFT or CANCELLED status can be deleted. Finalizes the record removal from the system. |
| `PUT` | `/api/v1/tenant/reimbursements/requests/{id}/status` | Update reimbursement request status | Transition a reimbursement request through its approval workflow. Valid transitions: DRAFT->SUBMITTED (submit), SUBMITTED->APPROVED (approve), SUBM... |
| `POST` | `/api/v1/tenant/reimbursements/requests/{requestId}/items` | Add reimbursement item | Add a new expense item to a reimbursement request. Each item represents a single expense with date, type (e.g., transportation, accommodation, meal... |
| `GET` | `/api/v1/tenant/reimbursements/requests/{requestId}/items` | List reimbursement items | Retrieve all expense items attached to a specific reimbursement request. Returns item details including expense date, type, amount, and receipt att... |
| `PUT` | `/api/v1/tenant/reimbursements/requests/{requestId}/items/{itemId}` | Update reimbursement item | Update an expense item's details including expense date, type, amount, description, or receipt URL. Only modifiable while the request is in DRAFT s... |
| `DELETE` | `/api/v1/tenant/reimbursements/requests/{requestId}/items/{itemId}` | Delete reimbursement item | Remove an expense item from a reimbursement request. The total request amount will be recalculated automatically. |
| `POST` | `/api/v1/tenant/reimbursements/types` | Create reimbursement type | Create a new reimbursement type for the company. Defines the category name, maximum claimable amount, and a description of eligible expenses. Used ... |
| `GET` | `/api/v1/tenant/reimbursements/types` | List reimbursement types | Retrieve a paginated list of reimbursement type configurations. Supports filtering by name or category. |
| `GET` | `/api/v1/tenant/reimbursements/types/{id}` | Get reimbursement type | Get detailed information about a specific reimbursement type, including its name, maximum amount, and description. |
| `PUT` | `/api/v1/tenant/reimbursements/types/{id}` | Update reimbursement type | Update a reimbursement type's name, maximum amount, or description. Changes apply to new requests only. |
| `DELETE` | `/api/v1/tenant/reimbursements/types/{id}` | Delete reimbursement type | Delete a reimbursement type from the system. Existing requests using this type will retain their type reference. |

### Platform: Companies
**Description:** Company/Tenant management
**Endpoints:** 11 | **Paths:** 8
**Methods:** DELETE=1 GET=2 POST=7 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/platform/companies` | Create a new company/tenant | Register a new company tenant. Also creates a company_admin user automatically. Request now includes admin_name, admin_email, admin_password fields... |
| `GET` | `/api/v1/platform/companies` | List all companies | Retrieve a paginated list of all registered companies (tenants) in the platform. Includes company status, contact information, subscription details... |
| `GET` | `/api/v1/platform/companies/{id}` | Get company by ID | Get detailed information about a specific company/tenant including its status, contact details, subscription plan, database connection health, and ... |
| `PUT` | `/api/v1/platform/companies/{id}` | Update company | Update a company's profile information including name, email, phone, address, and other contact details. |
| `DELETE` | `/api/v1/platform/companies/{id}` | Soft delete company (deactivate connection + deleted_at) | Soft-delete a company tenant. Deactivates the tenant database connection and sets the deleted_at timestamp. The company record is hidden from stand... |
| `POST` | `/api/v1/platform/companies/{id}/activate` | Activate a company/tenant (reactivate connection) | Reactivate a previously suspended company tenant. Re-establishes the database connection and sets the company status back to 'active'. All tenant A... |
| `POST` | `/api/v1/platform/companies/{id}/backup` | Trigger tenant backup (Phase 2) | Trigger an on-demand database backup for the specified company tenant. The backup is stored according to the platform's backup configuration. |
| `POST` | `/api/v1/platform/companies/{id}/restore` | Trigger tenant restore (Phase 2) | Restore a company tenant's database from a previously created backup. Requires a valid backup reference. |
| `POST` | `/api/v1/platform/companies/{id}/rotate-credentials` | Rotate tenant DB credentials (ALTER USER + update encrypted connection) | Rotate the tenant database credentials for a company. Runs ALTER USER on the DB server (dialect-aware), updates tenant_connections with the new pas... |
| `POST` | `/api/v1/platform/companies/{id}/suspend` | Suspend a company/tenant (deactivate connection) | Suspend a company tenant â€” deactivates the database connection, clears cached connections, and sets the company status to 'suspended'. All tenant... |
| `POST` | `/api/v1/platform/companies/{id}/terminate` | Terminate a company/tenant (drop database + remove connection) | Permanently terminate a company tenant. This drops the tenant database entirely, removes the connection record, and sets the company status to 'ter... |

### Platform: RBAC Management
**Description:** Role-based access control management for roles, permissions, and role-permission assignments
**Endpoints:** 10 | **Paths:** 6
**Methods:** DELETE=3 GET=3 POST=3 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/rbac/permissions` | List all permissions (resource + action) | Retrieve all available permissions in the system. Permissions are defined as resource.action pairs (e.g., company.create, user.view). |
| `POST` | `/api/v1/platform/rbac/permissions` | Create a new permission (resource.action) | Create a new permission in the format resource.action (e.g., report.view). Permissions can be assigned to roles via the role-permission assignment ... |
| `DELETE` | `/api/v1/platform/rbac/permissions/{id}` | Delete a permission (non-system only) | Delete a non-system permission. The permission will be removed from all role assignments. System permissions cannot be deleted. |
| `GET` | `/api/v1/platform/rbac/roles` | List all roles with their permissions | Retrieve all RBAC roles with their associated permissions. Roles are organized in a hierarchy with permission inheritance from parent roles. |
| `POST` | `/api/v1/platform/rbac/roles` | Create a new role | Create a new RBAC role with a name, description, and optional parent role. New roles inherit permissions from their parent role. |
| `GET` | `/api/v1/platform/rbac/roles/{id}` | Get role by ID with permissions | Get detailed information about a specific role including its name, description, parent role (if any), and all assigned permissions. |
| `PUT` | `/api/v1/platform/rbac/roles/{id}` | Update role (name, description, parent) | Update a role's name, description, or parent role assignment. Changes to system roles may be restricted. |
| `DELETE` | `/api/v1/platform/rbac/roles/{id}` | Delete a role (non-system roles only) | Delete a non-system role. Users assigned to this role will lose their associated permissions until reassigned. System roles (super_admin) cannot be... |
| `POST` | `/api/v1/platform/rbac/roles/{id}/permissions` | Assign a permission to a role | Assign a permission to a role. The permission becomes available to all users with that role (and child roles through inheritance). The enforcer aut... |
| `DELETE` | `/api/v1/platform/rbac/roles/{id}/permissions/{permissionId}` | Revoke a permission from a role | Revoke a permission from a role. The permission will no longer be available to users with that role. The enforcer auto-reloads to apply changes imm... |

### Platform: Packages
**Description:** Package management â€” bundle tenant modules with pricing, dependency validation, and publishing
**Endpoints:** 9 | **Paths:** 6
**Methods:** DELETE=1 GET=4 POST=3 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/packages` | List all packages (admin) | Retrieve a paginated list of all packages. Includes draft, published, and archived packages with their module associations. |
| `POST` | `/api/v1/platform/packages` | Create a new package | Create a new package that bundles tenant modules with pricing. Validates module dependencies before creation â€” all required dependencies (depends... |
| `GET` | `/api/v1/platform/packages/{id}` | Get package by ID | Get detailed information about a specific package including its modules, pricing, status, and sort order. |
| `PUT` | `/api/v1/platform/packages/{id}` | Update package | Update a package's metadata, pricing, status, or module associations. When updating modules, dependency validation is re-run. |
| `DELETE` | `/api/v1/platform/packages/{id}` | Delete package (soft-delete) | Soft-delete a package. Sets the deleted_at timestamp and removes module associations. The package is hidden from standard queries. |
| `POST` | `/api/v1/platform/packages/{id}/publish` | Publish package | Publish a package to make it visible on public endpoints. Validates that the package has at least one module and all module dependencies are fulfil... |
| `POST` | `/api/v1/platform/packages/{id}/unpublish` | Unpublish package | Unpublish a package. Sets status back to 'draft' and removes it from public endpoints. |
| `GET` | `/api/v1/platform/packages/{id}/validate` | Validate package module dependencies | Validate that all module dependencies within a package are fulfilled. Returns a detailed report showing each module's dependency status (resolved/u... |
| `GET` | `/api/v1/public/packages` | List published packages (public) | Retrieve a list of published packages for public display. No authentication required. Returns package name, description, price, and included module... |

### Tenant: RBAC Management
**Description:** Tenant RBAC role & permission management
**Endpoints:** 8 | **Paths:** 6
**Methods:** DELETE=1 GET=3 POST=1 PUT=3

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/rbac/permissions` | List tenant RBAC permissions | Ambil daftar permission RBAC tenant. |
| `GET` | `/api/v1/tenant/rbac/roles` | List tenant RBAC roles | Ambil daftar role RBAC tenant beserta permission-nya. |
| `POST` | `/api/v1/tenant/rbac/roles` | Create tenant RBAC role | Buat role RBAC tenant baru. |
| `PUT` | `/api/v1/tenant/rbac/roles/{id}` | Update tenant RBAC role | Perbarui nama/deskripsi role RBAC tenant. |
| `DELETE` | `/api/v1/tenant/rbac/roles/{id}` | Delete tenant RBAC role | Hapus role RBAC tenant (role system tidak dapat dihapus). |
| `PUT` | `/api/v1/tenant/rbac/roles/{id}/permissions` | Assign permissions to tenant role | Ganti (replace) daftar permission milik role RBAC tenant. |
| `GET` | `/api/v1/tenant/rbac/users` | List tenant RBAC users | Ambil daftar user tenant beserta role-nya. |
| `PUT` | `/api/v1/tenant/rbac/users/{id}/roles` | Assign roles to tenant user | Ganti (replace) daftar role milik user tenant. |

### Platform: Modules
**Description:** Module registration and activation management
**Endpoints:** 7 | **Paths:** 5
**Methods:** GET=3 POST=3 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/modules` | List all registered modules | Retrieve a paginated list of all registered system modules. Each module represents a functional area of the HRIS platform (e.g., Payroll, Attendanc... |
| `POST` | `/api/v1/platform/modules` | Register a new module | Register a new system module with its name, slug, version, and optional dependencies. Modules can be activated or deactivated per company tenant. |
| `GET` | `/api/v1/platform/modules/{id}` | Get module by ID | Get detailed information about a specific module including its version, dependencies, and activation status across companies. |
| `PUT` | `/api/v1/platform/modules/{id}` | Update module | Update a module's configuration, metadata, version, or feature flags. Changes apply globally across all tenants. |
| `POST` | `/api/v1/platform/modules/{id}/activate` | Activate module for a company | Activate a module for a specific company tenant. The module's features become available in that tenant's API and UI. |
| `GET` | `/api/v1/platform/modules/{id}/companies` | List companies using this module | Retrieve a list of companies that have this module activated. Shows activation date and module-specific settings per company. |
| `POST` | `/api/v1/platform/modules/{id}/deactivate` | Deactivate module for a company | Deactivate a module for a specific company tenant. The module's features are hidden from that tenant's API and UI. |

### Platform: Users
**Description:** Platform user management
**Endpoints:** 6 | **Paths:** 3
**Methods:** DELETE=1 GET=2 POST=1 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/users` | List all platform users | Retrieve a paginated list of platform user accounts. Now includes company_name field. Supports filtering by company, role, and search term. |
| `POST` | `/api/v1/platform/users` | Create a new platform user | Register a new platform user account. Assigns the user to a company with a specific role (super_admin or company_admin). The user will receive acce... |
| `DELETE` | `/api/v1/platform/users/{id}` | Delete platform user | Soft-delete a platform user account. Super admin users cannot be deleted. |
| `GET` | `/api/v1/platform/users/{id}` | Get platform user by ID | Ambil detail platform user berdasarkan ID. |
| `PUT` | `/api/v1/platform/users/{id}` | Update platform user | Perbarui nama/email/role/status platform user. |
| `PUT` | `/api/v1/platform/users/{id}/password` | Change user password | Change a user's password. Requires current password verification, new password (min 6 characters), and confirmation. |

### Platform: Licenses
**Description:** License management for companies
**Endpoints:** 5 | **Paths:** 2
**Methods:** DELETE=1 GET=2 POST=1 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/licenses` | List all licenses | Retrieve a paginated list of all license records. Licenses define the plan type, feature entitlements, and validity period for each company tenant. |
| `POST` | `/api/v1/platform/licenses` | Create a new license for company | Issue a new software license to a company tenant. Specifies the plan type (trial, subscription), start date, end date, and seat count. If a package... |
| `GET` | `/api/v1/platform/licenses/{id}` | Get license by ID | Get detailed information about a specific license including plan type, validity period, seat usage, and feature entitlements. |
| `PUT` | `/api/v1/platform/licenses/{id}` | Update license | Update license terms including plan upgrade/downgrade, extension of validity period, seat count adjustments, or license status changes. |
| `DELETE` | `/api/v1/platform/licenses/{id}` | Delete license | Hapus lisensi perusahaan. |

### Platform: Monitoring
**Description:** Platform and tenant health monitoring
**Endpoints:** 5 | **Paths:** 5
**Methods:** GET=5

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/monitoring/health` | Platform health check with database status | Platform health check endpoint providing detailed status of database connectivity, Redis cache health, and overall service uptime metrics. |
| `GET` | `/api/v1/platform/monitoring/pool` | Get connection pool statistics | Get real-time database connection pool statistics including open, idle, and wait counts for both platform and tenant pools. |
| `GET` | `/api/v1/platform/monitoring/seed-status` | Get seed data status | Get seed data status for all tenant master tables. Shows record counts and seeding status per table. Optionally filter by company_id. |
| `GET` | `/api/v1/platform/monitoring/tenants` | List all active tenant connections health | List all active tenant database connections with their health status, pool statistics (open/idle connections), and last activity timestamps. |
| `GET` | `/api/v1/platform/monitoring/tenants/{id}` | Get tenant connection health detail | Get detailed health information for a specific tenant database connection, including connection pool stats, query latency, and error counts. |

### Health
**Description:** Health check endpoints
**Endpoints:** 4 | **Paths:** 4
**Methods:** GET=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/docs` | Scalar API Documentation UI | Interactive API documentation powered by Scalar UI. Browse all available endpoints, test requests, and view response schemas directly from the brow... |
| `GET` | `/healthz` | Server health check | Simple health check endpoint for load balancer probes. Returns HTTP 200 with service status when the server is running. |
| `GET` | `/openapi.json` | OpenAPI 3.0 Specification | Download the complete OpenAPI 3.0 specification as JSON. Compatible with tools like Postman, Insomnia, Swagger Editor, and client code generators. |
| `GET` | `/readyz` | Readiness check | Readiness check endpoint for Kubernetes or container orchestration probes. Returns HTTP 200 when the server is ready to accept traffic. |

### Tenant: Packages
**Description:** Published package browsing for authenticated tenant users
**Endpoints:** 4 | **Paths:** 4
**Methods:** GET=2 POST=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/company-modules` | List active modules for current company | Retrieve the list of active modules for the authenticated company. Returns module slugs and names that are activated for this tenant. Used for UI m... |
| `GET` | `/api/v1/tenant/packages` | List published packages (tenant) | Retrieve a list of published packages for authenticated tenant users. Requires JWT Bearer Token. Returns the same data as the public endpoint but w... |
| `POST` | `/api/v1/tenant/packages/{id}/subscribe` | Subscribe to a package (create/renew license) | Subscribe the authenticated company to a published package. Creates a new license for the company associated with the specified package and auto-ac... |
| `POST` | `/api/v1/tenant/packages/{id}/unsubscribe` | Unsubscribe from a package (deactivate modules + suspend license) | Unsubscribe the authenticated company from a package. Deactivates all modules included in the package and suspends the active license associated wi... |

### Tenant: User Accounts
**Description:** Employee login account management
**Endpoints:** 3 | **Paths:** 2
**Methods:** GET=1 POST=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/user-accounts/employees/{employeeId}` | Get employee account status | Ambil status akun login employee (email, password_set, setup token). |
| `POST` | `/api/v1/tenant/user-accounts/employees/{employeeId}` | Create employee account | Buat akun login employee (kirim email setup password). |
| `POST` | `/api/v1/tenant/user-accounts/employees/{employeeId}/resend` | Resend account setup email | Kirim ulang email link setup password akun employee. |

### Platform: Auth
**Description:** Platform authentication (login, refresh token)
**Endpoints:** 2 | **Paths:** 2
**Methods:** POST=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/platform/login` | Platform admin login | Authenticate a platform admin user with email and password credentials. Returns a JWT access token (short-lived) and a refresh token (long-lived) f... |
| `POST` | `/api/v1/platform/refresh` | Refresh access token | Exchange a valid refresh token for a new access token. Use this endpoint to maintain session continuity without requiring the user to re-login. |

### Public
**Description:** Public endpoints — no authentication required
**Endpoints:** 2 | **Paths:** 2
**Methods:** GET=1 POST=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/public/account/setup-password` | Set account password via email link | Atur password akun login employee melalui link email (tanpa autentikasi). |
| `GET` | `/api/v1/public/companies/resolve` | Resolve company by hostname/subdomain | Public endpoint (no auth) to determine the company from the app URL hostname/subdomain (SaaS mode). Returns company id/name/slug/subdomain/domain/s... |

### Tenant Auth
**Description:** Tenant authentication (login, refresh token)
**Endpoints:** 2 | **Paths:** 2
**Methods:** POST=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/auth/login` | Tenant user login (employee or company admin) | Public login for tenant users (employees) stored in the tenant DB, with fallback to platform users (company_admin) bound to the company. Company id... |
| `POST` | `/api/v1/tenant/auth/refresh` | Refresh tenant access token | Public endpoint to exchange a refresh token for a new access token. |

### Tenant: Company
**Description:** Self-service company endpoints for the authenticated tenant user
**Endpoints:** 2 | **Paths:** 1
**Methods:** GET=1 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/companies/me` | Get current company detail | Retrieve the profile of the company the authenticated tenant user belongs to. Company is resolved from the tenant context (X-Tenant-ID / JWT claims). |
| `PUT` | `/api/v1/tenant/companies/me` | Update current company information | Update the tenant's own company profile (email, phone, address, NPWP, NIB). Company is resolved from the tenant context; name/subdomain/domain are ... |
