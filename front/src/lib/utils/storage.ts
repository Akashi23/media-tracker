import { v4 as uuidv4 } from 'uuid';
import type { Entry, MediaItem, GuestData } from '$types';
import { browser } from '$app/environment';

const GUEST_ID_KEY = 'guestId';
const MEDIA_STORE_KEY = 'media-store';

export const storage = {
	// Guest ID management
	getGuestId: (): string => {
		let guestId = localStorage.getItem(GUEST_ID_KEY);
		if (!guestId) {
			guestId = uuidv4();
			localStorage.setItem(GUEST_ID_KEY, guestId);
		}
		return guestId;
	},

	setGuestId: (guestId: string): void => {
		localStorage.setItem(GUEST_ID_KEY, guestId);
	},

	// Media store management
	getMediaStore: (): { entries: Entry[] } => {
		if (!browser) {
			return { entries: [] }; // SSR
		}

		const data = localStorage.getItem(MEDIA_STORE_KEY);
		if (data) {
			try {
				console.log("data", JSON.parse(data));
				return JSON.parse(data);
			} catch (error) {
				console.error('Failed to parse media store:', error);
			}
		}
		return { entries: [] };
	},

	setMediaStore: (data: { entries: Entry[] }): void => {
		console.log("data", data);
		localStorage.setItem(MEDIA_STORE_KEY, JSON.stringify(data));
	},

	// Entry management
	addEntry: (entry: Entry, media: MediaItem): void => {
		const store = storage.getMediaStore();
		entry.media = media;
		store.entries.push(entry);
		storage.setMediaStore(store);
	},

	updateEntry: (id: string, updates: Partial<Entry>): void => {
		const store = storage.getMediaStore();
		const index = store.entries.findIndex(entry => entry.id === id);
		if (index !== -1) {
			store.entries[index] = { ...store.entries[index], ...updates };
			storage.setMediaStore(store);
		}
	},

	removeEntry: (id: string): void => {
		const store = storage.getMediaStore();
		store.entries = store.entries.filter(entry => entry.id !== id);
		storage.setMediaStore(store);
	},

	getMediaById: (id: string): MediaItem | undefined => {
		const store = storage.getMediaStore();
		return store.entries.find(entry => entry.media?.id === id)?.media;
	},

	// Guest data management
	getGuestData: (): GuestData => {
		const guestId = storage.getGuestId();
		const store = storage.getMediaStore();
		return {
			guestId,
			entries: store.entries,
		};
	},

	setGuestData: (data: GuestData): void => {
		storage.setGuestId(data.guestId);
		storage.setMediaStore({
			entries: data.entries,
		});
	},

	// Clear all guest data
	clearGuestData: (): void => {
		localStorage.removeItem(GUEST_ID_KEY);
		localStorage.removeItem(MEDIA_STORE_KEY);
	},

	// Export/Import
	exportData: (): string => {
		const data = storage.getGuestData();
		return JSON.stringify(data, null, 2);
	},

	importData: (jsonData: string): boolean => {
		try {
			const data: GuestData = JSON.parse(jsonData);
			storage.setGuestData(data);
			return true;
		} catch (error) {
			console.error('Failed to import data:', error);
			return false;
		}
	}
};
