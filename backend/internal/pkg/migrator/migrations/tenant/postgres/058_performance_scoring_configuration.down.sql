-- 058_performance_scoring_configuration.down.sql
-- Rollback Phase 5 KPI Enhancement: Drop scoring configuration tables

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_evalcomp_component'
    ) THEN
        ALTER TABLE performance_evaluation_components DROP CONSTRAINT fk_perf_evalcomp_component;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_evalcomp_eval'
    ) THEN
        ALTER TABLE performance_evaluation_components DROP CONSTRAINT fk_perf_evalcomp_eval;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_orgcomp_component'
    ) THEN
        ALTER TABLE performance_organization_components DROP CONSTRAINT fk_perf_orgcomp_component;
    END IF;
END $$;

DROP TABLE IF EXISTS performance_evaluation_components;
DROP TABLE IF EXISTS performance_organization_components;
DROP TABLE IF EXISTS performance_components;