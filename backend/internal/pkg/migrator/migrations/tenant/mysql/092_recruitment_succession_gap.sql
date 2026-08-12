-- =============================================================================
-- Tenant Migration: 092_recruitment_succession_gap (MySQL)
-- =============================================================================
-- S-5 Succession Planning → Fallback External Recruitment
-- (docs/module-recruitment-strategic-layer-plan.md §4.5)
--
-- CI menandai posisi kunci tanpa successor siap (succession gap); Recruitment
-- menjadi fallback external: requisition dengan reason_type=SUCCESSION_GAP
-- menautkan ke posisi kunci tsb via succession_position_id. Recruitment TIDAK
-- menghitung readiness succession sendiri — ia membaca hasil CI melalui
-- interface narrow (wiring di cmd/server/main.go).
--
--   succession_position_id  CHAR(36)  NULL  → referensi posisi kunci (positions.id)
--                        yang membutuhkan fallback external recruitment
--
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_succession_position_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'succession_position_id'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN succession_position_id CHAR(36) NULL AFTER workforce_plan_id'
);
PREPARE stmt FROM @add_succession_position_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
