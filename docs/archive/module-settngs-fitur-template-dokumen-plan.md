> ✅ **Diarsipkan (2026-08-21)**: fitur ini sudah selesai secara kode — seluruh Phase 1-6 (§21) dan Acceptance Criteria (§22) tercentang `[x]`. 3 checkbox tersisa di Phase 4 ("Template DOCX rendering", "DOCX generation", "Generated document storage") adalah forward-reference basi ke Phase 5 yang sudah selesai, bukan gap nyata.
>
> ⚠️ **Sisa kerja — tugas operasional, bukan kode**: `[ ] LibreOffice Headless setup di server produksi` (§21 Phase 4) belum dikonfirmasi. Binary LibreOffice perlu di-install manual di server produksi (lihat §15 untuk auto-detect path per platform + env override `HRIS_STORAGE_LIBREOFFICE_PATH`). Ada mitigasi: engine alternatif pure-Go `docx2pdf` (`storage.pdf_engine = "docx2pdf"`) tersedia tanpa dependency LibreOffice bila setup produksi belum kelar — lihat §15 baris "Dua opsi engine (2026-08-14)". Verifikasi mana yang dipakai di produksi sebelum benar-benar menganggap deployment selesai.
>
> ✅ **Diarsipkan sebagai selesai** — item di atas adalah catatan operasional untuk tim infra/deployment, bukan alasan menunda status "selesai" dokumen ini.

# Plan Fitur: Template Dokumen

**Module:** Settings  
**Feature:** Template Dokumen  
**Stack:** Go Backend + Vue.js 3 + PrimeVue Frontend  
**Template Format:** Microsoft Word `.docx`  
**PDF Engine:** LibreOffice Headless

## 1. Tujuan

Membangun fitur **Template Dokumen** di dalam module **Settings** untuk memungkinkan administrator membuat, mengelola, dan menggunakan template dokumen HR yang dapat digenerate menjadi PDF.

Use case awal:

- Perjanjian Karyawan Kontrak
- SK Movement

Fitur dirancang agar dapat dikembangkan untuk dokumen HR lainnya tanpa perubahan besar pada arsitektur.

---

## 2. Scope

### 2.1 Template Management

Administrator dapat:

- Melihat daftar template.
- Membuat template baru.
- Mengubah template.
- Mengaktifkan/nonaktifkan template — **hanya 1 template aktif per jenis dokumen**.
- Melihat detail template.
- Membuat versi baru dari template.
- Melihat histori versi template.
- Melakukan preview template.
- Menghapus template yang belum digunakan.

### 2.2 Document Type

Template memiliki jenis dokumen.

Contoh:

```text
CONTRACT_AGREEMENT
MOVEMENT_SK
```

Struktur dibuat extensible sehingga nantinya dapat ditambahkan:

```text
EMPLOYMENT_LETTER
WARNING_LETTER
PROMOTION_SK
TRANSFER_SK
TERMINATION_LETTER
CERTIFICATE
```

Setiap jenis dokumen memiliki:

```text
maksimal 1 template ACTIVE
```

---

## 2.3 Aturan Template Aktif

### Satu Template Aktif per Jenis Dokumen

- Setiap jenis dokumen hanya boleh memiliki **1 template dengan status `Active`**.
- Mengaktifkan sebuah template otomatis **menonaktifkan template lain** dengan jenis dokumen yang sama (dalam satu transaksi).
- Template `Inactive` **tidak dapat digunakan** untuk Generate Document.
- Aturan dipaksa di dua lapisan:
  - **Database** — partial unique index `(document_type) WHERE status = 'ACTIVE'`.
  - **Service** — validasi + transaksi saat aktivasi.

> **Keputusan (2026-08-14): fitur template default (referensi) DIHAPUS.** Tidak ada lagi template default/referensi, alur "Gunakan Template Default", maupun "Edit Default Content". Template dibuat langsung dari nol (ACTIVE/INACTIVE). Migrasi 113 menghapus seed default + kolom `is_default` + index partial default; `CreateFromDefault`/`UpdateDefaultContent` dan route `/from-default`/`/default-content` dihapus dari backend, UI terkait dihapus dari frontend.

---

# 3. Struktur Menu

Template masuk ke dalam **Settings** yang sudah tersedia.

```text
Settings
├── General
├── Organization
├── ...
└── Template Dokumen
```

Menu:

```text
Template Dokumen
```

Icon PrimeVue yang direkomendasikan:

```text
pi-file-edit
```

---

# 4. Template List

Halaman:

```text
Settings → Template Dokumen
```

Tampilan menggunakan PrimeVue `DataTable`.

Kolom:

| Kolom | Keterangan |
|---|---|
| Name | Nama template |
| Code | Kode template |
| Document Type | Jenis dokumen |
| Version | Versi aktif |
| Status | Active / Inactive |
| Updated At | Waktu terakhir update |
| Action | Detail/Edit/Version/Preview |

Catatan:

