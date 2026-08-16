-- =============================================================================
-- 147_competency_assessment_responses.sql
-- Competency 360 Module — Assessment Response (plan generik §11).
-- Menyimpan jawaban mentah per rater per indicator. Immutable setelah
-- submit — perubahan hanya via mekanisme correction/reopen yang terkontrol.
-- =============================================================================

CREATE TABLE IF NOT EXISTS competency_assessment_responses (
    id           CHAR(36) PRIMARY KEY,
    rater_id     CHAR(36)     NOT NULL,
    indicator_id CHAR(36)     NOT NULL,
    rating_value SMALLINT     NOT NULL,
    comment      TEXT         NULL,
    submitted_at TIMESTAMP    NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_comp_response_rater_indicator (rater_id, indicator_id),
    INDEX idx_comp_response_indicator (indicator_id),
    CONSTRAINT fk_comp_response_rater     FOREIGN KEY (rater_id)     REFERENCES competency_assessment_raters(id)  ON DELETE CASCADE,
    CONSTRAINT fk_comp_response_indicator FOREIGN KEY (indicator_id) REFERENCES competency_indicators(id)          ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
