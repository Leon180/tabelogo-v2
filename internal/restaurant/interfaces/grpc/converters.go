package grpc

import (
	restaurantv1 "github.com/Leon180/tabelogo-v2/api/gen/restaurant/v1"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Domain to Proto converters

func toProtoLocation(loc *model.Location) *restaurantv1.Location {
	if loc == nil {
		return nil
	}
	return &restaurantv1.Location{
		Latitude:  loc.Latitude(),
		Longitude: loc.Longitude(),
	}
}

func toProtoRestaurant(r *model.Restaurant) *restaurantv1.Restaurant {
	if r == nil {
		return nil
	}

	return &restaurantv1.Restaurant{
		Id:           r.ID().String(),
		Name:         r.Name(),
		Source:       string(r.Source()),
		ExternalId:   r.ExternalID(),
		Address:      r.Address(),
		Location:     toProtoLocation(r.Location()),
		Rating:       r.Rating(),
		PriceRange:   r.PriceRange(),
		CuisineType:  r.CuisineType(),
		Phone:        r.Phone(),
		Website:      r.Website(),
		OpeningHours: r.OpeningHours(),
		ViewCount:    r.ViewCount(),
		CreatedAt:    timestamppb.New(r.CreatedAt()),
		UpdatedAt:    timestamppb.New(r.UpdatedAt()),
	}
}

func toProtoFavorite(f *model.Favorite) *restaurantv1.Favorite {
	if f == nil {
		return nil
	}

	var lastVisitedAt *timestamppb.Timestamp
	if f.LastVisitedAt() != nil {
		lastVisitedAt = timestamppb.New(*f.LastVisitedAt())
	}

	return &restaurantv1.Favorite{
		Id:            f.ID().String(),
		UserId:        f.UserID().String(),
		RestaurantId:  f.RestaurantID().String(),
		Notes:         f.Notes(),
		Tags:          f.Tags(),
		VisitCount:    int32(f.VisitCount()),
		LastVisitedAt: lastVisitedAt,
		CreatedAt:     timestamppb.New(f.CreatedAt()),
	}
}

func toProtoRestaurants(restaurants []*model.Restaurant) []*restaurantv1.Restaurant {
	result := make([]*restaurantv1.Restaurant, len(restaurants))
	for i, r := range restaurants {
		result[i] = toProtoRestaurant(r)
	}
	return result
}

func toProtoFavorites(favorites []*model.Favorite) []*restaurantv1.Favorite {
	result := make([]*restaurantv1.Favorite, len(favorites))
	for i, f := range favorites {
		result[i] = toProtoFavorite(f)
	}
	return result
}

// Proto to Domain converters

func fromProtoLocation(loc *restaurantv1.Location) (*model.Location, error) {
	if loc == nil {
		return nil, nil
	}
	return model.NewLocation(loc.Latitude, loc.Longitude)
}