- Hanya **1 template `Active`** per Document Type. Mengaktifkan template lain otomatis menonaktifkan template yang sedang aktif (dengan konfirmasi).

Contoh:

```text
Template Dokumen

[ + Template Baru ]                       [ Search ]

┌─────────────────────────────────────────────────────────────────────┐
│ Name                 │ Type        │ Version │ Status    │          │
├─────────────────────────────────────────────────────────────────────┤
│ Perjanjian PKWT      │ Contract    │ v3       │ Active    │          │
│ SK Movement          │ Movement SK │ v2       │ Active    │          │
└─────────────────────────────────────────────────────────────────────┘
```

---

# 5. Template Form

Form pembuatan template:

```text
Template Information

Name
[ Perjanjian Kerja Waktu Tertentu ]

Code
[ CONTRACT_AGREEMENT ]

Document Type
[ Contract Agreement                ▼ ]

Status
[ Inactive ]   ← hanya 1 template Active per jenis dokumen;
                  mengaktifkan otomatis menonaktifkan yang lain.

Template Content
┌──────────────────────────────────────────────┐
│              Microsoft Word                 │
│                                              │
│ PERJANJIAN KERJA WAKTU TERTENTU              │
│                                              │
│ Nomor: {{contract.number}}                   │
│                                              │
│ Nama: {{employee.name}}                      │
│ Jabatan: {{employee.position}}               │
│                                              │
└──────────────────────────────────────────────┘

[Cancel] [Preview] [Save]
```

---

# 6. Microsoft Word DOCX Template

Template dokumen dibuat menggunakan **Microsoft Word** dan di-upload ke aplikasi dalam format `.docx`. Tidak perlu membangun document editor Word di Vue.

Contoh isi template Word:

```text
PERJANJIAN KERJA WAKTU TERTENTU

Nomor: {{contract.number}}

Nama       : {{employee.name}}
NIK        : {{employee.nik}}
Jabatan    : {{employee.position}}
Organisasi : {{employee.organization.name}}

Perjanjian berlaku mulai tanggal {{contract.start_date}}
sampai dengan {{contract.end_date}}.
```

Administrator dapat menggunakan fitur Microsoft Word seperti:

- Font dan font size
- Bold / italic / underline
- Alignment dan paragraph spacing
- Table
- Image/logo
- Header/footer
- Page number
- Page break
- Signature
- Kop surat
- Numbering dan bullet list

Frontend Vue + PrimeVue hanya menangani management template, upload, versioning, variable reference, preview, dan document actions.

## 6.1 Upload Template

Form template menyediakan:

```text
Name
[ Perjanjian Kerja Waktu Tertentu ]

Code
[ CONTRACT_AGREEMENT ]

Document Type
[ Contract Agreement ▼ ]

Template File
[ Choose File ]

Perjanjian_PKWT.docx

[Cancel] [Upload]
```

Validasi minimal:

- Extension `.docx`.
- MIME type valid.
- Ukuran file sesuai limit.
- File valid sebagai Office Open XML.
- Placeholder dapat dideteksi.
- Variable yang digunakan harus terdaftar atau menghasilkan warning sesuai konfigurasi.

# 7. Variable / Placeholder System

Template menggunakan placeholder:

```text
{{employee.name}}
{{employee.nik}}
{{employee.position}}
{{employee.organization}}
{{contract.number}}
{{contract.start_date}}
{{contract.end_date}}
{{company.name}}
```

Jangan mengharuskan administrator mengetik placeholder secara manual.

Tambahkan fitur:

```text
[ Insert Variable ]
```

Contoh:

```text
Insert Variable

Employee
├── Employee Number
├── Name
├── NIK
├── Join Date
├── Position
└── Organization

Contract
├── Contract Number
├── Contract Type
├── Start Date
└── End Date

Movement
├── Movement Number
├── Movement Type
├── Effective Date
├── Previous Position
└── New Position

Company
├── Name
├── Address
└── Phone
```

Ketika administrator memilih:

```text
Employee → Name
```

editor akan memasukkan:

```text
{{employee.name}}
```

---

# 8. Template Preview (Preview ke PDF)

Preview dilakukan melalui backend dengan sample data. Karena template berbasis DOCX, preview mengikuti pipeline yang sama dengan dokumen final.

```text
DOCX Template
      ↓
Sample Data
      ↓
Resolve Variables
      ↓
Generated DOCX
      ↓
LibreOffice Headless
      ↓
PDF
      ↓
PDF Viewer
```

### 8.1 Preview Versi Tersimpan

Template list menyediakan tombol `[ Preview ]` yang menggunakan active version template.

### 8.2 Preview Draft

Opsional untuk enhancement: administrator dapat memilih file DOCX draft dan meminta preview sebelum version disimpan.

### 8.3 Sample Data

```text
{{employee.name}}        → John Fieldman
{{employee.nik}}         → 199001012015011001
{{employee.position}}    → HR Staff
{{contract.number}}      → CTR-2026-001
{{contract.start_date}}  → 2026-01-01
{{contract.end_date}}    → 2027-01-01
{{company.name}}         → PT Maju Bersama
```

