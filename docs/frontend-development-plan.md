# Frontend Development Plan — HRIS Platform

**Generated:** 26 July 2026
**Tech Stack:** Vue 3 + PrimeVue 4 + Tailwind CSS 4 + Vite + Axios

---

## A. Ringkasan Arsitektur Frontend

```
frontend/
├── platform-admin/          # Admin panel (Vue 3 + PrimeVue)
│   ├── src/
│   │   ├── App.vue          # Root - Toast + ConfirmDialog
│   │   ├── main.js          # Entry point
│   │   ├── router/index.js  # Routes + Auth guard
│   │   ├── layouts/         # AppLayout, HeaderBar, Sidebar
│   │   ├── views/           # Dashboard, Companies, Users, etc.
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

### Shared Stack
- **Framework:** Vue 3 (Composition API + `<script setup>`)
- **UI Library:** PrimeVue 4 (DataTable, Dialog, Tag, Button, InputText, Select, etc.)
- **Styling:** Tailwind CSS 4 + PrimeVue themes (Aura)
- **HTTP:** Axios (interceptors for JWT)
- **State:** Pinia (auth store)
- **Router:** Vue Router 4 (auth guard, lazy loading)
- **Icons:** PrimeIcons

---

## B. Phase 1: Platform Admin — Fundamental (MVP)

### B.1. Login Page ✅ (Existing)
- [x] Form login (email + password)
- [x] JWT token management (access + refresh)
- [x] Redirect ke dashboard setelah login
- [x] Auto-redirect ke login jika token expired

### B.2. Dashboard Page ✅ (Existing - Perlu Penyempurnaan)
- [x] KPI Cards (Total Companies, Active Tenants, Users, Modules, Connections, Health)
- [x] Recent Companies list
- [x] System Health (DB, Cache, Pool Stats)
- [x] Quick Actions
- [x] **Tambahkan:** Grafik/tren jumlah company per bulan (PrimeVue Chart + Chart.js bar chart, 6 months)
- [x] **Tambahkan:** Pool wait count / system load indicator dari monitoring/health endpoint
- [x] **Tambahkan:** Real-time updates via polling (refresh otomatis tiap 30s, cleanup onUnmounted)

### B.3. Companies Page ✅ (Existing - Enhanced)
- [x] DataTable daftar companies (pagination)
- [x] Create company dialog (name, slug, package selection)
- [x] Company detail / edit
- [x] Status management (Activate, Suspend, Terminate)
- [x] Status badges (Tag component)
- [x] **Filter by status** — quick filter chips: All / Active / Suspended / Terminated
- [x] **Filter by package** — Select dropdown filter dari daftar packages API
- [x] **Search by name/slug/phone/address** — client-side filter
- [x] **License info inline** — kolom plan_type dengan severity badge + edit dialog
- [x] **Tenant provisioning progress** — Backend: ProvisioningInfo DTO + repository methods + service update; Frontend: DB column dengan Tag status (Provisioned/Deactivated/Not Provisioned) + tooltip bilingual

### B.4. Users Page ✅ (Existing - Enhanced)
- [x] DataTable daftar platform users
- [x] Create / Edit user dialog
- [x] **Tambahkan:** Filter chips by role (All / Super Admin / Company Admin)
- [x] **Tambahkan:** Search by name/email (IconField + client-side filter)
- [x] **Tambahkan:** Dialog labels bilingual (t('users.name'), t('users.email'), t('common.password'))
- [x] Bulk actions — selection checkboxes, bulk action toolbar (Change Role + Delete), bulk confirmation dialog, super_admin protection
- [x] Backend: Added DELETE /api/v1/platform/users/:id endpoint (soft delete, super_admin cannot be deleted)

### B.5. Modules Page ✅ (Existing - Enhanced)
- [x] DataTable daftar modules
- [x] Create module dialog
- [x] Module detail
- [x] Activate/deactivate untuk company
- [x] Status badges
- [x] **Tambahkan:** Filter `?module_type=platform|tenant` — server-side filter via API query param
- [x] **Tambahkan:** Kolom `module_type` di tabel (badge: "Platform" / "Tenant")
- [x] **Tambahkan:** Kolom `depends_on` dengan tooltip untuk nilai panjang
- [x] **Tambahkan:** Quick filter chips: "All" | "Platform" | "Tenant"
- [x] **Tambahkan:** Client-side search by name/slug/description

### B.6. Licenses Page ✅ (Existing - Enhanced)
- [x] DataTable daftar licenses
- [x] Create license dialog
- [x] License detail
- [x] **Tambahkan:** Kolom `package_id` + `package_name` (integrasi package)
- [x] **Tambahkan:** Filter by package (Select dropdown dari API packages)
- [x] **Tambahkan:** Status filter chips (All / Active / Expired / Suspended)
- [x] **Tambahkan:** License key column + copy button (clipboard API)
- [x] **Tambahkan:** Expiration date warning (Tag: "Expiring Soon" jika ≤30 hari, "Expired" jika sudah lewat)

### B.7. Monitoring Page ✅ (Existing)
- [x] Platform health status
- [x] Database connectivity per tenant
- [x] Pool stats
- [x] **Tambahkan:** Auto-refresh toggle (ToggleSwitch + 30s polling + live indicator)
- [x] **Tambahkan:** Grafis pool utilization over time (Line chart: Open/In Use/Idle, rolling 20 samples)
- [x] **Tambahkan:** Alert thresholds (wait_count > 0, utilization >50%/>80%, unhealthy tenants, cache)

### B.8. Packages Page ✅ (BARU - Done)
- [x] DataTable daftar packages (CRUD)
  - Kolom: Name+Slug, Price (IDR format), Status (Tag), Module Count, Sort Order
  - Search by name/slug/description
  - Tooltip deskriptif pada action buttons (Edit, Publish, Unpublish, Validate, Delete)
- [x] Create Package Dialog
  - Form: Name, Slug (auto-generate dari name + highlight animasi), Description (TextArea), Price (InputNumber), Sort Order
  - Module selector dua kolom: daftar modul dengan ToggleSwitch, expand detail (deskripsi + depends_on)
  - Select All / Deselect All toggle untuk modul
  - Per-module: isMandatory toggle + sort_order input
- [x] Edit Package Dialog (pre-filled, module mandatory state preserved)
- [x] Publish / Unpublish button (dengan konfirmasi bilingual)
  - Backend validasi dependensi otomatis sebelum publish
- [x] Delete Package (soft delete, dengan konfirmasi bilingual)
- [x] Validate Dependencies button
  - Modal daftar dependensi per module: green (resolved) / red (unresolved)
- [x] Status badge: Draft (info/blue), Published (success/green)
- [x] Field-level validation error display (via getValidationErrors)
- [x] Bilingual support (EN/ID locale keys)
- [x] `useSlugify` composable (shared slugify logic antara Packages & Modules)
- [x] Slug animasi CSS global (slug-animation.css)

### B.9. RBAC Management Page 🔴 (BARU - Done)
- [x] Daftar roles (DataTable)
- [x] Create role dialog (name, slug auto-generate dari useSlugify, description)
- [x] Role detail page (permissions matrix — grouped by module/resource)
- [x] Permission assignment (ToggleSwitch + grouped by resource card layout + Select All per-group)
- [x] Delete role (non-system only, dengan konfirmasi bilingual)
- [x] SYSTEM badge untuk role sistem (tidak bisa dihapus)
- [x] Permission API endpoints: `/api/v1/platform/rbac/roles`, `/rbac/permissions`, assign/revoke
- [x] Auto-slug form (useSlugify composable, highlight animasi, sync icon)
- [x] Permission modal title fix (remove duplicate rbac section di locale files)
- [x] Module description enhancement — permissions grouped by resource dengan module header

### B.10. Profile Page ✅ (BARU - Done)
- [x] **User info display** — Card dengan Avatar, name, email, role (Tag), company, status (Tag), last login
- [x] **Change password form** — current password, new password, confirm new password
- [x] Client-side validation (required, min 6 chars, password mismatch check)
- [x] Server-side validation via API error per-field (`getValidationErrors`)
- [x] Backend: Added `ChangePassword` service method with bcrypt verify + hash
- [x] Backend: Added `PUT /users/:id/password` endpoint
- [x] Backend: Added `password.updated` locale key (EN/ID)
- [x] **User data persistence** — Auth store now saves user to localStorage (survives page refresh)
- [x] **HeaderBar** — Profile menu item wired to `/profile` route (was empty placeholder)
- [x] Route registered: `/profile` with auth guard
- [x] Bilingual (EN/ID): 12 locale keys in `profile.*` section
- [x] Date formatting respects language setting (`id-ID` vs `en-US`)

---

## C. Phase 2: Tenant — Module Views

### C.1. Layout & Navigation ✅ (Existing)
- [x] Sidebar with all module links
- [x] HeaderBar with user info
- [x] Responsive sidebar (collapsible)
- [x] 15 module routes registered
- [x] Dashboard with KPI cards + quick access

### C.2. Dashboard ✅ (Existing - Mock Data)
- [x] KPI Cards (Total Employees, Active Today, On Leave, Pending Approvals)
- [x] Quick Access Modules grid (12 modules)
- [x] Recent Activity (static)
- [x] Period filter (This Month/Quarter/Year)
- [ ] **Perbaiki:** Ganti mock data dengan real API calls
  - GET /api/v1/tenant/employees?per_page=1 (total count)
  - GET /api/v1/tenant/attendance/events (active today)
  - GET /api/v1/tenant/leave/requests (on leave count)
  - GET /api/v1/tenant/approval/tasks (pending approval count)

### C.3. Organization Management 🔴 (BARU)
**Backend:** 12 endpoints
- [ ] DataTable organizations (tree view dengan PrimeVue TreeTable atau OrganizationChart)
- [ ] CRUD organization (inline editing atau dialog)
- [ ] Organization Tree view (hierarki parent-child)
- [ ] Zones CRUD
- [ ] Positions CRUD
- [ ] Job Families CRUD

### C.4. Employee Management 🔴 (BARU)
**Backend:** 29 endpoints
- [ ] DataTable employees (with search, filter by department/status)
- [ ] Create Employee Wizard (multi-step form)
  - Step 1: Personal Data
  - Step 2: Address (MAIN/DOMICILE)
  - Step 3: Emergency Contact
  - Step 4: Family
  - Step 5: Education
  - Step 6: Work Experience
  - Step 7: Documents
  - Step 8: Insurance (BPJS)
  - Step 9: Employment Record
- [ ] Employee Detail Page (tab view per sub-module)
  - Tab: Profile, Addresses, Contacts, Family, Education, Experience, Documents, Insurance, Employment
- [ ] Edit per sub-module (inline dialog)
- [ ] Delete with confirmation

### C.5. Job Management 🔴 (BARU)
**Backend:** 88 endpoints (18 sub-entities)
- [ ] Tab layout per sub-entity
- [ ] Job Titles CRUD + Title Subs
- [ ] Job Values CRUD
- [ ] Job Objectives CRUD
- [ ] Job Identifications CRUD
- [ ] Responsibilities CRUD
- [ ] Working Activities CRUD
- [ ] Working Risks CRUD
- [ ] Job Scores (by organization)
- [ ] Competency Groups CRUD

### C.6. Competency Management 🔴 (BARU)
**Backend:** 35 endpoints (7 entities)
- [ ] Competencies CRUD
- [ ] Competency Values CRUD
- [ ] Competency Events + Targets CRUD
- [ ] Competency Scores + Details CRUD

### C.7. Employee Movement 🔴 (BARU)
**Backend:** 15 endpoints
- [ ] Movements list (DataTable with status workflow)
- [ ] Create Movement (type: promotion/demotion/mutation, etc.)
- [ ] Approve / Execute / Cancel workflow buttons
- [ ] Contracts CRUD (PKWT/PKWTT)

### C.8. Time & Attendance 🔴 (BARU)
**Backend:** 30 endpoints (10 entities)
- [ ] Company Settings
- [ ] Shifts CRUD
- [ ] Employee Shifts assignment
- [ ] Attendance Events (check-in/out log)
- [ ] Daily Sessions view
- [ ] Overtime Requests CRUD
- [ ] Location management (geofence)
- [ ] Exempt Positions

### C.9. Approval Engine 🔴 (BARU)
**Backend:** 15 endpoints
- [ ] Approval Flows CRUD (with multi-step)
- [ ] My Tasks (pending approvals list)
- [ ] Approval Instances history
- [ ] Approve / Reject actions

### C.10. Payroll & Compensation 🔴 (BARU)
**Backend:** 47 endpoints (21 entities)
- [ ] Salary Components CRUD
- [ ] Payroll Periods CRUD
- [ ] Payroll Runs (with status workflow)
- [ ] Employee Payroll Profiles
- [ ] BPJS Settings
- [ ] PPh21 Settings
- [ ] Tax Brackets
- [ ] Payslip view

### C.11. Leave & Time Off 🔴 (BARU)
**Backend:** 23 endpoints
- [ ] Leave Types CRUD
- [ ] Accrual Policies CRUD
- [ ] Leave Requests (create, approve, reject)
- [ ] Leave Balances per employee

### C.12. Performance Management 🔴 (BARU)
**Backend:** 34 endpoints
- [ ] Performance Periods CRUD
- [ ] Perspectives CRUD
- [ ] Templates CRUD
- [ ] Evaluations (with details)
- [ ] KPI/OKR Targets

### C.13. Recruitment & Onboarding (ATS) 🔴 (BARU)
**Backend:** 33 endpoints
- [ ] Job Requisitions CRUD
- [ ] Candidates CRUD
- [ ] Applications (status pipeline)
- [ ] Interviews scheduling
- [ ] Onboarding Tasks

### C.14. Reimbursement & Claim 🔴 (BARU)
**Backend:** 15 endpoints
- [ ] Reimbursement Types CRUD
- [ ] Requests (DRAFT → SUBMITTED → APPROVED → PAID)
- [ ] Items per request
- [ ] Status workflow buttons

### C.15. Training & Development 🔴 (BARU)
**Backend:** 35 endpoints
- [ ] Categories CRUD
- [ ] Courses CRUD
- [ ] Sessions CRUD
- [ ] Participants registration
- [ ] Evaluations
- [ ] Certificates

### C.16. Workforce Intelligence 🔴 (BARU)
**Backend:** 68 endpoints (analytics layer)
- [ ] Dashboard (KPI summaries)
- [ ] Headcount Planning
- [ ] Analytics dashboards (headcount, attendance, leave, payroll, etc.)
- [ ] Risk Dashboard
- [ ] Executive Dashboard
- [ ] Scenario Planning
- [ ] Organization Health metrics
- [ ] People Analytics charts

### C.17. Career Intelligence 🔴 (BARU)
**Backend:** 19 endpoints
- [ ] 9-Box Talent Grid visualization
- [ ] Talent Maps CRUD
- [ ] Career Interests
- [ ] Career Paths
- [ ] Gap Analysis view
- [ ] Succession Plans CRUD

### C.18. Package Subscription (Tenant) 🔴 (BARU)
- [ ] Browse published packages (GET /api/v1/public/packages)
  - Card layout: nama, deskripsi, harga, daftar module
  - Filter: `?module_type=tenant`
- [ ] Current subscription status
- [ ] Subscribe button → POST /api/v1/tenant/packages/:id/subscribe
- [ ] Unsubscribe button → POST /api/v1/tenant/packages/:id/unsubscribe
- [ ] Activated modules list (dari response subscribe)

---

## D. Shared Components yang Perlu Dibuat

### D.1. Reusable Components
- [ ] **DataTableWrapper** — wrapper PrimeVue DataTable dengan:
  - Sorting, filtering, pagination built-in
  - Search input
  - Export CSV
  - Column visibility toggle
- [ ] **ConfirmDialog** — reusable confirmation dengan PrimeVue ConfirmDialog
- [ ] **StatusBadge** — Tag component dengan severity mapping otomatis
- [ ] **FormDialog** — reusable form dialog (create/edit)
- [ ] **SearchInput** — input dengan icon search + debounce
- [ ] **FilterChips** — chip group untuk quick filter (All | Active | Inactive, dll)
- [ ] **LoadingSkeleton** — skeleton loading (PrimeVue Skeleton)
- [ ] **EmptyState** — komponen untuk empty table/message
- [ ] **ApiError** — error display component
- [ ] **PageHeader** — reusable page header (title, subtitle, actions)

### D.2. Shared Services
- [ ] **api.js** — Axios instance with interceptors ✅ (Existing)
  - Auto-attach JWT token
  - Auto-refresh on 401
  - Error toast notification
- [ ] **NotificationService** — wrapper toast untuk success/error/warning
- [ ] **FormattingUtils** — date, currency, number formatting
- [ ] **ValidationUtils** — form validation helpers (NIK, NPWP, dll)

---

## I. Bilingual Support (Dua Bahasa) di Frontend

### I.1. Arsitektur Bilingual Frontend

Backend sudah support bilingual (EN/ID) via middleware `Localize()` dan response helpers (`SuccessJSON`, `NotFound`, dll). Frontend perlu mengintegrasikan:

```
┌─────────────────────────────────────────────────────────┐
│  Frontend Bilingual Architecture                         │
│                                                          │
│  ┌───────────┐  ┌───────────────┐  ┌──────────────────┐  │
│  │ Language   │  │ Pinia Store   │  │ Axios Interceptor │  │
│  │ Switcher   │──▶ (language.js) │──▶ Accept-Language  │  │
│  │ (HeaderBar)│  │ lang: en/id   │  │ header otomatis  │  │
│  └───────────┘  └───────┬───────┘  └──────────────────┘  │
│                          │                                │
│                          ▼                                │
│  ┌───────────────────────────────────────────────────┐   │
│  │  Response Handler                                  │   │
│  │  • Toast/SnackBar → baca message.en atau message.id│   │
│  │  • Validation error → tampil sesuai bahasa aktif   │   │
│  │  • DataTable header → label dari locale map        │   │
│  └───────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────────────┐  ┌──────────────────────────────┐  │
│  │ useI18n (composable)│  │ localeMessages.json         │  │
│  │ t(key) → string    │  │ { "dashboard.title": {      │  │
│  │ t(key, 'id')       │  │   "en": "Dashboard",        │  │
│  └──────────────────┘  │   "id": "Dasbor"              │  │
│                         │ } }                            │  │
│                         └──────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### I.2. Langkah Implementasi

