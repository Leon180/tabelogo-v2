import axios from 'axios';
import type { Restaurant, UserFavorite } from '@/types/restaurant';

const RESTAURANT_SERVICE_URL = process.env.NEXT_PUBLIC_RESTAURANT_SERVICE_URL || 'http://localhost:8082';

const restaurantClient = axios.create({
    baseURL: RESTAURANT_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests
restaurantClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

/**
 * Get restaurant by ID
 */
export async function getRestaurant(id: string): Promise<Restaurant> {
    const response = await restaurantClient.get(`/restaurants/${id}`);
    return response.data;
}

/**
 * Get user's favorite restaurants
 */
export async function getFavorites(): Promise<UserFavorite[]> {
    const response = await restaurantClient.get('/favorites');
    return response.data;
}

/**
 * Add restaurant to favorites
 */
export async function addFavorite(restaurantId: string, notes?: string): Promise<UserFavorite> {
    const response = await restaurantClient.post('/favorites', {
        restaurant_id: restaurantId,
        notes,
    });
    return response.data;
}

/**
 * Remove restaurant from favorites
 */
export async function removeFavorite(favoriteId: string): Promise<void> {
    await restaurantClient.delete(`/favorites/${favoriteId}`);
}

export const restaurantService = {
    getRestaurant,
    getFavorites,
    addFavorite,
    removeFavorite,
};
