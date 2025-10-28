-- EstateSaleFinder.ai Database Schema
-- This script initializes the local development database
-- Auto-run when PostgreSQL container first starts

-- Run the main migration file
\i /docker-entrypoint-initdb.d/001_initial_schema.sql
