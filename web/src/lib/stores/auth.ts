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

	return {
		subscribe,

		async init() {
			if (!browser) return;

			update(state => ({ ...state, loading: true }));

			const token = localStorage.getItem('session_token');
			if (token) {
				try {
					const user = await api.getMe();
					set({ user, loading: false, initialized: true });
				} catch (error) {
					console.error('Failed to get user:', error);
					localStorage.removeItem('session_token');
					set({ user: null, loading: false, initialized: true });
				}
			} else {
				set({ user: null, loading: false, initialized: true });
			}
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

		async setToken(token: string) {
			if (!browser) return;

			localStorage.setItem('session_token', token);
			update(state => ({ ...state, loading: true }));

			try {
				const user = await api.getMe();
				set({ user, loading: false, initialized: true });
				goto('/dashboard');
			} catch (error) {
				console.error('Failed to get user:', error);
				localStorage.removeItem('session_token');
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

			if (browser) {
				localStorage.removeItem('session_token');
			}
			set({ user: null, loading: false, initialized: true });
			goto('/login');
		},

		reset() {
			set({ user: null, loading: false, initialized: false });
		}
	};
}

export const authStore = createAuthStore();
