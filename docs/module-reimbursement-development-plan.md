# Plan Pengembangan Reimbursement & Claim

> Sumber: `docs/analisis-modul-reimbursements.md` (analisis eksplorasi kode, 2026-08-16).
> Dokumen ini adalah rencana kerja lanjutan — bukan analisis ulang. Untuk detail struktur kode/tabel/state machine yang sudah ada, rujuk dokumen analisis tersebut; dokumen ini fokus pada **apa yang tersisa dan bagaimana mengerjakannya**.

---

# 1. Ringkasan Status

Modul Reimbursement **backend-complete**, **frontend belum ada sama sekali** (placeholder "Coming soon"), dan dua integrasi lintas modul belum dikerjakan.

| Layer | Status |
|---|---|
| Database & migrasi (Postgres + MySQL) | ✅ Selesai |
| Model / Repository / Service / Handler / Routes | ✅ Selesai (±60 test) |
| Integrasi Approval Engine (central) | ✅ Selesai |
| RBAC permissions | ✅ Selesai (5 permission ter-seed) |
| Frontend UI | ❌ Placeholder saja — **prioritas utama plan ini** |
| Integrasi payroll (pembayaran) | ❌ Belum ada — perlu keputusan produk dulu |
| Notifikasi outcome | ❌ Belum ada — masuk Notification Phase 5 |
| Upload lampiran bukti (receipt) | ⚠️ Field `receipt_url` ada, endpoint upload belum dipakai end-to-end |
| Konsistensi tipe kolom timestamp (model.go vs SQL) | ⚠️ Perlu diselaraskan |

`docs/project-completion-dashboard.md` menandai modul ini "✅ Complete" — itu keliru/menyesatkan karena hanya mengukur backend. Plan ini eksis untuk menutup gap tersebut.

---

# 2. Prinsip

- **Tidak membuat approval engine baru** — reimbursement sudah terintegrasi penuh dengan Central Approval Module (`approvalSvc.RegisterStatusHandler("reimbursement", ...)`), pola ini dipertahankan.
- **Tidak mengubah state machine backend** kecuali ditemukan bug nyata saat membangun FE — backend sudah punya ±60 test yang lulus, risiko regresi harus dihindari.
- **Ikuti pola FE module lain yang sudah selesai** (Attendance, Business Travel — lihat `docs/module-attendance-business-travel-development-plan.md` §54.6) sebagai referensi konvensi: Vue 3 `<script setup>`, PrimeVue, panggil `@/services/api` langsung tanpa store/API-client layer terpisah, upload dua-langkah via endpoint generik `POST /api/v1/tenant/uploads`.
- **Keputusan produk dulu, baru kode** untuk integrasi payroll (§5) — jangan berasumsi sepihak apakah `PAID` harus otomatis membuat payslip line item atau tetap manual.

---

# 3. Phase 1 — Frontend: Reimbursement Types (Master Data)

Prasyarat sebelum request bisa dibuat (setiap request butuh `request_type_id`).

- [ ] Halaman list `reimbursement_types` (table: code, name, description, is_active, actions).
- [ ] Dialog create/edit (code, name, description, is_active toggle).
- [ ] Soft-delete (konfirmasi).
- [ ] Endpoint yang dipakai: `POST/GET/PUT/DELETE /api/v1/tenant/reimbursements/types`, `GET /types/:id`.
- [ ] Gate permission: `reimbursement.create`/`update`/`delete` untuk aksi mutasi; `reimbursement.view` untuk akses halaman.

---

# 4. Phase 2 — Frontend: Reimbursement Requests

## 4.1 List

- [ ] Tabel daftar pengajuan dengan filter `employee_id` (untuk admin/HR melihat semua) dan `status`, pagination.
- [ ] Kolom: title, request type, total_amount, status (Tag berwarna), submitted_at.
- [ ] Untuk employee biasa: default filter ke pengajuan miliknya sendiri (pola sama seperti `AttendanceOvertime.vue` pakai `employeeId` dari `useMyEmployee()`).
- [ ] Endpoint: `GET /api/v1/tenant/reimbursements/requests`.

