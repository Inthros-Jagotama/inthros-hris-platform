-- ========================================================
-- Module: Key Performance Management (KPI / BSC)
-- Engine: InnoDB
-- Charset: utf8mb4
-- Description: Redesigned KPI database schema with standardized
--              English naming, Foreign Key constraints, and
--              optimized index strategy.
-- ========================================================

SET FOREIGN_KEY_CHECKS = 0;

-- --------------------------------------------------------
-- 1. MASTER & REFERENCE TABLES
-- --------------------------------------------------------

-- Master Perspectives (Financial, Customer, Internal Process, Learning & Growth, etc.)
CREATE TABLE IF NOT EXISTS `perspectives` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL,
  `description` TEXT DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Master BSC Categories / Types
CREATE TABLE IF NOT EXISTS `bsc_categories` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(150) NOT NULL,
  `is_perspective` TINYINT(1) NOT NULL DEFAULT 0,
  `is_work_program` TINYINT(1) NOT NULL DEFAULT 0,
  `is_subordinate_performance` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Evaluation Periods (e.g., Year 2026, Q1, Semester 1)
CREATE TABLE IF NOT EXISTS `evaluation_periods` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `period_code` VARCHAR(10) NOT NULL, -- e.g., 'Q1', 'H1', 'FY'
  `period_type` VARCHAR(20) NOT NULL, -- e.g., 'MONTHLY', 'QUARTERLY', 'ANNUAL'
  `year` YEAR NOT NULL,
  `start_date` DATE DEFAULT NULL,
  `end_date` DATE DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indicator Types (Positive vs Negative KPI Target Calculation)
