-- 058_performance_scoring_configuration.down.sql
-- Rollback Phase 5 KPI Enhancement: Drop scoring configuration tables

SET @drop_fk_evalcomp_component = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluation_components'
      AND CONSTRAINT_NAME = 'fk_perf_evalcomp_component'
  ),
  'ALTER TABLE performance_evaluation_components DROP FOREIGN KEY fk_perf_evalcomp_component',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_evalcomp_component;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_fk_evalcomp_eval = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluation_components'
      AND CONSTRAINT_NAME = 'fk_perf_evalcomp_eval'
  ),
  'ALTER TABLE performance_evaluation_components DROP FOREIGN KEY fk_perf_evalcomp_eval',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_evalcomp_eval;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_fk_orgcomp_component = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_organization_components'
      AND CONSTRAINT_NAME = 'fk_perf_orgcomp_component'
  ),
  'ALTER TABLE performance_organization_components DROP FOREIGN KEY fk_perf_orgcomp_component',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_orgcomp_component;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

DROP TABLE IF EXISTS performance_evaluation_components;
DROP TABLE IF EXISTS performance_organization_components;
DROP TABLE IF EXISTS performance_components;