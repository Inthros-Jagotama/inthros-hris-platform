-- 064_seed_performance_components.down.sql

DELETE FROM performance_components WHERE code IN ('KPI', 'PROGRAM', 'SUBORDINATE');
