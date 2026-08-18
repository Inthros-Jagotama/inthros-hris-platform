# Design: Brand Color Palette (Navy / Teal / Orange / Neutral)

## 1. Tujuan

Mengganti palet warna brand di app `frontend/tenant` (207 file `.vue`) dari warna default Tailwind (indigo sebagai warna utama link/aktif, gray sebagai teks) menjadi palet brand baru:

- **Navy `#1B2A41`** — teks utama (heading/primary text), elemen navigasi (Sidebar/AppLayout).
- **Teal `#1B7F93`** — brand sekunder: link, ikon/badge aktif, badge info.
- **Oranye `#F5941E`** — CTA utama, notifikasi, highlight penting. Dipakai terbatas, bukan pengganti massal.
- **Neutral (gray/slate existing)** — background, card, border, teks sekunder. **Tidak diubah.**

`frontend/platform-admin` **di luar scope** — dikerjakan terpisah nanti bila diperlukan.

## 2. Kondisi Saat Ini

- Tailwind v4 (`@import "tailwindcss"` di `frontend/tenant/src/assets/styles/main.css`), CSS-first — **tidak ada** `tailwind.config.js` maupun `@theme` block custom. Semua warna pakai skala default Tailwind (indigo/emerald/amber/rose/sky/gray/teal/orange dst).
- PrimeVue pakai preset **Aura** bawaan tanpa `definePreset` custom (`frontend/tenant/src/main.js`) — Button/Input/Menu/dll memakai token warna default Aura (biru).
- Dark mode: class-based (`.dark` di `<html>`, toggle via `frontend/tenant/src/stores/theme.js`), PrimeVue sinkron lewat `.p-dark`. **Harus tetap berfungsi** setelah perubahan.
- Warna semantic status (sukses=emerald, warning=amber, danger=rose, info=sky/blue) dipakai luas untuk badge/tag status — **dipertahankan, tidak disentuh** (keputusan user).
- Indigo dipakai di 34 file sebagai warna brand de-facto (link, ikon aktif, focus ring, badge info) — ini yang digantikan Teal.
- Teal sudah dipakai ringan di 14 file untuk kebutuhan lain (bukan status semantic) — aman di-override karena tidak bentrok makna.
- Pola teks utama yang berulang: `text-gray-800 dark:text-gray-100` (heading/label penting) muncul di hampir semua file — hanya varian **light-mode** (`text-gray-800`) yang diganti ke Navy; varian dark-mode (`dark:text-gray-100`) **tetap** (teks terang di atas background gelap, bukan navy gelap yang tak terbaca).
- Tidak ada file design-system dokumentasi maupun test frontend (dikonfirmasi saat eksplorasi) — tidak ada test yang perlu disesuaikan untuk perubahan visual ini.

## 3. Desain Token

### 3.1 Tailwind `@theme` (di `frontend/tenant/src/assets/styles/main.css`)

Tambahkan/override 3 skala warna (masing-masing 50–950, mengikuti struktur skala Tailwind bawaan supaya pola `bg-{color}-50 dark:bg-{color}-500/10 text-{color}-600 dark:text-{color}-400` yang sudah dipakai luas tetap jalan tanpa perlu restrukturisasi tiap komponen):

```css
@theme {
  --color-navy-50: ...;
  --color-navy-100: ...;
  ...
  --color-navy-500: #1B2A41;  /* base */
  ...
  --color-navy-950: ...;

  /* override skala teal bawaan supaya base match brand */
  --color-teal-500: #1B7F93;
  --color-teal-600: ...;
  ... (skala lain di-generate proporsional dari base)

  /* override skala orange bawaan supaya base match brand */
  --color-orange-500: #F5941E;
  --color-orange-600: ...;
  ...
}
```

Nilai shade 50–950 di-generate dari base hex memakai tint/shade standar (bukan dipilih manual satu-satu) supaya konsisten dengan look Tailwind default. Ditentukan saat implementasi (bagian dari langkah teknis di plan, bukan spec ini).

### 3.2 PrimeVue Preset

