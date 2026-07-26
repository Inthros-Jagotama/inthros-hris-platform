# Panduan & Standar UI/UX HRIS Enterprise Grade

> **Tech Stack:** Vue 3 (`<script setup>`) + PrimeVue v4+ + Tailwind CSS  
> **Konsep:** High-Density, Compact, Modal-First, & Single-Page Dashboard (Bebas Pindah Halaman)

---

## 1. Konsep Utama & Prinsip Desain

Dalam sistem HRIS skala Enterprise, produktivitas pengguna sangat bergantung pada kecepatan eksekusi data. **Prinsip Utama:** Pengguna tidak boleh sering berpindah halaman penuh (*page navigation/routing*) hanya untuk melihat detail, mengedit data, atau melakukan persetujuan (*approval*).

| Prinsip UX | Implementasi Teknikal | Manfaat Utama |
| :--- | :--- | :--- |
| **Modal & Dialog First** | `Dialog` (Modal) & `Drawer` (Slide-over) dari PrimeVue. | Mencegah perpindahan halaman & menjaga konteks pengguna di tabel utama. |
| **Single Viewport (Zero Scroll)** | Master-Detail (Split View) & Tabbed Interface. | Semua informasi penting terlihat tanpa perlu *scroll* panjang. |
| **High Density Layout** | Padding ketat (`py-1.5`, `px-3`), font size `12px`–`14px` (`text-xs`–`text-sm`). | Memaksimalkan data yang tampil dalam 1 layar. |

---

## 2. Hirarki Penggunaan Modal & Overlay (Kapan Menggunakan Apa?)

Agar tidak terjadi *modal-cluttering* (terlalu banyak popup menumpuk), ikuti aturan penggunaan overlay berikut:

### A. Compact `Dialog` (Modal Standar)
* **Penggunaan:** Form ringkas (1-2 kolom), konfirmasi aksi (*Approve/Reject*), ubah status, atau input catatan cepat.
* **Properti PrimeVue:** `<Dialog modal :style="{ width: '400px' }" header="...">`
* **Keunggulan:** Fokus instan ke aksi spesifik tanpa mengganggu layar latar belakang.

### B. Slide-over `Drawer` (Modal Samping)
* **Penggunaan:** Form data kompleks (misal: *Tambah Karyawan Baru*, *Detail Gaji & Komponen*).
* **Properti PrimeVue:** `<Drawer v-model:visible="open" position="right" class="!w-[500px]">`
* **Keunggulan:** Memberikan area vertikal lebih luas untuk form ber-tab tanpa menutupi seluruh layar utama.

### C. `Popover` / Context Menu
* **Penggunaan:** Menu tindakan cepat pada baris tabel (Edit, Hapus, Kirim Slip Gaji, Cetak ID).
* **Keunggulan:** Menghemat kolom pada `DataTable` sehingga tabel tidak terasa sesak.

---

## 3. Master Prompt untuk AI Assistant (Diperbarui)

Gunakan prompt acuan di bawah ini untuk disalin ke AI (ChatGPT / Claude / Gemini) agar seluruh kode yang dihasilkan selalu konsisten dengan standar HRIS Anda:

```text
[CONTEXT & ROLE]
Anda adalah Senior Full-Stack Frontend Engineer & UX Designer spesialis sistem Enterprise ERP/HRIS. Tugas Anda adalah membantu saya membangun aplikasi Human Resource Information System (HRIS) Enterprise Grade yang modern, minimalis, efisien ruang (high-density), user-friendly, bebas dari scroll berlebih, dan SANGAT MINIM PINDAH HALAMAN.

[TECH STACK]
- Framework: Vue 3 (Composition API / <script setup>)
- UI Library: PrimeVue (Gunakan skema komponen modern v4+, Pass-Through, atau Design Tokens)
- Utility Styling: Tailwind CSS
- Icons: PrimeIcons / Lucide Icons

[DESIGN & UX PRINCIPLES - HARUS DITURUTI]
1. Minimum Page Navigation & Modal-First Architecture:
   - UTAMAKAN penggunaan Dialog/Modal, Drawer (Slide-over), dan Popover alih-alih berpindah halaman/route baru.
   - Semua aksi CRUD (Create, Read, Update, Delete) dan Approval HARUS dilakukan di atas halaman utama menggunakan Modal/Drawer.
   - Pengguna harus tetap berada di tabel/dashboard utama saat melakukan operasi data.

2. Compact & High-Density First:
   - Hindari ruang kosong (whitespace) yang tidak perlu.
   - Gunakan padding dan margin yang ketat (misal: py-1.5, py-2, px-3, p-3).
   - Ukuran teks dominan adalah text-xs (12px) sampai text-sm (14px).
   - Gunakan border-radius yang kecil/presisi (rounded-md atau rounded-sm).

3. Zero / Minimal Scroll Architecture:
   - Manfaatkan layout 'h-screen' atau 'h-[calc(100vh-...)]' dengan 'overflow-hidden' pada container utama.
   - Gunakan pola Split-Screen / Master-Detail (Panel Kiri: DataTable, Panel Kanan: Detail View/Tabs).
   - Gunakan komponen 'Tabs' horizontal di dalam Modal/Drawer untuk memecah form berukuran besar.

4. Standar Komponen PrimeVue & Tailwind:
   - DataTable: Gunakan opsi compact/small (`p-datatable-sm`), scrollable dengan `scrollHeight="flex"`, virtual scroll jika perlu, dan pinned/frozen columns untuk aksi.
   - Dialog & Drawer: Selalu atur ukuran lebar eksplisit (misal: !w-[450px] untuk Dialog, !w-[550px] untuk Drawer).
   - Filters & Actions: Satukan Search Bar, Filter Select, dan Quick Action Buttons dalam satu baris horizontal menggunakan `Toolbar` atau Flexbox ringkas.
   - Warna Status: Gunakan Soft Badge dengan Tailwind (misal: `bg-emerald-50 text-emerald-700` untuk aktif/disetujui, `bg-amber-50 text-amber-700` untuk pending).

[OUTPUT FORMAT REQUIREMENT]
Setiap kali memberikan kodingan atau solusi UI/UX:
1. Tuliskan kode Vue 3 lengkap menggunakan Single File Component (SFC) <template> dan <script setup>.
2. Pastikan gabungan kelas Tailwind CSS dan atribut PrimeVue mematuhi aturan high-density & modal-first di atas.
3. Berikan penjelasan singkat mengenai struktur layout jika ada pola UX khusus yang diterapkan.
```

---

## 4. Standar Warna Badge Status

Gunakan palet warna status yang lembut (*soft badges*) agar tampilan tetap profesional dan tidak mencolok secara berlebihan:

* **Active / Approved:** `bg-emerald-50 text-emerald-700 border border-emerald-200`
* **Pending / Review:** `bg-amber-50 text-amber-700 border border-amber-200`
* **Draft / Info:** `bg-indigo-50 text-indigo-700 border border-indigo-200`
* **Inactive / Rejected:** `bg-rose-50 text-rose-700 border border-rose-200`
* **Draft:** `bg-gray-50 text-gray-700 border border-gray-200`
* **Published / Paid:** `bg-sky-50 text-sky-700 border border-sky-200`
* **Expired / Terminated:** `bg-red-50 text-red-700 border border-red-200`

---

## 5. Frontend Development Plan — Roadmap Implementasi

Dokumen ini berisi roadmap implementasi frontend untuk Platform Admin dan Tenant HRIS, disusun berdasarkan prioritas bisnis dan ketergantungan backend.

> **Tech Stack:** Vue 3 (Composition API + `<script setup>`) + PrimeVue 4 + Tailwind CSS 4 + Vite + Axios
> **Lihat juga:** `docs/frontend-development-plan.md` untuk versi lengkap (termasuk checklist item per halaman)

