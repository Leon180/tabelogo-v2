'use client';

import { Place } from '@/types/search';
import { PlaceCard } from './PlaceCard';

interface PlaceListProps {
    places: Place[];
    selectedPlaceId?: string;
    onPlaceClick?: (place: Place) => void;
    isLoading?: boolean;
}

export function PlaceList({ places, selectedPlaceId, onPlaceClick, isLoading = false }: PlaceListProps) {
    if (isLoading) {
        return (
            <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                    <div key={i} className="bg-zinc-800 rounded-lg overflow-hidden animate-pulse">
                        <div className="w-full h-48 bg-zinc-700" />
                        <div className="p-4 space-y-3">
                            <div className="h-6 bg-zinc-700 rounded w-3/4" />
                            <div className="h-4 bg-zinc-700 rounded w-1/2" />
                            <div className="h-4 bg-zinc-700 rounded w-full" />
                        </div>
                    </div>
                ))}
            </div>
        );
    }

    if (places.length === 0) {
        return (
            <div className="text-center py-12">
                <div className="text-zinc-500 mb-2">
                    <svg className="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </div>
                <h3 className="text-white font-medium mb-1">No restaurants found</h3>
                <p className="text-zinc-400 text-sm">Try adjusting your search criteria</p>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {places.map((place) => (
                <PlaceCard
                    key={place.id}
                    place={place}
                    isSelected={place.id === selectedPlaceId}
                    onClick={onPlaceClick}
                />
            ))}
        </div>
    );
}
