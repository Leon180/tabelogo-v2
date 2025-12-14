package application

import (
	"context"
	"testing"
	"time"

	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock Repository
type MockRestaurantRepository struct {
	mock.Mock
}

func (m *MockRestaurantRepository) Create(ctx context.Context, restaurant *model.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Restaurant, error) {
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

func (m *MockRestaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRestaurantRepository) Search(ctx context.Context, query string, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindByLocation(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, lat, lng, radiusKm, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) List(ctx context.Context, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRestaurantRepository) FindByCuisineType(ctx context.Context, cuisineType string, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, cuisineType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindBySource(ctx context.Context, source model.RestaurantSource, limit, offset int) ([]*model.Restaurant, error) {
	args := m.Called(ctx, source, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Restaurant), args.Error(1)
}

// Mock Favorite Repository
type MockFavoriteRepository struct {
	mock.Mock
}

func (m *MockFavoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	args := m.Called(ctx, favorite)
	return args.Error(0)
}

func (m *MockFavoriteRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Favorite, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Favorite), args.Error(1)
}

func (m *MockFavoriteRepository) FindByUserAndRestaurant(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error) {
	args := m.Called(ctx, userID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Favorite), args.Error(1)
}

func (m *MockFavoriteRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Favorite, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Favorite), args.Error(1)
}

func (m *MockFavoriteRepository) FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*model.Favorite, error) {
	args := m.Called(ctx, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Favorite), args.Error(1)
}

func (m *MockFavoriteRepository) Update(ctx context.Context, favorite *model.Favorite) error {
	args := m.Called(ctx, favorite)
	return args.Error(0)
}

func (m *MockFavoriteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFavoriteRepository) Exists(ctx context.Context, userID, restaurantID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, restaurantID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFavoriteRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFavoriteRepository) FindByTag(ctx context.Context, userID uuid.UUID, tag string) ([]*model.Favorite, error) {
	args := m.Called(ctx, userID, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Favorite), args.Error(1)
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

// Test CreateRestaurant
func TestRestaurantService_CreateRestaurant_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	req := CreateRestaurantRequest{
		Name:        "Sushi Dai",
		Source:      model.SourceGoogle,
		ExternalID:  "ChIJTest123",
		Address:     "Tokyo",
		Latitude:    35.6762,
		Longitude:   139.6503,
		Rating:      4.5,
		PriceRange:  "$$",
		CuisineType: "Sushi",
		Phone:       "03-1234-5678",
		Website:     "https://example.com",
	}

	mockRestaurantRepo.On("FindByExternalID", ctx, model.SourceGoogle, "ChIJTest123").
		Return(nil, domainerrors.ErrRestaurantNotFound)
	mockRestaurantRepo.On("Create", ctx, mock.AnythingOfType("*model.Restaurant")).
		Return(nil)

	restaurant, err := service.CreateRestaurant(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, restaurant)
	assert.Equal(t, "Sushi Dai", restaurant.Name())
	assert.Equal(t, model.SourceGoogle, restaurant.Source())
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_CreateRestaurant_DuplicateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	req := CreateRestaurantRequest{
		Name:       "Sushi Dai",
		Source:     model.SourceGoogle,
		ExternalID: "ChIJTest123",
		Address:    "Tokyo",
		Latitude:   35.6762,
		Longitude:  139.6503,
	}

	location, _ := model.NewLocation(35.6762, 139.6503)
	existingRestaurant := model.NewRestaurant("Existing", "Existing", model.SourceGoogle, "ChIJTest123", "Tokyo", location)

	mockRestaurantRepo.On("FindByExternalID", ctx, model.SourceGoogle, "ChIJTest123").
		Return(existingRestaurant, nil)

	restaurant, err := service.CreateRestaurant(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantAlreadyExists, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_CreateRestaurant_InvalidLocation(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	req := CreateRestaurantRequest{
		Name:       "Test",
		Source:     model.SourceGoogle,
		ExternalID: "test",
		Address:    "Test",
		Latitude:   999.0, // Invalid
		Longitude:  139.6503,
	}

	mockRestaurantRepo.On("FindByExternalID", ctx, model.SourceGoogle, "test").
		Return(nil, domainerrors.ErrRestaurantNotFound)

	restaurant, err := service.CreateRestaurant(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrInvalidLocation, err)
	assert.Nil(t, restaurant)
}

// Test GetRestaurant
func TestRestaurantService_GetRestaurant_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	expectedRestaurant := model.NewRestaurant("Sushi Dai", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).
		Return(expectedRestaurant, nil)

	restaurant, err := service.GetRestaurant(ctx, restaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, restaurant)
	assert.Equal(t, expectedRestaurant.Name(), restaurant.Name())
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_GetRestaurant_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).
		Return(nil, domainerrors.ErrRestaurantNotFound)

	restaurant, err := service.GetRestaurant(ctx, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantNotFound, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

// Test AddToFavorites
func TestRestaurantService_AddToFavorites_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(restaurant, nil)
	mockFavoriteRepo.On("Exists", ctx, userID, restaurantID).Return(false, nil)
	mockFavoriteRepo.On("Create", ctx, mock.AnythingOfType("*model.Favorite")).Return(nil)

	favorite, err := service.AddToFavorites(ctx, userID, restaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, favorite)
	assert.Equal(t, userID, favorite.UserID())
	assert.Equal(t, restaurantID, favorite.RestaurantID())
	mockRestaurantRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddToFavorites_AlreadyExists(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(restaurant, nil)
	mockFavoriteRepo.On("Exists", ctx, userID, restaurantID).Return(true, nil)

	favorite, err := service.AddToFavorites(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrFavoriteAlreadyExists, err)
	assert.Nil(t, favorite)
	mockRestaurantRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddToFavorites_RestaurantNotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).
		Return(nil, domainerrors.ErrRestaurantNotFound)

	favorite, err := service.AddToFavorites(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantNotFound, err)
	assert.Nil(t, favorite)
	mockRestaurantRepo.AssertExpectations(t)
}

// Test SearchRestaurants
func TestRestaurantService_SearchRestaurants_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant1 := model.NewRestaurant("Sushi Dai", "Tokyo", model.SourceGoogle, "test1", "Tokyo", location)
	restaurant2 := model.NewRestaurant("Sushi Saito", "Tokyo", model.SourceTabelog, "test2", "Tokyo", location)
	expectedRestaurants := []*model.Restaurant{restaurant1, restaurant2}

	mockRestaurantRepo.On("Search", ctx, "sushi", 20, 0).
		Return(expectedRestaurants, nil)

	restaurants, err := service.SearchRestaurants(ctx, "sushi", 20, 0)

	assert.NoError(t, err)
	assert.Len(t, restaurants, 2)
	assert.Equal(t, "Sushi Dai", restaurants[0].Name())
	assert.Equal(t, "Sushi Saito", restaurants[1].Name())
	mockRestaurantRepo.AssertExpectations(t)
}

// Test IsFavorite
func TestRestaurantService_IsFavorite(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()

	mockFavoriteRepo.On("Exists", ctx, userID, restaurantID).Return(true, nil)

	isFavorite, err := service.IsFavorite(ctx, userID, restaurantID)

	assert.NoError(t, err)
	assert.True(t, isFavorite)
	mockFavoriteRepo.AssertExpectations(t)
}

// Test GetRestaurantByExternalID
func TestRestaurantService_GetRestaurantByExternalID_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	location, _ := model.NewLocation(35.6762, 139.6503)
	expectedRestaurant := model.NewRestaurant("Sushi Dai", "Tokyo", model.SourceGoogle, "ChIJTest123", "Tokyo", location)

	mockRestaurantRepo.On("FindByExternalID", ctx, model.SourceGoogle, "ChIJTest123").
		Return(expectedRestaurant, nil)

	restaurant, err := service.GetRestaurantByExternalID(ctx, model.SourceGoogle, "ChIJTest123")

	assert.NoError(t, err)
	assert.NotNil(t, restaurant)
	assert.Equal(t, expectedRestaurant.Name(), restaurant.Name())
	assert.Equal(t, "ChIJTest123", restaurant.ExternalID())
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_GetRestaurantByExternalID_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()

	mockRestaurantRepo.On("FindByExternalID", ctx, model.SourceTabelog, "nonexistent").
		Return(nil, domainerrors.ErrRestaurantNotFound)

	restaurant, err := service.GetRestaurantByExternalID(ctx, model.SourceTabelog, "nonexistent")

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantNotFound, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

// Test UpdateRestaurant
func TestRestaurantService_UpdateRestaurant_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	existingRestaurant := model.NewRestaurant("Old Name", "Old Area", model.SourceGoogle, "test", "Old Address", location)

	req := UpdateRestaurantRequest{
		Name:        "New Name",
		Address:     "New Address",
		Rating:      4.8,
		PriceRange:  "$$$",
		CuisineType: "Japanese",
		Phone:       "03-9999-8888",
		Website:     "https://newsite.com",
		Latitude:    35.7000,
		Longitude:   139.7000,
	}

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(existingRestaurant, nil)
	mockRestaurantRepo.On("Update", ctx, mock.AnythingOfType("*model.Restaurant")).Return(nil)

	restaurant, err := service.UpdateRestaurant(ctx, restaurantID, req)

	assert.NoError(t, err)
	assert.NotNil(t, restaurant)
	assert.Equal(t, "New Name", restaurant.Name())
	assert.Equal(t, "New Address", restaurant.Address())
	assert.Equal(t, 4.8, restaurant.Rating())
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_UpdateRestaurant_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	req := UpdateRestaurantRequest{Name: "New Name"}

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).
		Return(nil, domainerrors.ErrRestaurantNotFound)

	restaurant, err := service.UpdateRestaurant(ctx, restaurantID, req)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantNotFound, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_UpdateRestaurant_InvalidLocation(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	existingRestaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)

	req := UpdateRestaurantRequest{
		Latitude:  999.0, // Invalid
		Longitude: 139.6503,
	}

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(existingRestaurant, nil)

	restaurant, err := service.UpdateRestaurant(ctx, restaurantID, req)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrInvalidLocation, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

