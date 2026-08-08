-- =============================================================================
-- Tenant Migration: 074_notification
-- =============================================================================
-- Notification Module Phase 1 (docs/module-notification-plan.md §3/§9):
--   No user-facing notification system exists anywhere in this codebase.
--   Other modules (Leave, Attendance, etc.) can only signal outcomes through
--   the Central Approval Module's internal callbacks, which nothing surfaces
--   to the recipient. This table is the storage for in-app notifications
--   keyed by platform user_id (not employee_id), matching the recipient
--   identity convention already used by the Approval module.

CREATE TABLE IF NOT EXISTS notifications (
    id                  CHAR(36) PRIMARY KEY,
    recipient_user_id   CHAR(36) NOT NULL,
    type                VARCHAR(50) NOT NULL,
    title               VARCHAR(255) NOT NULL,
    body                VARCHAR(1000) NOT NULL,
    reference_type      VARCHAR(50) NULL,
    reference_id        CHAR(36) NULL,
    is_read             BOOLEAN NOT NULL DEFAULT FALSE,
    read_at             DATETIME(6) NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_notification_recipient (recipient_user_id),
    INDEX idx_notification_recipient_unread (recipient_user_id, is_read),
    INDEX idx_notification_reference (reference_type, reference_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
