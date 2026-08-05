-- 055_performance_kpi_phase1.down.sql
-- Rollback Phase 1 KPI Enhancement

-- =========================================================================
-- 1. Drop FK constraints first (idempotent)
-- =========================================================================

SET @drop_fk_tmpl_period = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_templates'
      AND CONSTRAINT_NAME = 'fk_perf_tmpl_period'
  ),
  'ALTER TABLE performance_templates DROP FOREIGN KEY fk_perf_tmpl_period',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_tmpl_period;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_fk_detail_indicator = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluation_details'
      AND CONSTRAINT_NAME = 'fk_perf_detail_indicator'
  ),
  'ALTER TABLE performance_evaluation_details DROP FOREIGN KEY fk_perf_detail_indicator',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_detail_indicator;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 2. Drop indexes (idempotent)
-- =========================================================================

SET @drop_idx_tmpl_period = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND index_name = 'idx_perf_tmpl_period'
  ),
  'ALTER TABLE performance_templates DROP INDEX idx_perf_tmpl_period',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_tmpl_period;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_idx_tmpl_effective = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND index_name = 'idx_perf_tmpl_effective'
  ),
  'ALTER TABLE performance_templates DROP INDEX idx_perf_tmpl_effective',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_tmpl_effective;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_idx_ind_code = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND index_name = 'idx_perf_ind_code'
  ),
  'ALTER TABLE performance_indicators DROP INDEX idx_perf_ind_code',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_ind_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_idx_eval_rating = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND index_name = 'idx_perf_eval_rating'
  ),
  'ALTER TABLE performance_evaluations DROP INDEX idx_perf_eval_rating',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_eval_rating;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_idx_detail_indicator = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND index_name = 'idx_perf_detail_indicator'
  ),
  'ALTER TABLE performance_evaluation_details DROP INDEX idx_perf_detail_indicator',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_detail_indicator;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 3. Drop columns from performance_templates
-- =========================================================================

SET @drop_period_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'period_id'
  ),
  'ALTER TABLE performance_templates DROP COLUMN period_id',
  'DO 0'
);
PREPARE stmt FROM @drop_period_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_effective_date = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'effective_date'
  ),
  'ALTER TABLE performance_templates DROP COLUMN effective_date',
  'DO 0'
);
PREPARE stmt FROM @drop_effective_date;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_expired_date = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'expired_date'
  ),
  'ALTER TABLE performance_templates DROP COLUMN expired_date',
  'DO 0'
);
PREPARE stmt FROM @drop_expired_date;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 4. Drop columns from performance_indicators
-- =========================================================================

SET @drop_code = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'code'
  ),
  'ALTER TABLE performance_indicators DROP COLUMN code',
  'DO 0'
);
PREPARE stmt FROM @drop_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_formula_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'formula_type'
  ),
  'ALTER TABLE performance_indicators DROP COLUMN formula_type',
  'DO 0'
);
PREPARE stmt FROM @drop_formula_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_minimum_score = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'minimum_score'
  ),
  'ALTER TABLE performance_indicators DROP COLUMN minimum_score',
  'DO 0'
);
PREPARE stmt FROM @drop_minimum_score;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_maximum_score = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'maximum_score'
  ),
  'ALTER TABLE performance_indicators DROP COLUMN maximum_score',
  'DO 0'
);
PREPARE stmt FROM @drop_maximum_score;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_target_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'target_type'
  ),
  'ALTER TABLE performance_indicators DROP COLUMN target_type',
  'DO 0'
);
PREPARE stmt FROM @drop_target_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_is_required = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'is_required'
  ),
  'ALTER TABLE performance_indicators DROP COLUMN is_required',
  'DO 0'
);
PREPARE stmt FROM @drop_is_required;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 5. Drop columns from performance_evaluations
-- =========================================================================

SET @drop_rating_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'rating_id'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN rating_id',
  'DO 0'
);
PREPARE stmt FROM @drop_rating_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'submitted_at'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN submitted_at',
  'DO 0'
);
PREPARE stmt FROM @drop_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_approved_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'approved_at'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN approved_at',
  'DO 0'
);
PREPARE stmt FROM @drop_approved_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 6. Drop columns from performance_evaluation_details
-- =========================================================================

SET @drop_indicator_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'indicator_id'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN indicator_id',
  'DO 0'
);
PREPARE stmt FROM @drop_indicator_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_indicator_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'indicator_name'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN indicator_name',
  'DO 0'
);
PREPARE stmt FROM @drop_indicator_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_target = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'target'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN target',
  'DO 0'
);
PREPARE stmt FROM @drop_target;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'actual'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN actual',
  'DO 0'
);
PREPARE stmt FROM @drop_actual;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_achievement = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'achievement'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN achievement',
  'DO 0'
);
PREPARE stmt FROM @drop_achievement;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_remarks = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'remarks'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN remarks',
  'DO 0'
);
PREPARE stmt FROM @drop_remarks;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
