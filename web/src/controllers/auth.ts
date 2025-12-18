/**
 * Auth Controller
 *
 * Wires together AuthModel and AuthView:
 * - Subscribes to model changes and updates view
 * - Handles view events (login, logout) and updates model
 */

import { Controller } from '../lib/controller.js';
import { AuthModel } from '../models/auth.js';
import { AuthView } from '../views/auth.js';

export class AuthController extends Controller<AuthModel, AuthView> {
  constructor(model: AuthModel, view: AuthView) {
    super(model, view);
  }

  /**
   * Initialize the controller
   */
  override initialize(): void {
    super.initialize();

    // Handle login button click
    this.view.on('login', () => {
      this.model.login();
    });

    // Handle logout button click
    this.view.on('logout', async () => {
      await this.model.logout();
    });

    // Check authentication status on startup
    this.model.checkAuth();
  }

  /**
   * Handle model errors
   */
  protected override onModelError(error: Error): void {
    console.error('Auth error:', error);
    // Error is already in state and will be rendered by the view
  }
}
