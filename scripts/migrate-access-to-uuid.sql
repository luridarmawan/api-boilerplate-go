-- Migration script to change access table ID from integer to UUID
-- This script should be run manually before deploying the new code

-- Step 1: Create a backup table
CREATE TABLE access_backup AS SELECT * FROM access;

-- ignore this script
