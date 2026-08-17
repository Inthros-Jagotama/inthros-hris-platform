# 📚 Dokumentasi HRIS Platform

Index pusat seluruh dokumentasi proyek. Semua dokumen (kecuali `README.md` di root) tersimpan di folder `docs/` ini dan saling merujuk.

## 🌐 Navigasi Cepat

| Kategori | Dokumen | Isi |
|---|---|---|
| **Penggunaan** | [`api/api-usage-guide.md`](api/api-usage-guide.md) | Cara pakai API: menjalankan server, autentikasi, format response, contoh curl (KPI & OKR), error codes |
| **Arsitektur** | [`platform-architecture-design.md`](platform-architecture-design.md) | Satu-satunya dokumen arsitektur: modular monolith, multi-tenant, provisioning engine |
| **Deployment** | [`deployment-guide.md`](deployment-guide.md) | Deployment Subscription SaaS & On-Premise (`.lic` RSA), instalasi, troubleshooting |
| **API Reference** | [`openapi-report.md`](openapi-report.md) | Laporan OpenAPI komprehensif per modul *(generated via `make docs`)* |
| **Database** | [`database-schema.md`](database-schema.md) | Struktur database & ERD: Platform DB (11 tabel) + Tenant DB (199 tabel), relasi FK, konvensi kolom |
| **Laporan Teknis** | [`go-module-architecture-report.md`](go-module-architecture-report.md) | Laporan arsitektur Go module (entities, services, tests per modul) |
| **Progres Proyek** | [`project-completion-dashboard.md`](project-completion-dashboard.md) | Dashboard penyelesaian proyek (modul, test coverage, infrastruktur) |
| **Frontend** | [`frontend-development-plan.md`](frontend-development-plan.md) | Roadmap implementasi frontend Platform Admin & Tenant |
| **Frontend** | [`panduan-uiux-hris-enterprise.md`](panduan-uiux-hris-enterprise.md) | Standar UI/UX enterprise: modal-first, high-density, master prompt AI, warna badge |
| **Analisis** | [`analisis-blueprint-vs-existing.md`](analisis-blueprint-vs-existing.md) | Perbandingan blueprint vs existing Laravel app (inthros-web) |
| **Analisis** | [`job-management-score-analysis.md`](job-management-score-analysis.md) | Analisa perhitungan Job Management Score (dirujuk `calculator.go`) |
| **Plan Recruitment** | `module-recruitment-development-plan.md` *(di-archive: `docs/archive/`)* · [`module-recruitment-strategic-layer-plan.md`](module-recruitment-strategic-layer-plan.md) · [`flow/module-recruitment-flow.md`](flow/module-recruitment-flow.md) | Plan pengembangan Recruitment (**operasional**, selesai — diarsipkan) + Strategic Layer Integration (**Workforce/Career Intelligence ↔ Recruitment**, S-1 s.d. S-7) + alur pengguna. Plan modul lain dikelola di [`project-completion-dashboard.md`](project-completion-dashboard.md) → *Module Development Plan Status* |
| **Alur Pengguna** | [`flow/module-recruitment-flow.md`](flow/module-recruitment-flow.md) · [`flow/module-payroll-user-flow.md`](flow/module-payroll-user-flow.md) · [`flow/module-reimbursement-flow.md`](flow/module-reimbursement-flow.md) · [`flow/module-competency-360-flow.md`](flow/module-competency-360-flow.md) · [`flow/module-leave-flow.md`](flow/module-leave-flow.md) | Alur bisnis & runbook pengisian modul: Recruitment, Payroll, Reimbursement, Competency 360, Leave |

## 🔗 Alur Referensi Antar Dokumen

```text
docs/README.md (index)
   │
   ├──→ api/api-usage-guide.md ──────→ openapi-report.md
   ├──→ database-schema.md ──────────→ platform-architecture-design.md
   ├──→ platform-architecture-design.md ──→ deployment-guide.md
   ├──→ deployment-guide.md ────────→ platform-architecture-design.md
   ├──→ frontend-development-plan.md ─→ panduan-uiux-hris-enterprise.md
   ├──→ project-completion-dashboard.md
   └──→ go-module-architecture-report.md
```

