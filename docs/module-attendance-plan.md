> ⚠️ **Status vs. Plan ini**: dokumen ini ditulis seolah modul Attendance belum ada (greenfield). Setelah dicek ulang terhadap kode aktual, **seluruh tabel data (§2-§14) dan integrasi approval Overtime (§29-§30) sudah diimplementasikan sepenuhnya** — sama seperti pola Leave/Payroll. Namun **seluruh calculation/processing engine yang menjadi inti modul ini (session generation, lateness/overtime calculation, capture validation, integrasi Leave, correction workflow) belum ada sama sekali** — hanya CRUD tipis di atas tabel mentah. Ini gap yang lebih besar dibanding yang ditemukan di modul Leave. **Frontend masih placeholder "Coming soon"**. Lihat section **"Implementation Status"** di bagian bawah dokumen untuk status per-fase yang sudah diverifikasi terhadap kode, dan catatan blockquote (`>`) di beberapa section untuk koreksi spesifik.

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

> ❌ **Belum ada sama sekali.** Tidak ada tabel/model/migration/endpoint untuk `attendance_event_corrections` maupun `attendance_correction_requests` di manapun dalam codebase. Section 33-34 (Attendance Correction flow, tabel `attendance_correction_requests`) di bawah seluruhnya masih proposal murni, bukan sebagian terbangun.

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

> ❌ **Belum ada sama sekali.** `Service.CreateEvent` (`attendance/service.go:396-432`) hanya mem-parsing input dan menyimpan raw event dengan `ValidationStatus: PENDING` — tidak ada perhitungan jarak/geofence (tidak ada fungsi haversine/distance di manapun di `service.go`), tidak ada pemanggilan face verification, tidak ada validasi device. Field `DistanceM`/`IsInGeofence`/`ValidatedLocationID`/`FaceCaptureID` sudah ada di model tapi tidak pernah diisi oleh service — hanya bisa diisi manual oleh client atau tetap kosong. Section 18 di bawah ini seluruhnya masih proposal.

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

> ❌ **Gap paling kritis di modul ini, lebih besar dari gap manapun di Leave.** Tabel `attendance_sessions` sudah ada dan sudah bisa dibaca (`GetSession`/`ListSessions`), tapi **tidak ada satupun kode yang menulis ke tabel ini**. `CreateEvent` hanya menyimpan raw event, tidak pernah membuat/update `AttendanceSession`. Tidak ada `GenerateSession`/`CloseSession`/session-builder service di manapun dalam modul (hanya `service.go`, `repository.go`, `handler.go`, `dto.go`, `module.go`, `routes.go` — tidak ada file "calculation" terpisah). Tidak ada scheduled job untuk `ProcessDailyAttendance` (§44) atau `DetectMissingAttendance` (§45) — digrep di seluruh `backend/` dan tidak ditemukan. Section 20-24 di bawah (session processing, lateness, early leave, work minutes, cross-midnight) seluruhnya masih proposal murni. Karena Overtime (§29), Leave Integration (§26-27/§50), dan Payroll Integration (§49) semuanya bergantung pada session yang benar-benar ter-generate, ini harus jadi prioritas #1 sebelum fase-fase lain bisa berjalan penuh.

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

---

# 26. Leave Integration

> ❌ **Belum ada di kedua arah.** Field `AttendanceSession.LeaveRequestID`/`LeaveFraction` dan konstanta `SessionStatusLeave = "LEAVE"` sudah ada di model, tapi tidak ada kode yang mengisinya. Digrep, modul `leave/*.go` tidak punya referensi apapun ke attendance, dan modul `attendance/*.go` tidak punya referensi apapun ke leave. Ini konsisten dengan temuan dari sisi modul Leave sendiri (`docs/module-leave-plan.md` §24 Attendance Integration — juga ditandai belum ada). Section 27 (Leave Fraction) di bawah tetap proposal murni.

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
> Catatan: "Actual Overtime" vs "Calculated Overtime" (§31-32) **belum diimplementasikan** — itu bagian dari session calculation engine yang belum ada (lihat Section 19). Approval instance-nya sudah jalan, tapi perhitungan menit overtime aktual berdasarkan attendance masih proposal.

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