#### Langkah 1: Language Store (Pinia)
- **File:** `src/stores/language.js`
- State: `lang` ("en" | "id")
- Actions: `setLang(lang)`, `toggleLang()`
- Getters: `isID`, `isEN`
- Persist: localStorage key `hris_lang`
- Inisialisasi: cek localStorage → `Accept-Language` browser → default "en"

```js
// stores/language.js
export const useLanguageStore = defineStore('language', () => {
  const lang = ref(localStorage.getItem('hris_lang') || navigator.language?.startsWith('id') ? 'id' : 'en')
  
  function setLang(l) {
    lang.value = l
    localStorage.setItem('hris_lang', l)
  }
  function toggleLang() {
    setLang(lang.value === 'en' ? 'id' : 'en')
  }
  
  return { lang, setLang, toggleLang }
})
```

#### Langkah 2: Axios Interceptor — Auto Header Bahasa
- **File:** `src/services/api.js`
- Saat setiap request, attach header `Accept-Language` dari language store
- Saat response error, parse `error.errors` untuk menampilkan field validation errors dalam bahasa yang sesuai

```js
// api.js — request interceptor
api.interceptors.request.use((config) => {
  const langStore = useLanguageStore()
  config.headers['Accept-Language'] = langStore.lang
  return config
})
```

