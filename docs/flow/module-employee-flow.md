# Alur Pengisian Employee Management (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Employee Management** — CRUD data
karyawan (profil inti + 8 sub-profil), upload foto, sensitive field settings, statistik —
polas runbook seperti [`module-leave-flow.md`](module-leave-flow.md) &
[`module-reimbursement-flow.md`](module-reimbursement-flow.md).

- Lokasi kode: `backend/internal/modules/employee/` · `frontend/tenant/src/views/modules/employee/`
- Modul terkait: Organization · Settings (Religion, Marital Status, Education, Employment Status, Insurance, Bank) · Recruitment · Attendance · Leave · Training · Payroll · Employee Movement

---

## 1. Ringkasan Alur End-to-End

```
BUAT KARYAWAN                    SUB-PROFIL                        INTEGRASI
┌────────────────────────┐   ┌───────────────────────────────┐   ┌──────────────────────────┐
│ Data Inti              │   │ Addresses (alamat)            │   │ Recruitment (offer →      │
│  - Employee ID         │──▶│ Emergency Contacts (kontak    │──▶│   employee baru)          │
│  - NIK, Nama, Gender   │   │   darurat)                    │   │ Attendance (exempt,       │
│  - DOB, POB, Agama     │   │ Families (keluarga)           │   │   sessions, shifts)       │
│  - Status aktif/tidak  │   │ Educations (pendidikan)       │   │ Leave (balance, policy)   │
│  - Photo               │   │ Experiences (pengalaman)      │   │ Training (participant)    │
│                        │   │ Documents (dokumen)           │   │ Payroll (gaji, bank)      │
│                        │   │ Insurances (asuransi)         │   │ Employee Movement         │
│                        │   │ Bank Accounts (rekening)      │   │   (promosi, mutasi)       │
│                        │   │ Employments (riwayat jabatan) │   │ User Account (self-service)│
└────────────────────────┘   └───────────────────────────────┘   └──────────────────────────┘
```

- **Employee ID** bisa di-generate otomatis (AUTO), semi-otomatis (HYBRID), atau manual (MANUAL)
  — tergantung konfigurasi `EmployeeIDFormatProvider`.
- **Status karyawan:** `active` / `inactive`.
- **Sensitive Field Settings** — toggle enkripsi at-rest per field (NIK, email, dll.) via admin.
- **Recruited From Application** — karyawan dari offer recruitment eksternal punya referensi
  ke `job_applications` (traceability recruitment → employee).

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Employee | `employees` | Data inti karyawan (employee_id, NIK, nama, gender, DOB, POB, agama, status, foto, dll.) |
| Employee Address | `employee_addresses` | Alamat karyawan (tipe HOME/DORMITORY, alamat, provinsi/kabupaten/kecamatan/kelurahan, kode pos) |
| Emergency Contact | `emergency_contacts` | Kontak darurat (nama, hubungan, telepon, alamat) |
| Employee Family | `employee_families` | Data keluarga (NIK, nama, DOB, hubungan, pendidikan) |
| Employee Education | `employee_educations` | Riwayat pendidikan (institusi, jurusan, tahun lulus) |
| Employee Experience | `employee_experiences` | Pengalaman kerja (perusahaan, posisi, tahun mulai/akhir) |
| Employee Document | `employee_documents` | Dokumen karyawan (nama, file URL, catatan) |
| Employee Insurance | `employee_insurances` | Asuransi karyawan (nomor polis, tipe, referensi ke master Insurance) |
| Employee Bank Account | `employee_bank_accounts` | Rekening bank (nomor rekening, nama rekening, referensi ke master Bank) |
| Employment | `employments` | Riwayat jabatan (organisasi, posisi, status, SK, tanggal efektif) |

---

## 3. TAHAP 1 — BUAT KARYAWAN (Data Inti)

Menu **Employee → Employee List** (`/admin/employees`, `Employees.vue`) → tombol **Create** →
halaman create (`EmployeeForm.vue`).

### A. Form Data Inti

Isi field wajib & opsional:

