package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SSEHandler handles Server-Sent Events for job status streaming
type SSEHandler struct {
	jobRepo repositories.JobRepository
	logger  *zap.Logger
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(jobRepo repositories.JobRepository, logger *zap.Logger) *SSEHandler {
	return &SSEHandler{
		jobRepo: jobRepo,
		logger:  logger.With(zap.String("component", "sse_handler")),
	}
}

// StreamJobStatus streams job status updates via SSE
// GET /api/v1/spider/jobs/:job_id/stream
func (h *SSEHandler) StreamJobStatus(c *gin.Context) {
	jobIDStr := c.Param("job_id")
	jobID, err := models.ParseJobID(jobIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid job_id"})
		return
	}

	h.logger.Info("Starting SSE stream", zap.String("job_id", jobID.String()))

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	// Create a channel for job updates
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		h.logger.Error("Streaming not supported")
		c.JSON(500, gin.H{"error": "streaming not supported"})
		return
	}

	lastStatus := ""
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("SSE stream context done", zap.String("job_id", jobID.String()))
			return
		case <-c.Request.Context().Done():
			h.logger.Info("Client disconnected", zap.String("job_id", jobID.String()))
			return
		case <-ticker.C:
			// Fetch job status
			job, err := h.jobRepo.FindByID(ctx, jobID)
			if err != nil {
				h.sendEvent(c, "error", map[string]interface{}{
					"error": "job not found",
				})
				flusher.Flush()
				return
			}

			// Only send update if status changed
			currentStatus := string(job.Status())
			if currentStatus != lastStatus {
				h.sendJobUpdate(c, job)
				flusher.Flush()
				lastStatus = currentStatus
			}

			// If job is completed, send final update and close
			if job.IsCompleted() {
				h.sendEvent(c, "done", map[string]interface{}{
					"message": "job completed",
				})
				flusher.Flush()
				h.logger.Info("Job completed, closing SSE stream",
					zap.String("job_id", jobID.String()),
					zap.String("status", currentStatus),
				)
				return
			}
		}
	}
}

// sendJobUpdate sends a job status update event
func (h *SSEHandler) sendJobUpdate(c *gin.Context, job *models.ScrapingJob) {
	data := map[string]interface{}{
		"job_id":     job.ID().String(),
		"google_id":  job.GoogleID(),
		"status":     string(job.Status()),
		"created_at": job.CreatedAt().Format(time.RFC3339),
	}

	if job.StartedAt() != nil {
		data["started_at"] = job.StartedAt().Format(time.RFC3339)
	}

	if job.CompletedAt() != nil {
		data["completed_at"] = job.CompletedAt().Format(time.RFC3339)
		data["duration"] = job.Duration().Seconds()
	}

	if job.Status() == models.JobStatusCompleted {
		data["results_count"] = len(job.Results())
	}

	if job.Status() == models.JobStatusFailed {
		data["error"] = job.Error()
	}

	h.sendEvent(c, "update", data)
}

// sendEvent sends an SSE event
func (h *SSEHandler) sendEvent(c *gin.Context, event string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal SSE data", zap.Error(err))
		return
	}

	fmt.Fprintf(c.Writer, "event: %s\n", event)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(jsonData))
}
