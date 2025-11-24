// Restaurant types based on database schema
export interface Restaurant {
    id: string;
    name: string;
    source: 'google' | 'tabelog' | 'opentable';
    external_id: string;
    address: string;
    latitude: number;
    longitude: number;
    rating: number;
    price_range: string;
    cuisine_type: string;
    phone?: string;
    website?: string;
    opening_hours?: OpeningHours;
    metadata?: Record<string, any>;
    view_count?: number;
    photos?: Photo[];
}

export interface OpeningHours {
    weekday_text?: string[];
    periods?: Period[];
}

export interface Period {
    open: TimeOfDay;
    close: TimeOfDay;
}

export interface TimeOfDay {
    day: number;
    time: string;
}

export interface Photo {
    url: string;
    width: number;
    height: number;
    attribution?: string;
}

export interface UserFavorite {
    id: string;
    user_id: string;
    restaurant_id: string;
    notes?: string;
    tags?: string[];
    visit_count: number;
    last_visited_at?: string;
    created_at: string;
}
