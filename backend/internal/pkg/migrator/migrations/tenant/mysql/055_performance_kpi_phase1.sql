-- 055_performance_kpi_phase1.sql
-- Phase 1 KPI Enhancement: Add new columns to performance tables
-- Supports: Period linking, Formula types, Rating integration, KPI snapshots

-- =========================================================================
-- 1. Update performance_templates
-- Add: period_id, effective_date, expired_date
-- =========================================================================

-- Add period_id column
SET @add_period_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'period_id'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD COLUMN period_id CHAR(36) NULL AFTER organization_id'
);
PREPARE stmt FROM @add_period_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add effective_date column
SET @add_effective_date = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'effective_date'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD COLUMN effective_date DATE NULL AFTER status'
);
PREPARE stmt FROM @add_effective_date;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add expired_date column
SET @add_expired_date = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'expired_date'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD COLUMN expired_date DATE NULL AFTER effective_date'
);
PREPARE stmt FROM @add_expired_date;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index for period_id
SET @idx_tmpl_period = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND index_name = 'idx_perf_tmpl_period'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD INDEX idx_perf_tmpl_period (period_id)'
);
PREPARE stmt FROM @idx_tmpl_period;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index for effective_date
SET @idx_tmpl_effective = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND index_name = 'idx_perf_tmpl_effective'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD INDEX idx_perf_tmpl_effective (effective_date)'
);
PREPARE stmt FROM @idx_tmpl_effective;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add FK constraint for period_id
SET @fk_tmpl_period = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_templates'
      AND CONSTRAINT_NAME = 'fk_perf_tmpl_period'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD CONSTRAINT fk_perf_tmpl_period FOREIGN KEY (period_id) REFERENCES performance_periods(id) ON DELETE SET NULL'
);
PREPARE stmt FROM @fk_tmpl_period;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 2. Update performance_indicators
-- Add: code, formula_type, minimum_score, maximum_score, target_type, is_required
-- =========================================================================

-- Add code column
SET @add_code = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'code'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD COLUMN code VARCHAR(50) NULL AFTER perspective_id'
);
PREPARE stmt FROM @add_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add formula_type column
SET @add_formula_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'formula_type'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD COLUMN formula_type VARCHAR(30) NOT NULL DEFAULT ''MANUAL'' AFTER unit_of_measurement'
);
PREPARE stmt FROM @add_formula_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add minimum_score column
SET @add_minimum_score = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'minimum_score'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD COLUMN minimum_score DECIMAL(5,2) NOT NULL DEFAULT 0.00 AFTER formula_type'
);
PREPARE stmt FROM @add_minimum_score;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add maximum_score column
SET @add_maximum_score = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'maximum_score'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD COLUMN maximum_score DECIMAL(5,2) NOT NULL DEFAULT 100.00 AFTER minimum_score'
);
PREPARE stmt FROM @add_maximum_score;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add target_type column
SET @add_target_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'target_type'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD COLUMN target_type VARCHAR(30) NOT NULL DEFAULT ''NUMBER'' AFTER maximum_score'
);
PREPARE stmt FROM @add_target_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add is_required column
SET @add_is_required = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND column_name = 'is_required'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD COLUMN is_required TINYINT(1) NOT NULL DEFAULT 1 AFTER target_type'
);
PREPARE stmt FROM @add_is_required;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index for code
SET @idx_ind_code = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_indicators'
      AND index_name = 'idx_perf_ind_code'
  ),
  'DO 0',
  'ALTER TABLE performance_indicators ADD INDEX idx_perf_ind_code (code)'
);
PREPARE stmt FROM @idx_ind_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 3. Update performance_evaluations
-- Add: rating_id, submitted_at, approved_at
-- =========================================================================

-- Add rating_id column
SET @add_rating_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'rating_id'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN rating_id CHAR(36) NULL AFTER final_score'
);
PREPARE stmt FROM @add_rating_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add submitted_at column
SET @add_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'submitted_at'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN submitted_at TIMESTAMP(6) NULL AFTER actual_approved_at'
);
PREPARE stmt FROM @add_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add approved_at column
SET @add_approved_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'approved_at'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN approved_at TIMESTAMP(6) NULL AFTER submitted_at'
);
PREPARE stmt FROM @add_approved_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index for rating_id
SET @idx_eval_rating = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND index_name = 'idx_perf_eval_rating'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD INDEX idx_perf_eval_rating (rating_id)'
);
PREPARE stmt FROM @idx_eval_rating;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 4. Update performance_evaluation_details
-- Add: indicator_id, indicator_name (snapshot), target (snapshot), actual, achievement, remarks
-- =========================================================================

-- Add indicator_id column
SET @add_indicator_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'indicator_id'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN indicator_id CHAR(36) NULL AFTER perspective_id'
);
PREPARE stmt FROM @add_indicator_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add indicator_name column
SET @add_indicator_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'indicator_name'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN indicator_name VARCHAR(255) NULL AFTER indicator_id'
);
PREPARE stmt FROM @add_indicator_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add target column
SET @add_target = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'target'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN target DECIMAL(18,2) NOT NULL DEFAULT 0.00 AFTER indicator_name'
);
PREPARE stmt FROM @add_target;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add actual column
SET @add_actual = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'actual'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN actual DECIMAL(18,2) NOT NULL DEFAULT 0.00 AFTER target'
);
PREPARE stmt FROM @add_actual;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add achievement column
SET @add_achievement = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'achievement'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN achievement DECIMAL(5,2) NOT NULL DEFAULT 0.00 AFTER actual'
);
PREPARE stmt FROM @add_achievement;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add remarks column
SET @add_remarks = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'remarks'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN remarks TEXT NULL AFTER description'
);
PREPARE stmt FROM @add_remarks;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index for indicator_id
SET @idx_detail_indicator = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND index_name = 'idx_perf_detail_indicator'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD INDEX idx_perf_detail_indicator (indicator_id)'
);
PREPARE stmt FROM @idx_detail_indicator;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add FK constraint for indicator_id
SET @fk_detail_indicator = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluation_details'
      AND CONSTRAINT_NAME = 'fk_perf_detail_indicator'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD CONSTRAINT fk_perf_detail_indicator FOREIGN KEY (indicator_id) REFERENCES performance_indicators(id) ON DELETE SET NULL'
);
PREPARE stmt FROM @fk_detail_indicator;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
