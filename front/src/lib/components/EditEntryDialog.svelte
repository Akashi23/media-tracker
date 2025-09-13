<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { auth } from "$stores/auth";
    import { entriesApi } from "$utils/api";
    import { storage } from "$utils/storage";
    import type { Entry, Status, CreateEntryRequest } from "$types";

    export let open = false;
    export let entry: Entry | null = null;

    const dispatch = createEventDispatcher();

    // Form fields
    let status: Status = "planned";
    let rating: number | undefined;
    let reviewMd = "";
    let startedAt = "";
    let finishedAt = "";
    let progress: Record<string, any> = {};
    let loading = false;

    // Progress fields based on media type
    let pagesRead: number | undefined;
    let pagesTotal: number | undefined;
    let episodesSeen: number | undefined;
    let episodesTotal: number | undefined;
    let gamePercent: number | undefined;

    const statuses: Array<{ value: Status; label: string }> = [
        { value: "planned", label: "Planned" },
        { value: "in_progress", label: "In Progress" },
        { value: "completed", label: "Completed" },
        { value: "on_hold", label: "On Hold" },
        { value: "dropped", label: "Dropped" },
    ];

    // Track if form has been initialized to prevent overwriting user input
    let formInitialized = false;

    // Initialize form when entry changes
    $: if (entry && !formInitialized) {
        status = entry.status;
        rating = entry.rating;
        reviewMd = entry.review_md || "";
        startedAt = entry.started_at ? entry.started_at.split("T")[0] : "";
        finishedAt = entry.finished_at ? entry.finished_at.split("T")[0] : "";
        progress = entry.progress || {};

        // Initialize progress fields based on media type
        if (entry.media?.type === "book") {
            pagesRead = progress.pagesRead;
            pagesTotal = progress.pagesTotal;
        } else if (
            entry.media?.type === "anime" ||
            entry.media?.type === "tv"
        ) {
            episodesSeen = progress.episodesSeen;
            episodesTotal = progress.episodesTotal;
        } else if (entry.media?.type === "game") {
            gamePercent = progress.gamePercent;
        }

        formInitialized = true;
    }

    // Reset form initialization when dialog closes
    $: if (!open) {
        formInitialized = false;
    }

    function updateProgress() {
        if (entry?.media?.type === "book") {
            progress = {
                pagesRead,
                pagesTotal,
            };
        } else if (
            entry?.media?.type === "anime" ||
            entry?.media?.type === "tv"
        ) {
            progress = {
                episodesSeen,
                episodesTotal,
            };
        } else if (entry?.media?.type === "game") {
            progress = {
                gamePercent,
            };
        }
    }

    async function handleSubmit() {
        if (!entry) return;

        loading = true;
        try {
            updateProgress();

            const entryData: CreateEntryRequest = {
                media_id: entry.media_id,
                status,
                rating,
                review_md: reviewMd || undefined,
                progress:
                    Object.keys(progress).length > 0 ? progress : undefined,
                started_at: startedAt
                    ? new Date(startedAt).toISOString()
                    : undefined,
                finished_at: finishedAt
                    ? new Date(finishedAt).toISOString()
                    : undefined,
            };

            if ($auth.isAuthenticated && $auth.token) {
                const updatedEntry = await entriesApi.update(
                    entry.id,
                    entryData,
                    $auth.token
                );
                dispatch("entry-updated", updatedEntry);
            } else {
                // Guest mode - update locally
                const updatedEntry = {
                    ...entry,
                    ...entryData,
                    updated_at: new Date().toISOString(),
                };
                storage.updateEntry(entry.id, entryData);
                dispatch("entry-updated", updatedEntry);
            }

            open = false;
        } catch (error) {
            console.error("Failed to update entry:", error);
            alert("Failed to update entry. Please try again.");
        } finally {
            loading = false;
        }
    }

    function handleClose() {
        open = false;
        formInitialized = false;
    }

    function formatDate(dateStr?: string): string {
        if (!dateStr) return "";
        return new Date(dateStr).toISOString().split("T")[0];
    }
</script>

