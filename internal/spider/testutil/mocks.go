package testutil

import (
	"context"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
)

// MockJobRepository is a mock implementation of JobRepository for testing
type MockJobRepository struct {
	SaveFunc           func(ctx context.Context, job *models.ScrapingJob) error
	FindByIDFunc       func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error)
	FindByGoogleIDFunc func(ctx context.Context, googleID string) ([]*models.ScrapingJob, error)
	UpdateFunc         func(ctx context.Context, job *models.ScrapingJob) error
	DeleteFunc         func(ctx context.Context, id models.JobID) error
	FindPendingFunc    func(ctx context.Context, limit int) ([]*models.ScrapingJob, error)
}

func (m *MockJobRepository) Save(ctx context.Context, job *models.ScrapingJob) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, job)
	}
	return nil
}

func (m *MockJobRepository) FindByID(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockJobRepository) FindByGoogleID(ctx context.Context, googleID string) ([]*models.ScrapingJob, error) {
	if m.FindByGoogleIDFunc != nil {
		return m.FindByGoogleIDFunc(ctx, googleID)
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

func (m *MockJobRepository) FindPending(ctx context.Context, limit int) ([]*models.ScrapingJob, error) {
	if m.FindPendingFunc != nil {
		return m.FindPendingFunc(ctx, limit)
	}
	return nil, nil
}

// MockResultCacheRepository is a mock implementation of ResultCacheRepository
type MockResultCacheRepository struct {
	GetFunc    func(ctx context.Context, placeID string) (*models.CachedResult, error)
	SetFunc    func(ctx context.Context, placeID string, results []models.TabelogRestaurant, ttl time.Duration) error
	DeleteFunc func(ctx context.Context, placeID string) error
}

func (m *MockResultCacheRepository) Get(ctx context.Context, placeID string) (*models.CachedResult, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, placeID)
	}
	return nil, nil
}

func (m *MockResultCacheRepository) Set(ctx context.Context, placeID string, results []models.TabelogRestaurant, ttl time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, placeID, results, ttl)
	}
	return nil
}

func (m *MockResultCacheRepository) Delete(ctx context.Context, placeID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, placeID)
	}
	return nil
}

// MockJobProcessor is a mock implementation of JobProcessor
type MockJobProcessor struct {
	SubmitJobFunc func(ctx context.Context, jobID models.JobID) error
	StartFunc     func(ctx context.Context) error
	StopFunc      func(ctx context.Context) error
}

func (m *MockJobProcessor) SubmitJob(ctx context.Context, jobID models.JobID) error {
	if m.SubmitJobFunc != nil {
		return m.SubmitJobFunc(ctx, jobID)
	}
	return nil
}

func (m *MockJobProcessor) Start(ctx context.Context) error {
	if m.StartFunc != nil {
		return m.StartFunc(ctx)
	}
	return nil
}

func (m *MockJobProcessor) Stop(ctx context.Context) error {
	if m.StopFunc != nil {
		return m.StopFunc(ctx)
	}
	return nil
}