---

## 5.1. Ringkasan Arsitektur Frontend

```
frontend/
├── platform-admin/          # Admin panel (Vue 3 + PrimeVue)
│   ├── src/
│   │   ├── App.vue          # Root - Toast + ConfirmDialog
│   │   ├── main.js          # Entry point
│   │   ├── router/index.js  # Routes + Auth guard
│   │   ├── layouts/         # AppLayout, HeaderBar, Sidebar
│   │   ├── views/           # Dashboard, Companies, Users, dll
│   │   ├── components/      # FormRow, TextInput, ToggleSwitch, dll
│   │   ├── services/api.js  # Axios instance
│   │   └── stores/auth.js   # Pinia auth store
│   └── package.json
│
└── tenant/                  # Tenant dashboard (Vue 3 + PrimeVue)
    ├── src/
    │   ├── App.vue          # Root - Toast + ConfirmDialog
    │   ├── main.js          # Entry point
    │   ├── router/index.js  # Routes for all 15 modules
    │   ├── layouts/         # AppLayout, HeaderBar, Sidebar
    │   └── views/
    │       ├── Dashboard.vue
    │       └── modules/     # 15 module views (mostly placeholders)
    └── package.json
```

### Stack Bersama
- **Framework:** Vue 3 (Composition API + `<script setup>`)
- **UI Library:** PrimeVue 4 (DataTable, Dialog, Tag, Button, InputText, Select, dll)
- **Styling:** Tailwind CSS 4 + PrimeVue themes (Aura)
- **HTTP:** Axios (interceptors untuk JWT + refresh token)
- **State:** Pinia (auth store)
- **Router:** Vue Router 4 (auth guard, lazy loading)
- **Icons:** PrimeIcons

---

## 5.2. Phase 1 — Platform Admin (MVP)

### 5.2.1. Login Page ✅ (Existing)
- Form login (email + password)
- JWT token management (access + refresh)
- Redirect ke dashboard setelah login
- Auto-redirect ke login jika token expired

### 5.2.2. Dashboard Page ✅ (Existing — Perlu Penyempurnaan)
- KPI Cards (Total Companies, Active Tenants, Users, Modules, Connections, Health)
- Recent Companies list
- System Health (DB, Cache, Pool Stats)
- Quick Actions
- **⬜ Enhancement:** Grafik/tren jumlah company per bulan (PrimeVue Chart)
- **⬜ Enhancement:** Error rate / request count dari monitoring
- **⬜ Enhancement:** Auto-refresh tiap 30s via polling

### 5.2.3. Companies Page ✅ (Existing — CRUD)
- DataTable daftar companies (pagination)
- Create company dialog (name, slug, package selection)
- Company detail / edit via Dialog
- Status management (Activate, Suspend, Terminate)
- Status badges (Tag component)
- **⬜ Enhancement:** Filter by status (active/suspended/terminated)
- **⬜ Enhancement:** Filter by package (`package_id`)
- **⬜ Enhancement:** Search by name/slug
- **⬜ Enhancement:** License info inline (current package, plan type)
- **⬜ Enhancement:** Tenant provisioning progress indicator

### 5.2.4. Users Page ✅ (Existing — CRUD)
- DataTable daftar platform users
- Create / Edit user dialog
- **⬜ Enhancement:** Filter by role (super_admin, company_admin)
- **⬜ Enhancement:** Search by name/email
- **⬜ Enhancement:** Bulk actions

### 5.2.5. Modules Page ✅ (Existing — Perlu Update)
- DataTable daftar modules
- Create module dialog
- Module detail
- Activate/deactivate untuk company
- Status badges
- **⬜ Enhancement:** Filter `?module_type=platform|tenant` (dropdown)
- **⬜ Enhancement:** Kolom `module_type` di tabel (badge: "Platform" / "Tenant")
- **⬜ Enhancement:** Kolom `depends_on` yang bisa diexpand
- **⬜ Enhancement:** Quick filter chips: "All" | "Platform" | "Tenant"

