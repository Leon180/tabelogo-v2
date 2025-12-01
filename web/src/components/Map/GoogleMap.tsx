'use client';

import { APIProvider, Map, Marker, InfoWindow } from '@vis.gl/react-google-maps';
import { useState, useCallback } from 'react';
import type { Place } from '@/types/search';

interface GoogleMapProps {
    places: Place[];
    center?: { lat: number; lng: number };
    zoom?: number;
    onMarkerClick?: (place: Place) => void;
    onBoundsChanged?: (bounds: google.maps.LatLngBounds) => void;
}

export function GoogleMap({
    places,
    center = {
        lat: Number(process.env.NEXT_PUBLIC_DEFAULT_LAT) || 35.6762,
        lng: Number(process.env.NEXT_PUBLIC_DEFAULT_LNG) || 139.6503,
    },
    zoom = 13,
    onMarkerClick,
    onBoundsChanged,
}: GoogleMapProps) {
    const [selectedPlace, setSelectedPlace] = useState<Place | null>(null);
    const apiKey = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || '';

    const handleMarkerClick = useCallback((place: Place) => {
        setSelectedPlace(place);
        onMarkerClick?.(place);
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
                {/* Render place markers */}
                {places
                    .filter((place) => place.location) // Only render places with valid location
                    .map((place) => (
                        <Marker
                            key={place.id}
                            position={{ lat: place.location!.latitude, lng: place.location!.longitude }}
                            onClick={() => handleMarkerClick(place)}
                        />
                    ))}

                {/* Show info window for selected place */}
                {selectedPlace && selectedPlace.location && (
                    <InfoWindow
                        position={{
                            lat: selectedPlace.location.latitude,
                            lng: selectedPlace.location.longitude,
                        }}
                        onCloseClick={() => setSelectedPlace(null)}
                    >
                        <div className="p-2 min-w-[200px]">
                            <h3 className="font-semibold text-lg">{selectedPlace.displayName?.text || 'Unknown'}</h3>
                            <p className="text-sm text-gray-600">{selectedPlace.formattedAddress}</p>
                            {selectedPlace.rating && <p className="text-sm">⭐ {selectedPlace.rating}</p>}
                            {selectedPlace.priceLevel && <p className="text-sm">{selectedPlace.priceLevel}</p>}
                        </div>
                    </InfoWindow>
                )}
            </Map>
        </APIProvider>
    );
}
