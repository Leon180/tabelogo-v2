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
    photos: string[];
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
    from_cache?: boolean;
    cached_at?: string;
}

export interface ScrapeJobResponse {
    job_id: string;
    status: string;
}

export interface JobStatusResponse {
    job_id: string;
    google_id: string;
    status: 'PENDING' | 'RUNNING' | 'COMPLETED' | 'FAILED';
    results?: TabelogRestaurant[];
    error?: string;
    created_at: string;
    completed_at?: string;
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
 * Search Tabelog for similar restaurants using SSE for real-time updates
 * 
 * @param params - Search parameters including place name and area
 * @param onProgress - Optional callback for progress updates
 * @returns List of Tabelog restaurants matching the search
 */
export async function searchTabelog(
    params: SearchTabelogRequest,
    onProgress?: (status: string) => void
): Promise<SearchTabelogResponse> {
    // Initiate scrape job
    const response = await spiderClient.post<ScrapeJobResponse | SearchTabelogResponse>(
        '/api/v1/spider/scrape',
        params
    );

    // Check if we got cached results (200 OK with restaurants field)
    if (response.status === 200 && 'restaurants' in response.data) {
        console.log('✅ Got cached results:', response.data.restaurants.length, 'restaurants');
        return response.data as SearchTabelogResponse;
    }

    // Otherwise, we got a job ID (202 Accepted or 200 with job_id)
    const jobResponse = response.data as ScrapeJobResponse;
    const jobId = jobResponse.job_id;

    if (!jobId) {
        throw new ScrapingError('No job ID received from server');
    }

    console.log('🔄 Starting SSE stream for job:', jobId);

    // Subscribe to SSE stream for real-time updates
    return new Promise((resolve, reject) => {
        const eventSource = new EventSource(
            `http://localhost:18084/api/v1/spider/jobs/${jobId}/stream`
        );

        // Handle status updates
        eventSource.addEventListener('status', (event) => {
            const status: JobStatusResponse = JSON.parse(event.data);

            console.log('📡 SSE status update:', status.status);

            // Update progress
            if (onProgress) {
                onProgress(status.status);
            }

            // Handle completion
            if (status.status === 'COMPLETED') {
                eventSource.close();
                console.log('✅ Scraping completed, got', status.results?.length, 'results');
                resolve({
                    google_id: status.google_id,
                    restaurants: status.results || [],
                    total_found: status.results?.length || 0,
                });
            } else if (status.status === 'FAILED') {
                eventSource.close();
                console.error('❌ Scraping failed:', status.error);
                reject(new ScrapingError(status.error || 'Scraping failed'));
            }
        });

        // Handle errors
        eventSource.addEventListener('error', (event: any) => {
            const errorData = event.data ? JSON.parse(event.data) : null;
            eventSource.close();
            console.error('❌ SSE error:', errorData);
            reject(new SpiderServiceError(errorData?.error || 'SSE connection failed'));
        });

        // Handle connection errors
        eventSource.onerror = () => {
            eventSource.close();
            console.error('❌ SSE connection error');
            reject(new SpiderServiceError('Failed to connect to SSE stream'));
        };
    });
}

// ============================================
// Exported Service Object
// ============================================

export const spiderService = {
    searchTabelog,
};

export default spiderService;