### 8.4 PDF Viewer

Preview ditampilkan dalam PrimeVue `Dialog`/`Drawer` menggunakan PDF embed/iframe dengan action:

```text
[ Download PDF ] [ Print ] [ Close ]
```

# 9. Document Configuration

Layout utama dibuat di Microsoft Word. Konfigurasi aplikasi hanya menyimpan metadata yang memang diperlukan document pipeline.

Optional:

```text
paper_size
orientation
```

Contoh:

```text
A4
Portrait
```

Margin, header, footer, font, spacing, signature, dan layout lainnya mengikuti template Word.

# 10. Template Versioning

Template harus menggunakan versioning.

Contoh:

```text
Perjanjian PKWT

v1
v2
v3 ← Active
```

Ketika template aktif diubah:

```text
Current Version
      ↓
Create New Version
      ↓
v4
      ↓
Set Active
```

Jangan mengubah isi versi lama secara langsung.

Catatan: status **Active/Inactive** berada di level **template** (maksimal 1 Active per jenis dokumen), sedangkan **versi aktif** berada di level template (`active_version_id`).

---

# 11. Database Design

## `document_templates`

```text
id UUID PK
name
code
document_type
description
active_version_id
status              -- ACTIVE / INACTIVE
created_at
updated_at
deleted_at
```

## `document_template_versions`

```text
id UUID PK
template_id UUID FK
version
content
paper_size
orientation
margin_top
margin_right
margin_bottom
margin_left
created_by
created_at
```

## `document_template_variables`

Jika variable perlu dikonfigurasi secara dinamis:

```text
id UUID PK
template_id UUID FK
key
label
category
source
created_at
updated_at
```

Variable standar seperti:

```text
employee.name
employee.nik
contract.number
```

sebaiknya didefinisikan melalui registry di backend, bukan sepenuhnya melalui database.

### Constraints & Seed Data

```sql
-- Hanya 1 template aktif per jenis dokumen
CREATE UNIQUE INDEX uq_document_templates_active_per_type
  ON document_templates (document_type)
  WHERE status = 'ACTIVE';

```

> Fitur template default/referensi dihapus (2026-08-14): kolom `is_default` dan index partial default di-drop pada migrasi 113; tidak ada seed template default.

---

# 12. Generated Document

Buat entity terpisah untuk dokumen yang sudah digenerate.

## `generated_documents`

```text
id UUID PK
template_id UUID FK
template_version_id UUID FK
document_type
reference_type
reference_id
file_name
file_path
mime_type
generated_by
generated_at
created_at
```

Contoh:

```text
template
    ↓
Perjanjian PKWT v3
    ↓
Contract #CTR-2026-001
    ↓
Generated Document
    ↓
Perjanjian_PKWT_CTR-2026-001.pdf
```

---

# 13. Backend Go

**Template Dokumen adalah sub-feature dari module `setting`** — bukan module tenant terpisah. Route-nya didaftarkan di bawah group `/settings` pada module `setting` (sama seperti resource settings lain: zones, provinces, competencies, document-numbering), sehingga path API menjadi `/api/v1/tenant/settings/document-templates`.

Struktur feature (package flat, sesuai konvensi project):

```text
backend/internal/modules/documenttemplate/
    ├── model.go        -- DocumentTemplate, DocumentTemplateVersion, DocumentTemplateAudit, GeneratedDocument
    ├── repository.go   -- repo + NewTenantDBResolver
    ├── service.go      -- one-active-per-jenis atomik, versioning
    ├── handler.go      -- HTTP handlers
    ├── routes.go       -- RegisterRoutes(rg, handler) → sub-group "/document-templates"
    ├── dto.go          -- request/response DTO + binding tags
    ├── errors.go
    └── variables.go    -- static variable registry
```

Wiring ke module `setting`:

```text
setting.NewModule(dbManager, logger, numberingSvc)
    ├── dtRepo := documenttemplate.NewRepository(documenttemplate.NewTenantDBResolver(dbManager))
    ├── dtSvc  := documenttemplate.NewService(dtRepo, logger)
    ├── dtHandler := documenttemplate.NewHandler(dtSvc)
    └── RegisterRoutesWithNumbering(rg, handler, numberingHandler, dtHandler)
        └── documenttemplate.RegisterRoutes(settings, dtHandler)  → /settings/document-templates/*
```

Tidak ada lagi registrasi `documenttemplate.NewModule` di `cmd/server/main.go` (di-remove; module.go dihapus). Permission & menu dideklarasikan di `setting` module `Info()`.

---

# 14. API

Semua endpoint di bawah prefix tenant: `/api/v1/tenant/settings/document-templates`.

### Template

