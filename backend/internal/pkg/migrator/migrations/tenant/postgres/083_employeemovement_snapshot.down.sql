-- =============================================================================
-- Tenant Migration Down: 083_employeemovement_snapshot
-- =============================================================================

ALTER TABLE employee_movements
    DROP COLUMN IF EXISTS from_organization_name,
    DROP COLUMN IF EXISTS from_position_name,
    DROP COLUMN IF EXISTS from_employment_status_name,
    DROP COLUMN IF EXISTS to_organization_name,
    DROP COLUMN IF EXISTS to_position_name,
    DROP COLUMN IF EXISTS to_employment_status_name;
