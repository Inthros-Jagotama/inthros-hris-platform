-- Down Migration: 049_seed_job_value_education_experience
-- Menghapus HANYA 10 row tipe 'education' & 'experience' yang di-seed (by UUID acak)
-- agar row buatan user (lewat UI) tidak ikut terhapus saat rollback.
DELETE FROM job_management_values
WHERE id IN (
  '129d8602-5eda-4095-b403-2bdcafa80ae0',
  'df9019c2-e110-44c3-8815-82fdef3c700a',
  '8f2ef00c-4f6a-4904-88ed-bcdf4917476b',
  'c1791488-e91b-45cc-99a8-30ab0cb80951',
  'ed02d3d8-8eb8-4c38-a026-931e4a3998ec',
  '4e3e2ed3-131d-4ebc-a198-03e77b097cd1',
  '4ffa3804-d566-4986-9a53-5282e85506da',
  'dab6918a-1c0e-48b7-bfe7-e28eff2dca1f',
  '62cbaf52-e99a-41fa-8cc2-7e61eb3c4969',
  'd89056a0-39da-45a2-a643-04598c630888'
);
