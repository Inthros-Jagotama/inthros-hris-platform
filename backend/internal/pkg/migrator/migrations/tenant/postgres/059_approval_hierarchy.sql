-- 059_approval_hierarchy.sql
-- Approval Module Phase 1: Organization-hierarchy approver resolution +
-- step participation type (APPROVER vs WATCHER).

ALTER TABLE approval_flow_steps
    ADD COLUMN IF NOT EXISTS hierarchy_level INT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS participation_type VARCHAR(20) NOT NULL DEFAULT 'APPROVER';

-- ---------------------------------------------------------------------------
-- approval_flow_step_organizations — target Organization(s) untuk step dengan
-- approver_type ORGANIZATION/BOTH. Satu step bisa punya banyak Organization
-- (mendukung "lebih dari satu approver di satu level").
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_flow_step_organizations (
    id              CHAR(36) PRIMARY KEY,
    step_id         CHAR(36) NOT NULL,
    organization_id CHAR(36) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_approval_step_org_step ON approval_flow_step_organizations (step_id);
CREATE INDEX IF NOT EXISTS idx_approval_step_org_org ON approval_flow_step_organizations (organization_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_approval_step_org_step'
    ) THEN
        ALTER TABLE approval_flow_step_organizations
            ADD CONSTRAINT fk_approval_step_org_step
            FOREIGN KEY (step_id) REFERENCES approval_flow_steps(id)
            ON DELETE CASCADE;
    END IF;
END $$;
