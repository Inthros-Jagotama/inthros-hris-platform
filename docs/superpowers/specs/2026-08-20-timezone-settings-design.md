# Timezone Settings (Company Default & Zone Override)

## Purpose
Saat ini backend tidak punya konsep zona waktu sama sekali — semua timestamp transaksi (absensi, payroll, cuti, dll) memakai `time.Now()` server tanpa konversi. Fitur ini menambahkan setting zona waktu per tenant (company) dengan opsi override per Zone (cabang/wilayah pada modul `setting`), sehingga logika bisnis berbasis tanggal (batas "hari ini", cutoff, keterlambatan) dan tampilan timestamp mengikuti zona waktu yang relevan, bukan zona server.

## Scope
- Zona waktu dibatasi ke 3 zona Indonesia: WIB (`Asia/Jakarta`), WITA (`Asia/Makassar`), WIT (`Asia/Jayapura`). Nilai disimpan sebagai string IANA, label WIB/WITA/WIT hanya untuk tampilan.
- Dua level konfigurasi:
  1. **Company default** — wajib diisi, berlaku untuk seluruh tenant jika tidak ada override.
  2. **Zone override** — opsional, di-set per baris pada tabel `zones` yang sudah ada (modul `setting`). Employee mewarisi zona lewat rantai `Employee → Organization → Zone`.
- Rollout penerapan ke modul transaksi dilakukan bertahap (lihat bagian Rollout). Spec ini mencakup fondasi (data model, helper resolusi, API, UI) plus penerapan penuh di modul **Attendance** sebagai flagship. Modul lain (Payroll cutoff, Leave) menyusul di iterasi berikutnya menggunakan fondasi yang sama — tidak termasuk dalam implementasi awal.
- Tidak mencakup preferensi zona waktu per user individual (di luar scope; lihat Out of Scope).

## Out of Scope
- Timezone per-user (preferensi tampilan personal terlepas dari zona resmi perusahaan).
- Auto-deteksi zona waktu dari lokasi/IP.
- Dukungan zona waktu di luar Indonesia (bisa ditambah nanti tanpa ubah skema, karena kolom menyimpan string IANA bebas — tapi dropdown UI saat ini dikunci ke 3 opsi).
- Migrasi/backfill data historis lintas zona (timestamp lama tetap diinterpretasikan sesuai perilaku lama; lihat Migrasi Data).

## Data Model

### 1. `companies.timezone` (platform DB)
Migration baru di `backend/internal/pkg/migrator/migrations/platform/`, nomor lanjutan setelah migration terakhir di folder tersebut:
- Kolom `timezone VARCHAR(64) NOT NULL DEFAULT 'Asia/Jakarta'`.
- Existing rows diisi default `Asia/Jakarta` (opsi paling umum) saat migrasi.
- `.down.sql` menghapus kolom.

### 2. `zones.timezone` (tenant DB, mysql + postgres)
Migration baru di `backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/`, nomor lanjutan setelah `155_employee_id_format_settings.sql` (mis. `156_zone_timezone.sql` — nomor pasti ditentukan saat implementasi jika ada migration lain masuk lebih dulu):
- Kolom `timezone VARCHAR(64) NULL` (nullable — `NULL` berarti ikut default company).
- `.down.sql` menghapus kolom.

Kedua migration mengikuti konvensi existing: sequential numeric prefix, dual-dialect untuk tenant DB, disertai `.down.sql`. Tidak pakai AutoMigrate (sesuai kebiasaan proyek untuk modul tenant).

## Resolusi Zona Waktu

Package baru `backend/internal/pkg/timezone/`:

```go
// Resolve mengembalikan *time.Location efektif untuk sebuah organization,
// dengan urutan prioritas: Zone.timezone override -> Company.timezone default -> "Asia/Jakarta".
func Resolve(ctx context.Context, companyTimezone string, zoneTimezone *string) (*time.Location, error)
```

- Dipanggil dengan hasil lookup `Organization.ZoneID → Zone.timezone` (jika ada) dan `Company.timezone` (selalu ada, sudah NOT NULL).
- Hasil `time.LoadLocation` di-cache in-memory (map sederhana, key = nama IANA) karena hanya 3 kemungkinan nilai — menghindari syscall/file lookup berulang.
- Dipakai oleh:
  - **Business logic** yang menentukan batas tanggal (mis. `attendance.service` untuk menentukan "hari ini" saat clock-in/out dan perhitungan keterlambatan).
  - **Presentation layer** (response formatting) bila diperlukan menampilkan timestamp terformat sesuai zona (opsional, fase awal cukup expose `timezone` info di response agar frontend yang format).

