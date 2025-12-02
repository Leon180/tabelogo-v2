import axios, { AxiosError } from 'axios';
import type {
    QuickSearchRequest,
    QuickSearchResponse,
    AdvanceSearchRequest,
    AdvanceSearchResponse,
    ErrorResponse,
} from '@/types/search';

// API base URL - configured via environment variables
const MAP_SERVICE_URL = process.env.NEXT_PUBLIC_MAP_SERVICE_URL || 'http://localhost:8081';

// Create axios instance with base configuration
const mapClient = axios.create({
    baseURL: MAP_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    timeout: 10000, // 10 second timeout
});

// Add response interceptor for error handling
mapClient.interceptors.response.use(
    (response) => response,
    (error: AxiosError<ErrorResponse>) => {
        if (error.response) {
            // Server responded with error status
            const { status, data } = error.response;

            if (status === 429) {
                // Rate limit exceeded
                console.warn('Rate limit exceeded:', data);
                throw new RateLimitError(data.message, data.retry_after);
            } else if (status === 400) {
                // Bad request
                console.error('Invalid request:', data);
                throw new ValidationError(data.message);
            } else if (status >= 500) {
                // Server error
                console.error('Server error:', data);
                throw new ServerError(data.message);
            }
        } else if (error.request) {
            // Request made but no response
            console.error('Network error:', error.message);
            throw new NetworkError('Unable to connect to Map Service');
        }

        throw error;
    }
);

// ============================================
// Custom Error Classes
// ============================================

export class RateLimitError extends Error {
    constructor(message: string, public retryAfter?: number) {
        super(message);
        this.name = 'RateLimitError';
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
 * Quick Search - Get place details by place_id
 * Used when user clicks on a map marker
 * 
 * @param params - Quick search parameters
 * @returns Place details with source information (google/redis)
 */
export async function quickSearch(params: QuickSearchRequest): Promise<QuickSearchResponse> {
    const response = await mapClient.post<QuickSearchResponse>('/api/v1/map/quick_search', params);
    return response.data;
}

/**
 * Advance Search - Text search with filters and location bias
 * Used for the advance search form
 * 
 * @param params - Advance search parameters
 * @returns List of places matching search criteria
 */
export async function advanceSearch(params: AdvanceSearchRequest): Promise<AdvanceSearchResponse> {
    const response = await mapClient.post<AdvanceSearchResponse>('/api/v1/map/advance_search', params);
    return response.data;
}

// ============================================
// Exported Service Object
// ============================================

export const mapService = {
    quickSearch,
    advanceSearch,
};

export default mapService;