// Test DeleteRestaurant
func TestRestaurantService_DeleteRestaurant_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()

	mockRestaurantRepo.On("Delete", ctx, restaurantID).Return(nil)

	err := service.DeleteRestaurant(ctx, restaurantID)

	assert.NoError(t, err)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_DeleteRestaurant_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()

	mockRestaurantRepo.On("Delete", ctx, restaurantID).
		Return(domainerrors.ErrRestaurantNotFound)

	err := service.DeleteRestaurant(ctx, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantNotFound, err)
	mockRestaurantRepo.AssertExpectations(t)
}

// Test ListRestaurants
func TestRestaurantService_ListRestaurants_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant1 := model.NewRestaurant("Restaurant 1", "Tokyo", model.SourceGoogle, "test1", "Tokyo", location)
	restaurant2 := model.NewRestaurant("Restaurant 2", "Osaka", model.SourceTabelog, "test2", "Osaka", location)
	expectedRestaurants := []*model.Restaurant{restaurant1, restaurant2}

	mockRestaurantRepo.On("List", ctx, 20, 0).Return(expectedRestaurants, nil)

	restaurants, err := service.ListRestaurants(ctx, 20, 0)

	assert.NoError(t, err)
	assert.Len(t, restaurants, 2)
	assert.Equal(t, "Restaurant 1", restaurants[0].Name())
	assert.Equal(t, "Restaurant 2", restaurants[1].Name())
	mockRestaurantRepo.AssertExpectations(t)
}

