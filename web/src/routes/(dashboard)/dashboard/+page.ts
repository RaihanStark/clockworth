import { api } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const today = new Date().toISOString().split('T')[0]; // YYYY-MM-DD

	try {
		const [apiActivities, heatmapData] = await Promise.all([
			api.getActivities(today),
			api.getHeatmap()
		]);

		// Convert API activities to frontend format
		const activities = apiActivities.map((a) => ({
			id: a.id,
			title: a.title,
			startTime: a.start_time,
			endTime: a.end_time,
			category: a.category,
			roi: a.roi as 'High' | 'Neutral' | 'Waste'
		}));

		return {
			activities,
			heatmapData,
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
