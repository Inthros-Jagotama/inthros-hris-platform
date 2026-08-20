# Alur Settings & Document Template (Runbook)

Dokumen ini menjelaskan **cara pakai** modul **Settings** — master data referensi
(geografi, pendidikan, agama, status kerja, bank, asuransi, kompetensi, grading, salary grade,
hari libur, TER, PTKP) + **Document Template** (template dokumen PDF, versioning, generate)
+ **Document Numbering** (penomoran dokumen otomatis) + **Employee ID Format**
+ **Timezone Settings** (zona waktu default company & override per Zone).

- Lokasi kode: `backend/internal/modules/setting/` · `backend/internal/modules/documenttemplate/` · `backend/internal/pkg/timezone/`
- Module slug: `setting`
- Semua master data di Settings bersifat CRUD standar (tanpa workflow khusus).

---

## 1. Ringkasan Alur

```
SETTINGS (master data)           DOCUMENT TEMPLATES              NUMBERING
┌────────────────────────┐   ┌──────────────────────────┐   ┌──────────────────────┐
│ Zones, Provinces,       │   │ CRUD template             │   │ Konfigurasi format   │
│ Regencies, Districts,   │   │ → Version management      │   │ penomoran dokumen    │
│ Villages, Education,    │   │ → Activate/deactivate     │   │ → Preview numbering  │
│ Majors, Religion,       │   │ → Preview                 │   │ → Per document type  │
│ Marital Status,         │   │ → Generate PDF            │   └──────────────────────┘
│ Relationship Type,      │   └──────────────────────────┘
│ Employment Status,      │   EMPLOYEE ID FORMAT
│ Bank, Nationality,      │   ┌──────────────────────────┐
│ Job Family, Grading,    │   │ Konfigurasi format ID    │
│ Salary Grade, TER,      │   │ → auto/hybrid/manual     │
│ PTKP, Insurance,        │   │ → Preview                │
│ Company Holiday,        │   └──────────────────────────┘
│ Competency              │   TIMEZONE SETTINGS
└────────────────────────┘   ┌──────────────────────────┐
                              │ Company default timezone │
                              │ → Zone override (opsional)│
                              │ → Resolve: Zone>Company>WIB│
                              └──────────────────────────┘
```

---

## 2. Entitas Utama

### Settings Master Data

| Entitas | Deskripsi |
|---|---|
| Zone | Zona geografis (payroll/allowance zoning) + `timezone` (opsional): override IANA timezone (`Asia/Jakarta`/`Asia/Makassar`/`Asia/Jayapura`) untuk organization yang tergabung di zona ini |
| Province | Provinsi |
| Regency | Kabupaten/Kota |
| District | Kecamatan |
| Village | Desa/Kelurahan |
| Education | Jenjang pendidikan |
| EducationMajor | Jurusan pendidikan |
| Religion | Agama |
| MaritalStatus | Status pernikahan |
| RelationshipType | Hubungan keluarga |
| EmploymentStatus | Status kerja (active, inactive, dll) |
| Bank | Bank |
| Nationality | Kewarganegaraan |
| JobFamily | Keluarga jabatan |
| Grading | Grade/jenjang |
| SalaryGrade | Grade gaji |
| TER | Tarif Efektif Rata-rata (pajak) |
| PTKP | Penghasilan Tidak Kena Pajak |
| Insurance | Asuransi |
| CompanyHoliday | Hari libur perusahaan |
| Competency | Kompetensi (master — juga ada di modul Competency) |

### Document Template

| Entitas | Deskripsi |
|---|---|
| DocumentTemplate | Template dokumen (name, module_slug, movement_types, is_active) |
| DocumentTemplateVersion | Versi template (content, variables, is_active) |
| DocumentTemplateAudit | Audit trail perubahan template |
| GeneratedDocument | Dokumen yang sudah di-generate (ref_type, ref_id, file_url, template_id, version_id) |

---

## 3. SETTINGS — Master Data CRUD

Semua master data mengikuti pola yang sama: `GET /settings/<resource>`, `POST`, `GET /:id`, `PUT /:id`, `DELETE /:id`.