// Test FindRestaurantsByLocation
func TestRestaurantService_FindRestaurantsByLocation_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	location1, _ := model.NewLocation(35.6762, 139.6503)
	restaurant1 := model.NewRestaurant("Nearby Restaurant", "Tokyo", model.SourceGoogle, "test1", "Tokyo", location1)
	expectedRestaurants := []*model.Restaurant{restaurant1}

	mockRestaurantRepo.On("FindByLocation", ctx, 35.6762, 139.6503, 5.0, 10).
		Return(expectedRestaurants, nil)

	restaurants, err := service.FindRestaurantsByLocation(ctx, 35.6762, 139.6503, 5.0, 10)

	assert.NoError(t, err)
	assert.Len(t, restaurants, 1)
	assert.Equal(t, "Nearby Restaurant", restaurants[0].Name())
	mockRestaurantRepo.AssertExpectations(t)
}

// Test FindRestaurantsByCuisineType
func TestRestaurantService_FindRestaurantsByCuisineType_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant1 := model.NewRestaurant("Sushi Place", "Tokyo", model.SourceGoogle, "test1", "Tokyo", location)
	restaurant1.UpdateDetails("", "", "", "Japanese", "", "")
	expectedRestaurants := []*model.Restaurant{restaurant1}

	mockRestaurantRepo.On("FindByCuisineType", ctx, "Japanese", 20, 0).
		Return(expectedRestaurants, nil)

	restaurants, err := service.FindRestaurantsByCuisineType(ctx, "Japanese", 20, 0)

	assert.NoError(t, err)
	assert.Len(t, restaurants, 1)
	assert.Equal(t, "Sushi Place", restaurants[0].Name())
	mockRestaurantRepo.AssertExpectations(t)
}

