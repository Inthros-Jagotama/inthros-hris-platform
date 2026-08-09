-- =============================================================================
-- Tenant Migration: 076_notification_params
-- =============================================================================
-- Notification: tambah kolom `params` (JSON-encoded []string) untuk merender
-- ulang title/body notifikasi secara bilingual sesuai bahasa penerima saat
-- GET /notifications, bukan bahasa pengirim aksi yang memicu notifikasi.
-- Tanpa kolom ini, INSERT ke tabel notifications gagal (kolom tidak dikenal
-- GORM model) sehingga notifikasi baru berhenti tersimpan sama sekali.

SET @add_params = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'notifications'
      AND column_name = 'params'
  ),
  'DO 0',
  'ALTER TABLE notifications ADD COLUMN params TEXT NULL COMMENT ''JSON-encoded []string, body template placeholders'''
);
PREPARE stmt FROM @add_params;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
