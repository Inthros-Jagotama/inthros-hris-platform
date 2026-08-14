-- =============================================================================
-- Tenant Migration: 115_payroll_formula_engine
-- =============================================================================
-- Prasyarat struktural Formula Engine (docs/payroll/02-formula-engine.md):
-- kolom formula (untuk calculation_type FORMULA/PERCENTAGE) dan
-- reference_component_id (untuk calculation_type REFERENCE) di salary_components.

ALTER TABLE salary_components
    ADD COLUMN formula TEXT NULL,
    ADD COLUMN reference_component_id CHAR(36) NULL AFTER calculation_type;

CREATE INDEX idx_salary_comp_reference ON salary_components (reference_component_id);

ALTER TABLE salary_components
    ADD CONSTRAINT fk_salary_comp_reference FOREIGN KEY (reference_component_id)
        REFERENCES salary_components(id) ON DELETE SET NULL;
