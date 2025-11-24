import axios from 'axios';
import type { QuickSearchParams, AdvanceSearchParams } from '@/types/search';
import type { Restaurant } from '@/types/restaurant';

// API base URL - will be configured via environment variables
const MAP_SERVICE_URL = process.env.NEXT_PUBLIC_MAP_SERVICE_URL || 'http://localhost:8080';

const mapClient = axios.create({
    baseURL: MAP_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

/**
 * Quick Search - Get single restaurant by place_id
 * Used when user clicks on a map marker
 */
export async function quickSearch(params: QuickSearchParams): Promise<Restaurant> {
    const response = await mapClient.post('/quick_search', params);
    return response.data.result;
}

/**
 * Advance Search - Text search with filters
 * Used for the advance search form
 */
export async function advanceSearch(params: AdvanceSearchParams): Promise<Restaurant[]> {
    const response = await mapClient.post('/advance_search', params);
    return response.data.places || [];
}

export const mapService = {
    quickSearch,
    advanceSearch,
};
