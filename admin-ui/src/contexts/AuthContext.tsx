"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { setRefreshTokenFunction } from '@/lib/api/client';
import { decodeJWT, type JWTClaims } from '@/lib/jwt';

interface User {
  user_id: string;
  email: string;
  name: string;
  organization_id?: string;
  role?: string;
  permissions?: string[];
}

interface OrganizationContext {
  organization_id: string;
  organization_type: string; // manufacturer, hospital, distributor, dealer, supplier, imaging_center
  role: string;
}

interface AuthContextType {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  organizationContext: OrganizationContext | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (accessToken: string, refreshToken: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshAccessToken: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [refreshToken, setRefreshToken] = useState<string | null>(null);
  const [organizationContext, setOrganizationContext] = useState<OrganizationContext | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  // Load tokens from localStorage on mount
  useEffect(() => {
    const loadTokens = async () => {
      try {
        const storedAccessToken = localStorage.getItem('access_token');
        const storedRefreshToken = localStorage.getItem('refresh_token');

        if (storedAccessToken && storedRefreshToken) {
          // Validate token before using it
          if (!isTokenValid(storedAccessToken)) {
            console.log('[Auth] Stored token is invalid or expired');
            clearTokens();
            setIsLoading(false);
            return;
          }

          setAccessToken(storedAccessToken);
          setRefreshToken(storedRefreshToken);

          // Extract organization context from JWT
          extractOrganizationContext(storedAccessToken);

          // Fetch user info
          await fetchUserInfo(storedAccessToken);
        }
      } catch (error) {
        console.error('Failed to load tokens:', error);
        clearTokens();
      } finally {
        setIsLoading(false);
      }
    };

    loadTokens();
  }, []);

  // Register token refresh function with API client
  useEffect(() => {
    setRefreshTokenFunction(refreshAccessToken);
  }, [refreshToken]);

  // Extract organization context from JWT token
  const extractOrganizationContext = (token: string) => {
    const claims = decodeJWT(token);
    
    if (claims && claims.organization_id && claims.organization_type) {
      setOrganizationContext({
        organization_id: claims.organization_id,
        organization_type: claims.organization_type,
        role: claims.role,
      });
      
      console.log('[AUTH] Organization context extracted:', {
        org_id: claims.organization_id,
        org_type: claims.organization_type,
        role: claims.role,
      });
    } else {
      console.warn('[AUTH] Failed to extract organization context from token');
      setOrganizationContext(null);
    }
  };

  // Fetch user information
  const fetchUserInfo = async (token: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/v1/auth/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const userData = await response.json();
        setUser(userData);
        // Store user data in localStorage for components that need it
        localStorage.setItem('user', JSON.stringify(userData));
      } else {
        throw new Error('Failed to fetch user info');
      }
    } catch (error) {
      console.error('Failed to fetch user info:', error);
      throw error;
    }
  };

  // Login function
  const login = async (newAccessToken: string, newRefreshToken: string) => {
    setAccessToken(newAccessToken);
    setRefreshToken(newRefreshToken);

    // Store in localStorage
    localStorage.setItem('access_token', newAccessToken);
    localStorage.setItem('refresh_token', newRefreshToken);

    // Extract organization context from JWT
    extractOrganizationContext(newAccessToken);

    // Fetch user info
    await fetchUserInfo(newAccessToken);
  };

  // Logout function
  const logout = async () => {
    try {
      if (refreshToken) {
        await fetch(`${API_BASE_URL}/v1/auth/logout`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      clearTokens();
      router.push('/login');
    }
  };

  // Clear tokens
  const clearTokens = () => {
    setUser(null);
    setAccessToken(null);
    setRefreshToken(null);
    setOrganizationContext(null);
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user'); // Clear user data too
  };

  // Refresh access token
  const refreshAccessToken = async (): Promise<boolean> => {
    if (!refreshToken) {
      return false;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/v1/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (response.ok) {
        const data = await response.json();
        await login(data.access_token, data.refresh_token);
        return true;
      } else {
        clearTokens();
        return false;
      }
    } catch (error) {
      console.error('Token refresh error:', error);
      clearTokens();
      return false;
    }
  };

  const value: AuthContextType = {
    user,
    accessToken,
    refreshToken,
    organizationContext,
    isLoading,
    isAuthenticated: !!user && !!accessToken,
    login,
    logout,
    refreshAccessToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