| Resource | Endpoint Prefix | Catatan |
|---|---|---|
| Zones | `/settings/zones` | + field `timezone` opsional (override — lihat §9 Timezone Settings) |
| Provinces | `/settings/provinces` | |
| Regencies | `/settings/regencies` | |
| Districts | `/settings/districts` | |
| Villages | `/settings/villages` | + `GET /search` + `GET /:id/detail` |
| Education | `/settings/educations` | |
| Education Majors | `/settings/education-majors` | |
| Religion | `/settings/religions` | |
| Marital Status | `/settings/marital-statuses` | |
| Relationship Type | `/settings/relationship-types` | |
| Employment Status | `/settings/employment-statuses` | |
| Bank | `/settings/banks` | |
| Nationality | `/settings/nationalities` | |
| Job Family | `/settings/job-families` | |
| Grading | `/settings/gradings` | |
| Salary Grade | `/settings/salary-grades` | |
| TER | `/settings/ters` | Tarif Efektif Rata-rata |
| PTKP | `/settings/ptkps` | Penghasilan Tidak Kena Pajak |
| Insurance | `/settings/insurances` | |
| Company Holiday | `/settings/company-holidays` | |
| Competency | `/settings/competencies` | Master kompetensi |

---

## 4. DOCUMENT TEMPLATE — CRUD

1. **Buat template** — `POST /settings/document-templates`: name, module_slug, description.
2. **Daftar template** — `GET /settings/document-templates`.
3. **Detail template** — `GET /settings/document-templates/:id`.
4. **Update template** — `PUT /settings/document-templates/:id`.
5. **Hapus template** — `DELETE /settings/document-templates/:id`.

---

## 5. DOCUMENT TEMPLATE — Versioning

1. **Daftar versi** — `GET /settings/document-templates/:id/versions`.
2. **Buat versi** — `POST /settings/document-templates/:id/versions`: content (HTML/template string), variables (JSON).
3. **Detail versi** — `GET /settings/document-templates/:id/versions/:versionId`.

---

## 6. DOCUMENT TEMPLATE — Activate / Preview

1. **Activate** — `POST /settings/document-templates/:id/activate`: set versi aktif.
2. **Deactivate** — `POST /settings/document-templates/:id/deactivate`.
3. **Preview** — `POST /settings/document-templates/:id/preview`: preview dokumen dengan sample data.
4. **Movement types** — `GET /settings/document-templates/movement-types`: daftar tipe movement untuk filter.
5. **Variables** — `GET /settings/document-templates/variables`: daftar variabel yang tersedia.

---

## 7. DOCUMENT NUMBERING

| Endpoint | Deskripsi |
|---|---|
| `GET /settings/document-numbering` | Daftar format penomoran per document_type |
| `PUT /settings/document-numbering/:document_type` | Update format (prefix, counter, padding) |
| `GET /settings/document-numbering/:document_type/preview` | Preview nomor berikutnya |

---

## 8. EMPLOYEE ID FORMAT

| Endpoint | Deskripsi |
|---|---|
| `GET /settings/employee-id-format` | Konfigurasi mode (auto/hybrid/manual) |
| `PUT /settings/employee-id-format` | Update konfigurasi |
| `GET /settings/employee-id-format/preview` | Preview format ID |

---

## 9. TIMEZONE SETTINGS

Zona waktu tenant, dipakai untuk menentukan boundary tanggal (mis. "hari ini") pada modul
transaksi yang tidak punya input tanggal dari client, dan untuk menampilkan jam/tanggal
berjalan di header aplikasi.

**Cakupan:** dibatasi ke 3 zona Indonesia — WIB (`Asia/Jakarta`), WITA (`Asia/Makassar`),
WIT (`Asia/Jayapura`). Nilai disimpan sebagai string IANA; label WIB/WITA/WIT hanya untuk
tampilan.

### Dua level konfigurasi

1. **Company default** (`companies.timezone`, platform DB) — wajib diisi, berlaku untuk
   seluruh tenant jika tidak ada override.
