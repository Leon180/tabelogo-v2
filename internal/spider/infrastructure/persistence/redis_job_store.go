package persistence

import (
	"context"
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
	helper *RedisHelper
	logger *zap.Logger
	ttl    time.Duration
}

// NewRedisJobStore creates a new Redis job store
func NewRedisJobStore(client *redisclient.Client, logger *zap.Logger) repositories.JobRepository {
	return &RedisJobStore{
		client: client,
		helper: NewRedisHelper(client, logger),
		logger: logger.With(zap.String("component", "redis_job_store")),
		ttl:    24 * time.Hour, // 24 hours TTL
	}
}

// Save saves a scraping job
func (r *RedisJobStore) Save(ctx context.Context, job *models.ScrapingJob) error {
	key := r.jobKey(job.ID())

	// Save to Redis with TTL using helper
	if err := r.helper.SetJSON(ctx, key, job, r.ttl); err != nil {
		r.logger.Error("Failed to save job to Redis", zap.Error(err), zap.String("job_id", job.ID().String()))
		return err
	}

	// Add to index by Google ID for quick lookup
	indexKey := r.googleIDIndexKey(job.GoogleID())
	if err := r.helper.SetAdd(ctx, indexKey, job.ID().String()); err != nil {
		r.logger.Warn("Failed to add job to Google ID index", zap.Error(err))
	}
	r.helper.Expire(ctx, indexKey, r.ttl)

	// Add to pending jobs set if status is pending
	if job.Status() == models.JobStatusPending {
		pendingKey := r.pendingJobsKey()
		if err := r.helper.SetAdd(ctx, pendingKey, job.ID().String()); err != nil {
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

	var job models.ScrapingJob
	if err := r.helper.GetJSON(ctx, key, &job); err != nil {
		r.logger.Error("Failed to get job from Redis", zap.Error(err), zap.String("job_id", id.String()))
		return nil, err
	}

	return &job, nil
}

// FindByGoogleID finds jobs by Google Place ID
func (r *RedisJobStore) FindByGoogleID(ctx context.Context, googleID string) ([]*models.ScrapingJob, error) {
	indexKey := r.googleIDIndexKey(googleID)

	jobIDs, err := r.helper.SetMembers(ctx, indexKey)
	if err != nil {
		r.logger.Error("Failed to get job IDs from index", zap.Error(err), zap.String("google_id", googleID))
		return nil, err
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
	if err := r.helper.Delete(ctx, key); err != nil {
		r.logger.Error("Failed to delete job", zap.Error(err), zap.String("job_id", id.String()))
		return err
	}

	// Remove from Google ID index
	indexKey := r.googleIDIndexKey(job.GoogleID())
	r.helper.SetRemove(ctx, indexKey, id.String())

	// Remove from pending jobs set
	pendingKey := r.pendingJobsKey()
	r.helper.SetRemove(ctx, pendingKey, id.String())

	r.logger.Info("Job deleted", zap.String("job_id", id.String()))
	return nil
}

// FindPending finds all pending jobs
func (r *RedisJobStore) FindPending(ctx context.Context, limit int) ([]*models.ScrapingJob, error) {
	pendingKey := r.pendingJobsKey()

	// Get pending job IDs
	jobIDs, err := r.helper.SetMembers(ctx, pendingKey)
	if err != nil {
		r.logger.Error("Failed to get pending job IDs", zap.Error(err))
		return nil, err
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
			r.helper.SetRemove(ctx, pendingKey, jobIDStr)
			continue
		}
		job, err := r.FindByID(ctx, jobID)
		if err != nil {
			r.logger.Warn("Failed to get pending job", zap.Error(err), zap.String("job_id", jobIDStr))
			// Remove from pending set if not found
			r.helper.SetRemove(ctx, pendingKey, jobIDStr)
			continue
		}

		// Double-check status (in case it was updated)
		if job.Status() == models.JobStatusPending {
			jobs = append(jobs, job)
		} else {
			// Remove from pending set if status changed
			r.helper.SetRemove(ctx, pendingKey, jobIDStr)
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
