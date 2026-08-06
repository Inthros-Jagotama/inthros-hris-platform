-- Down migration for OKR Module

-- Drop indexes first
DROP INDEX IF EXISTS idx_okr_attach_user;
DROP INDEX IF EXISTS idx_okr_attach_detail;
DROP INDEX IF EXISTS idx_okr_comment_user;
DROP INDEX IF EXISTS idx_okr_comment_parent;
DROP INDEX IF EXISTS idx_okr_comment_eval;
DROP INDEX IF EXISTS idx_okr_prog_date;
DROP INDEX IF EXISTS idx_okr_prog_detail;
DROP INDEX IF EXISTS idx_okr_detail_kr;
DROP INDEX IF EXISTS idx_okr_detail_obj;
DROP INDEX IF EXISTS idx_okr_detail_eval;
DROP INDEX IF EXISTS idx_okr_eval_status;
DROP INDEX IF EXISTS idx_okr_eval_period;
DROP INDEX IF EXISTS idx_okr_eval_org;
DROP INDEX IF EXISTS idx_okr_eval_employee;
DROP INDEX IF EXISTS idx_okr_kr_sort;
DROP INDEX IF EXISTS idx_okr_kr_objective;
DROP INDEX IF EXISTS idx_okr_obj_sort;
DROP INDEX IF EXISTS idx_okr_obj_template;
DROP INDEX IF EXISTS idx_okr_tpl_status;
DROP INDEX IF EXISTS idx_okr_tpl_period;
DROP INDEX IF EXISTS idx_okr_tpl_org;

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS okr_attachments;
DROP TABLE IF EXISTS okr_comments;
DROP TABLE IF EXISTS okr_progress;
DROP TABLE IF EXISTS okr_evaluation_details;
DROP TABLE IF EXISTS okr_evaluations;
DROP TABLE IF EXISTS okr_key_results;
DROP TABLE IF EXISTS okr_objectives;
DROP TABLE IF EXISTS okr_templates;
