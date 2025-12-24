package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/alicebob/miniredis/v2"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func setupTestJobStore(t *testing.T) (*RedisJobStore, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)

	client := redisclient.NewClient(&redisclient.Options{
		Addr: mr.Addr(),
	})

	store := NewRedisJobStore(client, zap.NewNop()).(*RedisJobStore)

	return store, mr
}

func TestRedisJobStore_SaveAndFind(t *testing.T) {
	store, cleanup := setupTestJobStore(t)
	defer cleanup.Close()

	ctx := context.Background()

	// Create a job
	job := models.NewScrapingJob("test-google-id", "Tokyo", "Test Restaurant")

	// Save job
	err := store.Save(ctx, job)
	if err != nil {
		t.Fatalf("Failed to save job: %v", err)
	}

	// Find by ID
	found, err := store.FindByID(ctx, job.ID())
	if err != nil {
		t.Fatalf("Failed to find job: %v", err)
	}

	if found.ID().String() != job.ID().String() {
		t.Errorf("Expected job ID %s, got %s", job.ID().String(), found.ID().String())
	}

	if found.GoogleID() != "test-google-id" {
		t.Errorf("Expected Google ID test-google-id, got %s", found.GoogleID())
	}
}

func TestRedisJobStore_FindByGoogleID(t *testing.T) {
	store, mr := setupTestJobStore(t)
	defer mr.Close()

	ctx := context.Background()
	googleID := "test-google-id"

	// Create multiple jobs with same Google ID
	job1 := models.NewScrapingJob(googleID, "Tokyo", "Restaurant 1")
	job2 := models.NewScrapingJob(googleID, "Osaka", "Restaurant 2")

	store.Save(ctx, job1)
	store.Save(ctx, job2)

	// Find by Google ID
	jobs, err := store.FindByGoogleID(ctx, googleID)
	if err != nil {
		t.Fatalf("Failed to find jobs: %v", err)
	}

	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}
}

func TestRedisJobStore_Update(t *testing.T) {
	store, mr := setupTestJobStore(t)
	defer mr.Close()

	ctx := context.Background()

	// Create and save job
	job := models.NewScrapingJob("test-google-id", "Tokyo", "Test Restaurant")
	store.Save(ctx, job)

	// Update job status
	job.Start()
	err := store.Update(ctx, job)
	if err != nil {
		t.Fatalf("Failed to update job: %v", err)
	}

	// Verify update
	found, _ := store.FindByID(ctx, job.ID())
	if found.Status() != models.JobStatusRunning {
		t.Errorf("Expected status RUNNING, got %s", found.Status())
	}
}

func TestRedisJobStore_Delete(t *testing.T) {
	store, mr := setupTestJobStore(t)
	defer mr.Close()

	ctx := context.Background()

	// Create and save job
	job := models.NewScrapingJob("test-google-id", "Tokyo", "Test Restaurant")
	store.Save(ctx, job)

	// Delete job
	err := store.Delete(ctx, job.ID())
	if err != nil {
		t.Fatalf("Failed to delete job: %v", err)
	}

	// Verify deletion
	_, err = store.FindByID(ctx, job.ID())
	if err == nil {
		t.Error("Expected error when finding deleted job")
	}
}

func TestRedisJobStore_FindPending(t *testing.T) {
	store, mr := setupTestJobStore(t)
	defer mr.Close()

	ctx := context.Background()

	// Create pending jobs
	job1 := models.NewScrapingJob("google-1", "Tokyo", "Restaurant 1")
	job2 := models.NewScrapingJob("google-2", "Osaka", "Restaurant 2")
	job3 := models.NewScrapingJob("google-3", "Kyoto", "Restaurant 3")

	store.Save(ctx, job1)
	store.Save(ctx, job2)
	store.Save(ctx, job3)

	// Start one job
	job2.Start()
	store.Update(ctx, job2)

	// Find pending jobs
	pending, err := store.FindPending(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to find pending jobs: %v", err)
	}

	// Should find 2 pending jobs (job1 and job3)
	if len(pending) != 2 {
		t.Errorf("Expected 2 pending jobs, got %d", len(pending))
	}

	// Verify all are pending
	for _, job := range pending {
		if job.Status() != models.JobStatusPending {
			t.Errorf("Expected pending status, got %s", job.Status())
		}
	}
}

func TestRedisJobStore_TTL(t *testing.T) {
	store, mr := setupTestJobStore(t)
	defer mr.Close()

	ctx := context.Background()

	// Create and save job
	job := models.NewScrapingJob("test-google-id", "Tokyo", "Test Restaurant")
	store.Save(ctx, job)

	// Check TTL is set
	key := "spider:jobs:" + job.ID().String()
	ttl := store.client.TTL(ctx, key).Val()

	if ttl <= 0 {
		t.Error("Expected positive TTL")
	}

	// TTL should be approximately 24 hours
	expectedTTL := 24 * time.Hour
	if ttl < expectedTTL-time.Minute || ttl > expectedTTL+time.Minute {
		t.Errorf("Expected TTL around %v, got %v", expectedTTL, ttl)
	}
}