// Test IncrementRestaurantViewCount
func TestRestaurantService_IncrementRestaurantViewCount_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(restaurant, nil)
	mockRestaurantRepo.On("Update", ctx, mock.AnythingOfType("*model.Restaurant")).Return(nil)

	err := service.IncrementRestaurantViewCount(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), restaurant.ViewCount())
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_IncrementRestaurantViewCount_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).
		Return(nil, domainerrors.ErrRestaurantNotFound)

	err := service.IncrementRestaurantViewCount(ctx, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrRestaurantNotFound, err)
	mockRestaurantRepo.AssertExpectations(t)
}

// Test RemoveFromFavorites
func TestRestaurantService_RemoveFromFavorites_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favoriteID := uuid.New()
	favorite := model.ReconstructFavorite(favoriteID, userID, restaurantID, "", nil, 0, nil,
		time.Now(), time.Now(), nil)

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Delete", ctx, favoriteID).Return(nil)

	err := service.RemoveFromFavorites(ctx, userID, restaurantID)

	assert.NoError(t, err)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_RemoveFromFavorites_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).
		Return(nil, domainerrors.ErrFavoriteNotFound)

	err := service.RemoveFromFavorites(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrFavoriteNotFound, err)
	mockFavoriteRepo.AssertExpectations(t)
}

// Test GetUserFavorites
func TestRestaurantService_GetUserFavorites_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	favorite1 := model.NewFavorite(userID, uuid.New())
	favorite2 := model.NewFavorite(userID, uuid.New())
	expectedFavorites := []*model.Favorite{favorite1, favorite2}

	mockFavoriteRepo.On("FindByUserID", ctx, userID).Return(expectedFavorites, nil)

	favorites, err := service.GetUserFavorites(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, favorites, 2)
	assert.Equal(t, userID, favorites[0].UserID())
	assert.Equal(t, userID, favorites[1].UserID())
	mockFavoriteRepo.AssertExpectations(t)
}

// Test GetFavoriteByUserAndRestaurant
func TestRestaurantService_GetFavoriteByUserAndRestaurant_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	expectedFavorite := model.NewFavorite(userID, restaurantID)

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).
		Return(expectedFavorite, nil)

	favorite, err := service.GetFavoriteByUserAndRestaurant(ctx, userID, restaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, favorite)
	assert.Equal(t, userID, favorite.UserID())
	assert.Equal(t, restaurantID, favorite.RestaurantID())
	mockFavoriteRepo.AssertExpectations(t)
}

// Test UpdateFavoriteNotes
func TestRestaurantService_UpdateFavoriteNotes_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(nil)

	updatedFavorite, err := service.UpdateFavoriteNotes(ctx, userID, restaurantID, "Great sushi!")

	assert.NoError(t, err)
	assert.NotNil(t, updatedFavorite)
	assert.Equal(t, "Great sushi!", updatedFavorite.Notes())
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_UpdateFavoriteNotes_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).
		Return(nil, domainerrors.ErrFavoriteNotFound)

	favorite, err := service.UpdateFavoriteNotes(ctx, userID, restaurantID, "notes")

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrFavoriteNotFound, err)
	assert.Nil(t, favorite)
	mockFavoriteRepo.AssertExpectations(t)
}

