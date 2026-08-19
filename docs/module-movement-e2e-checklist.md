# Checklist Verifikasi E2E Manual — Employee Movement & Career Management

> 📅 Dibuat: 2026-08-10 · Basis: Testing Plan §10 `docs/archive/module-movement-plan.md` · Status: **SIAP DIEKSEKUSI** (checklist disusun dari implementasi aktual langkah 1–12 + enhancement P1/P2 s.d. item 14 — termasuk reporting §12.17 & dashboard §12.18; eksekusi menunggu environment tenant + akun dengan permission `employeemovement.*`).

---

## 0. Prasyarat

| # | Item | Keterangan |
|---|---|---|
| P1 | Backend berjalan | `make run` atau `go run ./cmd/server --config ./config/config.yaml` → `http://localhost:8080` |
| P2 | Frontend tenant berjalan | `cd frontend/tenant && npm run dev` → `http://localhost:5174` (proxy `/api` & `/uploads` → `:8080`) |
| P3 | Database & Redis | MySQL/Postgres tenant (migration s.d. **087** ter-apply — termasuk 083 snapshot, 084 audit, 085 documents, 086 career paths, 087 cancellation) + Redis aktif |
| P4 | Akun tenant | User dengan permission `employeemovement.*` (create/update/delete/execute) + `approval.*` (untuk approve/reject di Central Approval) |
| P5 | **Approval flow aktif untuk modul `employeemovement`** | Via FE **Approvals → Flows**: buat flow untuk module `employeemovement` + tambah step approver + **set active**. Tanpa ini `submit` gagal (`approval flow not configured...`, G-3). Verifikasi: `GET /api/v1/tenant/approval/active-flow?module=employeemovement` |
| P6 | Data pendukung | Minimal 1 employee **aktif + punya employment** (org A, posisi P1, status S1, `effective_end_date` kosong), 1–2 organization aktif (summary active), 1 employment status |

### Data yang dipakai (isi sebelum eksekusi)

| Nama | Nilai |
|---|---|
| Employee (nama/id) | |
| Organization asal / tujuan | |
| Position asal / tujuan | |
| Employment status asal / tujuan | |
| Nomor SK & tanggal | |
| Effective date (A) / (B) | |
| Flow approval (id) | |

---

## 1. Skenario A — Alur lengkap Movement (Promotion, inti G-1 + approval)

> Kolom **Cara** berisi jalur FE utama dan API alternatif (cURL). Login HR untuk create/submit/execute; login approver untuk approve/reject.

