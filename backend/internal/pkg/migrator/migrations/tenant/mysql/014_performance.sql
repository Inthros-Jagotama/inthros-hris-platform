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
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_period_code (period_code),
    INDEX idx_perf_period_year (year)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Performance Perspectives (Perspektif BSC: Financial, Customer, Internal, L&G)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_perspectives (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NULL,
    sort_order  SMALLINT     NOT NULL DEFAULT 0,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Performance Templates (Template BSC per organisasi/posisi)
-- =========================================================================
CREATE TABLE IF NOT EXISTS performance_templates (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    organization_id CHAR(36)     NOT NULL,
    name            VARCHAR(200) NOT NULL,
    description     TEXT         NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    created_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_tmpl_org (organization_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    created_at                TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_ind_tmpl (performance_template_id),
    INDEX idx_perf_ind_persp (perspective_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    created_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_eval_emp (employee_id),
    INDEX idx_perf_eval_org (organization_id),
    INDEX idx_perf_eval_period (period_id),
    INDEX idx_perf_eval_tmpl (template_id),
    INDEX idx_perf_eval_sup (supervisor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    created_at                  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_detail_eval (performance_evaluation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    created_at                  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_perf_tgt_eval (performance_evaluation_id),
    INDEX idx_perf_tgt_ind (indicator_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
