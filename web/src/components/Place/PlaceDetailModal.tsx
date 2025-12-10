'use client';

import { useEffect, useState } from 'react';
import { X, MapPin, Star, DollarSign, Clock, Phone, Globe, ExternalLink } from 'lucide-react';
import Image from 'next/image';
import { useRestaurantQuickSearch } from '@/hooks/useRestaurantSearch';
import { useQuickSearch } from '@/hooks/useMapSearch';
import { getPlacePhotoUrl, getPlaceholderImageUrl } from '@/lib/utils/getPlacePhotoUrl';
import { updateRestaurant } from '@/lib/api/restaurant-service';
import { getJapaneseName, quickSearch } from '@/lib/api/map-service';
import { searchTabelog, type TabelogRestaurant } from '@/lib/api/spider-service';
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
  // IMPORTANT: Use English (en) for area extraction and basic info
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
  const [tabelogResults, setTabelogResults] = useState<TabelogRestaurant[]>([]);
  const [scrapingProgress, setScrapingProgress] = useState<string>('');

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
      area: restaurant.area,
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
    console.log('🍜 Tabelog button clicked');

    if (!restaurantData?.restaurant) {
      setUpdateError('Restaurant data not available');
      return;
    }

    try {
      setIsUpdatingJaName(true);
      setUpdateError(null);

      // 1. Get Japanese name and addressComponents from Map Service
      console.log('📞 Fetching Japanese name and addressComponents from Map Service');
      const jaResponse = await quickSearch({
        place_id: placeId,
        language_code: 'ja',
        // Request addressComponents field
        api_mask: 'id,displayName'
      });

      const nameJa = jaResponse.result?.displayName?.text || place?.displayName?.text || '';
      console.log('✅ Japanese name:', nameJa);

      // 2. Extract place data with addressComponents from the Japanese API response
      const placeWithComponents = jaResponse.result || place;

      console.log('🚀 Preparing to call Spider Service:', {
        google_id: placeId,
        place_name: place?.displayName?.text || '',
        place_name_ja: nameJa,
      });

      // 4. Update restaurant with Japanese name AND area
      await updateRestaurant(restaurantData.restaurant.id, {
        name_ja: nameJa
      });

      console.log('✅ Updated restaurant with name_ja and area');
      console.log('Area:', place?.area);

      // 5. Call Spider Service to search Tabelog with SSE progress
      setScrapingProgress('Starting Tabelog search...');

      const tabelogResponse = await searchTabelog({
        google_id: placeId,
        place_name: place?.displayName?.text || '',
        place_name_ja: nameJa,
        area: place?.area || '', // Use extracted area from addressComponents
        max_results: 10
      }, (status) => {
        // Update progress based on job status
        const statusMessages: Record<string, string> = {
          'pending': 'Waiting to start...',
          'running': 'Scraping Tabelog...',
          'completed': 'Complete!',
          'failed': 'Failed'
        };
        setScrapingProgress(statusMessages[status] || status);
      });

      console.log('✅ Spider Service response:', {
        total_found: tabelogResponse.total_found,
        restaurants_count: tabelogResponse.restaurants?.length || 0,
        from_cache: tabelogResponse.from_cache
      });

      setTabelogResults(tabelogResponse.restaurants);
      setScrapingProgress(''); // Clear progress
      console.log(`✅ Found ${tabelogResponse.total_found} Tabelog restaurants`);
    } catch (error) {
      console.error('Failed to search Tabelog:', error);
      setUpdateError(error instanceof Error ? error.message : 'Failed to search Tabelog');
      setScrapingProgress(''); // Clear progress on error
    } finally {
      setIsUpdatingJaName(false);
    }
  };

  // Helper function to extract area from Google Maps addressComponents
  // Following v1 approach: extract locality (ward/city) from addressComponents
  function extractArea(address: string): string {
    // This is a fallback - ideally we should use addressComponents from place object
    // For now, extract first part of address (e.g., "Sugamo, Tokyo" -> "Sugamo")
    const parts = address.split(',');
    return parts[0]?.trim() || '';
  }

  // Extract locality from place's addressComponents (preferred method)
  // Following v1 approach: extract administrative_area_level_1 (e.g., "Tokyo")
  function extractLocalityFromPlace(place: any): string {
    console.log('🔍 Extracting locality from place:', {
      hasAddressComponents: !!place?.addressComponents,
      addressComponentsLength: place?.addressComponents?.length,
      formattedAddress: place?.formattedAddress
    });

    if (!place?.addressComponents) {
      const fallback = extractArea(place?.formattedAddress || '');
      console.log('⚠️ No addressComponents, using fallback:', fallback);
      return fallback;
    }

    // Find administrative_area_level_1 (prefecture/state level - e.g., "Tokyo")
    // This matches v1 behavior
    for (const component of place.addressComponents) {
      const types = component.types || [];

      console.log('📍 Checking component:', {
        types,
        longText: component.longText,
        shortText: component.shortText
      });

      // Priority: administrative_area_level_1 (Tokyo, Osaka, etc.)
      if (types.includes('administrative_area_level_1')) {
        const result = component.longText || component.shortText || '';
        console.log('✅ Found administrative_area_level_1:', result);
        return result;
      }
    }

    // Fallback to locality if administrative_area_level_1 not found
    for (const component of place.addressComponents) {
      const types = component.types || [];

      if (types.includes('locality')) {
        const result = component.longText || component.shortText || '';
        console.log('✅ Found locality (fallback):', result);
        return result;
      }
    }

    // Final fallback to formatted address
    const fallback = extractArea(place?.formattedAddress || '');
    console.log('⚠️ No area component found, using fallback:', fallback);
    return fallback;
  }

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

                    {/* Tabelog Results */}
                    {tabelogResults && tabelogResults.length > 0 && (
                      <div className="mt-4">
                        <h3 className="text-lg font-semibold text-white mb-3">
                          🍜 Tabelog Results ({tabelogResults.length})
                        </h3>
                        <div className="space-y-3 max-h-96 overflow-y-auto">
                          {tabelogResults.map((restaurant, index) => (
                            <a
                              key={index}
                              href={restaurant.link}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="block p-3 bg-zinc-800/50 hover:bg-zinc-700/50 rounded-lg transition-colors border border-zinc-700/50 hover:border-orange-600/50"
                            >
                              <div className="flex items-start justify-between gap-3">
                                <div className="flex-1 min-w-0">
                                  <h4 className="font-medium text-white truncate">
                                    {restaurant.name}
                                  </h4>
                                  {restaurant.types && restaurant.types.length > 0 && (
                                    <div className="flex flex-wrap gap-1 mt-1">
                                      {restaurant.types.slice(0, 3).map((type, i) => (
                                        <span
                                          key={i}
                                          className="text-xs px-2 py-0.5 bg-zinc-700/50 text-zinc-300 rounded"
                                        >
                                          {type}
                                        </span>
                                      ))}
                                    </div>
                                  )}
                                </div>
                                <div className="flex flex-col items-end gap-1">
                                  {restaurant.rating > 0 && (
                                    <div className="flex items-center gap-1">
                                      <span className="text-yellow-400">⭐</span>
                                      <span className="text-sm font-medium text-white">
                                        {restaurant.rating.toFixed(2)}
                                      </span>
                                    </div>
                                  )}
                                  {restaurant.rating_count > 0 && (
                                    <span className="text-xs text-zinc-400">
                                      {restaurant.rating_count} reviews
                                    </span>
                                  )}
                                </div>
                              </div>
                              {restaurant.photos && restaurant.photos.length > 0 && (
                                <div className="mt-2 flex gap-2 overflow-x-auto">
                                  {restaurant.photos.slice(0, 3).map((photo, i) => (
                                    <img
                                      key={i}
                                      src={photo}
                                      alt={`${restaurant.name} photo ${i + 1}`}
                                      className="h-16 w-16 object-cover rounded"
                                    />
                                  ))}
                                </div>
                              )}
                            </a>
                          ))}
                        </div>
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