```http
GET    /api/v1/tenant/settings/document-templates
POST   /api/v1/tenant/settings/document-templates
GET    /api/v1/tenant/settings/document-templates/{id}
PUT    /api/v1/tenant/settings/document-templates/{id}
DELETE /api/v1/tenant/settings/document-templates/{id}

POST   /api/v1/tenant/settings/document-templates/{id}/activate   ← otomatis menonaktifkan
                                                                    template lain sejenis
POST   /api/v1/tenant/settings/document-templates/{id}/deactivate
```

### Version

```http
GET  /api/v1/tenant/settings/document-templates/{id}/versions
POST /api/v1/tenant/settings/document-templates/{id}/versions
GET  /api/v1/tenant/settings/document-templates/{id}/versions/{versionId}
```

### Preview

```http
POST /api/v1/tenant/settings/document-templates/{id}/preview   ← preview VERSI AKTIF (data contoh)

POST /api/v1/tenant/settings/document-templates/preview-draft  ← preview DRAFT:
                                                               body: { document_type,
                                                                       content,
                                                                       paper_size,
                                                                       orientation,
                                                                       margins }

Response: { "pdf_url": "...", "file_name": "preview_<template>.pdf" }
```

### Variables

```http
GET /api/v1/tenant/settings/document-templates/variables
```

> **RBAC note:** karena route berada di bawah `/settings/`, middleware RBAC menurunkan resource `setting` (lihat `singularize("settings") → "setting"` di `authz.ResourceFromPath`), sehingga permission `setting.*` yang sudah di-seed di tenant RBAC langsung berlaku — Admin dapat akses penuh, Employee view-only. Ini memperbaiki masalah sebelumnya saat route berdiri sendiri di `/document-templates` yang menurunkan resource `document-template` (tidak pernah di-seed).

---

# 15. PDF Generation

Backend Go menangani:

```text
Template DOCX
   ↓
Load Template Version
   ↓
Load Reference Data
   ↓
Resolve Variables
   ↓
Generate DOCX
   ↓
LibreOffice Headless
   ↓
PDF
   ↓
Store PDF
   ↓
Return File
```

Gunakan **LibreOffice Headless** sebagai engine DOCX → PDF.

Contoh konsep:

```bash
libreoffice \
    --headless \
    --convert-to pdf \
    --outdir /tmp/output \
    document.docx
```

Untuk production, LibreOffice dapat ditempatkan sebagai service/container terpisah untuk isolasi dan scaling.

**Setup binary (implementasi saat ini):** `LibreOfficePDFService` auto-detect binary per platform — Windows mencoba `C:\Program Files\LibreOffice\program\soffice.exe` & `C:\Program Files (x86)\...` lalu nama `soffice(.exe)`/`libreoffice` di PATH; Linux/macOS mencoba `/usr/bin/libreoffice` dll. Bila tidak ditemukan, preview mengembalikan 503 `PDF_ENGINE_NOT_CONFIGURED` dengan pesan jelas. Override via env `HRIS_STORAGE_LIBREOFFICE_PATH` atau `storage.libreoffice_path` di config file. Catatan: installer LibreOffice Windows TIDAK menambahkan ke PATH — gunakan path penuh atau env var.

**Dua opsi engine (2026-08-14):** selain LibreOffice, tersedia `Docx2pdfPDFService` — implementasi pure-Go (`github.com/bobyeoh/docx2pdf-go` v0.3.0, MIT) tanpa dependency eksternal/LibreOffice. Pilih via `storage.pdf_engine` (`"libreoffice"` default | `"docx2pdf"`; env `HRIS_STORAGE_PDF_ENGINE`), helper `newPDFService(cfg)` di main.go. Keduanya tetap ada; service dibuat sesuai pilihan config. Prasyarat: proyek di-upgrade ke **Go 1.26.1** (go.mod + `golang:1.26-alpine` di Dockerfile) karena library menuntut Go ≥ 1.26.

## 15.1 PDF Service Abstraction

```go
type PDFService interface {
    ConvertDOCXToPDF(
        ctx context.Context,
        inputPath string,
        outputPath string,
    ) error
}
```

Implementasi awal:

```text
PDFService
    ↓
LibreOfficePDFService
```

# 16. Integrasi dengan Contract

Pada module Contract:

```text
Contract Detail

[ Generate Document ▼ ]
```

Pilihan:

```text
Generate Document
├── Perjanjian PKWT
└── ...
```

Flow:

```text
Contract
   ↓
Select Template (default pilihan: template aktif untuk CONTRACT_AGREEMENT)
   ↓
Get Template Version
   ↓
Resolve Contract + Employee Data
   ↓
Generate PDF
   ↓
Download / Preview / Print
```

Jika **belum ada template aktif** untuk `CONTRACT_AGREEMENT`, tampilkan pesan:

```text
Belum ada template aktif untuk jenis ini.
Buat template di Settings → Template Dokumen.
```

---

# 17. Integrasi dengan Employee Movement

Pada detail Movement:

```text
Employee Movement

[ Generate Document ▼ ]
```

Contoh:

```text
Generate Document
└── SK Movement   ← default pilihan: template aktif untuk MOVEMENT_SK
```

