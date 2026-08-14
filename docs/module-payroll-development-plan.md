# Payroll Module — Development Plan (Index)

> 📅 Split 2026-08-14: dokumen ini sebelumnya satu file (1716 baris). Setelah audit terhadap implementasi aktual (2026-08-12), dipecah menjadi beberapa sub-plan di [`docs/payroll/`](payroll/00-overview.md) agar lebih mudah dikelola per area kerja, mengikuti pola yang sama dengan `docs/superpowers/plans/` di modul recruitment.

## Mulai dari sini

👉 **[docs/payroll/00-overview.md](payroll/00-overview.md)** — tujuan, prinsip arsitektur, struktur module, security, multi-company, target arsitektur akhir, dan roadmap prioritas lengkap.

## Daftar sub-plan

| # | File | Isi | Status |
|---|---|---|---|
| 00 | [payroll/00-overview.md](payroll/00-overview.md) | Overview, prinsip, roadmap | referensi |
| 01 | [payroll/01-master-data-selesai.md](payroll/01-master-data-selesai.md) | Salary Component, Salary Structure, Payroll Policy, Database Structure, Approval — CRUD jalan, skema belum lengkap | 🔶 BE / 🚫 FE |
| 02 | [payroll/02-formula-engine.md](payroll/02-formula-engine.md) | Formula Engine | ❌ prioritas #1 |
| 03 | [payroll/03-payroll-run-snapshot.md](payroll/03-payroll-run-snapshot.md) | Payroll Period, Payroll Run, Calculation Order, Snapshot | ❌ prioritas #2 |
| 04 | [payroll/04-bpjs-engine.md](payroll/04-bpjs-engine.md) | BPJS Kesehatan/Ketenagakerjaan kalkulator | 🔶 prioritas #3 |
| 05 | [payroll/05-pph21-engine.md](payroll/05-pph21-engine.md) | PPh 21 rule engine & kalkulator | 🔶 prioritas #4 |
| 06 | [payroll/06-proration-workforce.md](payroll/06-proration-workforce.md) | Proration, integrasi Attendance/Leave/Overtime/Movement/Recruitment/Performance | ❌ prioritas #5-6 |
| 07 | [payroll/07-payslip-payment.md](payroll/07-payslip-payment.md) | Payslip generator, Payment/bank transfer | 🔶/❌ prioritas #7-8 |
| 08 | [payroll/08-reporting-testing.md](payroll/08-reporting-testing.md) | Reporting, Dashboard, Audit Trail, Testing Strategy | ❌ prioritas #9 |

## Ringkasan status (audit 2026-08-12)

Modul payroll **bukan greenfield di backend** — `backend/internal/modules/payroll/` sudah berisi ±6.289 baris kode, 21 GORM entity, 21 tabel, 48 handler, 71 repository function, 43 test. Seluruhnya CRUD master data + status-transition sederhana. **Calculation engine (Formula Engine, kalkulator BPJS/PPh21, eksekusi "Calculate" pada payroll run, payslip generator, payment) belum ada sama sekali** — ini prioritas implementasi berikutnya, detail lengkap ada di sub-plan di atas.

**Frontend (re-verifikasi 2026-08-14): 0% — belum ada implementasi sama sekali.** Halaman `Payroll.vue` hanya placeholder "Coming soon"; route dan menu sidebar sudah terdaftar tapi tidak mengarah ke UI apa pun. Bahkan untuk master data yang backend-nya sudah selesai (§01), belum ada satu pun form/list UI di frontend.
