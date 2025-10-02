import { writable } from 'svelte/store';

export interface Settings {
	hourlyWorth: number;
}

const defaultSettings: Settings = {
	hourlyWorth: 25
};

export const settings = writable<Settings>(defaultSettings);