CREATE TABLE IF NOT EXISTS `indicator_types` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) NOT NULL, -- e.g., 'MAXIMIZATION' / 'MINIMIZATION'
  `is_positive` TINYINT(1) NOT NULL DEFAULT 1, -- 1 if higher is better, 0 if lower is better
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Evaluation Statuses (Draft, Submitted, Approved, Rejected)
CREATE TABLE IF NOT EXISTS `ref_evaluation_statuses` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `status_key` VARCHAR(50) NOT NULL UNIQUE, -- e.g., 'DRAFT', 'SUBMITTED', 'APPROVED'
  `label` VARCHAR(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------
-- 2. TEMPLATE & CONFIGURATION TABLES
-- --------------------------------------------------------

-- Template BSC Header per Position
CREATE TABLE IF NOT EXISTS `bsc_templates` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `position_id` INT UNSIGNED NOT NULL,
  `ref_bsc_category_id` INT UNSIGNED NOT NULL,
  `weight` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
  `status` ENUM('DRAFT', 'PUBLISHED', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT',
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_bsc_templates_position` FOREIGN KEY (`position_id`) REFERENCES `util_position` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_bsc_templates_category` FOREIGN KEY (`ref_bsc_category_id`) REFERENCES `ref_bsc_categories` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Master KPI / Objectives linked to BSC Templates
CREATE TABLE IF NOT EXISTS `bsc_indicators` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `bsc_template_id` INT UNSIGNED NOT NULL,
  `perspective_id` INT UNSIGNED NOT NULL,
  `indicator_type_id` INT UNSIGNED NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT DEFAULT NULL,
  `weight` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_bsc_ind_template` FOREIGN KEY (`bsc_template_id`) REFERENCES `bsc_templates` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_bsc_ind_perspective` FOREIGN KEY (`perspective_id`) REFERENCES `ref_perspectives` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_bsc_ind_type` FOREIGN KEY (`indicator_type_id`) REFERENCES `ref_indicator_types` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------
-- 3. TRANSACTION / EVALUATION TABLES
-- --------------------------------------------------------

-- Performance Evaluation Header (Penilaian)
CREATE TABLE IF NOT EXISTS `kpi_evaluations` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `employee_id` INT UNSIGNED NOT NULL,
  `position_id` INT UNSIGNED NOT NULL,
  `evaluation_period_id` INT UNSIGNED NOT NULL,
  `supervisor_id` INT UNSIGNED DEFAULT NULL,
  `supervisor_position_id` INT UNSIGNED DEFAULT NULL,
  `final_score` DECIMAL(5,2) DEFAULT 0.00,
  `status_id` INT UNSIGNED NOT NULL,
  `plan_submitted_at` DATETIME DEFAULT NULL,
  `plan_approved_at` DATETIME DEFAULT NULL,
  `actual_submitted_at` DATETIME DEFAULT NULL,
  `actual_approved_at` DATETIME DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_eval_employee` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_eval_position` FOREIGN KEY (`position_id`) REFERENCES `util_position` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_eval_period` FOREIGN KEY (`evaluation_period_id`) REFERENCES `ref_evaluation_periods` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_eval_supervisor` FOREIGN KEY (`supervisor_id`) REFERENCES `employee` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_eval_status` FOREIGN KEY (`status_id`) REFERENCES `ref_evaluation_statuses` (`id`) ON DELETE RESTRICT,
  INDEX `idx_eval_employee_period` (`employee_id`, `evaluation_period_id`),
  INDEX `idx_eval_supervisor` (`supervisor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Evaluation BSC Group Detail (Penilaian BSC)
CREATE TABLE IF NOT EXISTS `evaluation_bsc_details` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `kpi_evaluation_id` INT UNSIGNED NOT NULL,
  `bsc_template_id` INT UNSIGNED DEFAULT NULL,
  `achievement_percentage` DECIMAL(5,2) DEFAULT 0.00,
  `weight` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
  `score` DECIMAL(5,2) DEFAULT 0.00,
  `description` VARCHAR(255) DEFAULT NULL,
  `is_perspective` TINYINT(1) NOT NULL DEFAULT 0,
  `is_work_program` TINYINT(1) NOT NULL DEFAULT 0,
  `is_subordinate_performance` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_eval_bsc_main` FOREIGN KEY (`kpi_evaluation_id`) REFERENCES `kpi_evaluations` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_eval_bsc_template` FOREIGN KEY (`bsc_template_id`) REFERENCES `bsc_templates` (`id`) ON DELETE SET NULL,
  INDEX `idx_eval_bsc_main` (`kpi_evaluation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Evaluation Target Items (Target / KPI Details)
CREATE TABLE IF NOT EXISTS `evaluation_targets` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `evaluation_bsc_detail_id` INT UNSIGNED NOT NULL,
  `bsc_indicator_id` INT UNSIGNED DEFAULT NULL, -- FK to master indicator template
  `indicator_type_id` INT UNSIGNED NOT NULL,
  `target_value` DECIMAL(12,2) DEFAULT 0.00,
  `actual_value` DECIMAL(12,2) DEFAULT 0.00,
  `unit_of_measurement` VARCHAR(50) DEFAULT NULL, -- e.g., '%', 'IDR', 'Count'
  `achievement_percentage` DECIMAL(5,2) DEFAULT 0.00,
  `weight` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
  `score` DECIMAL(5,2) DEFAULT 0.00,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_target_bsc_detail` FOREIGN KEY (`evaluation_bsc_detail_id`) REFERENCES `evaluation_bsc_details` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_target_indicator` FOREIGN KEY (`bsc_indicator_id`) REFERENCES `bsc_indicators` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_target_ind_type` FOREIGN KEY (`indicator_type_id`) REFERENCES `ref_indicator_types` (`id`) ON DELETE RESTRICT,
  INDEX `idx_target_bsc_detail` (`evaluation_bsc_detail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Evaluation Work Programs (Program Kerja)
CREATE TABLE IF NOT EXISTS `evaluation_work_programs` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `evaluation_bsc_detail_id` INT UNSIGNED NOT NULL,
  `indicator_type_id` INT UNSIGNED DEFAULT NULL,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT DEFAULT NULL,
  `target_value` DECIMAL(12,2) DEFAULT 0.00,
  `actual_value` DECIMAL(12,2) DEFAULT 0.00,
  `unit_of_measurement` VARCHAR(50) DEFAULT NULL,
  `achievement_percentage` DECIMAL(5,2) DEFAULT 0.00,
  `weight` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
  `score` DECIMAL(5,2) DEFAULT 0.00,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_work_prog_bsc_detail` FOREIGN KEY (`evaluation_bsc_detail_id`) REFERENCES `evaluation_bsc_details` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_work_prog_ind_type` FOREIGN KEY (`indicator_type_id`) REFERENCES `ref_indicator_types` (`id`) ON DELETE SET NULL,
  INDEX `idx_work_prog_bsc_detail` (`evaluation_bsc_detail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
