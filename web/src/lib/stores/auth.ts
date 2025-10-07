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

			// Create new initialization promise
			initPromise = (async () => {
				update(state => ({ ...state, loading: true }));

				// Try to fetch user - cookie is sent automatically via credentials: 'include'
				try {
					const user = await api.getMe();
					set({ user, loading: false, initialized: true });
				} catch (error) {
					// No valid session
					set({ user: null, loading: false, initialized: true });
				}
			})();

			await initPromise;
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
