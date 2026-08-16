-- 129_business_travel_expenses.down.sql
-- Rollback actual expenses & expense documents

DROP TABLE IF EXISTS business_travel_expense_documents;
DROP TABLE IF EXISTS business_travel_expenses;
