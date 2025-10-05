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

class ApiClient {
	private baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	private getHeaders(): HeadersInit {
		const headers: HeadersInit = {
			'Content-Type': 'application/json'
		};

		const token = localStorage.getItem('session_token');
		if (token) {
			headers['Authorization'] = `Bearer ${token}`;
		}

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

		localStorage.removeItem('session_token');
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
}

export const api = new ApiClient(API_URL);
