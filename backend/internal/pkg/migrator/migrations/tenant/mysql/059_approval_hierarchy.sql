-- 059_approval_hierarchy.sql
-- Approval Module Phase 1: Organization-hierarchy approver resolution +
-- step participation type (APPROVER vs WATCHER).

ALTER TABLE approval_flow_steps
    MODIFY COLUMN approver_type ENUM('SUPERVISOR', 'ROLE', 'USER', 'ORGANIZATION', 'BOTH') NOT NULL;

SET @add_hierarchy_level = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'approval_flow_steps'
      AND COLUMN_NAME = 'hierarchy_level'
  ),
  'DO 0',
  'ALTER TABLE approval_flow_steps ADD COLUMN hierarchy_level INT NULL DEFAULT 1 AFTER approver_type'
);
PREPARE stmt FROM @add_hierarchy_level;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_participation_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'approval_flow_steps'
      AND COLUMN_NAME = 'participation_type'
  ),
  'DO 0',
  'ALTER TABLE approval_flow_steps ADD COLUMN participation_type ENUM(''APPROVER'',''WATCHER'') NOT NULL DEFAULT ''APPROVER'' AFTER approval_mode'
);
PREPARE stmt FROM @add_participation_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- approval_flow_step_organizations — target Organization(s) untuk step dengan
-- approver_type ORGANIZATION/BOTH. Satu step bisa punya banyak Organization
-- (mendukung "lebih dari satu approver di satu level").
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_flow_step_organizations (
    id              CHAR(36) PRIMARY KEY,
    step_id         CHAR(36) NOT NULL,
    organization_id CHAR(36) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_approval_step_org_step (step_id),
    INDEX idx_approval_step_org_org (organization_id),

    CONSTRAINT fk_approval_step_org_step FOREIGN KEY (step_id) REFERENCES approval_flow_steps(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
