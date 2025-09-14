<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Collection } from "$types";

    export let collection: Collection;

    const dispatch = createEventDispatcher();

    function handleEdit() {
        dispatch("edit", collection);
    }

    function handleDelete() {
        if (confirm("Are you sure you want to delete this collection?")) {
            dispatch("delete", collection.id);
        }
    }

    function handleShare() {
        dispatch("share", collection);
    }

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString();
    }

    function getEntryCount(): number {
        return collection.entries?.length || 0;
    }

    function getPublicStatus(): string {
        return collection.is_public ? "Public" : "Private";
    }

    function getPublicStatusColor(): string {
        return collection.is_public
            ? "bg-green-100 text-green-800"
            : "bg-gray-100 text-gray-800";
    }
</script>

<div class="card hover:shadow-lg transition-shadow">
    <!-- Header -->
    <div class="flex justify-between items-start mb-4">
        <div class="flex-1">
            <h3 class="font-semibold text-lg text-gray-900 mb-1">
                {collection.title}
            </h3>
            <div class="flex items-center space-x-2 text-sm text-gray-600">
                <span>Created: {formatDate(collection.created_at)}</span>
                <span>•</span>
                <span>{getEntryCount()} entries</span>
            </div>
        </div>

        <div class="flex space-x-2">
            <button
                class="text-gray-400 hover:text-blue-600"
                on:click={handleShare}
                title="Share collection"
            >
                🔗
            </button>
            <button
                class="text-gray-400 hover:text-gray-600"
                on:click={handleEdit}
                title="Edit collection"
            >
                ✏️
            </button>
            <button
                class="text-gray-400 hover:text-red-600"
                on:click={handleDelete}
                title="Delete collection"
            >
                🗑️
            </button>
        </div>
    </div>

    <!-- Status -->
    <div class="flex justify-between items-center mb-3">
        <span
            class="px-2 py-1 rounded-full text-xs font-medium {getPublicStatusColor()}"
        >
            {getPublicStatus()}
        </span>
    </div>

    <!-- Entries Preview -->
    {#if collection.entries && collection.entries.length > 0}
        <div class="border-t pt-3">
            <div class="text-sm text-gray-600 mb-2">Recent entries:</div>
            <div class="space-y-2">
                {#each collection.entries.slice(0, 3) as entry}
                    <div class="flex items-center space-x-2 text-sm">
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
                        <span class="flex-1 truncate">
                            {entry.media?.title || "Unknown Title"}
                        </span>
                        {#if entry.rating}
                            <span class="text-yellow-500 text-xs">
                                ⭐ {entry.rating}/10
                            </span>
                        {/if}
                    </div>
                {/each}
                {#if collection.entries.length > 3}
                    <div class="text-xs text-gray-500">
                        +{collection.entries.length - 3} more entries
                    </div>
                {/if}
            </div>
        </div>
    {:else}
        <div class="border-t pt-3">
            <div class="text-sm text-gray-500 text-center py-2">
                No entries in this collection yet
            </div>
        </div>
    {/if}
</div>
