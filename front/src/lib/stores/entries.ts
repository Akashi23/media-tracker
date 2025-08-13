import { writable } from 'svelte/store';
import type { Entry, Status, MediaType } from '$types';

interface EntriesState {
	entries: Entry[];
	loading: boolean;
	error: string | null;
}

function createEntriesStore() {
	const { subscribe, set, update } = writable<EntriesState>({
		entries: [],
		loading: false,
		error: null
	});

	return {
		subscribe,
		setEntries: (entries: Entry[]) => {
			set({ entries, loading: false, error: null });
		},
		addEntry: (entry: Entry) => {
			update(state => ({
				...state,
				entries: [entry, ...state.entries]
			}));
		},
		updateEntry: (id: string, updates: Partial<Entry>) => {
			update(state => ({
				...state,
				entries: state.entries.map(entry =>
					entry.id === id ? { ...entry, ...updates } : entry
				)
			}));
		},
		removeEntry: (id: string) => {
			update(state => ({
				...state,
				entries: state.entries.filter(entry => entry.id !== id)
			}));
		},
		setLoading: (loading: boolean) => {
			update(state => ({ ...state, loading }));
		},
		setError: (error: string | null) => {
			update(state => ({ ...state, error }));
		},
		getByStatus: (status: Status) => {
			let entries: Entry[] = [];
			subscribe(state => {
				entries = state.entries.filter(entry => entry.status === status);
			})();
			return entries;
		},
		getByType: (type: MediaType) => {
			let entries: Entry[] = [];
			subscribe(state => {
				entries = state.entries.filter(entry => entry.media?.type === type);
			})();
			return entries;
		},
		clear: () => {
			set({ entries: [], loading: false, error: null });
		}
	};
}

export const entries = createEntriesStore();