### 5.2.6. Licenses Page ✅ (Existing — Perlu Update)
- DataTable daftar licenses
- Create license dialog
- License detail
- **⬜ Enhancement:** Kolom `package_id` + `package_name` (integrasi package)
- **⬜ Enhancement:** Filter by package
- **⬜ Enhancement:** Status filter (active/expired/suspended)
- **⬜ Enhancement:** License key copy button
- **⬜ Enhancement:** Expiration date warning (Tag: "Expiring Soon")

### 5.2.7. Monitoring Page ✅ (Existing)
- Platform health status
- Database connectivity per tenant
- Pool stats
- **⬜ Enhancement:** Auto-refresh toggle
- **⬜ Enhancement:** Grafik pool utilization over time
- **⬜ Enhancement:** Alert thresholds

### 5.2.8. Packages Page 🔴 (BARU — High Priority)
- **⬜** DataTable daftar packages (CRUD)
  - Kolom: Name, Slug, Price, Status (Draft/Published), Module Count, Sort Order
  - Filter `?module_type=platform|tenant`
  - Search by name
- **⬜** Create Package Dialog
  - Form: Name, Slug, Description, Price, Sort Order
  - Module selector: MultiSelect dengan search
    - Tampilkan module_name, module_slug, module_type
    - Validasi depends_on otomatis sebelum publish
  - IsMandatory toggle per module
  - Sort Order input per module
- **⬜** Edit Package Dialog (sama seperti create, pre-filled)
- **⬜** Publish / Unpublish button (dengan konfirmasi)
  - Validasi dependensi sebelum publish
  - Tampilkan hasil validasi (ModuleDependency list)
- **⬜** Delete Package (soft delete, dengan konfirmasi)
- **⬜** Validate Dependencies button
  - Modal daftar dependensi per module
  - Color-coded: resolved (hijau) / unresolved (merah)
- **⬜** Status badge: Draft (gray), Published (green)

### 5.2.9. RBAC Management Page 🔴 (BARU — Medium Priority)
- **⬜** Daftar roles (DataTable)
- **⬜** Create role dialog
- **⬜** Role detail page (permissions matrix)
- **⬜** Permission assignment (checkbox grid)
- **⬜** Delete role (non-system only, dengan konfirmasi)

---

## 5.3. Phase 2 — Tenant Module Views

### 5.3.1. Layout & Navigation ✅ (Existing)
- Sidebar with all module links
- HeaderBar dengan user info
- Responsive sidebar (collapsible)
- 15 module routes registered
- Dashboard dengan KPI cards + quick access

### 5.3.2. Dashboard ✅ (Existing — Mock Data)
- KPI Cards (Total Employees, Active Today, On Leave, Pending Approvals)
- Quick Access Modules grid (12 modules)
- Recent Activity (static)
- Period filter (This Month/Quarter/Year)
- **⬜** Ganti mock data dengan real API calls

### 5.3.3–5.3.17. Modul Tenant 🔴 (BARU)

