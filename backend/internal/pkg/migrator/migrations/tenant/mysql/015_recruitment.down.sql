-- 015_recruitment.down.sql
-- Rollback Recruitment & Onboarding (ATS) Module

DROP TABLE IF EXISTS onboarding_task_items;
DROP TABLE IF EXISTS employee_onboardings;
DROP TABLE IF EXISTS onboarding_task_templates;
DROP TABLE IF EXISTS interviews;
DROP TABLE IF EXISTS job_applications;
DROP TABLE IF EXISTS candidates;
DROP TABLE IF EXISTS job_requisitions;
