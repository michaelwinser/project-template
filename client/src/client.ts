import {
  User,
  HealthResponse,
  LogoutResponse,
  ErrorResponse,
  ApiError,
} from './types';

export interface ClientConfig {
  baseUrl: string;
  credentials?: RequestCredentials;
}

/**
 * API client for project-template server
 */
export class ApiClient {
  private readonly baseUrl: string;
  private readonly credentials: RequestCredentials;

  constructor(config: ClientConfig) {
    this.baseUrl = config.baseUrl.replace(/\/$/, ''); // Remove trailing slash
    this.credentials = config.credentials ?? 'include'; // Include cookies by default
  }

  /**
   * Check server health
   */
  async getHealth(): Promise<HealthResponse> {
    return this.request<HealthResponse>('GET', '/health');
  }

  /**
   * Get the login URL for OAuth
   * Note: This returns the URL to redirect to, not the actual login
   */
  getLoginUrl(): string {
    return `${this.baseUrl}/auth/login`;
  }

  /**
   * Get the current authenticated user
   * @throws ApiError if not authenticated
   */
  async getCurrentUser(): Promise<User> {
    return this.request<User>('GET', '/auth/me');
  }

  /**
   * Check if the user is authenticated
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
   * Log out the current user
   * @throws ApiError if not authenticated
   */
  async logout(): Promise<LogoutResponse> {
    return this.request<LogoutResponse>('POST', '/auth/logout');
  }

  /**
   * Make an HTTP request to the API
   */
  private async request<T>(
    method: string,
    path: string,
    body?: unknown
  ): Promise<T> {
    const url = `${this.baseUrl}${path}`;

    const headers: Record<string, string> = {
      'Accept': 'application/json',
    };

    if (body !== undefined) {
      headers['Content-Type'] = 'application/json';
    }

    const response = await fetch(url, {
      method,
      headers,
      credentials: this.credentials,
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
