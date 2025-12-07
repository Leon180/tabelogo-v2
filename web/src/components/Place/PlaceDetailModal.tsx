'use client';

import { useEffect, useState } from 'react';
import { X, MapPin, Star, DollarSign, Clock, Phone, Globe, ExternalLink } from 'lucide-react';
import Image from 'next/image';
import { useRestaurantQuickSearch } from '@/hooks/useRestaurantSearch';
import { useQuickSearch } from '@/hooks/useMapSearch';
import { getPlacePhotoUrl, getPlaceholderImageUrl } from '@/lib/utils/getPlacePhotoUrl';
import { updateRestaurant } from '@/lib/api/restaurant-service';
import { getJapaneseName } from '@/lib/api/map-service';
import type { Place } from '@/types/search';

interface PlaceDetailModalProps {
  placeId: string;
  isOpen: boolean;
  onClose: () => void;
}

export function PlaceDetailModal({ placeId, isOpen, onClose }: PlaceDetailModalProps) {
  // NEW: Use Restaurant Service (cache-first)
  const { data: restaurantData, isLoading: isRestaurantLoading, error: restaurantError } = useRestaurantQuickSearch(
    isOpen ? placeId : null
  );

  // FALLBACK: Use Map Service if Restaurant Service fails
  const { data: mapData, isLoading: isMapLoading, error: mapError } = useQuickSearch(
    isOpen && restaurantError ? { place_id: placeId, language_code: 'en' } : null
  );

  // Use Restaurant Service data if available, otherwise fall back to Map Service
  const data = restaurantData || mapData;
  const isLoading = isRestaurantLoading || isMapLoading;
  const error = restaurantError && mapError ? mapError : null;

  // Tabelog integration state
  const [isUpdatingJaName, setIsUpdatingJaName] = useState(false);
  const [japaneseName, setJapaneseName] = useState<string | null>(null);
  const [updateError, setUpdateError] = useState<string | null>(null);

  // Close modal on ESC key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = 'unset';
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  // Extract place data - Restaurant Service returns different format
  const place = restaurantData ? convertRestaurantToPlace(restaurantData.restaurant) : mapData?.result;

  // Helper function to convert Restaurant Service format to Place format
  function convertRestaurantToPlace(restaurant: any): Place {
    return {
      id: restaurant.external_id,
      displayName: { text: restaurant.name },
      formattedAddress: restaurant.address,
      location: {
        latitude: restaurant.latitude,
        longitude: restaurant.longitude,
      },
      rating: restaurant.rating,
      priceLevel: restaurant.price_range ? `PRICE_LEVEL_${restaurant.price_range}` : undefined,
      nationalPhoneNumber: restaurant.phone,
      websiteUri: restaurant.website,
      // Note: Restaurant Service doesn't include photos, opening hours yet
      // These will be added in future updates
    } as Place;
  }

  // Handle Tabelog button click
  const handleTabelogClick = async () => {
    if (!restaurantData?.restaurant) {
      setUpdateError('Restaurant data not available');
      return;
    }

    setIsUpdatingJaName(true);
    setUpdateError(null);

    try {
      // 1. Get Japanese name from Google Maps
      const nameJa = await getJapaneseName(placeId);
      if (!nameJa) {
        setUpdateError('Could not retrieve Japanese name from Google Maps');
        setIsUpdatingJaName(false);
        return;
      }

      setJapaneseName(nameJa);

      // 2. Update restaurant with Japanese name
      await updateRestaurant(restaurantData.restaurant.id, {
        name_ja: nameJa
      });

      // 3. TODO: Call Spider Service with name_ja
      // const tabelogResults = await searchTabelog({
      //     google_id: placeId,
      //     name_ja: nameJa,
      //     area: extractArea(place.formattedAddress)
      // });

      console.log(`✅ Updated restaurant with Japanese name: ${nameJa}`);
    } catch (error) {
      console.error('Failed to update restaurant:', error);
      setUpdateError(error instanceof Error ? error.message : 'Failed to update restaurant');
    } finally {
      setIsUpdatingJaName(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
      {/* Modal Container */}
      <div className="relative w-full max-w-4xl max-h-[90vh] bg-zinc-900 rounded-xl shadow-2xl overflow-hidden">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 z-10 p-2 bg-zinc-800/90 hover:bg-zinc-700 rounded-full transition-colors"
        >
          <X className="w-5 h-5 text-white" />
        </button>

        {/* Content */}
        <div className="overflow-y-auto max-h-[90vh]">
          {isLoading && (
            <div className="p-12 flex items-center justify-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
            </div>
          )}

          {error && (
            <div className="p-12 text-center">
              <p className="text-red-400 mb-2">Failed to load place details</p>
              <p className="text-zinc-500 text-sm">{error.message}</p>
            </div>
          )}

          {place && (
            <>
              {/* Photo Gallery */}
              <div className="relative w-full h-96 bg-zinc-800">
                {place.photos && place.photos.length > 0 ? (
                  <Image
                    src={getPlacePhotoUrl(place.photos[0].name, 800, 600)}
                    alt={place.displayName?.text || 'Restaurant'}
                    fill
                    className="object-cover"
                    sizes="(max-width: 1024px) 100vw, 800px"
                    unoptimized
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center bg-zinc-800">
                    <p className="text-zinc-500">No Image Available</p>
                  </div>
                )}

                {/* Cache Badge - Show if using Restaurant Service cache */}
                {restaurantData && (
                  <div className="absolute top-4 left-4 px-3 py-1 bg-green-500/90 text-white text-xs font-medium rounded-full flex items-center gap-1">
                    ⚡ Restaurant Service (Cached)
                  </div>
                )}
                {/* Fallback Badge - Show if using Map Service */}
                {!restaurantData && mapData?.source === 'redis' && (
                  <div className="absolute top-4 left-4 px-3 py-1 bg-blue-500/90 text-white text-xs font-medium rounded-full">
                    Map Service (Cached)
                  </div>
                )}
              </div>

              {/* Details Section */}
              <div className="p-6 space-y-6">
                {/* Title & Rating */}
                <div>
                  <h2 className="text-3xl font-bold text-white mb-2">
                    {place.displayName?.text || 'Unknown Place'}
                  </h2>

                  {place.rating && (
                    <div className="flex items-center gap-3">
                      <div className="flex items-center gap-1">
                        <Star className="w-5 h-5 text-amber-500 fill-amber-500" />
                        <span className="text-xl font-semibold text-white">{place.rating.toFixed(1)}</span>
                      </div>
                      {place.userRatingCount && (
                        <span className="text-zinc-400">({place.userRatingCount} reviews)</span>
                      )}
                    </div>
                  )}
                </div>

                {/* Quick Info Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {/* Price Level */}
                  {place.priceLevel && place.priceLevel !== 'PRICE_LEVEL_UNSPECIFIED' && (
                    <div className="flex items-start gap-3 p-4 bg-zinc-800/50 rounded-lg">
                      <DollarSign className="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="text-sm text-zinc-400">Price Level</p>
                        <p className="text-white font-medium">
                          {place.priceLevel.replace('PRICE_LEVEL_', '').replace('_', ' ')}
                        </p>
                      </div>
                    </div>
                  )}

                  {/* Opening Hours */}
                  {place.currentOpeningHours && (
                    <div className="flex items-start gap-3 p-4 bg-zinc-800/50 rounded-lg">
                      <Clock className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="text-sm text-zinc-400">Status</p>
                        <p className={`font-medium ${place.currentOpeningHours.openNow ? 'text-green-500' : 'text-red-500'}`}>
                          {place.currentOpeningHours.openNow ? 'Open Now' : 'Closed'}
                        </p>
                      </div>
                    </div>
                  )}

                  {/* Address */}
                  {place.formattedAddress && (
                    <div className="flex items-start gap-3 p-4 bg-zinc-800/50 rounded-lg md:col-span-2">
                      <MapPin className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="text-sm text-zinc-400">Address</p>
                        <p className="text-white">{place.formattedAddress}</p>
                      </div>
                    </div>
                  )}

                  {/* Phone */}
                  {place.nationalPhoneNumber && (
                    <div className="flex items-start gap-3 p-4 bg-zinc-800/50 rounded-lg">
                      <Phone className="w-5 h-5 text-purple-500 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="text-sm text-zinc-400">Phone</p>
                        <a href={`tel:${place.nationalPhoneNumber}`} className="text-white hover:text-amber-500">
                          {place.nationalPhoneNumber}
                        </a>
                      </div>
                    </div>
                  )}

                  {/* Website */}
                  {place.websiteUri && (
                    <div className="flex items-start gap-3 p-4 bg-zinc-800/50 rounded-lg">
                      <Globe className="w-5 h-5 text-cyan-500 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="text-sm text-zinc-400">Website</p>
                        <a
                          href={place.websiteUri}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-white hover:text-amber-500 flex items-center gap-1"
                        >
                          Visit Website <ExternalLink className="w-3 h-3" />
                        </a>
                      </div>
                    </div>
                  )}
                </div>

                {/* Tabelog Integration - Only show if restaurant data available */}
                {restaurantData?.restaurant && (
                  <div className="pt-4 border-t border-zinc-700">
                    <button
                      onClick={handleTabelogClick}
                      disabled={isUpdatingJaName}
                      className="w-full px-4 py-3 bg-gradient-to-r from-orange-600 to-orange-700 hover:from-orange-700 hover:to-orange-800 disabled:from-gray-600 disabled:to-gray-700 text-white font-semibold rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl disabled:cursor-not-allowed"
                    >
                      {isUpdatingJaName ? (
                        <span className="flex items-center justify-center gap-2">
                          <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                          Updating...
                        </span>
                      ) : (
                        '🍜 Search Tabelog'
                      )}
                    </button>

                    {/* Japanese Name Display */}
                    {japaneseName && (
                      <div className="mt-3 p-3 bg-green-900/30 border border-green-700/50 rounded-lg">
                        <p className="text-sm text-green-400 font-medium">
                          ✓ Japanese name: {japaneseName}
                        </p>
                      </div>
                    )}

                    {/* Error Display */}
                    {updateError && (
                      <div className="mt-3 p-3 bg-red-900/30 border border-red-700/50 rounded-lg">
                        <p className="text-sm text-red-400">
                          ✗ {updateError}
                        </p>
                      </div>
                    )}
                  </div>
                )}

                {/* Opening Hours Details */}
                {place.currentOpeningHours?.weekdayDescriptions && (
                  <div className="p-4 bg-zinc-800/50 rounded-lg">
                    <h3 className="text-white font-semibold mb-3">Opening Hours</h3>
                    <div className="space-y-1">
                      {place.currentOpeningHours.weekdayDescriptions.map((hours, index) => (
                        <p key={index} className="text-sm text-zinc-300">{hours}</p>
                      ))}
                    </div>
                  </div>
                )}

                {/* Google Maps Link */}
                {place.googleMapsUri && (
                  <a
                    href={place.googleMapsUri}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="block w-full py-3 px-4 bg-amber-500 hover:bg-amber-600 text-white font-medium text-center rounded-lg transition-colors"
                  >
                    View on Google Maps
                  </a>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