| Field | Wajib | Keterangan |
|---|---|---|
| Employee ID | ✅ | Auto-generate (AUTO/HYBRID) atau manual tergantung konfigurasi |
| NIK | | Nomor Induk Kependudukan (16 digit) |
| Family ID | | Kartu Keluarga |
| Nama Lengkap | ✅ | Nama karyawan |
| Nama Ibu | | Nama ibu kandung |
| Gender | | Laki-laki / Perempuan |
| Nationality Type | | WNI / WNA |
| Nationality | | Kewarganegaraan (referensi master Nationality) |
| Passport | | Nomor paspor (untuk WNA) |
| Tempat Lahir | | Kota/kabupaten kelahiran |
| Tanggal Lahir | | Format YYYY-MM-DD |
| No. Telepon | | Nomor telepon/HP |
| Email | | Email karyawan (unique) |
| LinkedIn | | Profil LinkedIn |
| Instagram | | Akun Instagram |
| Agama | | Referensi master Religion |
| Status Perkawinan | | Referensi master Marital Status |
| Status Karyawan | | `active` / `inactive` |

### B. Upload Foto

- Upload foto profil karyawan: `PUT /employees/:id/photo`.
- Hapus foto: `DELETE /employees/:id/photo`.

### C. Employee ID Generation

- **AUTO** — sistem generate employee_id otomatis (format dari setting).
- **HYBRID** — kombinasi prefix manual + nomor otomatis.
- **MANUAL** — user mengisi sendiri.
- Provider di-inject dari main.go via `SetEmployeeIDFormatProvider`.

---

## 4. TAHAP 2 — SUB-PROFIL KARYAWAN

Setelah employee dibuat, kelola sub-profil di halaman detail (`EmployeeDetail.vue`).
Setiap sub-profil punya form terpisah dengan CRUD-nya sendiri.

### A. Addresses (Alamat)

Form: `AddressForm.vue`. Tipe: `HOME` / `DORMITORY`.

- Alamat lengkat, Provinsi, Kabupaten, Kecamatan, Kelurahan, Kode Pos.
- Referensi: master Provinsi, Regency, District, Village (modul Settings).
- Endpoint: `POST /employees/:id/addresses`, `PUT /employees/:id/addresses/:addressId`, `DELETE /employees/:id/addresses/:addressId`.

### B. Emergency Contacts (Kontak Darurat)

Form: `ContactForm.vue`.

- Nama, Hubungan (referensi master Relationship Type), No. Telepon, Alamat.
- Endpoint: `POST /employees/:id/emergency-contacts`, `PUT /employees/:id/emergency-contacts/:contactId`, `DELETE /employees/:id/emergency-contacts/:contactId`.

### C. Families (Keluarga)

Form: `FamilyForm.vue`.

- NIK, Nama, Tanggal Lahir, Hubungan (referensi master Relationship Type),
  Pendidikan Terakhir (referensi master Education).
- Endpoint: `POST /employees/:id/families`, `PUT /employees/:id/families/:familyId`, `DELETE /employees/:id/families/:familyId`.

### D. Educations (Pendidikan)

Form: `EducationForm.vue`.

- Institusi, Jurusan (referensi master Education Major), Tahun Lulus.
- Referensi: master Education, Education Major (modul Settings).
- Endpoint: `POST /employees/:id/educations`, `PUT /employees/:id/educations/:educationId`, `DELETE /employees/:id/educations/:educationId`.

### E. Experiences (Pengalaman Kerja)

Form: `ExperienceForm.vue`.

- Nama Perusahaan, Posisi, Tahun Mulai, Tahun Selesai.
- Endpoint: `POST /employees/:id/experiences`, `PUT /employees/:id/experiences/:experienceId`, `DELETE /employees/:id/experiences/:experienceId`.

### F. Documents (Dokumen)

Form: `DocumentForm.vue`.

- Nama Dokumen, File (upload dua-langkah via `POST /uploads` → simpan URL), Catatan.
- Upload file: `POST /employees/:id/documents/upload`.
- Update file: `PUT /employees/:id/documents/:documentId/upload`.
- Endpoint: `POST /employees/:id/documents`, `PUT /employees/:id/documents/:documentId`, `DELETE /employees/:id/documents/:documentId`.

### G. Insurances (Asuransi)

