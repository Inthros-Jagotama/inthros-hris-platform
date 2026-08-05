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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_prog_detail ON performance_progress (evaluation_detail_id);
CREATE INDEX IF NOT EXISTS idx_perf_prog_date ON performance_progress (progress_date);
CREATE INDEX IF NOT EXISTS idx_perf_prog_created_by ON performance_progress (created_by);

-- Add FK constraint for evaluation_detail_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_prog_detail'
    ) THEN
        ALTER TABLE performance_progress
            ADD CONSTRAINT fk_perf_prog_detail
            FOREIGN KEY (evaluation_detail_id) REFERENCES performance_evaluation_details(id)
            ON DELETE CASCADE;
    END IF;
END $$;

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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_comment_eval ON performance_comments (evaluation_id);
CREATE INDEX IF NOT EXISTS idx_perf_comment_emp ON performance_comments (employee_id);
CREATE INDEX IF NOT EXISTS idx_perf_comment_created_by ON performance_comments (created_by);

-- Add FK constraint for evaluation_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_comment_eval'
    ) THEN
        ALTER TABLE performance_comments
            ADD CONSTRAINT fk_perf_comment_eval
            FOREIGN KEY (evaluation_id) REFERENCES performance_evaluations(id)
            ON DELETE CASCADE;
    END IF;
END $$;

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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_attach_detail ON performance_attachments (evaluation_detail_id);
CREATE INDEX IF NOT EXISTS idx_perf_attach_uploaded_by ON performance_attachments (uploaded_by);

-- Add FK constraint for evaluation_detail_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_attach_detail'
    ) THEN
        ALTER TABLE performance_attachments
            ADD CONSTRAINT fk_perf_attach_detail
            FOREIGN KEY (evaluation_detail_id) REFERENCES performance_evaluation_details(id)
            ON DELETE CASCADE;
    END IF;
END $$;

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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_perf_rating_code ON performance_ratings (code);
CREATE INDEX IF NOT EXISTS idx_perf_rating_sort ON performance_ratings (sort_order);

-- Add FK constraint from performance_evaluations to performance_ratings
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_eval_rating'
    ) THEN
        ALTER TABLE performance_evaluations
            ADD CONSTRAINT fk_perf_eval_rating
            FOREIGN KEY (rating_id) REFERENCES performance_ratings(id)
            ON DELETE SET NULL;
    END IF;
END $$;

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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_perf_formula_code ON performance_indicator_formulas (code);
CREATE INDEX IF NOT EXISTS idx_perf_formula_type ON performance_indicator_formulas (formula_type);
CREATE INDEX IF NOT EXISTS idx_perf_formula_sort ON performance_indicator_formulas (sort_order);

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
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_log_eval ON performance_logs (evaluation_id);
CREATE INDEX IF NOT EXISTS idx_perf_log_entity ON performance_logs (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_perf_log_action ON performance_logs (action);
CREATE INDEX IF NOT EXISTS idx_perf_log_user ON performance_logs (created_by);
CREATE INDEX IF NOT EXISTS idx_perf_log_created ON performance_logs (created_at);

-- Add FK constraint for evaluation_id (optional link)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_log_eval'
    ) THEN
        ALTER TABLE performance_logs
            ADD CONSTRAINT fk_perf_log_eval
            FOREIGN KEY (evaluation_id) REFERENCES performance_evaluations(id)
            ON DELETE SET NULL;
    END IF;
END $$;