2. **Zone override** (`zones.timezone`, tenant DB, opsional) — di-set per baris Zone.
   Employee mewarisi zona lewat rantai `Employee → Organization → Zone`. Kosongkan
   ("Ikut default perusahaan") untuk kembali ikut default company.

### Endpoint

| Endpoint | Deskripsi |
|---|---|
| `GET /settings/company/timezone` | Baca timezone default company |
| `PUT /settings/company/timezone` | Update timezone default company |
| `POST/PUT /settings/zones` | Field `timezone` (opsional) pada create/update Zone — override untuk zona tsb |
| `GET /attendance/timezone/me` | Zona waktu efektif user yang login (dipakai jam berjalan di header) |

### Resolusi zona waktu

Package `internal/pkg/timezone` (`Resolve`) menentukan zona efektif dengan prioritas:
**Zone override → Company default → fallback `Asia/Jakarta`**. Semua timestamp tetap
disimpan **UTC** di database — resolusi ini hanya dipakai untuk interpretasi
tanggal/tampilan, bukan penulisan timestamp.

**Penerapan saat ini:**
- Attendance: query "attendance hari ini" (dashboard/ringkasan tanpa input tanggal dari
  client) dan validasi clock-skew saat check-in/check-out (device time vs server, toleransi
  5 menit — melewati batas ditandai `INVALID`, tidak diblokir).
- Header aplikasi: jam & tanggal berjalan (`LiveClock.vue`) mengikuti zona efektif user
  yang login (via organization → zone → company).
- **Belum diterapkan** (fase berikutnya, per rollout plan): Payroll cutoff, Leave/cuti.

---

## 10. Integrasi Lintas Modul

| Modul | Peran Settings |
|---|---|
| **Employee** | Education, Religion, Marital Status, Bank, Nationality, Employment Status untuk data karyawan |
| **Organization** | Job Family, Grading untuk struktur posisi |
| **Job Management** | Job Family, Education, Grading untuk analisis jabatan |
| **Employee Movement** | Employment Status untuk movement from/to |
| **Payroll** | TER, PTKP, Salary Grade, Insurance untuk kalkulasi gaji |
| **Attendance** | Company Holiday untuk kalkulasi sesi; Timezone untuk boundary tanggal & clock-skew check |
| **All modules** | Document Template untuk generate PDF (SK, kontrak, dll) |

---

## 11. Peta Halaman UI

| Menu | Halaman |
|---|---|
| Settings (hub) | Settings hub page |
| Zones (+ timezone override) | `ZonesView.vue` |
| Company Profile (timezone default) | `CompanyDetail.vue` — bagian "Regional Settings" |
| Provinces / Regencies / Districts / Villages | Geografi pages |
| Education / Majors | Education pages |
| Religion / Marital Status / Relationship Type | Personal pages |
| Employment Status / Bank / Nationality | Employment pages |
| Job Family / Grading / Salary Grade | Compensation pages |
| TER / PTKP / Insurance | Tax & Insurance pages |
| Company Holiday | Holiday pages |
| Competency | Competency master pages |
| Document Templates | `DocumentTemplates.vue` |
| Document Numbering | Numbering config page |
| Employee ID Format | ID format config page |

---

## 12. Catatan Penting

- **Semua master data Settings** bersifat CRUD standar tanpa workflow — bisa di-setup sekali dan diubah sesuai kebutuhan.
- **Document Template** punya versioning — satu template bisa punya banyak versi, hanya 1 yang aktif.
- **Generated Document** disimpan sebagai referensi (ref_type + ref_id) → modul lain (Employee Movement, Contract) yang memanggil generate.
- **Employee ID Format** punya 3 mode: auto (system generate), hybrid (system generate + manual override), manual (user isi sendiri).
- **Document Numbering** dikonfigurasi per document_type (mis. SK_PROMOTION, CONTRACT, dll) dengan format prefix + counter + padding.
- **Timezone** dibatasi ke 3 zona Indonesia (WIB/WITA/WIT); semua timestamp tetap UTC di database — resolusi zona hanya untuk interpretasi tanggal & tampilan, bukan penulisan data.
