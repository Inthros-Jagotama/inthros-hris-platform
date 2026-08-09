-- =============================================================================
-- Tenant Migration (down): 082_employeemovement_status_rejected
-- =============================================================================
-- Rollback: tidak ada perubahan skema yang harus di-revert — kolom `status`
-- tetap VARCHAR(20) (nilai 'rejected' hanya not used lagi oleh aplikasi).
-- Jika migration up sempat mengubah kolom menjadi ENUM (kasus drift), kembalikan
-- ke daftar ENUM tanpa 'rejected'.
--
-- ⚠️ PERINGATAN: jika ada baris employee_movements yang berstatus 'rejected'
-- (data yang ditulis setelah up), konversi ke ENUM tanpa 'rejected' akan gagal
-- atau meng-truncate data. Konversi di-skip otomatis bila baris tsb ada —
-- operator harus menormalkan/hapus baris tersebut dulu sebelum rollback.

SET @is_varchar = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'employee_movements'
    AND COLUMN_NAME = 'status'
    AND DATA_TYPE = 'varchar'
);

-- Baris berstatus 'rejected' yang masih ada (mencegah data truncation).
SET @has_rejected_rows = (
  SELECT COUNT(*)
  FROM employee_movements
  WHERE status = 'rejected'
);

SET @revert_status = IF(
  @is_varchar = 0 AND @has_rejected_rows = 0,
  'ALTER TABLE employee_movements MODIFY COLUMN status ENUM(''draft'',''pending_approval'',''approved'',''executed'',''cancelled'') DEFAULT ''draft''',
  'DO 0'
);
PREPARE stmt FROM @revert_status;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
