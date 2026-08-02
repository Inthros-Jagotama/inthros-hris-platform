# Frontend Development Plan — HRIS Platform

**Generated:** 27 July 2026
**Last Updated:** 2 August 2026 (Organization Tree drag & drop fix — PrimeVue 4.5.5 API; Employee Form: 8/9 steps completed; Platform Admin: CompanyActions.vue reusable component)
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
- [x] **Company Detail Page** — `views/CompanyDetail.vue` (route `/companies/:id`): header card (nama, slug, status Tag) + 4 section cards (Basic Info, License, Database & Provisioning, Admin User) + tombol navigasi balik; tombol eye di list untuk masuk ke detail
- [x] **Rotate Credentials** — tombol (pi-key) di list & detail → ConfirmDialog → `POST /companies/:id/rotate-credentials` → dialog password auto-generated sekali lihat (readonly InputText + copy + warning amber); tampil hanya jika company provisioned & bukan terminated
- [x] **CompanyActions.vue (reusable)** — `components/CompanyActions.vue` mengelola Edit/Suspend/Activate/Terminate/Rotate + ConfirmDialog + edit dialog + rotate password dialog; dipakai di **list** (`mode="icons"`) & **detail** (`mode="buttons"`); emit `updated` → parent reload; tombol Edit/Suspend/Activate/Terminate disembunyikan untuk company `terminated` (konsisten list & detail)
- [x] **Edit Company** — dialog FormRow/TextInput (name/email/phone/address) + info lisensi saat ini, validasi field-error via `getValidationErrors`, `PUT /companies/:id`

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
- [x] **Company name + last login fetch** — Profile fetches enriched user data from `GET /api/v1/platform/users/:id` on mount (backfill company_name, updated last_login)

### B.11. Dark/Light Mode Toggle ✅ (BARU - Done)
- [x] **Theme store** — `stores/theme.js` dengan localStorage persist, deteksi system preference, auto-apply `.dark` + `.p-dark` ke `<html>`
- [x] **HeaderBar** — Theme toggle button (moon/sun icon) dengan tooltip bilingual
- [x] **Sidebar** — Dark mode classes (`dark:bg-gray-800`, `dark:border-gray-700`, dll)
- [x] **All views** — Dark mode classes di text, bg, border, cards, dialogs, tables
- [x] **Form dialogs** — Section headers, dividers, labels, input fields dark mode
- [x] **PrimeVue integration** — `.p-dark` selector untuk PrimeVue theme
- [x] **main.css** — `@custom-variant dark`, dark scrollbar, smooth transitions

### B.12. Skeleton Loading Components ✅ (BARU - Done)
- [x] **SkeletonTable.vue** — Reusable table skeleton dengan props columns (compound/tag/icons/checkbox/key-copy) + rows
- [x] **SkeletonCard.vue** — 6 card skeleton types: kpi, stat, metric, alert, sparkline, detail
- [x] **useSkeletonPage composable** — Standardized loading/loaded/error states dengan `wrapLoad()` helper
- [x] **Dashboard** — Skeleton KPI cards, chart, recent companies, system health dengan `<Transition>` fade
- [x] **Monitoring** — Skeleton untuk stat cards + table
- [x] **All remaining views** — Integrasi useSkeletonPage ke Users, Modules, Licenses, Packages, Rbac, Profile
- [x] **Transition fade** — `<Transition name="fade">` antara skeleton dan konten

---

## C. Phase 2: Tenant — Module Views

### C.1. Layout & Navigation ✅ (Existing - Enhanced)
- [x] Sidebar with all module links
- [x] HeaderBar with user info
- [x] Responsive sidebar (collapsible)
- [x] 15 module routes registered
- [x] Dashboard with KPI cards + quick access
- [x] **Dark mode** — All layout components dark mode classes (bg, text, border)
- [x] **Skeleton components** — Copied SkeletonCard, SkeletonTable, useSkeletonPage to tenant frontend
- [x] **Sidebar company name** — Dynamic company_name from `GET /api/v1/platform/users/:id` fetch + local ref (fallback: 'Company' / 'HRIS Platform')
- [x] **Sidebar PanelMenu dark hover** — Custom `:deep(.p-dark ...)` for green hover in dark mode

