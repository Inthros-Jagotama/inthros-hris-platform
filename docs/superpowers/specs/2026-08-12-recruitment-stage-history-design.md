# Recruitment Pipeline Stage History (G-5) — Design

> Ref: `docs/module-recruitment-development-plan.md` §G-5, §12 (M4), §16, §18.

## Goal

Aplikasi rekrutmen (`job_applications`) saat ini punya kolom `status` yang bisa
diubah ke nilai apa saja tanpa validasi (`UpdateApplicationStatus` menerima
string bebas), dan tidak ada jejak audit siapa mengubah status apa menjadi
apa dan kapan. G-5 menambahkan:

1. Validasi transisi status (state machine) di satu titik source-of-truth.
2. Tabel histori (`job_application_stage_histories`) yang mencatat setiap
   perubahan status, termasuk yang dipicu otomatis oleh `AcceptOffer` (G-3).
3. Endpoint untuk membaca histori sebuah aplikasi.

## Non-goals

- Tidak mengganti taxonomy `CandidateStatus` (tetap 8 nilai existing: `NEW,
  SCREENED, SHORTLISTED, INTERVIEWED, OFFERED, ACCEPTED, REJECTED,
  WITHDRAWN`) — taxonomy 10-stage yang disebut di plan lama (`APPLIED,
  SCREENING, ASSESSMENT, FINAL_REVIEW, HIRED`, dst.) tidak dipakai, supaya
  tidak perlu migrasi data status dan tidak mengubah logika existing yang
  bergantung padanya (auto `slots_filled`, offer linkage, dst — lihat
  `service.go:1145-1169` dan `service.go:690-745`).
- Tidak ada halaman FE baru di iterasi ini (pipeline board / kanban tetap di
  G-12 lanjutan, menunggu API ini tersedia).
- Tidak mengubah perilaku `AcceptOffer` selain menambah pencatatan history —
  business rule guard idempotensi (`wasAccepted`) tetap sama.

## Data Model

### Tabel baru: `recruitment_stages` (master, seeded)

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `code` | VARCHAR(20) NOT NULL UNIQUE | = nilai `CandidateStatus` (`NEW`, `SCREENED`, dst.) |
| `name` | VARCHAR(100) NOT NULL | label tampilan (mis. "New Application") |
| `sort_order` | INT NOT NULL DEFAULT 0 | urutan tampilan pipeline |
| `created_at`, `updated_at` | TIMESTAMP(6) | |

Diseed idempotent (8 baris, satu per `CandidateStatus`) lewat
`module.go Seed()` — pola sama dengan seed 10 onboarding task template
(cek by `code` sebelum insert, skip kalau sudah ada).

### Tabel baru: `job_application_stage_histories`

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `application_id` | CHAR(36) NOT NULL | FK → `job_applications`, index `idx_ash_app` |
| `from_stage_id` | CHAR(36) NULL | FK → `recruitment_stages`; NULL untuk baris pertama (aplikasi baru dibuat berstatus `NEW`, tidak ada "from") |
| `to_stage_id` | CHAR(36) NOT NULL | FK → `recruitment_stages` |
| `changed_by` | CHAR(36) NULL | user id aktor; NULL bila transisi dipicu otomatis tanpa actor (mis. job/system) |
| `notes` | TEXT NULL | disalin dari `rejection_reason`/`notes` request bila ada |
| `changed_at` | BIGINT NOT NULL | unix nano, pola existing (`applied_at`, dst.) |

Index: `idx_ash_app (application_id)`, `idx_ash_changed_at (changed_at)`.

Migration baru: `097_recruitment_stage_history` (mysql + postgres, up/down
idempotent) — nomor lanjutan dari `096_recruitment_employee_handoff`.

## State Machine

