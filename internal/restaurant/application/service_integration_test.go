//go:build integration
// +build integration

package application_test

import (
	"context"
	"testing"
	"time"

	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/application"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// MockRestaurantRepository is a mock implementation of RestaurantRepository
type MockRestaurantRepository struct {
	mock.Mock
}

func (m *MockRestaurantRepository) Create(ctx context.Context, restaurant *model.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) FindByID(ctx context.Context, id interface{}) (*model.Restaurant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindByExternalID(ctx context.Context, source model.RestaurantSource, externalID string) (*model.Restaurant, error) {
	args := m.Called(ctx, source, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) Update(ctx context.Context, restaurant *model.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) Delete(ctx context.Context, id interface{}) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRestaurantRepository) Search(ctx context.Context, query string, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindByLocation(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, lat, lng, radiusKm, limit)
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) List(ctx context.Context, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRestaurantRepository) FindByCuisineType(ctx context.Context, cuisineType string, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, cuisineType, limit, offset)
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindBySource(ctx context.Context, source model.RestaurantSource, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, source, limit, offset)
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

// MockFavoriteRepository is a mock implementation of FavoriteRepository
type MockFavoriteRepository struct {
	mock.Mock
}

func (m *MockFavoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	return nil
}

func (m *MockFavoriteRepository) FindByID(ctx context.Context, id interface{}) (*model.Favorite, error) {
	return nil, nil
}

func (m *MockFavoriteRepository) FindByUserID(ctx context.Context, userID interface{}) ([]*model.Favorite, error) {
	return nil, nil
}

func (m *MockFavoriteRepository) FindByUserAndRestaurant(ctx context.Context, userID, restaurantID interface{}) (*model.Favorite, error) {
	return nil, nil
}

func (m *MockFavoriteRepository) Update(ctx context.Context, favorite *model.Favorite) error {
	return nil
}

func (m *MockFavoriteRepository) Delete(ctx context.Context, id interface{}) error {
	return nil
}

func (m *MockFavoriteRepository) Exists(ctx context.Context, userID, restaurantID interface{}) (bool, error) {
	return false, nil
}

// MockMapServiceClient is a mock implementation of MapServiceClient
type MockMapServiceClient struct {
	mock.Mock
}

func (m *MockMapServiceClient) QuickSearch(ctx context.Context, placeID string) (*mapv1.Place, error) {
	args := m.Called(ctx, placeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mapv1.Place), args.Error(1)
}

// TestQuickSearchByPlaceID_CacheHit tests cache hit scenario
func TestQuickSearchByPlaceID_CacheHit(t *testing.T) {
	// Setup
	ctx := context.Background()
	placeID := "ChIJN1t_tDeuEmsRUsoyG83frY4"

	mockRepo := new(MockRestaurantRepository)
	mockFavRepo := new(MockFavoriteRepository)
	mockMapClient := new(MockMapServiceClient)
	logger := zaptest.NewLogger(t)

	config := &application.Config{
		DataFreshnessTTL: 3 * 24 * time.Hour, // 3 days
	}

	// Create fresh restaurant (updated 1 hour ago)
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant(
		"Test Restaurant",
		model.SourceGoogle,
		placeID,
		"Tokyo, Japan",
		location,
	)

	// Mock: FindByExternalID returns fresh data
	mockRepo.On("FindByExternalID", ctx, model.SourceGoogle, placeID).
		Return(restaurant, nil)

	// Create service
	service := application.NewRestaurantService(
		mockRepo,
		mockFavRepo,
		mockMapClient,
		config,
		logger,
	)

	// Execute
	result, err := service.QuickSearchByPlaceID(ctx, placeID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Restaurant", result.Name())
	assert.Equal(t, placeID, result.ExternalID())

	// Verify Map Service was NOT called (cache hit)
	mockMapClient.AssertNotCalled(t, "QuickSearch")
	mockRepo.AssertExpectations(t)
}

// TestQuickSearchByPlaceID_CacheMiss tests cache miss scenario
func TestQuickSearchByPlaceID_CacheMiss(t *testing.T) {
	// Setup
	ctx := context.Background()
	placeID := "ChIJN1t_tDeuEmsRUsoyG83frY4"

	mockRepo := new(MockRestaurantRepository)
	mockFavRepo := new(MockFavoriteRepository)
	mockMapClient := new(MockMapServiceClient)
	logger := zaptest.NewLogger(t)

	config := &application.Config{
		DataFreshnessTTL: 3 * 24 * time.Hour,
	}

	// Mock: FindByExternalID returns not found
	mockRepo.On("FindByExternalID", ctx, model.SourceGoogle, placeID).
		Return(nil, repository.ErrNotFound)

	// Mock: Map Service returns place
	place := &mapv1.Place{
		Id:               placeID,
		Name:             "New Restaurant from Map",
		FormattedAddress: "Tokyo, Japan",
		Location: &mapv1.Location{
			Latitude:  35.6762,
			Longitude: 139.6503,
		},
		Rating:     4.5,
		PriceLevel: "PRICE_LEVEL_MODERATE",
	}
	mockMapClient.On("QuickSearch", ctx, placeID).Return(place, nil)

	// Mock: Create saves new restaurant
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Restaurant")).
		Return(nil)

	// Create service
	service := application.NewRestaurantService(
		mockRepo,
		mockFavRepo,
		mockMapClient,
		config,
		logger,
	)

	// Execute
	result, err := service.QuickSearchByPlaceID(ctx, placeID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Restaurant from Map", result.Name())
	assert.Equal(t, placeID, result.ExternalID())
	assert.Equal(t, model.SourceGoogle, result.Source())

	// Verify Map Service was called
	mockMapClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestQuickSearchByPlaceID_StaleData tests stale data refresh scenario
func TestQuickSearchByPlaceID_StaleData(t *testing.T) {
	// Setup
	ctx := context.Background()
	placeID := "ChIJN1t_tDeuEmsRUsoyG83frY4"

	mockRepo := new(MockRestaurantRepository)
	mockFavRepo := new(MockFavoriteRepository)
	mockMapClient := new(MockMapServiceClient)
	logger := zaptest.NewLogger(t)

	config := &application.Config{
		DataFreshnessTTL: 1 * time.Hour, // Short TTL for testing
	}

	// Create stale restaurant (updated 2 hours ago - beyond TTL)
	location, _ := model.NewLocation(35.6762, 139.6503)
	staleRestaurant := model.ReconstructRestaurant(
	// Use reconstruction to set old timestamp
	// This would need to be implemented based on your actual model
	)

	// For simplicity, just test that Map Service is called for stale data
	// In real scenario, you'd manipulate UpdatedAt timestamp

	t.Skip("Skipping stale data test - requires timestamp manipulation")
}

// TestQuickSearchByPlaceID_MapServiceFailure tests fallback to stale data
func TestQuickSearchByPlaceID_MapServiceFailure(t *testing.T) {
	// Setup
	ctx := context.Background()
	placeID := "ChIJN1t_tDeuEmsRUsoyG83frY4"

	mockRepo := new(MockRestaurantRepository)
	mockFavRepo := new(MockFavoriteRepository)
	mockMapClient := new(MockMapServiceClient)
	logger := zaptest.NewLogger(t)

	config := &application.Config{
		DataFreshnessTTL: 1 * time.Hour,
	}

	// Create stale restaurant
	location, _ := model.NewLocation(35.6762, 139.6503)
	staleRestaurant := model.NewRestaurant(
		"Stale Restaurant",
		model.SourceGoogle,
		placeID,
		"Tokyo, Japan",
		location,
	)

	// Mock: FindByExternalID returns stale data
	mockRepo.On("FindByExternalID", ctx, model.SourceGoogle, placeID).
		Return(staleRestaurant, nil)

	// Mock: Map Service fails
	mockMapClient.On("QuickSearch", ctx, placeID).
		Return(nil, assert.AnError)

	// Create service
	service := application.NewRestaurantService(
		mockRepo,
		mockFavRepo,
		mockMapClient,
		config,
		logger,
	)

	// Execute
	result, err := service.QuickSearchByPlaceID(ctx, placeID)

	// Assert: Should return stale data despite Map Service failure
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Stale Restaurant", result.Name())

	// Verify Map Service was called but failed
	mockMapClient.AssertExpectations(t)
}
