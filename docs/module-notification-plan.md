# Notification Module Development Plan

## Implementation Status

| Phase | Status | Catatan |
|---|---|---|
| Phase 1 — Database Schema | ✅ Selesai | Migration `074_notification`, model `Notification`, repository CRUD dasar + test. |
| Phase 2 — Notifier Interface & Service Layer | ✅ Selesai | Service `Notify`/`ListNotifications`/`MarkAsRead`/`MarkAllAsRead`/`GetUnreadCount`, `module.go` dengan permissions `notification.view`/`notification.manage`. |
| Phase 3 — REST API | ✅ Selesai | Handler + routes untuk 4 endpoint, module terdaftar di `main.go` (priority 19). |
| Phase 4 — Integrasi Leave | ✅ Selesai | `leave.Notifier` + `SetNotifier`, `Notify` dipanggil di `HandleApprovalStatusChange` untuk APPROVED/REJECTED, resolusi `employee_id → user_id` via `employee_accounts`. |
| Phase 5 — Rollout ke Modul Lain | ⏳ Belum dimulai | |
| Phase 6 — Email/Push Delivery | ⏸️ Di luar cakupan | Butuh keputusan provider terpisah. |
| Phase 7 — Notification Preferences | ⏸️ Di luar cakupan | Ditunda sampai ada kebutuhan bisnis konkret. |
| Phase FE-1 — Store & Bell Dropdown | ⏳ Belum dimulai | Lihat §13. FE saat ini 0% — bell di `HeaderBar.vue` murni kosmetik (badge hardcoded `"3"`). |
| Phase FE-2 — Halaman Notifikasi Penuh | ⏳ Belum dimulai | Lihat §13. |

## 1. Objective

Membangun modul **Notification** yang bertanggung jawab untuk membuat, menyimpan, mengirim (in-app), dan melacak status baca/belum-baca notifikasi yang dipicu oleh modul lain di HRIS ini.

Modul ini murni **greenfield** — dikonfirmasi lewat pencarian menyeluruh terhadap codebase bahwa **tidak ada modul Notification sama sekali** di `backend/internal/modules/*` (tidak ada folder, tabel, atau service), dan **tidak ada satupun modul lain** (Leave, Payroll, Attendance, Performance, dll.) yang pernah terintegrasi dengan sistem notifikasi user-facing apapun. Setiap kali dokumen plan modul lain (`module-leave-plan.md` §23, `module-attendance-plan.md` §48/Phase 12) menyentuh topik notifikasi, itu selalu ditandai "ditunda — infrastruktur belum ada". Dokumen ini adalah jawaban atas ketergantungan tersebut.

Modul yang diperkirakan akan mengintegrasikan diri begitu Notification tersedia — semuanya modul yang sudah punya integrasi ke **Central Approval Module** (karena hasil approval adalah sumber notifikasi paling umum):

* Leave Management
* Attendance Management (Overtime & Correction)
* Payroll
* Performance Management (KPI & OKR)
* Employee Movement
* Reimbursement
* Recruitment

Prinsip utama:

* Seluruh Primary Key menggunakan UUID.
* Notification bersifat **tenant-scoped**, bukan platform-wide — mengikuti pola `NewTenantDBResolver` yang dipakai seluruh modul lain.
* Penerima notifikasi diidentifikasi dengan **platform `user_id`** (dari modul `useraccount`), bukan `employee_id` — konsisten dengan cara `approval` module mengidentifikasi approver (`GetUserIDsByOrganization` mengembalikan user ID, bukan employee ID).
* Notification module **tidak membuat approval engine atau workflow sendiri** — ia murni penerima event dari modul lain.
* Modul lain memanggil Notification lewat **narrow interface + adapter**, pola yang sama persis dengan `ApprovalEngine`/`AttendanceSessionUpdater` yang sudah mapan di codebase ini — bukan pola baru.
* Kegagalan pengiriman notifikasi **tidak boleh pernah menggagalkan aksi yang memicunya** (best-effort, log dan lanjut — sama seperti disiplin yang sudah dipakai di integrasi Leave→Attendance).
* Notification adalah jejak audit — tidak di-hard-delete.