> ❌ **Frontend belum dibangun.** `frontend/tenant/src/views/modules/Attendance.vue` hanya satu baris placeholder ("Attendance Module — Coming soon"), tidak ada komponen atau pemanggilan API. Struktur menu (Dashboard, Shifts, Schedules, Events, Overtime, Locations, Settings) memang sudah didefinisikan di sisi server (`module.go:76-91`), tapi tidak ada halaman frontend yang cocok dengan itu. Section 36-40 (Calendar, Employee/Manager/HR Dashboard) di bawah ini 0% dikerjakan.

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
> Yang **belum ada**: endpoint `check-in`/`check-out` khusus (§41 di bawah), endpoint session write, endpoint corrections (§34), endpoint reports (§40), dan endpoint approve/reject overtime khusus (approval sepenuhnya lewat endpoint generik Central Approval Module — lihat Section 29).

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
* Add correction request. ⏳ Deferred — seluruh correction workflow (§16, §33-34) belum mulai dikerjakan, revisit saat Phase 8 dimulai.
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
Location Validation
Face Validation
Device Validation
Time Validation
Duplicate Event Detection
```

---

## Phase 6 - Attendance Session

Develop:

```text
Session Generation
Check-in Mapping
Check-out Mapping
Work Minutes
Late
Early Leave
Day Off
Absent
Leave
Exempt
```

---

## Phase 7 - Overtime

Develop:

```text
Overtime Request
Central Approval Integration
Approved Overtime
Actual Overtime
Calculated Overtime
Attendance Integration
```

---

## Phase 8 - Correction

Develop:

```text
Missing Check-in
Missing Checkout
Attendance Correction
Central Approval
Session Recalculation
Audit
```

---

## Phase 9 - Leave Integration

Integrate:

```text
Approved Leave
        ↓
