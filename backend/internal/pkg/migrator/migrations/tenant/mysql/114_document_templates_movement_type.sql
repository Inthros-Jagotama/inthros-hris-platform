-- Template SK Movement kini bisa dibuat per jenis movement (promotion, mutation,
-- dll). movement_type kosong ('') = template umum (fallback semua jenis movement).
ALTER TABLE document_templates
    ADD COLUMN movement_type VARCHAR(50) NOT NULL DEFAULT '' AFTER type;

-- Index lookup aktif per (type, movement_type, status) — dipakai generator saat
-- memilih template (type + movement_type spesifik, lalu fallback movement_type='').
CREATE INDEX idx_document_templates_type_movement_status
    ON document_templates (type, movement_type, status);