---

# 2. Scope Decision — In-App Only untuk Tahap Awal

Tahap awal modul ini **hanya mencakup notifikasi in-app** (disimpan di database, ditampilkan sebagai feed/badge di frontend) — **bukan** email, push notification, atau SMS.

Alasan: tidak ada provider email/SMS/push apapun yang terintegrasi di codebase ini sama sekali (sama seperti tidak adanya infrastruktur scheduled-job yang sudah didokumentasikan sebagai gap di modul Attendance §44-45). Membangun integrasi delivery eksternal di awal akan menambah dependency infra (SMTP/provider push) yang merupakan keputusan terpisah dari sekadar "menyimpan dan menampilkan notifikasi". Email/push didorong ke fase belakangan (§8, Phase 6) dan eksplisit di luar cakupan pass pertama.

---

# 3. Database Schema

## 3.1 notifications

Tabel utama — satu baris per notifikasi per penerima.

```text
id                 UUID PK
recipient_user_id  UUID NOT NULL      -- platform user_id, dari useraccount
type               VARCHAR(100)       -- mis. LEAVE_APPROVED, OVERTIME_REJECTED
title              VARCHAR(255)
body               TEXT
reference_type     VARCHAR(50)        -- module slug sumber, mis. "leave", "attendance"
reference_id       UUID NULL          -- id record sumber (leave_request_id, dll.)
is_read            BOOLEAN DEFAULT FALSE
read_at            TIMESTAMP NULL
created_at         TIMESTAMP
```

Index yang direkomendasikan:

```text
idx_notification_recipient (recipient_user_id, created_at)
idx_notification_recipient_unread (recipient_user_id, is_read)
idx_notification_reference (reference_type, reference_id)
```

## 3.2 notification_preferences (opsional, fase belakangan)

Per-user per-type opt-in/opt-out. **Sengaja tidak dibangun di fase awal** — belum ada kebutuhan bisnis konkret untuk granularitas ini (mengikuti disiplin "jangan membangun surface area spekulatif" yang sudah dipakai konsisten di seluruh sesi ini untuk Leave/Attendance, mis. Leave Eligibility Rules yang ditunda karena alasan sama). Revisit jika ada requirement nyata muncul.

---

# 4. Migration Convention

Mengikuti pola migration yang sudah baku di seluruh codebase ini — **bukan pola baru**:

```text
backend/internal/pkg/migrator/migrations/tenant/mysql/074_notification.sql
backend/internal/pkg/migrator/migrations/tenant/mysql/074_notification.down.sql
backend/internal/pkg/migrator/migrations/tenant/postgres/074_notification.sql
backend/internal/pkg/migrator/migrations/tenant/postgres/074_notification.down.sql
```

`074` adalah nomor migration berikutnya yang tersedia (nomor tertinggi saat ini: `073_attendance_correction_requests`). Migration ini murni `CREATE TABLE IF NOT EXISTS` untuk kedua tabel di atas — tidak butuh guard idempotent dinamis (`information_schema` + `PREPARE`/`EXECUTE`) karena bukan `ALTER TABLE` pada tabel existing, cukup pola `CREATE TABLE IF NOT EXISTS` (mysql) / native (postgres) seperti migration tabel baru lainnya di codebase ini (mis. `leave_balance_transactions`, `attendance_correction_requests`).

---

# 5. Notifier Interface & Wiring

Ini adalah bagian paling penting dari plan ini: **jangan menciptakan pola integrasi baru**. Gunakan persis pola "narrow interface + adapter, wired di `main.go`" yang sudah dipakai berulang kali di codebase ini (`ApprovalEngine`, `HolidayProvider`, `AttendanceSessionUpdater`).

## 5.1 Interface di sisi Notification

