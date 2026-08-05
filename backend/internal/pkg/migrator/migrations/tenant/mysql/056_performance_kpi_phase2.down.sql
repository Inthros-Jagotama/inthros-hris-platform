-- 056_performance_kpi_phase2.down.sql
-- Rollback Phase 2 KPI Enhancement: Drop 6 new tables

-- Drop FK constraint from performance_evaluations to performance_ratings first
SET @drop_fk_eval_rating = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluations'
      AND CONSTRAINT_NAME = 'fk_perf_eval_rating'
  ),
  'ALTER TABLE performance_evaluations DROP FOREIGN KEY fk_perf_eval_rating',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_eval_rating;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS performance_logs;
DROP TABLE IF EXISTS performance_indicator_formulas;
DROP TABLE IF EXISTS performance_ratings;
DROP TABLE IF EXISTS performance_attachments;
DROP TABLE IF EXISTS performance_comments;
DROP TABLE IF EXISTS performance_progress;
