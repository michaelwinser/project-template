/**
 * Main Application
 *
 * Sets up all models, views, and controllers.
 * This is the composition root of the application.
 */

import { AuthModel } from './models/auth.js';
import { AuthView } from './views/auth.js';
import { AuthController } from './controllers/auth.js';

export class App {
  private authModel: AuthModel;
  private authView: AuthView;
  private authController: AuthController;

  constructor() {
    // Create models
    this.authModel = new AuthModel();

    // Create views
    this.authView = new AuthView('#auth-container');

    // Create controllers
    this.authController = new AuthController(this.authModel, this.authView);
  }

  /**
   * Start the application
   */
  start(): void {
    console.log('App starting...');
    this.authController.initialize();
    console.log('App started');
  }

  /**
   * Clean up resources
   */
  destroy(): void {
    this.authController.destroy();
  }
}
