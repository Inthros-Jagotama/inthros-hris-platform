-- 056_performance_kpi_phase2.sql
-- Phase 2 KPI Enhancement: Create 6 new tables
-- Tables: performance_progress, performance_comments, performance_attachments,
--         performance_ratings, performance_indicator_formulas, performance_logs

-- =========================================================================
-- 1. performance_progress - Progress monitoring untuk KPI evaluation details
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_progress (
    id CHAR(36) PRIMARY KEY,
    evaluation_detail_id CHAR(36) NOT NULL,
    progress_date DATE NOT NULL,
    actual_value DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    achievement DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    notes TEXT NULL,
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_perf_prog_detail (evaluation_detail_id),
    INDEX idx_perf_prog_date (progress_date),
    INDEX idx_perf_prog_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint for evaluation_detail_id
SET @fk_prog_detail = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_progress'
      AND CONSTRAINT_NAME = 'fk_perf_prog_detail'
  ),
  'DO 0',
  'ALTER TABLE performance_progress ADD CONSTRAINT fk_perf_prog_detail FOREIGN KEY (evaluation_detail_id) REFERENCES performance_evaluation_details(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_prog_detail;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 2. performance_comments - Komentar antara Employee dan Reviewer
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_comments (
    id CHAR(36) PRIMARY KEY,
    evaluation_id CHAR(36) NOT NULL,
    employee_id CHAR(36) NOT NULL,
    comment TEXT NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_perf_comment_eval (evaluation_id),
    INDEX idx_perf_comment_emp (employee_id),
    INDEX idx_perf_comment_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint for evaluation_id
SET @fk_comment_eval = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_comments'
      AND CONSTRAINT_NAME = 'fk_perf_comment_eval'
  ),
  'DO 0',
  'ALTER TABLE performance_comments ADD CONSTRAINT fk_perf_comment_eval FOREIGN KEY (evaluation_id) REFERENCES performance_evaluations(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_comment_eval;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 3. performance_attachments - Evidence/lampiran KPI
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_attachments (
    id CHAR(36) PRIMARY KEY,
    evaluation_detail_id CHAR(36) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(100) NULL,
    file_size BIGINT NULL,
    description TEXT NULL,
    uploaded_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_perf_attach_detail (evaluation_detail_id),
    INDEX idx_perf_attach_uploaded_by (uploaded_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint for evaluation_detail_id
SET @fk_attach_detail = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_attachments'
      AND CONSTRAINT_NAME = 'fk_perf_attach_detail'
  ),
  'DO 0',
  'ALTER TABLE performance_attachments ADD CONSTRAINT fk_perf_attach_detail FOREIGN KEY (evaluation_detail_id) REFERENCES performance_evaluation_details(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_attach_detail;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 4. performance_ratings - Master data rating performance
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_ratings (
    id CHAR(36) PRIMARY KEY,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    min_score DECIMAL(5,2) NOT NULL,
    max_score DECIMAL(5,2) NOT NULL,
    color VARCHAR(20) NULL,
    description TEXT NULL,
    sort_order SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_perf_rating_code (code),
    INDEX idx_perf_rating_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint from performance_evaluations to performance_ratings
SET @fk_eval_rating = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluations'
      AND CONSTRAINT_NAME = 'fk_perf_eval_rating'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD CONSTRAINT fk_perf_eval_rating FOREIGN KEY (rating_id) REFERENCES performance_ratings(id) ON DELETE SET NULL'
);
PREPARE stmt FROM @fk_eval_rating;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =========================================================================
-- 5. performance_indicator_formulas - Master formula kalkulasi KPI
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_indicator_formulas (
    id CHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    formula_type VARCHAR(30) NOT NULL,
    expression TEXT NULL,
    description TEXT NULL,
    sort_order SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_perf_formula_code (code),
    INDEX idx_perf_formula_type (formula_type),
    INDEX idx_perf_formula_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- 6. performance_logs - Audit trail aktivitas performance
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_logs (
    id CHAR(36) PRIMARY KEY,
    evaluation_id CHAR(36) NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id CHAR(36) NOT NULL,
    action VARCHAR(100) NOT NULL,
    old_values TEXT NULL,
    new_values TEXT NULL,
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_perf_log_eval (evaluation_id),
    INDEX idx_perf_log_entity (entity_type, entity_id),
    INDEX idx_perf_log_action (action),
    INDEX idx_perf_log_user (created_by),
    INDEX idx_perf_log_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint for evaluation_id (optional link)
SET @fk_log_eval = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_logs'
      AND CONSTRAINT_NAME = 'fk_perf_log_eval'
  ),
  'DO 0',
  'ALTER TABLE performance_logs ADD CONSTRAINT fk_perf_log_eval FOREIGN KEY (evaluation_id) REFERENCES performance_evaluations(id) ON DELETE SET NULL'
);
PREPARE stmt FROM @fk_log_eval;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
