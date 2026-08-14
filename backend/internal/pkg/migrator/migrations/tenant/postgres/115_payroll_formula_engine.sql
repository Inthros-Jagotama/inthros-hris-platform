-- =============================================================================
-- Tenant Migration: 115_payroll_formula_engine
-- =============================================================================
-- Prasyarat struktural Formula Engine (docs/payroll/02-formula-engine.md):
-- kolom formula (untuk calculation_type FORMULA/PERCENTAGE) dan
-- reference_component_id (untuk calculation_type REFERENCE) di salary_components.
-- Sebelumnya FORMULA/REFERENCE tidak punya tempat menyimpan datanya sama sekali.

ALTER TABLE salary_components
    ADD COLUMN formula TEXT NULL,
    ADD COLUMN reference_component_id CHAR(36) NULL;

-- Referensi komponen lain (self-FK) — dihapus set null agar master component
-- tetap bisa dihapus walau dirujuk.
CREATE INDEX IF NOT EXISTS idx_salary_comp_reference ON salary_components (reference_component_id);

ALTER TABLE salary_components
    ADD CONSTRAINT fk_salary_comp_reference FOREIGN KEY (reference_component_id)
        REFERENCES salary_components(id) ON DELETE SET NULL;
