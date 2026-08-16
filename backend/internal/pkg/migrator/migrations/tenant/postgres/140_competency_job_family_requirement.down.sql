-- 140_competency_job_family_requirement.down.sql

ALTER TABLE job_family_competencies
    DROP COLUMN IF EXISTS required_level,
    DROP COLUMN IF EXISTS weight;
