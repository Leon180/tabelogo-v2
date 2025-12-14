package usecases

import (
	"context"
	"fmt"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"go.uber.org/zap"
)

// GetJobStatusUseCase handles getting job status
type GetJobStatusUseCase struct {
	jobRepo repositories.JobRepository
	logger  *zap.Logger
}

// NewGetJobStatusUseCase creates a new use case
func NewGetJobStatusUseCase(
	jobRepo repositories.JobRepository,
	logger *zap.Logger,
) *GetJobStatusUseCase {
	return &GetJobStatusUseCase{
		jobRepo: jobRepo,
		logger:  logger.With(zap.String("usecase", "get_job_status")),
	}
}

// Execute executes the use case
func (uc *GetJobStatusUseCase) Execute(ctx context.Context, jobID string) (*models.ScrapingJob, error) {
	id, err := models.ParseJobID(jobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID: %w", err)
	}

	job, err := uc.jobRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	return job, nil
}
