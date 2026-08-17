# Plan Pengembangan Reimbursement & Claim

> 📅 Status: **✅ SELESAI (2026-08-16)** — seluruh fase dieksekusi (Phase 1–2 FE, Phase 4 notifikasi, Phase 5 approval popup, Phase 6 housekeeping). Satu-satunya item yang **sengaja tidak dibangun**: Phase 3 integrasi payroll (keputusan produk final — pembayaran dicatat manual di modul, §5).
> ✅ **Verifikasi arsip (2026-08-17):** backend (approval + RBAC + ±60 test) & frontend (4 view: hub, types, requests, detail) diverifikasi — lihat [`flow/module-reimbursement-flow.md`](flow/module-reimbursement-flow.md) untuk alur pengisian.

> Sumber: `docs/analisis-modul-reimbursements.md` (analisis eksplorasi kode, 2026-08-16).
> Dokumen ini adalah rencana kerja lanjutan — bukan analisis ulang. Untuk detail struktur kode/tabel/state machine yang sudah ada, rujuk dokumen analisis tersebut; dokumen ini fokus pada **apa yang tersisa dan bagaimana mengerjakannya**.

---

# 1. Ringkasan Status

Modul Reimbursement **selesai penuh** (backend + frontend + integrasi). Ringkasan per layer:

| Layer | Status |
|---|---|
| Database & migrasi (Postgres + MySQL) | ✅ Selesai |
| Model / Repository / Service / Handler / Routes | ✅ Selesai (±75 test) |
| Integrasi Approval Engine (central) | ✅ Selesai |
| RBAC permissions | ✅ Selesai (5 permission ter-seed) |
| Frontend UI | ✅ Selesai (2026-08-16) — hub, types, requests, detail (Phase 1–2) |
| Integrasi payroll (pembayaran) | 🚫 Tidak dibangun — keputusan produk (2026-08-16): pembayaran langsung di module reimbursement, tanpa linkage payroll |
| Notifikasi outcome | ✅ Selesai (2026-08-16) — `REIMBURSEMENT_APPROVED/REJECTED/PAID` (Phase 4) |
| Upload lampiran bukti (receipt) | ✅ Selesai (Phase 2 — dua-langkah via `POST /uploads`) |
| Konsistensi tipe kolom timestamp (model.go vs SQL) | ✅ Fixed (2026-08-16) — `int64` → `*time.Time` (Phase 6) |

---

# 2. Prinsip

- **Tidak membuat approval engine baru** — reimbursement sudah terintegrasi penuh dengan Central Approval Module (`approvalSvc.RegisterStatusHandler("reimbursement", ...)`), pola ini dipertahankan.
- **Tidak mengubah state machine backend** kecuali ditemukan bug nyata saat membangun FE — backend sudah punya ±60 test yang lulus, risiko regresi harus dihindari.
- **Ikuti pola FE module lain yang sudah selesai** (Attendance, Business Travel — lihat `docs/module-attendance-business-travel-development-plan.md` §54.6) sebagai referensi konvensi: Vue 3 `<script setup>`, PrimeVue, panggil `@/services/api` langsung tanpa store/API-client layer terpisah, upload dua-langkah via endpoint generik `POST /api/v1/tenant/uploads`.
- **Tidak ada integrasi payroll untuk pembayaran** — keputusan produk final (2026-08-16): `PAID` adalah flag manual yang dicatat langsung di module reimbursement; tidak ada linkage otomatis ke payslip/payroll run (§5).

---

# 3. Phase 1 — Frontend: Reimbursement Types (Master Data)

Prasyarat sebelum request bisa dibuat (setiap request butuh `request_type_id`).

✅ **Selesai (2026-08-16)** — `frontend/tenant/src/views/modules/reimbursement/ReimbursementTypes.vue`, route `/reimbursements/types`.

- [x] Halaman list `reimbursement_types` (table: code, name, description, is_active, actions).
- [x] Dialog create/edit (name, description, is_active toggle) — **`code` dihilangkan dari form (2026-08-16)**: backend auto-generate UPPER_SNAKE_CASE dari name dengan suffix numerik bila duplikat (`generateReimbursementTypeCode`, pola sama dengan Business Travel expense category/funding method); form tidak mengirim `code` lagi.
- [x] Soft-delete (konfirmasi).
- [x] Endpoint yang dipakai: `POST/GET/PUT/DELETE /api/v1/tenant/reimbursements/types`, `GET /types/:id`.
- [x] Gate permission: `reimbursement.create`/`update`/`delete` untuk aksi mutasi; `reimbursement.view` untuk akses halaman.

