/**
 * Base Model class for the MVC framework.
 *
 * Models are responsible for:
 * - Holding application state
 * - Business logic and data manipulation
 * - Emitting events when state changes
 *
 * Models should NOT:
 * - Access the DOM
 * - Know about Views or Controllers
 * - Import from views/ or controllers/
 */

import { EventEmitter } from './events.js';

export interface ModelEvents {
  change: void;
  error: Error;
}

export abstract class Model<
  State,
  Events extends ModelEvents = ModelEvents
> extends EventEmitter<Events> {
  protected state: State;

  constructor(initialState: State) {
    super();
    this.state = initialState;
  }

  /**
   * Get the current state (read-only)
   */
  getState(): Readonly<State> {
    return this.state;
  }

  /**
   * Update state and emit change event
   */
  protected setState(partial: Partial<State>): void {
    this.state = { ...this.state, ...partial };
    this.emit('change' as keyof Events, undefined as Events[keyof Events]);
  }

  /**
   * Emit an error event
   */
  protected emitError(error: Error): void {
    this.emit('error' as keyof Events, error as Events[keyof Events]);
  }
}
