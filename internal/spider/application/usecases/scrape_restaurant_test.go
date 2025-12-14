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

func TestScrapeRestaurantUseCase_Execute_Success(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	mockJobRepo := &testutil.MockJobRepository{
		SaveFunc: func(ctx context.Context, job *models.ScrapingJob) error {
			return nil
		},
	}
	mockProcessor := &testutil.MockJobProcessor{
		SubmitJobFunc: func(ctx context.Context, jobID models.JobID) error {
			return nil
		},
	}

	useCase := NewScrapeRestaurantUseCase(mockJobRepo, mockProcessor, logger)

	req := ScrapeRestaurantRequest{
		GoogleID:  "test-google-id",
		Area:      "Tokyo",
		PlaceName: "Test Restaurant",
	}

	// Act
	resp, err := useCase.Execute(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.JobID)
	assert.Equal(t, string(models.JobStatusPending), resp.Status)
}

func TestScrapeRestaurantUseCase_Execute_SaveError(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	expectedErr := errors.New("save failed")
	mockJobRepo := &testutil.MockJobRepository{
		SaveFunc: func(ctx context.Context, job *models.ScrapingJob) error {
			return expectedErr
		},
	}
	mockProcessor := &testutil.MockJobProcessor{}

	useCase := NewScrapeRestaurantUseCase(mockJobRepo, mockProcessor, logger)

	req := ScrapeRestaurantRequest{
		GoogleID:  "test-google-id",
		Area:      "Tokyo",
		PlaceName: "Test Restaurant",
	}

	// Act
	resp, err := useCase.Execute(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to save scraping job")
}

func TestScrapeRestaurantUseCase_Execute_SubmitError(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	mockJobRepo := &testutil.MockJobRepository{
		SaveFunc: func(ctx context.Context, job *models.ScrapingJob) error {
			return nil
		},
	}

	expectedErr := errors.New("submit failed")
	mockProcessor := &testutil.MockJobProcessor{
		SubmitJobFunc: func(ctx context.Context, jobID models.JobID) error {
			return expectedErr
		},
	}

	useCase := NewScrapeRestaurantUseCase(mockJobRepo, mockProcessor, logger)

	req := ScrapeRestaurantRequest{
		GoogleID:  "test-google-id",
		Area:      "Tokyo",
		PlaceName: "Test Restaurant",
	}

	// Act
	resp, err := useCase.Execute(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to submit scraping job")
}

func TestScrapeRestaurantUseCase_Execute_ContextCancellation(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	mockJobRepo := &testutil.MockJobRepository{
		SaveFunc: func(ctx context.Context, job *models.ScrapingJob) error {
			// Simulate slow operation
			<-ctx.Done()
			return ctx.Err()
		},
	}
	mockProcessor := &testutil.MockJobProcessor{}

	useCase := NewScrapeRestaurantUseCase(mockJobRepo, mockProcessor, logger)

	req := ScrapeRestaurantRequest{
		GoogleID:  "test-google-id",
		Area:      "Tokyo",
		PlaceName: "Test Restaurant",
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	resp, err := useCase.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
}
