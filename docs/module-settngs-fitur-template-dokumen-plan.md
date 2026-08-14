# Plan Fitur: Template Dokumen

**Module:** Settings  
**Feature:** Template Dokumen  
**Stack:** Go (Backend) + Vue.js 3 + PrimeVue (Frontend)

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
- Menggunakan **template default (referensi)** untuk membuat template baru.
- Mengubah konten **template default (referensi)**.
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
maksimal 1 template ACTIVE  +  1 template DEFAULT (referensi)
```

---

## 2.3 Aturan Template Aktif & Template Default

### Satu Template Aktif per Jenis Dokumen

- Setiap jenis dokumen hanya boleh memiliki **1 template dengan status `Active`**.
- Mengaktifkan sebuah template otomatis **menonaktifkan template lain** dengan jenis dokumen yang sama (dalam satu transaksi).
- Template `Inactive` **tidak dapat digunakan** untuk Generate Document.
- Aturan dipaksa di dua lapisan:
  - **Database** — partial unique index `(document_type) WHERE status = 'ACTIVE'`.
  - **Service** — validasi + transaksi saat aktivasi.

### Template Default (Referensi)

- Setiap jenis dokumen memiliki **1 template default** yang berperan sebagai **referensi** (starter template).
- Template default **tidak dapat digunakan langsung** untuk Generate Document dan **tidak dapat diaktifkan**.
- Jika belum memiliki template untuk suatu jenis dokumen, administrator menggunakan template default melalui alur:

```text
Template Default (referensi)
      ↓
[ Gunakan Template Default ]
      ↓
Salin konten ke template baru (draft)
      ↓
Simpan data template        ← WAJIB disimpan terlebih dahulu
      ↓
Template baru siap diedit & diaktifkan
      ↓
Aktifkan → menjadi 1-satunya template aktif untuk jenis tsb
```

- Data template **harus disimpan terlebih dahulu** sebelum bisa digunakan; template default hanya menyediakan konten awal, bukan template yang langsung dipakai.
- Konten template default **dapat diubah** oleh administrator (referensi diperbarui). Perubahan ini **tidak memengaruhi** template yang sudah dibuat dari default, karena template tersebut sudah menjadi salinan tersendiri.
- Template default bersifat **seeded** (satu per jenis dokumen), tidak dapat dihapus, dan dapat di-*reset* ke konten bawaan.

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
| Status | Active / Inactive / Reference (default) |
| Default | Penanda template default (referensi) |
| Updated At | Waktu terakhir update |
| Action | Detail/Edit/Version/Preview/Gunakan Default |

Catatan:

- Hanya **1 template `Active`** per Document Type. Mengaktifkan template lain otomatis menonaktifkan template yang sedang aktif (dengan konfirmasi).
- Baris bertanda **Default** adalah template referensi — statusnya `Reference`, tidak dapat diaktifkan/digunakan langsung.
- Jika suatu Document Type **belum memiliki template sama sekali**, tampilkan aksi `[ Gunakan Template Default ]` untuk jenis tersebut.

Contoh:

```text
Template Dokumen

[ + Template Baru ]                       [ Search ]

┌─────────────────────────────────────────────────────────────────────┐
│ Name                 │ Type        │ Version │ Status    │ Default  │
├─────────────────────────────────────────────────────────────────────┤
│ Perjanjian PKWT      │ Contract    │ v3       │ Active    │          │
│ Perjanjian PKWT (D)  │ Contract    │ -        │ Reference │ ★        │
│ SK Movement          │ Movement SK │ v2       │ Active    │          │
│ SK Movement (D)      │ Movement SK │ -        │ Reference │ ★        │
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

[ Ambil dari Template Default ]   ← opsional; mengisi Template Content
                                    dengan konten default jenis tsb
                                    (alur: salin → simpan → aktifkan)

Status
[ Inactive ]   ← hanya 1 template Active per jenis dokumen;
                  mengaktifkan otomatis menonaktifkan yang lain.
                  Template default (referensi) tidak dapat diaktifkan.

Template Content
┌──────────────────────────────────────────────┐
│              PrimeVue Editor                 │
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

# 6. PrimeVue Editor

Gunakan:

```vue
<Editor />
```