Jika **belum ada template aktif** untuk `MOVEMENT_SK`, tampilkan pesan yang sama seperti pada Contract (arahkan ke Settings → Template Dokumen).

Variable otomatis mengambil:

```text
Employee
Previous Organization
Previous Position
New Organization
New Position
Effective Date
Movement Number
```

---

# 18. Security & Permission

Permission didaftarkan di **module `setting`** (`Info().Permissions`) dengan nama resource `setting.document_template` (pola `module.resource.action` yang sama seperti `setting.zone.*`):

```text
setting.document_template.view
setting.document_template.create
setting.document_template.update
setting.document_template.delete
setting.document_template.activate
setting.document_template.deactivate
setting.document_template.version
```

Untuk generated document (Phase 5):

```text
setting.document_template.generate
setting.document_template.download
```

> **Catatan penting:** permission granular di atas adalah deklarasi untuk RBAC UI. Enforcement runtime tetap melalui `authz` middleware yang menurunkan resource dari path — karena route berada di `/settings/...`, permission yang dicek adalah `setting.view/create/update/delete` yang sudah di-seed di tenant RBAC (Admin: semua action; Employee: view). Jadi Admin tenant otomatis bisa mengelola template dokumen tanpa konfigurasi RBAC tambahan.

Gunakan mekanisme permission/RBAC yang sudah digunakan aplikasi.

---

# 19. Audit Log

Catat aktivitas penting:

```text
Template Created
Template Updated
Template Version Created
Template Activated
Template Deactivated
Template Deleted
Document Generated
Document Downloaded
```

Contoh:

```text
User
    ↓
Create Template Version
    ↓
Template: SK Movement
    ↓
Version: v3
    ↓
Timestamp
```

---

# 20. Frontend Components

Komponen yang disarankan:

```text
DocumentTemplateIndex.vue
DocumentTemplateForm.vue
DocumentTemplateUpload.vue
DocumentTemplateVariableList.vue
DocumentTemplatePreview.vue         ← preview hasil DOCX → PDF
DocumentTemplatePreviewDialog.vue   ← PDF viewer (iframe + Download/Print)
DocumentTemplateVersionDialog.vue
DocumentTemplateConfiguration.vue
```

PrimeVue component:

```text
DataTable
Dialog
Drawer
Button
InputText
Select
Textarea
Tag
Menu
Popover
Tabs
ConfirmDialog
Toast
ProgressSpinner
```

---

# 21. Tahapan Implementasi

## Phase 1 — Database & Backend Foundation ✅ SELESAI (2026-08-14, termasuk final review fix: MySQL locking pada Activate, partial-update tidak lagi menghapus field lain, soft-delete code reuse, case-insensitive search lintas DB, FK active_version_id via migrasi 111)

- [x] Migration `document_templates` — migrasi 110 (mysql+postgres), commit `6b12129c`
- [x] Partial unique index: 1 template aktif per jenis dokumen — Postgres partial index (migrasi 110); MySQL ditegakkan di service layer (Task 4)
- [x] Migration `document_template_versions` — migrasi 110
- [x] **Fitur template default dihapus** — migrasi 113 (mysql+postgres): DELETE seed default, DROP index partial default (postgres), DROP COLUMN `is_default`; `CreateFromDefault`/`UpdateDefaultContent`/`FindDefaultByType`/guard `IsDefault` + route `/from-default` & `/:id/default-content` + permission `set_default` dihapus dari backend
- [ ] Migration `document_template_variables` jika diperlukan — tidak dibuat; variable registry diimplementasikan sebagai static Go registry (Task 6), bukan tabel DB, sesuai catatan spec §11 ("sebaiknya didefinisikan melalui registry di backend")
- [x] Migration `generated_documents` — migrasi 110 (skema saja; belum ada writer, itu Phase 5)
- [x] Domain entity — `backend/internal/modules/documenttemplate/model.go`, commit `0313eb26`
- [x] Repository — `backend/internal/modules/documenttemplate/repository.go`, commit `1d48b854`
- [x] Service — `backend/internal/modules/documenttemplate/service.go`, commit `4182bc73` (one-active-per-jenis atomik via transaksi, versioning; from-default dihapus migrasi 113)
- [x] DTO — `backend/internal/modules/documenttemplate/dto.go`, commit `b5616b79`
- [x] API handler — `handler.go` + `routes.go`, commit `b5616b79`
- [x] **Terintegrasi ke module `setting`** — handler documenttemplate di-wire di `setting.NewModule` dan route didaftarkan di `/settings/document-templates` (bukan module tenant terpisah); registrasi `documenttemplate.NewModule` & `module.go` dihapus dari `main.go` (fix RBAC: route di `/settings/...` menurunkan resource `setting.*` yang sudah di-seed tenant)
- [x] Permission — daftar permission `setting.document_template.*` di `setting` module's `Info()`, commit `b5616b79`
- [x] Validation — `httputil.BindAndValidate` + tag binding di dto.go, commit `b5616b79`

