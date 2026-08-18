-- =============================================================================
-- Tenant Migration: 155_employee_id_format_settings (PostgreSQL)
-- =============================================================================
-- Pengaturan format generasi employee_id (nomor induk karyawan). Singleton
-- table (satu baris per tenant, dibuat otomatis oleh service saat pertama
-- kali dibaca — lihat pola attendance_company_settings).
--
-- generation_mode:
--   MANUAL  — employee_id wajib diisi manual saat create employee.
--   AUTO    — employee_id selalu di-generate otomatis, input klien diabaikan.
--   HYBRID  — pakai input klien jika diisi, generate otomatis jika kosong.
--
-- format_template mendukung placeholder {sequence}/{sequence:N}, {year},
-- {yy}, {month}, {month_roman} (lihat internal/pkg/numbering/format.go).
-- last_sequence & last_reset_key melacak nomor urut terakhir per periode
-- reset (reset_period).
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS employee_id_format_settings (
    id              CHAR(36) PRIMARY KEY,
    generation_mode VARCHAR(20) NOT NULL DEFAULT 'MANUAL',
    format_template VARCHAR(255) NOT NULL DEFAULT 'EMP-{year}-{sequence:4}',
    reset_period    VARCHAR(20) NOT NULL DEFAULT 'yearly',
    last_sequence   INT NOT NULL DEFAULT 0,
    last_reset_key  VARCHAR(16) NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
