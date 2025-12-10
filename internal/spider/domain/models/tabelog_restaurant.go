package models

// TabelogRestaurant represents a restaurant scraped from Tabelog
type TabelogRestaurant struct {
	Link        string   `json:"link"`
	Name        string   `json:"name"`
	Rating      float64  `json:"rating"`
	RatingCount int      `json:"rating_count"`
	Bookmarks   int      `json:"bookmarks"`
	Phone       string   `json:"phone"`
	Types       []string `json:"types"`
	Photos      []string `json:"photos"`
}

// NewTabelogRestaurant creates a new TabelogRestaurant
func NewTabelogRestaurant(
	link, name string,
	rating float64,
	ratingCount, bookmarks int,
	phone string,
	types, photos []string,
) *TabelogRestaurant {
	return &TabelogRestaurant{
		Link:        link,
		Name:        name,
		Rating:      rating,
		RatingCount: ratingCount,
		Bookmarks:   bookmarks,
		Phone:       phone,
		Types:       types,
		Photos:      photos,
	}
}

// AddPhotos adds photos to the restaurant
func (r *TabelogRestaurant) AddPhotos(photos []string) {
	r.Photos = append(r.Photos, photos...)
}