#### Langkah 3: Language Switcher di HeaderBar
- **File:** `src/layouts/HeaderBar.vue` (kedua apps)
- Tombol toggle EN/ID di pojok kanan HeaderBar
- Icon: globe + flag indicator
- Dropdown atau toggle button dengan transisi halus
- State dari language store

```vue
<!-- HeaderBar.vue -->
<Button 
  v-tooltip="langStore.lang === 'en' ? 'Bahasa Indonesia' : 'English'"
  @click="langStore.toggleLang()"
  class="p-button-text"
>
  <i class="pi pi-globe mr-1"></i>
  {{ langStore.lang === 'en' ? 'EN' : 'ID' }}
</Button>
```

#### Langkah 4: Composable `useI18n` untuk Static UI Text
- **File:** `src/composables/useI18n.js`
- Load file `src/locales/en.json` dan `src/locales/id.json`
- Fungsi `t(key, params?)` — ambil teks sesuai bahasa aktif
- Fallback ke English jika key tidak ditemukan di ID

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
  
  return { t, lang: computed(() => langStore.lang) }
}
```

#### Langkah 5: Locale Message Files
- **File:** `src/locales/en.json`
- **File:** `src/locales/id.json`
- Cakupan: sidebar menu, button labels, table headers, form labels, status texts, notifikasi, halaman login

```json
// locales/en.json
{
  "nav.dashboard": "Dashboard",
  "nav.companies": "Companies",
  "nav.users": "Users",
  "nav.modules": "Modules",
  "nav.licenses": "Licenses",
  "nav.packages": "Packages",
  "nav.monitoring": "Monitoring",
  "common.search": "Search...",
  "common.filter": "Filter",
  "common.create": "Create",
  "common.edit": "Edit",
  "common.delete": "Delete",
  "common.save": "Save",
  "common.cancel": "Cancel",
  "common.confirm": "Confirm",
  "common.status.active": "Active",
  "common.status.inactive": "Inactive",
  "common.status.draft": "Draft",
  "common.status.published": "Published",
  "common.loading": "Loading...",
  "common.no_data": "No data available",
  "auth.login.title": "Sign In",
  "auth.login.email": "Email",
  "auth.login.password": "Password",
  "auth.login.button": "Sign In",
  "auth.logout": "Sign Out",
  "message.success": "Success",
  "message.error": "Error",
  "message.warning": "Warning",
  "message.info": "Info"
}
```

```json
// locales/id.json
{
  "nav.dashboard": "Dasbor",
  "nav.companies": "Perusahaan",
  "nav.users": "Pengguna",
  "nav.modules": "Modul",
  "nav.licenses": "Lisensi",
  "nav.packages": "Paket",
  "nav.monitoring": "Pemantauan",
  "common.search": "Cari...",
  "common.filter": "Filter",
  "common.create": "Buat",
  "common.edit": "Ubah",
  "common.delete": "Hapus",
  "common.save": "Simpan",
  "common.cancel": "Batal",
  "common.confirm": "Konfirmasi",
  "common.status.active": "Aktif",
  "common.status.inactive": "Nonaktif",
  "common.status.draft": "Draf",
  "common.status.published": "Dipublikasikan",
  "common.loading": "Memuat...",
  "common.no_data": "Tidak ada data",
  "auth.login.title": "Masuk",
  "auth.login.email": "Email",
  "auth.login.password": "Kata Sandi",
  "auth.login.button": "Masuk",
  "auth.logout": "Keluar",
  "message.success": "Berhasil",
  "message.error": "Gagal",
  "message.warning": "Peringatan",
  "message.info": "Informasi"
}
```

#### Langkah 6: Response Handler — Parse Bilingual Message
- **File:** `src/services/responseHandler.js`
- Fungsi `parseMessage(response, lang)` — baca `data.message.en` atau `data.message.id`
- Fungsi `parseError(error, lang)` — baca `error.error.message.en/id`, tampilkan field errors
- Integrasi dengan toast/notification service

```js
// services/responseHandler.js
export function getMessage(response, lang = 'en') {
  if (!response?.message) return ''
  if (typeof response.message === 'string') return response.message
  return response.message[lang] || response.message.en || ''
}

