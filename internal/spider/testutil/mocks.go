package testutil

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
)

// MockJobRepository is a mock implementation of JobRepository for testing
type MockJobRepository struct {
	SaveFunc   func(ctx context.Context, job *models.ScrapingJob) error
	GetFunc    func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error)
	UpdateFunc func(ctx context.Context, job *models.ScrapingJob) error
	DeleteFunc func(ctx context.Context, id models.JobID) error
}

func (m *MockJobRepository) Save(ctx context.Context, job *models.ScrapingJob) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, job)
	}
	return nil
}

func (m *MockJobRepository) Get(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockJobRepository) Update(ctx context.Context, job *models.ScrapingJob) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, job)
	}
	return nil
}

func (m *MockJobRepository) Delete(ctx context.Context, id models.JobID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockResultCacheRepository is a mock implementation of ResultCacheRepository
type MockResultCacheRepository struct {
	GetFunc func(ctx context.Context, placeID string) (*models.CachedResult, error)
	SetFunc func(ctx context.Context, placeID string, results []*models.TabelogRestaurant) error
}

func (m *MockResultCacheRepository) Get(ctx context.Context, placeID string) (*models.CachedResult, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, placeID)
	}
	return nil, nil
}

func (m *MockResultCacheRepository) Set(ctx context.Context, placeID string, results []*models.TabelogRestaurant) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, placeID, results)
	}
	return nil
}

// MockJobProcessor is a mock implementation of JobProcessor
type MockJobProcessor struct {
	SubmitJobFunc func(ctx context.Context, jobID models.JobID) error
}

func (m *MockJobProcessor) SubmitJob(ctx context.Context, jobID models.JobID) error {
	if m.SubmitJobFunc != nil {
		return m.SubmitJobFunc(ctx, jobID)
	}
	return nil
}
