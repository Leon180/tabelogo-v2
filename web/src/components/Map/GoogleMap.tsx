'use client';

import { APIProvider, Map, Marker, InfoWindow } from '@vis.gl/react-google-maps';
import { useState, useCallback } from 'react';
import type { Restaurant } from '@/types/restaurant';

interface GoogleMapProps {
    restaurants: Restaurant[];
    center?: { lat: number; lng: number };
    zoom?: number;
    onMarkerClick?: (restaurant: Restaurant) => void;
    onBoundsChanged?: (bounds: google.maps.LatLngBounds) => void;
}

export function GoogleMap({
    restaurants,
    center = {
        lat: Number(process.env.NEXT_PUBLIC_DEFAULT_LAT) || 35.6762,
        lng: Number(process.env.NEXT_PUBLIC_DEFAULT_LNG) || 139.6503,
    },
    zoom = 13,
    onMarkerClick,
    onBoundsChanged,
}: GoogleMapProps) {
    const [selectedRestaurant, setSelectedRestaurant] = useState<Restaurant | null>(null);
    const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || '';

    const handleMarkerClick = useCallback((restaurant: Restaurant) => {
        setSelectedRestaurant(restaurant);
        onMarkerClick?.(restaurant);
    }, [onMarkerClick]);

    const handleMapBoundsChanged = useCallback((map: google.maps.Map) => {
        const bounds = map.getBounds();
        if (bounds) {
            onBoundsChanged?.(bounds);
        }
    }, [onBoundsChanged]);

    if (!apiKey) {
        return (
            <div className="flex items-center justify-center h-full bg-zinc-900 text-zinc-400">
                <p>Google Maps API key not configured. Please check ENV_CONFIG.md</p>
            </div>
        );
    }

    return (
        <APIProvider apiKey={apiKey}>
            <Map
                defaultCenter={center}
                defaultZoom={zoom}
                mapId="tabelogo-map"
                className="w-full h-full"
                onBoundsChanged={({ map }) => handleMapBoundsChanged(map)}
            >
                {/* Render restaurant markers */}
                {restaurants.map((restaurant) => (
                    <Marker
                        key={restaurant.id}
                        position={{ lat: restaurant.latitude, lng: restaurant.longitude }}
                        onClick={() => handleMarkerClick(restaurant)}
                    />
                ))}

                {/* Show info window for selected restaurant */}
                {selectedRestaurant && (
                    <InfoWindow
                        position={{
                            lat: selectedRestaurant.latitude,
                            lng: selectedRestaurant.longitude,
                        }}
                        onCloseClick={() => setSelectedRestaurant(null)}
                    >
                        <div className="p-2 min-w-[200px]">
                            <h3 className="font-semibold text-lg">{selectedRestaurant.name}</h3>
                            <p className="text-sm text-gray-600">{selectedRestaurant.cuisine_type}</p>
                            <p className="text-sm">⭐ {selectedRestaurant.rating}</p>
                            <p className="text-sm">{selectedRestaurant.price_range}</p>
                        </div>
                    </InfoWindow>
                )}
            </Map>
        </APIProvider>
    );
}
