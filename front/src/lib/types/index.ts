export type MediaType = 'video' | 'book' | 'anime' | 'game' | 'tv' | 'movie';

export type Status = 'planned' | 'in_progress' | 'completed' | 'on_hold' | 'dropped';

export interface User {
	id: string;
	email: string;
	name: string;
	created_at: string;
}

export interface MediaItem {
	id: string;
	type: MediaType;
	title: string;
	original_title?: string;
	year?: number;
	cover_url?: string;
	creators?: Record<string, any>;
	genres?: string[];
	duration?: number;
	metadata?: Record<string, any>;
	created_at: string;
}

export interface Entry {
	id: string;
	user_id: string;
	media_id: string;
	status: Status;
	rating?: number;
	review_md?: string;
	progress?: Record<string, any>;
	started_at?: string;
	finished_at?: string;
	updated_at: string;
	media?: MediaItem;
}

export interface Collection {
	id: string;
	user_id: string;
	title: string;
	is_public: boolean;
	created_at: string;
	entries?: Entry[];
}

export interface ShareToken {
	token: string;
	kind: string;
	target_id: string;
	created_at: string;
	expires_at?: string;
}

// Request/Response DTOs
export interface LoginRequest {
	email: string;
}

export interface CreateEntryRequest {
	media_id: string;
	status: Status;
	rating?: number;
	review_md?: string;
	progress?: Record<string, any>;
	started_at?: string;
	finished_at?: string;
}

export interface CreateMediaRequest {
	type: MediaType;
	title: string;
	original_title?: string;
	year?: number;
	cover_url?: string;
	creators?: Record<string, any>;
	genres?: string[];
	duration?: number;
	metadata?: Record<string, any>;
}

export interface CreateCollectionRequest {
	title: string;
	is_public: boolean;
	entry_ids?: string[];
}

export interface GuestSnapshotRequest {
	entries: Entry[];
	media: MediaItem[];
}

export interface MergeRequest {
	guest_entries: Entry[];
}

// Local storage types
export interface GuestData {
	guestId: string;
	entries: Entry[];
}
