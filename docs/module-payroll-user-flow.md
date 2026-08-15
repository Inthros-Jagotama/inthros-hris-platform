# Panduan Pengisian Module Payroll (Runbook)

> 📅 2026-08-15 · Dokumen ini **berdiri sendiri** — panduan praktis langkah demi langkah memakai module payroll, dari nol sampai laporan.
> Untuk detail teknis/arsitektur lihat [docs/payroll/00-overview.md](payroll/00-overview.md) dan sub-plan 01–08; versi ringkas alur ada di [docs/payroll/09-user-flow.md](payroll/09-user-flow.md).

---

## Gambaran: 3 Tahap

```text
TAHAP 1                     TAHAP 2                    TAHAP 3
SETUP (sekali)        →     PROSES (tiap periode) →    DISTRIBUSI
─────────────────           ────────────────────       ─────────────────
A. Master Data Gaji         D. Buat Payroll Run        F. Payslip
B. Pengaturan BPJS/Pajak    E. Calculate → Final       G. Pembayaran
C. Profil Karyawan                                    H. Laporan
```

> **Aturan emas:** Tahap 1 cukup dikerjakan sekali (atau saat ada perubahan). Tahap 2–3 diulang setiap periode gajian.

---

## TAHAP 1 — SETUP (dikerjakan sekali)

### A. Master Data Gaji

**1. Payroll Period** — menu `Settings → Payroll Periods`
- [ ] Klik **New Period** → isi **Tahun** (pilih tahun), **Bulan**, **Tanggal Mulai**, **Tanggal Selesai**, **As-Of Date**.
- [ ] As-Of Date = tanggal cut-off pembacaan data employment & tarif. Status biarkan **Terbuka (OPEN)**.
- [ ] Setelah periode beres, periode ditutup (**CLOSED**) agar terkunci.

**2. Salary Components** — menu `Settings → Salary Components`
- [ ] Komponen dikelompokkan per tipe (card): **Earning, Deduction, Employer Contribution, Information**.
- [ ] Klik **New Component** pada card tipe yang diinginkan → tipe otomatis terisi di judul modal.
- [ ] Isi **Nama** → **kode dibuat otomatis** dari inisial nama, contoh `Gaji Pokok` → `GP_20260815_143205_AB12`.
  ⚠️ Kode memakai `_` (bukan `-`) karena `-` dibaca operator minus oleh formula engine.
- [ ] **Tipe Perhitungan** (radio): `Fixed` / `Percentage` / `Formula` / `Reference` / `Manual`.
  - Untuk Formula/Percentage: tulis formula, **ketik `@`** di dalam field → cari komponen berdasarkan nama → klik, hanya kode yang masuk. Contoh: `BASIC_SALARY + TUNJANGAN_MAKAN`. Klik **Validate** untuk cek.
- [ ] Switch penting: **Taxable** (dasar PPh21), **BPJS Base** (dasar upah BPJS), **Recurring** (tiap periode), **Proratable** (diprorata saat join/resign), **Print on Structure** (tampil di payslip).
- [ ] **Komponen PPh21** (baris potongan hasil pajak di payslip — cukup 1 komponen DEDUCTION) dan **Pengurang PPh21** (mis. iuran pensiun — bisa lebih dari satu). Engine PPh21 mengikuti flag ini, jadi tidak perlu pilih komponen lagi di PPh21 Settings.
- [ ] **Display Order** (paling bawah) menentukan urutan tampil.

**3. Salary Structure** — menu `Settings → Salary Structure`
- [ ] **Tab Grade Components**: default komponen per grading, mis. Grade 5 → `GAJI_POKOK = 10.000.000`. Pilih grading (opsional), komponen, nominal, tanggal efektif, wajib/default.
- [ ] **Tab Employee Components**: override per karyawan bila ada karyawan yang nilainya beda dari grade-nya. Pilih employee, komponen, nominal, sumber (Manual / Grade Inherit / Formula / Adjustment).

### B. Pengaturan BPJS & Pajak

**4. BPJS Settings** — menu `Settings → BPJS Settings`
- [ ] Buat setting BPJS + **Rate Components** per program: Kesehatan, JHT, JKK, JKM, JP, JKP (tarif % employee/employer, batas maksimum/cap, class risiko JKK). Berlaku efektif per tanggal.

**5. PPh21 Settings** — menu `Settings → PPh21 Settings`
- [ ] **Setting PPh21**: **Metode Perhitungan** — pilih **TER (default, PP 58/2023)** atau Reguler (Annualized Gross); biaya jabatan (rate + maks), metode annualisasi, multiplier non-NPWP, pembulatan. *(Tidak perlu pilih komponen — komponen PPh21 & pengurang ditentukan dari flag di Salary Components.)*
- [ ] *Catatan metode TER:* Jan–Nov pajak = **bruto bulanan × tarif TER** (kategori A/B/C dari status PTKP, tarif otomatis dari tabel tarif efektif); Desember dihitung metode normal dikurangi potongan Jan–Nov.
- [ ] **Tarif TER**: tab TER — rentang bruto bulanan per kategori A/B/C (kelola tarif efektif; dipakai metode TER Jan–Nov).
- [ ] **PTKP Rates**: tab Pajak — kelola tabel `ptkps` (nama status, grup A/B/C, nominal tahunan). Satu-satunya pengelolaan PTKP; engine PPh21 membacanya langsung (tabel `pph21_ptkp_rates` dihapus).
- [ ] **Tax Brackets**: lapisan tarif progresif (5%, 15%, 25%, dst.).

### C. Profil Karyawan

