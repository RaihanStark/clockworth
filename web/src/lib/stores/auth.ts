import { writable } from 'svelte/store';
import { api, type User } from '$lib/api';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

interface AuthState {
	user: User | null;
	loading: boolean;
	initialized: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: false,
		initialized: false
	});

	let initPromise: Promise<void> | null = null;

	return {
		subscribe,

		async init() {
			if (!browser) return;

			// Get current state to check if already initialized
			let currentState: AuthState;
			const unsubscribe = subscribe(state => { currentState = state; })();

			// If already initialized, skip
			if (currentState!.initialized) {
				return;
			}

			// If already initializing, return the existing promise
			if (initPromise) {
				return initPromise;
			}

			// Create and assign promise atomically to prevent race condition
			return await (initPromise = (async () => {
				update(state => ({ ...state, loading: true }));

				// Retry logic for server errors
				let retries = 3;
				let delay = 300; // Start with 300ms

				while (retries > 0) {
					try {
						const user = await api.getMe();
						set({ user, loading: false, initialized: true });
						return;
					} catch (error) {
						const errorMessage = error instanceof Error ? error.message : 'Unknown error';

						// If it's a retryable error and we have retries left, retry
						const retryableErrors = ['SERVER_ERROR', 'REQUEST_CANCELLED', 'SERVICE_UNAVAILABLE'];
						if (retryableErrors.includes(errorMessage) && retries > 1) {
							retries--;
							await new Promise(resolve => setTimeout(resolve, delay));
							delay *= 2; // Exponential backoff
							continue;
						}

						// For UNAUTHORIZED or any other error after retries, treat as no session
						set({ user: null, loading: false, initialized: true });
						return;
					}
				}
			})());
		},

		async login() {
			try {
				const url = await api.getLoginUrl();
				if (browser) {
					window.location.href = url;
				}
			} catch (error) {
				console.error('Failed to get login URL:', error);
				throw error;
			}
		},

		async loginWithEmail(email: string, password: string) {
			update(state => ({ ...state, loading: true }));
			try {
				const user = await api.loginWithEmail(email, password);
				set({ user, loading: false, initialized: true });
				goto('/dashboard');
			} catch (error) {
				update(state => ({ ...state, loading: false }));
				throw error;
			}
		},

		async registerWithEmail(email: string, password: string, name?: string) {
			update(state => ({ ...state, loading: true }));
			try {
				const user = await api.registerWithEmail(email, password, name);
				set({ user, loading: false, initialized: true });
				goto('/dashboard');
			} catch (error) {
				update(state => ({ ...state, loading: false }));
				throw error;
			}
		},

		async setToken(_token: string) {
			if (!browser) return;

			// Cookie is already set by backend, just fetch user info
			update(state => ({ ...state, loading: true }));

			try {
				const user = await api.getMe();
				set({ user, loading: false, initialized: true });
				goto('/dashboard');
			} catch (error) {
				console.error('Failed to get user:', error);
				set({ user: null, loading: false, initialized: true });
				throw error;
			}
		},

		async logout() {
			try {
				await api.logout();
			} catch (error) {
				console.error('Failed to logout:', error);
			}

			// Cookie is cleared by backend
			set({ user: null, loading: false, initialized: true });
			goto('/login');
		},

		reset() {
			initPromise = null;
			set({ user: null, loading: false, initialized: false });
		}
	};
}

export const authStore = createAuthStore();
