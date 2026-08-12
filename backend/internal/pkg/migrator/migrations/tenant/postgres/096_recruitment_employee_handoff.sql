-- =============================================================================
-- Tenant Migration: 096_recruitment_employee_handoff (PostgreSQL)
-- =============================================================================
-- G-4 Recruitment → Employee / Employee Movement
-- (docs/module-recruitment-development-plan.md §G-4)
--
--   employee.recruited_from_application_id CHAR(36) NULL
--                        → referensi job_applications saat employee dibuat dari
--                          offer eksternal yang diterima; menelusuri
--                          Employee → Application → Requisition → Position
--   candidates.candidate_type VARCHAR(20) NOT NULL DEFAULT 'EXTERNAL'
--                        (EXTERNAL | INTERNAL) — kandidat internal menunjuk
--                        employee yang sudah ada (tidak dibuatkan employee baru)
--   candidates.employee_id CHAR(36) NULL
--                        → referensi employee untuk kandidat INTERNAL; saat
--                        offer internal diterima, hasil seleksi diteruskan ke
--                        Employee Movement (bukan employee baru)
--
-- Idempotent: ADD COLUMN IF NOT EXISTS.

ALTER TABLE employees
    ADD COLUMN IF NOT EXISTS recruited_from_application_id CHAR(36) NULL;

ALTER TABLE candidates
    ADD COLUMN IF NOT EXISTS candidate_type VARCHAR(20) NOT NULL DEFAULT 'EXTERNAL',
    ADD COLUMN IF NOT EXISTS employee_id CHAR(36) NULL;
