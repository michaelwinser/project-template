/**
 * Base Controller class for the MVC framework.
 *
 * Controllers are responsible for:
 * - Wiring Models and Views together
 * - Handling View events and updating Models
 * - Subscribing to Model changes and updating Views
 * - Application flow and coordination
 *
 * Controllers CAN:
 * - Import from models/
 * - Import from views/
 * - Orchestrate multiple models and views
 */

import { Model, ModelEvents } from './model.js';
import { View, ViewEvents } from './view.js';

export abstract class Controller<
  M extends Model<unknown, ModelEvents>,
  V extends View<ViewEvents>
> {
  protected model: M;
  protected view: V;
  private unsubscribers: Array<() => void> = [];

  constructor(model: M, view: V) {
    this.model = model;
    this.view = view;
  }

  /**
   * Initialize the controller - wire up event handlers
   * Subclasses should call super.initialize() and add their own wiring
   */
  initialize(): void {
    // Subscribe to model changes and re-render view
    const unsubChange = this.model.on('change', () => {
      this.onModelChange();
    });
    this.unsubscribers.push(unsubChange);

    // Subscribe to model errors
    const unsubError = this.model.on('error', (error) => {
      this.onModelError(error);
    });
    this.unsubscribers.push(unsubError);

    // Initial render
    this.onModelChange();
  }

  /**
   * Called when model emits 'change' event
   * Default behavior: render the view with current model state
   */
  protected onModelChange(): void {
    this.view.render(this.model.getState());
  }

  /**
   * Called when model emits 'error' event
   * Subclasses should override to handle errors appropriately
   */
  protected onModelError(error: Error): void {
    console.error('Model error:', error);
  }

  /**
   * Helper to subscribe to view events
   */
  protected onViewEvent<K extends keyof ViewEvents>(
    event: K,
    handler: (data: ViewEvents[K]) => void
  ): void {
    const unsub = this.view.on(event, handler);
    this.unsubscribers.push(unsub);
  }

  /**
   * Clean up all subscriptions
   */
  destroy(): void {
    this.unsubscribers.forEach(unsub => unsub());
    this.unsubscribers = [];
    this.view.destroy();
  }
}
