-- =============================================================================
-- Tenant Migration Down: 089_training_planning (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS training_documents;
DROP TABLE IF EXISTS training_session_costs;
DROP TABLE IF EXISTS training_mandatories;
DROP TABLE IF EXISTS training_course_prerequisites;
DROP TABLE IF EXISTS training_course_competencies;
DROP TABLE IF EXISTS training_course_objectives;
DROP TABLE IF EXISTS training_requests;
DROP TABLE IF EXISTS training_needs;
DROP TABLE IF EXISTS training_plan_items;
DROP TABLE IF EXISTS training_plans;
