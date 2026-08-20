# Alur Pengisian Organization Management (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Organization Management** — setup
struktur organisasi (tree hierarchy), Organization Summary (SK pendirian), versioning/snapshot,
history/audit trail, clone — plus master data terkait (Zones, Job Families, Gradings di modul
Settings; Job Management di modul terpisah) — pola runbook seperti
[`module-leave-flow.md`](module-leave-flow.md) & [`module-reimbursement-flow.md`](module-reimbursement-flow.md).

- Lokasi kode: `backend/internal/modules/organization/` · `frontend/tenant/src/views/modules/organization/`
- Modul terkait: Settings (Zones, Job Families, Gradings) · Job Management (Job Titles, Values, dll.)

---

## 1. Ringkasan Alur End-to-End

```
SETUP MASTER (sekali)              STRUKTUR ORGANISASI                 VERSIONING & AUDIT
┌────────────────────┐   ┌──────────────────────────────────┐   ┌──────────────────────────┐
│ Zones              │   │ Buat/edit node di tree            │   │ Create Version (snapshot) │
│ Job Families       │──▶│   Code, Nomenclature, Parent      │──▶│   Diff versions           │
│ Gradings           │   │   Zone, Job Family, Grading       │   │   Restore version         │
└────────────────────┘   │ Full Code otomatis (propagasi)    │   │   Clone tree              │
                         └──────────────────────────────────┘   │ History (audit trail)     │
                                                                └──────────────────────────┘
                                                                         │
                                                                         ▼
                                                                ┌──────────────────────────┐
                                                                │ Organization Summary      │
                                                                │   (SK pendirian)          │
                                                                │   Code, Decree, Status    │
                                                                │   active / inactive       │
                                                                └──────────────────────────┘
```

- **Organization** = node dalam tree hierarki (parent → child). Setiap node punya `code`,
  `full_code` (path dari root), `nomenclature` (nama), dan referensi ke Zone/Job Family/Grading.
- **Full Code Propagation** — saat kode node berubah, semua `full_code` descendant otomatis di-update.
- **Organization Summary** = SK pendirian yang mengelompokkan beberapa node organisasi
  (satu summary bisa punya banyak organization).
- **Versioning** = snapshot lengkap pohon organisasi pada waktu tertentu — bisa di-diff & di-restore.
- **History** = audit trail setiap perubahan (CREATE/UPDATE/DELETE) dengan old/new values.

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Organization | `organizations` | Node dalam tree organisasi (code, full_code, nomenclature, parent, zone, job family, grading, level, sort_order) |
| OrganizationHistory | `organization_histories` | Audit trail perubahan (action: CREATE/UPDATE/DELETE, old/new values JSON) |
| OrganizationVersion | `organization_versions` | Snapshot pohon organisasi (version_name, description, snapshot JSON, status, node_count) |
| Organization Summary | `organization_summaries` | SK pendirian (code, decree_no, decree_date, status: active/inactive) — relasi hasMany ke Organizations |
| **Master terkait (modul Settings):** | | |
| Zone | `zones` | Zona geografis (referensi di Organization) |
| Job Family | `job_families` | Keluarga pekerjaan (referensi di Organization) |
| Grading | `gradings` | Grading/level jabatan (referensi di Organization) |

---

## 3. TAHAP 1 — SETUP Master Data (dikerjakan sekali)

Master data berikut dikelola di modul **Settings** (`/admin/settings/`) — bukan di modul
Organization — tapi menjadi referensi wajib sebelum membuat node organisasi.

### A. Zones

Menu **Settings → Zones** (`/admin/settings/zones`, `ZonesView.vue`).

- CRUD zona geografis: nama, deskripsi.
- Endpoint: `GET/POST /settings/zones`, `GET/PUT/DELETE /settings/zones/:id`.

### B. Job Families

Menu **Settings → Job Families** (`/admin/settings/job-families`, `JobFamiliesView.vue`).

- CRUD keluarga pekerjaan: kode, nama, deskripsi.
- Endpoint: `GET/POST /settings/job-families`, `GET/PUT/DELETE /settings/job-families/:id`.

### C. Gradings

Menu **Settings → Gradings** (`/admin/settings/gradings`, `GradingsView.vue`).

- CRUD grading/level jabatan: kode, nama, level number, deskripsi.
- Endpoint: `GET/POST /settings/gradings`, `GET/PUT/DELETE /settings/gradings/:id`.

