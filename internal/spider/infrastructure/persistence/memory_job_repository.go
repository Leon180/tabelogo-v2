package persistence

import (
	"context"
	"fmt"
	"sync"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
)

// InMemoryJobRepository is an in-memory implementation of JobRepository
type InMemoryJobRepository struct {
	jobs map[string]*models.ScrapingJob
	mu   sync.RWMutex
}

// NewInMemoryJobRepository creates a new in-memory repository
func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{
		jobs: make(map[string]*models.ScrapingJob),
	}
}

// Save saves a job
func (r *InMemoryJobRepository) Save(ctx context.Context, job *models.ScrapingJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.jobs[job.ID().String()] = job
	return nil
}

// FindByID finds a job by ID
func (r *InMemoryJobRepository) FindByID(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id.String()]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", id.String())
	}

	return job, nil
}

// FindByGoogleID finds jobs by Google Place ID
func (r *InMemoryJobRepository) FindByGoogleID(ctx context.Context, googleID string) ([]*models.ScrapingJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var jobs []*models.ScrapingJob
	for _, job := range r.jobs {
		if job.GoogleID() == googleID {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

// Update updates a job
func (r *InMemoryJobRepository) Update(ctx context.Context, job *models.ScrapingJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.ID().String()]; !exists {
		return fmt.Errorf("job not found: %s", job.ID().String())
	}

	r.jobs[job.ID().String()] = job
	return nil
}

// Delete deletes a job
func (r *InMemoryJobRepository) Delete(ctx context.Context, id models.JobID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.jobs, id.String())
	return nil
}

// FindPending finds all pending jobs
func (r *InMemoryJobRepository) FindPending(ctx context.Context, limit int) ([]*models.ScrapingJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var pending []*models.ScrapingJob
	for _, job := range r.jobs {
		if job.Status() == models.JobStatusPending {
			pending = append(pending, job)
			if len(pending) >= limit {
				break
			}
		}
	}

	return pending, nil
}
