// User and authentication types (matches backend UserResponse)
export interface User {
    id: string;
    email: string;
    username: string;
    role: 'admin' | 'user' | 'guest';
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

// Matches backend LoginResponse
export interface AuthResponse {
    access_token: string;
    refresh_token: string;
    user: User;
}

// Matches backend ValidateTokenResponse
export interface ValidateTokenResponse {
    valid: boolean;
    user?: User;
}