{#if open && entry}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div
            class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto"
        >
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-xl font-bold">Edit Entry</h2>
                <button
                    class="text-gray-400 hover:text-gray-600"
                    on:click={handleClose}
                >
                    ✕
                </button>
            </div>

            <!-- Media Info (Read-only) -->
            <div class="mb-6 p-4 bg-gray-50 rounded-lg">
                <div class="flex items-center space-x-3">
                    <span class="text-2xl">
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
                    <div>
                        <h3 class="font-semibold text-lg">
                            {entry.media?.title || "Unknown Title"}
                        </h3>
                        {#if entry.media?.year}
                            <p class="text-sm text-gray-600">
                                ({entry.media.year})
                            </p>
                        {/if}
                    </div>
                </div>
            </div>

            <form on:submit|preventDefault={handleSubmit} class="space-y-4">
                <!-- Status -->
                <div>
                    <label
                        for="status"
                        class="block text-sm font-medium text-gray-700 mb-1"
                    >
                        Status
                    </label>
                    <select id="status" bind:value={status} class="input">
                        {#each statuses as statusOption}
                            <option value={statusOption.value}
                                >{statusOption.label}</option
                            >
                        {/each}
                    </select>
                </div>

                <!-- Rating -->
                <div>
                    <label
                        for="rating"
                        class="block text-sm font-medium text-gray-700 mb-1"
                    >
                        Rating (optional)
                    </label>
                    <input
                        id="rating"
                        type="number"
                        bind:value={rating}
                        class="input"
                        placeholder="8.5"
                        min="0"
                        max="10"
                        step="0.5"
                    />
                </div>

                <!-- Progress Fields -->
                {#if entry.media?.type === "book"}
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label
                                for="pagesRead"
                                class="block text-sm font-medium text-gray-700 mb-1"
                            >
                                Pages Read
                            </label>
                            <input
                                id="pagesRead"
                                type="number"
                                bind:value={pagesRead}
                                class="input"
                                placeholder="0"
                                min="0"
                            />
                        </div>
                        <div>
                            <label
                                for="pagesTotal"
                                class="block text-sm font-medium text-gray-700 mb-1"
                            >
                                Total Pages
                            </label>
                            <input
                                id="pagesTotal"
                                type="number"
                                bind:value={pagesTotal}
                                class="input"
                                placeholder="0"
                                min="0"
                            />
                        </div>
                    </div>
                {:else if entry.media?.type === "anime" || entry.media?.type === "tv"}
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label
                                for="episodesSeen"
                                class="block text-sm font-medium text-gray-700 mb-1"
                            >
                                Episodes Seen
                            </label>
                            <input
                                id="episodesSeen"
                                type="number"
                                bind:value={episodesSeen}
                                class="input"
                                placeholder="0"
                                min="0"
                            />
                        </div>
                        <div>
                            <label
                                for="episodesTotal"
                                class="block text-sm font-medium text-gray-700 mb-1"
                            >
                                Total Episodes
                            </label>
                            <input
                                id="episodesTotal"
                                type="number"
                                bind:value={episodesTotal}
                                class="input"
                                placeholder="0"
                                min="0"
                            />
                        </div>
                    </div>
                {:else if entry.media?.type === "game"}
                    <div>
                        <label
                            for="gamePercent"
                            class="block text-sm font-medium text-gray-700 mb-1"
                        >
                            Completion Percentage
                        </label>
                        <input
                            id="gamePercent"
                            type="number"
                            bind:value={gamePercent}
                            class="input"
                            placeholder="0"
                            min="0"
                            max="100"
                        />
                    </div>
                {/if}

                <!-- Dates -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label
                            for="startedAt"
                            class="block text-sm font-medium text-gray-700 mb-1"
                        >
                            Started Date (optional)
                        </label>
                        <input
                            id="startedAt"
                            type="date"
                            bind:value={startedAt}
                            class="input"
                        />
                    </div>
                    <div>
                        <label
                            for="finishedAt"
                            class="block text-sm font-medium text-gray-700 mb-1"
                        >
                            Finished Date (optional)
                        </label>
                        <input
                            id="finishedAt"
                            type="date"
                            bind:value={finishedAt}
                            class="input"
                        />
                    </div>
                </div>

                <!-- Review -->
                <div>
                    <label
                        for="review"
                        class="block text-sm font-medium text-gray-700 mb-1"
                    >
                        Review (optional)
                    </label>
                    <textarea
                        id="review"
                        bind:value={reviewMd}
                        class="input"
                        rows="4"
                        placeholder="Write your review in Markdown..."
                    />
                </div>

                <!-- Actions -->
                <div class="flex space-x-3 pt-4">
                    <button
                        type="submit"
                        class="btn btn-primary flex-1"
                        disabled={loading}
                    >
                        {loading ? "Updating..." : "Update Entry"}
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