Editor digunakan untuk membuat konten HTML dokumen.

Fitur toolbar minimal:

- Bold
- Italic
- Underline
- Font size
- Heading
- Text alignment
- Ordered list
- Bullet list
- Link
- Table
- Image
- Text color
- Background
- Undo/Redo

Editor harus mendukung:

- Paragraph
- Table
- Image/logo
- Alignment
- Signature area
- Page break

---

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

Preview menghasilkan **PDF sungguhan** dari konten template sehingga administrator dapat melihat hasil akhir persis seperti dokumen yang akan digenerate — langsung dari Settings, tanpa harus masuk ke module Contract/Movement.

### 8.1 Tombol Preview

Tersedia di dua tempat:

```text
1. Template List → aksi [ Preview ] per baris
   → preview menggunakan VERSI AKTIF template yang tersimpan.

2. Template Editor (form) → tombol [ Preview ] di toolbar
   → preview konten DRAFT yang sedang diedit (belum/tidak perlu disimpan).
```

### 8.2 Dua Mode Preview

```text
Preview Versi Tersimpan               Preview Draft
─────────────────────────             ──────────────────────────
Pakai active_version_id               Kirim konten editor saat ini
+ konfigurasi dokumen tersimpan       + konfigurasi dokumen (paper,
      ↓                               orientation, margin)
Resolve dengan data contoh                  ↓
      ↓                                Resolve dengan data contoh
Render HTML → PDF                            ↓
                                       Render HTML → PDF
```

- **Preview Versi Tersimpan** — memakai konten versi aktif + konfigurasi dokumen yang sudah disimpan.
- **Preview Draft** — mengirimkan konten editor dan konfigurasi saat itu juga, agar hasilnya akurat meskipun template **belum disimpan**.

### 8.3 Data Contoh (Sample Data)

Preview tidak memakai data asli, melainkan **data contoh** agar hasil variabel terlihat:

```text
{{employee.name}}       → Asep Ruswanda
{{employee.nik}}        → 199001012015011001
{{employee.position}}   → HR Staff
{{employee.organization}}→ HR Division
{{contract.number}}     → CTR-2026-001
{{contract.start_date}} → 2026-01-01
{{contract.end_date}}   → 2027-01-01
{{company.name}}        → PT Maju Bersama
```

### 8.4 Alur Preview ke PDF

```text
Template Editor / Template List
      ↓
[ Preview ]
      ↓
POST preview (versi tersimpan ATAU konten draft)
      ↓
Resolve Variables (data contoh)
      ↓
Render HTML (engine sama dengan PDF final)
      ↓
HTML → PDF (Headless Chromium)
      ↓
Tampilkan di Dialog PDF Viewer
      ↓
[ Download PDF ] [ Print ] [ Tutup ]
```

- Selama proses berjalan tampilkan `ProgressSpinner`/loading (generasi PDF membutuhkan beberapa detik).
- **WYSIWYG**: preview memakai engine rendering yang sama dengan Generate Document (Headless Chromium) sehingga hasil preview = hasil PDF final.

### 8.5 PDF Viewer

- Preview ditampilkan dalam `Dialog`/`Drawer` berukuran besar menggunakan `<iframe>`/PDF embed.
- Tombol aksi di viewer: **Download PDF** dan **Print**.
- Jika variabel dari `document_type` tidak memiliki sample data, tampilkan placeholder `-` dengan catatan kecil bahwa variabel kosong akan diisi data asli saat generate.

---

# 9. Document Configuration

Template perlu memiliki konfigurasi dokumen.

### Paper

```text
A4
A5
Letter
Legal
```

### Orientation

```text
Portrait
Landscape
```

### Margin

```text
Top
Right
Bottom
Left
```

### Header

- Logo
- Company name
- Address
- Custom header

### Footer

- Page number
- Document number
- Custom footer

---

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

Catatan: status **Active/Inactive** berada di level **template** (maksimal 1 Active per jenis dokumen), sedangkan **versi aktif** berada di level template (`active_version_id`). Template default (referensi) tidak memakai mekanisme versi aktif — kontennya diperbarui langsung sebagai referensi.

---

# 11. Database Design

