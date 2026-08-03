-- 014_performance.sql
-- Performance Management Module (KPI, OKR, BSC)
-- Tabel untuk balanced scorecard performance evaluation

-- =========================================================================
-- Performance Periods (Periode evaluasi: tahunan, kuartal, bulanan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_periods (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    period_code VARCHAR(10)  NOT NULL DEFAULT '',
    period_type VARCHAR(20)  NOT NULL DEFAULT '',
    year        SMALLINT     NOT NULL DEFAULT 0,
    start_date  DATE         NULL,
    end_date    DATE         NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_period_code ON performance_periods (period_code);

CREATE INDEX IF NOT EXISTS idx_perf_period_year ON performance_periods (year);

-- =========================================================================
-- Performance Perspectives (Perspektif BSC: Financial, Customer, Internal, L&G)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_perspectives (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NULL,
    sort_order  SMALLINT     NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =========================================================================
-- Performance Templates (Template BSC per organisasi/posisi)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_templates (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    organization_id CHAR(36)     NOT NULL,
    name            VARCHAR(200) NOT NULL,
    description     TEXT         NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_tmpl_org ON performance_templates (organization_id);

-- =========================================================================
-- Performance Indicators (Indikator KPI — linked to template)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_indicators (
    id                        CHAR(36)     NOT NULL PRIMARY KEY,
    performance_template_id   CHAR(36)     NOT NULL,
    perspective_id            CHAR(36)     NOT NULL,
    indicator_type            VARCHAR(20)  NOT NULL DEFAULT '',
    title                     VARCHAR(255) NOT NULL,
    description               TEXT         NULL,
    weight                    DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    target_value              DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    unit_of_measurement       VARCHAR(50)  NULL,
    sort_order                SMALLINT     NOT NULL DEFAULT 0,
    created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_ind_tmpl ON performance_indicators (performance_template_id);

CREATE INDEX IF NOT EXISTS idx_perf_ind_persp ON performance_indicators (perspective_id);

-- =========================================================================
-- Performance Evaluations (Evaluasi kinerja karyawan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_evaluations (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id         CHAR(36)     NOT NULL,
    organization_id     CHAR(36)     NOT NULL,
    period_id           CHAR(36)     NOT NULL,
    template_id         CHAR(36)     NOT NULL,
    supervisor_id       CHAR(36)     NULL,
    final_score         DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    status              VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    plan_submitted_at   BIGINT       NULL DEFAULT 0,
    plan_approved_at    BIGINT       NULL DEFAULT 0,
    actual_submitted_at BIGINT       NULL DEFAULT 0,
    actual_approved_at  BIGINT       NULL DEFAULT 0,
    notes               TEXT         NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_eval_emp ON performance_evaluations (employee_id);

CREATE INDEX IF NOT EXISTS idx_perf_eval_org ON performance_evaluations (organization_id);

CREATE INDEX IF NOT EXISTS idx_perf_eval_period ON performance_evaluations (period_id);

CREATE INDEX IF NOT EXISTS idx_perf_eval_tmpl ON performance_evaluations (template_id);

CREATE INDEX IF NOT EXISTS idx_perf_eval_sup ON performance_evaluations (supervisor_id);

-- =========================================================================
-- Performance Evaluation Details (Detail grup BSC dalam evaluasi)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_evaluation_details (
    id                          CHAR(36)     NOT NULL PRIMARY KEY,
    performance_evaluation_id   CHAR(36)     NOT NULL,
    perspective_id              CHAR(36)     NOT NULL,
    achievement_percentage      DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    weight                      DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    score                       DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    description                 VARCHAR(255) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_detail_eval ON performance_evaluation_details (performance_evaluation_id);

-- =========================================================================
-- Performance Targets (Target KPI individual — nilai target vs aktual)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_targets (
    id                          CHAR(36)     NOT NULL PRIMARY KEY,
    performance_evaluation_id   CHAR(36)     NOT NULL,
    indicator_id                CHAR(36)     NOT NULL,
    target_value                DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    actual_value                DECIMAL(12,2) NULL,
    unit_of_measurement         VARCHAR(50)  NULL,
    achievement_percentage      DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    weight                      DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    score                       DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perf_tgt_eval ON performance_targets (performance_evaluation_id);

CREATE INDEX IF NOT EXISTS idx_perf_tgt_ind ON performance_targets (indicator_id);
