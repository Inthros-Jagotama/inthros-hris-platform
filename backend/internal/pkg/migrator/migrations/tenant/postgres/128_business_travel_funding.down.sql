-- 128_business_travel_funding.down.sql
-- Rollback funding method master, fundings & funding documents

DROP TABLE IF EXISTS business_travel_funding_documents;
DROP TABLE IF EXISTS business_travel_fundings;
DROP TABLE IF EXISTS business_travel_funding_methods;
