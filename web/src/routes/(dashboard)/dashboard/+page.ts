import { api } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const today = new Date().toISOString().split('T')[0]; // YYYY-MM-DD

	// Retry helper for transient failures
	async function fetchWithRetry<T>(
		fn: () => Promise<T>,
		retries = 2,
		delay = 300
	): Promise<T | null> {
		for (let i = 0; i <= retries; i++) {
			try {
				return await fn();
			} catch (error) {
				if (i < retries) {
					await new Promise((resolve) => setTimeout(resolve, delay));
					delay *= 2;
				} else {
					console.error('Failed after retries:', error);
					return null;
				}
			}
		}
		return null;
	}

	try {
		const [apiActivities, heatmapData] = await Promise.all([
			fetchWithRetry(() => api.getActivities(today)),
			fetchWithRetry(() => api.getHeatmap())
		]);

		// Convert API activities to frontend format
		const activities = (apiActivities || []).map((a) => ({
			id: a.id,
			title: a.title,
			startTime: a.start_time,
			endTime: a.end_time,
			category: a.category,
			roi: a.roi as 'High' | 'Neutral' | 'Waste'
		}));

		return {
			activities,
			heatmapData: heatmapData || [],
			date: today
		};
	} catch (error) {
		console.error('Failed to load activities:', error);
		return {
			activities: [],
			heatmapData: [],
			date: today
		};
	}
};
