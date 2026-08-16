-- 140_competency_job_family_requirement.down.sql

ALTER TABLE job_family_competencies
    DROP COLUMN required_level,
    DROP COLUMN weight;
