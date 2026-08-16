-- 142_competency_score_employee_unique.down.sql

DROP INDEX IF EXISTS uk_comp_score_event_employee;

-- Kembalikan unique lama (organization_id). Sebelum eksekusi pastikan tidak
-- ada duplikat organization_id.
CREATE UNIQUE INDEX IF NOT EXISTS uk_comp_score_org
    ON competency_scores (organization_id);
