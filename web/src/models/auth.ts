/**
 * Auth Model
 *
 * Manages authentication state and user data.
 * Communicates with the server API for login/logout.
 */

import { Model, ModelEvents } from '../lib/model.js';
import { api, ApiError } from '../lib/api.js';
import { AuthState } from '../lib/auth-types.js';

export interface AuthModelEvents extends ModelEvents {
  change: void;
  error: Error;
  loginRequired: void;
}

const initialState: AuthState = {
  isAuthenticated: false,
  isLoading: true,
  user: null,
  error: null,
};

export class AuthModel extends Model<AuthState, AuthModelEvents> {
  constructor() {
    super(initialState);
  }

  /**
   * Check authentication status on startup
   */
  async checkAuth(): Promise<void> {
    this.setState({ isLoading: true, error: null });

    try {
      const user = await api.getCurrentUser();
      this.setState({
        isAuthenticated: true,
        isLoading: false,
        user,
      });
    } catch (error) {
      if (error instanceof ApiError && error.statusCode === 401) {
        this.setState({
          isAuthenticated: false,
          isLoading: false,
          user: null,
        });
      } else {
        this.setState({
          isLoading: false,
          error: error instanceof Error ? error.message : 'Unknown error',
        });
        this.emitError(error instanceof Error ? error : new Error('Unknown error'));
      }
    }
  }

  /**
   * Initiate login - redirects to OAuth provider
   */
  login(): void {
    window.location.href = api.getLoginUrl();
  }

  /**
   * Log out the current user
   */
  async logout(): Promise<void> {
    this.setState({ isLoading: true, error: null });

    try {
      await api.logout();
      this.setState({
        isAuthenticated: false,
        isLoading: false,
        user: null,
      });
    } catch (error) {
      this.setState({
        isLoading: false,
        error: error instanceof Error ? error.message : 'Logout failed',
      });
      this.emitError(error instanceof Error ? error : new Error('Logout failed'));
    }
  }

  /**
   * Clear any error state
   */
  clearError(): void {
    this.setState({ error: null });
  }
}
