package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.uber.org/zap"
)

// Scraper handles web scraping operations
type Scraper struct {
	config     *models.ScraperConfig
	logger     *zap.Logger
	areaMapper *AreaMapper
}

// NewScraper creates a new scraper
func NewScraper(config *models.ScraperConfig, logger *zap.Logger) *Scraper {
	return &Scraper{
		config:     config,
		logger:     logger.With(zap.String("component", "scraper")),
		areaMapper: NewAreaMapper(),
	}
}

// ScrapeRestaurants scrapes Tabelog for restaurants
func (s *Scraper) ScrapeRestaurants(area, placeName string) ([]models.TabelogRestaurant, error) {
	s.logger.Info("Starting restaurant scrape",
		zap.String("area", area),
		zap.String("place_name", placeName),
	)

	// Step 1: Get restaurant links
	links, err := s.scrapeLinks(area, placeName)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape links: %w", err)
	}

	if len(links) == 0 {
		return nil, fmt.Errorf("no restaurants found")
	}

	s.logger.Info("Found restaurant links",
		zap.Int("count", len(links)),
	)

	// Limit links
	if len(links) > s.config.MaxLinksToCollect() {
		links = links[:s.config.MaxLinksToCollect()]
	}

	// Step 2: Scrape details concurrently
	restaurants := make([]models.TabelogRestaurant, len(links))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := []error{}

	for i, link := range links {
		wg.Add(1)
		go func(index int, url string) {
			defer wg.Done()

			restaurant, err := s.scrapeRestaurantDetails(url)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				s.logger.Error("Failed to scrape restaurant",
					zap.String("url", url),
					zap.Error(err),
				)
				return
			}

			mu.Lock()
			restaurants[index] = *restaurant
			mu.Unlock()
		}(i, link)

		// Rate limiting
		time.Sleep(s.config.DelayBetweenRequests())
	}

	wg.Wait()

	// Filter out empty results
	validRestaurants := []models.TabelogRestaurant{}
	for _, r := range restaurants {
		if r.Name() != "" {
			validRestaurants = append(validRestaurants, r)
		}
	}

	if len(validRestaurants) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all scraping attempts failed: %v", errors[0])
	}

	s.logger.Info("Scraping completed",
		zap.Int("total", len(validRestaurants)),
		zap.Int("errors", len(errors)),
	)

	return validRestaurants, nil
}

// scrapeLinks scrapes restaurant links from search results
func (s *Scraper) scrapeLinks(area, placeName string) ([]string, error) {
	c := s.newCollector()

	links := []string{}
	c.OnHTML(".list-rst__rst-name-target", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			links = append(links, link)
		}
	})

	searchURL := s.buildSearchURL(area, placeName)
	s.logger.Info("Visiting Tabelog search URL",
		zap.String("url", searchURL),
		zap.String("area", area),
		zap.String("place_name", placeName),
	)

	err := c.Visit(searchURL)
	if err != nil {
		s.logger.Error("Failed to visit search URL",
			zap.String("url", searchURL),
			zap.Error(err),
		)
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
func (s *Scraper) buildSearchURL(area, placeName string) string {
	s.logger.Info("🔧 Building Tabelog search URL",
		zap.String("input_area", area),
		zap.String("input_place_name", placeName),
	)

	// Map Google Maps address to Tabelog area code
	// Example: "Meguro, Tokyo" -> "tokyo/A1316"
	// Example: "4-chōme-6-8 Komaba" -> "tokyo" (default if no match)
	tabelogArea := s.areaMapper.MapToTabelogArea(area)

	baseURL := fmt.Sprintf("https://tabelog.com/%s/rstLst/", tabelogArea)
	params := url.Values{}
	params.Add("vs", "1")
	params.Add("sk", placeName)
	params.Add("sw", placeName)

	finalURL := baseURL + "?" + params.Encode()

	s.logger.Info("✅ Tabelog URL constructed",
		zap.String("tabelog_area_code", tabelogArea),
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
