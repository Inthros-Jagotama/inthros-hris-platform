-- =============================================================================
-- 145_competency_indicators.sql
-- Competency 360 Module — Assessment Indicator (plan generik §6).
-- Indicator/behavioral statement milik sebuah competency (mis. Leadership:
-- "Memberikan arahan yang jelas"). Template dapat memilih subset indicator
-- lewat competency_assessment_template_indicators.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 145.1 Competency Indicators
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_indicators (
    id            CHAR(36) PRIMARY KEY,
    competency_id CHAR(36)      NOT NULL,
    code          VARCHAR(50)   NULL,
    statement     VARCHAR(1000) NOT NULL,
    description   TEXT          NULL,
    status        VARCHAR(20)   NOT NULL DEFAULT 'active',
    sort_order    INT           NOT NULL DEFAULT 0,
    created_at    TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_comp_indicator_comp (competency_id),
    INDEX idx_comp_indicator_status (status),
    CONSTRAINT fk_comp_indicator_comp FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 145.2 Template Indicators (indicator yang dipakai per template)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_assessment_template_indicators (
    id           CHAR(36)     PRIMARY KEY,
    template_id  CHAR(36)     NOT NULL,
    indicator_id CHAR(36)     NOT NULL,
    weight       DECIMAL(6,2) NOT NULL DEFAULT 1,
    sort_order   INT          NOT NULL DEFAULT 0,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_comp_tpl_indicator (template_id, indicator_id),
    INDEX idx_comp_tpl_ind_indicator (indicator_id),
    CONSTRAINT fk_comp_tpl_ind_tpl FOREIGN KEY (template_id)  REFERENCES competency_assessment_templates(id) ON DELETE CASCADE,
    CONSTRAINT fk_comp_tpl_ind_ind FOREIGN KEY (indicator_id) REFERENCES competency_indicators(id)           ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
