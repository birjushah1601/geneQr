/**
 * JWT Decoder Utility
 * Decodes JWT tokens to extract claims without verification
 * (verification happens on backend)
 */

export interface JWTClaims {
  user_id: string;
  email?: string;
  name?: string;
  organization_id: string;
  organization_type: string; // manufacturer, hospital, Channel Partner, etc.
  role: string;
  permissions?: string[];
  exp?: number;
  iat?: number;
}

/**
 * Decode JWT token and extract claims
 * @param token JWT access token
 * @returns Decoded claims or null if invalid
 */
export function decodeJWT(token: string): JWTClaims | null {
  try {
    // JWT format: header.payload.signature
    const parts = token.split('.');
    
    if (parts.length !== 3) {
      console.error('Invalid JWT format');
      return null;
    }

    // Decode payload (base64url)
    const payload = parts[1];
    
    // Replace URL-safe characters
    const base64 = payload.replace(/-/g, '+').replace(/_/g, '/');
    
    // Add padding if needed
    const paddedBase64 = base64.padEnd(
      base64.length + ((4 - (base64.length % 4)) % 4),
      '='
    );

    // Decode from base64
    const jsonString = atob(paddedBase64);
    
    // Parse JSON
    const claims = JSON.parse(jsonString) as JWTClaims;
    
    return claims;
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
}

/**
 * Check if JWT token is expired
 * @param token JWT access token
 * @returns true if expired, false otherwise
 */
export function isTokenExpired(token: string): boolean {
  const claims = decodeJWT(token);
  
  if (!claims || !claims.exp) {
    return true;
  }

  // exp is in seconds, Date.now() is in milliseconds
  const expirationTime = claims.exp * 1000;
  const currentTime = Date.now();

  return currentTime >= expirationTime;
}

/**
 * Get organization type from token
 * @param token JWT access token
 * @returns Organization type or null
 */
export function getOrganizationType(token: string): string | null {
  const claims = decodeJWT(token);
  return claims?.organization_type || null;
}

/**
 * Get organization ID from token
 * @param token JWT access token
 * @returns Organization ID or null
 */
export function getOrganizationID(token: string): string | null {
  const claims = decodeJWT(token);
  return claims?.organization_id || null;
}

/**
 * Get user role from token
 * @param token JWT access token
 * @returns User role or null
 */
export function getUserRole(token: string): string | null {
  const claims = decodeJWT(token);
  return claims?.role || null;
}
