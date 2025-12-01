'use client';

import { useState } from 'react';
import { GoogleMap } from '@/components/Map/GoogleMap';
import { AdvanceSearchForm } from '@/components/Search/AdvanceSearchForm';
import { CollapsibleSearchPanel } from '@/components/Search/CollapsibleSearchPanel';
import { PlaceList } from '@/components/Place/PlaceList';
import { ResizableSidebar } from '@/components/ui/resizable-sidebar';
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
  const [selectedPlaceId, setSelectedPlaceId] = useState<string | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const { user, logout } = useAuth();

  const handleSearchResults = (results: Place[]) => {
    setPlaces(results);
    setErrorMessage(null);
    setIsSearching(false);
    setSelectedPlaceId(null); // Clear selection on new search
  };

  const handleSearchError = (error: Error) => {
    setErrorMessage(error.message);
    setIsSearching(false);
    // Clear error after 5 seconds
    setTimeout(() => setErrorMessage(null), 5000);
  };

  const handlePlaceClick = (place: Place) => {
    setSelectedPlaceId(place.id);
    // Close search panel on mobile when place is selected
    if (window.innerWidth < 1024) {
      setIsSearchPanelOpen(false);
    }
  };

  const handleMarkerClick = (place: Place) => {
    setSelectedPlaceId(place.id);
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

      {/* Main Content: Resizable Search Panel + Map */}
      <div className="flex flex-1 overflow-hidden">
        {/* Resizable Sidebar */}
        <ResizableSidebar
          defaultWidth={384}
          minWidth={320}
          maxWidth={600}
          className={`
            ${isSearchPanelOpen ? 'translate-x-0' : '-translate-x-full'}
            lg:translate-x-0 transition-transform duration-300
            bg-zinc-900
            absolute lg:relative h-full z-10
            flex flex-col
          `}
        >
          {/* Floating Search Panel */}
          <div className="p-4">
            <CollapsibleSearchPanel title="Search Filters" defaultExpanded={true}>
              <div className="pt-4 space-y-4">
                <AdvanceSearchForm
                  mapBounds={mapBounds}
                  onResults={handleSearchResults}
                  onError={handleSearchError}
                />
              </div>
            </CollapsibleSearchPanel>
          </div>

          {/* Restaurant List */}
          <div className="flex-1 min-h-0 overflow-y-auto px-4 pb-4">
            {places.length > 0 ? (
              <div>
                <h3 className="text-sm font-medium text-zinc-400 mb-4 px-2">
                  Showing {places.length} restaurant{places.length !== 1 ? 's' : ''}
                </h3>
                <PlaceList
                  places={places}
                  selectedPlaceId={selectedPlaceId || undefined}
                  onPlaceClick={handlePlaceClick}
                  isLoading={isSearching}
                />
              </div>
            ) : (
              <div className="px-2 pt-4">
                <div className="bg-zinc-800/50 rounded-lg border border-zinc-700 p-6">
                  <h3 className="text-sm font-medium text-zinc-400 mb-2">Quick Tips</h3>
                  <ul className="text-sm text-zinc-500 space-y-1">
                    <li>• Pan the map to set search area</li>
                    <li>• Enter a search query (e.g., "sushi Tokyo")</li>
                    <li>• Use filters to refine results</li>
                    <li>• Click cards to center map</li>
                  </ul>
                </div>
              </div>
            )}
          </div>
        </ResizableSidebar>

        {/* Map Container */}
        <main className="flex-1 relative">
          <GoogleMap
            places={places}
            selectedPlaceId={selectedPlaceId || undefined}
            onMarkerClick={handleMarkerClick}
            onBoundsChanged={handleBoundsChanged}
          />
        </main>
      </div>
    </div>
  );
}
