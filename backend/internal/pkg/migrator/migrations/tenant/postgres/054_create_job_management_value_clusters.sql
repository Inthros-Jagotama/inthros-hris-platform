-- Migration: 054_create_job_management_value_clusters
-- Mapping tipe nilai jabatan (type) dengan cluster kompetensi (dari tabel competencies).
-- Dipakai card Kompetensi Teknis di form job management sebagai filter cluster yang valid
-- untuk tipe 'technical' — diatur di halaman Mapping Job Value → sub menu technical.

CREATE TABLE job_management_value_clusters (
    id UUID NOT NULL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    cluster VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_job_management_value_clusters_type_cluster UNIQUE (type, cluster)
);

CREATE INDEX idx_job_management_value_clusters_type ON job_management_value_clusters (type);
