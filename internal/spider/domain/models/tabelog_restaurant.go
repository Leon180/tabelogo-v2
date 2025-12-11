package models

// TabelogRestaurant represents a restaurant scraped from Tabelog
type TabelogRestaurant struct {
	link        string
	name        string
	rating      float64
	ratingCount int
	bookmarks   int
	phone       string
	types       []string
	photos      []string
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
		link:        link,
		name:        name,
		rating:      rating,
		ratingCount: ratingCount,
		bookmarks:   bookmarks,
		phone:       phone,
		types:       types,
		photos:      photos,
	}
}

// Link returns the Tabelog URL
func (r *TabelogRestaurant) Link() string {
	return r.link
}

// Name returns the restaurant name
func (r *TabelogRestaurant) Name() string {
	return r.name
}

// Rating returns the Tabelog rating
func (r *TabelogRestaurant) Rating() float64 {
	return r.rating
}

// RatingCount returns the number of ratings
func (r *TabelogRestaurant) RatingCount() int {
	return r.ratingCount
}

// Bookmarks returns the bookmark count
func (r *TabelogRestaurant) Bookmarks() int {
	return r.bookmarks
}

// Phone returns the phone number
func (r *TabelogRestaurant) Phone() string {
	return r.phone
}

// Types returns the restaurant types/categories
func (r *TabelogRestaurant) Types() []string {
	return r.types
}

// Photos returns the photo URLs
func (r *TabelogRestaurant) Photos() []string {
	return r.photos
}

// AddPhotos adds photos to the restaurant
func (r *TabelogRestaurant) AddPhotos(photos []string) {
	r.photos = append(r.photos, photos...)
}