Tambah `definePreset` baru di `frontend/tenant/src/main.js` (extend dari `Aura`), set token `primary` → skala `teal` (karena primary PrimeVue dipakai untuk Button/Input focus/Menu active — perannya paling dekat dengan "link, ikon aktif" bukan CTA). Efeknya otomatis menjalar ke seluruh komponen PrimeVue built-in (Button default, Checkbox, InputText focus ring, Menu active item, dll) tanpa perlu override manual per file.

CTA/notification oranye **tidak** dijadikan primary PrimeVue (supaya tidak "meng-oranye-kan" semua button/checkbox secara default) — oranye diterapkan eksplisit per elemen lewat class Tailwind (`bg-orange-500`, dll) hanya di titik CTA/notifikasi yang benar-benar dimaksud (misal tombol "Ajukan"/"Simpan" utama, badge notifikasi unread).

## 4. Aturan Pemetaan (dipakai konsisten saat migrasi per-file)

| Elemen | Sebelum | Sesudah |
|---|---|---|
| Heading/label utama (light mode) | `text-gray-800` | `text-navy-800` (dark mode tetap `dark:text-gray-100`, tidak diubah) |
| Link, ikon/background aktif, focus ring | `indigo-*` (`bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400`, `hover:border-indigo-300`, `focus-visible:ring-indigo-500/50`) | `teal-*` pola sama |
| Badge info | `indigo-*` badge | `teal-*` badge |
| Elemen navigasi (Sidebar/AppLayout: background/text/active item) | gray/indigo campuran | Navy sebagai warna dominan nav |
| CTA utama (tombol primary aksi penting, mis. submit/ajukan), badge notifikasi | bervariasi (default PrimeVue biru, atau tidak eksplisit) | `orange-500`/`orange-600`, dipakai selektif per komponen, bukan blanket replace |
| Status badge (sukses/warning/danger/info) | emerald/amber/rose/sky | **tidak diubah** |
| Background, card, border, teks sekunder | gray/slate | **tidak diubah** |

## 5. Urutan Eksekusi (Batch)

**Batch 1 — Fondasi + dampak luas** (dikerjakan & direview lebih dulu, karena mengubah tampilan besar tanpa harus sentuh tiap file):
1. Tambah `@theme` token (navy/teal/orange) di `main.css`.
2. Tambah `definePreset` PrimeVue (primary → teal) di `main.js`.
3. Migrasi `AppLayout.vue` + komponen Sidebar/navigasi terkait → Navy untuk background/text/active-nav-item.
4. Verifikasi build + smoke-check visual light & dark mode (dev server, cek beberapa halaman).

**Batch 2+ — Per modul** (`frontend/tenant/src/views/modules/*`), urutan berdasarkan direktori existing:
```
attendance → leave → payroll → performance → recruitment
→ reimbursement → training → employee → employeemovement
→ organization → job / jobvalues → approval → notification
→ competency → career-intelligence → workforce-intelligence
```
Tiap batch: grep literal `indigo-`/`text-gray-800` (light-mode only) di direktori modul tsb, terapkan aturan §4, build, commit terpisah per batch supaya mudah direview/di-revert per modul.

Komponen shared (`frontend/tenant/src/components/**`, dipakai lintas modul) dikerjakan **sebelum** batch modul dimulai (bagian akhir Batch 1 atau awal Batch 2) karena perubahannya berdampak ke banyak halaman sekaligus.

## 6. Di Luar Scope

- `frontend/platform-admin` (app terpisah).
- Warna status/semantic (sukses/warning/danger/info).
- Background/card/border/teks sekunder (neutral, tidak berubah).
- Redesain layout/komponen — murni penggantian warna, bukan restrukturisasi UI.

## 7. Risiko & Mitigasi

- **Override skala `teal`/`orange` bawaan Tailwind mengubah tempat lain yang kebetulan pakai warna itu** (14 file untuk teal, 4 file untuk orange) — diterima sebagai efek samping kecil & dicek manual saat Batch 1 (bukan blocker, karena bukan dipakai untuk semantic status).
- **Kontras Navy di dark mode** — dimitigasi dengan aturan eksplisit di §4: hanya varian light-mode yang diganti, dark-mode tetap pakai neutral terang.
- **PrimeVue primary → teal berdampak ke semua Button/Input default** — disengaja (itu tujuannya), tapi perlu smoke-test visual di Batch 1 sebelum lanjut ke batch berikutnya.