// Test AddFavoriteTag
func TestRestaurantService_AddFavoriteTag_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(nil)

	updatedFavorite, err := service.AddFavoriteTag(ctx, userID, restaurantID, "sushi")

	assert.NoError(t, err)
	assert.NotNil(t, updatedFavorite)
	assert.True(t, updatedFavorite.HasTag("sushi"))
	mockFavoriteRepo.AssertExpectations(t)
}

// Test RemoveFavoriteTag
func TestRestaurantService_RemoveFavoriteTag_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)
	favorite.AddTag("sushi")
	favorite.AddTag("expensive")

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(nil)

	updatedFavorite, err := service.RemoveFavoriteTag(ctx, userID, restaurantID, "sushi")

	assert.NoError(t, err)
	assert.NotNil(t, updatedFavorite)
	assert.False(t, updatedFavorite.HasTag("sushi"))
	assert.True(t, updatedFavorite.HasTag("expensive"))
	mockFavoriteRepo.AssertExpectations(t)
}

// Test AddFavoriteVisit
func TestRestaurantService_AddFavoriteVisit_Success(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(nil)

	updatedFavorite, err := service.AddFavoriteVisit(ctx, userID, restaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedFavorite)
	assert.Equal(t, 1, updatedFavorite.VisitCount())
	assert.NotNil(t, updatedFavorite.LastVisitedAt())
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddFavoriteVisit_NotFound(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).
		Return(nil, domainerrors.ErrFavoriteNotFound)

	favorite, err := service.AddFavoriteVisit(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, domainerrors.ErrFavoriteNotFound, err)
	assert.Nil(t, favorite)
	mockFavoriteRepo.AssertExpectations(t)
}

// Test error scenarios for better coverage

