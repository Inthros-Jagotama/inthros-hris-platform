-- Migration: 050_seed_job_value_cash_impact_environment_risk_relationship_frequency_asset
-- Menambahkan 7 tipe yang belum ter-seed ke tabel job_management_values:
--   cash         (Jumlah Uang)                          — 5 level
--   impact       (Dampak pada Hasil Akhir — Memiliki Wewenang Keuangan) — 6 level
--   environment  (Lingkungan Kerja)                     — 5 level
--   risk         (Resiko/Bahaya)                        — 5 level
--   relationship (Lingkup Hubungan Kerja)               — 5 level
--   frequency    (Frekuensi Hubungan Kerja)             — 5 level
--   asset        (Nilai Asset)                          — 8 level
--
-- Semua level & deskripsi mengikuti docs/seeder/JobManagementValuesTableSeeder.php.
-- Tipe authority/authority_unauthorized/impact_unauthorized/asset_authority/activity/
-- subordinate/communicating_influencing_skill/thinking_*/psikologi/technical/managerial/
-- education/experience sudah di-seed di migration 033–042 & 049.
--
-- Idempotent GANDA: row di-INSERT hanya jika (a) id (UUID v4 acak) belum ada DAN
-- (b) belum ada row dengan type+level yang sama. Pengaman (b) mencegah duplikat
-- pada tenant yang sudah membuat row manual untuk tipe-tipe ini via UI
-- (tabel tidak punya UNIQUE (type, level)).
--
-- Menggunakan UUID v4 acak (bukan pola '-4000-8000-') agar konsisten dengan
-- migration 043_standardize_job_value_uuids.

-- ── tipe 'cash' (Jumlah Uang) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'f273b716-bc4c-43fa-b803-8512ceec046a', 'cash', 1, '0 - 500 Jt', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'f273b716-bc4c-43fa-b803-8512ceec046a')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'cash' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '91ade1a9-408c-41f1-ad77-a76282aee729', 'cash', 2, '500 Jt - 1 M', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '91ade1a9-408c-41f1-ad77-a76282aee729')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'cash' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '84dee31f-e8aa-4c36-b2e6-25fef50ac373', 'cash', 3, '1 M - 5 M', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '84dee31f-e8aa-4c36-b2e6-25fef50ac373')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'cash' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'c12e94eb-1b60-4417-9b24-cbc7ffc27fe5', 'cash', 4, '5 M - 10 M', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'c12e94eb-1b60-4417-9b24-cbc7ffc27fe5')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'cash' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '536f5fd3-7e47-4529-a468-5411a54417a3', 'cash', 5, '> 10 M', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '536f5fd3-7e47-4529-a468-5411a54417a3')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'cash' AND level = 5);

-- ── tipe 'impact' (Dampak pada Hasil Akhir — Memiliki Wewenang Keuangan) — 6 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '78347b19-9ad6-4068-808a-a5fdb608f9da', 'impact', 1, 'Penyediaan jasa insidentil untuk penggunaan lainnya.', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '78347b19-9ad6-4068-808a-a5fdb608f9da')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'impact' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'e4bc6bec-f77d-43c4-9342-f5172cfceb69', 'impact', 2, 'Penyediaan layanan dukungan informasi/pencatatan atau pengoperasian sederhana peralatan pendukung.', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'e4bc6bec-f77d-43c4-9342-f5172cfceb69')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'impact' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'd1abd8c9-7ba8-44af-ad5f-344f00bf4f49', 'impact', 3, 'Pengoperasian proses atau peralatan yang berhubungan langsung dengan rantai nilai inti bisnis.', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'd1abd8c9-7ba8-44af-ad5f-344f00bf4f49')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'impact' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '2eb9ed35-64e1-48b7-ae02-ed758ec9e5f5', 'impact', 4, 'Layanan analitis, diagnostik, konsultasi, atau pengoperasian sistem penting yang sangat kompleks.', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '2eb9ed35-64e1-48b7-ae02-ed758ec9e5f5')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'impact' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '4a40df12-0a0b-4ccf-b4d8-4ad3f1bbd397', 'impact', 5, 'Memimpin area aktivitas/tim dalam parameter yang jelas atau memberi nasihat pada tingkat kebijakan.', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '4a40df12-0a0b-4ccf-b4d8-4ad3f1bbd397')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'impact' AND level = 5);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '904499a7-f9c9-43f6-b0c4-aaff6210f212', 'impact', 6, 'Menghasilkan operasi multi tim atau memastikan penyampaian program strategis serta kebijakan fungsional.', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '904499a7-f9c9-43f6-b0c4-aaff6210f212')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'impact' AND level = 6);

-- ── tipe 'environment' (Lingkungan Kerja) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '4116ebac-285d-4e74-8f88-bcec7c0ca4cd', 'environment', 1, 'Tenang, nyaman', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '4116ebac-285d-4e74-8f88-bcec7c0ca4cd')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'environment' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '73f8d76c-baf8-4952-a88b-8679765e8781', 'environment', 2, 'Cukup bising dan sibuk', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '73f8d76c-baf8-4952-a88b-8679765e8781')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'environment' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '889df711-cd1e-4cda-877d-372c2db4b95a', 'environment', 3, 'Bising, sibuk dan banyak gangguan', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '889df711-cd1e-4cda-877d-372c2db4b95a')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'environment' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'd5a09942-617b-45a2-91dc-05bc405ae19a', 'environment', 4, 'Banyak tantangan dan tekanan kerja', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'd5a09942-617b-45a2-91dc-05bc405ae19a')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'environment' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'd36d03f8-de8e-4178-aea9-83a028d7a4ad', 'environment', 5, 'Sangat menekan dan menegangkan', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'd36d03f8-de8e-4178-aea9-83a028d7a4ad')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'environment' AND level = 5);

