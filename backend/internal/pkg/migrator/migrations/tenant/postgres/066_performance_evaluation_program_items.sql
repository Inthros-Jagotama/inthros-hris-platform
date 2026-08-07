-- 066_performance_evaluation_program_items.sql
-- KPI Enhancement Phase 3: employee-authored Program items per evaluation
-- (no HR template — the employee proposes title/target/formula themselves
-- when the org's PROGRAM component is enabled).

CREATE TABLE IF NOT EXISTS performance_evaluation_program_items (
    id                          CHAR(36)      NOT NULL PRIMARY KEY,
    performance_evaluation_id   CHAR(36)      NOT NULL REFERENCES performance_evaluations(id) ON DELETE CASCADE,
    title                       VARCHAR(255)  NOT NULL,
    formula_type                VARCHAR(30)   NOT NULL DEFAULT 'MANUAL',
    target                      DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    actual                      DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    achievement                 DECIMAL(5,2)  NOT NULL DEFAULT 0.00,
    score                       DECIMAL(5,2)  NOT NULL DEFAULT 0.00,
    sort_order                  SMALLINT      NOT NULL DEFAULT 0,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_prog_eval ON performance_evaluation_program_items (performance_evaluation_id);