func TestRestaurantService_SearchRestaurants_Error(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	expectedErr := assert.AnError

	mockRestaurantRepo.On("Search", ctx, "sushi", 20, 0).Return(nil, expectedErr)

	restaurants, err := service.SearchRestaurants(ctx, "sushi", 20, 0)

	assert.Error(t, err)
	assert.Nil(t, restaurants)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_ListRestaurants_Error(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	expectedErr := assert.AnError

	mockRestaurantRepo.On("List", ctx, 20, 0).Return(nil, expectedErr)

	restaurants, err := service.ListRestaurants(ctx, 20, 0)

	assert.Error(t, err)
	assert.Nil(t, restaurants)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_FindRestaurantsByLocation_Error(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByLocation", ctx, 35.6762, 139.6503, 5.0, 10).
		Return(nil, expectedErr)

	restaurants, err := service.FindRestaurantsByLocation(ctx, 35.6762, 139.6503, 5.0, 10)

	assert.Error(t, err)
	assert.Nil(t, restaurants)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_FindRestaurantsByCuisineType_Error(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByCuisineType", ctx, "Japanese", 20, 0).
		Return(nil, expectedErr)

	restaurants, err := service.FindRestaurantsByCuisineType(ctx, "Japanese", 20, 0)

	assert.Error(t, err)
	assert.Nil(t, restaurants)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_GetUserFavorites_Error(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := assert.AnError

	mockFavoriteRepo.On("FindByUserID", ctx, userID).Return(nil, expectedErr)

	favorites, err := service.GetUserFavorites(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, favorites)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddToFavorites_ExistsError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(restaurant, nil)
	mockFavoriteRepo.On("Exists", ctx, userID, restaurantID).Return(false, expectedErr)

	favorite, err := service.AddToFavorites(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, favorite)
	mockRestaurantRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddToFavorites_CreateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(restaurant, nil)
	mockFavoriteRepo.On("Exists", ctx, userID, restaurantID).Return(false, nil)
	mockFavoriteRepo.On("Create", ctx, mock.AnythingOfType("*model.Favorite")).Return(expectedErr)

	favorite, err := service.AddToFavorites(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, favorite)
	mockRestaurantRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_RemoveFromFavorites_DeleteError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favoriteID := uuid.New()
	favorite := model.ReconstructFavorite(favoriteID, userID, restaurantID, "", nil, 0, nil,
		time.Now(), time.Now(), nil)
	expectedErr := assert.AnError

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Delete", ctx, favoriteID).Return(expectedErr)

	err := service.RemoveFromFavorites(ctx, userID, restaurantID)

	assert.Error(t, err)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_UpdateFavoriteNotes_UpdateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)
	expectedErr := assert.AnError

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(expectedErr)

	updatedFavorite, err := service.UpdateFavoriteNotes(ctx, userID, restaurantID, "Great sushi!")

	assert.Error(t, err)
	assert.Nil(t, updatedFavorite)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddFavoriteTag_UpdateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)
	expectedErr := assert.AnError

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(expectedErr)

	updatedFavorite, err := service.AddFavoriteTag(ctx, userID, restaurantID, "sushi")

	assert.Error(t, err)
	assert.Nil(t, updatedFavorite)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_RemoveFavoriteTag_UpdateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)
	favorite.AddTag("sushi")
	expectedErr := assert.AnError

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(expectedErr)

	updatedFavorite, err := service.RemoveFavoriteTag(ctx, userID, restaurantID, "sushi")

	assert.Error(t, err)
	assert.Nil(t, updatedFavorite)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_AddFavoriteVisit_UpdateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	favorite := model.NewFavorite(userID, restaurantID)
	expectedErr := assert.AnError

	mockFavoriteRepo.On("FindByUserAndRestaurant", ctx, userID, restaurantID).Return(favorite, nil)
	mockFavoriteRepo.On("Update", ctx, mock.AnythingOfType("*model.Favorite")).Return(expectedErr)

	updatedFavorite, err := service.AddFavoriteVisit(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, updatedFavorite)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_IsFavorite_Error(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	userID := uuid.New()
	restaurantID := uuid.New()
	expectedErr := assert.AnError

	mockFavoriteRepo.On("Exists", ctx, userID, restaurantID).Return(false, expectedErr)

	isFavorite, err := service.IsFavorite(ctx, userID, restaurantID)

	assert.Error(t, err)
	assert.False(t, isFavorite)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestRestaurantService_CreateRestaurant_CreateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	req := CreateRestaurantRequest{
		Name:       "Sushi Dai",
		Source:     model.SourceGoogle,
		ExternalID: "ChIJTest123",
		Address:    "Tokyo",
		Latitude:   35.6762,
		Longitude:  139.6503,
	}
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByExternalID", ctx, model.SourceGoogle, "ChIJTest123").
		Return(nil, domainerrors.ErrRestaurantNotFound)
	mockRestaurantRepo.On("Create", ctx, mock.AnythingOfType("*model.Restaurant")).
		Return(expectedErr)

	restaurant, err := service.CreateRestaurant(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_UpdateRestaurant_UpdateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	existingRestaurant := model.NewRestaurant("Old Name", "Old Area", model.SourceGoogle, "test", "Old Address", location)
	req := UpdateRestaurantRequest{Name: "New Name"}
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(existingRestaurant, nil)
	mockRestaurantRepo.On("Update", ctx, mock.AnythingOfType("*model.Restaurant")).Return(expectedErr)

	restaurant, err := service.UpdateRestaurant(ctx, restaurantID, req)

	assert.Error(t, err)
	assert.Nil(t, restaurant)
	mockRestaurantRepo.AssertExpectations(t)
}

func TestRestaurantService_IncrementRestaurantViewCount_UpdateError(t *testing.T) {
	mockRestaurantRepo := new(MockRestaurantRepository)
	mockFavoriteRepo := new(MockFavoriteRepository)
	logger := zap.NewNop()
	mockMapClient := new(MockMapServiceClient)
	config := &Config{DataFreshnessTTL: 3 * 24 * time.Hour}
	service := NewRestaurantService(mockRestaurantRepo, mockFavoriteRepo, mockMapClient, config, logger)

	ctx := context.Background()
	restaurantID := uuid.New()
	location, _ := model.NewLocation(35.6762, 139.6503)
	restaurant := model.NewRestaurant("Test", "Tokyo", model.SourceGoogle, "test", "Tokyo", location)
	expectedErr := assert.AnError

	mockRestaurantRepo.On("FindByID", ctx, restaurantID).Return(restaurant, nil)
	mockRestaurantRepo.On("Update", ctx, mock.AnythingOfType("*model.Restaurant")).Return(expectedErr)

	err := service.IncrementRestaurantViewCount(ctx, restaurantID)

	assert.Error(t, err)
	mockRestaurantRepo.AssertExpectations(t)
}
