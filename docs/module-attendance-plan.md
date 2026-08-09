> ⚠️ **Status vs. Plan ini**: dokumen ini ditulis seolah modul Attendance belum ada (greenfield). Setelah dicek ulang terhadap kode aktual, **seluruh tabel data (§2-§14) dan integrasi approval Overtime (§29-§30) sudah diimplementasikan sepenuhnya** — sama seperti pola Leave/Payroll. Awalnya **seluruh calculation/processing engine (session generation, capture validation) belum ada sama sekali**; per 2026-08-08, **geofence validation, duplicate-event detection, dan session generation/calculation (lateness, early-leave, work-minutes, cross-midnight) sudah diimplementasikan** (Phase 3-6). Per 2026-08-09, hal-hal berikut juga sudah terimplementasi: correction workflow (Phase 8), Leave integration (Phase 9), notifikasi outcome (Phase 12), dan **seluruh frontend (FE-1 s.d. FE-5)** — plus auto-resolve approval flow untuk Overtime & Correction (migration `079`) dan setting `allow_checkin_on_day_off` (migration `075`). Yang masih belum ada: Absent/Exempt detection (butuh scheduled job §44-45), Manager/HR Dashboard & Team Calendar (butuh cross-module employee read), Payroll integration (Phase 13), dan penerapan otomatis `WRONG_CHECKIN`/`WRONG_CHECKOUT` ke session (butuh perluasan logic seleksi Phase 6). Lihat section **"Implementation Status"** di bagian bawah dokumen untuk status per-fase yang sudah diverifikasi terhadap kode, dan catatan blockquote (`>`) di beberapa section untuk koreksi spesifik.

# Attendance Management Module Development Plan

## 1. Objective

Membangun modul **Attendance Management** untuk mengelola kehadiran karyawan secara terintegrasi dengan:

* Employee Management
* Organization Management
* Work Schedule / Shift Management
* Leave Management
* Approval Management
* Overtime Management
* Payroll
* Notification
* Performance Management

Modul Attendance bertanggung jawab terhadap:

```text
Attendance Configuration
        ↓
Work Schedule / Shift
        ↓
Attendance Capture
        ↓
Attendance Event
        ↓
Attendance Session
        ↓
Attendance Validation
        ↓
Attendance Correction
        ↓
Overtime
        ↓
Attendance Report
```

Prinsip utama:

* Seluruh ID menggunakan UUID.
* Raw attendance event tidak boleh diubah untuk menghilangkan histori.
* Attendance Session menjadi hasil pengolahan kehadiran per employee/per tanggal.
* Approval menggunakan **Central Approval Module**.
* Leave yang sudah approved dapat mempengaruhi Attendance Session.
* Overtime yang sudah approved dapat mempengaruhi Attendance Session.
* Attendance harus dapat menangani GPS, device, face verification, dan manual override.
* Attendance harus dapat diaudit.

---

# 2. Existing Database Structure

Berdasarkan `004_attendance.sql`, terdapat 10 tabel:

```text
attendance_company_settings
attendance_company_shifts
attendance_employee_shifts
attendance_locations
attendance_device_captures
attendance_face_captures
attendance_events
attendance_sessions
attendance_overtime_requests
attendance_exempt_positions
```

Relasi utama:

```text
attendance_company_settings

attendance_company_shifts
        │
        ▼
attendance_employee_shifts
        │
        ▼
employees
        │
        ▼
attendance_events
        │
        ▼
attendance_sessions
```

Dengan komponen tambahan:

```text
attendance_locations
attendance_device_captures
attendance_face_captures
attendance_overtime_requests
attendance_exempt_positions
```

> ✅ **Seluruh 10 tabel di atas sudah diimplementasikan sepenuhnya**, bukan proposal — model Go lengkap di `backend/internal/modules/attendance/model.go`, semua terdaftar di `AutoMigrate` (`module.go:104-115`), migration `004_attendance.sql`. `AttendanceOvertimeRequest` juga sudah punya `approval_instance_id` (migration `063_attendance_overtime_approval_instance.sql`, lihat Section 29). Section 3-14 di bawah ini masing-masing hanya perlu dibaca sebagai "sudah ada", bukan dianotasi satu-satu — koreksi substantif ada di Section 17-24 (capture/validation/calculation, belum ada) dan Section 29-30 (approval, sudah ada).

---

# 3. Attendance Company Settings

## Existing Table

```text
attendance_company_settings
```

Saat ini menyimpan:

* latitude
* longitude
* location requirement
* face requirement
* overtime enabled
* maximum distance
* late tolerance
* minimum overtime

### Existing Configuration

```text
Location Required
Face Required
Overtime Enabled
Max Distance
Late Tolerance
Minimum Overtime
```

---

## Enhancement

Tambahkan konfigurasi jika diperlukan:

```text
attendance_mode
work_day_start
work_day_end
timezone
allow_manual_attendance
allow_attendance_correction
late_calculation_type
early_leave_calculation_type
minimum_work_minutes
```

### Recommendation

Timezone sebaiknya disimpan pada level company/tenant jika belum tersedia pada konfigurasi tenant.

Attendance menggunakan timezone tersebut untuk menentukan:

```text
work_date
planned_start
planned_end
check-in
check-out
```

---

# 4. Attendance Company Shifts

## Existing Table

```text
attendance_company_shifts
```

Menyimpan:

```text
shift_name
check_in_time
check_out_time
is_cross_midnight
```

Contoh:

```text
Morning Shift

08:00 - 17:00
```

Night Shift:

```text
22:00 - 06:00
is_cross_midnight = 1
```

---

# 5. Shift Enhancement

Struktur shift sebaiknya diperluas untuk mendukung aturan attendance.

Tambahkan:

```text
break_start_time
break_end_time
break_minutes
late_tolerance_minutes
early_leave_tolerance_minutes
minimum_work_minutes
is_flexible
```

Jika perusahaan memiliki shift kompleks, sebaiknya jangan memasukkan seluruh aturan ke `attendance_company_shifts`.

Pertimbangkan tabel:

```text
attendance_shift_rules
```

untuk memisahkan master shift dengan aturan perhitungan.

---

# 6. Employee Shift Assignment

## Existing Table

```text
attendance_employee_shifts
```

Sudah mendukung:

```text
employee_id
attendance_shift_id
effective_date_from
effective_date_to
days_of_week_mask
is_day_off
```

Ini menjadi dasar penentuan shift employee pada tanggal tertentu.

---

# 7. Shift Assignment Enhancement

Perlu validasi agar tidak terjadi overlapping assignment.

Contoh tidak boleh:

```text
Employee A

01 Jan - 31 Jan
Shift A

15 Jan - 15 Feb
Shift B
```

karena terjadi overlap.

Tambahkan business validation:

```text
effective_date_from
<=
effective_date_to
```

dan:

```text
Employee + Date
=
Maximum 1 Active Shift
```

---

# 8. Attendance Locations

## Existing Table

```text
attendance_locations
```

Sudah mendukung:

```text
name
latitude
longitude
radius_m
```

Digunakan untuk geofencing.

Contoh:

```text
Head Office
Latitude  : ...
Longitude : ...
Radius    : 100m
```

---

# 9. Location Enhancement

Jika satu location memiliki aturan berbeda, tambahkan:

```text
is_active
allowed_for_checkin
allowed_for_checkout
```

Jika dibutuhkan:

```text
organization_id
```

agar lokasi dapat dibatasi untuk Organization tertentu.

---

# 10. Device Captures

## Existing Table

```text
attendance_device_captures
```

Menyimpan:

* device UUID
* device type
* OS
* model
* app version
* last seen

Digunakan sebagai identitas perangkat yang melakukan attendance.

---

# 11. Device Security

Tambahkan business rules:

```text
Device Registered
        ↓
Employee Allowed
        ↓
Attendance Capture
```

Jika diperlukan, buat tabel:

```text
attendance_employee_devices
```

untuk mapping:

```text
employee
    ↓
device
```

Sehingga dapat membatasi:

```text
Employee A
→ Device A

Employee B
→ Device B
```

---

# 12. Face Capture

## Existing Table

```text
attendance_face_captures
```

Sudah menyimpan:

```text
employee_id
captured_at
image_url
image_sha256
liveness_score
match_score
verified
provider
```

Ini sudah cukup sebagai histori hasil face verification.

---

# 13. Face Verification Flow

```text
Employee
    ↓
Capture Face
    ↓
Liveness Check
    ↓
Face Matching
    ↓
Verified
    ↓
Create Attendance Event
```

Jika gagal:

```text
Face Verification Failed
        ↓
Attendance Event INVALID
```

Jika face diwajibkan:

```text
verified = false
```

tidak boleh menghasilkan attendance valid.

---

# 14. Attendance Events

## Existing Table

```text
attendance_events
```

Ini adalah **raw attendance event**.

Event:

```text
CHECKIN
CHECKOUT
```

Informasi:

```text
employee
device
location
GPS
face
validation
time
```

---

# 15. Attendance Event Principle

`attendance_events` sebaiknya dianggap sebagai **immutable/raw event**.

Jangan melakukan:

```text
UPDATE event_time
```

untuk memperbaiki histori.

Jika terjadi kesalahan:

```text
Original Event
        ↓
Correction
        ↓
New Audit / Override
```

---

# 16. Attendance Event Correction

> ✅ **Sudah diimplementasikan (Phase 8, 2026-08-08)** — tabel `attendance_correction_requests` (migration `073_attendance_correction_requests`, mysql+postgres) + model `AttendanceCorrectionRequest` + CRUD + approval integration (module slug `"attendance"`) sudah ada; lihat Section 33-34 dan catatan Phase 8. `MISSING_CHECKIN`/`MISSING_CHECKOUT` diterapkan otomatis ke session saat approved; `WRONG_CHECKIN`/`WRONG_CHECKOUT` tercatat & bisa di-approve untuk audit tapi belum diterapkan otomatis (butuh perluasan logic seleksi Phase 6).

Saya merekomendasikan tabel baru:

```text
attendance_event_corrections
```

Field:

| Field                | Description                   |
| -------------------- | ----------------------------- |
| id                   | UUID                          |
| attendance_event_id  | Event                         |
| employee_id          | Employee                      |
| correction_type      | TIME / LOCATION / EVENT_TYPE  |
| old_value            | Nilai lama                    |
| new_value            | Nilai baru                    |
| reason               | Alasan                        |
| status               | Pending / Approved / Rejected |
| approval_instance_id | Central Approval              |
| corrected_by         | User                          |
| corrected_at         | Timestamp                     |

Namun jika correction dilakukan terhadap **Session**, bukan raw Event, gunakan pendekatan Session Correction.

---

# 17. Attendance Validation

> ✅ **Location/geofence validation sudah diimplementasikan (Phase 4, 2026-08-08)** — `geofence.go` (`haversineDistanceMeters`, `validateEventLocation`) + `applyEventValidation` di `CreateEvent`: event dicek terhadap seluruh `attendance_locations` terdaftar (radius masing-masing), fallback ke titik `attendance_company_settings.latitude/longitude/max_distance_meter`; di luar radius → `ValidationStatus = INVALID`. Field `DistanceM`/`IsInGeofence`/`ValidatedLocationID` kini terisi oleh service. Face Verification & Device Validation sengaja tetap belum ada — butuh provider face-matching & employee-device mapping yang belum ada (lihat Phase 4-5).

Setiap Event melalui validation:

```text
Employee
    ↓
Shift
    ↓
Time
    ↓
Location
    ↓
Device
    ↓
Face
    ↓
Validation
```

Status:

```text
PENDING
VALID
INVALID
OVERRIDDEN
```

---

# 18. Validation Rules

### Location

Jika:

```text
is_location_required = 1
```

maka:

```text
distance <= max_distance
```

atau:

```text
location berada dalam geofence
```

---

### Face

Jika:

```text
is_face_required = 1
```

maka:

```text
face verified = true
```

---

### Time

Check-in dibandingkan dengan:

```text
planned_start_local
```

Check-out dibandingkan dengan:

```text
planned_end_local
```

---

# 19. Attendance Sessions

> ✅ **Gap paling kritis sudah teratasi (Phase 6, 2026-08-08)** — `session.go` (`recalculateSession`) kini menulis `attendance_sessions` secara real-time setiap CHECKIN/CHECKOUT (di-wire dari `CreateEvent`): resolusi shift per tanggal, lateness/early-leave/work-minutes, cross-midnight (work_date = tanggal CHECKIN), DAY_OFF; `UpsertSession` baru ditambahkan. Yang masih belum ada dari daftar status di bawah: **Absent** detection (butuh scheduled job §44-45), **Exempt** (butuh cross-module employee read), dan **Payroll** integration (Phase 13) — Leave integration sudah ada di Phase 9.

## Existing Table

```text
attendance_sessions
```

Ini merupakan **hasil perhitungan attendance per employee per work date**.

Unique:

```text
employee_id
work_date
```

Status:

```text
OPEN
CLOSED
MISSING_CHECKIN
MISSING_CHECKOUT
ABSENT
DAY_OFF
EXEMPT
LEAVE
```

---

# 20. Session Processing

Flow:

```text
Employee
    ↓
Determine Shift
    ↓
Determine Work Date
    ↓
Determine Planned Start/End
    ↓
Read Attendance Events
    ↓
Select Valid Check-in
    ↓
Select Valid Check-out
    ↓
Calculate Attendance
    ↓
Create / Update Session
```

---

# 21. Session Calculation

Session menghitung:

```text
lateness_minutes
early_leave_minutes
work_minutes
break_minutes
overtime_minutes
```

---

## Lateness

Contoh:

```text
Planned:
08:00

Actual:
08:15
```

Maka:

```text
lateness = 15 minutes
```

Jika tolerance:

```text
late_tolerance = 10
```

maka:

```text
Calculated Late = 5 minutes
```

atau sesuai policy perusahaan.

---

# 22. Early Leave

Contoh:

```text
Planned Checkout:
17:00

Actual Checkout:
16:45
```

Maka:

```text
early_leave = 15 minutes
```

---

# 23. Work Minutes

Contoh:

```text
Check-in:
08:00

Check-out:
17:00

Break:
60 minutes
```

Maka:

```text
Work Minutes
= 9 hours - 1 hour
= 480 minutes
```

---

# 24. Cross Midnight Shift

Contoh:

```text
Shift:
22:00 - 06:00

Work Date:
01 Jan
```

Check-in:

```text
01 Jan 22:00
```

Check-out:

```text
02 Jan 06:00
```

Session tetap:

