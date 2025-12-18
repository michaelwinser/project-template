/**
 * Base View class for the MVC framework.
 *
 * Views are responsible for:
 * - Rendering state to the DOM
 * - Capturing user interactions
 * - Emitting UI events (clicks, form submissions, etc.)
 *
 * Views should NOT:
 * - Hold application state
 * - Perform business logic
 * - Import from models/ (only receive data via render())
 * - Import from controllers/
 */

import { EventEmitter } from './events.js';

export interface ViewEvents {
  [key: string]: unknown;
}

export abstract class View<Events extends ViewEvents = ViewEvents> extends EventEmitter<Events> {
  protected element: HTMLElement;
  private eventCleanup: Array<() => void> = [];

  constructor(container: HTMLElement | string) {
    super();
    if (typeof container === 'string') {
      const el = document.querySelector(container);
      if (!el) {
        throw new Error(`Container not found: ${container}`);
      }
      this.element = el as HTMLElement;
    } else {
      this.element = container;
    }
  }

  /**
   * Render the view. Called by Controller when model changes.
   * Subclasses should override this.
   */
  abstract render(data: unknown): void;

  /**
   * Clean up event listeners and resources
   */
  destroy(): void {
    this.eventCleanup.forEach(cleanup => cleanup());
    this.eventCleanup = [];
    this.removeAllListeners();
  }

  /**
   * Helper to add DOM event listeners with automatic cleanup
   */
  protected addDOMListener<K extends keyof HTMLElementEventMap>(
    element: HTMLElement,
    event: K,
    handler: (e: HTMLElementEventMap[K]) => void,
    options?: AddEventListenerOptions
  ): void {
    element.addEventListener(event, handler, options);
    this.eventCleanup.push(() => element.removeEventListener(event, handler, options));
  }

  /**
   * Helper to query within this view's element
   */
  protected query<T extends HTMLElement>(selector: string): T | null {
    return this.element.querySelector(selector);
  }

  /**
   * Helper to query all within this view's element
   */
  protected queryAll<T extends HTMLElement>(selector: string): NodeListOf<T> {
    return this.element.querySelectorAll(selector);
  }

  /**
   * Set innerHTML safely (use with caution - prefer DOM methods for user content)
   */
  protected setHTML(html: string): void {
    this.element.innerHTML = html;
  }

  /**
   * Create an element with attributes and children
   */
  protected createElement<K extends keyof HTMLElementTagNameMap>(
    tag: K,
    attrs?: Record<string, string>,
    children?: Array<HTMLElement | string>
  ): HTMLElementTagNameMap[K] {
    const el = document.createElement(tag);
    if (attrs) {
      Object.entries(attrs).forEach(([key, value]) => {
        el.setAttribute(key, value);
      });
    }
    if (children) {
      children.forEach(child => {
        if (typeof child === 'string') {
          el.appendChild(document.createTextNode(child));
        } else {
          el.appendChild(child);
        }
      });
    }
    return el;
  }
}