## Phase 2 — Template Management ✅ SELESAI (2026-08-14)

- [x] Template list — `frontend/tenant/src/views/settings/DocumentTemplatesView.vue`, commit `95cd6d26`
- [x] Search — search box dengan debounce 400ms, memanggil `?search=`
- [x] Filter document type — dropdown Contract Agreement / Movement SK
- [x] Create template — **halaman terpisah** `DocumentTemplateForm.vue` (`/settings/document-templates/new`; name, document_type, description) → `POST /document-templates`. **Code otomatis**: tidak ditampilkan di form; backend meng-generate `TMPL-{DOC_TYPE}-{RANDOM8}` (UUID pendek uppercase) saat `code` kosong, sehingga admin tidak perlu mengisi kode
- [x] Edit template — **halaman terpisah** `DocumentTemplateForm.vue` (`/settings/document-templates/:id/edit`; name, description; document_type readonly) → `PUT /document-templates/{id}`
- [x] Detail template — dialog info + preview konten versi aktif + tombol Versions
- [x] Active/inactive — tombol Activate/Deactivate per baris + konfirmasi, memanggil endpoint Phase 1 yang sudah teratomisasi (hanya 1 aktif per jenis)
- [x] Delete template — aksi per baris + ConfirmDeleteDialog → `DELETE /document-templates/{id}`
- [x] Version management — dialog list versi + buat versi baru (konten + paper/orientation/margin) + detail versi (`GET /{id}/versions`, `POST /{id}/versions`, `GET /{id}/versions/{versionId}`)

## Phase 3 — DOCX Template Management

- [x] Upload template `.docx` — form halaman `DocumentTemplateForm.vue` (`/settings/document-templates/new` & `/:id/edit`) + endpoint backend `POST /{id}/versions` mode multipart (field `file`); file disimpan ke `{uploadDir}/document_templates/{uuid}.docx`, `content` = path relatif, `file_name` = nama asli (kolom baru migrasi 112); response versi membawa `file_url`
- [x] DOCX file validation — ekstensi `.docx` + ukuran ≤ 10 MB (backend `documenttemplate.file_*` locale EN/ID; frontend validasi `.docx` sebelum submit)
- [x] Template version upload — dialog New Version di `DocumentTemplatesView.vue` sekarang upload `.docx` (multipart), bukan textarea HTML
- [x] Template download — link download file `.docx` tersimpan di form edit, detail template, dan detail versi
- [x] Placeholder detection — `docx.go` membaca file .docx (zip OOXML via stdlib `archive/zip`), ekstrak `{{key}}` dari semua XML `word/`, dikembalikan di response `POST /{id}/versions` (`placeholders`)
- [x] Variable validation — placeholder yang tidak terdaftar di `VariableRegistry()` ditolak (400 `documenttemplate.unknown_variables` berisi daftar); file non-zip ditolak (400 `documenttemplate.invalid_docx`); test `TestHandlerCreateVersionRejectsUnknownPlaceholders` & `TestHandlerCreateVersionRejectsInvalidDocx`
- [x] Variable reference UI — card "Variable Reference" di `DocumentTemplateForm.vue` menampilkan grup variabel dari `GET /variables` (Employee/Contract/Movement/Company)
- [x] Copy variable action — klik variabel → salin `{{key}}` ke clipboard (`navigator.clipboard` + toast); toast info jumlah variable terdeteksi setelah save
- [x] ~~Template default/reference flow~~ — fitur dihapus (migrasi 113)

> **Catatan:** Microsoft Word bukan bagian dari editor aplikasi; PrimeVue Editor/Quill tidak digunakan sebagai document authoring tool. Template dibuat menggunakan Microsoft Word dan di-upload sebagai `.docx`.

> **Perubahan implementasi (2026-08-14):** editor Quill (Phase 3 lama) dihapus dari `DocumentTemplateForm.vue` — card Template Content diganti card Template File (upload `.docx`), variable picker/page-break/table custom dihapus, CSS `.page-break` dihapus dari main.css. Document Configuration disederhanakan sesuai plan §9 (hanya `paper_size` + `orientation`; margin/header/footer mengikuti Word). Backend: `POST /{id}/versions` menerima multipart `.docx` (mode JSON `content` tetap didukung untuk backward compat), kolom `file_name` ditambah via migrasi 112 (mysql+postgres), `file_url` dihitung dari `content` bila berupa path file.

## Phase 4 — Preview & PDF

- [x] Variable resolver ({{key}} → sample data) — `docx.go` `resolveDocxVariables`
- [x] Sample data preview — `sampleData()` di `docx.go`
- [x] DOCX → PDF conversion — `LibreOfficePDFService` (`pdf_service.go`)
- [x] API preview versi tersimpan — `POST /{id}/preview` (resolve → convert → simpan ke `{uploadDir}/previews/`)
- [x] PDF viewer dialog di Settings — iframe dialog di `DocumentTemplatesView.vue`
- [x] Download PDF dari preview
- [x] Print PDF dari preview
- [x] Loading state selama generasi PDF
- [ ] Template DOCX rendering (Generate Document modul lain — Phase 5)
- [ ] DOCX generation (Generate Document modul lain — Phase 5)
- [ ] LibreOffice Headless setup di server produksi
- [ ] Generated document storage (bagian Generate Document — Phase 5)

