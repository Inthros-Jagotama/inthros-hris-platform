-- Template SK Movement kini bisa dibuat per jenis movement (promotion, mutation,
-- dll). movement_type kosong ('') = template umum (fallback semua jenis movement).
ALTER TABLE document_templates
    ADD COLUMN movement_type VARCHAR(50) NOT NULL DEFAULT '';

-- Unique ACTIVE per (type, movement_type): satu template aktif per jenis
-- movement (dan satu template umum per document type non-movement).
DROP INDEX IF EXISTS uq_document_templates_active_per_type;
CREATE UNIQUE INDEX uq_document_templates_active_per_type
    ON document_templates (type, movement_type)
    WHERE status = 'ACTIVE';

CREATE INDEX idx_document_templates_type_movement_status
    ON document_templates (type, movement_type, status);