| # | Modul | Endpoints | Kompleksitas | Fitur Utama |
|:-:|-------|:---------:|:------------:|-------------|
| 3 | **Organization** | 12 | 🟡 Medium | TreeTable organisasi, CRUD zones/positions/job families |
| 4 | **Employee** | 29 | 🔴 Complex | Wizard multi-step (9 sub-modul), tab detail per sub-modul |
| 5 | **Job Management** | 88 | 🔴 Complex | 18 sub-entities, tab layout, CRUD title/value/objective/risk/score |
| 6 | **Competency** | 35 | 🔴 Complex | CRUD competencies + values + events + scores |
| 7 | **Employee Movement** | 15 | 🟡 Medium | Workflow status (Approve/Execute/Cancel), Contracts PKWT/PKWTT |
| 8 | **Attendance** | 30 | 🟡 Medium | Shifts, geofence, overtime, check-in/out log |
| 9 | **Approval** | 15 | 🟡 Medium | Multi-step flows, My Tasks, Approve/Reject actions |
| 10 | **Payroll** | 47 | 🟡 Medium | Salary components, BPJS, PPh21, tax brackets, payslip view |
| 11 | **Leave** | 23 | 🟡 Medium | Leave types, accrual policies, balances, requests |
| 12 | **Performance** | 34 | 🔴 Complex | Periods, perspectives, templates, evaluations, KPI/OKR |
| 13 | **Recruitment** | 33 | 🔴 Complex | Job requisitions, candidates pipeline, interviews, onboarding |
| 14 | **Reimbursement** | 15 | 🟡 Medium | Workflow DRAFT→SUBMITTED→APPROVED→PAID |
| 15 | **Training** | 35 | 🟡 Medium | Categories, courses, sessions, participants, certificates |
| 16 | **Workforce Intelligence** | 68 | 🔴 Complex | Multiple dashboards (headcount, risk, executive), scenario planning |
| 17 | **Career Intelligence** | 19 | 🔴 Complex | 9-Box Talent Grid, talent maps, succession plans, gap analysis |
| 18 | **Package Subscription** | 3 | 🟡 Medium | Browse public packages, subscribe/unsubscribe, activated modules list |

> Detail lengkap setiap modul (termasuk sub-checklist fitur) ada di `docs/frontend-development-plan.md` Section C.

---

## 5.4. Shared Components Library

### 5.4.1. Komponen Reusable 🔴 (BARU)

| Komponen | Fungsi |
|----------|--------|
| **DataTableWrapper** | Wrapper PrimeVue DataTable dengan sorting, filter, pagination, search, export CSV, column visibility |
| **ConfirmDialog** | Reusable confirmation dengan PrimeVue ConfirmDialog |
| **StatusBadge** | Tag component dengan severity mapping otomatis berdasarkan status |
| **FormDialog** | Reusable form dialog (create/edit) |
| **SearchInput** | Input dengan icon search + debounce |
| **FilterChips** | Chip group untuk quick filter (All \| Active \| Inactive, dll) |
| **LoadingSkeleton** | Skeleton loading (PrimeVue Skeleton) |
| **EmptyState** | Komponen untuk empty table/message |
| **ApiError** | Error display component |
| **PageHeader** | Reusable page header (title, subtitle, actions) |

### 5.4.2. Services & Utilities

| Service | Status |
|---------|:------:|
| `api.js` — Axios instance dengan interceptor JWT + refresh token | ✅ Existing |
| `NotificationService` — wrapper toast untuk success/error/warning | ⬜ Perlu dibuat |
| `FormattingUtils` — date, currency, number formatting | ⬜ Perlu dibuat |
| `ValidationUtils` — form validation helpers (NIK, NPWP, dll) | ⬜ Perlu dibuat |

---

## 5.5. Prioritas Eksekusi

### Phase 1 — Platform Admin (MVP) — Estimasi: 2-3 minggu

| Prioritas | Feature | Kompleksitas |
|:---------:|---------|:------------:|
| 🚨 P0 | **Packages Page** (CRUD + Publish + Validasi Dependensi) | 🟡 Medium |
| 🚨 P0 | **Module Type Filter** di Modules Page | 🟢 Easy |
| 🔶 P1 | **Companies Filter** by status + package | 🟢 Easy |
| 🔶 P1 | **Licenses Package Integration** | 🟢 Easy |
| 🔶 P1 | **Dashboard Charts & Real-time** | 🟡 Medium |
| 🔵 P2 | **RBAC Management** | 🔴 Complex |

### Phase 2 — Tenant Core Modules — Estimasi: 4-6 minggu

