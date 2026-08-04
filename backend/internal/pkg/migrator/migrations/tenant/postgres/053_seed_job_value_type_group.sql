-- Migration: 053_seed_job_value_type_group
-- Melengkapi kolom type_group & description_group pada job_management_values
-- untuk semua tipe yang sudah ada (seed 033–050 & data tenant).
--
-- Mapping mengikuti codeMap legacy docs/seeder/JobManagementValuesTableSeeder.php:
--   type_group        = slug kelompok (dipakai FE untuk mengelompokkan card)
--   description_group = label per-tipe (nama tampilan tipe, sesuai codeMap legacy)
--
-- CATATAN: file diamend pasca-apply (checksum schema_migrations berbeda dari file).
-- Migrator skip by version tanpa validasi checksum; tenant yang sudah apply 053
-- diperbarui dengan menjalankan file ini langsung ke DB. Aman dijalankan ulang.
--
-- Idempotent secara alami (UPDATE — aman dijalankan ulang).

UPDATE job_management_values SET
    type_group = CASE type
        -- Pendidikan & Pengalaman
        WHEN 'education'   THEN 'education'
        WHEN 'experience'  THEN 'experience'
        -- Potensi Psikologi (migration 041)
        WHEN 'kecerdasan'            THEN 'psychological'
        WHEN 'innovation_creativity' THEN 'psychological'
        WHEN 'self_confidence'       THEN 'psychological'
        WHEN 'flexibility'           THEN 'psychological'
        WHEN 'tenacity'              THEN 'psychological'
        WHEN 'continuous_learning'   THEN 'psychological'
        -- Kompetensi Technical (migration 042)
        WHEN 'competency_based_human_resources_management' THEN 'technical'
        WHEN 'competency_development'                      THEN 'technical'
        WHEN 'people_development'                          THEN 'technical'
        WHEN 'career_management'                           THEN 'technical'
        WHEN 'hr_assessment'                               THEN 'technical'
        WHEN 'recruitement_selection'                      THEN 'technical'
        WHEN 'job_analysis_evaluation'                     THEN 'technical'
        WHEN 'organizational_development'                  THEN 'technical'
        WHEN 'human_resources_information_system'          THEN 'technical'
        WHEN 'workload_analysis'                           THEN 'technical'
        WHEN 'performance_apraisal'                        THEN 'technical'
        WHEN 'remuneration_manajemen'                      THEN 'technical'
        WHEN 'reward_punisment_management'                 THEN 'technical'
        WHEN 'health_safety_environment'                   THEN 'technical'
        WHEN 'hubungan_industrial'                         THEN 'technical'
        WHEN 'budgeting'                                   THEN 'technical'
        -- Kompetensi Managerial (migration 042)
        WHEN 'integrity'               THEN 'managerial'
        WHEN 'achievement_orientation' THEN 'managerial'
        WHEN 'building_partnership'    THEN 'managerial'
        WHEN 'planning_organizing'     THEN 'managerial'
        WHEN 'leadership'              THEN 'managerial'
        WHEN 'developing_others'       THEN 'managerial'
        -- Communication & Influencing Skill (migration 039)
        WHEN 'communicating_influencing_skill' THEN 'communication'
        -- Problem Solving (migration 040)
        WHEN 'thinking_environment' THEN 'problem_solving'
        WHEN 'thinking_chalenge'   THEN 'problem_solving'
        -- Financial (migration 036/037/050)
        WHEN 'cash'                  THEN 'financial'
        WHEN 'authority'             THEN 'financial'
        WHEN 'impact'                THEN 'financial'
        WHEN 'authority_unauthorized' THEN 'financial'
        WHEN 'impact_unauthorized'   THEN 'financial'
        -- Asset (migration 038/050)
        WHEN 'asset'            THEN 'asset'
        WHEN 'asset_authority'  THEN 'asset'
        -- Lainnya
        WHEN 'subordinate'  THEN 'subordinate'
        WHEN 'activity'     THEN 'activity'
        WHEN 'environment'  THEN 'environment'
        WHEN 'risk'         THEN 'risk'
        WHEN 'relationship' THEN 'relationship'
        WHEN 'frequency'    THEN 'frequency'
        ELSE type_group
    END,
    description_group = CASE type
        -- Pendidikan & Pengalaman
        WHEN 'education'   THEN 'Pendidikan'
        WHEN 'experience'  THEN 'Pengalaman Kerja'
        -- Potensi Psikologi (codeMap legacy)
        WHEN 'kecerdasan'            THEN 'Kecerdasan'
        WHEN 'innovation_creativity' THEN 'Innovation & Creativity'
        WHEN 'self_confidence'       THEN 'Self Confidence'
        WHEN 'flexibility'           THEN 'Flexibility'
        WHEN 'tenacity'              THEN 'Tenacity'
        WHEN 'continuous_learning'   THEN 'Continuous Learning'
        -- Kompetensi Technical (codeMap legacy)
        WHEN 'competency_based_human_resources_management' THEN 'Competency Based Human Resources Management'
        WHEN 'competency_development'                      THEN 'Competency Development'
        WHEN 'people_development'                          THEN 'People Development'
        WHEN 'career_management'                           THEN 'Career Management'
        WHEN 'hr_assessment'                               THEN 'HR Assessment'
        WHEN 'recruitement_selection'                      THEN 'Recruitement & Selection'
        WHEN 'job_analysis_evaluation'                     THEN 'Job Analysis & Evaluation'
        WHEN 'organizational_development'                  THEN 'Organizational Development'
        WHEN 'human_resources_information_system'          THEN 'Human Resources Information System'
        WHEN 'workload_analysis'                           THEN 'Workload Analysis'
        WHEN 'performance_apraisal'                        THEN 'Performance Apraisal'
        WHEN 'remuneration_manajemen'                      THEN 'Remuneration Manajemen'
        WHEN 'reward_punisment_management'                 THEN 'Reward & Punisment Management'
        WHEN 'health_safety_environment'                   THEN 'Health & Safety Environment'
        WHEN 'hubungan_industrial'                         THEN 'Hubungan Industrial'
        WHEN 'budgeting'                                   THEN 'Budgeting'
        -- Kompetensi Managerial (codeMap legacy)
        WHEN 'integrity'               THEN 'Integrity'
        WHEN 'achievement_orientation' THEN 'Achievement Orientation'
        WHEN 'building_partnership'    THEN 'Building Partnership'
        WHEN 'planning_organizing'     THEN 'Planning & Organizing'
        WHEN 'leadership'              THEN 'Leadership'
        WHEN 'developing_others'       THEN 'Developing Others'
        -- Communication & Influencing Skill (codeMap legacy)
        WHEN 'communicating_influencing_skill' THEN 'Communicating & Influencing Skill'
        -- Problem Solving (codeMap legacy)
        WHEN 'thinking_environment' THEN 'Thinking Environment'
        WHEN 'thinking_chalenge'   THEN 'Thinking Chalenge'
        -- Financial
        WHEN 'cash'                  THEN 'Jumlah Uang'
        WHEN 'authority'             THEN 'Wewenang'
        WHEN 'impact'                THEN 'Dampak pada Hasil Akhir (Memiliki Wewenang Keuangan)'
        WHEN 'authority_unauthorized' THEN 'Wewenang (Tidak Memiliki Wewenang Keuangan)'
        WHEN 'impact_unauthorized'   THEN 'Dampak pada Hasil Akhir (Tidak Memiliki Wewenang Keuangan)'
        -- Asset
        WHEN 'asset'           THEN 'Nilai Asset'
        WHEN 'asset_authority' THEN 'Wewenang Asset'
        -- Lainnya
        WHEN 'subordinate'  THEN 'Total Bawahan'
        WHEN 'activity'     THEN 'Aktifitas Fisik'
        WHEN 'environment'  THEN 'Lingkungan Kerja'
        WHEN 'risk'         THEN 'Resiko/Bahaya'
        WHEN 'relationship' THEN 'Lingkup Hubungan Kerja'
        WHEN 'frequency'    THEN 'Frekuensi Hubungan Kerja'
        ELSE description_group
    END
WHERE type IN (
    'education', 'experience',
    'kecerdasan', 'innovation_creativity', 'self_confidence', 'flexibility', 'tenacity', 'continuous_learning',
    'competency_based_human_resources_management', 'competency_development', 'people_development',
    'career_management', 'hr_assessment', 'recruitement_selection', 'job_analysis_evaluation',
    'organizational_development', 'human_resources_information_system', 'workload_analysis',
    'performance_apraisal', 'remuneration_manajemen', 'reward_punisment_management',
    'health_safety_environment', 'hubungan_industrial', 'budgeting',
    'integrity', 'achievement_orientation', 'building_partnership', 'planning_organizing',
    'leadership', 'developing_others',
    'communicating_influencing_skill',
    'thinking_environment', 'thinking_chalenge',
    'cash', 'authority', 'impact', 'authority_unauthorized', 'impact_unauthorized',
    'asset', 'asset_authority',
    'subordinate', 'activity', 'environment', 'risk', 'relationship', 'frequency'
);
