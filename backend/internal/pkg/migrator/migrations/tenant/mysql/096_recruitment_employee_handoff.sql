-- =============================================================================
-- Tenant Migration: 096_recruitment_employee_handoff (MySQL)
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
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_recruited_from_application_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employees'
      AND COLUMN_NAME = 'recruited_from_application_id'
  ),
  'DO 0',
  'ALTER TABLE employees ADD COLUMN recruited_from_application_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_recruited_from_application_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_candidate_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'candidate_type'
  ),
  'DO 0',
  'ALTER TABLE candidates ADD COLUMN candidate_type VARCHAR(20) NOT NULL DEFAULT ''EXTERNAL'''
);
PREPARE stmt FROM @add_candidate_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_candidate_employee_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'employee_id'
  ),
  'DO 0',
  'ALTER TABLE candidates ADD COLUMN employee_id CHAR(36) NULL AFTER candidate_type'
);
PREPARE stmt FROM @add_candidate_employee_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
