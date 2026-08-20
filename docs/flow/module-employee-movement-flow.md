# Alur Employee Movement & Contract Management (Runbook)

Dokumen ini menjelaskan **cara pakai** modul **Employee Movement & Contract Management** —
pergerakan karyawan (promosi, mutasi, demosi, perpanjangan kontrak, pensiun, offboarding),
kontrak kerja (PKWT/PKWTT), approval via Central Approval, audit trail, dokumen pendukung,
career timeline, eligibility check, dan reporting.

- Lokasi kode: `backend/internal/modules/employeemovement/`
- Halaman UI: `frontend/tenant/src/views/modules/employeemovement/`
- Integrasi: career paths ada di modul **Career Intelligence** (strategical); modul ini membaca career_paths untuk promotion eligibility

---

## 1. Ringkasan Alur End-to-End

```
DRAFT → SUBMIT → PENDING_APPROVAL → APPROVED → EXECUTE → EXECUTED
  │                        │              │                    │
  │                        ▼              ▼                    │
  │                   REJECTED       REJECTED                  │
  │                                                          │
  │              CANCELLATION (approved movement)             │
  │              ┌──────────────────────────────────┐         │
  │              │ CANCELLATION_PENDING → APPROVED/  │         │
  │              │   REJECTED → EXECUTED (cancelled) │         │
  │              └──────────────────────────────────┘         │
  └───────────────────────────────────────────────────────────┘

CONTRACT: DRAFT → ACTIVE → EXPIRED / EXTENDED / TERMINATED
```

- **Movement status:** `draft → pending_approval → approved → executed` · `rejected` · `cancelled` · `cancellation_pending` (pembatalan movement approved)
- **Contract status:** `active → expired / extended / terminated`

---

## 2. Entitas Utama

| Entitas | Deskripsi |
|---|---|
| EmployeeMovement | Riwayat pergerakan karyawan (type, from/to org/position/status, SK, effective date) |
| EmployeeMovementAudit | Audit trail perubahan status movement (CREATED/UPDATED/SUBMITTED/APPROVED/REJECTED/EXECUTED/CANCELLED) |
| EmployeeMovementDocument | Metadata dokumen pendukung movement (SK, surat, dll) |
| EmployeeContract | Kontrak kerja (PKWT/PKWTT/daily), periode aktif, perpanjangan |
| CareerPath / CareerPathStep | Jenjang karier (shared dengan Career Intelligence — modul ini hanya membaca untuk eligibility) |

---

## 3. SETUP

Tidak ada master data khusus modul ini. Prasyarat:
- **Organization** — posisi & struktur organisasi harus sudah ada
- **Employee** — data karyawan harus lengkap
- **Setting** — employment status, grading harus sudah ada
- **Career Paths** — dikonfigurasi di modul Career Intelligence (opsional, untuk promotion eligibility)

---

## 4. MOVEMENT — CRUD

1. **Buat movement** — `POST /movements`: `employee_id`, `movement_type` (promotion/demotion/mutation/contract_extension/status_change/retirement/offboarding), `from_*` / `to_*` fields (org, position, employment status), `decision_letter_number`, `decision_letter_date`, `effective_date`, `reason`.
2. **Daftar movement** — `GET /movements`: filter by `movement_type`, `status`, `search` (employee name/ID).
3. **Detail movement** — `GET /movements/:id`: termasuk snapshot nama master (from/to org name, position name, status name).
4. **Ubah movement** — `PUT /movements/:id`: hanya status `draft` yang bisa diubah.
5. **Hapus movement** — `DELETE /movements/:id`: hanya status `draft` yang bisa dihapus.

---

## 5. MOVEMENT — Approval Flow

1. **Submit** — `POST /movements/:id/submit` → status `pending_approval` + instance Central Approval dibuat. Modul `employeemovement`, auto-resolve.
2. **Approval callback** — `HandleApprovalStatusChange(documentID, status, note)`:
   - `approved` → status `approved` + notifikasi
   - `rejected` → status `rejected` + notifikasi
3. **Execute** — `POST /movements/:id/execute` → status `executed` + update employment data karyawan (org, position, status) + notifikasi + audit trail.

> ⚠️ Approval **hanya lewat Central Approval** (submit). Endpoint approve/reject manual tidak ada (keputusan plan §11.5 / G-5).

---

## 6. MOVEMENT — Cancellation

Movement yang sudah **approved** bisa dibatalkan melalui approval:

1. **Cancel** — `POST /movements/:id/cancel` → `cancellation_pending` + instance approval baru (terpisah dari submission).
2. **Callback** — `HandleCancellationStatusChange(documentID, status, note)`:
   - `approved` → `cancelled` + revert employment data + notifikasi
   - `rejected` → kembali ke `approved` + notifikasi

---

## 7. MOVEMENT — Audit Trail & Documents

