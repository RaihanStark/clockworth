import { writable } from 'svelte/store';

export type ROICategory = 'High' | 'Neutral' | 'Waste';

export interface Activity {
	id: string;
	title: string;
	startTime: string; // HH:mm format
	endTime: string; // HH:mm format
	category: string;
	roi: ROICategory;
}

const mockActivities: Activity[] = [
	{
		id: '1',
		title: 'Deep Work - Code Review',
		startTime: '09:00',
		endTime: '11:30',
		category: 'Development',
		roi: 'High'
	},
	{
		id: '2',
		title: 'Team Meeting',
		startTime: '14:00',
		endTime: '15:00',
		category: 'Meetings',
		roi: 'Neutral'
	},
	{
		id: '3',
		title: 'Social Media Browsing',
		startTime: '16:00',
		endTime: '17:30',
		category: 'Personal',
		roi: 'Waste'
	}
];

function createActivitiesStore() {
	const { subscribe, set, update } = writable<Activity[]>(mockActivities);

	return {
		subscribe,
		add: (activity: Activity) => update((activities) => [...activities, activity]),
		update: (id: string, updates: Partial<Activity>) =>
			update((activities) =>
				activities.map((a) => (a.id === id ? { ...a, ...updates } : a))
			),
		remove: (id: string) => update((activities) => activities.filter((a) => a.id !== id)),
		reset: () => set(mockActivities)
	};
}

export const activities = createActivitiesStore();
