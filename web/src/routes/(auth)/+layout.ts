import { redirect } from '@sveltejs/kit';
import { authStore } from '$lib/stores/auth';
import { get } from 'svelte/store';
import { browser } from '$app/environment';
import type { LayoutLoad } from './$types';

export const ssr = false; // Disable SSR for auth pages

export const load: LayoutLoad = async () => {
	if (!browser) return {};

	// Always ensure auth is initialized (safe to call multiple times)
	await authStore.init();

	// Get the current auth state
	const auth = get(authStore);

	// Redirect to dashboard if already logged in
	if (auth.user) {
		throw redirect(307, '/dashboard');
	}

	return {};
};
