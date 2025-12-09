package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/application/usecases"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SpiderHandler handles HTTP requests for spider service
type SpiderHandler struct {
	scrapeUseCase       *usecases.ScrapeRestaurantUseCase
	getJobStatusUseCase *usecases.GetJobStatusUseCase
	resultCache         repositories.ResultCacheRepository
	logger              *zap.Logger
}

// NewSpiderHandler creates a new HTTP handler
func NewSpiderHandler(
	scrapeUseCase *usecases.ScrapeRestaurantUseCase,
	getJobStatusUseCase *usecases.GetJobStatusUseCase,
	resultCache repositories.ResultCacheRepository,
	logger *zap.Logger,
) *SpiderHandler {
	return &SpiderHandler{
		scrapeUseCase:       scrapeUseCase,
		getJobStatusUseCase: getJobStatusUseCase,
		resultCache:         resultCache,
		logger:              logger.With(zap.String("component", "http_handler")),
	}
}

// ScrapeRequest is the request body for scraping
type ScrapeRequest struct {
	GoogleID  string `json:"google_id" binding:"required"`
	Area      string `json:"area" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}

// ScrapeResponse is the response for scraping
type ScrapeResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// Scrape handles POST /api/v1/spider/scrape
func (h *SpiderHandler) Scrape(c *gin.Context) {
	var req ScrapeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info("Received scrape request",
		zap.String("google_id", req.GoogleID),
		zap.String("area", req.Area),
		zap.String("place_name", req.PlaceName),
	)

	// Check cache first
	cached, err := h.resultCache.Get(c.Request.Context(), req.GoogleID)
	if err == nil && cached != nil {
		h.logger.Info("Returning cached results",
			zap.String("google_id", req.GoogleID),
			zap.Int("results_count", len(cached.Results)),
		)
		// Return cached results immediately
		results := make([]TabelogRestaurantDTO, len(cached.Results))
		for i, r := range cached.Results {
			results[i] = TabelogRestaurantDTO{
				Link:        r.Link(),
				Name:        r.Name(),
				Rating:      r.Rating(),
				RatingCount: r.RatingCount(),
				Bookmarks:   r.Bookmarks(),
				Phone:       r.Phone(),
				Types:       r.Types(),
				Photos:      r.Photos(),
			}
		}
		c.JSON(http.StatusOK, CachedResultsResponse{
			GoogleID:   req.GoogleID,
			Results:    results,
			TotalFound: len(results),
			FromCache:  true,
			CachedAt:   cached.CachedAt.Format(time.RFC3339),
		})
		return
	}

	// Start scraping job
	resp, err := h.scrapeUseCase.Execute(c.Request.Context(), usecases.ScrapeRestaurantRequest{
		GoogleID:  req.GoogleID,
		Area:      req.Area,
		PlaceName: req.PlaceName,
	})
	if err != nil {
		h.logger.Error("Scrape failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, ScrapeResponse{
		JobID:  resp.JobID,
		Status: resp.Status,
	})
}

// JobStatusResponse is the response for job status
type JobStatusResponse struct {
	JobID       string                 `json:"job_id"`
	GoogleID    string                 `json:"google_id"`
	Status      string                 `json:"status"`
	Results     []TabelogRestaurantDTO `json:"results,omitempty"`
	Error       string                 `json:"error,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	CompletedAt *string                `json:"completed_at,omitempty"`
}

// TabelogRestaurantDTO is the DTO for Tabelog restaurant
type TabelogRestaurantDTO struct {
	Link        string   `json:"link"`
	Name        string   `json:"name"`
	Rating      float64  `json:"rating"`
	RatingCount int      `json:"rating_count"`
	Bookmarks   int      `json:"bookmarks"`
	Phone       string   `json:"phone"`
	Types       []string `json:"types"`
	Photos      []string `json:"photos"`
}

