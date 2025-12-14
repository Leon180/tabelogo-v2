package usecases

import (
	"context"
	"fmt"

	"github.com/Leon180/tabelogo-v2/internal/spider/application/services"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
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
	jobRepo      repositories.JobRepository
	jobProcessor *services.JobProcessor
	logger       *zap.Logger
}

// NewScrapeRestaurantUseCase creates a new use case
func NewScrapeRestaurantUseCase(
	jobRepo repositories.JobRepository,
	jobProcessor *services.JobProcessor,
	logger *zap.Logger,
) *ScrapeRestaurantUseCase {
	return &ScrapeRestaurantUseCase{
		jobRepo:      jobRepo,
		jobProcessor: jobProcessor,
		logger:       logger.With(zap.String("usecase", "scrape_restaurant")),
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

	// Save job to repository
	if err := uc.jobRepo.Save(ctx, job); err != nil {
		uc.logger.Error("Failed to save job",
			zap.String("job_id", job.ID().String()),
			zap.String("google_id", req.GoogleID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to save scraping job for place '%s' (google_id: %s): %w", req.PlaceName, req.GoogleID, err)
	}

	// Submit to job processor (worker pool)
	if err := uc.jobProcessor.SubmitJob(ctx, job.ID()); err != nil {
		uc.logger.Error("Failed to submit job to processor",
			zap.String("job_id", job.ID().String()),
			zap.String("google_id", req.GoogleID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to submit scraping job %s for place '%s': %w", job.ID().String(), req.PlaceName, err)
	}

	uc.logger.Info("Job submitted to worker pool",
		zap.String("job_id", job.ID().String()),
	)

	return &ScrapeRestaurantResponse{
		JobID:  job.ID().String(),
		Status: string(job.Status()),
	}, nil
}
