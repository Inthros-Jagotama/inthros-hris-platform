-- Down: kembalikan pph21_ptkp_rates (schema migration 006) dan hapus kolom code.
CREATE TABLE IF NOT EXISTS pph21_ptkp_rates (
    id                   CHAR(36) PRIMARY KEY,
    ptkp_status          VARCHAR(20) NOT NULL,
    description          VARCHAR(255) NULL,
    annual_amount        DECIMAL(18, 2) NOT NULL DEFAULT 0,
    effective_start_date DATE NOT NULL,
    effective_end_date   DATE NULL,
    status               VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_by           CHAR(36) NULL,
    updated_by           CHAR(36) NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_pph21_ptkp_status_start UNIQUE (ptkp_status, effective_start_date)
);

CREATE INDEX IF NOT EXISTS idx_pph21_ptkp_effective ON pph21_ptkp_rates (effective_start_date, effective_end_date, status);

ALTER TABLE ptkps DROP COLUMN code;