export function getErrorMessage(error, lang = 'en') {
  if (!error?.error) return 'Unknown error'
  const err = error.error
  if (typeof err.message === 'string') return err.message
  return err.message?.[lang] || err.message?.en || err.code || 'Unknown error'
}

export function getValidationErrors(error, lang = 'en') {
  // error.errors berisi map field → array of strings dalam bahasa yang diminta
  // Backend sudah mengembalikan dalam bahasa yang sesuai dengan Accept-Language
  return error?.error?.errors || {}
}
```

#### Langkah 7: Notifikasi Toast Terintegrasi Bahasa
- **File:** `src/composables/useNotification.js`
- Wrapper untuk PrimeVue Toast
- Auto-read message dari response API dalam bahasa yang sesuai

```js
// composables/useNotification.js
export function useNotification() {
  const toast = useToast()
  const langStore = useLanguageStore()
  
  function showSuccess(response) {
    toast.add({
      severity: 'success',
      summary: langStore.lang === 'id' ? 'Berhasil' : 'Success',
      detail: getMessage(response, langStore.lang),
      life: 3000
    })
  }
  
  return { showSuccess, showError, showWarning }
}
```

### I.3. Struktur File Baru

```
frontend/platform-admin/src/
├── composables/
│   ├── useI18n.js          # BARU — t(key) untuk static text
│   └── useNotification.js  # BARU — toast dengan bilingual message
├── locales/
│   ├── en.json             # BARU — English static text
│   └── id.json             # BARU — Indonesian static text
├── services/
│   ├── api.js              # UPDATE — add Accept-Language interceptor
│   └── responseHandler.js  # BARU — parse bilingual response/error
├── stores/
│   ├── auth.js             # Existing
│   └── language.js         # BARU — Pinia store for language state
└── layouts/
    └── HeaderBar.vue       # UPDATE — add language switcher button

