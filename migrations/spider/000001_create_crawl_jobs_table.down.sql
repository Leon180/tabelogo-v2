-- Drop crawl_jobs table
DROP TRIGGER IF EXISTS update_crawl_jobs_updated_at ON crawl_jobs;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_crawl_jobs_config;
DROP INDEX IF EXISTS idx_crawl_jobs_active;
DROP INDEX IF EXISTS idx_crawl_jobs_created_at;
DROP INDEX IF EXISTS idx_crawl_jobs_next_run;
DROP INDEX IF EXISTS idx_crawl_jobs_region;
DROP INDEX IF EXISTS idx_crawl_jobs_priority;
DROP INDEX IF EXISTS idx_crawl_jobs_source;
DROP INDEX IF EXISTS idx_crawl_jobs_status;
DROP TABLE IF EXISTS crawl_jobs;
