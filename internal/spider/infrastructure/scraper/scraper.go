package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// Scraper handles web scraping operations
type Scraper struct {
	config         *models.ScraperConfig
	logger         *zap.Logger
	circuitBreaker *gobreaker.CircuitBreaker
	metrics        *metrics.SpiderMetrics
}

// NewScraper creates a new scraper
func NewScraper(config *models.ScraperConfig, metrics *metrics.SpiderMetrics, logger *zap.Logger) *Scraper {
	scraperLogger := logger.With(zap.String("component", "scraper"))
	return &Scraper{
		config:         config,
		logger:         scraperLogger,
		circuitBreaker: NewCircuitBreaker(scraperLogger, metrics, DefaultCircuitBreakerConfig()),
		metrics:        metrics,
	}
}

// ScrapeRestaurants scrapes Tabelog for restaurants
func (s *Scraper) ScrapeRestaurants(area, placeName string) ([]models.TabelogRestaurant, error) {
	// Track scrape duration
	startTime := time.Now()
	defer func() {
		s.metrics.RecordScrapeDuration("search", time.Since(startTime).Seconds())
	}()

	s.logger.Info("Starting restaurant scrape",
		zap.String("area", area),
		zap.String("place_name", placeName),
	)

	// Step 1: Scrape links
	links, err := s.scrapeLinks(area, placeName)
	if err != nil {
		s.metrics.RecordScrapeError("search_failed")
		return nil, fmt.Errorf("failed to scrape links: %w", err)
	}

	if len(links) == 0 {
		s.logger.Warn("No restaurant links found",
			zap.String("area", area),
			zap.String("place_name", placeName),
		)
		s.metrics.RecordScrapeError("no_results")
		return []models.TabelogRestaurant{}, nil
	}

	s.logger.Info("Found restaurant links",
		zap.Int("count", len(links)),
	)

	// Step 2: Scrape details for each link (with concurrency)
	var wg sync.WaitGroup
	resultsChan := make(chan models.TabelogRestaurant, len(links))
	errorsChan := make(chan error, len(links))

	for _, link := range links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Track detail scrape duration
			detailStart := time.Now()
			restaurant, err := s.scrapeRestaurantDetails(url)
			s.metrics.RecordScrapeDuration("details", time.Since(detailStart).Seconds())

			if err != nil {
				s.logger.Warn("Failed to scrape restaurant details",
					zap.String("url", url),
					zap.Error(err),
				)
				s.metrics.RecordScrapeError("details_failed")
				errorsChan <- err
				return
			}

			if restaurant != nil {
				resultsChan <- *restaurant
			}
		}(link)
	}

	wg.Wait()
	close(resultsChan)
	close(errorsChan)

	// Collect results
	var validRestaurants []models.TabelogRestaurant
	for restaurant := range resultsChan {
		validRestaurants = append(validRestaurants, restaurant)
	}

	// Collect errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(validRestaurants) == 0 && len(errors) > 0 {
		s.metrics.RecordScrapeError("all_failed")
		return nil, fmt.Errorf("all scraping attempts failed: %v", errors[0])
	}

	s.logger.Info("Scraping completed",
		zap.Int("total", len(validRestaurants)),
		zap.Int("errors", len(errors)),
	)

	return validRestaurants, nil
}

// scrapeLinks scrapes restaurant links from Tabelog search
func (s *Scraper) scrapeLinks(area, placeName string) ([]string, error) {
	var links []string
	var scrapeErr error

	// Execute with circuit breaker protection
	_, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		// Build search URL
		searchURL := s.buildSearchURL(area, placeName)

		s.logger.Info("Visiting Tabelog search URL",
			zap.String("url", searchURL),
			zap.String("area", area),
			zap.String("place_name", placeName),
		)

		// Create collector
		c := s.newCollector()

		// Collect restaurant links
		c.OnHTML("a.list-rst__rst-name-target", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			if link != "" {
				// Ensure absolute URL
				if !strings.HasPrefix(link, "http") {
					link = "https://tabelog.com" + link
				}
				links = append(links, link)
			}
		})

		// Handle errors
		c.OnError(func(r *colly.Response, err error) {
			s.logger.Error("Failed to visit search URL",
				zap.String("url", searchURL),
				zap.Error(err),
			)
			scrapeErr = err
		})

		// Visit the search page
		if err := c.Visit(searchURL); err != nil {
			return nil, err
		}

		c.Wait()

		if scrapeErr != nil {
			return nil, scrapeErr
		}

		return links, nil
	})
	if err != nil {
		// Check if it's a circuit breaker error
		if IsCircuitBreakerError(err) {
			s.logger.Warn("Circuit breaker is open, rejecting request",
				zap.String("area", area),
				zap.String("place_name", placeName),
			)
			return nil, fmt.Errorf("service temporarily unavailable (circuit breaker open): %w", err)
		}
		return nil, err
	}

	// Remove duplicates
	return removeDuplicates(links), nil
}