Attendance Session
```

Support:

```text
Full Day
Half Day
Hourly Leave
```

---

## Phase 10 - Dashboard & Calendar

Develop:

```text
Employee Dashboard
Manager Dashboard
HR Dashboard
Attendance Calendar
Team Calendar
```

---

## Phase 11 - Reports

Develop:

```text
Daily Attendance
Monthly Attendance
Late
Absent
Early Leave
Overtime
Missing Attendance
Correction
Attendance Anomaly
```

---

## Phase 12 - Notification

Integrate:

```text
Attendance Notification
Overtime Notification
Correction Notification
Missing Checkout Reminder
```

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

Diverifikasi langsung terhadap kode per 2026-08-08.

| Phase (§55) | Status | Catatan |
|---|---|---|
| Phase 1 - Database Review & Enhancement | ✅ Selesai (2026-08-08) | Seluruh 10 tabel ada (§2), `approval_instance_id` sudah ada di overtime (migration `063`). Schema/FK/timestamp-type direview — tidak ada bug seperti kasus Leave. Index harian yang kurang di `attendance_events` ditambahkan lewat migration `071_attendance_phase1_event_index` (`idx_att_event_employee_time`). Shift rules table, employee-device mapping, correction table, timezone column — sengaja ditunda, tidak ada kebutuhan konkret saat ini |
| Phase 2 - Attendance Configuration | 🔶 Sebagian (2026-08-08) | Settings/Shifts/Locations/Exempt Positions CRUD lengkap. Shift Rules ditunda (belum ada requirement konkret, sama alasan dengan Phase 1). Devices dan Face Configuration sengaja ditunda — `attendance_device_captures`/`attendance_face_captures` tidak punya repository method sama sekali (tabel mati total), baru masuk akal dibangun bersamaan dengan capture validation engine (Phase 4-5) |
| Phase 3 - Shift Management | 🔶 Sebagian (2026-08-08) | CRUD shift + employee-shift assignment lengkap. Overlap validation (§7) ditemukan benar-benar belum ada — sekarang diperbaiki via `CountOverlappingEmployeeShifts` + validasi `effective_date_from <= effective_date_to` di `CreateEmployeeShift`/`UpdateEmployeeShift`. `DaysOfWeekMask`/`IsCrossMidnight` masih sekadar field pass-through — belum dikonsumsi calculation engine manapun (nunggu Phase 6) |
| Phase 4 - Attendance Capture | 🔶 Sebagian (2026-08-08) | `POST /events` generik (CHECKIN/CHECKOUT satu endpoint, bukan endpoint terpisah — lihat §41). GPS + **Geofence validation kini diimplementasikan** (`geofence.go`, `applyEventValidation`) — event di luar radius jadi `INVALID`. Face Verification & Device Validation sengaja tetap belum ada: tidak ada face-matching provider maupun employee-device mapping (keduanya butuh keputusan/komponen di luar cakupan Phase 4) |
| Phase 5 - Attendance Validation | ❌ Belum ada | Tidak ada location/face/device/time validation logic di `service.go` — lihat Section 17 |
| Phase 6 - Attendance Session | ❌ Belum ada | **Gap paling kritis.** Tidak ada session generation/calculation engine sama sekali — lihat Section 19 |
| Phase 7 - Overtime | 🔶 Sebagian | Approval integration ke Central Approval Module sudah selesai (Section 29). "Actual Overtime"/"Calculated Overtime" berdasarkan attendance (§31-32) belum ada karena bergantung pada session calculation (Phase 6) |
| Phase 8 - Correction | ❌ Belum ada | Tidak ada tabel/model/endpoint correction — lihat Section 16 |
| Phase 9 - Leave Integration | ❌ Belum ada | Field `LeaveRequestID`/`LeaveFraction` ada di model tapi tidak diisi kode manapun — lihat Section 26 |
| Phase 10 - Dashboard & Calendar | ❌ Belum ada | Tidak ada endpoint dashboard/calendar; frontend hanya stub |
| Phase 11 - Reports | ❌ Belum ada | Tidak ada endpoint report apapun |
| Phase 12 - Notification | ❌ Belum ada | Tidak ditemukan pemanggilan Notification module dari modul Attendance |
| Phase 13 - Payroll Integration | ❌ Belum ada | Tidak ada data untuk diekspos ke Payroll karena session calculation (Phase 6) belum ada |
| Phase 14 - Testing | ❔ Belum diverifikasi detail | Belum dicek apakah ada test file untuk CRUD/approval integration (pola lain di codebase biasanya punya `approval_integration_test.go`) — perlu verifikasi terpisah sebelum diklaim |

**Frontend**: ❌ belum dimulai — `Attendance.vue` hanya placeholder "Coming soon", tidak ada halaman di bawah `views/modules/attendance/`.

**Rekomendasi urutan lanjutan** (berbeda bentuk dari Leave — di sini ada satu blocker struktural, bukan satu gap P0 yang berdiri sendiri):
1. **Session generation/calculation engine (Phase 6)** — prioritas mutlak pertama. Overtime aktual (Phase 7), Leave Integration (Phase 9), dan Payroll Integration (Phase 13) semuanya secara struktural bergantung pada `attendance_sessions` benar-benar ter-generate; membangun fase-fase itu sebelum Phase 6 akan menghasilkan kode yang tidak punya data untuk diproses.
2. Capture validation (Phase 4-5) — geofence sudah selesai (2026-08-08). Face/device validation masih perlu keputusan produk (provider face-matching, kebutuhan device-restriction) sebelum bisa dikerjakan, tapi tidak memblokir Phase 6 karena `ValidationStatus = PENDING` untuk keduanya tidak menghalangi event dibaca oleh session calculation.
3. Dedicated check-in/check-out endpoints (pisah dari `POST /events` generik) — supaya validasi Phase 4 punya tempat spesifik untuk dipasang per aksi, sama seperti rekomendasi Leave untuk endpoint `submit`/`cancel` khusus.
4. Leave Integration (Phase 9) — setelah Phase 6 ada, ini jadi straightforward: baca leave request APPROVED_FINAL, set session status LEAVE.
5. Frontend dasar — baru masuk akal setelah data yang ditampilkan (sessions dengan status yang benar) benar-benar ada.

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
