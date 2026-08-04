-- Migration: 051_job_management_scores_is_complete
-- Menambah kolom is_complete & completed_at pada job_management_scores
-- (menyamai perilaku legacy JobValueCalculator.persistResults — dokumen 8.4/8.7.2).

ALTER TABLE job_management_scores
    ADD COLUMN is_complete TINYINT(1) NOT NULL DEFAULT 0 AFTER calculated_at,
    ADD COLUMN completed_at TIMESTAMP NULL DEFAULT NULL AFTER is_complete;
