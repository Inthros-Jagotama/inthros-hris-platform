-- =============================================================================
-- Tenant Migration: 158_career_path_step_eligibility_thresholds (MySQL)
-- =============================================================================
-- career_path_requirements (roadmap module-career-intelligence-plan.md §9 #6):
-- ambang batas skor eligibility promosi (performance/competency/OKR) yang
-- sebelumnya hardcoded global (>=80 untuk ketiganya, lihat
-- employeemovement/service.go eligibilityMin*Score) kini bisa dikonfigurasi
-- PER LANGKAH career path -- kolom nullable, NULL berarti pakai default
-- global. Ditambahkan ke career_path_steps (bukan tabel baru) konsisten
-- dengan skema terpadu 086/087 (Employee Movement + Career Intelligence
-- berbagi tabel yang sama).

ALTER TABLE career_path_steps
    ADD COLUMN min_performance_score DECIMAL(5,2) NULL AFTER requirements,
    ADD COLUMN min_competency_score  DECIMAL(5,2) NULL AFTER min_performance_score,
    ADD COLUMN min_okr_score         DECIMAL(5,2) NULL AFTER min_competency_score;
