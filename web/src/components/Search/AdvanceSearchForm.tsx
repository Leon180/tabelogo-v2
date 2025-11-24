'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import type { SearchFilters } from '@/types/search';

interface AdvanceSearchFormProps {
    onSearch: (filters: SearchFilters) => void;
    isLoading?: boolean;
}

export function AdvanceSearchForm({ onSearch, isLoading = false }: AdvanceSearchFormProps) {
    const [query, setQuery] = useState('');
    const [minRating, setMinRating] = useState(0);
    const [openNow, setOpenNow] = useState(false);
    const [rankBy, setRankBy] = useState<'relevance' | 'distance'>('relevance');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSearch({
            query,
            minRating,
            openNow,
            rankBy,
        });
    };

    const handleClear = () => {
        setQuery('');
        setMinRating(0);
        setOpenNow(false);
        setRankBy('relevance');
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
                />
            </div>

            {/* Rank Preference Switch */}
            <div className="flex items-center justify-between">
                <Label htmlFor="rank-by" className="cursor-pointer">
                    Relevance First
                </Label>
                <Switch
                    id="rank-by"
                    checked={rankBy === 'relevance'}
                    onCheckedChange={(checked) => setRankBy(checked ? 'relevance' : 'distance')}
                />
            </div>

            {/* Minimum Rating Select */}
            <div className="space-y-2">
                <Label htmlFor="min-rating">Minimum Rating</Label>
                <Select value={minRating.toString()} onValueChange={(value) => setMinRating(Number(value))}>
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

            {/* Action Buttons */}
            <div className="flex gap-2">
                <Button
                    type="submit"
                    className="flex-1 bg-amber-500 hover:bg-amber-600"
                    disabled={isLoading}
                >
                    {isLoading ? 'Searching...' : 'Search'}
                </Button>
                <Button
                    type="button"
                    variant="outline"
                    onClick={handleClear}
                    disabled={isLoading}
                >
                    Clear
                </Button>
            </div>
        </form>
    );
}