---

# 4. Phase 2 — Frontend: Reimbursement Requests

## 4.1 List

✅ **Selesai (2026-08-16)** — `Reimbursements.vue` (ditulis ulang dari placeholder) sebagai halaman list `GET /reimbursements`.

- [x] Tabel daftar pengajuan dengan filter `employee_id` (untuk admin/HR melihat semua) dan `status`, pagination.
- [x] Kolom: title, request type, total_amount, status (Tag berwarna), submitted_at.
- [x] Untuk employee biasa: default filter ke pengajuan miliknya sendiri (pola sama seperti `AttendanceOvertime.vue` pakai `employeeId` dari `useMyEmployee()`). Admin/HR (gate `hasPermission('reimbursement.approve')`) melihat semua + dropdown filter employee.
- [x] Endpoint: `GET /api/v1/tenant/reimbursements/requests`.

## 4.2 Create / Edit Draft

✅ **Selesai (2026-08-16)** — dialog create di list + dialog edit di `ReimbursementRequestDetail.vue` (hanya tampil saat DRAFT).

> **Fix atribusi employee (2026-08-16):** `CreateReimbursementRequest` sebelumnya menyimpan `employee_id` = `ctx "user_id"` (UUID akun user dari JWT), padahal `employee_id` adalah UUID record employee yang berbeda (`employee_accounts` memetakan `user_id` ↔ `employee_id`). Akibatnya request milik karyawan tidak pernah muncul di list filter miliknya sendiri. Service kini resolve `user_id → employee_id` via `employee_accounts` (pola sama dengan `/user-accounts/me`), dan migration **137** (postgres+mysql, belum dijalankan) backfill data lama yang salah atribusi.

- [x] Dialog/halaman create: title, description, request_type, currency (default IDR).
- [x] Setelah draft dibuat, tambah **item biaya** (expense_date, expense_type, description, amount, receipt_url) — form dinamis multi-row, mengikuti pola "add sub-resource" seperti Business Travel Activities/Schedules.
- [x] Edit title/description/currency **hanya saat status DRAFT** (guardrail sudah ada di backend, FE tinggal menyembunyikan tombol edit di status lain).
- [x] Endpoint: `POST/PUT /api/v1/tenant/reimbursements/requests`, `GET/POST/PUT/DELETE /requests/:id/items` & `/items/:itemId`.

## 4.3 Upload Bukti (Receipt)

✅ **Selesai (2026-08-16)** — tombol upload per item di detail (draft only), pola dua-langkah.

- [x] Per item biaya: tombol upload (pola dua-langkah sama seperti Business Travel funding/expense document — lihat `docs/module-attendance-business-travel-development-plan.md` §54.4): `POST /api/v1/tenant/uploads` dulu untuk dapat URL, baru simpan URL itu ke `receipt_url` saat create/update item.
- [x] **Ini menutup gap yang sudah teridentifikasi di analisis (§5.3/§7.5)** — backend sudah siap (field `receipt_url` ada), FE-nya yang belum pernah dibangun.

## 4.4 Detail View

✅ **Selesai (2026-08-16)** — `ReimbursementRequestDetail.vue` (route `/reimbursements/:id`): header + ringkasan (total, request type, currency, paid) + info + daftar item + aksi status.

