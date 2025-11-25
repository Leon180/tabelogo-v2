'use client';

import { useState } from 'react';
import { GoogleMap } from '@/components/Map/GoogleMap';
import { AdvanceSearchForm } from '@/components/Search/AdvanceSearchForm';
import { Button } from '@/components/ui/button';
import { Menu, X } from 'lucide-react';
import type { Restaurant } from '@/types/restaurant';
import type { SearchFilters } from '@/types/search';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

// Mock data for initial development
const mockRestaurants: Restaurant[] = [
  {
    id: '1',
    name: 'Sukiyabashi Jiro',
    source: 'google',
    external_id: 'ChIJ...',
    address: 'Tokyo, Ginza',
    latitude: 35.6708,
    longitude: 139.7634,
    rating: 4.8,
    price_range: '$$$$',
    cuisine_type: 'Sushi',
  },
  {
    id: '2',
    name: 'Narisawa',
    source: 'google',
    external_id: 'ChIJ...',
    address: 'Tokyo, Minato',
    latitude: 35.6654,
    longitude: 139.7236,
    rating: 4.7,
    price_range: '$$$$',
    cuisine_type: 'French-Japanese',
  },
];

export default function HomePage() {
  const [restaurants, setRestaurants] = useState<Restaurant[]>(mockRestaurants);
  const [isSearchPanelOpen, setIsSearchPanelOpen] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const { user, logout } = useAuth();

  const handleSearch = async (filters: SearchFilters) => {
    setIsLoading(true);
    try {
      // TODO: Call advance search API
      console.log('Search filters:', filters);
      // const results = await mapService.advanceSearch({...});
      // setRestaurants(results);
    } catch (error) {
      console.error('Search error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleMarkerClick = (restaurant: Restaurant) => {
    console.log('Marker clicked:', restaurant);
    // TODO: Show restaurant details or navigate to detail page
  };

  const handleBoundsChanged = (bounds: google.maps.LatLngBounds) => {
    // TODO: Update search area based on map bounds
    console.log('Map bounds changed:', bounds.toJSON());
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
                Use the form below to search for restaurants, or click on markers on the map.
              </p>
            </div>

            <AdvanceSearchForm onSearch={handleSearch} isLoading={isLoading} />

            <div className="pt-6 border-t border-zinc-800">
              <h3 className="text-sm font-medium text-zinc-400 mb-2">Quick Tips</h3>
              <ul className="text-sm text-zinc-500 space-y-1">
                <li>• Click markers on the map for quick info</li>
                <li>• Pan the map to explore different areas</li>
                <li>• Use filters to refine your search</li>
              </ul>
            </div>
          </div>
        </aside>

        {/* Map Container */}
        <main className="flex-1 relative">
          <GoogleMap
            restaurants={restaurants}
            onMarkerClick={handleMarkerClick}
            onBoundsChanged={handleBoundsChanged}
          />
        </main>
      </div>
    </div>
  );
}
