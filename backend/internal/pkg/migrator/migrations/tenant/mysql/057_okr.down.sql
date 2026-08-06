-- Down migration for OKR Module

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS okr_attachments;
DROP TABLE IF EXISTS okr_comments;
DROP TABLE IF EXISTS okr_progress;
DROP TABLE IF EXISTS okr_evaluation_details;
DROP TABLE IF EXISTS okr_evaluations;
DROP TABLE IF EXISTS okr_key_results;
DROP TABLE IF EXISTS okr_objectives;
DROP TABLE IF EXISTS okr_templates;
