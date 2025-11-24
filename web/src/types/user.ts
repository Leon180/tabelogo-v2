// User and authentication types
export interface User {
    id: string;
    email: string;
    username: string;
    role: 'admin' | 'user' | 'guest';
    is_active: boolean;
    email_verified: boolean;
    created_at: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
    username: string;
}

export interface AuthResponse {
    access_token: string;
    refresh_token: string;
    user_id: string;
    username: string;
}