**6. Payroll Profile** — menu `Settings → Employee Payroll Profiles`
- [ ] Pilih employee, group code, frekuensi (MONTHLY/WEEKLY/DAILY), metode bayar (BANK_TRANSFER/CASH/CHEQUE), currency, aktif, tanggal efektif.
- [ ] Hanya profil **aktif** (is_payroll_active + ACTIVE + dalam rentang tanggal efektif) yang dihitung.

**7. Data pendukung** (harus dibuat setelah payroll profile):
- [ ] **Bank Profiles** — `Settings → Bank Profiles`: rekening untuk transfer gaji (akan disnapshot saat payment batch dibuat).
- [ ] **BPJS Profiles** — `Settings → BPJS Profiles`: nomor kepesertaan Kesehatan/Ketenagakerjaan, class risiko JKK, JP aktif.
- [ ] **Tax Profiles** — `Settings → Tax Profiles`: NPWP, status PTKP, metode pajak (GROSS/GROSS_UP/NETT).

> ✅ **Cek siap proses:** periode OPEN ada · komponen gaji yang dipakai berstatus ACTIVE · minimal 1 employee punya payroll profile aktif · setting BPJS/PPh21 berlaku di as-of date periode.

---

## TAHAP 2 — PROSES (setiap periode)

Menu **Payroll** (sidebar).

**8. Buat Run**
- [ ] Klik **New Run** → pilih **Periode** + **Metode Prorasi**:
  - `CALENDAR_DAYS` — proporsi hari kalender (default)
  - `WORKING_DAYS` — proporsi hari kerja
  - `FIXED_30_DAYS` — selalu / 30 hari
  - `ATTENDANCE_DAYS` — proporsi hari hadir
- [ ] Status run: **DRAFT**.

**9. Calculate**
- [ ] Klik **Calculate** → sistem menghitung otomatis (urut):
  1. Struktur gaji (grade → override employee → adjustment periode)
  2. Formula komponen (variabel built-in: `WORKING_DAYS`, `WORKED_DAYS`, `ABSENCE_DAYS`, `UNPAID_LEAVE_DAYS`, `OVERTIME_HOURS`)
  3. Prorasi join/resign tengah bulan
  4. BPJS (dari komponen `is_bpjs_base` + rate component)
  5. PPh21 (dari komponen `is_taxable` − pengurang `is_pph21_deductible` − iuran BPJS + PTKP + bracket)
- [ ] Hasil **disnapshot** — mengubah data sumber setelah ini tidak mengubah angka run. Status → CALCULATED/REVIEWED.
- [ ] Ada yang salah? Perbaiki data → **Calculate ulang** (snapshot lama dihapus, aman).

**10. Review → Approve → Lock**
- [ ] Cek angka di tab **Employees / Items / Reports**.
- [ ] **Approve** (APPROVED) → **Lock** (LOCKED, terkunci, tidak bisa dihitung ulang). **Cancel** untuk membatalkan run.

---

## TAHAP 3 — DISTRIBUSI

**11. Payslip** — tab **Payslips** di detail run
- [ ] **Generate Payslips** (satu per employee, nomor `SLP-<periode>-<seq>`).
- [ ] **Publish** → payslip bisa dilihat karyawan (status PUBLISHED); **Cancel** membatalkan.
- [ ] Klik payslip → lihat **HTML payslip** (earning, potongan, iuran perusahaan, gaji bersih, biaya perusahaan).

**12. Pembayaran** — tab **Payments**
- [ ] **Create Batch** → satu baris per employee, nominal = **gaji bersih**, rekening dari **snapshot bank profile**. Karyawan tanpa rekening dilewati & dilaporkan.
- [ ] Ubah status: **PENDING → PROCESSING → PAID** (atau FAILED / REVERSED) dengan timestamp & referensi.
- [ ] **Export CSV** untuk file transfer bank.

**13. Laporan** — tab **Reports**
- [ ] **Summary** (total employee, gross, potongan, kontribusi perusahaan, net, biaya perusahaan)
- [ ] **Detail** (per employee per komponen)
- [ ] **BPJS** (upah dasar + kontribusi per program) · **Tax** (PKP + PPh21) · **Bank** (daftar transfer)
- [ ] **Overview** = dashboard agregat run.

---

## Ringkasan Status

| Objek | Status | Keterangan |
|---|---|---|
| Payroll Period | OPEN / CLOSED | dipakai saat OPEN |
| Payroll Run | DRAFT → CALCULATED → REVIEWED → APPROVED → LOCKED | LOCKED = terkunci; CANCELLED = batal |
| Payslip | DRAFT → PUBLISHED → CANCELLED | setelah run final |
| Payment | PENDING → PROCESSING → PAID / FAILED / REVERSED | setelah payslip |

## Troubleshooting Singkat

| Gejala | Kemungkinan penyebab |
|---|---|
| Employee tidak muncul di run | Profil payroll tidak aktif / tanggal efektif tidak mencakup as-of date |
| Komponen tidak terhitung | Status komponen bukan ACTIVE, atau tidak terpasang di struktur grade/employee |
| PPh21 tidak muncul | Belum ada tax profile (status PTKP) / setting PPh21 / komponen `is_taxable` |
| BPJS tidak muncul | Belum ada setting BPJS + rate component aktif / komponen `is_bpjs_base` |
| Kode komponen aneh di formula | Pakai `_` bukan `-` agar tidak jadi operator minus |
| Angka berubah sendiri | Data sudah di-snapshot — hitung ulang hanya bila run belum LOCKED |
| Rekening berubah setelah run final | Batch pembayaran tetap pakai rekening lama (snapshot) — aman |
