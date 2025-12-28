package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LatencyConfig controls response delay simulation
type LatencyConfig struct {
	Enabled    bool          // Enable latency simulation
	MinLatency time.Duration // Minimum latency
	MaxLatency time.Duration // Maximum latency
}

// Simulate adds a random delay if latency is enabled
func (lc *LatencyConfig) Simulate() {
	if !lc.Enabled {
		return
	}

	// Random delay between min and max
	delay := lc.MinLatency
	if lc.MaxLatency > lc.MinLatency {
		diff := lc.MaxLatency - lc.MinLatency
		delay += time.Duration(rand.Int63n(int64(diff)))
	}

	time.Sleep(delay)
}

// MockMapService provides mock Google Maps API endpoints
type MockMapService struct {
	router   *gin.Engine
	testData *TestDataSet
	latency  *LatencyConfig
}

// TestDataSet holds mock restaurant data
type TestDataSet struct {
	Places []Place `json:"places"`
}

// Place represents a mock place/restaurant
type Place struct {
	ID               string                 `json:"id"`
	DisplayName      map[string]string      `json:"displayName"`
	FormattedAddress string                 `json:"formattedAddress"`
	Location         map[string]float64     `json:"location"`
	Rating           float64                `json:"rating"`
	UserRatingCount  int                    `json:"userRatingCount,omitempty"`
	PriceLevel       string                 `json:"priceLevel,omitempty"`
	Types            []string               `json:"types,omitempty"`
	Photos           []Photo                `json:"photos,omitempty"`
	OpeningHours     map[string]interface{} `json:"regularOpeningHours,omitempty"`
	Phone            string                 `json:"nationalPhoneNumber,omitempty"`
	Website          string                 `json:"websiteUri,omitempty"`
}

// Photo represents a place photo
type Photo struct {
	Name     string `json:"name"`
	WidthPx  int    `json:"widthPx"`
	HeightPx int    `json:"heightPx"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	service := NewMockMapService()

	log.Println("🎭 Mock Map Service starting on :8085")
	log.Println("📍 Endpoints:")
	log.Println("   - GET  /health")
	log.Println("   - POST /v1/places:searchText")
	log.Println("   - GET  /v1/places/:placeId")
	log.Println("   - GET  /v1/:photoName/media")

	if err := service.router.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}

// NewMockMapService creates a new mock service instance
func NewMockMapService() *MockMapService {
	router := gin.Default()

	// Load test data
	testData := loadTestData()

	// Load latency configuration
	latencyConfig := loadLatencyConfig()

	service := &MockMapService{
		router:   router,
		testData: testData,
		latency:  latencyConfig,
	}

	service.setupRoutes()
	return service
}

func (s *MockMapService) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "mock-map-service",
			"version": "1.0.0",
		})
	})

	// Mock Google Places API - Text Search
	s.router.POST("/v1/places:searchText", s.handleTextSearch)

	// Mock Google Places API - Place Details
	s.router.GET("/v1/places/:placeId", s.handlePlaceDetails)

	// Mock Google Places API - Photo
	s.router.GET("/v1/:photoName/media", s.handlePhoto)
}

// handleTextSearch handles text search requests
func (s *MockMapService) handleTextSearch(c *gin.Context) {
	// Simulate latency
	s.latency.Simulate()

	var req struct {
		TextQuery string `json:"textQuery"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	log.Printf("🔍 Text search: %s", req.TextQuery)

	// Search in test data
	places := s.searchPlaces(req.TextQuery)

	c.JSON(200, gin.H{
		"places": places,
	})
}

// handlePlaceDetails handles place details requests
func (s *MockMapService) handlePlaceDetails(c *gin.Context) {
	// Simulate latency
	s.latency.Simulate()

	placeId := c.Param("placeId")

	log.Printf("📍 Place details: %s", placeId)

	place := s.getPlaceById(placeId)
	if place == nil {
		c.JSON(404, gin.H{"error": "place not found"})
		return
	}

	c.JSON(200, place)
}

// handlePhoto returns a simple 1x1 pixel PNG
func (s *MockMapService) handlePhoto(c *gin.Context) {
	photoName := c.Param("photoName")

	log.Printf("📷 Photo request: %s", photoName)

	// Simple 1x1 pixel PNG (transparent)
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	c.Data(200, "image/png", pngData)
}

