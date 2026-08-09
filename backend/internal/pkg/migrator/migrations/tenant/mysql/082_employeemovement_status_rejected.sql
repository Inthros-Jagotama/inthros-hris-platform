-- =============================================================================
-- Tenant Migration: 082_employeemovement_status_rejected
-- =============================================================================
-- Employee Movement: menambahkan nilai `rejected` ke status pergerakan.
--
-- Kolom `status` di employee_movements sudah VARCHAR(20) sejak migration 012
-- (bukan ENUM), sehingga nilai 'rejected' (8 karakter) tertampung tanpa
-- perubahan skema.
--
-- Migration ini bersifat dokumentasi + verifikasi: memastikan kolom tetap
-- VARCHAR (bukan ENUM) yang mengizinkan nilai 'rejected'. Jika suatu saat
-- kolom diubah menjadi ENUM, migration ini akan menghasilkan pesan peringatan
-- agar operator menambah nilai 'rejected' ke daftar ENUM.

SET @is_varchar = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'employee_movements'
    AND COLUMN_NAME = 'status'
    AND DATA_TYPE = 'varchar'
);

-- Jika bukan VARCHAR (misal sudah ENUM), tambahkan 'rejected' ke daftar nilai.
SET @alter_status = IF(
  @is_varchar = 0,
  CONCAT(
    'ALTER TABLE employee_movements MODIFY COLUMN status ',
    'ENUM(''draft'',''pending_approval'',''approved'',''rejected'',''executed'',''cancelled'') ',
    'DEFAULT ''draft'''
  ),
  'DO 0'
);
PREPARE stmt FROM @alter_status;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
