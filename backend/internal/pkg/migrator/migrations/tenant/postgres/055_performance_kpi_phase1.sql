-- 055_performance_kpi_phase1.sql
-- Phase 1 KPI Enhancement: Add new columns to performance tables
-- Supports: Period linking, Formula types, Rating integration, KPI snapshots

-- =========================================================================
-- 1. Update performance_templates
-- Add: period_id, effective_date, expired_date
-- =========================================================================
ALTER TABLE performance_templates
    ADD COLUMN IF NOT EXISTS period_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS effective_date DATE NULL,
    ADD COLUMN IF NOT EXISTS expired_date DATE NULL;

CREATE INDEX IF NOT EXISTS idx_perf_tmpl_period ON performance_templates (period_id);
CREATE INDEX IF NOT EXISTS idx_perf_tmpl_effective ON performance_templates (effective_date);

-- Add FK constraint for period_id (idempotent)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_tmpl_period'
    ) THEN
        ALTER TABLE performance_templates
            ADD CONSTRAINT fk_perf_tmpl_period
            FOREIGN KEY (period_id) REFERENCES performance_periods(id)
            ON DELETE SET NULL;
    END IF;
END $$;

-- =========================================================================
-- 2. Update performance_indicators
-- Add: code, formula_type, minimum_score, maximum_score, target_type, is_required
-- =========================================================================
ALTER TABLE performance_indicators
    ADD COLUMN IF NOT EXISTS code VARCHAR(50) NULL,
    ADD COLUMN IF NOT EXISTS formula_type VARCHAR(30) NOT NULL DEFAULT 'MANUAL',
    ADD COLUMN IF NOT EXISTS minimum_score DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS maximum_score DECIMAL(5,2) NOT NULL DEFAULT 100.00,
    ADD COLUMN IF NOT EXISTS target_type VARCHAR(30) NOT NULL DEFAULT 'NUMBER',
    ADD COLUMN IF NOT EXISTS is_required SMALLINT NOT NULL DEFAULT 1;

CREATE INDEX IF NOT EXISTS idx_perf_ind_code ON performance_indicators (code);

-- =========================================================================
-- 3. Update performance_evaluations
-- Add: rating_id, submitted_at, approved_at
-- =========================================================================
ALTER TABLE performance_evaluations
    ADD COLUMN IF NOT EXISTS rating_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS submitted_at TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP NULL;

CREATE INDEX IF NOT EXISTS idx_perf_eval_rating ON performance_evaluations (rating_id);

-- =========================================================================
-- 4. Update performance_evaluation_details
-- Add: indicator_id, indicator_name (snapshot), target (snapshot), actual, achievement, remarks
-- =========================================================================
ALTER TABLE performance_evaluation_details
    ADD COLUMN IF NOT EXISTS indicator_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS indicator_name VARCHAR(255) NULL,
    ADD COLUMN IF NOT EXISTS target DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS actual DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS achievement DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS remarks TEXT NULL;

CREATE INDEX IF NOT EXISTS idx_perf_detail_indicator ON performance_evaluation_details (indicator_id);

-- Add FK constraint for indicator_id (idempotent)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_detail_indicator'
    ) THEN
        ALTER TABLE performance_evaluation_details
            ADD CONSTRAINT fk_perf_detail_indicator
            FOREIGN KEY (indicator_id) REFERENCES performance_indicators(id)
            ON DELETE SET NULL;
    END IF;
END $$;
