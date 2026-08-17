# Alur Pengisian Reimbursement & Claim (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Reimbursement** dari setup master data
sampai request dibayar — pola runbook seperti [`module-payroll-user-flow.md`](module-payroll-user-flow.md).

- Plan pengembangan: `module-reimbursement-development-plan.md` *(di-archive: `docs/archive/`)* — selesai (2026-08-16)
- Lokasi kode: `backend/internal/modules/reimbursement/` · `frontend/tenant/src/views/modules/reimbursement/`
- Daftar endpoint + contoh curl: [`../api/api-usage-guide.md`](../api/api-usage-guide.md) → §8.2 (tabel Reimbursements)

---

## 1. Ringkasan Alur End-to-End

```
SETUP (sekali)                    PENGAJUAN (per request)                    PEMBAYARAN
┌──────────────┐     ┌──────────────────────────────────────────────────┐   ┌──────────┐
│ Types        │────▶│ DRAFT → SUBMITTED → (Approval) → APPROVED ──▶ PAID│   │ Approval │
│ (master)     │     │   │                          │                   │   │ via      │
└──────────────┘     │   └── CANCELLED              └── REJECTED         │   │ Approvals│
                     └──────────────────────────────────────────────────┘   └──────────┘
```

- **Status request:** `DRAFT → SUBMITTED → APPROVED → PAID` · terminal: `REJECTED`, `CANCELLED`
- **Approval** tidak punya tombol sendiri di halaman reimbursement — approver menindaklanjuti lewat
  halaman **Approvals** generik (push-based callback, modul `reimbursement`).

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Reimbursement Type | `reimbursement_types` | Master jenis biaya (transport, makan, medis, dll.) — prasyarat request |
| Reimbursement Request | `reimbursement_requests` | Pengajuan biaya karyawan (request_type, currency, amount, note) |
| Request Item | `reimbursement_request_items` | Rincian biaya per request (deskripsi, amount, receipt) |

---

## 3. TAHAP 1 — SETUP (dikerjakan sekali)

### A. Master Reimbursement Types

Menu **Reimbursement → Types** (`/reimbursements/types`, `ReimbursementTypes.vue`).

- CRUD jenis biaya: kode, nama, deskripsi, status aktif.
- Setiap request **wajib** memilih `request_type_id` — tanpa type, request tidak bisa dibuat.
- Endpoint: `GET/POST /reimbursements/types`, `GET/PUT/DELETE /reimbursements/types/:id`.

---

## 4. TAHAP 2 — PENGAJUAN (setiap request)

### A. Buat Request (DRAFT)

- Dari **All Requests** (`/reimbursements/all`) atau **My Requests** (`/reimbursements/my-requests`)
  → tombol buat → dialog create (`ReimbursementRequestDetail.vue`).
- Isi: request type, currency, jumlah (amount), catatan.
- Request baru berstatus **DRAFT** — belum masuk approval/notifikasi.

### B. Tambah Item & Upload Bukti

- Di detail request (`/reimbursements/:id`): tambah **item** biaya (deskripsi + jumlah).
- Upload bukti (receipt) per item — **hanya saat DRAFT**, pola dua-langkah:
  `POST /api/v1/tenant/uploads` (file) → simpan `receipt_url` ke item.
- Item juga bisa di-edit/dihapus selama DRAFT.

### C. Submit (SUBMITTED)

- Tombol **Submit** → `PUT /requests/:id/status` `{"status":"SUBMITTED"}`:
  - Membuat **instance approval** di Central Approval (modul `reimbursement`); flow aktif
    di-auto-resolve; bila flow tidak dapat di-resolve → fallback status `SUBMITTED` biasa.
  - `submitted_at` di-set otomatis.

### D. Approval (via halaman Approvals)

- Approver melihat task di **Approvals** generik → **Approve / Reject**.
- Keputusan dipropagasi balik via callback `HandleApprovalStatusChange`:
  - `APPROVED` → status `APPROVED`
  - `REJECTED` → status `REJECTED` (+ catatan di instance approval)
- Notifikasi ke employee pengaju: `REIMBURSEMENT_APPROVED` / `REIMBURSEMENT_REJECTED`.

### E. Pembayaran (PAID)

- HR/Finance membayar langsung di modul (keputusan produk: **tanpa linkage payroll**):
  `PUT /requests/:id/status` `{"status":"PAID", "payment_method": "BANK_TRANSFER|CASH|CHEQUE", ...}`
- Kolom pembayaran dicatat langsung: `payment_method`, `payment_reference`, `payment_note`,
  `paid_at`, `paid_amount`.
- Notifikasi `REIMBURSEMENT_PAID` ke employee pengaju.

### F. Batal (CANCELLED)

- Request yang belum diputuskan bisa dibatalkan → `CANCELLED`.

---

## 5. Ringkasan Status & Transisi

| Status | Makna | Transisi masuk |
|---|---|---|
| `DRAFT` | Baru dibuat, belum diajukan | create |
| `SUBMITTED` | Diajukan (instance approval dibuat) | submit |
| `APPROVED` | Disetujui approver | callback approval |
| `PAID` | Dibayar (manual, tanpa payroll) | pay (HR) |
| `REJECTED` | Ditolak approver | callback approval |
| `CANCELLED` | Dibatalkan | cancel |

> Endpoint transisi tunggal: `PUT /requests/:id/status` — `status` `oneof=SUBMITTED APPROVED REJECTED PAID CANCELLED`.
> `APPROVED`/`REJECTED` umumnya datang dari callback approval (bukan di-set manual lewat UI).

---

## 6. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Central Approval** | Instance approval modul `reimbursement`; approver bertindak di halaman Approvals |
| **Notification** | `REIMBURSEMENT_APPROVED` / `REIMBURSEMENT_REJECTED` / `REIMBURSEMENT_PAID` ke pengaju |
| **Uploads** | Upload bukti receipt dua-langkah via `POST /uploads` |
| Payroll | 🚫 **Tidak terintegrasi** — pembayaran dicatat manual di modul (keputusan produk 2026-08-16) |

---

## 7. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Reimbursement (hub) | `Reimbursements.vue` | Kartu menu: All Requests / My Requests / Types |
| Reimbursement → All Requests | `ReimbursementRequests.vue` | List semua request (filter status/employee) |
| Reimbursement → My Requests | `ReimbursementRequests.vue` | List request milik sendiri |
| Reimbursement → Detail | `ReimbursementRequestDetail.vue` | Ringkasan + item + upload + aksi status (submit/cancel/pay) |
| Reimbursement → Types | `ReimbursementTypes.vue` | Master jenis biaya |

---

## 8. Endpoint API Utama

Semua di bawah `/api/v1/tenant/reimbursement/`.

| Area | Endpoint |
|---|---|
| Types | `GET/POST /types`, `GET/PUT/DELETE /types/:id` |
| Requests | `GET/POST /requests`, `GET/PUT/DELETE /requests/:id`, `PUT /requests/:id/status` |
| Items | `POST /requests/:id/items`, `PUT/DELETE /requests/:id/items/:itemId` |
| Summary | `GET /reimbursements` (ringkasan/hub) |

---

## 9. Catatan Penting

- **Draft vs Ajukan**: request DRAFT belum masuk approval/notifikasi sampai tombol Submit ditekan.
- **Approval di halaman Approvals** — tidak ada tombol approve/reject di halaman reimbursement.
- **PAID manual tanpa payroll** — `payment_method` `oneof=BANK_TRANSFER CASH CHEQUE`.
- **Receipt upload** hanya saat DRAFT (dua-langkah via `POST /uploads`).
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
