-- =============================================================================
-- 144_competency_assessment_templates.sql
-- Competency 360 Module — Assessment Template (plan generik §5, §10).
-- Template menentukan APA yang dinilai (competency + required level + weight)
-- dan BAGAIMANA rater di-weights (rater type config) untuk sebuah assessment.
-- Template TIDAK menggantikan competency requirement dari position/job family.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 144.1 Competency Assessment Templates
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_assessment_templates (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(50)  NOT NULL,
    description TEXT         NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'active',
    scale_id    CHAR(36)     NULL,
    created_by  CHAR(36)     NULL,
    updated_by  CHAR(36)     NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_comp_tpl_code (code),
    INDEX idx_comp_tpl_status (status),
    INDEX idx_comp_tpl_scale (scale_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 144.2 Template Competencies (apa yang dinilai dalam template)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_assessment_template_competencies (
    id             CHAR(36)     PRIMARY KEY,
    template_id    CHAR(36)     NOT NULL,
    competency_id  CHAR(36)     NOT NULL,
    required_level SMALLINT     NULL,
    weight         DECIMAL(6,2) NOT NULL DEFAULT 1,
    sort_order     INT          NOT NULL DEFAULT 0,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_comp_tpl_comp (template_id, competency_id),
    INDEX idx_comp_tpl_comp_competency (competency_id),
    CONSTRAINT fk_comp_tpl_comp_tpl  FOREIGN KEY (template_id)   REFERENCES competency_assessment_templates(id) ON DELETE CASCADE,
    CONSTRAINT fk_comp_tpl_comp_comp FOREIGN KEY (competency_id) REFERENCES competencies(id)                    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 144.3 Template Rater Types (rater weight config per template — §10)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_assessment_template_rater_types (
    id          CHAR(36)     PRIMARY KEY,
    template_id CHAR(36)     NOT NULL,
    rater_type  VARCHAR(20)  NOT NULL,
    weight      DECIMAL(6,2) NOT NULL DEFAULT 0,
    min_rater   INT          NOT NULL DEFAULT 1,
    max_rater   INT          NULL,
    required    TINYINT(1)   NOT NULL DEFAULT 0,
    anonymous   TINYINT(1)   NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_comp_tpl_rater_type (template_id, rater_type),
    INDEX idx_comp_tpl_rater_type (rater_type),
    CONSTRAINT fk_comp_tpl_rater_tpl FOREIGN KEY (template_id) REFERENCES competency_assessment_templates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
