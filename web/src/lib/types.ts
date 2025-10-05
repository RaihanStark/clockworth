export type ROICategory = 'High' | 'Neutral' | 'Waste';

export interface Activity {
	id: string;
	title: string;
	startTime: string; // HH:mm format
	endTime: string; // HH:mm format
	category: string;
	roi: ROICategory;
}
