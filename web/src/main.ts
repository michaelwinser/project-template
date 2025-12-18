/**
 * Entry point for the web application
 */

import { App } from './app.js';

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', () => {
  const app = new App();
  app.start();

  // Expose for debugging in development
  if (import.meta.env?.DEV) {
    (window as unknown as { app: App }).app = app;
  }
});
