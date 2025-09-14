import { writable } from 'svelte/store';
import type { Collection } from '$types';

interface CollectionsState {
    collections: Collection[];
    loading: boolean;
    error: string | null;
}

function createCollectionsStore() {
    const { subscribe, set, update } = writable<CollectionsState>({
        collections: [],
        loading: false,
        error: null,
    });

    return {
        subscribe,
        setCollections: (collections: Collection[]) =>
            update(state => ({ ...state, collections, error: null })),

        addCollection: (collection: Collection) =>
            update(state => ({
                ...state,
                collections: [...state.collections, collection],
                error: null
            })),

        updateCollection: (id: string, updates: Partial<Collection>) =>
            update(state => ({
                ...state,
                collections: state.collections.map(c =>
                    c.id === id ? { ...c, ...updates } : c
                ),
                error: null
            })),

        removeCollection: (id: string) =>
            update(state => ({
                ...state,
                collections: state.collections.filter(c => c.id !== id),
                error: null
            })),

        setLoading: (loading: boolean) =>
            update(state => ({ ...state, loading })),

        setError: (error: string | null) =>
            update(state => ({ ...state, error })),

        reset: () => set({
            collections: [],
            loading: false,
            error: null,
        })
    };
}

export const collections = createCollectionsStore();
