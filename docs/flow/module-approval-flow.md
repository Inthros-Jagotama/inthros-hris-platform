# Alur Central Approval Engine (Runbook)

Dokumen ini menjelaskan **cara pakai** modul **Approval Engine** — alur persetujuan
multi-step (ANY_ONE, ALL, N_OF_M), instance tracking, task management untuk approver,
dan integrasi dengan semua modul yang memerlukan approval.

- Lokasi kode: `backend/internal/modules/approval/`
- Halaman UI: `frontend/tenant/src/views/modules/approval/`
- Module slug: `approval`

---

## 1. Ringkasan Alur End-to-End

```
CONFIGURATION                    REQUEST                           ACTION
┌──────────────────┐   ┌──────────────────────────────┐   ┌──────────────────────┐
│ Flow Definition   │   │ Module creates instance       │   │ Approver:            │
│   module slug     │   │   → Step 1 pending            │   │   APPROVE / REJECT   │
│   step sequence   │──▶│   → Step 2 pending (next)     │──▶│   → next step atau   │
│   approver rules  │   │   → Callback ke module         │   │     callback ke modul│
│   org scope       │   │     (APPROVED / REJECTED)      │   │                      │
└──────────────────┘   └──────────────────────────────┘   └──────────────────────┘
```

---

## 2. Entitas Utama

| Entitas | Deskripsi |
|---|---|
| ApprovalFlow | Definisi alur approval (module slug, flow name, is_active) |
| ApprovalFlowStep | Langkah dalam alur (sequence, approver type: USER/ROLE/MANAGER/HR, approval_mode: ANY_ONE/ALL/N_OF_M) |
| ApprovalFlowStepOrganization | Scope organisasi per step (opsional — step berlaku untuk org tertentu) |
| ApprovalInstance | Instance approval yang dibuat saat module meminta approval (module_slug, document_id, status) |
| ApprovalAction | Aksi approve/reject per step (instance_id, step_id, approver_id, action, note) |
| ApprovalTask | Task approval per approver (instance_id, step_id, approver_id, status: PENDING/DONE) |

---

## 3. SETUP — Approval Flows

| Menu | Endpoint | Deskripsi |
|---|---|---|
| List Flows | `GET /approval/flows` | Daftar semua alur approval |
| Create Flow | `POST /approval/flows` | Buat alur baru (module slug, name) |
| Detail Flow | `GET /approval/flows/:flowId` | Detail alur + steps |
| Update Flow | `PUT /approval/flows/:flowId` | Ubah alur |
| Delete Flow | `DELETE /approval/flows/:flowId` | Hapus alur |
| Active Flow | `GET /approval/active-flow` | Alur aktif untuk module tertentu |
| Available Modules | `GET /approval/available-modules` | Daftar module yang bisa pakai approval |

### Flow Steps

| Endpoint | Deskripsi |
|---|---|
| `POST /approval/flows/:flowId/steps` | Tambah langkah (sequence, approver_type, approval_mode, approver_user_id/role_id) |
| `GET /approval/flows/:flowId/steps` | Daftar langkah |
| `PUT /approval/flows/:flowId/steps/:stepId` | Ubah langkah |
| `DELETE /approval/flows/:flowId/steps/:stepId` | Hapus langkah |

> ⚠️ Flow settings hanya bisa diakses dengan permission `approval.settings.<action>` (bukan `approval.view`).

---

## 4. INSTANCE — Request Approval

1. **Module membuat instance** — `POST /approval/instances`: `module_slug`, `document_id`, `document_type`, `requested_by`, `title`, `description`.
2. **Instance status awal** — `PENDING` + task dibuat untuk step pertama.
3. **Approval steps berjalan berurutan** — step N selesai → step N+1 mulai (atau langsung callback jika step terakhir).

---

## 5. ACTION — Approve/Reject

1. **Lihat task pending** — `GET /approval/tasks/pending`: daftar task yang perlu saya approve.
2. **Lihat task selesai** — `GET /approval/tasks/done`: daftar task yang sudah saya kerjakan.
3. **Submit action** — `POST /approval/instances/:id/actions`: `action` (APPROVED/REJECTED), `step_id`, `note`.
4. **Approval mode**:
   - `ANY_ONE` — cukup 1 approver approve → step selesai
   - `ALL` — semua approver harus approve → step selesai
   - `N_OF_M` — minimal N dari M approver harus approve → step selesai
5. **Callback** — setelah step terakhir selesai, instance status berubah → callback ke module (`HandleApprovalStatusChange`).

---

## 6. INSTANCE — Monitoring

| Endpoint | Deskripsi |
|---|---|
| `GET /approval/instances` | Daftar semua instance (filter by module, status) |
| `GET /approval/instances/:id` | Detail instance + steps + actions |
| `POST /approval/instances/:id/cancel` | Batalkan instance (hanya sebelum selesai) |

---

## 7. Ringkasan Status

| Entitas | Status |
|---|---|
| Instance | `PENDING → APPROVED / REJECTED / CANCELLED` |
| Task | `PENDING → DONE` |
| Flow | `is_active: true/false` |

---

## 8. Integrasi Lintas Modul

| Modul | Penggunaan |
|---|---|
| **Employee Movement** | Submit movement → approval → callback approve/reject |
| **Attendance** | Koreksi sesi, lembur (self & assigned) |
| **Leave** | Pengajuan cuti |
| **Training** | Training request |
| **Recruitment** | Recruitment request |
| **Reimbursement** | Pengajuan reimbursement |
| **Competency** | 360 assessment approval |
| **Business Travel** | Travel approval, settlement |

> Modul lain membuat instance via `POST /approval/instances` dengan `module_slug` masing-masing.

---

## 9. Peta Halaman UI

| Menu | Halaman |
|---|---|
| Approvals (hub) | `Approvals.vue` |
| My Tasks | `Approvals.vue` (tab pending tasks) |
| Approval Flows | `ApprovalFlows.vue` |
| Instances | `Approvals.vue` (tab instances) |
| Movement Detail | `detail/MovementDetail.vue` |

---

## 10. Endpoint API Utama

Semua di bawah `/api/v1/tenant/approval/`.

| Area | Endpoint |
|---|---|
| Flows | `GET/POST /flows`, `GET/PUT/DELETE /flows/:flowId` |
| Active Flow | `GET /active-flow` |
| Modules | `GET /available-modules` |
| Steps | `GET/POST /flows/:flowId/steps`, `PUT/DELETE .../steps/:stepId` |
| Instances | `GET/POST /instances`, `GET /instances/:id`, `POST /instances/:id/cancel` |
| Actions | `POST /instances/:id/actions` |
| Tasks | `GET /tasks/pending`, `GET /tasks/done` |

---

## 11. Catatan Penting

- **Central Approval** adalah engine generic — semua modul yang butuh approval menggunakan engine ini.
- **Approval mode** (ANY_ONE/ALL/N_OF_M) dikonfigurasi per step, memungkinkan fleksibilitas alur.
- **Org scope** per step memungkinkan step approval hanya berlaku untuk karyawan organisasi tertentu.
- **Callback mechanism** — setelah keputusan akhir, engine memanggil `HandleApprovalStatusChange` di module terkait.
- **Permission gate** — flow settings (create/update/delete) memerlukan `approval.settings.<action>`, bukan `approval.view`.
- **Module slug** harus unik dan konsisten dengan yang didaftarkan di module terkait.
