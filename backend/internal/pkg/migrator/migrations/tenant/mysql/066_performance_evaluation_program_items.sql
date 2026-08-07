-- 066_performance_evaluation_program_items.sql
-- KPI Enhancement Phase 3: employee-authored Program items per evaluation
-- (no HR template — the employee proposes title/target/formula themselves
-- when the org's PROGRAM component is enabled).

CREATE TABLE IF NOT EXISTS performance_evaluation_program_items (
    id                          CHAR(36)      NOT NULL PRIMARY KEY,
    performance_evaluation_id   CHAR(36)      NOT NULL,
    title                       VARCHAR(255)  NOT NULL,
    formula_type                VARCHAR(30)   NOT NULL DEFAULT 'MANUAL',
    target                      DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    actual                      DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    achievement                 DECIMAL(5,2)  NOT NULL DEFAULT 0.00,
    score                       DECIMAL(5,2)  NOT NULL DEFAULT 0.00,
    sort_order                  SMALLINT      NOT NULL DEFAULT 0,
    created_at                  TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                  TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_prog_eval (performance_evaluation_id),
    CONSTRAINT fk_perf_prog_eval FOREIGN KEY (performance_evaluation_id)
        REFERENCES performance_evaluations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
