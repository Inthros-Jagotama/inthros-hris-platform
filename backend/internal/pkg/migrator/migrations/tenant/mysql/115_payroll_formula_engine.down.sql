ALTER TABLE salary_components
    DROP FOREIGN KEY fk_salary_comp_reference;

DROP INDEX idx_salary_comp_reference ON salary_components;

ALTER TABLE salary_components
    DROP COLUMN reference_component_id,
    DROP COLUMN formula;