> **Catatan migrasi (sesuai konvensi project):**
> - Tabel `document_templates` **sudah ada** (dibuat di `migrator/migrations/tenant/{postgres,mysql}/011_settings.sql`) dengan skema `id, name, type, content, is_active` → migration baru bersifat **ALTER/ekstensi**, bukan CREATE baru.
> - Module memakai **GORM AutoMigrate** (pola module `setting`): kolom baru otomatis ditambahkan ke tabel existing. Partial unique index & seed tidak bisa AutoMigrate → dibuat via file SQL migrator baru `012_document_templates.sql`.

## `document_templates` (ALTER dari tabel existing 011_settings)

```text
id UUID PK                       ← existing
name                             ← existing
code (baru)
type                             ← existing, dipakai sebagai jenis dokumen
                                   (atau tambah kolom baru document_type, mapping 1:1)
description (baru)
content                          ← existing (NULLABLE; konten baru disimpan di versions;
                                     data lama dimigrasi ke version v1)
active_version_id (baru)
status              -- ACTIVE / INACTIVE / REFERENCE (default) (baru)
is_default          -- BOOLEAN, penanda template default (referensi) (baru)
created_at / updated_at          ← existing
deleted_at (baru, soft delete)
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
-- Hanya 1 template aktif per jenis dokumen (PostgreSQL)
CREATE UNIQUE INDEX uq_document_templates_active_per_type
  ON document_templates (document_type)
  WHERE status = 'ACTIVE';

-- Hanya 1 template default (referensi) per jenis dokumen (PostgreSQL)
CREATE UNIQUE INDEX uq_document_templates_default_per_type
  ON document_templates (document_type)
  WHERE is_default = TRUE;
```

Catatan portabilitas:

- Partial unique index (`WHERE`) **hanya didukung PostgreSQL** dan tidak bisa dibuat GORM AutoMigrate → dibuat via file SQL migrator baru `012_document_templates.sql` (Postgres).
- **MySQL**: tanpa partial index → validasi "1 aktif per jenis" di **service layer** (query template aktif per jenis sebelum aktivasi, dalam transaksi).

Seed data: satu template default per jenis dokumen (`CONTRACT_AGREEMENT`, `MOVEMENT_SK`, dst.) berisi konten contoh lengkap dengan placeholder.

Template default:

- `is_default = TRUE`, `status = REFERENCE` (tidak pernah dipakai Generate Document).
- Tidak dapat dihapus/diaktifkan; konten dapat diubah (referensi diperbarui).
- Tidak dihitung dalam aturan "1 template aktif per jenis".

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

Catatan: `generated_documents` adalah tabel **baru** (belum ada) — dibuat via GORM AutoMigrate. File PDF disimpan mengikuti pola upload existing (`internal/pkg/upload`, path `/uploads/attachments/...` seperti `file_url` pada module employeemovement).

---

# 13. Backend Go

Ikuti konvensi module project: **package flat** di `backend/internal/modules/` (bukan Clean Architecture subfolder) — pola yang sama dengan `setting`, `employeemovement`, dll.

```text
backend/internal/modules/documenttemplate/
├── module.go         ← implementasi module.Module (Info, RegisterRoutes, Migrate, Seed, Permissions)
├── model.go          ← GORM model: DocumentTemplate, DocumentTemplateVersion, GeneratedDocument, DocumentTemplateAudit
├── repository.go     ← GORM + tenant resolver
├── service.go        ← bisnis: 1-aktif-per-jenis, default (copy), versioning
├── render_service.go ← resolusi variabel + render HTML
├── pdf_service.go    ← HTML → PDF (chromedp)
├── handler.go        ← Gin handlers (httputil)
├── routes.go         ← registrasi route
├── dto.go            ← request/response + tag validasi
└── errors.go         ← error domain (mis. DuplicateActiveTemplateError)
```

### Registrasi Module

Daftarkan di `cmd/server/main.go` (tenantModules), setelah module setting:

```go
module.ModuleRegistration{
    Module:   documenttemplate.NewModule(dbManager, l),
    TargetDB: module.TargetTenant,
    Priority: 16,
},
```

### Tenancy & Migrasi