Dikonfirmasi dari kode existing: `CreateOffer` (`service.go:467-492`) **tidak**
mensyaratkan aplikasi sudah berstatus `OFFERED` — hanya validasi aplikasi
ada. Artinya recruiter bisa membuat & mengirim offer kapan saja tanpa
melewati status `OFFERED` di aplikasi lebih dulu, sehingga `AcceptOffer`
bisa memicu transisi ke `ACCEPTED` dari status non-terminal manapun (bukan
hanya dari `OFFERED`).

Transisi valid (dari → daftar tujuan yang diizinkan):

```
NEW, SCREENED, SHORTLISTED, INTERVIEWED, OFFERED
    → REJECTED, WITHDRAWN     (boleh dari status non-terminal manapun)
    → ACCEPTED                (boleh dari status non-terminal manapun — offer
                                bisa dibuat/diterima tanpa melewati OFFERED
                                lebih dulu)

Progresi maju satu-langkah (khusus, tidak boleh lompat):
NEW         → SCREENED
SCREENED    → SHORTLISTED
SHORTLISTED → INTERVIEWED
INTERVIEWED → OFFERED

ACCEPTED, REJECTED, WITHDRAWN → (terminal, tidak ada transisi keluar)
```

Implementasi: satu tabel `map[CandidateStatus][]CandidateStatus` berisi
tujuan yang diizinkan per status asal (union dari progresi maju + `{ACCEPTED,
REJECTED, WITHDRAWN}` untuk tiap status non-terminal) — tidak perlu
membedakan "dipanggil dari mana" (manual endpoint vs `AcceptOffer`), cukup
satu aturan transisi yang berlaku untuk kedua caller.

Transisi ke status yang **sama** (`from == to`) diperlakukan sebagai no-op:
tidak menulis baris history baru, tidak error. Ini diperlukan supaya
`AcceptOffer` pada offer kedua di aplikasi yang statusnya sudah `ACCEPTED`
tetap idempoten seperti perilaku sekarang (`wasAccepted` guard di
`service.go:726`).

Transisi yang tidak ada di daftar (termasuk dari status terminal, atau
lompat tahap seperti `NEW → ACCEPTED`) mengembalikan error:
`fmt.Errorf("invalid status transition: %s -> %s", from, to)`.

## Implementation Point: Single Source of Truth

Tambah method privat baru di `service.go`:

```go
// transitionApplicationStatus memvalidasi transisi, menulis history, dan
// meng-update a.Status + timestamp stage (logika existing dari
// UpdateApplicationStatus dipindah ke sini). Dipanggil oleh
// UpdateApplicationStatus (manual, via endpoint) dan AcceptOffer (otomatis).
// Tidak melakukan repo.UpdateApplication — caller yang menyimpan, supaya
// caller bisa menggabungkan perubahan lain (mis. RejectionReason/Notes) dalam
// satu update.
func (s *Service) transitionApplicationStatus(ctx context.Context, a *JobApplication, newStatus CandidateStatus, changedBy *uuid.UUID, notes string) error
```

- `UpdateApplicationStatus` (existing, dipakai handler `PUT
  /applications/:id/status`) memanggil helper ini, lalu tetap melakukan
  `repo.UpdateApplication`.
- `AcceptOffer` (`service.go:725-732`) memanggil helper ini alih-alih
  langsung `a.Status = CandStatusAccepted; a.AcceptedAt = &now`.
- Kalau helper mengembalikan error (transisi invalid) di jalur
  `AcceptOffer`, **offer accept tetap gagal** (bukan best-effort) — beda
  dengan downstream best-effort yang sudah ada (update requisition
  `slots_filled` tetap warning-only), karena transisi status aplikasi
  adalah bagian inti dari accept, bukan efek samping.

`changed_by` diambil di handler lewat `c.GetString("user_id")` (pola
`performance/okr_handler.go:33`), diteruskan ke service sebagai
`*uuid.UUID` (nil kalau kosong/parse gagal — mis. dipicu dari `AcceptOffer`
tanpa request HTTP eksplisit, pakai user id dari konteks accept-offer yang
sedang berjalan bila ada, else nil).