> ⚠️ Folder `archive/`, `backlog/`, `seeder/`, `tmp/` **DIABAIKAN** — lihat [🚫 Folder yang Diabaikan](#-folder-yang-diabaikan) di bawah.

## 🗂️ Struktur Folder

```text
docs/
├── README.md                        # ← Index ini (mulai baca dari sini)
├── api/
│   └── api-usage-guide.md           # Panduan penggunaan API
├── flow/
│   ├── module-recruitment-flow.md   # Alur bisnis Recruitment
│   ├── module-payroll-user-flow.md  # Runbook pengisian Payroll (setup → proses → distribusi)
│   ├── module-reimbursement-flow.md # Runbook pengisian Reimbursement (type → request → approval → PAID)
│   ├── module-competency-360-flow.md# Runbook Competency 360 (master → event → rater → hasil)
│   └── module-leave-flow.md         # Runbook pengisian Leave (setup → request → approval → saldo)
├── database-schema.md               # Struktur database & ERD
├── platform-architecture-design.md  # Arsitektur utama
├── deployment-guide.md              # Panduan deployment
├── openapi-report.md                # Generated: make docs
├── go-module-architecture-report.md
├── project-completion-dashboard.md
├── frontend-development-plan.md
├── panduan-uiux-hris-enterprise.md
├── analisis-blueprint-vs-existing.md
├── job-management-score-analysis.md
├── module-recruitment-strategic-layer-plan.md
├── archive/                         # 🚫 DIABAIKAN — plan lama yang sudah selesai
├── backlog/                         # 🚫 DIABAIKAN — ide/backlog, bukan dokumentasi aktif
├── seeder/                          # 🚫 DIABAIKAN — file seeder/export sementara
└── tmp/                             # 🚫 DIABAIKAN — file temp & node_modules (besar)
```

## 🚫 Folder yang Diabaikan

Folder berikut **tidak boleh dibaca/diindeks** saat memproses dokumentasi (termasuk oleh AI/tooling/link checker):

| Folder | Isi | Alasan Diabaikan |
|---|---|---|
| `docs/archive/` | Plan/rencana modul yang sudah selesai diimplementasikan | Referensi historis — gunakan dokumen aktif di `docs/` |
| `docs/backlog/` | Ide/backlog fitur | Bukan dokumentasi resmi — tidak dipakai |
| `docs/seeder/` | File seeder/export data sementara | File kerja — tidak relevan untuk dokumentasi |
| `docs/tmp/` | File temp, `inthros-web/`, `node_modules/` | File sementara & dependency (besar) — memperlambat pencarian |

**Aturan:**
- Jangan sertakan isi folder-folder di atas dalam ringkasan/indeks/tabel referensi.
- Jangan jalankan pencarian (`grep`/`rg`) di dalamnya — gunakan flag exclude (mis. `--glob '!docs/tmp/**'`, `--glob '!docs/archive/**'`).
- Jika ada dokumen di `archive/` yang masih relevan, pindahkan ke root `docs/` (dengan penamaan kebab-case) dan daftarkan di tabel navigasi di atas — jangan dibaca dari lokasi lama.

## 📝 Standar Penamaan Dokumen

Semua file dokumentasi mengikuti konvensi **`kebab-case` (huruf kecil + tanda hubung)**, kecuali `README.md` (nama khusus Git/GitHub).

| ✅ Benar | ❌ Salah |
|---|---|
| `project-completion-dashboard.md` | `PROJECT_COMPLETION_DASHBOARD.md` |
| `phase-1-completion-report.md` | `Phase-1-Completion-Report.md` |
| `api-usage-guide.md` | `ApiUsageGuide.md` / `api_usage_guide.md` |
| `deployment-guide.md` | `Deployment Guide.md` |

Aturan: huruf kecil semua, pisahkan kata dengan `-` (bukan `_` atau spasi), tanpa huruf kapital kecuali `README.md`.

## ⚙️ Maintenance

- **`openapi-report.md`** di-generate otomatis: `cd backend && make docs` (menjalankan `check-openapi` + `generate_openapi_report.py`).
- **`openapi.json`** diperbarui secara manual lewat script inject saat endpoint baru ditambahkan ke routes (pola: `scripts/inject_*_openapi.py`, mis. `inject_missing_endpoints_openapi.py` untuk modul lama, `inject_performance_kpi_okr_openapi.py` untuk KPI/OKR, `inject_training_employee_openapi.py` untuk Training P0-P2 + Employee Movement enhancement + Attendance overtime + Career Intelligence paths + KPI duplicate + Workforce candidate-search). Setelah inject, jalankan `make check-openapi` (harus 0 missing) lalu `make docs` untuk regenerate report.
- **`database-schema.md`** di-generate dari migrasi SQL: `cd backend && make db-docs`. Verifikasi sinkronisasi tanpa menimpa file: `make check-db-docs` (menjalankan `scripts/check_database_schema_doc.py`). Check memvalidasi **kedua dialect** — daftar tabel di dokumen harus sama dengan migrasi `postgres/` **dan** `mysql/` (199 tabel masing-masing).
- **`go-module-architecture-report.md`** di-generate dari analisis statis kode Go: `cd backend && make arch-report` (menjalankan `scripts/generate_go_module_report.py`) — menghitung entities (struct di `model.go`), service/repo/handler methods, route registrations, dan `func Test` per modul (tenant, platform, shared pkg). Verifikasi sinkronisasi tanpa menimpa file: `make check-arch-report` (menjalankan `scripts/check_go_module_report.py`).
- **Gap check lintas dokumentasi**: jalankan `python scripts/check_doc_gaps.py` dari root proyek — memeriksa (1) link internal rusak, (2) referensi inline `docs/...` yang menunjuk file tidak ada, dan (3) daftar file markdown di `docs/` root untuk dicek terhadap tabel navigasi di atas. Folder yang diabaikan (`tmp/`, `archive/`, `backlog/`, `seeder/`) otomatis dikecualikan.
- Tambahkan dokumen baru di folder `docs/` dengan penamaan **kebab-case** (lihat standar di atas), lalu perbarui index ini dan referensi silang di dokumen terkait.
- Saat merename/memindahkan file, perbarui **semua** referensi (gunakan `git grep '<nama-lama>'` untuk menemukannya).
- Semua dokumen di luar `docs/` (kecuali `README.md` root) harus dipindahkan ke sini — lihat juga `README.md` root bagian "Dokumentasi Lainnya".
