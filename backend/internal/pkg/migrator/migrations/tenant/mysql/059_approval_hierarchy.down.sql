-- 059_approval_hierarchy.down.sql

DROP TABLE IF EXISTS approval_flow_step_organizations;

SET @drop_participation_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'approval_flow_steps'
      AND COLUMN_NAME = 'participation_type'
  ),
  'ALTER TABLE approval_flow_steps DROP COLUMN participation_type',
  'DO 0'
);
PREPARE stmt FROM @drop_participation_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_hierarchy_level = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'approval_flow_steps'
      AND COLUMN_NAME = 'hierarchy_level'
  ),
  'ALTER TABLE approval_flow_steps DROP COLUMN hierarchy_level',
  'DO 0'
);
PREPARE stmt FROM @drop_hierarchy_level;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

ALTER TABLE approval_flow_steps
    MODIFY COLUMN approver_type ENUM('SUPERVISOR', 'ROLE', 'USER') NOT NULL;
