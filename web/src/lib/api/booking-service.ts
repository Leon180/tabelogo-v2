import axios from 'axios';
import type { Booking, CreateBookingRequest } from '@/types/booking';

const BOOKING_SERVICE_URL = process.env.NEXT_PUBLIC_BOOKING_SERVICE_URL || 'http://localhost:8083';

const bookingClient = axios.create({
    baseURL: BOOKING_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests
bookingClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

/**
 * Create a new booking
 */
export async function createBooking(data: CreateBookingRequest): Promise<Booking> {
    const response = await bookingClient.post('/bookings', data);
    return response.data;
}

/**
 * Get user's bookings
 */
export async function getBookings(): Promise<Booking[]> {
    const response = await bookingClient.get('/bookings');
    return response.data;
}

/**
 * Get booking by ID
 */
export async function getBooking(id: string): Promise<Booking> {
    const response = await bookingClient.get(`/bookings/${id}`);
    return response.data;
}

/**
 * Cancel a booking
 */
export async function cancelBooking(id: string): Promise<Booking> {
    const response = await bookingClient.patch(`/bookings/${id}/cancel`);
    return response.data;
}

export const bookingService = {
    createBooking,
    getBookings,
    getBooking,
    cancelBooking,
};
