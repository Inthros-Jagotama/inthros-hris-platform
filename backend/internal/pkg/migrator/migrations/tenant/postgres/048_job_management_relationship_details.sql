-- Migration: 048_job_management_relationship_details
-- Tabel detail Hubungan Kerja: rincian aktivitas per job_management_relationships.
-- Mengikuti standar tabel detail yang sudah ada (009_job_management, 047).

CREATE TABLE IF NOT EXISTS job_management_relationship_details (
    id                             CHAR(36) PRIMARY KEY,
    job_management_relationship_id CHAR(36) NOT NULL,
    organization_id                CHAR(36) NULL,
    activity                       TEXT NULL,
    created_by                     CHAR(36) NULL,
    updated_by                     CHAR(36) NULL,
    created_at                     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_jmrd_relationship ON job_management_relationship_details (job_management_relationship_id);
CREATE INDEX IF NOT EXISTS idx_jmrd_organization ON job_management_relationship_details (organization_id);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmrd_relationship') THEN
        ALTER TABLE job_management_relationship_details
            ADD CONSTRAINT fk_jmrd_relationship
            FOREIGN KEY (job_management_relationship_id)
            REFERENCES job_management_relationships(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmrd_organization') THEN
        ALTER TABLE job_management_relationship_details
            ADD CONSTRAINT fk_jmrd_organization
            FOREIGN KEY (organization_id)
            REFERENCES organizations(id) ON DELETE SET NULL;
    END IF;
END $$;
