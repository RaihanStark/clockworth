import type { PageLoad } from './$types';
import { api } from '$lib/api';
import { settings } from '$lib/stores/settings';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async ({ parent }) => {
	try {
		// Ensure user is authenticated
		await parent();
		
		// Fetch hourly worth from API
		const hourlyWorth = await api.getHourlyWorth();
		
		// Update the store with the fetched value
		settings.update((s) => ({ ...s, hourlyWorth }));
		
		return {
			hourlyWorth
		};
	} catch (error) {
		// If there's an error (like user not authenticated), redirect to login
		redirect(302, '/login');
	}
};