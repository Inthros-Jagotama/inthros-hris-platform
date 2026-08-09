-- =============================================================================
-- Tenant Migration: 076_notification_params
-- =============================================================================
-- Notification: tambah kolom `params` (JSON-encoded []string) untuk merender
-- ulang title/body notifikasi secara bilingual sesuai bahasa penerima saat
-- GET /notifications, bukan bahasa pengirim aksi yang memicu notifikasi.
-- Tanpa kolom ini, INSERT ke tabel notifications gagal (kolom tidak dikenal
-- GORM model) sehingga notifikasi baru berhenti tersimpan sama sekali.

ALTER TABLE notifications ADD COLUMN IF NOT EXISTS params TEXT NULL;
