package main

import (
	"fmt"
	"log"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/zap"
)

func main() {
	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create scraper config
	config := models.NewScraperConfig()

	// Create scraper
	s := scraper.NewScraper(config, logger)

	// Test cases
	testCases := []struct {
		name      string
		area      string
		placeName string
	}{
		{
			name:      "Test 1: English name",
			area:      "tokyo",
			placeName: "Afuri Ramen",
		},
		{
			name:      "Test 2: Japanese name",
			area:      "tokyo",
			placeName: "阿夫利",
		},
		{
			name:      "Test 3: English - Ichiran",
			area:      "tokyo",
			placeName: "Ichiran Ramen",
		},
		{
			name:      "Test 4: Japanese - Ichiran",
			area:      "tokyo",
			placeName: "一蘭",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n========================================\n")
		fmt.Printf("%s\n", tc.name)
		fmt.Printf("Area: %s, Place Name: %s\n", tc.area, tc.placeName)
		fmt.Printf("========================================\n")

		restaurants, err := s.ScrapeRestaurants(tc.area, tc.placeName)
		if err != nil {
			log.Printf("❌ ERROR: %v\n", err)
			continue
		}

		fmt.Printf("✅ Found %d restaurants\n\n", len(restaurants))

		for i, r := range restaurants {
			if i >= 3 { // Only show first 3 results
				fmt.Printf("... and %d more results\n", len(restaurants)-3)
				break
			}

			fmt.Printf("%d. %s\n", i+1, r.Name())
			fmt.Printf("   Rating: %.2f (%d reviews, %d bookmarks)\n", r.Rating(), r.RatingCount(), r.Bookmarks())
			fmt.Printf("   Types: %v\n", r.Types())
			fmt.Printf("   Link: %s\n", r.Link())
			fmt.Printf("\n")
		}
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("Test Summary\n")
	fmt.Printf("========================================\n")
	fmt.Printf("If English names work: ✅ No need to store Japanese names\n")
	fmt.Printf("If English names fail: ❌ Need to add Japanese name field\n")
}