- **Tenant-scoped**: ikuti pola `NewTenantDBResolver(dbManager)` — baca `company_id` dari context; semua repo memakai tenant DB.
- **Migrate**: `Module.Migrate(db)` memanggil `db.AutoMigrate(&DocumentTemplate{}, &DocumentTemplateVersion{}, &GeneratedDocument{}, &DocumentTemplateAudit{})` (pola module setting). AutoMigrate otomatis menambah kolom baru ke tabel `document_templates` existing.
- **SQL migrator**: file baru `012_document_templates.sql` di `migrator/migrations/tenant/{postgres,mysql}/` untuk partial unique index (Postgres) & seed template default.

---

# 14. API

Base path mengikuti konvensi project: **`/api/v1/tenant/document-templates`** (route didaftarkan via `rg.Group("/document-templates")` di `routes.go`, di dalam tenant router `/api/v1/tenant`).

Konvensi response & error (helper `httputil`):

```text
Sukses list:   { "success": true, "data": [...], "total": N, "page": P }
Sukses detail: { "success": true, "data": {...} }
Error:         { "success": false, "error": { "code": "...", "message": "..." } }
Pesan error bilingual EN/ID via locale key di internal/pkg/httputil/locale.go
Validasi body: httputil.BindAndValidate + tag validasi di dto.go
```

Catatan route (Gin): daftarkan route statis (`from-default`, `preview-draft`, `variables`) **sebelum** route `/:id` agar tidak tertangkap parameter — pola yang sama dengan `GET /search` pada villages.

### Template

```http
GET    /api/v1/tenant/document-templates
POST   /api/v1/tenant/document-templates
GET    /api/v1/tenant/document-templates/{id}
PUT    /api/v1/tenant/document-templates/{id}
DELETE /api/v1/tenant/document-templates/{id}

POST   /api/v1/tenant/document-templates/{id}/activate   ← otomatis menonaktifkan
                                                          template lain sejenis
POST   /api/v1/tenant/document-templates/{id}/deactivate
POST   /api/v1/tenant/document-templates/from-default   ← body: { "document_type": "..." }
                                                          salin konten default →
                                                          template baru (draft, wajib disimpan)
```

### Version

```http
GET  /api/v1/tenant/document-templates/{id}/versions
POST /api/v1/tenant/document-templates/{id}/versions
GET  /api/v1/tenant/document-templates/{id}/versions/{versionId}
```

### Preview

```http
POST /api/v1/tenant/document-templates/{id}/preview   ← preview VERSI AKTIF (data contoh)

POST /api/v1/tenant/document-templates/preview-draft  ← preview DRAFT:
                                                       body: { document_type,
                                                               content,
                                                               paper_size,
                                                               orientation,
                                                               margins }

Response: { "pdf_url": "...", "file_name": "preview_<template>.pdf" }
```

### Variables

```http
GET /api/v1/tenant/document-templates/variables
```

---

# 15. PDF Generation

Backend Go menangani:

```text
Template
   ↓
Load Template Version
   ↓
Load Reference Data
   ↓
Resolve Variables
   ↓
Generate HTML
   ↓
HTML → PDF
   ↓
Store PDF
   ↓
Return File
```

### PDF Engine (keputusan dependency — BARU)

Project saat ini **belum memiliki** library PDF/headless browser. Rekomendasi: **`chromedp`** (pure Go, Chrome DevTools Protocol) — memakai `Page.printToPDF` dari Headless Chromium sehingga hasil render HTML/CSS identik dengan preview (WYSIWYG).

```text
Dependency baru:
  github.com/chromedp/chromedp
  + binary Chromium/Chrome di environment server
  (Docker image: tambahkan chromium; konfigurasi env mis. CHROME_PATH)

Alternatif:
  - go-rod  (high-level CDP driver, API lebih sederhana)
  - gotenberg (service container terpisah — opsi jika tidak ingin binary Chrome di app server)
```

Catatan infrastruktur:

- Preview dan Generate Document memakai **engine yang sama** (chromedp + Chrome) agar hasil konsisten.
- Generate PDF dilakukan di backend; frontend hanya menampilkan hasil (`pdf_url`).

---

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
Buat template dari default di Settings → Template Dokumen (alur Gunakan Template Default).
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

