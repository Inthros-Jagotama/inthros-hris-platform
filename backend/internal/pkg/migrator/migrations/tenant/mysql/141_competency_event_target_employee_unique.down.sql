-- 141_competency_event_target_employee_unique.down.sql

ALTER TABLE competency_event_targets
    DROP INDEX uk_comp_event_target_employee,
    ADD UNIQUE KEY uk_comp_event_target (competency_event_id, organization_id);
