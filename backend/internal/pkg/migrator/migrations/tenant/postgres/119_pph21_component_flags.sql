-- =============================================================================
-- Tenant Migration: 119_pph21_component_flags
-- =============================================================================
-- Sumber kebenaran komponen PPh21 pindah dari setting ke salary_components:
--   - is_pph21_component  : komponen "wadah" baris potongan hasil pajak (1 komponen).
--   - is_pph21_deductible : komponen pengurang penghasilan bruto (mis. iuran
--                           pensiun — bisa lebih dari satu).
-- Setting PPh21 tidak lagi memilih komponen (pph21_component_id dihapus),
-- dan log kalkulasi mencatat pengurang non-BPJS (pension_deductible_monthly).

-- ---------------------------------------------------------------------------
-- 119.1 salary_components: flag komponen PPh21
-- ---------------------------------------------------------------------------
ALTER TABLE salary_components
    ADD COLUMN is_pph21_component SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN is_pph21_deductible SMALLINT NOT NULL DEFAULT 0;

-- ---------------------------------------------------------------------------
-- 119.2 pph21_settings: hapus pph21_component_id + FK-nya
-- ---------------------------------------------------------------------------
ALTER TABLE pph21_settings
    DROP CONSTRAINT IF EXISTS fk_pph21_setting_component;

ALTER TABLE pph21_settings
    DROP COLUMN IF EXISTS pph21_component_id;

-- ---------------------------------------------------------------------------
-- 119.3 pph21_calculation_logs: pengurang non-BPJS (iuran pensiun dll.)
-- ---------------------------------------------------------------------------
ALTER TABLE pph21_calculation_logs
    ADD COLUMN pension_deductible_monthly DECIMAL(18, 2) NOT NULL DEFAULT 0;