Jika **belum ada template aktif** untuk `MOVEMENT_SK`, tampilkan pesan yang sama seperti pada Contract (arahkan ke alur *Gunakan Template Default*).

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

Tambahkan permission (konvensi `module.resource.action` — lowercase, underscore, prefix module; didaftarkan di `Info().Permissions`):

```text
documenttemplate.template.view
documenttemplate.template.create
documenttemplate.template.update
documenttemplate.template.delete
documenttemplate.template.preview
documenttemplate.template.version
documenttemplate.template.activate
documenttemplate.template.deactivate
documenttemplate.template.set_default
```

Untuk generated document:

```text
documenttemplate.generated.view
documenttemplate.generated.generate
documenttemplate.generated.download
```

Gunakan mekanisme permission/RBAC yang sudah digunakan aplikasi (RBAC middleware + daftar permission di `module.go`).

---

# 19. Audit Log

Project **tidak memiliki sistem audit terpusat** — pola existing adalah tabel append-only per entity (mis. `OrganizationHistory` di module organization, `CandidateConsent` di recruitment). Ikuti pola tersebut: buat tabel **`document_template_audits`** (append-only, via AutoMigrate):

```text
id UUID PK
template_id CHAR(36) FK
version_id CHAR(36) NULL
action VARCHAR(50)     -- CREATED / UPDATED / VERSION_CREATED / ACTIVATED /
                        -- DEACTIVATED / DEFAULT_UPDATED / CREATED_FROM_DEFAULT /
                        -- GENERATED / DOWNLOADED
actor_id CHAR(36)
payload JSONB NULL     -- snapshot ringkas (nama template, nomor versi, dll)
created_at TIMESTAMP
```

Catat aktivitas penting (nilai kolom `action`):

```text
Template Created
Template Updated
Template Version Created
Template Activated
Template Deactivated
Template Deleted
Template Default Updated
Template Created from Default
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

Ikuti konvensi views settings: satu view utama per resource di `frontend/tenant/src/views/settings/` (pola `RelationshipTypesView.vue` — DataTable lazy + Dialog + SkeletonTable + ConfirmDeleteDialog + FormRow/TextInput + useI18n + axios `api`). Dialog dipisah bila terlalu besar:

```text
frontend/tenant/src/views/settings/
├── DocumentTemplatesView.vue          ← halaman utama (list, status, default, aksi)
├── DocumentTemplateFormDialog.vue     ← form + PrimeVue Editor + variable picker
├── DocumentTemplateVersionDialog.vue  ← histori/versi
└── DocumentTemplatePreviewDialog.vue  ← PDF viewer (iframe + Download/Print)
```

Registrasi & i18n (wajib sinkron):

- **Router**: tambah route `settings/document-templates` di `frontend/tenant/src/router/index.js` dengan `meta { titleKey, descKey, icon, module: 'setting' }` (pola route settings lain).
- **SettingsIndex.vue**: tambah card sub-menu "Template Dokumen" (grup baru, mis. `dokumen`).
- **Locale**: tambah key `document_templates.*` (title, description, kolom, aksi, pesan) di `frontend/tenant/src/locales/en.json` dan `id.json`.

PrimeVue component:

```text
DataTable
Dialog
Drawer
Editor
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

## Phase 1 — Database & Backend Foundation

- [ ] Module `documenttemplate` + registrasi di `cmd/server/main.go` (Priority 16)
- [ ] GORM model + AutoMigrate (`document_templates` ALTER, `document_template_versions`, `generated_documents`, `document_template_audits`)
- [ ] SQL migrator `012_document_templates.sql`: partial unique index (Postgres) + seed template default
- [ ] Tenant resolver + repository
- [ ] Service (1-aktif-per-jenis, default copy, versioning)
- [ ] DTO + validasi (BindAndValidate)
- [ ] API handler + routes (`/api/v1/tenant/document-templates`)
- [ ] Permission `documenttemplate.*` + RBAC
- [ ] httputil locale keys (EN/ID)

## Phase 2 — Template Management

