package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetJobStatusUseCase_Execute_Success(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	job := testutil.CreateTestJob()
	job.Start()
	job.Complete([]models.TabelogRestaurant{*testutil.CreateTestRestaurant()})

	mockJobRepo := &testutil.MockJobRepository{
		FindByIDFunc: func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
			return job, nil
		},
	}

	useCase := NewGetJobStatusUseCase(mockJobRepo, logger)

	// Act
	result, err := useCase.Execute(context.Background(), job.ID().String())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, job.ID(), result.ID())
	assert.Equal(t, models.JobStatusCompleted, result.Status())
}

func TestGetJobStatusUseCase_Execute_InvalidJobID(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	mockJobRepo := &testutil.MockJobRepository{}

	useCase := NewGetJobStatusUseCase(mockJobRepo, logger)

	// Act
	result, err := useCase.Execute(context.Background(), "invalid-uuid")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid job ID")
}

func TestGetJobStatusUseCase_Execute_JobNotFound(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	mockJobRepo := &testutil.MockJobRepository{
		FindByIDFunc: func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
			return nil, errors.New("not found")
		},
	}

	useCase := NewGetJobStatusUseCase(mockJobRepo, logger)
	jobID := models.NewJobID()

	// Act
	result, err := useCase.Execute(context.Background(), jobID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "job not found")
}

func TestGetJobStatusUseCase_Execute_ContextCancellation(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	mockJobRepo := &testutil.MockJobRepository{
		FindByIDFunc: func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	}

	useCase := NewGetJobStatusUseCase(mockJobRepo, logger)
	jobID := models.NewJobID()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	result, err := useCase.Execute(ctx, jobID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetJobStatusUseCase_Execute_DifferentStatuses(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*models.ScrapingJob)
		status models.JobStatus
	}{
		{
			name:   "pending job",
			setup:  func(job *models.ScrapingJob) {},
			status: models.JobStatusPending,
		},
		{
			name: "running job",
			setup: func(job *models.ScrapingJob) {
				job.Start()
			},
			status: models.JobStatusRunning,
		},
		{
			name: "completed job",
			setup: func(job *models.ScrapingJob) {
				job.Start()
				job.Complete([]models.TabelogRestaurant{})
			},
			status: models.JobStatusCompleted,
		},
		{
			name: "failed job",
			setup: func(job *models.ScrapingJob) {
				job.Start()
				job.Fail(errors.New("test error"))
			},
			status: models.JobStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zap.NewNop()
			job := testutil.CreateTestJob()
			tt.setup(job)

			mockJobRepo := &testutil.MockJobRepository{
				FindByIDFunc: func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
					return job, nil
				},
			}

			useCase := NewGetJobStatusUseCase(mockJobRepo, logger)

			// Act
			result, err := useCase.Execute(context.Background(), job.ID().String())

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.status, result.Status())
		})
	}
}
