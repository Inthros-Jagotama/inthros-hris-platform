# HRIS Platform — Docker Environment

Jalankan dari direktori `docker/`.

## Menjalankan Stack

```bash
# PostgreSQL + PgBouncer + Redis + API
docker compose --profile postgres up -d

# ATAU MySQL
docker compose --profile mysql up -d

# Hentikan
docker compose --profile postgres --profile mysql down

# Reset database (menghapus volume agar terinisialisasi ulang dari awal)
docker compose --profile mysql down -v
docker compose --profile mysql up -d --force-recreate

docker compose --profile postgres down -v
docker compose --profile postgres up -d --force-recreate
```

## Template Dokumen (LibreOffice DOCX → PDF)

Image `api` sudah meng-bundle **LibreOffice** (`libreoffice-writer` pada runtime
Debian) sehingga fitur Settings → Template Dokumen (preview & generate document)
langsung berfungsi. File upload (template `.docx`, `previews/`,
`generated_documents/`) disimpan di volume `uploads_data` agar persist saat
container di-recreate.

Verifikasi / konversi manual via helper container:

```bash
# Cek versi LibreOffice
docker compose --profile tools run --rm libreoffice --version

# Konversi manual: DOCX di dalam volume uploads → PDF (hasil di volume yang sama)
docker compose --profile tools run --rm libreoffice --headless \
  --convert-to pdf /uploads/contoh.docx --outdir /uploads
```

> Catatan: `soffice` dipanggil langsung oleh proses backend (`os/exec`), jadi
> LibreOffice harus berada di image yang sama dengan server (bukan container
> terpisah). Helper `libreoffice` di atas hanya untuk testing/konversi manual —
> service `api` melakukan konversinya sendiri.
