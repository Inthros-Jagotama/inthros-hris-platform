-- Down Migration: 050_seed_job_value_cash_impact_environment_risk_relationship_frequency_asset
-- Menghapus HANYA 39 row (7 tipe) yang di-seed (by UUID v4 acak) agar row buatan
-- user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  -- cash (5)
  'f273b716-bc4c-43fa-b803-8512ceec046a',
  '91ade1a9-408c-41f1-ad77-a76282aee729',
  '84dee31f-e8aa-4c36-b2e6-25fef50ac373',
  'c12e94eb-1b60-4417-9b24-cbc7ffc27fe5',
  '536f5fd3-7e47-4529-a468-5411a54417a3',
  -- impact (6)
  '78347b19-9ad6-4068-808a-a5fdb608f9da',
  'e4bc6bec-f77d-43c4-9342-f5172cfceb69',
  'd1abd8c9-7ba8-44af-ad5f-344f00bf4f49',
  '2eb9ed35-64e1-48b7-ae02-ed758ec9e5f5',
  '4a40df12-0a0b-4ccf-b4d8-4ad3f1bbd397',
  '904499a7-f9c9-43f6-b0c4-aaff6210f212',
  -- environment (5)
  '4116ebac-285d-4e74-8f88-bcec7c0ca4cd',
  '73f8d76c-baf8-4952-a88b-8679765e8781',
  '889df711-cd1e-4cda-877d-372c2db4b95a',
  'd5a09942-617b-45a2-91dc-05bc405ae19a',
  'd36d03f8-de8e-4178-aea9-83a028d7a4ad',
  -- risk (5)
  '4ee7a915-d12e-487a-af9c-bf855955ec29',
  '4440f46a-0c9a-4d42-88fa-a5e066db608e',
  'c6aed4f7-0857-46f3-b2e0-0e0b81800636',
  '8a4f3fbe-606d-4dfe-b18a-ea0b7be747fe',
  '06f5eacb-1146-43e8-8cf0-2e88736b8023',
  -- relationship (5)
  'f82c3392-086e-4d2f-84d5-6bf9707aa61a',
  '372d6da2-de55-4ffa-bc71-3931d6d0630b',
  'cb650d11-43c1-4bd4-b19c-3888d4df4b69',
  '59f3ab76-db58-4b41-97a1-6b62042306dd',
  '330e38ce-341f-46af-820b-8fd2e0bee4da',
  -- frequency (5)
  '5c7476eb-c662-43ba-b37a-305d63c7cb21',
  '1aa3cd37-6918-48de-975d-9ec5676e0db7',
  '3a387f2c-0ec1-4369-b47e-5bc54254861e',
  'd3f8deba-771e-4258-b018-468654b27969',
  '2e777fa5-7949-4d63-a630-e743b10cec3a',
  -- asset (8)
  'a73fa37d-bf83-458c-a8dc-d80f4dd0cdbe',
  'a2f3ec3f-1722-4271-9e76-6b856cac274f',
  '643f0773-0d29-4fcd-8d08-536dab90c8cc',
  '46a04f2e-2e3e-4757-b81a-9aaa1f99fb6f',
  '303def52-0977-4fc3-82ba-3a8d9c18cc19',
  '91b7b107-b7b0-4d1f-b50e-5863246084ae',
  '34785d5b-4c61-4efe-a71f-4d2b3e6a1b0d',
  '182b8d48-8b33-4629-b56c-7b58a6eaf4ca'
);
