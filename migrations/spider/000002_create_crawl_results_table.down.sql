-- Drop crawl_results table
DROP TRIGGER IF EXISTS update_crawl_results_updated_at ON crawl_results;
DROP INDEX IF EXISTS idx_crawl_results_parsed_data;
DROP INDEX IF EXISTS idx_crawl_results_raw_data;
DROP INDEX IF EXISTS idx_crawl_results_source_external_id;
DROP INDEX IF EXISTS idx_crawl_results_created_at;
DROP INDEX IF EXISTS idx_crawl_results_checksum;
DROP INDEX IF EXISTS idx_crawl_results_processed;
DROP INDEX IF EXISTS idx_crawl_results_status;
DROP INDEX IF EXISTS idx_crawl_results_source;
DROP INDEX IF EXISTS idx_crawl_results_job_id;
DROP TABLE IF EXISTS crawl_results;
