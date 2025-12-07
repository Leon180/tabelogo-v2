package models

import (
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the status of a scraping job
type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

// JobID is a unique identifier for a scraping job
type JobID struct {
	value uuid.UUID
}

// NewJobID creates a new JobID
func NewJobID() JobID {
	return JobID{value: uuid.New()}
}

// ParseJobID parses a string into a JobID
func ParseJobID(s string) (JobID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return JobID{}, err
	}
	return JobID{value: id}, nil
}

// String returns the string representation of the JobID
func (id JobID) String() string {
	return id.value.String()
}

// ScrapingJob is the aggregate root for a scraping job
type ScrapingJob struct {
	id          JobID
	googleID    string
	area        string
	placeName   string
	status      JobStatus
	results     []TabelogRestaurant
	errorMsg    string
	createdAt   time.Time
	startedAt   *time.Time
	completedAt *time.Time
}

// NewScrapingJob creates a new scraping job
func NewScrapingJob(googleID, area, placeName string) *ScrapingJob {
	return &ScrapingJob{
		id:        NewJobID(),
		googleID:  googleID,
		area:      area,
		placeName: placeName,
		status:    JobStatusPending,
		results:   []TabelogRestaurant{},
		createdAt: time.Now(),
	}
}

// ID returns the job ID
func (j *ScrapingJob) ID() JobID {
	return j.id
}

// GoogleID returns the Google Place ID
func (j *ScrapingJob) GoogleID() string {
	return j.googleID
}

// Area returns the search area
func (j *ScrapingJob) Area() string {
	return j.area
}

// PlaceName returns the place name
func (j *ScrapingJob) PlaceName() string {
	return j.placeName
}

// Status returns the job status
func (j *ScrapingJob) Status() JobStatus {
	return j.status
}

// Results returns the scraping results
func (j *ScrapingJob) Results() []TabelogRestaurant {
	return j.results
}

// Error returns the error message
func (j *ScrapingJob) Error() string {
	return j.errorMsg
}

// CreatedAt returns the creation time
func (j *ScrapingJob) CreatedAt() time.Time {
	return j.createdAt
}

// StartedAt returns the start time
func (j *ScrapingJob) StartedAt() *time.Time {
	return j.startedAt
}

// CompletedAt returns the completion time
func (j *ScrapingJob) CompletedAt() *time.Time {
	return j.completedAt
}

// Start marks the job as running
func (j *ScrapingJob) Start() {
	j.status = JobStatusRunning
	now := time.Now()
	j.startedAt = &now
}

// Complete marks the job as completed with results
func (j *ScrapingJob) Complete(results []TabelogRestaurant) {
	j.status = JobStatusCompleted
	j.results = results
	now := time.Now()
	j.completedAt = &now
}

// Fail marks the job as failed with an error
func (j *ScrapingJob) Fail(err error) {
	j.status = JobStatusFailed
	j.errorMsg = err.Error()
	now := time.Now()
	j.completedAt = &now
}

// IsCompleted returns true if the job is completed or failed
func (j *ScrapingJob) IsCompleted() bool {
	return j.status == JobStatusCompleted || j.status == JobStatusFailed
}

// Duration returns the job duration
func (j *ScrapingJob) Duration() time.Duration {
	if j.startedAt == nil {
		return 0
	}
	if j.completedAt == nil {
		return time.Since(*j.startedAt)
	}
	return j.completedAt.Sub(*j.startedAt)
}
