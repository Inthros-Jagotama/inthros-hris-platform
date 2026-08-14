ALTER TABLE salary_components
    DROP CONSTRAINT IF EXISTS fk_salary_comp_reference;

DROP INDEX IF EXISTS idx_salary_comp_reference;

ALTER TABLE salary_components
    DROP COLUMN formula,
    DROP COLUMN reference_component_id;
