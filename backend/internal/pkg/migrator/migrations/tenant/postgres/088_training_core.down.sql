-- =============================================================================
-- Tenant Migration Down: 088_training_core
-- =============================================================================

-- Assessment results & assessments
DROP INDEX IF EXISTS uk_trn_assess_res_attempt;
DROP INDEX IF EXISTS idx_trn_assess_res_part;
DROP INDEX IF EXISTS idx_trn_assess_res_assessment;
DROP TABLE IF EXISTS training_assessment_results;

DROP INDEX IF EXISTS idx_trn_assess_session;
DROP TABLE IF EXISTS training_assessments;

-- Attendances
DROP INDEX IF EXISTS uk_trn_att_part_date;
DROP INDEX IF EXISTS idx_trn_att_part;
DROP TABLE IF EXISTS training_attendances;

-- Materials enhancement
ALTER TABLE training_materials DROP COLUMN IF EXISTS available_from;
ALTER TABLE training_materials DROP COLUMN IF EXISTS is_required;
ALTER TABLE training_materials DROP COLUMN IF EXISTS description;

-- Participants enhancement
DROP INDEX IF EXISTS uk_trn_part_session_employee;
ALTER TABLE training_participants DROP COLUMN IF EXISTS remarks;
ALTER TABLE training_participants DROP COLUMN IF EXISTS passed;
ALTER TABLE training_participants DROP COLUMN IF EXISTS final_score;
ALTER TABLE training_participants DROP COLUMN IF EXISTS completion_date;
ALTER TABLE training_participants DROP COLUMN IF EXISTS completion_status;
ALTER TABLE training_participants DROP COLUMN IF EXISTS approved_at;
ALTER TABLE training_participants DROP COLUMN IF EXISTS registered_at;
ALTER TABLE training_participants DROP COLUMN IF EXISTS registration_status;

-- Sessions enhancement
DROP INDEX IF EXISTS idx_trn_sess_provider;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS registration_deadline;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS meeting_url;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS end_datetime;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS start_datetime;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS provider_id;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS delivery_mode;
ALTER TABLE training_sessions DROP COLUMN IF EXISTS provider_type;

-- Courses enhancement
ALTER TABLE training_courses DROP COLUMN IF EXISTS is_mandatory;
ALTER TABLE training_courses DROP COLUMN IF EXISTS delivery_type;
ALTER TABLE training_courses DROP COLUMN IF EXISTS course_type;

-- Session trainers
DROP INDEX IF EXISTS idx_trn_sess_trn_trainer;
DROP INDEX IF EXISTS idx_trn_sess_trn_session;
DROP TABLE IF EXISTS training_session_trainers;

-- Trainers & providers
DROP INDEX IF EXISTS idx_trn_trainer_provider;
DROP INDEX IF EXISTS idx_trn_trainer_emp;
DROP TABLE IF EXISTS training_trainers;

DROP INDEX IF EXISTS idx_trn_provider_type;
DROP INDEX IF EXISTS idx_trn_provider_code;
DROP TABLE IF EXISTS training_providers;
