package models

// TabelogRestaurantDTO is the data transfer object for TabelogRestaurant
// Used for JSON serialization/deserialization (e.g., Redis cache)
type TabelogRestaurantDTO struct {
	Link        string   `json:"link"`
	Name        string   `json:"name"`
	Rating      float64  `json:"rating"`
	RatingCount int      `json:"rating_count"`
	Bookmarks   int      `json:"bookmarks"`
	Phone       string   `json:"phone"`
	Types       []string `json:"types"`
	Photos      []string `json:"photos"`
}

// ToDTO converts TabelogRestaurant domain model to DTO
func (r *TabelogRestaurant) ToDTO() TabelogRestaurantDTO {
	return TabelogRestaurantDTO{
		Link:        r.link,
		Name:        r.name,
		Rating:      r.rating,
		RatingCount: r.ratingCount,
		Bookmarks:   r.bookmarks,
		Phone:       r.phone,
		Types:       r.types,
		Photos:      r.photos,
	}
}

// ToDomain converts DTO to TabelogRestaurant domain model
func (dto TabelogRestaurantDTO) ToDomain() *TabelogRestaurant {
	return &TabelogRestaurant{
		link:        dto.Link,
		name:        dto.Name,
		rating:      dto.Rating,
		ratingCount: dto.RatingCount,
		bookmarks:   dto.Bookmarks,
		phone:       dto.Phone,
		types:       dto.Types,
		photos:      dto.Photos,
	}
}
