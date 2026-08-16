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
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_comp_tpl_code ON competency_assessment_templates (code);

CREATE INDEX IF NOT EXISTS idx_comp_tpl_status ON competency_assessment_templates (status);

CREATE INDEX IF NOT EXISTS idx_comp_tpl_scale ON competency_assessment_templates (scale_id);

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
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_comp_tpl_comp UNIQUE (template_id, competency_id),

    CONSTRAINT fk_comp_tpl_comp_tpl    FOREIGN KEY (template_id)   REFERENCES competency_assessment_templates(id) ON DELETE CASCADE,
    CONSTRAINT fk_comp_tpl_comp_comp   FOREIGN KEY (competency_id) REFERENCES competencies(id)                    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comp_tpl_comp_competency ON competency_assessment_template_competencies (competency_id);

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
    required    BOOLEAN      NOT NULL DEFAULT FALSE,
    anonymous   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_comp_tpl_rater_type UNIQUE (template_id, rater_type),

    CONSTRAINT fk_comp_tpl_rater_tpl FOREIGN KEY (template_id) REFERENCES competency_assessment_templates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comp_tpl_rater_type ON competency_assessment_template_rater_types (rater_type);

-- ---------------------------------------------------------------------------
-- 144.4 Link Event ke Template
-- `competency_events` (migration 008) mendapat template_id: event menentukan
-- template assessment mana yang dipakai untuk menilai subject-nya.
-- ---------------------------------------------------------------------------
ALTER TABLE competency_events
    ADD COLUMN IF NOT EXISTS template_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_comp_event_template ON competency_events (template_id);

ALTER TABLE competency_events
    ADD CONSTRAINT fk_comp_event_template FOREIGN KEY (template_id) REFERENCES competency_assessment_templates(id) ON DELETE SET NULL;
