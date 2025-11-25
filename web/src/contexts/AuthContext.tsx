'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { authService } from '@/lib/api/auth-service';
import type { User, LoginRequest, RegisterRequest } from '@/types/user';

interface AuthContextType {
    user: User | null;
    isLoading: boolean;
    login: (data: LoginRequest) => Promise<void>;
    register: (data: RegisterRequest) => Promise<void>;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        checkAuth();
    }, []);

    const checkAuth = async () => {
        try {
            const token = localStorage.getItem('access_token');
            if (token) {
                const userData = await authService.validateToken();
                setUser(userData);
            }
        } catch (error) {
            console.error('Auth check failed:', error);
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
        } finally {
            setIsLoading(false);
        }
    };

    const login = async (data: LoginRequest) => {
        try {
            const response = await authService.login(data);
            // After login, validate token to get full user details if needed, 
            // or use the user data from response if available.
            // The AuthResponse has user_id and username, but User type has more fields.
            // Let's fetch full user profile or construct it.
            // For now, let's try to validate token to get the User object.
            const userData = await authService.validateToken();
            setUser(userData);
            router.push('/');
        } catch (error) {
            throw error;
        }
    };

    const register = async (data: RegisterRequest) => {
        try {
            await authService.register(data);
            // Auto login after register? Or redirect to login?
            // For now, redirect to login
            router.push('/auth/login');
        } catch (error) {
            throw error;
        }
    };

    const logout = async () => {
        try {
            await authService.logout();
            setUser(null);
            router.push('/auth/login');
        } catch (error) {
            console.error('Logout failed:', error);
        }
    };

    return (
        <AuthContext.Provider value={{ user, isLoading, login, register, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