```go
// notification/service.go
type Service struct {
    repo   *Repository
    logger *zap.Logger
}

func (s *Service) Notify(ctx context.Context, recipientUserID uuid.UUID, notifType, title, body, referenceType string, referenceID uuid.UUID) error {
    n := &Notification{
        RecipientUserID: recipientUserID,
        Type:            notifType,
        Title:           title,
        Body:            body,
        ReferenceType:   &referenceType,
        ReferenceID:     &referenceID,
    }
    return s.repo.CreateNotification(ctx, n)
}
```

## 5.2 Interface di sisi consumer (contoh: Leave)

Setiap modul consumer mendefinisikan interface sempit miliknya sendiri (persis seperti `leave.AttendanceSessionUpdater` yang didefinisikan oleh Leave, bukan oleh Attendance):

```go
// leave/service.go
type Notifier interface {
    Notify(ctx context.Context, recipientUserID uuid.UUID, notifType, title, body, referenceType string, referenceID uuid.UUID) error
}

func (s *Service) SetNotifier(n Notifier) {
    s.notifier = n
}
```

`notification.Service` otomatis memenuhi interface ini karena Go interface bersifat struktural — sama seperti `attendance.Service` yang memenuhi `leave.AttendanceSessionUpdater` tanpa import eksplisit satu sama lain.

## 5.3 Wiring di main.go

```go
notificationSvc := notification.NewService(notificationRepo, l.Named("notification"))

leaveSvc.SetNotifier(notificationSvc)
attendanceSvc.SetNotifier(notificationSvc)
payrollSvc.SetNotifier(notificationSvc)
// dst. — satu baris per modul consumer, sama seperti wiring ApprovalEngine
```

---

# 6. Service Layer

```text
CreateNotification(ctx, recipientUserID, type, title, body, referenceType, referenceID) — dipanggil lewat Notify(), bukan endpoint HTTP
ListNotifications(ctx, recipientUserID, isRead *bool, page, perPage) — paginated, filter opsional oleh status baca
MarkAsRead(ctx, notificationID, recipientUserID)
MarkAllAsRead(ctx, recipientUserID)
GetUnreadCount(ctx, recipientUserID)
```

`recipientUserID` di setiap method **wajib** diambil dari `authctx.GetUserID(ctx)` sisi handler (bukan dari parameter yang bisa dipalsukan client) — user hanya boleh membaca/menandai notifikasi miliknya sendiri.

---

# 7. API Plan

```http
GET   /api/v1/tenant/notifications                 -- list milik user yang login, paginated, ?is_read=
GET   /api/v1/tenant/notifications/unread-count     -- badge counter
PATCH /api/v1/tenant/notifications/:id/read         -- tandai satu sebagai dibaca
POST  /api/v1/tenant/notifications/read-all         -- tandai semua sebagai dibaca
```

**Sengaja tidak ada endpoint `POST /notifications` untuk create** — notifikasi hanya dibuat oleh modul backend lain lewat `Notifier.Notify(...)`, bukan oleh API client secara langsung. Ini mencegah notification module disalahgunakan sebagai generic messaging API dan menjaga setiap notifikasi selalu punya `reference_type`/`reference_id` yang jelas asal-usulnya.

---

# 8. Consumer Integration — Titik Pemicu Konkret

Berdasarkan kode yang sudah ada dari pekerjaan sesi ini terhadap Leave dan Attendance, berikut titik pemicu `Notify(...)` yang bisa langsung dipasang begitu modul ini ada — semuanya sudah punya call site untuk `ApprovalEngine`, tinggal ditambahkan pemanggilan `Notify` di titik yang sama:

```text
leave.Service.HandleApprovalStatusChange
    → APPROVED: Notify(employee, "LEAVE_APPROVED", ...)
    → REJECTED: Notify(employee, "LEAVE_REJECTED", ...)

attendance.Service.HandleApprovalStatusChange (overtime)
    → APPROVED: Notify(employee, "OVERTIME_APPROVED", ...)
    → REJECTED: Notify(employee, "OVERTIME_REJECTED", ...)

attendance.Service.HandleApprovalStatusChange (correction)
    → APPROVED: Notify(employee, "CORRECTION_APPROVED", ...)
    → REJECTED: Notify(employee, "CORRECTION_REJECTED", ...)

payroll (approval flow serupa, pola sama)
performance/kpi, performance/okr (approval flow serupa, pola sama)
```

