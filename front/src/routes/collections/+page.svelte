<script lang="ts">
    import { onMount } from "svelte";
    import { auth } from "$stores/auth";
    import { collections } from "$stores/collections";
    import { entries } from "$stores/entries";
    import { collectionsApi, entriesApi } from "$utils/api";
    import { storage } from "$utils/storage";
    import CollectionCard from "$components/CollectionCard.svelte";
    import CollectionDialog from "$components/CollectionDialog.svelte";
    import type { Collection, Entry } from "$types";

    let showCreateDialog = false;
    let showEditDialog = false;
    let showShareModal = false;
    let shareUrl = "";
    let selectedCollection: Collection | null = null;
    let loading = false;

    onMount(async () => {
        // Initialize auth state from localStorage
        auth.init();

        // If user is authenticated, sync local entries with backend
        if ($auth.isAuthenticated && $auth.token) {
            await syncLocalEntries();
        }

        // Load entries first, then collections
        await loadEntries();
        await loadCollections();
    });

    async function loadEntries() {
        if ($auth.isAuthenticated && $auth.token) {
            try {
                // Load all entries from backend
                const entriesList = await entriesApi.list($auth.token);
                entries.setEntries(entriesList);
            } catch (error) {
                console.error("Failed to load entries:", error);
            }
        } else {
            // Load from local storage for guest mode
            const store = storage.getMediaStore();
            entries.setEntries(store.entries);
        }
    }

    async function syncLocalEntries() {
        try {
            console.log("Syncing local entries with backend...");

            // Get local entries from storage
            const store = storage.getMediaStore();
            const localEntries = store.entries;

            if (localEntries.length === 0) {
                console.log("No local entries to sync");
                return;
            }

            console.log("Found", localEntries.length, "local entries to sync");

            // Convert local entries to SyncEntryRequest format
            const entriesToSync = localEntries.map((entry) => ({
                media: {
                    type: entry.media?.type || "video",
                    title: entry.media?.title || "Unknown Title",
                    original_title: entry.media?.original_title,
                    year: entry.media?.year,
                    cover_url: entry.media?.cover_url,
                    creators: entry.media?.creators,
                    genres: entry.media?.genres || [],
                    duration: entry.media?.duration,
                    metadata: entry.media?.metadata,
                },
                status: entry.status,
                rating: entry.rating,
                review_md: entry.review_md,
                progress: entry.progress,
                started_at: entry.started_at,
                finished_at: entry.finished_at,
            }));

            // Sync with server
            const response = await entriesApi.sync(entriesToSync, $auth.token!);
            console.log("Sync completed:", response.count, "entries synced");

            if (response.errors && response.errors.length > 0) {
                console.warn("Sync errors:", response.errors);
            }

            // Clear local storage after successful sync
            storage.clearEntries();
            console.log("Local entries cleared after sync");
        } catch (error) {
            console.error("Failed to sync local entries:", error);
        }
    }

    async function loadCollections() {
        if ($auth.isAuthenticated && $auth.token) {
            loading = true;
            try {
                const collectionsList = await collectionsApi.list($auth.token);
                collections.setCollections(collectionsList);
            } catch (error) {
                console.error("Failed to load collections:", error);
                collections.setError("Failed to load collections");
            } finally {
                loading = false;
            }
        } else {
            // Guest mode - load from local storage
            const store = storage.getMediaStore();
            collections.setCollections(store.collections);
        }
    }

    function handleCreateCollection() {
        selectedCollection = null;
        showCreateDialog = true;
    }

    function handleEditCollection(collection: Collection) {
        selectedCollection = collection;
        showEditDialog = true;
    }

    async function handleDeleteCollection(id: string) {
        if ($auth.isAuthenticated && $auth.token) {
            try {
                await collectionsApi.delete(id, $auth.token);
                // Reload collections from server to get updated data
                await loadCollections();
            } catch (error) {
                console.error("Failed to delete collection:", error);
                alert("Failed to delete collection. Please try again.");
            }
        } else {
            // Guest mode - remove locally
            storage.removeCollection(id);
            collections.removeCollection(id);
        }
    }

    async function handleShareCollection(collection: Collection) {
        if ($auth.isAuthenticated && $auth.token) {
            try {
                const response = await collectionsApi.createShare(
                    collection.id,
                    $auth.token
                );
                shareUrl = `${window.location.origin}${response.share_url}`;
                showShareModal = true;
            } catch (error) {
                console.error("Failed to create share link:", error);
                alert("Failed to create share link. Please try again.");
            }
        } else {
            alert("Sharing is only available for registered users");
        }
    }

    async function copyShareUrl() {
        try {
            await navigator.clipboard.writeText(shareUrl);
            alert("Share link copied to clipboard!");
        } catch (error) {
            console.error("Failed to copy to clipboard:", error);
        }
    }

    function closeShareModal() {
        showShareModal = false;
        shareUrl = "";
    }

    async function handleCollectionCreated(event: CustomEvent<Collection>) {
        const collection = event.detail;

        // For authenticated users, reload collections from server to get updated data
        if ($auth.isAuthenticated && $auth.token) {
            await loadCollections();
        } else {
            // For guest mode, add to store and local storage
            collections.addCollection(collection);
            storage.addCollection(collection);
        }

        showCreateDialog = false;
    }

    async function handleCollectionUpdated(event: CustomEvent<Collection>) {
        const collection = event.detail;

        // For authenticated users, reload collections from server to get updated data
        if ($auth.isAuthenticated && $auth.token) {
            await loadCollections();
        } else {
            // For guest mode, update store and local storage
            collections.updateCollection(collection.id, collection);
            storage.updateCollection(collection.id, collection);
        }

        showEditDialog = false;
        selectedCollection = null;
    }