### C.2. Dashboard ✅ (Existing - Enhanced)
- [x] KPI Cards (Total Employees, Active Today, On Leave, Pending Approvals)
- [x] Quick Access Modules grid (12 modules)
- [x] Recent Activity (static)
- [x] Period filter (This Month/Quarter/Year)
- [x] **Dark mode** — All cards, modules, activity `dark:` classes
- [x] **Bilingual** — All static text converted to `t()` locale lookups via `computed()`
- [x] **KPI labels bilingual** — `dashboard.kpi_total_employees`, `dashboard.kpi_active_today`, dll
- [ ] **Ganti mock data dengan real API calls**
  - GET /api/v1/tenant/employees?per_page=1 (total count)
  - GET /api/v1/tenant/attendance/events (active today)
  - GET /api/v1/tenant/leave/requests (on leave count)
  - GET /api/v1/tenant/approval/tasks (pending approval count)

### C.3. Tenant Authentication & Login ✅ (BARU - Done)
- [x] **Login page** — `views/Login.vue` bilingual + dark mode (emerald brand, gradient bg)
- [x] **Auth store** — `stores/auth.js` dengan login/logout/refresh, localStorage persist (prefix `tenant_`)
- [x] **API service** — `services/api.js` axios instance dengan Language header + Auth token interceptor + 401 auto-refresh
- [x] **Router guard** — Navigation guard `beforeEach`: redirect ke `/login` jika tidak authenticated, redirect ke `/dashboard` jika sudah login
- [x] **HeaderBar** — User menu dengan Profile + Logout, display `authState.user?.name`
- [x] **HeaderBar dropdown fix** — Menggunakan `menu.toggle($event)` (PrimeVue API) bukan boolean ref
- [x] **Bilingual** — Auth locale keys (title, email, password, button, validation, profile, logout) EN/ID
- [x] **Same backend** — Tenant users login via `/api/v1/platform/login` (PlatformUser dengan company_id)
- [x] **User data enrichment** — `onMounted` fetch `GET /api/v1/platform/users/:id` untuk company_name + last_login

### C.4. Tenant Profile Page ✅ (BARU - Done)
- [x] **User info card** — Avatar (emerald), name, email, role (Tag), company, status (Tag), last login
- [x] **Change password form** — Current/New/Confirm password dengan PrimeVue Password component (toggleMask + feedback)
- [x] **Client-side validation** — Required current_password, min 6 chars, confirm match
- [x] **Server validation** — `getValidationErrors()` + Toast notification
- [x] **API** — `PUT /api/v1/platform/users/:id/password` untuk change password
- [x] **Enriched data** — Fetch `GET /api/v1/platform/users/:id` untuk company_name + updated last_login
- [x] **Bilingual** — 14+ locale keys (profile.*) EN/ID
- [x] **Dark mode** — All text, bg, Avatar, Tag dark: classes
- [x] **Route** — `/profile` with auth guard
- [x] **Router breadcrumb** — Static `title: 'Profile'` untuk HeaderBar breadcrumb

### C.5. Tenant Dark/Light Mode ✅ (BARU - Done)
- [x] **Theme store** — Copied from platform-admin: localStorage persist, system preference, `.dark` + `.p-dark` auto-apply
- [x] **main.css** — `@custom-variant dark`, dark scrollbar, dark DataTable + PanelMenu overrides
- [x] **AppLayout** — `dark:bg-gray-900` root + main bg
- [x] **HeaderBar** — Theme toggle button (moon/sun), dark mode bg/text/border
- [x] **Sidebar** — `dark:` classes, PanelMenu dark hover fix, tenant label dark mode
- [x] **Dashboard** — All KPI cards, module grid, activity, skeleton `dark:` classes
- [x] **Locale keys** — `dashboard.light_mode` / `dashboard.dark_mode` EN/ID