Semua titik ini butuh resolusi `employee_id → user_id` (lewat `useraccount`), yang sudah jadi konvensi baku di `approval` module — bukan konsep baru.

---

# 9. Development Phases

## Phase 1 - Database Schema ✅ Selesai

* ✅ Migration `074_notification` (mysql + postgres, up + down) — tabel `notifications` dengan index `idx_notification_recipient`, `idx_notification_recipient_unread`, `idx_notification_reference`.
* ✅ Model Go `Notification` (`backend/internal/modules/notification/model.go`) — UUID PK dengan `BeforeCreate`, `TableName()`.
* ✅ Repository CRUD dasar (`backend/internal/modules/notification/repository.go`): `CreateNotification`, `FindNotificationByID`, `ListNotificationsByRecipient` (filter `is_read`, paginasi), `UpdateNotification`.
* ✅ Test repository-level (`repository_test.go`) dengan sqlite in-memory: create+find, list dengan filter unread & paginasi.
* Belum ada `Notifier` interface, service layer, handler, routes, atau registrasi `main.go` — sesuai batas scope Phase 1 (murni DB layer), dilanjutkan di Phase 2.

## Phase 2 - Notifier Interface & Service Layer ✅ Selesai

* ✅ `notification.Service` (`backend/internal/modules/notification/service.go`) dengan `Notify`, `ListNotifications`, `MarkAsRead` (validasi kepemilikan recipient), `MarkAllAsRead`, `GetUnreadCount`.
* ✅ `module.go` — `NewModule`/`NewModuleWithService`, `Info()` dengan permissions `notification.view`/`notification.manage`, `Migrate()` (`AutoMigrate(&Notification{})`). `RegisterRoutes` masih no-op karena handler baru dibangun di Phase 3.
* ✅ Test service-level (`service_test.go`): `Notify` menghasilkan notifikasi unread dengan reference tersimpan, `MarkAsRead` menolak user lain, `MarkAllAsRead` + `GetUnreadCount` konsisten.
* Belum ada handler, routes, atau registrasi module di `main.go` — sesuai batas scope Phase 2, dilanjutkan di Phase 3.

## Phase 3 - REST API ✅ Selesai

* ✅ Handler (`backend/internal/modules/notification/handler.go`) untuk 4 endpoint di §7 — `recipientUserID` selalu diambil dari `authctx.GetUserID(ctx)`, tidak pernah dari parameter client.
* ✅ Routes (`backend/internal/modules/notification/routes.go`) — `GET /notifications`, `GET /notifications/unread-count`, `PATCH /notifications/:id/read`, `POST /notifications/read-all`.
* ✅ `module.go` — `RegisterRoutes` kini memanggil handler (sebelumnya no-op di Phase 2).
* ✅ Registrasi module di `cmd/server/main.go` (`notification.NewModule(dbManager, l)`, priority 19, tenant-scoped).
* Wiring `Notifier` ke modul consumer (Leave, dst.) belum dilakukan — itu Phase 4/5.

## Phase 4 - Integrasi Consumer Pertama (Leave) ✅ Selesai

* ✅ `leave.Notifier` interface (`backend/internal/modules/leave/service.go`) + `SetNotifier` — `notification.Service` memenuhinya secara struktural, tanpa import eksplisit satu sama lain.
* ✅ `leave.Repository.FindUserIDByEmployeeID` — resolusi `employee_id → user_id` lewat tabel `employee_accounts`, mengikuti konvensi yang sama dengan `approval.GetUserIDsByOrganization`.
* ✅ `HandleApprovalStatusChange` memanggil `Notify` lewat helper `notifyLeaveOutcome` untuk transisi ke `APPROVED_FINAL` (`LEAVE_APPROVED`) dan `REJECTED_FINAL` (`LEAVE_REJECTED`) — best-effort: jika notifier belum wired, employee tidak punya user account, atau `Notify` gagal, hanya di-log dan tidak menggagalkan approval itu sendiri.
* ✅ Wiring di `cmd/server/main.go`: `notificationSvc` dikonstruksi di awal lalu `leaveSvc.SetNotifier(notificationSvc)` sebelum module leave/notification di-mount.
* ✅ Test integrasi (`notifier_integration_test.go`): notify terkirim ke `recipient_user_id` yang benar saat approved, notify di-skip saat employee tidak punya user account, dan kegagalan `Notify` tidak menggagalkan `HandleApprovalStatusChange`.
* Dipilih sebagai modul consumer pertama karena Leave adalah modul paling matang dari pekerjaan sesi ini (approval integration + balance + calculation engine semua sudah selesai) — integrasi notifikasi jadi validasi paling murah risikonya.