| # | Langkah | Cara (FE / API) | Hasil yang diharapkan | ✅/❌ | Catatan |
|---|---|---|---|---|---|
| A1 | Persiapan | FE: Employees → buka employee, pastikan ada employment aktif | `GET /employees/:id` → `employments[]` berisi 1 record tanpa `effective_end_date` | | |
| A2 | Create draft promotion | FE: Career → **Movements** → **+ New Movement** → employee, type *Promotion*, `to_position` (org/posisi), SK number+date, `effective_date` (boleh future, mis. +7 hari), reason → Save<br>API: `POST /api/v1/tenant/employee-movements/movements` | Toast sukses; row muncul status **draft**; respons enriched: `employee_name`, `employee_code`, `to_organization_name`, `to_position_name` terisi (G-4) | | |
| A3 | Validasi per tipe (negatif, G-7) | ⚠️ Validasi FE memblokir Save lebih dulu — untuk menguji 400 backend, panggil **API langsung**: `POST /movements` dengan `movement_type=mutation` tanpa `to_organization_id`/`to_position_id` | Error **400 `VALIDATION_ERROR`** + field error (mis. `to_organization_id`); via FE cukup pastikan tombol Save ter-block dengan pesan `field_required` | | Backend re-validasi saat create & update |
| A4 | Submit | FE: tombol **Submit** pada baris draft → konfirmasi<br>API: `POST /movements/:id/submit` body `{}` (tanpa `flow_id` → auto-resolve G-3) | Status → **pending_approval**; instance approval dibuat (module `employeemovement`); approver mendapat **task/notifikasi assignment dari Central Approval** (modul tidak mengirim `MOVEMENT_SUBMITTED` manual, §3.8.2) | | |
| A5 | Approve (Central Approval) | FE (login approver): **Approvals** → instance → **Approve**<br>API: `POST /approval/instances/:id/actions` body `{"action":"APPROVE"}` | Instance `APPROVED`; push-callback → movement status **approved** + `approved_at` terisi; notifikasi `MOVEMENT_APPROVED` | | |
| A6 | Reject path (verifikasi status `rejected`, §11.4) | Buat movement lain → submit → **Reject** (API `{"action":"REJECT"}`) | Movement status **`rejected`** (bukan `cancelled`); FE badge merah "Ditolak"; notifikasi `MOVEMENT_REJECTED` | | |
| A7 | Execute (manual oleh HR, §11.1) | FE: tombol **Execute** pada baris approved → konfirmasi<br>API: `POST /movements/:id/execute` | Status → **executed**; notifikasi `MOVEMENT_EXECUTED`; `executed_at` terisi | | |
| A8 | **Verifikasi employment (inti G-1)** | `GET /api/v1/tenant/employees/:id` → `employments[]` | (a) Employment lama: `effective_end_date` = `effective_date − 1`; (b) Employment **baru**: `organization_id`/`position_id`/`employment_status_id` = nilai `to_*`, `decision_letter_number` = SK, `effective_date` sesuai; (c) movement punya `to_employment_id` | | **Poin penting** — transaksi employment |
| A9 | Cancel path | Buat movement → submit → **Cancel** (FE tombol ×, API `POST /movements/:id/cancel`) | Status → **cancelled** | | |
| A10 | Delete path | Buat draft → **Delete** (FE ikon hapus, API `DELETE /movements/:id`) | Draft hilang dari list | | Hanya draft |

### Negatif / edge yang disarankan

- **Submit tanpa flow aktif** → error jelas (`approval flow not configured...`), bukan 500.
- **Execute sebelum approved** → ditolak (status harus `approved`).
- **Double execute** → ditolak karena status bukan lagi `approved` (validasi status di service); `CloseEmployment` juga punya guard `effective_end_date IS NULL` agar tidak menimpa employment yang sudah tertutup.
- **Approve manual** → endpoint `POST /movements/:id/approve` **sudah tidak ada** (404), satu pintu approval = submit (G-5).

---

## 2. Skenario B — Offboarding / Retirement (employee non-aktif, §11.3)

| # | Langkah | Cara | Hasil yang diharapkan | ✅/❌ |
|---|---|---|---|---|
| B1 | Create + submit + approve | Sama seperti A2→A5, tipe *Offboarding* atau *Retirement* (**tanpa** `to_*`) | Validasi lolos (G-7: tipe ini boleh tanpa to_*); status `approved` | |
| B2 | Execute | `POST /movements/:id/execute` | Employment aktif **ditutup tanpa employment baru**; `GET /employees/:id` → `status = "inactive"` | |

---

## 3. Skenario C — CRUD Kontrak (halaman Contracts)

| # | Langkah | Cara (FE / API) | Hasil yang diharapkan | ✅/❌ | Catatan |
|---|---|---|---|---|---|
| C1 | Create PKWT + upload | FE: Career → **Contracts** → **+ New Contract** → employee, `contract_number`, type *PKWT*, `start_date`, `end_date` (wajib PKWT), SK, **upload dokumen** (`POST /uploads` FormData `file` → `data.url` disimpan ke `document_url`) → Save<br>API: `POST /contracts` | Row muncul status **active**; kolom dokumen menampilkan link lampiran; `extension_count = 0` | | |
| C2 | List + filter | Filter status + search (`contract_number`/nama/kode) | Hasil sesuai filter (backend `ListContracts` param `status` & `search`) | | |
| C3 | Edit | FE tombol edit: ubah `end_date` / status → Save | Perubahan tersimpan; respons `ContractResponse` terisi | | |
| C4 | **Extension chain (G-6)** | Create kontrak B dengan `previous_contract_id = A` (via FE dropdown kontrak milik employee) | `extension_count B = A.extension_count + 1`; **A berubah status `extended`**; `previous_contract_number` tampil di respons B | | Coba 2x bertingkat → count 1, 2 |
| C5 | Delete | Hapus kontrak (bukan yang jadi `previous`) → konfirmasi | Kontrak hilang dari list | | |