frontend/tenant/src/
├── composables/
│   ├── useI18n.js          # BARU (copy atau shared)
│   └── useNotification.js  # BARU
├── locales/
│   ├── en.json
│   └── id.json
├── services/
│   ├── api.js              # UPDATE
│   └── responseHandler.js  # BARU
├── stores/
│   └── language.js         # BARU
└── layouts/
    └── HeaderBar.vue       # UPDATE — add language switcher
```

### I.4. Prioritas Implementasi Bilingual

| Langkah | Task | Kompleksitas | Status |
|:-------:|------|:------------:|:------:|
| 1 | Language Store (Pinia) + localStorage persistence | 🟢 Easy | ⬜ TODO |
| 2 | Axios interceptor `Accept-Language` header | 🟢 Easy | ⬜ TODO |
| 3 | Language Switcher button di HeaderBar | 🟢 Easy | ⬜ TODO |
| 4 | Composable `useI18n` + locale files EN/ID | 🟡 Medium | ⬜ TODO |
| 5 | Response handler + bilingual toast notification | 🟡 Medium | ⬜ TODO |
| 6 | Integrasi bilingual di semua view (label, header, filter) | 🟡 Medium | ⬜ TODO |

### I.5. Alur Lengkap Bilingual Request → Response

```text
[User pilih bahasa ID]
    ↓
