//go:generate mockgen -destination=../../testutil/mocks/mock_job_repository.go -package=mocks github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories JobRepository
package repositories

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
)

// JobRepository defines the interface for scraping job persistence
type JobRepository interface {
	// Save saves a scraping job
	Save(ctx context.Context, job *models.ScrapingJob) error

	// FindByID finds a job by ID
	FindByID(ctx context.Context, id models.JobID) (*models.ScrapingJob, error)

	// FindByGoogleID finds jobs by Google Place ID
	FindByGoogleID(ctx context.Context, googleID string) ([]*models.ScrapingJob, error)

	// Update updates a job
	Update(ctx context.Context, job *models.ScrapingJob) error

	// Delete deletes a job
	Delete(ctx context.Context, id models.JobID) error

	// FindPending finds all pending jobs
	FindPending(ctx context.Context, limit int) ([]*models.ScrapingJob, error)
}
