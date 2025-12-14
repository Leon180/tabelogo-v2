package models

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScrapingJob(t *testing.T) {
	// Arrange
	googleID := "test-google-id"
	area := "Tokyo"
	placeName := "Test Restaurant"

	// Act
	job := NewScrapingJob(googleID, area, placeName)

	// Assert
	assert.NotNil(t, job)
	assert.NotEqual(t, uuid.Nil, job.ID())
	assert.Equal(t, googleID, job.GoogleID())
	assert.Equal(t, area, job.Area())
	assert.Equal(t, placeName, job.PlaceName())
	assert.Equal(t, JobStatusPending, job.Status())
	assert.Empty(t, job.Results())
	assert.Empty(t, job.Error())
	assert.False(t, job.CreatedAt().IsZero())
	assert.Nil(t, job.StartedAt())
	assert.Nil(t, job.CompletedAt())
}

func TestScrapingJob_Start(t *testing.T) {
	// Arrange
	job := NewScrapingJob("test-id", "Tokyo", "Test")

	// Act
	job.Start()

	// Assert
	assert.Equal(t, JobStatusRunning, job.Status())
	assert.NotNil(t, job.StartedAt())
	assert.False(t, job.StartedAt().IsZero())
}

func TestScrapingJob_Complete(t *testing.T) {
	// Arrange
	job := NewScrapingJob("test-id", "Tokyo", "Test")
	results := []TabelogRestaurant{
		*NewTabelogRestaurant("https://tabelog.com/1", "Restaurant 1", 3.5, 100, 50, "03-1234-5678", []string{"Japanese"}, []string{}),
		*NewTabelogRestaurant("https://tabelog.com/2", "Restaurant 2", 4.0, 200, 75, "03-8765-4321", []string{"Sushi"}, []string{}),
	}

	// Act
	job.Complete(results)

	// Assert
	assert.Equal(t, JobStatusCompleted, job.Status())
	assert.Len(t, job.Results(), 2)
	assert.Empty(t, job.Error())
	assert.NotNil(t, job.CompletedAt())
}

func TestScrapingJob_Complete_EmptyResults(t *testing.T) {
	// Arrange
	job := NewScrapingJob("test-id", "Tokyo", "Test")

	// Act
	job.Complete([]TabelogRestaurant{})

	// Assert
	assert.Equal(t, JobStatusCompleted, job.Status())
	assert.Empty(t, job.Results())
	assert.Empty(t, job.Error())
}

func TestScrapingJob_Fail(t *testing.T) {
	// Arrange
	job := NewScrapingJob("test-id", "Tokyo", "Test")
	err := errors.New("scraping failed: connection timeout")

	// Act
	job.Fail(err)

	// Assert
	assert.Equal(t, JobStatusFailed, job.Status())
	assert.NotEmpty(t, job.Error())
	assert.Equal(t, err.Error(), job.Error())
	assert.Empty(t, job.Results())
	assert.NotNil(t, job.CompletedAt())
}

func TestScrapingJob_StateTransitions(t *testing.T) {
	tests := []struct {
		name           string
		transitions    func(*ScrapingJob)
		expectedStatus JobStatus
	}{
		{
			name: "pending -> running -> completed",
			transitions: func(job *ScrapingJob) {
				job.Start()
				job.Complete([]TabelogRestaurant{})
			},
			expectedStatus: JobStatusCompleted,
		},
		{
			name: "pending -> running -> failed",
			transitions: func(job *ScrapingJob) {
				job.Start()
				job.Fail(errors.New("error"))
			},
			expectedStatus: JobStatusFailed,
		},
		{
			name: "pending -> failed (direct)",
			transitions: func(job *ScrapingJob) {
				job.Fail(errors.New("validation error"))
			},
			expectedStatus: JobStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			job := NewScrapingJob("test-id", "Tokyo", "Test")
			initialStatus := job.Status()

			// Act
			tt.transitions(job)

			// Assert
			assert.Equal(t, JobStatusPending, initialStatus)
			assert.Equal(t, tt.expectedStatus, job.Status())
		})
	}
}