languageStore.lang = 'id'
    ↓
Axios request → header Accept-Language: id
    ↓
Backend Localize() middleware detect lang = ID
    ↓
Backend response helpers (SuccessJSON, ErrorJSON, NotFound, dll)
    ↓ panggil tCtx(c, 'success.created')
    ↓ localeMessages['success.created']['id'] = "Berhasil dibuat"
    ↓
Response JSON:
{
  "success": true,
  "data": { ... },
  "message": "Berhasil dibuat"   ← otomatis bahasa Indonesia!
}
    ↓
Response handler → toast.show("Berhasil dibuat")
```

> **Catatan:** Backend sudah fully support bilingual (80+ message pairs, EN/ID) di `internal/pkg/httputil/locale.go`. Frontend hanya perlu:
> 1. Kirim header `Accept-Language` yang benar
> 2. Parse response `message` field (yang sudah dalam bahasa yang diminta)
> 3. Static UI text via locale files untuk label/title yang tidak dari API

---

## E. Prioritas Eksekusi

### Phase 1 — Platform Admin (MVP) — Estimasi: 2-3 minggu
| Priority | Feature | Kompleksitas | Status |
|:--------:|---------|:------------:|:------:|
| P0 | Packages Page (CRUD + Publish) | 🟡 Medium | ✅ Done |
| P0 | Package dependency validation UI | 🟡 Medium | ✅ Done |
| P0 | Modules filter by `module_type` | 🟢 Easy | ✅ Done |
| P1 | Companies filter by status + package | 🟢 Easy | ✅ Done |
| P1 | Licenses package integration | 🟢 Easy | ✅ Done |
| P1 | Dashboard charts & real-time | 🟡 Medium | ✅ Done |
| P2 | RBAC Management | 🔴 Complex | ✅ Done |
| P2 | Profile Page + Change Password | 🟢 Easy | ✅ Done |

### Phase 2 — Tenant Core Modules — Estimasi: 4-6 minggu
| Priority | Feature | Kompleksitas |
|:--------:|---------|:------------:|
| P0 | Organization Management | 🟡 Medium |
| P0 | Employee Management (Wizard) | 🔴 Complex |
| P1 | Leave & Attendance | 🟡 Medium |
| P1 | Payroll (read-only payslip) | 🟡 Medium |
| P2 | Job Management | 🔴 Complex |
| P2 | Competency Management | 🔴 Complex |

### Phase 3 — Tenant Advanced Modules — Estimasi: 4-6 minggu
| Priority | Feature | Kompleksitas |
|:--------:|---------|:------------:|
| P1 | Performance Management | 🔴 Complex |
| P1 | Recruitment (ATS Pipeline) | 🔴 Complex |
| P2 | Approval Engine (Flow Builder) | 🔴 Complex |
| P2 | Employee Movement (Workflow) | 🟡 Medium |
| P2 | Training Management | 🟡 Medium |
| P2 | Reimbursement | 🟡 Medium |

### Phase 4 — Intelligence & Subscription — Estimasi: 2-3 minggu
| Priority | Feature | Kompleksitas |
|:--------:|---------|:------------:|
| P1 | Workforce Intelligence Dashboards | 🔴 Complex |
| P1 | Career Intelligence (9-box Grid) | 🔴 Complex |
| P1 | Package Subscription (Tenant) | 🟡 Medium |

---

## F. Struktur File yang Diharapkan (Final)

```
frontend/platform-admin/src/
├── components/
│   ├── form/           # FormRow, InputLabel, TextInput, SelectLabel, etc.
│   ├── layout/         # AppLayout, Sidebar, HeaderBar
│   └── shared/         # DataTableWrapper, ConfirmDialog, StatusBadge, etc.
├── composables/        # useApi, useNotification, usePagination, etc.
├── router/index.js
├── services/api.js
├── stores/             # auth.js, companies.js, packages.js, etc.
├── utils/              # formatters, validators (NIK, NPWP, etc.)
└── views/
    ├── Dashboard.vue
    ├── Companies.vue
    ├── Users.vue
    ├── Modules.vue
    ├── Licenses.vue
    ├── Packages.vue         # BARU
    ├── Monitoring.vue
    ├── Profile.vue          # BARU
    └── Rbac.vue             # BARU