| Prioritas | Feature | Kompleksitas |
|:---------:|---------|:------------:|
| 🚨 P0 | **Organization Management** (TreeTable) | 🟡 Medium |
| 🚨 P0 | **Employee Management** (Wizard 9 step) | 🔴 Complex |
| 🔶 P1 | **Leave & Attendance** | 🟡 Medium |
| 🔶 P1 | **Payroll** (read-only payslip) | 🟡 Medium |
| 🔵 P2 | **Job Management** (18 sub-entities) | 🔴 Complex |
| 🔵 P2 | **Competency Management** | 🔴 Complex |

### Phase 3 — Tenant Advanced Modules — Estimasi: 4-6 minggu

| Prioritas | Feature | Kompleksitas |
|:---------:|---------|:------------:|
| 🔶 P1 | **Performance Management** | 🔴 Complex |
| 🔶 P1 | **Recruitment (ATS Pipeline)** | 🔴 Complex |
| 🔵 P2 | **Approval Engine** (Flow Builder) | 🔴 Complex |
| 🔵 P2 | **Employee Movement** (Workflow) | 🟡 Medium |
| 🔵 P2 | **Training Management** | 🟡 Medium |
| 🔵 P2 | **Reimbursement** | 🟡 Medium |

### Phase 4 — Intelligence & Subscription — Estimasi: 2-3 minggu

| Prioritas | Feature | Kompleksitas |
|:---------:|---------|:------------:|
| 🔶 P1 | **Workforce Intelligence** Dashboards | 🔴 Complex |
| 🔶 P1 | **Career Intelligence** (9-box Grid) | 🔴 Complex |
| 🔶 P1 | **Package Subscription** (Tenant) | 🟡 Medium |

---

## 5.6. Struktur File Final (Target)

```
frontend/platform-admin/src/
├── components/
│   ├── form/           # FormRow, InputLabel, TextInput, SelectLabel, dll
│   ├── layout/         # AppLayout, Sidebar, HeaderBar
│   └── shared/         # DataTableWrapper, ConfirmDialog, StatusBadge, dll
├── composables/        # useApi, useNotification, usePagination, dll
├── router/index.js
├── services/api.js
├── stores/             # auth.js, companies.js, packages.js, dll
├── utils/              # formatters, validators (NIK, NPWP, dll)
└── views/
    ├── Dashboard.vue
    ├── Companies.vue
    ├── Users.vue
    ├── Modules.vue
    ├── Licenses.vue
    ├── Packages.vue         # BARU
    ├── Monitoring.vue
    └── Rbac.vue             # BARU

frontend/tenant/src/
├── components/
│   ├── layout/         # AppLayout, Sidebar, HeaderBar
│   └── shared/         # DataTableWrapper, FormDialog, ConfirmDialog
├── composables/        # useTenantApi, useModuleAccess, dll
├── router/index.js
├── services/api.js     # tenant-specific axios instance
├── stores/             # auth, subscription
├── utils/              # formatters
└── views/
    ├── Dashboard.vue              # Perlu real API
    └── modules/
        ├── Organizations.vue      # BARU
        ├── Employees.vue          # BARU (Wizard)
        ├── JobManagement.vue      # BARU
        ├── Competencies.vue       # BARU
        ├── EmployeeMovements.vue
        ├── Attendance.vue         # BARU
        ├── Approvals.vue          # BARU
        ├── Payroll.vue            # BARU
        ├── Leave.vue              # BARU
        ├── Performance.vue        # BARU
        ├── Recruitment.vue        # BARU
        ├── Reimbursements.vue     # BARU
        ├── Training.vue           # BARU
        ├── WorkforceIntelligence.vue
        ├── CareerIntelligence.vue
        └── Packages.vue           # BARU (Subscription)
```

---

## 5.7. Backend API Reference per Module

### Platform Admin API