func TestScrapingJob_IsCompleted(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*ScrapingJob)
		expected bool
	}{
		{
			name:     "pending job is not completed",
			setup:    func(job *ScrapingJob) {},
			expected: false,
		},
		{
			name: "running job is not completed",
			setup: func(job *ScrapingJob) {
				job.Start()
			},
			expected: false,
		},
		{
			name: "completed job is completed",
			setup: func(job *ScrapingJob) {
				job.Complete([]TabelogRestaurant{})
			},
			expected: true,
		},
		{
			name: "failed job is completed",
			setup: func(job *ScrapingJob) {
				job.Fail(errors.New("error"))
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			job := NewScrapingJob("test-id", "Tokyo", "Test")
			tt.setup(job)

			// Act
			result := job.IsCompleted()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScrapingJob_Duration(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*ScrapingJob)
		validate func(*testing.T, time.Duration)
	}{
		{
			name:  "not started job has zero duration",
			setup: func(job *ScrapingJob) {},
			validate: func(t *testing.T, d time.Duration) {
				assert.Equal(t, time.Duration(0), d)
			},
		},
		{
			name: "running job has positive duration",
			setup: func(job *ScrapingJob) {
				job.Start()
				time.Sleep(10 * time.Millisecond)
			},
			validate: func(t *testing.T, d time.Duration) {
				assert.Greater(t, d, time.Duration(0))
			},
		},
		{
			name: "completed job has fixed duration",
			setup: func(job *ScrapingJob) {
				job.Start()
				time.Sleep(10 * time.Millisecond)
				job.Complete([]TabelogRestaurant{})
			},
			validate: func(t *testing.T, d time.Duration) {
				assert.Greater(t, d, time.Duration(0))
				assert.Less(t, d, 1*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			job := NewScrapingJob("test-id", "Tokyo", "Test")
			tt.setup(job)

			// Act
			duration := job.Duration()

			// Assert
			tt.validate(t, duration)
		})
	}
}

func TestScrapingJob_MarshalJSON(t *testing.T) {
	// Arrange
	job := NewScrapingJob("test-google-id", "Tokyo", "Test Restaurant")
	results := []TabelogRestaurant{
		*NewTabelogRestaurant("https://tabelog.com/1", "Restaurant 1", 3.5, 100, 50, "03-1234-5678", []string{"Japanese"}, []string{"photo1.jpg"}),
	}
	job.Start()
	job.Complete(results)

	// Act
	data, err := json.Marshal(job)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify we can unmarshal it back
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, job.ID().String(), unmarshaled["id"])
	assert.Equal(t, job.GoogleID(), unmarshaled["google_id"])
	assert.Equal(t, string(JobStatusCompleted), unmarshaled["status"])
}

func TestScrapingJob_UnmarshalJSON(t *testing.T) {
	// Arrange
	jobID := uuid.New().String()
	jsonData := `{
		"id": "` + jobID + `",
		"google_id": "test-google-id",
		"area": "Tokyo",
		"place_name": "Test Restaurant",
		"status": "COMPLETED",
		"results": [],
		"created_at": "2025-01-01T00:00:00Z"
	}`

	// Act
	var job ScrapingJob
	err := json.Unmarshal([]byte(jsonData), &job)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, jobID, job.ID().String())
	assert.Equal(t, "test-google-id", job.GoogleID())
	assert.Equal(t, "Tokyo", job.Area())
	assert.Equal(t, "Test Restaurant", job.PlaceName())
	assert.Equal(t, JobStatusCompleted, job.Status())
}

func TestParseJobID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "valid UUID",
			input:     "550e8400-e29b-41d4-a716-446655440000",
			wantError: false,
		},
		{
			name:      "invalid UUID",
			input:     "invalid-uuid",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			jobID, err := ParseJobID(tt.input)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, jobID.String())
			}
		})
	}
}