Penyimpanan tetap UTC di DB (tidak berubah dari perilaku saat ini) — perubahan hanya di titik *interpretasi* tanggal untuk logika bisnis dan *tampilan*, bukan di penulisan timestamp.

## Migrasi Data (Timestamp Historis)
Timestamp yang sudah ada sebelum fitur ini tidak diubah/dikonversi. Karena backend sebelumnya tidak melakukan konversi zona apapun, secara efektif timestamp lama diperlakukan sebagai sudah dalam UTC apa adanya (tidak ada asumsi baru yang dipaksakan ke data lama). Ini didokumentasikan sebagai catatan known-limitation, bukan bug yang diperbaiki dalam spec ini.

## Backend API
Extend modul `setting` yang sudah ada (`backend/internal/modules/setting/`):

- `GET /api/v1/tenant/settings/company/timezone` → `{ "timezone": "Asia/Jakarta" }` (baca dari Company).
- `PUT /api/v1/tenant/settings/company/timezone` → update `companies.timezone`. Validasi: harus salah satu dari 3 nilai yang diizinkan (`Asia/Jakarta`, `Asia/Makassar`, `Asia/Jayapura`).
- Endpoint Zone yang sudah ada (`create`/`update` Zone) ditambah field `timezone` opsional pada DTO, dengan validasi sama (salah satu dari 3 nilai, atau `null`/kosong untuk "ikut default company").

## Frontend
- Halaman Company Settings (di `frontend/tenant/src/views/settings/`, lokasi tepat mengikuti pola existing seperti `CompanyHolidaysView.vue`): tambah field dropdown "Zona Waktu Perusahaan" dengan 3 opsi berlabel WIB/WITA/WIT, wajib diisi.
- Halaman/form Zone yang sudah ada: tambah field dropdown opsional "Override Zona Waktu" dengan 3 opsi + opsi "Ikut default perusahaan" (mengirim `null`).
- Tidak perlu route baru — field ditambahkan ke form yang sudah ada.

## Penerapan di Modul Attendance (flagship)
Temuan saat investigasi detail: clock-in/out (`CreateEvent` di `service.go`, logika sesi di `session.go`) **sudah** menerima `EventTimeLocal` dari client, dan "work date" sudah ditentukan dari tanggal lokal clock-in versi client — bukan dari `time.Now()` server. 6 titik `time.Now()` yang teridentifikasi di awal eksplorasi ternyata berada di alur approval overtime/koreksi (timestamp instan, bukan penentu boundary tanggal), sehingga tidak relevan untuk resolusi zona waktu.

Titik integrasi flagship yang sebenarnya butuh resolusi zona waktu server-side adalah query **"attendance hari ini"** (dashboard/ringkasan) yang tidak punya input tanggal dari client untuk disandarkan:
1. Tambah `Repository.ResolveOrganizationTimezone(ctx, organizationID)` yang join `organizations → zones` lalu fallback ke `Company.timezone`.
2. Query "hari ini" memakai `time.Now().In(loc)` alih-alih `time.Now()` mentah.
3. Penyimpanan timestamp tetap UTC (tidak berubah).

## Rollout Modul Lain (di luar implementasi awal)
Fondasi (`pkg/timezone`, kolom DB, API, UI) dibangun generik agar modul lain tinggal memanggil `timezone.Resolve` di titik yang relevan:
- Fase 2: Payroll cutoff & perhitungan lembur.
- Fase 3: Leave/cuti (tanggal pengajuan, approval).
- Fase 4: Modul sisanya + audit menyeluruh pemakaian `time.Now()` (60 file teridentifikasi saat eksplorasi) untuk memastikan tidak ada boundary tanggal lain yang terlewat.

Setiap fase adalah unit kerja terpisah, bisa di-plan dan di-implementasi independen setelah fondasi ini selesai.

## Testing
- Backend unit test untuk `timezone.Resolve`: override zone menang atas default company; fallback ke company saat zone timezone `nil`; fallback ke `Asia/Jakarta` saat company timezone kosong (defensif, meski kolom NOT NULL).
- Backend unit test untuk validasi endpoint (menolak nilai timezone di luar 3 opsi yang diizinkan).
- Backend test untuk Attendance: clock-in mendekati tengah malam menghasilkan tanggal absensi yang benar sesuai zona employee (kasus WIB vs WIT beda ~2 jam, uji boundary lintas hari).
- Manual verification: set company timezone ke WITA, buat Zone dengan override WIT, assign Organization ke Zone tsb, lalu clock-in karyawan di organization itu mendekati jam 23:00–01:00 WIT dan pastikan tanggal absensi tercatat sesuai WIT, bukan WITA/UTC.