</script>

<svelte:head>
    <title>Collections - Media Tracker</title>
</svelte:head>

<div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
        <div>
            <h1 class="text-3xl font-bold text-gray-900">Collections</h1>
            <p class="text-gray-600 mt-1">
                {#if $auth.isAuthenticated}
                    Organize your media entries into collections
                {:else}
                    Collections are available for registered users
                {/if}
            </p>
        </div>

        {#if $auth.isAuthenticated}
            <button class="btn btn-primary" on:click={handleCreateCollection}>
                + Create Collection
            </button>
        {/if}
    </div>

    <!-- Collections List -->
    {#if loading}
        <div class="flex justify-center py-12">
            <div
                class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
            />
        </div>
    {:else if $collections.collections.length === 0}
        <div class="text-center py-12">
            <div class="text-gray-400 text-6xl mb-4">📚</div>
            <h3 class="text-lg font-medium text-gray-900 mb-2">
                {#if $auth.isAuthenticated}
                    No collections yet
                {:else}
                    Collections require authentication
                {/if}
            </h3>
            <p class="text-gray-600 mb-4">
                {#if $auth.isAuthenticated}
                    Create your first collection to organize your media entries.
                {:else}
                    Sign in to create and manage collections.
                {/if}
            </p>
            {#if $auth.isAuthenticated}
                <button
                    class="btn btn-primary"
                    on:click={handleCreateCollection}
                >
                    Create Your First Collection
                </button>
            {/if}
        </div>
    {:else}
        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {#each $collections.collections as collection (collection.id)}
                <CollectionCard
                    {collection}
                    on:edit={(e) => handleEditCollection(e.detail)}
                    on:delete={(e) => handleDeleteCollection(e.detail)}
                    on:share={(e) => handleShareCollection(e.detail)}
                />
            {/each}
        </div>
    {/if}
</div>

<!-- Create Collection Dialog -->
{#if showCreateDialog}
    <CollectionDialog
        bind:open={showCreateDialog}
        availableEntries={$entries.entries}
        on:collection-created={handleCollectionCreated}
    />
{/if}

<!-- Edit Collection Dialog -->
{#if showEditDialog}
    <CollectionDialog
        bind:open={showEditDialog}
        collection={selectedCollection}
        availableEntries={$entries.entries}
        on:collection-updated={handleCollectionUpdated}
    />
{/if}

<!-- Share Modal -->
{#if showShareModal}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-bold text-gray-900">
                    Share Collection
                </h3>
                <button
                    class="text-gray-400 hover:text-gray-600"
                    on:click={closeShareModal}
                >
                    ✕
                </button>
            </div>

            <div class="mb-4">
                <p class="text-gray-700 mb-3">
                    Your collection has been shared! Copy the link below:
                </p>
                <div class="flex items-center space-x-2">
                    <input
                        type="text"
                        value={shareUrl}
                        readonly
                        class="input flex-1 text-sm"
                    />
                    <button
                        class="btn btn-secondary px-3"
                        on:click={copyShareUrl}
                    >
                        Copy
                    </button>
                </div>
            </div>

            <div class="flex justify-end space-x-3">
                <button class="btn btn-secondary" on:click={closeShareModal}>
                    Close
                </button>
            </div>
        </div>
    </div>
{/if}
