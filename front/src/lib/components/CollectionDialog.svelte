<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { auth } from "$stores/auth";
    import { collectionsApi } from "$utils/api";
    import { storage } from "$utils/storage";
    import type { Collection, CreateCollectionRequest, Entry } from "$types";

    export let open = false;
    export let collection: Collection | null = null;
    export let availableEntries: Entry[] = [];

    const dispatch = createEventDispatcher();

    // Form fields
    let title = "";
    let isPublic = false;
    let selectedEntryIds: string[] = [];
    let loading = false;

    // Track if form has been initialized
    let formInitialized = false;

    // Initialize form when collection changes
    $: if (collection && !formInitialized) {
        title = collection.title;
        isPublic = collection.is_public;
        selectedEntryIds = collection.entries?.map((e) => e.id) || [];
        formInitialized = true;
    }

    // Reset form initialization when dialog closes
    $: if (!open) {
        formInitialized = false;
    }

    async function handleSubmit() {
        if (!title.trim()) return;

        loading = true;
        try {
            const collectionData: CreateCollectionRequest = {
                title: title.trim(),
                is_public: isPublic,
                entry_ids:
                    selectedEntryIds.length > 0 ? selectedEntryIds : undefined,
            };

            if (collection) {
                // Update existing collection
                if ($auth.isAuthenticated && $auth.token) {
                    const updatedCollection = await collectionsApi.update(
                        collection.id,
                        collectionData,
                        $auth.token
                    );
                    dispatch("collection-updated", updatedCollection);
                } else {
                    // Guest mode - update locally (simplified)
                    const updatedCollection = {
                        ...collection,
                        ...collectionData,
                    };
                    dispatch("collection-updated", updatedCollection);
                }
            } else {
                // Create new collection
                if ($auth.isAuthenticated && $auth.token) {
                    const newCollection = await collectionsApi.create(
                        collectionData,
                        $auth.token
                    );
                    dispatch("collection-created", newCollection);
                } else {
                    // Guest mode - create locally (simplified)
                    const newCollection = {
                        id: crypto.randomUUID(),
                        user_id: storage.getGuestId(),
                        ...collectionData,
                        created_at: new Date().toISOString(),
                        entries: availableEntries.filter((e) =>
                            selectedEntryIds.includes(e.id)
                        ),
                    };
                    dispatch("collection-created", newCollection);
                }
            }

            // Reset form
            title = "";
            isPublic = false;
            selectedEntryIds = [];
            open = false;
        } catch (error) {
            console.error("Failed to save collection:", error);
            alert("Failed to save collection. Please try again.");
        } finally {
            loading = false;
        }
    }

    function handleClose() {
        open = false;
        formInitialized = false;
        title = "";
        isPublic = false;
        selectedEntryIds = [];
    }

    function toggleEntry(entryId: string) {
        if (selectedEntryIds.includes(entryId)) {
            selectedEntryIds = selectedEntryIds.filter((id) => id !== entryId);
        } else {
            selectedEntryIds = [...selectedEntryIds, entryId];
        }
    }

    function isEntrySelected(entryId: string): boolean {
        return selectedEntryIds.includes(entryId);
    }
</script>

{#if open}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div
            class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto"
        >
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-xl font-bold">
                    {collection ? "Edit Collection" : "Create Collection"}
                </h2>
                <button
                    class="text-gray-400 hover:text-gray-600"
                    on:click={handleClose}
                >
                    ✕
                </button>
            </div>

            <form on:submit|preventDefault={handleSubmit} class="space-y-4">
                <!-- Title -->
                <div>
                    <label
                        for="title"
                        class="block text-sm font-medium text-gray-700 mb-1"
                    >
                        Collection Title
                    </label>
                    <input
                        id="title"
                        type="text"
                        bind:value={title}
                        class="input"
                        placeholder="Enter collection title..."
                        required
                    />
                </div>

                <!-- Public/Private -->
                <div>
                    <label class="flex items-center space-x-2">
                        <input
                            type="checkbox"
                            bind:checked={isPublic}
                            class="rounded border-gray-300"
                        />
                        <span class="text-sm font-medium text-gray-700">
                            Make this collection public
                        </span>
                    </label>
                    <p class="text-xs text-gray-500 mt-1">
                        Public collections can be shared with others via a link
                    </p>
                </div>

                <!-- Entry Selection -->
                {#if availableEntries.length > 0}
                    <div>
                        <div
                            class="block text-sm font-medium text-gray-700 mb-2"
                        >
                            Select Entries (optional)
                        </div>
                        <div
                            class="border rounded-lg p-3 max-h-48 overflow-y-auto"
                        >
                            {#each availableEntries as entry}
                                <label
                                    class="flex items-center space-x-3 py-2 hover:bg-gray-50 rounded px-2"
                                >
                                    <input
                                        id="entry-{entry.id}"
                                        type="checkbox"
                                        checked={isEntrySelected(entry.id)}
                                        on:change={() => toggleEntry(entry.id)}
                                        class="rounded border-gray-300"
                                    />
                                    <span class="text-lg">
                                        {entry.media?.type === "movie"
                                            ? "🎬"
                                            : entry.media?.type === "tv"
                                            ? "📺"
                                            : entry.media?.type === "anime"
                                            ? "🌸"
                                            : entry.media?.type === "book"
                                            ? "📚"
                                            : entry.media?.type === "game"
                                            ? "🎮"
                                            : "📹"}
                                    </span>
                                    <div class="flex-1">
                                        <div class="font-medium text-sm">
                                            {entry.media?.title ||
                                                "Unknown Title"}
                                        </div>
                                        {#if entry.media?.year}
                                            <div class="text-xs text-gray-500">
                                                ({entry.media.year})
                                            </div>
                                        {/if}
                                    </div>
                                    {#if entry.rating}
                                        <span class="text-yellow-500 text-xs">
                                            ⭐ {entry.rating}/10
                                        </span>
                                    {/if}
                                </label>
                            {/each}
                        </div>
                        <p class="text-xs text-gray-500 mt-1">
                            {selectedEntryIds.length} entries selected
                        </p>
                    </div>
                {/if}

                <!-- Actions -->
                <div class="flex space-x-3 pt-4">
                    <button
                        type="submit"
                        class="btn btn-primary flex-1"
                        disabled={loading}
                    >
                        {loading
                            ? collection
                                ? "Updating..."
                                : "Creating..."
                            : collection
                            ? "Update Collection"
                            : "Create Collection"}
                    </button>
                    <button
                        type="button"
                        class="btn btn-secondary"
                        on:click={handleClose}
                    >
                        Cancel
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}
