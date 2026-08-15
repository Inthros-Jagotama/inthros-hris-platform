-- Rollback 119: kembalikan pph21_component_id + FK, hapus flag & kolom log.

ALTER TABLE pph21_calculation_logs
    DROP COLUMN IF EXISTS pension_deductible_monthly;

ALTER TABLE pph21_settings
    ADD COLUMN pph21_component_id CHAR(36) NULL;

ALTER TABLE pph21_settings
    ADD CONSTRAINT fk_pph21_setting_component FOREIGN KEY (pph21_component_id)
        REFERENCES salary_components(id) ON DELETE CASCADE;

ALTER TABLE salary_components
    DROP COLUMN IF EXISTS is_pph21_deductible,
    DROP COLUMN IF EXISTS is_pph21_component;
