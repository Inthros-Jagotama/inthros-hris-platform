# Payroll — Panduan Alur Pengisian (User Flow)

> 📅 2026-08-15 · Sinkron dengan implementasi aktual (backend + frontend payroll).
> Dokumen ini menjawab: **"dari mana mulai, isi apa dulu, lalu apa?"** — alur praktis memakai module payroll, bukan spesifikasi teknis (lihat sub-plan 01–08 untuk detail).

## Gambaran Umum — 3 Tahap Besar

```text
TAHAP 1                          TAHAP 2                        TAHAP 3
SETUP MASTER DATA        →       JALANKAN PAYROLL RUN    →      DISTRIBUSI HASIL
──────────────────              ────────────────────           ─────────────────
1. Payroll Period                1. Buat Run (pilih periode)     1. Generate Payslip
2. Salary Components              2. Calculate → Snapshot         2. Publish Payslip
3. Salary Structure               3. Review / Approve / Lock      3. Payment Batch (transfer)
4. BPJS Settings                  4. (Ulang) Re-Calculate         4. Export CSV Bank
5. PPh21 Settings                                                   5. Reports & Dashboard
6. Profil Payroll Karyawan
   (+ Bank / BPJS / Tax)
```

> ⚠️ **Aturan emas:** Tahap 1 cukup diisi **sekali** (atau saat ada perubahan). Tahap 2–3 diulang **setiap periode** gajian.

---

## TAHAP 1 — Setup Master Data (isi sekali, lalu pelihara)

Urutan pengisian penting karena tahap berikutnya memakai data tahap sebelumnya.

### 1.1 Payroll Period — `Settings → Payroll Periods`
- Buat periode gajian: **Tahun** (pilih tahun), **Bulan**, **Tanggal Mulai**, **Tanggal Selesai**, dan **As-Of Date** (tanggal cut-off untuk membaca data employment & tarif).
- Status periode: **Terbuka (OPEN)** untuk dipakai run; **Ditutup (CLOSED)** setelah selesai (terkunci).
- Kode periode dibuat otomatis dari tahun+bulan.

### 1.2 Salary Components — `Settings → Salary Components`
- Komponen gaji dikelompokkan dalam **card per tipe**: Earning (penghasilan), Deduction (potongan), Employer Contribution (kontribusi perusahaan), Information.
- Klik **"+ New Component"** pada card tipe yang diinginkan → tipe otomatis terisi (tampil di judul modal).
- Isi **Nama** → **kode otomatis** dibuat dari inisial nama (`Gaji Pokok` → `GP_20260815_143205_AB12`). Kode memakai `_` (bukan `-`) agar tidak konflik dengan operator minus di formula.
- **Tipe Perhitungan** (radio):
  - `Fixed` — nominal tetap.
  - `Percentage` / `Formula` — tulis formula. **Ketik `@` di dalam field formula** lalu cari komponen berdasarkan nama; klik → hanya kode komponen yang masuk (`BASIC_SALARY + TUNJANGAN_MAKAN`).
  - `Reference` — ambil nilai dari komponen lain (pilih dari daftar).
  - `Manual` — diisi manual per periode.
- Pengaturan switch: **Taxable** (masuk dasar PPh21), **BPJS Base** (masuk upah dasar BPJS), **Recurring** (dibayar tiap periode), **Proratable** (diprorata saat join/resign tengah bulan), **Print on Structure** (tampil di payslip). **Komponen PPh21** (baris potongan hasil pajak — 1 komponen) dan **Pengurang PPh21** (mis. iuran pensiun — bisa banyak) juga diatur di sini; engine PPh21 mengikuti flag ini.
- **Urutan (Display Order)** di paling bawah menentukan urutan tampil.
- Klik **Validate** untuk memastikan formula valid.

### 1.3 Salary Structure — `Settings → Salary Structure`
- **Tab Komponen Grade**: default komponen per **grading** (mis. Grade 5 dapat `GAJI_POKOK = 10.000.000`) — pilih grading (opsional), komponen, nominal, tanggal efektif, flag wajib/default.
- **Tab Komponen Karyawan**: **override per employee** bila perlu (mis. karyawan tertentu dapat lebih tinggi) — pilih employee, komponen, nominal, sumber (manual/warisan grade/formula/penyesuaian), tanggal efektif.

