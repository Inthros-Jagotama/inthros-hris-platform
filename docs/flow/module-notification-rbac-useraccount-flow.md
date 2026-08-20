# Alur Notification, RBAC & User Account (Runbook)

Dokumen ini menjelaskan **cara pakai** modul **Notification** (notifikasi push),
**RBAC** (role-based access control), dan **User Account** (akun karyawan) — tiga modul
infrastruktur yang mendukung seluruh platform.

- Lokasi kode: `backend/internal/modules/notification/` · `rbac/` · `useraccount/`
- Module slug: `notification` · `rbac` · `useraccount`

---

## 1. Ringkasan Alur

```
RBAC                            USER ACCOUNT                     NOTIFICATION
┌──────────────────┐   ┌──────────────────────────┐   ┌──────────────────────┐
│ Create Role       │   │ Create employee account   │   │ Module mengirim       │
│ Assign Permissions│──▶│ → Setup email sent         │   │ notifikasi            │
│ Assign User→Role  │   │ → User activates account  │   │ → User melihat list   │
│                   │   │                           │   │ → Mark as read        │
└──────────────────┘   └──────────────────────────┘   └──────────────────────┘
```

---

## 2. MODUL A — RBAC (Role-Based Access Control)

### Entitas

| Entitas | Deskripsi |
|---|---|
| Role | Role/posisi (name, description) |
| Permission | Permission yang tersedia (module.action) |
| RolePermission | Mapping role ↔ permission |
| UserRole | Mapping user ↔ role |

### Setup

1. **Buat role** — `POST /rbac/roles`: name, description.
2. **Lihat permissions** — `GET /rbac/permissions`: daftar semua permission yang tersedia (auto-generated dari modules).
3. **Assign permissions ke role** — `PUT /rbac/roles/:id/permissions`: `permission_ids[]`.
4. **Assign user ke role** — `PUT /rbac/users/:id/roles`: `role_ids[]`.
5. **Daftar roles** — `GET /rbac/roles`.
6. **Ubah role** — `PUT /rbac/roles/:id`.
7. **Hapus role** — `DELETE /rbac/roles/:id`.
8. **Daftar users** — `GET /rbac/users`.

### Endpoint

| Area | Endpoint |
|---|---|
| Roles | `GET/POST /rbac/roles`, `PUT/DELETE /rbac/roles/:id` |
| Permissions | `GET /rbac/permissions` |
| Role Permissions | `PUT /rbac/roles/:id/permissions` |
| Users | `GET /rbac/users` |
| User Roles | `PUT /rbac/users/:id/roles` |

### Catatan

- **Permissions** di-generate otomatis dari `module.Permissions()` saat startup — tidak perlu dibuat manual.
- **Role assignment** bersifat replace (PUT) — array permission_ids/role_ids menggantikan yang lama.
- **RBAC middleware** mengecek permission berdasarkan path (`resource.action`) secara global untuk semua tenant routes.

---

## 3. MODUL B — User Account

### Entitas

| Entitas | Deskripsi |
|---|---|
| UserAccount | Akun karyawan (user_id, employee_id, status) |

### Setup & Management

1. **Lihat akun saya** — `GET /user-accounts/me`: employee_id milik user yang login.
2. **Buat akun karyawan** — `POST /user-accounts/employees/:employeeId`: buat akun + kirim setup email.
3. **Status akun** — `GET /user-accounts/employees/:employeeId`: status akun karyawan.
4. **Kirim ulang email** — `POST /user-accounts/employees/:employeeId/resend`: kirim ulang setup email.

### Endpoint

| Area | Endpoint |
|---|---|
| My Account | `GET /user-accounts/me` |
| Create Account | `POST /user-accounts/employees/:employeeId` |
| Account Status | `GET /user-accounts/employees/:employeeId` |
| Resend Email | `POST /user-accounts/employees/:employeeId/resend` |

### Catatan

- **Akun karyawan** dibuat dari employee_id — setelah employee data lengkap.
- **Setup email** dikirim otomatis saat akun dibuat / di-resend.
- **Employee ID format** dikonfigurasi di modul Settings.

---

## 4. MODUL C — Notification

### Entitas

| Entitas | Deskripsi |
|---|---|
| Notification | Notifikasi untuk user (title, message, type, reference_type, reference_id, is_read) |

### Endpoint

| Area | Endpoint |
|---|---|
| List Notifications | `GET /notifications` |
| Unread Count | `GET /notifications/unread-count` |
| Mark as Read | `PATCH /notifications/:id/read` |
| Mark All as Read | `POST /notifications/read-all` |

### Cara Kerja

1. **Module mengirim notifikasi** — via `Service.notifyMovementOutcome()` atau pola serupa di modul lain.
2. **User melihat notifikasi** — `GET /notifications`: daftar notifikasi (unread dulu).
3. **Mark read** — `PATCH /notifications/:id/read` atau `POST /notifications/read-all`.

### Catatan

- **Notifikasi** bersifat push dari berbagai modul (movement, attendance, leave, training, approval, dll).
- **Reference type + ID** memungkinkan user klik notifikasi → langsung ke halaman terkait.
- **Unread count** dipakai untuk badge counter di UI.

---

## 5. Integrasi Lintas Modul

| Modul | Interaksi |
|---|---|
| **Semua modul** | RBAC middleware mengecek permission untuk setiap request |
| **Semua modul** | Mengirim notifikasi saat ada perubahan status (approval, movement, leave, dll) |
| **Employee Movement** | `notifyMovementOutcome()` mengirim notifikasi saat movement approved/rejected/executed |
| **Employee Movement** | `notifyContractEvent()` mengirim notifikasi saat kontrak expired/extended |
| **Platform Admin** | User management di platform-admin (bukan tenant) |

---

## 6. Catatan Penting

- **RBAC** adalah infrastruktur global — semua modul mendaftarkan permissions di `module.Permissions()`.
- **Permission naming**: `<module>.<action>` (mis. `employee.view`, `employeemovement.create`, `approval.settings.update`).
- **Submenu permissions** (mis. `approval.settings.create`) diperiksa terpisah dari module-level permission untuk endpoint admin-config.
- **User Account** terpisah dari User Management di platform-admin — ini untuk tenant-scoped account.
- **Notification** belum punya UI khusus — biasanya ditampilkan via header/bell icon di layout utama.
