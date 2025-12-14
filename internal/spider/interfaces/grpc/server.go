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
	// Prioritize Japanese name if available, otherwise use English name
	searchName := req.PlaceNameJa
	if searchName == "" {
		searchName = req.PlaceName
	}

	s.logger.Info("Searching Tabelog",
		zap.String("google_id", req.GoogleId),
		zap.String("search_name", searchName),
		zap.String("place_name_ja", req.PlaceNameJa),
		zap.String("place_name", req.PlaceName),
		zap.String("area", req.Area),
	)

	// Scrape Tabelog
	results, err := s.scraper.ScrapeRestaurants(searchName, req.Area)
	if err != nil {
		s.logger.Error("Failed to scrape Tabelog",
			zap.Error(err),
			zap.String("search_name", searchName),
			zap.String("area", req.Area),
		)
		return nil, status.Errorf(codes.Internal, "failed to scrape Tabelog: %v", err)
	}

	// Convert to proto
	protoRestaurants := make([]*spiderv1.TabelogRestaurant, 0, len(results))
	for _, r := range results {
		protoRestaurants = append(protoRestaurants, &spiderv1.TabelogRestaurant{
			Link:        r.Link(),
			Name:        r.Name(),
			Rating:      r.Rating(),
			RatingCount: int32(r.RatingCount()),
			Bookmarks:   int32(r.Bookmarks()),
			Phone:       r.Phone(),
			Types:       r.Types(),
		})
	}

	s.logger.Info("Tabelog search completed",
		zap.String("google_id", req.GoogleId),
		zap.Int("results_count", len(protoRestaurants)),
	)

	return &spiderv1.SearchSimilarRestaurantsResponse{
		GoogleId:    req.GoogleId,
		Restaurants: protoRestaurants,
		TotalFound:  int32(len(protoRestaurants)),
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
	photos, err := s.scraper.ScrapePhotos(req.TabelogLink)
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
