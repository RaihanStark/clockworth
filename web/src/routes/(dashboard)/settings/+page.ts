import { api } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const settingsData = await api.getSettings();
		return {
			hourlyWorth: settingsData.hourly_worth
		};
	} catch (error) {
		console.error('Failed to load settings:', error);
		return {
			hourlyWorth: 0 // default value
		};
	}
};