- **Audit trail** — `GET /movements/:id/audits`: setiap perubahan status tercatat (action, old/new status, old/new data JSON, reason, actor).
- **Documents** — `GET/POST /movements/:id/documents`, `DELETE .../documents/:documentId`: metadata dokumen pendukung (file upload via endpoint generik).
- **Generate Document** — `POST /movements/:id/generate-document`: PDF SK Movement dari template aktif (DocumentTemplate module).
- **Generated Documents** — `GET /movements/:id/generated-documents`.

---

## 8. MOVEMENT — Career Timeline & Eligibility

- **Career History** — `GET /employees/:employeeId/career-history`: timeline riwayat karier dari employment + movements.
- **Movement Eligibility** — `GET /employees/:employeeId/movement-eligibility`: aturan default (tenure, performance, competency).
- **Promotion Eligibility** — `GET /employees/:employeeId/promotion-eligibility`: aturan promosi (tenure, minimum service, performance, competency, OKR) + next step dari career path.
- **Movements by Employee** — `GET /employees/:employeeId/movements`.

---

## 9. CONTRACT — CRUD

1. **Buat kontrak** — `POST /contracts`: `employee_id`, `contract_number`, `contract_type` (pkwt/pkwtt/daily), `start_date`, `end_date` (nullable untuk pkwtt), `decision_letter_number`.
2. **Daftar kontrak** — `GET /contracts`: filter by `status`, `search`.
3. **Detail kontrak** — `GET /contracts/:id`.
4. **Ubah kontrak** — `PUT /contracts/:id`: update fields, `extension_count` bertambah saat perpanjangan.
5. **Hapus kontrak** — `DELETE /contracts/:id`.
6. **Contracts by Employee** — `GET /employees/:employeeId/contracts`.
7. **Generate Document** — `POST /contracts/:id/generate-document`: PDF Perjanjian Kerja dari template aktif.

---

## 10. REPORTS & DASHBOARD

| Endpoint | Deskripsi |
|---|---|
| `GET /reports/movements` | Laporan pergerakan (filter: dateFrom/To, org, position, employee, type, status) |
| `GET /reports/contracts` | Laporan kontrak (summary by type/status) |
| `GET /dashboard` | Dashboard HR (total movements, contracts by status, pending approvals) |

---

## 11. Ringkasan Status

| Entitas | Status |
|---|---|
| Movement | `draft → pending_approval → approved → executed` · `rejected` · `cancelled` · `cancellation_pending` |
| Movement Audit | `CREATED · UPDATED · SUBMITTED · APPROVED · REJECTED · EXECUTED · CANCELLED · CANCELLATION_REQUESTED · CANCELLATION_REJECTED` |
| Contract | `active → expired / extended / terminated` |

---

## 12. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Central Approval** | Approval movement (submit + cancellation) |
| **Career Intelligence** | Career path dibaca untuk promotion eligibility |
| **Employee** | Employment data di-update saat movement executed |
| **Organization** | Posisi asal/tujuan movement |
| **Notification** | Notifikasi perubahan status movement & contract |
| **Performance** | Performance score dibaca untuk eligibility check |
| **Competency** | Competency score dibaca untuk eligibility check |
| **Document Template** | Generate PDF SK movement & contract |

---

## 13. Peta Halaman UI

| Menu | Halaman |
|---|---|
| Movements & Contracts (hub) | `EmployeeMovementReports.vue` |
| Movements | `EmployeeMovements.vue` |
| Contracts | `EmployeeContracts.vue` |
| Approval Detail (Movement) | `MovementDetail.vue` (di modul approval) |

---

## 14. Endpoint API Utama

Semua di bawah `/api/v1/tenant/employee-movements/`.

| Area | Endpoint |
|---|---|
| Movements CRUD | `GET/POST /movements`, `GET/PUT/DELETE /movements/:id` |
| Movement Workflow | `POST /movements/:id/submit`, `POST .../execute`, `POST .../cancel` |
| Movement Audit | `GET /movements/:id/audits` |
| Movement Documents | `GET/POST /movements/:id/documents`, `DELETE .../documents/:documentId` |
| Generate Document | `POST /movements/:id/generate-document`, `GET .../generated-documents` |
| Employee Movement | `GET /employees/:employeeId/movements`, `GET .../career-history`, `GET .../movement-eligibility`, `GET .../promotion-eligibility` |
| Reports | `GET /reports/movements`, `GET /reports/contracts`, `GET /dashboard` |
| Contracts CRUD | `GET/POST /contracts`, `GET/PUT/DELETE /contracts/:id` |
| Contract Documents | `POST /contracts/:id/generate-document`, `GET .../generated-documents` |
| Employee Contract | `GET /employees/:employeeId/contracts` |

---

## 15. Catatan Penting

- **Career paths** dimiliki modul Career Intelligence; modul ini hanya **membaca** career_paths/career_path_steps untuk promotion eligibility.
- **Movement Snapshot** (from/to org/position/status name) disimpan saat pembuatan agar histori tidak berubah ketika master data diubah.
- **Cancellation** memerlukan approval terpisah dari submission (instance approval berbeda).
- **Approval hanya via Central Approval** — tidak ada endpoint approve/reject manual.
- **Document generation** bergantung pada template aktif di modul Document Template; jika tidak ada template aktif, generate gagal.
