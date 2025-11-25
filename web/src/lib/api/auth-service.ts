import axios, { AxiosError } from 'axios';
import type { LoginRequest, RegisterRequest, AuthResponse, User, ValidateTokenResponse } from '@/types/user';

const AUTH_SERVICE_URL = process.env.NEXT_PUBLIC_AUTH_SERVICE_URL || 'http://localhost:8080';

const authClient = axios.create({
    baseURL: AUTH_SERVICE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Flag to prevent multiple simultaneous refresh requests
let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value?: unknown) => void;
    reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null, token: string | null = null) => {
    failedQueue.forEach(({ resolve, reject }) => {
        if (error) {
            reject(error);
        } else {
            resolve(token);
        }
    });
    failedQueue = [];
};

// Helper to clear auth tokens
const clearTokens = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
};

// Add token to requests if available
authClient.interceptors.request.use((config) => {
    // Skip auth header for login, register, and refresh endpoints
    if (config.url?.includes('/login') ||
        config.url?.includes('/register') ||
        config.url?.includes('/refresh')) {
        return config;
    }

    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Add response interceptor to handle token expiration
authClient.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        const originalRequest = error.config as any;

        // If error is 401 and we haven't tried to refresh yet
        if (error.response?.status === 401 && !originalRequest._retry) {
            if (isRefreshing) {
                // If already refreshing, queue this request
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                }).then(token => {
                    originalRequest.headers.Authorization = `Bearer ${token}`;
                    return authClient(originalRequest);
                }).catch(err => {
                    return Promise.reject(err);
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            const refreshToken = localStorage.getItem('refresh_token');
            if (!refreshToken) {
                // No refresh token, clear storage and reject
                clearTokens();
                return Promise.reject(error);
            }

            try {
                const response = await axios.post(`${AUTH_SERVICE_URL}/api/v1/auth/refresh`, {
                    refresh_token: refreshToken,
                });

                const { access_token, refresh_token: new_refresh_token } = response.data;

                // Store new tokens
                localStorage.setItem('access_token', access_token);
                localStorage.setItem('refresh_token', new_refresh_token);

                // Update authorization header
                originalRequest.headers.Authorization = `Bearer ${access_token}`;

                processQueue(null, access_token);
                return authClient(originalRequest);
            } catch (refreshError) {
                processQueue(refreshError as Error, null);
                // Refresh failed, clear tokens
                clearTokens();
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

/**
 * Register a new user
 */
export async function register(data: RegisterRequest): Promise<User> {
    const response = await authClient.post('/api/v1/auth/register', data);
    return response.data.user;
}

/**
 * Login user
 */
export async function login(data: LoginRequest): Promise<AuthResponse> {
    const response = await authClient.post<AuthResponse>('/api/v1/auth/login', data);

    // Store tokens
    if (response.data.access_token) {
        localStorage.setItem('access_token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
    }

    return response.data;
}

/**
 * Logout user - clears tokens from storage
 */
export async function logout(): Promise<void> {
    clearTokens();
}

/**
 * Validate current token and return user data if valid
 */
export async function validateToken(): Promise<User | null> {
    try {
        const response = await authClient.get<ValidateTokenResponse>('/api/v1/auth/validate');
        return response.data.user || null;
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
