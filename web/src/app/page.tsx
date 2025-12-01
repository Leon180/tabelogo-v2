'use client';

import { useState } from 'react';
import { GoogleMap } from '@/components/Map/GoogleMap';
import { AdvanceSearchForm } from '@/components/Search/AdvanceSearchForm';
import { Button } from '@/components/ui/button';
import { Menu, X } from 'lucide-react';
import type { Place } from '@/types/search';
import type { MapBounds } from '@/hooks/useMapSearch';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

export default function HomePage() {
  const [places, setPlaces] = useState<Place[]>([]);
  const [mapBounds, setMapBounds] = useState<MapBounds | null>(null);
  const [isSearchPanelOpen, setIsSearchPanelOpen] = useState(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const { user, logout } = useAuth();

  const handleSearchResults = (results: Place[]) => {
    console.log('📍 Received search results:', results);
    setPlaces(results);
    setErrorMessage(null);
  };

  const handleSearchError = (error: Error) => {
    console.error('🚨 Search error:', error);
    setErrorMessage(error.message);
    // Clear error after 5 seconds
    setTimeout(() => setErrorMessage(null), 5000);
  };

  const handleMarkerClick = (place: any) => {
    console.log('Marker clicked:', place);
    // TODO: Show place details modal with Quick Search
  };

  const handleBoundsChanged = (bounds: google.maps.LatLngBounds) => {
    const boundsJson = bounds.toJSON();
    const mapBounds: MapBounds = {
      north: boundsJson.north,
      south: boundsJson.south,
      east: boundsJson.east,
      west: boundsJson.west,
    };
    setMapBounds(mapBounds);
    console.log('🗺️ Map bounds updated:', mapBounds);
  };

  return (
    <div className="flex flex-col h-screen bg-zinc-950">
      {/* Navigation Bar */}
      <nav className="bg-zinc-900 border-b border-zinc-800 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setIsSearchPanelOpen(!isSearchPanelOpen)}
            className="lg:hidden"
          >
            {isSearchPanelOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
          <h1 className="text-2xl font-bold text-amber-500">Tabelogo</h1>
        </div>
        <div className="flex items-center gap-4">
          {user ? (
            <div className="flex items-center gap-4">
              <span className="text-zinc-300 text-sm hidden sm:inline-block">Hi, {user.username}</span>
              <Button variant="ghost" className="text-zinc-300 hover:text-white" onClick={() => logout()}>
                Logout
              </Button>
            </div>
          ) : (
            <Link href="/auth/login">
              <Button variant="ghost" className="text-zinc-300 hover:text-white">
                Login
              </Button>
            </Link>
          )}
        </div>
      </nav>

      {/* Error Toast */}
      {errorMessage && (
        <div className="fixed top-20 right-4 z-50 p-4 bg-red-500/90 text-white rounded-lg shadow-lg max-w-md">
          <p className="font-medium">❌ Error</p>
          <p className="text-sm mt-1">{errorMessage}</p>
        </div>
      )}

      {/* Main Content: Search Panel + Map */}
      <div className="flex flex-1 overflow-hidden">
        {/* Search Panel (Sidebar) */}
        <aside
          className={`
            ${isSearchPanelOpen ? 'translate-x-0' : '-translate-x-full'}
            lg:translate-x-0 transition-transform duration-300
            w-full lg:w-96 bg-zinc-900 border-r border-zinc-800
            overflow-y-auto p-6
            absolute lg:relative h-full z-10
          `}
        >
          <div className="space-y-6">
            <div>
              <h2 className="text-xl font-semibold text-white mb-4">Search Restaurants</h2>
              <p className="text-sm text-zinc-400 mb-6">
                Search for restaurants using the Map Service API. Results will appear on the map.
              </p>
            </div>

            <AdvanceSearchForm
              mapBounds={mapBounds}
              onResults={handleSearchResults}
              onError={handleSearchError}
            />

            <div className="pt-6 border-t border-zinc-800">
              <h3 className="text-sm font-medium text-zinc-400 mb-2">Quick Tips</h3>
              <ul className="text-sm text-zinc-500 space-y-1">
                <li>• Pan the map to set search area</li>
                <li>• Enter a search query (e.g., "sushi Tokyo")</li>
                <li>• Use filters to refine results</li>
                <li>• Click markers for details</li>
              </ul>
            </div>

            {/* Results Summary */}
            {places.length > 0 && (
              <div className="pt-6 border-t border-zinc-800">
                <h3 className="text-sm font-medium text-zinc-400 mb-2">Results</h3>
                <p className="text-sm text-zinc-300">
                  Showing {places.length} restaurant{places.length !== 1 ? 's' : ''} on the map
                </p>
              </div>
            )}
          </div>
        </aside>

        {/* Map Container */}
        <main className="flex-1 relative">
          <GoogleMap
            places={places}
            onMarkerClick={handleMarkerClick}
            onBoundsChanged={handleBoundsChanged}
          />
        </main>
      </div>
    </div>
  );
}
