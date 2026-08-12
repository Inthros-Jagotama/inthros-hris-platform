# Panduan & Standar UI/UX HRIS Enterprise Grade

> **Tech Stack:** Vue 3 (`<script setup>`) + PrimeVue v4+ + Tailwind CSS  
> **Konsep:** High-Density, Compact, Modal-First, Single-Page Dashboard (Bebas Pindah Halaman), Dark Mode, & Bilingual EN/ID  
> **Implementasi aktual:** `frontend/platform-admin` & `frontend/tenant` (pola komponen & kelas Tailwind di bawah mengikuti kode yang sudah berjalan)

---

## 1. Konsep Utama & Prinsip Desain

Dalam sistem HRIS skala Enterprise, produktivitas pengguna sangat bergantung pada kecepatan eksekusi data. **Prinsip Utama:** Pengguna tidak boleh sering berpindah halaman penuh (*page navigation/routing*) hanya untuk melihat detail, mengedit data, atau melakukan persetujuan (*approval*).

| Prinsip UX | Implementasi Teknikal | Manfaat Utama |
| :--- | :--- | :--- |
| **Modal & Dialog First** | `Dialog` (Modal) dari PrimeVue — pola dominan di codebase (CRUD & approval). `Drawer` (Slide-over) disediakan sebagai opsi form kompleks namun belum banyak dipakai. | Mencegah perpindahan halaman & menjaga konteks pengguna di tabel utama. |
| **Single Viewport (Zero Scroll)** | Root `flex h-screen overflow-hidden` (AppLayout) + sub-halaman detail via route (bukan TabView — komponen `Tabs` belum dipakai di codebase). | Semua informasi penting terlihat tanpa perlu *scroll* panjang; detail dibuka di halaman/route terpisah. |
| **High Density Layout** | Padding ketat (`py-1.5`, `px-3`), font size `12px`–`14px` (`text-xs`–`text-sm`), `rounded-md`/`rounded-sm`. | Memaksimalkan data yang tampil dalam 1 layar. |
| **Dark Mode** | Variant `dark:` di semua komponen + selector `.p-dark` (PrimeVue Aura) + store `theme.js` (localStorage + system pref). | Konsisten antara mode terang/gelap; tidak ada halaman yang terlihat rusak di dark mode. |
| **Bilingual EN/ID** | `useI18n` + locale `en.json`/`id.json` + language switcher di HeaderBar (header `Accept-Language`). | Seluruh label/status/pesan mengikuti bahasa aktif tanpa reload. |

---

## 2. Hirarki Penggunaan Modal & Overlay (Kapan Menggunakan Apa?)

Agar tidak terjadi *modal-cluttering* (terlalu banyak popup menumpuk), ikuti aturan penggunaan overlay berikut:

### A. Compact `Dialog` (Modal Standar)
* **Penggunaan:** Form ringkas (1-2 kolom), konfirmasi aksi (*Approve/Reject*), ubah status, atau input catatan cepat.
* **Properti PrimeVue:** `<Dialog modal :style="{ width: '400px' }" header="...">`
* **Keunggulan:** Fokus instan ke aksi spesifik tanpa mengganggu layar latar belakang.

### B. Slide-over `Drawer` (Modal Samping)
* **Penggunaan:** Form data kompleks (misal: *Tambah Karyawan Baru*, *Detail Gaji & Komponen*) — **standar opsional**; belum banyak dipakai di codebase saat ini (dominasi `Dialog`).
* **Properti PrimeVue:** `<Drawer v-model:visible="open" position="right" class="!w-[500px]">`
* **Keunggulan:** Memberikan area vertikal lebih luas untuk form panjang tanpa menutupi seluruh layar utama.

### C. `Popover` / Context Menu
* **Penggunaan:** Panel notifikasi di HeaderBar (bell dropdown — `Popover`), menu tindakan cepat pada baris tabel (Edit, Hapus, dst.). Menu pengguna di HeaderBar memakai PrimeVue `<Menu popup>` (bukan `Popover`).
* **Keunggulan:** Menghemat kolom pada `DataTable` sehingga tabel tidak terasa sesak; panel ringan tanpa perpindahan halaman.

---

## 3. Master Prompt untuk AI Assistant (Diperbarui)

Gunakan prompt acuan di bawah ini untuk disalin ke AI (ChatGPT / Claude / Gemini) agar seluruh kode yang dihasilkan selalu konsisten dengan standar HRIS Anda:

```text
[CONTEXT & ROLE]
Anda adalah Senior Full-Stack Frontend Engineer & UX Designer spesialis sistem Enterprise ERP/HRIS. Tugas Anda adalah membantu saya membangun aplikasi Human Resource Information System (HRIS) Enterprise Grade yang modern, minimalis, efisien ruang (high-density), user-friendly, bebas dari scroll berlebih, dan SANGAT MINIM PINDAH HALAMAN.

[TECH STACK]
- Framework: Vue 3 (Composition API / <script setup>)
- UI Library: PrimeVue (Gunakan skema komponen modern v4+, Pass-Through, atau Design Tokens)
- Utility Styling: Tailwind CSS
- Icons: PrimeIcons / Lucide Icons

[DESIGN & UX PRINCIPLES - HARUS DITURUTI]
1. Minimum Page Navigation & Modal-First Architecture:
   - UTAMAKAN penggunaan Dialog/Modal, Drawer (Slide-over), dan Popover alih-alih berpindah halaman/route baru.
   - Semua aksi CRUD (Create, Read, Update, Delete) dan Approval HARUS dilakukan di atas halaman utama menggunakan Modal/Drawer.
   - Pengguna harus tetap berada di tabel/dashboard utama saat melakukan operasi data.

2. Compact & High-Density First:
   - Hindari ruang kosong (whitespace) yang tidak perlu.
   - Gunakan padding dan margin yang ketat (misal: py-1.5, py-2, px-3, p-3).
   - Ukuran teks dominan adalah text-xs (12px) sampai text-sm (14px).
   - Gunakan border-radius yang kecil/presisi (rounded-md atau rounded-sm).

3. Zero / Minimal Scroll Architecture:
   - Manfaatkan layout 'h-screen' atau 'h-[calc(100vh-...)]' dengan 'overflow-hidden' pada container utama.
   - Gunakan pola Split-Screen / Master-Detail (Panel Kiri: DataTable, Panel Kanan: Detail View/Tabs).
   - Gunakan komponen 'Tabs' horizontal di dalam Modal/Drawer untuk memecah form berukuran besar.

4. Standar Komponen PrimeVue & Tailwind:
   - DataTable: Gunakan opsi compact/small (`p-datatable-sm`), scrollable dengan `scrollHeight="flex"`, virtual scroll jika perlu, dan pinned/frozen columns untuk aksi.
   - Dialog: Selalu atur ukuran lebar eksplisit (misal: `:style="{ width: '450px' }"`).
   - Filters & Actions: Satukan Search Bar, Filter Select, dan Quick Action Buttons dalam satu baris horizontal menggunakan Flexbox ringkas (pola filter chips `All | Active | ...`).
   - Warna Status: UTAMAKAN `<Tag :severity="...">` PrimeVue (success/danger/warn/info/secondary) — ini pola dominan di codebase; Soft Badge Tailwind (`bg-emerald-50 text-emerald-700`) hanya untuk elemen non-status/kustom.

5. Dark Mode:
   - Setiap komponen WAJIB punya variant `dark:` (bg, text, border) — jangan pernah menulis warna terang saja.
   - Root memakai `dark:bg-gray-900`/`dark:bg-gray-800`; selector `.p-dark` menangani komponen PrimeVue.
   - Test visual dark mode di setiap halaman baru sebelum dianggap selesai.

6. Bilingual EN/ID:
   - JANGAN hardcode teks label/status/pesan — gunakan `t('...')` + key di `en.json`/`id.json`.
   - Setiap key baru WAJIB ada di kedua locale (en & id) dengan isi yang benar.
   - Format tanggal/angka mengikuti bahasa aktif (`id-ID` vs `en-US`).

7. Empty State:
   - Tabel/panel kosong WAJIB menampilkan empty state: ikon `pi-inbox` + teks singkat (lihat `#empty` template DataTable).

[OUTPUT FORMAT REQUIREMENT]
Setiap kali memberikan kodingan atau solusi UI/UX:
1. Tuliskan kode Vue 3 lengkap menggunakan Single File Component (SFC) <template> dan <script setup>.
2. Pastikan gabungan kelas Tailwind CSS dan atribut PrimeVue mematuhi aturan high-density & modal-first di atas.
3. Berikan penjelasan singkat mengenai struktur layout jika ada pola UX khusus yang diterapkan.
```

---

## 4. Standar Warna Badge Status

**Pola dominan di codebase: `<Tag :severity="...">` PrimeVue** — mapping severity via fungsi helper per halaman (mis. `statusSeverity(status)`):

| Status | `severity` PrimeVue | Setara Soft Badge Tailwind |
| :--- | :--- | :--- |
| **Approved / Active / Success** | `success` | `bg-emerald-50 text-emerald-700 border border-emerald-200` |
| **Pending / PENDING_APPROVAL / Review** | `warn` | `bg-amber-50 text-amber-700 border border-amber-200` |
| **Submitted / Info / Mutation** | `info` | `bg-indigo-50 text-indigo-700 border border-indigo-200` |
| **Rejected / Terminated / Offboarding** | `danger` | `bg-rose-50 text-rose-700 border border-rose-200` |
| **Draft / Cancelled / Retired** | `secondary` | `bg-gray-50 text-gray-700 border border-gray-200` |
| **Published / Paid** | `info` atau `success` | `bg-sky-50 text-sky-700 border border-sky-200` |
| **Expired** | `danger` | `bg-red-50 text-red-700 border border-red-200` |

Contoh pola aktual (Leave): `APPROVED_FINAL → 'success'`, `REJECTED_FINAL → 'danger'`, `PENDING_APPROVAL → 'warn'`, `SUBMITTED → 'info'`, `DRAFT/CANCELLED → 'secondary'`.  
Gunakan ukuran ringkas: `class="!text-xs !px-1.5 !py-0.5"`. Soft Badge Tailwind hanya untuk elemen non-status/kustom (mis. highlight step aktif).

---

> 📄 **Roadmap implementasi frontend** (Phase 1–4, checklist per halaman, bilingual support) telah dipindah ke [`docs/frontend-development-plan.md`](frontend-development-plan.md) — dokumen ini hanya berisi standar UI/UX.

*Dokumen ini akan diupdate seiring progres implementasi frontend.*
