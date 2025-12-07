// Map Service API Types - matching backend implementation

// ============================================
// Quick Search Types
// ============================================

export interface QuickSearchRequest {
    place_id: string;
    language_code: 'en' | 'ja' | 'zh-TW';
    api_mask?: string;
}

export interface QuickSearchResponse {
    source: 'google' | 'redis';
    cached_at?: string;
    result: Place;
}

// ============================================
// Advance Search Types
// ============================================

export interface AdvanceSearchRequest {
    text_query: string;
    location_bias: LocationBias;
    max_result_count: number;
    min_rating?: number;
    open_now?: boolean;
    rank_preference: 'DISTANCE' | 'RELEVANCE';
    language_code: 'en' | 'ja' | 'zh-TW';
    api_mask?: string;
}

export interface LocationBias {
    rectangle: Rectangle;
}

export interface Rectangle {
    low: LatLng;
    high: LatLng;
}

export interface LatLng {
    latitude: number;
    longitude: number;
}

export interface AdvanceSearchResponse {
    places: Place[];
    total_count: number;
    search_metadata: SearchMetadata;
}

export interface SearchMetadata {
    text_query: string;
    search_time_ms: number;
}

// ============================================
// Place Model (from Google Places API)
// ============================================

export interface Place {
    id: string;
    displayName?: DisplayName;
    formattedAddress?: string;
    location?: Location;
    rating?: number;
    userRatingCount?: number;
    priceLevel?: string;
    websiteUri?: string;
    regularOpeningHours?: OpeningHours;
    currentOpeningHours?: OpeningHours;
    photos?: Photo[];
    reviews?: Review[];
    types?: string[];
    // Google Maps addressComponents for extracting area information
    addressComponents?: Array<{
        longText?: string;
        shortText?: string;
        types: string[];
    }>;
    nationalPhoneNumber?: string;
    internationalPhoneNumber?: string;
    googleMapsUri?: string;
}

export interface DisplayName {
    text: string;
    languageCode: string;
}

export interface Location {
    latitude: number;
    longitude: number;
}

export interface OpeningHours {
    openNow?: boolean;
    weekdayDescriptions?: string[];
}

export interface Photo {
    name: string;
    widthPx: number;
    heightPx: number;
    authorAttributions?: AuthorAttribution[];
}

export interface AuthorAttribution {
    displayName: string;
    uri?: string;
    photoUri?: string;
}

export interface Review {
    name: string;
    relativePublishTimeDescription?: string;
    rating?: number;
    text?: TextContent;
    originalText?: TextContent;
    authorAttribution?: AuthorAttribution;
    publishTime?: string;
}

export interface TextContent {
    text: string;
    languageCode: string;
}

// ============================================
// Error Response
// ============================================

export interface ErrorResponse {
    error: string;
    message: string;
    retry_after?: number; // For 429 rate limit errors
}

// ============================================
// Utility Types
// ============================================

export interface MapBounds {
    north: number;
    south: number;
    east: number;
    west: number;
}

// ============================================
// Legacy types for backward compatibility
// ============================================

export type QuickSearchParams = QuickSearchRequest;
export type AdvanceSearchParams = AdvanceSearchRequest;
