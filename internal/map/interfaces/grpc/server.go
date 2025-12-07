package grpc

import (
	"context"

	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
	"github.com/Leon180/tabelogo-v2/internal/map/application/usecases"
	"github.com/Leon180/tabelogo-v2/internal/map/domain/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the MapService gRPC server
type Server struct {
	mapv1.UnimplementedMapServiceServer
	quickSearchUC   *usecases.QuickSearchUseCase
	advanceSearchUC *usecases.AdvanceSearchUseCase
	logger          *zap.Logger
}

// NewServer creates a new gRPC server instance
func NewServer(
	quickSearchUC *usecases.QuickSearchUseCase,
	advanceSearchUC *usecases.AdvanceSearchUseCase,
	logger *zap.Logger,
) *Server {
	return &Server{
		quickSearchUC:   quickSearchUC,
		advanceSearchUC: advanceSearchUC,
		logger:          logger.With(zap.String("component", "grpc_server")),
	}
}

// QuickSearch implements the QuickSearch RPC method
func (s *Server) QuickSearch(
	ctx context.Context,
	req *mapv1.QuickSearchRequest,
) (*mapv1.QuickSearchResponse, error) {
	s.logger.Info("QuickSearch gRPC called",
		zap.String("place_id", req.PlaceId),
		zap.String("language", req.LanguageCode),
	)

	// Validate request
	if req.PlaceId == "" {
		return nil, status.Error(codes.InvalidArgument, "place_id is required")
	}
	if req.LanguageCode == "" {
		return nil, status.Error(codes.InvalidArgument, "language_code is required")
	}

	// Convert gRPC request to domain model
	domainReq := &models.QuickSearchRequest{
		PlaceID:      req.PlaceId,
		LanguageCode: req.LanguageCode,
		APIMask:      req.ApiMask,
	}

	// Execute use case
	result, err := s.quickSearchUC.Execute(ctx, domainReq)
	if err != nil {
		s.logger.Error("QuickSearch failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to search place")
	}

	// Convert domain response to gRPC response
	resp := &mapv1.QuickSearchResponse{
		Source: result.Source,
	}

	if result.CachedAt != nil {
		resp.CachedAt = result.CachedAt.Unix()
	}

	// Convert result to Place proto
	if placeData, ok := result.Result.(map[string]interface{}); ok {
		resp.Place = convertToProtoPlace(placeData)
	}

	return resp, nil
}

// AdvanceSearch implements the AdvanceSearch RPC method
func (s *Server) AdvanceSearch(
	ctx context.Context,
	req *mapv1.AdvanceSearchRequest,
) (*mapv1.AdvanceSearchResponse, error) {
	s.logger.Info("AdvanceSearch gRPC called",
		zap.String("query", req.TextQuery),
		zap.Int32("max_results", req.MaxResultCount),
	)

	// Validate request
	if req.TextQuery == "" {
		return nil, status.Error(codes.InvalidArgument, "text_query is required")
	}
	if req.LocationBias == nil {
		return nil, status.Error(codes.InvalidArgument, "location_bias is required")
	}
	if req.MaxResultCount < 1 || req.MaxResultCount > 20 {
		return nil, status.Error(codes.InvalidArgument, "max_result_count must be between 1 and 20")
	}

	// Convert gRPC request to domain model
	domainReq := &models.AdvanceSearchRequest{
		TextQuery: req.TextQuery,
		LocationBias: models.LocationBias{
			Rectangle: models.Rectangle{
				Low: models.Coordinates{
					Latitude:  req.LocationBias.Rectangle.Low.Latitude,
					Longitude: req.LocationBias.Rectangle.Low.Longitude,
				},
				High: models.Coordinates{
					Latitude:  req.LocationBias.Rectangle.High.Latitude,
					Longitude: req.LocationBias.Rectangle.High.Longitude,
				},
			},
		},
		MaxResultCount: int(req.MaxResultCount),
		MinRating:      req.MinRating,
		OpenNow:        req.OpenNow,
		RankPreference: req.RankPreference,
		LanguageCode:   req.LanguageCode,
		APIMask:        req.ApiMask,
	}

	// Execute use case
	result, err := s.advanceSearchUC.Execute(ctx, domainReq)
	if err != nil {
		s.logger.Error("AdvanceSearch failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to search places")
	}

	// Convert domain response to gRPC response
	resp := &mapv1.AdvanceSearchResponse{
		TotalCount: int32(result.TotalCount),
		Metadata: &mapv1.SearchMetadata{
			TextQuery:    result.SearchMetadata.TextQuery,
			SearchTimeMs: result.SearchMetadata.SearchTimeMs,
		},
	}

	// Convert places
	for _, placeData := range result.Places {
		if placeMap, ok := placeData.(map[string]interface{}); ok {
			resp.Places = append(resp.Places, convertToProtoPlace(placeMap))
		}
	}

	return resp, nil
}

// BatchGetPlaces implements the BatchGetPlaces RPC method
func (s *Server) BatchGetPlaces(
	ctx context.Context,
	req *mapv1.BatchGetPlacesRequest,
) (*mapv1.BatchGetPlacesResponse, error) {
	s.logger.Info("BatchGetPlaces gRPC called",
		zap.Int("place_count", len(req.PlaceIds)),
	)

	// Validate request
	if len(req.PlaceIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "place_ids cannot be empty")
	}
	if len(req.PlaceIds) > 50 {
		return nil, status.Error(codes.InvalidArgument, "maximum 50 place_ids allowed")
	}

	resp := &mapv1.BatchGetPlacesResponse{
		Places: make([]*mapv1.Place, 0, len(req.PlaceIds)),
	}

	// Fetch each place
	for _, placeID := range req.PlaceIds {
		domainReq := &models.QuickSearchRequest{
			PlaceID:      placeID,
			LanguageCode: req.LanguageCode,
			APIMask:      req.ApiMask,
		}

		result, err := s.quickSearchUC.Execute(ctx, domainReq)
		if err != nil {
			s.logger.Warn("Failed to fetch place in batch",
				zap.String("place_id", placeID),
				zap.Error(err),
			)
			resp.FailedCount++
			continue
		}

		// Convert to proto place
		if placeData, ok := result.Result.(map[string]interface{}); ok {
			resp.Places = append(resp.Places, convertToProtoPlace(placeData))
			resp.SuccessCount++
		} else {
			resp.FailedCount++
		}
	}

	s.logger.Info("BatchGetPlaces completed",
		zap.Int32("success", resp.SuccessCount),
		zap.Int32("failed", resp.FailedCount),
	)

	return resp, nil
}

