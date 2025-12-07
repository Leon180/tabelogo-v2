import axios, { AxiosError } from 'axios';

// ============================================
// Types
// ============================================

export interface Restaurant {
    id: string;
    name: string;
    name_ja?: string;
    source: string;
    external_id: string;
    address: string;
    latitude: number;
    longitude: number;
    rating: number;
    price_range: string;
    cuisine_type: string;
    phone: string;
    website: string;
    view_count: number;
    created_at: string;
    updated_at: string;
}

export interface RestaurantResponse {
    restaurant: Restaurant;
}

export interface ErrorResponse {
    error: string;
    message: string;
}

// ============================================
// Configuration
// ============================================

const RESTAURANT_SERVICE_URL = process.env.NEXT_PUBLIC_RESTAURANT_SERVICE_URL || 'http://localhost:18082';

// Create axios instance
const restaurantClient = axios.create({
    baseURL: RESTAURANT_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    timeout: 10000,
});

// Response interceptor for error handling
restaurantClient.interceptors.response.use(
    (response) => response,
    (error: AxiosError<ErrorResponse>) => {
        if (error.response) {
            const { status, data } = error.response;

            if (status === 404) {
                console.warn('Restaurant not found:', data);
                throw new NotFoundError(data.message);
            } else if (status === 400) {
                console.error('Invalid request:', data);
                throw new ValidationError(data.message);
            } else if (status >= 500) {
                console.error('Server error:', data);
                throw new ServerError(data.message);
            }
        } else if (error.request) {
            console.error('Network error:', error.message);
            throw new NetworkError('Unable to connect to Restaurant Service');
        }

        throw error;
    }
);

// ============================================
// Custom Error Classes
// ============================================

export class NotFoundError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'NotFoundError';
    }
}

export class ValidationError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'ValidationError';
    }
}

export class ServerError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'ServerError';
    }
}

export class NetworkError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'NetworkError';
    }
}

// ============================================
// API Functions
// ============================================

/**
 * Quick Search by Place ID - Get restaurant details with cache-first strategy
 * 
 * Benefits:
 * - 80% cost reduction (cached responses)
 * - Faster response times (<50ms for cache hits)
 * - Automatic fallback to Map Service if needed
 * 
 * @param placeId - Google Place ID
 * @returns Restaurant details with cache metadata
 */
export async function quickSearchByPlaceId(placeId: string): Promise<RestaurantResponse> {
    const response = await restaurantClient.get<RestaurantResponse>(
        `/api/v1/restaurants/quick-search/${placeId}`
    );

    // Log cache status for monitoring
    const cacheStatus = response.headers['x-cache-status'];
    const dataSource = response.headers['x-data-source'];
    const dataAge = response.headers['x-data-age'];

    if (process.env.NODE_ENV === 'development') {
        console.log(`[Restaurant Service] Cache: ${cacheStatus}, Source: ${dataSource}, Age: ${dataAge}`);
    }

    return response.data;
}

// ============================================
// Exported Service Object
// ============================================

export const restaurantService = {
    quickSearchByPlaceId,
};

export default restaurantService;