Form: `InsuranceForm.vue`.

- Asuransi (referensi master Insurance), Nomor Polis, Tipe.
- Endpoint: `POST /employees/:id/insurances`, `PUT /employees/:id/insurances/:insuranceId`, `DELETE /employees/:id/insurances/:insuranceId`.

### H. Bank Accounts (Rekening Bank)

Form: `BankAccountForm.vue`.

- Bank (referensi master Bank), Nomor Rekening, Nama Rekening.
- Endpoint: `POST /employees/:id/banks`, `PUT /employees/:id/banks/:bankId`, `DELETE /employees/:id/banks/:bankId`.

### I. Employments (Riwayat Jabatan)

Form: `EmploymentForm.vue`.

- Organisasi (referensi Organization), Posisi, Status Employment (referensi master Employment Status),
  Nomor SK, Tanggal SK, Tanggal Efektif, Tanggal Akhir Efektif.
- Endpoint: `POST /employees/:id/employments`, `PUT /employees/:id/employments/:employmentId`, `DELETE /employees/:id/employments/:employmentId`.

> ⚠️ **Employment** adalah data historis — bisa ada banyak record per karyawan (riwayat
> promosi/mutasi). Employment terakhir yang dianggap sebagai posisi saat ini.

---

## 5. TAHAP 3 — STATISTIK & SENSITIVE FIELDS

### A. Statistik Karyawan

- **Gender Stats** — distribusi karyawan berdasarkan gender: `GET /employees/stats/gender`.
- **Employment Status Stats** — distribusi berdasarkan status employment: `GET /employees/stats/employment-status`.

### B. Sensitive Field Settings

Menu **Employee → Sensitive Field Settings** (admin only).

- Toggle enkripsi at-rest per field sensitif (NIK, email, no. telepon, dll.).
- Permission khusus: `setting.sensitive-fields.view` / `setting.sensitive-fields.manage`.
- Endpoint: `GET /employees/settings/sensitive-fields`, `PUT /employees/settings/sensitive-fields/:fieldKey`.

> ⚠️ **Enkripsi** bersifat opt-in per field — hanya field yang di-toggle yang dienkripsi.
> Data terenkripsi hanya bisa dilihat oleh user dengan permission `setting.sensitive-fields.view`.

---

## 6. Ringkasan Status & Transisi

| Entitas | Status | Transisi |
|---|---|---|
| **Employee** | `active` / `inactive` | manual via `PUT /employees/:id` |
| **Employment** | Tidak ada status — data historis | create/update/delete per riwayat jabatan |

---

## 7. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Organization** | Employee ditugaskan ke organization node (via `organization_id` di Employment) |
| **Settings** | Master data: Religion, Marital Status, Relationship Type, Education, Education Major, Employment Status, Insurance, Bank, Nationality, Province, Regency, District, Village |
| **Recruitment** | Employee baru bisa dibuat dari offer recruitment eksternal (`recruited_from_application_id`) |
| **Attendance** | Employee terdaftar di shift, session, dan exempt positions |
| **Leave** | Employee punya balance cuti per type |
| **Training** | Employee terdaftar sebagai participant training |
| **Payroll** | Employee punya data gaji, tunjangan, potongan, rekening bank |
| **Employee Movement** | Promosi, mutasi, rotasi mengubah employment record |
| **User Account** | Employee login dan mengakses self-service profile (`/user-accounts/me`) |
| **Workforce Intelligence** | Analisis workforce per employee |
| **Career Intelligence** | Career path per employee |

---

