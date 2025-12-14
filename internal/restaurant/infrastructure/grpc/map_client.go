package grpc

import (
	"context"
	"fmt"
	"time"

	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// MapServiceClient wraps the gRPC client for Map Service
type MapServiceClient struct {
	client  mapv1.MapServiceClient
	logger  *zap.Logger
	timeout time.Duration
}

// NewMapServiceClient creates a new Map Service client
func NewMapServiceClient(conn *grpc.ClientConn, logger *zap.Logger, timeout time.Duration) *MapServiceClient {
	return &MapServiceClient{
		client:  mapv1.NewMapServiceClient(conn),
		logger:  logger,
		timeout: timeout,
	}
}

// QuickSearch calls Map Service QuickSearch RPC
func (c *MapServiceClient) QuickSearch(ctx context.Context, placeID string) (*mapv1.Place, error) {
	c.logger.Info("Calling Map Service QuickSearch",
		zap.String("place_id", placeID),
	)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &mapv1.QuickSearchRequest{
		PlaceId:      placeID,
		LanguageCode: "en", // Default to English for area extraction
		// IMPORTANT: Must include addressComponents for area extraction
		ApiMask: "id,displayName,formattedAddress,location,rating,priceLevel,photos,currentOpeningHours,addressComponents",
	}

	resp, err := c.client.QuickSearch(ctx, req)
	if err != nil {
		c.logger.Error("Map Service QuickSearch failed",
			zap.String("place_id", placeID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("map service quick search failed: %w", err)
	}

	if resp.Place == nil {
		c.logger.Warn("Map Service returned no place",
			zap.String("place_id", placeID),
		)
		return nil, fmt.Errorf("place not found: %s", placeID)
	}

	c.logger.Info("Map Service QuickSearch succeeded",
		zap.String("place_id", placeID),
		zap.String("place_name", resp.Place.Name),
	)

	return resp.Place, nil
}

// BatchGetPlaces calls Map Service BatchGetPlaces RPC
func (c *MapServiceClient) BatchGetPlaces(ctx context.Context, placeIDs []string) ([]*mapv1.Place, error) {
	c.logger.Info("Calling Map Service BatchGetPlaces",
		zap.Int("count", len(placeIDs)),
	)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &mapv1.BatchGetPlacesRequest{
		PlaceIds: placeIDs,
	}

	resp, err := c.client.BatchGetPlaces(ctx, req)
	if err != nil {
		c.logger.Error("Map Service BatchGetPlaces failed",
			zap.Int("count", len(placeIDs)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("map service batch get places failed: %w", err)
	}

	c.logger.Info("Map Service BatchGetPlaces succeeded",
		zap.Int("requested", len(placeIDs)),
		zap.Int("returned", len(resp.Places)),
	)

	return resp.Places, nil
}

// AdvanceSearch calls Map Service AdvanceSearch RPC
func (c *MapServiceClient) AdvanceSearch(ctx context.Context, req *mapv1.AdvanceSearchRequest) ([]*mapv1.Place, error) {
	c.logger.Info("Calling Map Service AdvanceSearch",
		zap.String("text_query", req.TextQuery),
	)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.AdvanceSearch(ctx, req)
	if err != nil {
		c.logger.Error("Map Service AdvanceSearch failed",
			zap.String("text_query", req.TextQuery),
			zap.Error(err),
		)
		return nil, fmt.Errorf("map service advance search failed: %w", err)
	}

	c.logger.Info("Map Service AdvanceSearch succeeded",
		zap.String("text_query", req.TextQuery),
		zap.Int("results", len(resp.Places)),
	)

	return resp.Places, nil
}
