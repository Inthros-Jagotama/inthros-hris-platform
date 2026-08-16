-- 127_business_travel_expense_master.down.sql
-- Rollback expense category master & expense plans

DROP TABLE IF EXISTS business_travel_expense_plans;
DROP TABLE IF EXISTS business_travel_expense_categories;
