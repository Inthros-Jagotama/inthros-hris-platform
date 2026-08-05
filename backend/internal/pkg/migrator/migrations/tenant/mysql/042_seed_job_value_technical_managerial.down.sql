-- Down Migration: 042_seed_job_value_technical_managerial
-- Menghapus row seed level Technical dan Managerial.
--
-- Mencakup dua generasi seed via kolom `type` (selalu ada di semua tahap
-- rollback — TIDAK memakai type_group karena kolom tsb di-drop oleh 052 down
-- yang berjalan lebih dulu pada rollback fresh install):
--   1. type literal 'technical'/'managerial'  → seed baru (up migration 042)
--   2. slug per kompetensi (16 technical + 6 managerial) → seed lama, tenant
--      yang belum di-reseed (mis. 'budgeting', 'integrity', dst.)
-- Aman dijalankan kapan pun: fresh rollback, manual down, maupun tenant lama.

DELETE FROM job_management_values
WHERE type IN (
    'technical', 'managerial',
    'competency_based_human_resources_management', 'competency_development', 'people_development',
    'career_management', 'hr_assessment', 'recruitement_selection', 'job_analysis_evaluation',
    'organizational_development', 'human_resources_information_system', 'workload_analysis',
    'performance_apraisal', 'remuneration_manajemen', 'reward_punisment_management',
    'health_safety_environment', 'hubungan_industrial', 'budgeting',
    'integrity', 'achievement_orientation', 'building_partnership', 'planning_organizing',
    'leadership', 'developing_others'
);
