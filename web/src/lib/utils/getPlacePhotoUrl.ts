/**
 * Generate Google Places Photo URL from photo name
 * @param photoName - Photo name from Google Places API (e.g., "places/ChIJ...")
 * @param maxWidth - Maximum width in pixels (default: 400)
 * @param maxHeight - Maximum height in pixels (default: 300)
 * @returns Full URL to fetch the photo
 */
export function getPlacePhotoUrl(
    photoName: string,
    maxWidth: number = 400,
    maxHeight: number = 300
): string {
    const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY;

    if (!apiKey) {
        console.error('Google Maps API key not configured');
        return '';
    }

    // Google Places Photo API endpoint
    return `https://places.googleapis.com/v1/${photoName}/media?key=${apiKey}&maxWidthPx=${maxWidth}&maxHeightPx=${maxHeight}`;
}

/**
 * Get placeholder image URL for places without photos
 */
export function getPlaceholderImageUrl(): string {
    // Using a simple gradient placeholder
    return 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="400" height="300"%3E%3Crect width="400" height="300" fill="%23374151"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-family="sans-serif" font-size="18" fill="%239CA3AF"%3ENo Image%3C/text%3E%3C/svg%3E';
}