## Phase 5 — Module Integration

### Contract

- [x] Generate Contract Agreement — `POST /contracts/:id/generate-document` (employeemovement) → shared `documenttemplate.Generator`
- [x] Template selection — otomatis memakai template ACTIVE untuk `CONTRACT_AGREEMENT` (max 1 per jenis; tanpa template aktif → pesan jelas)
- [x] Contract data mapping — `contract.number`, `contract.start_date`, `contract.end_date` dari `employee_contracts`
- [x] Employee data mapping — `employee.name`/`employee.nik` dari `employees`, `employee.position`/`employee.organization` dari employment aktif terakhir
- [x] Generated document history — `GET /contracts/:id/generated-documents` (tabel `generated_documents`, migration 110)

### Movement

- [x] Generate SK Movement — `POST /movements/:id/generate-document` (hanya status approved/executed)
- [x] Template selection — otomatis memakai template ACTIVE untuk `MOVEMENT_SK`
- [x] Movement data mapping — `movement.number`, `movement.effective_date`, `movement.previous_position`, `movement.new_position` dari snapshot movement
- [x] Employee data mapping — `employee.*` + posisi/org tujuan (snapshot)
- [x] Generated document history — `GET /movements/:id/generated-documents`

### Shared Document Generator (documenttemplate)

- [x] `Generator.Generate` — template aktif + versi aktif → resolve variabel → DOCX→PDF (LibreOffice) → simpan ke `{uploadDir}/generated_documents/` → catat `generated_documents` + audit `DOCUMENT_GENERATED`
- [x] `company.name`/`company.address` diisi Generator via `CompanyProvider` (platform DB, di-wire main.go)
- [x] Wiring: `employeeMovementSvc.SetDocumentGenerator(documentGeneratorAdapter{gen})` — narrow interface + adapter (pola sama ApprovalEngine/CareerExecutor)
- [x] Frontend: tombol Generate + PDF viewer (iframe/Download/Print) + histori di `EmployeeMovements.vue` (detail dialog) & `EmployeeContracts.vue` (per baris)

## Phase 6 — Testing ✅ SELESAI (2026-08-14)

### Backend

- [x] Template CRUD test — `service_test.go` (`TestServiceCreateRejectsInvalidDocumentType`, `TestServiceCreateAutoGeneratesCodeWhenEmpty`, `TestServiceCreateRejectsDuplicateCode`, `TestServiceDeleteThenRecreateWithSameCode`), `repository_test.go` (`TestRepositoryCreateAndGetByID`, `TestRepositoryGetByIDNotFound`, `TestRepositoryListPaginationAndSearch`, `TestRepositorySoftDeleteExcludesFromList`), `handler_test.go` (`TestHandlerCreateAndList`)
- [x] Versioning test — `service_test.go` (`TestServiceCreateVersionIncrementsAndSetsActiveVersion`, `TestServiceCreateVersionRejectsNonexistentTemplate`), `repository_test.go` (`TestRepositoryVersionsCreateListNextNumber`), `handler_test.go` (`TestHandlerUpdateDefaultContentAndVersionDetail`)
- [x] Variable resolver test — `docx_test.go` (`TestExtractDocxPlaceholders`, `TestExtractDocxPlaceholdersRejectsNonZip`, `TestUnknownPlaceholders`, `TestResolveDocxVariables`, `TestSampleDataCoversRegistry`, `TestResolveDocxVariablesToleratesBadCRC`), `variables_test.go` (`TestVariableRegistryHasExpectedCategories`, `TestVariableRegistryKeysAreDotted`)
- [x] Permission test — `setting/permission_test.go` (`TestModulePermissionsDeclareDocumentTemplate` verifikasi deklarasi `setting.document_template.*` di `Info().Permissions` + `TestModuleMenusIncludeDocumentTemplates`)
- [x] PDF generation test — `pdf_service_test.go` (4 test resolver binary LibreOffice), `docx2pdf_service_test.go` (`TestDocx2pdfPDFServiceConvertsDocx` — konversi DOCX→PDF nyata tanpa LibreOffice), `handler_test.go` (`TestHandlerPreviewWithMockPDF`, `TestHandlerPreviewWithoutPDFEngine`)
- [x] Generated document test — `generator_test.go` (`TestGeneratorGenerateCreatesPDFAndRecord`, `TestGeneratorGenerateNoActiveTemplate`, `TestGeneratorGenerateWithoutPDFEngine`) + `employeemovement/document_test.go` (generate movement/contract + history)

### Frontend (baru — vitest + @vue/test-utils + jsdom)

