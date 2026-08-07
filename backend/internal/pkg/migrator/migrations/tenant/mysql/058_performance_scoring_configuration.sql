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
    sort_order INT NOT NULL DEFAULT 0,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE INDEX idx_perf_comp_code (code),
    INDEX idx_perf_comp_sort (sort_order),
    INDEX idx_perf_comp_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- 2. performance_organization_components - Konfigurasi bobot komponen per Organization
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_organization_components (
    id CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    component_id CHAR(36) NOT NULL,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    is_enabled TINYINT(1) NOT NULL DEFAULT 1,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_perf_orgcomp_org (organization_id),
    INDEX idx_perf_orgcomp_comp (component_id),
    UNIQUE INDEX idx_perf_orgcomp_unique (organization_id, component_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint for component_id
SET @fk_orgcomp_component = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_organization_components'
      AND CONSTRAINT_NAME = 'fk_perf_orgcomp_component'
  ),
  'DO 0',
  'ALTER TABLE performance_organization_components ADD CONSTRAINT fk_perf_orgcomp_component FOREIGN KEY (component_id) REFERENCES performance_components(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_orgcomp_component;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_perf_evalcomp_eval (evaluation_id),
    INDEX idx_perf_evalcomp_comp (component_id),
    UNIQUE INDEX idx_perf_evalcomp_unique (evaluation_id, component_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add FK constraint for evaluation_id
SET @fk_evalcomp_eval = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluation_components'
      AND CONSTRAINT_NAME = 'fk_perf_evalcomp_eval'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_components ADD CONSTRAINT fk_perf_evalcomp_eval FOREIGN KEY (evaluation_id) REFERENCES performance_evaluations(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_evalcomp_eval;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add FK constraint for component_id
SET @fk_evalcomp_component = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'performance_evaluation_components'
      AND CONSTRAINT_NAME = 'fk_perf_evalcomp_component'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_components ADD CONSTRAINT fk_perf_evalcomp_component FOREIGN KEY (component_id) REFERENCES performance_components(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_evalcomp_component;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
