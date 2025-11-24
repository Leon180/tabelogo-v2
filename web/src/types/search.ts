// Search-related types
export interface SearchFilters {
    query: string;
    minRating: number;
    openNow: boolean;
    rankBy: 'relevance' | 'distance';
}

export interface QuickSearchParams {
    place_id: string;
    api_mask?: string;
    language_code: 'en' | 'ja';
}

export interface AdvanceSearchParams {
    text_query: string;
    low_latitude: number;
    low_longitude: number;
    high_latitude: number;
    high_longitude: number;
    max_result_count: number;
    min_rating: number;
    open_now: boolean;
    rank_preference: 'RELEVANCE' | 'DISTANCE';
    language_code: 'en' | 'ja';
    api_mask?: string;
}

export interface MapBounds {
    north: number;
    south: number;
    east: number;
    west: number;
}
