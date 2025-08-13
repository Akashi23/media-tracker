import { writable } from 'svelte/store';
import type { User } from '$types';

interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
	isGuest: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		token: null,
		isAuthenticated: false,
		isGuest: true
	});

	return {
		subscribe,
		login: (user: User, token: string) => {
			set({
				user,
				token,
				isAuthenticated: true,
				isGuest: false
			});
			localStorage.setItem('auth_token', token);
			localStorage.setItem('user', JSON.stringify(user));
		},
		logout: () => {
			set({
				user: null,
				token: null,
				isAuthenticated: false,
				isGuest: true
			});
			localStorage.removeItem('auth_token');
			localStorage.removeItem('user');
		},
		init: () => {
			const token = localStorage.getItem('auth_token');
			const userStr = localStorage.getItem('user');
			
			if (token && userStr) {
				try {
					const user = JSON.parse(userStr);
					set({
						user,
						token,
						isAuthenticated: true,
						isGuest: false
					});
				} catch (error) {
					console.error('Failed to parse user data:', error);
					localStorage.removeItem('auth_token');
					localStorage.removeItem('user');
				}
			}
		},
		setGuest: (isGuest: boolean) => {
			update(state => ({ ...state, isGuest }));
		}
	};
}

export const auth = createAuthStore();