### C.6. Tenant Shared Components ✅ (BARU - Done)
- [x] **SkeletonTable.vue** — Copied from platform-admin with compound/tag/icons/checkbox/key-copy columns
- [x] **SkeletonCard.vue** — 6 card skeleton types (kpi/stat/metric/alert/sparkline/detail)
- [x] **useSkeletonPage.js** — Copied composable with loading/loaded/error states
- [x] **ConfirmDeleteDialog.vue** — Custom dialog component yang tetap terbuka selama API call:
  - Props: visible, title, message, loading, errorMsg, cancelLabel, confirmLabel
  - Emits: update:visible, confirm, cancel
  - `closable`, `dismissable-mask`, `close-on-escape` disable saat loading
  - Error message ditampilkan sebagai red banner di dalam dialog (bukan toast)
  - Dialog hanya close setelah sukses response dari API

### C.6a. DateInput Bilingual + FormatDate Utility ✅ (BARU - Done)
**Files:**
- `utils/formatDate.js` — **NEW** Bilingual date formatting: `formatDate(value, lang)` → "30 July 2026" (EN) / "30 Juli 2026" (ID); `formatDateShort(value)` → "30/07/2026"
- `utils/primevueLocale.js` — **NEW** PrimeVue locale configs EN/ID for DatePicker (month/day names, today/clear buttons, first day of week)
- `components/DateInput.vue` — **UPDATED** Added `useI18n` import + reactive `:locale="primeLocale"` prop to PrimeDatePicker (calendar popup shows bilingual month/day names)

