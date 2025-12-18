/**
 * API client for browser
 *
 * This is a browser-specific implementation that mirrors the
 * TypeScript client library in client/src/client.ts
 */

export interface User {
  id: string;
  email: string;
  name?: string;
  picture?: string;
}

export interface HealthResponse {
  status: 'healthy' | 'degraded' | 'unhealthy';
  timestamp: string;
  version?: string;
}

export interface LogoutResponse {
  success: boolean;
  message?: string;
}

export interface ErrorResponse {
  error: string;
  message: string;
  details?: Record<string, unknown>;
}

export class ApiError extends Error {
  constructor(
    public readonly statusCode: number,
    public readonly errorCode: string,
    message: string,
    public readonly details?: Record<string, unknown>
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Browser API client
 * Uses same-origin requests (no baseUrl needed since served from same server)
 */
export class ApiClient {
  /**
   * Check server health
   */
  async getHealth(): Promise<HealthResponse> {
    return this.request<HealthResponse>('GET', '/health');
  }

  /**
   * Get the login URL
   */
  getLoginUrl(): string {
    return '/auth/login';
  }

  /**
   * Get the current authenticated user
   */
  async getCurrentUser(): Promise<User> {
    return this.request<User>('GET', '/auth/me');
  }

  /**
   * Check if authenticated
   */
  async isAuthenticated(): Promise<boolean> {
    try {
      await this.getCurrentUser();
      return true;
    } catch (error) {
      if (error instanceof ApiError && error.statusCode === 401) {
        return false;
      }
      throw error;
    }
  }

  /**
   * Log out
   */
  async logout(): Promise<LogoutResponse> {
    return this.request<LogoutResponse>('POST', '/auth/logout');
  }

  /**
   * Make an HTTP request
   */
  private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const headers: Record<string, string> = {
      'Accept': 'application/json',
    };

    if (body !== undefined) {
      headers['Content-Type'] = 'application/json';
    }

    const response = await fetch(path, {
      method,
      headers,
      credentials: 'same-origin',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      let errorResponse: ErrorResponse;
      try {
        errorResponse = await response.json();
      } catch {
        errorResponse = {
          error: 'unknown_error',
          message: response.statusText,
        };
      }

      throw new ApiError(
        response.status,
        errorResponse.error,
        errorResponse.message,
        errorResponse.details
      );
    }

    return response.json();
  }
}

// Singleton instance for convenience
export const api = new ApiClient();