frontend/tenant/src/
├── components/
│   ├── layout/         # AppLayout, Sidebar, HeaderBar
│   └── shared/         # DataTableWrapper, FormDialog, ConfirmDialog
├── composables/        # useTenantApi, useModuleAccess, etc.
├── router/index.js
├── services/api.js     # tenant-specific axios instance
├── stores/             # auth, subscription
├── utils/              # formatters
└── views/
    ├── Dashboard.vue           # Perlu real API
    └── modules/
        ├── Organizations.vue   # BARU
        ├── Employees.vue       # BARU (Wizard)
        ├── JobManagement.vue   # BARU
        ├── Competencies.vue    # BARU
        ├── EmployeeMovements.vue
        ├── Attendance.vue      # BARU
        ├── Approvals.vue       # BARU
        ├── Payroll.vue         # BARU
        ├── Leave.vue           # BARU
        ├── Performance.vue     # BARU
        ├── Recruitment.vue     # BARU
        ├── Reimbursements.vue  # BARU
        ├── Training.vue        # BARU
        ├── WorkforceIntelligence.vue
        ├── CareerIntelligence.vue
        └── Packages.vue        # BARU (Subscription)
```

---

## G. Backend API Reference per Module

### Platform Admin API Summary
| Module | Endpoints | Auth | RBAC |
|--------|:---------:|:----:|:----:|
| Auth | 2 | 🔓 | - |
| Users | 5 (+ password) | ✅ | ✅ |
| Companies | 10 | ✅ | ✅ |
| Modules | 7 | ✅ | ✅ |
| Licenses | 4 | ✅ | ✅ |
| **Packages** | **9** | ✅ | ✅ |
| Monitoring | 3 | ✅ | ✅ |
| RBAC | 9 | ✅ | ✅ |

### Tenant API Summary
| Module | Endpoints |
|--------|:---------:|
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
| **Package Subscription** | **3** |
| **TOTAL** | **~502** |

---

## H. Catatan Teknis

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

### State Management
- **Pinia** untuk auth store ✅ (Existing)
- **Pinia** untuk companies cache
- **Pinia** untuk packages store
- **Local component state** (`ref` / `reactive`) untuk form data
- **Route query** (`?page=1&search=...`) untuk filter persistence

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
    // Jika tidak ada field-level errors, tampilkan toast biasa
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

**FormRow component** akan menampilkan error text otomatis:
```vue
<FormRow label="Email" :errors="errors?.email">
  <TextInput v-model="form.email" :class="{ 'p-invalid': errors?.email }" />
</FormRow>
```

---

*Dokumen ini akan diupdate seiring progres implementasi frontend.*