> ⚠️ **Pastikan master di atas sudah diisi** sebelum membuat node organisasi — Zone, Job Family,
> dan Grading dipilih saat membuat/mengedit organization node.

---

## 4. TAHAP 2 — STRUKTUR ORGANISASI (Tree)

Menu **Organization → Organization Tree** (`/admin/organizations`, `Organizations.vue`).

### A. Membuat Node Organisasi

1. Klik **"Add Root"** untuk membuat node tingkat atas (tanpa parent), atau klik **"+"** pada
   node yang sudah ada untuk menambah child.
2. Isi form:
   - **Code** — kode singkat node (unik per parent).
   - **Nomenclature** — nama/judul organisasi.
   - **Zone** — pilih zona dari master (opsional).
   - **Job Family** — pilih keluarga pekerjaan (opsional).
   - **Grading** — pilih grading (opsional).
3. **Full Code** dihitung otomatis oleh backend: gabungan code dari root ke node ini
   (mis. `01.02.03` untuk root→child→grandchild).
4. **Level** dihitung otomatis berdasarkan kedalaman di tree.
5. **Sort Order** bisa diatur untuk pengurutan sibling.

### B. Mengedit Node

- Klik ikon edit pada node → ubah field yang diperlukan.
- **Perubahan Code** memicu **Full Code Propagation** — semua `full_code` descendant
  otomatis di-update (backend memproses secara rekursif).

### C. Menghapus Node

- Klik ikon hapus → konfirmasi.
- Node dengan child **tidak bisa dihapus** sampai semua child dihapus terlebih dahulu.

### D. Tampilan

- **Table View** — tampilan tree table (expand/collapse).
- **Chart View** — tampilan org chart visual (`OrgChartView.vue`).
- **Search** — filter node berdasarkan nama/kode.

---

## 5. TAHAP 3 — ORGANIZATION SUMMARY (SK Pendirian)

Menu **Organization → Organization Summary** (`/admin/organizations/summary`, `OrganizationSummary.vue`).

### A. Buat Summary

1. Isi: **Code** (7 karakter), **Decree No** (nomor SK), **Decree Date** (tanggal SK).
2. Status default: `inactive`.
3. **Clone From** (opsional) — pilih summary lain untuk menyalin semua node organisasi
   ke summary baru (termasuk seluruh tree hierarchy).

### B. Kelola Summary

- **Status**: `active` / `inactive` — tentukan SK pendirian mana yang berlaku.
- **Org Count** — jumlah node organisasi yang terkait (ditampilkan otomatis).
- **Stats** — endpoint `GET /organization-summaries/stats` menampilkan ringkasan.

### C. Hubungan Summary ↔ Organization

- Satu **Organization Summary** bisa memiliki banyak **Organization** nodes
  (relasi hasMany via `organization_summary_id`).
- Saat kode organization berubah, propagasi `full_code` juga memperbarui
 所有 descendant yang terkait summary yang sama.

---

## 6. TAHAP 4 — VERSIONING & AUDIT TRAIL

### A. Create Version (Snapshot)

Menu **Organization → Versions** (di dalam Organizations.vue).

- Klik **"Create Version"** → isi version name + description.
- Sistem mengambil snapshot JSON dari seluruh pohon organisasi saat ini.
- **Status version**: `DRAFT → ACTIVE → ARCHIVED`
- **Node Count** dihitung otomatis.
- **Parent Version ID** — versi sebelumnya (opsional, untuk tracking lineage).

### B. Diff Versions

- Pilih dua version → lihat perbedaan (node yang ditambah/dihapus/diubah).
- Endpoint: `GET /versions/:id/diff/:targetId`.

### C. Restore Version

- Pilih version → klik **"Restore"** → pohon organisasi dikembalikan ke kondisi snapshot.
- Perubahan tercatat di history (CREATE/DELETE/UPDATE untuk setiap node yang berubah).
- Endpoint: `POST /versions/:id/restore`.

### D. History (Audit Trail)

- Endpoint: `GET /organizations/history`.
- Menampilkan setiap perubahan: action (CREATE/UPDATE/DELETE), old values, new values,
  changed by (user), timestamp.
- Filterable by organization_id.

### E. Clone Tree

- Clone seluruh pohon organisasi ke summary baru.
- Endpoint: `POST /organizations/clone`.

---

## 7. Ringkasan Status & Transisi