- [x] Halaman/dialog detail: info request + daftar item biaya + total + status + riwayat approval (bisa reuse pola tampilan dari `Approvals.vue`'s document detail popup, atau link ke sana).
- [x] Endpoint: `GET /api/v1/tenant/reimbursements/requests/:id` (preload items).

> Catatan: `GET /requests/:id` tidak preload items — FE memanggil `GET /requests/:id/items` terpisah. Approval history tetap di-popup Approvals generik (`Approvals.vue` sudah punya case `reimbursement`).

## 4.5 Aksi Status

✅ **Selesai (2026-08-16)** — Submit/Cancel di detail; Pay manual disertakan (keputusan produk: **opsi manual**, lihat §5).

- [x] Tombol **Submit** (DRAFT → SUBMITTED/PENDING_APPROVAL) — **backend ditambah auto-resolve flow aktif** (`GetActiveFlowIDForModule` ditambahkan ke interface `ApprovalEngine` reimbursement + dipanggil saat submit tanpa `flow_id`, pola sama persis dengan Business Travel), sehingga submit FE tanpa `flow_id` tetap masuk Approval inbox.
- [x] Tombol **Cancel** (status manapun → CANCELLED), dengan konfirmasi.
- [x] Approve/Reject **tidak perlu dibuat di sini** — sudah tertangani lewat halaman Approvals generik (`Approvals.vue`) via `HandleApprovalStatusChange` (push-based callback, §5.2 analisis). Cukup pastikan detail dokumen reimbursement muncul dengan benar di popup approval (lihat §7 di bawah).
- [x] Tombol **Pay** (APPROVED → PAID) — keputusan produk diambil: **opsi manual** (flag manual, tanpa linkage payroll), tombol di-gate `reimbursement.approve`, dengan hint UI bahwa integrasi payroll tidak termasuk.
- [x] Endpoint: `PUT /api/v1/tenant/reimbursements/requests/:id/status`.

---

# 5. Phase 3 — Integrasi Payroll (TIDAK DIBANGUN)

**Keputusan produk final (2026-08-16): pembayaran reimbursement TIDAK terintegrasi dengan payroll.** Pembayaran dicatat langsung di module reimbursement — `PAID` adalah flag manual (`paid_at` + `paid_amount`), tanpa linkage otomatis ke payslip/payroll run. Opsi integrasi (`SalaryEmployeeAdjustment` one-off, pola Business Travel §54.8) **tidak akan dikerjakan**; modul ini berdiri sendiri.

Implementasi yang sudah ada dan sesuai keputusan ini:

- [x] Tombol **Pay** (APPROVED → PAID) di halaman detail — langsung di module, gate permission `reimbursement.approve`.
- [x] **Form pembayaran** di halaman detail (dialog "Catat Pembayaran"): jumlah dibayar (default = total, editable untuk partial), metode (`BANK_TRANSFER`/`CASH`/`CHEQUE` — konvensi nilai sama dengan payroll), referensi pembayaran, dan catatan. Kirim via `PUT /requests/:id/status` dengan `{ status: PAID, amount, payment_method, payment_reference, payment_note }`.
- [x] Backend: kolom baru `payment_method` (varchar 50), `payment_reference` (varchar 200), `payment_note` (varchar 500) di `reimbursement_requests` — **migration 138** (postgres + mysql + down, **tidak dijalankan**). `UpdateReimbursementRequestStatus` case `PAID` meng-set `status = PAID`, `paid_at`, `paid_amount` + detail pembayaran; parameter detail dibuat **variadic** agar call-site lama (test ±25) tetap kompatibel. Validasi `oneof=BANK_TRANSFER CASH CHEQUE` di DTO. **Tidak ada kode payroll** di module (terverifikasi).
- [x] Halaman detail menampilkan info pembayaran saat PAID: jumlah, tanggal, metode (label lokal), referensi, catatan.
- [x] UI hint di halaman detail saat status APPROVED: "Menandai dibayar mencatat pembayaran di modul ini. Integrasi ke payroll tidak termasuk." (`reimbursement.manual_pay_hint`, EN+ID).
- [x] Notifikasi `REIMBURSEMENT_PAID` dikirim saat status PAID (manual, sesuai keputusan ini).
- [x] Test: `TestService_PayReimbursementRequest_RecordsPaymentDetails` (persist method/reference/note + reload dari repo).

Catatan: bila di masa depan pembayaran perlu masuk payroll, pola `SalaryEmployeeAdjustment` one-off sudah terbukti di Business Travel (§54.8) dan bisa dipakai tanpa mengubah module payroll — tetapi keputusan saat ini adalah **tidak** membangunnya.

---

# 6. Phase 4 — Notifikasi

Sesuai `docs/module-notification-plan.md` Phase 5 (sebelumnya belum dikerjakan untuk reimbursement).

✅ **Selesai (2026-08-16)** — mengikuti pola Leave/Business Travel (commit `ab5646b0`).

- [x] Wire `Service.SetNotifier(notificationSvc)` untuk `reimbursement.Service` di `main.go`.
- [x] Notifikasi submit → pending approval (opsional, tergantung kebutuhan) — **tidak dibuat** (bukan kebutuhan; `APPROVAL_TASK_ASSIGNED` sudah di-push oleh Approval module sendiri).
- [x] `REIMBURSEMENT_APPROVED` / `REIMBURSEMENT_REJECTED` — dikirim dari `HandleApprovalStatusChange` saat callback approval datang.
- [x] `REIMBURSEMENT_PAID` — dikirim saat status berubah ke PAID (manual, sesuai keputusan §5).
- [x] Tambahkan template pesan (EN+ID) ke `backend/internal/modules/notification/i18n.go` — **wajib**, tanpa ini `Notify()` tetap jalan tapi hanya menampilkan raw type string (lihat pola yang sama dikerjakan untuk Business Travel, commit `ab5646b0`).
- [x] `FindUserIDByEmployeeID` ditambahkan ke `reimbursement.Repository` (pola sama dengan leave) untuk resolve employee → user.
- [x] Test: `notifier_integration_test.go` (6 test) + 2 test auto-resolve flow di `approval_integration_test.go`.

---

# 7. Phase 5 — Approval Detail Popup

Reimbursement sudah terintegrasi ke Central Approval (module slug `"reimbursement"`).

✅ **Selesai (sebelumnya, diverifikasi 2026-08-16)** — case `'reimbursement'` sudah ada di `documentEndpointFor` mengarah ke `/api/v1/tenant/reimbursements/requests/${documentId}` (commit `6ab24b37`); field flat otomatis tampil via fallback generik `documentFields`. Tidak ada perubahan yang diperlukan.

---

# 8. Phase 6 — Housekeeping

- [x] **Selaraskan tipe kolom timestamp** antara `model.go` (int64 unix-nano) dan migrasi SQL (`TIMESTAMP`).

  **Fix (2026-08-16):** 7 field aksi di `ReimbursementRequest` (`supervisor_action_at`, `hr_action_at`, `paid_at`, `submitted_at`, `approved_at`, `rejected_at`, `cancelled_at`) diubah dari `int64` (unix-nano) → `*time.Time` nullable, dan service menulis `&now` (nil → `NULL`). **Pemicu konfirmasi:** endpoint `POST/GET /api/v1/tenant/reimbursements/requests` error MySQL 1292 — GORM menulis `0` ke kolom `TIMESTAMP(6) NULL` (`0000-00-00 00:00:00` ditolak strict mode). Kolom SQL sudah `TIMESTAMP(6) NULL DEFAULT NULL` (MySQL) / `TIMESTAMP NULL` (Postgres) sehingga **tidak perlu migration file baru**. `unixNanoToTimePtr` dihapus; DTO response (sudah `*time.Time`) diisi langsung. Catatan: field `*time.Time` **tanpa** gorm tag eksplisit — tag `type:timestamp(6)` justru memecah scan `*time.Time` di driver SQLite test (terverifikasi), dan AutoMigrate modul tenant memang tidak pernah jalan (SQL migration authoritative).
- [x] Update `docs/project-completion-dashboard.md` supaya status "Reimbursement & Claim" mencerminkan kondisi FE yang sebenarnya (jangan biarkan "✅ Complete" menyesatkan pembaca lain) — row modul diperbarui (16 Agu 2026), status notification-plan, dan test count 60 → 75.

---

# 9. Urutan Kerja yang Disarankan

Status per 2026-08-16:

1. **Phase 1** (Types master data) — ✅ selesai.
2. **Phase 2** (Request CRUD + items + upload + detail + submit/cancel) — ✅ selesai.
3. **Phase 5** (Approval detail popup) — ✅ selesai (sudah ada sebelumnya, diverifikasi).
4. **Phase 4** (Notifikasi) — ✅ selesai.
5. **Phase 3** (Payroll integration) — 🚫 **tidak dibangun** (keputusan produk final 2026-08-16): pembayaran langsung di module reimbursement, tanpa linkage payroll (§5).
6. **Phase 6** (Housekeeping) — dashboard doc ✅; selarasan timestamp `model.go` vs SQL ✅ **fixed (2026-08-16)** — `int64` → `*time.Time` (lihat §8).

---

# 10. Referensi

- Analisis lengkap: `docs/analisis-modul-reimbursements.md`.
- Pola integrasi payroll one-off adjustment (§5.2 di sini): `docs/module-attendance-business-travel-development-plan.md` §54.8.
- Pola upload dua-langkah & approval detail popup: `docs/module-attendance-business-travel-development-plan.md` §54.4, §54.3.
- Rencana notifikasi keseluruhan: `docs/module-notification-plan.md`.
