package usecases

import (
	"context"
	"fmt"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/zap"
)

// ScrapeRestaurantRequest is the request for scraping a restaurant
type ScrapeRestaurantRequest struct {
	GoogleID  string
	Area      string
	PlaceName string
}

// ScrapeRestaurantResponse is the response for scraping a restaurant
type ScrapeRestaurantResponse struct {
	JobID  string
	Status string
}

// ScrapeRestaurantUseCase handles restaurant scraping
type ScrapeRestaurantUseCase struct {
	jobRepo repositories.JobRepository
	scraper *scraper.Scraper
	logger  *zap.Logger
}

// NewScrapeRestaurantUseCase creates a new use case
func NewScrapeRestaurantUseCase(
	jobRepo repositories.JobRepository,
	scraper *scraper.Scraper,
	logger *zap.Logger,
) *ScrapeRestaurantUseCase {
	return &ScrapeRestaurantUseCase{
		jobRepo: jobRepo,
		scraper: scraper,
		logger:  logger.With(zap.String("usecase", "scrape_restaurant")),
	}
}

// Execute executes the use case
func (uc *ScrapeRestaurantUseCase) Execute(ctx context.Context, req ScrapeRestaurantRequest) (*ScrapeRestaurantResponse, error) {
	uc.logger.Info("Starting scrape job",
		zap.String("google_id", req.GoogleID),
		zap.String("area", req.Area),
		zap.String("place_name", req.PlaceName),
	)

	// Create job
	job := models.NewScrapingJob(req.GoogleID, req.Area, req.PlaceName)

	// Save job
	if err := uc.jobRepo.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to save job: %w", err)
	}

	// Start job (in real implementation, this would be queued)
	go uc.processJob(context.Background(), job)

	return &ScrapeRestaurantResponse{
		JobID:  job.ID().String(),
		Status: string(job.Status()),
	}, nil
}

// processJob processes a scraping job
func (uc *ScrapeRestaurantUseCase) processJob(ctx context.Context, job *models.ScrapingJob) {
	job.Start()
	if err := uc.jobRepo.Update(ctx, job); err != nil {
		uc.logger.Error("Failed to update job status",
			zap.String("job_id", job.ID().String()),
			zap.Error(err),
		)
	}

	// Scrape
	restaurants, err := uc.scraper.ScrapeRestaurants(job.Area(), job.PlaceName())
	if err != nil {
		job.Fail(err)
		uc.logger.Error("Scraping failed",
			zap.String("job_id", job.ID().String()),
			zap.Error(err),
		)
	} else {
		job.Complete(restaurants)
		uc.logger.Info("Scraping completed",
			zap.String("job_id", job.ID().String()),
			zap.Int("results", len(restaurants)),
		)
	}

	// Update job
	if err := uc.jobRepo.Update(ctx, job); err != nil {
		uc.logger.Error("Failed to update job",
			zap.String("job_id", job.ID().String()),
			zap.Error(err),
		)
	}
}