## Phase 5 - Rollout ke Modul Lain

* Attendance (Overtime + Correction), Payroll, Performance/KPI, Performance/OKR, Employee Movement, Reimbursement — masing-masing menambahkan `Notifier` interface lokal + wiring, mengikuti pola Phase 4.

## Phase 6 - Email/Push Delivery (di luar cakupan pass ini)

* Butuh keputusan provider (SMTP/SES/FCM/dll.) — keputusan infra terpisah, bukan sekadar penambahan kode. Deferred eksplisit.

## Phase 7 - Notification Preferences (di luar cakupan pass ini)

* Tabel `notification_preferences` (§3.2) — deferred sampai ada kebutuhan bisnis konkret.

---

# 10. Priority

| Feature                          | Priority |
| --------------------------------- | -------- |
| Database Schema (Phase 1)         | P0       |
| Notifier Interface & Service (2)  | P0       |
| REST API (Phase 3)                | P0       |
| Leave Integration (Phase 4)       | P0       |
| Rollout ke modul lain (Phase 5)   | P1       |
| Email/Push Delivery (Phase 6)     | P2       |
| Notification Preferences (Phase 7)| P2       |

---

# 11. Final Architecture

```text
                    ┌─────────────────────┐
                    │  Notification Module │
                    └──────────┬───────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
        Notifier          notifications      REST API
        Interface            table          (list/read)
              │
   ┌──────────┼──────────┬──────────┬──────────┐
   ▼          ▼          ▼          ▼          ▼
 Leave    Attendance   Payroll  Performance   dst.
   │          │            │          │
   ▼          ▼            ▼          ▼
       Central Approval Module
       (sumber utama trigger notifikasi)
```

---

# 12. Design Principles

1. Notification module bertanggung jawab terhadap penyimpanan dan pelacakan status baca notifikasi — bukan business logic modul lain.
2. Notification bersifat tenant-scoped, mengikuti `NewTenantDBResolver` seperti seluruh modul lain.
3. Penerima diidentifikasi dengan `user_id` (platform), bukan `employee_id` — konsisten dengan konvensi `approval`.
4. Modul lain memanggil Notification lewat narrow interface + adapter (pola `ApprovalEngine`/`AttendanceSessionUpdater`), bukan import langsung.
5. Kegagalan `Notify()` tidak boleh menggagalkan aksi pemicunya — best-effort, log dan lanjut.
6. Tidak ada endpoint create notifikasi lewat API publik — hanya lewat pemanggilan internal antar-modul.
7. Notification adalah jejak audit — tidak di-hard-delete.
8. Seluruh ID menggunakan UUID.
9. Email/push/SMS dan notification preferences adalah fase terpisah di masa depan, bukan bagian dari cakupan awal modul ini.
10. Setiap notifikasi harus punya `reference_type`/`reference_id` yang jelas menunjuk ke record sumbernya, untuk audit dan navigasi dari frontend.

---

# 13. Frontend Implementation Plan

## 13.1 Ringkasan & Prinsip

Backend sudah lengkap sampai Phase 4 (schema, service, 4 endpoint REST, Leave sudah jadi consumer pertama). Frontend saat ini **0%** — dikonfirmasi lewat pemeriksaan langsung:

* `frontend/tenant/src/layouts/HeaderBar.vue:146-156` punya bell icon, tapi murni kosmetik: `Button icon="pi pi-bell"` + `Badge value="3"` dengan angka **hardcoded**, tanpa `@click`, tanpa dropdown/panel, tanpa pemanggilan API sama sekali.
* Tidak ada route apapun yang menyebut "notification" di `router/index.js`.
* Tidak ada key `notification` di `locales/en.json` maupun `locales/id.json`.
* Tidak ada file `Notifications.vue` atau sejenisnya di `views/modules/`.
* Tidak ada pola `setInterval`/polling apapun di codebase FE saat ini — refresh unread-count akan jadi pola pertama semacam ini.

Prinsip pengerjaan FE ini:

* **In-app only**, selaras dengan §2 (Scope Decision) backend — tidak ada browser push/service worker/email di FE.
* **Reuse pola yang sudah ada, jangan menciptakan pola baru**:
  * List/pagination → contoh `Approvals.vue` (`loadPendingTasks()`: `api.get(url, {params:{page,per_page}})`, baca `res.data.data`/`res.data.total`, try/catch/finally + toast error, render lewat `DataTable`/`Column` + `SkeletonTable` saat loading).
  * Global lightweight state → contoh `stores/activeModules.js`: singleton `reactive()` di luar fungsi (supaya sama-sama dipakai semua consumer), dibungkus composable `useXxx()`, dengan guard `loaded` dan method `reset()` dipanggil dari router guard saat logout. **Bukan** Pinia `defineStore` — proyek ini tidak memakainya sama sekali.
  * Panggilan API → langsung lewat `services/api.js` (axios instance, `baseURL` kosong, setiap call site menulis path lengkap `/api/v1/tenant/...` sendiri; interceptor Bearer token + `X-Tenant-ID` + refresh-token sudah otomatis).
* Modul `notification` **subscribable** (`IsCore: false` di `module.go` `Info()`, sama seperti Attendance/Leave/dll.) — jadi halaman penuh & entri sidebar tetap perlu gating `meta.module:'notification'` + `useActiveModules().hasModule('notification')` mengikuti convention module lain. Bell icon di header, karena bagian dari layout global (selalu tampil di semua halaman), tetap ditampilkan terlepas dari gating tersebut — cukup source datanya (unread count/list) yang no-op/kosong kalau module tidak disubscribe (backend akan mengembalikan list kosong, bukan error, untuk tenant yang belum punya data notifikasi).
* Tidak ada endpoint create notifikasi di FE — notifikasi murni ditampilkan/ditandai-dibaca, sesuai §7/§12.6 backend.

## 13.2 Komponen & Routing

* **`layouts/HeaderBar.vue`** — ganti bell kosmetik jadi interaktif:
  * `@click` men-toggle PrimeVue `OverlayPanel`/`Popover` berisi daftar notifikasi terbaru (mis. 5-10 item terakhir, ringkas: title + waktu relatif + indikator belum-dibaca), plus link "Lihat semua" ke halaman penuh.
  * `Badge` value diambil dari store unread-count (§13.3), bukan hardcoded — disembunyikan atau tidak dirender kalau count 0.
  * Klik item di dropdown → `PATCH /:id/read` lalu refresh store; kalau notifikasi punya `reference_type`/`reference_id` yang sudah punya halaman detail di FE (lihat §13.5), navigasi ke sana.
* **`views/modules/Notifications.vue`** — halaman daftar penuh:
  * `DataTable` lazy-paginated (pola sama dengan `Employees.vue`: `:totalRecords`, `:first`, `:rows`, `@page`) memanggil `GET /api/v1/tenant/notifications` dengan `page`/`per_page`/`is_read`.
  * Filter status baca (tab atau dropdown: Semua / Belum Dibaca) via query param `is_read`.
  * Tombol "Tandai semua dibaca" → `POST /read-all`, refresh list + store.
  * Klik baris → `PATCH /:id/read` (kalau belum dibaca) + navigasi kondisional sama seperti dropdown.