// convertToProtoPlace converts a map[string]interface{} to a proto Place
func convertToProtoPlace(data map[string]interface{}) *mapv1.Place {
	place := &mapv1.Place{}

	// Extract basic fields
	if id, ok := data["id"].(string); ok {
		place.Id = id
	}
	if name, ok := data["displayName"].(map[string]interface{}); ok {
		if text, ok := name["text"].(string); ok {
			place.Name = text
		}
	}
	if addr, ok := data["formattedAddress"].(string); ok {
		place.FormattedAddress = addr
	}

	// Extract location
	if loc, ok := data["location"].(map[string]interface{}); ok {
		place.Location = &mapv1.Location{
			Latitude:  getFloat64(loc, "latitude"),
			Longitude: getFloat64(loc, "longitude"),
		}
	}

	// Extract rating
	if rating, ok := data["rating"].(float64); ok {
		place.Rating = rating
	}
	if userRatingsTotal, ok := data["userRatingCount"].(float64); ok {
		place.UserRatingsTotal = int32(userRatingsTotal)
	}

	// Extract business status
	if status, ok := data["businessStatus"].(string); ok {
		place.BusinessStatus = status
	}

	// Extract phone number
	if phone, ok := data["internationalPhoneNumber"].(string); ok {
		place.PhoneNumber = phone
	}

	// Extract website
	if website, ok := data["websiteUri"].(string); ok {
		place.Website = website
	}

	// Extract types
	if types, ok := data["types"].([]interface{}); ok {
		for _, t := range types {
			if typeStr, ok := t.(string); ok {
				place.Types = append(place.Types, typeStr)
			}
		}
	}

	// Extract opening hours
	// Try currentOpeningHours first (for places that are currently open/closed)
	if hours, ok := data["currentOpeningHours"].(map[string]interface{}); ok {
		place.OpeningHours = &mapv1.OpeningHours{}
		if openNow, ok := hours["openNow"].(bool); ok {
			place.OpeningHours.OpenNow = openNow
		}
		if weekdayText, ok := hours["weekdayDescriptions"].([]interface{}); ok {
			for _, day := range weekdayText {
				if dayStr, ok := day.(string); ok {
					place.OpeningHours.WeekdayText = append(place.OpeningHours.WeekdayText, dayStr)
				}
			}
		}
	} else if hours, ok := data["regularOpeningHours"].(map[string]interface{}); ok {
		// Fallback to regularOpeningHours if currentOpeningHours is not available
		place.OpeningHours = &mapv1.OpeningHours{}
		if openNow, ok := hours["openNow"].(bool); ok {
			place.OpeningHours.OpenNow = openNow
		}
		if weekdayText, ok := hours["weekdayDescriptions"].([]interface{}); ok {
			for _, day := range weekdayText {
				if dayStr, ok := day.(string); ok {
					place.OpeningHours.WeekdayText = append(place.OpeningHours.WeekdayText, dayStr)
				}
			}
		}
	}

	// Extract photos
	if photos, ok := data["photos"].([]interface{}); ok {
		for _, p := range photos {
			if photoMap, ok := p.(map[string]interface{}); ok {
				photo := &mapv1.Photo{}
				if name, ok := photoMap["name"].(string); ok {
					photo.Name = name
				}
				if width, ok := photoMap["widthPx"].(float64); ok {
					photo.WidthPx = int32(width)
				}
				if height, ok := photoMap["heightPx"].(float64); ok {
					photo.HeightPx = int32(height)
				}
				place.Photos = append(place.Photos, photo)
			}
		}
	}

	// Extract editorial summary
	if summary, ok := data["editorialSummary"].(map[string]interface{}); ok {
		if text, ok := summary["text"].(string); ok {
			place.EditorialSummary = text
		}
	}

	return place
}

// getFloat64 safely extracts a float64 from a map
func getFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}
