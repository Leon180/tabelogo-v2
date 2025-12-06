package application

import (
	"context"

	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
)

// MapServiceClient defines the interface for Map Service gRPC client
type MapServiceClient interface {
	QuickSearch(ctx context.Context, placeID string) (*mapv1.Place, error)
}
