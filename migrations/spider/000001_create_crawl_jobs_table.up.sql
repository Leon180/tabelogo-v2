-- Create crawl_jobs table
CREATE TABLE IF NOT EXISTS crawl_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source VARCHAR(50) NOT NULL,  -- 'tabelog', 'google_maps', 'instagram', etc.
    region VARCHAR(100),  -- Target region/area for crawling
    job_type VARCHAR(50) DEFAULT 'full',  -- 'full', 'incremental', 'update'
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled', 'paused')),
    priority INT DEFAULT 5 CHECK (priority >= 1 AND priority <= 10),  -- 1 is highest priority
    total_pages INT,
    completed_pages INT DEFAULT 0,
    total_items INT,
    completed_items INT DEFAULT 0,
    success_count INT DEFAULT 0,
    error_count INT DEFAULT 0,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    config JSONB,  -- Crawl configuration (rate limit, user agents, proxy settings, etc.)
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    next_run_at TIMESTAMP,  -- For scheduled jobs
    error_message TEXT,
    error_details JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_crawl_jobs_status ON crawl_jobs(status);
CREATE INDEX idx_crawl_jobs_source ON crawl_jobs(source);
CREATE INDEX idx_crawl_jobs_priority ON crawl_jobs(priority, created_at);
CREATE INDEX idx_crawl_jobs_region ON crawl_jobs(region);
CREATE INDEX idx_crawl_jobs_next_run ON crawl_jobs(next_run_at) WHERE status = 'pending';
CREATE INDEX idx_crawl_jobs_created_at ON crawl_jobs(created_at DESC);

-- Composite index for active jobs
CREATE INDEX idx_crawl_jobs_active ON crawl_jobs(status, priority, created_at)
    WHERE status IN ('pending', 'running');

-- Create GIN index for JSONB columns
CREATE INDEX idx_crawl_jobs_config ON crawl_jobs USING GIN(config);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create updated_at trigger
CREATE TRIGGER update_crawl_jobs_updated_at BEFORE UPDATE ON crawl_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE crawl_jobs IS 'Crawl job management and tracking';
COMMENT ON COLUMN crawl_jobs.source IS 'Data source: tabelog, google_maps, instagram, etc.';
COMMENT ON COLUMN crawl_jobs.job_type IS 'Job type: full, incremental, update';
COMMENT ON COLUMN crawl_jobs.status IS 'Job status: pending, running, completed, failed, cancelled, paused';
COMMENT ON COLUMN crawl_jobs.priority IS 'Job priority: 1 (highest) to 10 (lowest)';
COMMENT ON COLUMN crawl_jobs.config IS 'Crawl configuration in JSON format (rate limits, proxies, etc.)';
COMMENT ON COLUMN crawl_jobs.next_run_at IS 'Scheduled time for next run (for recurring jobs)';
