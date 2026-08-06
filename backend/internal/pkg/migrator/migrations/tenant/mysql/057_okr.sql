-- OKR Module Tables
-- Phase 2: New Tables for OKR (Objective & Key Results)

-- 1. okr_templates - Template OKR untuk setiap Organization
CREATE TABLE IF NOT EXISTS okr_templates (
    id CHAR(36) PRIMARY KEY,
    organization_id CHAR(36) NOT NULL,
    period_id CHAR(36),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status SMALLINT NOT NULL DEFAULT 0,
    effective_date DATE,
    expired_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_okr_tpl_org (organization_id),
    INDEX idx_okr_tpl_period (period_id),
    INDEX idx_okr_tpl_status (status),
    CONSTRAINT fk_okr_tpl_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_okr_tpl_period FOREIGN KEY (period_id) REFERENCES performance_periods(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. okr_objectives - Objective yang ingin dicapai
CREATE TABLE IF NOT EXISTS okr_objectives (
    id CHAR(36) PRIMARY KEY,
    template_id CHAR(36) NOT NULL,
    code VARCHAR(50),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_okr_obj_template (template_id),
    INDEX idx_okr_obj_sort (sort_order),
    CONSTRAINT fk_okr_obj_template FOREIGN KEY (template_id) REFERENCES okr_templates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. okr_key_results - Target terukur dari setiap Objective
CREATE TABLE IF NOT EXISTS okr_key_results (
    id CHAR(36) PRIMARY KEY,
    objective_id CHAR(36) NOT NULL,
    code VARCHAR(50),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    target_type VARCHAR(30) NOT NULL DEFAULT 'NUMBER',
    target_value DECIMAL(18,2) NOT NULL DEFAULT 0,
    unit VARCHAR(50),
    formula_type VARCHAR(30) NOT NULL DEFAULT 'HIGHER_BETTER',
    weight DECIMAL(5,2) NOT NULL DEFAULT 0,
    minimum_score DECIMAL(5,2) NOT NULL DEFAULT 0,
    maximum_score DECIMAL(5,2) NOT NULL DEFAULT 100,
    sort_order INT NOT NULL DEFAULT 0,
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_okr_kr_objective (objective_id),
    INDEX idx_okr_kr_sort (sort_order),
    CONSTRAINT fk_okr_kr_objective FOREIGN KEY (objective_id) REFERENCES okr_objectives(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. okr_evaluations - Header penilaian OKR
CREATE TABLE IF NOT EXISTS okr_evaluations (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36) NOT NULL,
    organization_id CHAR(36) NOT NULL,
    period_id CHAR(36) NOT NULL,
    template_id CHAR(36),
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    submitted_at TIMESTAMP NULL,
    submitted_by CHAR(36),
    approved_at TIMESTAMP NULL,
    approved_by CHAR(36),
    final_score DECIMAL(5,2) NOT NULL DEFAULT 0,
    rating_id CHAR(36),
    reviewer_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_okr_eval_employee (employee_id),
    INDEX idx_okr_eval_org (organization_id),
    INDEX idx_okr_eval_period (period_id),
    INDEX idx_okr_eval_status (status),
    CONSTRAINT fk_okr_eval_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_okr_eval_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_okr_eval_period FOREIGN KEY (period_id) REFERENCES performance_periods(id) ON DELETE CASCADE,
    CONSTRAINT fk_okr_eval_template FOREIGN KEY (template_id) REFERENCES okr_templates(id) ON DELETE SET NULL,
    CONSTRAINT fk_okr_eval_rating FOREIGN KEY (rating_id) REFERENCES performance_ratings(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. okr_evaluation_details - Snapshot Objective dan Key Result saat evaluasi
CREATE TABLE IF NOT EXISTS okr_evaluation_details (
    id CHAR(36) PRIMARY KEY,
    evaluation_id CHAR(36) NOT NULL,
    objective_id CHAR(36),
    key_result_id CHAR(36),
    objective_title VARCHAR(255) NOT NULL,
    key_result_title VARCHAR(255) NOT NULL,
    objective_weight DECIMAL(5,2) NOT NULL DEFAULT 0,
    key_result_weight DECIMAL(5,2) NOT NULL DEFAULT 0,
    target_value DECIMAL(18,2) NOT NULL DEFAULT 0,
    target_type VARCHAR(30) NOT NULL DEFAULT 'NUMBER',
    unit VARCHAR(50),
    formula_type VARCHAR(30) NOT NULL DEFAULT 'HIGHER_BETTER',
    actual_value DECIMAL(18,2) NOT NULL DEFAULT 0,
    achievement DECIMAL(5,2) NOT NULL DEFAULT 0,
    score DECIMAL(5,2) NOT NULL DEFAULT 0,
    remarks TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_okr_detail_eval (evaluation_id),
    INDEX idx_okr_detail_obj (objective_id),
    INDEX idx_okr_detail_kr (key_result_id),
    CONSTRAINT fk_okr_detail_eval FOREIGN KEY (evaluation_id) REFERENCES okr_evaluations(id) ON DELETE CASCADE,
    CONSTRAINT fk_okr_detail_obj FOREIGN KEY (objective_id) REFERENCES okr_objectives(id) ON DELETE SET NULL,
    CONSTRAINT fk_okr_detail_kr FOREIGN KEY (key_result_id) REFERENCES okr_key_results(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. okr_progress - Progress Check-in
CREATE TABLE IF NOT EXISTS okr_progress (
    id CHAR(36) PRIMARY KEY,
    evaluation_detail_id CHAR(36) NOT NULL,
    progress_date DATE NOT NULL,
    actual_value DECIMAL(18,2) NOT NULL DEFAULT 0,
    achievement DECIMAL(5,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_okr_prog_detail (evaluation_detail_id),
    INDEX idx_okr_prog_date (progress_date),
    CONSTRAINT fk_okr_prog_detail FOREIGN KEY (evaluation_detail_id) REFERENCES okr_evaluation_details(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. okr_comments - Komentar antara Employee dan Reviewer
CREATE TABLE IF NOT EXISTS okr_comments (
    id CHAR(36) PRIMARY KEY,
    evaluation_id CHAR(36) NOT NULL,
    parent_id CHAR(36),
    comment TEXT NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_okr_comment_eval (evaluation_id),
    INDEX idx_okr_comment_parent (parent_id),
    INDEX idx_okr_comment_user (created_by),
    CONSTRAINT fk_okr_comment_eval FOREIGN KEY (evaluation_id) REFERENCES okr_evaluations(id) ON DELETE CASCADE,
    CONSTRAINT fk_okr_comment_parent FOREIGN KEY (parent_id) REFERENCES okr_comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 8. okr_attachments - Evidence pencapaian Key Result
CREATE TABLE IF NOT EXISTS okr_attachments (
    id CHAR(36) PRIMARY KEY,
    evaluation_detail_id CHAR(36) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(100),
    file_size BIGINT,
    description TEXT,
    uploaded_by CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_okr_attach_detail (evaluation_detail_id),
    INDEX idx_okr_attach_user (uploaded_by),
    CONSTRAINT fk_okr_attach_detail FOREIGN KEY (evaluation_detail_id) REFERENCES okr_evaluation_details(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
