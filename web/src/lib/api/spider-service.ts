import axios from 'axios';

// Spider Service Client Configuration
const spiderClient = axios.create({
    baseURL: 'http://localhost:18084',  // Spider Service port
    timeout: 30000,  // 30s timeout for scraping operations
    headers: {
        'Content-Type': 'application/json',
    },
});

// ============================================
// TypeScript Interfaces
// ============================================

export interface TabelogRestaurant {
    link: string;
    name: string;
    rating: number;
    rating_count: number;
    bookmarks: number;
    phone: string;
    types: string[];
}

export interface SearchTabelogRequest {
    google_id: string;
    area: string;
    place_name: string;
    place_name_ja?: string;
    max_results?: number;
}

export interface SearchTabelogResponse {
    google_id: string;
    restaurants: TabelogRestaurant[];
    total_found: number;
}

// ============================================
// Error Classes
// ============================================

export class SpiderServiceError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'SpiderServiceError';
    }
}

export class ScrapingError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'ScrapingError';
    }
}

// ============================================
// Response Interceptor
// ============================================

spiderClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response) {
            const { status, data } = error.response;

            if (status === 500) {
                throw new ScrapingError(data.message || 'Failed to scrape Tabelog');
            } else if (status === 400) {
                throw new SpiderServiceError(data.message || 'Invalid request parameters');
            }
        } else if (error.request) {
            throw new SpiderServiceError('Unable to connect to Spider Service');
        }

        throw error;
    }
);

// ============================================
// API Functions
// ============================================

/**
 * Search Tabelog for similar restaurants
 * Uses Spider Service to scrape Tabelog website
 * 
 * @param params - Search parameters including place name and area
 * @returns List of Tabelog restaurants matching the search
 */
export async function searchTabelog(
    params: SearchTabelogRequest
): Promise<SearchTabelogResponse> {
    const response = await spiderClient.post<SearchTabelogResponse>(
        '/api/v1/spider/scrape',
        params
    );
    return response.data;
}

// ============================================
// Exported Service Object
// ============================================

export const spiderService = {
    searchTabelog,
};

export default spiderService;