```text
work_date = 01 Jan
```

Ini harus menjadi bagian dari `AttendanceCalculationService`.

---

# 25. Day Off

Jika employee tidak memiliki jadwal kerja:

```text
status = DAY_OFF
```

Attendance event pada day off dapat tetap disimpan sebagai event.

Jika bekerja pada day off:

```text
is_overtime_day = 1
```

jika memenuhi policy overtime.

> ✅ **Setting `allow_checkin_on_day_off` (migration `075_attendance_allow_checkin_day_off`, 2026-08-08):** `attendance_company_settings` kini punya toggle `allow_checkin_on_day_off` (default `1` = diizinkan). Jika dimatikan, check-in pada hari libur / tanpa jadwal shift (`IsDayOff` / tidak ada assignment) **ditolak** oleh `CreateEvent` (cek di `session.go`), sehingga event tidak tersimpan; jika diizinkan, event tetap tersimpan dan session menjadi `DAY_OFF` (atau `is_overtime_day` jika memenuhi policy). Toggle ini tampil di halaman `AttendanceSettings.vue` (card "Pengaturan Perusahaan").

---

# 26. Leave Integration

> ✅ **Sudah diimplementasikan (Phase 9, 2026-08-08)** — lihat catatan Phase 9: `leave.AttendanceSessionUpdater` (interface di `leave/service.go`) memberi Leave cara memanggil Attendance tanpa import langsung, `attendance.Service.ApplyApprovedLeave` mengimplementasikannya, diwire di `main.go` (`leaveSvc.SetAttendanceSessionUpdater(attendanceSvc)`). Saat leave `APPROVED_FINAL`, tiap `LeaveRequestDetail` mendorong session jadi `LEAVE` (full day) atau mencatat `LeaveFraction` saja (half-day/hourly dengan attendance riil — session `CLOSED` tidak ditimpa, sesuai §27). Section 27 (Leave Fraction) bukan lagi proposal murni.

Attendance harus terintegrasi dengan Leave Management.

Jika Leave:

```text
APPROVED_FINAL
```

maka Attendance Session dapat menjadi:

```text
LEAVE
```

Contoh:

```text
Leave
01 Jan - 03 Jan

↓

Attendance Sessions

01 Jan → LEAVE
02 Jan → LEAVE
03 Jan → LEAVE
```

Leave yang:

```text
PENDING_APPROVAL
REJECTED
```

tidak boleh dianggap sebagai approved leave pada attendance final.

---

# 27. Leave Fraction

Existing field:

```text
leave_fraction
```

dapat digunakan untuk:

```text
1.00 = Full Day
0.50 = Half Day
```

Contoh:

```text
Half Day Leave
+
Half Day Attendance
```

Session dapat menghasilkan:

```text
leave_fraction = 0.50
work_minutes = 240
```

---

# 28. Exempt Position

## Existing Table

```text
attendance_exempt_positions
```

Menggunakan:

```text
organization_id
```

Jika Organization exempt:

```text
status = EXEMPT
```

Employee tidak dihitung sebagai absent berdasarkan attendance normal.

---

# 29. Overtime

> ✅ **Sudah diimplementasikan sepenuhnya** — satu-satunya bagian dari dokumen ini yang benar-benar sudah production-ready. `attendance/service.go` punya `ApprovalEngine` interface (`CreateApprovalInstance`, `GetApprovalInstanceStatus`, service.go:25-28), `SetApprovalEngine` (service.go:41-43), `HandleApprovalStatusChange` push-callback (service.go:585-615), dan `AttendanceOvertimeRequest.ApprovalInstanceID` (model.go:352, migration `063_attendance_overtime_approval_instance.sql`). Wired di `cmd/server/main.go:523-531` via `attendanceSvc.SetApprovalEngine(sharedApprovalEngine)` + `approvalSvc.RegisterStatusHandler("attendance", ...)` — identik dengan pola Leave/Payroll. Tidak ada endpoint approve/reject khusus di Attendance sendiri; approval sepenuhnya lewat endpoint generik Central Approval Module, sesuai desain di bawah.
>
> ✅ **Auto-resolve flow (2026-08-09):** `CreateOvertimeRequest` **dan** `CreateCorrectionRequest` kini auto-resolve active flow module `"attendance"` via `GetActiveFlowIDForModule` (pola sama dengan Leave) jika client tidak mengirim `flow_id` — request lembur yang dibuat tanpa `flow_id` (termasuk dari FE) sekarang benar-benar masuk Central Approval Module (status `PENDING_APPROVAL`) alih-alih diam di `SUBMITTED`, sehingga muncul di daftar approval + mendapat notifikasi.
>
> ✅ **Actual/Calculated Overtime sudah diimplementasikan (Phase 7, 2026-08-08)** — `applyOvertimeCalculation` dipanggil dari `HandleApprovalStatusChange` saat APPROVED: `actual_minutes = actual checkout − planned checkout` (dari session hari itu), `calculated_minutes = min(actual_minutes, requested_minutes)`; migration `072_attendance_overtime_actual_calculated` menambah kedua kolom. Detail di catatan Phase 7 di bawah.
>
> ✅ **Migration `079_attendance_overtime_status_pending_approval` (2026-08-09):** kolom `status` di `attendance_overtime_requests` (MySQL) adalah `ENUM('SUBMITTED','APPROVED','REJECTED')` yang tidak memuat `PENDING_APPROVAL` — sejak auto-resolve flow di atas, request yang masuk approval disimpan dengan status tersebut dan MySQL menolak insert (`Error 1265 Data truncated`). Migration 079 menambah `PENDING_APPROVAL` ke enum (Postgres no-op, sudah `VARCHAR(255)`).
>
> ✅ **Overtime dua alur (2026-08-09; backend + FE selesai 2026-08-10):** request lembur menjadi dua tahap: (1) planned request → approval → isian aktual (jam mulai/selesai, catatan, lampiran) → approval kedua; (2) alur tugaskan bawahan (notifikasi → isian aktual → approval). Kedua alur bisa dibatalkan sebelum isian aktual. Backend ✅ (migration 080/081, upload generik, service dua-alur, notifikasi) + FE ✅ (§32b.7: assign/isi aktual/batal, kolom & badge baru, deep-link notifikasi, dropdown assignable-employees + `employee_code`, konfirmasi batal pakai `ConfirmDeleteDialog`, detail approval ber-group). Detail lengkap: **§32b** di bawah.

## Existing Table

```text
attendance_overtime_requests
```

Saat ini memiliki approval:

```text
approved_by
approved_at
approval_note
```

### Important Enhancement

Karena HRIS Anda akan menggunakan **Central Approval Module**, overtime harus diubah menjadi:

```text
attendance_overtime_requests
        │
        ▼
approval_instance_id
        │
        ▼
Central Approval Module
```

Tidak perlu membuat workflow approval khusus di Attendance.

---

# 30. Overtime Flow

```text
Employee
    ↓
Create Overtime Request
    ↓
Submit
    ↓
Central Approval Module
    ↓
Approved
    ↓
Attendance Calculation
    ↓
Overtime Minutes
```

Jika rejected:

```text
No Approved Overtime
```

---

# 31. Overtime Calculation

Existing:

```text
requested_minutes
```

Jangan langsung menganggap:

```text
requested_minutes = actual overtime
```

Bedakan:

```text
requested_minutes
approved_minutes
actual_minutes
calculated_minutes
```

Disarankan menambahkan:

```text
approved_start_time
approved_end_time
actual_start_time
actual_end_time
actual_minutes
calculated_minutes
```

Jika overtime diatur berdasarkan attendance:

```text
Actual Overtime
=
Actual Checkout
-
Planned Checkout
```

kemudian dibatasi oleh:

```text
Approved Overtime
```

---

# 32. Overtime Calculation Example

Shift:

```text
08:00 - 17:00
```

Actual checkout:

```text
19:00
```

Employee request:

```text
17:00 - 19:00
120 minutes
```

Approved:

```text
120 minutes
```

Maka:

```text
Actual Overtime = 120
Approved Overtime = 120
Calculated Overtime = 120
```

---

# 32b. RENCANA — Overtime Dua Alur: Request → Isian Aktual (SELF & ASSIGNED)

> ✅ **Status: SELESAI (per 2026-08-10) — backend tahap 1-3 ✅, FE (tahap 4) ✅.** Enhancement alur lembur menjadi dua tahap (planned request → isian aktual) dengan dua alur: lembur sendiri (SELF) dan tugaskan bawahan (ASSIGNED).
> ✅ **Selesai (2026-08-10):** migration `080` (13 kolom baru + enum status `WAITING_ACTUAL`/`ACTUAL_SUBMITTED`/`CANCELLED`), migration `081` backfill `APPROVED`→`WAITING_ACTUAL`, endpoint upload generik `POST /uploads` (`internal/pkg/upload` + config `storage.upload_dir`), service `AssignOvertimeRequest`/`SubmitActualOvertime`/`CancelOvertimeRequest`, dispatch dua instance di `HandleApprovalStatusChange` + kalkulasi dari isian aktual di instance #2, notifikasi `OVERTIME_ASSIGNED`/`OVERTIME_ACTUAL_APPROVED`/`OVERTIME_ACTUAL_REJECTED` (i18n en/id), test service dua-alur (`overtime_two_flow_test.go`).
> ✅ **FE §32b.7 selesai (2026-08-10):** `AttendanceOvertime.vue` ditulis ulang (dialog Ajukan/Tugaskan/Isi Aktual/Batal pakai `ConfirmDeleteDialog`, kolom `flow_type`/`assigned_by`/jam aktual/link lampiran, badge status baru + i18n `attendance.status_*`), dropdown assign memakai endpoint `assignable-employees` (bawahan efektif, label nama + `employee_code`), deep-link notifikasi `OVERTIME_*` di `Notifications.vue`, detail persetujuan lembur di `Approvals.vue` ditata ber-group (Informasi Lembur + Detail Aktual), perbaikan `localTime` (NaN saat `work_date` RFC3339) & reset cache `useMyEmployee` saat login/logout. Keputusan desain yang sudah dikonfirmasi user: (1) isian aktual di alur SELF **perlu approval kedua**; (2) penugasan di alur ASSIGNED **langsung kirim notifikasi ke bawahan tanpa approval penugasan** — approval hanya di isian aktual; (3) approver isian aktual = **Central Approval Module** (instance approval kedua, module slug `"attendance"` yang sama); (4) lampiran memakai **endpoint upload generik baru** (bukan `attachment_url` teks seperti leave).

## 32b.1 Requirement

### Alur 1 — Lembur sendiri (SELF)

```text
Karyawan mengajukan lembur untuk dirinya (work_date, jam planned, reason)
        ↓
Central Approval Module — instance #1 (request/planned)
        ↓
APPROVED → status WAITING_ACTUAL
        ↓
Karyawan mengisi AKTUAL: jam mulai aktual, jam selesai aktual, catatan, file lampiran
        ↓
Central Approval Module — instance #2 (isian aktual)
        ↓
APPROVED → selesai (status APPROVED final)
```

### Alur 2 — Tugaskan bawahan (ASSIGNED)

```text
Atasan membuat penugasan lembur untuk bawahan (planned + assigned_employee_id)
        ↓
Kirim notifikasi OVERTIME_ASSIGNED ke bawahan (TANPA approval penugasan)
        ↓
status = WAITING_ACTUAL
        ↓
Bawahan mengisi AKTUAL (jam mulai aktual, jam selesai aktual, catatan, lampiran)
        ↓
Central Approval Module — instance #2 (isian aktual)
        ↓
APPROVED → selesai (status APPROVED final)
```

### Aturan lintas alur