// GetJobStatus handles GET /api/v1/spider/jobs/:job_id
func (h *SpiderHandler) GetJobStatus(c *gin.Context) {
	jobIDStr := c.Param("job_id")

	job, err := h.getJobStatusUseCase.Execute(c.Request.Context(), jobIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	resp := JobStatusResponse{
		JobID:     job.ID().String(),
		GoogleID:  job.GoogleID(),
		Status:    string(job.Status()),
		CreatedAt: job.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	if job.Error() != "" {
		resp.Error = job.Error()
	}

	if job.CompletedAt() != nil {
		completedAt := job.CompletedAt().Format("2006-01-02T15:04:05Z07:00")
		resp.CompletedAt = &completedAt
	}

	if len(job.Results()) > 0 {
		resp.Results = make([]TabelogRestaurantDTO, len(job.Results()))
		for i, r := range job.Results() {
			resp.Results[i] = TabelogRestaurantDTO{
				Link:        r.Link(),
				Name:        r.Name(),
				Rating:      r.Rating(),
				RatingCount: r.RatingCount(),
				Bookmarks:   r.Bookmarks(),
				Phone:       r.Phone(),
				Types:       r.Types(),
				Photos:      r.Photos(),
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

// CachedResultsResponse is the response for cached results
type CachedResultsResponse struct {
	GoogleID   string                 `json:"google_id"`
	Results    []TabelogRestaurantDTO `json:"results"`
	TotalFound int                    `json:"total_found"`
	FromCache  bool                   `json:"from_cache"`
	CachedAt   string                 `json:"cached_at"`
}

// StreamJobStatus handles GET /api/v1/spider/jobs/:job_id/stream
// Streams job status updates via Server-Sent Events (SSE)
func (h *SpiderHandler) StreamJobStatus(c *gin.Context) {
	jobIDStr := c.Param("job_id")

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	h.logger.Info("Starting SSE stream", zap.String("job_id", jobIDStr))

	// Create ticker for polling job status
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Stream updates
	for {
		select {
		case <-ticker.C:
			job, err := h.getJobStatusUseCase.Execute(c.Request.Context(), jobIDStr)
			if err != nil {
				// Send error event
				c.SSEvent("error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
				c.Writer.Flush()
				return
			}

			// Build status response
			status := JobStatusResponse{
				JobID:     job.ID().String(),
				GoogleID:  job.GoogleID(),
				Status:    string(job.Status()),
				CreatedAt: job.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			}

			if job.Error() != "" {
				status.Error = job.Error()
			}

			if job.CompletedAt() != nil {
				completedAt := job.CompletedAt().Format("2006-01-02T15:04:05Z07:00")
				status.CompletedAt = &completedAt
			}

			if len(job.Results()) > 0 {
				status.Results = make([]TabelogRestaurantDTO, len(job.Results()))
				for i, r := range job.Results() {
					status.Results[i] = TabelogRestaurantDTO{
						Link:        r.Link(),
						Name:        r.Name(),
						Rating:      r.Rating(),
						RatingCount: r.RatingCount(),
						Bookmarks:   r.Bookmarks(),
						Phone:       r.Phone(),
						Types:       r.Types(),
						Photos:      r.Photos(),
					}
				}

				// Cache results when completed
				if job.Status() == "completed" {
					if err := h.resultCache.Set(c.Request.Context(), job.GoogleID(), job.Results(), 24*time.Hour); err != nil {
						h.logger.Error("Failed to cache results", zap.Error(err))
					}
				}
			}

			// Send status event
			data, _ := json.Marshal(status)
			c.SSEvent("status", string(data))
			c.Writer.Flush()

			// Close stream if job is done
			if job.Status() == "completed" || job.Status() == "failed" {
				h.logger.Info("Job finished, closing SSE stream",
					zap.String("job_id", jobIDStr),
					zap.String("status", string(job.Status())),
				)
				return
			}

		case <-c.Request.Context().Done():
			h.logger.Info("Client disconnected from SSE stream", zap.String("job_id", jobIDStr))
			return
		}
	}
}

// ErrorResponse is the error response
type ErrorResponse struct {
	Error string `json:"error"`
}
