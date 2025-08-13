import type { 
	User, 
	Entry, 
	MediaItem, 
	Collection, 
	LoginRequest, 
	CreateEntryRequest, 
	CreateMediaRequest,
	CreateCollectionRequest,
	GuestSnapshotRequest,
	MergeRequest
} from '$types';

const API_BASE = '/api';

class ApiError extends Error {
	constructor(message: string, public status: number) {
		super(message);
		this.name = 'ApiError';
	}
}

async function request<T>(
	endpoint: string, 
	options: RequestInit = {}
): Promise<T> {
	const url = `${API_BASE}${endpoint}`;
	const config: RequestInit = {
		headers: {
			'Content-Type': 'application/json',
			...options.headers
		},
		...options
	};

	const response = await fetch(url, config);
	
	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: 'Unknown error' }));
		throw new ApiError(error.error || 'Request failed', response.status);
	}

	return response.json();
}

// Auth API
export const authApi = {
	login: (data: LoginRequest) => 
		request<{ token: string }>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	
	logout: (token: string) => 
		request('/auth/logout', {
			method: 'POST',
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	getProfile: (token: string) => 
		request<User>('/auth/me', {
			headers: { Authorization: `Bearer ${token}` }
		})
};

// Media API
export const mediaApi = {
	create: (data: CreateMediaRequest, token: string) => 
		request<MediaItem>('/media', {
			method: 'POST',
			body: JSON.stringify(data),
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	search: (query: string, type?: string) => 
		request<MediaItem[]>(`/media/search?q=${encodeURIComponent(query)}${type ? `&type=${type}` : ''}`)
};

// Entries API
export const entriesApi = {
	list: (token: string, params?: { type?: string; status?: string }) => {
		const searchParams = new URLSearchParams();
		if (params?.type) searchParams.append('type', params.type);
		if (params?.status) searchParams.append('status', params.status);
		
		return request<Entry[]>(`/entries?${searchParams.toString()}`, {
			headers: { Authorization: `Bearer ${token}` }
		});
	},
	
	create: (data: CreateEntryRequest, token: string) => 
		request<Entry>('/entries', {
			method: 'POST',
			body: JSON.stringify(data),
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	get: (id: string, token: string) => 
		request<Entry>(`/entries/${id}`, {
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	update: (id: string, data: CreateEntryRequest, token: string) => 
		request<Entry>(`/entries/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data),
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	delete: (id: string, token: string) => 
		request(`/entries/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` }
		})
};

// Collections API
export const collectionsApi = {
	list: (token: string) => 
		request<Collection[]>('/collections', {
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	create: (data: CreateCollectionRequest, token: string) => 
		request<Collection>('/collections', {
			method: 'POST',
			body: JSON.stringify(data),
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	get: (id: string, token: string) => 
		request<Collection>(`/collections/${id}`, {
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	update: (id: string, data: Partial<CreateCollectionRequest>, token: string) => 
		request<Collection>(`/collections/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data),
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	delete: (id: string, token: string) => 
		request(`/collections/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` }
		}),
	
	createShare: (id: string, token: string) => 
		request<{ share_url: string }>(`/collections/${id}/share`, {
			method: 'POST',
			headers: { Authorization: `Bearer ${token}` }
		})
};

// Guest API
export const guestApi = {
	createSnapshot: (data: GuestSnapshotRequest) => 
		request<{ share_url: string }>('/guest/snapshot', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	
	merge: (data: MergeRequest, token: string) => 
		request('/guest/merge', {
			method: 'POST',
			body: JSON.stringify(data),
			headers: { Authorization: `Bearer ${token}` }
		})
};

// Public API
export const publicApi = {
	getShare: (token: string) => 
		request<any>(`/s/${token}`)
};