---

## 4. Skenario D — Detail dialog & Deep-link (langkah 12)

| # | Langkah | Cara | Hasil yang diharapkan | ✅/❌ |
|---|---|---|---|---|
| D1 | Buka detail movement | FE Movements → klik ikon **mata** pada baris | Dialog menampilkan: Tag tipe + Tag status (badge `rejected` merah untuk status rejected), employee (nama+kode), section **Dari → Ke** (nama org/posisi/status), nomor SK (mono), tanggal SK/efektif/created, alasan & catatan, `approved_at`/`executed_at` bila ada | |
| D2 | Aksi dari dalam dialog | Pada draft: **Submit** dari footer dialog detail | Confirm dialog terbuka → confirm → status berubah, detail tertutup, list ter-refresh | |
| D3 | Deep-link notifikasi | FE **Notifications** → klik notifikasi `MOVEMENT_*` | Pindah ke `/admin/career/movements`; notifikasi ditandai read | |
| D4 | Deep-link approval | FE **Approvals** → modal instance employeemovement | Detail di-load dari `GET /employee-movements/movements/:id` (case `employeemovement` sudah ada) | |

---

## 5. Skenario E — Movement Report & Contract Report (halaman Reports, plan §12.17)

> Prasyarat: module `employeemovement` aktif (menu Reports hanya muncul di Sidebar Career bila module aktif & permission `employeemovement.view`).
>
> **Catatan kepemilikan Career Paths (2026-08-10):** Career Paths kini milik modul **Career Intelligence** (strategical) — endpoint `/api/v1/tenant/career-intelligence/paths` (ladder-style), FE route `/career-intelligence/paths`, sidebar di grup Strategic, permission `careerintelligence.view`. Modul Employee Movement hanya membaca career paths untuk promotion eligibility. Plan lengkap: `docs/module-career-intelligence-plan.md` §4/§6/§7.

| # | Langkah | Cara (FE / API) | Hasil yang diharapkan | ✅/❌ | Catatan |
|---|---|---|---|---|---|
| E1 | Buka halaman Reports | FE: Career → **Reports** (route `/admin/career/reports`, menu sidebar `employee_movement.reports`)<br>API: `GET /reports/movements` + `GET /reports/contracts` (di-load paralel saat mount) | Halaman menampilkan: filter (periode, org, posisi, tipe, status), kartu statistik per tipe movement, breakdown per status, dan kartu Contract Report (Active/Expired/Extended/Terminated) | | |
| E2 | **Movement Report — tanpa filter** | API: `GET /reports/movements` (tanpa param) | `data.total` = jumlah seluruh movement; `data.by_type` = jumlah per tipe (promotion/demotion/mutation/contract_extension/status_change/retirement/offboarding/other — hanya kunci yang ada datanya); `data.by_status` = jumlah per status (draft/pending_approval/approved/rejected/executed/cancelled/cancellation_pending) | | Konsistensi: `sum(by_type) == total == sum(by_status)` |
| E3 | **Filter periode** | FE: isi `date_from` & `date_to` → Refresh<br>API: `GET /reports/movements?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD` | Total ter-filter ke movement dengan `effective_date` dalam rentang (inklusif); kartu by_type/by_status ikut ter-filter | | |
| E4 | **Filter organisasi / posisi** | FE: dropdown Organisasi & Posisi<br>API: `GET /reports/movements?organization_id=...&position_id=...` | Filter mencocokkan `to_*` **ATAU** `from_*` (report mencakup kedua arah); hasil turun sesuai | | |
| E5 | **Filter tipe + status** | FE: dropdown Tipe & Status<br>API: `GET /reports/movements?movement_type=promotion&status=approved` | Kombinasi filter sekaligus → total sesuai; `by_type`/`by_status` terfilter | | |
| E6 | **Periode terbalik → 400** | API: `GET /reports/movements?date_from=2026-08-10&date_to=2026-08-01` | Error **400 VALIDATION_ERROR** (`date_from cannot be after date_to`) — bukan 200 dengan 0 baris | | |
| E7 | **UUID invalid → 400** | API: `GET /reports/movements?organization_id=abc` (atau position/employee) | Error **400 VALIDATION_ERROR** (`invalid organization_id`) | | |
| E8 | **Contract Report** | API: `GET /reports/contracts` | `data.total` = jumlah kontrak; `data.by_status` (active/expired/extended/terminated); `data.expiring` = kontrak **active** yang berakhir dalam ≤30 hari (subset active — bukan bucket terpisah, tidak menambah total) | | |
| E9 | Konsistensi kartu FE vs API | Bandingkan kartu di halaman Reports dengan output API di atas | Angka sama; kartu per tipe yang kosong menampilkan 0 (FE menampilkan semua bucket, backend hanya mengembalikan bucket berdata) | | |

