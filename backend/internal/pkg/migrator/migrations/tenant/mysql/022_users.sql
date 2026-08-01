-- =============================================================================
-- Tenant Migration: 022_users
-- =============================================================================
-- Tabel users (Level 2 Tenant RBAC — identitas user internal tenant).
-- Melengkapi skema RBAC tenant (011_settings) yang sudah memiliki roles,
-- permissions, role_has_permissions, model_has_roles, model_has_permissions.
-- model_id pada model_has_roles/model_has_permissions menunjuk ke users.id.

CREATE TABLE IF NOT EXISTS users (
    id            CHAR(36) PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active     SMALLINT NOT NULL DEFAULT 1,
    last_login_at TIMESTAMP NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP NULL,

    UNIQUE KEY uk_users_email (email),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
