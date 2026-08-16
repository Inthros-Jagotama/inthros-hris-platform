-- 142_competency_score_employee_unique.down.sql

ALTER TABLE competency_scores
    DROP INDEX uk_comp_score_event_employee,
    ADD UNIQUE KEY uk_comp_score_org (organization_id);