## 4.2 Create / Edit Draft

- [ ] Dialog/halaman create: title, description, request_type, currency (default IDR).
- [ ] Setelah draft dibuat, tambah **item biaya** (expense_date, expense_type, description, amount, receipt_url) — form dinamis multi-row, mengikuti pola "add sub-resource" seperti Business Travel Activities/Schedules.
- [ ] Edit title/description/currency **hanya saat status DRAFT** (guardrail sudah ada di backend, FE tinggal menyembunyikan tombol edit di status lain).
- [ ] Endpoint: `POST/PUT /api/v1/tenant/reimbursements/requests`, `GET/POST/PUT/DELETE /requests/:id/items` & `/items/:itemId`.

## 4.3 Upload Bukti (Receipt)

- [ ] Per item biaya: tombol upload (pola dua-langkah sama seperti Business Travel funding/expense document — lihat `docs/module-attendance-business-travel-development-plan.md` §54.4): `POST /api/v1/tenant/uploads` dulu untuk dapat URL, baru simpan URL itu ke `receipt_url` saat create/update item.
- [ ] **Ini menutup gap yang sudah teridentifikasi di analisis (§5.3/§7.5)** — backend sudah siap (field `receipt_url` ada), FE-nya yang belum pernah dibangun.

## 4.4 Detail View

