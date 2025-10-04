import { redirect } from '@sveltejs/kit';
import { authStore } from '$lib/stores/auth';
import { get } from 'svelte/store';
import { browser } from '$app/environment';
import type { LayoutLoad } from './$types';

export const ssr = false; // Disable SSR for dashboard

export const load: LayoutLoad = async () => {
	if (!browser) return {};

	const auth = get(authStore);

	// Wait for auth to initialize if not yet done
	if (!auth.initialized) {
		await authStore.init();
	}

	const updatedAuth = get(authStore);

	if (!updatedAuth.user) {
		throw redirect(307, '/login');
	}

	return {
		user: updatedAuth.user
	};
};