| Module | Endpoints | Auth | RBAC |
|--------|:---------:|:----:|:----:|
| Auth | 2 | 🔓 Public | - |
| Users | 4 | ✅ JWT | ✅ |
| Companies | 10 | ✅ JWT | ✅ |
| Modules | 7 | ✅ JWT | ✅ |
| Licenses | 4 | ✅ JWT | ✅ |
| **Packages** | **9** | ✅ JWT | ✅ |
| Monitoring | 3 | ✅ JWT | ✅ |
| RBAC | 9 | ✅ JWT | ✅ |
| **Total** | **48** | | |

### Tenant API

| Modul | Endpoints |
|-------|:---------:|
| Organization | 12 |
| Employee | 29 |
| Job Management | 88 |
| Competency | 36 |
| Employee Movement | 15 |
| Attendance | 30 |
| Approval | 15 |
| Payroll | 47 |
| Leave | 23 |
| Performance | 34 |
| Recruitment | 33 |
| Reimbursement | 15 |
| Training | 35 |
| Workforce Intelligence | 68 |
| Career Intelligence | 19 |
| Package Subscription | 3 |
| **Total** | **~502** |

---

## 5.8. Catatan Teknis Implementasi

### PrimeVue 4 Components yang Paling Sering Digunakan
- `DataTable` — tabel dengan sorting, filter, pagination
- `Dialog` — modal dialog (create/edit)
- `Form` — form layout
- `InputText` / `InputNumber` — text & number input
- `Select` — dropdown
- `MultiSelect` — multi-select (module picker)
- `Button` — dengan icon + loading state
- `Tag` — status badge
- `Toast` — notifikasi
- `ConfirmDialog` — konfirmasi aksi
- `TabView` / `Tabs` — tab layout
- `Panel` / `Card` — container
- `Skeleton` — loading state
- `TreeTable` — tree data (organization)
- `Chart` — grafik (dashboard)

### Tailwind CSS Utility yang Sering Digunakan
- `space-y-*` / `space-x-*` — spacing antar children
- `grid grid-cols-*` — grid layout
- `flex items-center justify-between` — flex layout
- `p-*` / `m-*` — padding/margin
- `bg-*` / `text-*` / `border-*` — colors
- `rounded-*` — border radius
- `shadow-*` — shadow
- `hover:*` — hover state
- `transition-*` — animasi
- `animate-pulse` — skeleton loading

### State Management Pattern
| Store | Scope | Status |
|-------|-------|:------:|
| `auth.js` (Pinia) | Platform Admin | ✅ Existing |
| `companies.js` (Pinia) | Platform Admin | ⬜ Perlu dibuat |
| `packages.js` (Pinia) | Platform Admin | ⬜ Perlu dibuat |
| Local `ref` / `reactive` | Per component | ✅ Standar |
| Route query (`?page=1&search=...`) | Filter persistence | ✅ Standar |

### Error Handling Pattern

#### General Error
```js
try {
  const res = await api.get('/endpoint')
} catch (err) {
  toast.add({ severity: 'error', summary: 'Error', detail: err.message })
}
```

#### Field-Level Validation Error ✅ (Implemented)

Backend mengembalikan validation error dengan format:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "fields": {
      "name": ["This field is required"],
      "email": ["Must be a valid email address", "Already exists"]
    },
    "message": "Validation failed"
  }
}
```

Frontend menangani dengan:
```js
import { getValidationErrors } from '@/services/responseHandler'

const errors = ref({})

try {
  await api.post('/endpoint', form.value)
} catch (e) {
  errors.value = getValidationErrors(e)
  if (Object.keys(errors.value).length === 0) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.error?.message })
  }
}
```

**`getValidationErrors()`** (`services/responseHandler.js`):
- Mendukung key `errors` (format lama) dan `fields` (format baru)
- Array values otomatis di-implode dengan `', '` (koma + spasi)
  - `["Required"]` → `"Required"`
  - `["Required", "Already exists"]` → `"Required, Already exists"`
- Output: `{ name: "Required", email: "Required, Already exists" }` — langsung cocok untuk prop `:errors` di `FormRow`

**FormRow component** menampilkan error text otomatis:
```vue
<FormRow label="Email" :errors="errors?.email">
  <TextInput v-model="form.email" :class="{ 'p-invalid': errors?.email }" />
