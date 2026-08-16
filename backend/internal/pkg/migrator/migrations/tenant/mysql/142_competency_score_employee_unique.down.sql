-- 142_competency_score_employee_unique.down.sql

ALTER TABLE competency_scores
    DROP FOREIGN KEY fk_comp_score_org,
    DROP INDEX uk_comp_score_event_employee,
    DROP INDEX idx_comp_score_org,
    ADD UNIQUE KEY uk_comp_score_org (organization_id);

ALTER TABLE competency_scores
    ADD CONSTRAINT fk_comp_score_org FOREIGN KEY (organization_id)
        REFERENCES organizations(id) ON DELETE CASCADE;