- [ ] Halaman/dialog detail: info request + daftar item biaya + total + status + riwayat approval (bisa reuse pola tampilan dari `Approvals.vue`'s document detail popup, atau link ke sana).
- [ ] Endpoint: `GET /api/v1/tenant/reimbursements/requests/:id` (preload items).

## 4.5 Aksi Status

- [ ] Tombol **Submit** (DRAFT → SUBMITTED/PENDING_APPROVAL), munculkan pilihan `flow_id` kalau approval module butuh (opsional, auto-resolve kalau tidak dipilih — cek pola `SubmitBusinessTravelRequest.flow_id`).
- [ ] Tombol **Cancel** (status manapun → CANCELLED), dengan konfirmasi.
- [ ] Approve/Reject **tidak perlu dibuat di sini** — sudah tertangani lewat halaman Approvals generik (`Approvals.vue`) via `HandleApprovalStatusChange` (push-based callback, §5.2 analisis). Cukup pastikan detail dokumen reimbursement muncul dengan benar di popup approval (lihat §7 di bawah).
- [ ] Tombol **Pay** (APPROVED → PAID) — **tunda sampai keputusan Phase 5 (§5) dibuat**, supaya tidak membangun UI untuk alur yang mungkin berubah jadi otomatis.
- [ ] Endpoint: `PUT /api/v1/tenant/reimbursements/requests/:id/status`.

---

# 5. Phase 3 — Integrasi Payroll (Keputusan Produk Diperlukan)

**Belum boleh dikerjakan sebelum ada keputusan produk yang jelas** — analisis menemukan `PAID` saat ini murni flag manual, tanpa linkage ke payslip/payroll run manapun.

Pertanyaan yang perlu dijawab dulu:

1. Apakah reimbursement yang `PAID` harus otomatis muncul sebagai line item di payslip berikutnya, atau tetap dianggap "sudah dibayar di luar sistem payroll" (transfer manual, dicatat manual)?
2. Kalau otomatis: masuk sebagai `SalaryEmployeeAdjustment` one-off (pola yang sama dipakai Business Travel, lihat `docs/module-attendance-business-travel-development-plan.md` §54.8) sudah tersedia sebagai extension point tanpa perlu mengubah module payroll — reuse pola itu.
3. Kalau tetap manual: cukup pastikan field `paid_amount`/`paid_at`/catatan referensi transfer ditampilkan jelas di FE, tidak perlu kerja backend tambahan.

**Rekomendasi**: mulai dari opsi manual (3) dulu untuk Phase 2 FE selesai lebih cepat, opsi otomatis (2) sebagai iterasi berikutnya kalau memang dibutuhkan — pola integrasinya sudah terbukti jalan di Business Travel jadi risikonya rendah kalau mau dibangun belakangan.

---

# 6. Phase 4 — Notifikasi

Sesuai `docs/module-notification-plan.md` Phase 5 (belum dikerjakan untuk reimbursement).

- [ ] Wire `Service.SetNotifier(notificationSvc)` untuk `reimbursement.Service` di `main.go` (cek dulu apakah sudah ada — analisis tidak menyebutkan ini di-wire).
- [ ] Notifikasi submit → pending approval (opsional, tergantung kebutuhan).
- [ ] `REIMBURSEMENT_APPROVED` / `REIMBURSEMENT_REJECTED` — dikirim dari `HandleApprovalStatusChange` saat callback approval datang.
- [ ] `REIMBURSEMENT_PAID` — dikirim saat status berubah ke PAID (baik manual maupun otomatis, tergantung hasil keputusan §5).
- [ ] Tambahkan template pesan (EN+ID) ke `backend/internal/modules/notification/i18n.go` — **wajib**, tanpa ini `Notify()` tetap jalan tapi hanya menampilkan raw type string (lihat pola yang sama dikerjakan untuk Business Travel, commit `ab5646b0`).

---

# 7. Phase 5 — Approval Detail Popup

Reimbursement sudah terintegrasi ke Central Approval (module slug `"reimbursement"`), tapi perlu diverifikasi apakah detail dokumennya sudah tampil benar di popup approval generik.

- [ ] Cek `frontend/tenant/src/views/modules/approval/Approvals.vue` — fungsi `documentEndpointFor(module, documentId)` — pastikan case `'reimbursement'` sudah ada dan mengarah ke endpoint yang benar (`/api/v1/tenant/reimbursements/requests/${documentId}`).
- [ ] Kalau belum ada case-nya, tambahkan — polanya sama persis dengan yang baru dikerjakan untuk Business Travel (commit `6ab24b37`): field flat (title, description, total_amount, status, dst) otomatis tampil lewat fallback generik `documentFields`, tidak wajib bikin markup khusus kecuali mau tampilan lebih rapi.

---

# 8. Phase 6 — Housekeeping

- [ ] **Selaraskan tipe kolom timestamp** antara `model.go` (int64 unix-nano) dan migrasi SQL (`TIMESTAMP`) — analisis §3 menandai ini sebagai inkonsistensi yang bisa jadi masalah tergantung `AutoMigrate` vs SQL migration mana yang jalan duluan. Perlu investigasi lebih dulu (bukan langsung diubah) karena berpotensi breaking change pada modul yang sudah production-tested.
- [ ] Update `docs/project-completion-dashboard.md` supaya status "Reimbursement & Claim" mencerminkan kondisi FE yang sebenarnya (jangan biarkan "✅ Complete" menyesatkan pembaca lain).

---

# 9. Urutan Kerja yang Disarankan

1. **Phase 1** (Types master data) — prasyarat, kecil, cepat.
2. **Phase 2** (Request CRUD + items + upload + detail + submit/cancel) — inti pekerjaan, paling besar.
3. **Phase 5** (Approval detail popup) — cepat, bisa dikerjakan paralel dengan Phase 2 karena tidak saling bergantung.
4. **Phase 4** (Notifikasi) — setelah Phase 2 jalan, supaya ada event nyata untuk ditest.
5. **Phase 3** (Payroll integration) — **setelah keputusan produk didapat**, bisa dikerjakan kapan saja setelah itu.
6. **Phase 6** (Housekeeping) — kapan saja, tidak blocking, tapi jangan dilupakan.

---

# 10. Referensi

- Analisis lengkap: `docs/analisis-modul-reimbursements.md`.
- Pola integrasi payroll one-off adjustment (§5.2 di sini): `docs/module-attendance-business-travel-development-plan.md` §54.8.
- Pola upload dua-langkah & approval detail popup: `docs/module-attendance-business-travel-development-plan.md` §54.4, §54.3.
- Rencana notifikasi keseluruhan: `docs/module-notification-plan.md`.