</FormRow>
```

### Integrasi dengan Backend Bilingual Response

Backend telah mendukung **bilingual penuh** (English + Bahasa Indonesia) via:
- Middleware `Localize()` — auto-detect dari header `Accept-Language`
- Response helpers — `SuccessJSON`, `CreatedJSON`, `NotFound`, `ErrorJSON`, dll
- Message catalog — 80+ pasangan EN/ID di `internal/pkg/httputil/locale.go`
- Validasi error — field-specific errors dengan pesan bilingual

Semua API endpoint mengembalikan response dengan format:
```json
{
  "success": true,
  "data": { ... },
  "message": "Created successfully"   ← otomatis dalam bahasa yang diminta!
}
```

### Implementasi Frontend untuk Bilingual

#### 1. Language Store (Pinia)
```js
// stores/language.js
export const useLanguageStore = defineStore('language', () => {
  const lang = ref(localStorage.getItem('hris_lang') || 'en')
  function setLang(l) { lang.value = l; localStorage.setItem('hris_lang', l) }
  function toggleLang() { setLang(lang.value === 'en' ? 'id' : 'en') }
  return { lang, setLang, toggleLang }
})
```

#### 2. Axios Interceptor — Kirim Header Tiap Request
```js
api.interceptors.request.use((config) => {
  const langStore = useLanguageStore()
  config.headers['Accept-Language'] = langStore.lang
  return config
})
```

#### 3. Language Switcher di HeaderBar
- Tombol toggle EN/ID di pojok kanan
- Icon `pi pi-globe` + label "EN" / "ID"
- State dari language store

```vue
<Button @click="langStore.toggleLang()" class="p-button-text">
  <i class="pi pi-globe mr-1"></i>
  {{ langStore.lang === 'en' ? 'EN' : 'ID' }}
</Button>
```

#### 4. Composable `useI18n` untuk Static UI Text
```js
// composables/useI18n.js
import en from '@/locales/en.json'
import id from '@/locales/id.json'
const messages = { en, id }
export function useI18n() {
  const langStore = useLanguageStore()
  function t(key, params) {
    const text = (messages[langStore.lang]?.[key] || messages.en[key] || key)
    return params ? text.replace(/%s/g, params) : text
  }
  return { t }
}
```

#### 5. Locale Files
- `src/locales/en.json` — English static text (sidebar menu, button labels, table headers)
- `src/locales/id.json` — Indonesian static text

#### 6. Response Handler
```js
export function getMessage(response, lang = 'en') {
  if (!response?.message) return ''
  if (typeof response.message === 'string') return response.message
  return response.message[lang] || response.message.en || ''
}
```

### Alur Lengkap
```text
[User klik ID] → languageStore.lang = 'id'
    → Axios header Accept-Language: id
    → Backend Localize() middleware detek bahasa
    → Backend response helper panggil tCtx(c, 'success.created')
    → localeMessages['success.created']['id'] = "Berhasil dibuat"
    → Response JSON: { "message": "Berhasil dibuat" }
    → Toast/banner tampilkan: "Berhasil dibuat"
```

### Struktur File Baru
```
frontend/*/src/
├── composables/useI18n.js          # t(key) untuk static text
├── composables/useNotification.js  # toast dengan bilingual message
├── locales/en.json                 # English static text
├── locales/id.json                 # Indonesian static text
├── services/responseHandler.js     # parse bilingual response/error
├── stores/language.js              # Pinia store
└── layouts/HeaderBar.vue           # + language switcher button
```

> **Detail lengkap:** Lihat `docs/frontend-development-plan.md` Section I — Bilingual Support.

---

*Dokumen ini akan diupdate seiring progres implementasi frontend.*
