package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/zap"
)

// JobProcessor handles async job processing with worker pool
type JobProcessor struct {
	jobRepo     repositories.JobRepository
	resultCache repositories.ResultCacheRepository
	scraper     *scraper.Scraper
	metrics     *metrics.SpiderMetrics
	logger      *zap.Logger
	workerCount int
	jobQueue    chan models.JobID
	stopChan    chan struct{}
	wg          sync.WaitGroup
	rateLimiter *DynamicRateLimiter
}

// NewJobProcessor creates a new job processor
func NewJobProcessor(
	jobRepo repositories.JobRepository,
	resultCache repositories.ResultCacheRepository,
	scraper *scraper.Scraper,
	metrics *metrics.SpiderMetrics,
	logger *zap.Logger,
	workerCount int,
) *JobProcessor {
	return &JobProcessor{
		jobRepo:     jobRepo,
		resultCache: resultCache,
		scraper:     scraper,
		metrics:     metrics,
		logger:      logger.With(zap.String("component", "job_processor")),
		workerCount: workerCount,
		jobQueue:    make(chan models.JobID, 100), // Buffer of 100 jobs
		stopChan:    make(chan struct{}),
		rateLimiter: NewDynamicRateLimiter(30, logger), // 30 req/min
	}
}

// Start starts the worker pool
func (p *JobProcessor) Start(ctx context.Context) {
	p.logger.Info("Starting job processor", zap.Int("workers", p.workerCount))

	// Set worker pool size metric
	p.metrics.SetWorkerPoolSize(p.workerCount)

	// Start worker goroutines
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer func() {
				// Recover from panic first
				if r := recover(); r != nil {
					p.logger.Error("Worker panic recovered",
						zap.Int("worker_id", workerID),
						zap.Any("panic", r),
						zap.Stack("stack"),
					)
					p.metrics.RecordScrapeError("worker_panic")
				}
				// Always call Done() last, after panic recovery
				p.wg.Done()
			}()

			p.worker(ctx, workerID)
		}(i)
	}
	// Start job fetcher
	p.wg.Add(1)
	go p.jobFetcher(ctx)

	p.logger.Info("Job processor started")
}

// Stop gracefully stops the worker pool
func (p *JobProcessor) Stop(ctx context.Context) error {
	p.logger.Info("Stopping job processor...")

	// Signal stop
	close(p.stopChan)

	// Wait for all workers to finish with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		p.logger.Info("Job processor stopped gracefully")
		return nil
	case <-ctx.Done():
		p.logger.Warn("Job processor stop timeout")
		return fmt.Errorf("shutdown timeout")
	}
}

// SubmitJob submits a job to the queue
func (p *JobProcessor) SubmitJob(ctx context.Context, jobID models.JobID) error {
	select {
	case p.jobQueue <- jobID:
		p.logger.Info("Job submitted to queue", zap.String("job_id", jobID.String()))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("job queue is full")
	}
}

// worker processes jobs from the queue
func (p *JobProcessor) worker(ctx context.Context, workerID int) {
	// Note: wg.Done() is called in the goroutine's defer, not here
	logger := p.logger.With(zap.Int("worker_id", workerID))
	logger.Info("Worker started")

	for {
		select {
		case <-p.stopChan:
			logger.Info("Worker stopping")
			return
		case <-ctx.Done():
			logger.Info("Worker context cancelled")
			return
		case jobID := <-p.jobQueue:
			p.processJob(ctx, jobID, logger)
		}
	}
}