### 1.4 BPJS Settings — `Settings → BPJS Settings`
- Buat setting BPJS (program: **Kesehatan, JHT, JKK, JKM, JP, JKP**) + **Rate Components** per program: tarif % employee/employer, batas maksimal (cap), dan class risiko JKK.
- Effective-dated: tarif lama tetap berlaku untuk periode sebelum tanggal efektif baru.

### 1.5 PPh21 Settings — `Settings → PPh21 Settings`
- **Setting PPh21**: **Metode Perhitungan** — **TER** (default, PP 58/2023: Jan–Nov = bruto × tarif efektif per kategori PTKP; Desember = normal − YTD) atau **Reguler** (annualized gross). Biaya jabatan (rate + maks), metode annualisasi, multiplier non-NPWP, pembulatan. *(Komponen potongan & pengurang PPh21 tidak dipilih di sini — diambil dari flag di Salary Components.)*
- **Tarif TER**: tab TER — tarif efektif per kategori A/B/C (rentang bruto bulanan) untuk metode TER Jan–Nov.
- **PTKP Rates**: tab Pajak — kelola tabel `ptkps` (name, group A/B/C, nominal tahunan). Satu-satunya pengelolaan PTKP & dipakai langsung engine (`pph21_ptkp_rates` dihapus — migration 121).
- **Tax Brackets**: lapisan tarif progresif (5%, 15%, 25%, dst.).

### 1.6 Profil Payroll Karyawan — `Settings → Employee Payroll Profiles`
- Pilih **Employee**, **Payroll Group Code**, **Frequency** (MONTHLY/WEEKLY/DAILY), **Payment Method** (BANK_TRANSFER/CASH/CHEQUE), **Currency** (IDR), aktif/nonaktif, dan **tanggal efektif**.
- Hanya profil **aktif** (is_payroll_active + ACTIVE + dalam rentang tanggal efektif) yang ikut dihitung run.
- Data pendukung per employee — **`Settings → Bank Profiles`** (rekening untuk pembayaran — disnapshot saat payment batch dibuat), **`Settings → BPJS Profiles`** (nomor kepesertaan kesehatan/ketenagakerjaan + class risiko JKK), **`Settings → Tax Profiles`** (NPWP, status PTKP, metode pajak). Semuanya terikat ke **Profil Payroll** karyawan (pilih employee → pilih profil payroll-nya).

> ✅ **Cek kelengkapan sebelum lanjut:** periode OPEN sudah ada; komponen gaji yang dipakai struktur sudah `ACTIVE`; minimal satu employee punya payroll profile aktif; setting BPJS/PPh21 aktif berlaku di as-of date periode.

---

## TAHAP 2 — Jalankan Payroll Run (per periode)

Menu **Payroll** (sidebar) → daftar run.

### 2.1 Buat Run
- Klik **"+ New Run"** → pilih **Periode** (yang sudah dibuat) dan **Metode Prorasi**:
  - `CALENDAR_DAYS` — proporsi hari kalender (default).
  - `WORKING_DAYS` — proporsi hari kerja.
  - `FIXED_30_DAYS` — selalu dibagi 30 hari.
  - `ATTENDANCE_DAYS` — proporsi hari hadir.
- Status run: **DRAFT**.

### 2.2 Calculate
- Klik **Calculate** → sistem menghitung seluruh employee dengan payroll profile aktif:
  1. Struktur gaji (grade default → override employee → adjustment periode).
  2. Evaluasi formula komponen (variabel built-in tersedia: `WORKING_DAYS`, `WORKED_DAYS`, `ABSENCE_DAYS`, `UNPAID_LEAVE_DAYS`, `OVERTIME_HOURS` dari Attendance/Leave/Overtime).
  3. Prorasi join/resign tengah bulan sesuai metode run.
  4. **BPJS** (dari komponen `is_bpjs_base` + rate component) → item potongan & kontribusi.
  5. **PPh21** (dari komponen `is_taxable` − pengurang `is_pph21_deductible` − iuran BPJS + PTKP + bracket) → item potongan + log.
