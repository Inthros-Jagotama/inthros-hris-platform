-- 144_competency_assessment_templates.down.sql

-- Hapus kolom template_id dari competency_events (144.4) sebelum drop tabel.
ALTER TABLE competency_events
    DROP COLUMN template_id;

DROP TABLE IF EXISTS competency_assessment_template_rater_types;

DROP TABLE IF EXISTS competency_assessment_template_competencies;

DROP TABLE IF EXISTS competency_assessment_templates;