// searchPlaces searches for places matching the query
func (s *MockMapService) searchPlaces(query string) []Place {
	var results []Place
	queryLower := strings.ToLower(query)

	for _, place := range s.testData.Places {
		displayName := place.DisplayName["text"]
		if strings.Contains(strings.ToLower(displayName), queryLower) {
			results = append(results, place)
		}
	}

	// Return max 20 results
	if len(results) > 20 {
		results = results[:20]
	}

	return results
}

// getPlaceById gets a place by ID
func (s *MockMapService) getPlaceById(id string) *Place {
	for _, place := range s.testData.Places {
		if place.ID == id {
			return &place
		}
	}
	return nil
}

// loadTestData loads test data from file or returns default data
func loadTestData() *TestDataSet {
	// Try to load from file
	data, err := os.ReadFile("testdata/places.json")
	if err != nil {
		log.Println("⚠️  No test data file found, using default data")
		return getDefaultTestData()
	}

	var testData TestDataSet
	if err := json.Unmarshal(data, &testData); err != nil {
		log.Printf("⚠️  Error parsing test data: %v, using default data", err)
		return getDefaultTestData()
	}

	log.Printf("✅ Loaded %d places from test data", len(testData.Places))
	return &testData
}

// getDefaultTestData returns default mock data
func getDefaultTestData() *TestDataSet {
	return &TestDataSet{
		Places: []Place{
			{
				ID: "mock_tokyo_ramen_1",
				DisplayName: map[string]string{
					"text":         "一蘭拉麵 (Mock)",
					"languageCode": "ja",
				},
				FormattedAddress: "東京都新宿區歌舞伎町1-1-1",
				Location: map[string]float64{
					"latitude":  35.6938,
					"longitude": 139.7034,
				},
				Rating:          4.5,
				UserRatingCount: 5000,
				PriceLevel:      "PRICE_LEVEL_MODERATE",
				Types:           []string{"restaurant", "ramen_restaurant", "food"},
				Phone:           "03-1234-5678",
				Website:         "https://example.com/mock-ramen",
				Photos: []Photo{
					{
						Name:     "places/mock_tokyo_ramen_1/photos/photo1",
						WidthPx:  1200,
						HeightPx: 800,
					},
				},
				OpeningHours: map[string]interface{}{
					"openNow": true,
					"weekdayDescriptions": []string{
						"Monday: 11:00 AM – 10:00 PM",
						"Tuesday: 11:00 AM – 10:00 PM",
						"Wednesday: 11:00 AM – 10:00 PM",
						"Thursday: 11:00 AM – 10:00 PM",
						"Friday: 11:00 AM – 11:00 PM",
						"Saturday: 11:00 AM – 11:00 PM",
						"Sunday: 11:00 AM – 9:00 PM",
					},
				},
			},
			{
				ID: "mock_osaka_sushi_1",
				DisplayName: map[string]string{
					"text":         "すしざんまい (Mock)",
					"languageCode": "ja",
				},
				FormattedAddress: "大阪府大阪市中央區道頓堀1-1-1",
				Location: map[string]float64{
					"latitude":  34.6686,
					"longitude": 135.5004,
				},
				Rating:          4.7,
				UserRatingCount: 3000,
				PriceLevel:      "PRICE_LEVEL_EXPENSIVE",
				Types:           []string{"restaurant", "sushi_restaurant", "food"},
				Phone:           "06-1234-5678",
				Website:         "https://example.com/mock-sushi",
				Photos: []Photo{
					{
						Name:     "places/mock_osaka_sushi_1/photos/photo1",
						WidthPx:  1200,
						HeightPx: 800,
					},
				},
			},
			{
				ID: "mock_kyoto_tempura_1",
				DisplayName: map[string]string{
					"text":         "天ぷら京都 (Mock)",
					"languageCode": "ja",
				},
				FormattedAddress: "京都府京都市東山區祇園町1-1-1",
				Location: map[string]float64{
					"latitude":  35.0036,
					"longitude": 135.7681,
				},
				Rating:          4.6,
				UserRatingCount: 2000,
				PriceLevel:      "PRICE_LEVEL_EXPENSIVE",
				Types:           []string{"restaurant", "japanese_restaurant", "food"},
				Phone:           "075-1234-5678",
				Website:         "https://example.com/mock-tempura",
			},
		},
	}
}
