USE psy_data;
-- Disable foreign key check
SET FOREIGN_KEY_CHECKS = 0;

-- Truncate all tables
-- Replace 'table_name' with the actual names of your tables
TRUNCATE surveys;
TRUNCATE questions;
TRUNCATE answers;
TRUNCATE patients;
TRUNCATE survey_results;
-- Repeat the TRUNCATE statement for each table you want to truncate

-- Enable foreign key check back
SET FOREIGN_KEY_CHECKS = 1;