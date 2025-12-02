'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { useAdvanceSearch, mapBoundsToLocationBias, type MapBounds } from '@/hooks/useMapSearch';
import type { Place } from '@/types/search';

interface AdvanceSearchFormProps {
    mapBounds: MapBounds | null;
    onResults: (places: Place[]) => void;
    onError: (error: Error) => void;
}

export function AdvanceSearchForm({ mapBounds, onResults, onError }: AdvanceSearchFormProps) {
    const [query, setQuery] = useState('');
    const [minRating, setMinRating] = useState(0);
    const [openNow, setOpenNow] = useState(false);
    const [rankBy, setRankBy] = useState<'DISTANCE' | 'RELEVANCE'>('RELEVANCE');

    const { mutate, isPending, isError, error, data, isRateLimited, retryAfter } = useAdvanceSearch();

    // Handle successful response
    useEffect(() => {
        if (data) {
            console.log('✅ API Response:', data);
            onResults(data.places);
        }
    }, [data, onResults]);

    // Handle errors
    useEffect(() => {
        if (isError && error) {
            console.error('❌ API Error:', error);
            onError(error);
        }
    }, [isError, error, onError]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!query.trim()) {
            onError(new Error('Please enter a search query'));
            return;
        }

        if (!mapBounds) {
            onError(new Error('Map bounds not available'));
            return;
        }

        console.log('🔍 Submitting search:', { query, minRating, openNow, rankBy });
        console.log('📍 Map bounds:', mapBounds);

        // Build the request - only include defined fields
        const request: any = {
            text_query: query,
            location_bias: mapBoundsToLocationBias(mapBounds),
            max_result_count: 20,
            rank_preference: rankBy,
            language_code: 'en' as const,
        };

        // Only add optional fields if they have meaningful values
        if (minRating > 0) {
            request.min_rating = minRating;
        }
        if (openNow) {
            request.open_now = true;
        }

        console.log('📤 API Request:', request);

        // Call the API
        mutate(request);
    };

    const handleClear = () => {
        setQuery('');
        setMinRating(0);
        setOpenNow(false);
        setRankBy('RELEVANCE');
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {/* Search Input */}
            <div className="space-y-2">
                <Label htmlFor="search-query">Search</Label>
                <Input
                    id="search-query"
                    type="text"
                    placeholder="Search for restaurants..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="w-full"
                    disabled={isPending}
                />
            </div>

            {/* Operating Hours Switch */}
            <div className="flex items-center justify-between">
                <Label htmlFor="open-now" className="cursor-pointer">
                    Open Now
                </Label>
                <Switch
                    id="open-now"
                    checked={openNow}
                    onCheckedChange={setOpenNow}
                    disabled={isPending}
                />
            </div>

            {/* Rank Preference Switch */}
            <div className="flex items-center justify-between">
                <Label htmlFor="rank-by" className="cursor-pointer">
                    Relevance First
                </Label>
                <Switch
                    id="rank-by"
                    checked={rankBy === 'RELEVANCE'}
                    onCheckedChange={(checked) => setRankBy(checked ? 'RELEVANCE' : 'DISTANCE')}
                    disabled={isPending}
                />
            </div>

            {/* Minimum Rating Select */}
            <div className="space-y-2">
                <Label htmlFor="min-rating">Minimum Rating</Label>
                <Select
                    value={minRating.toString()}
                    onValueChange={(value) => setMinRating(Number(value))}
                    disabled={isPending}
                >
                    <SelectTrigger id="min-rating">
                        <SelectValue placeholder="Select rating" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="0">Any Rating</SelectItem>
                        <SelectItem value="3">3+ Stars</SelectItem>
                        <SelectItem value="4">4+ Stars</SelectItem>
                        <SelectItem value="4.5">4.5+ Stars</SelectItem>
                        <SelectItem value="5">5 Stars</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            {/* Rate Limit Warning */}
            {isRateLimited && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-md">
                    <p className="text-sm text-red-400">
                        ⚠️ Rate limit exceeded. Please wait {retryAfter ? `${Math.ceil((retryAfter - Date.now() / 1000) / 60)} minutes` : 'a moment'} before searching again.
                    </p>
                </div>
            )}

            {/* Error Message */}
            {isError && !isRateLimited && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-md">
                    <p className="text-sm text-red-400">
                        ❌ {error?.message || 'Search failed. Please try again.'}
                    </p>
                </div>
            )}

            {/* Results Count */}
            {data && (
                <div className="p-3 bg-green-500/10 border border-green-500/20 rounded-md">
                    <p className="text-sm text-green-400">
                        ✅ Found {data.total_count} result{data.total_count !== 1 ? 's' : ''} in {data.search_metadata.search_time_ms}ms
                    </p>
                </div>
            )}

            {/* Action Buttons */}
            <div className="flex gap-2">
                <Button
                    type="submit"
                    className="flex-1 bg-amber-500 hover:bg-amber-600"
                    disabled={isPending || isRateLimited}
                >
                    {isPending ? 'Searching...' : 'Search'}
                </Button>
                <Button
                    type="button"
                    variant="outline"
                    onClick={handleClear}
                    disabled={isPending}
                >
                    Clear
                </Button>
            </div>
        </form>
    );
}
