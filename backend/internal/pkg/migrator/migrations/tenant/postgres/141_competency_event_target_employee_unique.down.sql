-- 141_competency_event_target_employee_unique.down.sql

DROP INDEX IF EXISTS uk_comp_event_target_employee;

-- Kembalikan unique lama (event, organization). Sama seperti .sql-nya,
-- sebelum eksekusi pastikan tidak ada duplikat (event, organization).
CREATE UNIQUE INDEX IF NOT EXISTS uk_comp_event_target
    ON competency_event_targets (competency_event_id, organization_id);
