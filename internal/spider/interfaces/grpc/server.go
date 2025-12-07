package grpc

import (
	"context"

	spiderv1 "github.com/Leon180/tabelogo-v2/api/gen/spider/v1"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SpiderServer implements the Spider Service gRPC server
type SpiderServer struct {
	spiderv1.UnimplementedSpiderServiceServer
	scraper *scraper.Scraper
	logger  *zap.Logger
}

// NewSpiderServer creates a new Spider gRPC server
func NewSpiderServer(
	scraper *scraper.Scraper,
	logger *zap.Logger,
) *SpiderServer {
	return &SpiderServer{
		scraper: scraper,
		logger:  logger.With(zap.String("component", "grpc_server")),
	}
}

// SearchSimilarRestaurants searches Tabelog for similar restaurants
func (s *SpiderServer) SearchSimilarRestaurants(
	ctx context.Context,
	req *spiderv1.SearchSimilarRestaurantsRequest,
) (*spiderv1.SearchSimilarRestaurantsResponse, error) {
	s.logger.Info("SearchSimilarRestaurants called",
		zap.String("google_id", req.GoogleId),
		zap.String("area", req.Area),
		zap.String("place_name", req.PlaceName),
		zap.Int32("max_results", req.MaxResults),
	)

	// Validate request
	if req.Area == "" || req.PlaceName == "" {
		return nil, status.Error(codes.InvalidArgument, "area and place_name are required")
	}

	// Set default max results
	maxResults := req.MaxResults
	if maxResults == 0 {
		maxResults = 10
	} else if maxResults > 20 {
		maxResults = 20
	}

	// Scrape restaurants synchronously
	restaurants, err := s.scraper.ScrapeRestaurants(
		ctx,
		req.GoogleId,
		req.Area,
		req.PlaceName,
	)
	if err != nil {
		s.logger.Error("Failed to scrape restaurants",
			zap.String("google_id", req.GoogleId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to scrape restaurants")
	}

	// Limit results
	totalFound := len(restaurants)
	if len(restaurants) > int(maxResults) {
		restaurants = restaurants[:maxResults]
	}

	// Convert to proto
	protoRestaurants := make([]*spiderv1.TabelogRestaurant, len(restaurants))
	for i, r := range restaurants {
		protoRestaurants[i] = toProtoRestaurant(&r)
	}

	s.logger.Info("SearchSimilarRestaurants succeeded",
		zap.String("google_id", req.GoogleId),
		zap.Int("total_found", totalFound),
		zap.Int("returned", len(protoRestaurants)),
	)

	return &spiderv1.SearchSimilarRestaurantsResponse{
		GoogleId:    req.GoogleId,
		Restaurants: protoRestaurants,
		TotalFound:  int32(totalFound),
	}, nil
}

// GetRestaurantPhotos retrieves photos from a Tabelog restaurant page
func (s *SpiderServer) GetRestaurantPhotos(
	ctx context.Context,
	req *spiderv1.GetRestaurantPhotosRequest,
) (*spiderv1.GetRestaurantPhotosResponse, error) {
	s.logger.Info("GetRestaurantPhotos called",
		zap.String("google_id", req.GoogleId),
		zap.String("link", req.TabelogLink),
		zap.String("name", req.Name),
	)

	// Validate request
	if req.TabelogLink == "" {
		return nil, status.Error(codes.InvalidArgument, "tabelog_link is required")
	}

	// Scrape photos
	photos, err := s.scraper.ScrapePhotos(ctx, req.TabelogLink)
	if err != nil {
		s.logger.Error("Failed to scrape photos",
			zap.String("google_id", req.GoogleId),
			zap.String("link", req.TabelogLink),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to scrape photos")
	}

	s.logger.Info("GetRestaurantPhotos succeeded",
		zap.String("google_id", req.GoogleId),
		zap.Int("photo_count", len(photos)),
	)

	return &spiderv1.GetRestaurantPhotosResponse{
		GoogleId: req.GoogleId,
		Link:     req.TabelogLink,
		Name:     req.Name,
		Photos:   photos,
	}, nil
}

// toProtoRestaurant converts domain model to proto
func toProtoRestaurant(r *models.TabelogRestaurant) *spiderv1.TabelogRestaurant {
	return &spiderv1.TabelogRestaurant{
		Link:        r.Link(),
		Name:        r.Name(),
		Rating:      r.Rating(),
		RatingCount: int32(r.RatingCount()),
		Bookmarks:   int32(r.Bookmarks()),
		Phone:       r.Phone(),
		Types:       r.Types(),
	}
}
