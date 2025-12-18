/**
 * Auth View
 *
 * Renders authentication UI:
 * - Loading state
 * - Login button (when not authenticated)
 * - User info + logout button (when authenticated)
 * - Error messages
 */

import { View, ViewEvents } from '../lib/view.js';
import { AuthState } from '../lib/auth-types.js';

export interface AuthViewEvents extends ViewEvents {
  login: void;
  logout: void;
}

export class AuthView extends View<AuthViewEvents> {
  constructor(container: HTMLElement | string) {
    super(container);
  }

  render(state: AuthState): void {
    // Clear previous content and event listeners
    this.element.innerHTML = '';

    if (state.isLoading) {
      this.renderLoading();
      return;
    }

    if (state.error) {
      this.renderError(state.error);
    }

    if (state.isAuthenticated && state.user) {
      this.renderUserInfo(state.user.name || state.user.email, state.user.picture);
    } else {
      this.renderLoginButton();
    }
  }

  private renderLoading(): void {
    const loading = this.createElement('div', { class: 'auth-loading' }, ['Loading...']);
    this.element.appendChild(loading);
  }

  private renderError(message: string): void {
    const error = this.createElement('div', { class: 'auth-error' }, [message]);
    this.element.appendChild(error);
  }

  private renderLoginButton(): void {
    const container = this.createElement('div', { class: 'auth-login' });

    const button = this.createElement('button', {
      class: 'btn btn-primary',
      type: 'button',
    }, ['Sign in with Google']);

    this.addDOMListener(button, 'click', () => {
      this.emit('login', undefined);
    });

    container.appendChild(button);
    this.element.appendChild(container);
  }

  private renderUserInfo(name: string, picture?: string): void {
    const container = this.createElement('div', { class: 'auth-user' });

    // User info
    const userInfo = this.createElement('div', { class: 'user-info' });

    if (picture) {
      const img = this.createElement('img', {
        class: 'user-avatar',
        src: picture,
        alt: name,
      });
      userInfo.appendChild(img);
    }

    const nameEl = this.createElement('span', { class: 'user-name' }, [name]);
    userInfo.appendChild(nameEl);

    container.appendChild(userInfo);

    // Logout button
    const logoutBtn = this.createElement('button', {
      class: 'btn btn-secondary',
      type: 'button',
    }, ['Sign out']);

    this.addDOMListener(logoutBtn, 'click', () => {
      this.emit('logout', undefined);
    });

    container.appendChild(logoutBtn);
    this.element.appendChild(container);
  }
}
