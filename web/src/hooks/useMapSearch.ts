import { useMutation, useQuery, type UseMutationResult, type UseQueryResult } from '@tanstack/react-query';
import { mapService, RateLimitError, ValidationError } from '@/lib/api/map-service';
import type {
    QuickSearchRequest,
    QuickSearchResponse,
    AdvanceSearchRequest,
    AdvanceSearchResponse,
} from '@/types/search';

// ============================================
// Quick Search Hook
// ============================================

export interface UseQuickSearchOptions {
    enabled?: boolean;
    onSuccess?: (data: QuickSearchResponse) => void;
    onError?: (error: Error) => void;
}

/**
 * Hook for Quick Search (place details by ID)
 * Used when user clicks on a map marker
 */
export function useQuickSearch(
    request: QuickSearchRequest | null,
    options?: UseQuickSearchOptions
): UseQueryResult<QuickSearchResponse, Error> {
    return useQuery({
        queryKey: ['quick-search', request?.place_id, request?.language_code],
        queryFn: () => {
            if (!request) throw new Error('No request provided');
            return mapService.quickSearch(request);
        },
        enabled: options?.enabled !== false && !!request,
        staleTime: 5 * 60 * 1000, // 5 minutes
        retry: (failureCount, error) => {
            // Don't retry on rate limit or validation errors
            if (error instanceof RateLimitError || error instanceof ValidationError) {
                return false;
            }
            return failureCount < 2;
        },
    });
}

// ============================================
// Advance Search Hook
// ============================================

export interface UseAdvanceSearchResult {
    mutate: (request: AdvanceSearchRequest) => void;
    mutateAsync: (request: AdvanceSearchRequest) => Promise<AdvanceSearchResponse>;
    data?: AdvanceSearchResponse;
    error: Error | null;
    isError: boolean;
    isIdle: boolean;
    isPending: boolean;
    isSuccess: boolean;
    reset: () => void;
    isRateLimited: boolean;
    retryAfter?: number;
}

/**
 * Hook for Advance Search (text search with filters)
 * Used for the advance search form
 */
export function useAdvanceSearch(): UseAdvanceSearchResult {
    const mutation = useMutation({
        mutationFn: (request: AdvanceSearchRequest) => mapService.advanceSearch(request),
        retry: (failureCount, error) => {
            // Don't retry on rate limit or validation errors
            if (error instanceof RateLimitError || error instanceof ValidationError) {
                return false;
            }
            return failureCount < 2;
        },
    });

    const isRateLimited = mutation.error instanceof RateLimitError;
    const retryAfter = mutation.error instanceof RateLimitError
        ? mutation.error.retryAfter
        : undefined;

    return {
        mutate: mutation.mutate,
        mutateAsync: mutation.mutateAsync,
        data: mutation.data,
        error: mutation.error,
        isError: mutation.isError,
        isIdle: mutation.isIdle,
        isPending: mutation.isPending,
        isSuccess: mutation.isSuccess,
        reset: mutation.reset,
        isRateLimited,
        retryAfter,
    };
}

// ============================================
// Helper Hook for Map Bounds
// ============================================

export interface MapBounds {
    north: number;
    south: number;
    east: number;
    west: number;
}

/**
 * Convert map bounds to location bias for Advance Search
 */
export function mapBoundsToLocationBias(bounds: MapBounds) {
    return {
        rectangle: {
            low: {
                latitude: bounds.south,
                longitude: bounds.west,
            },
            high: {
                latitude: bounds.north,
                longitude: bounds.east,
            },
        },
    };
}