---

## 6. Skenario F — HR Dashboard (plan §12.18)

> Prasyarat: module `employeemovement` aktif — section "Employee Movement" di Dashboard utama hanya dirender bila module aktif (`hasModule`).

| # | Langkah | Cara (FE / API) | Hasil yang diharapkan | ✅/❌ | Catatan |
|---|---|---|---|---|---|
| F1 | Load Dashboard utama | FE: Dashboard → tunggu section "Employee Movement" muncul (skeleton saat loading) | Section tampil berisi: kartu movement per tipe, highlight **Pending Approval** + **Effective This Month**, ringkasan kontrak (Active/Expiring<30d/Expired) | | |
| F2 | **Dashboard API** | API: `GET /employee-movements/dashboard` | `data.movement_by_type` (jumlah movement per tipe, semua status) · `data.pending_approval` (status `pending_approval`) · `data.effective_this_month` (effective_date di bulan berjalan, server-time) · `data.contracts.{active,expiring,expired}` | | |
| F3 | Konsistensi dashboard vs report | Bandingkan `pending_approval` dashboard dengan `GET /reports/movements?status=pending_approval` → total; `effective_this_month` dengan filter `date_from=YYYY-MM-01&date_to=<akhir bulan>` | Angka konsisten (dashboard memakai agregasi yang sama dgn report §12.17) | | |
| F4 | Tombol **View Reports** | FE: klik tombol di section dashboard | Pindah ke `/admin/career/reports` | | |
| F5 | **Module-gating** | Nonaktifkan module `employeemovement` (tenant module mgmt) → reload Dashboard | Section "Employee Movement" **tidak muncul**; halaman Reports juga tidak ada di Sidebar | | |
| F6 | Best-effort fetch | (Opsional) hentikan server sementara lalu load Dashboard | Section tidak merusak Dashboard — data kartu kosong + skeleton selesai; tidak ada unhandled error | | |

---

## 7. Regresi & lintas bahasa

| # | Item | Hasil yang diharapkan | ✅/❌ |
|---|---|---|---|
| R1 | Switch bahasa EN/ID | Label tipe/status/dialog detail/kontrak tampil bilingual benar (mis. "Ditolak" / "Rejected") | |
| R2 | Filter kombinasi Movements | Type + status + search sekaligus → total records konsisten | |
| R3 | Build & konsol | `npm run build` bersih; tidak ada error/⚠️ di console browser | |
| R4 | Permission aksi | ⚠️ *Catatan: saat ini FE menampilkan aksi tanpa cek `hasPermission` per tombol — belum jadi gate UI; verifikasi batasan dilakukan di sisi backend (authz). Perbaikan FE bisa dijadwalkan sebagai enhancement.* | |
| R5 | Reports & dashboard gated module | Halaman Reports (`/admin/career/reports`) & section Dashboard hanya tampil saat module `employeemovement` aktif (F5); tanpa module → item tidak ada di Sidebar, direct-URL tetap di route-guard | | |