* **Route baru di `router/index.js`**: path `notifications`, `name: 'Notifications'`, `meta: { titleKey: 'notification.title', descKey: 'notification.description', icon: 'pi pi-bell', module: 'notification' }` — mengikuti pola module-gated route yang sudah ada (router guard otomatis redirect kalau `hasModule('notification')` false).
* **Sidebar** — tambahkan entri baru mengikuti pola `Sidebar.vue` module lain (`moduleSlug: 'notification'`, `permission: 'notification.view'`).
* **Locales** — tambahkan namespace `notification.*` baru di `en.json`/`id.json` (title, description, empty-state, mark_all_read, dst.) — belum ada sama sekali saat ini.

## 13.3 State: `stores/notifications.js`

Composable baru mengikuti bentuk persis `activeModules.js`:

```js
// module-level singleton, di luar fungsi — shared oleh semua consumer
const state = reactive({ unreadCount: 0, recentItems: [], loaded: false })

export function useNotifications() {
  async function fetchUnreadCount() { /* GET /unread-count, set state.unreadCount */ }
  async function fetchRecent() { /* GET /notifications?per_page=10, set state.recentItems */ }
  async function refresh() { await Promise.all([fetchUnreadCount(), fetchRecent()]) }
  function reset() { state.unreadCount = 0; state.recentItems = []; state.loaded = false }
  return { state, fetchUnreadCount, fetchRecent, refresh, reset }
}
```

`reset()` dipanggil dari router guard yang sama yang sudah memanggil `useActiveModules().reset()` saat logout (cek lokasi persis di router guard, jangan bikin hook logout baru).

## 13.4 Polling

Karena tidak ada infrastruktur push/websocket, unread-count di-refresh via `setInterval` (mis. tiap 60 detik), dipasang **satu kali di level store/HeaderBar.vue** (bukan per-komponen, supaya tidak numpuk banyak timer kalau bell/store dipakai di beberapa tempat), dengan `clearInterval` saat unmount/logout. Ini pola `setInterval` **pertama** di codebase FE — jaga sesederhana mungkin, jangan bangun generic polling utility dulu sebelum ada kebutuhan kedua yang nyata.

## 13.5 Data Shape & Navigasi Referensi

Setiap notifikasi punya `reference_type`/`reference_id` (§3.1). FE tidak perlu membangun routing map lengkap untuk semua kemungkinan `reference_type` di awal — cukup dukung yang sudah punya halaman FE nyata (mis. `leave` → halaman Leave My Requests, kalau ada request ID yang bisa di-deep-link). Untuk `reference_type` yang modul FE-nya sendiri masih placeholder (attendance, payroll, dll.), notifikasi cukup ditampilkan + bisa ditandai dibaca, tanpa navigasi — jangan membangun deep-link spekulatif ke halaman yang belum ada.

Response envelope sama dengan modul lain (`success`/`data`/`page`/`per_page`/`total`) — parsing FE harus konsisten dengan cara `Approvals.vue` membaca `res.data.data`/`res.data.total`.

## 13.6 Development Phases (FE)

* **Phase FE-1 — Store & Bell Dropdown**: `stores/notifications.js` (unread-count + recent list + polling) dan bell interaktif di `HeaderBar.vue` (OverlayPanel dropdown, badge dinamis, mark-as-read per item). Nilai tertinggi untuk end user karena langsung terlihat di semua halaman.
* **Phase FE-2 — Halaman Notifikasi Penuh**: `Notifications.vue` (list paginated, filter `is_read`, mark-all-as-read) + route + entri sidebar + locale keys.
* **Phase FE-3 (opsional, nanti)** — deep-link ke halaman detail per `reference_type` begitu makin banyak modul FE (Attendance, Payroll, dll.) punya halaman detail sendiri. **Belum dikerjakan sekarang** — dicatat sebagai perluasan alami, bukan bagian dari FE-1/FE-2.
* **Di luar cakupan eksplisit**: browser push notification/service worker, UI notification preferences (backend §3.2/Phase 7 juga di luar cakupan), email digest.