// processJob processes a single job
func (p *JobProcessor) processJob(ctx context.Context, jobID models.JobID, logger *zap.Logger) {
	logger = logger.With(zap.String("job_id", jobID.String()))
	logger.Info("Processing job")

	// Track job start time
	jobStartTime := time.Now()

	// Get job from repository
	job, err := p.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		logger.Error("Failed to get job", zap.Error(err))
		return
	}

	// Mark job as running
	job.Start()
	p.metrics.RecordJob("running")
	if err := p.jobRepo.Update(ctx, job); err != nil {
		logger.Error("Failed to update job status", zap.Error(err))
	}

	// Wait for rate limiter
	if err := p.rateLimiter.Wait(ctx); err != nil {
		logger.Error("Rate limiter error", zap.Error(err))
		job.Fail(fmt.Errorf("rate limiter error: %w", err))
		p.jobRepo.Update(ctx, job)
		p.metrics.RecordJob("failed")
		p.metrics.RecordJobDuration("failed", time.Since(jobStartTime).Seconds())
		return
	}

	// Scrape restaurants
	startTime := time.Now()
	results, err := p.scraper.ScrapeRestaurants(job.Area(), job.PlaceName())
	duration := time.Since(startTime)

	if err != nil {
		logger.Error("Scraping failed",
			zap.Error(err),
			zap.Duration("duration", duration),
		)

		// Check if it's a rate limit error (429)
		if isRateLimitError(err) {
			p.rateLimiter.OnRateLimitHit()
			logger.Warn("Rate limit hit, backing off")
		}

		job.Fail(err)
		p.jobRepo.Update(ctx, job)
		p.metrics.RecordJob("failed")
		p.metrics.RecordJobDuration("failed", time.Since(jobStartTime).Seconds())
		return
	}

	// On success, restore rate limiter
	p.rateLimiter.OnSuccess()

	logger.Info("Scraping completed",
		zap.Int("results_count", len(results)),
		zap.Duration("duration", duration),
	)

	// Convert to slice of pointers for Complete method
	resultPtrs := make([]models.TabelogRestaurant, len(results))
	for i := range results {
		resultPtrs[i] = results[i]
	}

	// Mark job as completed
	job.Complete(resultPtrs)
	if err := p.jobRepo.Update(ctx, job); err != nil {
		logger.Error("Failed to update job", zap.Error(err))
		return
	}

	// Record metrics
	p.metrics.RecordJob("completed")
	p.metrics.RecordJobDuration("completed", time.Since(jobStartTime).Seconds())
	p.metrics.RecordRestaurantsScraped("success", len(resultPtrs))

	// Cache results
	logger.Info("Attempting to cache results",
		zap.String("google_id", job.GoogleID()),
		zap.Int("results_count", len(resultPtrs)),
	)

	if err := p.resultCache.Set(ctx, job.GoogleID(), resultPtrs, 24*time.Hour); err != nil {
		logger.Error("Failed to cache results",
			zap.Error(err),
			zap.String("google_id", job.GoogleID()),
		)
	} else {
		logger.Info("Successfully cached results",
			zap.String("google_id", job.GoogleID()),
			zap.Int("results_count", len(resultPtrs)),
		)
	}

	logger.Info("Job completed successfully")
}

// jobFetcher periodically fetches pending jobs
func (p *JobProcessor) jobFetcher(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	p.logger.Info("Job fetcher started")

	for {
		select {
		case <-p.stopChan:
			p.logger.Info("Job fetcher stopping")
			return
		case <-ctx.Done():
			p.logger.Info("Job fetcher context cancelled")
			return
		case <-ticker.C:
			p.fetchPendingJobs(ctx)
		}
	}
}

// fetchPendingJobs fetches and submits pending jobs
func (p *JobProcessor) fetchPendingJobs(ctx context.Context) {
	jobs, err := p.jobRepo.FindPending(ctx, 10)
	if err != nil {
		p.logger.Error("Failed to fetch pending jobs", zap.Error(err))
		return
	}

	if len(jobs) == 0 {
		return
	}

	p.logger.Info("Found pending jobs", zap.Int("count", len(jobs)))

	for _, job := range jobs {
		select {
		case p.jobQueue <- job.ID():
			// Job submitted
		default:
			p.logger.Warn("Job queue full, skipping job", zap.String("job_id", job.ID().String()))
		}
	}
}

// isRateLimitError checks if error is a rate limit error
func isRateLimitError(err error) bool {
	// Check for common rate limit error patterns
	errStr := err.Error()
	return contains(errStr, "429") ||
		contains(errStr, "Too Many Requests") ||
		contains(errStr, "rate limit")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
