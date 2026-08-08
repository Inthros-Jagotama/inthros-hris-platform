# Notification Module Development Plan

## Implementation Status

| Phase | Status | Catatan |
|---|---|---|
| Phase 1 — Database Schema | ✅ Selesai | Migration `074_notification`, model `Notification`, repository CRUD dasar + test. |
| Phase 2 — Notifier Interface & Service Layer | ✅ Selesai | Service `Notify`/`ListNotifications`/`MarkAsRead`/`MarkAllAsRead`/`GetUnreadCount`, `module.go` dengan permissions `notification.view`/`notification.manage`. |
| Phase 3 — REST API | ⏳ Belum dimulai | |
| Phase 4 — Integrasi Leave | ⏳ Belum dimulai | |
| Phase 5 — Rollout ke Modul Lain | ⏳ Belum dimulai | |
| Phase 6 — Email/Push Delivery | ⏸️ Di luar cakupan | Butuh keputusan provider terpisah. |
| Phase 7 — Notification Preferences | ⏸️ Di luar cakupan | Ditunda sampai ada kebutuhan bisnis konkret. |

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

## Phase 3 - REST API

* Handler + routes untuk 4 endpoint di §7.
* Registrasi module di `main.go`.

## Phase 4 - Integrasi Consumer Pertama (Leave)

* `leave.Notifier` interface + `SetNotifier`.
* Panggil `Notify` di `HandleApprovalStatusChange` untuk APPROVED/REJECTED.
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
