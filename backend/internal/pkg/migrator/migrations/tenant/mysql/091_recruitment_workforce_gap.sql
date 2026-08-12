-- =============================================================================
-- Tenant Migration: 091_recruitment_workforce_gap (MySQL)
-- =============================================================================
-- S-1 Workforce Gap → Requisition
-- (docs/module-recruitment-strategic-layer-plan.md §4.1)
--
-- Recruitment membaca hiring need dari Workforce Intelligence melalui interface
-- narrow (wiring di cmd/server/main.go) — Recruitment TIDAK menghitung gap
-- sendiri. Kolom berikut hanya menyimpan referensi & alasan pembuatan
-- requisition dari sisi strategis:
--
--   reason_type        VARCHAR(30)  NULL  → alasan pembuatan requisition
--                        (NEW_POSITION | REPLACEMENT | EXPANSION | WORKFORCE_GAP)
--   workforce_gap_id   CHAR(36)     NULL  → referensi gap analysis Workforce
--                        Intelligence yang memicu requisition ini
--   workforce_plan_id  CHAR(36)     NULL  → referensi headcount plan (period)
--                        tempat gap teridentifikasi
--
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_reason_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'reason_type'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN reason_type VARCHAR(30) NULL AFTER approved_by'
);
PREPARE stmt FROM @add_reason_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_workforce_gap_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'workforce_gap_id'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN workforce_gap_id CHAR(36) NULL AFTER reason_type'
);
PREPARE stmt FROM @add_workforce_gap_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_workforce_plan_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'workforce_plan_id'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN workforce_plan_id CHAR(36) NULL AFTER workforce_gap_id'
);
PREPARE stmt FROM @add_workforce_plan_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
