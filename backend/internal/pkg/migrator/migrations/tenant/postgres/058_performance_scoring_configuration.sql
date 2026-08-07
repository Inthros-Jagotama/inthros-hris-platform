-- 058_performance_scoring_configuration.sql
-- Phase 5 KPI Enhancement: Performance Scoring Configuration
-- Tables: performance_components, performance_organization_components,
--         performance_evaluation_components

-- =========================================================================
-- 1. performance_components - Master komponen penilaian (KPI, Work Program, dst)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_components (
    id CHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_perf_comp_code ON performance_components (code);
CREATE INDEX IF NOT EXISTS idx_perf_comp_sort ON performance_components (sort_order);
CREATE INDEX IF NOT EXISTS idx_perf_comp_deleted ON performance_components (deleted_at);

-- =========================================================================
-- 2. performance_organization_components - Konfigurasi bobot komponen per Organization
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_organization_components (
    id CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    component_id CHAR(36) NOT NULL,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_orgcomp_org ON performance_organization_components (organization_id);
CREATE INDEX IF NOT EXISTS idx_perf_orgcomp_comp ON performance_organization_components (component_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_perf_orgcomp_unique ON performance_organization_components (organization_id, component_id);

-- Add FK constraint for component_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_orgcomp_component'
    ) THEN
        ALTER TABLE performance_organization_components
            ADD CONSTRAINT fk_perf_orgcomp_component
            FOREIGN KEY (component_id) REFERENCES performance_components(id)
            ON DELETE CASCADE;
    END IF;
END $$;

-- =========================================================================
-- 3. performance_evaluation_components - Snapshot hasil perhitungan komponen
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_evaluation_components (
    id CHAR(36) PRIMARY KEY,
    evaluation_id CHAR(36) NOT NULL,
    component_id CHAR(36) NOT NULL,
    component_name VARCHAR(100) NOT NULL,
    score DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    final_score DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    calculated_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_evalcomp_eval ON performance_evaluation_components (evaluation_id);
CREATE INDEX IF NOT EXISTS idx_perf_evalcomp_comp ON performance_evaluation_components (component_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_perf_evalcomp_unique ON performance_evaluation_components (evaluation_id, component_id);

-- Add FK constraint for evaluation_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_evalcomp_eval'
    ) THEN
        ALTER TABLE performance_evaluation_components
            ADD CONSTRAINT fk_perf_evalcomp_eval
            FOREIGN KEY (evaluation_id) REFERENCES performance_evaluations(id)
            ON DELETE CASCADE;
    END IF;
END $$;

-- Add FK constraint for component_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_evalcomp_component'
    ) THEN
        ALTER TABLE performance_evaluation_components
            ADD CONSTRAINT fk_perf_evalcomp_component
            FOREIGN KEY (component_id) REFERENCES performance_components(id)
            ON DELETE CASCADE;
    END IF;
END $$;