## API

### `GET /api/v1/tenant/recruitment/applications/:id/history`

Response: list `StageHistoryResponse` terurut `changed_at` ASC:

```json
{
  "success": true,
  "data": [
    {
      "id": "...",
      "from_stage": null,
      "to_stage": {"code": "NEW", "name": "New Application"},
      "changed_by": null,
      "notes": null,
      "changed_at": 1755000000000000000
    },
    {
      "id": "...",
      "from_stage": {"code": "NEW", "name": "New Application"},
      "to_stage": {"code": "SCREENED", "name": "Screened"},
      "changed_by": "<user-uuid>",
      "notes": "passed initial screening",
      "changed_at": 1755003600000000000
    }
  ]
}
```

Baris pertama (histori "NEW" awal, `from_stage: null`) ditulis saat
`CreateApplication` dipanggil — supaya histori lengkap dari awal tanpa
gap, bukan hanya mulai dari transisi pertama.

## Error Handling

- Transisi invalid → `400 Bad Request`, pesan bilingual mengikuti pola
  handler existing (lihat `approval.EmitRoutingError` di G-1 sebagai
  referensi format, tapi ini bukan approval routing error — pakai error
  biasa yang di-map ke 400 di handler, pola `UpdateOffer`/`DeleteOffer`
  yang menolak non-DRAFT).
- `GET .../history` untuk `application_id` yang tidak ada → `404`.
- Migration tetap idempotent (guard `IF NOT EXISTS` / cek existing index),
  pola migration 093-096.

## Testing Plan

- **Migration**: up/down idempotent test (pola existing repo test untuk
  migration, kalau ada; kalau tidak ada pola test migration khusus, cukup
  smoke-test lewat service test yang jalan di atas skema hasil migrasi).
- **Service** (`service_test.go`):
  - Setiap transisi valid dari daftar di atas → sukses, history baris
    baru tertulis dengan `from`/`to` benar.
  - Transisi invalid (mis. `NEW → ACCEPTED`, `ACCEPTED → SCREENED`) → error,
    tidak ada history baru ditulis, `a.Status` tidak berubah.
  - Transisi ke status sama → no-op, tidak ada history baru, tidak error.
  - `CreateApplication` menulis baris history pertama (`from: null, to:
    NEW`).
  - `AcceptOffer` menulis history `<status aplikasi saat itu> → ACCEPTED`
    (bisa dari `NEW`/`SCREENED`/`SHORTLISTED`/`INTERVIEWED`/`OFFERED` —
    lihat State Machine) dan tetap idempoten pada offer kedua (tidak
    menulis history duplikat kalau `wasAccepted`).
  - `GetApplicationHistory` mengembalikan list terurut `changed_at` ASC.
- **Handler** (`handler_test.go`): `GET .../history` sukses + 404 untuk id
  tak ditemukan.

## Files Touched (ringkasan — detail path & kode ada di implementation plan)

- `backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/097_recruitment_stage_history.sql` (+ `.down.sql`)
- `backend/internal/modules/recruitment/model.go` — `RecruitmentStage`, `ApplicationStageHistory` GORM entity
- `backend/internal/modules/recruitment/repository.go` — `CreateStageHistory`, `ListStageHistoryByApplication`, `FindStageByCode`
- `backend/internal/modules/recruitment/service.go` — `transitionApplicationStatus`, `GetApplicationHistory`, ubah `UpdateApplicationStatus` + `AcceptOffer` + `CreateApplication`
- `backend/internal/modules/recruitment/dto.go` — `StageHistoryResponse`
- `backend/internal/modules/recruitment/handler.go` + `routes.go` — `GET /applications/:id/history`
- `backend/internal/modules/recruitment/module.go` — seed 8 `recruitment_stages` + `AutoMigrate` 2 entity baru
- `docs/module-recruitment-development-plan.md` — update status G-5 setelah selesai
