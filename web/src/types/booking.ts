// Booking types based on database schema
export interface Booking {
    id: string;
    user_id: string;
    restaurant_id: string;
    booking_date: string;
    party_size: number;
    status: 'pending' | 'confirmed' | 'cancelled' | 'completed' | 'no_show';
    external_booking_id?: string;
    external_service?: string;
    special_requests?: string;
    customer_name: string;
    customer_phone: string;
    customer_email: string;
    notes?: string;
    confirmation_code?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateBookingRequest {
    restaurant_id: string;
    booking_date: string;
    party_size: number;
    special_requests?: string;
    customer_name: string;
    customer_phone: string;
    customer_email: string;
}
