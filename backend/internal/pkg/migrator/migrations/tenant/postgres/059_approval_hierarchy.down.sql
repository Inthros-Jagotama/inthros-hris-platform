-- 059_approval_hierarchy.down.sql

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_approval_step_org_step'
    ) THEN
        ALTER TABLE approval_flow_step_organizations DROP CONSTRAINT fk_approval_step_org_step;
    END IF;
END $$;

DROP TABLE IF EXISTS approval_flow_step_organizations;

ALTER TABLE approval_flow_steps
    DROP COLUMN IF EXISTS participation_type,
    DROP COLUMN IF EXISTS hierarchy_level;
