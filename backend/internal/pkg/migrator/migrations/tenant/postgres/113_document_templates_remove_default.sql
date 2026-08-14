-- Hapus seluruh fitur template default (referensi): seed rows dihapus dan
-- kolom is_default tidak lagi dipakai (keputusan: template dibuat dari nol,
-- tanpa default template).
DELETE FROM document_templates WHERE is_default = TRUE;

DROP INDEX IF EXISTS uq_document_templates_default_per_type;

ALTER TABLE document_templates
    DROP COLUMN is_default;