| Entitas | Status | Transisi |
|---|---|---|
| **Organization Version** | `DRAFT → ACTIVE → ARCHIVED` | manual via `PUT /versions/:id` |
| **Organization Summary** | `active` / `inactive` | manual via `PUT /organization-summaries/:id` |
| **Organization History** | `CREATE` / `UPDATE` / `DELETE` | otomatis saat CRUD organization |

---

## 8. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Settings (Zones)** | Referensi zona geografis untuk organization node |
| **Settings (Job Families)** | Referensi keluarga pekerjaan untuk organization node |
| **Settings (Gradings)** | Referensi grading/level jabatan untuk organization node |
| **Employee** | Employee ditugaskan ke organization node (via `organization_id` di employee) |
| **Job Management** | Job Titles, Values, dll. menggunakan organization sebagai konteks |
| **Recruitment** | Requisition ditautkan ke organization (posisi dalam struktur) |
| **Attendance** | Exempt positions & shift bisa difilter per organization |
| **Leave** | Balance & policy bisa difilter per organization |
| **Training** | Session & participant bisa difilter per organization |
| **Payroll** | Struktur organisasi menjadi dasar penggajian |
| **Workforce Intelligence** | Analisis workforce per organization node |

---

## 9. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Organization (hub) | — | Bagian dari Settings/Admin |
| Organization → Tree | `Organizations.vue` | CRUD node tree + table/chart view + search + versioning + history |
| Organization → Summary | `OrganizationSummary.vue` | CRUD SK pendirian + status active/inactive + clone + stats |
| Organization → Org Chart | `OrgChartView.vue` | Visualisasi org chart |
| Organization → Tree Table | `OrgTreeTable.vue` | Tampilan tree table (alternatif) |
| Settings → Zones | `ZonesView.vue` | CRUD zona geografis |
| Settings → Job Families | `JobFamiliesView.vue` | CRUD keluarga pekerjaan |
| Settings → Gradings | `GradingsView.vue` | CRUD grading/level jabatan |

---

## 10. Endpoint API Utama

### Organization (`/api/v1/tenant/organizations/`)

| Area | Endpoint |
|---|---|
| CRUD | `POST /organizations`, `GET /organizations`, `GET /organizations/:id`, `PUT /organizations/:id`, `DELETE /organizations/:id` |
| History | `GET /organizations/history` |
| Versions | `POST /organizations/versions`, `GET /organizations/versions`, `GET /organizations/versions/:id`, `GET /organizations/versions/:id/diff/:targetId`, `POST /organizations/versions/:id/restore` |
| Clone | `POST /organizations/clone` |

### Organization Summary (`/api/v1/tenant/organization-summaries/`)

| Area | Endpoint |
|---|---|
| CRUD | `POST /organization-summaries`, `GET /organization-summaries`, `GET /organization-summaries/stats`, `GET /organization-summaries/:id`, `PUT /organization-summaries/:id`, `DELETE /organization-summaries/:id` |

### Master Terkait (`/api/v1/tenant/settings/`)

| Area | Endpoint |
|---|---|
| Zones | `GET/POST /settings/zones`, `GET/PUT/DELETE /settings/zones/:id` |
| Job Families | `GET/POST /settings/job-families`, `GET/PUT/DELETE /settings/job-families/:id` |
| Gradings | `GET/POST /settings/gradings`, `GET/PUT/DELETE /settings/gradings/:id` |

---

## 11. Catatan Penting

- **Full Code Propagation** — saat kode node berubah, semua `full_code` descendant otomatis
  di-update. Proses ini dilakukan backend secara rekursif (bukan batch job).
- **Organization Summary (SK Pendirian)** adalah kontainer untuk struktur organisasi —
  satu summary bisa punya banyak node, dan status `active`/`inactive` menentukan SK yang berlaku.
- **Clone** pada summary menyalin seluruh tree structure ke summary baru — useful untuk
  membuat draft reorganisasi tanpa mengubah struktur aktif.
- **Versioning** menyimpan snapshot JSON pohon organisasi — bisa di-diff (lihat perbedaan)
  dan di-restore (kembalikan ke kondisi tertentu).
- **History** mencatat setiap perubahan dengan old/new values JSON — useful untuk audit trail.
- **Isi master data (Zones, Job Families, Gradings) terlebih dahulu** sebelum membuat
  node organisasi — field ini dipilih saat create/edit organization.
- **Job Management** (Job Titles, Values, Objectives, dll.) dikelola di modul terpisah
  (`/admin/job-management/`) — bukan bagian dari modul Organization, tapi berkaitan erat.
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
