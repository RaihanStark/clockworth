const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface User {
	id: string;
	email: string;
	name: string;
	picture: string;
}

export interface LoginResponse {
	url: string;
}

export interface MeResponse {
	id: string;
	email: string;
	name: string;
	picture: string;
}

export interface SettingsResponse {
	hourly_worth: number;
}

export interface UpdateHourlyWorthRequest {
	hourly_worth: number;
}

export interface Activity {
	id: string;
	title: string;
	start_time: string; // HH:MM format
	end_time: string; // HH:MM format
	category: string;
	roi: string;
	activity_date: string; // YYYY-MM-DD format
	created_at: string;
	updated_at: string;
}

export interface CreateActivityRequest {
	title: string;
	start_time: string;
	end_time: string;
	category: string;
	roi: string;
	activity_date: string;
}

export interface UpdateActivityRequest {
	title: string;
	start_time: string;
	end_time: string;
	category: string;
	roi: string;
}

export interface HeatmapDay {
	date: string;
	status: 'positive' | 'negative' | 'neutral';
}

class ApiClient {
	private baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	private getHeaders(): HeadersInit {
		const headers: HeadersInit = {
			'Content-Type': 'application/json'
		};

		// Cookie is automatically sent with credentials: 'include'
		// No need to manually set Authorization header

		return headers;
	}

	async getLoginUrl(): Promise<string> {
		const response = await fetch(`${this.baseUrl}/auth/google/login`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
			throw new Error('Failed to get login URL');
		}

		const data: LoginResponse = await response.json();
		return data.url;
	}

	async getMe(): Promise<User> {
		const response = await fetch(`${this.baseUrl}/api/me`, {
			method: 'GET',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			// 408 = Request cancelled (page refresh) - treat as retryable
			if (response.status === 408) {
				throw new Error('REQUEST_CANCELLED');
			}
			// 401 = Not authenticated
			if (response.status === 401) {
				throw new Error('UNAUTHORIZED');
			}
			// For 500+ errors, throw a different error so we can retry
			if (response.status >= 500) {
				throw new Error('SERVER_ERROR');
			}
			throw new Error('Failed to get user info');
		}

		const data: MeResponse = await response.json();
		return {
			id: data.id,
			email: data.email,
			name: data.name,
			picture: data.picture
		};
	}

	async logout(): Promise<void> {
		const response = await fetch(`${this.baseUrl}/api/logout`, {
			method: 'POST',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to logout');
		}

		// Cookie is cleared by the backend
	}

	async getSettings(): Promise<SettingsResponse> {
		const response = await fetch(`${this.baseUrl}/api/settings`, {
			method: 'GET',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to get settings');
		}

		return await response.json();
	}

	async updateHourlyWorth(hourlyWorth: number): Promise<SettingsResponse> {
		const response = await fetch(`${this.baseUrl}/api/settings/hourly-worth`, {
			method: 'PUT',
			headers: this.getHeaders(),
			credentials: 'include',
			body: JSON.stringify({ hourly_worth: hourlyWorth })
		});

		if (!response.ok) {
			throw new Error('Failed to update hourly worth');
		}

		return await response.json();
	}

	// Activities API
	async getActivities(date: string): Promise<Activity[]> {
		const response = await fetch(`${this.baseUrl}/api/activities?date=${date}`, {
			method: 'GET',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to get activities');
		}

		return await response.json();
	}

	async getActivity(id: string): Promise<Activity> {
		const response = await fetch(`${this.baseUrl}/api/activities/${id}`, {
			method: 'GET',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to get activity');
		}

		return await response.json();
	}

	async createActivity(request: CreateActivityRequest): Promise<Activity> {
		const response = await fetch(`${this.baseUrl}/api/activities`, {
			method: 'POST',
			headers: this.getHeaders(),
			credentials: 'include',
			body: JSON.stringify(request)
		});

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to create activity' }));
			throw new Error(errorData.message || 'Failed to create activity');
		}

		return await response.json();
	}

	async updateActivity(id: string, request: UpdateActivityRequest): Promise<Activity> {
		const response = await fetch(`${this.baseUrl}/api/activities/${id}`, {
			method: 'PUT',
			headers: this.getHeaders(),
			credentials: 'include',
			body: JSON.stringify(request)
		});

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to update activity' }));
			throw new Error(errorData.message || 'Failed to update activity');
		}

		return await response.json();
	}

	async deleteActivity(id: string): Promise<void> {
		const response = await fetch(`${this.baseUrl}/api/activities/${id}`, {
			method: 'DELETE',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to delete activity');
		}
	}

	async getHeatmap(): Promise<HeatmapDay[]> {
		const response = await fetch(`${this.baseUrl}/api/activities/heatmap`, {
			method: 'GET',
			headers: this.getHeaders(),
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to get heatmap data');
		}

		return await response.json();
	}
}

export const api = new ApiClient(API_URL);
