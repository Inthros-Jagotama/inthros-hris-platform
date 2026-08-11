-- =============================================================================
-- Tenant Migration Down: 088_training_core (MySQL)
-- =============================================================================

-- Assessment results & assessments
DROP TABLE IF EXISTS training_assessment_results;
DROP TABLE IF EXISTS training_assessments;

-- Attendances
DROP TABLE IF EXISTS training_attendances;

-- Materials enhancement
SET @drop_mat_available_from = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_materials' AND COLUMN_NAME='available_from'),
  'ALTER TABLE training_materials DROP COLUMN available_from',
  'DO 0'
);
PREPARE stmt FROM @drop_mat_available_from; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_mat_is_required = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_materials' AND COLUMN_NAME='is_required'),
  'ALTER TABLE training_materials DROP COLUMN is_required',
  'DO 0'
);
PREPARE stmt FROM @drop_mat_is_required; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_mat_description = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_materials' AND COLUMN_NAME='description'),
  'ALTER TABLE training_materials DROP COLUMN description',
  'DO 0'
);
PREPARE stmt FROM @drop_mat_description; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Participants enhancement
SET @drop_part_uk = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND INDEX_NAME='uk_trn_part_active_key'),
  'ALTER TABLE training_participants DROP INDEX uk_trn_part_active_key',
  'DO 0'
);
PREPARE stmt FROM @drop_part_uk; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_active_key = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='active_session_employee'),
  'ALTER TABLE training_participants DROP COLUMN active_session_employee',
  'DO 0'
);
PREPARE stmt FROM @drop_active_key; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_remarks = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='remarks'),
  'ALTER TABLE training_participants DROP COLUMN remarks',
  'DO 0'
);
PREPARE stmt FROM @drop_part_remarks; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_passed = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='passed'),
  'ALTER TABLE training_participants DROP COLUMN passed',
  'DO 0'
);
PREPARE stmt FROM @drop_part_passed; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_final_score = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='final_score'),
  'ALTER TABLE training_participants DROP COLUMN final_score',
  'DO 0'
);
PREPARE stmt FROM @drop_part_final_score; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_completion_date = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='completion_date'),
  'ALTER TABLE training_participants DROP COLUMN completion_date',
  'DO 0'
);
PREPARE stmt FROM @drop_part_completion_date; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_completion_status = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='completion_status'),
  'ALTER TABLE training_participants DROP COLUMN completion_status',
  'DO 0'
);
PREPARE stmt FROM @drop_part_completion_status; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_approved_at = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='approved_at'),
  'ALTER TABLE training_participants DROP COLUMN approved_at',
  'DO 0'
);
PREPARE stmt FROM @drop_part_approved_at; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_registered_at = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='registered_at'),
  'ALTER TABLE training_participants DROP COLUMN registered_at',
  'DO 0'
);
PREPARE stmt FROM @drop_part_registered_at; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_part_reg_status = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='registration_status'),
  'ALTER TABLE training_participants DROP COLUMN registration_status',
  'DO 0'
);
PREPARE stmt FROM @drop_part_reg_status; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Sessions enhancement
SET @drop_sess_provider_index = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND INDEX_NAME='idx_trn_sess_provider'),
  'ALTER TABLE training_sessions DROP INDEX idx_trn_sess_provider',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_provider_index; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_reg_deadline = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='registration_deadline'),
  'ALTER TABLE training_sessions DROP COLUMN registration_deadline',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_reg_deadline; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_meeting_url = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='meeting_url'),
  'ALTER TABLE training_sessions DROP COLUMN meeting_url',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_meeting_url; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_end_datetime = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='end_datetime'),
  'ALTER TABLE training_sessions DROP COLUMN end_datetime',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_end_datetime; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_start_datetime = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='start_datetime'),
  'ALTER TABLE training_sessions DROP COLUMN start_datetime',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_start_datetime; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_provider_id = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='provider_id'),
  'ALTER TABLE training_sessions DROP COLUMN provider_id',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_provider_id; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_delivery_mode = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='delivery_mode'),
  'ALTER TABLE training_sessions DROP COLUMN delivery_mode',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_delivery_mode; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_sess_provider_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='provider_type'),
  'ALTER TABLE training_sessions DROP COLUMN provider_type',
  'DO 0'
);
PREPARE stmt FROM @drop_sess_provider_type; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Courses enhancement
SET @drop_course_mandatory = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_courses' AND COLUMN_NAME='is_mandatory'),
  'ALTER TABLE training_courses DROP COLUMN is_mandatory',
  'DO 0'
);
PREPARE stmt FROM @drop_course_mandatory; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_course_delivery_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_courses' AND COLUMN_NAME='delivery_type'),
  'ALTER TABLE training_courses DROP COLUMN delivery_type',
  'DO 0'
);
PREPARE stmt FROM @drop_course_delivery_type; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_course_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_courses' AND COLUMN_NAME='course_type'),
  'ALTER TABLE training_courses DROP COLUMN course_type',
  'DO 0'
);
PREPARE stmt FROM @drop_course_type; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Session trainers
DROP TABLE IF EXISTS training_session_trainers;

-- Trainers & providers
DROP TABLE IF EXISTS training_trainers;
DROP TABLE IF EXISTS training_providers;