// scrapeRestaurantDetails scrapes details for a single restaurant
func (s *Scraper) scrapeRestaurantDetails(link string) (*models.TabelogRestaurant, error) {
	c := s.newCollector()

	data := make(map[string][]string)

	// Scrape basic info
	c.OnHTML("#container", func(e *colly.HTMLElement) {
		data["name"] = []string{e.ChildText("h2.display-name")}
		data["rating"] = []string{e.ChildText(".rdheader-rating__score b.c-rating__val")}
		data["ratingCount"] = []string{e.ChildText(".rdheader-rating__review-target .num")}
		data["bookmarks"] = []string{e.ChildText(".rdheader-rating__hozon-target .num")}
		data["phone"] = []string{e.ChildText(".rstinfo-table__tel-num")}
	})

	// Scrape types
	types := []string{}
	c.OnHTML(".rdheader-subinfo__item", func(e *colly.HTMLElement) {
		if e.ChildText(".rdheader-subinfo__item-title") == "ジャンル：" {
			e.ForEach(".linktree__parent-target-text", func(_ int, el *colly.HTMLElement) {
				types = append(types, strings.TrimSpace(el.Text))
			})
		}
	})

	err := c.Visit(link)
	if err != nil {
		s.logger.Error("Failed to visit URL", zap.String("url", link), zap.Error(err))
		return nil, err
	}

	// Parse data
	name := getFirst(data["name"])
	rating := parseFloat(getFirst(data["rating"]))
	ratingCount := parseInt(getFirst(data["ratingCount"]))
	bookmarks := parseInt(getFirst(data["bookmarks"]))
	phone := getFirst(data["phone"])

	return models.NewTabelogRestaurant(
		link,
		name,
		rating,
		ratingCount,
		bookmarks,
		phone,
		types,
		[]string{}, // Photos will be scraped separately
	), nil
}

// ScrapePhotos scrapes photos for a restaurant
func (s *Scraper) ScrapePhotos(link string) ([]string, error) {
	c := s.newCollector()

	photos := []string{}
	c.OnHTML(".rstdtl-photo-list__item", func(e *colly.HTMLElement) {
		photo := e.ChildAttr(".rstdtl-photo-list__img", "src")
		if photo != "" {
			photos = append(photos, photo)
		}
	})

	photoURL := link + "dtlphotolst"
	err := c.Visit(photoURL)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

// newCollector creates a new colly collector with configuration
func (s *Scraper) newCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.Async(false),
	)

	// Random User-Agent
	extensions.RandomUserAgent(c)

	// Set timeout
	c.SetRequestTimeout(s.config.Timeout())

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       s.config.DelayBetweenRequests(),
		RandomDelay: 200 * time.Millisecond,
	})

	return c
}

// buildSearchURL builds the Tabelog search URL
// Following v1: area is administrative_area_level_1 (e.g., "Tokyo") converted to lowercase
func (s *Scraper) buildSearchURL(area, placeName string) string {
	s.logger.Info("🔧 Building Tabelog search URL",
		zap.String("input_area", area),
		zap.String("input_place_name", placeName),
	)

	// V1 approach: just lowercase the area
	// Example: "Tokyo" -> "tokyo"
	tabelogArea := strings.ToLower(strings.TrimSpace(area))

	baseURL := fmt.Sprintf("https://tabelog.com/%s/rstLst/", tabelogArea)
	params := url.Values{}
	params.Add("vs", "1")
	params.Add("sk", placeName)
	params.Add("sw", placeName)

	finalURL := baseURL + "?" + params.Encode()

	s.logger.Info("✅ Tabelog URL constructed",
		zap.String("tabelog_area", tabelogArea),
		zap.String("final_url", finalURL),
	)

	return finalURL
}

// Helper functions

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func getFirst(slice []string) string {
	if len(slice) > 0 {
		return strings.TrimSpace(slice[0])
	}
	return ""
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
