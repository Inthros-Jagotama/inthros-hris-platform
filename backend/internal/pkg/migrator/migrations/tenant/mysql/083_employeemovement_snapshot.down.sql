-- =============================================================================
-- Tenant Migration Down: 083_employeemovement_snapshot
-- =============================================================================

SET @drop_from_org_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'from_organization_name'
  ),
  'ALTER TABLE employee_movements DROP COLUMN from_organization_name',
  'DO 0'
);
PREPARE stmt FROM @drop_from_org_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_from_pos_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'from_position_name'
  ),
  'ALTER TABLE employee_movements DROP COLUMN from_position_name',
  'DO 0'
);
PREPARE stmt FROM @drop_from_pos_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_from_status_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'from_employment_status_name'
  ),
  'ALTER TABLE employee_movements DROP COLUMN from_employment_status_name',
  'DO 0'
);
PREPARE stmt FROM @drop_from_status_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_to_org_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'to_organization_name'
  ),
  'ALTER TABLE employee_movements DROP COLUMN to_organization_name',
  'DO 0'
);
PREPARE stmt FROM @drop_to_org_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_to_pos_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'to_position_name'
  ),
  'ALTER TABLE employee_movements DROP COLUMN to_position_name',
  'DO 0'
);
PREPARE stmt FROM @drop_to_pos_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_to_status_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'to_employment_status_name'
  ),
  'ALTER TABLE employee_movements DROP COLUMN to_employment_status_name',
  'DO 0'
);
PREPARE stmt FROM @drop_to_status_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