---

## 8. Kriteria Penerimaan (acceptance dari plan §10)

- [ ] Transisi status benar: `draft → pending_approval → approved → executed`; `approved → rejected`; `pending_approval/approved → cancelled`; `draft → edit/delete`.
- [ ] `ExecuteMovement` benar-benar menutup employment lama (`effective_date − 1`) & membuat employment baru (G-1), offboarding/retirement men-set employee `inactive` (B2).
- [ ] Satu pintu approval — tidak ada `POST /movements/:id/approve` (G-5).
- [ ] Auto-resolve flow saat submit tanpa `flow_id` (G-3).
- [ ] Notifikasi `MOVEMENT_SUBMITTED/APPROVED/REJECTED/EXECUTED` (G-2).
- [ ] Contract extension count berantai + previous `extended` (G-6).
- [ ] Validasi per tipe mengembalikan 400 `VALIDATION_ERROR` (G-7).
- [ ] Dua halaman FE (`/admin/career/movements`, `/admin/career/contracts`), detail dialog, badge `rejected`, deep-link notifikasi (G-8 + langkah 10–12).
- [ ] **Movement Report** (`GET /reports/movements`) mengembalikan `total`/`by_type`/`by_status` yang konsisten (`sum(by_type) == total == sum(by_status)`), filter periode/org/posisi/tipe/status benar, dan periode terbalik / UUID invalid → 400 (E2–E7, §12.17).
- [ ] **Contract Report** (`GET /reports/contracts`) mengembalikan `by_status` + `expiring` sebagai subset `active` ≤ 30 hari (E8, §12.17).
- [ ] **HR Dashboard** (`GET /employee-movements/dashboard`) menampilkan movement by type, `pending_approval`, `effective_this_month`, ringkasan kontrak; konsisten dengan report & hanya tampil saat module aktif (F2–F5, §12.18).
- [ ] Unit/service/integration test backend **PASS** (sebagai pelengkap verifikasi manual).

---

## 9. Bukti yang disarankan untuk disimpan

- Screenshot tiap langkah (khususnya A8: `employments[]` sebelum/sesudah execute, D1 detail dialog, E1 halaman Reports, dan F1 section Dashboard).
- ID movement / contract / approval instance yang diuji.
- Output cURL untuk langkah kunci (A4 submit, A5 approve, A7 execute, A8 cek employment) + output JSON `GET /reports/movements`, `GET /reports/contracts`, `GET /employee-movements/dashboard`.
- Catatan anomali (bug, pesan error, perilaku tidak sesuai ekspektasi) → laporkan untuk diperbaiki.

---

## 10. Referensi Endpoint (ringkas)

```text
Movement : POST/GET  /api/v1/tenant/employee-movements/movements
           GET/PUT/DELETE /api/v1/tenant/employee-movements/movements/:id
           POST  /api/v1/tenant/employee-movements/movements/:id/submit|execute|cancel
Contract : POST/GET  /api/v1/tenant/employee-movements/contracts
           GET/PUT/DELETE /api/v1/tenant/employee-movements/contracts/:id
Approval : POST  /api/v1/tenant/approval/instances/:id/actions   {"action":"APPROVE"|"REJECT","note":"..."}
           GET   /api/v1/tenant/approval/instances?module=employeemovement&status=pending
           GET   /api/v1/tenant/approval/active-flow?module=employeemovement
Employee : GET   /api/v1/tenant/employees/:id                     → status + employments[]
Upload   : POST  /api/v1/tenant/uploads                            (FormData file → data.url)
Notif    : GET   /api/v1/tenant/notifications                      (type = MOVEMENT_*)
Report   : GET   /api/v1/tenant/employee-movements/reports/movements   ?date_from&date_to&organization_id&position_id&employee_id&movement_type&status
           GET   /api/v1/tenant/employee-movements/reports/contracts    → {total, by_status, expiring}
Dashboard: GET   /api/v1/tenant/employee-movements/dashboard           → {movement_by_type, pending_approval, effective_this_month, contracts}
```