-- ── tipe 'risk' (Resiko/Bahaya) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '4ee7a915-d12e-487a-af9c-bf855955ec29', 'risk', 1, 'Risiko minimum dan bebas bahaya', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '4ee7a915-d12e-487a-af9c-bf855955ec29')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'risk' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '4440f46a-0c9a-4d42-88fa-a5e066db608e', 'risk', 2, 'Risiko kecil dengan sedikit ancaman bahaya', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '4440f46a-0c9a-4d42-88fa-a5e066db608e')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'risk' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'c6aed4f7-0857-46f3-b2e0-0e0b81800636', 'risk', 3, 'Risiko dan bahaya yang cukup besar', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'c6aed4f7-0857-46f3-b2e0-0e0b81800636')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'risk' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '8a4f3fbe-606d-4dfe-b18a-ea0b7be747fe', 'risk', 4, 'Risiko dan bahaya tinggi', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '8a4f3fbe-606d-4dfe-b18a-ea0b7be747fe')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'risk' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '06f5eacb-1146-43e8-8cf0-2e88736b8023', 'risk', 5, 'Ancaman bahaya besar yang mematikan', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '06f5eacb-1146-43e8-8cf0-2e88736b8023')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'risk' AND level = 5);

-- ── tipe 'relationship' (Lingkup Hubungan Kerja) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'f82c3392-086e-4d2f-84d5-6bf9707aa61a', 'relationship', 1, 'Unit Kerja', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'f82c3392-086e-4d2f-84d5-6bf9707aa61a')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'relationship' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '372d6da2-de55-4ffa-bc71-3931d6d0630b', 'relationship', 2, 'Antar Unit Kerja', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '372d6da2-de55-4ffa-bc71-3931d6d0630b')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'relationship' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'cb650d11-43c1-4bd4-b19c-3888d4df4b69', 'relationship', 3, 'Lingkup Internal secara Nasional', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cb650d11-43c1-4bd4-b19c-3888d4df4b69')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'relationship' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '59f3ab76-db58-4b41-97a1-6b62042306dd', 'relationship', 4, 'Lingkup Eksternal secara Nasional', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '59f3ab76-db58-4b41-97a1-6b62042306dd')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'relationship' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '330e38ce-341f-46af-820b-8fd2e0bee4da', 'relationship', 5, 'Lingkup Internasional', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '330e38ce-341f-46af-820b-8fd2e0bee4da')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'relationship' AND level = 5);

-- ── tipe 'frequency' (Frekuensi Hubungan Kerja) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '5c7476eb-c662-43ba-b37a-305d63c7cb21', 'frequency', 1, 'Sesekali', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '5c7476eb-c662-43ba-b37a-305d63c7cb21')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'frequency' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '1aa3cd37-6918-48de-975d-9ec5676e0db7', 'frequency', 2, 'Kadang - kadang', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '1aa3cd37-6918-48de-975d-9ec5676e0db7')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'frequency' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '3a387f2c-0ec1-4369-b47e-5bc54254861e', 'frequency', 3, 'Cukup Sering', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '3a387f2c-0ec1-4369-b47e-5bc54254861e')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'frequency' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'd3f8deba-771e-4258-b018-468654b27969', 'frequency', 4, 'Sering', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'd3f8deba-771e-4258-b018-468654b27969')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'frequency' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '2e777fa5-7949-4d63-a630-e743b10cec3a', 'frequency', 5, 'Sangat Sering', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '2e777fa5-7949-4d63-a630-e743b10cec3a')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'frequency' AND level = 5);

-- ── tipe 'asset' (Nilai Asset) — 8 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'a73fa37d-bf83-458c-a8dc-d80f4dd0cdbe', 'asset', 1, '0 - 1 Jt', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'a73fa37d-bf83-458c-a8dc-d80f4dd0cdbe')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'a2f3ec3f-1722-4271-9e76-6b856cac274f', 'asset', 2, '1 - 10 Jt', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'a2f3ec3f-1722-4271-9e76-6b856cac274f')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '643f0773-0d29-4fcd-8d08-536dab90c8cc', 'asset', 3, '10 - 50 Jt', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '643f0773-0d29-4fcd-8d08-536dab90c8cc')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '46a04f2e-2e3e-4757-b81a-9aaa1f99fb6f', 'asset', 4, '50 - 100 Jt', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '46a04f2e-2e3e-4757-b81a-9aaa1f99fb6f')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '303def52-0977-4fc3-82ba-3a8d9c18cc19', 'asset', 5, '100 - 250 Jt', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '303def52-0977-4fc3-82ba-3a8d9c18cc19')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 5);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '91b7b107-b7b0-4d1f-b50e-5863246084ae', 'asset', 6, '250 - 500 Jt', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '91b7b107-b7b0-4d1f-b50e-5863246084ae')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 6);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '34785d5b-4c61-4efe-a71f-4d2b3e6a1b0d', 'asset', 7, '500 Jt - 1 M', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '34785d5b-4c61-4efe-a71f-4d2b3e6a1b0d')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 7);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '182b8d48-8b33-4629-b56c-7b58a6eaf4ca', 'asset', 8, '> 1 M', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '182b8d48-8b33-4629-b56c-7b58a6eaf4ca')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'asset' AND level = 8);
