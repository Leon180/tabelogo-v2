package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisJobStore implements JobRepository using Redis
type RedisJobStore struct {
	client *redisclient.Client
	logger *zap.Logger
	ttl    time.Duration
}

// NewRedisJobStore creates a new Redis job store
func NewRedisJobStore(client *redisclient.Client, logger *zap.Logger) repositories.JobRepository {
	return &RedisJobStore{
		client: client,
		logger: logger.With(zap.String("component", "redis_job_store")),
		ttl:    24 * time.Hour, // 24 hours TTL
	}
}

// Save saves a scraping job
func (r *RedisJobStore) Save(ctx context.Context, job *models.ScrapingJob) error {
	key := r.jobKey(job.ID())

	// Convert job to JSON
	data, err := json.Marshal(job)
	if err != nil {
		r.logger.Error("Failed to marshal job", zap.Error(err), zap.String("job_id", job.ID().String()))
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Save to Redis with TTL
	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		r.logger.Error("Failed to save job to Redis", zap.Error(err), zap.String("job_id", job.ID().String()))
		return fmt.Errorf("failed to save job: %w", err)
	}

	// Add to index by Google ID for quick lookup
	indexKey := r.googleIDIndexKey(job.GoogleID())
	if err := r.client.SAdd(ctx, indexKey, job.ID().String()).Err(); err != nil {
		r.logger.Warn("Failed to add job to Google ID index", zap.Error(err))
	}
	r.client.Expire(ctx, indexKey, r.ttl)

	// Add to pending jobs set if status is pending
	if job.Status() == models.JobStatusPending {
		pendingKey := r.pendingJobsKey()
		if err := r.client.SAdd(ctx, pendingKey, job.ID().String()).Err(); err != nil {
			r.logger.Warn("Failed to add job to pending set", zap.Error(err))
		}
	}

	r.logger.Info("Job saved to Redis",
		zap.String("job_id", job.ID().String()),
		zap.String("google_id", job.GoogleID()),
		zap.String("status", string(job.Status())),
	)

	return nil
}

// FindByID finds a job by ID
func (r *RedisJobStore) FindByID(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
	key := r.jobKey(id)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redisclient.Nil {
		return nil, fmt.Errorf("job not found: %s", id.String())
	}
	if err != nil {
		r.logger.Error("Failed to get job from Redis", zap.Error(err), zap.String("job_id", id.String()))
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	var job models.ScrapingJob
	if err := json.Unmarshal(data, &job); err != nil {
		r.logger.Error("Failed to unmarshal job", zap.Error(err), zap.String("job_id", id.String()))
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

// FindByGoogleID finds jobs by Google Place ID
func (r *RedisJobStore) FindByGoogleID(ctx context.Context, googleID string) ([]*models.ScrapingJob, error) {
	indexKey := r.googleIDIndexKey(googleID)

	jobIDs, err := r.client.SMembers(ctx, indexKey).Result()
	if err != nil {
		r.logger.Error("Failed to get job IDs from index", zap.Error(err), zap.String("google_id", googleID))
		return nil, fmt.Errorf("failed to get job IDs: %w", err)
	}

	jobs := make([]*models.ScrapingJob, 0, len(jobIDs))
	for _, jobIDStr := range jobIDs {
		jobID, err := models.ParseJobID(jobIDStr)
		if err != nil {
			r.logger.Warn("Failed to parse job ID", zap.Error(err), zap.String("job_id", jobIDStr))
			continue
		}
		job, err := r.FindByID(ctx, jobID)
		if err != nil {
			r.logger.Warn("Failed to get job", zap.Error(err), zap.String("job_id", jobIDStr))
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// Update updates a job
func (r *RedisJobStore) Update(ctx context.Context, job *models.ScrapingJob) error {
	// Check if job exists
	key := r.jobKey(job.ID())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check job existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("job not found: %s", job.ID().String())
	}

	// Update is same as Save
	return r.Save(ctx, job)
}

// Delete deletes a job
func (r *RedisJobStore) Delete(ctx context.Context, id models.JobID) error {
	// Get job first to clean up indexes
	job, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from main storage
	key := r.jobKey(id)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("Failed to delete job", zap.Error(err), zap.String("job_id", id.String()))
		return fmt.Errorf("failed to delete job: %w", err)
	}

	// Remove from Google ID index
	indexKey := r.googleIDIndexKey(job.GoogleID())
	r.client.SRem(ctx, indexKey, id.String())

	// Remove from pending jobs set
	pendingKey := r.pendingJobsKey()
	r.client.SRem(ctx, pendingKey, id.String())

	r.logger.Info("Job deleted", zap.String("job_id", id.String()))
	return nil
}

// FindPending finds all pending jobs
func (r *RedisJobStore) FindPending(ctx context.Context, limit int) ([]*models.ScrapingJob, error) {
	pendingKey := r.pendingJobsKey()

	// Get pending job IDs
	jobIDs, err := r.client.SMembers(ctx, pendingKey).Result()
	if err != nil {
		r.logger.Error("Failed to get pending job IDs", zap.Error(err))
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}

	// Limit results
	if limit > 0 && len(jobIDs) > limit {
		jobIDs = jobIDs[:limit]
	}

	// Fetch jobs
	jobs := make([]*models.ScrapingJob, 0, len(jobIDs))
	for _, jobIDStr := range jobIDs {
		jobID, err := models.ParseJobID(jobIDStr)
		if err != nil {
			r.logger.Warn("Failed to parse job ID", zap.Error(err), zap.String("job_id", jobIDStr))
			// Remove from pending set if invalid
			r.client.SRem(ctx, pendingKey, jobIDStr)
			continue
		}
		job, err := r.FindByID(ctx, jobID)
		if err != nil {
			r.logger.Warn("Failed to get pending job", zap.Error(err), zap.String("job_id", jobIDStr))
			// Remove from pending set if not found
			r.client.SRem(ctx, pendingKey, jobIDStr)
			continue
		}

		// Double-check status (in case it was updated)
		if job.Status() == models.JobStatusPending {
			jobs = append(jobs, job)
		} else {
			// Remove from pending set if status changed
			r.client.SRem(ctx, pendingKey, jobIDStr)
		}
	}

	return jobs, nil
}

// Helper methods for Redis keys
func (r *RedisJobStore) jobKey(id models.JobID) string {
	return fmt.Sprintf("spider:jobs:%s", id.String())
}

func (r *RedisJobStore) googleIDIndexKey(googleID string) string {
	return fmt.Sprintf("spider:jobs:google:%s", googleID)
}

func (r *RedisJobStore) pendingJobsKey() string {
	return "spider:jobs:pending"
}
