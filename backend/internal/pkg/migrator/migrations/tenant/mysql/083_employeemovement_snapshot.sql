-- =============================================================================
-- Tenant Migration: 083_employeemovement_snapshot
-- =============================================================================
-- Employee Movement Enhancement (plan §12.5): Movement Snapshot.
--
-- Menyimpan nama Organization / Position / Employment Status pada saat
-- movement dibuat, sehingga histori movement tidak berubah ketika master
-- data (nomenclature organisasi, title posisi, nama status) diubah namanya.
--
-- Kolom baru (semua nullable — movement lama tetap valid tanpa snapshot):
--   from_organization_name
--   from_position_name
--   from_employment_status_name
--   to_organization_name
--   to_position_name
--   to_employment_status_name
--
-- Foreign key (from_*/to_*_id) tetap dipertahankan untuk relasi & navigasi.

SET @add_from_org_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'from_organization_name'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN from_organization_name VARCHAR(255) NULL'
);
PREPARE stmt FROM @add_from_org_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_from_pos_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'from_position_name'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN from_position_name VARCHAR(255) NULL'
);
PREPARE stmt FROM @add_from_pos_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_from_status_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'from_employment_status_name'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN from_employment_status_name VARCHAR(255) NULL'
);
PREPARE stmt FROM @add_from_status_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_to_org_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'to_organization_name'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN to_organization_name VARCHAR(255) NULL'
);
PREPARE stmt FROM @add_to_org_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_to_pos_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'to_position_name'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN to_position_name VARCHAR(255) NULL'
);
PREPARE stmt FROM @add_to_pos_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_to_status_name = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'to_employment_status_name'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN to_employment_status_name VARCHAR(255) NULL'
);
PREPARE stmt FROM @add_to_status_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