* **Batal sebelum isian aktual**: kedua alur bisa dibatalkan selama status masih `PENDING_APPROVAL` (instance #1 belum approve) atau `WAITING_ACTUAL` (sudah approve tapi belum isi aktual). Setelah `ACTUAL_SUBMITTED` tidak bisa dibatalkan oleh pemohon — hanya lewat jalur approval instance #2 (REJECTED/CANCELLED).
* **Perubahan semantik**: status terminal `APPROVED` final hanya tercapai **setelah isian aktual disetujui** — berbeda dari perilaku lama yang men-set `APPROVED` saat request disetujui.

## 32b.2 State Machine (usulan)

```text
                        ┌───────────────┐
                        │   SUBMITTED   │  ← CreateOvertimeRequest (SELF), auto-resolve flow
                        └──────┬────────┘
                               │ instance #1 (Central Approval)
                    ┌──────────┴──────────┐
                    ▼                     ▼
        ┌──────────────────┐   ┌──────────────────┐
        │  PENDING_APPROVAL │   │     REJECTED      │  ← instance #1 REJECTED / CANCELLED
        └──────┬───────────┘   └──────────────────┘
               │ APPROVED (instance #1)
               ▼
        ┌──────────────────┐   ← entri langsung juga untuk ASSIGNED:
        │  WAITING_ACTUAL   │     AssignOvertimeRequest (tanpa instance #1)
        └──────┬───────────┘
               │ SubmitActualOvertime → instance #2 (Central Approval)
               ▼
        ┌──────────────────┐
        │ ACTUAL_SUBMITTED  │
        └──────┬───────────┘
     ┌─────────┴──────────┐
     ▼                    ▼
┌──────────┐        ┌──────────┐
│ APPROVED │        │ REJECTED │  ← instance #2 REJECTED / CANCELLED
│  (final) │        └──────────┘
└──────────┘
```

* Status baru: `WAITING_ACTUAL`, `ACTUAL_SUBMITTED`, `CANCELLED` (MySQL enum perlu `MODIFY COLUMN` — pola sama migration `079`; Postgres `VARCHAR(255)` sudah muat).
* **Dispatch dua instance**: `HandleApprovalStatusChange` (module `"attendance"`) saat ini hanya mencocokkan `documentID` ke `approval_instance_id` (lihat Phase 8/29) — harus diperluas agar juga mencocokkan `actual_approval_instance_id`, supaya callback instance #1 vs #2 tidak tertukar.

## 32b.3 Field Baru (migration ~080, mysql + postgres) — ✅ selesai (migration 080, 2026-08-10)

Tabel `attendance_overtime_requests`:

| Field | Tipe | Keterangan |
|---|---|---|
| `flow_type` | varchar (SELF / ASSIGNED) | default `SELF` — pembeda alur |
| `assigned_by` | uuid NULL | pembuat penugasan (alur 2); NULL untuk SELF |
| `assigned_at` | timestamp NULL | waktu penugasan |
| `actual_start_time_local` | timestamp NULL | jam mulai aktual (isian) |
| `actual_end_time_local` | timestamp NULL | jam selesai aktual (isian) |
| `actual_note` | varchar NULL | catatan isian aktual |
| `attachment_url` | varchar NULL | URL file lampiran (hasil endpoint upload) |
| `actual_approval_instance_id` | uuid NULL | instance approval kedua |
| `actual_submitted_at` | timestamp NULL | waktu submit isian aktual |
| `actual_approved_by` / `actual_approved_at` | uuid / timestamp NULL | approver & waktu approval aktual |
| `cancelled_by` / `cancelled_at` | uuid / timestamp NULL | jejak pembatalan |

Catatan: `actual_minutes` (kolom sudah ada) diisi dari `actual_end_time_local − actual_start_time_local` saat submit isian aktual — **bukan lagi snapshot otomatis dari checkout session** seperti perilaku `applyOvertimeCalculation` saat ini (perubahan semantik yang disengaja sesuai requirement isian aktual manual).

## 32b.4 API Plan — ✅ selesai (2026-08-10)

```http
POST /api/v1/tenant/attendance/overtime-requests                 ← SELF (seperti sekarang; flow_type default SELF)
GET  /api/v1/tenant/attendance/overtime-requests/assignable-employees  ← daftar karyawan bawahan efektif (dropdown assign; nama + employee_code)
POST /api/v1/tenant/attendance/overtime-requests/assign          ← ASSIGNED (assigned_employee_id, work_date, start/end, reason)
POST /api/v1/tenant/attendance/overtime-requests/:id/actual      ← submit isian aktual (actual_start/end, actual_note, attachment_url)
POST /api/v1/tenant/attendance/overtime-requests/:id/cancel      ← batal sebelum isian aktual (PENDING_APPROVAL / WAITING_ACTUAL)
POST /api/v1/tenant/uploads                                      ← endpoint upload generik baru (multipart → URL publik)
GET  /api/v1/tenant/attendance/overtime-requests/:id             ← detail + status + isian aktual (existing, diperkaya)
```

* **Route order**: `/overtime-requests/assign` harus didaftarkan **sebelum** `/:id` (gin wildcard conflict).
* **Upload generik**: modul/paket bersama baru (bukan terikat resource employee seperti `PUT /employees/:id/photo` / `POST /employees/:id/documents/upload` yang sudah ada) — simpan ke direktori upload (perlu config `storage.upload_dir`), validasi tipe & ukuran, nama file UUID, kembalikan URL. Konsumen pertama: lampiran lembur; bisa dipakai modul lain (mis. `attachment_url` leave) belakangan.
* **Cancel + instance aktif**: request `PENDING_APPROVAL` dibatalkan lewat jalur approval (instance #1 CANCELLED — sama seperti pemetaan CANCELLED → REJECTED di handler saat ini); `WAITING_ACTUAL` dibatalkan langsung (tidak ada instance aktif). Perlu verifikasi apakah `ApprovalEngine` butuh method cancel instance (interface saat ini: `CreateApprovalInstance`/`GetApprovalInstanceStatus`/`GetActiveFlowIDForModule`) atau cukup endpoint generik modul Approval.

## 32b.5 Notifikasi Baru (katalog i18n en/id) — ✅ selesai (2026-08-10)

| Type | Penerima | Trigger |
|---|---|---|
| `OVERTIME_ASSIGNED` | bawahan (employee yang ditugaskan) | `AssignOvertimeRequest` dibuat |
| `OVERTIME_ACTUAL_SUBMITTED` | approver | task-assigned otomatis dari engine pusat (tidak perlu wiring baru) — atau eksplisit bila perlu |
| `OVERTIME_ACTUAL_APPROVED` / `OVERTIME_ACTUAL_REJECTED` | karyawan/pengisi aktual | instance #2 APPROVED / REJECTED |

## 32b.6 Service Layer (perubahan di `attendance/service.go`) — ✅ selesai (2026-08-10)

* `CreateOvertimeRequest` — SELF: perilaku existing (auto-resolve flow), status `PENDING_APPROVAL`, `flow_type=SELF`.
* `AssignOvertimeRequest` (baru) — ASSIGNED: simpan planned + `assigned_by` (dari `authctx.GetUserID` → resolusi employee), status langsung `WAITING_ACTUAL` (tanpa instance #1), kirim `OVERTIME_ASSIGNED` ke bawahan (best-effort).
* `handleOvertimeApprovalStatusChange` — dispatch dua jalur: (a) `documentID == *o.ApprovalInstanceID` → instance #1: APPROVED → `WAITING_ACTUAL` (bukan `APPROVED` final), REJECTED/CANCELLED → `REJECTED`/`CANCELLED`; (b) `documentID == *o.ActualApprovalInstanceID` → instance #2: APPROVED → `APPROVED` final + `applyOvertimeCalculation` dari isian aktual + update session; REJECTED/CANCELLED → `REJECTED`/`CANCELLED`.
* `SubmitActualOvertime` (baru) — validasi status `WAITING_ACTUAL` & kepemilikan (pengisi = employee request, baik SELF maupun bawahan pada ASSIGNED); simpan actual fields; hitung `actual_minutes`; auto-resolve flow → instance #2; status `ACTUAL_SUBMITTED`.
* `CancelOvertimeRequest` (baru) — validasi status ∈ {`PENDING_APPROVAL`, `WAITING_ACTUAL`}; catat `cancelled_by/at`; cancel instance aktif bila ada; status `CANCELLED`.
* `applyOvertimeCalculation` — pemicu dipindah ke approval instance #2, dan `actual` diambil dari isian (`actual_start/end`) bukan snapshot checkout.

## 32b.7 Frontend (`AttendanceOvertime.vue`) — ✅ selesai (2026-08-10)

* Dua aksi pembuatan: **Ajukan Lembur** (SELF, form existing) dan **Tugaskan Lembur** (ASSIGNED — dropdown pilih karyawan dari endpoint `assignable-employees`, label nama + `employee_code`).
* Status `WAITING_ACTUAL`: tampil form isian aktual (TimeInput jam mulai/selesai, textarea catatan, upload lampiran → `POST /uploads` lalu simpan URL) + tombol **Simpan Aktual**.
* Tombol **Batal** pada status `PENDING_APPROVAL` / `WAITING_ACTUAL` (kedua alur) — pakai `ConfirmDeleteDialog` (seragam dengan halaman lain).
* Tabel: kolom baru `flow_type` (badge "Self"/"Ditugaskan"), `assigned_by`, jam aktual, link lampiran; badge status baru (`WAITING_ACTUAL`/`ACTUAL_SUBMITTED`) + i18n `attendance.status_*`.
* Notifikasi `OVERTIME_*` → deep-link ke halaman lembur (pola `Notifications.vue` `handleRowClick`, §13.5 notification plan).
* Detail persetujuan lembur di `Approvals.vue`: renderer khusus modul attendance (Informasi Lembur + Detail Aktual, header group + garis pemisah).

## 32b.8 Data Migration (existing requests) — ✅ selesai (migration 081, 2026-08-10)

* Request yang sudah `APPROVED` (alur lama, tanpa isian aktual): **backfill → `WAITING_ACTUAL`** — diputuskan & dieksekusi: `UPDATE ... WHERE status='APPROVED' AND actual_submitted_at IS NULL` (mysql + postgres). Diverifikasi di tenant dev: 4 request lama berhasil jadi `WAITING_ACTUAL`.
* Request dengan `actual_minutes` snapshot lama (dari approval tanpa isian manual): ikut backfill `WAITING_ACTUAL` karena actual di alur baru diisi manual.

## 32b.9 Testing Plan

* Unit/service: transisi state tiap aksi (create/assign/approve instance 1/submit actual/approve instance 2/cancel), validasi status yang dilarang, perhitungan `actual_minutes`/`calculated_minutes`, dispatch dua instance di `HandleApprovalStatusChange`.
* Integration: instance approval ganda per request, notifikasi `OVERTIME_ASSIGNED`/`OVERTIME_ACTUAL_*`, upload lampiran.
* FE: build bersih + verifikasi manual alur 1 & alur 2 end-to-end.

## 32b.10 Urutan Implementasi

1. ✅ Migration `080` (kolom baru + enum status) + `081` backfill — **selesai (2026-08-10)**, diterapkan ke tenant dev (mysql + postgres).
2. ✅ Backend: model/DTO + endpoint upload generik `POST /uploads` — **selesai** (`internal/pkg/upload`, config `storage.upload_dir`, locale keys `upload.*`).
3. ✅ Backend: service (assign/submit actual/cancel + dispatch dua instance) + notifikasi baru — **selesai** (+ test service `overtime_two_flow_test.go`).
4. ✅ FE: form dua alur + isian aktual + batal + kolom baru — **selesai** (`AttendanceOvertime.vue` ditulis ulang, `Notifications.vue` deep-link, locale keys, `Approvals.vue` detail ber-group, `ConfirmDeleteDialog` untuk batal, dropdown assignable-employees + `employee_code`).
5. 🔶 Test + review: unit/service ✅ (build + seluruh test module pass, migrator integration MySQL pass) + FE build ✅; verifikasi manual E2E alur 1 & alur 2 belum.

---

# 33. Attendance Correction

Attendance correction dibutuhkan untuk:

* Missing check-in
* Missing check-out
* Wrong check-in
* Wrong check-out
* Device failure
* GPS failure
* Face failure

Flow:

```text
Employee / HR
      ↓
Correction Request
      ↓
Approval Module
      ↓
Approved
      ↓
Recalculate Session
```

---

# 34. Recommended New Table

## attendance_correction_requests

```text
id
employee_id
attendance_session_id
correction_type
requested_checkin
requested_checkout
reason
status
approval_instance_id
created_by
approved_at
created_at
updated_at
```

---

# 35. Attendance Status

Recommended final status:

```text
PRESENT
LATE
EARLY_LEAVE
LATE_AND_EARLY_LEAVE
ABSENT
MISSING_CHECKIN
MISSING_CHECKOUT
DAY_OFF
LEAVE
EXEMPT
HOLIDAY
```

Namun jangan mengganti status existing secara langsung tanpa migration analysis.

Alternatif yang lebih baik:

```text
session.status
```

menyimpan kondisi utama:

```text
CLOSED
ABSENT
DAY_OFF
LEAVE
EXEMPT
```

sedangkan:

```text
lateness_minutes > 0
early_leave_minutes > 0
```

menentukan attribute:

```text
LATE
EARLY_LEAVE
```

---

> ✅ **Frontend sudah dibangun (2026-08-10).** `Attendance.vue` kini hub dengan menu card (Overtime, Corrections, Admin) + summary cards (present/late/missing checkout/leave) + kalender bulan berjalan + sesi hari ini. Halaman pendukung lengkap: AttendanceAdmin, Settings, Shifts, EmployeeShifts, Locations, ExemptPositions, Overtime, Corrections, Events, Sessions, Reports (12 file `views/modules/Attendance*.vue`). §37 (Employee Dashboard) sebagian terwakili di hub; §40 Reports ada (laporan sesi); §36 (Calendar) & §38/§39 (Manager/HR Dashboard) belum punya halaman khusus.

# 36. Attendance Calendar

Menu:

```text
Attendance Calendar
```

Menampilkan:

```text
Present
Late
Absent
Leave
Day Off
Holiday
Overtime
```

Untuk:

* Employee
* Team
* Organization
* Company

---

# 37. Employee Attendance Dashboard

Menampilkan:

```text
Attendance Summary

Present       : 20
Late          : 2
Absent        : 1
Leave         : 2
Day Off       : 5
Overtime      : 8h
```

---

# 38. Manager Dashboard

Manager dapat melihat:

```text
Today's Attendance

Present
Late
Absent
Missing Checkout
On Leave
Overtime
```

Filter:

```text
Organization
Date
Employee
Status
```

---

# 39. HR Dashboard

HR dashboard:

```text
Attendance Rate
Late Rate
Absence Rate
Leave Rate
Overtime Hours
Missing Attendance
```

Breakdown:

```text
Company
Organization
Position
Employee
Date
Shift
```

---

# 40. Attendance Reports

Reports:

```text
Daily Attendance
Monthly Attendance
Late Report
Absence Report
Early Leave Report
Overtime Report
Missing Attendance
Attendance Correction
Location Validation
Face Verification
Device Attendance
```

Export:

```text
Excel
CSV
PDF
```

---

# 41. API Plan

> ⚠️ Endpoint di bawah adalah **rencana**, bukan kondisi aktual. Endpoint yang **sudah ada** di `attendance/routes.go` hari ini (prefix `/api/v1/tenant/attendance`):
> ```http
> GET/PUT        /settings
> POST/GET        /shifts
> GET/PUT/DELETE  /shifts/:id
> POST/GET        /employee-shifts
> GET/PUT/DELETE  /employee-shifts/:id
> POST/GET        /locations
> GET/PUT/DELETE  /locations/:id
> POST/GET        /events          ← endpoint generik untuk CHECKIN & CHECKOUT, bukan endpoint terpisah
> GET             /events/:id
> GET             /sessions        ← read-only, tidak ada endpoint yang menulis session
> GET             /sessions/detail
> POST/GET        /overtime-requests
> GET             /overtime-requests/:id
> POST/GET        /exempt-positions
> GET/PUT/DELETE  /exempt-positions/:id
> ```
> Yang **sudah ada tambahan** sejak catatan di atas: corrections (§34: `POST/GET /corrections`, `GET /corrections/:id`), reports (§40: `GET /reports/sessions`), calendar & summary (`GET /calendar`, `GET /summary`), serta endpoint lembur dua-alur (`/assignable-employees`, `/assign`, `/:id/actual`, `/:id/cancel`). Yang **belum ada**: endpoint `check-in`/`check-out` khusus (§41 di bawah), endpoint session write (session ditulis real-time oleh `CreateEvent`), dan endpoint approve/reject overtime khusus (approval sepenuhnya lewat endpoint generik Central Approval Module — lihat Section 29).

## Attendance Settings

```http
GET    /api/v1/tenant/attendance/settings
PUT    /api/v1/tenant/attendance/settings
```

## Shifts

```http
GET    /api/v1/tenant/attendance/shifts
POST   /api/v1/tenant/attendance/shifts
GET    /api/v1/tenant/attendance/shifts/{id}
PUT    /api/v1/tenant/attendance/shifts/{id}
DELETE /api/v1/tenant/attendance/shifts/{id}
```

## Employee Shift

```http
GET    /api/v1/tenant/attendance/employee-shifts
POST   /api/v1/tenant/attendance/employee-shifts
PUT    /api/v1/tenant/attendance/employee-shifts/{id}
DELETE /api/v1/tenant/attendance/employee-shifts/{id}
```

## Attendance Capture

```http
POST /api/v1/tenant/attendance/check-in
POST /api/v1/tenant/attendance/check-out
```

## Attendance Events

```http
GET /api/v1/tenant/attendance/events
GET /api/v1/tenant/attendance/events/{id}
```

## Sessions

```http
GET /api/v1/tenant/attendance/sessions
GET /api/v1/tenant/attendance/sessions/{id}
```

## Overtime

```http
GET  /api/v1/tenant/attendance/overtime
POST /api/v1/tenant/attendance/overtime
GET  /api/v1/tenant/attendance/overtime/{id}
POST /api/v1/tenant/attendance/overtime/{id}/submit
```

Approval dilakukan melalui Central Approval Module.

## Correction

```http
GET  /api/v1/tenant/attendance/corrections
POST /api/v1/tenant/attendance/corrections
GET  /api/v1/tenant/attendance/corrections/{id}
POST /api/v1/tenant/attendance/corrections/{id}/submit
```

---

# 42. Service Layer

Recommended services:

```text
AttendanceCaptureService
AttendanceValidationService
AttendanceCalculationService
AttendanceSessionService
AttendanceShiftService
AttendanceLocationService
AttendanceDeviceService
AttendanceFaceService
AttendanceCorrectionService
AttendanceOvertimeService
AttendanceReportService
```

---

# 43. Attendance Calculation Engine

Pisahkan calculation engine dari controller.

```text
AttendanceCalculationService
        │
        ├── resolveShift()
        ├── resolveWorkDate()
        ├── resolvePlannedTime()
        ├── calculateLate()
        ├── calculateEarlyLeave()
        ├── calculateWorkMinutes()
        ├── calculateBreak()
        └── calculateOvertime()
```

Dengan demikian logic dapat digunakan untuk:

* Real-time attendance
* Daily closing
* Recalculation
* Correction
* Payroll

---

# 44. Daily Attendance Processing

Sistem dapat memiliki scheduled job:

```text
ProcessDailyAttendance
```

Flow:

```text
Employee
    ↓
Resolve Schedule
    ↓
Create Session
    ↓
Process Events
    ↓
Process Leave
    ↓
Process Holiday
    ↓
Process Day Off
    ↓
Calculate Attendance
    ↓
Close Session
```

---

# 45. Missing Attendance Processing

Scheduled job:

```text
DetectMissingAttendance
```

Contoh:

```text
Shift:
08:00 - 17:00

Tidak ada CHECKIN

↓

MISSING_CHECKIN
```

Jika ada check-in tetapi tidak ada checkout:

```text
MISSING_CHECKOUT
```

---

# 46. Recalculation

Sistem harus mendukung:

```text
Recalculate Attendance
```

Contoh:

```text
Attendance Session
    ↓
Shift berubah
    ↓
Recalculate
    ↓
Late / Work / Overtime berubah
```

Recalculation diperlukan ketika:

* Shift diperbaiki
* Leave approved setelah attendance dibuat
* Correction approved
* Overtime approved
* Company policy berubah

---

# 47. Audit

Tambahkan audit untuk perubahan penting:

```text
Attendance Event
Attendance Session
Attendance Correction
Attendance Settings
Shift
Employee Shift
Overtime
```

Minimal mencatat:

```text
who
when
what changed
old value
new value
reason
```

---

# 48. Notification

Integrasi Notification Module.

### Employee

* Check-in success
* Check-in failed
* Check-out success
* Missing checkout reminder
* Overtime approved
* Overtime rejected
* Correction approved
* Correction rejected

### Manager

* Employee absent
* Late employee
* Overtime request
* Correction request

### HR

* Missing attendance
* Attendance anomalies
* Correction request

---

# 49. Payroll Integration

Attendance menyediakan data:

```text
work_minutes
overtime_minutes
late_minutes
early_leave_minutes
unpaid_leave_days
absence_days
```

Payroll melakukan perhitungan:

```text
salary deduction
overtime payment
attendance allowance
```

Attendance tidak sebaiknya menghitung nominal payroll.

---

# 50. Leave Integration

Leave menjadi salah satu sumber Attendance Session.

```text
Approved Leave
      ↓
Attendance Session
      ↓
LEAVE
```

Tidak boleh:

```text
Pending Leave
      ↓
LEAVE
```

---

# 51. Data Integrity

### Attendance Event

```text
employee_id
event_time
event_type
```

harus memiliki index yang baik untuk query harian.

### Session

Existing:

```text
UNIQUE(employee_id, work_date)
```

dipertahankan.

### Shift

Pastikan assignment tidak overlap.

### Device

Existing:

```text
UNIQUE(device_uuid)
```

dipertahankan.

---

# 52. Seeder

Seeder yang direkomendasikan:

```text
AttendanceCompanySettingsSeeder
AttendanceShiftSeeder
AttendanceLocationSeeder
AttendanceExemptPositionSeeder
```

Jika project memiliki pola seeder master yang sudah ada, ikuti struktur tersebut.

Seluruh ID menggunakan UUID.

---

# 53. Testing Plan

## Unit Test

### Shift

* Normal shift
* Cross midnight
* Day off
* Multiple schedules

### Attendance Calculation

* On time
* Late
* Early leave
* Late + early leave
* Missing check-in
* Missing checkout
* Full day
* Half day
* Leave
* Holiday
* Day off
* Overtime

### Location

* Inside geofence
* Outside geofence
* GPS inaccurate

### Face

* Verified
* Failed
* Low liveness
* Low match score

---

# 54. Feature Test

## Check-in

```text
Employee
→ Check-in
→ Validate Location
→ Validate Face
→ Create Event
→ Update Session
```

## Check-out

```text
Employee
→ Check-out
→ Validate
→ Create Event
→ Calculate Session
```

## Leave

```text
Leave Approved
→ Attendance Session
→ LEAVE
```

## Overtime

```text
Request Overtime
→ Central Approval
→ Approved
→ Calculate Overtime
```

## Correction

```text
Correction
→ Central Approval
→ Approved
→ Recalculate Session
```

---

# 55. Development Phases

## Phase 1 - Database Review & Enhancement

* Review existing attendance schema. ✅ Seluruh 10 tabel + kolom cross-checked terhadap `attendance/model.go` — tidak ada type mismatch seperti bug `deleted_at` yang ditemukan di Leave; semua `deleted_at`/timestamp column sudah bertipe benar.
* Review foreign keys. ✅ Direview — FK utama (`employees`, `attendance_company_shifts`) sudah lengkap dengan `ON DELETE CASCADE`/`SET NULL` yang tepat. `attendance_events.overtime_request_id`/`validated_location_id` tidak punya FK constraint (hanya index) — dicatat sebagai temuan minor, tidak diperbaiki karena bukan bug yang terbukti (bisa jadi disengaja untuk menghindari cascade coupling), bukan seperti kasus Leave yang jelas salah tipe data.
* Add missing indexes. ✅ `attendance_events` sebelumnya hanya punya index tunggal `employee_id` — tidak memenuhi kebutuhan §51 ("index yang baik untuk query harian"). Ditambahkan `idx_att_event_employee_time (employee_id, event_time_local)` via migration `071_attendance_phase1_event_index`. `attendance_sessions` (`UNIQUE(employee_id, work_date)`) dan `attendance_device_captures` (`UNIQUE(device_uuid)`) sudah benar sejak awal, tidak perlu perubahan.
* Review timezone handling. 🔶 Direview — tidak ada kolom timezone di `attendance_company_settings` maupun modul lain manapun (`setting`, dll). Ini murni proposal opsional di plan asli ("jika belum tersedia pada konfigurasi tenant") dan merupakan keputusan desain lintas-modul (bukan spesifik Attendance) — sengaja tidak ditambahkan sekarang tanpa keputusan eksplisit apakah timezone disimpan di level tenant/company yang mana.
* Add shift rule table jika diperlukan. ⏳ Deferred — tidak ada business rule konkret yang membutuhkannya saat ini (sama seperti pola deferral di Leave §22 Eligibility).
* Add employee-device mapping jika diperlukan. ⏳ Deferred — belum ada kebutuhan device-restriction konkret.
* Add correction request. ✅ Selesai (Phase 8, 2026-08-08) — migration `073_attendance_correction_requests` + model + CRUD + approval integration (module slug `attendance`) + FE `AttendanceCorrections.vue` (auto-resolve flow, input jam pakai TimeInput). Lihat §33-34.
* Integrate approval_instance_id ke overtime. ✅ Sudah ada sejak sebelumnya — migration `063_attendance_overtime_approval_instance.sql`, field `AttendanceOvertimeRequest.ApprovalInstanceID` (model.go:319).
* Review existing approval fields. ✅ `approved_by`/`approved_at`/`approval_note` masih ada di model tapi sudah tidak jadi penggerak workflow utama — approval sepenuhnya lewat `approval_instance_id` + Central Approval Module (lihat Section 29).

---

## Phase 2 - Attendance Configuration

Develop:

```text
Attendance Settings   ✅ UpsertCompanySetting/GetCompanySetting (service.go:49,76)
Shift                 ✅ CRUD lengkap (service.go:88-171)
Shift Rules           ⏳ Deferred — sama seperti keputusan Phase 1, belum ada business rule konkret yang butuh field break/tolerance terpisah dari attendance_company_shifts
Locations             ✅ CRUD lengkap (service.go:301-384)
Devices               ❌ Bukan sekadar "belum dikonfigurasi" — attendance_device_captures sama sekali tidak punya repository method maupun handler; tabelnya mati total, tidak ada satupun kode yang menulis atau membaca device capture
Face Configuration    ❌ Sama seperti Devices — attendance_face_captures juga tidak punya repository method/handler, tabel histori mati total
Exempt Positions      ✅ CRUD lengkap (service.go:666-746)
```

> 🔶 **Sebagian selesai, sisanya sengaja ditunda.** Settings/Shift/Locations/Exempt Positions sudah CRUD penuh. Shift Rules ditunda dengan alasan sama seperti Phase 1 (belum ada requirement konkret). Devices dan Face Configuration **bukan gap konfigurasi biasa** — kedua tabelnya (`attendance_device_captures`, `attendance_face_captures`) tidak punya repository method sama sekali, jadi membangun CRUD/config di atasnya sekarang akan jadi surface area tanpa konsumen, karena keduanya baru berguna setelah capture validation engine (Section 17, Phase 4-5) benar-benar memvalidasi device/face saat check-in/check-out. Revisit bersamaan dengan Phase 4-5, bukan sekarang.

---

## Phase 3 - Shift Management

Develop:

```text
Create Shift        ✅ service.go:88-171
Update Shift         ✅
Delete Shift         ✅
Assign Shift         ✅ CreateEmployeeShift/UpdateEmployeeShift (service.go:183-287)
Effective Date       ✅ EffectiveDateFrom/To sudah tervalidasi (from <= to) sejak perubahan ini
Day of Week          🔶 DaysOfWeekMask disimpan, tapi tidak ada logic decode/interpretasi bit-to-weekday di manapun (lihat catatan Leave §41 soal hal yang sama) — field pass-through saja
Day Off              🔶 IsDayOff disimpan sebagai flag, belum dikonsumsi oleh calculation engine manapun (karena calculation engine sendiri belum ada — Section 19)
Cross Midnight        ⏳ Field `IsCrossMidnight` ada di `attendance_company_shifts`, tapi logic session cross-midnight (§24) belum ada — bergantung pada Phase 6
```

> ✅ **Overlap validation (§7) — gap nyata, sekarang sudah diperbaiki.** Sebelumnya `CreateEmployeeShift`/`UpdateEmployeeShift` tidak melakukan validasi apapun terhadap `effective_date_from <= effective_date_to` maupun terhadap assignment lain yang overlap untuk employee yang sama — persis seperti yang diperingatkan §7, tapi belum diimplementasikan. Ditambahkan `CountOverlappingEmployeeShifts` (repository) + validasi di kedua service method (mengizinkan `effective_date_to` NULL sebagai open-ended). Ini satu-satunya gap Phase 3 yang genuinely fixable tanpa bergantung pada calculation engine — Day of Week mask decoding dan Cross Midnight session logic tetap bergantung pada Phase 6.

---

## Phase 4 - Attendance Capture

Develop:

```text
Check-in            🔶 Lewat POST /events generik (event_type=CHECKIN), bukan endpoint khusus — lihat §41
Check-out            🔶 Sama, lewat POST /events (event_type=CHECKOUT)
GPS                  ✅ latitude/longitude sudah diterima dan disimpan sejak awal
Geofence             ✅ Baru diimplementasikan — lihat catatan di bawah
Face Verification    ⏳ Tetap belum ada — tidak ada face-matching provider, hanya field pass-through (lihat Phase 2)
Device Validation     ⏳ Tetap belum ada — tidak ada employee-device mapping (lihat Phase 1/2, sengaja ditunda)
Attendance Event      ✅ CreateEvent sudah lengkap (raw event + sekarang location validation)
```

> ✅ **Geofence validation — gap nyata, sekarang diimplementasikan.** `CreateEvent` sebelumnya menyimpan `DistanceM`/`IsInGeofence`/`ValidatedLocationID` sebagai kosong selamanya (§17). Ditambahkan `geofence.go` (`haversineDistanceMeters`, `validateEventLocation`) + `applyEventValidation` di service: jika `IsLocationRequired`, event dicek terhadap seluruh `attendance_locations` terdaftar (radius masing-masing), fallback ke titik tunggal `attendance_company_settings.latitude/longitude/max_distance_meter` jika belum ada location terdaftar. Event yang berada di luar geofence → `ValidationStatus = INVALID`.
>
> Face verification dan device validation **sengaja tetap tidak diimplementasikan** — bukan kelalaian. Face verification butuh provider face-matching eksternal yang tidak ada di codebase ini (§13 tetap murni proposal); menandai event sebagai VALID padahal face tidak pernah benar-benar diverifikasi akan menyesatkan, jadi ketika `IsFaceRequired = true`, event sengaja dibiarkan `PENDING` bukan `VALID`. Device validation butuh `attendance_employee_devices` yang sudah sengaja ditunda di Phase 1/2 (tidak ada kebutuhan konkret) — memvalidasi tanpa tabel itu tidak mungkin.

---

## Phase 5 - Attendance Validation

Develop:

```text
Location Validation       ✅ Selesai di Phase 4 (geofence.go)
Face Validation             ⏳ Tetap belum ada — tidak ada face-matching provider (lihat Phase 4)
Device Validation           ⏳ Tetap belum ada — tidak ada employee-device mapping (lihat Phase 1/2)
Time Validation              ⏳ Belum ada — butuh resolusi shift per employee/tanggal (planned_start/end_local), yaitu bagian dari session calculation engine (Phase 6), bukan validation berdiri sendiri. Membangun versi parsial sekarang berisiko salah karena interpretasi `DaysOfWeekMask`/day-off/cross-midnight belum ada konvensinya di manapun
Duplicate Event Detection    ✅ Baru diimplementasikan — lihat catatan di bawah
```

> ✅ **Duplicate Event Detection — gap nyata, sekarang diimplementasikan.** Sebelumnya `CreateEvent` menerima event apapun tanpa cek urutan — employee bisa check-in dua kali berturut-turut tanpa check-out di antaranya, atau check-out tanpa pernah check-in. Ditambahkan `checkEventSequence` (service.go) + `FindLastEventForEmployee` (repository): CHECKIN ditolak jika event terakhir employee juga CHECKIN (belum check-out), CHECKOUT ditolak jika event terakhir bukan CHECKIN terbuka (atau belum ada event sama sekali). Ini murni cek urutan raw event — sengaja tidak menyentuh resolusi shift/work-date, itu tetap tanggung jawab Phase 6.
>
> Time Validation sengaja **tidak** diimplementasikan sebagian di sini karena akan butuh resolusi shift (`attendance_employee_shifts` + `attendance_company_shifts`) yang benar termasuk `DaysOfWeekMask`/cross-midnight — komponen yang sama yang bikin Phase 6 jadi gap paling kritis di modul ini. Membangun versi Time Validation yang naif sekarang (tanpa resolusi shift yang benar) berisiko menghasilkan validasi yang salah, lebih buruk daripada tidak ada validasi sama sekali.

---

## Phase 6 - Attendance Session

Develop:

```text
Session Generation    ✅ recalculateSession (session.go), triggered synchronously from CreateEvent
Check-in Mapping       ✅ selectCheckinCheckout — first CHECKIN of the work date
Check-out Mapping      ✅ first CHECKOUT following that CHECKIN (crosses midnight correctly)
Work Minutes           ✅ workMinutesBetween (checkout - checkin - break_minutes, floored at 0)
Late                   ✅ minutesLate (raw lateness minus late_tolerance_minutes, floored at 0)
Early Leave            ✅ minutesEarly (planned_end - actual checkout, floored at 0)
Day Off                ✅ IsDayOff on the employee's shift assignment, or no assignment at all (§25)
Absent                 ⏳ Not implemented — needs the scheduled ProcessDailyAttendance/DetectMissingAttendance job (§44-45) to proactively mark days with zero events; this engine only reacts to events that actually happened
Leave                  ⏳ Not implemented — belongs to Phase 9 (needs a read into the Leave module)
Exempt                 ⏳ Not implemented — needs the employee's organization_id, which needs a cross-module read into the employee module that doesn't exist anywhere in this codebase yet
```

> ✅ **The module's critical gap is closed.** New `session.go`: `recalculateSession(ctx, employeeID, workDate)` resolves the employee's active `attendance_employee_shifts` row for that date, loads `attendance_company_shifts`, reads the work date's events (`FindEventsForWorkDate` — a 2-day window catching cross-midnight checkouts), picks the first CHECKIN and the CHECKOUT that follows it, and computes lateness/early-leave/work-minutes before upserting `attendance_sessions` via new `UpsertSession`. Wired into `CreateEvent`: every check-in/check-out now actually recalculates that day's session, so `attendance_sessions` is no longer a dead table.
>
> **Work-date attribution for cross-midnight shifts (§24)** uses the CHECKIN event's local calendar date, not the CHECKOUT's — `CreateEvent` passes `lastEvent`'s date (the open check-in being closed) as the work date when a CHECKOUT arrives, exactly matching §24's worked example.
>
> **Timezone note**: shift `check_in_time`/`check_out_time` (e.g. `"08:00:00"`) carry no timezone of their own. Since this codebase has no company/tenant timezone setting (§3, deferred in Phase 1), planned start/end are anchored to whichever event's `event_time_local` offset is available for that work date — the only source of "local" time this system has. This is a pragmatic choice given the constraint, not a permanent design; a real timezone setting would replace it cleanly.
>
> **Deliberately not implemented**: Absent detection and Exempt status need components explicitly out of scope for this pass — a scheduled job (§44-45, no cron infrastructure was found anywhere in this codebase) and a cross-module employee/organization read (no established interface, same category of gap noted for Leave). Leave integration is its own phase (Phase 9). None of these block session generation for events that actually happen, which was the P0 gap.

---

## Phase 7 - Overtime

Develop:

```text
Overtime Request           ✅ Sudah ada sejak awal (CRUD lengkap)
Central Approval Integration   ✅ Sudah ada sejak awal — lihat Section 29
Approved Overtime               ✅ Status/ApprovedAt/ApprovedBy sudah ada sejak awal
Actual Overtime                 ✅ Baru diimplementasikan — lihat catatan di bawah
Calculated Overtime             ✅ Baru diimplementasikan — lihat catatan di bawah
Attendance Integration          ✅ Baru diimplementasikan — session di-update dengan data overtime saat approved
```

> ✅ **Actual/Calculated Overtime dan Attendance Integration — gap nyata, sekarang diimplementasikan.** `AttendanceOvertimeRequest` sebelumnya hanya punya `requested_minutes`; §31 secara eksplisit meminta field terpisah untuk `actual_minutes`/`calculated_minutes` supaya `requested_minutes` tidak langsung dianggap sebagai overtime final. Ditambahkan kedua kolom via migration `072_attendance_overtime_actual_calculated` (mysql+postgres), lalu `applyOvertimeCalculation` (service.go) dipanggil dari `HandleApprovalStatusChange` saat status APPROVED: `Actual Overtime = actual checkout - planned checkout` (dari session hari itu, sesuai formula §31), lalu `Calculated Overtime = min(actual, requested_minutes)` (dibatasi oleh yang diminta, sesuai §32's worked example). Session yang sama juga di-update: `IsOvertimeDay`, `OvertimeRequestID`, `ApprovedOvertimeStartLocal/EndLocal`, `OvertimeMinutes` — field-field yang sudah ada di model sejak awal tapi sebelumnya tidak pernah diisi.
>
> Jika approval terjadi sebelum employee benar-benar checkout (session belum punya `CheckoutEventID`), `actual_minutes`/`calculated_minutes` sengaja dibiarkan `NULL` — ini snapshot best-effort di waktu approval, bukan proses yang re-trigger otomatis saat checkout terlambat terjadi (itu di luar cakupan; approval ulang bukan alur yang didukung endpoint ini).

---

## Phase 8 - Correction

Develop:

```text
Missing Check-in       ✅ Diimplementasikan — lihat catatan di bawah
Missing Checkout        ✅ Diimplementasikan — lihat catatan di bawah
Attendance Correction    ✅ Tabel + model + CRUD baru — lihat catatan di bawah
Central Approval          ✅ Menggunakan ApprovalEngine yang sama dengan Overtime (module slug "attendance")
Session Recalculation      ✅ recalculateSession dipanggil setelah correction diterapkan
Audit                       ✅ Correction event baru diberi ValidationStatus OVERRIDDEN, raw event asli tidak pernah diubah (§15)
```

> ✅ **Correction workflow — gap nyata, sekarang diimplementasikan (sebagian).** Tabel baru `attendance_correction_requests` (migration `073_attendance_correction_requests`, mysql+postgres) + model `AttendanceCorrectionRequest` + CRUD (`CreateCorrectionRequest`/`GetCorrectionRequestByID`/`ListCorrectionRequests`) + endpoint `POST/GET /corrections`, `GET /corrections/:id`. Approval memakai `ApprovalEngine` yang sama dengan Overtime (module slug `"attendance"` yang sama) — `HandleApprovalStatusChange` sekarang dispatch ke overtime request atau correction request tergantung mana yang cocok dengan `documentID`, karena keduanya berbagi satu module slug.
>
> **Hanya `MISSING_CHECKIN`/`MISSING_CHECKOUT` yang diterapkan otomatis ke session saat approved** — `applyCorrectionToSession` membuat `attendance_event` baru dengan `ValidationStatus = OVERRIDDEN` (raw event lama tetap utuh, sesuai §15) lalu memanggil `recalculateSession`. `WRONG_CHECKIN`/`WRONG_CHECKOUT` **sengaja tidak diterapkan otomatis**: `selectCheckinCheckout` (Phase 6) selalu mengambil CHECKIN pertama dan CHECKOUT pertama-setelahnya di hari itu — menyisipkan event kedua tidak bisa dijamin menggantikan event yang salah tanpa entah mengecualikan event asli dari seleksi (belum ada tracking untuk itu) atau memutasi event asli (dilarang §15). Request untuk kedua tipe ini tetap tercatat dan bisa di-approve untuk keperluan audit, tapi penerapan ke session masih perlu dilakukan manual di luar alur ini sampai logic seleksi diperluas.

---

## Phase 9 - Leave Integration

Integrate:

```text
Approved Leave           ✅ Diimplementasikan — lihat catatan di bawah
        ↓
Attendance Session         ✅
```

Support:

```text
Full Day        ✅ leave_fraction 1.0 → session.status = LEAVE
Half Day          ✅ leave_fraction 0.5 tercatat, status session TIDAK ditimpa jika sudah CLOSED (ada attendance riil)
Hourly Leave       ✅ Leave module sudah menghasilkan `LeaveRequestDetail.DayFraction` yang benar untuk hourly (Leave Phase 3); Attendance hanya menerima fraction apa adanya, tidak perlu logic tambahan
```

> ✅ **Leave Integration — gap nyata, sekarang diimplementasikan.** Sebelumnya field `AttendanceSession.LeaveRequestID`/`LeaveFraction` ada di model tapi tidak pernah diisi kode manapun, dan Leave tidak pernah memanggil Attendance. Sekarang: `leave.AttendanceSessionUpdater` (interface baru di `leave/service.go`) memberi Leave cara memanggil Attendance tanpa import langsung — `attendance.Service` mengimplementasikan `ApplyApprovedLeave(ctx, employeeID, workDate, leaveRequestID, dayFraction)` yang cocok dengan interface tersebut, diwire di `main.go` via `leaveSvc.SetAttendanceSessionUpdater(attendanceSvc)` (arah sebaliknya dari pola `ApprovalEngine` yang biasa — di sini Leave yang mendefinisikan interface dan memanggil Attendance, bukan Attendance memanggil modul lain).
>
> Saat leave request menjadi `APPROVED_FINAL`, `applyAttendanceIntegration` (leave/service.go) membaca `LeaveRequestDetail` (satu baris per tanggal, dengan `DayFraction` yang sudah benar untuk full/half/hourly leave sejak Leave Phase 3) dan memanggil `ApplyApprovedLeave` untuk tiap tanggal. Di sisi Attendance, `ApplyApprovedLeave` (session.go) membuat session baru jika belum ada (sehingga hari yang murni cuti tanpa check-in apapun tetap tercatat, mengatasi keterbatasan Phase 6 yang hanya bereaksi terhadap event), atau meng-update session yang sudah ada — **tapi tidak menimpa status `CLOSED`** (hari dengan attendance nyata), hanya mencatat `LeaveFraction`/`LeaveRequestID`, sesuai contoh half-day leave + half-day attendance di §27.
>
> Kegagalan integrasi ini bersifat best-effort/non-fatal — jika updater belum di-wire atau satu tanggal gagal, di-log dan lanjut ke tanggal berikutnya, tidak menggagalkan approval leave itu sendiri.

---

## Phase 10 - Dashboard & Calendar

Develop:

```text
Employee Dashboard    ✅ Backend selesai — GetEmployeeSummary, lihat catatan di bawah
Manager Dashboard      ⏳ Belum ada — butuh cross-module read employee/organization (siapa yang lapor ke manager ini), belum ada interface untuk itu
HR Dashboard            ⏳ Sama seperti Manager Dashboard — agregasi lintas-organization butuh employee/organization read yang belum ada
Attendance Calendar      ✅ Backend selesai — GetEmployeeCalendar, lihat catatan di bawah
Team Calendar             ⏳ Sama seperti Manager/HR Dashboard — butuh tahu anggota tim/organization
```

> ✅ **Employee Dashboard/Calendar (backend) — diimplementasikan.** `GetEmployeeCalendar` (`GET /attendance/calendar?employee_id=&from=&to=`) mengembalikan sesi dalam rentang tanggal (§36-37). `GetEmployeeSummary` (`GET /attendance/summary?employee_id=&from=&to=`) mengagregasi sesi tersebut menjadi hitungan present/late/missing-checkin/missing-checkout/day-off/leave-days + total work/overtime minutes, sesuai contoh ringkasan di §37.
>
> **Absent sengaja tidak dihitung** di `SummaryResponse` — `SessionStatusAbsent` tidak pernah di-set di manapun di codebase ini (tidak ada scheduled job `ProcessDailyAttendance`/`DetectMissingAttendance`, §44-45), jadi menampilkan hitungan Absent akan selalu 0 dan menyesatkan seolah-olah absensi benar-benar dilacak.
>
> **Manager Dashboard, HR Dashboard, dan Team Calendar sengaja tidak diimplementasikan** — ketiganya butuh mengetahui "siapa bawahan manager ini" atau "siapa anggota organization ini", yaitu cross-module read ke modul employee/organization yang tidak ada interface-nya di manapun di codebase ini (kategori gap yang sama dengan Exempt status di Phase 6 dan employee-active check di Leave). Employee-level dashboard/calendar tidak butuh itu karena `employee_id` sudah diberikan langsung oleh caller.
>
> **Frontend tetap belum dikerjakan** — endpoint backend di atas belum punya UI konsumen; `Attendance.vue` masih placeholder.

---

## Phase 11 - Reports

Develop:

```text
Daily Attendance        ✅ GET /attendance/reports/sessions?from=X&to=X (same date) — lihat catatan di bawah
Monthly Attendance       ✅ Endpoint sama, rentang from/to sebulan
Late                      🔶 Data tersedia (lateness_minutes per session), tapi belum ada endpoint terpisah khusus "Late Report" — filter dilakukan di sisi consumer terhadap hasil /reports/sessions
Absent                    ❌ Tidak bisa dilaporkan — SessionStatusAbsent tidak pernah di-set (lihat Phase 6/10, butuh scheduled job §44-45)
Early Leave               🔶 Sama seperti Late — data (early_leave_minutes) tersedia per session, belum ada endpoint filter khusus
Overtime                  🔶 Data tersedia lewat GET /overtime-requests (sudah ada sejak awal), belum ada filter rentang tanggal khusus "report"
Missing Attendance        🔶 Sama seperti Late — status MISSING_CHECKIN/MISSING_CHECKOUT ada di data /reports/sessions, difilter di consumer
Correction                🔶 Data tersedia lewat GET /corrections (Phase 8), belum ada filter rentang tanggal khusus
Attendance Anomaly        ⏳ Tidak diimplementasikan — "anomaly" butuh definisi bisnis yang belum ada (threshold apa yang dianggap anomali?), bukan sekadar query
```

> ✅ **Tenant-wide session report — gap nyata, sekarang diimplementasikan.** Sebelumnya endpoint session (`GET /sessions`, Phase 10) hanya bisa difilter per-employee, tidak ada cara melihat semua employee dalam satu tanggal/rentang sekaligus. Ditambahkan `FindSessionsInRange` (repository, tanpa filter employee) + `GetAttendanceReport` (service) + `GET /attendance/reports/sessions?from=&to=`. Ini **berbeda dari Manager/HR Dashboard yang sengaja ditunda di Phase 10** — laporan tenant-wide tidak perlu tahu "siapa bawahan siapa", cukup mengambil seluruh session dalam rentang tanggal, jadi tidak terhalang oleh ketiadaan cross-module employee/organization read.
>
> **Late/Early Leave/Missing Attendance tidak diberi endpoint terpisah** — datanya (`lateness_minutes`, `early_leave_minutes`, `status`) sudah ada di setiap baris hasil `/reports/sessions`; menambah endpoint duplikat hanya untuk memfilter kolom yang sama akan jadi surface area berlebih. Overtime dan Correction reports memakai endpoint list yang sudah ada dari Phase 7/8 (belum punya filter rentang tanggal khusus "report", tapi datanya sudah bisa diambil).
>
> **Absent tetap tidak bisa dilaporkan** (sama seperti Phase 10) karena statusnya tidak pernah di-set. **Attendance Anomaly sengaja tidak dikerjakan** — ini bukan gap implementasi teknis, tapi butuh keputusan bisnis dulu (definisi "anomali" apa: pola check-in tidak wajar? lokasi berbeda-beda? belum ada spesifikasi). **Export (Excel/CSV/PDF) tetap belum ada** — itu presentation-layer, di luar cakupan backend murni.

---

## Phase 12 - Notification ✅ Sebagian Selesai (2026-08-08)

Integrate:

```text
Attendance Notification    ⏳ N/A — tidak ada notifikasi umum di luar overtime/correction
Overtime Notification       ✅ Diimplementasikan — lihat catatan di bawah
Correction Notification      ✅ Diimplementasikan — lihat catatan di bawah
Missing Checkout Reminder     ⏳ Tetap belum ada — butuh scheduled job (§44-45), tidak ada infrastruktur cron di codebase ini
```

> ✅ **Modul Notification sekarang ada** (`docs/module-notification-plan.md`, dibangun terpisah Phase 1-5) — blocker infrastruktur lintas-modul yang sebelumnya menahan seluruh Phase 12 ini sudah teratasi. Attendance jadi consumer kedua setelah Leave (Notification Phase 5, rollout).
>
> **Overtime & Correction Notification — diimplementasikan.** `attendance.Notifier` interface baru (`service.go`) + `SetNotifier`, dipenuhi secara struktural oleh `notification.Service` (pola narrow-interface-plus-adapter yang sama dengan `ApprovalEngine`). `attendance.Repository.FindUserIDByEmployeeID` resolusi `employee_id → user_id` lewat `employee_accounts`, sama seperti yang dipakai Leave. `handleOvertimeApprovalStatusChange` dan `handleCorrectionApprovalStatusChange` memanggil `notifyRequestOutcome` pada transisi APPROVED (`OVERTIME_APPROVED`/`CORRECTION_APPROVED`) dan REJECTED (`OVERTIME_REJECTED`/`CORRECTION_REJECTED`) — best-effort: notifier belum wired, employee tanpa user account, atau `Notify` gagal hanya di-log, tidak menggagalkan approval. Diwire di `cmd/server/main.go` via `attendanceSvc.SetNotifier(notificationSvc)`.
>
> **Missing Checkout Reminder tetap belum ada** — bukan lagi soal ketiadaan modul Notification, tapi butuh scheduled job proaktif (§44-45) yang mendeteksi sesi `MISSING_CHECKOUT` tanpa event pemicu apapun; tidak ada infrastruktur cron/scheduler di codebase ini (kategori gap yang sama dengan Absent detection di Phase 6/10-11). **Attendance Notification generik** (di luar overtime/correction) tidak punya trigger konkret lain yang teridentifikasi di plan ini, jadi tidak dibangun sebagai surface area spekulatif.

---

## Phase 13 - Payroll Integration

Expose:

```text
Work Minutes
Overtime Minutes
Absence
Unpaid Leave
Late
Early Leave
```

---

## Phase 14 - Testing

* Unit Tests
* Feature Tests
* Integration Tests
* Calculation Tests
* Approval Tests
* Leave Integration Tests
* Overtime Tests
* Correction Tests
* Payroll Integration Tests

---

# 56. Priority

| Feature                | Priority |
| ---------------------- | -------- |
| Attendance Settings    | P0       |
| Shift                  | P0       |
| Employee Shift         | P0       |
| Check-in               | P0       |
| Check-out              | P0       |
| Attendance Event       | P0       |
| Attendance Session     | P0       |
| Attendance Calculation | P0       |
| Location Validation    | P0       |
| Leave Integration      | P0       |
| Overtime               | P0       |
| Approval Integration   | P0       |
| Attendance Correction  | P1       |
| Face Verification      | P1       |
| Device Management      | P1       |
| Attendance Calendar    | P1       |
| Notification           | P1       |
| Payroll Integration    | P1       |
| Advanced Analytics     | P2       |

---

# 57. Recommended Final Database Structure

```text
attendance_company_settings
        │
        ├── attendance_company_shifts
        │       │
        │       └── attendance_employee_shifts
        │
        ├── attendance_locations
        │
        └── attendance_device_captures
                    │
                    ▼
             attendance_events
                    │
                    ▼
             attendance_sessions
                    │
        ┌───────────┼────────────┐
        ▼           ▼            ▼
      Leave      Overtime    Correction
        │           │            │
        ▼           ▼            ▼
  Leave Module  Approval      Approval
                Module         Module
```

Additional tables jika enhancement diterapkan:

```text
attendance_shift_rules
attendance_employee_devices
attendance_event_corrections
attendance_correction_requests
```

---

# 58. Final Architecture

```text
                        Attendance Module
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
 Configuration          Capture Engine        Calculation
        │                     │                     │
        │                     ▼                     │
        │              Attendance Events            │
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              ▼
                     Attendance Session
                              │
          ┌───────────────────┼───────────────────┐
          │                   │                   │
          ▼                   ▼                   ▼
        Leave             Overtime           Correction
          │                   │                   │
          ▼                   ▼                   ▼
   Leave Management    Approval Module      Approval Module
                              │
                              ▼
                       Notification
                              │
                              ▼
                           Payroll
```

---

# Implementation Status

Diverifikasi langsung terhadap kode per 2026-08-09.

| Phase (§55) | Status | Catatan |
|---|---|---|
| Phase 1 - Database Review & Enhancement | ✅ Selesai (2026-08-08) | Seluruh 10 tabel ada (§2), `approval_instance_id` sudah ada di overtime (migration `063`). Schema/FK/timestamp-type direview — tidak ada bug seperti kasus Leave. Index harian yang kurang di `attendance_events` ditambahkan lewat migration `071_attendance_phase1_event_index` (`idx_att_event_employee_time`). Shift rules table, employee-device mapping, correction table, timezone column — sengaja ditunda, tidak ada kebutuhan konkret saat ini |
| Phase 2 - Attendance Configuration | 🔶 Sebagian (2026-08-08) | Settings/Shifts/Locations/Exempt Positions CRUD lengkap. Shift Rules ditunda (belum ada requirement konkret, sama alasan dengan Phase 1). Devices dan Face Configuration sengaja ditunda — `attendance_device_captures`/`attendance_face_captures` tidak punya repository method sama sekali (tabel mati total), baru masuk akal dibangun bersamaan dengan capture validation engine (Phase 4-5) |
| Phase 3 - Shift Management | 🔶 Sebagian (2026-08-08) | CRUD shift + employee-shift assignment lengkap. Overlap validation (§7) ditemukan benar-benar belum ada — sekarang diperbaiki via `CountOverlappingEmployeeShifts` + validasi `effective_date_from <= effective_date_to` di `CreateEmployeeShift`/`UpdateEmployeeShift`. `DaysOfWeekMask`/`IsCrossMidnight` masih sekadar field pass-through — belum dikonsumsi calculation engine manapun (nunggu Phase 6) |
| Phase 4 - Attendance Capture | 🔶 Sebagian (2026-08-08) | `POST /events` generik (CHECKIN/CHECKOUT satu endpoint, bukan endpoint terpisah — lihat §41). GPS + **Geofence validation kini diimplementasikan** (`geofence.go`, `applyEventValidation`) — event di luar radius jadi `INVALID`. Face Verification & Device Validation sengaja tetap belum ada: tidak ada face-matching provider maupun employee-device mapping (keduanya butuh keputusan/komponen di luar cakupan Phase 4) |
| Phase 5 - Attendance Validation | 🔶 Sebagian (2026-08-08) | Location Validation selesai di Phase 4. **Duplicate Event Detection kini diimplementasikan** (`checkEventSequence` + `FindLastEventForEmployee`) — menolak CHECKIN ganda tanpa CHECKOUT dan CHECKOUT tanpa CHECKIN terbuka. Face/Device Validation tetap belum ada (butuh provider/mapping yang belum ada). Time Validation sengaja ditunda ke Phase 6 karena butuh resolusi shift yang benar (DaysOfWeekMask/cross-midnight) |
| Phase 6 - Attendance Session | 🔶 Sebagian (2026-08-08) | **Gap paling kritis kini teratasi.** `session.go` (`recalculateSession`) menghasilkan/update `attendance_sessions` secara real-time setiap CHECKIN/CHECKOUT — resolusi shift, lateness/early-leave/work-minutes, cross-midnight (work_date = tanggal CHECKIN), DAY_OFF. Belum ada: Absent detection (butuh scheduled job §44-45), Exempt (butuh cross-module read ke employee/organization), Leave integration (Phase 9) |
| Phase 7 - Overtime | ✅ Selesai (2026-08-08) | Approval integration sudah ada sejak awal (Section 29). **Actual/Calculated Overtime kini diimplementasikan**: `applyOvertimeCalculation` dipanggil saat approval, membaca session hari itu untuk `actual_minutes` (aktual checkout vs planned checkout) dan `calculated_minutes` (dibatasi `requested_minutes`), migration `072_attendance_overtime_actual_calculated` menambah kedua kolom. Session juga diupdate dengan `IsOvertimeDay`/`OvertimeMinutes`/dll. **Auto-resolve flow + enum status (2026-08-09)**: `CreateOvertimeRequest` auto-resolve active flow module `"attendance"` jika `flow_id` kosong (lihat §29); migration `079_attendance_overtime_status_pending_approval` menambah `PENDING_APPROVAL` ke enum status MySQL (Postgres sudah VARCHAR). ⏳→🔶 **Rencana dua alur request→isian aktual (SELF & ASSIGNED): lihat §32b — backend selesai (2026-08-10), FE belum** |
| Phase 8 - Correction | 🔶 Sebagian (2026-08-08) | Tabel `attendance_correction_requests` + model + CRUD + approval integration baru dibangun (migration `073`). `HandleApprovalStatusChange` sekarang dispatch overtime vs correction berdasarkan `documentID`. `MISSING_CHECKIN`/`MISSING_CHECKOUT` diterapkan otomatis ke session saat approved (event baru OVERRIDDEN + recalculate). `WRONG_CHECKIN`/`WRONG_CHECKOUT` tercatat & bisa di-approve tapi **tidak** diterapkan otomatis — butuh perluasan logic seleksi checkin/checkout di Phase 6 |
| Phase 9 - Leave Integration | ✅ Selesai (2026-08-08) | `leave.AttendanceSessionUpdater` (interface baru) + `attendance.Service.ApplyApprovedLeave` diwire di `main.go`. Saat leave `APPROVED_FINAL`, tiap `LeaveRequestDetail` mendorong session Attendance jadi `LEAVE` (atau mencatat `LeaveFraction` saja jika session sudah `CLOSED` karena ada attendance nyata, sesuai §27) — termasuk membuat session baru untuk hari yang murni cuti tanpa event apapun |
| Phase 10 - Dashboard & Calendar | 🔶 Sebagian (2026-08-08) | Employee Calendar/Summary (backend) selesai — `GetEmployeeCalendar`/`GetEmployeeSummary` (`GET /attendance/calendar`, `GET /attendance/summary`). Manager/HR Dashboard dan Team Calendar sengaja belum dibangun — butuh cross-module read employee/organization yang belum ada. Frontend tetap belum dimulai |
| Phase 11 - Reports | 🔶 Sebagian (2026-08-08) | `GET /attendance/reports/sessions?from=&to=` (tenant-wide, semua employee) selesai — mencakup Daily/Monthly Attendance. Late/Early Leave/Missing Attendance datanya sudah ada di respons yang sama (difilter di consumer, bukan endpoint terpisah). Overtime/Correction reports pakai endpoint list Phase 7/8 yang sudah ada. Absent tetap tidak bisa (status tidak pernah di-set). Attendance Anomaly & Export (Excel/CSV/PDF) sengaja tidak dikerjakan — masing-masing butuh definisi bisnis / presentation layer di luar cakupan |
| Phase 12 - Notification | 🔶 Sebagian (2026-08-08) | Modul Notification kini ada (`docs/module-notification-plan.md`). `attendance.Notifier` + `SetNotifier` diimplementasikan; `Notify` dipanggil saat Overtime/Correction request APPROVED/REJECTED (`OVERTIME_APPROVED/REJECTED`, `CORRECTION_APPROVED/REJECTED`), best-effort. Missing Checkout Reminder tetap belum ada — butuh scheduled job (§44-45) yang belum ada infrastrukturnya |
| Phase 13 - Payroll Integration | ❌ Belum ada | Tidak ada data untuk diekspos ke Payroll karena session calculation (Phase 6) belum ada |
| Phase 14 - Testing | ✅ Sebagian (2026-08-09) | Test file sudah ada & hijau (`go test ./internal/modules/attendance/`): `approval_integration_test.go` (approval overtime+correction, auto-resolve flow), `notifier_integration_test.go` (OVERTIME/CORRECTION APPROVED/REJECTED), `session_test.go`, `service_test.go`, `repository_test.go`, `handler_test.go`, `helpers_test.go`, `calendar_test.go`, `attendance_integration_test.go` |

**Frontend**: ✅ selesai (FE-1 s.d. FE-5, 2026-08-09) — 12 halaman di `frontend/tenant/src/views/modules/`: `Attendance.vue`, `AttendanceAdmin.vue`, `AttendanceSettings.vue`, `AttendanceShifts.vue`, `AttendanceEmployeeShifts.vue`, `AttendanceLocations.vue`, `AttendanceExemptPositions.vue`, `AttendanceOvertime.vue`, `AttendanceCorrections.vue`, `AttendanceEvents.vue`, `AttendanceSessions.vue`, `AttendanceReports.vue`. Detail di section "Frontend Implementation Plan" di bawah.

**Rekomendasi urutan lanjutan** (blocker struktural utama sudah teratasi per 2026-08-08):
1. ~~Session generation/calculation engine (Phase 6)~~ ✅ Selesai (2026-08-08) — `recalculateSession` sekarang men-generate `attendance_sessions` secara real-time.
2. ~~Overtime actual/calculated minutes (Phase 7)~~ ✅ Selesai (2026-08-08) — `applyOvertimeCalculation` dipanggil saat approval, membaca session hari itu.
3. ~~Correction workflow — Missing Check-in/Checkout (Phase 8)~~ ✅ Selesai (2026-08-08). WRONG_CHECKIN/WRONG_CHECKOUT masih perlu perluasan logic seleksi di Phase 6 sebelum bisa diterapkan otomatis.
4. ~~Leave Integration (Phase 9)~~ ✅ Selesai (2026-08-08) — `leaveSvc.SetAttendanceSessionUpdater(attendanceSvc)`, session ter-update otomatis saat leave disetujui.
5. ~~Employee Dashboard/Calendar backend (Phase 10)~~ ✅ Selesai (2026-08-08) — `GetEmployeeSummary`/`GetEmployeeCalendar`. Manager/HR Dashboard + Team Calendar masih perlu cross-module employee/organization read.
6. ~~Tenant-wide session report (Phase 11)~~ ✅ Selesai (2026-08-08) — `GET /attendance/reports/sessions`. Attendance Anomaly masih butuh definisi bisnis, Export masih presentation-layer.
7. **Payroll Integration (Phase 13)** — sekarang jadi kandidat berikutnya yang paling bernilai: session sudah punya `WorkMinutes`/`OvertimeMinutes`/`LeaveFraction` lengkap untuk diekspos ke Payroll, tinggal endpoint/query agregasi per periode.
8. Scheduled job untuk Absent/Missing detection (§44-45) — perlu keputusan infra (cron/scheduler) yang belum ada polanya di codebase ini.
9. ~~Modul Notification terpusat (Phase 12)~~ ✅ Selesai sebagian (2026-08-08) — Overtime/Correction notification terwire. Missing Checkout Reminder masih menunggu scheduled job (§44-45).
10. Dedicated check-in/check-out endpoints (pisah dari `POST /events` generik) — supaya validasi per-aksi (Phase 4-5) punya tempat spesifik dipasang.
11. ~~Frontend dasar~~ ✅ Selesai (2026-08-09) — FE-1 s.d. FE-5; data session yang ditampilkan sekarang benar-benar berarti.

---

# Frontend Implementation Plan ✅ Selesai (FE-1 s.d. FE-5, 2026-08-09)

Ditambahkan 2026-08-08. Backend Attendance sudah selesai >90% (lihat Implementation Status di atas) — sisi frontend awalnya 0%: `frontend/tenant/src/views/modules/Attendance.vue` hanya placeholder satu baris ("Attendance Module — Coming soon"), tidak ada komponen atau pemanggilan API apapun. Sidebar entry (`layouts/Sidebar.vue:346`, `moduleSlug: 'attendance'`, `permission: 'attendance.view'`) dan router stub sudah ada sejak awal. Section ini adalah rencana untuk mengisi gap tersebut, mengikuti konvensi FE yang sudah baku di repo ini — bukan pola baru.

> ✅ **Seluruh 5 phase FE (FE-1 s.d. FE-5) sudah diimplementasikan per 2026-08-09.** 11 halaman baru: `Attendance.vue` (My Dashboard, diisi), `AttendanceAdmin.vue`, `AttendanceSettings.vue`, `AttendanceShifts.vue`, `AttendanceEmployeeShifts.vue`, `AttendanceLocations.vue`, `AttendanceExemptPositions.vue`, `AttendanceOvertime.vue`, `AttendanceCorrections.vue`, `AttendanceEvents.vue`, `AttendanceSessions.vue`, `AttendanceReports.vue`. Detail per-phase ada di masing-masing blockquote FE-1 s.d. FE-5 di bawah. Yang tetap di luar cakupan (lihat "Eksplisit di luar cakupan rencana FE ini" di bawah §FE-3.5): Manager/HR Dashboard, Absent detection di UI, notification bell wiring, dan export Excel/CSV/PDF — semuanya blocked oleh gap backend atau modul lain, bukan sesuatu yang FE bisa selesaikan sendiri.

## FE-1. Ringkasan & Prinsip

* **Tidak ada dokumen konvensi FE terpisah** — pola diambil dari modul tenant yang sudah punya UI nyata: `views/modules/Employees.vue` + `EmployeeForm.vue` (list+form+pagination+search — template utama), `Organizations.vue`, dan `Approvals.vue`/`ApprovalFlows.vue` (untuk halaman dengan state/tab lebih kompleks). **`Leave.vue` dan `Payroll.vue` masih placeholder sama seperti Attendance** — jangan dipakai sebagai referensi meskipun secara domain lebih mirip.
* **Tidak ada service-layer/Pinia store per modul.** Komponen memanggil `services/api.js` (axios instance dengan interceptor auth/tenant) langsung dengan path REST mentah, mis. `api.get('/api/v1/tenant/employees', { params })`. Error ditangani lewat `services/responseHandler.js` (`getMessage`, `getErrorMessage`, `getValidationErrors`, `isValidationError`). Attendance mengikuti pola yang sama — tidak membuat `attendanceStore.js` atau `attendanceApi.js` terpisah.
* **State cukup `ref`/`reactive` lokal per view** (`<script setup>`), bukan Pinia store baru — Pinia di repo ini hanya untuk concern lintas-modul (`stores/auth.js`, `stores/activeModules.js`, dll.).
* **Approval Overtime/Correction TIDAK dibangun ulang di Attendance.** Backend approval sepenuhnya lewat Central Approval Module (module slug `"attendance"` generik, endpoint approve/reject ada di modul Approval, bukan di `attendance/routes.go`). FE Attendance hanya perlu link-out ke halaman detail approval instance yang sudah ada (pola sama seperti yang seharusnya dilakukan Leave), bukan tombol approve/reject sendiri.
* **Permission gating** pakai composable yang sudah ada — `useAuth().hasPermission(slug)` (`stores/auth.js`), dipakai inline `v-if="hasPermission('attendance.create')"`. Tidak ada directive `v-can`. Slug permission Attendance yang tersedia (`backend/internal/modules/attendance/module.go`): `attendance.view`, `attendance.create`, `attendance.update`, `attendance.delete` — **tidak ada** `attendance.approve` terpisah, karena approve lewat modul Approval.
* **Tidak ada link otomatis backend↔frontend.** `moduleSlug`/`permission` di router & sidebar FE adalah string yang di-hardcode manual, harus dicocokkan sendiri dengan `Info().Slug`/`Info().Permissions` backend — bukan digenerate dari manifest manapun. Setiap halaman baru di bawah wajib memakai slug modul `'attendance'` yang sama persis.

## FE-2. Halaman & Routing

Semua route baru sebagai sibling di bawah path `attendance/...`, `meta.module: 'attendance'`, kebab-case path, PascalCase `name`, `meta.titleKey`/`descKey` untuk i18n, sub-halaman pakai `meta.backRoute`/`backLabelKey` mengikuti konvensi routing yang sudah ada di `router/index.js`. Router guard yang sudah ada (`router/index.js:330-346`, cek `useActiveModules().hasModule(meta.module)`) otomatis berlaku tanpa perubahan.

| Halaman | Route (usulan) | Endpoint backend | Permission |
|---|---|---|---|
| `Attendance.vue` (My Attendance dashboard: ringkasan + kalender + tombol check-in/out) | `attendance` (existing stub, diisi) | `GET /calendar`, `GET /summary`, `POST /events` | `attendance.view` (create implisit via check-in/out sebagai aksi diri sendiri) |
| `AttendanceSettings.vue` | `attendance/settings` | `GET/PUT /settings` | `attendance.update` |
| `AttendanceShifts.vue` + `AttendanceShiftForm.vue` | `attendance/shifts` | `POST/GET /shifts`, `GET/PUT/DELETE /shifts/:id` | `attendance.view`/`create`/`update`/`delete` |
| `AttendanceEmployeeShifts.vue` + form | `attendance/employee-shifts` | `POST/GET /employee-shifts`, `GET/PUT/DELETE /employee-shifts/:id` | sama pola di atas |
| `AttendanceLocations.vue` + form | `attendance/locations` | `POST/GET /locations`, `GET/PUT/DELETE /locations/:id` | sama pola di atas |
| `AttendanceExemptPositions.vue` + form | `attendance/exempt-positions` | `POST/GET /exempt-positions`, `GET/PUT/DELETE /exempt-positions/:id` | sama pola di atas |
| `AttendanceEvents.vue` (read-only, audit) | `attendance/events` | `GET /events`, `GET /events/:id` | `attendance.view` |
| `AttendanceSessions.vue` (read-only, per-employee) | `attendance/sessions` | `GET /sessions`, `GET /sessions/detail` | `attendance.view` |
| `AttendanceOvertime.vue` + `AttendanceOvertimeForm.vue` | `attendance/overtime` | `POST/GET /overtime-requests`, `GET /overtime-requests/:id` | `attendance.view`/`create` |
| `AttendanceCorrections.vue` + form | `attendance/corrections` | `POST/GET /corrections`, `GET /corrections/:id` | `attendance.view`/`create` |
| `AttendanceReports.vue` (tenant-wide, admin/HR) | `attendance/reports` | `GET /reports/sessions` | `attendance.view` (direkomendasikan dibatasi lebih lanjut di FE lewat role/permission HR jika tersedia) |

Catatan: `POST /events` adalah endpoint generik untuk CHECKIN & CHECKOUT (`event_type` di body), bukan dua endpoint terpisah — FE cukup satu handler check-in/out yang menentukan `event_type` berdasarkan status sesi terakhir (dari `GET /summary` atau `GET /calendar` hari berjalan).

## FE-3. Development Phases (FE)

Mengikuti gaya granular §55 (backend), diurutkan berdasarkan nilai bagi pengguna dan ketergantungan data:

**Phase FE-1 — My Attendance Dashboard ✅ Selesai (2026-08-08)**
Halaman utama employee: ringkasan (present/late/missing/leave-days dari `GetEmployeeSummary`), kalender bulan berjalan (`GetEmployeeCalendar`), tombol Check-in/Check-out dengan capture GPS browser (`navigator.geolocation`) yang mengirim `POST /events`. Nilai tertinggi untuk end user biasa — tanpa ini modul Attendance FE benar-benar 0% dipakai sehari-hari.

> ✅ **Diimplementasikan.** `frontend/tenant/src/views/modules/Attendance.vue` (diisi, sebelumnya placeholder): kartu ringkasan (present/late/missing-checkout/leave-days), daftar sesi bulan berjalan dengan status tag, dan tombol Check-in/Check-out yang menentukan `event_type` dari status sesi hari ini (`OPEN` → tombol Check-out, selain itu → Check-in) lalu mengirim `POST /attendance/events` dengan GPS dari `navigator.geolocation`. Menggunakan `services/api.js`/`responseHandler.js` yang sama dengan modul lain — tidak ada service-layer/store baru, sesuai FE-1 di atas.
>
> **Blocker yang ditemukan saat implementasi dan cara mengatasinya**: `GetEmployeeSummary`/`GetEmployeeCalendar`/`CreateEvent` semuanya butuh `employee_id` sebagai parameter dari client — tidak ada resolusi otomatis dari user yang login di modul manapun (`performance`'s `GetCurrentEmployeeContextByUserID` adalah pola serupa tapi khusus modul performance, tidak reusable). Tanpa ini, dashboard tidak bisa tahu "employee_id saya sendiri". Ditambahkan endpoint baru **`GET /api/v1/tenant/user-accounts/me`** (`useraccount` module — `Repository.FindAccountByUserID`, `Service.GetMyAccount`, `Handler.GetMyAccount`) yang resolusi `employee_accounts` lewat `authctx.GetUserID(ctx)`, arah sebaliknya dari `FindAccountByEmployeeID` yang sudah ada. Ini bukan endpoint khusus Attendance — sengaja diletakkan di `useraccount` karena modul itu yang memiliki tabel `employee_accounts`, agar fitur self-service modul lain di masa depan (mis. Leave "My Requests") bisa reuse endpoint yang sama tanpa duplikasi logic resolusi.
>
> Belum ada test otomatis untuk `GetMyAccount` — package `useraccount` belum punya test harness (`setupTestDB`) sama sekali di codebase ini sebelum perubahan ini, jadi menambah satu test untuk method ini akan berarti membangun seluruh harness baru di luar cakupan Phase FE-1; diverifikasi manual lewat `go build`/`go vet` saja. Revisit jika `useraccount` mendapat test harness di kesempatan lain.

**Phase FE-2 — Admin Configuration ✅ Selesai (2026-08-08)**
`AttendanceSettings`, `AttendanceShifts` (+form), `AttendanceEmployeeShifts` (+form, termasuk validasi overlap yang sudah dilakukan backend — tampilkan error validasi apa adanya lewat `isValidationError`/`getValidationErrors`), `AttendanceLocations` (+form, mis. input lat/long manual atau pick-on-map jika ada komponen map di repo — cek dulu sebelum menambah dependency baru), `AttendanceExemptPositions` (+form). Semua CRUD standar mengikuti pola `Employees.vue`.

> ✅ **Diimplementasikan.** 6 halaman baru di `frontend/tenant/src/views/modules/`: `AttendanceAdmin.vue` (index kartu, pola sama dengan `SettingsIndex.vue`), `AttendanceSettings.vue` (form tunggal GET/PUT, tanpa list), `AttendanceShifts.vue`, `AttendanceLocations.vue`, `AttendanceExemptPositions.vue` (masing-masing list + Dialog create/edit, mengikuti pola compact `NationalitiesView.vue` — DataTable server-pagination + `Dialog` inline, bukan route form terpisah seperti `EmployeeForm.vue` yang khusus untuk wizard kompleks), dan `AttendanceEmployeeShifts.vue` (list + Dialog, dengan dropdown employee/shift dan `DateInput` untuk `effective_date_from`/`to`). Tombol "Admin" (gated `hasPermission('attendance.update')`) ditambahkan di header `Attendance.vue` menuju `/attendance/admin`. Routing: 6 route baru di `router/index.js` di bawah `attendance/...` dengan `meta.backRoute` mengarah balik ke `/attendance/admin`.
>
> **Pilihan pola: Dialog inline, bukan `<Feature>Form.vue` terpisah** — plan FE-2 awal menyebut "+form" secara generik; setelah diperiksa, precedent CRUD sederhana yang sebenarnya paling umum di codebase ini adalah pola compact satu-file (`NationalitiesView.vue`, `ReligionsView.vue`, dll — DataTable + `Dialog` inline), bukan route/file form terpisah seperti `EmployeeForm.vue` (yang khusus untuk wizard 11-step). Shift/Location/ExemptPosition/EmployeeShift semuanya entity sederhana (4-6 field), jadi mengikuti pola Dialog inline lebih konsisten dengan mayoritas modul lain dan menghindari proliferasi file untuk form yang trivial.
>
> **Overlap validation (§7)** tidak diberi UI khusus di FE — backend (`CountOverlappingEmployeeShifts`, Phase 3 backend) sudah menolak assignment yang overlap lewat `VALIDATION_ERROR`; FE `AttendanceEmployeeShifts.vue` cukup menampilkan pesan error itu apa adanya lewat `getValidationErrors`, tidak mencoba mendeteksi ulang overlap di client.
>
> ✅ **Map picker untuk lokasi — diimplementasikan (2026-08-09).** `AttendanceLocations.vue` memakai **Leaflet** (import `leaflet` + `leaflet/dist/leaflet.css` + marker assets): peta interaktif di dalam Dialog form, klik peta mengisi `latitude`/`longitude`, lingkaran `L.circle` menggambarkan jangkauan `radius_m`, plus **pencarian alamat via Nominatim** (`https://nominatim.openstreetmap.org/search` — debounce, dropdown hasil dengan keyboard nav, hasil tampil di atas peta, z-index dropdown diperbaiki agar tidak tertutup peta). Modal lokasi diperlebar agar peta nyaman. Keputusan awal "tanpa map" diubah setelah kebutuhan klik lat/long + visualisasi radius menjadi requirement eksplisit.
>
> **Sidebar tetap satu entry** (`Sidebar.vue:346`, tidak diubah) — `operationsItems` di sidebar tidak mendukung dropdown children (`coreHRItems`/`talentItems` yang mendukung), jadi navigasi ke halaman admin baru dilakukan lewat tombol "Admin" di dalam `Attendance.vue` menuju index kartu `AttendanceAdmin.vue`, bukan lewat perubahan struktur sidebar — konsisten dengan pola `/settings` yang juga satu sidebar entry menuju index kartu.

**Phase FE-3 — Overtime & Correction Requests ✅ Selesai (2026-08-09)**
`AttendanceOvertime`/`AttendanceCorrections` (+form masing-masing): create request + list status (SUBMITTED/PENDING_APPROVAL/APPROVED/REJECTED). Approve/reject **tidak dibangun di sini** — tampilkan status + link-out ke halaman detail approval instance module Approval yang sudah ada.

> ✅ **Diimplementasikan.** `AttendanceOvertime.vue` (list + Dialog create, field `work_date`/`start_time`/`end_time`/`requested_minutes`/`reason`, kolom `calculated_minutes` ditampilkan read-only dari hasil `applyOvertimeCalculation` backend) dan `AttendanceCorrections.vue` (list + Dialog create, 4 `correction_type` dari model backend). Keduanya difilter `employee_id` milik user yang login (lewat `useMyEmployee` composable, lihat di bawah) — halaman ini murni self-service, bukan admin view lintas-employee. **Enrichment submitter (2026-08-09)**: backend `ListOvertimeRequests` selalu menambahkan `employee_name`/`organization_name` via `Repository.GetEmployeeInfoByIDs` (batch lookup per halaman, tanpa N+1) — kolom nama karyawan & organisasi tampil di tabel lembur; berguna untuk list tenant-wide (admin) dan tetap terisi untuk self-only. Tombol "View in Approvals" mengarah ke `/approvals` (halaman Approval Module yang sudah ada) — **bukan deep-link ke instance spesifik**, karena `Approvals.vue` tidak punya dukungan query-param untuk membuka satu instance langsung; membangun itu di luar cakupan plan Attendance FE ini (milik Approval module).
>
> **`AttendanceCorrections.vue` me-resolve `attendance_session_id` dari tanggal** — backend `CreateCorrectionRequest` mewajibkan `attendance_session_id`, bukan `work_date` mentah, jadi form memanggil `GET /attendance/sessions/detail?employee_id=&work_date=` saat tanggal dipilih untuk mendapatkan session yang sesuai (menampilkan pesan "sesi tidak ditemukan" jika tidak ada). `WRONG_CHECKIN`/`WRONG_CHECKOUT` tetap bisa diajukan lewat form ini untuk keperluan audit (sesuai catatan Phase 8 backend), meskipun backend belum menerapkannya otomatis ke session.
>
> **Dua file baru diekstrak** karena sekarang dipakai di 3 tempat (Attendance.vue, AttendanceOvertime.vue, AttendanceCorrections.vue) — bukan duplikasi lagi:
> - `composables/useMyEmployee.js`: membungkus `GET /user-accounts/me` dengan cache module-level (sekali per sesi browser), menggantikan fetch inline yang sebelumnya cuma ada di `Attendance.vue`.
> - `utils/localTime.js` (`toLocalISOString`, `localDateTimeISOString`): logic ISO+offset yang sebelumnya inline di `Attendance.vue` untuk `event_time_local`, sekarang dipakai ulang untuk `start_time_local`/`end_time_local` (Overtime) dan `requested_checkin`/`requested_checkout` (Correction) — semuanya di-parse backend sebagai RFC3339 (`time.Parse(time.RFC3339, ...)`), butuh offset eksplisit, bukan sekadar `YYYY-MM-DDTHH:mm:ss`.
>
> Route baru: `attendance/overtime`, `attendance/corrections` (sibling `/attendance`, bukan di bawah `/attendance/admin` — kedua halaman ini untuk semua employee, bukan admin-only). Tombol akses ditambahkan di header `Attendance.vue` di samping tombol "Admin".

**Phase FE-4 — Events & Sessions (Read-only Audit Views) ✅ Selesai (2026-08-09)**
`AttendanceEvents`, `AttendanceSessions`: tabel read-only untuk audit/troubleshooting HR (raw event log, session detail per employee/tanggal). Prioritas lebih rendah dari FE-1/2/3 karena bukan alur kerja harian.

> ✅ **Diimplementasikan.** `AttendanceEvents.vue` (tabel raw event: employee, event_type, event_time_local, validation_status, distance_m, is_in_geofence) dan `AttendanceSessions.vue` (tabel session: employee, work_date, status, lateness/early-leave/work/overtime minutes) — keduanya tenant-wide (bukan "punya saya sendiri" seperti Overtime/Correction FE-3), dengan dropdown filter employee opsional (kosong = semua employee, sesuai perilaku backend `ListEvents`/`ListSessions` saat `employee_id` kosong). Ditambahkan sebagai 2 kartu baru di `AttendanceAdmin.vue` (bukan di dashboard utama seperti Overtime/Correction) — audit tooling ini untuk HR/admin, konsisten dengan penempatan Settings/Shifts/dll.
>
> Label status session (`status_open`/`status_closed`/dst.) **reuse key i18n yang sudah ada dari Phase FE-1** (`attendance.status_*`), tidak diduplikasi.
>
> Route baru: `attendance/events`, `attendance/sessions`, sibling di bawah `/attendance/admin` dengan `backRoute` kembali ke situ.

**Phase FE-5 — Tenant-wide Reports ✅ Selesai (2026-08-09)**
`AttendanceReports`: `GET /reports/sessions?from=&to=` dengan filter rentang tanggal, tabel semua employee. Late/Early Leave/Missing Attendance ditampilkan sebagai kolom dari respons yang sama (bukan endpoint terpisah — backend memang tidak menyediakannya, lihat Phase 11 backend).

> ✅ **Diimplementasikan — sekaligus phase FE terakhir, seluruh Frontend Implementation Plan (FE-1 s.d. FE-5) sekarang selesai.** `AttendanceReports.vue`: filter `from`/`to` (`DateInput`, default bulan berjalan), tabel semua employee dengan kolom yang sama seperti `AttendanceSessions.vue` (status, lateness/early-leave/work/overtime minutes).
>
> **Pagination client-side, bukan server-side** — beda dari halaman list lain di modul ini. `GET /reports/sessions` (`service.GetAttendanceReport`) mengembalikan `[]SessionResponse` langsung tanpa amplop `page`/`total`/paginasi (lihat `service.go:716-726` — murni `FindSessionsInRange` tanpa limit/offset), jadi `AttendanceReports.vue` memakai `DataTable` `paginator` bawaan PrimeVue (client-side, tanpa `lazy`) alih-alih pola `lazy` + `@page` yang dipakai `AttendanceEvents`/`AttendanceSessions`/dll. Ini bukan penyimpangan dari konvensi — backend memang tidak menyediakan paginasi untuk endpoint ini.
>
> Ditambahkan sebagai kartu ke-6 di `AttendanceAdmin.vue`, route `attendance/reports`.

**Eksplisit di luar cakupan rencana FE ini:**
* **Manager Dashboard, HR Dashboard, Team Calendar** — backend belum ada (blocked oleh cross-module employee/organization read yang belum ada interface-nya, sama seperti dicatat di backend Phase 10). Tidak ada FE yang bisa dibangun sampai backend-nya ada.
* **Absent detection di UI mana pun** — status `ABSENT` tidak pernah di-set backend (butuh scheduled job §44-45), jadi FE tidak boleh menampilkan hitungan/badge Absent yang akan selalu 0 dan menyesatkan.
* **Notification bell wiring** — bell icon sudah ada secara kosmetik (`layouts/HeaderBar.vue:146-156`, badge hardcoded "3", belum ada dropdown/API). Attendance sekarang memicu notifikasi (`OVERTIME_APPROVED/REJECTED`, `CORRECTION_APPROVED/REJECTED` — lihat Phase 12 di atas), tapi mewujudkan bell jadi dropdown fungsional adalah scope `docs/module-notification-plan.md`, bukan bagian dari plan Attendance FE ini.
* **Export Excel/CSV/PDF** pada Reports — presentation-layer tambahan, di luar cakupan pass pertama FE.

## FE-4. Catatan Teknis

* **Response envelope**: seluruh endpoint pakai `httputil.SuccessJSON`/pola paginated response yang sama dengan modul lain — FE harus mem-parsing respons dengan cara yang sama seperti `Employees.vue` menangani `api.get(...)`, jangan berasumsi shape berbeda untuk Attendance.
* **Tabel & pagination**: PrimeVue `DataTable` dengan `lazy` + server-side pagination (`:totalRecords`, `:first`, `:rows`, `@page`), loading state pakai komponen `SkeletonTable` yang sudah ada (bukan built-in loading DataTable), search pakai debounce `setTimeout` — semua meniru `Employees.vue`.
* **Form & validasi**: tidak ada vee-validate/yup/zod di dependency — validasi manual + error dari backend ditangkap lewat try/catch dan `responseHandler.js`. PrimeVue input components dengan `v-model`.
* **GPS capture** (Phase FE-1) adalah kebutuhan baru yang belum ada polanya di modul lain manapun di FE ini — perlu `navigator.geolocation.getCurrentPosition` langsung di komponen, tidak ada composable existing untuk itu; kalau dipakai di lebih dari satu tempat (mis. juga dibutuhkan saat submit correction dengan lokasi), pertimbangkan composable `useGeolocation()` kecil saat itu terjadi, bukan dibuat spekulatif di awal.

---

# 59. Design Principles

1. `attendance_events` merupakan raw attendance event.
2. Raw event tidak boleh diubah untuk menghilangkan histori.
3. `attendance_sessions` merupakan hasil kalkulasi attendance.
4. Satu employee hanya memiliki satu session untuk satu `work_date`.
5. Shift menentukan planned working time.
6. Cross-midnight shift harus tetap memiliki satu `work_date`.
7. Leave yang approved dapat mengubah status session menjadi `LEAVE`.
8. Overtime harus menggunakan Central Approval Module.
9. Attendance Correction juga menggunakan Central Approval Module.
10. Attendance tidak memiliki approval engine sendiri.
11. Location validation dilakukan sebelum event menjadi valid.
12. Face validation dilakukan jika diwajibkan oleh configuration.
13. Device dapat divalidasi jika perusahaan menggunakan device restriction.
14. Attendance calculation harus dipisahkan dari controller.
15. Attendance harus dapat direcalculate.
16. Semua perubahan penting harus dapat diaudit.
17. Balance/quota bukan tanggung jawab Attendance; Leave Management yang bertanggung jawab.
18. Nominal payroll bukan tanggung jawab Attendance; Payroll yang menghitung.
19. Attendance menyediakan data faktual untuk Payroll.
20. Seluruh ID menggunakan UUID.
21. Current Attendance Session dapat diperbarui, tetapi raw event dan histori koreksi harus tetap terlacak.
22. Configuration dapat berbeda sesuai kebutuhan Company/Organization.
23. Attendance harus dapat digunakan oleh Employee, Manager, dan HR dengan hak akses berbeda.
24. Modul harus dapat dikembangkan untuk biometric, GPS, device, dan metode attendance lainnya tanpa mengubah konsep inti `attendance_events`.