- Hasil **disnapshot** (tidak berubah walau data sumber diubah belakangan). Status run → **CALCULATED/REVIEWED**.
- Perlu ubah sesuatu? **Ulangi Calculate** — snapshot lama dihapus & dihitung ulang dengan aman.

### 2.3 Review → Approve → Lock
- Cek angka di tab **Employees/Items/Reports**; bila ada yang salah, perbaiki data → Calculate ulang.
- **Approve** → status **APPROVED**; **Lock** → **LOCKED** (terkunci, tidak bisa dihitung ulang lagi).
- **Cancel** membatalkan run.

---

## TAHAP 3 — Distribusi Hasil (payslip, pembayaran, laporan)

### 3.1 Generate & Publish Payslip
- Tab **Payslips** pada run → **Generate Payslips** (satu payslip per employee, nomor `SLP-<periode>-<seq>`).
- **Publish** → employee bisa melihat (status PUBLISHED). **Cancel** membatalkan payslip.
- Klik payslip → **HTML payslip** (earning, potongan, iuran perusahaan, gaji bersih, biaya perusahaan).

### 3.2 Payment Batch (transfer bank)
- Tab **Payments** → **Create Batch**: satu baris per employee, nominal = **gaji bersih (net)**, rekening diambil dari **snapshot bank profile** (perubahan rekening setelah run final tidak mengubah batch). Employee tanpa rekening dilewati & dilaporkan.
- Ubah status: **PENDING → PROCESSING → PAID** (atau **FAILED** / **REVERSED**) dengan timestamp & referensi.
- **Export CSV** untuk file transfer bank.

### 3.3 Reports & Dashboard
- Tab **Reports** pada run:
  - **Summary** — total employee, gross, potongan, kontribusi perusahaan, net, biaya perusahaan.
  - **Detail** — per employee per komponen.
  - **BPJS** — upah dasar + kontribusi employee/employer per program.
  - **Tax** — penghasilan kena pajak + PPh21 per employee.
  - **Bank** — daftar pembayaran untuk transfer.
- Tab **Overview** menampilkan dashboard agregat run.

---

## Ringkasan Status

| Objek | Status | Alur |
|---|---|---|
| Payroll Period | OPEN / CLOSED | periode dipakai saat OPEN |
| Payroll Run | DRAFT → CALCULATED → REVIEWED → APPROVED → LOCKED (CANCELLED) | Calculate berulang sampai final; LOCKED = terkunci |
| Payslip | DRAFT → PUBLISHED → CANCELLED | setelah run final |
| Payment | PENDING → PROCESSING → PAID / FAILED / REVERSED | setelah payslip |

## FAQ / Tips Penting

1. **Kode komponen pakai `_`**, jangan `-` — `-` dibaca operator minus oleh formula engine.
2. **Data berubah setelah Calculate?** Run sudah disnapshot — angka tidak berubah; hitung ulang bila perlu (sebelum LOCKED).
3. **Employee tidak muncul di run?** Cek payroll profile-nya: `is_payroll_active`, status ACTIVE, dan tanggal efektif mencakup as-of date periode.
4. **Komponen tidak terhitung?** Pastikan status komponen `ACTIVE`, komponen terpasang di struktur grade/employee (atau adjustment), dan formula valid.
5. **PPh21 tidak muncul?** Employee butuh tax profile (status PTKP) + setting PPh21 aktif + minimal satu komponen `is_taxable`.
6. **BPJS tidak muncul?** Setting BPJS + rate component aktif di as-of date, dan ada komponen `is_bpjs_base`.
7. **Rekening berubah setelah run final?** Batch pembayaran tetap memakai rekening lama (snapshot) — aman.
8. **Belum tersedia:** payslip PDF, payroll journal, payment reconciliation (backend/endpoint belum ada).
