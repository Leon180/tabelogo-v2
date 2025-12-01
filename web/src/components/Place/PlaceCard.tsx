'use client';

import Image from 'next/image';
import { Place } from '@/types/search';
import { getPlacePhotoUrl, getPlaceholderImageUrl } from '@/lib/utils/getPlacePhotoUrl';

interface PlaceCardProps {
    place: Place;
    isSelected?: boolean;
    onClick?: (place: Place) => void;
}

export function PlaceCard({ place, isSelected = false, onClick }: PlaceCardProps) {
    const photoUrl = place.photos?.[0]?.name
        ? getPlacePhotoUrl(place.photos[0].name, 400, 300)
        : getPlaceholderImageUrl();

    const priceLevel = place.priceLevel?.replace('PRICE_LEVEL_', '').toLowerCase();
    const isOpenNow = place.currentOpeningHours?.openNow;

    const handleClick = () => {
        onClick?.(place);
    };

    return (
        <div
            onClick={handleClick}
            className={`
        bg-zinc-800 rounded-lg overflow-hidden cursor-pointer
        transition-all duration-200 hover:bg-zinc-750 hover:shadow-lg
        border-2 ${isSelected ? 'border-amber-500' : 'border-transparent'}
      `}
        >
            {/* Photo */}
            <div className="relative w-full h-48 bg-zinc-900">
                <Image
                    src={photoUrl}
                    alt={place.displayName?.text || 'Restaurant'}
                    fill
                    className="object-cover"
                    sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                    unoptimized // For external URLs
                />

                {/* Open/Closed Badge */}
                {isOpenNow !== undefined && (
                    <div className={`
            absolute top-2 right-2 px-2 py-1 rounded text-xs font-medium
            ${isOpenNow ? 'bg-green-500 text-white' : 'bg-red-500 text-white'}
          `}>
                        {isOpenNow ? 'Open Now' : 'Closed'}
                    </div>
                )}
            </div>

            {/* Content */}
            <div className="p-4">
                {/* Name */}
                <h3 className="text-white font-semibold text-lg mb-2 line-clamp-1">
                    {place.displayName?.text || 'Unknown Place'}
                </h3>

                {/* Rating & Price */}
                <div className="flex items-center gap-3 mb-2">
                    {place.rating && (
                        <div className="flex items-center gap-1">
                            <span className="text-amber-500">⭐</span>
                            <span className="text-white font-medium">{place.rating.toFixed(1)}</span>
                            {place.userRatingCount && (
                                <span className="text-zinc-400 text-sm">({place.userRatingCount})</span>
                            )}
                        </div>
                    )}

                    {priceLevel && priceLevel !== 'unspecified' && (
                        <div className="text-zinc-400 text-sm">
                            {priceLevel === 'free' && 'Free'}
                            {priceLevel === 'inexpensive' && '$'}
                            {priceLevel === 'moderate' && '$$'}
                            {priceLevel === 'expensive' && '$$$'}
                            {priceLevel === 'very_expensive' && '$$$$'}
                        </div>
                    )}
                </div>

                {/* Address */}
                {place.formattedAddress && (
                    <p className="text-zinc-400 text-sm line-clamp-2">
                        {place.formattedAddress}
                    </p>
                )}
            </div>
        </div>
    );
}