## 8. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Employee (hub) | `Employees.vue` | Daftar karyawan + search + filter + tombol Create |
| Employee → Create | `EmployeeForm.vue` | Form buat karyawan baru |
| Employee → Detail | `EmployeeDetail.vue` | Detail karyawan + tab sub-profil (Addresses, Emergency Contacts, Families, Educations, Experiences, Documents, Insurances, Bank Accounts, Employments) |
| Employee → Edit | `EmployeeForm.vue` | Form edit data inti karyawan |
| Employee → Sub-forms | `AddressForm.vue`, `ContactForm.vue`, `FamilyForm.vue`, `EducationForm.vue`, `ExperienceForm.vue`, `DocumentForm.vue`, `InsuranceForm.vue`, `BankAccountForm.vue`, `EmploymentForm.vue` | Form CRUD per sub-profil |
| Employee → List Card | `ListCard.vue` | Kartu tampilan alternatif |
| Employee → Detail Row | `DetailRow.vue` | Baris detail di expansion |
| Employee → Payroll Profiles | `PayrollProfilesForm.vue` | Profil payroll karyawan |
| Employee → Salary Structure | `SalaryStructureForm.vue` | Struktur gaji |
| Employee → Account | `AccountForm.vue` | Form akun user |
| Employee → Personal | `PersonalForm.vue` | Form data pribadi |

---

## 9. Endpoint API Utama

Semua di bawah `/api/v1/tenant/employees/`.

| Area | Endpoint |
|---|---|
| **CRUD** | `POST /employees`, `GET /employees`, `GET /employees/:id`, `PUT /employees/:id`, `DELETE /employees/:id` |
| **Photo** | `PUT /employees/:id/photo`, `DELETE /employees/:id/photo` |
| **Stats** | `GET /employees/stats/gender`, `GET /employees/stats/employment-status` |
| **Addresses** | `POST /employees/:id/addresses`, `PUT /employees/:id/addresses/:addressId`, `DELETE /employees/:id/addresses/:addressId` |
| **Emergency Contacts** | `POST /employees/:id/emergency-contacts`, `PUT /employees/:id/emergency-contacts/:contactId`, `DELETE /employees/:id/emergency-contacts/:contactId` |
| **Families** | `POST /employees/:id/families`, `PUT /employees/:id/families/:familyId`, `DELETE /employees/:id/families/:familyId` |
| **Educations** | `POST /employees/:id/educations`, `PUT /employees/:id/educations/:educationId`, `DELETE /employees/:id/educations/:educationId` |
| **Experiences** | `POST /employees/:id/experiences`, `PUT /employees/:id/experiences/:experienceId`, `DELETE /employees/:id/experiences/:experienceId` |
| **Documents** | `POST /employees/:id/documents`, `POST /employees/:id/documents/upload`, `PUT /employees/:id/documents/:documentId`, `PUT /employees/:id/documents/:documentId/upload`, `DELETE /employees/:id/documents/:documentId` |
| **Insurances** | `POST /employees/:id/insurances`, `PUT /employees/:id/insurances/:insuranceId`, `DELETE /employees/:id/insurances/:insuranceId` |
| **Bank Accounts** | `POST /employees/:id/banks`, `PUT /employees/:id/banks/:bankId`, `DELETE /employees/:id/banks/:bankId` |
| **Employments** | `POST /employees/:id/employments`, `PUT /employees/:id/employments/:employmentId`, `DELETE /employees/:id/employments/:employmentId` |
| **Sensitive Fields** | `GET /employees/settings/sensitive-fields`, `PUT /employees/settings/sensitive-fields/:fieldKey` |

---

## 10. Catatan Penting

- **Employee ID** bisa di-generate otomatis (AUTO), semi-otomatis (HYBRID), atau manual (MANUAL) —
  tergantung konfigurasi `EmployeeIDFormatProvider` yang di-inject dari main.go.
- **Employment** adalah data historis — bisa ada banyak record per karyawan (riwayat promosi/mutasi).
  Employment terakhir (berdasarkan `effective_date` terbaru) dianggap sebagai posisi saat ini.
- **Recruited From Application** — employee dari offer recruitment eksternal punya referensi
  `recruited_from_application_id` yang menelusuri ke Application → Requisition → Position.
- **Sensitive Field Settings** bersifat opt-in per field — hanya field yang di-toggle yang
  dienkripsi. Data terenkripsi hanya bisa dilihat oleh user dengan permission khusus.
- **Photo** di-upload via `PUT /employees/:id/photo` (multipart form).
- **Document** di-upload dua-langkah: `POST /uploads` (file) → simpan URL ke document record.
- **Isi master data (Settings) terlebih dahulu** — Religion, Marital Status, Education,
  Employment Status, Insurance, Bank, dll. dipilih saat create/edit employee dan sub-profil.
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
