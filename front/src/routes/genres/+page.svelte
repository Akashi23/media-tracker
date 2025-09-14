<script lang="ts">
    import { onMount } from "svelte";
    import { auth } from "$stores/auth";
    import { entries } from "$stores/entries";
    import { entriesApi } from "$utils/api";
    import { storage } from "$utils/storage";
    import type { Entry } from "$types";

    let allGenres: string[] = [];
    let genreStats: Record<string, { count: number; entries: Entry[] }> = {};
    let selectedGenre: string | null = null;
    let loading = true;
    let searchTerm = "";

    onMount(async () => {
        // Initialize auth state from localStorage
        auth.init();

        await loadGenres();
        loading = false;
    });

    async function loadGenres() {
        try {
            let entriesList: Entry[] = [];

            if ($auth.isAuthenticated && $auth.token) {
                // Load entries from backend
                entriesList = await entriesApi.list($auth.token);
            } else {
                // Load entries from local storage for guest mode
                const store = storage.getMediaStore();
                entriesList = store.entries;
            }

            // Extract all unique genres
            const genreSet = new Set<string>();
            entriesList.forEach((entry) => {
                if (entry.media?.genres) {
                    entry.media.genres.forEach((genre) => {
                        if (genre.trim()) {
                            genreSet.add(genre.trim());
                        }
                    });
                }
            });

            allGenres = Array.from(genreSet).sort();

            // Calculate genre statistics
            genreStats = {};
            allGenres.forEach((genre) => {
                const genreEntries = entriesList.filter((entry) =>
                    entry.media?.genres?.includes(genre)
                );
                genreStats[genre] = {
                    count: genreEntries.length,
                    entries: genreEntries,
                };
            });
        } catch (error) {
            console.error("Failed to load genres:", error);
        }
    }

    function getFilteredGenres(): string[] {
        if (!searchTerm.trim()) {
            return allGenres;
        }

        return allGenres.filter((genre) =>
            genre.toLowerCase().includes(searchTerm.toLowerCase())
        );
    }

    function selectGenre(genre: string) {
        selectedGenre = selectedGenre === genre ? null : genre;
    }

    function getGenreColor(genre: string): string {
        // Generate a consistent color for each genre
        const colors = [
            "bg-blue-100 text-blue-800",
            "bg-green-100 text-green-800",
            "bg-purple-100 text-purple-800",
            "bg-pink-100 text-pink-800",
            "bg-yellow-100 text-yellow-800",
            "bg-indigo-100 text-indigo-800",
            "bg-red-100 text-red-800",
            "bg-teal-100 text-teal-800",
            "bg-orange-100 text-orange-800",
            "bg-cyan-100 text-cyan-800",
        ];

        const index = genre.charCodeAt(0) % colors.length;
        return colors[index];
    }

    function getStatusColor(status: string): string {
        switch (status) {
            case "completed":
                return "bg-green-100 text-green-800";
            case "in_progress":
                return "bg-blue-100 text-blue-800";
            case "planned":
                return "bg-yellow-100 text-yellow-800";
            case "on_hold":
                return "bg-orange-100 text-orange-800";
            case "dropped":
                return "bg-red-100 text-red-800";
            default:
                return "bg-gray-100 text-gray-800";
        }
    }
</script>

