-- Create crawl_results table
CREATE TABLE IF NOT EXISTS crawl_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES crawl_jobs(id) ON DELETE CASCADE,
    external_id VARCHAR(255) NOT NULL,  -- ID from external source
    source VARCHAR(50) NOT NULL,  -- 'tabelog', 'google_maps', 'instagram', etc.
    url TEXT,  -- Source URL
    raw_data JSONB NOT NULL,  -- Raw crawled data (reviews, ratings, photos, hours, etc.)
    parsed_data JSONB,  -- Parsed and normalized data
    checksum VARCHAR(64),  -- MD5/SHA256 hash for duplicate detection
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processed', 'failed', 'duplicate')),
    processed BOOLEAN DEFAULT FALSE,
    processing_error TEXT,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_crawl_results_job_id ON crawl_results(job_id);
CREATE INDEX idx_crawl_results_source ON crawl_results(source);
CREATE INDEX idx_crawl_results_status ON crawl_results(status);
CREATE INDEX idx_crawl_results_processed ON crawl_results(processed) WHERE processed = FALSE;
CREATE INDEX idx_crawl_results_checksum ON crawl_results(checksum);
CREATE INDEX idx_crawl_results_created_at ON crawl_results(created_at DESC);

-- Unique index to prevent duplicates from same source
CREATE UNIQUE INDEX idx_crawl_results_source_external_id
    ON crawl_results(source, external_id);

-- Create GIN indexes for JSONB columns
CREATE INDEX idx_crawl_results_raw_data ON crawl_results USING GIN(raw_data);
CREATE INDEX idx_crawl_results_parsed_data ON crawl_results USING GIN(parsed_data);

-- Create updated_at trigger
CREATE TRIGGER update_crawl_results_updated_at BEFORE UPDATE ON crawl_results
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE crawl_results IS 'Individual crawl results from external sources (Google Maps, Tabelog, Instagram)';
COMMENT ON COLUMN crawl_results.job_id IS 'Reference to crawl_jobs.id';
COMMENT ON COLUMN crawl_results.external_id IS 'ID from external source (e.g., Google Place ID, Tabelog restaurant ID)';
COMMENT ON COLUMN crawl_results.raw_data IS 'Raw crawled data in JSON format (reviews, ratings, photos, hours, menu, etc.)';
COMMENT ON COLUMN crawl_results.parsed_data IS 'Parsed and normalized data ready for processing';
COMMENT ON COLUMN crawl_results.checksum IS 'Hash for duplicate detection';
COMMENT ON COLUMN crawl_results.status IS 'Processing status: pending, processed, failed, duplicate';
