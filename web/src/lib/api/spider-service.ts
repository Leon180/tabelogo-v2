import axios from 'axios';

// Spider Service Client Configuration
const spiderClient = axios.create({
    baseURL: 'http://localhost:18084',  // Spider Service port
    timeout: 30000,  // 30s timeout for scraping operations
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add Authorization header to all requests
spiderClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
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

            if (status === 401) {
                // Authentication failed - token invalid or expired
                throw new SpiderServiceError('Authentication required. Please login again.');
            } else if (status === 500) {
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

    // Get token for SSE connection
    const token = localStorage.getItem('access_token');

    // Use Fetch API with ReadableStream for SSE (supports Authorization headers)
    return new Promise(async (resolve, reject) => {
        try {
            const fetchResponse = await fetch(
                `http://localhost:18084/api/v1/spider/jobs/${jobId}/stream`,
                {
                    headers: {
                        'Accept': 'text/event-stream',
                        'Cache-Control': 'no-cache',
                        'Authorization': token ? `Bearer ${token}` : '',
                    },
                }
            );

            if (!fetchResponse.ok) {
                throw new SpiderServiceError(`SSE connection failed: ${fetchResponse.status}`);
            }

            if (!fetchResponse.body) {
                throw new SpiderServiceError('Response body is null');
            }

            const reader = fetchResponse.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';
            let lastStatus: JobStatusResponse | null = null;

            const processLine = (line: string) => {
                if (line.startsWith('data:')) {
                    const data = line.slice(5).trim();
                    try {
                        const parsed = JSON.parse(data);

                        if (parsed.job_id) {
                            const status: JobStatusResponse = parsed;
                            lastStatus = status;

                            console.log('📡 SSE update:', status.status);

                            if (onProgress) {
                                onProgress(status.status);
                            }

                            if (status.status === 'COMPLETED') {
                                console.log('✅ Scraping completed, got', status.results?.length, 'results');
                                reader.cancel();
                                resolve({
                                    google_id: status.google_id,
                                    restaurants: status.results || [],
                                    total_found: status.results?.length || 0,
                                });
                            } else if (status.status === 'FAILED') {
                                console.error('❌ Scraping failed:', status.error);
                                reader.cancel();
                                reject(new ScrapingError(status.error || 'Scraping failed'));
                            }
                        }
                    } catch (e) {
                        console.error('Failed to parse SSE data:', e);
                    }
                }
            };

            // Read stream
            const readStream = async () => {
                try {
                    while (true) {
                        const { done, value } = await reader.read();

                        if (done) {
                            console.log('📡 SSE stream ended');
                            if (lastStatus?.status === 'COMPLETED') {
                                resolve({
                                    google_id: lastStatus.google_id,
                                    restaurants: lastStatus.results || [],
                                    total_found: lastStatus.results?.length || 0,
                                });
                            } else {
                                reject(new ScrapingError('Stream ended without completion'));
                            }
                            break;
                        }

                        buffer += decoder.decode(value, { stream: true });
                        const lines = buffer.split('\n');
                        buffer = lines.pop() || '';

                        for (const line of lines) {
                            if (line.trim()) {
                                processLine(line);
                            }
                        }
                    }
                } catch (error) {
                    console.error('❌ SSE stream error:', error);
                    reject(new ScrapingError('Stream reading failed'));
                }
            };

            readStream();

        } catch (error: any) {
            console.error('❌ Failed to establish SSE connection:', error);
            reject(new SpiderServiceError(error.message || 'Failed to connect to SSE stream'));
        }
    });
}

// ============================================
// Exported Service Object
// ============================================

export const spiderService = {
    searchTabelog,
};

export default spiderService;