**Bug fix:** `openEditDialog()` di OrganizationSummary.vue — Changed from `new Date(item.decree_date + 'T00:00:00')` (timezone-ambiguous, caused empty field) to `item.decree_date || null` (raw string, parsed safely by DateInput's `toDate()` via `new Date(y, m-1, d)` without timezone issues)

### C.6b. ConfirmDeleteDialog Refactoring — 20+ Tenant Views ✅ (BARU - Done)
- Mengganti `confirm.require()` (PrimeVue ConfirmDialog) dengan `ConfirmDeleteDialog` custom component
- Pattern baru: `confirmDelete(item)` → set `deleteTarget`, clear `deleteError`, buka dialog
- Pattern baru: `handleDelete()` async → `deleting=true`, API call, close dialog on success, show error inside dialog on failure
- Semua `import { useConfirm }`, `const confirm = useConfirm()`, `<ConfirmDialog />` dihapus

**Files refactored:**
| No | File | Perubahan |
|:--:|------|-----------|
| 1 | `components/ConfirmDeleteDialog.vue` | **NEW** — Custom dialog component |
| 2 | `settings/ZonesView.vue` | Full konversi dari `confirm.require()` |
| 3 | `modules/Organizations.vue` | Full konversi dari `confirm.require()` |
| 4 | `modules/OrganizationSummary.vue` | Remove leftover `<ConfirmDialog />` + import |
| 5 | `settings/ProvincesView.vue` | Remove leftover `<ConfirmDialog />` + import |
| 6 | `settings/RegenciesView.vue` | Same |
| 7 | `settings/DistrictsView.vue` | Same |
| 8 | `settings/VillagesView.vue` | Same |
| 9 | `settings/EducationsView.vue` | Same |
| 10 | `settings/ReligionsView.vue` | Same |
| 11 | `settings/MaritalStatusesView.vue` | Same |
| 12 | `settings/RelationshipTypesView.vue` | Same |
| 13 | `settings/EmploymentStatusesView.vue` | Same |
| 14 | `settings/NationalitiesView.vue` | Same |
| 15 | `settings/BanksView.vue` | Remove duplicate `handleDelete` + leftover `<ConfirmDialog />` |
| 16 | `settings/JobFamiliesView.vue` | Remove leftover `<ConfirmDialog />` + import |
| 17 | `settings/SalaryGradesView.vue` | Same |
| 18 | `settings/TersView.vue` | Same |
| 19 | `settings/PtkpsView.vue` | Same |
| 20 | `settings/GradingsView.vue` | Remove duplicate `handleDelete` + leftover old code fragment (wrong API path `/gradings/` → `/settings/gradings/`) |

### C.7. Setting Module — All 19 Reference CRUDs ✅ (BARU - Done)
**Backend:** `backend/internal/modules/setting/` — packages for zones, provinces, regencies, districts, villages, educations, education_majors, religions, marital_statuses, relationship_types, banks, employment_statuses, nationalities, job_families, salary_grades, ters, ptkps, insurances, company_holidays + 5 legacy endpoints
- [x] **Zones CRUD** — Code, Name, Region, IsActive, SortOrder
- [x] **Provinces CRUD** — Code, Name, SortOrder
- [x] **Regencies CRUD** — Code, Name, ProvinceID, SortOrder
- [x] **Districts CRUD** — Code, Name, RegencyID, SortOrder
- [x] **Villages CRUD** — Code, Name, DistrictID, SortOrder
- [x] **Educations CRUD** — Code, Name, SortOrder
- [x] **Education Majors CRUD** — Code, Name, SortOrder — Backend: model + repo + service + handler + routes + module; Frontend: EducationMajorsView.vue + route + SettingsIndex card + locale keys (EN/ID); OpenAPI: 5 endpoints injected + report regenerated (371 paths, 684 endpoints, 439 schemas); **Seeder: `seedEducationMajors` pakai kode 3 digit (001–020)** — UUID deterministik per kode, idempotent (20 jurusan)
- [x] **MaritalStatuses CRUD** — Code, Name, SortOrder
- [x] **RelationshipTypes CRUD** — Code, Name, SortOrder
- [x] **Banks CRUD** — Code, Name, SortOrder
- [x] **EmploymentStatuses CRUD** — Code, Name, SortOrder
- [x] **Nationalities CRUD** — Code, Name, SortOrder
- [x] **JobFamilies CRUD** — Code, Name, SortOrder
- [x] **SalaryGrades CRUD** — Code, Name, Grade, MinSalary, MaxSalary, SortOrder
- [x] **TER CRUD** — Group, BrutoMin, BrutoMax, Rate, SortOrder
- [x] **PTKP CRUD** — Name, Group, PTKP Amount, SortOrder
- [x] **Insurances CRUD** — Code, Name, SortOrder — Backend: model + repo + service + handler + routes + module; Frontend: InsurancesView.vue + route + sidebar + locale keys (EN/ID); OpenAPI: 5 endpoints injected + report regenerated (371 paths, 684 endpoints, 439 schemas)
- [x] **Company Holidays CRUD** — HolidayDate (unique), Name, Description, IsActive — Backend: model + repo + service + handler + routes + module (permissions `setting.company_holiday.*`); Frontend: CompanyHolidaysView.vue + route + SettingsIndex card + locale keys (EN/ID); OpenAPI: 5 endpoints + 3 schemas injected (373 paths, 689 endpoints, 442 schemas). **Catatan:** tabel `company_holidays` TANPA kolom updated_at/deleted_at → Update pakai `Updates(map)`, Delete hard-delete (`Unscoped`); model `IsActive` sengaja tanpa tag `default:true` (gotcha GORM false→true). **Fix bilingual:** label field description pakai key `desc`/`desc_placeholder` terpisah dari page description
- [x] **Company Holidays — Calendar View** — Toggle Table/Calendar di header (pola Organizations.vue): DatePicker inline (primevue/datepicker ^4.5.5) + slot `#date` badge rose penuh (`bg-rose-500 text-white` di angka tanggal libur) + panel legend "Holidays This Month" + klik tanggal: ada libur → edit, belum ada → modal tambah dengan tanggal terisi (`openDialog(null, key)`); `viewYear/viewMonth` ikut `@month-change`/`@year-change` (payload `{ month: 1-indexed, year }`) agar legend tidak basi; `loadCalendar()` fetch `per_page=500` (maxPerPage=500) + refresh setelah save/delete. **Bug fix ×2:** (1) slot `#date` PrimeVue v4 menerima OBJECT `{ day, month 0-indexed, year }` bukan instance Date → pakai `date.day` + `toDateKey()` dual-mode (Date & object slot) — mencegah `date.getDate is not a function`; (2) normalisasi key `normDateKey()` (`slice(0,10)` → selalu `YYYY-MM-DD`) dipakai di `holidayMap` + filter/sort `monthHolidays` agar badge tetap tampil apa pun format `holiday_date` dari API (dengan waktu/TZ)

**Frontend — Pattern Companies.vue (applied to all 19 views):**
- [x] **Search bar** (IconField + InputText) — client-side filter by code/name
- [x] **DataTable** — `p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden` styling
- [x] **SkeletonTable** — Proper `skeletonColumns` array config (not just column count)
- [x] **filteredItems computed** — search by code/name + client-side pagination
- [x] **Dialog section header** — Icon indigo (`pi pi-map-marker`/`pi pi-book`/`pi pi-globe`/`pi pi-heart`/`pi pi-users`/`pi pi-building`/`pi pi-briefcase`/`pi pi-flag`/`pi pi-star`/`pi pi-dollar`)
- [x] **Actions** — `v-tooltip.left` with bilingual tooltips
- [x] **Sortable columns** — All data columns sortable
- [x] **Footer** — `ml-auto` layout (Cancel + Save buttons)
- [x] **Dead code cleaned** — `onPage`, `page`, `total` refs removed
- [x] **Data loading** — All via `?per_page=200` (client-side pagination & filter)
- [x] **Dark mode classes** — All bg, text, border, dialog elements
- [x] **Bilingual** — All labels via `t()` locale lookups
- [x] **ConfirmDeleteDialog** (custom) — Delete confirmation dengan custom component (bukan PrimeVue ConfirmDialog), tetap terbuka selama API call, error tampil inline
- [x] **Locale keys** — 6+ keys per entity (title, description, create, edit, delete, confirm_delete)

### C.8. Organization Management ✅ (BARU - Done)
**Backend:** 12 endpoints
- [x] **TreeTable view** — Organization hierarchy dengan PrimeVue TreeTable (semua level, dinamis — recursive `OrgTreeTable.vue`)
- [x] **CRUD organization** — Dialog create (dengan parent selection), edit, delete
- [x] **Organization Tree** — Hierarki parent-child, expandable, striped rows
- [x] **Create with parent** — Otomatis set parent_id + parent label banner
- [x] **Edit** — Pre-filled form, parent bisa diubah (dropdown parent di dialog + eksklusi self/descendant)
- [x] **Tree view drag & drop (pindah parent)** — `views/modules/Organizations.vue` mode Tree: PrimeVue Tree `:draggable-nodes`/`:droppable-nodes`/`validate-drop` + `@node-drop` (`dragNode`/`dropNode`/`dropPosition`/`accept`); parent baru dari `dropPosition` (0 = drop ON node → parent = dropNode; ±1 = drop sebelum/sesudah → sibling → parent = parent dropNode; dropNode null = root); guard anti-siklus (self/descendant → toast error + `accept()` + reload); parent sama → settle + reload; parent beda → `accept()` + `PUT /organizations/:id { parent_id }` + toast sukses + reload; `moving` ref guard concurrency (drag dinonaktifkan selama proses). **Bug fix:** kode lama pakai props `:draggable`/`droppable` + `event.node`/`event.dropIndex` — tidak ada di PrimeVue 4.5.5 (API asli: `draggableNodes`/`droppableNodes`/`validateDrop` + `dragNode`/`dropPosition`/`accept`), sehingga drag tidak pernah aktif
- [x] **Delete** — ConfirmDialog bilingual (warning: child nodes juga terhapus)
- [x] **Skeleton loading** — 5 row skeleton saat loading
- [x] **Bilingual + dark mode** — t() locale keys + dark: classes
- [x] **Locale keys:** 20+ keys `organization.*` EN/ID
- [ ] **Positions CRUD** — ⏸️ **Postponed**

### C.9. Employee Management ✅ (8/9 Steps — Sisanya ⏸️ Postponed)
**Backend:** 34+ endpoints (29 CRUD + photo upload/delete + documents upload)
- [x] **DataTable employees** — List with search, pagination, sorting
- [x] **Create/Edit Employee** — Two-column layout (nav sidebar + form content)
  - [x] Step 1: **Personal Data** (`PersonalForm.vue`) — Photo upload with crop/resize, nationality type (WNI/WNA) conditional fields (NPWP, Passport, FamilyID), IG & LinkedIn social media fields, RadioLabel for gender
  - [x] Step 2: **Address** (`AddressForm.vue`) — Village autocomplete with cascading (province/regency/district auto-fill), type RadioLabel (MAIN/DOMICILE), inline delete with ConfirmDeleteDialog
  - [x] Step 3: **Emergency Contact** (`ContactForm.vue`) — Relationship type SelectLabel, delete with ConfirmDeleteDialog
  - [x] Step 4: **Family** (`FamilyForm.vue`) — **DataTable + Dialog pattern** — category Tag, edit inline
  - [x] Step 5: **Education** (`EducationForm.vue`) — **DataTable + Dialog pattern** — education Tag, edit inline
  - [x] Step 6: **Work Experience** (`ExperienceForm.vue`) — **DataTable + Dialog pattern** — edit inline
  - [x] Step 7: **Documents** (`DocumentForm.vue`) — **DataTable + Dialog + File Upload** — multipart FormData upload (POST/PUT /documents/upload), file validation (PDF/DOC/XLSX/JPG/PNG max 10MB), download link in table
  - [x] Step 8: **Insurance (BPJS)** (`InsuranceForm.vue`) — **DataTable + Dialog pattern** — insurance_id SelectLabel (options dari `/settings/insurances`, label **tanpa kode**), edit inline. Refactor: field `category` → `insurance_id` (relasi ke `insurances`, migration 026/027 + FK `fk_empins_insurance`), kolom `name` dihapus (migration 028) — nama tampil dari relasi via `insurance_name` (preload `Insurances.Insurance`), fallback `getInsuranceLabel()`
  - [x] Step 9: **Bank Account** (`BankAccountForm.vue` — ex-`BankProfileForm.vue`) — **DataTable + Dialog pattern** — bank SelectLabel (options dari `/settings/banks`, label **tanpa kode**), edit inline. **Rename komponen**: `BankProfileForm.vue` → `BankAccountForm.vue` (import & tag `<BankAccountForm>` di `EmployeeForm.vue` diupdate), label 'Bank Profile' → 'Bank Account' (locale `wizard_step_bank`/`tab_bank`)
  - [ ] Step 10: **Employment Record** (`EmploymentForm.vue`) — ⏸️ **Postponed**
- [ ] **Employee Detail Page** — Tab view — ⏸️ **Postponed**

**Key Features:**
- ✅ Step persistence via URL query param (`?step=N`) — survives page refresh
- ✅ Personal data auto-saves → redirects to edit mode (employeeId preserved)
- ✅ All forms use **DataTable + Dialog + Edit** pattern (Family, Education, Experience, Documents, Insurance)
- ✅ **Insurance options label tanpa kode** — dropdown asuransi (EmployeeForm `insuranceOptions`) hanya menampilkan nama (mis. "BPJS Kesehatan", bukan "BPJS Kesehatan (01)")
- ✅ **Bank Account rename** — komponen `BankProfileForm.vue` di-rename jadi `BankAccountForm.vue` + label 'Bank Profile' → 'Bank Account' (locale `wizard_step_bank`/`tab_bank`) — konsisten dengan penamaan label baru
- ✅ **Document file upload** — file picker + FormData + validation + download link
- ✅ **Photo upload** — crop/resize with cropper.js, multipart upload
- ✅ **Village autocomplete** — cascading province/regency/district auto-fill
- ✅ **ConfirmDeleteDialog** — consistent pattern across all forms
- ✅ **Bilingual + dark mode** — all labels, tooltips, toasts
- ✅ **FormRow/TextInput/SelectLabel** — reusable components

### C.10. Job Management 🔴 (BARU)
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

### C.11. Competency Management 🔴 (BARU)
**Backend:** 35 endpoints (7 entities)
- [ ] Competencies CRUD
- [ ] Competency Values CRUD
- [ ] Competency Events + Targets CRUD
- [ ] Competency Scores + Details CRUD

### C.12. Employee Movement 🔴 (BARU)
**Backend:** 15 endpoints
- [ ] Movements list (DataTable with status workflow)
- [ ] Create Movement (type: promotion/demotion/mutation, etc.)
- [ ] Approve / Execute / Cancel workflow buttons
- [ ] Contracts CRUD (PKWT/PKWTT)

### C.13. Time & Attendance 🔴 (BARU)
**Backend:** 30 endpoints (10 entities)
- [ ] Company Settings
- [ ] Shifts CRUD
- [ ] Employee Shifts assignment
- [ ] Attendance Events (check-in/out log)
- [ ] Daily Sessions view
- [ ] Overtime Requests CRUD
- [ ] Location management (geofence)
- [ ] Exempt Positions

### C.14. Approval Engine 🔴 (BARU)
**Backend:** 15 endpoints
- [ ] Approval Flows CRUD (with multi-step)
- [ ] My Tasks (pending approvals list)
- [ ] Approval Instances history
- [ ] Approve / Reject actions

### C.15. Payroll & Compensation 🔴 (BARU)
**Backend:** 47 endpoints (21 entities)
- [ ] Salary Components CRUD
- [ ] Payroll Periods CRUD
- [ ] Payroll Runs (with status workflow)
- [ ] Employee Payroll Profiles
- [ ] BPJS Settings
- [ ] PPh21 Settings
- [ ] Tax Brackets
- [ ] Payslip view

### C.16. Leave & Time Off 🔴 (BARU)
**Backend:** 23 endpoints
- [ ] Leave Types CRUD
- [ ] Accrual Policies CRUD
- [ ] Leave Requests (create, approve, reject)
- [ ] Leave Balances per employee

### C.17. Performance Management 🔴 (BARU)
**Backend:** 34 endpoints
- [ ] Performance Periods CRUD
- [ ] Perspectives CRUD
- [ ] Templates CRUD
- [ ] Evaluations (with details)
- [ ] KPI/OKR Targets

### C.18. Recruitment & Onboarding (ATS) 🔴 (BARU)
**Backend:** 33 endpoints
- [ ] Job Requisitions CRUD
- [ ] Candidates CRUD
- [ ] Applications (status pipeline)
- [ ] Interviews scheduling
- [ ] Onboarding Tasks

### C.19. Reimbursement & Claim 🔴 (BARU)
**Backend:** 15 endpoints
- [ ] Reimbursement Types CRUD
- [ ] Requests (DRAFT → SUBMITTED → APPROVED → PAID)
- [ ] Items per request
- [ ] Status workflow buttons

### C.20. Training & Development 🔴 (BARU)
**Backend:** 35 endpoints
- [ ] Categories CRUD
- [ ] Courses CRUD
- [ ] Sessions CRUD
- [ ] Participants registration
- [ ] Evaluations
- [ ] Certificates

### C.21. Workforce Intelligence 🔴 (BARU)
**Backend:** 68 endpoints (analytics layer)
- [ ] Dashboard (KPI summaries)
- [ ] Headcount Planning
- [ ] Analytics dashboards (headcount, attendance, leave, payroll, etc.)
- [ ] Risk Dashboard
- [ ] Executive Dashboard
- [ ] Scenario Planning
- [ ] Organization Health metrics
- [ ] People Analytics charts

### C.22. Career Intelligence 🔴 (BARU)
**Backend:** 19 endpoints
- [ ] 9-Box Talent Grid visualization
- [ ] Talent Maps CRUD
- [ ] Career Interests
- [ ] Career Paths
- [ ] Gap Analysis view
- [ ] Succession Plans CRUD

### C.23. Package Subscription (Tenant) 🔴 (BARU)
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
| 1 | Language Store (Pinia) + localStorage persistence | 🟢 Easy | ✅ Done |
| 2 | Axios interceptor `Accept-Language` header | 🟢 Easy | ✅ Done |
| 3 | Language Switcher button di HeaderBar | 🟢 Easy | ✅ Done |
| 4 | Composable `useI18n` + locale files EN/ID | 🟡 Medium | ✅ Done |
| 5 | Response handler + bilingual toast notification | 🟡 Medium | ✅ Done |
| 6 | Integrasi bilingual di semua view (label, header, filter) | 🟡 Medium | ✅ Done |

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

### Phase 2 — Tenant Foundation & Auth — Estimasi: 1 minggu ✅ Selesai
| Priority | Feature | Kompleksitas | Status |
|:--------:|---------|:------------:|:------:|
| P0 | Login Page + Auth Store + API Service | 🟢 Easy | ✅ Done |
| P0 | Profile Page + Change Password | 🟢 Easy | ✅ Done |
| P0 | Dark/Light Mode Toggle | 🟡 Medium | ✅ Done |
| P0 | Skeleton Components (Table/Card) | 🟡 Medium | ✅ Done |
| P0 | Sidebar Company Name (dynamic) | 🟢 Easy | ✅ Done |
| P0 | Dashboard Dark Mode + Bilingual | 🟡 Medium | ✅ Done |

### Phase 3 — Tenant Core Modules — Estimasi: 4-6 minggu
| Priority | Feature | Kompleksitas | Status |
|:--------:|---------|:------------:|:------:|
| P0 | Organization Management (Tree + CRUD) | 🟡 Medium | ✅ Done |
| P0 | Setting Module — All 19 Reference CRUDs (incl. TER & PTKP & Company Holidays) | 🟡 Medium | ✅ Done |
| P0 | Employee Management (Wizard) | 🔴 Complex | 🟡 8/9 Steps ✅ |
| P1 | Leave & Attendance | 🟡 Medium | 🔴 TODO |
| P1 | Payroll (read-only payslip) | 🟡 Medium | 🔴 TODO |
| P2 | Job Management | 🔴 Complex | 🔴 TODO |
| P2 | Competency Management | 🔴 Complex | 🔴 TODO |

### Phase 4 — Tenant Advanced Modules — Estimasi: 4-6 minggu
| Priority | Feature | Kompleksitas | Status |
|:--------:|---------|:------------:|:------:|
| P1 | Performance Management | 🔴 Complex | 🔴 TODO |
| P1 | Recruitment (ATS Pipeline) | 🔴 Complex | 🔴 TODO |
| P2 | Approval Engine (Flow Builder) | 🔴 Complex | 🔴 TODO |
| P2 | Employee Movement (Workflow) | 🟡 Medium | 🔴 TODO |
| P2 | Training Management | 🟡 Medium | 🔴 TODO |
| P2 | Reimbursement | 🟡 Medium | 🔴 TODO |

### Phase 5 — Intelligence & Subscription — Estimasi: 2-3 minggu
| Priority | Feature | Kompleksitas | Status |
|:--------:|---------|:------------:|:------:|
| P1 | Workforce Intelligence Dashboards | 🔴 Complex | 🔴 TODO |
| P1 | Career Intelligence (9-box Grid) | 🔴 Complex | 🔴 TODO |
| P1 | Package Subscription (Tenant) | 🟡 Medium | 🔴 TODO |

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