<svelte:head>
    <title>Genres - Media Tracker</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
    <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Genres</h1>
        <p class="text-gray-600">Explore your media by genre</p>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <div
                class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"
            />
        </div>
    {:else}
        <!-- Search and Stats -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <div
                class="flex flex-col md:flex-row md:items-center md:justify-between gap-4"
            >
                <div class="flex-1">
                    <input
                        type="text"
                        placeholder="Search genres..."
                        bind:value={searchTerm}
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                </div>
                <div class="text-sm text-gray-600">
                    {allGenres.length} unique genres found
                </div>
            </div>
        </div>

        <!-- Genres Grid -->
        <div
            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 mb-6"
        >
            {#each getFilteredGenres() as genre}
                <button
                    on:click={() => selectGenre(genre)}
                    class="p-4 rounded-lg border-2 transition-all duration-200 hover:shadow-md {selectedGenre ===
                    genre
                        ? 'border-blue-500 bg-blue-50'
                        : 'border-gray-200 bg-white hover:border-gray-300'}"
                >
                    <div class="flex items-center justify-between mb-2">
                        <span class="font-medium text-gray-900 truncate"
                            >{genre}</span
                        >
                        <span class="text-sm text-gray-500"
                            >{genreStats[genre]?.count || 0}</span
                        >
                    </div>
                    <div class="flex items-center justify-between">
                        <span class="text-xs text-gray-500">entries</span>
                        <div
                            class="w-2 h-2 rounded-full {selectedGenre === genre
                                ? 'bg-blue-500'
                                : 'bg-gray-300'}"
                        />
                    </div>
                </button>
            {/each}
        </div>

        <!-- Selected Genre Details -->
        {#if selectedGenre && genreStats[selectedGenre]}
            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="flex items-center justify-between mb-4">
                    <h2 class="text-xl font-semibold text-gray-900">
                        {selectedGenre}
                    </h2>
                    <span
                        class="px-3 py-1 {getGenreColor(
                            selectedGenre
                        )} text-sm font-medium rounded-full"
                    >
                        {genreStats[selectedGenre].count} entries
                    </span>
                </div>

                <!-- Genre Entries -->
                <div class="space-y-3">
                    {#each genreStats[selectedGenre].entries as entry}
                        <div
                            class="flex items-center space-x-4 p-3 bg-gray-50 rounded-lg"
                        >
                            <div
                                class="w-12 h-12 bg-gray-200 rounded flex items-center justify-center flex-shrink-0"
                            >
                                <span class="text-sm font-medium text-gray-600">
                                    {entry.media?.type
                                        ?.charAt(0)
                                        .toUpperCase() || "M"}
                                </span>
                            </div>

                            <div class="flex-1 min-w-0">
                                <h3
                                    class="text-sm font-medium text-gray-900 truncate"
                                >
                                    {entry.media?.title || "Unknown Title"}
                                </h3>
                                <p class="text-xs text-gray-500">
                                    {entry.media?.year
                                        ? entry.media.year
                                        : "Unknown Year"} • {entry.media
                                        ?.type || "Unknown Type"}
                                </p>
                            </div>

                            <div class="flex items-center space-x-2">
                                <span
                                    class="px-2 py-1 text-xs font-medium rounded-full {getStatusColor(
                                        entry.status
                                    )}"
                                >
                                    {entry.status.replace("_", " ")}
                                </span>
                                {#if entry.rating}
                                    <span class="text-sm text-gray-600">
                                        {entry.rating}/10
                                    </span>
                                {/if}
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}

        <!-- Genre Statistics -->
        <div class="mt-6 bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">
                Genre Statistics
            </h3>

            {#if allGenres.length > 0}
                <div class="space-y-3">
                    {#each allGenres.slice(0, 10) as genre}
                        <div class="flex items-center justify-between">
                            <div class="flex items-center space-x-3">
                                <span
                                    class="px-2 py-1 text-xs font-medium rounded-full {getGenreColor(
                                        genre
                                    )}"
                                >
                                    {genre}
                                </span>
                            </div>
                            <div class="flex items-center space-x-2">
                                <div class="w-24 bg-gray-200 rounded-full h-2">
                                    <div
                                        class="bg-blue-500 h-2 rounded-full"
                                        style="width: {((genreStats[genre]
                                            ?.count || 0) /
                                            Math.max(
                                                ...Object.values(
                                                    genreStats
                                                ).map((s) => s.count)
                                            )) *
                                            100}%"
                                    />
                                </div>
                                <span
                                    class="text-sm font-medium text-gray-900 w-8 text-right"
                                >
                                    {genreStats[genre]?.count || 0}
                                </span>
                            </div>
                        </div>
                    {/each}
                </div>
            {:else}
                <p class="text-gray-500 text-center py-8">
                    No genres found. Add some media with genres to see them
                    here.
                </p>
            {/if}
        </div>
    {/if}
</div>
