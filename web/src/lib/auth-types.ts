/**
 * Shared auth types
 *
 * These types are shared between Models, Views, and Controllers.
 * Keeping them in lib/ maintains proper MVC boundaries.
 */

import { User } from './api.js';

export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  error: string | null;
}
