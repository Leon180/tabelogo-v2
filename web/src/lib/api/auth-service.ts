import axios from 'axios';
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '@/types/user';

const AUTH_SERVICE_URL = process.env.NEXT_PUBLIC_AUTH_SERVICE_URL || 'http://localhost:8080';

const authClient = axios.create({
    baseURL: AUTH_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests if available
authClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

/**
 * Register a new user
 */
export async function register(data: RegisterRequest): Promise<User> {
    const response = await authClient.post('/api/v1/auth/register', data);
    return response.data;
}

/**
 * Login user
 */
export async function login(data: LoginRequest): Promise<AuthResponse> {
    const response = await authClient.post('/api/v1/auth/login', data);

    // Store tokens
    if (response.data.access_token) {
        localStorage.setItem('access_token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
    }

    return response.data;
}

/**
 * Logout user
 */
export async function logout(): Promise<void> {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
}

/**
 * Validate current token
 */
export async function validateToken(): Promise<User | null> {
    try {
        const response = await authClient.get('/api/v1/auth/validate');
        return response.data;
    } catch (error) {
        return null;
    }
}

export const authService = {
    register,
    login,
    logout,
    validateToken,
};