- [ ] Template list
- [ ] Search
- [ ] Filter document type
- [ ] Create template
- [ ] Edit template
- [ ] Detail template
- [ ] Active/inactive — hanya 1 aktif per jenis (aktivasi otomatis menonaktifkan yang lain)
- [ ] Set/ubah template default (referensi)
- [ ] Alur "Gunakan Template Default" (salin → simpan → aktifkan)
- [ ] Delete template
- [ ] Version management

## Phase 3 — Template Editor

- [ ] Integrasi PrimeVue Editor
- [ ] Toolbar customization
- [ ] Variable picker
- [ ] Insert variable
- [ ] Table support
- [ ] Image/logo
- [ ] Page break
- [ ] Document configuration

## Phase 4 — Preview & PDF

- [ ] Template rendering
- [ ] Variable resolver
- [ ] Sample data preview
- [ ] HTML rendering
- [ ] Dependency chromedp + Chrome binary (env `CHROME_PATH`)
- [ ] PDF generation (Headless Chromium via chromedp)
- [ ] API preview versi tersimpan (`{id}/preview`)
- [ ] API preview draft (`/preview-draft`)
- [ ] PDF viewer dialog (iframe) di Settings
- [ ] Download PDF dari preview
- [ ] Print PDF dari preview
- [ ] Loading state selama generasi PDF
- [ ] Generated document storage

## Phase 5 — Module Integration

### Contract

- [ ] Generate Contract Agreement
- [ ] Template selection
- [ ] Contract data mapping
- [ ] Employee data mapping
- [ ] Generated document history

### Movement

- [ ] Generate SK Movement
- [ ] Template selection
- [ ] Movement data mapping
- [ ] Employee data mapping
- [ ] Generated document history

## Phase 6 — Testing

### Backend

- [ ] Template CRUD test
- [ ] Versioning test
- [ ] Variable resolver test
- [ ] Permission test
- [ ] PDF generation test
- [ ] Generated document test

### Frontend

- [ ] Template list test
- [ ] Create template test
- [ ] Editor test
- [ ] Variable insertion test
- [ ] Preview test
- [ ] Versioning test
- [ ] PDF download test

---

# 22. Acceptance Criteria

Fitur dianggap selesai apabila:

- [ ] Administrator dapat membuat template dokumen dari Settings.
- [ ] Administrator dapat menggunakan PrimeVue Editor.
- [ ] Administrator dapat memasukkan variable ke template.
- [ ] Variable dapat di-resolve berdasarkan data employee/contract/movement.
- [ ] Template memiliki versioning.
- [ ] Template aktif dapat digunakan oleh module Contract.
- [ ] Template aktif dapat digunakan oleh module Employee Movement.
- [ ] Contract dapat menghasilkan Perjanjian PDF.
- [ ] Movement dapat menghasilkan SK PDF.
- [ ] PDF dapat di-preview langsung dari Settings (versi tersimpan maupun draft).
- [ ] Preview menghasilkan PDF sungguhan dengan engine yang sama dengan Generate Document (WYSIWYG).
- [ ] PDF dapat di-download.
- [ ] PDF dapat di-print.
- [ ] Dokumen yang sudah dihasilkan menyimpan referensi template dan versinya.
- [ ] Perubahan template tidak mengubah dokumen yang sudah diterbitkan.
- [ ] Hanya ada 1 template aktif per jenis dokumen; mengaktifkan template lain otomatis menonaktifkan template sebelumnya.
- [ ] Setiap jenis dokumen memiliki 1 template default (referensi) yang dapat diubah kontennya.
- [ ] Template default tidak dapat digunakan/diaktifkan langsung; harus disalin & disimpan sebagai template baru terlebih dahulu.
- [ ] Alur "Gunakan Template Default" menghasilkan template baru yang tersimpan sebelum bisa diaktifkan.
- [ ] Permission dan audit log diterapkan.

---

# 23. Target Arsitektur Akhir

```text
                         SETTINGS
                            │
                     TEMPLATE DOKUMEN
                            │
              ┌─────────────┴─────────────┐
              │                           │
        Contract Template          Movement Template
              │                           │
              └─────────────┬─────────────┘
                            │
                    Template Version
                            │
                    PrimeVue Editor
                            │
                    {{variables}}
                            │
                       Go Backend
                            │
                  Variable Resolution
                            │
                     HTML Rendering
                            │
                    Headless Chromium
                            │
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