- [x] Test framework — `vitest` v4 + `@vue/test-utils` v2 + `jsdom` di devDependencies; script `npm test` (`vitest run`); konfigurasi di `vite.config.js` (`test.environment = 'jsdom'`, `globals`, `setupFiles: tests/setup.js`); stub global PrimeVue (DataTable/Column render-function yang meneruskan `data` row ke body slot, Button dengan `emits` agar tidak double-fire, Dialog/ConfirmDeleteDialog pakai prop `visible`, Select merender options)
- [x] Template list test — `tests/DocumentTemplatesView.test.js`: load list saat mounted, pesan kosong, filter/status label, search
- [x] Create template test — `tests/DocumentTemplateForm.test.js`: validasi nama wajib, validasi document_type, submit → POST template + POST versi (FormData) + redirect
- [x] DOCX upload/validation test — tolak non-`.docx`, tolak > 10MB, valid `.docx` diterima
- [x] Variable insertion test — variable reference dimuat dari `GET /variables`, klik variabel → clipboard `{{key}}` + toast
- [x] Preview test — klik preview → POST `/{id}/preview` → iframe `pdf_url` tampil; error ditampilkan bila gagal
- [x] Versioning test — dialog daftar versi (`GET /{id}/versions`), detail versi, buat versi baru (upload file + paper/orientation)
- [x] PDF download test — tombol download link `:href` + `:download` pada preview & file template (via stub `<a>` assertion di detail/version dialog)
- [x] Activate/deactivate/delete test — konfirmasi dialog → POST activate/deactivate + DELETE → reload + toast

---

# 22. Acceptance Criteria

Fitur dianggap selesai apabila:

- [x] Administrator dapat membuat template dokumen dari Settings.
- [x] Administrator dapat upload template Microsoft Word `.docx`.
- [x] Administrator dapat melihat dan menyalin daftar variable untuk digunakan di Microsoft Word.
- [x] Variable dapat di-resolve berdasarkan data employee/contract/movement.
- [x] Template memiliki versioning.
- [x] Template aktif dapat digunakan oleh module Contract.
- [x] Template aktif dapat digunakan oleh module Employee Movement.
- [x] Contract dapat menghasilkan Perjanjian PDF.
- [x] Movement dapat menghasilkan SK PDF.
- [x] PDF dapat di-preview langsung dari Settings (versi tersimpan maupun draft).
- [x] Preview menggunakan pipeline DOCX → PDF yang sama dengan Generate Document.
- [x] PDF dapat di-download.
- [x] PDF dapat di-print.
- [x] Dokumen yang sudah dihasilkan menyimpan referensi template dan versinya.
- [x] Perubahan template tidak mengubah dokumen yang sudah diterbitkan.
- [x] Hanya ada 1 template aktif per jenis dokumen; mengaktifkan template lain otomatis menonaktifkan template sebelumnya.
- [x] Permission dan audit log diterapkan.

---

# 23. Target Arsitektur Akhir

```text
                         SETTINGS
                            │
                     TEMPLATE DOKUMEN
                            │
                     Word DOCX Template
                            │
              ┌─────────────┴─────────────┐
              │                           │
        Contract Template          Movement Template
              │                           │
              └─────────────┬─────────────┘
                            │
                    Template Version
                            │
                       Go Backend
                            │
                  Variable Registry
                            │
                  Variable Resolver
                            │
                       DOCX Output
                            │
                            ▼
                 LibreOffice Headless
                            │
                            ▼
                           PDF
                            │
             ┌──────────────┼──────────────┐
             │              │              │
           Preview       Download         Print
```

## Prinsip Desain

`Template Dokumen` tetap berada di **Settings**, sedangkan proses **Generate Document** tetap berada di module bisnis masing-masing.

Dengan demikian:

- **Settings** bertanggung jawab atas template, variable, konfigurasi, dan versioning.
- **Contract** bertanggung jawab menentukan kapan Perjanjian Kontrak dibuat dan data contract yang digunakan.
- **Employee Movement** bertanggung jawab menentukan kapan SK Movement dibuat dan data movement yang digunakan.
- **Document Generator** bertanggung jawab melakukan rendering dan menghasilkan PDF.
- Dokumen yang sudah diterbitkan menyimpan **template version** yang digunakan sehingga perubahan template tidak mengubah dokumen lama.


## Keputusan Teknis

> **Microsoft Word `.docx` digunakan sebagai template authoring tool.**

> **Vue.js + PrimeVue digunakan untuk management template, upload, versioning, variable reference, preview, dan document actions.**

> **Go digunakan sebagai backend, variable resolver, DOCX processor, dan orchestration document generation.**

> **LibreOffice Headless digunakan sebagai engine konversi DOCX → PDF.**

Pendekatan ini dipilih karena paling sesuai untuk dokumen formal seperti **SK, perjanjian kontrak, surat mutasi, surat promosi, dan dokumen HR lainnya**, sekaligus memungkinkan tim HR mengatur layout menggunakan Microsoft Word yang sudah familiar.
